package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// DebugCollector collects debugging information
type DebugCollector struct {
	logger *Logger
	info   map[string]interface{}
}

// NewDebugCollector creates a new debug collector
func NewDebugCollector(logger *Logger) *DebugCollector {
	return &DebugCollector{
		logger: logger,
		info:   make(map[string]interface{}),
	}
}

// CollectSystemInfo collects system information
func (d *DebugCollector) CollectSystemInfo() {
	systemInfo := map[string]interface{}{
		"go_version":    runtime.Version(),
		"go_os":         runtime.GOOS,
		"go_arch":       runtime.GOARCH,
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	// Get environment variables
	env := make(map[string]string)
	for _, envVar := range os.Environ() {
		if strings.HasPrefix(envVar, "BR2_") || strings.HasPrefix(envVar, "GO") {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
	}
	systemInfo["environment"] = env

	d.info["system"] = systemInfo
}

// CollectLogInfo collects log file information
func (d *DebugCollector) CollectLogInfo(logPath string) {
	logInfo := map[string]interface{}{
		"path": logPath,
	}

	if _, err := os.Stat(logPath); err == nil {
		// File exists
		if content, err := os.ReadFile(logPath); err == nil {
			logInfo["size"] = len(content)
			logInfo["exists"] = true

			// Get last few lines
			lines := strings.Split(string(content), "\n")
			if len(lines) > 10 {
				logInfo["last_lines"] = lines[len(lines)-11:]
			} else {
				logInfo["last_lines"] = lines
			}
		}
	} else {
		logInfo["exists"] = false
		logInfo["error"] = err.Error()
	}

	d.info["logs"] = logInfo
}

// CollectConfigInfo collects configuration information
func (d *DebugCollector) CollectConfigInfo(configPath string) {
	configInfo := map[string]interface{}{
		"path": configPath,
	}

	if content, err := os.ReadFile(configPath); err == nil {
		configInfo["content"] = string(content)
		configInfo["size"] = len(content)
		configInfo["readable"] = true
	} else {
		configInfo["readable"] = false
		configInfo["error"] = err.Error()
	}

	d.info["config"] = configInfo
}

// CaptureErrorContext captures error context information
func (d *DebugCollector) CaptureErrorContext(message string, context map[string]interface{}) error {
	errorInfo := map[string]interface{}{
		"message":   message,
		"context":   context,
		"timestamp": time.Now().Format(time.RFC3339),
		"stack":     d.GenerateStackTrace(),
	}

	d.info["error"] = errorInfo
	return nil
}

// GenerateStackTrace generates a stack trace
func (d *DebugCollector) GenerateStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// ViewLogs views and filters log content
func (d *DebugCollector) ViewLogs(logPath, levelFilter, componentFilter string) (string, error) {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return "", fmt.Errorf("log file does not exist: %s", logPath)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("failed to read log file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var filteredLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Apply filters
		if levelFilter != "" && !strings.Contains(strings.ToUpper(line), strings.ToUpper(levelFilter)) {
			continue
		}

		if componentFilter != "" && !strings.Contains(line, "["+componentFilter+"]") {
			continue
		}

		filteredLines = append(filteredLines, line)
	}

	return strings.Join(filteredLines, "\n"), nil
}

// ValidateConfig validates configuration content
func (d *DebugCollector) ValidateConfig(configContent string) error {
	// Basic validation - check for basic structure
	if strings.TrimSpace(configContent) == "" {
		return fmt.Errorf("empty configuration")
	}

	// Check for basic YAML structure
	if !strings.Contains(configContent, ":") {
		return fmt.Errorf("invalid configuration format: missing key-value pairs")
	}

	// Try to parse as basic key-value format
	lines := strings.Split(configContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for unclosed brackets/braces/quotes
		if strings.Contains(line, "[") && !strings.Contains(line, "]") {
			return fmt.Errorf("invalid configuration: unclosed bracket")
		}
		if strings.Contains(line, "{") && !strings.Contains(line, "}") {
			return fmt.Errorf("invalid configuration: unclosed brace")
		}
		if strings.Count(line, "\"")%2 != 0 {
			return fmt.Errorf("invalid configuration: unclosed quote")
		}
	}

	return nil
}

// InspectBuildArtifact inspects a build artifact
func (d *DebugCollector) InspectBuildArtifact(artifactPath string) (string, error) {
	info, err := os.Stat(artifactPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat artifact: %v", err)
	}

	result := fmt.Sprintf("Artifact: %s\n", filepath.Base(artifactPath))
	result += fmt.Sprintf("Size: %d bytes\n", info.Size())
	result += fmt.Sprintf("Modified: %s\n", info.ModTime().Format(time.RFC3339))

	// Try to read first few bytes
	if file, err := os.Open(artifactPath); err == nil {
		defer file.Close()
		buf := make([]byte, 256)
		if n, err := file.Read(buf); err == nil {
			result += fmt.Sprintf("First %d bytes: %s\n", n, string(buf[:n]))
		}
	}

	return result, nil
}

// GenerateReport generates a debug report
func (d *DebugCollector) GenerateReport() string {
	var report strings.Builder

	report.WriteString("=== Forge OS Debug Report ===\n")
	report.WriteString(fmt.Sprintf("Generated at: %s\n\n", time.Now().Format(time.RFC3339)))

	if systemInfo, ok := d.info["system"]; ok {
		report.WriteString("=== System Information ===\n")
		d.appendMapToReport(&report, systemInfo.(map[string]interface{}))
		report.WriteString("\n")
	}

	if logInfo, ok := d.info["logs"]; ok {
		report.WriteString("=== Log Information ===\n")
		d.appendMapToReport(&report, logInfo.(map[string]interface{}))
		report.WriteString("\n")
	}

	if configInfo, ok := d.info["config"]; ok {
		report.WriteString("=== Configuration Information ===\n")
		d.appendMapToReport(&report, configInfo.(map[string]interface{}))
		report.WriteString("\n")
	}

	if errorInfo, ok := d.info["error"]; ok {
		report.WriteString("=== Error Context ===\n")
		d.appendMapToReport(&report, errorInfo.(map[string]interface{}))
		report.WriteString("\n")
	}

	return report.String()
}

// GenerateDiagnosticReport generates a comprehensive diagnostic report
func (d *DebugCollector) GenerateDiagnosticReport() string {
	var report strings.Builder

	report.WriteString("=== Forge OS Diagnostic Report ===\n")
	report.WriteString(fmt.Sprintf("Generated at: %s\n\n", time.Now().Format(time.RFC3339)))

	d.CollectSystemInfo()

	report.WriteString(d.GenerateReport())

	// Add additional diagnostic information
	report.WriteString("=== Recommendations ===\n")
	if runtime.NumCPU() < 4 {
		report.WriteString("- Consider using a machine with more CPU cores for better build performance\n")
	}

	if systemInfo, ok := d.info["system"].(map[string]interface{}); ok {
		if env, ok := systemInfo["environment"].(map[string]string); ok {
			if _, hasBr2Dl := env["BR2_DL_DIR"]; !hasBr2Dl {
				report.WriteString("- Set BR2_DL_DIR environment variable to cache Buildroot downloads\n")
			}
		}
	}

	return report.String()
}

// ExecuteDebugCommand executes a debug command
func (d *DebugCollector) ExecuteDebugCommand(command string) error {
	switch command {
	case "system_info":
		d.CollectSystemInfo()
	case "stack_trace":
		d.info["stack_trace"] = d.GenerateStackTrace()
	default:
		return fmt.Errorf("unknown debug command: %s", command)
	}
	return nil
}

// AnalyzeError analyzes an error message
func (d *DebugCollector) AnalyzeError(errorMsg string) string {
	analysis := fmt.Sprintf("Error Analysis: %s\n", errorMsg)

	// Common error patterns
	if strings.Contains(errorMsg, "command not found") {
		analysis += "This indicates a missing dependency. Check if required tools are installed.\n"
	}

	if strings.Contains(errorMsg, "permission denied") {
		analysis += "This indicates insufficient permissions. Try running with elevated privileges.\n"
	}

	if strings.Contains(errorMsg, "no space left on device") {
		analysis += "This indicates disk space exhaustion. Free up disk space and try again.\n"
	}

	return analysis
}

// ExportDebugData exports debug data as JSON
func (d *DebugCollector) ExportDebugData() (string, error) {
	data := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      d.info,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal debug data: %v", err)
	}

	return string(jsonData), nil
}

// Reset clears all collected debug information
func (d *DebugCollector) Reset() {
	d.info = make(map[string]interface{})
}

// appendMapToReport appends a map to the report
func (d *DebugCollector) appendMapToReport(report *strings.Builder, data map[string]interface{}) {
	for key, value := range data {
		report.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}
}
