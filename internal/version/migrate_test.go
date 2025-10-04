package version

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrateTestSuite struct {
	suite.Suite
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateTestSuite))
}

func (s *MigrateTestSuite) TestSchemaVersionDetection() {
	tests := []struct {
		name     string
		config   string
		expected *SchemaVersion
		hasError bool
	}{
		{
			name:     "schema version 1.0",
			config:   "schema_version: \"1.0\"",
			expected: &SchemaVersion{Major: 1, Minor: 0},
			hasError: false,
		},
		{
			name:     "schema version 1.1",
			config:   "schema_version: \"1.1\"",
			expected: &SchemaVersion{Major: 1, Minor: 1},
			hasError: false,
		},
		{
			name:     "no schema version",
			config:   "name: test",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid schema version",
			config:   "schema_version: \"invalid\"",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := DetectSchemaVersion(tt.config)
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

func (s *MigrateTestSuite) TestMigrationPathCalculation() {
	tests := []struct {
		name     string
		from     *SchemaVersion
		to       *SchemaVersion
		expected []MigrationStep
		hasError bool
	}{
		{
			name:     "no migration needed",
			from:     &SchemaVersion{Major: 1, Minor: 0},
			to:       &SchemaVersion{Major: 1, Minor: 0},
			expected: []MigrationStep{},
			hasError: false,
		},
		{
			name: "minor version upgrade",
			from: &SchemaVersion{Major: 1, Minor: 0},
			to:   &SchemaVersion{Major: 1, Minor: 1},
			expected: []MigrationStep{
				{From: &SchemaVersion{Major: 1, Minor: 0}, To: &SchemaVersion{Major: 1, Minor: 1}, Description: "Upgrade from 1.0 to 1.1"},
			},
			hasError: false,
		},
		{
			name: "major version upgrade",
			from: &SchemaVersion{Major: 1, Minor: 0},
			to:   &SchemaVersion{Major: 2, Minor: 0},
			expected: []MigrationStep{
				{From: &SchemaVersion{Major: 1, Minor: 0}, To: &SchemaVersion{Major: 2, Minor: 0}, Description: "Upgrade from 1.0 to 2.0"},
			},
			hasError: false,
		},
		{
			name:     "downgrade not allowed",
			from:     &SchemaVersion{Major: 2, Minor: 0},
			to:       &SchemaVersion{Major: 1, Minor: 0},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := CalculateMigrationPath(tt.from, tt.to)
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

func (s *MigrateTestSuite) TestForgeYmlMigrationFromV1ToV2() {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name: "migrate v1 to v2 with missing fields",
			input: `schema_version: "1.0"
name: test-project
architecture: x86_64`,
			expected: `schema_version: "2.0"
name: test-project
architecture: x86_64
buildroot_version: stable
kernel_version: lts`,
			hasError: false,
		},
		{
			name: "migrate v1 to v2 with some fields present",
			input: `schema_version: "1.0"
name: test-project
architecture: arm
buildroot_version: 2023.02`,
			expected: `schema_version: "2.0"
name: test-project
architecture: arm
buildroot_version: 2023.02
kernel_version: lts`,
			hasError: false,
		},
		{
			name: "migrate v1 to v2 with all fields present",
			input: `schema_version: "1.0"
name: test-project
architecture: aarch64
buildroot_version: stable
kernel_version: latest`,
			expected: `schema_version: "2.0"
name: test-project
architecture: aarch64
buildroot_version: stable
kernel_version: latest`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := MigrateV1ToV2(tt.input)
			if tt.hasError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}

func (s *MigrateTestSuite) TestBackwardCompatibility() {
	tests := []struct {
		name     string
		config   string
		expected bool
	}{
		{
			name: "v1.0 config compatible",
			config: `schema_version: "1.0"
name: test`,
			expected: true,
		},
		{
			name: "v2.0 config compatible",
			config: `schema_version: "2.0"
name: test`,
			expected: true,
		},
		{
			name: "future version not compatible",
			config: `schema_version: "3.0"
name: test`,
			expected: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := IsBackwardCompatible(tt.config)
			s.Equal(tt.expected, result)
		})
	}
}

func (s *MigrateTestSuite) TestMigrationRollback() {
	tests := []struct {
		name     string
		steps    []MigrationStep
		expected bool
	}{
		{
			name: "successful rollback",
			steps: []MigrationStep{
				{From: &SchemaVersion{Major: 1, Minor: 0}, To: &SchemaVersion{Major: 1, Minor: 1}},
			},
			expected: true,
		},
		{
			name:     "no steps to rollback",
			steps:    []MigrationStep{},
			expected: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := CanRollback(tt.steps)
			s.Equal(tt.expected, result)
		})
	}
}

func (s *MigrateTestSuite) TestMigrationDryRun() {
	tests := []struct {
		name     string
		config   string
		expected []string
		hasError bool
	}{
		{
			name: "dry run v1 to v2",
			config: `schema_version: "1.0"
name: test-project
architecture: x86_64`,
			expected: []string{
				"Would update schema_version from 1.0 to 2.0",
				"Would add buildroot_version: stable",
				"Would add kernel_version: lts",
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := DryRunMigration(tt.config)
			if tt.hasError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, result)
			}
		})
	}
}

func (s *MigrateTestSuite) TestBreakingChangeWarnings() {
	tests := []struct {
		name        string
		from        *SchemaVersion
		to          *SchemaVersion
		expectedMsg string
		hasWarning  bool
	}{
		{
			name:        "major version change warning",
			from:        &SchemaVersion{Major: 1, Minor: 0},
			to:          &SchemaVersion{Major: 2, Minor: 0},
			expectedMsg: "Major version upgrade from 1.0 to 2.0 may contain breaking changes",
			hasWarning:  true,
		},
		{
			name:        "minor version no warning",
			from:        &SchemaVersion{Major: 1, Minor: 0},
			to:          &SchemaVersion{Major: 1, Minor: 1},
			expectedMsg: "",
			hasWarning:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			msg, hasWarning := GetBreakingChangeWarning(tt.from, tt.to)
			s.Equal(tt.hasWarning, hasWarning)
			if tt.hasWarning {
				s.Equal(tt.expectedMsg, msg)
			}
		})
	}
}
