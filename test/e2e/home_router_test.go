//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HomeRouterE2ETestSuite struct {
	suite.Suite
	tempDir string
}

func TestHomeRouterE2ETestSuite(t *testing.T) {
	suite.Run(t, new(HomeRouterE2ETestSuite))
}

func (s *HomeRouterE2ETestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-e2e-home-router-*")
	s.Require().NoError(err)
}

func (s *HomeRouterE2ETestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *HomeRouterE2ETestSuite) TestHomeRouterUseCase() {
	projectName := "home-router"
	projectDir := filepath.Join(s.tempDir, projectName)

	// Step 1: Create project with networking template
	s.T().Log("Step 1: Creating home router project")
	err := runForgeCommand(s.tempDir, "new", projectName, "--template", "networking", "--arch", "x86_64")
	s.NoError(err)

	// Verify project structure
	s.DirExists(projectDir)
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	s.FileExists(forgeYmlPath)

	// Step 2: Add required packages for home router
	s.T().Log("Step 2: Adding router packages")
	err = runForgeCommand(projectDir, "add", "package", "dnsmasq")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "wireguard")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "iptables")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "hostapd")
	s.NoError(err)

	// Step 3: Add features
	s.T().Log("Step 3: Adding router features")
	err = runForgeCommand(projectDir, "add", "feature", "firewall")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "vpn-gateway")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "web-dashboard")
	s.NoError(err)

	// Step 4: Verify forge.yml was updated
	s.T().Log("Step 4: Verifying configuration")
	content, err := os.ReadFile(forgeYmlPath)
	s.NoError(err)
	contentStr := string(content)

	// Check packages were added
	s.Contains(contentStr, "dnsmasq")
	s.Contains(contentStr, "wireguard")
	s.Contains(contentStr, "iptables")
	s.Contains(contentStr, "hostapd")

	// Check features were added
	s.Contains(contentStr, "firewall")
	s.Contains(contentStr, "vpn-gateway")
	s.Contains(contentStr, "web-dashboard")

	// Step 5: Attempt build (will fail in test environment but should be accepted)
	s.T().Log("Step 5: Testing build command")
	err = runForgeCommand(projectDir, "build", "--timeout", "30s")
	s.Error(err, "build should fail in test environment")
	s.Contains(err.Error(), "build failed")

	s.T().Log("Home router use case E2E test completed")
}
