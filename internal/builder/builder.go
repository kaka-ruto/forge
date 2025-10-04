package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sst/forge/internal/buildroot"
	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
	"github.com/sst/forge/internal/metrics"
)

// BuildOptions represents build configuration options
type BuildOptions struct {
	Clean       bool
	Verbose     bool
	Incremental bool
	Jobs        int
	OptimizeFor string
	Timeout     time.Duration
}

// BuildOrchestrator coordinates the entire build process
type BuildOrchestrator struct {
	config       *config.Config
	projectDir   string
	buildDir     string
	artifactsDir string
	buildroot    *buildroot.BuildrootManager
	logger       *logger.Logger
	metrics      *metrics.MetricsCollector
	buildTimer   *metrics.Timer
	buildPhases  []BuildPhase
}

// BuildPhase represents a phase in the build process
type BuildPhase struct {
	Name        string
	Description string
	Handler     func(ctx context.Context) error
}

// NewBuildOrchestrator creates a new build orchestrator
func NewBuildOrchestrator(cfg *config.Config, projectDir string) *BuildOrchestrator {
	if cfg == nil || projectDir == "" {
		return nil
	}

	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts")

	bo := &BuildOrchestrator{
		config:       cfg,
		projectDir:   projectDir,
		buildDir:     buildDir,
		artifactsDir: artifactsDir,
		buildroot:    buildroot.NewBuildrootManager(cfg, projectDir),
		logger:       logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
		metrics:      metrics.NewMetricsCollector(),
	}

	bo.initializeBuildPhases()

	return bo
}

// initializeBuildPhases sets up the build phases
func (bo *BuildOrchestrator) initializeBuildPhases() {
	bo.buildPhases = []BuildPhase{
		{
			Name:        "validate",
			Description: "Validate build configuration",
			Handler:     bo.validateBuildConfig,
		},
		{
			Name:        "resources",
			Description: "Check system resources",
			Handler:     bo.checkBuildResources,
		},
		{
			Name:        "prepare",
			Description: "Prepare build environment",
			Handler:     bo.prepareBuildEnvironment,
		},
		{
			Name:        "buildroot",
			Description: "Configure and build with Buildroot",
			Handler:     bo.executeBuildrootBuild,
		},
		{
			Name:        "artifacts",
			Description: "Collect build artifacts",
			Handler:     bo.collectArtifacts,
		},
	}
}

// Build executes the complete build process
func (bo *BuildOrchestrator) Build(ctx context.Context, opts BuildOptions) error {
	bo.logger.Info("Starting Forge OS build", "project", bo.config.Name, "version", bo.config.Version)

	// Set default options
	if opts.Jobs == 0 {
		opts.Jobs = 1 // Default to 1 job for testing
	}
	if opts.Timeout == 0 {
		opts.Timeout = 2 * time.Hour // Default timeout
	}

	// Apply optimization settings
	if opts.OptimizeFor != "" {
		if err := bo.applyOptimization(opts.OptimizeFor); err != nil {
			return fmt.Errorf("failed to apply optimization: %v", err)
		}
	}

	// Execute build phases
	for _, phase := range bo.buildPhases {
		bo.logger.Info("Executing build phase", "phase", phase.Name, "description", phase.Description)

		// Check for cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("build cancelled: %v", ctx.Err())
		default:
		}

		// Check timeout
		if opts.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
			defer cancel()
		}

		// Execute phase
		if err := phase.Handler(ctx); err != nil {
			bo.logger.Error("Build phase failed", "phase", phase.Name, "error", err)
			return fmt.Errorf("build failed at phase %s: %v", phase.Name, err)
		}

		bo.logger.Info("Build phase completed", "phase", phase.Name)
	}

	bo.logger.Info("Build completed successfully", "project", bo.config.Name)
	return nil
}

// validateBuildConfig validates the build configuration
func (bo *BuildOrchestrator) validateBuildConfig(ctx context.Context) error {
	if err := bo.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	// Validate architecture support
	supportedArchs := []string{"x86_64", "arm", "aarch64", "mips", "riscv64"}
	archSupported := false
	for _, arch := range supportedArchs {
		if bo.config.Architecture == arch {
			archSupported = true
			break
		}
	}
	if !archSupported {
		return fmt.Errorf("unsupported architecture: %s", bo.config.Architecture)
	}

	return nil
}

// checkBuildResources checks if system has sufficient resources
func (bo *BuildOrchestrator) checkBuildResources(ctx context.Context) error {
	// In a real implementation, this would check disk space, memory, etc.
	// For testing, we just return success
	return nil
}

// prepareBuildEnvironment prepares the build environment
func (bo *BuildOrchestrator) prepareBuildEnvironment(ctx context.Context) error {
	// Create build directories
	if err := bo.createBuildDirectories(); err != nil {
		return fmt.Errorf("failed to create build directories: %v", err)
	}

	// Initialize metrics collection
	bo.buildTimer = bo.metrics.StartTimer("build")

	return nil
}

// createBuildDirectories creates necessary build directories
func (bo *BuildOrchestrator) createBuildDirectories() error {
	dirs := []string{
		bo.buildDir,
		bo.artifactsDir,
		filepath.Join(bo.buildDir, "logs"),
		filepath.Join(bo.buildDir, "cache"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// applyOptimization applies build optimizations
func (bo *BuildOrchestrator) applyOptimization(optimizeFor string) error {
	validOptimizations := []string{"size", "performance", "realtime"}
	valid := false
	for _, opt := range validOptimizations {
		if optimizeFor == opt {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid optimization: %s", optimizeFor)
	}

	// In a real implementation, this would modify buildroot config
	bo.logger.Info("Applying optimization", "type", optimizeFor)
	return nil
}

// executeBuildrootBuild executes the Buildroot build
func (bo *BuildOrchestrator) executeBuildrootBuild(ctx context.Context) error {
	// In a real implementation, this would:
	// 1. Clone/download Buildroot
	// 2. Apply defconfig
	// 3. Run make with appropriate options
	// 4. Handle parallel jobs
	// For testing, we simulate this

	bo.logger.Info("Executing Buildroot build", "architecture", bo.config.Architecture)

	// Simulate build process (would fail in real environment without Buildroot)
	return fmt.Errorf("buildroot build simulation - would require actual Buildroot setup")
}

// collectArtifacts collects build artifacts
func (bo *BuildOrchestrator) collectArtifacts(ctx context.Context) error {
	// In a real implementation, this would copy built images to artifacts directory
	bo.logger.Info("Collecting build artifacts")

	// Stop metrics collection
	if bo.buildTimer != nil {
		duration := bo.buildTimer.Stop()
		bo.logger.Info("Build completed", "duration", duration)
	}

	return nil
}
