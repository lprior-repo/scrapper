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
	ID         string                 `json:"id,omitempty"`
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

// GraphPath represents a path through the graph
type GraphPath struct {
	Nodes         []*Node         `json:"nodes"`
	Relationships []*Relationship `json:"relationships"`
	Length        int             `json:"length"`
}

// convertStringToInt64 converts string to int64 with validation (Pure Core)
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

// prepareNodeParameters prepares parameters for node operations (Pure Core)
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

// prepareRelationshipParameters prepares parameters for relationship operations (Pure Core)
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

	// Always include properties parameter, even if empty
	if request.Properties != nil {
		params["properties"] = request.Properties
	} else {
		params["properties"] = make(map[string]interface{})
	}

	return params
}

// processNodeResult processes database result into Node struct (Pure Core)
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

// processRelationshipResult processes database result into Relationship struct (Pure Core)
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

// parseNodeID parses a string to an int64, returning 0 for invalid strings (Pure Core)
func parseNodeID(s string) int64 {
	if s == "" {
		return 0
	}

	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// processPathResult processes database result into GraphPath struct (Pure Core)
func processPathResult(record map[string]interface{}) (*GraphPath, error) {
	nodesData, found := record["nodes"]
	if !found {
		return nil, fmt.Errorf("nodes not found in record")
	}

	relationshipsData, found := record["relationships"]
	if !found {
		return nil, fmt.Errorf("relationships not found in record")
	}

	lengthData, found := record["length"]
	if !found {
		return nil, fmt.Errorf("length not found in record")
	}

	// Process nodes
	nodesList := nodesData.([]interface{})
	nodes := make([]*Node, len(nodesList))
	for i, nodeData := range nodesList {
		nodeMap := nodeData.(map[string]interface{})
		node := &Node{
			ID:         fmt.Sprintf("%d", nodeMap["id"].(int64)),
			Labels:     convertToStringSlice(nodeMap["labels"].([]interface{})),
			Properties: nodeMap["properties"].(map[string]interface{}),
		}
		nodes[i] = node
	}

	// Process relationships
	relsList := relationshipsData.([]interface{})
	relationships := make([]*Relationship, len(relsList))
	for i, relData := range relsList {
		relMap := relData.(map[string]interface{})
		rel := &Relationship{
			ID:         fmt.Sprintf("%d", relMap["id"].(int64)),
			Type:       relMap["type"].(string),
			FromID:     fmt.Sprintf("%d", relMap["from_id"].(int64)),
			ToID:       fmt.Sprintf("%d", relMap["to_id"].(int64)),
			Properties: relMap["properties"].(map[string]interface{}),
		}
		relationships[i] = rel
	}

	return &GraphPath{
		Nodes:         nodes,
		Relationships: relationships,
		Length:        int(lengthData.(int64)),
	}, nil
}

// convertToStringSlice converts interface{} slice to string slice (Pure Core)
func convertToStringSlice(input []interface{}) []string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = v.(string)
	}
	return result
}