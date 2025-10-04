package buildroot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sst/forge/internal/config"
)

// BuildrootManager manages Buildroot operations
type BuildrootManager struct {
	config     *config.Config
	projectDir string
	buildDir   string
}

// NewBuildrootManager creates a new Buildroot manager
func NewBuildrootManager(cfg *config.Config, projectDir string) *BuildrootManager {
	return &BuildrootManager{
		config:     cfg,
		projectDir: projectDir,
		buildDir:   filepath.Join(projectDir, "build"),
	}
}

// DownloadBuildroot downloads and extracts Buildroot
func (bm *BuildrootManager) DownloadBuildroot() error {
	version := bm.config.Buildroot.Version
	if version == "" {
		version = "stable"
	}

	// Create build directory
	if err := os.MkdirAll(bm.buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %v", err)
	}

	// Check if Buildroot is already downloaded
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	if _, err := os.Stat(buildrootDir); err == nil {
		// Buildroot already exists, skip download
		return nil
	}

	// Determine download URL
	var downloadURL string
	if version == "stable" {
		downloadURL = "https://buildroot.org/downloads/buildroot-latest.tar.gz"
	} else {
		downloadURL = fmt.Sprintf("https://buildroot.org/downloads/buildroot-%s.tar.gz", version)
	}

	// Download Buildroot
	tarPath := filepath.Join(bm.buildDir, "buildroot.tar.gz")
	if err := bm.downloadFile(downloadURL, tarPath); err != nil {
		return fmt.Errorf("failed to download Buildroot: %v", err)
	}

	// Extract Buildroot
	if err := bm.extractTarGz(tarPath, bm.buildDir); err != nil {
		return fmt.Errorf("failed to extract Buildroot: %v", err)
	}

	// Rename extracted directory to buildroot
	extractedDir := bm.findExtractedBuildrootDir()
	if extractedDir == "" {
		return fmt.Errorf("could not find extracted Buildroot directory")
	}

	if err := os.Rename(filepath.Join(bm.buildDir, extractedDir), buildrootDir); err != nil {
		return fmt.Errorf("failed to rename Buildroot directory: %v", err)
	}

	// Clean up tar file
	os.Remove(tarPath)

	return nil
}

// GenerateConfig generates Buildroot .config file from forge configuration
func (bm *BuildrootManager) GenerateConfig() error {
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")

	// Start with default configuration
	if err := bm.runMake(buildrootDir, "defconfig"); err != nil {
		return fmt.Errorf("failed to generate default config: %v", err)
	}

	// Apply architecture-specific configuration
	if err := bm.applyArchitectureConfig(); err != nil {
		return fmt.Errorf("failed to apply architecture config: %v", err)
	}

	// Apply package configuration
	if err := bm.applyPackageConfig(); err != nil {
		return fmt.Errorf("failed to apply package config: %v", err)
	}

	// Apply kernel configuration
	if err := bm.applyKernelConfig(); err != nil {
		return fmt.Errorf("failed to apply kernel config: %v", err)
	}

	// Apply feature configuration
	if err := bm.applyFeatureConfig(); err != nil {
		return fmt.Errorf("failed to apply feature config: %v", err)
	}

	return nil
}

// Build executes the Buildroot build process
func (bm *BuildrootManager) Build() error {
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")

	// Run make with parallel jobs
	if err := bm.runMake(buildrootDir, fmt.Sprintf("-j%d", bm.getParallelJobs())); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	return nil
}

// GetOutputDir returns the Buildroot output directory
func (bm *BuildrootManager) GetOutputDir() string {
	return filepath.Join(bm.buildDir, "buildroot", "output")
}

// GetImagesDir returns the Buildroot images directory
func (bm *BuildrootManager) GetImagesDir() string {
	return filepath.Join(bm.GetOutputDir(), "images")
}

// downloadFile downloads a file from URL to local path
func (bm *BuildrootManager) downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractTarGz extracts a tar.gz file to destination directory
func (bm *BuildrootManager) extractTarGz(tarPath, destDir string) error {
	cmd := exec.Command("tar", "-xzf", tarPath, "-C", destDir)
	return cmd.Run()
}

// findExtractedBuildrootDir finds the extracted Buildroot directory name
func (bm *BuildrootManager) findExtractedBuildrootDir() string {
	entries, err := os.ReadDir(bm.buildDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "buildroot") {
			return entry.Name()
		}
	}

	return ""
}

// runMake executes make command in the specified directory
func (bm *BuildrootManager) runMake(dir string, args ...string) error {
	cmd := exec.Command("make", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// getParallelJobs returns the number of parallel jobs to use for building
func (bm *BuildrootManager) getParallelJobs() int {
	// Use number of CPU cores
	return 4 // TODO: Detect actual CPU count
}

// applyArchitectureConfig applies architecture-specific Buildroot configuration
func (bm *BuildrootManager) applyArchitectureConfig() error {
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	configPath := filepath.Join(buildrootDir, ".config")

	arch := bm.config.Architecture
	var configLines []string

	switch arch {
	case "x86_64":
		configLines = []string{
			"BR2_x86_64=y",
			"BR2_ARCH=\"x86_64\"",
		}
	case "arm":
		configLines = []string{
			"BR2_arm=y",
			"BR2_ARCH=\"arm\"",
		}
	case "aarch64":
		configLines = []string{
			"BR2_aarch64=y",
			"BR2_ARCH=\"aarch64\"",
		}
	case "mips":
		configLines = []string{
			"BR2_mips=y",
			"BR2_ARCH=\"mips\"",
		}
	default:
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	return bm.appendConfigLines(configPath, configLines)
}

// applyPackageConfig applies package configuration to Buildroot
func (bm *BuildrootManager) applyPackageConfig() error {
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	configPath := filepath.Join(buildrootDir, ".config")

	var configLines []string

	for _, pkg := range bm.config.Packages {
		switch pkg {
		case "openssh":
			configLines = append(configLines, "BR2_PACKAGE_OPENSSH=y")
		case "wpa_supplicant":
			configLines = append(configLines, "BR2_PACKAGE_WPA_SUPPLICANT=y")
		case "dhcpcd":
			configLines = append(configLines, "BR2_PACKAGE_DHCPCD=y")
		case "mosquitto":
			configLines = append(configLines, "BR2_PACKAGE_MOSQUITTO=y")
		case "python3":
			configLines = append(configLines, "BR2_PACKAGE_PYTHON3=y")
		case "i2c-tools":
			configLines = append(configLines, "BR2_PACKAGE_I2C_TOOLS=y")
		case "openvpn":
			configLines = append(configLines, "BR2_PACKAGE_OPENVPN=y")
		case "iptables":
			configLines = append(configLines, "BR2_PACKAGE_IPTABLES=y")
		case "fail2ban":
			configLines = append(configLines, "BR2_PACKAGE_FAIL2BAN=y")
		case "modbus":
			configLines = append(configLines, "BR2_PACKAGE_LIBMODBUS=y")
		case "chrony":
			configLines = append(configLines, "BR2_PACKAGE_CHRONY=y")
		case "rsyslog":
			configLines = append(configLines, "BR2_PACKAGE_RSYSLOG=y")
		case "xorg-server":
			configLines = append(configLines, "BR2_PACKAGE_XORG7=y")
		case "chromium":
			configLines = append(configLines, "BR2_PACKAGE_CHROMIUM=y")
		case "xterm":
			configLines = append(configLines, "BR2_PACKAGE_XTERM=y")
		case "fluxbox":
			configLines = append(configLines, "BR2_PACKAGE_FLUXBOX=y")
			// Add more package mappings as needed
		}
	}

	return bm.appendConfigLines(configPath, configLines)
}

// applyKernelConfig applies kernel configuration to Buildroot
func (bm *BuildrootManager) applyKernelConfig() error {
	// Kernel configuration is handled separately in the kernel config file
	// This method can be extended to set kernel-related Buildroot options
	return nil
}

// applyFeatureConfig applies feature configuration to Buildroot
func (bm *BuildrootManager) applyFeatureConfig() error {
	buildrootDir := filepath.Join(bm.buildDir, "buildroot")
	configPath := filepath.Join(buildrootDir, ".config")

	var configLines []string

	for _, feature := range bm.config.Features {
		switch feature {
		case "systemd":
			configLines = append(configLines, "BR2_INIT_SYSTEMD=y")
		case "network":
			// Network features are enabled by default, but we can add specific configs
			configLines = append(configLines, "BR2_SYSTEM_ENABLE_NLS=y")
			// Add more feature mappings as needed
		}
	}

	return bm.appendConfigLines(configPath, configLines)
}

// appendConfigLines appends configuration lines to the Buildroot .config file
func (bm *BuildrootManager) appendConfigLines(configPath string, lines []string) error {
	if len(lines) == 0 {
		return nil
	}

	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}
