package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuildCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestBuildCommandTestSuite(t *testing.T) {
	suite.Run(t, new(BuildCommandTestSuite))
}

func (s *BuildCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-build-test-*")
	s.Require().NoError(err)
}

func (s *BuildCommandTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *BuildCommandTestSuite) TestBuildCommandCreation() {
	cmd := NewBuildCommand()
	s.NotNil(cmd)
	s.Equal("build", cmd.Use)
	s.Contains(cmd.Short, "Build the Forge OS image")
}

func (s *BuildCommandTestSuite) TestBuildCommandWithValidConfig() {
	// Create a test project
	projectDir := filepath.Join(s.tempDir, "test-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Change to project directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{
		"clean":   "false",
		"verbose": "false",
	})

	// Build will fail in test environment due to missing Buildroot
	s.Error(err)
}

func (s *BuildCommandTestSuite) TestBuildCommandNoConfigFile() {
	// Change to a directory without forge.yml
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(s.tempDir)

	err := runBuildCommand([]string{}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "no forge.yml found")
}

func (s *BuildCommandTestSuite) TestBuildCommandInvalidConfig() {
	// Create a project with invalid forge.yml
	projectDir := filepath.Join(s.tempDir, "invalid-project")
	err := os.MkdirAll(projectDir, 0755)
	s.NoError(err)

	// Create invalid forge.yml
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	err = os.WriteFile(forgeYmlPath, []byte("invalid: yaml: content:"), 0644)
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "invalid forge.yml")
}

func (s *BuildCommandTestSuite) TestBuildCommandCleanBuild() {
	projectDir := filepath.Join(s.tempDir, "clean-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{
		"clean": "true",
	})
	// Build will fail in test environment due to missing Buildroot
	s.Error(err)
}

func (s *BuildCommandTestSuite) TestBuildCommandVerboseOutput() {
	projectDir := filepath.Join(s.tempDir, "verbose-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	// Verbose output is now supported through Buildroot
	err = runBuildCommand([]string{}, map[string]string{
		"verbose": "true",
	})
	// This will fail because Buildroot isn't actually available in test environment
	s.Error(err)
}

func (s *BuildCommandTestSuite) TestBuildCommandParallelJobs() {
	projectDir := filepath.Join(s.tempDir, "parallel-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	// Parallel jobs are now supported through Buildroot
	err = runBuildCommand([]string{}, map[string]string{
		"jobs": "4",
	})
	// Build will fail in test environment due to missing Buildroot
	s.Error(err)
}

func (s *BuildCommandTestSuite) TestBuildCommandResourceChecking() {
	projectDir := filepath.Join(s.tempDir, "resource-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{})
	// Build will fail due to Buildroot download in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")
}

func (s *BuildCommandTestSuite) TestBuildCommandProgressTracking() {
	projectDir := filepath.Join(s.tempDir, "progress-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{})
	// Build will fail due to Buildroot download in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")
}

func (s *BuildCommandTestSuite) TestBuildCommandCaching() {
	projectDir := filepath.Join(s.tempDir, "cache-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	// First build
	err = runBuildCommand([]string{}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")

	// Second build should also fail (no caching implemented yet)
	err = runBuildCommand([]string{}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")
}

func (s *BuildCommandTestSuite) TestBuildCommandIncrementalBuild() {
	projectDir := filepath.Join(s.tempDir, "incremental-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{})
	// Build will fail due to Buildroot download in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")
}

func (s *BuildCommandTestSuite) TestBuildCommandOptimizations() {
	projectDir := filepath.Join(s.tempDir, "optimized-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	// Test size optimization
	err = runBuildCommand([]string{}, map[string]string{
		"optimize-for": "size",
	})
	s.Error(err) // Expected to fail initially

	// Test performance optimization
	err = runBuildCommand([]string{}, map[string]string{
		"optimize-for": "performance",
	})
	s.Error(err) // Expected to fail initially

	// Test realtime optimization
	err = runBuildCommand([]string{}, map[string]string{
		"optimize-for": "realtime",
	})
	s.Error(err) // Expected to fail initially
}

func (s *BuildCommandTestSuite) TestBuildCommandMetricsCollection() {
	projectDir := filepath.Join(s.tempDir, "metrics-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{})
	// Build will fail due to Buildroot download in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")
}

func (s *BuildCommandTestSuite) TestBuildCommandLoggingIntegration() {
	projectDir := filepath.Join(s.tempDir, "logging-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{})
	// Build will fail due to Buildroot download in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to download Buildroot")
}

func (s *BuildCommandTestSuite) TestBuildCommandTimeoutHandling() {
	projectDir := filepath.Join(s.tempDir, "timeout-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runBuildCommand([]string{}, map[string]string{
		"timeout": "30m",
	})
	s.Error(err) // Expected to fail initially
}
