package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DeployCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestDeployCommandTestSuite(t *testing.T) {
	suite.Run(t, new(DeployCommandTestSuite))
}

func (s *DeployCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-deploy-cmd-*")
	s.Require().NoError(err)
}

func (s *DeployCommandTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DeployCommandTestSuite) TestDeployCommandCreation() {
	cmd := NewDeployCommand()
	s.NotNil(cmd)
	s.Equal("deploy [target]", cmd.Use)
	s.Contains(cmd.Short, "Deploy the Forge OS image")
}

func (s *DeployCommandTestSuite) TestDeployCommandUSB() {
	projectDir := filepath.Join(s.tempDir, "usb-deploy-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate build
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	// Create a mock device file
	mockDevice := filepath.Join(s.tempDir, "mock-sdb")
	err = os.WriteFile(mockDevice, []byte("mock device"), 0644)
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"usb"}, map[string]string{
		"device": mockDevice,
	})
	s.NoError(err)
}

func (s *DeployCommandTestSuite) TestDeployCommandSDCard() {
	projectDir := filepath.Join(s.tempDir, "sd-deploy-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate build
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	// Create a mock device file
	mockDevice := filepath.Join(s.tempDir, "mock-mmcblk0")
	err = os.WriteFile(mockDevice, []byte("mock device"), 0644)
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"sd"}, map[string]string{
		"device": mockDevice,
	})
	s.NoError(err)
}

func (s *DeployCommandTestSuite) TestDeployCommandRemote() {
	projectDir := filepath.Join(s.tempDir, "remote-deploy-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate build
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"remote"}, map[string]string{
		"host": "192.168.1.100",
		"user": "pi",
	})
	s.NoError(err)
}

func (s *DeployCommandTestSuite) TestDeployCommandNoBuildArtifacts() {
	projectDir := filepath.Join(s.tempDir, "no-build-deploy-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"usb"}, map[string]string{
		"device": "/dev/sdb",
	})
	s.Error(err)
	s.Contains(err.Error(), "no build artifacts found")
}

func (s *DeployCommandTestSuite) TestDeployCommandNoConfigFile() {
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(s.tempDir)

	err := runDeployCommand([]string{"usb"}, map[string]string{
		"device": "/dev/sdb",
	})
	s.Error(err)
	s.Contains(err.Error(), "no forge.yml found")
}

func (s *DeployCommandTestSuite) TestDeployCommandInvalidTarget() {
	projectDir := filepath.Join(s.tempDir, "invalid-target-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate build
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"invalid"}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "unsupported deployment target")
}

func (s *DeployCommandTestSuite) TestDeployCommandUSBNoDevice() {
	projectDir := filepath.Join(s.tempDir, "usb-no-device-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate build
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"usb"}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "device not specified")
}

func (s *DeployCommandTestSuite) TestDeployCommandRemoteNoHost() {
	projectDir := filepath.Join(s.tempDir, "remote-no-host-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate build
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runDeployCommand([]string{"remote"}, map[string]string{})
	s.Error(err)
	s.Contains(err.Error(), "host not specified")
}
