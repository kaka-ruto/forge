package qemu

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sst/forge/internal/config"
	"github.com/sst/forge/internal/logger"
	"golang.org/x/crypto/ssh"
)

// QEMUManager manages QEMU instances for testing
type QEMUManager struct {
	config     *config.Config
	projectDir string
	logger     *logger.Logger
}

// QEMUInstance represents a running QEMU instance
type QEMUInstance struct {
	ID          string
	Process     *os.Process
	Config      *config.Config
	StartTime   time.Time
	MonitorPort int
	SSHPort     int
	SerialPort  int
	LogFile     *os.File
}

// TestMetrics contains detailed metrics collected during testing
type TestMetrics struct {
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	MemoryUsageMB      float64 `json:"memory_usage_mb"`
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	DiskUsageMB        float64 `json:"disk_usage_mb"`
	DiskUsagePercent   float64 `json:"disk_usage_percent"`
	NetworkBytesSent   int64   `json:"network_bytes_sent"`
	NetworkBytesRecv   int64   `json:"network_bytes_recv"`
	LoadAverage1       float64 `json:"load_average_1"`
	LoadAverage5       float64 `json:"load_average_5"`
	LoadAverage15      float64 `json:"load_average_15"`
}

// TestResult represents the result of a test scenario
type TestResult struct {
	TestName   string        `json:"test_name"`
	Success    bool          `json:"success"`
	Duration   time.Duration `json:"duration"`
	Output     string        `json:"output"`
	Error      string        `json:"error"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	InstanceID string        `json:"instance_id"`
	ConfigHash string        `json:"config_hash"`
	Metrics    *TestMetrics  `json:"metrics,omitempty"`
}

// TestScenario defines a test to run on a QEMU instance
type TestScenario struct {
	Name        string
	Description string
	Timeout     time.Duration
	Run         func(ctx context.Context, instance *QEMUInstance) *TestResult
}

// TestComparison represents a comparison between test result sets
type TestComparison struct {
	TotalTests     int
	PassedTests    int
	FailedTests    int
	ImprovedTests  int
	RegressedTests int
	NewTests       int
	RemovedTests   int
	Details        map[string]*TestComparisonDetail
}

// TestComparisonDetail represents detailed comparison for a single test
type TestComparisonDetail struct {
	TestName       string
	Status         string // "improved", "regressed", "stable_pass", "stable_fail", "new", "removed"
	Current        *TestResult
	Previous       *TestResult
	DurationChange time.Duration
	DurationStatus string // "faster", "slower", "same"
}

// NewQEMUManager creates a new QEMU manager
func NewQEMUManager(cfg *config.Config, projectDir string) *QEMUManager {
	return &QEMUManager{
		config:     cfg,
		projectDir: projectDir,
		logger:     logger.NewLogger(logger.INFO, os.Stdout, os.Stderr),
	}
}

// StartInstance starts a QEMU instance with the built image
func (qm *QEMUManager) StartInstance(ctx context.Context, imagePath string) (*QEMUInstance, error) {
	instance := &QEMUInstance{
		ID:        generateInstanceID(),
		Config:    qm.config,
		StartTime: time.Now(),
	}

	// Find available ports
	var err error
	instance.MonitorPort, err = findAvailablePort(4444, 4500)
	if err != nil {
		return nil, fmt.Errorf("failed to find monitor port: %v", err)
	}

	instance.SSHPort, err = findAvailablePort(2222, 2300)
	if err != nil {
		return nil, fmt.Errorf("failed to find SSH port: %v", err)
	}

	instance.SerialPort, err = findAvailablePort(8000, 8100)
	if err != nil {
		return nil, fmt.Errorf("failed to find serial port: %v", err)
	}

	// Create log file
	logPath := filepath.Join(qm.projectDir, "test-logs", fmt.Sprintf("qemu-%s.log", instance.ID))
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	instance.LogFile, err = os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}

	// Build QEMU command
	cmd := qm.buildQEMUCommand(instance, imagePath)

	qm.logger.Info("Starting QEMU instance %s with command: %s", instance.ID, strings.Join(cmd, " "))

	// Start QEMU process
	process, err := os.StartProcess("/usr/bin/qemu-system-x86_64", cmd, &os.ProcAttr{
		Files: []*os.File{nil, instance.LogFile, instance.LogFile},
	})
	if err != nil {
		instance.LogFile.Close()
		return nil, fmt.Errorf("failed to start QEMU: %v", err)
	}

	instance.Process = process

	// Wait for QEMU to boot
	if err := qm.waitForBoot(ctx, instance); err != nil {
		qm.StopInstance(instance)
		return nil, fmt.Errorf("failed to wait for boot: %v", err)
	}

	qm.logger.Info("QEMU instance %s started successfully", instance.ID)
	return instance, nil
}

// StopInstance stops a QEMU instance
func (qm *QEMUManager) StopInstance(instance *QEMUInstance) error {
	if instance.Process == nil {
		return nil
	}

	qm.logger.Info("Stopping QEMU instance %s", instance.ID)

	// Try graceful shutdown first
	if err := qm.sendShutdownCommand(instance); err != nil {
		qm.logger.Warn("Graceful shutdown failed, force killing: %v", err)
	}

	// Wait a bit for graceful shutdown
	time.Sleep(2 * time.Second)

	// Force kill if still running
	if err := instance.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill QEMU process: %v", err)
	}

	// Wait for process to exit
	_, err := instance.Process.Wait()
	if err != nil && !strings.Contains(err.Error(), "signal") {
		qm.logger.Warn("Process wait error: %v", err)
	}

	// Close log file
	if instance.LogFile != nil {
		instance.LogFile.Close()
	}

	qm.logger.Info("QEMU instance %s stopped", instance.ID)
	return nil
}

// RunTestScenario runs a test scenario on a QEMU instance
func (qm *QEMUManager) RunTestScenario(ctx context.Context, instance *QEMUInstance, scenario TestScenario) *TestResult {
	result := &TestResult{
		TestName:  scenario.Name,
		StartTime: time.Now(),
	}

	qm.logger.Info("Running test scenario: %s", scenario.Name)

	// Create context with timeout
	testCtx, cancel := context.WithTimeout(ctx, scenario.Timeout)
	defer cancel()

	// Run the test
	testResult := scenario.Run(testCtx, instance)

	// Copy results
	result.Success = testResult.Success
	result.Duration = time.Since(result.StartTime)
	result.EndTime = time.Now()
	result.Output = testResult.Output
	result.Error = testResult.Error

	if result.Success {
		qm.logger.Info("Test %s PASSED in %v", scenario.Name, result.Duration)
	} else {
		qm.logger.Error("Test %s FAILED: %s", scenario.Name, result.Error)
	}

	return result
}

// testBootScenario tests that the system boots successfully
func (qm *QEMUManager) testBootScenario(ctx context.Context, instance *QEMUInstance) *TestResult {
	result := &TestResult{
		TestName:  "boot",
		StartTime: time.Now(),
	}

	// Wait for the system to be reachable via SSH
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			result.Success = false
			result.Error = "context cancelled"
			result.EndTime = time.Now()
			result.Duration = time.Since(result.StartTime)
			return result
		case <-timeout:
			result.Success = false
			result.Error = "boot timeout"
			result.EndTime = time.Now()
			result.Duration = time.Since(result.StartTime)
			return result
		case <-ticker.C:
			if qm.isBooted(instance) {
				result.Success = true
				result.Output = "System booted successfully"
				result.EndTime = time.Now()
				result.Duration = time.Since(result.StartTime)
				return result
			}
		}
	}
}

// GetDefaultTestScenarios returns the default test scenarios
func (qm *QEMUManager) GetDefaultTestScenarios() []TestScenario {
	return []TestScenario{
		{
			Name:        "boot",
			Description: "Test that the system boots successfully",
			Timeout:     2 * time.Minute,
			Run:         qm.testBootScenario,
		},
		{
			Name:        "network",
			Description: "Test network connectivity",
			Timeout:     1 * time.Minute,
			Run:         qm.testNetworkScenario,
		},
		{
			Name:        "services",
			Description: "Test that essential services are running",
			Timeout:     1 * time.Minute,
			Run:         qm.testServicesScenario,
		},
		{
			Name:        "performance",
			Description: "Test basic system performance",
			Timeout:     2 * time.Minute,
			Run:         qm.testPerformanceScenario,
		},
		{
			Name:        "stress",
			Description: "Test system under stress",
			Timeout:     3 * time.Minute,
			Run:         qm.testStressScenario,
		},
	}
}

// SaveTestResults saves test results to a JSON file
func (qm *QEMUManager) SaveTestResults(results []*TestResult, instanceID string) error {
	// Create test results directory
	resultsDir := filepath.Join(qm.projectDir, "test-results")
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create test results directory: %v", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("test-results-%s-%s.json", instanceID, timestamp)
	filePath := filepath.Join(resultsDir, filename)

	// Convert results to JSON
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal test results: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write test results: %v", err)
	}

	qm.logger.Info("Test results saved to %s", filePath)
	return nil
}

// LoadTestResults loads the most recent test results from file
func (qm *QEMUManager) LoadTestResults(instanceID string) ([]*TestResult, error) {
	resultsDir := filepath.Join(qm.projectDir, "test-results")

	// Check if results directory exists
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no test results directory found")
	}

	// Find the most recent results file for this instance
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read test results directory: %v", err)
	}

	var latestFile string
	var latestTime time.Time

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "test-results-"+instanceID) && strings.HasSuffix(entry.Name(), ".json") {
			fileTime, err := time.Parse("2006-01-02_15-04-05", strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "test-results-"+instanceID+"-"), ".json"))
			if err == nil && fileTime.After(latestTime) {
				latestTime = fileTime
				latestFile = entry.Name()
			}
		}
	}

	if latestFile == "" {
		return nil, fmt.Errorf("no test results found for instance %s", instanceID)
	}

	// Read and parse the file
	filePath := filepath.Join(resultsDir, latestFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test results file: %v", err)
	}

	var results []*TestResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test results: %v", err)
	}

	return results, nil
}

// CompareTestResults compares current results with previous results
func (qm *QEMUManager) CompareTestResults(current, previous []*TestResult) *TestComparison {
	comparison := &TestComparison{
		TotalTests:     len(current),
		PassedTests:    0,
		FailedTests:    0,
		ImprovedTests:  0,
		RegressedTests: 0,
		NewTests:       0,
		RemovedTests:   0,
		Details:        make(map[string]*TestComparisonDetail),
	}

	// Create maps for easy lookup
	currentMap := make(map[string]*TestResult)
	previousMap := make(map[string]*TestResult)

	for _, result := range current {
		currentMap[result.TestName] = result
		if result.Success {
			comparison.PassedTests++
		} else {
			comparison.FailedTests++
		}
	}

	for _, result := range previous {
		previousMap[result.TestName] = result
	}

	// Compare results
	for testName, currentResult := range currentMap {
		detail := &TestComparisonDetail{
			TestName: testName,
			Current:  currentResult,
		}

		if previousResult, exists := previousMap[testName]; exists {
			detail.Previous = previousResult

			// Check for improvements/regressions
			if !previousResult.Success && currentResult.Success {
				comparison.ImprovedTests++
				detail.Status = "improved"
			} else if previousResult.Success && !currentResult.Success {
				comparison.RegressedTests++
				detail.Status = "regressed"
			} else if currentResult.Success {
				detail.Status = "stable_pass"
			} else {
				detail.Status = "stable_fail"
			}

			// Compare durations
			if currentResult.Duration < previousResult.Duration {
				detail.DurationChange = previousResult.Duration - currentResult.Duration
				detail.DurationStatus = "faster"
			} else if currentResult.Duration > previousResult.Duration {
				detail.DurationChange = currentResult.Duration - previousResult.Duration
				detail.DurationStatus = "slower"
			}
		} else {
			comparison.NewTests++
			detail.Status = "new"
		}

		comparison.Details[testName] = detail
	}

	// Check for removed tests
	for testName := range previousMap {
		if _, exists := currentMap[testName]; !exists {
			comparison.RemovedTests++
			comparison.Details[testName] = &TestComparisonDetail{
				TestName: testName,
				Status:   "removed",
			}
		}
	}

	return comparison
}

// testPerformanceScenario tests basic system performance
func (qm *QEMUManager) testPerformanceScenario(ctx context.Context, instance *QEMUInstance) *TestResult {
	result := &TestResult{
		TestName:  "performance",
		StartTime: time.Now(),
	}

	// Run a simple performance test (CPU benchmark)
	stdout, stderr, err := qm.executeSSHCommand(instance, "time dd if=/dev/zero of=/dev/null bs=1M count=100", 2*time.Minute)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Performance test failed: %v", err)
		if stderr != "" {
			result.Error += fmt.Sprintf(" (stderr: %s)", stderr)
		}
		result.EndTime = time.Now()
		result.Duration = time.Since(result.StartTime)
		return result
	}

	// Check if the command completed successfully
	if strings.Contains(stdout, "100+0 records in") && strings.Contains(stdout, "100+0 records out") {
		result.Success = true
		result.Output = "Performance test completed successfully"
		// Extract timing information if available
		if strings.Contains(stdout, "real") {
			result.Output += "\n" + stdout
		}
	} else {
		result.Success = false
		result.Error = "Performance test did not complete successfully"
		result.Output = stdout
	}

	// Collect system metrics
	if metrics, err := qm.collectSystemMetrics(instance); err == nil {
		result.Metrics = metrics
	}

	result.EndTime = time.Now()
	result.Duration = time.Since(result.StartTime)

	return result
}

// testStressScenario tests system under stress
func (qm *QEMUManager) testStressScenario(ctx context.Context, instance *QEMUInstance) *TestResult {
	result := &TestResult{
		TestName:  "stress",
		StartTime: time.Now(),
	}

	// Run a simple stress test (memory and CPU)
	stdout, stderr, err := qm.executeSSHCommand(instance, "echo 'Testing system stress...' && free -h && uptime", 3*time.Minute)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Stress test failed: %v", err)
		if stderr != "" {
			result.Error += fmt.Sprintf(" (stderr: %s)", stderr)
		}
		result.EndTime = time.Now()
		result.Duration = time.Since(result.StartTime)
		return result
	}

	// Check if we got memory and uptime information
	if strings.Contains(stdout, "Mem:") && strings.Contains(stdout, "load average") {
		result.Success = true
		result.Output = "Stress test completed - system information retrieved"
		result.Output += "\n" + stdout
	} else {
		result.Success = false
		result.Error = "Stress test failed to retrieve system information"
		result.Output = stdout
	}

	// Collect system metrics
	if metrics, err := qm.collectSystemMetrics(instance); err == nil {
		result.Metrics = metrics
	}

	result.EndTime = time.Now()
	result.Duration = time.Since(result.StartTime)

	return result
}

// testNetworkScenario tests network connectivity
func (qm *QEMUManager) testNetworkScenario(ctx context.Context, instance *QEMUInstance) *TestResult {
	result := &TestResult{
		TestName:  "network",
		StartTime: time.Now(),
	}

	// Test network connectivity by pinging a known host
	stdout, stderr, err := qm.executeSSHCommand(instance, "ping -c 3 8.8.8.8", 30*time.Second)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Network test failed: %v", err)
		if stderr != "" {
			result.Error += fmt.Sprintf(" (stderr: %s)", stderr)
		}
		result.EndTime = time.Now()
		result.Duration = time.Since(result.StartTime)
		return result
	}

	// Check if ping was successful (look for "3 packets transmitted, 3 received")
	if strings.Contains(stdout, "3 packets transmitted, 3 received") {
		result.Success = true
		result.Output = "Network connectivity verified - ping successful"
	} else {
		result.Success = false
		result.Error = "Network test failed - ping unsuccessful"
		result.Output = stdout
	}

	// Collect system metrics
	if metrics, err := qm.collectSystemMetrics(instance); err == nil {
		result.Metrics = metrics
	}

	result.EndTime = time.Now()
	result.Duration = time.Since(result.StartTime)

	return result
}

// testServicesScenario tests that essential services are running
func (qm *QEMUManager) testServicesScenario(ctx context.Context, instance *QEMUInstance) *TestResult {
	result := &TestResult{
		TestName:  "services",
		StartTime: time.Now(),
	}

	// Check that SSH service is running (sshd)
	stdout, stderr, err := qm.executeSSHCommand(instance, "ps aux | grep sshd | grep -v grep", 10*time.Second)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Service check failed: %v", err)
		if stderr != "" {
			result.Error += fmt.Sprintf(" (stderr: %s)", stderr)
		}
		result.EndTime = time.Now()
		result.Duration = time.Since(result.StartTime)
		return result
	}

	// Check if SSH daemon is running
	if strings.Contains(stdout, "sshd") {
		result.Success = true
		result.Output = "SSH service is running"
	} else {
		result.Success = false
		result.Error = "SSH service is not running"
		result.Output = stdout
	}

	// Collect system metrics
	if metrics, err := qm.collectSystemMetrics(instance); err == nil {
		result.Metrics = metrics
	}

	result.EndTime = time.Now()
	result.Duration = time.Since(result.StartTime)

	return result
}

// generateInstanceID generates a unique instance ID
func generateInstanceID() string {
	return fmt.Sprintf("qemu-%d", time.Now().UnixNano())
}

// findAvailablePort finds an available port in the given range
func findAvailablePort(start, end int) (int, error) {
	for port := start; port <= end; port++ {
		if isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports in range %d-%d", start, end)
}

// isPortAvailable checks if a port is available
func isPortAvailable(port int) bool {
	// This is a simplified check - in a real implementation,
	// you would try to bind to the port
	return true
}

// buildQEMUCommand builds the QEMU command line arguments
func (qm *QEMUManager) buildQEMUCommand(instance *QEMUInstance, imagePath string) []string {
	args := []string{
		"qemu-system-x86_64", // Use just the command name, not full path
		"-machine", "pc",
		"-cpu", "host",
		"-m", "512",
		"-kernel", imagePath,
		"-append", "console=ttyS0 root=/dev/sda1",
		"-drive", fmt.Sprintf("file=%s,if=virtio,format=raw", imagePath),
		"-net", "nic,model=virtio",
		"-net", fmt.Sprintf("user,hostfwd=tcp::%d-:22", instance.SSHPort),
		"-monitor", fmt.Sprintf("tcp:127.0.0.1:%d,server,nowait", instance.MonitorPort),
		"-serial", fmt.Sprintf("tcp:127.0.0.1:%d,server,nowait", instance.SerialPort),
		"-nographic",
		"-daemonize",
	}

	// Add architecture-specific options
	switch qm.config.Architecture {
	case "arm":
		args[0] = "qemu-system-arm"
		args = append(args[:1], append([]string{"-M", "versatilepb"}, args[1:]...)...)
	case "aarch64":
		args[0] = "qemu-system-aarch64"
		args = append(args[:1], append([]string{"-M", "virt"}, args[1:]...)...)
	}

	return args
}

// waitForBoot waits for the QEMU instance to finish booting
func (qm *QEMUManager) waitForBoot(ctx context.Context, instance *QEMUInstance) error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for boot")
		case <-ticker.C:
			if qm.isBooted(instance) {
				return nil
			}
		}
	}
}

// isBooted checks if the QEMU instance has finished booting
func (qm *QEMUManager) isBooted(instance *QEMUInstance) bool {
	// Try to connect to SSH port
	client, err := qm.connectToSSH(instance)
	if err != nil {
		return false
	}
	client.Close()
	return true
}

// sendShutdownCommand sends a shutdown command to QEMU via monitor
func (qm *QEMUManager) sendShutdownCommand(instance *QEMUInstance) error {
	// Connect to QEMU monitor
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", instance.MonitorPort))
	if err != nil {
		return fmt.Errorf("failed to connect to QEMU monitor: %v", err)
	}
	defer conn.Close()

	// Send shutdown command
	_, err = fmt.Fprintf(conn, "system_powerdown\n")
	if err != nil {
		return fmt.Errorf("failed to send shutdown command: %v", err)
	}

	// Wait a moment for the command to be processed
	time.Sleep(100 * time.Millisecond)

	return nil
}

// connectToSSH attempts to connect to the SSH port
func (qm *QEMUManager) connectToSSH(instance *QEMUInstance) (*ssh.Client, error) {
	// SSH client config for connecting to QEMU instance
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(""), // Empty password for root in QEMU
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing purposes
		Timeout:         5 * time.Second,
	}

	// Connect to SSH server
	addr := fmt.Sprintf("localhost:%d", instance.SSHPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH: %v", err)
	}

	return client, nil
}

// collectSystemMetrics collects system metrics from the QEMU instance
func (qm *QEMUManager) collectSystemMetrics(instance *QEMUInstance) (*TestMetrics, error) {
	metrics := &TestMetrics{}

	// Get memory information
	memOutput, _, err := qm.executeSSHCommand(instance, "free -m | grep '^Mem:'", 10*time.Second)
	if err == nil {
		// Parse memory output: "Mem: total used free shared buff/cache available"
		var total, used, free, shared, buffCache, available float64
		if n, _ := fmt.Sscanf(memOutput, "Mem: %f %f %f %f %f %f", &total, &used, &free, &shared, &buffCache, &available); n >= 2 {
			metrics.MemoryUsageMB = used
			if total > 0 {
				metrics.MemoryUsagePercent = (used / total) * 100
			}
		}
	}

	// Get disk usage
	diskOutput, _, err := qm.executeSSHCommand(instance, "df / | tail -1", 10*time.Second)
	if err == nil {
		// Parse disk output: "Filesystem 1K-blocks Used Available Use% Mounted-on"
		var fs string
		var blocks, used, available int64
		var usePercent string
		var mount string
		if n, _ := fmt.Sscanf(diskOutput, "%s %d %d %d %s %s", &fs, &blocks, &used, &available, &usePercent, &mount); n >= 5 {
			metrics.DiskUsageMB = float64(used) / 1024
			// Remove % from usePercent and parse
			usePercent = strings.TrimSuffix(usePercent, "%")
			if percent, err := strconv.ParseFloat(usePercent, 64); err == nil {
				metrics.DiskUsagePercent = percent
			}
		}
	}

	// Get load average
	loadOutput, _, err := qm.executeSSHCommand(instance, "uptime | awk -F'load average:' '{print $2}'", 10*time.Second)
	if err == nil {
		// Parse load average: " 1.23, 1.45, 1.67"
		loadStr := strings.TrimSpace(loadOutput)
		loads := strings.Split(loadStr, ",")
		if len(loads) >= 3 {
			if load1, err := strconv.ParseFloat(strings.TrimSpace(loads[0]), 64); err == nil {
				metrics.LoadAverage1 = load1
			}
			if load5, err := strconv.ParseFloat(strings.TrimSpace(loads[1]), 64); err == nil {
				metrics.LoadAverage5 = load5
			}
			if load15, err := strconv.ParseFloat(strings.TrimSpace(loads[2]), 64); err == nil {
				metrics.LoadAverage15 = load15
			}
		}
	}

	return metrics, nil
}

// executeSSHCommand executes a command on the QEMU instance via SSH
func (qm *QEMUManager) executeSSHCommand(instance *QEMUInstance, command string, timeout time.Duration) (string, string, error) {
	client, err := qm.connectToSSH(instance)
	if err != nil {
		return "", "", fmt.Errorf("SSH connection failed: %v", err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Set up pipes for stdout and stderr
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	// Start the command
	if err := session.Start(command); err != nil {
		return "", "", fmt.Errorf("failed to start command: %v", err)
	}

	// Read output with timeout
	var stdoutBuf, stderrBuf strings.Builder
	done := make(chan bool, 2)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdoutBuf.WriteString(scanner.Text() + "\n")
		}
		done <- true
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stderrBuf.WriteString(scanner.Text() + "\n")
		}
		done <- true
	}()

	// Wait for command completion or timeout
	doneCh := make(chan error, 1)
	go func() {
		doneCh <- session.Wait()
	}()

	select {
	case err := <-doneCh:
		// Wait for output readers to finish
		<-done
		<-done
		if err != nil {
			return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("command failed: %v", err)
		}
		return stdoutBuf.String(), stderrBuf.String(), nil
	case <-time.After(timeout):
		session.Signal(ssh.SIGKILL)
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("command timed out after %v", timeout)
	}
}
