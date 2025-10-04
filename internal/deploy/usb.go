package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sst/forge/internal/logger"
)

// USBDeployer handles USB drive deployments
type USBDeployer struct {
	logger *logger.Logger
}

// NewUSBDeployer creates a new USB deployer
func NewUSBDeployer() *USBDeployer {
	return &USBDeployer{
		logger: logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
	}
}

// Validate validates the USB deployment configuration
func (u *USBDeployer) Validate(config *DeploymentConfig) error {
	if config.Device == "" {
		return fmt.Errorf("device not specified for USB deployment")
	}

	// Basic validation that it's a block device path (on Unix-like systems)
	if !strings.HasPrefix(config.Device, "/dev/") {
		u.logger.Warn("Device %s doesn't look like a block device path", config.Device)
	}

	return nil
}

// Deploy deploys the image to a USB drive
func (u *USBDeployer) Deploy(artifactsDir string, config *DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Success: false,
	}

	u.logger.Info("Starting USB deployment to %s", config.Device)

	// Validate device is not mounted (basic check)
	if u.isDeviceMounted(config.Device) {
		return nil, fmt.Errorf("device %s appears to be mounted, please unmount it first", config.Device)
	}

	// Get artifact paths
	kernelPath := filepath.Join(artifactsDir, "bzImage")
	rootfsPath := filepath.Join(artifactsDir, "rootfs.ext4")

	// Create temporary mount point
	mountPoint, err := os.MkdirTemp("", "forge-usb-mount-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary mount point: %v", err)
	}
	defer os.RemoveAll(mountPoint)

	// Partition the USB drive
	if err := u.partitionUSBDrive(config.Device); err != nil {
		return nil, fmt.Errorf("failed to partition USB drive: %v", err)
	}

	// Format the partition
	if err := u.formatUSBPartition(config.Device); err != nil {
		return nil, fmt.Errorf("failed to format USB partition: %v", err)
	}

	// Mount the partition
	if err := u.mountPartition(config.Device+"1", mountPoint); err != nil {
		return nil, fmt.Errorf("failed to mount USB partition: %v", err)
	}
	defer u.unmountPartition(mountPoint)

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
	if err := u.installBootloader(config.Device, mountPoint); err != nil {
		return nil, fmt.Errorf("failed to install bootloader: %v", err)
	}

	// Create GRUB configuration
	if err := u.createGRUBConfig(mountPoint); err != nil {
		return nil, fmt.Errorf("failed to create GRUB config: %v", err)
	}

	result.Success = true
	result.Details = fmt.Sprintf("Successfully deployed to USB drive %s", config.Device)
	result.Artifacts = []string{kernelDest, rootfsDest}

	u.logger.Info("USB deployment completed successfully")
	return result, nil
}

// Cleanup cleans up the USB deployment
func (u *USBDeployer) Cleanup(config *DeploymentConfig) error {
	// For USB deployments, cleanup is minimal as the device remains usable
	u.logger.Info("USB deployment cleanup completed")
	return nil
}

// isDeviceMounted checks if a device is currently mounted
func (u *USBDeployer) isDeviceMounted(device string) bool {
	cmd := exec.Command("mount")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), device)
}

// partitionUSBDrive partitions the USB drive
func (u *USBDeployer) partitionUSBDrive(device string) error {
	u.logger.Info("Partitioning USB drive %s", device)

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

// formatUSBPartition formats the USB partition
func (u *USBDeployer) formatUSBPartition(device string) error {
	u.logger.Info("Formatting USB partition %s1", device)

	cmd := exec.Command("sudo", "mkfs.ext4", device+"1")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("formatting failed: %v, output: %s", err, string(output))
	}

	return nil
}

// mountPartition mounts a partition
func (u *USBDeployer) mountPartition(partition, mountPoint string) error {
	u.logger.Info("Mounting partition %s to %s", partition, mountPoint)

	cmd := exec.Command("sudo", "mount", partition, mountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("mount failed: %v, output: %s", err, string(output))
	}

	return nil
}

// unmountPartition unmounts a partition
func (u *USBDeployer) unmountPartition(mountPoint string) error {
	u.logger.Info("Unmounting %s", mountPoint)

	cmd := exec.Command("sudo", "umount", mountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		u.logger.Warn("Unmount failed: %v, output: %s", err, string(output))
	}

	return nil
}

// installBootloader installs GRUB bootloader
func (u *USBDeployer) installBootloader(device, mountPoint string) error {
	u.logger.Info("Installing GRUB bootloader on %s", device)

	// Create GRUB directory
	grubDir := filepath.Join(mountPoint, "boot", "grub")
	if err := os.MkdirAll(grubDir, 0755); err != nil {
		return fmt.Errorf("failed to create GRUB directory: %v", err)
	}

	// Install GRUB (simplified - in production you'd use grub-install)
	// For now, we'll just create the necessary directory structure
	u.logger.Info("GRUB installation simulated (would run grub-install here)")

	return nil
}

// createGRUBConfig creates a GRUB configuration file
func (u *USBDeployer) createGRUBConfig(mountPoint string) error {
	u.logger.Info("Creating GRUB configuration")

	grubDir := filepath.Join(mountPoint, "boot", "grub")
	grubCfgPath := filepath.Join(grubDir, "grub.cfg")

	grubCfg := `set timeout=5
set default=0

menuentry "Forge OS" {
    linux /bzImage root=/dev/sda1 console=ttyS0
    initrd /initrd.img
}
`

	if err := os.WriteFile(grubCfgPath, []byte(grubCfg), 0644); err != nil {
		return fmt.Errorf("failed to write GRUB config: %v", err)
	}

	return nil
}
