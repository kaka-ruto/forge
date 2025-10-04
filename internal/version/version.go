package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
}

// String returns the string representation of the version
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare compares two versions
// Returns -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}
	return 0
}

// ParseVersion parses a version string into a Version struct
func ParseVersion(s string) (*Version, error) {
	// Remove 'v' prefix if present
	s = strings.TrimPrefix(s, "v")

	// Split by dots
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid version format: %s", s)
	}

	// Parse each part
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// KernelType represents the type of kernel version
type KernelType string

const (
	KernelTypeLatest   KernelType = "latest"
	KernelTypeLTS      KernelType = "lts"
	KernelTypeSpecific KernelType = "specific"
)

// KernelVersion represents a kernel version specification
type KernelVersion struct {
	Version string
	Type    KernelType
}

// String returns the string representation of the kernel version
func (k *KernelVersion) String() string {
	return k.Version
}

// ParseKernelVersion parses a kernel version string
func ParseKernelVersion(s string) (*KernelVersion, error) {
	switch strings.ToLower(s) {
	case "latest":
		return &KernelVersion{Version: s, Type: KernelTypeLatest}, nil
	case "lts":
		return &KernelVersion{Version: s, Type: KernelTypeLTS}, nil
	default:
		// Check if it's a valid version number (x.y.z format)
		versionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
		if versionRegex.MatchString(s) {
			return &KernelVersion{Version: s, Type: KernelTypeSpecific}, nil
		}
		return nil, fmt.Errorf("invalid kernel version: %s", s)
	}
}

// IsLTSKernel checks if a kernel version is an LTS version
func IsLTSKernel(version string) bool {
	ltsVersions := []string{
		"5.15.0", "5.10.0", "5.4.0", "4.19.0", "4.14.0",
		// Add more LTS versions as needed
	}

	for _, lts := range ltsVersions {
		if strings.HasPrefix(version, lts) {
			return true
		}
	}
	return false
}

// ValidateBuildrootVersion validates a Buildroot version string
func ValidateBuildrootVersion(version string) (string, error) {
	switch strings.ToLower(version) {
	case "latest", "stable":
		return version, nil
	default:
		// Check if it's a valid YYYY.MM format
		versionRegex := regexp.MustCompile(`^\d{4}\.\d{2}$`)
		if versionRegex.MatchString(version) {
			return version, nil
		}
		return "", fmt.Errorf("invalid Buildroot version: %s", version)
	}
}

// CheckCompatibility checks if two versions are compatible
func CheckCompatibility(forgeVersion, configVersion *Version) bool {
	// Major versions must match for compatibility
	if forgeVersion.Major != configVersion.Major {
		return false
	}
	// Minor versions are backward compatible
	return true
}

// SchemaVersion represents a schema version
type SchemaVersion struct {
	Major int
	Minor int
}

// String returns the string representation of the schema version
func (s *SchemaVersion) String() string {
	return fmt.Sprintf("%d.%d", s.Major, s.Minor)
}

// ParseSchemaVersion parses a schema version string
func ParseSchemaVersion(s string) (*SchemaVersion, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid schema version format: %s", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major schema version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor schema version: %s", parts[1])
	}

	return &SchemaVersion{
		Major: major,
		Minor: minor,
	}, nil
}

// IsUpgradeAvailable checks if an upgrade is available
func IsUpgradeAvailable(current, latest *Version) bool {
	return current.Compare(latest) < 0
}

// GetDeprecationWarning returns a deprecation warning if the version is deprecated
func GetDeprecationWarning(version *Version) (string, bool) {
	// For now, consider versions < 1.0.0 as deprecated
	if version.Major < 1 {
		return fmt.Sprintf("Version %s is deprecated. Please upgrade to version 1.0.0 or later.", version.String()), true
	}
	return "", false
}

// HasBreakingChanges checks if there are breaking changes between versions
func HasBreakingChanges(oldVersion, newVersion *Version) (bool, string) {
	if oldVersion.Major != newVersion.Major {
		return true, "major"
	}
	return false, ""
}

// VersionInfo contains all version information
type VersionInfo struct {
	ForgeVersion     *Version
	BuildrootVersion string
	KernelVersion    *KernelVersion
	GoVersion        string
	BuildTimestamp   string
	GitCommit        string
}

// ParseVersionFile parses version information from a file content
func ParseVersionFile(content string) (*VersionInfo, error) {
	lines := strings.Split(content, "\n")
	info := &VersionInfo{}
	foundValidLine := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if value == "" {
			return nil, fmt.Errorf("empty value for key: %s", key)
		}

		foundValidLine = true

		switch key {
		case "forge_version":
			version, err := ParseVersion(value)
			if err != nil {
				return nil, fmt.Errorf("invalid forge version: %v", err)
			}
			info.ForgeVersion = version
		case "buildroot_version":
			info.BuildrootVersion = value
		case "kernel_version":
			kernelVersion, err := ParseKernelVersion(value)
			if err != nil {
				return nil, fmt.Errorf("invalid kernel version: %v", err)
			}
			info.KernelVersion = kernelVersion
		case "go_version":
			info.GoVersion = value
		case "build_timestamp":
			info.BuildTimestamp = value
		case "git_commit":
			info.GitCommit = value
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	if !foundValidLine {
		return nil, fmt.Errorf("no valid version information found")
	}

	return info, nil
}
