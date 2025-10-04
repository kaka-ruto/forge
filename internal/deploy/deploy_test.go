package deploy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type DeployTestSuite struct {
	suite.Suite
	config       *config.Config
	tempDir      string
	artifactsDir string
}

func TestDeployTestSuite(t *testing.T) {
	suite.Run(t, new(DeployTestSuite))
}

func (s *DeployTestSuite) SetupTest() {
	s.config = &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{},
		Features:      []string{},
	}

	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-deploy-test-*")
	s.Require().NoError(err)

	s.artifactsDir = filepath.Join(s.tempDir, "artifacts")
	err = os.MkdirAll(s.artifactsDir, 0755)
	s.Require().NoError(err)

	// Create mock artifacts
	kernelPath := filepath.Join(s.artifactsDir, "bzImage")
	rootfsPath := filepath.Join(s.artifactsDir, "rootfs.ext4")

	err = os.WriteFile(kernelPath, []byte("mock kernel"), 0644)
	s.Require().NoError(err)

	err = os.WriteFile(rootfsPath, []byte("mock rootfs"), 0644)
	s.Require().NoError(err)
}

func (s *DeployTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DeployTestSuite) TestNewDeploymentOrchestrator() {
	orchestrator := NewDeploymentOrchestrator(s.config)
	s.NotNil(orchestrator)
	s.NotNil(orchestrator.deployers)
}

func (s *DeployTestSuite) TestRegisterDeployer() {
	orchestrator := NewDeploymentOrchestrator(s.config)
	deployer := NewUSBDeployer()

	orchestrator.RegisterDeployer(TargetUSB, deployer)
	s.Contains(orchestrator.deployers, TargetUSB)
}

func (s *DeployTestSuite) TestGetAvailableTargets() {
	orchestrator := NewDeploymentOrchestrator(s.config)
	deployer := NewUSBDeployer()

	orchestrator.RegisterDeployer(TargetUSB, deployer)
	targets := orchestrator.GetAvailableTargets()

	s.Contains(targets, TargetUSB)
}

func (s *DeployTestSuite) TestValidateDeployment() {
	orchestrator := NewDeploymentOrchestrator(s.config)
	deployer := NewUSBDeployer()
	orchestrator.RegisterDeployer(TargetUSB, deployer)

	config := &DeploymentConfig{
		Target: TargetUSB,
		Device: "/dev/sdb",
	}

	err := orchestrator.ValidateDeployment(config)
	s.NoError(err) // Should pass validation even if device doesn't exist (basic validation)
}

func (s *DeployTestSuite) TestValidateDeploymentInvalidTarget() {
	orchestrator := NewDeploymentOrchestrator(s.config)

	config := &DeploymentConfig{
		Target: TargetUSB, // Not registered
	}

	err := orchestrator.ValidateDeployment(config)
	s.Error(err)
	s.Contains(err.Error(), "no deployer registered")
}

func (s *DeployTestSuite) TestExecuteDeploymentDryRun() {
	orchestrator := NewDeploymentOrchestrator(s.config)
	deployer := NewUSBDeployer()
	orchestrator.RegisterDeployer(TargetUSB, deployer)

	config := &DeploymentConfig{
		Target: TargetUSB,
		Device: "/dev/sdb",
		DryRun: true,
	}

	result, err := orchestrator.ExecuteDeployment(s.artifactsDir, config)
	s.NoError(err)
	s.True(result.Success)
	s.Contains(result.Details, "Dry run completed")
}

func (s *DeployTestSuite) TestExecuteDeploymentMissingArtifacts() {
	orchestrator := NewDeploymentOrchestrator(s.config)

	config := &DeploymentConfig{
		Target: TargetUSB,
		Device: "/dev/sdb",
	}

	// Test with non-existent artifacts directory
	result, err := orchestrator.ExecuteDeployment("/nonexistent", config)
	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "artifacts directory does not exist")
}

func (s *DeployTestSuite) TestValidateArtifacts() {
	err := ValidateArtifacts(s.artifactsDir)
	s.NoError(err)
}

func (s *DeployTestSuite) TestValidateArtifactsMissingKernel() {
	// Remove kernel file
	kernelPath := filepath.Join(s.artifactsDir, "bzImage")
	os.Remove(kernelPath)

	err := ValidateArtifacts(s.artifactsDir)
	s.Error(err)
	s.Contains(err.Error(), "required artifact missing")
}

func (s *DeployTestSuite) TestUSBDeployerValidate() {
	deployer := NewUSBDeployer()

	// Test valid config
	config := &DeploymentConfig{
		Device: "/dev/sdb",
	}
	err := deployer.Validate(config)
	s.NoError(err)

	// Test missing device
	config.Device = ""
	err = deployer.Validate(config)
	s.Error(err)
	s.Contains(err.Error(), "device not specified")
}

func (s *DeployTestSuite) TestRemoteDeployerValidate() {
	deployer := NewRemoteDeployer()

	// Test valid config
	config := &DeploymentConfig{
		Host: "192.168.1.100",
		User: "root",
	}
	err := deployer.Validate(config)
	s.NoError(err)

	// Test missing host
	config.Host = ""
	err = deployer.Validate(config)
	s.Error(err)
	s.Contains(err.Error(), "host not specified")
}
