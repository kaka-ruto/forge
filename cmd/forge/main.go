package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "forge",
		Short:   "Forge OS - Framework for creating custom embedded Linux operating systems",
		Long:    `Forge OS is a framework that follows the "Rails for embedded Linux" philosophy, providing convention over configuration for rapid development of custom embedded Linux systems.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add subcommands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Placeholder commands - will be implemented with TDD
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Forge OS project",
	Long:  `Create a new Forge OS project with the specified name and template.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("forge new command - not yet implemented")
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Forge OS image",
	Long:  `Build the Forge OS image using Buildroot with the current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("forge build command - not yet implemented")
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the Forge OS image in QEMU",
	Long:  `Test the built Forge OS image by running it in QEMU emulator.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("forge test command - not yet implemented")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy [target]",
	Short: "Deploy the Forge OS image",
	Long:  `Deploy the built Forge OS image to various targets (USB, SD card, remote, etc.).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("forge deploy command - not yet implemented")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show detailed version information including Forge OS version, Buildroot version, and build details.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Forge OS %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system and diagnose issues",
	Long:  `Check the development environment, dependencies, and diagnose potential issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("forge doctor command - not yet implemented")
	},
}
