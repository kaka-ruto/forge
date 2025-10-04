package packages

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sst/forge/internal/config"
	"github.com/stretchr/testify/suite"
)

type PackagesTestSuite struct {
	suite.Suite
	config  *config.Config
	manager *PackageManager
}

func TestPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(PackagesTestSuite))
}

func (s *PackagesTestSuite) SetupTest() {
	s.config = &config.Config{
		SchemaVersion: "1.0",
		Name:          "test-project",
		Version:       "0.1.0",
		Architecture:  "x86_64",
		Template:      "minimal",
		Packages:      []string{},
		Features:      []string{},
	}

	s.manager = NewPackageManager(s.config)
}

func (s *PackagesTestSuite) TestNewPackageManager() {
	manager := NewPackageManager(s.config)
	s.NotNil(manager)
	s.Equal(s.config, manager.config)
	s.NotNil(manager.packages)
	s.NotNil(manager.categories)
}

func (s *PackagesTestSuite) TestGetPackageInfo() {
	// Test existing package
	pkg, err := s.manager.GetPackageInfo("openssh")
	s.NoError(err)
	s.NotNil(pkg)
	s.Equal("openssh", pkg.Name)
	s.Equal("9.3", pkg.Version)
	s.Equal("network", pkg.Category)

	// Test non-existing package
	pkg, err = s.manager.GetPackageInfo("nonexistent")
	s.Error(err)
	s.Nil(pkg)
	s.Contains(err.Error(), "package nonexistent not found")
}

func (s *PackagesTestSuite) TestListPackages() {
	packages := s.manager.ListPackages()
	s.NotEmpty(packages)

	// Check that packages are sorted
	for i := 1; i < len(packages); i++ {
		s.True(packages[i-1].Name <= packages[i].Name)
	}

	// Check some known packages
	names := make(map[string]bool)
	for _, pkg := range packages {
		names[pkg.Name] = true
	}

	s.True(names["busybox"])
	s.True(names["openssh"])
	s.True(names["python3"])
	s.True(names["nginx"])
}

func (s *PackagesTestSuite) TestListPackagesByCategory() {
	// Test network category
	networkPackages := s.manager.ListPackagesByCategory("network")
	s.NotEmpty(networkPackages)

	for _, pkg := range networkPackages {
		s.Equal("network", pkg.Category)
	}

	// Check that network packages are sorted
	for i := 1; i < len(networkPackages); i++ {
		s.True(networkPackages[i-1].Name <= networkPackages[i].Name)
	}

	// Test libs category
	libsPackages := s.manager.ListPackagesByCategory("libs")
	s.NotEmpty(libsPackages)

	for _, pkg := range libsPackages {
		s.Equal("libs", pkg.Category)
	}
}

func (s *PackagesTestSuite) TestGetCategories() {
	categories := s.manager.GetCategories()
	s.NotEmpty(categories)

	// Check that categories are sorted
	for i := 1; i < len(categories); i++ {
		s.True(categories[i-1] <= categories[i])
	}

	// Check expected categories
	categoryMap := make(map[string]bool)
	for _, cat := range categories {
		categoryMap[cat] = true
	}

	s.True(categoryMap["core"])
	s.True(categoryMap["network"])
	s.True(categoryMap["security"])
	s.True(categoryMap["libs"])
	s.True(categoryMap["languages"])
	s.True(categoryMap["system"])
}

func (s *PackagesTestSuite) TestIsValidPackage() {
	// Test valid packages
	s.True(s.manager.IsValidPackage("busybox"))
	s.True(s.manager.IsValidPackage("openssh"))
	s.True(s.manager.IsValidPackage("python3"))

	// Test invalid packages
	s.False(s.manager.IsValidPackage("nonexistent"))
	s.False(s.manager.IsValidPackage(""))
}

func (s *PackagesTestSuite) TestPackageDependencies() {
	// Test package with dependencies
	pkg, err := s.manager.GetPackageInfo("mosquitto")
	s.NoError(err)
	s.Contains(pkg.Dependencies, "openssl")
	s.Contains(pkg.Dependencies, "libwebsockets")

	// Test package with no dependencies
	pkg, err = s.manager.GetPackageInfo("busybox")
	s.NoError(err)
	s.Empty(pkg.Dependencies)
}

func (s *PackagesTestSuite) TestPackageCategories() {
	// Test core packages
	corePackages := s.manager.ListPackagesByCategory("core")
	s.NotEmpty(corePackages)

	// Test security packages
	securityPackages := s.manager.ListPackagesByCategory("security")
	s.NotEmpty(securityPackages)

	// Test languages packages
	langPackages := s.manager.ListPackagesByCategory("languages")
	s.NotEmpty(langPackages)
}

func (s *PackagesTestSuite) TestResolveDependencies() {
	// Test simple dependency resolution
	result := s.manager.ResolveDependencies([]string{"mosquitto"})
	s.NotEmpty(result.Packages)
	s.Empty(result.Missing)
	s.Empty(result.Circular)

	// Check that dependencies come before the package
	pkgIndex := -1
	opensslIndex := -1
	libwebsocketsIndex := -1

	for i, pkg := range result.Packages {
		switch pkg {
		case "mosquitto":
			pkgIndex = i
		case "openssl":
			opensslIndex = i
		case "libwebsockets":
			libwebsocketsIndex = i
		}
	}

	if pkgIndex != -1 && opensslIndex != -1 {
		s.True(pkgIndex > opensslIndex, "mosquitto should come after openssl (dependencies first)")
	}
	if pkgIndex != -1 && libwebsocketsIndex != -1 {
		s.True(pkgIndex > libwebsocketsIndex, "mosquitto should come after libwebsockets (dependencies first)")
	}
}

func (s *PackagesTestSuite) TestResolveDependenciesMissingPackage() {
	result := s.manager.ResolveDependencies([]string{"nonexistent"})
	s.Empty(result.Packages)
	s.Contains(result.Missing, "nonexistent")
}

func (s *PackagesTestSuite) TestGetDependencyTree() {
	tree := s.manager.GetDependencyTree([]string{"mosquitto"})
	s.NotEmpty(tree["mosquitto"])

	deps := tree["mosquitto"]
	s.Contains(deps, "openssl")
	s.Contains(deps, "libwebsockets")
	s.Contains(deps, "zlib") // transitive dependency
}

func (s *PackagesTestSuite) TestValidatePackageSet() {
	// Test valid package set
	result := s.manager.ValidatePackageSet([]string{"busybox", "openssh"})
	s.NotEmpty(result.Packages)
	s.Empty(result.Missing)
	s.Empty(result.Circular)
	s.Empty(result.Conflicts)
}

func (s *PackagesTestSuite) TestGetRecommendedPackages() {
	recommended := s.manager.GetRecommendedPackages([]string{"mosquitto"})
	s.NotEmpty(recommended)

	// Should include dependencies
	recommendedMap := make(map[string]bool)
	for _, pkg := range recommended {
		recommendedMap[pkg] = true
	}

	s.True(recommendedMap["openssl"])
	s.True(recommendedMap["libwebsockets"])
	s.True(recommendedMap["zlib"])
}

func (s *PackagesTestSuite) TestComplexDependencyResolution() {
	// Test with multiple packages that share dependencies
	result := s.manager.ResolveDependencies([]string{"nginx", "openvpn"})
	s.NotEmpty(result.Packages)
	s.Empty(result.Missing)
	s.Empty(result.Circular)

	// Check that openssl appears only once
	opensslCount := 0
	for _, pkg := range result.Packages {
		if pkg == "openssl" {
			opensslCount++
		}
	}
	s.Equal(1, opensslCount, "openssl should appear only once in the installation order")
}

func (s *PackagesTestSuite) TestInstallPackages() {
	// Create a temporary directory for testing
	tempDir := "/tmp/forge-test-buildroot"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Create a mock Buildroot config
	configPath := filepath.Join(tempDir, ".config")
	configContent := `# Test Buildroot config
BR2_PACKAGE_BUSYBOX=y
# BR2_PACKAGE_OPENSSH is not set
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Test installing openssh
	results := s.manager.InstallPackages([]string{"openssh"}, tempDir)
	s.Len(results, 3) // zlib + openssl + openssh dependencies

	// Check that all packages were installed successfully
	for _, result := range results {
		s.True(result.Success, "Package %s should install successfully", result.Package)
	}

	// Check that all expected packages are present
	packages := make(map[string]int)
	for i, result := range results {
		packages[result.Package] = i
	}

	s.Contains(packages, "zlib")
	s.Contains(packages, "openssl")
	s.Contains(packages, "openssh")

	// Check that dependencies come before dependents
	s.True(packages["zlib"] < packages["openssl"], "zlib should come before openssl")
	s.True(packages["openssl"] < packages["openssh"], "openssl should come before openssh")

	// Check that the config was updated
	updatedConfig, err := os.ReadFile(configPath)
	s.NoError(err)
	s.Contains(string(updatedConfig), "BR2_PACKAGE_OPENSSL=y")
	s.Contains(string(updatedConfig), "BR2_PACKAGE_OPENSSH=y")
}

func (s *PackagesTestSuite) TestInstallPackagesMissingBuildroot() {
	results := s.manager.InstallPackages([]string{"openssh"}, "/nonexistent")
	s.Len(results, 3) // zlib + openssl + openssh

	for _, result := range results {
		s.False(result.Success)
		s.Contains(result.Error, "Buildroot config not found")
	}
}

func (s *PackagesTestSuite) TestUninstallPackages() {
	// Create a temporary directory for testing
	tempDir := "/tmp/forge-test-buildroot-uninstall"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Create a mock Buildroot config with packages enabled
	configPath := filepath.Join(tempDir, ".config")
	configContent := `# Test Buildroot config
BR2_PACKAGE_BUSYBOX=y
BR2_PACKAGE_OPENSSL=y
BR2_PACKAGE_OPENSSH=y
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Test uninstalling openssh
	results := s.manager.UninstallPackages([]string{"openssh"}, tempDir)
	s.Len(results, 1)

	s.Equal("openssh", results[0].Package)
	s.True(results[0].Success)

	// Check that the config was updated
	updatedConfig, err := os.ReadFile(configPath)
	s.NoError(err)
	s.Contains(string(updatedConfig), "# BR2_PACKAGE_OPENSSH is not set")
}

func (s *PackagesTestSuite) TestGetPackageServices() {
	// Test services for different packages
	s.Equal([]string{"sshd"}, s.manager.getPackageServices("openssh"))
	s.Equal([]string{"mosquitto"}, s.manager.getPackageServices("mosquitto"))
	s.Equal([]string{"nginx"}, s.manager.getPackageServices("nginx"))
	s.Equal([]string{}, s.manager.getPackageServices("busybox"))
}

func (s *PackagesTestSuite) TestGeneratePackageConfig() {
	tempDir := "/tmp/forge-test-config"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Test config generation for different packages
	sshConfigs, _ := s.manager.generatePackageConfig("openssh", tempDir)
	s.Contains(sshConfigs, "/etc/ssh/sshd_config")

	mqttConfigs, _ := s.manager.generatePackageConfig("mosquitto", tempDir)
	s.Contains(mqttConfigs, "/etc/mosquitto/mosquitto.conf")

	nginxConfigs, _ := s.manager.generatePackageConfig("nginx", tempDir)
	s.Contains(nginxConfigs, "/etc/nginx/nginx.conf")
}
