package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PackagesCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestPackagesCommandTestSuite(t *testing.T) {
	suite.Run(t, new(PackagesCommandTestSuite))
}

func (s *PackagesCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-packages-test-*")
	s.Require().NoError(err)
}

func (s *PackagesCommandTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *PackagesCommandTestSuite) TestPackagesCommandCreation() {
	cmd := NewPackagesCommand()
	s.NotNil(cmd)
	s.Equal("packages", cmd.Use)
	s.Contains(cmd.Short, "Manage Forge OS packages")
}

func (s *PackagesCommandTestSuite) TestPackagesInstallCommand() {
	cmd := newPackagesInstallCommand()
	s.NotNil(cmd)
	s.Equal("install [packages...]", cmd.Use)
	s.Contains(cmd.Short, "Install packages")
}

func (s *PackagesCommandTestSuite) TestPackagesUninstallCommand() {
	cmd := newPackagesUninstallCommand()
	s.NotNil(cmd)
	s.Equal("uninstall [packages...]", cmd.Use)
	s.Contains(cmd.Short, "Uninstall packages")
}

func (s *PackagesCommandTestSuite) TestPackagesListCommandCreation() {
	cmd := newPackagesListCommand()
	s.NotNil(cmd)
	s.Equal("list [category]", cmd.Use)
	s.Contains(cmd.Short, "List available packages")
}

func (s *PackagesCommandTestSuite) TestPackagesInfoCommandCreation() {
	cmd := newPackagesInfoCommand()
	s.NotNil(cmd)
	s.Equal("info [package]", cmd.Use)
	s.Contains(cmd.Short, "Show package information")
}

func (s *PackagesCommandTestSuite) TestPackagesInstallCommandWithValidProject() {
	// Create a project with build directory
	projectDir := filepath.Join(s.tempDir, "test-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Create build directory with Buildroot
	buildDir := filepath.Join(projectDir, "build")
	os.MkdirAll(buildDir, 0755)

	// Create mock Buildroot directory
	buildrootDir := filepath.Join(buildDir, "buildroot")
	os.MkdirAll(buildrootDir, 0755)

	// Create mock config
	configPath := filepath.Join(buildrootDir, ".config")
	configContent := `# Buildroot config
BR2_PACKAGE_BUSYBOX=y
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runPackagesInstallCommand([]string{"busybox"}, map[string]interface{}{
		"buildroot": "",
	})
	// Should work (busybox is already enabled)
	s.NoError(err)
}

func (s *PackagesCommandTestSuite) TestPackagesInstallCommandNoProject() {
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(s.tempDir)

	err := runPackagesInstallCommand([]string{"busybox"}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "no forge.yml found")
}

func (s *PackagesCommandTestSuite) TestPackagesListCommand() {
	// Create a minimal project
	projectDir := filepath.Join(s.tempDir, "list-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runPackagesListCommand([]string{}, map[string]interface{}{})
	s.NoError(err)
}

func (s *PackagesCommandTestSuite) TestPackagesInfoCommand() {
	// Create a minimal project
	projectDir := filepath.Join(s.tempDir, "info-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runPackagesInfoCommand([]string{"busybox"}, map[string]interface{}{})
	s.NoError(err)
}

func (s *PackagesCommandTestSuite) TestPackagesInfoCommandInvalidPackage() {
	// Create a minimal project
	projectDir := filepath.Join(s.tempDir, "info-invalid-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runPackagesInfoCommand([]string{"nonexistent"}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "package nonexistent not found")
}

// Wrapper functions for testing
func runPackagesInstallCommand(args []string, flags map[string]interface{}) error {
	cmd := newPackagesInstallCommand()
	if buildroot, ok := flags["buildroot"].(string); ok && buildroot != "" {
		cmd.Flags().Set("buildroot", buildroot)
	}
	return runPackagesInstallCommandE(cmd, args)
}

func runPackagesListCommand(args []string, flags map[string]interface{}) error {
	cmd := newPackagesListCommand()
	return runPackagesListCommandE(cmd, args)
}

func runPackagesInfoCommand(args []string, flags map[string]interface{}) error {
	cmd := newPackagesInfoCommand()
	return runPackagesInfoCommandE(cmd, args)
}
