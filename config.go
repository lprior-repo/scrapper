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

// LoadConfig loads configuration from environment and defaults
func LoadConfig() (*Config, error) {
	config := &Config{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		GraphDB:     getDefaultGraphServiceConfig(),
	}

	applyEnvironmentOverrides(config)
	return config, nil
}

// applyEnvironmentOverrides applies environment variable overrides to config
func applyEnvironmentOverrides(config *Config) {
	if config == nil {
		panic("Config cannot be nil")
	}
	if provider := os.Getenv("GRAPH_DB_PROVIDER"); provider != "" {
		config.GraphDB.Provider = provider
	}

	applyNeo4jOverrides(config)
	applyNeptuneOverrides(config)
}

// applyNeo4jOverrides applies Neo4j environment variable overrides
func applyNeo4jOverrides(config *Config) {
	if config == nil {
		panic("Config cannot be nil")
	}
	if uri := os.Getenv("NEO4J_URI"); uri != "" {
		config.GraphDB.Neo4j.URI = uri
	}
	if username := os.Getenv("NEO4J_USERNAME"); username != "" {
		config.GraphDB.Neo4j.Username = username
	}
	if password := os.Getenv("NEO4J_PASSWORD"); password != "" {
		config.GraphDB.Neo4j.Password = password
	}
}

// applyNeptuneOverrides applies Neptune environment variable overrides
func applyNeptuneOverrides(config *Config) {
	if config == nil {
		panic("Config cannot be nil")
	}
	if endpoint := os.Getenv("NEPTUNE_ENDPOINT"); endpoint != "" {
		config.GraphDB.Neptune.Endpoint = endpoint
	}
	if region := os.Getenv("NEPTUNE_REGION"); region != "" {
		config.GraphDB.Neptune.Region = region
	}
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(filename string) (*Config, error) {
	if filename == "" {
		panic("Filename cannot be empty")
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfigToFile saves configuration to a JSON file
func SaveConfigToFile(config *Config, filename string) error {
	if config == nil {
		panic("Config cannot be nil")
	}
	if filename == "" {
		panic("Filename cannot be empty")
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getEnvOrDefault retrieves environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if key == "" {
		panic("Environment variable key cannot be empty")
	}
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDefaultGraphServiceConfig returns default graph service configuration
func getDefaultGraphServiceConfig() GraphServiceConfig {
	environment := getEnvOrDefault("ENVIRONMENT", "development")

	if environment == envProduction {
		return GraphServiceConfig{
			Provider: providerNeptune,
			Neptune: struct {
				Endpoint string `json:"endpoint"`
				Region   string `json:"region"`
			}{
				Endpoint: getEnvOrDefault("NEPTUNE_ENDPOINT", ""),
				Region:   getEnvOrDefault("NEPTUNE_REGION", "us-east-1"),
			},
		}
	}

	return GraphServiceConfig{
		Provider: providerNeo4j,
		Neo4j: struct {
			URI      string `json:"uri"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			URI:      getEnvOrDefault("NEO4J_URI", "bolt://localhost:7687"),
			Username: getEnvOrDefault("NEO4J_USERNAME", defaultNeo4jUsername),
			Password: getEnvOrDefault("NEO4J_PASSWORD", "password"),
		},
	}
}

// checkIsProduction checks if the application is running in production
func checkIsProduction(config Config) bool {
	if config.Environment == "" {
		panic("Config environment cannot be empty")
	}
	return config.Environment == envProduction
}

// checkIsDevelopment checks if the application is running in development
func checkIsDevelopment(config Config) bool {
	if config.Environment == "" {
		panic("Config environment cannot be empty")
	}
	return config.Environment == "development"
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
