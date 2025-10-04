//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IndustrialIoTGatewayE2ETestSuite struct {
	suite.Suite
	tempDir string
}

func TestIndustrialIoTGatewayE2ETestSuite(t *testing.T) {
	suite.Run(t, new(IndustrialIoTGatewayE2ETestSuite))
}

func (s *IndustrialIoTGatewayE2ETestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-e2e-industrial-iot-*")
	s.Require().NoError(err)
}

func (s *IndustrialIoTGatewayE2ETestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *IndustrialIoTGatewayE2ETestSuite) TestIndustrialIoTGatewayUseCase() {
	projectName := "industrial-gateway"
	projectDir := filepath.Join(s.tempDir, projectName)

	// Step 1: Create project with industrial template
	s.T().Log("Step 1: Creating industrial IoT gateway project")
	err := runForgeCommand(s.tempDir, "new", projectName, "--template", "industrial", "--arch", "arm")
	s.NoError(err)

	// Verify project structure
	s.DirExists(projectDir)
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	s.FileExists(forgeYmlPath)

	// Step 2: Add industrial protocols
	s.T().Log("Step 2: Adding industrial protocol packages")
	err = runForgeCommand(projectDir, "add", "package", "modbus")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "mqtt")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "package", "node-red")
	s.NoError(err)

	// Step 3: Add monitoring and reliability features
	s.T().Log("Step 3: Adding monitoring features")
	err = runForgeCommand(projectDir, "add", "feature", "auto-updates")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "monitoring")
	s.NoError(err)
	err = runForgeCommand(projectDir, "add", "feature", "watchdog")
	s.NoError(err)

	// Step 4: Verify configuration
	s.T().Log("Step 4: Verifying configuration")
	content, err := os.ReadFile(forgeYmlPath)
	s.NoError(err)
	contentStr := string(content)

	// Check packages
	s.Contains(contentStr, "modbus")
	s.Contains(contentStr, "mqtt")
	s.Contains(contentStr, "node-red")

	// Check features
	s.Contains(contentStr, "auto-updates")
	s.Contains(contentStr, "monitoring")
	s.Contains(contentStr, "watchdog")

	// Check architecture
	s.Contains(contentStr, "architecture: arm")

	// Step 5: Test build command
	s.T().Log("Step 5: Testing build command")
	err = runForgeCommand(projectDir, "build", "--timeout", "30s")
	s.Error(err, "build should fail in test environment")
	s.Contains(err.Error(), "build failed")

	s.T().Log("Industrial IoT gateway use case E2E test completed")
}
