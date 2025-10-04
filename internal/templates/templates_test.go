package templates

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type TemplatesTestSuite struct {
	suite.Suite
	tempDir string
}

func TestTemplatesTestSuite(t *testing.T) {
	suite.Run(t, new(TemplatesTestSuite))
}

func (s *TemplatesTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-templates-test-*")
	s.Require().NoError(err)
}

func (s *TemplatesTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *TemplatesTestSuite) TestNewTemplateManager() {
	tm := NewTemplateManager()
	s.NotNil(tm)
	s.NotNil(tm.templates)
}

func (s *TemplatesTestSuite) TestGetTemplate() {
	tm := NewTemplateManager()

	// Test existing template
	template, err := tm.GetTemplate("minimal")
	s.NoError(err)
	s.NotNil(template)
	s.Equal("minimal", template.Name)
	s.Equal("Minimal Linux system with BusyBox", template.Description)

	// Test non-existing template
	template, err = tm.GetTemplate("nonexistent")
	s.Error(err)
	s.Nil(template)
	s.Contains(err.Error(), "template 'nonexistent' not found")
}

func (s *TemplatesTestSuite) TestListTemplates() {
	tm := NewTemplateManager()
	templates := tm.ListTemplates()

	s.NotNil(templates)
	s.Contains(templates, "minimal")
	s.Contains(templates, "networking")
	s.Contains(templates, "iot")
	s.Contains(templates, "security")
	s.Contains(templates, "industrial")
	s.Contains(templates, "kiosk")

	// Should have 6 templates
	s.Len(templates, 6)
}

func (s *TemplatesTestSuite) TestApplyTemplate() {
	tm := NewTemplateManager()
	projectDir := filepath.Join(s.tempDir, "test-project")

	err := tm.ApplyTemplate("minimal", projectDir, nil)
	s.NoError(err)

	// Check that forge.yml was created
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	s.FileExists(forgeYmlPath)

	// Check that README.md was created
	readmePath := filepath.Join(projectDir, "README.md")
	s.FileExists(readmePath)

	// Verify README content
	content, err := os.ReadFile(readmePath)
	s.NoError(err)
	s.Contains(string(content), "# test-project")
	s.Contains(string(content), "forge build")
}

func (s *TemplatesTestSuite) TestApplyTemplateNetworking() {
	tm := NewTemplateManager()
	projectDir := filepath.Join(s.tempDir, "networking-project")

	err := tm.ApplyTemplate("networking", projectDir, nil)
	s.NoError(err)

	// Check that network interface config was created
	interfacesPath := filepath.Join(projectDir, "overlays/rootfs/etc/network/interfaces")
	s.FileExists(interfacesPath)

	// Verify content
	content, err := os.ReadFile(interfacesPath)
	s.NoError(err)
	s.Contains(string(content), "auto eth0")
	s.Contains(string(content), "iface eth0 inet dhcp")
}

func (s *TemplatesTestSuite) TestApplyTemplateIOT() {
	tm := NewTemplateManager()
	projectDir := filepath.Join(s.tempDir, "iot-project")

	err := tm.ApplyTemplate("iot", projectDir, nil)
	s.NoError(err)

	// Check that mosquitto config was created
	mosquittoPath := filepath.Join(projectDir, "overlays/rootfs/etc/mosquitto/mosquitto.conf")
	s.FileExists(mosquittoPath)

	// Verify content
	content, err := os.ReadFile(mosquittoPath)
	s.NoError(err)
	s.Contains(string(content), "listener 1883")
}

func (s *TemplatesTestSuite) TestApplyTemplateKiosk() {
	tm := NewTemplateManager()
	projectDir := filepath.Join(s.tempDir, "kiosk-project")

	err := tm.ApplyTemplate("kiosk", projectDir, nil)
	s.NoError(err)

	// Check that X11 config was created
	xorgPath := filepath.Join(projectDir, "overlays/rootfs/etc/X11/xorg.conf")
	s.FileExists(xorgPath)

	// Verify content
	content, err := os.ReadFile(xorgPath)
	s.NoError(err)
	s.Contains(string(content), "Driver \"modesetting\"")
}

func (s *TemplatesTestSuite) TestApplyTemplateInvalid() {
	tm := NewTemplateManager()
	projectDir := filepath.Join(s.tempDir, "invalid-project")

	err := tm.ApplyTemplate("nonexistent", projectDir, nil)
	s.Error(err)
	s.Contains(err.Error(), "template 'nonexistent' not found")
}

func (s *TemplatesTestSuite) TestValidateTemplate() {
	tm := NewTemplateManager()

	// Valid template
	template := &Template{
		Name: "test",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "test",
			Version:       "1.0.0",
			Architecture:  "x86_64",
			Template:      "minimal",
		},
	}
	err := tm.ValidateTemplate(template)
	s.NoError(err)

	// Invalid template - no name
	template.Name = ""
	err = tm.ValidateTemplate(template)
	s.Error(err)
	s.Contains(err.Error(), "template name is required")

	// Invalid template - no config
	template.Name = "test"
	template.Config = nil
	err = tm.ValidateTemplate(template)
	s.Error(err)
	s.Contains(err.Error(), "template config is required")
}

func (s *TemplatesTestSuite) TestTemplateContentSubstitution() {
	tm := NewTemplateManager()
	projectDir := filepath.Join(s.tempDir, "substitution-project")

	// Apply template
	err := tm.ApplyTemplate("minimal", projectDir, nil)
	s.NoError(err)

	// Check forge.yml content
	forgeYmlPath := filepath.Join(projectDir, "forge.yml")
	content, err := os.ReadFile(forgeYmlPath)
	s.NoError(err)

	// Should contain the substituted project name
	s.Contains(string(content), "name: substitution-project")
	s.Contains(string(content), "architecture: x86_64")
}

func (s *TemplatesTestSuite) TestTemplateCategories() {
	tm := NewTemplateManager()
	templates := tm.ListTemplates()

	// Check categories
	s.Equal("basic", templates["minimal"].Category)
	s.Equal("networking", templates["networking"].Category)
	s.Equal("iot", templates["iot"].Category)
	s.Equal("security", templates["security"].Category)
	s.Equal("industrial", templates["industrial"].Category)
	s.Equal("desktop", templates["kiosk"].Category)
}

func (s *TemplatesTestSuite) TestTemplatePackages() {
	tm := NewTemplateManager()
	templates := tm.ListTemplates()

	// Check that templates have appropriate packages
	minimal := templates["minimal"]
	s.Empty(minimal.Config.Packages)

	networking := templates["networking"]
	s.Contains(networking.Config.Packages, "openssh")
	s.Contains(networking.Config.Packages, "wpa_supplicant")

	iot := templates["iot"]
	s.Contains(iot.Config.Packages, "mosquitto")
	s.Contains(iot.Config.Packages, "python3")

	security := templates["security"]
	s.Contains(security.Config.Packages, "openssh")
	s.Contains(security.Config.Packages, "openvpn")

	industrial := templates["industrial"]
	s.Contains(industrial.Config.Packages, "modbus")

	kiosk := templates["kiosk"]
	s.Contains(kiosk.Config.Packages, "xorg-server")
	s.Contains(kiosk.Config.Packages, "chromium")
}

func (s *TemplatesTestSuite) TestTemplateKernelConfigs() {
	tm := NewTemplateManager()
	templates := tm.ListTemplates()

	// Check kernel configurations
	iot := templates["iot"]
	s.Contains(iot.Config.Kernel.Config, "I2C")
	s.Contains(iot.Config.Kernel.Config, "SPI")
	s.Contains(iot.Config.Kernel.Config, "GPIO")

	industrial := templates["industrial"]
	s.Contains(industrial.Config.Kernel.Config, "PREEMPT_RT")
	s.Contains(industrial.Config.Kernel.Config, "HIGH_RES_TIMERS")
}

func (s *TemplatesTestSuite) TestTemplateFeatures() {
	tm := NewTemplateManager()
	templates := tm.ListTemplates()

	// Check features
	networking := templates["networking"]
	s.Contains(networking.Config.Features, "network")

	security := templates["security"]
	s.Contains(security.Config.Features, "systemd")

	industrial := templates["industrial"]
	s.Contains(industrial.Config.Features, "systemd")

	kiosk := templates["kiosk"]
	s.Contains(kiosk.Config.Features, "systemd")
}
