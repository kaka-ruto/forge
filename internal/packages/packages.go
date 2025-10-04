package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
)

// PackageInfo represents information about a package
type PackageInfo struct {
	Name         string
	Version      string
	Description  string
	Dependencies []string
	Conflicts    []string
	Category     string
	Size         int64  // Size in bytes
	BuildrootPkg string // Corresponding Buildroot package name
}

// PackageManager manages OS packages for Forge projects
type PackageManager struct {
	config     *config.Config
	logger     *logger.Logger
	packages   map[string]*PackageInfo
	categories map[string][]string // Category -> package names
}

// NewPackageManager creates a new package manager
func NewPackageManager(cfg *config.Config) *PackageManager {
	pm := &PackageManager{
		config:     cfg,
		logger:     logger.NewLogger(logger.INFO, nil, nil),
		packages:   make(map[string]*PackageInfo),
		categories: make(map[string][]string),
	}

	pm.initializePackageDatabase()
	return pm
}

// initializePackageDatabase sets up the package database with known packages
func (pm *PackageManager) initializePackageDatabase() {
	// Core system packages
	pm.addPackage(&PackageInfo{
		Name:         "busybox",
		Version:      "1.36.1",
		Description:  "Swiss Army Knife of Embedded Linux",
		Dependencies: []string{},
		Category:     "core",
		BuildrootPkg: "BR2_PACKAGE_BUSYBOX",
	})

	pm.addPackage(&PackageInfo{
		Name:         "openssh",
		Version:      "9.3",
		Description:  "OpenSSH connectivity tools",
		Dependencies: []string{"openssl", "zlib"},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_OPENSSH",
	})

	pm.addPackage(&PackageInfo{
		Name:         "openssl",
		Version:      "3.1.4",
		Description:  "OpenSSL cryptography and SSL/TLS toolkit",
		Dependencies: []string{},
		Category:     "security",
		BuildrootPkg: "BR2_PACKAGE_OPENSSL",
	})

	pm.addPackage(&PackageInfo{
		Name:         "zlib",
		Version:      "1.2.13",
		Description:  "Compression library",
		Dependencies: []string{},
		Category:     "libs",
		BuildrootPkg: "BR2_PACKAGE_ZLIB",
	})

	pm.addPackage(&PackageInfo{
		Name:         "python3",
		Version:      "3.11.5",
		Description:  "Python programming language",
		Dependencies: []string{"libffi", "expat"},
		Category:     "languages",
		BuildrootPkg: "BR2_PACKAGE_PYTHON3",
	})

	pm.addPackage(&PackageInfo{
		Name:         "libffi",
		Version:      "3.4.4",
		Description:  "Foreign Function Interface library",
		Dependencies: []string{},
		Category:     "libs",
		BuildrootPkg: "BR2_PACKAGE_LIBFFI",
	})

	pm.addPackage(&PackageInfo{
		Name:         "expat",
		Version:      "2.5.0",
		Description:  "XML parsing library",
		Dependencies: []string{},
		Category:     "libs",
		BuildrootPkg: "BR2_PACKAGE_EXPAT",
	})

	pm.addPackage(&PackageInfo{
		Name:         "mosquitto",
		Version:      "2.0.18",
		Description:  "MQTT broker",
		Dependencies: []string{"openssl", "libwebsockets"},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_MOSQUITTO",
	})

	pm.addPackage(&PackageInfo{
		Name:         "libwebsockets",
		Version:      "4.3.2",
		Description:  "WebSocket library",
		Dependencies: []string{"openssl", "zlib"},
		Category:     "libs",
		BuildrootPkg: "BR2_PACKAGE_LIBWEBSOCKETS",
	})

	pm.addPackage(&PackageInfo{
		Name:         "wpa_supplicant",
		Version:      "2.10",
		Description:  "WiFi supplicant",
		Dependencies: []string{"openssl", "dbus"},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_WPA_SUPPLICANT",
	})

	pm.addPackage(&PackageInfo{
		Name:         "dbus",
		Version:      "1.14.8",
		Description:  "D-Bus message bus system",
		Dependencies: []string{"expat"},
		Category:     "system",
		BuildrootPkg: "BR2_PACKAGE_DBUS",
	})

	pm.addPackage(&PackageInfo{
		Name:         "dhcpcd",
		Version:      "9.4.1",
		Description:  "DHCP client daemon",
		Dependencies: []string{},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_DHCPCD",
	})

	pm.addPackage(&PackageInfo{
		Name:         "nginx",
		Version:      "1.24.0",
		Description:  "HTTP and reverse proxy server",
		Dependencies: []string{"openssl", "zlib", "pcre"},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_NGINX",
	})

	pm.addPackage(&PackageInfo{
		Name:         "pcre",
		Version:      "8.45",
		Description:  "Perl Compatible Regular Expressions",
		Dependencies: []string{},
		Category:     "libs",
		BuildrootPkg: "BR2_PACKAGE_PCRE",
	})

	pm.addPackage(&PackageInfo{
		Name:         "openvpn",
		Version:      "2.6.6",
		Description:  "OpenVPN VPN daemon",
		Dependencies: []string{"openssl", "lzo"},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_OPENVPN",
	})

	pm.addPackage(&PackageInfo{
		Name:         "lzo",
		Version:      "2.10",
		Description:  "LZO compression library",
		Dependencies: []string{},
		Category:     "libs",
		BuildrootPkg: "BR2_PACKAGE_LZO",
	})

	pm.addPackage(&PackageInfo{
		Name:         "iptables",
		Version:      "1.8.9",
		Description:  "Linux kernel firewall",
		Dependencies: []string{},
		Category:     "network",
		BuildrootPkg: "BR2_PACKAGE_IPTABLES",
	})

	pm.addPackage(&PackageInfo{
		Name:         "fail2ban",
		Version:      "1.0.2",
		Description:  "Intrusion prevention system",
		Dependencies: []string{"python3", "iptables"},
		Category:     "security",
		BuildrootPkg: "BR2_PACKAGE_FAIL2BAN",
	})
}

// addPackage adds a package to the database
func (pm *PackageManager) addPackage(pkg *PackageInfo) {
	pm.packages[pkg.Name] = pkg
	pm.categories[pkg.Category] = append(pm.categories[pkg.Category], pkg.Name)
}

// GetPackageInfo returns information about a package
func (pm *PackageManager) GetPackageInfo(name string) (*PackageInfo, error) {
	pkg, exists := pm.packages[name]
	if !exists {
		return nil, fmt.Errorf("package %s not found", name)
	}
	return pkg, nil
}

// ListPackages returns a list of all available packages
func (pm *PackageManager) ListPackages() []*PackageInfo {
	var packages []*PackageInfo
	for _, pkg := range pm.packages {
		packages = append(packages, pkg)
	}

	// Sort by name for consistent output
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages
}

// ListPackagesByCategory returns packages in a specific category
func (pm *PackageManager) ListPackagesByCategory(category string) []*PackageInfo {
	var packages []*PackageInfo
	for _, name := range pm.categories[category] {
		if pkg, exists := pm.packages[name]; exists {
			packages = append(packages, pkg)
		}
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages
}

// GetCategories returns all available package categories
func (pm *PackageManager) GetCategories() []string {
	var categories []string
	for category := range pm.categories {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// IsValidPackage checks if a package name is valid
func (pm *PackageManager) IsValidPackage(name string) bool {
	_, exists := pm.packages[name]
	return exists
}

// DependencyResolution represents the result of dependency resolution
type DependencyResolution struct {
	Packages  []string // Packages in installation order
	Missing   []string // Packages that couldn't be found
	Circular  []string // Packages involved in circular dependencies
	Conflicts []string // Packages with conflicts
}

// ResolveDependencies resolves dependencies for a list of packages
func (pm *PackageManager) ResolveDependencies(packageNames []string) *DependencyResolution {
	result := &DependencyResolution{
		Packages:  make([]string, 0),
		Missing:   make([]string, 0),
		Circular:  make([]string, 0),
		Conflicts: make([]string, 0),
	}

	// Check for missing packages first
	for _, name := range packageNames {
		if !pm.IsValidPackage(name) {
			result.Missing = append(result.Missing, name)
		}
	}

	if len(result.Missing) > 0 {
		return result
	}

	// Collect all packages that need to be installed (including dependencies)
	allPackages := make(map[string]bool)
	for _, name := range packageNames {
		pm.collectAllDependencies(name, allPackages)
	}

	// Perform topological sort
	visited := make(map[string]bool)
	tempVisited := make(map[string]bool)
	installOrder := make([]string, 0)

	for pkg := range allPackages {
		if !visited[pkg] {
			if pm.topologicalSort(pkg, visited, tempVisited, &installOrder) {
				result.Circular = append(result.Circular, pkg)
				return result
			}
		}
	}

	result.Packages = installOrder
	result.Packages = installOrder
	return result
}

// hasCircularDependency performs DFS to detect circular dependencies and build topological order
func (pm *PackageManager) hasCircularDependency(pkgName string, visited, recursionStack map[string]bool, installOrder *[]string) bool {
	visited[pkgName] = true
	recursionStack[pkgName] = true

	pkg, _ := pm.packages[pkgName]
	for _, dep := range pkg.Dependencies {
		if !visited[dep] {
			if pm.hasCircularDependency(dep, visited, recursionStack, installOrder) {
				return true
			}
		} else if recursionStack[dep] {
			// Found circular dependency
			return true
		}
	}

	recursionStack[pkgName] = false
	// Add to install order after processing dependencies (dependencies first)
	*installOrder = append(*installOrder, pkgName)
	return false
}

// GetDependencyTree returns a tree representation of package dependencies
func (pm *PackageManager) GetDependencyTree(packageNames []string) map[string][]string {
	tree := make(map[string][]string)

	for _, name := range packageNames {
		if pkg, exists := pm.packages[name]; exists {
			tree[name] = pm.getAllDependencies(pkg)
		}
	}

	return tree
}

// getAllDependencies recursively gets all dependencies for a package
func (pm *PackageManager) getAllDependencies(pkg *PackageInfo) []string {
	visited := make(map[string]bool)
	var deps []string

	pm.collectDependencies(pkg.Name, visited, &deps)
	return deps
}

// collectDependencies recursively collects all dependencies
func (pm *PackageManager) collectDependencies(pkgName string, visited map[string]bool, deps *[]string) {
	if visited[pkgName] {
		return
	}

	visited[pkgName] = true

	pkg, exists := pm.packages[pkgName]
	if !exists {
		return
	}

	for _, dep := range pkg.Dependencies {
		pm.collectDependencies(dep, visited, deps)
	}

	*deps = append(*deps, pkgName)
}

// collectAllDependencies collects all packages that need to be installed for a given package
func (pm *PackageManager) collectAllDependencies(pkgName string, allPackages map[string]bool) {
	if allPackages[pkgName] {
		return
	}

	allPackages[pkgName] = true

	pkg, exists := pm.packages[pkgName]
	if !exists {
		return
	}

	for _, dep := range pkg.Dependencies {
		pm.collectAllDependencies(dep, allPackages)
	}
}

// topologicalSort performs topological sort using DFS
func (pm *PackageManager) topologicalSort(pkgName string, visited, tempVisited map[string]bool, result *[]string) bool {
	visited[pkgName] = true
	tempVisited[pkgName] = true

	pkg, exists := pm.packages[pkgName]
	if !exists {
		tempVisited[pkgName] = false
		return false
	}

	for _, dep := range pkg.Dependencies {
		if tempVisited[dep] {
			// Circular dependency detected
			tempVisited[pkgName] = false
			return true
		}
		if !visited[dep] {
			if pm.topologicalSort(dep, visited, tempVisited, result) {
				tempVisited[pkgName] = false
				return true
			}
		}
	}

	tempVisited[pkgName] = false
	// Add to result after all dependencies are processed
	*result = append(*result, pkgName)
	return false
}

// ValidatePackageSet checks if a set of packages can be installed together
func (pm *PackageManager) ValidatePackageSet(packageNames []string) *DependencyResolution {
	result := pm.ResolveDependencies(packageNames)

	if len(result.Missing) > 0 || len(result.Circular) > 0 {
		return result
	}

	// Check for conflicts
	installed := make(map[string]bool)
	for _, pkgName := range result.Packages {
		pkg := pm.packages[pkgName]

		// Check conflicts
		for _, conflict := range pkg.Conflicts {
			if installed[conflict] {
				result.Conflicts = append(result.Conflicts, fmt.Sprintf("%s conflicts with %s", pkgName, conflict))
			}
		}

		installed[pkgName] = true
	}

	return result
}

// GetRecommendedPackages returns recommended packages for a given set of packages
func (pm *PackageManager) GetRecommendedPackages(packageNames []string) []string {
	recommended := make(map[string]bool)
	added := make(map[string]bool)

	// Add explicitly requested packages
	for _, name := range packageNames {
		added[name] = true
	}

	// Add all dependencies (including transitive)
	for _, name := range packageNames {
		pm.addAllDependencies(name, recommended, added)
	}

	var result []string
	for pkg := range recommended {
		result = append(result, pkg)
	}

	sort.Strings(result)
	return result
}

// addAllDependencies recursively adds all dependencies
func (pm *PackageManager) addAllDependencies(pkgName string, recommended, added map[string]bool) {
	pkg, exists := pm.packages[pkgName]
	if !exists {
		return
	}

	for _, dep := range pkg.Dependencies {
		if !added[dep] {
			recommended[dep] = true
			pm.addAllDependencies(dep, recommended, added)
		}
	}
}

// InstallationResult represents the result of a package installation
type InstallationResult struct {
	Package     string
	Success     bool
	Error       string
	ConfigFiles []string // Configuration files created/modified
	Services    []string // Services that need to be started
}

// InstallPackages installs a set of packages
func (pm *PackageManager) InstallPackages(packageNames []string, buildrootDir string) []*InstallationResult {
	results := []*InstallationResult{}

	// Resolve dependencies
	depResult := pm.ResolveDependencies(packageNames)
	if len(depResult.Missing) > 0 {
		for _, missing := range depResult.Missing {
			results = append(results, &InstallationResult{
				Package: missing,
				Success: false,
				Error:   "package not found",
			})
		}
		return results
	}

	if len(depResult.Circular) > 0 {
		for _, circular := range depResult.Circular {
			results = append(results, &InstallationResult{
				Package: circular,
				Success: false,
				Error:   "circular dependency detected",
			})
		}
		return results
	}

	// Install packages in dependency order
	for _, pkgName := range depResult.Packages {
		result := pm.installPackage(pkgName, buildrootDir)
		results = append(results, result)
	}

	return results
}

// installPackage installs a single package
func (pm *PackageManager) installPackage(pkgName, buildrootDir string) *InstallationResult {
	result := &InstallationResult{
		Package: pkgName,
	}

	pkg, exists := pm.packages[pkgName]
	if !exists {
		result.Success = false
		result.Error = "package not found"
		return result
	}

	// Simulate package installation by enabling it in Buildroot config
	configPath := pm.getBuildrootConfigPath(buildrootDir)
	if configPath == "" {
		result.Success = false
		result.Error = "Buildroot config not found"
		return result
	}

	// Enable the package in Buildroot config
	err := pm.enablePackageInBuildroot(configPath, pkg.BuildrootPkg)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to enable package in Buildroot: %v", err)
		return result
	}

	// Generate configuration files if needed
	configFiles, err := pm.generatePackageConfig(pkgName, buildrootDir)
	if err != nil {
		pm.logger.Warn("Failed to generate config for %s: %v", pkgName, err)
	}

	// Get services that need to be started
	services := pm.getPackageServices(pkgName)

	result.Success = true
	result.ConfigFiles = configFiles
	result.Services = services

	return result
}

// getBuildrootConfigPath finds the Buildroot config file
func (pm *PackageManager) getBuildrootConfigPath(buildrootDir string) string {
	// Look for .config file in buildroot directory
	configPath := filepath.Join(buildrootDir, ".config")
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	// Look for config files in output directory
	outputDir := filepath.Join(buildrootDir, "output")
	configPath = filepath.Join(outputDir, ".config")
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	return ""
}

// enablePackageInBuildroot enables a package in the Buildroot config
func (pm *PackageManager) enablePackageInBuildroot(configPath, buildrootPkg string) error {
	// Read the current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	modified := false

	// Look for the package config line and enable it
	for i, line := range lines {
		if strings.HasPrefix(line, "# "+buildrootPkg+" is not set") {
			lines[i] = buildrootPkg + "=y"
			modified = true
			break
		} else if strings.HasPrefix(line, buildrootPkg+"=") {
			// Already configured, make sure it's enabled
			if !strings.Contains(line, "=y") {
				lines[i] = buildrootPkg + "=y"
				modified = true
			}
			break
		}
	}

	if !modified {
		// Add the config line if it doesn't exist
		lines = append(lines, buildrootPkg+"=y")
	}

	// Write back the modified config
	newContent := strings.Join(lines, "\n")
	return os.WriteFile(configPath, []byte(newContent), 0644)
}

// generatePackageConfig generates configuration files for a package
func (pm *PackageManager) generatePackageConfig(pkgName, buildrootDir string) ([]string, error) {
	var configFiles []string

	switch pkgName {
	case "openssh":
		configFiles = append(configFiles, pm.generateSSHConfig(buildrootDir)...)
	case "mosquitto":
		configFiles = append(configFiles, pm.generateMosquittoConfig(buildrootDir)...)
	case "nginx":
		configFiles = append(configFiles, pm.generateNginxConfig(buildrootDir)...)
	case "wpa_supplicant":
		configFiles = append(configFiles, pm.generateWPAConfig(buildrootDir)...)
	}

	return configFiles, nil
}

// generateSSHConfig generates SSH configuration
func (pm *PackageManager) generateSSHConfig(buildrootDir string) []string {
	// This would generate SSH config files
	// For now, just return the expected config files
	return []string{"/etc/ssh/sshd_config", "/etc/ssh/ssh_config"}
}

// generateMosquittoConfig generates Mosquitto MQTT broker configuration
func (pm *PackageManager) generateMosquittoConfig(buildrootDir string) []string {
	return []string{"/etc/mosquitto/mosquitto.conf"}
}

// generateNginxConfig generates Nginx web server configuration
func (pm *PackageManager) generateNginxConfig(buildrootDir string) []string {
	return []string{"/etc/nginx/nginx.conf", "/etc/nginx/sites-enabled/default"}
}

// generateWPAConfig generates WPA supplicant configuration
func (pm *PackageManager) generateWPAConfig(buildrootDir string) []string {
	return []string{"/etc/wpa_supplicant.conf"}
}

// getPackageServices returns services that need to be started for a package
func (pm *PackageManager) getPackageServices(pkgName string) []string {
	switch pkgName {
	case "openssh":
		return []string{"sshd"}
	case "mosquitto":
		return []string{"mosquitto"}
	case "nginx":
		return []string{"nginx"}
	case "dhcpcd":
		return []string{"dhcpcd"}
	case "wpa_supplicant":
		return []string{"wpa_supplicant"}
	default:
		return []string{}
	}
}

// UninstallPackages uninstalls a set of packages
func (pm *PackageManager) UninstallPackages(packageNames []string, buildrootDir string) []*InstallationResult {
	results := []*InstallationResult{}

	for _, pkgName := range packageNames {
		result := pm.uninstallPackage(pkgName, buildrootDir)
		results = append(results, result)
	}

	return results
}

// uninstallPackage uninstalls a single package
func (pm *PackageManager) uninstallPackage(pkgName, buildrootDir string) *InstallationResult {
	result := &InstallationResult{
		Package: pkgName,
	}

	pkg, exists := pm.packages[pkgName]
	if !exists {
		result.Success = false
		result.Error = "package not found"
		return result
	}

	// Disable the package in Buildroot config
	configPath := pm.getBuildrootConfigPath(buildrootDir)
	if configPath == "" {
		result.Success = false
		result.Error = "Buildroot config not found"
		return result
	}

	err := pm.disablePackageInBuildroot(configPath, pkg.BuildrootPkg)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to disable package in Buildroot: %v", err)
		return result
	}

	result.Success = true
	return result
}

// disablePackageInBuildroot disables a package in the Buildroot config
func (pm *PackageManager) disablePackageInBuildroot(configPath, buildrootPkg string) error {
	// Read the current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	modified := false

	// Look for the package config line and disable it
	for i, line := range lines {
		if strings.HasPrefix(line, buildrootPkg+"=y") {
			lines[i] = "# " + buildrootPkg + " is not set"
			modified = true
			break
		}
	}

	if modified {
		// Write back the modified config
		newContent := strings.Join(lines, "\n")
		return os.WriteFile(configPath, []byte(newContent), 0644)
	}

	return nil
}
