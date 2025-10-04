package buildroot

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type BuildrootTestSuite struct {
	suite.Suite
	tempDir string
	config  *config.Config
}

func TestBuildrootTestSuite(t *testing.T) {
	suite.Run(t, new(BuildrootTestSuite))
}

func (s *BuildrootTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-buildroot-test-*")
	s.Require().NoError(err)

	// Create a test config
	s.config = &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Buildroot: config.BuildrootConfig{
			Version: "stable",
		},
		Packages: []string{},
		Features: []string{},
	}
}

func (s *BuildrootTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *BuildrootTestSuite) TestNewBuildrootManager() {
	bm := NewBuildrootManager(s.config, s.tempDir)
	s.NotNil(bm)
	s.Equal(s.config, bm.config)
	s.Equal(s.tempDir, bm.projectDir)
	s.Equal(filepath.Join(s.tempDir, "build"), bm.buildDir)
}

func (s *BuildrootTestSuite) TestDownloadBuildroot() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// This test would actually download Buildroot, which is slow and requires internet
	// For now, we'll just test that the method doesn't panic and creates the build directory
	err := bm.DownloadBuildroot()
	// We expect this to fail in test environment without internet, but the directory should be created
	if err != nil {
		s.Contains(err.Error(), "failed to download") // Expected in test environment
	}

	// Check that build directory was created
	_, err = os.Stat(bm.buildDir)
	s.NoError(err)
}

func (s *BuildrootTestSuite) TestGenerateConfig() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create mock Buildroot directory structure
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	err := os.MkdirAll(buildrootDir, 0755)
	s.NoError(err)

	// Create a mock .config file
	configPath := filepath.Join(buildrootDir, ".config")
	err = os.WriteFile(configPath, []byte("# Test config\n"), 0644)
	s.NoError(err)

	// Test config generation
	err = bm.GenerateConfig()
	// This will fail because we don't have a real Buildroot setup, but it should not panic
	s.Error(err) // Expected to fail without real Buildroot
}

func (s *BuildrootTestSuite) TestBuild() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create mock Buildroot directory structure
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	err := os.MkdirAll(buildrootDir, 0755)
	s.NoError(err)

	// Test build (will fail without real Buildroot)
	err = bm.Build()
	s.Error(err) // Expected to fail without real Buildroot
}

func (s *BuildrootTestSuite) TestGetOutputDir() {
	bm := NewBuildrootManager(s.config, s.tempDir)
	expected := filepath.Join(s.tempDir, "build", "buildroot", "output")
	s.Equal(expected, bm.GetOutputDir())
}

func (s *BuildrootTestSuite) TestGetImagesDir() {
	bm := NewBuildrootManager(s.config, s.tempDir)
	expected := filepath.Join(s.tempDir, "build", "buildroot", "output", "images")
	s.Equal(expected, bm.GetImagesDir())
}

func (s *BuildrootTestSuite) TestApplyArchitectureConfig() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create mock config file
	configPath := filepath.Join(bm.buildDir, "buildroot", ".config")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	s.NoError(err)
	err = os.WriteFile(configPath, []byte("# Base config\n"), 0644)
	s.NoError(err)

	// Test x86_64 architecture
	s.config.Architecture = "x86_64"
	err = bm.applyArchitectureConfig()
	s.NoError(err)

	content, err := os.ReadFile(configPath)
	s.NoError(err)
	s.Contains(string(content), "BR2_x86_64=y")
	s.Contains(string(content), "BR2_ARCH=\"x86_64\"")
}

func (s *BuildrootTestSuite) TestApplyPackageConfig() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create mock config file
	configPath := filepath.Join(bm.buildDir, "buildroot", ".config")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	s.NoError(err)
	err = os.WriteFile(configPath, []byte("# Base config\n"), 0644)
	s.NoError(err)

	// Test with some packages
	s.config.Packages = []string{"openssh", "python3"}
	err = bm.applyPackageConfig()
	s.NoError(err)

	content, err := os.ReadFile(configPath)
	s.NoError(err)
	s.Contains(string(content), "BR2_PACKAGE_OPENSSH=y")
	s.Contains(string(content), "BR2_PACKAGE_PYTHON3=y")
}

func (s *BuildrootTestSuite) TestApplyFeatureConfig() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create mock config file
	configPath := filepath.Join(bm.buildDir, "buildroot", ".config")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	s.NoError(err)
	err = os.WriteFile(configPath, []byte("# Base config\n"), 0644)
	s.NoError(err)

	// Test with systemd feature
	s.config.Features = []string{"systemd"}
	err = bm.applyFeatureConfig()
	s.NoError(err)

	content, err := os.ReadFile(configPath)
	s.NoError(err)
	s.Contains(string(content), "BR2_INIT_SYSTEMD=y")
}

func (s *BuildrootTestSuite) TestAppendConfigLines() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create the buildroot directory structure
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	err := os.MkdirAll(buildrootDir, 0755)
	s.NoError(err)

	configPath := filepath.Join(buildrootDir, ".config")
	err = os.WriteFile(configPath, []byte("# Base config\n"), 0644)
	s.NoError(err)

	lines := []string{"BR2_TEST1=y", "BR2_TEST2=y"}
	err = bm.appendConfigLines(configPath, lines)
	s.NoError(err)

	content, err := os.ReadFile(configPath)
	s.NoError(err)
	s.Contains(string(content), "BR2_TEST1=y")
	s.Contains(string(content), "BR2_TEST2=y")
}

func (s *BuildrootTestSuite) TestFindExtractedBuildrootDir() {
	bm := NewBuildrootManager(s.config, s.tempDir)

	// Create mock extracted directory
	extractedDir := filepath.Join(bm.buildDir, "buildroot-2023.11")
	err := os.MkdirAll(extractedDir, 0755)
	s.NoError(err)

	result := bm.findExtractedBuildrootDir()
	s.Equal("buildroot-2023.11", result)
}

func (s *BuildrootTestSuite) TestGetParallelJobs() {
	bm := NewBuildrootManager(s.config, s.tempDir)
	jobs := bm.getParallelJobs()
	s.Equal(4, jobs) // Currently hardcoded to 4
}
