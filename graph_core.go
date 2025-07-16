package main

import (
	"fmt"
	"strconv"

	"github.com/samber/lo"
)

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
	Type       string                 `json:"type"`
	Query      string                 `json:"query"`
	Parameters map[string]interface{} `json:"parameters"`
}

// GraphServiceConfig holds configuration for graph services
type GraphServiceConfig struct {
	Provider string `json:"provider"`
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

// Neo4jConfig holds Neo4j-specific configuration
type Neo4jConfig struct {
	URI      string
	Username string
	Password string
}

// GraphOperationRequest represents a request for graph operations
type GraphOperationRequest struct {
	Operation  string                 `json:"operation"`
	Label      string                 `json:"label,omitempty"`
	NodeID     string                 `json:"node_id,omitempty"`
	FromID     string                 `json:"from_id,omitempty"`
	ToID       string                 `json:"to_id,omitempty"`
	RelType    string                 `json:"rel_type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Query      string                 `json:"query,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// GraphOperationResponse represents the response from graph operations
type GraphOperationResponse struct {
	Success bool                     `json:"success"`
	Error   string                   `json:"error,omitempty"`
	Node    *Node                    `json:"node,omitempty"`
	Rel     *Relationship            `json:"relationship,omitempty"`
	Results []map[string]interface{} `json:"results,omitempty"`
}

// validateNodeCreation validates node creation parameters
func validateNodeCreation(request GraphOperationRequest) error {
	if request.Label == "" {
		return fmt.Errorf("node label is required")
	}
	if request.Properties == nil {
		return fmt.Errorf("node properties cannot be nil")
	}
	return nil
}

// validateNodeRetrieval validates node retrieval parameters  
func validateNodeRetrieval(request GraphOperationRequest) error {
	if request.NodeID == "" {
		return fmt.Errorf("node ID is required")
	}
	if !isValidNodeID(request.NodeID) {
		return fmt.Errorf("invalid node ID format")
	}
	return nil
}

// validateRelationshipCreation validates relationship creation parameters
func validateRelationshipCreation(request GraphOperationRequest) error {
	if request.FromID == "" {
		return fmt.Errorf("from node ID is required")
	}
	if request.ToID == "" {
		return fmt.Errorf("to node ID is required")
	}
	if request.RelType == "" {
		return fmt.Errorf("relationship type is required")
	}
	if !isValidNodeID(request.FromID) {
		return fmt.Errorf("invalid from node ID format")
	}
	if !isValidNodeID(request.ToID) {
		return fmt.Errorf("invalid to node ID format")
	}
	return nil
}

// isValidNodeID checks if a node ID is valid
func isValidNodeID(id string) bool {
	if id == "" {
		return false
	}
	_, err := strconv.ParseInt(id, 10, 64)
	return err == nil
}

// buildCreateNodeQuery builds a Cypher query for node creation
func buildCreateNodeQuery(request GraphOperationRequest) string {
	if request.Label == "" {
		panic("Label cannot be empty")
	}
	return fmt.Sprintf("CREATE (n:%s) SET n = $properties RETURN id(n) as id, labels(n) as labels, properties(n) as properties", request.Label)
}

// buildGetNodeQuery builds a Cypher query for node retrieval
func buildGetNodeQuery() string {
	return "MATCH (n) WHERE id(n) = $id RETURN id(n) as id, labels(n) as labels, properties(n) as properties"
}

// buildUpdateNodeQuery builds a Cypher query for node update
func buildUpdateNodeQuery() string {
	return "MATCH (n) WHERE id(n) = $id SET n += $properties"
}

// buildDeleteNodeQuery builds a Cypher query for node deletion
func buildDeleteNodeQuery() string {
	return "MATCH (n) WHERE id(n) = $id DETACH DELETE n"
}

// buildCreateRelationshipQuery builds a Cypher query for relationship creation
func buildCreateRelationshipQuery(request GraphOperationRequest) string {
	if request.RelType == "" {
		panic("Relationship type cannot be empty")
	}
	return fmt.Sprintf(`
		MATCH (from) WHERE id(from) = $fromID
		MATCH (to) WHERE id(to) = $toID
		CREATE (from)-[r:%s]->(to)
		SET r = $properties
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`, request.RelType)
}

// buildGetRelationshipQuery builds a Cypher query for relationship retrieval
func buildGetRelationshipQuery() string {
	return `
		MATCH (from)-[r]->(to) WHERE id(r) = $id
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`
}

// buildDeleteRelationshipQuery builds a Cypher query for relationship deletion
func buildDeleteRelationshipQuery() string {
	return "MATCH ()-[r]-() WHERE id(r) = $id DELETE r"
}

// convertStringToInt64 converts string to int64 with validation
func convertStringToInt64(s string) int64 {
	if s == "" {
		return 0
	}
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

// prepareNodeParameters prepares parameters for node operations
func prepareNodeParameters(request GraphOperationRequest) map[string]interface{} {
	params := make(map[string]interface{})
	
	if request.NodeID != "" {
		params["id"] = convertStringToInt64(request.NodeID)
	}
	if request.Properties != nil {
		params["properties"] = request.Properties
	}
	
	return params
}

// prepareRelationshipParameters prepares parameters for relationship operations
func prepareRelationshipParameters(request GraphOperationRequest) map[string]interface{} {
	params := make(map[string]interface{})
	
	if request.NodeID != "" {
		params["id"] = convertStringToInt64(request.NodeID)
	}
	if request.FromID != "" {
		params["fromID"] = convertStringToInt64(request.FromID)
	}
	if request.ToID != "" {
		params["toID"] = convertStringToInt64(request.ToID)
	}
	if request.Properties != nil {
		params["properties"] = request.Properties
	}
	
	return params
}

// processNodeResult processes database result into Node struct
func processNodeResult(record map[string]interface{}) (*Node, error) {
	id, found := record["id"]
	if !found {
		return nil, fmt.Errorf("node ID not found in record")
	}

	labels, found := record["labels"]
	if !found {
		return nil, fmt.Errorf("node labels not found in record")
	}

	properties, found := record["properties"]
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

// processRelationshipResult processes database result into Relationship struct
func processRelationshipResult(record map[string]interface{}) (*Relationship, error) {
	requiredFields := []string{"id", "type", "from_id", "to_id", "properties"}
	
	for _, field := range requiredFields {
		if _, found := record[field]; !found {
			return nil, fmt.Errorf("relationship %s not found in record", field)
		}
	}

	return &Relationship{
		ID:         fmt.Sprintf("%d", record["id"].(int64)),
		Type:       record["type"].(string),
		FromID:     fmt.Sprintf("%d", record["from_id"].(int64)),
		ToID:       fmt.Sprintf("%d", record["to_id"].(int64)),
		Properties: record["properties"].(map[string]interface{}),
	}, nil
}