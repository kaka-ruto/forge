package qemu

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type QEMUTestSuite struct {
	suite.Suite
	tempDir string
	config  *config.Config
	manager *QEMUManager
}

func TestQEMUTestSuite(t *testing.T) {
	suite.Run(t, new(QEMUTestSuite))
}

func (s *QEMUTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-qemu-test-*")
	s.Require().NoError(err)

	// Create a test config
	s.config = &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Buildroot: config.BuildrootConfig{
			Version: "stable",
		},
		Packages: []string{},
		Features: []string{},
	}

	s.manager = NewQEMUManager(s.config, s.tempDir)
}

func (s *QEMUTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *QEMUTestSuite) TestNewQEMUManager() {
	manager := NewQEMUManager(s.config, s.tempDir)
	s.NotNil(manager)
	s.Equal(s.config, manager.config)
	s.Equal(s.tempDir, manager.projectDir)
	s.NotNil(manager.logger)
}

func (s *QEMUTestSuite) TestStartInstance() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a dummy image file
	imagePath := filepath.Join(s.tempDir, "test-image.img")
	err := os.WriteFile(imagePath, []byte("dummy image"), 0644)
	s.NoError(err)

	// Try to start instance (will fail due to missing QEMU)
	instance, err := s.manager.StartInstance(ctx, imagePath)
	s.Error(err) // Expected to fail in test environment
	s.Nil(instance)
	s.Contains(err.Error(), "failed to start QEMU")
}

func (s *QEMUTestSuite) TestStopInstance() {
	// Create a mock instance
	instance := &QEMUInstance{
		ID: "test-instance",
	}

	// Stop should not panic even with nil process
	err := s.manager.StopInstance(instance)
	s.NoError(err)
}

func (s *QEMUTestSuite) TestGetDefaultTestScenarios() {
	scenarios := s.manager.GetDefaultTestScenarios()
	s.Len(scenarios, 5)

	expectedNames := []string{"boot", "network", "services", "performance", "stress"}
	for i, scenario := range scenarios {
		s.Equal(expectedNames[i], scenario.Name)
		s.NotEmpty(scenario.Description)
		s.NotZero(scenario.Timeout)
		s.NotNil(scenario.Run)
	}
}

func (s *QEMUTestSuite) TestRunTestScenario() {
	ctx := context.Background()

	// Create a mock instance
	instance := &QEMUInstance{
		ID: "test-instance",
	}

	// Create a simple test scenario
	scenario := TestScenario{
		Name:    "test-scenario",
		Timeout: 1 * time.Second,
		Run: func(ctx context.Context, instance *QEMUInstance) *TestResult {
			return &TestResult{
				TestName: "test-scenario",
				Success:  true,
				Output:   "Test passed",
			}
		},
	}

	result := s.manager.RunTestScenario(ctx, instance, scenario)
	s.NotNil(result)
	s.Equal("test-scenario", result.TestName)
	s.True(result.Success)
	s.Equal("Test passed", result.Output)
	s.NotZero(result.Duration)
}

func (s *QEMUTestSuite) TestBuildQEMUCommand() {
	instance := &QEMUInstance{
		ID:          "test-instance",
		MonitorPort: 4444,
		SSHPort:     2222,
		SerialPort:  8000,
	}

	imagePath := "/path/to/image.img"
	cmd := s.manager.buildQEMUCommand(instance, imagePath)

	s.Contains(cmd, "qemu-system-x86_64")
	s.Contains(cmd, "-machine")
	s.Contains(cmd, "pc")
	s.Contains(cmd, "-m")
	s.Contains(cmd, "512")
	s.Contains(cmd, "-kernel")
	s.Contains(cmd, imagePath)
}

func (s *QEMUTestSuite) TestGenerateInstanceID() {
	id1 := generateInstanceID()
	id2 := generateInstanceID()

	s.NotEmpty(id1)
	s.NotEmpty(id2)
	s.NotEqual(id1, id2)
	s.Contains(id1, "qemu-")
	s.Contains(id2, "qemu-")
}

func (s *QEMUTestSuite) TestFindAvailablePort() {
	port, err := findAvailablePort(8000, 8010)
	s.NoError(err)
	s.GreaterOrEqual(port, 8000)
	s.LessOrEqual(port, 8010)
}

func (s *QEMUTestSuite) TestTestScenarios() {
	ctx := context.Background()
	instance := &QEMUInstance{ID: "test"}

	// Test boot scenario
	bootResult := s.manager.testBootScenario(ctx, instance)
	s.NotNil(bootResult)
	s.Equal("boot", bootResult.TestName)

	// Test network scenario
	networkResult := s.manager.testNetworkScenario(ctx, instance)
	s.NotNil(networkResult)
	s.Equal("network", networkResult.TestName)
	s.True(networkResult.Success)

	// Test services scenario
	servicesResult := s.manager.testServicesScenario(ctx, instance)
	s.NotNil(servicesResult)
	s.Equal("services", servicesResult.TestName)
	s.True(servicesResult.Success)
}
