package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewNewCommand creates the new command
func NewNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new Forge OS project",
		Long: `Create a new Forge OS project with the specified name and template.

The project will be created with a forge.yml configuration file and all
necessary directory structure based on the chosen template.`,
		Args: cobra.ExactArgs(1),
		RunE: runNewCommandE,
	}

	cmd.Flags().StringP("template", "t", "minimal", "Project template (minimal, networking, iot, security, industrial, kiosk)")
	cmd.Flags().StringP("arch", "a", "x86_64", "Target architecture (x86_64, arm, aarch64, mips)")
	cmd.Flags().Bool("git", true, "Initialize git repository")

	return cmd
}

// runNewCommandE is the cobra command handler for the new command
func runNewCommandE(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	template, _ := cmd.Flags().GetString("template")
	arch, _ := cmd.Flags().GetString("arch")
	initGit, _ := cmd.Flags().GetBool("git")

	return runNewCommand([]string{projectName}, map[string]string{
		"template": template,
		"arch":     arch,
		"git":      fmt.Sprintf("%t", initGit),
	})
}

// runNewCommand executes the new project creation logic
func runNewCommand(args []string, flags map[string]string) error {
	if len(args) != 1 {
		return fmt.Errorf("project name is required")
	}

	projectName := args[0]
	template := flags["template"]
	arch := flags["arch"]

	// Create project directory
	projectDir := filepath.Join(".", projectName)

	return createProjectStructure(projectDir, template, arch)
}

// createProjectStructure creates the project directory structure and files
func createProjectStructure(projectDir, template, arch string) error {
	// Validate template
	validTemplates := map[string]bool{
		"minimal":    true,
		"networking": true,
		"iot":        true,
		"security":   true,
		"industrial": true,
		"kiosk":      true,
	}

	if !validTemplates[template] {
		return fmt.Errorf("invalid template: %s", template)
	}

	// Validate architecture
	validArches := map[string]bool{
		"x86_64":  true,
		"arm":     true,
		"aarch64": true,
		"mips":    true,
	}

	if !validArches[arch] {
		return fmt.Errorf("invalid architecture: %s", arch)
	}

	// Check if directory already exists
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory %s already exists", projectDir)
	}

	// Create project directory
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Create forge.yml
	if err := createForgeYml(projectDir, template, arch); err != nil {
		return fmt.Errorf("failed to create forge.yml: %v", err)
	}

	// Create README.md
	if err := createReadme(projectDir, template); err != nil {
		return fmt.Errorf("failed to create README.md: %v", err)
	}

	// Create .gitignore
	if err := createGitignore(projectDir); err != nil {
		return fmt.Errorf("failed to create .gitignore: %v", err)
	}

	// Initialize git repository if requested
	// TODO: Implement git initialization

	return nil
}

// createForgeYml creates the forge.yml configuration file
func createForgeYml(projectDir, template, arch string) error {
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")

	content := fmt.Sprintf(`schema_version: "1.0"
name: "%s"
version: "0.1.0"
architecture: "%s"
template: "%s"

buildroot:
  version: "stable"

kernel:
  version: "latest"

packages: []

features: []

overlays: {}
`, filepath.Base(projectDir), arch, template)

	return os.WriteFile(forgeYmlPath, []byte(content), 0644)
}

// createReadme creates the README.md file
func createReadme(projectDir, template string) error {
	readmePath := filepath.Join(projectDir, "README.md")

	content := fmt.Sprintf(`# %s

This is a Forge OS project created with the %s template.

## Getting Started

1. Build the project:
   `+"```bash"+`
   forge build
   `+"```"+`

2. Test in QEMU:
   `+"```bash"+`
   forge test
   `+"```"+`

3. Deploy to target:
   `+"```bash"+`
   forge deploy usb --device /dev/sdb
   `+"```"+`

## Configuration

Edit `+"`forge.yml`"+` to customize your OS configuration.

## Documentation

For more information, visit: https://forge-os.dev
`, filepath.Base(projectDir), template)

	return os.WriteFile(readmePath, []byte(content), 0644)
}

// createGitignore creates the .gitignore file
func createGitignore(projectDir string) error {
	gitignorePath := filepath.Join(projectDir, ".gitignore")

	content := `# Build artifacts
*.img
*.qcow2
*.raw
*.iso
build/
output/
dl/
.ccache/
.cache/

# Logs
*.log
.forge/logs/

# Temporary files
*.tmp
*.swp
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
`

	return os.WriteFile(gitignorePath, []byte(content), 0644)
}
