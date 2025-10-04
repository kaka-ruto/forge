package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sst/forge/internal/config"
)

// Template represents a Forge OS project template
type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Config      *config.Config    `json:"config"`
	Files       map[string]string `json:"files"`      // filename -> content
	Overlays    map[string]string `json:"overlays"`   // path -> content
	PostHooks   []string          `json:"post_hooks"` // commands to run after creation
}

// TemplateManager manages built-in templates
type TemplateManager struct {
	templates map[string]*Template
}

// NewTemplateManager creates a new template manager with built-in templates
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*Template),
	}

	// Register built-in templates
	tm.registerMinimalTemplate()
	tm.registerNetworkingTemplate()
	tm.registerIOTTemplate()
	tm.registerSecurityTemplate()
	tm.registerIndustrialTemplate()
	tm.registerKioskTemplate()

	return tm
}

// GetTemplate returns a template by name
func (tm *TemplateManager) GetTemplate(name string) (*Template, error) {
	template, exists := tm.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found", name)
	}
	return template, nil
}

// ListTemplates returns all available templates
func (tm *TemplateManager) ListTemplates() map[string]*Template {
	return tm.templates
}

// ApplyTemplate applies a template to a project directory
func (tm *TemplateManager) ApplyTemplate(templateName, projectDir string, data map[string]interface{}) error {
	template, err := tm.GetTemplate(templateName)
	if err != nil {
		return err
	}

	// Set default template data
	if data == nil {
		data = make(map[string]interface{})
	}
	projectName := filepath.Base(projectDir)
	if _, exists := data["ProjectName"]; !exists {
		data["ProjectName"] = projectName
	}
	if _, exists := data["Architecture"]; !exists {
		data["Architecture"] = "x86_64" // Default architecture
	}

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Generate forge.yml from template config with substitution
	if err := tm.generateConfigFile(template.Config, filepath.Join(projectDir, "forge.yml"), data); err != nil {
		return fmt.Errorf("failed to save forge.yml: %v", err)
	}

	// Create template files with substitution
	for filename, content := range template.Files {
		if err := tm.generateTemplateFile(filename, content, projectDir, data); err != nil {
			return err
		}
	}

	// Create overlay files with substitution
	for overlayPath, content := range template.Overlays {
		if err := tm.generateTemplateFile(overlayPath, content, projectDir, data); err != nil {
			return err
		}
	}

	// Execute post-hooks (for future use)
	for _, hook := range template.PostHooks {
		// TODO: Execute post-hooks
		_ = hook
	}

	return nil
}

// generateConfigFile generates a config file with template substitution
func (tm *TemplateManager) generateConfigFile(cfg *config.Config, filePath string, data map[string]interface{}) error {
	// Create a copy of the config for modification
	configCopy := *cfg

	// Substitute template variables in config
	if name, ok := data["ProjectName"].(string); ok {
		configCopy.Name = strings.ReplaceAll(configCopy.Name, "{{.ProjectName}}", name)
	}
	if arch, ok := data["Architecture"].(string); ok {
		configCopy.Architecture = strings.ReplaceAll(configCopy.Architecture, "{{.Architecture}}", arch)
	}

	// Validate the config before saving
	if err := configCopy.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return config.SaveConfig(&configCopy, filePath)
}

// generateTemplateFile generates a template file with substitution
func (tm *TemplateManager) generateTemplateFile(filename, content, projectDir string, data map[string]interface{}) error {
	filePath := filepath.Join(projectDir, filename)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %v", filename, err)
	}

	// Parse and execute template
	tmpl, err := template.New(filename).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", filename, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", filename, err)
	}

	return nil
}

// ValidateTemplate validates a template configuration
func (tm *TemplateManager) ValidateTemplate(template *Template) error {
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}
	if template.Config == nil {
		return fmt.Errorf("template config is required")
	}
	if err := template.Config.Validate(); err != nil {
		return fmt.Errorf("invalid template config: %v", err)
	}
	return nil
}

// registerMinimalTemplate registers the minimal template
func (tm *TemplateManager) registerMinimalTemplate() {
	template := &Template{
		Name:        "minimal",
		Description: "Minimal Linux system with BusyBox",
		Category:    "basic",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "{{.ProjectName}}",
			Version:       "0.1.0",
			Architecture:  "{{.Architecture}}",
			Template:      "minimal",
			Buildroot: config.BuildrootConfig{
				Version: "stable",
			},
			Kernel: config.KernelConfig{
				Version: "latest",
			},
			Packages: []string{},
			Features: []string{},
			Overlays: map[string]interface{}{},
		},
		Files: map[string]string{
			"README.md": `# {{.ProjectName}}

This is a minimal Forge OS project.

## Building

` + "```bash" + `
forge build
` + "```" + `

## Testing

` + "```bash" + `
forge test
` + "```" + `
`,
		},
		Overlays:  map[string]string{},
		PostHooks: []string{},
	}

	tm.templates["minimal"] = template
}

// registerNetworkingTemplate registers the networking template
func (tm *TemplateManager) registerNetworkingTemplate() {
	template := &Template{
		Name:        "networking",
		Description: "Network-enabled system with SSH and wireless support",
		Category:    "networking",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "{{.ProjectName}}",
			Version:       "0.1.0",
			Architecture:  "{{.Architecture}}",
			Template:      "networking",
			Buildroot: config.BuildrootConfig{
				Version: "stable",
			},
			Kernel: config.KernelConfig{
				Version: "latest",
			},
			Packages: []string{"openssh", "wpa_supplicant", "dhcpcd"},
			Features: []string{"network"},
			Overlays: map[string]interface{}{},
		},
		Files: map[string]string{
			"README.md": `# {{.ProjectName}}

This is a networking-enabled Forge OS project with SSH and wireless support.

## Building

` + "```bash" + `
forge build
` + "```" + `

## Testing

` + "```bash" + `
forge test
` + "```" + `

## Connecting

SSH is enabled by default. Connect with:
` + "```bash" + `
ssh root@<ip-address>
` + "```" + `
`,
			"overlays/rootfs/etc/network/interfaces": `auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
`,
		},
		Overlays:  map[string]string{},
		PostHooks: []string{},
	}

	tm.templates["networking"] = template
}

// registerIOTTemplate registers the IoT template
func (tm *TemplateManager) registerIOTTemplate() {
	template := &Template{
		Name:        "iot",
		Description: "IoT system with sensors, MQTT, and embedded features",
		Category:    "iot",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "{{.ProjectName}}",
			Version:       "0.1.0",
			Architecture:  "{{.Architecture}}",
			Template:      "iot",
			Buildroot: config.BuildrootConfig{
				Version: "stable",
			},
			Kernel: config.KernelConfig{
				Version: "latest",
				Config: map[string]string{
					"I2C":  "y",
					"SPI":  "y",
					"GPIO": "y",
				},
			},
			Packages: []string{"mosquitto", "python3", "i2c-tools"},
			Features: []string{},
			Overlays: map[string]interface{}{},
		},
		Files: map[string]string{
			"README.md": `# {{.ProjectName}}

This is an IoT Forge OS project with MQTT and sensor support.

## Features

- MQTT broker (Mosquitto)
- Python 3 runtime
- I2C/SPI/GPIO support
- Sensor interfaces

## Building

` + "```bash" + `
forge build
` + "```" + `

## Testing

` + "```bash" + `
forge test
` + "```" + `
`,
			"overlays/rootfs/etc/mosquitto/mosquitto.conf": `listener 1883
allow_anonymous true
`,
		},
		Overlays:  map[string]string{},
		PostHooks: []string{},
	}

	tm.templates["iot"] = template
}

// registerSecurityTemplate registers the security template
func (tm *TemplateManager) registerSecurityTemplate() {
	template := &Template{
		Name:        "security",
		Description: "Security-focused system with VPN and hardening",
		Category:    "security",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "{{.ProjectName}}",
			Version:       "0.1.0",
			Architecture:  "{{.Architecture}}",
			Template:      "security",
			Buildroot: config.BuildrootConfig{
				Version: "stable",
			},
			Kernel: config.KernelConfig{
				Version: "latest",
			},
			Packages: []string{"openssh", "openvpn", "iptables", "fail2ban"},
			Features: []string{"systemd"},
			Overlays: map[string]interface{}{},
		},
		Files: map[string]string{
			"README.md": `# {{.ProjectName}}

This is a security-focused Forge OS project with VPN and hardening features.

## Security Features

- OpenVPN support
- SSH hardening
- Firewall (iptables)
- Fail2Ban intrusion prevention
- systemd init system

## Building

` + "```bash" + `
forge build
` + "```" + `

## Testing

` + "```bash" + `
forge test
` + "```" + `
`,
		},
		Overlays:  map[string]string{},
		PostHooks: []string{},
	}

	tm.templates["security"] = template
}

// registerIndustrialTemplate registers the industrial template
func (tm *TemplateManager) registerIndustrialTemplate() {
	template := &Template{
		Name:        "industrial",
		Description: "Industrial control system with Modbus and real-time features",
		Category:    "industrial",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "{{.ProjectName}}",
			Version:       "0.1.0",
			Architecture:  "{{.Architecture}}",
			Template:      "industrial",
			Buildroot: config.BuildrootConfig{
				Version: "stable",
			},
			Kernel: config.KernelConfig{
				Version: "latest",
				Config: map[string]string{
					"PREEMPT_RT":      "y",
					"HIGH_RES_TIMERS": "y",
				},
			},
			Packages: []string{"modbus", "chrony", "rsyslog"},
			Features: []string{"systemd"},
			Overlays: map[string]interface{}{},
		},
		Files: map[string]string{
			"README.md": `# {{.ProjectName}}

This is an industrial control Forge OS project with real-time features.

## Industrial Features

- Modbus protocol support
- Real-time kernel (PREEMPT_RT)
- NTP synchronization (Chrony)
- System logging (rsyslog)
- systemd init system

## Building

` + "```bash" + `
forge build
` + "```" + `

## Testing

` + "```bash" + `
forge test
` + "```" + `
`,
		},
		Overlays:  map[string]string{},
		PostHooks: []string{},
	}

	tm.templates["industrial"] = template
}

// registerKioskTemplate registers the kiosk template
func (tm *TemplateManager) registerKioskTemplate() {
	template := &Template{
		Name:        "kiosk",
		Description: "Kiosk system with Chromium browser and X11",
		Category:    "desktop",
		Config: &config.Config{
			SchemaVersion: "1.0",
			Name:          "{{.ProjectName}}",
			Version:       "0.1.0",
			Architecture:  "{{.Architecture}}",
			Template:      "kiosk",
			Buildroot: config.BuildrootConfig{
				Version: "stable",
			},
			Kernel: config.KernelConfig{
				Version: "latest",
			},
			Packages: []string{"xorg-server", "chromium", "xterm", "fluxbox"},
			Features: []string{"systemd"},
			Overlays: map[string]interface{}{},
		},
		Files: map[string]string{
			"README.md": `# {{.ProjectName}}

This is a kiosk Forge OS project with web browser and desktop environment.

## Desktop Features

- X11 window system
- Chromium web browser
- Fluxbox window manager
- Kiosk-mode configuration

## Building

` + "```bash" + `
forge build
` + "```" + `

## Testing

` + "```bash" + `
forge test
` + "```" + `
`,
			"overlays/rootfs/etc/X11/xorg.conf": `Section "Device"
    Identifier "Card0"
    Driver "modesetting"
EndSection

Section "Screen"
    Identifier "Screen0"
    Device "Card0"
    DefaultDepth 24
EndSection
`,
		},
		Overlays:  map[string]string{},
		PostHooks: []string{},
	}

	tm.templates["kiosk"] = template
}
