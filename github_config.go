package main

import (
	"fmt"
	"os"
	"strings"
)

// GitHubTokenConfig handles secure token management
type GitHubTokenConfig struct {
	tokenEnvVar string
	orgEnvVar   string
}

// createGitHubTokenConfig creates a new token configuration
func createGitHubTokenConfig() GitHubTokenConfig {
	return GitHubTokenConfig{
		tokenEnvVar: "GITHUB_TOKEN",
		orgEnvVar:   "GITHUB_ORG",
	}
}

// validateGitHubToken validates a GitHub token format
func validateGitHubToken(token string) error {
	if token == "" {
		return fmt.Errorf("GitHub token cannot be empty")
	}

	// Basic token format validation
	if len(token) < 20 {
		return fmt.Errorf("GitHub token appears to be too short")
	}

	// Check for common token prefixes
	validPrefixes := []string{"ghp_", "gho_", "ghu_", "ghs_", "ghr_"}
	hasValidPrefix := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(token, prefix) {
			hasValidPrefix = true
			break
		}
	}

	if !hasValidPrefix {
		return fmt.Errorf("GitHub token does not have a valid prefix (ghp_, gho_, ghu_, ghs_, ghr_)")
	}

	return nil
}

// getGitHubTokenFromEnv securely retrieves GitHub token from environment
func getGitHubTokenFromEnv(config GitHubTokenConfig) (string, error) {
	if config.tokenEnvVar == "" {
		panic("Token environment variable name cannot be empty")
	}

	token := os.Getenv(config.tokenEnvVar)
	if token == "" {
		return "", fmt.Errorf("GitHub token not found in environment variable %s", config.tokenEnvVar)
	}

	// Validate token format
	if err := validateGitHubToken(token); err != nil {
		return "", fmt.Errorf("invalid GitHub token: %w", err)
	}

	return token, nil
}

// getGitHubOrgFromEnv retrieves GitHub organization from environment
func getGitHubOrgFromEnv(config GitHubTokenConfig) (string, error) {
	if config.orgEnvVar == "" {
		panic("Organization environment variable name cannot be empty")
	}

	org := os.Getenv(config.orgEnvVar)
	if org == "" {
		return "", fmt.Errorf("GitHub organization not found in environment variable %s", config.orgEnvVar)
	}

	// Validate organization name
	if err := validateGitHubOrgName(org); err != nil {
		return "", fmt.Errorf("invalid GitHub organization: %w", err)
	}

	return org, nil
}

// validateGitHubOrgName validates GitHub organization name format
func validateGitHubOrgName(org string) error {
	if err := validateOrgNameBasics(org); err != nil {
		return err
	}

	trimmed := strings.TrimSpace(org)
	if err := validateOrgNameLength(trimmed); err != nil {
		return err
	}

	if err := validateOrgNameCharacters(trimmed); err != nil {
		return err
	}

	return validateOrgNameHyphens(trimmed)
}

// validateOrgNameBasics validates basic org name requirements
func validateOrgNameBasics(org string) error {
	if org == "" {
		return fmt.Errorf("organization name cannot be empty")
	}

	if strings.TrimSpace(org) == "" {
		return fmt.Errorf("organization name cannot be whitespace only")
	}

	return nil
}

// validateOrgNameLength validates organization name length
func validateOrgNameLength(trimmed string) error {
	if len(trimmed) > 39 {
		return fmt.Errorf("organization name too long (max 39 characters)")
	}
	return nil
}

// validateOrgNameCharacters validates organization name characters
func validateOrgNameCharacters(trimmed string) error {
	for _, char := range trimmed {
		if !isValidOrgChar(char) {
			return fmt.Errorf("organization name contains invalid characters")
		}
	}
	return nil
}

// isValidOrgChar checks if a character is valid for GitHub org names
func isValidOrgChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '-'
}

// validateOrgNameHyphens validates hyphen placement in org name
func validateOrgNameHyphens(trimmed string) error {
	if strings.HasPrefix(trimmed, "-") || strings.HasSuffix(trimmed, "-") {
		return fmt.Errorf("organization name cannot start or end with hyphen")
	}
	return nil
}

// maskToken masks a GitHub token for safe logging
func maskToken(token string) string {
	if token == "" {
		return ""
	}

	if len(token) <= 8 {
		return "****"
	}

	// Show first 4 and last 4 characters
	return token[:4] + "****" + token[len(token)-4:]
}

// createSecureGitHubScanConfig creates scan config with secure token handling
func createSecureGitHubScanConfig() (GitHubScanConfig, error) {
	tokenConfig := createGitHubTokenConfig()

	// Get token securely from environment
	token, err := getGitHubTokenFromEnv(tokenConfig)
	if err != nil {
		return GitHubScanConfig{}, fmt.Errorf("failed to get GitHub token: %w", err)
	}

	// Get organization from environment
	org, err := getGitHubOrgFromEnv(tokenConfig)
	if err != nil {
		return GitHubScanConfig{}, fmt.Errorf("failed to get GitHub organization: %w", err)
	}

	// Parse other configuration
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
