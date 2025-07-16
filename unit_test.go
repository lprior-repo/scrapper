package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestParseNodeID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "valid positive number",
			input:    "123",
			expected: 123,
		},
		{
			name:     "valid zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "valid negative number",
			input:    "-456",
			expected: -456,
		},
		{
			name:     "invalid string",
			input:    "invalid",
			expected: 0,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "very large number",
			input:    "9223372036854775807",
			expected: 9223372036854775807,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := parseNodeID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_KEY_1",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "environment variable empty",
			key:          "TEST_KEY_2",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "environment variable not set",
			key:          "TEST_KEY_3",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty default value",
			key:          "TEST_KEY_4",
			defaultValue: "",
			envValue:     "env_value",
			expected:     "env_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Clean up any existing env var
			_ = os.Unsetenv(tt.key)

			// Set env var if needed
			if tt.envValue != "" {
				_ = os.Setenv(tt.key, tt.envValue)
				defer func() { _ = os.Unsetenv(tt.key) }()
			}

			result := getEnvOrDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		config    *Config
		shouldErr bool
		errMsg    string
	}{
		{
			name: "valid neo4j config",
			config: &Config{
				Environment: "development",
				GraphDB: GraphServiceConfig{
					Provider: "neo4j",
					Neo4j: struct {
						URI      string `json:"uri"`
						Username string `json:"username"`
						Password string `json:"password"`
					}{
						URI:      "bolt://localhost:7687",
						Username: "neo4j",
						Password: "password",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "valid neptune config",
			config: &Config{
				Environment: "production",
				GraphDB: GraphServiceConfig{
					Provider: "neptune",
					Neptune: struct {
						Endpoint string `json:"endpoint"`
						Region   string `json:"region"`
					}{
						Endpoint: "wss://test.neptune.amazonaws.com:8182/gremlin",
						Region:   "us-east-1",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "empty provider",
			config: &Config{
				Environment: "development",
				GraphDB: GraphServiceConfig{
					Provider: "",
				},
			},
			shouldErr: true,
			errMsg:    "graph database provider is required",
		},
		{
			name: "unsupported provider",
			config: &Config{
				Environment: "development",
				GraphDB: GraphServiceConfig{
					Provider: "unsupported",
				},
			},
			shouldErr: true,
			errMsg:    "unsupported graph database provider",
		},
		{
			name: "neo4j missing uri",
			config: &Config{
				Environment: "development",
				GraphDB: GraphServiceConfig{
					Provider: "neo4j",
					Neo4j: struct {
						URI      string `json:"uri"`
						Username string `json:"username"`
						Password string `json:"password"`
					}{
						URI:      "",
						Username: "neo4j",
						Password: "password",
					},
				},
			},
			shouldErr: true,
			errMsg:    "Neo4j URI is required",
		},
		{
			name: "neo4j missing username",
			config: &Config{
				Environment: "development",
				GraphDB: GraphServiceConfig{
					Provider: "neo4j",
					Neo4j: struct {
						URI      string `json:"uri"`
						Username string `json:"username"`
						Password string `json:"password"`
					}{
						URI:      "bolt://localhost:7687",
						Username: "",
						Password: "password",
					},
				},
			},
			shouldErr: true,
			errMsg:    "Neo4j username is required",
		},
		{
			name: "neo4j missing password",
			config: &Config{
				Environment: "development",
				GraphDB: GraphServiceConfig{
					Provider: "neo4j",
					Neo4j: struct {
						URI      string `json:"uri"`
						Username string `json:"username"`
						Password string `json:"password"`
					}{
						URI:      "bolt://localhost:7687",
						Username: "neo4j",
						Password: "",
					},
				},
			},
			shouldErr: true,
			errMsg:    "Neo4j password is required",
		},
		{
			name: "neptune missing endpoint",
			config: &Config{
				Environment: "production",
				GraphDB: GraphServiceConfig{
					Provider: "neptune",
					Neptune: struct {
						Endpoint string `json:"endpoint"`
						Region   string `json:"region"`
					}{
						Endpoint: "",
						Region:   "us-east-1",
					},
				},
			},
			shouldErr: true,
			errMsg:    "Neptune endpoint is required",
		},
		{
			name: "neptune missing region",
			config: &Config{
				Environment: "production",
				GraphDB: GraphServiceConfig{
					Provider: "neptune",
					Neptune: struct {
						Endpoint string `json:"endpoint"`
						Region   string `json:"region"`
					}{
						Endpoint: "wss://test.neptune.amazonaws.com:8182/gremlin",
						Region:   "",
					},
				},
			},
			shouldErr: true,
			errMsg:    "Neptune region is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateConfig(*tt.config)

			if tt.shouldErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigEnvironmentChecks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		environment           string
		expectedIsProduction  bool
		expectedIsDevelopment bool
	}{
		{
			name:                  "production environment",
			environment:           "production",
			expectedIsProduction:  true,
			expectedIsDevelopment: false,
		},
		{
			name:                  "development environment",
			environment:           "development",
			expectedIsProduction:  false,
			expectedIsDevelopment: true,
		},
		{
			name:                  "staging environment",
			environment:           "staging",
			expectedIsProduction:  false,
			expectedIsDevelopment: false,
		},
		{
			name:                  "empty environment",
			environment:           "",
			expectedIsProduction:  false,
			expectedIsDevelopment: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := &Config{
				Environment: tt.environment,
			}

			assert.Equal(t, tt.expectedIsProduction, checkIsProduction(*config))
			assert.Equal(t, tt.expectedIsDevelopment, checkIsDevelopment(*config))
		})
	}
}

func TestGetDefaultGraphServiceConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		environment string
		expected    string
	}{
		{
			name:        "production environment should use neptune",
			environment: "production",
			expected:    "neptune",
		},
		{
			name:        "development environment should use neo4j",
			environment: "development",
			expected:    "neo4j",
		},
		{
			name:        "staging environment should use neo4j",
			environment: "staging",
			expected:    "neo4j",
		},
		{
			name:        "empty environment should use neo4j",
			environment: "",
			expected:    "neo4j",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Clean up any existing env var
			_ = os.Unsetenv("ENVIRONMENT")

			// Set environment if needed
			if tt.environment != "" {
				_ = os.Setenv("ENVIRONMENT", tt.environment)
				defer func() { _ = os.Unsetenv("ENVIRONMENT") }()
			}

			config := getDefaultGraphServiceConfig()
			assert.Equal(t, tt.expected, config.Provider)
		})
	}
}
