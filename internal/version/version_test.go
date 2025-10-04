package version

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type VersionTestSuite struct {
	suite.Suite
}

func TestVersionTestSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

func (s *VersionTestSuite) TestForgeVersionDetection() {
	tests := []struct {
		name     string
		input    string
		expected *Version
		hasError bool
	}{
		{
			name:  "valid semantic version",
			input: "1.2.3",
			expected: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			hasError: false,
		},
		{
			name:  "valid version with v prefix",
			input: "v1.2.3",
			expected: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			hasError: false,
		},
		{
			name:     "invalid version",
			input:    "invalid",
			expected: nil,
			hasError: true,
		},
		{
			name:     "empty version",
			input:    "",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := ParseVersion(tt.input)
			if tt.hasError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}

func (s *VersionTestSuite) TestVersionComparison() {
	tests := []struct {
		name     string
		v1       *Version
		v2       *Version
		expected int // -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
	}{
		{
			name:     "equal versions",
			v1:       &Version{Major: 1, Minor: 2, Patch: 3},
			v2:       &Version{Major: 1, Minor: 2, Patch: 3},
			expected: 0,
		},
		{
			name:     "v1 greater major version",
			v1:       &Version{Major: 2, Minor: 0, Patch: 0},
			v2:       &Version{Major: 1, Minor: 5, Patch: 9},
			expected: 1,
		},
		{
			name:     "v1 lesser major version",
			v1:       &Version{Major: 1, Minor: 0, Patch: 0},
			v2:       &Version{Major: 2, Minor: 0, Patch: 0},
			expected: -1,
		},
		{
			name:     "v1 greater minor version",
			v1:       &Version{Major: 1, Minor: 3, Patch: 0},
			v2:       &Version{Major: 1, Minor: 2, Patch: 9},
			expected: 1,
		},
		{
			name:     "v1 greater patch version",
			v1:       &Version{Major: 1, Minor: 2, Patch: 4},
			v2:       &Version{Major: 1, Minor: 2, Patch: 3},
			expected: 1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := tt.v1.Compare(tt.v2)
			s.Equal(tt.expected, result)
		})
	}
}

func (s *VersionTestSuite) TestBuildrootVersionPinning() {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "valid buildroot version",
			input:    "2023.02",
			expected: "2023.02",
			hasError: false,
		},
		{
			name:     "latest version",
			input:    "latest",
			expected: "latest",
			hasError: false,
		},
		{
			name:     "stable version",
			input:    "stable",
			expected: "stable",
			hasError: false,
		},
		{
			name:     "invalid version format",
			input:    "invalid-format",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := ValidateBuildrootVersion(tt.input)
			if tt.hasError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}

func (s *VersionTestSuite) TestKernelVersionSelection() {
	tests := []struct {
		name     string
		input    string
		expected *KernelVersion
		hasError bool
	}{
		{
			name:  "latest kernel",
			input: "latest",
			expected: &KernelVersion{
				Version: "latest",
				Type:    KernelTypeLatest,
			},
			hasError: false,
		},
		{
			name:  "LTS kernel",
			input: "lts",
			expected: &KernelVersion{
				Version: "lts",
				Type:    KernelTypeLTS,
			},
			hasError: false,
		},
		{
			name:  "specific kernel version",
			input: "5.15.0",
			expected: &KernelVersion{
				Version: "5.15.0",
				Type:    KernelTypeSpecific,
			},
			hasError: false,
		},
		{
			name:     "invalid kernel version",
			input:    "invalid",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := ParseKernelVersion(tt.input)
			if tt.hasError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}

func (s *VersionTestSuite) TestLTSKernelDetection() {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "LTS kernel 5.15",
			version:  "5.15.0",
			expected: true,
		},
		{
			name:     "LTS kernel 5.10",
			version:  "5.10.0",
			expected: true,
		},
		{
			name:     "non-LTS kernel 5.16",
			version:  "5.16.0",
			expected: false,
		},
		{
			name:     "latest keyword",
			version:  "latest",
			expected: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := IsLTSKernel(tt.version)
			s.Equal(tt.expected, result)
		})
	}
}

func (s *VersionTestSuite) TestVersionCompatibilityChecking() {
	tests := []struct {
		name        string
		forgeVer    string
		configVer   string
		expected    bool
		description string
	}{
		{
			name:        "compatible versions",
			forgeVer:    "1.0.0",
			configVer:   "1.0.0",
			expected:    true,
			description: "Same versions are compatible",
		},
		{
			name:        "major version mismatch",
			forgeVer:    "2.0.0",
			configVer:   "1.0.0",
			expected:    false,
			description: "Different major versions are incompatible",
		},
		{
			name:        "minor version compatible",
			forgeVer:    "1.1.0",
			configVer:   "1.0.0",
			expected:    true,
			description: "Higher minor versions are backward compatible",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			forgeV, _ := ParseVersion(tt.forgeVer)
			configV, _ := ParseVersion(tt.configVer)
			result := CheckCompatibility(forgeV, configV)
			s.Equal(tt.expected, result, tt.description)
		})
	}
}

func (s *VersionTestSuite) TestForgeYmlSchemaVersioning() {
	tests := []struct {
		name     string
		input    string
		expected *SchemaVersion
		hasError bool
	}{
		{
			name:  "valid schema version",
			input: "1.0",
			expected: &SchemaVersion{
				Major: 1,
				Minor: 0,
			},
			hasError: false,
		},
		{
			name:     "invalid schema version",
			input:    "invalid",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := ParseSchemaVersion(tt.input)
			if tt.hasError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}

func (s *VersionTestSuite) TestVersionUpgradeDetection() {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected bool
	}{
		{
			name:     "upgrade available",
			current:  "1.0.0",
			latest:   "1.1.0",
			expected: true,
		},
		{
			name:     "no upgrade needed",
			current:  "1.1.0",
			latest:   "1.1.0",
			expected: false,
		},
		{
			name:     "downgrade not detected",
			current:  "1.1.0",
			latest:   "1.0.0",
			expected: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			currentV, _ := ParseVersion(tt.current)
			latestV, _ := ParseVersion(tt.latest)
			result := IsUpgradeAvailable(currentV, latestV)
			s.Equal(tt.expected, result)
		})
	}
}

func (s *VersionTestSuite) TestDeprecationWarnings() {
	tests := []struct {
		name        string
		version     string
		expectedMsg string
		hasWarning  bool
	}{
		{
			name:        "deprecated version",
			version:     "0.9.0",
			expectedMsg: "Version 0.9.0 is deprecated",
			hasWarning:  true,
		},
		{
			name:        "current version",
			version:     "1.0.0",
			expectedMsg: "",
			hasWarning:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			v, _ := ParseVersion(tt.version)
			msg, hasWarning := GetDeprecationWarning(v)
			s.Equal(tt.hasWarning, hasWarning)
			if tt.hasWarning {
				s.Contains(msg, tt.expectedMsg)
			}
		})
	}
}

func (s *VersionTestSuite) TestBreakingChangeDetection() {
	tests := []struct {
		name       string
		oldVer     string
		newVer     string
		expected   bool
		breaksType string
	}{
		{
			name:       "major version change",
			oldVer:     "1.0.0",
			newVer:     "2.0.0",
			expected:   true,
			breaksType: "major",
		},
		{
			name:       "minor version change",
			oldVer:     "1.0.0",
			newVer:     "1.1.0",
			expected:   false,
			breaksType: "",
		},
		{
			name:       "patch version change",
			oldVer:     "1.0.0",
			newVer:     "1.0.1",
			expected:   false,
			breaksType: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			oldV, _ := ParseVersion(tt.oldVer)
			newV, _ := ParseVersion(tt.newVer)
			result, breaksType := HasBreakingChanges(oldV, newV)
			s.Equal(tt.expected, result)
			s.Equal(tt.breaksType, breaksType)
		})
	}
}

func (s *VersionTestSuite) TestVersionFileParsing() {
	tests := []struct {
		name     string
		content  string
		expected *VersionInfo
		hasError bool
	}{
		{
			name:    "valid version file",
			content: "forge_version=1.0.0\nbuildroot_version=2023.02\nkernel_version=5.15.0\ngo_version=1.21.0\nbuild_timestamp=2023-01-01T00:00:00Z\ngit_commit=abc123\n",
			expected: &VersionInfo{
				ForgeVersion:     &Version{Major: 1, Minor: 0, Patch: 0},
				BuildrootVersion: "2023.02",
				KernelVersion:    &KernelVersion{Version: "5.15.0", Type: KernelTypeSpecific},
				GoVersion:        "1.21.0",
				BuildTimestamp:   "2023-01-01T00:00:00Z",
				GitCommit:        "abc123",
			},
			hasError: false,
		},
		{
			name:     "invalid version file",
			content:  "invalid content",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := ParseVersionFile(tt.content)
			if tt.hasError {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}
