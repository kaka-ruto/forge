package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CleanupTestSuite struct {
	suite.Suite
	tempDir string
}

func TestCleanupTestSuite(t *testing.T) {
	suite.Run(t, new(CleanupTestSuite))
}

func (s *CleanupTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-cleanup-test-*")
	s.Require().NoError(err)
}

func (s *CleanupTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *CleanupTestSuite) TestResourceCleanerCreation() {
	cleaner := NewResourceCleaner()
	s.NotNil(cleaner)
}

func (s *CleanupTestSuite) TestCleanupDirectory() {
	// Create test directory with files
	testDir := filepath.Join(s.tempDir, "cleanup_test")
	err := os.MkdirAll(testDir, 0755)
	s.NoError(err)

	// Create some files
	oldFile := filepath.Join(testDir, "old_file.txt")
	err = os.WriteFile(oldFile, []byte("old content"), 0644)
	s.NoError(err)

	// Set modification time to 2 days ago
	oldTime := time.Now().Add(-48 * time.Hour)
	err = os.Chtimes(oldFile, oldTime, oldTime)
	s.NoError(err)

	// Create a new file
	newFile := filepath.Join(testDir, "new_file.txt")
	err = os.WriteFile(newFile, []byte("new content"), 0644)
	s.NoError(err)

	cleaner := NewResourceCleaner()

	// Cleanup files older than 24 hours
	cleaned, err := cleaner.CleanupDirectory(testDir, 24*time.Hour)
	s.NoError(err)
	s.Greater(cleaned, int64(0))

	// Check that old file is gone
	_, err = os.Stat(oldFile)
	s.True(os.IsNotExist(err))

	// Check that new file still exists
	_, err = os.Stat(newFile)
	s.NoError(err)
}

func (s *CleanupTestSuite) TestCleanupDirectoryDryRun() {
	// Create test directory with files
	testDir := filepath.Join(s.tempDir, "dry_run_test")
	err := os.MkdirAll(testDir, 0755)
	s.NoError(err)

	// Create old file
	oldFile := filepath.Join(testDir, "old_file.txt")
	err = os.WriteFile(oldFile, []byte("old content"), 0644)
	s.NoError(err)

	oldTime := time.Now().Add(-48 * time.Hour)
	err = os.Chtimes(oldFile, oldTime, oldTime)
	s.NoError(err)

	cleaner := NewResourceCleaner()

	// Dry run cleanup
	wouldClean, err := cleaner.CleanupDirectoryDryRun(testDir, 24*time.Hour)
	s.NoError(err)
	s.Greater(wouldClean, int64(0))

	// File should still exist
	_, err = os.Stat(oldFile)
	s.NoError(err)
}

func (s *CleanupTestSuite) TestCleanupBuildArtifacts() {
	// Create test build directory
	buildDir := filepath.Join(s.tempDir, "build")
	err := os.MkdirAll(buildDir, 0755)
	s.NoError(err)

	// Create some build artifacts
	artifacts := []string{
		"test.img",
		"output.qcow2",
		"disk.raw",
		"build.iso",
		"debug.log",
		"temp.tmp",
	}

	for _, artifact := range artifacts {
		filePath := filepath.Join(buildDir, artifact)
		err = os.WriteFile(filePath, []byte("artifact content"), 0644)
		s.NoError(err)
	}

	// Create a file that should not be cleaned
	keepFile := filepath.Join(buildDir, "config.yml")
	err = os.WriteFile(keepFile, []byte("config content"), 0644)
	s.NoError(err)

	cleaner := NewResourceCleaner()

	// Cleanup build artifacts
	err = cleaner.CleanupBuildArtifacts(buildDir)
	s.NoError(err)

	// Check that artifacts are gone
	for _, artifact := range artifacts {
		filePath := filepath.Join(buildDir, artifact)
		_, err = os.Stat(filePath)
		s.True(os.IsNotExist(err), "Artifact %s should be cleaned", artifact)
	}

	// Check that config file still exists
	_, err = os.Stat(keepFile)
	s.NoError(err)
}

func (s *CleanupTestSuite) TestCleanupCache() {
	// Create test cache directory
	cacheDir := filepath.Join(s.tempDir, "cache")
	err := os.MkdirAll(cacheDir, 0755)
	s.NoError(err)

	// Create files to fill cache
	for i := 0; i < 10; i++ {
		filePath := filepath.Join(cacheDir, fmt.Sprintf("cache_file_%d.dat", i))
		// Create 10MB files (simulating large cache)
		data := make([]byte, 10*1024*1024)
		for j := range data {
			data[j] = byte(i)
		}
		err = os.WriteFile(filePath, data, 0644)
		s.NoError(err)

		// Set different modification times
		modTime := time.Now().Add(time.Duration(-i) * time.Hour)
		err = os.Chtimes(filePath, modTime, modTime)
		s.NoError(err)
	}

	cleaner := NewResourceCleaner()

	// Cleanup cache to 50MB limit (each file is 10MB, so should keep ~5 files)
	err = cleaner.CleanupCache(cacheDir, 50) // 50MB limit
	s.NoError(err)

	// Check that some files were cleaned
	entries, err := os.ReadDir(cacheDir)
	s.NoError(err)
	s.Less(len(entries), 10, "Some cache files should have been cleaned")
}

func (s *CleanupTestSuite) TestCleanupForgeDirectories() {
	// Create test .forge directory structure
	forgeDir := filepath.Join(s.tempDir, ".forge")
	metricsDir := filepath.Join(forgeDir, "metrics")
	logsDir := filepath.Join(forgeDir, "logs")
	cacheDir := filepath.Join(forgeDir, "cache")

	err := os.MkdirAll(metricsDir, 0755)
	s.NoError(err)
	err = os.MkdirAll(logsDir, 0755)
	s.NoError(err)
	err = os.MkdirAll(cacheDir, 0755)
	s.NoError(err)

	// Create old files (31 days old)
	oldTime := time.Now().Add(-31 * 24 * time.Hour)

	oldMetricFile := filepath.Join(metricsDir, "metrics_old.json")
	err = os.WriteFile(oldMetricFile, []byte("{}"), 0644)
	s.NoError(err)
	err = os.Chtimes(oldMetricFile, oldTime, oldTime)
	s.NoError(err)

	oldLogFile := filepath.Join(logsDir, "old.log")
	err = os.WriteFile(oldLogFile, []byte("old log"), 0644)
	s.NoError(err)
	err = os.Chtimes(oldLogFile, oldTime, oldTime)
	s.NoError(err)

	oldCacheFile := filepath.Join(cacheDir, "old.cache")
	err = os.WriteFile(oldCacheFile, []byte("old cache"), 0644)
	s.NoError(err)
	err = os.Chtimes(oldCacheFile, oldTime, oldTime)
	s.NoError(err)

	// Create new files
	newMetricFile := filepath.Join(metricsDir, "metrics_new.json")
	err = os.WriteFile(newMetricFile, []byte("{}"), 0644)
	s.NoError(err)

	cleaner := NewResourceCleaner()

	// Cleanup forge directories
	err = cleaner.CleanupForgeDirectories(s.tempDir)
	s.NoError(err)

	// Check that old files are gone
	_, err = os.Stat(oldMetricFile)
	s.True(os.IsNotExist(err))

	_, err = os.Stat(oldLogFile)
	s.True(os.IsNotExist(err))

	_, err = os.Stat(oldCacheFile)
	s.True(os.IsNotExist(err))

	// Check that new files still exist
	_, err = os.Stat(newMetricFile)
	s.NoError(err)
}

func (s *CleanupTestSuite) TestCleanupNonexistentDirectory() {
	cleaner := NewResourceCleaner()

	// Try to cleanup non-existent directory
	cleaned, err := cleaner.CleanupDirectory("/nonexistent/path", 24*time.Hour)
	s.Error(err)
	s.Equal(int64(0), cleaned)
}

func (s *CleanupTestSuite) TestCleanupCacheInvalidDirectory() {
	cleaner := NewResourceCleaner()

	// Try to cleanup cache on a file instead of directory
	err := cleaner.CleanupCache("/etc/hosts", 10)
	s.Error(err)
}
