package main

import (
	"fmt"
	"time"
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
	UseTopics    bool
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

	errors = append(errors, validateGitHubStringFields(config)...)
	errors = append(errors, validateGitHubNumericFields(config)...)

	return errors
}

// validateGitHubStringFields validates string fields in GitHub configuration (Pure Core)
func validateGitHubStringFields(config GitHubConfig) []ValidationError {
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

	return errors
}

// validateGitHubNumericFields validates numeric fields in GitHub configuration (Pure Core)
func validateGitHubNumericFields(config GitHubConfig) []ValidationError {
	var errors []ValidationError

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

	errors = append(errors, validateNeo4jStringFields(config)...)
	errors = append(errors, validateNeo4jTimeoutField(config)...)

	return errors
}

// validateNeo4jStringFields validates string fields in Neo4j configuration (Pure Core)
func validateNeo4jStringFields(config Neo4jConfig) []ValidationError {
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

	return errors
}

// validateNeo4jTimeoutField validates timeout field in Neo4j configuration (Pure Core)
func validateNeo4jTimeoutField(config Neo4jConfig) []ValidationError {
	var errors []ValidationError

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

