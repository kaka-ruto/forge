package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/packages"
)

// NewPackagesCommand creates the packages command
func NewPackagesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages",
		Short: "Manage Forge OS packages",
		Long:  `Install, uninstall, and manage packages for your Forge OS project.`,
	}

	cmd.AddCommand(
		newPackagesInstallCommand(),
		newPackagesUninstallCommand(),
		newPackagesListCommand(),
		newPackagesInfoCommand(),
	)

	return cmd
}

func newPackagesInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [packages...]",
		Short: "Install packages",
		Long:  `Install one or more packages and their dependencies.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  runPackagesInstallCommandE,
	}

	cmd.Flags().StringP("buildroot", "b", "", "Path to Buildroot directory (auto-detect if not specified)")

	return cmd
}

func newPackagesUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall [packages...]",
		Short: "Uninstall packages",
		Long:  `Uninstall one or more packages.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  runPackagesUninstallCommandE,
	}

	cmd.Flags().StringP("buildroot", "b", "", "Path to Buildroot directory (auto-detect if not specified)")

	return cmd
}

func newPackagesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [category]",
		Short: "List available packages",
		Long:  `List all available packages or packages in a specific category.`,
		RunE:  runPackagesListCommandE,
	}

	return cmd
}

func newPackagesInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [package]",
		Short: "Show package information",
		Long:  `Show detailed information about a specific package.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runPackagesInfoCommandE,
	}

	return cmd
}

func runPackagesInstallCommandE(cmd *cobra.Command, args []string) error {
	// Check if we're in a Forge project directory
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no forge.yml found - not in a Forge project directory")
	}

	// Load configuration
	cfg, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Create package manager
	pm := packages.NewPackageManager(cfg)

	// Get Buildroot directory
	buildrootDir, err := getBuildrootDir(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("Installing packages: %s\n", strings.Join(args, ", "))
	fmt.Printf("Buildroot directory: %s\n\n", buildrootDir)

	// Install packages
	results := pm.InstallPackages(args, buildrootDir)

	// Display results
	successCount := 0
	for _, result := range results {
		if result.Success {
			fmt.Printf("✓ %s installed successfully\n", result.Package)
			if len(result.Services) > 0 {
				fmt.Printf("  Services: %s\n", strings.Join(result.Services, ", "))
			}
			if len(result.ConfigFiles) > 0 {
				fmt.Printf("  Config files: %s\n", strings.Join(result.ConfigFiles, ", "))
			}
			successCount++
		} else {
			fmt.Printf("✗ %s failed: %s\n", result.Package, result.Error)
		}
	}

	fmt.Printf("\nInstalled %d/%d packages successfully\n", successCount, len(results))

	if successCount < len(results) {
		fmt.Println("\nNote: You may need to rebuild your Forge OS image for changes to take effect.")
		fmt.Println("Run 'forge build' to rebuild.")
	}

	return nil
}

func runPackagesUninstallCommandE(cmd *cobra.Command, args []string) error {
	// Check if we're in a Forge project directory
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no forge.yml found - not in a Forge project directory")
	}

	// Load configuration
	cfg, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Create package manager
	pm := packages.NewPackageManager(cfg)

	// Get Buildroot directory
	buildrootDir, err := getBuildrootDir(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("Uninstalling packages: %s\n", strings.Join(args, ", "))
	fmt.Printf("Buildroot directory: %s\n\n", buildrootDir)

	// Uninstall packages
	results := pm.UninstallPackages(args, buildrootDir)

	// Display results
	successCount := 0
	for _, result := range results {
		if result.Success {
			fmt.Printf("✓ %s uninstalled successfully\n", result.Package)
			successCount++
		} else {
			fmt.Printf("✗ %s failed: %s\n", result.Package, result.Error)
		}
	}

	fmt.Printf("\nUninstalled %d/%d packages successfully\n", successCount, len(results))

	if successCount > 0 {
		fmt.Println("\nNote: You may need to rebuild your Forge OS image for changes to take effect.")
		fmt.Println("Run 'forge build' to rebuild.")
	}

	return nil
}

func runPackagesListCommandE(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Create package manager
	pm := packages.NewPackageManager(cfg)

	if len(args) == 0 {
		// List all categories
		categories := pm.GetCategories()
		fmt.Println("Available package categories:")
		for _, category := range categories {
			fmt.Printf("  %s\n", category)
		}
		fmt.Println("\nUse 'forge packages list <category>' to see packages in a category.")
	} else {
		// List packages in category
		category := args[0]
		pkgs := pm.ListPackagesByCategory(category)
		if len(pkgs) == 0 {
			return fmt.Errorf("no packages found in category '%s'", category)
		}

		fmt.Printf("Packages in category '%s':\n", category)
		for _, pkg := range pkgs {
			fmt.Printf("  %-15s %s\n", pkg.Name, pkg.Description)
		}
	}

	return nil
}

func runPackagesInfoCommandE(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Create package manager
	pm := packages.NewPackageManager(cfg)

	pkgName := args[0]
	pkg, err := pm.GetPackageInfo(pkgName)
	if err != nil {
		return err
	}

	fmt.Printf("Package: %s\n", pkg.Name)
	fmt.Printf("Version: %s\n", pkg.Version)
	fmt.Printf("Category: %s\n", pkg.Category)
	fmt.Printf("Description: %s\n", pkg.Description)

	if len(pkg.Dependencies) > 0 {
		fmt.Printf("Dependencies: %s\n", strings.Join(pkg.Dependencies, ", "))
	} else {
		fmt.Println("Dependencies: none")
	}

	if len(pkg.Conflicts) > 0 {
		fmt.Printf("Conflicts: %s\n", strings.Join(pkg.Conflicts, ", "))
	}

	fmt.Printf("Buildroot package: %s\n", pkg.BuildrootPkg)

	return nil
}

func getBuildrootDir(cmd *cobra.Command) (string, error) {
	// Check if Buildroot directory was specified
	buildrootDir, _ := cmd.Flags().GetString("buildroot")
	if buildrootDir != "" {
		if _, err := os.Stat(buildrootDir); os.IsNotExist(err) {
			return "", fmt.Errorf("Buildroot directory does not exist: %s", buildrootDir)
		}
		return buildrootDir, nil
	}

	// Auto-detect Buildroot directory
	// Look for build directory in current project
	buildDir := "build"
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return "", fmt.Errorf("no build directory found - run 'forge build' first")
	}

	// Look for Buildroot directory in build
	buildrootPath := filepath.Join(buildDir, "buildroot")
	if _, err := os.Stat(buildrootPath); err == nil {
		return buildrootPath, nil
	}

	// Look for Buildroot in the build directory structure
	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return "", fmt.Errorf("failed to read build directory: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "buildroot-") && entry.IsDir() {
			return filepath.Join(buildDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("Buildroot directory not found - specify with --buildroot flag")
}
