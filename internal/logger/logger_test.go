package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
	tempDir string
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-logger-test-*")
	s.Require().NoError(err)
}

func (s *LoggerTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *LoggerTestSuite) TestLogLevelFiltering() {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool // whether debug messages should be logged
	}{
		{
			name:     "debug level allows debug",
			level:    DEBUG,
			expected: true,
		},
		{
			name:     "info level blocks debug",
			level:    INFO,
			expected: false,
		},
		{
			name:     "warn level blocks debug",
			level:    WARN,
			expected: false,
		},
		{
			name:     "error level blocks debug",
			level:    ERROR,
			expected: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var buf bytes.Buffer
			logger := NewLogger(tt.level, &buf, &buf)

			// Always try to log a debug message to test filtering
			logger.Debug("debug test message")

			output := buf.String()
			if tt.expected {
				s.Contains(output, "debug test message")
			} else {
				s.NotContains(output, "debug test message")
			}
		})
	}
}

func (s *LoggerTestSuite) TestLogOutputFormatting() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	logger.Info("test message")

	output := buf.String()
	s.Contains(output, "INFO")
	s.Contains(output, "test message")
	s.Contains(output, time.Now().Format("2006-01-02")) // Date should be present
}

func (s *LoggerTestSuite) TestStructuredLogging() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)
	logger.SetFormat(JSON)

	logger.Info("test message")

	output := buf.String()
	s.Contains(output, `"level":"info"`)
	s.Contains(output, `"message":"test message"`)
	s.Contains(output, `"timestamp":`)
}

func (s *LoggerTestSuite) TestLogFileCreation() {
	logFile := filepath.Join(s.tempDir, "test.log")

	logger, err := NewFileLogger(DEBUG, logFile)
	s.NoError(err)
	s.NotNil(logger)

	logger.Info("test log message")
	logger.Close()

	// Check if file was created
	s.FileExists(logFile)

	// Check file contents
	content, err := os.ReadFile(logFile)
	s.NoError(err)
	s.Contains(string(content), "test log message")
}

func (s *LoggerTestSuite) TestLogFileRotation() {
	logFile := filepath.Join(s.tempDir, "test.log")

	logger, err := NewFileLogger(DEBUG, logFile)
	s.NoError(err)

	// Write enough logs to potentially trigger rotation
	for i := 0; i < 100; i++ {
		logger.Info("test message", i)
	}

	logger.Close()

	// Check if log file exists and has content
	s.FileExists(logFile)
	content, err := os.ReadFile(logFile)
	s.NoError(err)
	s.Contains(string(content), "test message")
}

func (s *LoggerTestSuite) TestConcurrentLogging() {
	logFile := filepath.Join(s.tempDir, "concurrent.log")

	logger, err := NewFileLogger(DEBUG, logFile)
	s.NoError(err)

	// Log from multiple goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Info("goroutine", id, "message", j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	logger.Close()

	// Check file contents
	content, err := os.ReadFile(logFile)
	s.NoError(err)
	contentStr := string(content)

	// Debug: print content
	s.T().Logf("Log content length: %d", len(contentStr))
	s.T().Logf("Log content preview: %s", contentStr[:min(500, len(contentStr))])

	// Should contain logs from all goroutines
	for i := 0; i < 10; i++ {
		expected := fmt.Sprintf("goroutine %d", i)
		s.Contains(contentStr, expected, "Should contain log for goroutine %d", i)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *LoggerTestSuite) TestLogContext() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	logger.WithField("component", "builder").Info("build started")

	output := buf.String()
	s.Contains(output, "component")
	s.Contains(output, "builder")
	s.Contains(output, "build started")
}

func (s *LoggerTestSuite) TestColoredOutput() {
	var errBuf bytes.Buffer
	logger := NewLogger(DEBUG, os.Stdout, &errBuf)
	logger.SetColored(true)

	logger.Error("error message")

	output := errBuf.String()
	// ANSI color codes should be present
	s.Contains(output, "\x1b[")
}

func (s *LoggerTestSuite) TestLogFilePathConfiguration() {
	customPath := filepath.Join(s.tempDir, "custom.log")

	logger, err := NewFileLogger(DEBUG, customPath)
	s.NoError(err)
	s.NotNil(logger)

	logger.Info("test message")
	logger.Close()

	s.FileExists(customPath)
}

func (s *LoggerTestSuite) TestGlobalLogger() {
	// Save original logger
	original := globalLogger

	// Set a test logger
	var buf bytes.Buffer
	testLogger := NewLogger(DEBUG, &buf, &buf)
	SetGlobalLogger(testLogger)

	// Test global logging functions
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	output := buf.String()
	s.Contains(output, "debug message")
	s.Contains(output, "info message")
	s.Contains(output, "warn message")
	s.Contains(output, "error message")

	// Restore original logger
	SetGlobalLogger(original)
}

func (s *LoggerTestSuite) TestLoggerWithFields() {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf, &buf)

	entry := logger.WithFields(Fields{
		"component": "builder",
		"operation": "compile",
	})

	entry.Info("compilation completed")

	output := buf.String()
	s.Contains(output, "component")
	s.Contains(output, "builder")
	s.Contains(output, "operation")
	s.Contains(output, "compile")
	s.Contains(output, "compilation completed")
}

func (s *LoggerTestSuite) TestLogLevelString() {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "debug"},
		{INFO, "info"},
		{WARN, "warn"},
		{ERROR, "error"},
	}

	for _, tt := range tests {
		s.Equal(tt.expected, tt.level.String())
	}
}

func (s *LoggerTestSuite) TestParseLogLevel() {
	tests := []struct {
		input    string
		expected LogLevel
		hasError bool
	}{
		{"debug", DEBUG, false},
		{"DEBUG", DEBUG, false},
		{"info", INFO, false},
		{"warn", WARN, false},
		{"error", ERROR, false},
		{"invalid", DEBUG, true}, // defaults to DEBUG on error
	}

	for _, tt := range tests {
		result, err := ParseLogLevel(tt.input)
		if tt.hasError {
			s.Error(err)
		} else {
			s.NoError(err)
			s.Equal(tt.expected, result)
		}
	}
}

func (s *LoggerTestSuite) TestLoggerClose() {
	logFile := filepath.Join(s.tempDir, "close-test.log")

	logger, err := NewFileLogger(DEBUG, logFile)
	s.NoError(err)

	logger.Info("before close")
	logger.Close()

	// Should not panic on subsequent calls
	logger.Close()

	// File should still exist
	s.FileExists(logFile)
}
