package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// CleanupConfig holds configuration for cleanup operations
type CleanupConfig struct {
	Ports           []int
	ProcessPatterns []string
	DockerServices  []string
	Timeout         time.Duration
	Verbose         bool
}

// getDefaultCleanupConfig returns default cleanup configuration
func getDefaultCleanupConfig() CleanupConfig {
	return CleanupConfig{
		Ports: []int{8081, 3000, 7474, 7687, 9090, 9091, 9092, 9093},
		ProcessPatterns: []string{
			// Exclude the current process ID to avoid self-termination
			"overseer(?!.*" + strconv.Itoa(os.Getpid()) + ")",
			"overseer api(?!.*" + strconv.Itoa(os.Getpid()) + ")",
			"go run . api",
			"bun.*server",
			"bun.*dev",
			"node.*server",
			"vite",
		},
		DockerServices: []string{
			// Docker services cleanup completely disabled
		},
		Timeout: 30 * time.Second,
		Verbose: false,
	}
}

// performCleanStartup performs comprehensive cleanup for fresh startup
func performCleanStartup(config CleanupConfig) error {
	printVerboseMessage(config.Verbose, "🧹 Starting comprehensive cleanup for fresh environment...")

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	errors := executeCleanupSteps(ctx, config)
	logCleanupErrors(errors, config.Verbose)

	printVerboseMessage(config.Verbose, "✅ Cleanup completed successfully!")
	return nil
}

func executeCleanupSteps(ctx context.Context, config CleanupConfig) []error {
	var errors []error

	steps := []func() error{
		func() error { return killProcessesByPattern(ctx, config.ProcessPatterns, config.Verbose) },
		func() error { return killProcessesByPort(ctx, config.Ports, config.Verbose) },
		func() error { return stopDockerServices(ctx, config.DockerServices, config.Verbose) },
		func() error { return cleanupTempFiles(config.Verbose) },
		func() error { return waitForPortsFree(ctx, config.Ports, config.Verbose) },
	}

	stepNames := []string{
		"kill processes",
		"kill processes by port",
		"stop Docker services",
		"cleanup temp files",
		"wait for ports to be free",
	}

	for i, step := range steps {
		if err := step(); err != nil {
			errors = append(errors, fmt.Errorf("failed to %s: %w", stepNames[i], err))
		}
	}

	return errors
}

func logCleanupErrors(errors []error, verbose bool) {
	if len(errors) == 0 {
		return
	}

	for _, err := range errors {
		if verbose {
			fmt.Printf("⚠️  Warning: %v\n", err)
		}
	}
}

func printVerboseMessage(verbose bool, message string) {
	if verbose {
		fmt.Println(message)
	}
}

// killProcessesByPattern kills processes matching given patterns
func killProcessesByPattern(ctx context.Context, patterns []string, verbose bool) error {
	for _, pattern := range patterns {
		if err := killProcessesByPatternSingle(ctx, pattern, verbose); err != nil {
			printVerboseMessage(verbose, fmt.Sprintf("   Error killing processes for pattern %s: %v", pattern, err))
		}
	}
	return nil
}

func killProcessesByPatternSingle(ctx context.Context, pattern string, verbose bool) error {
	printVerboseMessage(verbose, fmt.Sprintf("🔍 Killing processes matching pattern: %s", pattern))

	if shouldSkipPattern(pattern) {
		return nil
	}

	cmd := exec.CommandContext(ctx, "pkill", "-f", pattern)
	output, err := cmd.CombinedOutput()

	return handleKillCommandResult(cmd, output, err, pattern, verbose)
}

func shouldSkipPattern(pattern string) bool {
	return strings.Contains(pattern, strconv.Itoa(os.Getpid()))
}

func handleKillCommandResult(cmd *exec.Cmd, output []byte, err error, pattern string, verbose bool) error {
	if err != nil {
		if cmd.ProcessState.ExitCode() == 1 {
			printVerboseMessage(verbose, fmt.Sprintf("   No processes found for pattern: %s", pattern))
			return nil
		}
		return err
	}

	if verbose && len(output) > 0 {
		fmt.Printf("   Killed processes: %s\n", strings.TrimSpace(string(output)))
	}

	return nil
}

// killProcessesByPort kills processes using specific ports
func killProcessesByPort(ctx context.Context, ports []int, verbose bool) error {
	for _, port := range ports {
		if err := killProcessesByPortSingle(ctx, port, verbose); err != nil {
			printVerboseMessage(verbose, fmt.Sprintf("   Error killing processes on port %d: %v", port, err))
		}
	}
	return nil
}

func killProcessesByPortSingle(ctx context.Context, port int, verbose bool) error {
	printVerboseMessage(verbose, fmt.Sprintf("🔍 Checking port %d for processes", port))

	pids, err := findProcessesByPort(ctx, port, verbose)
	if err != nil {
		return nil // No processes found is OK
	}

	return killProcessesByPids(pids, port, verbose)
}

func findProcessesByPort(ctx context.Context, port int, verbose bool) ([]int, error) {
	cmd := exec.CommandContext(ctx, "lsof", "-ti", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()

	if err != nil {
		printVerboseMessage(verbose, fmt.Sprintf("   No processes found on port %d", port))
		return nil, err
	}

	return parseProcessIds(string(output))
}

func parseProcessIds(output string) ([]int, error) {
	pidStrings := strings.Fields(output)
	var pids []int

	for _, pidStr := range pidStrings {
		pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}

	return pids, nil
}

func killProcessesByPids(pids []int, port int, verbose bool) error {
	for _, pid := range pids {
		if err := killProcessByPid(pid, port, verbose); err != nil {
			printVerboseMessage(verbose, fmt.Sprintf("   Error killing PID %d: %v", pid, err))
		}
	}
	return nil
}

func killProcessByPid(pid, port int, verbose bool) error {
	printVerboseMessage(verbose, fmt.Sprintf("   Killing process %d on port %d", pid, port))

	if err := syscallKillProcess(pid, false); err != nil {
		printVerboseMessage(verbose, fmt.Sprintf("   Graceful kill failed for PID %d, trying force kill", pid))
		return syscallKillProcess(pid, true)
	}

	return nil
}

// syscallKillProcess kills a process by PID
func syscallKillProcess(pid int, force bool) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if force {
		return process.Kill()
	}

	// Try graceful termination
	return process.Signal(os.Interrupt)
}

// stopDockerServices stops Docker services
func stopDockerServices(ctx context.Context, services []string, verbose bool) error {
	if verbose {
		fmt.Println("🐳 Docker service cleanup disabled")
	}

	// Completely skip Docker cleanup to ensure Neo4j stays running
	return nil
}


// cleanupTempFiles removes temporary files and directories
func cleanupTempFiles(verbose bool) error {
	if verbose {
		fmt.Println("🗑️  Cleaning up temporary files...")
	}

	tempPaths := []string{
		"coverage.out",
		"coverage.html",
		"gosec-report.json",
		"bin/",
		"ui/node_modules/.cache/",
		"ui/dist/",
		"ui/coverage/",
		"ui/stryker-tmp/",
		"test_api.go", // Clean up our test file
	}

	for _, path := range tempPaths {
		if verbose {
			fmt.Printf("   Removing: %s\n", path)
		}

		if err := os.RemoveAll(path); err != nil {
			if verbose {
				fmt.Printf("   Warning: Failed to remove %s: %v\n", path, err)
			}
		}
	}

	return nil
}

// waitForPortsFree waits for ports to become available
func waitForPortsFree(ctx context.Context, ports []int, verbose bool) error {
	if verbose {
		fmt.Println("⏳ Waiting for ports to become free...")
	}

	for _, port := range ports {
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("timeout waiting for port %d to become free", port)
			default:
			}

			// Check if port is still in use
			cmd := exec.CommandContext(ctx, "lsof", "-ti", fmt.Sprintf(":%d", port))
			err := cmd.Run()

			if err != nil {
				// Port is free
				if verbose {
					fmt.Printf("   Port %d is now free\n", port)
				}
				break
			}

			// Wait a bit before checking again
			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}


// Emergency cleanup function that can be called from main
func emergencyCleanup() {
	fmt.Println("🚨 Emergency cleanup triggered...")

	config := getDefaultCleanupConfig()
	config.Verbose = true
	config.Timeout = 10 * time.Second // Shorter timeout for emergency
	// No Docker services to stop, even in emergency cleanup
	config.DockerServices = []string{}

	if err := performCleanStartup(config); err != nil {
		log.Printf("Emergency cleanup failed: %v", err)
	}
}

