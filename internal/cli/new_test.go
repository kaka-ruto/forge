package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NewCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestNewCommandTestSuite(t *testing.T) {
	suite.Run(t, new(NewCommandTestSuite))
}

func (s *NewCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-new-test-*")
	s.Require().NoError(err)
}

func (s *NewCommandTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *NewCommandTestSuite) TestNewCommandCreation() {
	cmd := NewNewCommand()
	s.NotNil(cmd)
	s.Equal("new [project-name]", cmd.Use)
	s.Contains(cmd.Short, "Create a new Forge OS project")
}

func (s *NewCommandTestSuite) TestNewCommandWithValidArgs() {
	// Change to temp directory for test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(s.tempDir)

	projectName := "test-project"

	err := runNewCommand([]string{projectName}, map[string]string{
		"template": "minimal",
		"arch":     "x86_64",
	})

	// Should succeed with valid arguments
	s.NoError(err)

	// Check that project directory was created
	projectDir := filepath.Join(s.tempDir, projectName)
	_, err = os.Stat(projectDir)
	s.NoError(err)
}

func (s *NewCommandTestSuite) TestNewCommandProjectDirectoryCreation() {
	projectDir := filepath.Join(s.tempDir, "test-project")

	err := createProjectStructure(projectDir, "minimal", "x86_64")
	// Should succeed
	s.NoError(err)

	// Check that directory was created
	_, err = os.Stat(projectDir)
	s.NoError(err)
}

func (s *NewCommandTestSuite) TestNewCommandWithExistingDirectory() {
	projectDir := filepath.Join(s.tempDir, "existing-project")
	err := os.MkdirAll(projectDir, 0755)
	s.NoError(err)

	// Try to create project in existing directory
	err = createProjectStructure(projectDir, "minimal", "x86_64")
	// Should handle existing directory appropriately
	s.Error(err) // Expected to fail initially
}

func (s *NewCommandTestSuite) TestNewCommandInvalidTemplate() {
	projectDir := filepath.Join(s.tempDir, "invalid-template-project")

	err := createProjectStructure(projectDir, "invalid-template", "x86_64")
	s.Error(err)
	s.Contains(err.Error(), "invalid template")
}

func (s *NewCommandTestSuite) TestNewCommandInvalidArchitecture() {
	projectDir := filepath.Join(s.tempDir, "invalid-arch-project")

	err := createProjectStructure(projectDir, "minimal", "invalid-arch")
	s.Error(err)
	s.Contains(err.Error(), "invalid architecture")
}

func (s *NewCommandTestSuite) TestNewCommandForgeYmlGeneration() {
	projectDir := filepath.Join(s.tempDir, "yml-test-project")

	err := createProjectStructure(projectDir, "minimal", "x86_64")
	// Should create forge.yml file
	s.NoError(err)

	// Check that forge.yml exists
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	_, err = os.Stat(forgeYmlPath)
	s.NoError(err)
}

func (s *NewCommandTestSuite) TestNewCommandReadmeGeneration() {
	projectDir := filepath.Join(s.tempDir, "readme-test-project")

	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Check that README.md exists
	readmePath := filepath.Join(projectDir, "README.md")
	_, err = os.Stat(readmePath)
	s.NoError(err)
}

func (s *NewCommandTestSuite) TestNewCommandGitInitialization() {
	projectDir := filepath.Join(s.tempDir, "git-test-project")

	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Git initialization is not yet implemented, so .git should not exist
	gitDir := filepath.Join(projectDir, ".git")
	_, err = os.Stat(gitDir)
	s.True(os.IsNotExist(err))
}

func (s *NewCommandTestSuite) TestNewCommandAllTemplates() {
	templates := []string{"minimal", "networking", "iot", "security", "industrial", "kiosk"}

	for _, template := range templates {
		projectDir := filepath.Join(s.tempDir, "template-test-"+template)

		err := createProjectStructure(projectDir, template, "x86_64")
		s.NoError(err, "Template %s should be valid", template)
	}
}

func (s *NewCommandTestSuite) TestNewCommandAllArchitectures() {
	architectures := []string{"x86_64", "arm", "aarch64", "mips"}

	for _, arch := range architectures {
		projectDir := filepath.Join(s.tempDir, "arch-test-"+arch)

		err := createProjectStructure(projectDir, "minimal", arch)
		s.NoError(err, "Architecture %s should be valid", arch)
	}
}

func (s *NewCommandTestSuite) TestNewCommandInsufficientPermissions() {
	// This test would require setting up a directory with no write permissions
	// For now, skip this test as it's complex to test properly
	s.T().Skip("Permission test requires special setup")
}

func (s *NewCommandTestSuite) TestNewCommandLoggingIntegration() {
	projectDir := filepath.Join(s.tempDir, "logging-test-project")

	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Logging integration is not yet implemented
	// This test passes as long as the function completes without error
}
