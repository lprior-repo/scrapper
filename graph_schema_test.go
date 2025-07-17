package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGraphSchema(t *testing.T) {
	schema := getGraphSchema()

	// Test that schema has indexes and constraints
	assert.NotEmpty(t, schema.Indexes, "Schema should have indexes")
	assert.NotEmpty(t, schema.Constraints, "Schema should have constraints")

	// Test that required indexes exist
	indexNames := make(map[string]bool)
	for _, index := range schema.Indexes {
		indexNames[index.Name] = true
	}

	requiredIndexes := []string{
		"idx_organization_name",
		"idx_repository_full_name",
		"idx_user_username",
		"idx_team_name",
	}

	for _, required := range requiredIndexes {
		assert.True(t, indexNames[required], "Required index %s should exist", required)
	}

	// Test that required constraints exist
	constraintNames := make(map[string]bool)
	for _, constraint := range schema.Constraints {
		constraintNames[constraint.Name] = true
	}

	requiredConstraints := []string{
		"constraint_organization_name_unique",
		"constraint_repository_full_name_unique",
		"constraint_user_username_unique",
	}

	for _, required := range requiredConstraints {
		assert.True(t, constraintNames[required], "Required constraint %s should exist", required)
	}
}

func TestBuildCreateIndexQuery(t *testing.T) {
	tests := []struct {
		name     string
		index    IndexDefinition
		expected string
	}{
		{
			name: "single property index",
			index: IndexDefinition{
				Name:       "idx_test",
				Label:      "TestNode",
				Properties: []string{"property1"},
				Type:       "btree",
			},
			expected: "CREATE INDEX idx_test IF NOT EXISTS FOR (n:TestNode) ON (n.property1)",
		},
		{
			name: "composite index",
			index: IndexDefinition{
				Name:       "idx_composite",
				Label:      "TestNode",
				Properties: []string{"prop1", "prop2"},
				Type:       "composite",
			},
			expected: "CREATE INDEX idx_composite IF NOT EXISTS FOR (n:TestNode) ON (n.prop1, n.prop2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildCreateIndexQuery(tt.index)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildCreateConstraintQuery(t *testing.T) {
	tests := []struct {
		name       string
		constraint ConstraintDefinition
		expected   string
	}{
		{
			name: "unique constraint single property",
			constraint: ConstraintDefinition{
				Name:       "constraint_test_unique",
				Label:      "TestNode",
				Properties: []string{"property1"},
				Type:       "unique",
			},
			expected: "CREATE CONSTRAINT constraint_test_unique IF NOT EXISTS FOR (n:TestNode) REQUIRE n.property1 IS UNIQUE",
		},
		{
			name: "unique constraint multiple properties",
			constraint: ConstraintDefinition{
				Name:       "constraint_composite_unique",
				Label:      "TestNode",
				Properties: []string{"prop1", "prop2"},
				Type:       "unique",
			},
			expected: "CREATE CONSTRAINT constraint_composite_unique IF NOT EXISTS FOR (n:TestNode) REQUIRE (n.prop1, n.prop2) IS UNIQUE",
		},
		{
			name: "exists constraint",
			constraint: ConstraintDefinition{
				Name:       "constraint_test_exists",
				Label:      "TestNode",
				Properties: []string{"property1"},
				Type:       "exists",
			},
			expected: "CREATE CONSTRAINT constraint_test_exists IF NOT EXISTS FOR (n:TestNode) REQUIRE n.property1 IS NOT NULL",
		},
		{
			name: "node key constraint",
			constraint: ConstraintDefinition{
				Name:       "constraint_test_key",
				Label:      "TestNode",
				Properties: []string{"prop1", "prop2"},
				Type:       "key",
			},
			expected: "CREATE CONSTRAINT constraint_test_key IF NOT EXISTS FOR (n:TestNode) REQUIRE (n.prop1, n.prop2) IS NODE KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildCreateConstraintQuery(tt.constraint)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildDropQueries(t *testing.T) {
	t.Run("drop index", func(t *testing.T) {
		result := buildDropIndexQuery("idx_test")
		expected := "DROP INDEX idx_test IF EXISTS"
		assert.Equal(t, expected, result)
	})

	t.Run("drop constraint", func(t *testing.T) {
		result := buildDropConstraintQuery("constraint_test")
		expected := "DROP CONSTRAINT constraint_test IF EXISTS"
		assert.Equal(t, expected, result)
	})
}

func TestBuildCreateConstraintQueryPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "Unknown constraint type")
		} else {
			t.Error("Expected panic for unknown constraint type")
		}
	}()

	constraint := ConstraintDefinition{
		Name:       "test",
		Label:      "Test",
		Properties: []string{"prop"},
		Type:       "invalid_type",
	}

	buildCreateConstraintQuery(constraint)
}

func TestIndexDefinitionValidation(t *testing.T) {
	t.Run("valid index definition", func(t *testing.T) {
		index := IndexDefinition{
			Name:       "idx_valid",
			Label:      "ValidNode",
			Properties: []string{"property1"},
			Type:       "btree",
		}

		assert.NotEmpty(t, index.Name)
		assert.NotEmpty(t, index.Label)
		assert.NotEmpty(t, index.Properties)
		assert.NotEmpty(t, index.Type)
	})
}

func TestConstraintDefinitionValidation(t *testing.T) {
	t.Run("valid constraint definition", func(t *testing.T) {
		constraint := ConstraintDefinition{
			Name:       "constraint_valid",
			Label:      "ValidNode",
			Properties: []string{"property1"},
			Type:       "unique",
		}

		assert.NotEmpty(t, constraint.Name)
		assert.NotEmpty(t, constraint.Label)
		assert.NotEmpty(t, constraint.Properties)
		assert.NotEmpty(t, constraint.Type)
	})
}
