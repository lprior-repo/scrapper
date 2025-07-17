package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
)

// GitHubScanConfig represents configuration for GitHub scanning
type GitHubScanConfig struct {
	Token        string `json:"token"`
	Organization string `json:"organization"`
	MaxRepos     int    `json:"max_repos"`
	MaxTeams     int    `json:"max_teams"`
	OutputFile   string `json:"output_file"`
}

// GitHubScanResult represents the result of a GitHub scan
type GitHubScanResult struct {
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
	Data    GitHubOrgData `json:"data,omitempty"`
	Summary ScanSummary   `json:"summary"`
}

// ScanSummary provides a summary of the scan results
type ScanSummary struct {
	TotalRepos          int      `json:"total_repos"`
	ReposWithCodeowners int      `json:"repos_with_codeowners"`
	TotalTeams          int      `json:"total_teams"`
	UniqueOwners        []string `json:"unique_owners"`
	APICallsUsed        int      `json:"api_calls_used"`
}

// scanGitHubOrganization scans a GitHub organization for codeowners and teams
func scanGitHubOrganization(ctx context.Context, config GitHubScanConfig) GitHubScanResult {
	// Validate configuration
	if err := validateScanConfig(config); err != nil {
		return GitHubScanResult{
			Success: false,
			Error:   fmt.Sprintf("invalid configuration: %v", err),
		}
	}

	// Create GitHub client
	client, err := createGitHubClient(ctx, config.Token, config.Organization)
	if err != nil {
		return GitHubScanResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create GitHub client: %v", err),
		}
	}

	// Prepare batch request
	request := BatchRequest{
		Organization: config.Organization,
		MaxRepos:     config.MaxRepos,
		MaxTeams:     config.MaxTeams,
	}

	// Fetch all data with minimal API calls
	data, err := fetchAllOrgData(ctx, client, request)
	if err != nil {
		return GitHubScanResult{
			Success: false,
			Error:   fmt.Sprintf("failed to fetch organization data: %v", err),
		}
	}

	// Generate summary
	summary := generateScanSummary(data)

	// Save to file if requested
	if config.OutputFile != "" {
		if err := saveOrgDataToFile(config.OutputFile, data); err != nil {
			return GitHubScanResult{
				Success: false,
				Error:   fmt.Sprintf("failed to save data to file: %v", err),
			}
		}
	}

	return GitHubScanResult{
		Success: true,
		Data:    data,
		Summary: summary,
	}
}

// validateScanConfig validates the scan configuration
func validateScanConfig(config GitHubScanConfig) error {
	if config.Token == "" {
		return fmt.Errorf("GitHub token is required")
	}
	if config.Organization == "" {
		return fmt.Errorf("organization name is required")
	}
	if config.MaxRepos < 0 {
		return fmt.Errorf("max repos cannot be negative")
	}
	if config.MaxTeams < 0 {
		return fmt.Errorf("max teams cannot be negative")
	}
	return nil
}

// generateScanSummary generates a summary of the scan results
func generateScanSummary(data GitHubOrgData) ScanSummary {
	reposWithCodeowners := lo.Filter(data.Repos, func(repo GitHubRepo, _ int) bool {
		return repo.HasCodeownersFile
	})

	// Extract all unique owners from all repos
	allOwners := []string{}
	for _, repo := range data.Repos {
		if repo.CodeownersContent != "" {
			entries := parseCodeownersContent(repo.CodeownersContent)
			owners := extractUniqueOwners(entries)
			allOwners = append(allOwners, owners...)
		}
	}

	uniqueOwners := lo.Uniq(allOwners)

	return ScanSummary{
		TotalRepos:          len(data.Repos),
		ReposWithCodeowners: len(reposWithCodeowners),
		TotalTeams:          len(data.Teams),
		UniqueOwners:        uniqueOwners,
		APICallsUsed:        data.APICallCount,
	}
}

// analyzeCodeownersCoverage analyzes CODEOWNERS coverage patterns
func analyzeCodeownersCoverage(repos []GitHubRepo) map[string]int {
	patternCounts := make(map[string]int)

	for _, repo := range repos {
		if repo.CodeownersContent != "" {
			entries := parseCodeownersContent(repo.CodeownersContent)
			for _, entry := range entries {
				patternCounts[entry.Pattern]++
			}
		}
	}

	return patternCounts
}

// saveOrgDataToFile saves organization data to a file
func saveOrgDataToFile(filename string, data GitHubOrgData) error {
	if filename == "" {
		panic("Filename cannot be empty")
	}

	// Marshal data to JSON
	jsonData, err := marshalOrgData(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// marshalOrgData marshals organization data to JSON
func marshalOrgData(data GitHubOrgData) ([]byte, error) {
	// This is a pure function that prepares data for writing
	return json.MarshalIndent(data, "", "  ")
}
