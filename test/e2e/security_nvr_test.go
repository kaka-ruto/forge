//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SecurityNVRE2ETestSuite struct {
	suite.Suite
	tempDir string
}

func TestSecurityNVRE2ETestSuite(t *testing.T) {
	suite.Run(t, new(SecurityNVRE2ETestSuite))
}

func (s *SecurityNVRE2ETestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-e2e-security-nvr-*")
	s.Require().NoError(err)
}

func (s *SecurityNVRE2ETestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *SecurityNVRE2ETestSuite) TestSecurityNVRUseCase() {
	projectName := "security-nvr"
	projectDir := filepath.Join(s.tempDir, projectName)

	// Step 1: Create project with security template
	s.T().Log("Step 1: Creating security NVR project")
	err := runForgeCommand(s.tempDir, "new", projectName, "--template", "security", "--arch", "x86_64")
	s.NoError(err)

	// Step 2: Add video and storage packages
	s.T().Log("Step 2: Adding NVR packages")
	err = runForgeCommand(projectDir, "add", "package", "ffmpeg")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "motion")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "nginx")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "samba")
	s.NoError(err)

	// Step 3: Add features
	s.T().Log("Step 3: Adding security features")
	err = runForgeCommand(projectDir, "add", "feature", "web-server")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "ssh-hardening")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "firewall")
	s.NoError(err)

	// Step 4: Verify configuration
	content, err := os.ReadFile(filepath.Join(projectDir, "forge.yml"))
	s.NoError(err)
	contentStr := string(content)

	s.Contains(contentStr, "ffmpeg")
	s.Contains(contentStr, "motion")
	s.Contains(contentStr, "nginx")
	s.Contains(contentStr, "samba")
	s.Contains(contentStr, "web-server")
	s.Contains(contentStr, "ssh-hardening")
	s.Contains(contentStr, "firewall")

	// Step 5: Test build
	err = runForgeCommand(projectDir, "build", "--timeout", "30s")
	s.Error(err, "build should fail in test environment")

	s.T().Log("Security NVR use case E2E test completed")
}
