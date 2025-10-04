package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ResourceCleaner handles resource cleanup operations
type ResourceCleaner struct{}

// NewResourceCleaner creates a new resource cleaner
func NewResourceCleaner() *ResourceCleaner {
	return &ResourceCleaner{}
}

// CleanupDirectory cleans up files in a directory older than the specified age
func (r *ResourceCleaner) CleanupDirectory(dir string, maxAge time.Duration) (int64, error) {
	var totalCleaned int64

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filePath := fmt.Sprintf("%s/%s", dir, entry.Name())
			size := info.Size()

			if err := os.Remove(filePath); err != nil {
				continue // Skip files that can't be removed
			}

			totalCleaned += size
		}
	}

	return totalCleaned, nil
}

// CleanupDirectoryDryRun simulates cleanup without actually removing files
func (r *ResourceCleaner) CleanupDirectoryDryRun(dir string, maxAge time.Duration) (int64, error) {
	var totalWouldClean int64

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			totalWouldClean += info.Size()
		}
	}

	return totalWouldClean, nil
}

// CleanupBuildArtifacts cleans up build artifacts
func (r *ResourceCleaner) CleanupBuildArtifacts(buildDir string) error {
	// Remove common build artifacts
	patterns := []string{
		"*.img",
		"*.qcow2",
		"*.raw",
		"*.iso",
		"*.log",
		"*.tmp",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(buildDir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				// Log error but continue
				continue
			}
		}
	}

	return nil
}

// CleanupCache cleans up cache directories
func (r *ResourceCleaner) CleanupCache(cacheDir string, maxSizeMB int) error {
	info, err := os.Stat(cacheDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", cacheDir)
	}

	// Get all files with their info
	type fileInfo struct {
		path    string
		size    int64
		modTime time.Time
	}

	var files []fileInfo
	err = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, fileInfo{
				path:    path,
				size:    info.Size(),
				modTime: info.ModTime(),
			})
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Calculate total size
	var totalSize int64
	for _, file := range files {
		totalSize += file.size
	}

	maxSizeBytes := int64(maxSizeMB) * 1024 * 1024
	if totalSize <= maxSizeBytes {
		return nil // Cache is within limits
	}

	// Sort files by modification time (oldest first)
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i].modTime.After(files[j].modTime) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	// Remove oldest files until under limit
	for _, file := range files {
		if totalSize <= maxSizeBytes {
			break
		}

		if err := os.Remove(file.path); err != nil {
			continue // Skip files that can't be removed
		}

		totalSize -= file.size
	}

	return nil
}

// CleanupForgeDirectories cleans up Forge-specific directories
func (r *ResourceCleaner) CleanupForgeDirectories(forgeDir string) error {
	// Clean up .forge directories
	patterns := []string{
		".forge/metrics/*.json",
		".forge/logs/*.log",
		".forge/cache/*",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(forgeDir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			// Remove files older than 30 days
			if time.Since(info.ModTime()) > 30*24*time.Hour {
				os.Remove(match)
			}
		}
	}

	return nil
}
