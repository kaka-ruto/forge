package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the main Forge OS configuration
type Config struct {
	SchemaVersion string                 `yaml:"schema_version" validate:"required"`
	Name          string                 `yaml:"name" validate:"required"`
	Version       string                 `yaml:"version" validate:"required"`
	Architecture  string                 `yaml:"architecture" validate:"required"`
	Template      string                 `yaml:"template" validate:"required"`
	Buildroot     BuildrootConfig        `yaml:"buildroot"`
	Kernel        KernelConfig           `yaml:"kernel"`
	Packages      []string               `yaml:"packages"`
	Features      []string               `yaml:"features"`
	Overlays      map[string]interface{} `yaml:"overlays"`
}

// BuildrootConfig represents Buildroot-specific configuration
type BuildrootConfig struct {
	Version string `yaml:"version"`
}

// KernelConfig represents kernel-specific configuration
type KernelConfig struct {
	Version string            `yaml:"version"`
	Config  map[string]string `yaml:"config"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.SchemaVersion == "" {
		return fmt.Errorf("schema_version is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}
	if c.Architecture == "" {
		return fmt.Errorf("architecture is required")
	}
	if c.Template == "" {
		return fmt.Errorf("template is required")
	}

	// Validate architecture
	validArchs := []string{"x86_64", "arm", "aarch64", "riscv64", "i386", "armv7", "armv5", "mips"}
	if !contains(validArchs, c.Architecture) {
		return fmt.Errorf("invalid architecture: %s (valid: %s)", c.Architecture, strings.Join(validArchs, ", "))
	}

	// Validate template
	validTemplates := []string{"minimal", "networking", "iot", "security", "industrial", "kiosk"}
	if !contains(validTemplates, c.Template) {
		return fmt.Errorf("invalid template: %s (valid: %s)", c.Template, strings.Join(validTemplates, ", "))
	}

	return nil
}

// LoadConfig loads and parses a forge.yml configuration file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to a forge.yml file
func SaveConfig(config *Config, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetBuildrootDefconfig generates a Buildroot defconfig from the configuration
func (c *Config) GetBuildrootDefconfig() (string, error) {
	var defconfig strings.Builder

	// Base configuration
	defconfig.WriteString("# Forge OS Buildroot defconfig\n")
	defconfig.WriteString("# Generated from forge.yml\n\n")

	// Architecture-specific settings
	switch c.Architecture {
	case "x86_64":
		defconfig.WriteString("BR2_x86_64=y\n")
		defconfig.WriteString("BR2_ARCH=\"x86_64\"\n")
	case "arm":
		defconfig.WriteString("BR2_arm=y\n")
		defconfig.WriteString("BR2_ARCH=\"arm\"\n")
	case "aarch64":
		defconfig.WriteString("BR2_aarch64=y\n")
		defconfig.WriteString("BR2_ARCH=\"aarch64\"\n")
	case "riscv64":
		defconfig.WriteString("BR2_riscv=y\n")
		defconfig.WriteString("BR2_ARCH=\"riscv\"\n")
	case "i386":
		defconfig.WriteString("BR2_i386=y\n")
		defconfig.WriteString("BR2_ARCH=\"i386\"\n")
	case "armv7":
		defconfig.WriteString("BR2_arm=y\n")
		defconfig.WriteString("BR2_ARCH=\"arm\"\n")
		defconfig.WriteString("BR2_ARM_CPU_ARMV7A=y\n")
	case "armv5":
		defconfig.WriteString("BR2_arm=y\n")
		defconfig.WriteString("BR2_ARCH=\"arm\"\n")
		defconfig.WriteString("BR2_ARM_CPU_ARMV5=y\n")
	}

	// Toolchain
	defconfig.WriteString("BR2_TOOLCHAIN_BUILDROOT_GLIBC=y\n")

	// Template-specific packages
	switch c.Template {
	case "minimal":
		defconfig.WriteString("BR2_PACKAGE_BUSYBOX=y\n")
		defconfig.WriteString("BR2_TARGET_ROOTFS_EXT2=y\n")
	case "networking":
		defconfig.WriteString("BR2_PACKAGE_BUSYBOX=y\n")
		defconfig.WriteString("BR2_PACKAGE_DROPBEAR=y\n")
		defconfig.WriteString("BR2_PACKAGE_WPA_SUPPLICANT=y\n")
		defconfig.WriteString("BR2_TARGET_ROOTFS_EXT2=y\n")
	case "iot":
		defconfig.WriteString("BR2_PACKAGE_BUSYBOX=y\n")
		defconfig.WriteString("BR2_PACKAGE_MOSQUITTO=y\n")
		defconfig.WriteString("BR2_TARGET_ROOTFS_EXT2=y\n")
	case "security":
		defconfig.WriteString("BR2_PACKAGE_BUSYBOX=y\n")
		defconfig.WriteString("BR2_PACKAGE_DROPBEAR=y\n")
		defconfig.WriteString("BR2_PACKAGE_OPENVPN=y\n")
		defconfig.WriteString("BR2_TARGET_ROOTFS_EXT2=y\n")
	case "industrial":
		defconfig.WriteString("BR2_PACKAGE_BUSYBOX=y\n")
		defconfig.WriteString("BR2_PACKAGE_MODBUS=y\n")
		defconfig.WriteString("BR2_TARGET_ROOTFS_EXT2=y\n")
	case "kiosk":
		defconfig.WriteString("BR2_PACKAGE_BUSYBOX=y\n")
		defconfig.WriteString("BR2_PACKAGE_XORG7=y\n")
		defconfig.WriteString("BR2_PACKAGE_CHROMIUM=y\n")
		defconfig.WriteString("BR2_TARGET_ROOTFS_EXT2=y\n")
	}

	// Additional packages
	for _, pkg := range c.Packages {
		defconfig.WriteString(fmt.Sprintf("BR2_PACKAGE_%s=y\n", strings.ToUpper(pkg)))
	}

	// Features
	for _, feature := range c.Features {
		switch feature {
		case "systemd":
			defconfig.WriteString("BR2_INIT_SYSTEMD=y\n")
		case "sysvinit":
			defconfig.WriteString("BR2_INIT_SYSV=y\n")
		case "network":
			defconfig.WriteString("BR2_SYSTEM_ENABLE_NLS=y\n")
		case "debug":
			defconfig.WriteString("BR2_ENABLE_DEBUG=y\n")
		}
	}

	return defconfig.String(), nil
}

// GetKernelConfig generates kernel configuration from the config
func (c *Config) GetKernelConfig() (string, error) {
	var kernelConfig strings.Builder

	kernelConfig.WriteString("# Forge OS Kernel Configuration\n")
	kernelConfig.WriteString("# Generated from forge.yml\n\n")

	// Basic kernel config based on architecture
	switch c.Architecture {
	case "x86_64":
		kernelConfig.WriteString("CONFIG_64BIT=y\n")
		kernelConfig.WriteString("CONFIG_X86_64=y\n")
	case "arm":
		kernelConfig.WriteString("CONFIG_ARM=y\n")
	case "aarch64":
		kernelConfig.WriteString("CONFIG_ARM64=y\n")
	}

	// Template-specific kernel features
	switch c.Template {
	case "minimal":
		kernelConfig.WriteString("CONFIG_EMBEDDED=y\n")
	case "networking":
		kernelConfig.WriteString("CONFIG_NET=y\n")
		kernelConfig.WriteString("CONFIG_INET=y\n")
	case "iot":
		kernelConfig.WriteString("CONFIG_EMBEDDED=y\n")
		kernelConfig.WriteString("CONFIG_I2C=y\n")
		kernelConfig.WriteString("CONFIG_SPI=y\n")
	}

	// Custom kernel config from forge.yml
	for key, value := range c.Kernel.Config {
		kernelConfig.WriteString(fmt.Sprintf("CONFIG_%s=%s\n", key, value))
	}

	return kernelConfig.String(), nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
