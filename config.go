package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	providerNeo4j        = "neo4j"
	providerNeptune      = "neptune"
	envProduction        = "production"
	envDevelopment       = "development"
	defaultNeo4jUsername = "neo4j"
)

// Config holds the application configuration
type Config struct {
	Environment string             `json:"environment"`
	GraphDB     GraphServiceConfig `json:"graph_db"`
}

// LoadConfig loads configuration from environment and defaults (Orchestrator)
func LoadConfig() (*Config, error) {
	envVars := readEnvironmentVariables()
	return buildConfigFromEnvironment(envVars), nil
}

// buildConfigFromEnvironment builds configuration from environment variables (Pure Core)
func buildConfigFromEnvironment(envVars map[string]string) *Config {
	config := &Config{
		Environment: getStringOrDefault(envVars, "ENVIRONMENT", "development"),
		GraphDB:     buildDefaultGraphServiceConfig(envVars),
	}

	return applyEnvironmentOverrides(config, envVars)
}

// applyEnvironmentOverrides applies environment variable overrides to config (Pure Core)
func applyEnvironmentOverrides(config *Config, envVars map[string]string) *Config {
	if config == nil {
		panic("Config cannot be nil")
	}
	if provider := getStringOrDefault(envVars, "GRAPH_DB_PROVIDER", ""); provider != "" {
		config.GraphDB.Provider = provider
	}

	config = applyNeo4jOverrides(config, envVars)
	config = applyNeptuneOverrides(config, envVars)
	return config
}

// applyNeo4jOverrides applies Neo4j environment variable overrides (Pure Core)
func applyNeo4jOverrides(config *Config, envVars map[string]string) *Config {
	if config == nil {
		panic("Config cannot be nil")
	}
	if uri := getStringOrDefault(envVars, "NEO4J_URI", ""); uri != "" {
		config.GraphDB.Neo4j.URI = uri
	}
	if username := getStringOrDefault(envVars, "NEO4J_USERNAME", ""); username != "" {
		config.GraphDB.Neo4j.Username = username
	}
	if password := getStringOrDefault(envVars, "NEO4J_PASSWORD", ""); password != "" {
		config.GraphDB.Neo4j.Password = password
	}
	return config
}

// applyNeptuneOverrides applies Neptune environment variable overrides (Pure Core)
func applyNeptuneOverrides(config *Config, envVars map[string]string) *Config {
	if config == nil {
		panic("Config cannot be nil")
	}
	if endpoint := getStringOrDefault(envVars, "NEPTUNE_ENDPOINT", ""); endpoint != "" {
		config.GraphDB.Neptune.Endpoint = endpoint
	}
	if region := getStringOrDefault(envVars, "NEPTUNE_REGION", ""); region != "" {
		config.GraphDB.Neptune.Region = region
	}
	return config
}

// LoadConfigFromFile loads configuration from a JSON file (Orchestrator)
func LoadConfigFromFile(filename string) (*Config, error) {
	if filename == "" {
		panic("Filename cannot be empty")
	}
	data, err := readConfigFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return buildConfigFromJSON(data)
}

// SaveConfigToFile saves configuration to a JSON file (Orchestrator)
func SaveConfigToFile(config *Config, filename string) error {
	if config == nil {
		panic("Config cannot be nil")
	}
	if filename == "" {
		panic("Filename cannot be empty")
	}
	data, err := serializeConfigToJSON(config)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	return writeConfigFile(filename, data)
}

// getStringOrDefault retrieves string from map or returns default value (Pure Core)
func getStringOrDefault(envVars map[string]string, key, defaultValue string) string {
	if key == "" {
		panic("Environment variable key cannot be empty")
	}
	if value, exists := envVars[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

// buildDefaultGraphServiceConfig returns default graph service configuration (Pure Core)
func buildDefaultGraphServiceConfig(envVars map[string]string) GraphServiceConfig {
	environment := getStringOrDefault(envVars, "ENVIRONMENT", "development")

	config := GraphServiceConfig{
		Neo4j: struct {
			URI      string `json:"uri"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			URI:      getStringOrDefault(envVars, "NEO4J_URI", "bolt://localhost:7687"),
			Username: getStringOrDefault(envVars, "NEO4J_USERNAME", defaultNeo4jUsername),
			Password: getStringOrDefault(envVars, "NEO4J_PASSWORD", "password"),
		},
		Neptune: struct {
			Endpoint string `json:"endpoint"`
			Region   string `json:"region"`
		}{
			Endpoint: getStringOrDefault(envVars, "NEPTUNE_ENDPOINT", ""),
			Region:   getStringOrDefault(envVars, "NEPTUNE_REGION", "us-east-1"),
		},
	}

	if environment == envProduction {
		config.Provider = providerNeptune
	} else {
		config.Provider = providerNeo4j
	}

	return config
}

// validateProductionEnvironment checks if the application is running in production (Pure Core)
func validateProductionEnvironment(config Config) bool {
	if config.Environment == "" {
		panic("Config environment cannot be empty")
	}
	return config.Environment == envProduction
}

// validateDevelopmentEnvironment checks if the application is running in development (Pure Core)
func validateDevelopmentEnvironment(config Config) bool {
	if config.Environment == "" {
		panic("Config environment cannot be empty")
	}
	return config.Environment == "development"
}

// checkIsProduction is an alias for validateProductionEnvironment (for backwards compatibility)
func checkIsProduction(config Config) bool {
	return validateProductionEnvironment(config)
}

// checkIsDevelopment is an alias for validateDevelopmentEnvironment (for backwards compatibility)
func checkIsDevelopment(config Config) bool {
	return validateDevelopmentEnvironment(config)
}

// getDefaultGraphServiceConfig is an alias for buildDefaultGraphServiceConfig with environment vars
func getDefaultGraphServiceConfig() GraphServiceConfig {
	envVars := map[string]string{
		"ENVIRONMENT":      getEnvOrDefault("ENVIRONMENT", "development"),
		"NEO4J_URI":        getEnvOrDefault("NEO4J_URI", ""),
		"NEO4J_USERNAME":   getEnvOrDefault("NEO4J_USERNAME", ""),
		"NEO4J_PASSWORD":   getEnvOrDefault("NEO4J_PASSWORD", ""),
		"NEPTUNE_ENDPOINT": getEnvOrDefault("NEPTUNE_ENDPOINT", ""),
		"NEPTUNE_REGION":   getEnvOrDefault("NEPTUNE_REGION", ""),
	}
	return buildDefaultGraphServiceConfig(envVars)
}

// validateConfig validates the configuration
func validateConfig(config Config) error {
	if config.GraphDB.Provider == "" {
		return fmt.Errorf("graph database provider is required")
	}

	return validateGraphDBProvider(config)
}

// validateGraphDBProvider validates the graph database provider configuration
func validateGraphDBProvider(config Config) error {
	if config.GraphDB.Provider == "" {
		panic("GraphDB provider cannot be empty")
	}

	switch config.GraphDB.Provider {
	case providerNeo4j:
		return validateNeo4jConfig(config)
	case providerNeptune:
		return validateNeptuneConfig(config)
	default:
		return fmt.Errorf("unsupported graph database provider: %s", config.GraphDB.Provider)
	}
}

// validateNeo4jConfig validates Neo4j configuration
func validateNeo4jConfig(config Config) error {
	if config.GraphDB.Neo4j.URI == "" {
		return fmt.Errorf("Neo4j URI is required")
	}
	if config.GraphDB.Neo4j.Username == "" {
		return fmt.Errorf("Neo4j username is required")
	}
	if config.GraphDB.Neo4j.Password == "" {
		return fmt.Errorf("Neo4j password is required")
	}
	return nil
}

// validateNeptuneConfig validates Neptune configuration
func validateNeptuneConfig(config Config) error {
	if config.GraphDB.Neptune.Endpoint == "" {
		return fmt.Errorf("Neptune endpoint is required")
	}
	if config.GraphDB.Neptune.Region == "" {
		return fmt.Errorf("Neptune region is required")
	}
	return nil
}

// buildConfigFromJSON builds configuration from JSON data (Pure Core)
func buildConfigFromJSON(data []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &config, nil
}

// serializeConfigToJSON serializes configuration to JSON (Pure Core)
func serializeConfigToJSON(config *Config) ([]byte, error) {
	if config == nil {
		panic("Config cannot be nil")
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}
	return data, nil
}

// readEnvironmentVariables reads environment variables (Impure Shell)
func readEnvironmentVariables() map[string]string {
	envVars := make(map[string]string)
	envKeys := []string{
		"ENVIRONMENT", "GRAPH_DB_PROVIDER",
		"NEO4J_URI", "NEO4J_USERNAME", "NEO4J_PASSWORD",
		"NEPTUNE_ENDPOINT", "NEPTUNE_REGION",
	}

	for _, key := range envKeys {
		if value := os.Getenv(key); value != "" {
			envVars[key] = value
		}
	}
	return envVars
}

// readConfigFile reads configuration file (Impure Shell)
func readConfigFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return data, nil
}

// writeConfigFile writes configuration file (Impure Shell)
func writeConfigFile(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
