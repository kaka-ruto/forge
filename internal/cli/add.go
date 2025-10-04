package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/config"
)

// NewAddCommand creates the add command
func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add packages or features to a Forge OS project",
		Long:  `Add packages or features to your Forge OS project configuration.`,
	}

	cmd.AddCommand(
		newAddPackageCommand(),
		newAddFeatureCommand(),
	)

	return cmd
}

func newAddPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package [package-name]",
		Short: "Add a package to the project",
		Long: `Add a package to the forge.yml configuration.
The package will be included in the next build.`,
		Args: cobra.ExactArgs(1),
		RunE: runAddPackageCommandE,
	}

	return cmd
}

func newAddFeatureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature [feature-name]",
		Short: "Add a feature to the project",
		Long: `Add a feature to the forge.yml configuration.
Features provide pre-configured functionality and may include additional packages.`,
		Args: cobra.ExactArgs(1),
		RunE: runAddFeatureCommandE,
	}

	return cmd
}

func runAddPackageCommandE(cmd *cobra.Command, args []string) error {
	packageName := args[0]
	return runAddPackageCommand([]string{packageName}, map[string]interface{}{})
}

func runAddFeatureCommandE(cmd *cobra.Command, args []string) error {
	featureName := args[0]
	return runAddFeatureCommand([]string{featureName}, map[string]interface{}{})
}

func runAddPackageCommand(args []string, flags map[string]interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("package name is required")
	}

	packageName := args[0]

	// Load current config
	cfg, err := config.LoadConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("failed to load forge.yml: %v", err)
	}

	// Check if package already exists
	for _, pkg := range cfg.Packages {
		if pkg == packageName {
			return fmt.Errorf("package '%s' is already added to the project", packageName)
		}
	}

	// Add package
	cfg.Packages = append(cfg.Packages, packageName)

	// Save config
	if err := config.SaveConfig(cfg, "forge.yml"); err != nil {
		return fmt.Errorf("failed to save forge.yml: %v", err)
	}

	fmt.Printf("Added package '%s' to forge.yml\n", packageName)
	return nil
}

func runAddFeatureCommand(args []string, flags map[string]interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("feature name is required")
	}

	featureName := args[0]

	// Load current config
	cfg, err := config.LoadConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("failed to load forge.yml: %v", err)
	}

	// Check if feature already exists
	for _, feature := range cfg.Features {
		if feature == featureName {
			return fmt.Errorf("feature '%s' is already added to the project", featureName)
		}
	}

	// Add feature
	cfg.Features = append(cfg.Features, featureName)

	// Save config
	if err := config.SaveConfig(cfg, "forge.yml"); err != nil {
		return fmt.Errorf("failed to save forge.yml: %v", err)
	}

	fmt.Printf("Added feature '%s' to forge.yml\n", featureName)
	return nil
}
