package main

import "fmt"

// SchemaDefinition holds the graph schema definitions
type SchemaDefinition struct {
	Indexes     []IndexDefinition      `json:"indexes"`
	Constraints []ConstraintDefinition `json:"constraints"`
}

// IndexDefinition represents a graph index
type IndexDefinition struct {
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Properties []string `json:"properties"`
	Type       string   `json:"type"` // "btree", "text", "composite"
}

// ConstraintDefinition represents a graph constraint
type ConstraintDefinition struct {
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Properties []string `json:"properties"`
	Type       string   `json:"type"` // "unique", "exists", "key"
}

// getGraphSchema returns the complete graph schema definition (Pure Core)
func getGraphSchema() SchemaDefinition {
	return SchemaDefinition{
		Indexes: []IndexDefinition{
			{
				Name:       "idx_organization_name",
				Label:      "Organization",
				Properties: []string{"name"},
				Type:       "btree",
			},
			{
				Name:       "idx_repository_full_name",
				Label:      "Repository",
				Properties: []string{"full_name"},
				Type:       "btree",
			},
			{
				Name:       "idx_repository_name",
				Label:      "Repository",
				Properties: []string{"name"},
				Type:       "btree",
			},
			{
				Name:       "idx_user_username",
				Label:      "User",
				Properties: []string{"username"},
				Type:       "btree",
			},
			{
				Name:       "idx_user_email",
				Label:      "User",
				Properties: []string{"email"},
				Type:       "btree",
			},
			{
				Name:       "idx_team_name",
				Label:      "Team",
				Properties: []string{"name"},
				Type:       "btree",
			},
			{
				Name:       "idx_org_repo_composite",
				Label:      "Repository",
				Properties: []string{"organization", "name"},
				Type:       "composite",
			},
			{
				Name:       "idx_codeowners_text",
				Label:      "Repository",
				Properties: []string{"codeowners_content"},
				Type:       "text",
			},
		},
		Constraints: []ConstraintDefinition{
			{
				Name:       "constraint_organization_name_unique",
				Label:      "Organization",
				Properties: []string{"name"},
				Type:       "unique",
			},
			{
				Name:       "constraint_repository_full_name_unique",
				Label:      "Repository",
				Properties: []string{"full_name"},
				Type:       "unique",
			},
			{
				Name:       "constraint_user_username_unique",
				Label:      "User",
				Properties: []string{"username"},
				Type:       "unique",
			},
			{
				Name:       "constraint_user_email_unique",
				Label:      "User",
				Properties: []string{"email"},
				Type:       "unique",
			},
			{
				Name:       "constraint_team_name_org_unique",
				Label:      "Team",
				Properties: []string{"name", "organization"},
				Type:       "unique",
			},
		},
	}
}

// buildCreateIndexQuery creates a Cypher query for index creation (Pure Core)
func buildCreateIndexQuery(index IndexDefinition) string {
	if len(index.Properties) == 1 {
		return fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (n.%s)",
			index.Name, index.Label, index.Properties[0])
	}

	// Composite index
	properties := ""
	for i, prop := range index.Properties {
		if i > 0 {
			properties += ", "
		}
		properties += fmt.Sprintf("n.%s", prop)
	}

	return fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (%s)",
		index.Name, index.Label, properties)
}

// buildCreateConstraintQuery creates a Cypher query for constraint creation (Pure Core)
func buildCreateConstraintQuery(constraint ConstraintDefinition) string {
	switch constraint.Type {
	case "unique":
		if len(constraint.Properties) == 1 {
			return fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE n.%s IS UNIQUE",
				constraint.Name, constraint.Label, constraint.Properties[0])
		}
		// Multi-property unique constraint
		properties := ""
		for i, prop := range constraint.Properties {
			if i > 0 {
				properties += ", "
			}
			properties += fmt.Sprintf("n.%s", prop)
		}
		return fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE (%s) IS UNIQUE",
			constraint.Name, constraint.Label, properties)
	case "exists":
		return fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE n.%s IS NOT NULL",
			constraint.Name, constraint.Label, constraint.Properties[0])
	case "key":
		properties := ""
		for i, prop := range constraint.Properties {
			if i > 0 {
				properties += ", "
			}
			properties += fmt.Sprintf("n.%s", prop)
		}
		return fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE (%s) IS NODE KEY",
			constraint.Name, constraint.Label, properties)
	default:
		panic(fmt.Sprintf("Unknown constraint type: %s", constraint.Type))
	}
}

// buildDropIndexQuery creates a Cypher query for index removal (Pure Core)
func buildDropIndexQuery(indexName string) string {
	return fmt.Sprintf("DROP INDEX %s IF EXISTS", indexName)
}

// buildDropConstraintQuery creates a Cypher query for constraint removal (Pure Core)
func buildDropConstraintQuery(constraintName string) string {
	return fmt.Sprintf("DROP CONSTRAINT %s IF EXISTS", constraintName)
}
