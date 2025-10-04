package cicd

import (
	"fmt"
	"strings"

	"github.com/sst/forge/internal/config"
)

// GitHubActionsGenerator generates GitHub Actions workflows
type GitHubActionsGenerator struct{}

// NewGitHubActionsGenerator creates a new GitHub Actions generator
func NewGitHubActionsGenerator() *GitHubActionsGenerator {
	return &GitHubActionsGenerator{}
}

// Generate generates a GitHub Actions workflow
func (g *GitHubActionsGenerator) Generate(ciConfig *CIConfig, forgeConfig *config.Config) (*PipelineResult, error) {
	result := &PipelineResult{
		Files:   make(map[string]string),
		Success: false,
	}

	workflowName := g.getWorkflowName(ciConfig.PipelineType)
	workflowFile := fmt.Sprintf(".github/workflows/%s.yml", workflowName)

	workflowContent := g.generateWorkflow(ciConfig, forgeConfig, workflowName)

	result.Files[workflowFile] = workflowContent
	result.MainFile = workflowFile
	result.Success = true

	return result, nil
}

// Validate validates the GitHub Actions configuration
func (g *GitHubActionsGenerator) Validate(ciConfig *CIConfig) error {
	if ciConfig.Provider != ProviderGitHubActions {
		return fmt.Errorf("generator only supports GitHub Actions")
	}

	return ValidateCIConfig(ciConfig)
}

// GetSupportedTriggers returns supported GitHub Actions triggers
func (g *GitHubActionsGenerator) GetSupportedTriggers() []string {
	return []string{"push", "pull_request", "schedule", "workflow_dispatch", "release"}
}

// getWorkflowName returns the workflow name based on pipeline type
func (g *GitHubActionsGenerator) getWorkflowName(pipelineType PipelineType) string {
	switch pipelineType {
	case PipelineBuild:
		return "build"
	case PipelineTest:
		return "test"
	case PipelineDeploy:
		return "deploy"
	case PipelineFull:
		return "ci"
	default:
		return "forge"
	}
}

// generateWorkflow generates the GitHub Actions workflow YAML
func (g *GitHubActionsGenerator) generateWorkflow(ciConfig *CIConfig, forgeConfig *config.Config, workflowName string) string {
	var workflow strings.Builder

	workflow.WriteString(fmt.Sprintf(`name: %s

on:
`, g.getDisplayName(ciConfig.PipelineType)))

	// Add triggers
	workflow.WriteString(g.generateTriggers(ciConfig))

	workflow.WriteString(`
jobs:
`)

	// Add jobs based on pipeline type
	switch ciConfig.PipelineType {
	case PipelineBuild:
		workflow.WriteString(g.generateBuildJob(ciConfig, forgeConfig))
	case PipelineTest:
		workflow.WriteString(g.generateTestJob(ciConfig, forgeConfig))
	case PipelineDeploy:
		workflow.WriteString(g.generateDeployJob(ciConfig, forgeConfig))
	case PipelineFull:
		workflow.WriteString(g.generateFullCIJobs(ciConfig, forgeConfig))
	}

	return workflow.String()
}

// getDisplayName returns a display name for the pipeline type
func (g *GitHubActionsGenerator) getDisplayName(pipelineType PipelineType) string {
	switch pipelineType {
	case PipelineBuild:
		return "Build"
	case PipelineTest:
		return "Test"
	case PipelineDeploy:
		return "Deploy"
	case PipelineFull:
		return "CI"
	default:
		return "Forge OS"
	}
}

// generateTriggers generates the workflow triggers
func (g *GitHubActionsGenerator) generateTriggers(ciConfig *CIConfig) string {
	var triggers strings.Builder

	for _, trigger := range ciConfig.Triggers {
		switch trigger {
		case "push":
			triggers.WriteString("  push:\n")
			if len(ciConfig.Branches) > 0 {
				triggers.WriteString(fmt.Sprintf("    branches: %s\n", g.formatList(ciConfig.Branches)))
			}
			if len(ciConfig.Tags) > 0 {
				triggers.WriteString(fmt.Sprintf("    tags: %s\n", g.formatList(ciConfig.Tags)))
			}
		case "pull_request":
			triggers.WriteString("  pull_request:\n")
			if len(ciConfig.Branches) > 0 {
				triggers.WriteString(fmt.Sprintf("    branches: %s\n", g.formatList(ciConfig.Branches)))
			}
		case "schedule":
			triggers.WriteString("  schedule:\n")
			triggers.WriteString("    - cron: '0 2 * * 1'  # Weekly on Monday at 2 AM UTC\n")
		case "workflow_dispatch":
			triggers.WriteString("  workflow_dispatch:\n")
		case "release":
			triggers.WriteString("  release:\n")
			triggers.WriteString("    types: [published]\n")
		}
	}

	return triggers.String()
}

// generateBuildJob generates a build job
func (g *GitHubActionsGenerator) generateBuildJob(ciConfig *CIConfig, forgeConfig *config.Config) string {
	return fmt.Sprintf(`  build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Forge OS
      run: |
        go install ./cmd/forge

    - name: Create project
      run: |
        forge new test-project --template minimal --arch %s

    - name: Build image
      run: |
        cd test-project
        forge build

    - name: Upload build artifacts
      uses: actions/upload-artifact@v3
      with:
        name: forge-image
        path: test-project/build/artifacts/
`, forgeConfig.Architecture)
}

// generateTestJob generates a test job
func (g *GitHubActionsGenerator) generateTestJob(ciConfig *CIConfig, forgeConfig *config.Config) string {
	return fmt.Sprintf(`  test:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Forge OS
      run: |
        go install ./cmd/forge

    - name: Run tests
      run: |
        go test ./...

    - name: Create test project
      run: |
        forge new test-project --template minimal --arch %s

    - name: Build and test image
      run: |
        cd test-project
        forge build
        forge test --duration 30s

    - name: Upload test results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: test-results
        path: test-project/test-results/
`, forgeConfig.Architecture)
}

// generateDeployJob generates a deploy job
func (g *GitHubActionsGenerator) generateDeployJob(ciConfig *CIConfig, forgeConfig *config.Config) string {
	return fmt.Sprintf(`  deploy:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Forge OS
      run: |
        go install ./cmd/forge

    - name: Create project
      run: |
        forge new test-project --template minimal --arch %s

    - name: Build image
      run: |
        cd test-project
        forge build

    - name: Deploy to test environment
      run: |
        cd test-project
        forge deploy remote --host test.example.com --user forge --dry-run
`, forgeConfig.Architecture)
}

// generateFullCIJobs generates a full CI pipeline with build, test, and deploy
func (g *GitHubActionsGenerator) generateFullCIJobs(ciConfig *CIConfig, forgeConfig *config.Config) string {
	return fmt.Sprintf(`  build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Forge OS
      run: |
        go install ./cmd/forge

    - name: Create project
      run: |
        forge new test-project --template minimal --arch %s

    - name: Build image
      run: |
        cd test-project
        forge build

    - name: Upload build artifacts
      uses: actions/upload-artifact@v3
      with:
        name: forge-image
        path: test-project/build/artifacts/

  test:
    runs-on: ubuntu-latest
    needs: build
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Forge OS
      run: |
        go install ./cmd/forge

    - name: Run unit tests
      run: |
        go test ./...

    - name: Download build artifacts
      uses: actions/download-artifact@v3
      with:
        name: forge-image
        path: ./artifacts

    - name: Create test project
      run: |
        forge new test-project --template minimal --arch %s

    - name: Copy artifacts
      run: |
        cp -r artifacts/* test-project/build/artifacts/

    - name: Run integration tests
      run: |
        cd test-project
        forge test --duration 60s

    - name: Upload test results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: test-results
        path: test-project/test-results/

  deploy:
    runs-on: ubuntu-latest
    needs: [build, test]
    if: github.ref == 'refs/heads/main' && github.event_name == 'push'
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Download build artifacts
      uses: actions/download-artifact@v3
      with:
        name: forge-image
        path: ./artifacts

    - name: Create deployment project
      run: |
        forge new deploy-project --template minimal --arch %s

    - name: Copy artifacts for deployment
      run: |
        cp -r artifacts/* deploy-project/build/artifacts/

    - name: Deploy (dry run)
      run: |
        cd deploy-project
        forge deploy remote --host deploy.example.com --user forge --dry-run
`, forgeConfig.Architecture, forgeConfig.Architecture, forgeConfig.Architecture)
}

// formatList formats a list for YAML
func (g *GitHubActionsGenerator) formatList(items []string) string {
	if len(items) == 0 {
		return "[]"
	}

	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("'%s'", item)
	}

	return fmt.Sprintf("[%s]", strings.Join(quoted, ", "))
}
