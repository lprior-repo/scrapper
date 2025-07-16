package main

import (
	"context"
	"fmt"
)

// GraphService defines the interface for graph database operations
type GraphService interface {
	// Connection management
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Health(ctx context.Context) error

	// Node operations
	CreateNode(ctx context.Context, label string, properties map[string]interface{}) (*Node, error)
	GetNode(ctx context.Context, id string) (*Node, error)
	UpdateNode(ctx context.Context, id string, properties map[string]interface{}) error
	DeleteNode(ctx context.Context, id string) error

	// Relationship operations
	CreateRelationship(ctx context.Context, fromID, toID, relType string, properties map[string]interface{}) (*Relationship, error)
	GetRelationship(ctx context.Context, id string) (*Relationship, error)
	DeleteRelationship(ctx context.Context, id string) error

	// Query operations
	ExecuteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error)
	ExecuteReadQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error)
	ExecuteWriteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error)

	// Batch operations
	ExecuteBatch(ctx context.Context, operations []BatchOperation) error

	// Utility operations
	ClearAll(ctx context.Context) error
}

// Node represents a graph node
type Node struct {
	ID         string                 `json:"id"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

// Relationship represents a graph relationship
type Relationship struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	FromID     string                 `json:"from_id"`
	ToID       string                 `json:"to_id"`
	Properties map[string]interface{} `json:"properties"`
}

// BatchOperation represents a batch operation
type BatchOperation struct {
	Type       string                 `json:"type"`       // "create_node", "create_relationship", "delete_node", etc.
	Query      string                 `json:"query"`      // Cypher query for the operation
	Parameters map[string]interface{} `json:"parameters"` // Parameters for the query
}

// GraphServiceConfig holds configuration for graph services
type GraphServiceConfig struct {
	Provider string `json:"provider"` // "neo4j" or "neptune"
	Neo4j    struct {
		URI      string `json:"uri"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"neo4j"`
	Neptune struct {
		Endpoint string `json:"endpoint"`
		Region   string `json:"region"`
	} `json:"neptune"`
}

// NewGraphService creates a new graph service based on configuration
func NewGraphService(config GraphServiceConfig) (GraphService, error) {
	switch config.Provider {
	case providerNeo4j:
		return NewNeo4jService(Neo4jConfig{
			URI:      config.Neo4j.URI,
			Username: config.Neo4j.Username,
			Password: config.Neo4j.Password,
		})
	case providerNeptune:
		return NewNeptuneService(NeptuneConfig{
			Endpoint: config.Neptune.Endpoint,
			Region:   config.Neptune.Region,
		})
	default:
		return nil, fmt.Errorf("unsupported graph service provider: %s", config.Provider)
	}
}

// NeptuneConfig holds Neptune-specific configuration
type NeptuneConfig struct {
	Endpoint string
	Region   string
}
