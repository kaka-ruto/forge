package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LogsCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestLogsCommandTestSuite(t *testing.T) {
	suite.Run(t, new(LogsCommandTestSuite))
}

func (s *LogsCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-logs-test-*")
	s.Require().NoError(err)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	s.tempDir = oldDir + "/" + s.tempDir
	os.Chdir(s.tempDir)
}

func (s *LogsCommandTestSuite) TearDownTest() {
	os.Chdir("/")
	os.RemoveAll(s.tempDir)
}

func (s *LogsCommandTestSuite) TestNewLogsCommand() {
	cmd := NewLogsCommand()
	s.NotNil(cmd)
	s.Equal("logs", cmd.Use)
	s.Contains(cmd.Short, "View and manage Forge OS logs")
}

func (s *LogsCommandTestSuite) TestLogsCommandNoLogs() {
	err := runLogsCommand([]string{}, map[string]interface{}{
		"level":     "",
		"component": "",
		"follow":    false,
		"tail":      50,
	})
	s.NoError(err)
}

func (s *LogsCommandTestSuite) TestLogsCommandWithLogFile() {
	// Create a test log file
	logContent := `2024-01-01 10:00:00 [INFO] Build started
2024-01-01 10:00:01 [DEBUG] Loading configuration
2024-01-01 10:00:02 [WARN] Configuration warning
2024-01-01 10:00:03 [ERROR] Build failed
2024-01-01 10:00:04 [INFO] Build completed`

	err := os.WriteFile("forge.log", []byte(logContent), 0644)
	s.NoError(err)

	err = runLogsCommand([]string{}, map[string]interface{}{
		"level":     "",
		"component": "",
		"follow":    false,
		"tail":      50,
	})
	s.NoError(err)
}

func (s *LogsCommandTestSuite) TestLogsCommandWithForgeLogsDir() {
	// Create .forge/logs directory
	logsDir := ".forge/logs"
	err := os.MkdirAll(logsDir, 0755)
	s.NoError(err)

	// Create a log file in the directory
	logContent := `2024-01-01 10:00:00 [INFO] Test log entry`
	err = os.WriteFile(filepath.Join(logsDir, "build.log"), []byte(logContent), 0644)
	s.NoError(err)

	err = runLogsCommand([]string{}, map[string]interface{}{
		"level":     "",
		"component": "",
		"follow":    false,
		"tail":      50,
	})
	s.NoError(err)
}

func (s *LogsCommandTestSuite) TestFilterLogLines() {
	lines := []string{
		"2024-01-01 10:00:00 [INFO] Info message",
		"2024-01-01 10:00:01 [DEBUG] Debug message",
		"2024-01-01 10:00:02 [ERROR] Error message",
		"",
		"2024-01-01 10:00:03 [INFO] Another info",
	}

	// Test no filter
	filtered := filterLogLines(lines, map[string]interface{}{
		"level":     "",
		"component": "",
	})
	s.Len(filtered, 4) // Excludes empty line

	// Test level filter
	filtered = filterLogLines(lines, map[string]interface{}{
		"level":     "ERROR",
		"component": "",
	})
	s.Len(filtered, 1)
	s.Contains(filtered[0], "[ERROR]")

	// Test component filter (not implemented in detail yet)
	filtered = filterLogLines(lines, map[string]interface{}{
		"level":     "",
		"component": "test",
	})
	s.Len(filtered, 0) // No lines contain "test"
}
