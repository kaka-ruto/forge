package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sst/forge/internal/cicd"
)

// NewCICDCommand creates the cicd command
func NewCICDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cicd",
		Short: "Generate CI/CD pipelines",
		Long:  `Generate CI/CD pipeline configurations for various platforms (GitHub Actions, GitLab CI, etc.).`,
	}

	cmd.AddCommand(
		newCICDGenerateCommand(),
		newCICDListCommand(),
	)

	return cmd
}

func newCICDGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [provider] [type]",
		Short: "Generate CI/CD pipeline",
		Long:  `Generate a CI/CD pipeline configuration for the specified provider and pipeline type.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runCICDGenerateCommandE,
	}

	cmd.Flags().StringP("output", "o", ".forge/cicd", "Output directory for generated files")
	cmd.Flags().StringSliceP("triggers", "t", []string{"push", "pull_request"}, "Pipeline triggers")
	cmd.Flags().StringSliceP("branches", "b", []string{"main", "master", "develop"}, "Branches to run on")
	cmd.Flags().StringSliceP("tags", "T", []string{"v*"}, "Tags to run on")
	cmd.Flags().Bool("dry-run", false, "Show what would be generated without writing files")

	return cmd
}

func newCICDListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available CI/CD providers",
		Long:  `List all available CI/CD providers and their supported pipeline types.`,
		RunE:  runCICDListCommandE,
	}

	return cmd
}

func runCICDGenerateCommandE(cmd *cobra.Command, args []string) error {
	provider := args[0]
	pipelineType := args[1]

	// Check if we're in a Forge project directory
	if _, err := os.Stat("forge.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no forge.yml found - not in a Forge project directory")
	}

	// Load configuration
	cfg, err := loadForgeConfig("forge.yml")
	if err != nil {
		return fmt.Errorf("invalid forge.yml: %v", err)
	}

	// Parse provider
	var ciProvider cicd.CIProvider
	switch strings.ToLower(provider) {
	case "github", "github-actions", "gha":
		ciProvider = cicd.ProviderGitHubActions
	case "gitlab", "gitlab-ci":
		ciProvider = cicd.ProviderGitLabCI
	case "jenkins":
		ciProvider = cicd.ProviderJenkins
	case "circle", "circleci":
		ciProvider = cicd.ProviderCircleCI
	case "travis", "travisci":
		ciProvider = cicd.ProviderTravisCI
	default:
		return fmt.Errorf("unsupported CI provider: %s", provider)
	}

	// Parse pipeline type
	var ciPipelineType cicd.PipelineType
	switch strings.ToLower(pipelineType) {
	case "build":
		ciPipelineType = cicd.PipelineBuild
	case "test":
		ciPipelineType = cicd.PipelineTest
	case "deploy":
		ciPipelineType = cicd.PipelineDeploy
	case "full", "ci", "all":
		ciPipelineType = cicd.PipelineFull
	default:
		return fmt.Errorf("unsupported pipeline type: %s", pipelineType)
	}

	// Get flags
	outputDir, _ := cmd.Flags().GetString("output")
	triggers, _ := cmd.Flags().GetStringSlice("triggers")
	branches, _ := cmd.Flags().GetStringSlice("branches")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Create CI configuration
	ciConfig := &cicd.CIConfig{
		Provider:     ciProvider,
		PipelineType: ciPipelineType,
		Triggers:     triggers,
		Branches:     branches,
		Tags:         tags,
		Environments: make(map[string]string),
		Secrets:      make(map[string]string),
	}

	// Create CI orchestrator
	orchestrator := cicd.NewCIOrchestrator(cfg)

	// Register generators
	switch ciProvider {
	case cicd.ProviderGitHubActions:
		orchestrator.RegisterGenerator(ciProvider, cicd.NewGitHubActionsGenerator())
	default:
		return fmt.Errorf("generator not implemented for provider: %s", ciProvider)
	}

	// Generate pipeline
	result, err := orchestrator.GeneratePipeline(ciConfig)
	if err != nil {
		return fmt.Errorf("failed to generate pipeline: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("pipeline generation failed: %s", result.Error)
	}

	// Display results
	fmt.Printf("Generated %s pipeline for %s\n", ciPipelineType, ciProvider)
	fmt.Printf("Main file: %s\n", result.MainFile)
	fmt.Printf("Generated files:\n")
	for filePath := range result.Files {
		fmt.Printf("  - %s\n", filePath)
	}

	if dryRun {
		fmt.Printf("\nDry run - files not written to disk\n")
		fmt.Printf("Use without --dry-run to write files\n")
		return nil
	}

	// Write files
	if err := orchestrator.WritePipelineFiles(result, outputDir); err != nil {
		return fmt.Errorf("failed to write pipeline files: %v", err)
	}

	fmt.Printf("\nPipeline files written to: %s\n", outputDir)
	fmt.Printf("You can now commit these files to your repository.\n")

	return nil
}

func runCICDListCommandE(cmd *cobra.Command, args []string) error {
	fmt.Println("Available CI/CD Providers:")
	fmt.Println()

	providers := []struct {
		name        string
		provider    cicd.CIProvider
		description string
	}{
		{"GitHub Actions", cicd.ProviderGitHubActions, "GitHub's built-in CI/CD platform"},
		{"GitLab CI", cicd.ProviderGitLabCI, "GitLab's integrated CI/CD"},
		{"Jenkins", cicd.ProviderJenkins, "Popular open-source CI server"},
		{"CircleCI", cicd.ProviderCircleCI, "Cloud-based CI/CD platform"},
		{"Travis CI", cicd.ProviderTravisCI, "Hosted CI/CD for GitHub projects"},
	}

	for _, p := range providers {
		fmt.Printf("  %-15s %s\n", p.name, p.description)
	}

	fmt.Println()
	fmt.Println("Pipeline Types:")
	fmt.Println("  build    - Build pipeline only")
	fmt.Println("  test     - Test pipeline only")
	fmt.Println("  deploy   - Deploy pipeline only")
	fmt.Println("  full     - Complete CI/CD pipeline (build + test + deploy)")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("  forge cicd generate <provider> <type> [flags]")
	fmt.Println("  forge cicd generate github-actions full")
	fmt.Println("  forge cicd generate gitlab-ci build --output .gitlab-ci")

	return nil
}
