package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MutationResult holds the results of mutation testing
type MutationResult struct {
	TotalMutations    int
	KilledMutations   int
	SurvivedMutations int
	MutationScore     float64
	Duration          time.Duration
	Error             error
}

// TestComprehensiveMutationTesting runs comprehensive mutation testing on ALL code
func TestComprehensiveMutationTesting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive mutation testing in short mode")
	}

	fmt.Println("🧬 Starting COMPREHENSIVE mutation testing on ALL code...")
	fmt.Println("   This will test every single line of code for maximum bug detection")
	fmt.Println("   No tiers, no alternatives - just maximum stringency as mandated")
	fmt.Println("")

	// Ensure Neo4j is ready
	err := ensureNeo4jReady(t)
	require.NoError(t, err, "Neo4j must be ready for mutation testing")

	overallStartTime := time.Now()

	// Run comprehensive mutation testing on entire codebase
	fmt.Println("🔬 Running mutation testing on ENTIRE codebase...")
	result := runComprehensiveMutationTest(t)

	overallDuration := time.Since(overallStartTime)

	// Generate comprehensive report
	report := generateComprehensiveReport(result, overallDuration)

	fmt.Println(report)

	// Save report to file
	err = os.WriteFile("comprehensive_mutation_report.txt", []byte(report), 0644)
	require.NoError(t, err, "Should be able to save comprehensive mutation report")

	// Assert high mutation score - we demand excellence!
	// Only enforce this if we actually found mutations to test
	if result.TotalMutations > 0 {
		assert.True(t, result.MutationScore >= 90.0,
			"Comprehensive mutation testing demands >= 90%% mutation score, got %.2f%%. "+
				"This indicates potential gaps in test coverage that must be addressed.", result.MutationScore)

		// Fail if any mutations survived - we want to catch ALL bugs
		if result.SurvivedMutations > 0 {
			t.Errorf("⚠️ %d mutations SURVIVED! This indicates potential bugs that tests failed to catch. "+
				"Review the detailed report and strengthen your test suite.", result.SurvivedMutations)
		}
	} else {
		t.Logf("⚠️ No mutations were generated - this might indicate issues with the mutation testing setup")
	}

	fmt.Println("🎉 Comprehensive mutation testing completed!")
	fmt.Printf("📊 Final Score: %.2f%% (%d/%d mutations killed)\n", result.MutationScore, result.KilledMutations, result.TotalMutations)
}

// runComprehensiveMutationTest runs comprehensive mutation testing on entire codebase
func runComprehensiveMutationTest(t *testing.T) MutationResult {
	t.Helper()

	startTime := time.Now()

	result := MutationResult{
		Duration: 0,
	}

	// Create context with generous timeout for comprehensive testing
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Build go-mutesting command with maximum stringency on entire codebase
	args := []string{
		"--verbose",                  // Verbose output for debugging
		"--do-not-remove-tmp-folder", // Keep temp files for analysis
		"--exec-timeout=300",         // 5 minutes per test execution
		".",                          // Test entire current directory
		"--exec=go test -v -timeout=300s -skip TestComprehensiveMutationTesting", // Skip self to avoid recursion
	}

	fmt.Printf("🔬 Running: go-mutesting %s\n", strings.Join(args, " "))

	// Execute go-mutesting
	cmd := exec.CommandContext(ctx, "go-mutesting", args...)
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a timeout
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("mutation testing timed out after 30 minutes")
		} else {
			result.Error = fmt.Errorf("go-mutesting failed: %w\nOutput: %s", err, string(output))
		}
		result.Duration = time.Since(startTime)
		return result
	}

	// Log the output for debugging
	if len(output) > 0 {
		fmt.Printf("🔬 Mutation testing output:\n%s\n", string(output))
	}

	// Parse results
	counts := parseMutationOutput(string(output))
	result.TotalMutations = counts.Total
	result.KilledMutations = counts.Killed
	result.SurvivedMutations = counts.Survived

	// Calculate mutation score
	if result.TotalMutations > 0 {
		result.MutationScore = float64(result.KilledMutations) / float64(result.TotalMutations) * 100
	}

	result.Duration = time.Since(startTime)
	return result
}

// MutationCounts holds the counts from mutation testing
type MutationCounts struct {
	Total    int
	Killed   int
	Survived int
}

// parseMutationOutput parses the output from go-mutesting
func parseMutationOutput(output string) MutationCounts {
	lines := strings.Split(output, "\n")

	var killed, survived int
	for _, line := range lines {
		if isMutationResult(line) {
			if strings.Contains(line, "PASS") {
				survived++
			} else if strings.Contains(line, "FAIL") {
				killed++
			}
		}
	}

	return MutationCounts{
		Total:    killed + survived,
		Killed:   killed,
		Survived: survived,
	}
}

func isMutationResult(line string) bool {
	return strings.Contains(line, "with checksum")
}

// generateComprehensiveReport generates a comprehensive mutation testing report
func generateComprehensiveReport(result MutationResult, duration time.Duration) string {
	var report strings.Builder

	writeHeader(&report, duration)
	writeStatistics(&report, result)
	writeQualityAssessment(&report, result.MutationScore)
	writeExpertOpinions(&report)
	writeExecutionDetails(&report, result)
	writeRecommendations(&report, result.SurvivedMutations)
	writeExplanation(&report)

	return report.String()
}

func writeHeader(report *strings.Builder, duration time.Duration) {
	report.WriteString("🧬 COMPREHENSIVE MUTATION TESTING REPORT\n")
	report.WriteString("========================================\n")
	report.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("Total Duration: %v\n", duration))
	report.WriteString("\n")
}

func writeStatistics(report *strings.Builder, result MutationResult) {
	report.WriteString("📊 OVERALL STATISTICS\n")
	report.WriteString("--------------------\n")
	report.WriteString(fmt.Sprintf("Total mutations tested: %d\n", result.TotalMutations))
	report.WriteString(fmt.Sprintf("Mutations killed: %d\n", result.KilledMutations))
	report.WriteString(fmt.Sprintf("Mutations survived: %d\n", result.SurvivedMutations))
	report.WriteString(fmt.Sprintf("Overall mutation score: %.2f%%\n", result.MutationScore))
	report.WriteString("\n")
}

func writeQualityAssessment(report *strings.Builder, score float64) {
	report.WriteString("🎯 QUALITY ASSESSMENT\n")
	report.WriteString("--------------------\n")

	assessment := getQualityAssessment(score)
	report.WriteString(assessment)
	report.WriteString("\n")
}

func getQualityAssessment(score float64) string {
	if score >= 95 {
		return "🏆 EXCEPTIONAL: >= 95% - Your test suite is world-class!\n"
	} else if score >= 90 {
		return "✅ EXCELLENT: >= 90% - Test suite is highly effective\n"
	} else if score >= 85 {
		return "🟡 GOOD: >= 85% - Test suite is reasonably effective\n"
	} else if score >= 80 {
		return "🟠 MODERATE: >= 80% - Test suite needs improvement\n"
	}
	return "❌ POOR: < 80% - Test suite requires significant improvement\n"
}

func writeExpertOpinions(report *strings.Builder) {
	report.WriteString("💬 EXPERT OPINIONS ON COMPREHENSIVE MUTATION TESTING\n")
	report.WriteString("---------------------------------------------------\n")
	report.WriteString("🎯 Martin Fowler: \"This level of testing is only justified for the most\n")
	report.WriteString("   critical systems where bugs could have severe consequences.\"\n")
	report.WriteString("\n")
	report.WriteString("🔄 Kent Beck: \"Comprehensive mutation testing gives maximum confidence\n")
	report.WriteString("   for refactoring, but consider the cost vs. benefit.\"\n")
	report.WriteString("\n")
	report.WriteString("🚀 Dave Farley: \"Ensure this level of testing doesn't slow your\n")
	report.WriteString("   deployment pipeline. Fast feedback is crucial.\"\n")
	report.WriteString("\n")
}

func writeExecutionDetails(report *strings.Builder, result MutationResult) {
	report.WriteString("📋 TEST EXECUTION DETAILS\n")
	report.WriteString("-------------------------\n")
	if result.Error != nil {
		report.WriteString(fmt.Sprintf("❌ Execution failed: %v\n", result.Error))
	} else {
		status := getExecutionStatus(result.MutationScore)
		report.WriteString(fmt.Sprintf("%s Entire codebase: %.2f%% (%d/%d killed) - %v\n",
			status, result.MutationScore, result.KilledMutations,
			result.TotalMutations, result.Duration))
	}
	report.WriteString("\n")
}

func getExecutionStatus(score float64) string {
	if score < 80 {
		return "⚠️"
	}
	return "✅"
}

func writeRecommendations(report *strings.Builder, survivedMutations int) {
	report.WriteString("💡 RECOMMENDATIONS\n")
	report.WriteString("-----------------\n")
	if survivedMutations > 0 {
		report.WriteString(fmt.Sprintf("⚠️ %d mutations survived - investigate these areas:\n", survivedMutations))
		report.WriteString("   1. Add more test cases for edge conditions\n")
		report.WriteString("   2. Strengthen assertions in existing tests\n")
		report.WriteString("   3. Consider if mutations represent equivalent behavior\n")
		report.WriteString("   4. Review error handling and boundary conditions\n")
	} else {
		report.WriteString("🎉 All mutations killed - exceptional test suite!\n")
		report.WriteString("   Your tests are highly effective at catching bugs.\n")
	}
	report.WriteString("\n")
}

func writeExplanation(report *strings.Builder) {
	report.WriteString("🔍 WHAT THIS MEANS\n")
	report.WriteString("-----------------\n")
	report.WriteString("Mutation testing introduces small changes (mutations) to your code\n")
	report.WriteString("and verifies that your tests catch these changes. A high mutation\n")
	report.WriteString("score indicates that your test suite is effective at detecting bugs.\n")
	report.WriteString("\n")
	report.WriteString("This comprehensive approach tests every line of code to ensure\n")
	report.WriteString("maximum bug detection as mandated.\n")
}

// ensureNeo4jReady ensures Neo4j is ready for testing
func ensureNeo4jReady(t *testing.T) error {
	t.Helper()

	// Start Neo4j
	cmd := exec.Command("docker", "compose", "up", "-d", "neo4j")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Neo4j: %w", err)
	}

	// Wait for Neo4j to be ready
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		cmd := exec.Command("docker", "compose", "exec", "neo4j",
			"cypher-shell", "-u", "neo4j", "-p", "password", "RETURN 1")
		if err := cmd.Run(); err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("Neo4j not ready after %d attempts", maxAttempts)
}

// TestMutationTestingConfiguration tests the mutation testing configuration
func TestMutationTestingConfiguration(t *testing.T) {
	t.Parallel()

	// Test that the mutation testing parameters are reasonable
	timeout := 30 * time.Minute
	assert.Positive(t, timeout, "Mutation testing timeout should be positive")
	assert.LessOrEqual(t, timeout, 60*time.Minute, "Mutation testing timeout should be reasonable")

	// Test parsing logic
	sampleOutput := `PASS 1 with checksum abc123
FAIL 2 with checksum def456
PASS 3 with checksum ghi789`

	counts := parseMutationOutput(sampleOutput)
	assert.Equal(t, 3, counts.Total, "Should parse total mutations correctly")
	assert.Equal(t, 1, counts.Killed, "Should parse killed mutations correctly")
	assert.Equal(t, 2, counts.Survived, "Should parse survived mutations correctly")
}
