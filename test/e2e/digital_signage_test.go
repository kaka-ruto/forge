//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DigitalSignageE2ETestSuite struct {
	suite.Suite
	tempDir string
}

func TestDigitalSignageE2ETestSuite(t *testing.T) {
	suite.Run(t, new(DigitalSignageE2ETestSuite))
}

func (s *DigitalSignageE2ETestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-e2e-digital-signage-*")
	s.Require().NoError(err)
}

func (s *DigitalSignageE2ETestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DigitalSignageE2ETestSuite) TestDigitalSignageUseCase() {
	projectName := "digital-signage"
	projectDir := filepath.Join(s.tempDir, projectName)

	// Step 1: Create project with kiosk template
	s.T().Log("Step 1: Creating digital signage project")
	err := runForgeCommand(s.tempDir, "new", projectName, "--template", "kiosk", "--arch", "x86_64")
	s.NoError(err)

	// Verify project structure
	s.DirExists(projectDir)
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	s.FileExists(forgeYmlPath)

	// Step 2: Add display and browser packages
	s.T().Log("Step 2: Adding display packages")
	err = runForgeCommand(projectDir, "add", "package", "xorg-server")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "chromium")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "plymouth")
	s.NoError(err)

	// Step 3: Add management features
	s.T().Log("Step 3: Adding management features")
	err = runForgeCommand(projectDir, "add", "feature", "auto-updates")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "remote-management")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "web-dashboard")
	s.NoError(err)

	// Step 4: Verify configuration
	s.T().Log("Step 4: Verifying configuration")
	content, err := os.ReadFile(forgeYmlPath)
	s.NoError(err)
	contentStr := string(content)

	// Check packages
	s.Contains(contentStr, "xorg-server")
	s.Contains(contentStr, "chromium")
	s.Contains(contentStr, "plymouth")

	// Check features
	s.Contains(contentStr, "auto-updates")
	s.Contains(contentStr, "remote-management")
	s.Contains(contentStr, "web-dashboard")

	// Check template
	s.Contains(contentStr, "template: kiosk")

	// Step 5: Test build command
	s.T().Log("Step 5: Testing build command")
	err = runForgeCommand(projectDir, "build", "--timeout", "30s")
	s.Error(err, "build should fail in test environment")
	s.Contains(err.Error(), "build failed")

	s.T().Log("Digital signage use case E2E test completed")
}
