package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// loadConfigFromEnv loads configuration from environment variables
func loadConfigFromEnv() AppConfig {
	return AppConfig{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Port:        getIntEnvOrDefault("HTTP_PORT", 8081),
		GitHub:      loadGitHubConfig(),
		Neo4j:       loadNeo4jConfig(),
		Server:      loadServerConfig(),
	}
}

// loadGitHubConfig loads GitHub configuration from environment
func loadGitHubConfig() GitHubConfig {
	return GitHubConfig{
		Token:        os.Getenv("GITHUB_TOKEN"),
		BaseURL:      getEnvOrDefault("GITHUB_BASE_URL", "https://api.github.com"),
		UserAgent:    getEnvOrDefault("GITHUB_USER_AGENT", "overseer-codeowners-scanner/1.0"),
		Timeout:      getDurationEnvOrDefault("GITHUB_TIMEOUT", 30*time.Second),
		MaxRetries:   getIntEnvOrDefault("GITHUB_MAX_RETRIES", 3),
		RateLimitMin: getIntEnvOrDefault("GITHUB_RATE_LIMIT_MIN", 100),
	}
}

// loadNeo4jConfig loads Neo4j configuration from environment
func loadNeo4jConfig() Neo4jConfig {
	return Neo4jConfig{
		URI:      getEnvOrDefault("NEO4J_URI", "bolt://localhost:7687"),
		Username: getEnvOrDefault("NEO4J_USERNAME", "neo4j"),
		Password: getEnvOrDefault("NEO4J_PASSWORD", "password"),
		Database: getEnvOrDefault("NEO4J_DATABASE", "neo4j"),
		Timeout:  getDurationEnvOrDefault("NEO4J_TIMEOUT", 30*time.Second),
	}
}

// loadServerConfig loads server configuration from environment
func loadServerConfig() ServerConfig {
	return ServerConfig{
		ReadTimeout:    getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:   getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:    getDurationEnvOrDefault("SERVER_IDLE_TIMEOUT", 60*time.Second),
		MaxHeaderBytes: getIntEnvOrDefault("SERVER_MAX_HEADER_BYTES", 1<<20),
	}
}

// getEnvOrDefault gets environment variable or returns default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnvOrDefault gets int environment variable or returns default
func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getDurationEnvOrDefault gets duration environment variable or returns default
func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// validateConfiguration validates the loaded configuration
func validateConfiguration(config AppConfig) error {
	validationErrors := validateAppConfig(config)
	if len(validationErrors) > 0 {
		return fmt.Errorf("configuration validation failed: %d errors found", len(validationErrors))
	}
	return nil
}