package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/samber/lo"
)

// AppConfig represents the complete application configuration
type AppConfig struct {
	Environment string
	Port        int
	GitHub      GitHubConfig
	Neo4j       Neo4jConfig
	Server      ServerConfig
}

// GitHubConfig represents GitHub API configuration
type GitHubConfig struct {
	Token        string
	BaseURL      string
	UserAgent    string
	Timeout      time.Duration
	MaxRetries   int
	RateLimitMin int
}

// Neo4jConfig represents Neo4j database configuration
type Neo4jConfig struct {
	URI      string
	Username string
	Password string
	Database string
	Timeout  time.Duration
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

// ValidationError represents configuration validation errors
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

// Error implements the error interface for ValidationError
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// buildDefaultAppConfig creates default application configuration (Pure Core)
func buildDefaultAppConfig() AppConfig {
	return AppConfig{
		Environment: "development",
		Port:        8081,
		GitHub:      buildDefaultGitHubConfig(),
		Neo4j:       buildDefaultNeo4jConfig(),
		Server:      buildDefaultServerConfig(),
	}
}

// buildDefaultGitHubConfig creates default GitHub configuration (Pure Core)
func buildDefaultGitHubConfig() GitHubConfig {
	return GitHubConfig{
		Token:        "",
		BaseURL:      "https://api.github.com",
		UserAgent:    "overseer-codeowners-scanner/1.0",
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RateLimitMin: 100,
	}
}

// buildDefaultNeo4jConfig creates default Neo4j configuration (Pure Core)
func buildDefaultNeo4jConfig() Neo4jConfig {
	return Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
		Database: "neo4j",
		Timeout:  30 * time.Second,
	}
}

// buildDefaultServerConfig creates default server configuration (Pure Core)
func buildDefaultServerConfig() ServerConfig {
	return ServerConfig{
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
}

// applyEnvironmentOverrides applies environment variable overrides to config (Pure Core)
func applyEnvironmentOverrides(config AppConfig, envVars map[string]string) AppConfig {
	validateConfigNotNil(config)
	validateEnvVarsNotNil(envVars)

	overridden := config

	// Apply environment overrides
	if env := getStringValue(envVars, "ENVIRONMENT"); env != "" {
		overridden.Environment = env
	}

	if port := getIntValue(envVars, "PORT"); port > 0 {
		overridden.Port = port
	}

	// Apply GitHub overrides
	overridden.GitHub = applyGitHubOverrides(overridden.GitHub, envVars)

	// Apply Neo4j overrides
	overridden.Neo4j = applyNeo4jOverrides(overridden.Neo4j, envVars)

	return overridden
}

// applyGitHubOverrides applies GitHub environment overrides (Pure Core)
func applyGitHubOverrides(config GitHubConfig, envVars map[string]string) GitHubConfig {
	validateEnvVarsNotNil(envVars)

	overridden := config

	if token := getStringValue(envVars, "GITHUB_TOKEN"); token != "" {
		overridden.Token = token
	}

	if baseURL := getStringValue(envVars, "GITHUB_BASE_URL"); baseURL != "" {
		overridden.BaseURL = baseURL
	}

	if userAgent := getStringValue(envVars, "GITHUB_USER_AGENT"); userAgent != "" {
		overridden.UserAgent = userAgent
	}

	if timeout := getDurationValue(envVars, "GITHUB_TIMEOUT"); timeout > 0 {
		overridden.Timeout = timeout
	}

	if maxRetries := getIntValue(envVars, "GITHUB_MAX_RETRIES"); maxRetries > 0 {
		overridden.MaxRetries = maxRetries
	}

	if rateLimitMin := getIntValue(envVars, "GITHUB_RATE_LIMIT_MIN"); rateLimitMin > 0 {
		overridden.RateLimitMin = rateLimitMin
	}

	return overridden
}

// applyNeo4jOverrides applies Neo4j environment overrides (Pure Core)
func applyNeo4jOverrides(config Neo4jConfig, envVars map[string]string) Neo4jConfig {
	validateEnvVarsNotNil(envVars)

	overridden := config

	if uri := getStringValue(envVars, "NEO4J_URI"); uri != "" {
		overridden.URI = uri
	}

	if username := getStringValue(envVars, "NEO4J_USERNAME"); username != "" {
		overridden.Username = username
	}

	if password := getStringValue(envVars, "NEO4J_PASSWORD"); password != "" {
		overridden.Password = password
	}

	if database := getStringValue(envVars, "NEO4J_DATABASE"); database != "" {
		overridden.Database = database
	}

	if timeout := getDurationValue(envVars, "NEO4J_TIMEOUT"); timeout > 0 {
		overridden.Timeout = timeout
	}

	return overridden
}

// validateAppConfig validates the complete application configuration (Pure Core)
func validateAppConfig(config AppConfig) []ValidationError {
	var errors []ValidationError

	// Validate environment
	if config.Environment == "" {
		errors = append(errors, ValidationError{
			Field:   "Environment",
			Message: "cannot be empty",
			Value:   config.Environment,
		})
	}

	// Validate port
	if config.Port <= 0 || config.Port > 65535 {
		errors = append(errors, ValidationError{
			Field:   "Port",
			Message: "must be between 1 and 65535",
			Value:   config.Port,
		})
	}

	// Validate GitHub config
	githubErrors := validateGitHubConfig(config.GitHub)
	errors = append(errors, githubErrors...)

	// Validate Neo4j config
	neo4jErrors := validateNeo4jConfig(config.Neo4j)
	errors = append(errors, neo4jErrors...)

	// Validate server config
	serverErrors := validateServerConfig(config.Server)
	errors = append(errors, serverErrors...)

	return errors
}

// validateGitHubConfig validates GitHub configuration (Pure Core)
func validateGitHubConfig(config GitHubConfig) []ValidationError {
	var errors []ValidationError

	if config.Token == "" {
		errors = append(errors, ValidationError{
			Field:   "GitHub.Token",
			Message: "cannot be empty",
			Value:   config.Token,
		})
	}

	if config.BaseURL == "" {
		errors = append(errors, ValidationError{
			Field:   "GitHub.BaseURL",
			Message: "cannot be empty",
			Value:   config.BaseURL,
		})
	}

	if config.UserAgent == "" {
		errors = append(errors, ValidationError{
			Field:   "GitHub.UserAgent",
			Message: "cannot be empty",
			Value:   config.UserAgent,
		})
	}

	if config.Timeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "GitHub.Timeout",
			Message: "must be positive",
			Value:   config.Timeout,
		})
	}

	if config.MaxRetries < 0 {
		errors = append(errors, ValidationError{
			Field:   "GitHub.MaxRetries",
			Message: "cannot be negative",
			Value:   config.MaxRetries,
		})
	}

	if config.RateLimitMin < 0 {
		errors = append(errors, ValidationError{
			Field:   "GitHub.RateLimitMin",
			Message: "cannot be negative",
			Value:   config.RateLimitMin,
		})
	}

	return errors
}

// validateNeo4jConfig validates Neo4j configuration (Pure Core)
func validateNeo4jConfig(config Neo4jConfig) []ValidationError {
	var errors []ValidationError

	if config.URI == "" {
		errors = append(errors, ValidationError{
			Field:   "Neo4j.URI",
			Message: "cannot be empty",
			Value:   config.URI,
		})
	}

	if config.Username == "" {
		errors = append(errors, ValidationError{
			Field:   "Neo4j.Username",
			Message: "cannot be empty",
			Value:   config.Username,
		})
	}

	if config.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "Neo4j.Password",
			Message: "cannot be empty",
			Value:   config.Password,
		})
	}

	if config.Database == "" {
		errors = append(errors, ValidationError{
			Field:   "Neo4j.Database",
			Message: "cannot be empty",
			Value:   config.Database,
		})
	}

	if config.Timeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "Neo4j.Timeout",
			Message: "must be positive",
			Value:   config.Timeout,
		})
	}

	return errors
}

// validateServerConfig validates server configuration (Pure Core)
func validateServerConfig(config ServerConfig) []ValidationError {
	var errors []ValidationError

	if config.ReadTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "Server.ReadTimeout",
			Message: "must be positive",
			Value:   config.ReadTimeout,
		})
	}

	if config.WriteTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "Server.WriteTimeout",
			Message: "must be positive",
			Value:   config.WriteTimeout,
		})
	}

	if config.IdleTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "Server.IdleTimeout",
			Message: "must be positive",
			Value:   config.IdleTimeout,
		})
	}

	if config.MaxHeaderBytes <= 0 {
		errors = append(errors, ValidationError{
			Field:   "Server.MaxHeaderBytes",
			Message: "must be positive",
			Value:   config.MaxHeaderBytes,
		})
	}

	return errors
}

// getStringValue extracts string value from environment variables map (Pure Core)
func getStringValue(envVars map[string]string, key string) string {
	validateEnvVarsNotNil(envVars)
	validateKeyNotEmpty(key)

	value, exists := envVars[key]
	if !exists {
		return ""
	}

	return value
}

// getIntValue extracts integer value from environment variables map (Pure Core)
func getIntValue(envVars map[string]string, key string) int {
	validateEnvVarsNotNil(envVars)
	validateKeyNotEmpty(key)

	value := getStringValue(envVars, key)
	if value == "" {
		return 0
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return parsed
}

// getDurationValue extracts duration value from environment variables map (Pure Core)
func getDurationValue(envVars map[string]string, key string) time.Duration {
	validateEnvVarsNotNil(envVars)
	validateKeyNotEmpty(key)

	value := getStringValue(envVars, key)
	if value == "" {
		return 0
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}

	return parsed
}

// readEnvironmentVariables reads all environment variables (Impure Shell)
func readEnvironmentVariables() map[string]string {
	environ := os.Environ()

	envVars := lo.Reduce(environ, func(acc map[string]string, item string, _ int) map[string]string {
		parts := lo.Filter([]string{item}, func(s string, _ int) bool {
			return len(s) > 0
		})

		if len(parts) == 0 {
			return acc
		}

		keyValue := parts[0]
		splitIdx := findFirstEquals(keyValue)

		if splitIdx == -1 {
			return acc
		}

		key := keyValue[:splitIdx]
		value := keyValue[splitIdx+1:]

		acc[key] = value
		return acc
	}, make(map[string]string))

	return envVars
}

// findFirstEquals finds the first equals sign in a string (Pure Core)
func findFirstEquals(s string) int {
	for i, char := range s {
		if char == '=' {
			return i
		}
	}
	return -1
}

// Validation helper functions (Pure Core)
func validateConfigNotNil(config AppConfig) {
	// Basic validation - config is a value type so can't be nil
	if config.Environment == "" && config.Port == 0 {
		panic("AppConfig appears to be zero value")
	}
}

func validateEnvVarsNotNil(envVars map[string]string) {
	if envVars == nil {
		panic("Environment variables map cannot be nil")
	}
}

func validateKeyNotEmpty(key string) {
	if key == "" {
		panic("Key cannot be empty")
	}
}

// loadAppConfig loads the complete application configuration (Orchestrator)
func loadAppConfig() (AppConfig, error) {
	defaultConfig := buildDefaultAppConfig()
	envVars := readEnvironmentVariables()

	config := applyEnvironmentOverrides(defaultConfig, envVars)

	validationErrors := validateAppConfig(config)
	if len(validationErrors) > 0 {
		return config, fmt.Errorf("configuration validation failed: %d errors found", len(validationErrors))
	}

	return config, nil
}

// isProductionEnvironment checks if running in production (Pure Core)
func isProductionEnvironment(config AppConfig) bool {
	return config.Environment == "production"
}

// isDevelopmentEnvironment checks if running in development (Pure Core)
func isDevelopmentEnvironment(config AppConfig) bool {
	return config.Environment == "development"
}
