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
	if provider := os.Getenv("GRAPH_DB_PROVIDER"); provider != "" {
		config.GraphDB.Provider = provider
	}

	applyNeo4jOverrides(config)
	applyNeptuneOverrides(config)
}

// applyNeo4jOverrides applies Neo4j environment variable overrides
func applyNeo4jOverrides(config *Config) {
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
	if endpoint := os.Getenv("NEPTUNE_ENDPOINT"); endpoint != "" {
		config.GraphDB.Neptune.Endpoint = endpoint
	}
	if region := os.Getenv("NEPTUNE_REGION"); region != "" {
		config.GraphDB.Neptune.Region = region
	}
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(filename string) (*Config, error) {
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
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
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

// IsProduction checks if the application is running in production
func (c *Config) IsProduction() bool {
	return c.Environment == envProduction
}

// IsDevelopment checks if the application is running in development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.GraphDB.Provider == "" {
		return fmt.Errorf("graph database provider is required")
	}

	return c.validateGraphDBProvider()
}

// validateGraphDBProvider validates the graph database provider configuration
func (c *Config) validateGraphDBProvider() error {
	switch c.GraphDB.Provider {
	case providerNeo4j:
		return c.validateNeo4jConfig()
	case providerNeptune:
		return c.validateNeptuneConfig()
	default:
		return fmt.Errorf("unsupported graph database provider: %s", c.GraphDB.Provider)
	}
}

// validateNeo4jConfig validates Neo4j configuration
func (c *Config) validateNeo4jConfig() error {
	if c.GraphDB.Neo4j.URI == "" {
		return fmt.Errorf("Neo4j URI is required")
	}
	if c.GraphDB.Neo4j.Username == "" {
		return fmt.Errorf("Neo4j username is required")
	}
	if c.GraphDB.Neo4j.Password == "" {
		return fmt.Errorf("Neo4j password is required")
	}
	return nil
}

// validateNeptuneConfig validates Neptune configuration
func (c *Config) validateNeptuneConfig() error {
	if c.GraphDB.Neptune.Endpoint == "" {
		return fmt.Errorf("Neptune endpoint is required")
	}
	if c.GraphDB.Neptune.Region == "" {
		return fmt.Errorf("Neptune region is required")
	}
	return nil
}
