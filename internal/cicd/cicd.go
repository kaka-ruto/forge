package cicd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
)

// CIProvider represents different CI/CD platforms
type CIProvider string

const (
	ProviderGitHubActions CIProvider = "github-actions"
	ProviderGitLabCI      CIProvider = "gitlab-ci"
	ProviderJenkins       CIProvider = "jenkins"
	ProviderCircleCI      CIProvider = "circleci"
	ProviderTravisCI      CIProvider = "travis-ci"
)

// PipelineType represents different pipeline types
type PipelineType string

const (
	PipelineBuild  PipelineType = "build"
	PipelineTest   PipelineType = "test"
	PipelineDeploy PipelineType = "deploy"
	PipelineFull   PipelineType = "full"
)

// CIConfig holds CI/CD configuration
type CIConfig struct {
	Provider     CIProvider
	PipelineType PipelineType
	Triggers     []string          // Events that trigger the pipeline
	Branches     []string          // Branches to run on
	Tags         []string          // Tags to run on
	Environments map[string]string // Environment variables
	Secrets      map[string]string // Secret names (not values)
}

// PipelineResult represents the result of pipeline generation
type PipelineResult struct {
	Files    map[string]string // File path -> content
	MainFile string            // Main pipeline file
	Success  bool
	Error    string
}

// PipelineGenerator interface for generating CI/CD pipelines
type PipelineGenerator interface {
	Generate(config *CIConfig, forgeConfig *config.Config) (*PipelineResult, error)
	Validate(config *CIConfig) error
	GetSupportedTriggers() []string
}

// CIOrchestrator manages CI/CD pipeline generation
type CIOrchestrator struct {
	config     *config.Config
	logger     *logger.Logger
	generators map[CIProvider]PipelineGenerator
}

// NewCIOrchestrator creates a new CI/CD orchestrator
func NewCIOrchestrator(cfg *config.Config) *CIOrchestrator {
	return &CIOrchestrator{
		config:     cfg,
		logger:     logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
		generators: make(map[CIProvider]PipelineGenerator),
	}
}

// RegisterGenerator registers a pipeline generator for a CI provider
func (co *CIOrchestrator) RegisterGenerator(provider CIProvider, generator PipelineGenerator) {
	co.generators[provider] = generator
}

// GeneratePipeline generates a CI/CD pipeline for the specified provider
func (co *CIOrchestrator) GeneratePipeline(ciConfig *CIConfig) (*PipelineResult, error) {
	generator, exists := co.generators[ciConfig.Provider]
	if !exists {
		return nil, fmt.Errorf("no generator registered for provider: %s", ciConfig.Provider)
	}

	// Validate configuration
	if err := generator.Validate(ciConfig); err != nil {
		return nil, fmt.Errorf("invalid CI configuration: %v", err)
	}

	co.logger.Info("Generating %s pipeline for %s", ciConfig.PipelineType, ciConfig.Provider)
	result, err := generator.Generate(ciConfig, co.config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pipeline: %v", err)
	}

	if result.Success {
		co.logger.Info("Pipeline generated successfully")
	} else {
		co.logger.Error("Pipeline generation failed: %s", result.Error)
	}

	return result, nil
}

// WritePipelineFiles writes the generated pipeline files to disk
func (co *CIOrchestrator) WritePipelineFiles(result *PipelineResult, outputDir string) error {
	if !result.Success {
		return fmt.Errorf("cannot write files for failed pipeline generation")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Write each file
	for filePath, content := range result.Files {
		fullPath := filepath.Join(outputDir, filePath)

		// Create directory for file if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", fullPath, err)
		}

		co.logger.Info("Wrote pipeline file: %s", fullPath)
	}

	return nil
}

// GetAvailableProviders returns all available CI providers
func (co *CIOrchestrator) GetAvailableProviders() []CIProvider {
	var providers []CIProvider
	for provider := range co.generators {
		providers = append(providers, provider)
	}
	return providers
}

// ValidateCIConfig validates a CI configuration
func ValidateCIConfig(config *CIConfig) error {
	if config.Provider == "" {
		return fmt.Errorf("CI provider not specified")
	}

	if config.PipelineType == "" {
		return fmt.Errorf("pipeline type not specified")
	}

	// Validate pipeline type
	validTypes := []PipelineType{PipelineBuild, PipelineTest, PipelineDeploy, PipelineFull}
	typeValid := false
	for _, t := range validTypes {
		if config.PipelineType == t {
			typeValid = true
			break
		}
	}
	if !typeValid {
		return fmt.Errorf("invalid pipeline type: %s", config.PipelineType)
	}

	return nil
}

// DefaultCIConfig returns a default CI configuration for a provider
func DefaultCIConfig(provider CIProvider) *CIConfig {
	config := &CIConfig{
		Provider:     provider,
		PipelineType: PipelineFull,
		Triggers:     []string{"push", "pull_request"},
		Branches:     []string{"main", "master", "develop"},
		Tags:         []string{"v*"},
		Environments: make(map[string]string),
		Secrets:      make(map[string]string),
	}

	// Set provider-specific defaults
	switch provider {
	case ProviderGitHubActions:
		config.Environments["GITHUB_TOKEN"] = "${{ secrets.GITHUB_TOKEN }}"
	case ProviderGitLabCI:
		config.Environments["CI_REGISTRY_USER"] = "$CI_REGISTRY_USER"
		config.Environments["CI_REGISTRY_PASSWORD"] = "$CI_REGISTRY_PASSWORD"
	}

	return config
}
