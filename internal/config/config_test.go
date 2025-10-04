package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	tempDir string
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-config-test-*")
	s.Require().NoError(err)
}

func (s *ConfigTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ConfigTestSuite) TestLoadConfigValid() {
	configPath := filepath.Join(s.tempDir, "forge.yml")
	configContent := `schema_version: "1.0"
name: "test-project"
version: "0.1.0"
architecture: "x86_64"
template: "minimal"

buildroot:
  version: "stable"

kernel:
  version: "latest"

packages: []
features: []
overlays: {}
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	s.NoError(err)

	config, err := LoadConfig(configPath)
	s.NoError(err)
	s.NotNil(config)
	s.Equal("1.0", config.SchemaVersion)
	s.Equal("test-project", config.Name)
	s.Equal("x86_64", config.Architecture)
	s.Equal("minimal", config.Template)
}

func (s *ConfigTestSuite) TestLoadConfigInvalidYAML() {
	configPath := filepath.Join(s.tempDir, "invalid.yml")
	err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644)
	s.NoError(err)

	config, err := LoadConfig(configPath)
	s.Error(err)
	s.Nil(config)
	s.Contains(err.Error(), "failed to parse config file")
}

func (s *ConfigTestSuite) TestLoadConfigMissingFile() {
	configPath := filepath.Join(s.tempDir, "missing.yml")

	config, err := LoadConfig(configPath)
	s.Error(err)
	s.Nil(config)
	s.Contains(err.Error(), "failed to read config file")
}

func (s *ConfigTestSuite) TestConfigValidation() {
	// Valid config
	config := &Config{
		SchemaVersion: "1.0",
		Name:          "test",
		Version:       "1.0.0",
		Architecture:  "x86_64",
		Template:      "minimal",
	}
	err := config.Validate()
	s.NoError(err)

	// Missing required fields
	config.SchemaVersion = ""
	err = config.Validate()
	s.Error(err)
	s.Contains(err.Error(), "schema_version is required")

	config.SchemaVersion = "1.0"
	config.Name = ""
	err = config.Validate()
	s.Error(err)
	s.Contains(err.Error(), "name is required")

	config.Name = "test"
	config.Architecture = "invalid"
	err = config.Validate()
	s.Error(err)
	s.Contains(err.Error(), "invalid architecture")

	config.Architecture = "x86_64"
	config.Template = "invalid"
	err = config.Validate()
	s.Error(err)
	s.Contains(err.Error(), "invalid template")
}

func (s *ConfigTestSuite) TestSaveConfig() {
	configPath := filepath.Join(s.tempDir, "output.yml")
	config := &Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Buildroot: BuildrootConfig{
			Version: "stable",
		},
		Kernel: KernelConfig{
			Version: "latest",
		},
		Packages: []string{},
		Features: []string{},
		Overlays: map[string]interface{}{},
	}

	err := SaveConfig(config, configPath)
	s.NoError(err)

	// Verify file was created and can be loaded
	loaded, err := LoadConfig(configPath)
	s.NoError(err)
	s.Equal(config.Name, loaded.Name)
	s.Equal(config.Architecture, loaded.Architecture)
}

func (s *ConfigTestSuite) TestGetBuildrootDefconfig() {
	config := &Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{"curl", "wget"},
		Features:      []string{"systemd"},
	}

	defconfig, err := config.GetBuildrootDefconfig()
	s.NoError(err)
	s.NotEmpty(defconfig)
	s.Contains(defconfig, "BR2_x86_64=y")
	s.Contains(defconfig, "BR2_PACKAGE_BUSYBOX=y")
	s.Contains(defconfig, "BR2_PACKAGE_CURL=y")
	s.Contains(defconfig, "BR2_PACKAGE_WGET=y")
	s.Contains(defconfig, "BR2_INIT_SYSTEMD=y")
}

func (s *ConfigTestSuite) TestGetBuildrootDefconfigTemplates() {
	templates := []string{"minimal", "networking", "iot", "security", "industrial", "kiosk"}

	for _, template := range templates {
		config := &Config{
			SchemaVersion: "1.0",
			Name:          "test",
			Version:       "1.0.0",
			Architecture:  "x86_64",
			Template:      template,
		}

		defconfig, err := config.GetBuildrootDefconfig()
		s.NoError(err)
		s.NotEmpty(defconfig)
		s.Contains(defconfig, "# Forge OS Buildroot defconfig")
	}
}

func (s *ConfigTestSuite) TestGetBuildrootDefconfigArchitectures() {
	architectures := []string{"x86_64", "arm", "aarch64", "riscv64", "i386", "armv7", "armv5"}

	for _, arch := range architectures {
		config := &Config{
			SchemaVersion: "1.0",
			Name:          "test",
			Version:       "1.0.0",
			Architecture:  arch,
			Template:      "minimal",
		}

		defconfig, err := config.GetBuildrootDefconfig()
		s.NoError(err)
		s.NotEmpty(defconfig)
	}
}

func (s *ConfigTestSuite) TestGetKernelConfig() {
	config := &Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Kernel: KernelConfig{
			Version: "latest",
			Config: map[string]string{
				"DEBUG_INFO": "y",
				"NETWORKING": "y",
			},
		},
	}

	kernelConfig, err := config.GetKernelConfig()
	s.NoError(err)
	s.NotEmpty(kernelConfig)
	s.Contains(kernelConfig, "# Forge OS Kernel Configuration")
	s.Contains(kernelConfig, "CONFIG_64BIT=y")
	s.Contains(kernelConfig, "CONFIG_EMBEDDED=y")
	s.Contains(kernelConfig, "CONFIG_DEBUG_INFO=y")
	s.Contains(kernelConfig, "CONFIG_NETWORKING=y")
}

func (s *ConfigTestSuite) TestGetKernelConfigTemplates() {
	templates := []string{"minimal", "networking", "iot"}

	for _, template := range templates {
		config := &Config{
			SchemaVersion: "1.0",
			Name:          "test",
			Version:       "1.0.0",
			Architecture:  "x86_64",
			Template:      template,
		}

		kernelConfig, err := config.GetKernelConfig()
		s.NoError(err)
		s.NotEmpty(kernelConfig)
		s.Contains(kernelConfig, "# Forge OS Kernel Configuration")
	}
}

func (s *ConfigTestSuite) TestConfigWithComplexFeatures() {
	config := &Config{
		SchemaVersion: "1.0",
		Name:          "complex-project",
		Version:       "1.0.0",
		Architecture:  "aarch64",
		Template:      "networking",
		Buildroot: BuildrootConfig{
			Version: "2023.11",
		},
		Kernel: KernelConfig{
			Version: "6.6",
			Config: map[string]string{
				"USB_SUPPORT": "y",
				"WIRELESS":    "y",
			},
		},
		Packages: []string{"openssh", "nginx", "python3"},
		Features: []string{"systemd", "network", "debug"},
		Overlays: map[string]interface{}{
			"rootfs": map[string]interface{}{
				"files": []string{"custom-script.sh"},
			},
		},
	}

	// Test validation
	err := config.Validate()
	s.NoError(err)

	// Test defconfig generation
	defconfig, err := config.GetBuildrootDefconfig()
	s.NoError(err)
	s.Contains(defconfig, "BR2_aarch64=y")
	s.Contains(defconfig, "BR2_PACKAGE_BUSYBOX=y")
	s.Contains(defconfig, "BR2_PACKAGE_DROPBEAR=y")
	s.Contains(defconfig, "BR2_PACKAGE_OPENSSH=y")
	s.Contains(defconfig, "BR2_PACKAGE_NGINX=y")
	s.Contains(defconfig, "BR2_PACKAGE_PYTHON3=y")
	s.Contains(defconfig, "BR2_INIT_SYSTEMD=y")
	s.Contains(defconfig, "BR2_SYSTEM_ENABLE_NLS=y")
	s.Contains(defconfig, "BR2_ENABLE_DEBUG=y")

	// Test kernel config generation
	kernelConfig, err := config.GetKernelConfig()
	s.NoError(err)
	s.Contains(kernelConfig, "CONFIG_ARM64=y")
	s.Contains(kernelConfig, "CONFIG_NET=y")
	s.Contains(kernelConfig, "CONFIG_INET=y")
	s.Contains(kernelConfig, "CONFIG_USB_SUPPORT=y")
	s.Contains(kernelConfig, "CONFIG_WIRELESS=y")
}
