package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/deploy"
)

// NewDeployCommand creates the deploy command
func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [target]",
		Short: "Deploy the Forge OS image",
		Long:  `Deploy the built Forge OS image to various targets (USB, SD card, remote, etc.).`,
		RunE:  runDeployCommandE,
	}

	cmd.Flags().String("device", "", "Target device (e.g., /dev/sdb)")
	cmd.Flags().String("host", "", "Remote host for deployment")
	cmd.Flags().String("user", "root", "SSH user for remote deployment")
	cmd.Flags().Int("port", 22, "SSH port for remote deployment")
	cmd.Flags().String("key", "", "SSH key path for remote deployment")
	cmd.Flags().Bool("dry-run", false, "Show what would be done without actually deploying")

	return cmd
}

func runDeployCommandE(cmd *cobra.Command, args []string) error {
	device, _ := cmd.Flags().GetString("device")
	host, _ := cmd.Flags().GetString("host")
	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetInt("port")
	key, _ := cmd.Flags().GetString("key")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	flags := map[string]string{
		"device": device,
		"host":   host,
		"user":   user,
		"port":   strconv.Itoa(port),
		"key":    key,
	}

	if dryRun {
		flags["dry-run"] = "true"
	}

	return runDeployCommand(args, flags)
}

// runDeployCommand executes the deploy logic
func runDeployCommand(args []string, flags map[string]string) error {
	if len(args) == 0 {
		return fmt.Errorf("deployment target not specified (usb, sd, remote)")
	}

	target := args[0]

	// Check if we're in a Forge project directory
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no forge.yml found - not in a Forge project directory")
	}

	// Load and validate configuration
	cfg, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Check for build artifacts
	buildDir := "build"
	artifactsDir := filepath.Join(buildDir, "artifacts", "images")
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		return fmt.Errorf("no build artifacts found - run 'forge build' first")
	}

	// Create deployment orchestrator
	orchestrator := deploy.NewDeploymentOrchestrator(cfg)

	// Register deployers
	orchestrator.RegisterDeployer(deploy.TargetUSB, deploy.NewUSBDeployer())
	orchestrator.RegisterDeployer(deploy.TargetRemote, deploy.NewRemoteDeployer())

	// Parse deployment target
	var deployTarget deploy.DeploymentTarget
	switch target {
	case "usb":
		deployTarget = deploy.TargetUSB
	case "sd":
		deployTarget = deploy.TargetSDCard
	case "remote":
		deployTarget = deploy.TargetRemote
	default:
		return fmt.Errorf("unsupported deployment target: %s", target)
	}

	// Create deployment config
	deployConfig := &deploy.DeploymentConfig{
		Target:       deployTarget,
		Device:       flags["device"],
		Host:         flags["host"],
		User:         flags["user"],
		KeyPath:      flags["key"],
		DryRun:       flags["dry-run"] == "true",
		ValidateOnly: false,
	}

	// Parse port if provided
	if portStr, ok := flags["port"]; ok && portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			deployConfig.Port = port
		}
	}

	// Validate deployment configuration
	if err := orchestrator.ValidateDeployment(deployConfig); err != nil {
		return fmt.Errorf("invalid deployment configuration: %v", err)
	}

	// Execute deployment
	result, err := orchestrator.ExecuteDeployment(artifactsDir, deployConfig)
	if err != nil {
		return fmt.Errorf("deployment failed: %v", err)
	}

	// Display results
	if result.Success {
		fmt.Printf("✓ Deployment completed successfully\n")
		if result.RemoteInfo != nil && result.RemoteInfo.AccessURL != "" {
			fmt.Printf("Access URL: %s\n", result.RemoteInfo.AccessURL)
		}
	} else {
		fmt.Printf("✗ Deployment failed: %s\n", result.Error)
	}

	return nil
}
