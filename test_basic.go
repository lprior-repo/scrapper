package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestBasicConfigFunctions tests basic configuration functions
func TestBasicConfigFunctions(t *testing.T) {
	// Test getStringOrDefault
	envVars := map[string]string{
		"TEST_KEY": "test_value",
	}
	
	result := getStringOrDefault(envVars, "TEST_KEY", "default")
	assert.Equal(t, "test_value", result)
	
	result = getStringOrDefault(envVars, "MISSING_KEY", "default")
	assert.Equal(t, "default", result)
}

// TestBasicGraphTypes tests basic graph types
func TestBasicGraphTypes(t *testing.T) {
	// Test Node creation
	node := &Node{
		ID:     "1",
		Labels: []string{"TestLabel"},
		Properties: map[string]interface{}{
			"name": "test",
		},
	}
	
	assert.Equal(t, "1", node.ID)
	assert.Equal(t, []string{"TestLabel"}, node.Labels)
	assert.Equal(t, "test", node.Properties["name"])
}

// TestBasicUtilityFunctions tests basic utility functions
func TestBasicUtilityFunctions(t *testing.T) {
	// Test convertStringToInt64
	result := convertStringToInt64("123")
	assert.Equal(t, int64(123), result)
	
	result = convertStringToInt64("")
	assert.Equal(t, int64(0), result)
	
	result = convertStringToInt64("invalid")
	assert.Equal(t, int64(0), result)
}