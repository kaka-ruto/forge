package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	// DEBUG level for detailed diagnostic information
	DEBUG LogLevel = iota
	// INFO level for general informational messages
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		return "unknown"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(s string) (LogLevel, error) {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn", "warning":
		return WARN, nil
	case "error":
		return ERROR, nil
	default:
		return DEBUG, fmt.Errorf("invalid log level: %s", s)
	}
}

// LogFormat represents the format of log output
type LogFormat int

const (
	// TEXT format for human-readable logs
	TEXT LogFormat = iota
	// JSON format for structured logs
	JSON
)

// Fields represents key-value pairs for structured logging
type Fields map[string]interface{}

// Logger represents a logging instance
type Logger struct {
	level     LogLevel
	format    LogFormat
	colored   bool
	out       io.Writer
	errOut    io.Writer
	file      *os.File
	mu        sync.Mutex
	fields    Fields
	component string
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, out, errOut io.Writer) *Logger {
	return &Logger{
		level:  level,
		format: TEXT,
		out:    out,
		errOut: errOut,
		fields: make(Fields),
	}
}

// NewFileLogger creates a logger that writes to a file
func NewFileLogger(level LogLevel, filename string) (*Logger, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	logger := NewLogger(level, file, file)
	logger.file = file
	return logger, nil
}

// SetFormat sets the log output format
func (l *Logger) SetFormat(format LogFormat) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.format = format
}

// SetColored enables or disables colored output
func (l *Logger) SetColored(colored bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.colored = colored
}

// SetComponent sets the component name for this logger
func (l *Logger) SetComponent(component string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.component = component
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:     l.level,
		format:    l.format,
		colored:   l.colored,
		out:       l.out,
		errOut:    l.errOut,
		file:      l.file,
		fields:    make(Fields),
		component: l.component,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields Fields) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:     l.level,
		format:    l.format,
		colored:   l.colored,
		out:       l.out,
		errOut:    l.errOut,
		file:      l.file,
		fields:    make(Fields),
		component: l.component,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// log writes a log entry
func (l *Logger) log(level LogLevel, message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	// Format message
	if len(args) > 0 {
		// If message contains format verbs, use sprintf
		if strings.Contains(message, "%") {
			message = fmt.Sprintf(message, args...)
		} else {
			// Otherwise, join all arguments with spaces
			allParts := []string{message}
			for _, arg := range args {
				allParts = append(allParts, fmt.Sprintf("%v", arg))
			}
			message = strings.Join(allParts, " ")
		}
	}

	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Component: l.component,
		Fields:    l.fields,
	}

	var output string
	var writer io.Writer

	switch l.format {
	case JSON:
		// Create a custom struct for JSON marshaling with string level
		jsonEntry := struct {
			Level     string    `json:"level"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			Component string    `json:"component,omitempty"`
			Fields    Fields    `json:"fields,omitempty"`
		}{
			Level:     level.String(),
			Message:   entry.Message,
			Timestamp: entry.Timestamp,
			Component: entry.Component,
			Fields:    entry.Fields,
		}

		jsonData, err := json.Marshal(jsonEntry)
		if err != nil {
			// Fallback to text format
			output = l.formatText(entry)
		} else {
			output = string(jsonData) + "\n"
		}
	default:
		output = l.formatText(entry)
	}

	// Choose output writer based on level
	if level >= ERROR {
		writer = l.errOut
	} else {
		writer = l.out
	}

	fmt.Fprint(writer, output)
}

// formatText formats a log entry as text
func (l *Logger) formatText(entry LogEntry) string {
	var builder strings.Builder

	// Timestamp
	builder.WriteString(entry.Timestamp.Format("2006-01-02 15:04:05"))

	// Level
	levelStr := strings.ToUpper(entry.Level.String())
	if l.colored {
		levelStr = colorizeLevel(entry.Level, levelStr)
	}
	builder.WriteString(fmt.Sprintf(" [%s]", levelStr))

	// Component
	if entry.Component != "" {
		builder.WriteString(fmt.Sprintf(" [%s]", entry.Component))
	}

	// Fields
	if len(entry.Fields) > 0 {
		builder.WriteString(" {")
		first := true
		for k, v := range entry.Fields {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%s=%v", k, v))
			first = false
		}
		builder.WriteString("}")
	}

	// Message
	builder.WriteString(fmt.Sprintf(" %s\n", entry.Message))

	return builder.String()
}

// colorizeLevel adds ANSI color codes to log level
func colorizeLevel(level LogLevel, text string) string {
	switch level {
	case DEBUG:
		return fmt.Sprintf("\x1b[36m%s\x1b[0m", text) // Cyan
	case INFO:
		return fmt.Sprintf("\x1b[32m%s\x1b[0m", text) // Green
	case WARN:
		return fmt.Sprintf("\x1b[33m%s\x1b[0m", text) // Yellow
	case ERROR:
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", text) // Red
	default:
		return text
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(DEBUG, message, args...)
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(INFO, message, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(WARN, message, args...)
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(ERROR, message, args...)
}

// Close closes the logger and any open files
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// LogEntry represents a single log entry
type LogEntry struct {
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Component string    `json:"component,omitempty"`
	Fields    Fields    `json:"fields,omitempty"`
}

// Global logger instance
var globalLogger *Logger

func init() {
	// Initialize with a default logger
	globalLogger = NewLogger(INFO, os.Stdout, os.Stderr)
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// Debug logs a debug message using the global logger
func Debug(message string, args ...interface{}) {
	globalLogger.log(DEBUG, message, args...)
}

// Info logs an info message using the global logger
func Info(message string, args ...interface{}) {
	globalLogger.log(INFO, message, args...)
}

// Warn logs a warning message using the global logger
func Warn(message string, args ...interface{}) {
	globalLogger.log(WARN, message, args...)
}

// Error logs an error message using the global logger
func Error(message string, args ...interface{}) {
	globalLogger.log(ERROR, message, args...)
}
