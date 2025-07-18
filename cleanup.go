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
	if config.Verbose {
		fmt.Println("🧹 Starting comprehensive cleanup for fresh environment...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	var errors []error

	// 1. Kill processes by pattern
	if err := killProcessesByPattern(ctx, config.ProcessPatterns, config.Verbose); err != nil {
		errors = append(errors, fmt.Errorf("failed to kill processes: %w", err))
	}

	// 2. Kill processes by port
	if err := killProcessesByPort(ctx, config.Ports, config.Verbose); err != nil {
		errors = append(errors, fmt.Errorf("failed to kill processes by port: %w", err))
	}

	// 3. Stop Docker services
	if err := stopDockerServices(ctx, config.DockerServices, config.Verbose); err != nil {
		errors = append(errors, fmt.Errorf("failed to stop Docker services: %w", err))
	}

	// 4. Clean up temporary files
	if err := cleanupTempFiles(config.Verbose); err != nil {
		errors = append(errors, fmt.Errorf("failed to cleanup temp files: %w", err))
	}

	// 5. Wait for ports to be free
	if err := waitForPortsFree(ctx, config.Ports, config.Verbose); err != nil {
		errors = append(errors, fmt.Errorf("ports still occupied: %w", err))
	}

	if len(errors) > 0 {
		for _, err := range errors {
			if config.Verbose {
				fmt.Printf("⚠️  Warning: %v\n", err)
			}
		}
		// Don't return error for warnings - continue with startup
	}

	if config.Verbose {
		fmt.Println("✅ Cleanup completed successfully!")
	}

	return nil
}

// killProcessesByPattern kills processes matching given patterns
func killProcessesByPattern(ctx context.Context, patterns []string, verbose bool) error {
	for _, pattern := range patterns {
		if verbose {
			fmt.Printf("🔍 Killing processes matching pattern: %s\n", pattern)
		}

		// Skip if pattern is for the current process
		if strings.Contains(pattern, strconv.Itoa(os.Getpid())) {
			continue
		}

		// Use pkill with pattern matching
		cmd := exec.CommandContext(ctx, "pkill", "-f", pattern)
		output, err := cmd.CombinedOutput()

		if err != nil {
			// pkill returns exit code 1 if no processes found - this is OK
			if cmd.ProcessState.ExitCode() == 1 {
				if verbose {
					fmt.Printf("   No processes found for pattern: %s\n", pattern)
				}
				continue
			}
			if verbose {
				fmt.Printf("   Error killing processes for pattern %s: %v\n", pattern, err)
			}
		} else if verbose && len(output) > 0 {
			fmt.Printf("   Killed processes: %s\n", strings.TrimSpace(string(output)))
		}
	}

	return nil
}

// killProcessesByPort kills processes using specific ports
func killProcessesByPort(ctx context.Context, ports []int, verbose bool) error {
	for _, port := range ports {
		if verbose {
			fmt.Printf("🔍 Checking port %d for processes\n", port)
		}

		// Find processes using the port
		cmd := exec.CommandContext(ctx, "lsof", "-ti", fmt.Sprintf(":%d", port))
		output, err := cmd.Output()

		if err != nil {
			// lsof returns exit code 1 if no processes found - this is OK
			if verbose {
				fmt.Printf("   No processes found on port %d\n", port)
			}
			continue
		}

		// Kill each PID found
		pids := strings.Fields(string(output))
		for _, pidStr := range pids {
			pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
			if err != nil {
				continue
			}

			if verbose {
				fmt.Printf("   Killing process %d on port %d\n", pid, port)
			}

			// Try graceful termination first
			if err := syscallKillProcess(pid, false); err != nil {
				if verbose {
					fmt.Printf("   Graceful kill failed for PID %d, trying force kill\n", pid)
				}
				// Force kill if graceful fails
				_ = syscallKillProcess(pid, true)
			}
		}
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

// isDockerAvailable checks if Docker is available
func isDockerAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "docker", "version")
	return cmd.Run() == nil
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
		"ui/mutation-report.html",
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

// initCleanStartup initializes cleanup before starting services
func initCleanStartup(verbose bool) error {
	config := getDefaultCleanupConfig()
	config.Verbose = verbose

	return performCleanStartup(config)
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

// Graceful shutdown cleanup
func gracefulShutdown() {
	fmt.Println("🛑 Graceful shutdown initiated...")

	config := getDefaultCleanupConfig()
	config.Verbose = true

	if err := performCleanStartup(config); err != nil {
		log.Printf("Graceful shutdown cleanup failed: %v", err)
	}
}
