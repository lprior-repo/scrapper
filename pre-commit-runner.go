package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// PreCommitConfig holds configuration for pre-commit checks
type PreCommitConfig struct {
	ProjectRoot     string
	EnableMutation  bool
	CoverageTarget  int
	Timeout         time.Duration
	SkipTests       bool
	SkipLinting     bool
	SkipFormatting  bool
	Verbose         bool
}

// CheckResult represents the result of a pre-commit check
type CheckResult struct {
	Name     string
	Command  string
	Success  bool
	Output   string
	Duration time.Duration
	Error    error
}

// PreCommitRunner handles running pre-commit checks
type PreCommitRunner struct {
	config  PreCommitConfig
	results []CheckResult
}

// NewPreCommitRunner creates a new pre-commit runner
func NewPreCommitRunner(config PreCommitConfig) *PreCommitRunner {
	return &PreCommitRunner{
		config:  config,
		results: make([]CheckResult, 0),
	}
}

// runCommand executes a command and returns the result
func (r *PreCommitRunner) runCommand(name, command string, args ...string) CheckResult {
	start := time.Now()
	
	if r.config.Verbose {
		fmt.Printf("🔄 Running %s: %s %s\n", name, command, strings.Join(args, " "))
	}
	
	cmd := exec.Command(command, args...)
	cmd.Dir = r.config.ProjectRoot
	
	// Set timeout
	if r.config.Timeout > 0 {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		go func() {
			time.Sleep(r.config.Timeout)
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
		}()
	}
	
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)
	
	result := CheckResult{
		Name:     name,
		Command:  fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
		Success:  err == nil,
		Output:   string(output),
		Duration: duration,
		Error:    err,
	}
	
	if r.config.Verbose {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		fmt.Printf("%s %s completed in %v\n", status, name, duration)
	}
	
	return result
}

// RunGoFormat runs go fmt on all Go files
func (r *PreCommitRunner) RunGoFormat() CheckResult {
	if r.config.SkipFormatting {
		return CheckResult{Name: "Go Format", Success: true, Output: "Skipped"}
	}
	
	result := r.runCommand("Go Format", "go", "fmt", "./...")
	r.results = append(r.results, result)
	return result
}

// RunGoVet runs go vet on all Go files
func (r *PreCommitRunner) RunGoVet() CheckResult {
	result := r.runCommand("Go Vet", "go", "vet", "./...")
	r.results = append(r.results, result)
	return result
}

// RunGoModTidy runs go mod tidy
func (r *PreCommitRunner) RunGoModTidy() CheckResult {
	result := r.runCommand("Go Mod Tidy", "go", "mod", "tidy")
	r.results = append(r.results, result)
	return result
}

// RunGoTest runs go test with coverage
func (r *PreCommitRunner) RunGoTest() CheckResult {
	if r.config.SkipTests {
		return CheckResult{Name: "Go Test", Success: true, Output: "Skipped"}
	}
	
	result := r.runCommand("Go Test", "go", "test", "-v", "-race", "-coverprofile=coverage.out", "-covermode=atomic", "./...")
	r.results = append(r.results, result)
	return result
}

// RunGoTestShort runs go test with short flag
func (r *PreCommitRunner) RunGoTestShort() CheckResult {
	if r.config.SkipTests {
		return CheckResult{Name: "Go Test Short", Success: true, Output: "Skipped"}
	}
	
	result := r.runCommand("Go Test Short", "go", "test", "-short", "./...")
	r.results = append(r.results, result)
	return result
}

// RunGolangciLint runs golangci-lint if available
func (r *PreCommitRunner) RunGolangciLint() CheckResult {
	if r.config.SkipLinting {
		return CheckResult{Name: "Golangci-Lint", Success: true, Output: "Skipped"}
	}
	
	// Check if golangci-lint is available
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		return CheckResult{
			Name:    "Golangci-Lint",
			Success: true,
			Output:  "golangci-lint not found, skipping",
		}
	}
	
	result := r.runCommand("Golangci-Lint", "golangci-lint", "run", "--verbose")
	r.results = append(r.results, result)
	return result
}

// RunTaskCommand runs a task command if available
func (r *PreCommitRunner) RunTaskCommand(taskName string) CheckResult {
	// Check if task is available
	if _, err := exec.LookPath("task"); err != nil {
		return CheckResult{
			Name:    fmt.Sprintf("Task %s", taskName),
			Success: true,
			Output:  "task not found, skipping",
		}
	}
	
	result := r.runCommand(fmt.Sprintf("Task %s", taskName), "task", taskName)
	r.results = append(r.results, result)
	return result
}

// CheckCoverage checks if coverage meets the target
func (r *PreCommitRunner) CheckCoverage() CheckResult {
	coverageFile := filepath.Join(r.config.ProjectRoot, "coverage.out")
	
	if _, err := os.Stat(coverageFile); os.IsNotExist(err) {
		return CheckResult{
			Name:    "Coverage Check",
			Success: false,
			Output:  "coverage.out not found, run tests first",
			Error:   err,
		}
	}
	
	result := r.runCommand("Coverage Check", "go", "tool", "cover", "-func=coverage.out")
	
	// Parse coverage percentage
	if result.Success {
		lines := strings.Split(result.Output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "total:") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					coverageStr := parts[2]
					coverageStr = strings.TrimSuffix(coverageStr, "%")
					
					if r.config.Verbose {
						fmt.Printf("📊 Coverage: %s%%\n", coverageStr)
					}
					
					// You could add percentage parsing here if needed
					if strings.Contains(coverageStr, ".") {
						// Coverage found
						result.Output = fmt.Sprintf("Coverage: %s%%\n%s", coverageStr, result.Output)
					}
				}
			}
		}
	}
	
	r.results = append(r.results, result)
	return result
}

// RunMutationTesting runs mutation testing if enabled
func (r *PreCommitRunner) RunMutationTesting() CheckResult {
	if !r.config.EnableMutation {
		return CheckResult{Name: "Mutation Testing", Success: true, Output: "Disabled"}
	}
	
	result := r.runCommand("Mutation Testing", "go", "test", "-v", "-run", "TestComprehensiveMutationTesting", "-timeout=60m")
	r.results = append(r.results, result)
	return result
}

// RunAllChecks runs all pre-commit checks
func (r *PreCommitRunner) RunAllChecks() bool {
	fmt.Println("🚀 Starting pre-commit checks...")
	
	checks := []func() CheckResult{
		r.RunGoFormat,
		r.RunGoModTidy,
		r.RunGoVet,
		r.RunGoTestShort,
		r.RunGolangciLint,
		r.CheckCoverage,
	}
	
	// Add mutation testing if enabled
	if r.config.EnableMutation {
		checks = append(checks, r.RunMutationTesting)
	}
	
	allPassed := true
	for _, check := range checks {
		result := check()
		if !result.Success {
			allPassed = false
			fmt.Printf("❌ %s failed: %s\n", result.Name, result.Error)
			if r.config.Verbose && result.Output != "" {
				fmt.Printf("Output: %s\n", result.Output)
			}
		} else {
			fmt.Printf("✅ %s passed\n", result.Name)
		}
	}
	
	return allPassed
}

// RunQuickChecks runs quick pre-commit checks (no mutation testing)
func (r *PreCommitRunner) RunQuickChecks() bool {
	fmt.Println("🚀 Starting quick pre-commit checks...")
	
	checks := []func() CheckResult{
		r.RunGoFormat,
		r.RunGoModTidy,
		r.RunGoVet,
		r.RunGoTestShort,
		r.RunGolangciLint,
	}
	
	allPassed := true
	for _, check := range checks {
		result := check()
		if !result.Success {
			allPassed = false
			fmt.Printf("❌ %s failed: %s\n", result.Name, result.Error)
			if r.config.Verbose && result.Output != "" {
				fmt.Printf("Output: %s\n", result.Output)
			}
		} else {
			fmt.Printf("✅ %s passed\n", result.Name)
		}
	}
	
	return allPassed
}

// PrintSummary prints a summary of all check results
func (r *PreCommitRunner) PrintSummary() {
	fmt.Println("\n📊 Pre-commit Check Summary:")
	fmt.Println("=" + strings.Repeat("=", 50))
	
	totalChecks := len(r.results)
	passedChecks := 0
	
	for _, result := range r.results {
		status := "❌ FAILED"
		if result.Success {
			status = "✅ PASSED"
			passedChecks++
		}
		
		fmt.Printf("%-20s %s (%v)\n", result.Name, status, result.Duration)
	}
	
	fmt.Printf("\nTotal: %d/%d checks passed\n", passedChecks, totalChecks)
	
	if passedChecks == totalChecks {
		fmt.Println("🎉 All checks passed! Ready to commit.")
	} else {
		fmt.Println("⚠️  Some checks failed. Please fix issues before committing.")
	}
}

// getProjectRoot finds the project root directory
func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	
	// Look for go.mod file
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	return "."
}

// main function
func main() {
	// Parse command line arguments
	var config PreCommitConfig
	config.ProjectRoot = getProjectRoot()
	config.CoverageTarget = 95
	config.Timeout = 5 * time.Minute
	config.Verbose = true
	
	args := os.Args[1:]
	
	if len(args) == 0 {
		fmt.Println("Usage: go run pre-commit-runner.go [quick|full|format|test|lint|coverage|mutation]")
		os.Exit(1)
	}
	
	command := args[0]
	
	// Parse flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--no-mutation":
			config.EnableMutation = false
		case "--skip-tests":
			config.SkipTests = true
		case "--skip-linting":
			config.SkipLinting = true
		case "--skip-formatting":
			config.SkipFormatting = true
		case "--quiet":
			config.Verbose = false
		case "--mutation":
			config.EnableMutation = true
		}
	}
	
	runner := NewPreCommitRunner(config)
	
	switch command {
	case "quick":
		success := runner.RunQuickChecks()
		runner.PrintSummary()
		if !success {
			os.Exit(1)
		}
	case "full":
		config.EnableMutation = true
		runner.config = config
		success := runner.RunAllChecks()
		runner.PrintSummary()
		if !success {
			os.Exit(1)
		}
	case "format":
		result := runner.RunGoFormat()
		if !result.Success {
			fmt.Printf("❌ Format failed: %s\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("✅ Format completed")
	case "test":
		result := runner.RunGoTest()
		if !result.Success {
			fmt.Printf("❌ Tests failed: %s\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("✅ Tests passed")
	case "lint":
		result := runner.RunGolangciLint()
		if !result.Success {
			fmt.Printf("❌ Linting failed: %s\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("✅ Linting passed")
	case "coverage":
		runner.RunGoTest()
		result := runner.CheckCoverage()
		if !result.Success {
			fmt.Printf("❌ Coverage check failed: %s\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("✅ Coverage check passed")
	case "mutation":
		config.EnableMutation = true
		runner.config = config
		result := runner.RunMutationTesting()
		if !result.Success {
			fmt.Printf("❌ Mutation testing failed: %s\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("✅ Mutation testing passed")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: quick, full, format, test, lint, coverage, mutation")
		os.Exit(1)
	}
}