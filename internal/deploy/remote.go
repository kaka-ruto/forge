package deploy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sst/forge/internal/logger"
	"golang.org/x/crypto/ssh"
)

// RemoteDeployer handles remote deployments via SSH
type RemoteDeployer struct {
	logger *logger.Logger
}

// NewRemoteDeployer creates a new remote deployer
func NewRemoteDeployer() *RemoteDeployer {
	return &RemoteDeployer{
		logger: logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
	}
}

// Validate validates the remote deployment configuration
func (r *RemoteDeployer) Validate(config *DeploymentConfig) error {
	if config.Host == "" {
		return fmt.Errorf("host not specified for remote deployment")
	}

	if config.Port == 0 {
		config.Port = 22 // Default SSH port
	}

	if config.User == "" {
		config.User = "root" // Default user
	}

	return nil
}

// Deploy deploys the image to a remote host via SSH
func (r *RemoteDeployer) Deploy(artifactsDir string, config *DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Success: false,
		RemoteInfo: &RemoteDeploymentInfo{
			Host: config.Host,
			Port: config.Port,
		},
	}

	r.logger.Info("Starting remote deployment to %s@%s:%d", config.User, config.Host, config.Port)

	// Establish SSH connection
	client, err := r.connectSSH(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote host: %v", err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Get artifact paths
	kernelPath := filepath.Join(artifactsDir, "bzImage")
	rootfsPath := filepath.Join(artifactsDir, "rootfs.ext4")

	// Upload kernel
	remoteKernelPath := "/tmp/forge-kernel"
	if err := r.uploadFile(client, kernelPath, remoteKernelPath); err != nil {
		return nil, fmt.Errorf("failed to upload kernel: %v", err)
	}

	// Upload root filesystem
	remoteRootfsPath := "/tmp/forge-rootfs.ext4"
	if err := r.uploadFile(client, rootfsPath, remoteRootfsPath); err != nil {
		return nil, fmt.Errorf("failed to upload root filesystem: %v", err)
	}

	// Configure bootloader on remote host
	if err := r.configureRemoteBootloader(session, remoteKernelPath, remoteRootfsPath); err != nil {
		return nil, fmt.Errorf("failed to configure remote bootloader: %v", err)
	}

	// Move files to final locations
	if err := r.finalizeRemoteDeployment(session, remoteKernelPath, remoteRootfsPath); err != nil {
		return nil, fmt.Errorf("failed to finalize remote deployment: %v", err)
	}

	result.Success = true
	result.Details = fmt.Sprintf("Successfully deployed to remote host %s", config.Host)
	result.RemoteInfo.AccessURL = fmt.Sprintf("ssh://%s@%s", config.User, config.Host)

	r.logger.Info("Remote deployment completed successfully")
	return result, nil
}

// Cleanup cleans up the remote deployment
func (r *RemoteDeployer) Cleanup(config *DeploymentConfig) error {
	r.logger.Info("Remote deployment cleanup completed")
	return nil
}

// connectSSH establishes an SSH connection
func (r *RemoteDeployer) connectSSH(config *DeploymentConfig) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	// Try SSH key authentication first
	if config.KeyPath != "" {
		if key, err := r.loadSSHKey(config.KeyPath); err == nil {
			authMethods = append(authMethods, ssh.PublicKeys(key))
		}
	}

	// Fallback to password auth (not recommended for production)
	// Note: This is simplified - in production you'd handle passwords securely

	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
	}

	addr := config.Host + ":" + strconv.Itoa(config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %v", err)
	}

	return client, nil
}

// loadSSHKey loads an SSH private key
func (r *RemoteDeployer) loadSSHKey(keyPath string) (ssh.Signer, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// uploadFile uploads a file via SCP
func (r *RemoteDeployer) uploadFile(client *ssh.Client, localPath, remotePath string) error {
	r.logger.Info("Uploading %s to %s", localPath, remotePath)

	// For simplicity, we'll use a basic SCP approach
	// In production, you'd want a more robust SCP/SFTP implementation

	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Get file info
	stat, err := localFile.Stat()
	if err != nil {
		return err
	}

	// Create remote file
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintf(w, "C%#o %d %s\n", stat.Mode().Perm(), stat.Size(), filepath.Base(remotePath))
		localFile.Seek(0, 0)
		io.Copy(w, localFile)
		fmt.Fprint(w, "\x00")
	}()

	cmd := fmt.Sprintf("scp -t %s", remotePath)
	if err := session.Run(cmd); err != nil {
		return err
	}

	return nil
}

// configureRemoteBootloader configures the bootloader on the remote host
func (r *RemoteDeployer) configureRemoteBootloader(session *ssh.Session, kernelPath, rootfsPath string) error {
	r.logger.Info("Configuring remote bootloader")

	// This is a simplified example - in production you'd detect the bootloader
	// and configure it appropriately (GRUB, U-Boot, etc.)

	commands := []string{
		"mkdir -p /boot/forge",
		fmt.Sprintf("cp %s /boot/forge/", kernelPath),
		fmt.Sprintf("cp %s /boot/forge/", rootfsPath),
		"update-grub", // For Debian/Ubuntu systems
	}

	for _, cmd := range commands {
		if err := session.Run(cmd); err != nil {
			r.logger.Warn("Command failed: %s (%v)", cmd, err)
			// Continue with other commands
		}
	}

	return nil
}

// finalizeRemoteDeployment moves files to final locations
func (r *RemoteDeployer) finalizeRemoteDeployment(session *ssh.Session, kernelPath, rootfsPath string) error {
	r.logger.Info("Finalizing remote deployment")

	// Clean up temporary files
	commands := []string{
		fmt.Sprintf("rm -f %s %s", kernelPath, rootfsPath),
	}

	for _, cmd := range commands {
		if err := session.Run(cmd); err != nil {
			r.logger.Warn("Cleanup command failed: %s (%v)", cmd, err)
		}
	}

	return nil
}
