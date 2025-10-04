package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/cli"
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
	rootCmd.AddCommand(cli.NewNewCommand())
	rootCmd.AddCommand(cli.NewBuildCommand())
	rootCmd.AddCommand(cli.NewTestCommand())
	rootCmd.AddCommand(cli.NewDeployCommand())
	rootCmd.AddCommand(cli.NewPackagesCommand())
	rootCmd.AddCommand(cli.NewVersionCommand())
	rootCmd.AddCommand(cli.NewDoctorCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
