//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NetworkingE2ETestSuite struct {
	suite.Suite
	tempDir string
}

func TestNetworkingE2ETestSuite(t *testing.T) {
	suite.Run(t, new(NetworkingE2ETestSuite))
}

func (s *NetworkingE2ETestSuite) SetupTest() {
	s.tempDir = "/Users/kaka/Code/go/forge"
}

func (s *NetworkingE2ETestSuite) TearDownTest() {
	testProjects := []string{"networking-test"}
	for _, project := range testProjects {
		projectDir := filepath.Join(s.tempDir, project)
		os.RemoveAll(projectDir)
	}
}

func (s *NetworkingE2ETestSuite) TestNetworkingTemplateWorkflow() {
	projectName := "networking-test"
	projectDir := filepath.Join(s.tempDir, projectName)

	// Create project with networking template
	err := runForgeCommand(s.tempDir, "new", projectName, "--template", "networking", "--arch", "x86_64")
	s.NoError(err)

	// Verify project structure
	s.DirExists(projectDir)
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	s.FileExists(forgeYmlPath)

	// Check that network interfaces file was created
	interfacesPath := filepath.Join(projectDir, "overlays/rootfs/etc/network/interfaces")
	s.FileExists(interfacesPath)

	// Verify forge.yml contains networking packages
	content, err := os.ReadFile(forgeYmlPath)
	s.NoError(err)
	contentStr := string(content)
	s.Contains(contentStr, "openssh")
	s.Contains(contentStr, "wpa_supplicant")

	// Test build
	err = runForgeCommand(projectDir, "build", "--timeout", "30s")
	s.Error(err, "build should fail in test environment")

	s.T().Log("Networking template E2E test completed")
}
