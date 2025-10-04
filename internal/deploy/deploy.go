package deploy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
)

// DeploymentTarget represents different deployment targets
type DeploymentTarget string

const (
	TargetUSB    DeploymentTarget = "usb"
	TargetSDCard DeploymentTarget = "sd"
	TargetRemote DeploymentTarget = "remote"
	TargetCloud  DeploymentTarget = "cloud"
)

// DeploymentConfig holds configuration for deployment
type DeploymentConfig struct {
	Target       DeploymentTarget
	Device       string       // For USB/SD card deployments
	Host         string       // For remote deployments
	User         string       // SSH user for remote deployments
	Port         int          // SSH port for remote deployments
	KeyPath      string       // SSH key path
	CloudConfig  *CloudConfig // For cloud deployments
	ValidateOnly bool         // Only validate, don't deploy
	DryRun       bool         // Show what would be done
}

// CloudConfig holds cloud deployment configuration
type CloudConfig struct {
	Provider     string // aws, gcp, azure, etc.
	Region       string
	InstanceType string
	ImageName    string
	Tags         map[string]string
}

// DeploymentResult represents the result of a deployment operation
type DeploymentResult struct {
	Success    bool
	Error      string
	Details    string
	Artifacts  []string // Files/artifacts created
	RemoteInfo *RemoteDeploymentInfo
}

// RemoteDeploymentInfo contains information about remote deployments
type RemoteDeploymentInfo struct {
	Host       string
	Port       int
	ImageID    string // For cloud deployments
	InstanceID string // For cloud deployments
	PublicIP   string
	AccessURL  string
}

// Deployer interface for different deployment strategies
type Deployer interface {
	Validate(config *DeploymentConfig) error
	Deploy(artifactsDir string, config *DeploymentConfig) (*DeploymentResult, error)
	Cleanup(config *DeploymentConfig) error
}

// DeploymentOrchestrator manages the deployment process
type DeploymentOrchestrator struct {
	config    *config.Config
	logger    *logger.Logger
	deployers map[DeploymentTarget]Deployer
}

// NewDeploymentOrchestrator creates a new deployment orchestrator
func NewDeploymentOrchestrator(cfg *config.Config) *DeploymentOrchestrator {
	return &DeploymentOrchestrator{
		config:    cfg,
		logger:    logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
		deployers: make(map[DeploymentTarget]Deployer),
	}
}

// RegisterDeployer registers a deployer for a specific target
func (do *DeploymentOrchestrator) RegisterDeployer(target DeploymentTarget, deployer Deployer) {
	do.deployers[target] = deployer
}

// ValidateDeployment validates a deployment configuration
func (do *DeploymentOrchestrator) ValidateDeployment(deployConfig *DeploymentConfig) error {
	deployer, exists := do.deployers[deployConfig.Target]
	if !exists {
		return fmt.Errorf("no deployer registered for target: %s", deployConfig.Target)
	}

	return deployer.Validate(deployConfig)
}

// ExecuteDeployment executes a deployment
func (do *DeploymentOrchestrator) ExecuteDeployment(artifactsDir string, deployConfig *DeploymentConfig) (*DeploymentResult, error) {
	// Validate artifacts directory exists
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("artifacts directory does not exist: %s", artifactsDir)
	}

	// Check for required artifacts
	requiredArtifacts := []string{"bzImage", "rootfs.ext4"}
	for _, artifact := range requiredArtifacts {
		artifactPath := filepath.Join(artifactsDir, artifact)
		if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("required artifact not found: %s", artifactPath)
		}
	}

	deployer, exists := do.deployers[deployConfig.Target]
	if !exists {
		return nil, fmt.Errorf("no deployer registered for target: %s", deployConfig.Target)
	}

	if deployConfig.DryRun {
		do.logger.Info("DRY RUN: Would deploy to %s", deployConfig.Target)
		return &DeploymentResult{
			Success: true,
			Details: fmt.Sprintf("Dry run completed for target %s", deployConfig.Target),
		}, nil
	}

	do.logger.Info("Starting deployment to %s", deployConfig.Target)
	result, err := deployer.Deploy(artifactsDir, deployConfig)
	if err != nil {
		return nil, fmt.Errorf("deployment failed: %v", err)
	}

	if result.Success {
		do.logger.Info("Deployment completed successfully")
	} else {
		do.logger.Error("Deployment failed: %s", result.Error)
	}

	return result, nil
}

// CleanupDeployment cleans up a deployment
func (do *DeploymentOrchestrator) CleanupDeployment(deployConfig *DeploymentConfig) error {
	deployer, exists := do.deployers[deployConfig.Target]
	if !exists {
		return fmt.Errorf("no deployer registered for target: %s", deployConfig.Target)
	}

	return deployer.Cleanup(deployConfig)
}

// GetAvailableTargets returns all available deployment targets
func (do *DeploymentOrchestrator) GetAvailableTargets() []DeploymentTarget {
	var targets []DeploymentTarget
	for target := range do.deployers {
		targets = append(targets, target)
	}
	return targets
}

// ValidateArtifacts validates that required artifacts exist
func ValidateArtifacts(artifactsDir string) error {
	requiredArtifacts := []string{"bzImage", "rootfs.ext4"}

	for _, artifact := range requiredArtifacts {
		path := filepath.Join(artifactsDir, artifact)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("required artifact missing: %s", path)
		}
	}

	return nil
}

// CopyArtifact copies an artifact to a destination
func CopyArtifact(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// Ensure destination file is synced
	return destFile.Sync()
}
