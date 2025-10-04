package resources

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
)

// DiskSpaceInfo represents disk space information
type DiskSpaceInfo struct {
	Path           string
	TotalBytes     int64
	AvailableBytes int64
}

// IsValid checks if the disk space info is valid
func (d DiskSpaceInfo) IsValid() bool {
	return d.TotalBytes > 0 && d.AvailableBytes >= 0 && d.AvailableBytes <= d.TotalBytes
}

// MemoryInfo represents memory information
type MemoryInfo struct {
	TotalBytes     int64
	AvailableBytes int64
}

// IsValid checks if the memory info is valid
func (m MemoryInfo) IsValid() bool {
	return m.TotalBytes > 0 && m.AvailableBytes >= 0
}

// ResourceRequirements represents resource requirements for different project types
type ResourceRequirements struct {
	MinDiskSpaceGB         int
	MinMemoryGB            int
	RecommendedDiskSpaceGB int
	RecommendedMemoryGB    int
}

// IsValid checks if the resource requirements are valid
func (r ResourceRequirements) IsValid() bool {
	return r.MinDiskSpaceGB > 0 && r.MinMemoryGB > 0 &&
		r.RecommendedDiskSpaceGB >= r.MinDiskSpaceGB &&
		r.RecommendedMemoryGB >= r.MinMemoryGB
}

// ResourceLimits represents resource limits
type ResourceLimits struct {
	MaxMemoryGB   int
	MaxDiskGB     int
	MaxBuildTime  int // seconds
	MaxConcurrent int
}

// IsValid checks if the resource limits are valid
func (r ResourceLimits) IsValid() bool {
	return r.MaxMemoryGB > 0 && r.MaxDiskGB > 0 && r.MaxBuildTime > 0 && r.MaxConcurrent > 0
}

// ResourceUsage represents current resource usage
type ResourceUsage struct {
	Timestamp   time.Time
	MemoryBytes int64
	CPUPercent  float64
	DiskBytes   int64
}

// ResourceQuota represents resource quotas
type ResourceQuota struct {
	MemoryLimitGB  int
	DiskLimitGB    int
	TimeLimitHours int
}

// IsValid checks if the resource quota is valid
func (r ResourceQuota) IsValid() bool {
	return r.MemoryLimitGB > 0 && r.DiskLimitGB > 0 && r.TimeLimitHours > 0
}

// IsWithinLimits checks if the usage is within the quota limits
func (r ResourceQuota) IsWithinLimits(usage ResourceUsage) bool {
	memoryGB := float64(usage.MemoryBytes) / (1024 * 1024 * 1024)
	diskGB := float64(usage.DiskBytes) / (1024 * 1024 * 1024)

	return memoryGB <= float64(r.MemoryLimitGB) && diskGB <= float64(r.DiskLimitGB)
}

// BuildResourceEstimate represents estimated build resource usage
type BuildResourceEstimate struct {
	EstimatedTimeSeconds int
	EstimatedMemoryGB    int
	EstimatedDiskGB      int
}

// ResourceChecker handles resource checking operations
type ResourceChecker struct{}

// NewResourceChecker creates a new resource checker
func NewResourceChecker() *ResourceChecker {
	return &ResourceChecker{}
}

// CheckDiskSpace checks disk space for a given path
func (r *ResourceChecker) CheckDiskSpace(path string) (DiskSpaceInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return DiskSpaceInfo{}, fmt.Errorf("failed to check disk space for %s: %v", path, err)
	}

	totalBytes := int64(stat.Blocks) * int64(stat.Bsize)
	availableBytes := int64(stat.Bavail) * int64(stat.Bsize)

	return DiskSpaceInfo{
		Path:           path,
		TotalBytes:     totalBytes,
		AvailableBytes: availableBytes,
	}, nil
}

// CheckMemory checks available memory
func (r *ResourceChecker) CheckMemory() (MemoryInfo, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// This is a simplified implementation
	// In a real implementation, you'd use syscall to get system memory info
	totalBytes := int64(8 * 1024 * 1024 * 1024) // Assume 8GB for testing
	availableBytes := totalBytes - int64(memStats.Alloc)

	return MemoryInfo{
		TotalBytes:     totalBytes,
		AvailableBytes: availableBytes,
	}, nil
}

// GetCPUCount returns the number of CPU cores
func (r *ResourceChecker) GetCPUCount() int {
	return runtime.NumCPU()
}

// EstimateRequirements estimates resource requirements for a project type
func (r *ResourceChecker) EstimateRequirements(projectType string) ResourceRequirements {
	switch projectType {
	case "minimal":
		return ResourceRequirements{
			MinDiskSpaceGB:         5,
			MinMemoryGB:            2,
			RecommendedDiskSpaceGB: 10,
			RecommendedMemoryGB:    4,
		}
	case "networking":
		return ResourceRequirements{
			MinDiskSpaceGB:         10,
			MinMemoryGB:            4,
			RecommendedDiskSpaceGB: 20,
			RecommendedMemoryGB:    8,
		}
	case "iot":
		return ResourceRequirements{
			MinDiskSpaceGB:         8,
			MinMemoryGB:            2,
			RecommendedDiskSpaceGB: 15,
			RecommendedMemoryGB:    4,
		}
	case "security":
		return ResourceRequirements{
			MinDiskSpaceGB:         15,
			MinMemoryGB:            4,
			RecommendedDiskSpaceGB: 30,
			RecommendedMemoryGB:    8,
		}
	case "industrial":
		return ResourceRequirements{
			MinDiskSpaceGB:         12,
			MinMemoryGB:            4,
			RecommendedDiskSpaceGB: 25,
			RecommendedMemoryGB:    8,
		}
	case "kiosk":
		return ResourceRequirements{
			MinDiskSpaceGB:         20,
			MinMemoryGB:            4,
			RecommendedDiskSpaceGB: 50,
			RecommendedMemoryGB:    8,
		}
	default:
		return ResourceRequirements{
			MinDiskSpaceGB:         10,
			MinMemoryGB:            4,
			RecommendedDiskSpaceGB: 20,
			RecommendedMemoryGB:    8,
		}
	}
}

// ValidateRequirements validates if the system meets requirements for a project type
func (r *ResourceChecker) ValidateRequirements(projectType string) error {
	reqs := r.EstimateRequirements(projectType)

	// Check disk space
	diskInfo, err := r.CheckDiskSpace("/")
	if err != nil {
		return fmt.Errorf("failed to check disk space: %v", err)
	}

	minDiskBytes := int64(reqs.MinDiskSpaceGB) * 1024 * 1024 * 1024
	if diskInfo.AvailableBytes < minDiskBytes {
		return fmt.Errorf("insufficient disk space: %d GB available, %d GB required",
			diskInfo.AvailableBytes/(1024*1024*1024), reqs.MinDiskSpaceGB)
	}

	// Check memory
	memInfo, err := r.CheckMemory()
	if err != nil {
		return fmt.Errorf("failed to check memory: %v", err)
	}

	minMemBytes := int64(reqs.MinMemoryGB) * 1024 * 1024 * 1024
	if memInfo.AvailableBytes < minMemBytes {
		return fmt.Errorf("insufficient memory: %d GB available, %d GB required",
			memInfo.AvailableBytes/(1024*1024*1024), reqs.MinMemoryGB)
	}

	return nil
}

// GetResourceWarnings returns resource-related warnings
func (r *ResourceChecker) GetResourceWarnings() []string {
	warnings := make([]string, 0)

	// Check disk space
	diskInfo, err := r.CheckDiskSpace("/")
	if err == nil {
		availableGB := diskInfo.AvailableBytes / (1024 * 1024 * 1024)
		if availableGB < 20 {
			warnings = append(warnings, fmt.Sprintf("Low disk space: %d GB available", availableGB))
		}
	}

	// Check memory
	memInfo, err := r.CheckMemory()
	if err == nil {
		availableGB := memInfo.AvailableBytes / (1024 * 1024 * 1024)
		if availableGB < 4 {
			warnings = append(warnings, fmt.Sprintf("Low memory: %d GB available", availableGB))
		}
	}

	return warnings
}

// ResourceMonitor monitors resource usage during operations
type ResourceMonitor struct {
	running bool
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor() *ResourceMonitor {
	return &ResourceMonitor{}
}

// Start starts resource monitoring
func (r *ResourceMonitor) Start() error {
	r.running = true
	return nil
}

// Stop stops resource monitoring
func (r *ResourceMonitor) Stop() {
	r.running = false
}

// GetCurrentUsage returns current resource usage
func (r *ResourceMonitor) GetCurrentUsage() ResourceUsage {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return ResourceUsage{
		Timestamp:   time.Now(),
		MemoryBytes: int64(memStats.Alloc),
		CPUPercent:  0.0, // Simplified - would need more complex CPU monitoring
		DiskBytes:   0,   // Simplified - would need disk I/O monitoring
	}
}

// CheckAlerts checks for resource alerts
func (r *ResourceMonitor) CheckAlerts() []string {
	alerts := make([]string, 0)

	usage := r.GetCurrentUsage()

	// Check memory usage
	memoryGB := float64(usage.MemoryBytes) / (1024 * 1024 * 1024)
	if memoryGB > 6 { // Alert if using more than 6GB
		alerts = append(alerts, fmt.Sprintf("High memory usage: %.1f GB", memoryGB))
	}

	// Check CPU usage (simplified)
	if usage.CPUPercent > 90 {
		alerts = append(alerts, fmt.Sprintf("High CPU usage: %.1f%%", usage.CPUPercent))
	}

	return alerts
}

// BuildEstimator estimates build resource requirements
type BuildEstimator struct{}

// NewBuildEstimator creates a new build estimator
func NewBuildEstimator() *BuildEstimator {
	return &BuildEstimator{}
}

// EstimateBuildResources estimates resources needed for building a project type
func (b *BuildEstimator) EstimateBuildResources(projectType string) BuildResourceEstimate {
	switch projectType {
	case "minimal":
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 600, // 10 minutes
			EstimatedMemoryGB:    2,
			EstimatedDiskGB:      5,
		}
	case "networking":
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 1800, // 30 minutes
			EstimatedMemoryGB:    4,
			EstimatedDiskGB:      10,
		}
	case "iot":
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 1200, // 20 minutes
			EstimatedMemoryGB:    3,
			EstimatedDiskGB:      8,
		}
	case "security":
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 2400, // 40 minutes
			EstimatedMemoryGB:    6,
			EstimatedDiskGB:      15,
		}
	case "industrial":
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 2100, // 35 minutes
			EstimatedMemoryGB:    5,
			EstimatedDiskGB:      12,
		}
	case "kiosk":
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 3600, // 60 minutes
			EstimatedMemoryGB:    8,
			EstimatedDiskGB:      25,
		}
	default:
		return BuildResourceEstimate{
			EstimatedTimeSeconds: 1800,
			EstimatedMemoryGB:    4,
			EstimatedDiskGB:      10,
		}
	}
}

// ResourceTracker tracks resource usage history
type ResourceTracker struct {
	history []ResourceUsage
}

// NewResourceTracker creates a new resource tracker
func NewResourceTracker() *ResourceTracker {
	return &ResourceTracker{
		history: make([]ResourceUsage, 0),
	}
}

// RecordUsage records resource usage
func (r *ResourceTracker) RecordUsage(usage ResourceUsage) {
	r.history = append(r.history, usage)
}

// GetHistory returns the resource usage history
func (r *ResourceTracker) GetHistory() []ResourceUsage {
	return r.history
}
