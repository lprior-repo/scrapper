package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/samber/lo"
)

// Neo4jService implements GraphService for Neo4j
type Neo4jService struct {
	driver neo4j.DriverWithContext
	config Neo4jConfig
}

// NewNeo4jService creates a new Neo4j service
func NewNeo4jService(config Neo4jConfig) (*Neo4jService, error) {
	return &Neo4jService{
		config: config,
	}, nil
}

// Connect establishes connection to Neo4j
func (s *Neo4jService) Connect(ctx context.Context) error {
	driver, err := createNeo4jDriver(s.config)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	s.driver = driver
	return s.Health(ctx)
}

// Close closes the Neo4j connection
func (s *Neo4jService) Close(ctx context.Context) error {
	if s.driver != nil {
		return s.driver.Close(ctx)
	}
	return nil
}

// Health checks Neo4j connection health
func (s *Neo4jService) Health(ctx context.Context) error {
	if s.driver == nil {
		return fmt.Errorf("driver not initialized")
	}
	return verifyNeo4jConnection(ctx, s.driver)
}

// CreateNode creates a new node in Neo4j
func (s *Neo4jService) CreateNode(ctx context.Context, label string, properties map[string]interface{}) (*Node, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := fmt.Sprintf("CREATE (n:%s) SET n = $properties RETURN id(n) as id, labels(n) as labels, properties(n) as properties", label)

	result, err := session.Run(ctx, query, map[string]interface{}{
		"properties": properties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get created node: %w", err)
	}

	return extractNodeFromRecord(record)
}

// GetNode retrieves a node by ID
func (s *Neo4jService) GetNode(ctx context.Context, id string) (*Node, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := "MATCH (n) WHERE id(n) = $id RETURN id(n) as id, labels(n) as labels, properties(n) as properties"

	result, err := session.Run(ctx, query, map[string]interface{}{
		"id": parseNodeID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	return extractNodeFromRecord(record)
}

// UpdateNode updates a node's properties
func (s *Neo4jService) UpdateNode(ctx context.Context, id string, properties map[string]interface{}) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := "MATCH (n) WHERE id(n) = $id SET n += $properties"

	_, err := session.Run(ctx, query, map[string]interface{}{
		"id":         parseNodeID(id),
		"properties": properties,
	})
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	return nil
}

// DeleteNode deletes a node by ID
func (s *Neo4jService) DeleteNode(ctx context.Context, id string) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := "MATCH (n) WHERE id(n) = $id DETACH DELETE n"

	_, err := session.Run(ctx, query, map[string]interface{}{
		"id": parseNodeID(id),
	})
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}

// CreateRelationship creates a relationship between two nodes
func (s *Neo4jService) CreateRelationship(ctx context.Context, fromID, toID, relType string, properties map[string]interface{}) (*Relationship, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := fmt.Sprintf(`
		MATCH (from) WHERE id(from) = $fromID
		MATCH (to) WHERE id(to) = $toID
		CREATE (from)-[r:%s]->(to)
		SET r = $properties
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`, relType)

	result, err := session.Run(ctx, query, map[string]interface{}{
		"fromID":     parseNodeID(fromID),
		"toID":       parseNodeID(toID),
		"properties": properties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get created relationship: %w", err)
	}

	return extractRelationshipFromRecord(record)
}

// GetRelationship retrieves a relationship by ID
func (s *Neo4jService) GetRelationship(ctx context.Context, id string) (*Relationship, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := `
		MATCH (from)-[r]->(to) WHERE id(r) = $id
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"id": parseNodeID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("relationship not found: %w", err)
	}

	return extractRelationshipFromRecord(record)
}

// DeleteRelationship deletes a relationship by ID
func (s *Neo4jService) DeleteRelationship(ctx context.Context, id string) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := "MATCH ()-[r]-() WHERE id(r) = $id DELETE r"

	_, err := session.Run(ctx, query, map[string]interface{}{
		"id": parseNodeID(id),
	})
	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	return nil
}

// ExecuteQuery executes a general query
func (s *Neo4jService) ExecuteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return extractResultRecords(ctx, result)
}

// ExecuteReadQuery executes a read-only query
func (s *Neo4jService) ExecuteReadQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute read query: %w", err)
	}

	return extractResultRecords(ctx, result)
}

// ExecuteWriteQuery executes a write query
func (s *Neo4jService) ExecuteWriteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute write query: %w", err)
	}

	return extractResultRecords(ctx, result)
}

// ExecuteBatch executes multiple operations in a transaction
func (s *Neo4jService) ExecuteBatch(ctx context.Context, operations []BatchOperation) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		for _, op := range operations {
			_, err := tx.Run(ctx, op.Query, op.Parameters)
			if err != nil {
				return nil, fmt.Errorf("failed to execute batch operation: %w", err)
			}
		}
		return nil, nil
	})
	return err
}

// ClearAll removes all nodes and relationships
func (s *Neo4jService) ClearAll(ctx context.Context) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	query := "MATCH (n) DETACH DELETE n"

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("failed to clear all data: %w", err)
	}

	return nil
}

func extractNodeFromRecord(record *neo4j.Record) (*Node, error) {
	id, found := record.Get("id")
	if !found {
		return nil, fmt.Errorf("node ID not found in record")
	}

	labels, found := record.Get("labels")
	if !found {
		return nil, fmt.Errorf("node labels not found in record")
	}

	properties, found := record.Get("properties")
	if !found {
		return nil, fmt.Errorf("node properties not found in record")
	}

	labelStrings := lo.Map(labels.([]interface{}), func(label interface{}, _ int) string {
		return label.(string)
	})

	return &Node{
		ID:         fmt.Sprintf("%d", id.(int64)),
		Labels:     labelStrings,
		Properties: properties.(map[string]interface{}),
	}, nil
}

func extractRelationshipFromRecord(record *neo4j.Record) (*Relationship, error) {
	fields, err := extractRelationshipFields(record)
	if err != nil {
		return nil, err
	}

	return buildRelationshipFromFields(fields)
}

func extractRelationshipFields(record *neo4j.Record) (map[string]interface{}, error) {
	requiredFields := []string{"id", "type", "from_id", "to_id", "properties"}
	fields := make(map[string]interface{})

	for _, field := range requiredFields {
		value, found := record.Get(field)
		if !found {
			return nil, fmt.Errorf("relationship %s not found in record", field)
		}
		fields[field] = value
	}

	return fields, nil
}

func buildRelationshipFromFields(fields map[string]interface{}) (*Relationship, error) {
	return &Relationship{
		ID:         fmt.Sprintf("%d", fields["id"].(int64)),
		Type:       fields["type"].(string),
		FromID:     fmt.Sprintf("%d", fields["from_id"].(int64)),
		ToID:       fmt.Sprintf("%d", fields["to_id"].(int64)),
		Properties: fields["properties"].(map[string]interface{}),
	}, nil
}

func extractResultRecords(ctx context.Context, result neo4j.ResultWithContext) ([]map[string]interface{}, error) {
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect results: %w", err)
	}

	return lo.Map(records, func(record *neo4j.Record, _ int) map[string]interface{} {
		return record.AsMap()
	}), nil
}

func parseNodeID(id string) int64 {
	nodeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}
	return nodeID
}
