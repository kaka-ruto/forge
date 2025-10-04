package cli

import (
	"runtime"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Show detailed version information including Forge OS version, Buildroot version, and build details.`,
		RunE:  runVersionCommandE,
	}

	cmd.Flags().Bool("verbose", false, "Show verbose version information")

	return cmd
}

func runVersionCommandE(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	cmd.Printf("Forge OS %s\n", "dev")
	cmd.Printf("Go Version: %s\n", runtime.Version())
	cmd.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	if verbose {
		cmd.Printf("Buildroot Version: %s\n", "stable")
		cmd.Printf("Kernel Version: %s\n", "latest")
	}

	return nil
}
