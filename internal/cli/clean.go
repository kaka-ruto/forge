package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewCleanCommand creates the clean command
func NewCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean build artifacts and cache",
		Long:  `Remove build artifacts, cache files, and temporary data to free up disk space.`,
		RunE:  runCleanCommandE,
	}

	cmd.Flags().Bool("all", false, "Remove all artifacts including cache and logs")
	cmd.Flags().Bool("cache", false, "Remove download cache")
	cmd.Flags().Bool("builds", false, "Remove build artifacts")
	cmd.Flags().Bool("logs", false, "Remove log files")
	cmd.Flags().Bool("dry-run", false, "Show what would be deleted without actually deleting")

	return cmd
}

func runCleanCommandE(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	cache, _ := cmd.Flags().GetBool("cache")
	builds, _ := cmd.Flags().GetBool("builds")
	logs, _ := cmd.Flags().GetBool("logs")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	return runCleanCommand(args, map[string]interface{}{
		"all":     all,
		"cache":   cache,
		"builds":  builds,
		"logs":    logs,
		"dry-run": dryRun,
	})
}

func runCleanCommand(args []string, flags map[string]interface{}) error {
	var pathsToClean []string

	// Determine what to clean
	if flags["all"].(bool) {
		pathsToClean = []string{
			"build/",
			"output/",
			"dl/",
			".ccache/",
			".cache/",
			"forge.log",
			".forge/logs/",
			".forge/cache/",
			".forge/metrics/",
		}
	} else {
		if flags["cache"].(bool) {
			pathsToClean = append(pathsToClean, "dl/", ".ccache/", ".cache/", ".forge/cache/")
		}
		if flags["builds"].(bool) {
			pathsToClean = append(pathsToClean, "build/", "output/")
		}
		if flags["logs"].(bool) {
			pathsToClean = append(pathsToClean, "forge.log", ".forge/logs/")
		}
	}

	if len(pathsToClean) == 0 {
		fmt.Println("Nothing to clean. Use --all, --cache, --builds, or --logs flags.")
		return nil
	}

	// Collect files to delete
	var filesToDelete []string
	var totalSize int64

	for _, path := range pathsToClean {
		size, files := collectFilesToClean(path, flags["dry-run"].(bool))
		filesToDelete = append(filesToDelete, files...)
		totalSize += size
	}

	if len(filesToDelete) == 0 {
		fmt.Println("No files to clean.")
		return nil
	}

	// Show what will be deleted
	fmt.Printf("Found %d files to clean (%s)\n", len(filesToDelete), formatSize(totalSize))

	if flags["dry-run"].(bool) {
		fmt.Println("\nDry run - would delete:")
	} else {
		fmt.Println("\nDeleting:")
	}

	for _, file := range filesToDelete {
		fmt.Printf("  %s\n", file)
	}

	// Perform deletion
	if !flags["dry-run"].(bool) {
		for _, file := range filesToDelete {
			if err := os.RemoveAll(file); err != nil {
				fmt.Printf("Warning: failed to delete %s: %v\n", file, err)
			}
		}
		fmt.Printf("\nCleaned %d files, freed %s\n", len(filesToDelete), formatSize(totalSize))
	}

	return nil
}

func collectFilesToClean(path string, dryRun bool) (int64, []string) {
	var files []string
	var totalSize int64

	// Check if path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return 0, files
	}

	if info.IsDir() {
		// Walk directory
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if !info.IsDir() {
				files = append(files, filePath)
				totalSize += info.Size()
			}
			return nil
		})
		// Include the directory itself
		files = append(files, path)
	} else {
		// Single file
		files = append(files, path)
		totalSize += info.Size()
	}

	return totalSize, files
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
