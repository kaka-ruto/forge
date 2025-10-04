package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ResourcesTestSuite struct {
	suite.Suite
	tempDir string
}

func TestResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(ResourcesTestSuite))
}

func (s *ResourcesTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-resources-test-*")
	s.Require().NoError(err)
}

func (s *ResourcesTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ResourcesTestSuite) TestDiskSpaceChecking() {
	checker := NewResourceChecker()

	// Test disk space check
	info, err := checker.CheckDiskSpace("/")
	s.NoError(err)
	s.Greater(info.AvailableBytes, int64(0))
	s.Greater(info.TotalBytes, int64(0))
	s.True(info.AvailableBytes <= info.TotalBytes)
}

func (s *ResourcesTestSuite) TestAvailableDiskSpaceCalculation() {
	checker := NewResourceChecker()

	info, err := checker.CheckDiskSpace("/")
	s.NoError(err)

	// Available space should be less than total space
	s.True(info.AvailableBytes < info.TotalBytes)
	// Available space should be positive
	s.True(info.AvailableBytes > 0)
}

func (s *ResourcesTestSuite) TestMemoryAvailabilityChecking() {
	checker := NewResourceChecker()

	info, err := checker.CheckMemory()
	s.NoError(err)
	s.Greater(info.TotalBytes, int64(0))
	s.GreaterOrEqual(info.AvailableBytes, int64(0))
}

func (s *ResourcesTestSuite) TestCPUCoreDetection() {
	checker := NewResourceChecker()

	cores := checker.GetCPUCount()
	s.Greater(cores, 0)
	s.LessOrEqual(cores, 128) // Reasonable upper bound
}

func (s *ResourcesTestSuite) TestResourceRequirementEstimation() {
	checker := NewResourceChecker()

	// Test minimal requirements
	reqs := checker.EstimateRequirements("minimal")
	s.Greater(reqs.MinDiskSpaceGB, 0)
	s.Greater(reqs.MinMemoryGB, 0)

	// Test networking requirements (should be higher)
	reqsNetworking := checker.EstimateRequirements("networking")
	s.GreaterOrEqual(reqsNetworking.MinDiskSpaceGB, reqs.MinDiskSpaceGB)
	s.GreaterOrEqual(reqsNetworking.MinMemoryGB, reqs.MinMemoryGB)
}

func (s *ResourcesTestSuite) TestResourceLimitEnforcement() {
	checker := NewResourceChecker()

	// Test with current system resources
	err := checker.ValidateRequirements("minimal")
	// This might fail on systems with insufficient resources, but shouldn't panic
	if err != nil {
		s.Contains(err.Error(), "insufficient")
	}
}

func (s *ResourcesTestSuite) TestCleanupOperations() {
	// Create some test files
	testDir := filepath.Join(s.tempDir, "test_cleanup")
	err := os.MkdirAll(testDir, 0755)
	s.NoError(err)

	// Create some files
	for i := 0; i < 5; i++ {
		filePath := filepath.Join(testDir, fmt.Sprintf("test_file_%d.txt", i))
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		s.NoError(err)
	}

	cleaner := NewResourceCleaner()

	// Test cleanup
	cleaned, err := cleaner.CleanupDirectory(testDir, 0) // Clean all files
	s.NoError(err)
	s.Greater(cleaned, int64(0))

	// Check that files are gone
	files, err := os.ReadDir(testDir)
	s.NoError(err)
	s.Equal(0, len(files))
}

func (s *ResourcesTestSuite) TestResourceWarnings() {
	checker := NewResourceChecker()

	warnings := checker.GetResourceWarnings()
	// Warnings slice should exist (may be empty)
	s.NotNil(warnings)
}

func (s *ResourcesTestSuite) TestResourceCheckerCreation() {
	checker := NewResourceChecker()
	s.NotNil(checker)
}

func (s *ResourcesTestSuite) TestDiskSpaceInfoValidation() {
	info := DiskSpaceInfo{
		Path:           "/",
		TotalBytes:     1000000,
		AvailableBytes: 500000,
	}

	s.True(info.IsValid())
	s.Equal("/", info.Path)
	s.Equal(int64(1000000), info.TotalBytes)
	s.Equal(int64(500000), info.AvailableBytes)
}

func (s *ResourcesTestSuite) TestMemoryInfoValidation() {
	info := MemoryInfo{
		TotalBytes:     8000000,
		AvailableBytes: 4000000,
	}

	s.True(info.IsValid())
	s.Equal(int64(8000000), info.TotalBytes)
	s.Equal(int64(4000000), info.AvailableBytes)
}

func (s *ResourcesTestSuite) TestResourceRequirementsValidation() {
	reqs := ResourceRequirements{
		MinDiskSpaceGB:         10,
		MinMemoryGB:            4,
		RecommendedDiskSpaceGB: 20,
		RecommendedMemoryGB:    8,
	}

	s.True(reqs.IsValid())
	s.Equal(10, reqs.MinDiskSpaceGB)
	s.Equal(4, reqs.MinMemoryGB)
	s.Equal(20, reqs.RecommendedDiskSpaceGB)
	s.Equal(8, reqs.RecommendedMemoryGB)
}

func (s *ResourcesTestSuite) TestResourceCleanerCreation() {
	cleaner := NewResourceCleaner()
	s.NotNil(cleaner)
}

func (s *ResourcesTestSuite) TestCleanupDryRun() {
	// Create test directory with files
	testDir := filepath.Join(s.tempDir, "dry_run_test")
	err := os.MkdirAll(testDir, 0755)
	s.NoError(err)

	// Create test file
	testFile := filepath.Join(testDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	s.NoError(err)

	cleaner := NewResourceCleaner()

	// Test dry run
	wouldClean, err := cleaner.CleanupDirectoryDryRun(testDir, 0)
	s.NoError(err)
	s.Greater(wouldClean, int64(0))

	// File should still exist
	_, err = os.Stat(testFile)
	s.NoError(err)
}

func (s *ResourcesTestSuite) TestResourceMonitoring() {
	monitor := NewResourceMonitor()

	// Start monitoring
	err := monitor.Start()
	s.NoError(err)

	// Get current usage
	usage := monitor.GetCurrentUsage()
	s.NotNil(usage)
	s.GreaterOrEqual(usage.MemoryBytes, int64(0))
	s.GreaterOrEqual(usage.CPUPercent, 0.0)

	// Stop monitoring
	monitor.Stop()
}

func (s *ResourcesTestSuite) TestResourceLimits() {
	limits := ResourceLimits{
		MaxMemoryGB:   8,
		MaxDiskGB:     50,
		MaxBuildTime:  7200, // 2 hours
		MaxConcurrent: 4,
	}

	s.True(limits.IsValid())
	s.Equal(8, limits.MaxMemoryGB)
	s.Equal(50, limits.MaxDiskGB)
	s.Equal(7200, limits.MaxBuildTime)
	s.Equal(4, limits.MaxConcurrent)
}

func (s *ResourcesTestSuite) TestResourceAlerting() {
	monitor := NewResourceMonitor()
	alerts := monitor.CheckAlerts()

	// Alerts should be a valid slice
	s.NotNil(alerts)
}

func (s *ResourcesTestSuite) TestBuildResourceEstimation() {
	estimator := NewBuildEstimator()

	// Test estimation for different project types
	for _, projectType := range []string{"minimal", "networking", "iot", "security", "industrial", "kiosk"} {
		estimate := estimator.EstimateBuildResources(projectType)
		s.NotNil(estimate)
		s.Greater(estimate.EstimatedTimeSeconds, 0)
		s.Greater(estimate.EstimatedMemoryGB, 0)
		s.Greater(estimate.EstimatedDiskGB, 0)
	}
}

func (s *ResourcesTestSuite) TestResourceHistoryTracking() {
	tracker := NewResourceTracker()

	// Record some resource usage
	tracker.RecordUsage(ResourceUsage{
		Timestamp:   time.Now(),
		MemoryBytes: 1000000,
		CPUPercent:  50.0,
		DiskBytes:   500000,
	})

	history := tracker.GetHistory()
	s.GreaterOrEqual(len(history), 1)
}

func (s *ResourcesTestSuite) TestResourceQuotaManagement() {
	quota := ResourceQuota{
		MemoryLimitGB:  4,
		DiskLimitGB:    10,
		TimeLimitHours: 2,
	}

	s.True(quota.IsValid())

	// Test quota checking
	usage := ResourceUsage{
		MemoryBytes: 2 * 1024 * 1024 * 1024, // 2GB
		DiskBytes:   5 * 1024 * 1024 * 1024, // 5GB
	}

	s.True(quota.IsWithinLimits(usage))
}
