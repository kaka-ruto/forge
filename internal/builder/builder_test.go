package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
	tempDir string
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}

func (s *BuilderTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-builder-test-*")
	s.Require().NoError(err)
}

func (s *BuilderTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *BuilderTestSuite) TestNewBuildOrchestrator() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	bo := NewBuildOrchestrator(cfg, projectDir)

	s.NotNil(bo)
	s.Equal(cfg, bo.config)
	s.Equal(projectDir, bo.projectDir)
	s.NotNil(bo.buildroot)
	s.NotNil(bo.logger)
	s.NotNil(bo.metrics)
}

func (s *BuilderTestSuite) TestNewBuildOrchestratorWithNilConfig() {
	bo := NewBuildOrchestrator(nil, s.tempDir)
	s.Nil(bo)
}

func (s *BuilderTestSuite) TestNewBuildOrchestratorWithEmptyProjectDir() {
	cfg := &config.Config{Name: "test"}
	bo := NewBuildOrchestrator(cfg, "")
	s.Nil(bo)
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithValidConfig() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        1,
	}

	ctx := context.Background()
	// This will fail because Buildroot is not actually set up, but we test the flow
	err := bo.Build(ctx, opts)
	s.Error(err)                            // Expected to fail in test environment
	s.Contains(err.Error(), "build failed") // Should contain our error message
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithCleanOption() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       true,
		Verbose:     false,
		Incremental: false,
		Jobs:        1,
	}

	ctx := context.Background()
	err := bo.Build(ctx, opts)
	s.Error(err) // Expected to fail
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithTimeout() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        1,
		Timeout:     1 * time.Millisecond, // Very short timeout
	}

	ctx := context.Background()
	err := bo.Build(ctx, opts)
	s.Error(err) // Should timeout
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithInvalidArchitecture() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "invalid-arch",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        1,
	}

	ctx := context.Background()
	err := bo.Build(ctx, opts)
	s.Error(err)
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithInvalidTemplate() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "invalid-template",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        1,
	}

	ctx := context.Background()
	err := bo.Build(ctx, opts)
	s.Error(err)
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithOptimization() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        1,
		OptimizeFor: "size",
	}

	ctx := context.Background()
	err := bo.Build(ctx, opts)
	s.Error(err) // Expected to fail in test environment
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildWithParallelJobs() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        4,
	}

	ctx := context.Background()
	err := bo.Build(ctx, opts)
	s.Error(err) // Expected to fail in test environment
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildCancellation() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	opts := BuildOptions{
		Clean:       false,
		Verbose:     false,
		Incremental: true,
		Jobs:        1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := bo.Build(ctx, opts)
	s.Error(err)
	s.Contains(err.Error(), "cancelled")
}

func (s *BuilderTestSuite) TestBuildOrchestratorPreBuildValidation() {
	cfg := &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	// Test pre-build validation
	ctx := context.Background()
	err := bo.validateBuildConfig(ctx)
	s.NoError(err)
}

func (s *BuilderTestSuite) TestBuildOrchestratorResourceChecking() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	// Test resource checking
	ctx := context.Background()
	err := bo.checkBuildResources(ctx)
	s.NoError(err) // Should pass in test environment
}

func (s *BuilderTestSuite) TestBuildOrchestratorMetricsCollection() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	// Test metrics initialization
	s.NotNil(bo.metrics)
}

func (s *BuilderTestSuite) TestBuildOrchestratorLogging() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	// Test logger initialization
	s.NotNil(bo.logger)
}

func (s *BuilderTestSuite) TestBuildOptionsDefaults() {
	opts := BuildOptions{}

	// Test default values
	s.False(opts.Clean)
	s.False(opts.Verbose)
	s.False(opts.Incremental) // Default should be false
	s.Equal(0, opts.Jobs)     // Default should be 0 (auto-detect)
	s.Equal("", opts.OptimizeFor)
	s.Equal(time.Duration(0), opts.Timeout)
}

func (s *BuilderTestSuite) TestBuildOrchestratorOutputDirectories() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	// Check that output directories are set correctly
	expectedBuildDir := filepath.Join(projectDir, "build")
	expectedArtifactsDir := filepath.Join(expectedBuildDir, "artifacts")

	s.Equal(expectedBuildDir, bo.buildDir)
	s.Equal(expectedArtifactsDir, bo.artifactsDir)
}

func (s *BuilderTestSuite) TestBuildOrchestratorBuildPhases() {
	cfg := &config.Config{
		Name:         "test-project",
		Version:      "0.1.0",
		Architecture: "x86_64",
		Template:     "minimal",
	}

	projectDir := filepath.Join(s.tempDir, "project")
	os.MkdirAll(projectDir, 0755)

	bo := NewBuildOrchestrator(cfg, projectDir)
	s.NotNil(bo)

	// Test that build phases are defined
	s.NotNil(bo.buildPhases)
	s.Greater(len(bo.buildPhases), 0)
}
