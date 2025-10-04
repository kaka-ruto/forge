package cli

import (
	"os"
	"testing"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type AddCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestAddCommandTestSuite(t *testing.T) {
	suite.Run(t, new(AddCommandTestSuite))
}

func (s *AddCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-add-test-*")
	s.Require().NoError(err)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	s.tempDir = oldDir + "/" + s.tempDir
	os.Chdir(s.tempDir)
}

func (s *AddCommandTestSuite) TearDownTest() {
	os.Chdir("/")
	os.RemoveAll(s.tempDir)
}

func (s *AddCommandTestSuite) TestNewAddCommand() {
	cmd := NewAddCommand()
	s.NotNil(cmd)
	s.Equal("add", cmd.Use)
	s.Contains(cmd.Short, "Add packages or features")
}

func (s *AddCommandTestSuite) TestAddPackageCommand() {
	// Create a test forge.yml
	testConfig := &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "1.0.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{"busybox"},
		Features:      []string{},
	}

	err := config.SaveConfig(testConfig, "forge.yml")
	s.NoError(err)

	// Test adding a package
	err = runAddPackageCommand([]string{"nginx"}, map[string]interface{}{})
	s.NoError(err)

	// Verify package was added
	cfg, err := config.LoadConfig("forge.yml")
	s.NoError(err)
	s.Contains(cfg.Packages, "nginx")
	s.Contains(cfg.Packages, "busybox") // Original package should still be there
}

func (s *AddCommandTestSuite) TestAddPackageDuplicate() {
	// Create a test forge.yml
	testConfig := &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "1.0.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{"nginx"},
		Features:      []string{},
	}

	err := config.SaveConfig(testConfig, "forge.yml")
	s.NoError(err)

	// Try to add the same package again
	err = runAddPackageCommand([]string{"nginx"}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "package 'nginx' is already added")
}

func (s *AddCommandTestSuite) TestAddPackageNoConfig() {
	// Don't create forge.yml
	err := runAddPackageCommand([]string{"nginx"}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "failed to load forge.yml")
}

func (s *AddCommandTestSuite) TestAddFeatureCommand() {
	// Create a test forge.yml
	testConfig := &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "1.0.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{},
		Features:      []string{"network"},
	}

	err := config.SaveConfig(testConfig, "forge.yml")
	s.NoError(err)

	// Test adding a feature
	err = runAddFeatureCommand([]string{"firewall"}, map[string]interface{}{})
	s.NoError(err)

	// Verify feature was added
	cfg, err := config.LoadConfig("forge.yml")
	s.NoError(err)
	s.Contains(cfg.Features, "firewall")
	s.Contains(cfg.Features, "network") // Original feature should still be there
}

func (s *AddCommandTestSuite) TestAddFeatureDuplicate() {
	// Create a test forge.yml
	testConfig := &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "1.0.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{},
		Features:      []string{"firewall"},
	}

	err := config.SaveConfig(testConfig, "forge.yml")
	s.NoError(err)

	// Try to add the same feature again
	err = runAddFeatureCommand([]string{"firewall"}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "feature 'firewall' is already added")
}

func (s *AddCommandTestSuite) TestAddFeatureNoConfig() {
	// Don't create forge.yml
	err := runAddFeatureCommand([]string{"firewall"}, map[string]interface{}{})
	s.Error(err)
	s.Contains(err.Error(), "failed to load forge.yml")
}
