package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/templates"
)

// NewListCommand creates the list command
func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates and packages",
		Long:  `List available project templates and packages for Forge OS projects.`,
	}

	cmd.AddCommand(
		newListTemplatesCommand(),
		newListPackagesCommand(),
	)

	return cmd
}

func newListTemplatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List available project templates",
		Long:  `List all available project templates with their descriptions.`,
		RunE:  runListTemplatesCommandE,
	}

	return cmd
}

func newListPackagesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages [category]",
		Short: "List available packages",
		Long: `List all available packages or packages in a specific category.
If no category is specified, lists all categories.`,
		RunE: runListPackagesCommandE,
	}

	return cmd
}

func runListTemplatesCommandE(cmd *cobra.Command, args []string) error {
	return runListTemplatesCommand(args, map[string]interface{}{})
}

func runListPackagesCommandE(cmd *cobra.Command, args []string) error {
	return runListPackagesCommand(args, map[string]interface{}{})
}

func runListTemplatesCommand(args []string, flags map[string]interface{}) error {
	tm := templates.NewTemplateManager()
	templates := tm.ListTemplates()

	// Sort templates by name for consistent output
	var names []string
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Available templates:")
	fmt.Println()

	for _, name := range names {
		template := templates[name]
		fmt.Printf("  %s\n", name)
		fmt.Printf("    %s\n", template.Description)
		fmt.Println()
	}

	return nil
}

func runListPackagesCommand(args []string, flags map[string]interface{}) error {
	// For now, just show a placeholder message
	// This would integrate with the packages system
	fmt.Println("Available package categories:")
	fmt.Println()
	fmt.Println("  networking    - Network-related packages")
	fmt.Println("  security      - Security and encryption packages")
	fmt.Println("  iot          - Internet of Things packages")
	fmt.Println("  multimedia   - Audio/video packages")
	fmt.Println("  development  - Development tools")
	fmt.Println("  utilities    - System utilities")
	fmt.Println()
	fmt.Println("Use 'forge packages list <category>' to see packages in a category.")
	fmt.Println("Use 'forge packages install <package>' to install a package.")

	return nil
}
