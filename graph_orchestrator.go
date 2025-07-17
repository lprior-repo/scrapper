package main

import (
	"context"
	"fmt"
)

// createNode creates a new node in the graph database
func createNode(ctx context.Context, conn GraphConnection, request GraphOperationRequest) (*Node, error) {
	if err := validateNodeCreation(request); err != nil {
		return nil, err
	}

	query := buildCreateNodeQuery(request)
	params := prepareNodeParameters(request)

	results, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results returned from create node query")
	}

	return processNodeResult(results[0])
}

// getNode retrieves a node from the graph database
func getNode(ctx context.Context, conn GraphConnection, nodeID string) (*Node, error) {
	request := GraphOperationRequest{NodeID: nodeID}
	if err := validateNodeRetrieval(request); err != nil {
		return nil, err
	}

	query := buildGetNodeQuery()
	params := prepareNodeParameters(request)

	results, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("node not found")
	}

	return processNodeResult(results[0])
}

// updateNode updates a node in the graph database
func updateNode(ctx context.Context, conn GraphConnection, nodeID string, properties map[string]interface{}) error {
	request := GraphOperationRequest{
		NodeID:     nodeID,
		Properties: properties,
	}

	if err := validateNodeRetrieval(request); err != nil {
		return err
	}

	query := buildUpdateNodeQuery()
	params := prepareNodeParameters(request)

	_, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	return nil
}

// deleteNode deletes a node from the graph database
func deleteNode(ctx context.Context, conn GraphConnection, nodeID string) error {
	request := GraphOperationRequest{NodeID: nodeID}
	if err := validateNodeRetrieval(request); err != nil {
		return err
	}

	query := buildDeleteNodeQuery()
	params := prepareNodeParameters(request)

	_, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}

// createRelationship creates a relationship between two nodes
func createRelationship(ctx context.Context, conn GraphConnection, request GraphOperationRequest) (*Relationship, error) {
	if err := validateRelationshipCreation(request); err != nil {
		return nil, err
	}

	if err := validateRelationshipBusinessRules(request); err != nil {
		return nil, err
	}

	query := buildCreateRelationshipQuery(request)
	params := prepareRelationshipParameters(request)

	results, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results returned from create relationship query")
	}

	return processRelationshipResult(results[0])
}

// getRelationship retrieves a relationship from the graph database
func getRelationship(ctx context.Context, conn GraphConnection, relationshipID string) (*Relationship, error) {
	request := GraphOperationRequest{NodeID: relationshipID}
	if err := validateRelationshipUpdate(request); err != nil {
		return nil, err
	}

	query := buildGetRelationshipQuery()
	params := prepareRelationshipParameters(request)

	results, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("relationship not found")
	}

	return processRelationshipResult(results[0])
}

// RelationshipUpdateRequest holds data for updating a relationship
type RelationshipUpdateRequest struct {
	ID         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties"`
}

// updateRelationship updates a relationship in the graph database
func updateRelationship(ctx context.Context, conn GraphConnection, request RelationshipUpdateRequest) (*Relationship, error) {
	operationRequest := GraphOperationRequest{
		NodeID:     request.ID,
		Properties: request.Properties,
	}

	if err := validateRelationshipUpdate(operationRequest); err != nil {
		return nil, err
	}

	query := buildUpdateRelationshipQuery()
	params := prepareRelationshipParameters(operationRequest)

	results, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update relationship: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("relationship not found")
	}

	return processRelationshipResult(results[0])
}

// deleteRelationship deletes a relationship from the graph database
func deleteRelationship(ctx context.Context, conn GraphConnection, relationshipID string) error {
	request := GraphOperationRequest{NodeID: relationshipID}
	if err := validateRelationshipUpdate(request); err != nil {
		return err
	}

	query := buildDeleteRelationshipQuery()
	params := prepareRelationshipParameters(request)

	_, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	return nil
}

// executeCustomQuery executes a custom Cypher query
func executeCustomQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	return executeQuery(ctx, conn, query, params)
}

// executeCustomWriteQuery executes a custom write Cypher query
func executeCustomWriteQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	return executeWriteQuery(ctx, conn, query, params)
}

// clearAll clears all data from the graph database
func clearAll(ctx context.Context, conn GraphConnection) error {
	query := "MATCH (n) DETACH DELETE n"
	_, err := executeQuery(ctx, conn, query, nil)
	return err
}
