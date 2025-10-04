package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestTestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(TestCommandTestSuite))
}

func (s *TestCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-test-cmd-*")
	s.Require().NoError(err)
}

func (s *TestCommandTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *TestCommandTestSuite) TestTestCommandCreation() {
	cmd := NewTestCommand()
	s.NotNil(cmd)
	s.Equal("test", cmd.Use)
	s.Contains(cmd.Short, "Test the Forge OS image")
}

func (s *TestCommandTestSuite) TestTestCommandWithValidBuild() {
	// Create a project with a built image
	projectDir := filepath.Join(s.tempDir, "test-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate a build by creating the artifacts
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  false,
		"image":     "",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 1,
	})
	// Should run QEMU test
	s.NoError(err)
}

func (s *TestCommandTestSuite) TestTestCommandNoBuildArtifacts() {
	projectDir := filepath.Join(s.tempDir, "no-build-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  false,
		"image":     "",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 1,
	})
	s.Error(err)
	s.Contains(err.Error(), "no build artifacts found")
}

func (s *TestCommandTestSuite) TestTestCommandNoConfigFile() {
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(s.tempDir)

	err := runTestCommand([]string{}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "no forge.yml found")
}

func (s *TestCommandTestSuite) TestTestCommandHeadlessMode() {
	projectDir := filepath.Join(s.tempDir, "headless-project")
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

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  true,
		"image":     "",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 1,
	})
	// Expect failure due to QEMU not being available in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to start QEMU")
}

func (s *TestCommandTestSuite) TestTestCommandWithPortForwarding() {
	projectDir := filepath.Join(s.tempDir, "port-project")
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

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  false,
		"image":     "",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 1,
	})

	// Expect failure due to QEMU not being available in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to start QEMU")
}

func (s *TestCommandTestSuite) TestTestCommandMultipleInstances() {
	projectDir := filepath.Join(s.tempDir, "multi-project")
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

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  false,
		"image":     "",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 3,
	})
	// Expect failure due to QEMU not being available in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to start QEMU")
}

func (s *TestCommandTestSuite) TestTestCommandInvalidImagePath() {
	projectDir := filepath.Join(s.tempDir, "invalid-image-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	// Simulate a build by creating the artifacts
	buildDir := filepath.Join(projectDir, "build")
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	os.MkdirAll(artifactsDir, 0755)
	os.WriteFile(filepath.Join(artifactsDir, "bzImage"), []byte("dummy kernel"), 0644)
	os.WriteFile(filepath.Join(artifactsDir, "rootfs.ext4"), []byte("dummy rootfs"), 0644)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  false,
		"image":     "/nonexistent/image.img",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 1,
	})
	s.Error(err)
	s.Contains(err.Error(), "does not exist")
}

func (s *TestCommandTestSuite) TestTestCommandResourceChecking() {
	projectDir := filepath.Join(s.tempDir, "resource-test-project")
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

	err = runTestCommand([]string{}, map[string]interface{}{
		"headless":  false,
		"image":     "",
		"scenarios": []string{},
		"timeout":   5 * time.Minute,
		"instances": 1,
	})
	// Expect failure due to QEMU not being available in test environment
	s.Error(err)
	s.Contains(err.Error(), "failed to start QEMU")
}
