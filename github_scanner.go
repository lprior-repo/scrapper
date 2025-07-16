package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

// runGitHubScanner runs the GitHub organization scanner
func runGitHubScanner(ctx context.Context) error {
	// Parse configuration from environment
	config, err := parseGitHubScanConfig()
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	
	// Validate configuration
	if err := validateScanConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	
	fmt.Printf("🔍 Scanning GitHub organization: %s\n", config.Organization)
	fmt.Printf("📊 Target: %d repos, %d teams (0 = unlimited)\n", config.MaxRepos, config.MaxTeams)
	fmt.Printf("🎯 Goal: ≤50 API calls\n\n")
	
	// Estimate API calls needed
	estimatedCalls := calculateAPICallsNeeded(config.MaxRepos, config.MaxTeams)
	if estimatedCalls > 50 {
		fmt.Printf("⚠️  Warning: Estimated %d API calls (exceeds 50 target)\n", estimatedCalls)
	}
	
	// Run the scan
	startTime := time.Now()
	result := scanGitHubOrganization(ctx, config)
	duration := time.Since(startTime)
	
	// Display results
	if !result.Success {
		return fmt.Errorf("scan failed: %s", result.Error)
	}
	
	displayScanResults(result, duration)
	
	return nil
}

// parseGitHubScanConfig parses configuration from environment variables
func parseGitHubScanConfig() (GitHubScanConfig, error) {
	token := getEnvOrDefault("GITHUB_TOKEN", "")
	if token == "" {
		return GitHubScanConfig{}, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}
	
	org := getEnvOrDefault("GITHUB_ORG", "")
	if org == "" {
		return GitHubScanConfig{}, fmt.Errorf("GITHUB_ORG environment variable is required")
	}
	
	maxRepos := parseIntEnv("GITHUB_MAX_REPOS", 0)
	maxTeams := parseIntEnv("GITHUB_MAX_TEAMS", 0)
	outputFile := getEnvOrDefault("GITHUB_OUTPUT_FILE", "")
	
	return GitHubScanConfig{
		Token:        token,
		Organization: org,
		MaxRepos:     maxRepos,
		MaxTeams:     maxTeams,
		OutputFile:   outputFile,
	}, nil
}

// parseIntEnv parses an integer from environment variable
func parseIntEnv(key string, defaultValue int) int {
	if key == "" {
		panic("Environment variable key cannot be empty")
	}
	
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return parsed
}

// displayScanResults displays the scan results in a formatted way
func displayScanResults(result GitHubScanResult, duration time.Duration) {
	fmt.Printf("✅ Scan completed in %v\n\n", duration)
	
	summary := result.Summary
	
	fmt.Printf("📈 SUMMARY\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🏢 Organization: %s\n", result.Data.Organization)
	fmt.Printf("📚 Total Repositories: %d\n", summary.TotalRepos)
	fmt.Printf("📋 Repos with CODEOWNERS: %d (%.1f%%)\n", 
		summary.ReposWithCodeowners, 
		calculatePercentage(summary.ReposWithCodeowners, summary.TotalRepos))
	fmt.Printf("👥 Total Teams: %d\n", summary.TotalTeams)
	fmt.Printf("🔗 Unique Owners: %d\n", len(summary.UniqueOwners))
	fmt.Printf("🌐 API Calls Used: %d/50\n", summary.APICallsUsed)
	
	if summary.APICallsUsed <= 50 {
		fmt.Printf("🎯 Target achieved! Used %d/50 API calls\n", summary.APICallsUsed)
	} else {
		fmt.Printf("⚠️  Exceeded target: %d/50 API calls\n", summary.APICallsUsed)
	}
	
	fmt.Printf("\n📊 DETAILED ANALYSIS\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Show repos without CODEOWNERS
	reposWithout := findReposWithoutCodeowners(result.Data)
	if len(reposWithout) > 0 {
		fmt.Printf("📭 Repositories without CODEOWNERS (%d):\n", len(reposWithout))
		for i, repo := range reposWithout {
			if i < 10 { // Show first 10
				fmt.Printf("   • %s\n", repo)
			}
		}
		if len(reposWithout) > 10 {
			fmt.Printf("   ... and %d more\n", len(reposWithout)-10)
		}
	}
	
	// Show teams by privacy
	teamsByPrivacy := findTeamsByPrivacy(result.Data)
	fmt.Printf("\n👥 Teams by Privacy:\n")
	for privacy, teams := range teamsByPrivacy {
		fmt.Printf("   • %s: %d teams\n", privacy, len(teams))
	}
	
	// Show most common owners
	if len(summary.UniqueOwners) > 0 {
		fmt.Printf("\n🏆 Sample Owners:\n")
		for i, owner := range summary.UniqueOwners {
			if i < 5 { // Show first 5
				fmt.Printf("   • %s\n", owner)
			}
		}
		if len(summary.UniqueOwners) > 5 {
			fmt.Printf("   ... and %d more\n", len(summary.UniqueOwners)-5)
		}
	}
	
	// Show optimization recommendations
	fmt.Printf("\n💡 OPTIMIZATION\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	recommendations := optimizeAPIUsage(summary.TotalRepos, summary.TotalTeams)
	fmt.Printf("%s\n", recommendations)
}

// calculatePercentage calculates percentage safely
func calculatePercentage(part, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return (float64(part) / float64(total)) * 100.0
}

// demonstrateGitHubScanner demonstrates the GitHub scanner functionality
func demonstrateGitHubScanner(ctx context.Context) error {
	fmt.Printf("🚀 GitHub Organization Scanner Demo\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Check if we have required environment variables
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Printf("❌ Demo requires GITHUB_TOKEN environment variable\n")
		fmt.Printf("💡 Set GITHUB_TOKEN=your_token to run the demo\n")
		return nil
	}
	
	if os.Getenv("GITHUB_ORG") == "" {
		fmt.Printf("❌ Demo requires GITHUB_ORG environment variable\n")
		fmt.Printf("💡 Set GITHUB_ORG=organization_name to run the demo\n")
		return nil
	}
	
	// Run the actual scanner
	return runGitHubScanner(ctx)
}