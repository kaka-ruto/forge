package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/templates"
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
	// Check if directory already exists
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory %s already exists", projectDir)
	}

	// Use template manager to create project
	tm := templates.NewTemplateManager()

	// Prepare template data
	data := map[string]interface{}{
		"ProjectName":  filepath.Base(projectDir),
		"Architecture": arch,
	}

	// Apply template
	if err := tm.ApplyTemplate(template, projectDir, data); err != nil {
		return fmt.Errorf("failed to apply template: %v", err)
	}

	// Create .gitignore
	if err := createGitignore(projectDir); err != nil {
		return fmt.Errorf("failed to create .gitignore: %v", err)
	}

	// Initialize git repository if requested
	// TODO: Implement git initialization

	return nil
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
