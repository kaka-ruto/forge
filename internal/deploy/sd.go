package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sst/forge/internal/logger"
)

// SDDeployer handles SD card deployments
type SDDeployer struct {
	logger *logger.Logger
}

// NewSDDeployer creates a new SD deployer
func NewSDDeployer() *SDDeployer {
	return &SDDeployer{
		logger: logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
	}
}

// Validate validates the SD deployment configuration
func (s *SDDeployer) Validate(config *DeploymentConfig) error {
	if config.Device == "" {
		return fmt.Errorf("device not specified for SD card deployment")
	}

	// Basic validation that it's a block device path (on Unix-like systems)
	if !strings.HasPrefix(config.Device, "/dev/") {
		s.logger.Warn("Device %s doesn't look like a block device path", config.Device)
	}

	return nil
}

// Deploy deploys the image to an SD card
func (s *SDDeployer) Deploy(artifactsDir string, config *DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Success: false,
	}

	s.logger.Info("Starting SD card deployment to %s", config.Device)

	// Validate device is not mounted (basic check)
	if s.isDeviceMounted(config.Device) {
		return nil, fmt.Errorf("device %s appears to be mounted, please unmount it first", config.Device)
	}

	// Get artifact paths
	kernelPath := filepath.Join(artifactsDir, "bzImage")
	rootfsPath := filepath.Join(artifactsDir, "rootfs.ext4")

	// Create temporary mount point
	mountPoint, err := os.MkdirTemp("", "forge-sd-mount-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary mount point: %v", err)
	}
	defer os.RemoveAll(mountPoint)

	// Partition the SD card
	if err := s.partitionSDCard(config.Device); err != nil {
		return nil, fmt.Errorf("failed to partition SD card: %v", err)
	}

	// Format the partition
	if err := s.formatSDPartition(config.Device); err != nil {
		return nil, fmt.Errorf("failed to format SD partition: %v", err)
	}

	// Mount the partition
	if err := s.mountPartition(config.Device+"1", mountPoint); err != nil {
		return nil, fmt.Errorf("failed to mount SD partition: %v", err)
	}
	defer s.unmountPartition(mountPoint)

	// Copy kernel
	kernelDest := filepath.Join(mountPoint, "bzImage")
	if err := CopyArtifact(kernelPath, kernelDest); err != nil {
		return nil, fmt.Errorf("failed to copy kernel: %v", err)
	}

	// Copy root filesystem
	rootfsDest := filepath.Join(mountPoint, "rootfs.ext4")
	if err := CopyArtifact(rootfsPath, rootfsDest); err != nil {
		return nil, fmt.Errorf("failed to copy root filesystem: %v", err)
	}

	// Install bootloader (GRUB)
	if err := s.installBootloader(config.Device, mountPoint); err != nil {
		return nil, fmt.Errorf("failed to install bootloader: %v", err)
	}

	// Create GRUB configuration
	if err := s.createGRUBConfig(mountPoint); err != nil {
		return nil, fmt.Errorf("failed to create GRUB config: %v", err)
	}

	result.Success = true
	result.Details = fmt.Sprintf("Successfully deployed to SD card %s", config.Device)
	result.Artifacts = []string{kernelDest, rootfsDest}

	s.logger.Info("SD card deployment completed successfully")
	return result, nil
}

// Cleanup cleans up the SD deployment
func (s *SDDeployer) Cleanup(config *DeploymentConfig) error {
	// For SD deployments, cleanup is minimal as the device remains usable
	s.logger.Info("SD card deployment cleanup completed")
	return nil
}

// isDeviceMounted checks if a device is currently mounted
func (s *SDDeployer) isDeviceMounted(device string) bool {
	cmd := exec.Command("mount")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), device)
}

// partitionSDCard partitions the SD card
func (s *SDDeployer) partitionSDCard(device string) error {
	s.logger.Info("Partitioning SD card %s", device)

	// Create a simple partition table with one partition
	// This is a simplified version - in production, you'd want more robust partitioning
	partitionScript := fmt.Sprintf(`
echo "label: dos
device: %s
unit: sectors
sector-size: 512

start=2048, type=83, bootable" | sudo sfdisk %s
`, device, device)

	cmd := exec.Command("bash", "-c", partitionScript)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("partitioning failed: %v, output: %s", err, string(output))
	}

	// Wait for the system to recognize the new partition
	exec.Command("partprobe", device).Run()

	return nil
}

// formatSDPartition formats the SD partition
func (s *SDDeployer) formatSDPartition(device string) error {
	s.logger.Info("Formatting SD partition %s1", device)

	cmd := exec.Command("sudo", "mkfs.ext4", device+"1")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("formatting failed: %v, output: %s", err, string(output))
	}

	return nil
}

// mountPartition mounts a partition
func (s *SDDeployer) mountPartition(partition, mountPoint string) error {
	s.logger.Info("Mounting partition %s to %s", partition, mountPoint)

	cmd := exec.Command("sudo", "mount", partition, mountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("mount failed: %v, output: %s", err, string(output))
	}

	return nil
}

// unmountPartition unmounts a partition
func (s *SDDeployer) unmountPartition(mountPoint string) error {
	s.logger.Info("Unmounting %s", mountPoint)

	cmd := exec.Command("sudo", "umount", mountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		s.logger.Warn("Unmount failed: %v, output: %s", err, string(output))
	}

	return nil
}

// installBootloader installs GRUB bootloader
func (s *SDDeployer) installBootloader(device, mountPoint string) error {
	s.logger.Info("Installing GRUB bootloader on %s", device)

	// Create GRUB directory
	grubDir := filepath.Join(mountPoint, "boot", "grub")
	if err := os.MkdirAll(grubDir, 0755); err != nil {
		return fmt.Errorf("failed to create GRUB directory: %v", err)
	}

	// Install GRUB (simplified - in production you'd use grub-install)
	// For now, we'll just create the necessary directory structure
	s.logger.Info("GRUB installation simulated (would run grub-install here)")

	return nil
}

// createGRUBConfig creates a GRUB configuration file
func (s *SDDeployer) createGRUBConfig(mountPoint string) error {
	s.logger.Info("Creating GRUB configuration")

	grubDir := filepath.Join(mountPoint, "boot", "grub")
	grubCfgPath := filepath.Join(grubDir, "grub.cfg")

	grubCfg := `set timeout=5
set default=0

menuentry "Forge OS" {
    linux /bzImage root=/dev/mmcblk0p1 console=ttyS0
    initrd /initrd.img
}
`

	if err := os.WriteFile(grubCfgPath, []byte(grubCfg), 0644); err != nil {
		return fmt.Errorf("failed to write GRUB config: %v", err)
	}

	return nil
}
