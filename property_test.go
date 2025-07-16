package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

// Property-based tests for pure functions

func TestParseNodeIDProperties(t *testing.T) {
	t.Parallel()
	t.Run("parseNodeID_round_trip", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate a random int64
			original := rapid.Int64().Draw(t, "original")

			// Convert to string and parse back
			str := fmt.Sprintf("%d", original)
			parsed := parseNodeID(str)

			// Property: parsing a valid int64 string should return original value
			assert.Equal(t, original, parsed)
		})
	})

	t.Run("parseNodeID_invalid_strings_return_zero", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate random strings that are not valid integers
			invalidStr := rapid.StringMatching(`[a-zA-Z]+`).Draw(t, "invalid")

			// Property: invalid strings should always return 0
			result := parseNodeID(invalidStr)
			assert.Equal(t, int64(0), result)
		})
	})

	t.Run("parseNodeID_empty_string", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Property: empty string should always return 0
			result := parseNodeID("")
			assert.Equal(t, int64(0), result)
		})
	})
}

func TestGetEnvOrDefaultProperties(t *testing.T) {
	t.Parallel()
	t.Run("getEnvOrDefault_with_env_set", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate random key and values with unique prefix
			keyNum := rapid.Uint64Range(0, 999999).Draw(t, "keyNum")
			key := fmt.Sprintf("TEST_PROP_%d", keyNum)
			envValue := rapid.String().Filter(func(s string) bool {
				return s != "" && !strings.ContainsRune(s, '\x00') && len(s) < 100
			}).Draw(t, "envValue")
			defaultValue := rapid.String().Draw(t, "defaultValue")

			// Clean up
			defer func() { _ = os.Unsetenv(key) }()

			// Set environment variable
			_ = os.Setenv(key, envValue)

			// Property: when env var is set to non-empty value, it should return the env value
			result := getEnvOrDefault(key, defaultValue)
			assert.Equal(t, envValue, result)
		})
	})

	t.Run("getEnvOrDefault_without_env_set", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate random key and default value with unique prefix
			keyNum := rapid.Uint64Range(0, 999999).Draw(t, "keyNum")
			key := fmt.Sprintf("TEST_UNSET_%d", keyNum)
			defaultValue := rapid.String().Draw(t, "defaultValue")

			// Ensure env var is not set
			_ = os.Unsetenv(key)

			// Property: when env var is not set, it should return default
			result := getEnvOrDefault(key, defaultValue)
			assert.Equal(t, defaultValue, result)
		})
	})

	t.Run("getEnvOrDefault_empty_env_returns_default", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate random key and default value with unique prefix
			keyNum := rapid.Uint64Range(0, 999999).Draw(t, "keyNum")
			key := fmt.Sprintf("TEST_EMPTY_%d", keyNum)
			defaultValue := rapid.String().Draw(t, "defaultValue")

			// Clean up
			defer func() { _ = os.Unsetenv(key) }()

			// Set environment variable to empty string
			_ = os.Setenv(key, "")

			// Property: when env var is empty, it should return default
			result := getEnvOrDefault(key, defaultValue)
			assert.Equal(t, defaultValue, result)
		})
	})
}

func TestConfigValidationProperties(t *testing.T) {
	t.Parallel()
	t.Run("valid_neo4j_config_always_passes", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate valid Neo4j config components
			uri := rapid.StringMatching(`bolt://[a-zA-Z0-9.-]+:[0-9]+`).Draw(t, "uri")
			username := rapid.StringMatching(`[a-zA-Z0-9_]+`).Draw(t, "username")
			password := rapid.StringMatching(`[a-zA-Z0-9_@#$%]+`).Draw(t, "password")

			config := &Config{
				Environment: envDevelopment,
				GraphDB: GraphServiceConfig{
					Provider: providerNeo4j,
					Neo4j: struct {
						URI      string `json:"uri"`
						Username string `json:"username"`
						Password string `json:"password"`
					}{
						URI:      uri,
						Username: username,
						Password: password,
					},
				},
			}

			// Property: valid Neo4j config should always pass validation
			err := validateConfig(*config)
			assert.NoError(t, err)
		})
	})

	t.Run("valid_neptune_config_always_passes", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate valid Neptune config components
			endpoint := rapid.StringMatching(`wss://[a-zA-Z0-9.-]+\.neptune\.amazonaws\.com:[0-9]+/gremlin`).Draw(t, "endpoint")
			region := rapid.SampledFrom([]string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}).Draw(t, "region")

			config := &Config{
				Environment: "production",
				GraphDB: GraphServiceConfig{
					Provider: "neptune",
					Neptune: struct {
						Endpoint string `json:"endpoint"`
						Region   string `json:"region"`
					}{
						Endpoint: endpoint,
						Region:   region,
					},
				},
			}

			// Property: valid Neptune config should always pass validation
			err := validateConfig(*config)
			assert.NoError(t, err)
		})
	})

	t.Run("empty_provider_always_fails", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			config := &Config{
				Environment: rapid.String().Draw(t, "environment"),
				GraphDB: GraphServiceConfig{
					Provider: "",
				},
			}

			// Property: empty provider should always fail validation
			err := validateConfig(*config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "graph database provider is required")
		})
	})

	t.Run("unsupported_provider_always_fails", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate invalid provider names
			invalidProvider := rapid.StringMatching(`[a-zA-Z]+`).
				Filter(func(s string) bool { return s != providerNeo4j && s != providerNeptune }).
				Draw(t, "provider")

			config := &Config{
				Environment: rapid.String().Draw(t, "environment"),
				GraphDB: GraphServiceConfig{
					Provider: invalidProvider,
				},
			}

			// Property: unsupported provider should always fail validation
			err := validateConfig(*config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "unsupported graph database provider")
		})
	})
}

func TestConfigEnvironmentProperties(t *testing.T) {
	t.Parallel()
	t.Run("production_environment_check", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			config := &Config{
				Environment: "production",
			}

			// Property: production environment should always return true for IsProduction
			assert.True(t, checkIsProduction(*config))
			assert.False(t, checkIsDevelopment(*config))
		})
	})

	t.Run("development_environment_check", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			config := &Config{
				Environment: envDevelopment,
			}

			// Property: development environment should always return true for IsDevelopment
			assert.False(t, checkIsProduction(*config))
			assert.True(t, checkIsDevelopment(*config))
		})
	})

	t.Run("non_production_non_development", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate environment names that are not production or development
			env := rapid.StringMatching(`[a-zA-Z]+`).
				Filter(func(s string) bool { return s != envProduction && s != envDevelopment }).
				Draw(t, "environment")

			config := &Config{
				Environment: env,
			}

			// Property: non-production/non-development should return false for both
			assert.False(t, checkIsProduction(*config))
			assert.False(t, checkIsDevelopment(*config))
		})
	})
}

func TestStringToInt64Properties(t *testing.T) {
	t.Parallel()
	t.Run("valid_int64_strings", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate valid int64 values
			original := rapid.Int64().Draw(t, "original")

			// Convert to string using standard library
			str := strconv.FormatInt(original, 10)

			// Parse back using our function
			parsed := parseNodeID(str)

			// Property: parsing valid int64 strings should be consistent
			assert.Equal(t, original, parsed)
		})
	})

	t.Run("parseNodeID_idempotent", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate any string
			str := rapid.String().Draw(t, "str")

			// Parse twice
			first := parseNodeID(str)
			second := parseNodeID(str)

			// Property: parseNodeID should be idempotent
			assert.Equal(t, first, second)
		})
	})
}

func TestBatchOperationProperties(t *testing.T) {
	t.Parallel()
	t.Run("batch_operation_structure", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate batch operation components
			opType := rapid.SampledFrom([]string{"create_node", "create_relationship", "delete_node"}).Draw(t, "type")
			query := rapid.String().Draw(t, "query")

			// Create batch operation
			op := BatchOperation{
				Type:       opType,
				Query:      query,
				Parameters: make(map[string]interface{}),
			}

			// Property: batch operation should maintain its structure
			assert.Equal(t, opType, op.Type)
			assert.Equal(t, query, op.Query)
			assert.NotNil(t, op.Parameters)
		})
	})
}

func TestNodeAndRelationshipProperties(t *testing.T) {
	t.Parallel()
	t.Run("node_structure_properties", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate node components
			id := rapid.String().Draw(t, "id")
			label := rapid.String().Draw(t, "label")

			node := &Node{
				ID:         id,
				Labels:     []string{label},
				Properties: make(map[string]interface{}),
			}

			// Property: node should maintain its structure
			assert.Equal(t, id, node.ID)
			assert.Contains(t, node.Labels, label)
			assert.NotNil(t, node.Properties)
		})
	})

	t.Run("relationship_structure_properties", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			// Generate relationship components
			id := rapid.String().Draw(t, "id")
			fromID := rapid.String().Draw(t, "fromID")
			toID := rapid.String().Draw(t, "toID")
			relType := rapid.String().Draw(t, "type")

			rel := &Relationship{
				ID:         id,
				Type:       relType,
				FromID:     fromID,
				ToID:       toID,
				Properties: make(map[string]interface{}),
			}

			// Property: relationship should maintain its structure
			assert.Equal(t, id, rel.ID)
			assert.Equal(t, relType, rel.Type)
			assert.Equal(t, fromID, rel.FromID)
			assert.Equal(t, toID, rel.ToID)
			assert.NotNil(t, rel.Properties)
		})
	})
}
