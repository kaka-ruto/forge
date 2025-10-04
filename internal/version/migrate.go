package version

import (
	"fmt"
	"regexp"
	"strings"
)

// MigrationStep represents a single migration step
type MigrationStep struct {
	From        *SchemaVersion
	To          *SchemaVersion
	Description string
}

// String returns the string representation of a migration step
func (m *MigrationStep) String() string {
	return fmt.Sprintf("%s: %s", m.Description, fmt.Sprintf("%s -> %s", m.From.String(), m.To.String()))
}

// DetectSchemaVersion extracts the schema version from a config string
func DetectSchemaVersion(config string) (*SchemaVersion, error) {
	// Look for schema_version field
	re := regexp.MustCompile(`schema_version:\s*["']?([^"'\s]+)["']?`)
	matches := re.FindStringSubmatch(config)
	if len(matches) < 2 {
		return nil, fmt.Errorf("schema_version not found in config")
	}

	return ParseSchemaVersion(matches[1])
}

// CalculateMigrationPath calculates the migration path from one version to another
func CalculateMigrationPath(from, to *SchemaVersion) ([]MigrationStep, error) {
	if from.Compare(to) == 0 {
		return []MigrationStep{}, nil
	}

	if from.Compare(to) > 0 {
		return nil, fmt.Errorf("downgrade not supported: %s to %s", from.String(), to.String())
	}

	var steps []MigrationStep

	// For now, we only support direct migration
	// In the future, this could support multi-step migrations
	step := MigrationStep{
		From:        from,
		To:          to,
		Description: fmt.Sprintf("Upgrade from %s to %s", from.String(), to.String()),
	}
	steps = append(steps, step)

	return steps, nil
}

// Compare compares two schema versions
// Returns -1 if v < other, 0 if v == other, 1 if v > other
func (s *SchemaVersion) Compare(other *SchemaVersion) int {
	if s.Major != other.Major {
		if s.Major > other.Major {
			return 1
		}
		return -1
	}
	if s.Minor != other.Minor {
		if s.Minor > other.Minor {
			return 1
		}
		return -1
	}
	return 0
}

// MigrateV1ToV2 migrates a forge.yml config from version 1.0 to 2.0
func MigrateV1ToV2(config string) (string, error) {
	lines := strings.Split(config, "\n")
	var result []string
	schemaUpdated := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Update schema version
		if strings.Contains(trimmed, "schema_version") {
			result = append(result, strings.Replace(line, "\"1.0\"", "\"2.0\"", 1))
			schemaUpdated = true
			continue
		}

		result = append(result, line)
	}

	// Add missing fields if not present
	if !schemaUpdated {
		return "", fmt.Errorf("schema_version not found in config")
	}

	// Add missing fields at the end
	var newLines []string

	// Add buildroot_version if missing
	if !strings.Contains(config, "buildroot_version") {
		newLines = append(newLines, "buildroot_version: stable")
	}

	// Add kernel_version if missing
	if !strings.Contains(config, "kernel_version") {
		newLines = append(newLines, "kernel_version: lts")
	}

	// Append new lines at the end
	if len(newLines) > 0 {
		result = append(result, newLines...)
	}

	return strings.Join(result, "\n"), nil
}

// IsBackwardCompatible checks if a config is backward compatible
func IsBackwardCompatible(config string) bool {
	version, err := DetectSchemaVersion(config)
	if err != nil {
		return false
	}

	// For now, we support versions 1.0 and 2.0
	// In the future, this could check against a compatibility matrix
	return version.Major <= 2
}

// CanRollback checks if a migration can be rolled back
func CanRollback(steps []MigrationStep) bool {
	// For now, we support rollback for any migration
	// In the future, this could check if each step is reversible
	return true
}

// DryRunMigration performs a dry run of migration and returns the changes that would be made
func DryRunMigration(config string) ([]string, error) {
	var changes []string

	fromVersion, err := DetectSchemaVersion(config)
	if err != nil {
		return nil, err
	}

	toVersion := &SchemaVersion{Major: 2, Minor: 0} // Target version

	if fromVersion.Compare(toVersion) == 0 {
		return []string{"No migration needed"}, nil
	}

	changes = append(changes, fmt.Sprintf("Would update schema_version from %s to %s", fromVersion.String(), toVersion.String()))

	// Check what fields would be added
	if !strings.Contains(config, "buildroot_version") {
		changes = append(changes, "Would add buildroot_version: stable")
	}

	if !strings.Contains(config, "kernel_version") {
		changes = append(changes, "Would add kernel_version: lts")
	}

	return changes, nil
}

// GetBreakingChangeWarning returns a warning if there are breaking changes
func GetBreakingChangeWarning(from, to *SchemaVersion) (string, bool) {
	if from.Major != to.Major {
		return fmt.Sprintf("Major version upgrade from %s to %s may contain breaking changes", from.String(), to.String()), true
	}
	return "", false
}
