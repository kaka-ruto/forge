package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/build"
	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/resources"
)

// NewBuildCommand creates the build command
func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build the Forge OS image",
		Long: `Build the Forge OS image using Buildroot with the current configuration.

This command will:
- Validate the forge.yml configuration
- Check system resources
- Download and configure Buildroot
- Build the complete OS image
- Generate build artifacts and reports`,
		RunE: runBuildCommandE,
	}

	cmd.Flags().BoolP("clean", "c", false, "Perform a clean build (remove previous build artifacts)")
	cmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	cmd.Flags().BoolP("incremental", "i", true, "Enable incremental builds (use cache when possible)")
	cmd.Flags().IntP("jobs", "j", 0, "Number of parallel build jobs (0 = auto-detect)")
	cmd.Flags().String("optimize-for", "", "Optimize build for specific use case (size, performance, realtime)")
	cmd.Flags().String("timeout", "2h", "Build timeout duration")

	return cmd
}

// runBuildCommandE is the cobra command handler for the build command
func runBuildCommandE(cmd *cobra.Command, args []string) error {
	clean, _ := cmd.Flags().GetBool("clean")
	verbose, _ := cmd.Flags().GetBool("verbose")
	incremental, _ := cmd.Flags().GetBool("incremental")
	jobs, _ := cmd.Flags().GetInt("jobs")
	optimizeFor, _ := cmd.Flags().GetString("optimize-for")
	timeout, _ := cmd.Flags().GetString("timeout")

	return runBuildCommand(args, map[string]string{
		"clean":        fmt.Sprintf("%t", clean),
		"verbose":      fmt.Sprintf("%t", verbose),
		"incremental":  fmt.Sprintf("%t", incremental),
		"jobs":         fmt.Sprintf("%d", jobs),
		"optimize-for": optimizeFor,
		"timeout":      timeout,
	})
}

// runBuildCommand executes the build logic
func runBuildCommand(args []string, flags map[string]string) error {
	// Check if we're in a Forge project directory
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no forge.yml found - not in a Forge project directory")
	}

	// Load and validate configuration
	config, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Check system resources
	if err := checkBuildResources(config); err != nil {
		return fmt.Errorf("resource check failed: %v", err)
	}

	// Get project directory
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Create build orchestrator
	bo := build.NewBuildOrchestrator(config, projectDir)

	// Parse build options
	opts := build.BuildOptions{
		Clean:       flags["clean"] == "true",
		Verbose:     flags["verbose"] == "true",
		Incremental: flags["incremental"] == "true",
	}

	// Parse jobs
	if jobs := flags["jobs"]; jobs != "" && jobs != "0" {
		// TODO: Parse jobs value
		opts.Jobs = 4 // Default for now
	}

	// Parse optimization
	if optimizeFor := flags["optimize-for"]; optimizeFor != "" {
		opts.OptimizeFor = optimizeFor
	}

	// Parse timeout
	if timeoutStr := flags["timeout"]; timeoutStr != "" && timeoutStr != "2h" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			opts.Timeout = timeout
		}
	}

	// Execute build with orchestrator
	ctx := context.Background()
	if err := bo.Build(ctx, opts); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	return nil
}

// loadForgeConfig loads and validates the forge.yml configuration
func loadForgeConfig(configPath string) (*config.Config, error) {
	return config.LoadConfig(configPath)
}

// checkBuildResources checks if the system has sufficient resources for building
func checkBuildResources(config *config.Config) error {
	checker := resources.NewResourceChecker()
	return checker.ValidateRequirements(config.Template)
}

// generateBuildReport generates a build report
func generateBuildReport(outputDir string) error {
	reportPath := filepath.Join(outputDir, "build-report.txt")
	report := "Forge OS Build Report\n"
	report += "===================\n"
	report += "Build completed successfully\n"
	report += "\n"
	report += fmt.Sprintf("Output directory: %s\n", outputDir)
	report += fmt.Sprintf("Images directory: %s\n", filepath.Join(outputDir, "images"))

	return os.WriteFile(reportPath, []byte(report), 0644)
}
