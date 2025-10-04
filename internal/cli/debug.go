package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// NewDebugCommand creates the debug command
func NewDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug Forge OS projects and environment",
		Long:  `Gather diagnostic information and debug project configurations.`,
		RunE:  runDebugCommandE,
	}

	cmd.Flags().Bool("config", false, "Validate and show configuration details")
	cmd.Flags().Bool("env", false, "Show environment information")
	cmd.Flags().Bool("system", false, "Show system information")

	return cmd
}

func runDebugCommandE(cmd *cobra.Command, args []string) error {
	config, _ := cmd.Flags().GetBool("config")
	env, _ := cmd.Flags().GetBool("env")
	system, _ := cmd.Flags().GetBool("system")

	// If no specific flags, show general debug info
	if !config && !env && !system {
		config, env, system = true, true, true
	}

	return runDebugCommand(args, map[string]interface{}{
		"config": config,
		"env":    env,
		"system": system,
	})
}

func runDebugCommand(args []string, flags map[string]interface{}) error {
	fmt.Println("Forge OS Debug Information")
	fmt.Println("==========================")

	if flags["system"].(bool) {
		showSystemInfo()
		fmt.Println()
	}

	if flags["env"].(bool) {
		showEnvironmentInfo()
		fmt.Println()
	}

	if flags["config"].(bool) {
		showConfigInfo()
		fmt.Println()
	}

	return nil
}

func showSystemInfo() {
	fmt.Println("System Information:")
	fmt.Printf("  OS: %s\n", runtime.GOOS)
	fmt.Printf("  Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("  CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("  Go Version: %s\n", runtime.Version())

	// Memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  Memory: %d KB used\n", m.Alloc/1024)
}

func showEnvironmentInfo() {
	fmt.Println("Environment Information:")
	fmt.Printf("  Current Directory: %s\n", getCurrentDir())

	// Check for required tools
	tools := []string{"git", "make", "gcc", "wget"}
	for _, tool := range tools {
		if path, err := exec.LookPath(tool); err == nil {
			fmt.Printf("  %s: %s\n", tool, path)
		} else {
			fmt.Printf("  %s: not found\n", tool)
		}
	}

	// Environment variables
	envVars := []string{"PATH", "GOROOT", "GOPATH", "BR2_DL_DIR"}
	for _, env := range envVars {
		if value := os.Getenv(env); value != "" {
			fmt.Printf("  %s: %s\n", env, value)
		} else {
			fmt.Printf("  %s: not set\n", env)
		}
	}
}

func showConfigInfo() {
	fmt.Println("Configuration Information:")

	// Check if forge.yml exists
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		fmt.Println("  No forge.yml found in current directory")
		return
	}

	fmt.Println("  forge.yml found")

	// Try to load and validate config
	// This would integrate with the config package
	fmt.Println("  Configuration validation: TODO")
}

func getCurrentDir() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return "unknown"
}
