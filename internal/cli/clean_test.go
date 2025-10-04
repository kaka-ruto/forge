package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CleanCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestCleanCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CleanCommandTestSuite))
}

func (s *CleanCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-clean-test-*")
	s.Require().NoError(err)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	s.tempDir = oldDir + "/" + s.tempDir
	os.Chdir(s.tempDir)
}

func (s *CleanCommandTestSuite) TearDownTest() {
	os.Chdir("/")
	os.RemoveAll(s.tempDir)
}

func (s *CleanCommandTestSuite) TestNewCleanCommand() {
	cmd := NewCleanCommand()
	s.NotNil(cmd)
	s.Equal("clean", cmd.Use)
	s.Contains(cmd.Short, "Clean build artifacts")
}

func (s *CleanCommandTestSuite) TestCleanCommandNothingToClean() {
	err := runCleanCommand([]string{}, map[string]interface{}{
		"all":     false,
		"cache":   false,
		"builds":  false,
		"logs":    false,
		"dry-run": false,
	})
	s.NoError(err)
}

func (s *CleanCommandTestSuite) TestCleanCommandWithFiles() {
	// Create some test files and directories
	err := os.MkdirAll("build", 0755)
	s.NoError(err)
	err = os.WriteFile("build/test.img", []byte("test"), 0644)
	s.NoError(err)

	err = os.MkdirAll("dl", 0755)
	s.NoError(err)
	err = os.WriteFile("dl/package.tar.gz", []byte("test"), 0644)
	s.NoError(err)

	err = os.WriteFile("forge.log", []byte("test log"), 0644)
	s.NoError(err)

	// Test dry run
	err = runCleanCommand([]string{}, map[string]interface{}{
		"all":     true,
		"cache":   false,
		"builds":  false,
		"logs":    false,
		"dry-run": true,
	})
	s.NoError(err)

	// Files should still exist
	s.FileExists("build/test.img")
	s.FileExists("dl/package.tar.gz")
	s.FileExists("forge.log")

	// Test actual clean
	err = runCleanCommand([]string{}, map[string]interface{}{
		"all":     true,
		"cache":   false,
		"builds":  false,
		"logs":    false,
		"dry-run": false,
	})
	s.NoError(err)

	// Files should be gone
	s.NoFileExists("build/test.img")
	s.NoFileExists("dl/package.tar.gz")
	s.NoFileExists("forge.log")
}

func (s *CleanCommandTestSuite) TestCleanCommandSelective() {
	// Create test files
	err := os.MkdirAll("build", 0755)
	s.NoError(err)
	err = os.WriteFile("build/test.img", []byte("test"), 0644)
	s.NoError(err)

	err = os.MkdirAll("dl", 0755)
	s.NoError(err)
	err = os.WriteFile("dl/package.tar.gz", []byte("test"), 0644)
	s.NoError(err)

	// Clean only builds
	err = runCleanCommand([]string{}, map[string]interface{}{
		"all":     false,
		"cache":   false,
		"builds":  true,
		"logs":    false,
		"dry-run": false,
	})
	s.NoError(err)

	// Build files should be gone, cache files should remain
	s.NoFileExists("build/test.img")
	s.FileExists("dl/package.tar.gz")
}

func (s *CleanCommandTestSuite) TestFormatSize() {
	// Test various sizes
	s.Equal("0 B", formatSize(0))
	s.Equal("512 B", formatSize(512))
	s.Equal("1.0 KB", formatSize(1024))
	s.Equal("1.5 MB", formatSize(1024*1024+512*1024))
	s.Equal("2.0 GB", formatSize(2*1024*1024*1024))
}
