package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
	"github.com/sst/forge/internal/metrics"
)

// BuildHook represents a build hook that can be executed at different stages
type BuildHook struct {
	Name        string
	Description string
	Stage       HookStage
	Command     string
	Timeout     time.Duration
	WorkingDir  string
	Environment map[string]string
}

// HookStage represents the stage at which a hook is executed
type HookStage int

const (
	// HookStagePreBuild executes before the build starts
	HookStagePreBuild HookStage = iota
	// HookStagePostBuild executes after the build completes successfully
	HookStagePostBuild
	// HookStageBuildFailure executes when the build fails
	HookStageBuildFailure
	// HookStagePrePhase executes before each build phase
	HookStagePrePhase
	// HookStagePostPhase executes after each build phase
	HookStagePostPhase
)

// HookManager manages build hooks
type HookManager struct {
	hooks   []BuildHook
	logger  *logger.Logger
	metrics *metrics.MetricsCollector
}

// NewHookManager creates a new hook manager
func NewHookManager(logger *logger.Logger, metrics *metrics.MetricsCollector) *HookManager {
	return &HookManager{
		hooks:   []BuildHook{},
		logger:  logger,
		metrics: metrics,
	}
}

// AddHook adds a build hook
func (hm *HookManager) AddHook(hook BuildHook) {
	hm.hooks = append(hm.hooks, hook)
}

// LoadHooksFromConfig loads hooks from the forge.yml configuration
func (hm *HookManager) LoadHooksFromConfig(config *config.Config, projectDir string) error {
	// In a real implementation, this would parse hooks from config.Overlays or a dedicated hooks section
	// For now, we'll look for hook scripts in the project directory

	hookDirs := []string{
		filepath.Join(projectDir, "hooks"),
		filepath.Join(projectDir, ".forge", "hooks"),
	}

	for _, hookDir := range hookDirs {
		if _, err := os.Stat(hookDir); os.IsNotExist(err) {
			continue
		}

		// Load pre-build hooks
		preBuildDir := filepath.Join(hookDir, "pre-build")
		if err := hm.loadHooksFromDir(preBuildDir, HookStagePreBuild, projectDir); err != nil {
			return fmt.Errorf("failed to load pre-build hooks: %v", err)
		}

		// Load post-build hooks
		postBuildDir := filepath.Join(hookDir, "post-build")
		if err := hm.loadHooksFromDir(postBuildDir, HookStagePostBuild, projectDir); err != nil {
			return fmt.Errorf("failed to load post-build hooks: %v", err)
		}

		// Load failure hooks
		failureDir := filepath.Join(hookDir, "failure")
		if err := hm.loadHooksFromDir(failureDir, HookStageBuildFailure, projectDir); err != nil {
			return fmt.Errorf("failed to load failure hooks: %v", err)
		}
	}

	return nil
}

// loadHooksFromDir loads hook scripts from a directory
func (hm *HookManager) loadHooksFromDir(dir string, stage HookStage, projectDir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read hook directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !isExecutable(entry) {
			continue
		}

		hookPath := filepath.Join(dir, entry.Name())
		hook := BuildHook{
			Name:        entry.Name(),
			Description: fmt.Sprintf("Hook script: %s", entry.Name()),
			Stage:       stage,
			Command:     hookPath,
			Timeout:     5 * time.Minute, // Default 5 minute timeout
			WorkingDir:  projectDir,
			Environment: map[string]string{
				"FORGE_PROJECT_DIR": projectDir,
				"FORGE_HOOK_STAGE":  stage.String(),
			},
		}

		hm.AddHook(hook)
	}

	return nil
}

// isExecutable checks if a file is executable
func isExecutable(entry os.DirEntry) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}

	mode := info.Mode()
	return !mode.IsDir() && (mode.Perm()&0111 != 0)
}

// ExecuteHooks executes all hooks for a given stage
func (hm *HookManager) ExecuteHooks(ctx context.Context, stage HookStage, phase string) error {
	var relevantHooks []BuildHook
	for _, hook := range hm.hooks {
		if hook.Stage == stage || (stage == HookStagePrePhase || stage == HookStagePostPhase) {
			relevantHooks = append(relevantHooks, hook)
		}
	}

	if len(relevantHooks) == 0 {
		return nil
	}

	hm.logger.Info("Executing hooks", "stage", stage.String(), "count", len(relevantHooks))

	for _, hook := range relevantHooks {
		if err := hm.executeHook(ctx, hook, phase); err != nil {
			hm.logger.Error("Hook execution failed", "hook", hook.Name, "error", err)
			// For pre/post phase hooks, we might want to continue on failure
			// For build hooks, we might want to fail the build
			if stage == HookStagePreBuild || stage == HookStagePostBuild || stage == HookStageBuildFailure {
				return fmt.Errorf("hook %s failed: %v", hook.Name, err)
			}
		}
	}

	return nil
}

// executeHook executes a single hook
func (hm *HookManager) executeHook(ctx context.Context, hook BuildHook, phase string) error {
	hm.logger.Info("Executing hook", "name", hook.Name, "command", hook.Command)

	// Start timing
	timer := hm.metrics.StartTimer(fmt.Sprintf("hook_%s", hook.Name))
	defer timer.Stop()

	// Prepare command
	cmd := exec.CommandContext(ctx, hook.Command)
	cmd.Dir = hook.WorkingDir

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range hook.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	if phase != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PHASE=%s", phase))
	}

	// Set up timeout if specified
	if hook.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, hook.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, hook.Command)
		cmd.Dir = hook.WorkingDir
		cmd.Env = os.Environ()
		for key, value := range hook.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
		if phase != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PHASE=%s", phase))
		}
	}

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		hm.logger.Error("Hook failed", "name", hook.Name, "output", string(output), "error", err)
		return fmt.Errorf("hook execution failed: %v", err)
	}

	hm.logger.Info("Hook completed successfully", "name", hook.Name)
	if len(output) > 0 {
		hm.logger.Debug("Hook output", "name", hook.Name, "output", string(output))
	}

	return nil
}

// GetHooks returns all registered hooks
func (hm *HookManager) GetHooks() []BuildHook {
	return hm.hooks
}

// GetHooksByStage returns hooks for a specific stage
func (hm *HookManager) GetHooksByStage(stage HookStage) []BuildHook {
	var hooks []BuildHook
	for _, hook := range hm.hooks {
		if hook.Stage == stage {
			hooks = append(hooks, hook)
		}
	}
	return hooks
}

// String returns the string representation of a hook stage
func (hs HookStage) String() string {
	switch hs {
	case HookStagePreBuild:
		return "pre-build"
	case HookStagePostBuild:
		return "post-build"
	case HookStageBuildFailure:
		return "build-failure"
	case HookStagePrePhase:
		return "pre-phase"
	case HookStagePostPhase:
		return "post-phase"
	default:
		return "unknown"
	}
}
