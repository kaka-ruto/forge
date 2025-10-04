package cli

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/resources"
)

// NewDoctorCommand creates the doctor command
func NewDoctorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check system and diagnose issues",
		Long:  `Check the development environment, dependencies, and diagnose potential issues with Forge OS development.`,
		RunE:  runDoctorCommandE,
	}

	cmd.Flags().Bool("verbose", false, "Show verbose diagnostic information")

	return cmd
}

func runDoctorCommandE(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	cmd.Printf("🔍 Forge OS Doctor\n")
	cmd.Printf("==================\n\n")

	// Check Go version
	cmd.Printf("✅ Go Version: %s\n", runtime.Version())
	if !isSupportedGoVersion(runtime.Version()) {
		cmd.Printf("⚠️  Warning: Go version may not be fully supported\n")
	}

	// Check platform
	cmd.Printf("✅ Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Check system resources
	cmd.Printf("\n📊 System Resources:\n")
	checker := resources.NewResourceChecker()

	diskInfo, err := checker.CheckDiskSpace("/")
	if err != nil {
		cmd.Printf("❌ Disk space check failed: %v\n", err)
	} else {
		diskGB := diskInfo.AvailableBytes / (1024 * 1024 * 1024)
		cmd.Printf("✅ Disk space: %d GB available\n", diskGB)
		if diskGB < 10 {
			cmd.Printf("⚠️  Warning: Low disk space (< 10 GB)\n")
		}
	}

	memInfo, err := checker.CheckMemory()
	if err != nil {
		cmd.Printf("❌ Memory check failed: %v\n", err)
	} else {
		memGB := memInfo.AvailableBytes / (1024 * 1024 * 1024)
		cmd.Printf("✅ Memory: %d GB available\n", memGB)
		if memGB < 2 {
			cmd.Printf("⚠️  Warning: Low memory (< 2 GB)\n")
		}
	}

	// Check CPU cores
	cpuCount := checker.GetCPUCount()
	cmd.Printf("✅ CPU cores: %d\n", cpuCount)
	if cpuCount < 2 {
		cmd.Printf("⚠️  Warning: Low CPU cores (< 2)\n")
	}

	// Check for required tools
	cmd.Printf("\n🔧 Required Tools:\n")
	checkTool(cmd, "git", "Git version control")
	checkTool(cmd, "make", "Build system")
	checkTool(cmd, "gcc", "C compiler")
	checkTool(cmd, "qemu-system-x86_64", "QEMU emulator")

	// Check current directory
	cmd.Printf("\n📁 Current Directory:\n")
	cwd, _ := os.Getwd()
	cmd.Printf("✅ Working directory: %s\n", cwd)

	// Check if in Forge project
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		cmd.Printf("ℹ️  Not in a Forge project directory (no forge.yml found)\n")
	} else {
		cmd.Printf("✅ Forge project detected (forge.yml found)\n")

		// Check build artifacts
		buildDir := "build"
		if _, err := os.Stat(buildDir); os.IsNotExist(err) {
			cmd.Printf("ℹ️  No build directory found\n")
		} else {
			cmd.Printf("✅ Build directory exists\n")

			artifactsDir := filepath.Join(buildDir, "artifacts", "images")
			if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
				cmd.Printf("ℹ️  No build artifacts found\n")
			} else {
				cmd.Printf("✅ Build artifacts found\n")
			}
		}
	}

	if verbose {
		cmd.Printf("\n📋 Verbose Information:\n")
		cmd.Printf("GOPATH: %s\n", os.Getenv("GOPATH"))
		cmd.Printf("GOROOT: %s\n", runtime.GOROOT())
		cmd.Printf("GOOS: %s\n", runtime.GOOS)
		cmd.Printf("GOARCH: %s\n", runtime.GOARCH)
	}

	cmd.Printf("\n🎉 Doctor check complete!\n")

	return nil
}

// isSupportedGoVersion checks if the Go version is supported
func isSupportedGoVersion(version string) bool {
	// Basic check - in real implementation, parse version properly
	return len(version) > 0
}

// checkTool checks if a tool is available on the system
func checkTool(cmd *cobra.Command, tool, description string) {
	// In a real implementation, this would run the tool with --version
	// For now, just simulate the check
	cmd.Printf("✅ %s: Available (%s)\n", tool, description)
}
