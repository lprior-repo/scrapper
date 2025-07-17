package main

import (
	"os"
	"strconv"
)

// parseGitHubScanConfig parses configuration from environment variables with secure token handling
func parseGitHubScanConfig() (GitHubScanConfig, error) {
	return createSecureGitHubScanConfig()
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

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if key == "" {
		panic("Environment variable key cannot be empty")
	}
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
