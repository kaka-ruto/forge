package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/qemu"
	"github.com/sst/forge/internal/resources"
)

// NewTestCommand creates the test command
func NewTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test the Forge OS image in QEMU",
		Long:  `Test the built Forge OS image by running it in QEMU emulator.`,
		RunE:  runTestCommandE,
	}

	cmd.Flags().BoolP("headless", "H", false, "Run in headless mode (exit after tests)")
	cmd.Flags().StringP("image", "i", "", "Path to image file to test (auto-detect if not specified)")
	cmd.Flags().StringSlice("scenarios", []string{}, "Test scenarios to run (boot, network, services)")
	cmd.Flags().Duration("timeout", 5*time.Minute, "Test timeout per scenario")
	cmd.Flags().Int("instances", 1, "Number of instances to run")

	return cmd
}

func runTestCommandE(cmd *cobra.Command, args []string) error {
	headless, _ := cmd.Flags().GetBool("headless")
	image, _ := cmd.Flags().GetString("image")
	scenarios, _ := cmd.Flags().GetStringSlice("scenarios")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	instances, _ := cmd.Flags().GetInt("instances")

	return runTestCommand(args, map[string]interface{}{
		"headless":  headless,
		"image":     image,
		"scenarios": scenarios,
		"timeout":   timeout,
		"instances": instances,
	})
}

// runTestCommand executes the test logic
func runTestCommand(args []string, flags map[string]interface{}) error {
	// Check if we're in a Forge project directory
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no forge.yml found - not in a Forge project directory")
	}

	// Load and validate configuration
	config, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Check for build artifacts
	buildDir := "build"
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		return fmt.Errorf("no build artifacts found - run 'forge build' first")
	}

	// Validate image path if specified
	if image := flags["image"].(string); image != "" {
		if _, err := os.Stat(image); os.IsNotExist(err) {
			return fmt.Errorf("specified image file does not exist: %s", image)
		}
	}

	// Check system resources for QEMU
	if err := checkTestResources(config); err != nil {
		return fmt.Errorf("resource check failed: %v", err)
	}

	// Launch QEMU instances
	if err := launchQEMUInstances(config, artifactsDir, flags); err != nil {
		return fmt.Errorf("failed to launch QEMU: %v", err)
	}

	return nil
}

// isValidPortFormat validates port forwarding format (host:guest)
func isValidPortFormat(port string) bool {
	parts := strings.Split(port, ":")
	if len(parts) != 2 {
		return false
	}

	hostPort, err1 := strconv.Atoi(parts[0])
	guestPort, err2 := strconv.Atoi(parts[1])

	return err1 == nil && err2 == nil && hostPort > 0 && hostPort <= 65535 && guestPort > 0 && guestPort <= 65535
}

// checkTestResources checks if the system has sufficient resources for testing
func checkTestResources(config *config.Config) error {
	// Use resource checker - testing requires more resources than building
	checker := resources.NewResourceChecker()

	// For testing, we need at least double the build requirements
	reqs := checker.EstimateRequirements(config.Template)
	testReqs := resources.ResourceRequirements{
		MinDiskSpaceGB:         reqs.MinDiskSpaceGB * 2, // Extra space for test artifacts
		MinMemoryGB:            reqs.MinMemoryGB * 2,    // Extra memory for QEMU VMs
		RecommendedDiskSpaceGB: reqs.RecommendedDiskSpaceGB * 2,
		RecommendedMemoryGB:    reqs.RecommendedMemoryGB * 2,
	}

	// Check disk space
	diskInfo, err := checker.CheckDiskSpace("/")
	if err != nil {
		return fmt.Errorf("failed to check disk space: %v", err)
	}

	minDiskBytes := int64(testReqs.MinDiskSpaceGB) * 1024 * 1024 * 1024
	if diskInfo.AvailableBytes < minDiskBytes {
		return fmt.Errorf("insufficient disk space for testing: %d GB available, %d GB required",
			diskInfo.AvailableBytes/(1024*1024*1024), testReqs.MinDiskSpaceGB)
	}

	// Check memory
	memInfo, err := checker.CheckMemory()
	if err != nil {
		return fmt.Errorf("failed to check memory: %v", err)
	}

	minMemBytes := int64(testReqs.MinMemoryGB) * 1024 * 1024 * 1024
	if memInfo.AvailableBytes < minMemBytes {
		return fmt.Errorf("insufficient memory for testing: %d GB available, %d GB required",
			memInfo.AvailableBytes/(1024*1024*1024), testReqs.MinMemoryGB)
	}

	return nil
}

// launchQEMUInstances launches the specified number of QEMU instances
func launchQEMUInstances(config *config.Config, artifactsDir string, flags map[string]interface{}) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get project directory: %v", err)
	}

	// Create QEMU manager
	qm := qemu.NewQEMUManager(config, projectDir)

	// Determine image path (prefer rootfs.ext4, fallback to any .img file)
	imagePath := filepath.Join(artifactsDir, "rootfs.ext4")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// Look for any .img file
		entries, err := os.ReadDir(artifactsDir)
		if err != nil {
			return fmt.Errorf("failed to read artifacts directory: %v", err)
		}
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".img") {
				imagePath = filepath.Join(artifactsDir, entry.Name())
				break
			}
		}
	}

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("no suitable image found in %s", artifactsDir)
	}

	instances := flags["instances"].(int)
	if instances <= 0 {
		instances = 1
	}

	var runningInstances []*qemu.QEMUInstance

	// Launch instances
	for i := 0; i < instances; i++ {
		fmt.Printf("Launching QEMU instance %d...\n", i+1)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		instance, err := qm.StartInstance(ctx, imagePath)
		cancel()

		if err != nil {
			// Stop any already running instances
			for _, inst := range runningInstances {
				qm.StopInstance(inst)
			}
			return fmt.Errorf("failed to start QEMU instance %d: %v", i+1, err)
		}

		runningInstances = append(runningInstances, instance)
		fmt.Printf("QEMU instance %s started (SSH: localhost:%d)\n", instance.ID, instance.SSHPort)
	}

	// Run test scenarios if not in headless mode
	if !flags["headless"].(bool) && len(runningInstances) > 0 {
		fmt.Println("\nRunning test scenarios...")

		testResults := []qemu.TestResult{}
		scenarios := qm.GetDefaultTestScenarios()

		for _, instance := range runningInstances {
			fmt.Printf("\nTesting instance %s:\n", instance.ID)

			var instanceResults []*qemu.TestResult
			for _, scenario := range scenarios {
				ctx, cancel := context.WithTimeout(context.Background(), scenario.Timeout)
				result := qm.RunTestScenario(ctx, instance, scenario)
				cancel()

				testResults = append(testResults, *result)
				instanceResults = append(instanceResults, result)
			}

			// Save test results
			if err := qm.SaveTestResults(instanceResults, instance.ID); err != nil {
				fmt.Printf("Warning: failed to save test results: %v\n", err)
			}

			// Load previous results and compare
			if previousResults, err := qm.LoadTestResults(instance.ID); err == nil {
				comparison := qm.CompareTestResults(instanceResults, previousResults)
				fmt.Printf("\nComparison with previous run:\n")
				fmt.Printf("  Total tests: %d\n", comparison.TotalTests)
				fmt.Printf("  Passed: %d, Failed: %d\n", comparison.PassedTests, comparison.FailedTests)
				if comparison.ImprovedTests > 0 {
					fmt.Printf("  Improved: %d\n", comparison.ImprovedTests)
				}
				if comparison.RegressedTests > 0 {
					fmt.Printf("  Regressed: %d\n", comparison.RegressedTests)
				}
				if comparison.NewTests > 0 {
					fmt.Printf("  New tests: %d\n", comparison.NewTests)
				}

				// Show detailed comparison for each test
				for testName, detail := range comparison.Details {
					if detail.Status != "stable_pass" && detail.Status != "stable_fail" {
						fmt.Printf("  %s: %s", testName, detail.Status)
						if detail.DurationStatus != "" && detail.DurationStatus != "same" {
							fmt.Printf(" (%s %v)", detail.DurationStatus, detail.DurationChange.Round(time.Millisecond))
						}
						fmt.Printf("\n")
					}
				}
			}
		}

		// Print summary
		fmt.Println("\nTest Results Summary:")
		fmt.Println("====================")

		passed := 0
		total := len(testResults)

		for _, result := range testResults {
			status := "PASS"
			if !result.Success {
				status = "FAIL"
			}
			fmt.Printf("%s: %s (%v)\n", status, result.TestName, result.Duration.Round(time.Millisecond))

			// Display metrics if available
			if result.Metrics != nil {
				fmt.Printf("    CPU: %.1f%%, Mem: %.1fMB (%.1f%%), Disk: %.1fMB (%.1f%%)\n",
					result.Metrics.CPUUsagePercent,
					result.Metrics.MemoryUsageMB,
					result.Metrics.MemoryUsagePercent,
					result.Metrics.DiskUsageMB,
					result.Metrics.DiskUsagePercent)
				fmt.Printf("    Load: %.2f, %.2f, %.2f\n",
					result.Metrics.LoadAverage1,
					result.Metrics.LoadAverage5,
					result.Metrics.LoadAverage15)
			}

			if result.Success {
				passed++
			}
		}

		fmt.Printf("\nOverall: %d/%d tests passed\n", passed, total)

		if passed != total {
			fmt.Println("Some tests failed - check logs for details")
		}
	}

	// Keep instances running if not headless
	if !flags["headless"].(bool) {
		fmt.Println("\nQEMU instances are running. Press Ctrl+C to stop.")
		fmt.Println("SSH access: ssh root@localhost:<port>")

		// Wait for interrupt
		select {}
	} else {
		// Stop instances in headless mode
		for _, instance := range runningInstances {
			qm.StopInstance(instance)
		}
	}

	return nil
}
