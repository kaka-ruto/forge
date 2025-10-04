package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
	"github.com/sst/forge/internal/metrics"
	"github.com/stretchr/testify/suite"
)

type HooksTestSuite struct {
	suite.Suite
	tempDir string
	logger  *logger.Logger
	metrics *metrics.MetricsCollector
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

func (s *HooksTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-hooks-test-*")
	s.Require().NoError(err)

	s.logger = logger.NewLogger(logger.INFO, os.Stdout, os.Stderr)
	s.metrics = metrics.NewMetricsCollector()
}

func (s *HooksTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *HooksTestSuite) TestNewHookManager() {
	hm := NewHookManager(s.logger, s.metrics)
	s.NotNil(hm)
	s.NotNil(hm.hooks)
	s.Equal(s.logger, hm.logger)
	s.Equal(s.metrics, hm.metrics)
	s.Empty(hm.hooks)
}

func (s *HooksTestSuite) TestAddHook() {
	hm := NewHookManager(s.logger, s.metrics)

	hook := BuildHook{
		Name:        "test-hook",
		Description: "Test hook",
		Stage:       HookStagePreBuild,
		Command:     "echo hello",
		Timeout:     1 * time.Minute,
	}

	hm.AddHook(hook)

	s.Len(hm.hooks, 1)
	s.Equal(hook, hm.hooks[0])
}

func (s *HooksTestSuite) TestHookStageString() {
	s.Equal("pre-build", HookStagePreBuild.String())
	s.Equal("post-build", HookStagePostBuild.String())
	s.Equal("build-failure", HookStageBuildFailure.String())
	s.Equal("pre-phase", HookStagePrePhase.String())
	s.Equal("post-phase", HookStagePostPhase.String())
	s.Equal("unknown", HookStage(999).String())
}

func (s *HooksTestSuite) TestGetHooks() {
	hm := NewHookManager(s.logger, s.metrics)

	hook1 := BuildHook{Name: "hook1", Stage: HookStagePreBuild}
	hook2 := BuildHook{Name: "hook2", Stage: HookStagePostBuild}

	hm.AddHook(hook1)
	hm.AddHook(hook2)

	hooks := hm.GetHooks()
	s.Len(hooks, 2)
	s.Contains(hooks, hook1)
	s.Contains(hooks, hook2)
}

func (s *HooksTestSuite) TestGetHooksByStage() {
	hm := NewHookManager(s.logger, s.metrics)

	hook1 := BuildHook{Name: "hook1", Stage: HookStagePreBuild}
	hook2 := BuildHook{Name: "hook2", Stage: HookStagePostBuild}
	hook3 := BuildHook{Name: "hook3", Stage: HookStagePreBuild}

	hm.AddHook(hook1)
	hm.AddHook(hook2)
	hm.AddHook(hook3)

	preBuildHooks := hm.GetHooksByStage(HookStagePreBuild)
	s.Len(preBuildHooks, 2)
	s.Contains(preBuildHooks, hook1)
	s.Contains(preBuildHooks, hook3)

	postBuildHooks := hm.GetHooksByStage(HookStagePostBuild)
	s.Len(postBuildHooks, 1)
	s.Equal(hook2, postBuildHooks[0])
}

func (s *HooksTestSuite) TestExecuteHooksNoHooks() {
	hm := NewHookManager(s.logger, s.metrics)

	ctx := context.Background()
	err := hm.ExecuteHooks(ctx, HookStagePreBuild, "")
	s.NoError(err)
}

func (s *HooksTestSuite) TestExecuteHooksWithSimpleCommand() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create a simple test script
	scriptPath := filepath.Join(s.tempDir, "test-hook.sh")
	scriptContent := `#!/bin/bash
echo "Hook executed"
exit 0
`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	s.NoError(err)

	hook := BuildHook{
		Name:       "test-hook",
		Stage:      HookStagePreBuild,
		Command:    scriptPath,
		WorkingDir: s.tempDir,
		Timeout:    10 * time.Second,
	}
	hm.AddHook(hook)

	ctx := context.Background()
	err = hm.ExecuteHooks(ctx, HookStagePreBuild, "")
	s.NoError(err)
}

func (s *HooksTestSuite) TestExecuteHooksWithFailingCommand() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create a failing test script
	scriptPath := filepath.Join(s.tempDir, "failing-hook.sh")
	scriptContent := `#!/bin/bash
echo "Hook failed" >&2
exit 1
`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	s.NoError(err)

	hook := BuildHook{
		Name:       "failing-hook",
		Stage:      HookStagePreBuild,
		Command:    scriptPath,
		WorkingDir: s.tempDir,
		Timeout:    10 * time.Second,
	}
	hm.AddHook(hook)

	ctx := context.Background()
	err = hm.ExecuteHooks(ctx, HookStagePreBuild, "")
	s.Error(err)
	s.Contains(err.Error(), "hook execution failed")
}

func (s *HooksTestSuite) TestExecuteHooksWithTimeout() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create a slow test script
	scriptPath := filepath.Join(s.tempDir, "slow-hook.sh")
	scriptContent := `#!/bin/bash
sleep 5
echo "Hook completed"
`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	s.NoError(err)

	hook := BuildHook{
		Name:       "slow-hook",
		Stage:      HookStagePreBuild,
		Command:    scriptPath,
		WorkingDir: s.tempDir,
		Timeout:    1 * time.Second, // Very short timeout
	}
	hm.AddHook(hook)

	ctx := context.Background()
	err = hm.ExecuteHooks(ctx, HookStagePreBuild, "")
	s.Error(err) // Should timeout
}

func (s *HooksTestSuite) TestExecuteHooksWithEnvironmentVariables() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create a script that checks environment variables
	scriptPath := filepath.Join(s.tempDir, "env-hook.sh")
	scriptContent := `#!/bin/bash
if [ "$FORGE_PROJECT_DIR" = "` + s.tempDir + `" ]; then
    echo "Environment variable set correctly"
    exit 0
else
    echo "Environment variable not set" >&2
    exit 1
fi
`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	s.NoError(err)

	hook := BuildHook{
		Name:       "env-hook",
		Stage:      HookStagePreBuild,
		Command:    scriptPath,
		WorkingDir: s.tempDir,
		Timeout:    10 * time.Second,
		Environment: map[string]string{
			"FORGE_PROJECT_DIR": s.tempDir,
		},
	}
	hm.AddHook(hook)

	ctx := context.Background()
	err = hm.ExecuteHooks(ctx, HookStagePreBuild, "")
	s.NoError(err)
}

func (s *HooksTestSuite) TestLoadHooksFromConfigNoHooksDir() {
	hm := NewHookManager(s.logger, s.metrics)

	cfg := &config.Config{Name: "test"}
	err := hm.LoadHooksFromConfig(cfg, s.tempDir)
	s.NoError(err)
	s.Empty(hm.hooks)
}

func (s *HooksTestSuite) TestLoadHooksFromConfigWithHooks() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create hooks directory structure
	hooksDir := filepath.Join(s.tempDir, "hooks")
	preBuildDir := filepath.Join(hooksDir, "pre-build")
	os.MkdirAll(preBuildDir, 0755)

	// Create a test hook script
	hookPath := filepath.Join(preBuildDir, "test-hook.sh")
	hookContent := `#!/bin/bash
echo "Test hook executed"
`
	err := os.WriteFile(hookPath, []byte(hookContent), 0755)
	s.NoError(err)

	cfg := &config.Config{Name: "test"}
	err = hm.LoadHooksFromConfig(cfg, s.tempDir)
	s.NoError(err)
	s.Len(hm.hooks, 1)

	hook := hm.hooks[0]
	s.Equal("test-hook.sh", hook.Name)
	s.Equal(HookStagePreBuild, hook.Stage)
	s.Equal(hookPath, hook.Command)
	s.Equal(s.tempDir, hook.WorkingDir)
}

func (s *HooksTestSuite) TestLoadHooksFromConfigMultipleStages() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create hooks directory structure
	hooksDir := filepath.Join(s.tempDir, "hooks")
	preBuildDir := filepath.Join(hooksDir, "pre-build")
	postBuildDir := filepath.Join(hooksDir, "post-build")
	failureDir := filepath.Join(hooksDir, "failure")

	os.MkdirAll(preBuildDir, 0755)
	os.MkdirAll(postBuildDir, 0755)
	os.MkdirAll(failureDir, 0755)

	// Create hook scripts
	preHookPath := filepath.Join(preBuildDir, "pre-hook.sh")
	postHookPath := filepath.Join(postBuildDir, "post-hook.sh")
	failHookPath := filepath.Join(failureDir, "fail-hook.sh")

	hookContent := `#!/bin/bash
echo "Hook executed"
`

	err := os.WriteFile(preHookPath, []byte(hookContent), 0755)
	s.NoError(err)
	err = os.WriteFile(postHookPath, []byte(hookContent), 0755)
	s.NoError(err)
	err = os.WriteFile(failHookPath, []byte(hookContent), 0755)
	s.NoError(err)

	cfg := &config.Config{Name: "test"}
	err = hm.LoadHooksFromConfig(cfg, s.tempDir)
	s.NoError(err)
	s.Len(hm.hooks, 3)

	// Check that hooks are loaded with correct stages
	stages := make(map[HookStage]int)
	for _, hook := range hm.hooks {
		stages[hook.Stage]++
	}

	s.Equal(1, stages[HookStagePreBuild])
	s.Equal(1, stages[HookStagePostBuild])
	s.Equal(1, stages[HookStageBuildFailure])
}

func (s *HooksTestSuite) TestIsExecutable() {
	// Create a regular file
	regularFile := filepath.Join(s.tempDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("test"), 0644)
	s.NoError(err)

	// Create an executable file
	execFile := filepath.Join(s.tempDir, "executable.sh")
	err = os.WriteFile(execFile, []byte("#!/bin/bash\necho test"), 0755)
	s.NoError(err)

	// Create a directory
	testDir := filepath.Join(s.tempDir, "testdir")
	os.MkdirAll(testDir, 0755)

	// Test entries
	entries, err := os.ReadDir(s.tempDir)
	s.NoError(err)

	for _, entry := range entries {
		switch entry.Name() {
		case "regular.txt":
			s.False(isExecutable(entry))
		case "executable.sh":
			s.True(isExecutable(entry))
		case "testdir":
			s.False(isExecutable(entry)) // Directories are not executable in this context
		}
	}
}

func (s *HooksTestSuite) TestExecuteHooksCancellation() {
	hm := NewHookManager(s.logger, s.metrics)

	// Create a slow script
	scriptPath := filepath.Join(s.tempDir, "slow-hook.sh")
	scriptContent := `#!/bin/bash
sleep 10
`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	s.NoError(err)

	hook := BuildHook{
		Name:       "slow-hook",
		Stage:      HookStagePreBuild,
		Command:    scriptPath,
		WorkingDir: s.tempDir,
		Timeout:    30 * time.Second, // Long timeout, but we'll cancel
	}
	hm.AddHook(hook)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = hm.ExecuteHooks(ctx, HookStagePreBuild, "")
	s.Error(err)
	s.Contains(err.Error(), "context canceled")
}
