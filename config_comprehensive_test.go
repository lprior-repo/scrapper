package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Comprehensive tests to catch mutations and improve mutation score

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	
	// Test basic functionality
	config, err := LoadConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "development", config.Environment)
	assert.NotEmpty(t, config.GraphDB.Provider)
}

func TestLoadConfigWithEnvironmentOverrides(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name     string
		envVars  map[string]string
		expected func(*Config)
	}{
		{
			name: "environment override",
			envVars: map[string]string{
				"ENVIRONMENT": "production",
			},
			expected: func(c *Config) {
				assert.Equal(t, "production", c.Environment)
			},
		},
		{
			name: "graph provider override",
			envVars: map[string]string{
				"GRAPH_DB_PROVIDER": "neptune",
			},
			expected: func(c *Config) {
				assert.Equal(t, "neptune", c.GraphDB.Provider)
			},
		},
		{
			name: "neo4j overrides",
			envVars: map[string]string{
				"NEO4J_URI":      "bolt://custom:7687",
				"NEO4J_USERNAME": "custom_user",
				"NEO4J_PASSWORD": "custom_pass",
			},
			expected: func(c *Config) {
				assert.Equal(t, "bolt://custom:7687", c.GraphDB.Neo4j.URI)
				assert.Equal(t, "custom_user", c.GraphDB.Neo4j.Username)
				assert.Equal(t, "custom_pass", c.GraphDB.Neo4j.Password)
			},
		},
		{
			name: "neptune overrides",
			envVars: map[string]string{
				"NEPTUNE_ENDPOINT": "wss://custom.neptune.amazonaws.com:8182/gremlin",
				"NEPTUNE_REGION":   "us-west-2",
			},
			expected: func(c *Config) {
				assert.Equal(t, "wss://custom.neptune.amazonaws.com:8182/gremlin", c.GraphDB.Neptune.Endpoint)
				assert.Equal(t, "us-west-2", c.GraphDB.Neptune.Region)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// Clean up and set environment variables
			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
				defer func(k string) { _ = os.Unsetenv(k) }(key)
			}
			
			config, err := LoadConfig()
			require.NoError(t, err)
			tt.expected(config)
		})
	}
}

func TestApplyEnvironmentOverrides(t *testing.T) {
	t.Parallel()
	
	config := &Config{
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
	}
	
	// Test with provider override
	_ = os.Setenv("GRAPH_DB_PROVIDER", "neptune")
	defer func() { _ = os.Unsetenv("GRAPH_DB_PROVIDER") }()
	
	applyEnvironmentOverrides(config)
	assert.Equal(t, "neptune", config.GraphDB.Provider)
}

func TestApplyNeo4jOverrides(t *testing.T) {
	t.Parallel()
	
	// Test each override individually
	tests := []struct {
		name    string
		envVar  string
		envVal  string
		checker func(*Config)
	}{
		{
			name:   "URI override",
			envVar: "NEO4J_URI",
			envVal: "bolt://custom:7687",
			checker: func(c *Config) {
				assert.Equal(t, "bolt://custom:7687", c.GraphDB.Neo4j.URI)
			},
		},
		{
			name:   "Username override",
			envVar: "NEO4J_USERNAME",
			envVal: "custom_user",
			checker: func(c *Config) {
				assert.Equal(t, "custom_user", c.GraphDB.Neo4j.Username)
			},
		},
		{
			name:   "Password override",
			envVar: "NEO4J_PASSWORD",
			envVal: "custom_pass",
			checker: func(c *Config) {
				assert.Equal(t, "custom_pass", c.GraphDB.Neo4j.Password)
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// Reset config
			testConfig := &Config{
				GraphDB: GraphServiceConfig{
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
			}
			
			_ = os.Setenv(tt.envVar, tt.envVal)
			defer func() { _ = os.Unsetenv(tt.envVar) }()
			
			applyNeo4jOverrides(testConfig)
			tt.checker(testConfig)
		})
	}
}

func TestApplyNeptuneOverrides(t *testing.T) {
	t.Parallel()
	
	// Test each override individually
	tests := []struct {
		name    string
		envVar  string
		envVal  string
		checker func(*Config)
	}{
		{
			name:   "Endpoint override",
			envVar: "NEPTUNE_ENDPOINT",
			envVal: "wss://custom.neptune.amazonaws.com:8182/gremlin",
			checker: func(c *Config) {
				assert.Equal(t, "wss://custom.neptune.amazonaws.com:8182/gremlin", c.GraphDB.Neptune.Endpoint)
			},
		},
		{
			name:   "Region override",
			envVar: "NEPTUNE_REGION",
			envVal: "us-west-2",
			checker: func(c *Config) {
				assert.Equal(t, "us-west-2", c.GraphDB.Neptune.Region)
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// Reset config
			testConfig := &Config{
				GraphDB: GraphServiceConfig{
					Neptune: struct {
						Endpoint string `json:"endpoint"`
						Region   string `json:"region"`
					}{
						Endpoint: "wss://test.neptune.amazonaws.com:8182/gremlin",
						Region:   "us-east-1",
					},
				},
			}
			
			_ = os.Setenv(tt.envVar, tt.envVal)
			defer func() { _ = os.Unsetenv(tt.envVar) }()
			
			applyNeptuneOverrides(testConfig)
			tt.checker(testConfig)
		})
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Parallel()
	
	// Test successful file loading
	t.Run("successful file loading", func(t *testing.T) {
		// Create a temporary config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.json")
		
		testConfig := &Config{
			Environment: "test",
			GraphDB: GraphServiceConfig{
				Provider: "neo4j",
				Neo4j: struct {
					URI      string `json:"uri"`
					Username string `json:"username"`
					Password string `json:"password"`
				}{
					URI:      "bolt://test:7687",
					Username: "test_user",
					Password: "test_pass",
				},
			},
		}
		
		data, err := json.MarshalIndent(testConfig, "", "  ")
		require.NoError(t, err)
		
		err = os.WriteFile(configFile, data, 0644)
		require.NoError(t, err)
		
		// Test loading
		config, err := LoadConfigFromFile(configFile)
		require.NoError(t, err)
		assert.Equal(t, "test", config.Environment)
		assert.Equal(t, "neo4j", config.GraphDB.Provider)
		assert.Equal(t, "bolt://test:7687", config.GraphDB.Neo4j.URI)
	})
	
	// Test file not found
	t.Run("file not found", func(t *testing.T) {
		_, err := LoadConfigFromFile("/nonexistent/config.json")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config file")
	})
	
	// Test invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "invalid.json")
		
		err := os.WriteFile(configFile, []byte("invalid json"), 0644)
		require.NoError(t, err)
		
		_, err = LoadConfigFromFile(configFile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config file")
	})
}

func TestSaveConfigToFile(t *testing.T) {
	t.Parallel()
	
	// Test successful file saving
	t.Run("successful file saving", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.json")
		
		testConfig := &Config{
			Environment: "test",
			GraphDB: GraphServiceConfig{
				Provider: "neo4j",
				Neo4j: struct {
					URI      string `json:"uri"`
					Username string `json:"username"`
					Password string `json:"password"`
				}{
					URI:      "bolt://test:7687",
					Username: "test_user",
					Password: "test_pass",
				},
			},
		}
		
		err := SaveConfigToFile(testConfig, configFile)
		require.NoError(t, err)
		
		// Verify file was created and contains correct data
		data, err := os.ReadFile(configFile)
		require.NoError(t, err)
		
		var loadedConfig Config
		err = json.Unmarshal(data, &loadedConfig)
		require.NoError(t, err)
		
		assert.Equal(t, testConfig.Environment, loadedConfig.Environment)
		assert.Equal(t, testConfig.GraphDB.Provider, loadedConfig.GraphDB.Provider)
	})
	
	// Test write error (invalid path)
	t.Run("write error", func(t *testing.T) {
		testConfig := &Config{
			Environment: "test",
			GraphDB: GraphServiceConfig{
				Provider: "neo4j",
			},
		}
		
		err := SaveConfigToFile(testConfig, "/nonexistent/path/config.json")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write config file")
	})
}

func TestGetDefaultGraphServiceConfigComprehensive(t *testing.T) {
	t.Parallel()
	
	// Test production environment path
	t.Run("production environment", func(t *testing.T) {
		_ = os.Setenv("ENVIRONMENT", "production")
		defer func() { _ = os.Unsetenv("ENVIRONMENT") }()
		
		config := getDefaultGraphServiceConfig()
		assert.Equal(t, "neptune", config.Provider)
		assert.NotEmpty(t, config.Neptune.Region) // Should have default region
	})
	
	// Test non-production environment path
	t.Run("non-production environment", func(t *testing.T) {
		_ = os.Setenv("ENVIRONMENT", "development")
		defer func() { _ = os.Unsetenv("ENVIRONMENT") }()
		
		config := getDefaultGraphServiceConfig()
		assert.Equal(t, "neo4j", config.Provider)
		assert.NotEmpty(t, config.Neo4j.URI)
		assert.NotEmpty(t, config.Neo4j.Username)
		assert.NotEmpty(t, config.Neo4j.Password)
	})
	
	// Test with custom environment variables
	t.Run("with custom environment variables", func(t *testing.T) {
		_ = os.Setenv("ENVIRONMENT", "production")
		_ = os.Setenv("NEPTUNE_ENDPOINT", "custom-endpoint")
		_ = os.Setenv("NEPTUNE_REGION", "custom-region")
		defer func() { 
			_ = os.Unsetenv("ENVIRONMENT") 
			_ = os.Unsetenv("NEPTUNE_ENDPOINT")
			_ = os.Unsetenv("NEPTUNE_REGION")
		}()
		
		config := getDefaultGraphServiceConfig()
		assert.Equal(t, "neptune", config.Provider)
		assert.Equal(t, "custom-endpoint", config.Neptune.Endpoint)
		assert.Equal(t, "custom-region", config.Neptune.Region)
	})
}

func TestValidateGraphDBProvider(t *testing.T) {
	t.Parallel()
	
	// Test valid Neo4j provider
	t.Run("valid neo4j provider", func(t *testing.T) {
		config := &Config{
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
		}
		
		err := validateGraphDBProvider(*config)
		assert.NoError(t, err)
	})
	
	// Test valid Neptune provider
	t.Run("valid neptune provider", func(t *testing.T) {
		config := &Config{
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
		}
		
		err := validateGraphDBProvider(*config)
		assert.NoError(t, err)
	})
	
	// Test invalid provider
	t.Run("invalid provider", func(t *testing.T) {
		config := &Config{
			GraphDB: GraphServiceConfig{
				Provider: "invalid",
			},
		}
		
		err := validateGraphDBProvider(*config)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported graph database provider")
	})
}

func TestConfigEnvironmentVariableEdgeCases(t *testing.T) {
	t.Parallel()
	
	// Test empty environment variable handling
	t.Run("empty environment variables", func(t *testing.T) {
		// Set empty environment variables
		_ = os.Setenv("GRAPH_DB_PROVIDER", "")
		_ = os.Setenv("NEO4J_URI", "")
		_ = os.Setenv("NEO4J_USERNAME", "")
		_ = os.Setenv("NEO4J_PASSWORD", "")
		_ = os.Setenv("NEPTUNE_ENDPOINT", "")
		_ = os.Setenv("NEPTUNE_REGION", "")
		
		defer func() {
			_ = os.Unsetenv("GRAPH_DB_PROVIDER")
			_ = os.Unsetenv("NEO4J_URI")
			_ = os.Unsetenv("NEO4J_USERNAME")
			_ = os.Unsetenv("NEO4J_PASSWORD")
			_ = os.Unsetenv("NEPTUNE_ENDPOINT")
			_ = os.Unsetenv("NEPTUNE_REGION")
		}()
		
		config := &Config{
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
		}
		
		applyEnvironmentOverrides(config)
		
		// Empty values should not override existing values
		assert.Equal(t, "neo4j", config.GraphDB.Provider)
		assert.Equal(t, "bolt://localhost:7687", config.GraphDB.Neo4j.URI)
		assert.Equal(t, "neo4j", config.GraphDB.Neo4j.Username)
		assert.Equal(t, "password", config.GraphDB.Neo4j.Password)
	})
}