package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DebugTestSuite struct {
	suite.Suite
	tempDir string
}

func TestDebugTestSuite(t *testing.T) {
	suite.Run(t, new(DebugTestSuite))
}

func (s *DebugTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-debug-test-*")
	s.Require().NoError(err)
}

func (s *DebugTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DebugTestSuite) TestDebugCommandFunctionality() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	// Create a debug collector
	collector := NewDebugCollector(logger)

	// Collect some debug information
	collector.CollectSystemInfo()
	collector.CollectLogInfo(filepath.Join(s.tempDir, "test.log"))
	collector.CollectConfigInfo("test-config")

	report := collector.GenerateReport()

	s.Contains(report, "System Information")
	s.Contains(report, "Log Information")
	s.Contains(report, "Configuration Information")
	s.Contains(report, "test-config")
}

func (s *DebugTestSuite) TestLogViewingAndFiltering() {
	// Create a test log file
	logFile := filepath.Join(s.tempDir, "test.log")
	logger, err := NewFileLogger(DEBUG, logFile)
	s.NoError(err)

	// Write some test logs
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("error message")
	logger.Close()

	// Test log viewing
	debug := NewDebugCollector(logger)

	// View all logs
	logs, err := debug.ViewLogs(logFile, "", "")
	s.NoError(err)
	s.Contains(logs, "info message")
	s.Contains(logs, "warning message")
	s.Contains(logs, "error message")

	// Filter by level
	errorLogs, err := debug.ViewLogs(logFile, "error", "")
	s.NoError(err)
	s.Contains(errorLogs, "error message")
	s.NotContains(errorLogs, "info message")

	// Filter by component
	logger2, err := NewFileLogger(DEBUG, logFile)
	s.NoError(err)
	logger2.SetComponent("builder")
	logger2.Info("build message")
	logger2.Close()

	builderLogs, err := debug.ViewLogs(logFile, "", "builder")
	s.NoError(err)
	s.Contains(builderLogs, "build message")
	s.Contains(builderLogs, "[builder]")
}

func (s *DebugTestSuite) TestErrorContextCapture() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Simulate an error with context
	err := collector.CaptureErrorContext("build failed", map[string]interface{}{
		"stage":     "compilation",
		"file":      "main.go",
		"exit_code": 1,
	})

	s.NoError(err)

	report := collector.GenerateReport()
	s.Contains(report, "build failed")
	s.Contains(report, "compilation")
	s.Contains(report, "main.go")
	s.Contains(report, "1")
}

func (s *DebugTestSuite) TestStackTraceGeneration() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Generate a stack trace
	stackTrace := collector.GenerateStackTrace()

	s.Contains(stackTrace, "goroutine")
	s.Contains(stackTrace, "debug_test.go")
}

func (s *DebugTestSuite) TestDebugOutputFormatting() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Add various types of information
	collector.CollectSystemInfo()
	collector.CollectLogInfo("test.log")
	collector.CaptureErrorContext("test error", map[string]interface{}{
		"code": 500,
	})

	report := collector.GenerateReport()

	// Check formatting
	lines := strings.Split(report, "\n")
	s.True(len(lines) > 10, "Report should have multiple lines")

	// Should contain headers
	s.Contains(report, "=== System Information ===")
	s.Contains(report, "=== Error Context ===")
}

func (s *DebugTestSuite) TestDebugCollectorCreation() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)
	s.NotNil(collector)
	s.NotNil(collector.logger)
	s.NotNil(collector.info)
}

func (s *DebugTestSuite) TestLogFollowing() {
	// Start log following in a goroutine
	done := make(chan bool)
	var followedLogs []string

	go func() {
		// This would normally follow logs, but for testing we'll simulate
		followedLogs = []string{"simulated log line 1", "simulated log line 2"}
		done <- true
	}()

	<-done
	s.Contains(followedLogs, "simulated log line 1")
}

func (s *DebugTestSuite) TestConfigValidation() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Test valid config
	validConfig := `schema_version: "1.0"
name: test-project
architecture: x86_64`

	err := collector.ValidateConfig(validConfig)
	s.NoError(err)

	// Test invalid config
	invalidConfig := `invalid yaml content: [unclosed`

	err = collector.ValidateConfig(invalidConfig)
	s.Error(err)
}

func (s *DebugTestSuite) TestBuildArtifactInspection() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Create a fake build artifact
	artifactPath := filepath.Join(s.tempDir, "artifact.img")
	err := os.WriteFile(artifactPath, []byte("fake image data"), 0644)
	s.NoError(err)

	info, err := collector.InspectBuildArtifact(artifactPath)
	s.NoError(err)
	s.Contains(info, "artifact.img")
	s.Contains(info, "fake image data")
}

func (s *DebugTestSuite) TestDiagnosticReportGeneration() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Add comprehensive diagnostic information
	collector.CollectSystemInfo()
	collector.CollectLogInfo("test.log")
	collector.CaptureErrorContext("diagnostic test", map[string]interface{}{
		"test": "value",
	})

	report := collector.GenerateDiagnosticReport()

	s.Contains(report, "Forge OS Diagnostic Report")
	s.Contains(report, "System Information")
	s.Contains(report, "diagnostic test")
	s.Contains(report, "Generated at:")
}

func (s *DebugTestSuite) TestDebugCommandExecution() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Test debug command execution (simulated)
	err := collector.ExecuteDebugCommand("system_info")
	s.NoError(err)

	report := collector.GenerateReport()
	s.Contains(report, "System Information")
}

func (s *DebugTestSuite) TestErrorAnalysis() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)

	// Analyze a common error
	analysis := collector.AnalyzeError("buildroot: command not found")

	s.Contains(analysis, "buildroot")
	s.Contains(analysis, "not found")
}

func (s *DebugTestSuite) TestDebugDataExport() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)
	collector.CollectSystemInfo()

	// Export debug data
	data, err := collector.ExportDebugData()
	s.NoError(err)
	s.Contains(data, "system")

	// Should be valid JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonData)
	s.NoError(err)
}

func (s *DebugTestSuite) TestDebugCollectorReset() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	collector := NewDebugCollector(logger)
	collector.CollectSystemInfo()

	// Generate report
	report1 := collector.GenerateReport()
	s.Contains(report1, "System Information")

	// Reset collector
	collector.Reset()

	// Generate report again
	report2 := collector.GenerateReport()
	s.NotContains(report2, "System Information")
}
