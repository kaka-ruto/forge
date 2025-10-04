package cicd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type CICDTestSuite struct {
	suite.Suite
	config    *config.Config
	tempDir   string
	outputDir string
}

func TestCICDTestSuite(t *testing.T) {
	suite.Run(t, new(CICDTestSuite))
}

func (s *CICDTestSuite) SetupTest() {
	s.config = &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{},
		Features:      []string{},
	}

	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-cicd-test-*")
	s.Require().NoError(err)

	s.outputDir = filepath.Join(s.tempDir, "output")
	err = os.MkdirAll(s.outputDir, 0755)
	s.Require().NoError(err)
}

func (s *CICDTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *CICDTestSuite) TestNewCIOrchestrator() {
	orchestrator := NewCIOrchestrator(s.config)
	s.NotNil(orchestrator)
	s.NotNil(orchestrator.generators)
}

func (s *CICDTestSuite) TestRegisterGenerator() {
	orchestrator := NewCIOrchestrator(s.config)
	generator := NewGitHubActionsGenerator()

	orchestrator.RegisterGenerator(ProviderGitHubActions, generator)
	s.Contains(orchestrator.generators, ProviderGitHubActions)
}

func (s *CICDTestSuite) TestGetAvailableProviders() {
	orchestrator := NewCIOrchestrator(s.config)
	generator := NewGitHubActionsGenerator()

	orchestrator.RegisterGenerator(ProviderGitHubActions, generator)
	providers := orchestrator.GetAvailableProviders()

	s.Contains(providers, ProviderGitHubActions)
}

func (s *CICDTestSuite) TestValidateCIConfig() {
	// Valid config
	config := &CIConfig{
		Provider:     ProviderGitHubActions,
		PipelineType: PipelineFull,
	}
	err := ValidateCIConfig(config)
	s.NoError(err)

	// Invalid provider
	config.Provider = ""
	err = ValidateCIConfig(config)
	s.Error(err)

	// Invalid pipeline type
	config.Provider = ProviderGitHubActions
	config.PipelineType = ""
	err = ValidateCIConfig(config)
	s.Error(err)
}

func (s *CICDTestSuite) TestDefaultCIConfig() {
	config := DefaultCIConfig(ProviderGitHubActions)
	s.Equal(ProviderGitHubActions, config.Provider)
	s.Equal(PipelineFull, config.PipelineType)
	s.Contains(config.Triggers, "push")
	s.Contains(config.Branches, "main")
}

func (s *CICDTestSuite) TestGeneratePipeline() {
	orchestrator := NewCIOrchestrator(s.config)
	generator := NewGitHubActionsGenerator()
	orchestrator.RegisterGenerator(ProviderGitHubActions, generator)

	ciConfig := &CIConfig{
		Provider:     ProviderGitHubActions,
		PipelineType: PipelineBuild,
		Triggers:     []string{"push"},
		Branches:     []string{"main"},
	}

	result, err := orchestrator.GeneratePipeline(ciConfig)
	s.NoError(err)
	s.True(result.Success)
	s.Contains(result.Files, ".github/workflows/build.yml")
	s.Equal(".github/workflows/build.yml", result.MainFile)
}

func (s *CICDTestSuite) TestWritePipelineFiles() {
	// Create a mock result
	result := &PipelineResult{
		Files: map[string]string{
			"test.yml": "name: test\non: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n    - run: echo hello",
		},
		MainFile: "test.yml",
		Success:  true,
	}

	orchestrator := NewCIOrchestrator(s.config)
	err := orchestrator.WritePipelineFiles(result, s.outputDir)
	s.NoError(err)

	// Check that file was written
	outputFile := filepath.Join(s.outputDir, "test.yml")
	s.FileExists(outputFile)

	// Check content
	content, err := os.ReadFile(outputFile)
	s.NoError(err)
	s.Contains(string(content), "name: test")
}

func (s *CICDTestSuite) TestGitHubActionsGeneratorValidate() {
	generator := NewGitHubActionsGenerator()

	// Valid config
	config := &CIConfig{
		Provider:     ProviderGitHubActions,
		PipelineType: PipelineBuild,
	}
	err := generator.Validate(config)
	s.NoError(err)

	// Wrong provider
	config.Provider = ProviderGitLabCI
	err = generator.Validate(config)
	s.Error(err)
}

func (s *CICDTestSuite) TestGitHubActionsGeneratorGetSupportedTriggers() {
	generator := NewGitHubActionsGenerator()
	triggers := generator.GetSupportedTriggers()

	s.Contains(triggers, "push")
	s.Contains(triggers, "pull_request")
	s.Contains(triggers, "schedule")
}

func (s *CICDTestSuite) TestGitHubActionsGeneratorGenerate() {
	generator := NewGitHubActionsGenerator()

	ciConfig := &CIConfig{
		Provider:     ProviderGitHubActions,
		PipelineType: PipelineBuild,
		Triggers:     []string{"push"},
		Branches:     []string{"main"},
	}

	result, err := generator.Generate(ciConfig, s.config)
	s.NoError(err)
	s.True(result.Success)
	s.Contains(result.Files, ".github/workflows/build.yml")

	workflowContent := result.Files[".github/workflows/build.yml"]
	s.Contains(workflowContent, "name: Build")
	s.Contains(workflowContent, "on:")
	s.Contains(workflowContent, "push:")
	s.Contains(workflowContent, "jobs:")
	s.Contains(workflowContent, "runs-on: ubuntu-latest")
}
