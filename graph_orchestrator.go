package main

import (
	"context"
	"fmt"
)

// createNode creates a new node using Pure Core + Impure Shell pattern
func createNode(ctx context.Context, conn GraphConnection, request GraphOperationRequest) (*Node, error) {
	if err := validateNodeCreation(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := buildCreateNodeQuery(request)
	params := prepareNodeParameters(request)

	results, err := executeWriteQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected number of results: %d", len(results))
	}

	return processNodeResult(results[0])
}

// getNode retrieves a node using Pure Core + Impure Shell pattern
func getNode(ctx context.Context, conn GraphConnection, nodeID string) (*Node, error) {
	request := GraphOperationRequest{
		Operation: "get_node",
		NodeID:    nodeID,
	}

	if err := validateNodeRetrieval(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := buildGetNodeQuery()
	params := prepareNodeParameters(request)

	results, err := executeReadQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("node not found")
	}

	return processNodeResult(results[0])
}

// updateNode updates a node using Pure Core + Impure Shell pattern
func updateNode(ctx context.Context, conn GraphConnection, nodeID string, properties map[string]interface{}) error {
	request := GraphOperationRequest{
		Operation:  "update_node",
		NodeID:     nodeID,
		Properties: properties,
	}

	if err := validateNodeRetrieval(request); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if properties == nil {
		return fmt.Errorf("properties cannot be nil")
	}

	query := buildUpdateNodeQuery()
	params := prepareNodeParameters(request)

	_, err := executeWriteQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	return nil
}

// deleteNode deletes a node using Pure Core + Impure Shell pattern
func deleteNode(ctx context.Context, conn GraphConnection, nodeID string) error {
	request := GraphOperationRequest{
		Operation: "delete_node",
		NodeID:    nodeID,
	}

	if err := validateNodeRetrieval(request); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	query := buildDeleteNodeQuery()
	params := prepareNodeParameters(request)

	_, err := executeWriteQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}

// createRelationship creates a relationship using Pure Core + Impure Shell pattern
func createRelationship(ctx context.Context, conn GraphConnection, request GraphOperationRequest) (*Relationship, error) {
	if err := validateRelationshipCreation(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := buildCreateRelationshipQuery(request)
	params := prepareRelationshipParameters(request)

	results, err := executeWriteQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected number of results: %d", len(results))
	}

	return processRelationshipResult(results[0])
}

// getRelationship retrieves a relationship using Pure Core + Impure Shell pattern
func getRelationship(ctx context.Context, conn GraphConnection, relationshipID string) (*Relationship, error) {
	request := GraphOperationRequest{
		Operation: "get_relationship",
		NodeID:    relationshipID, // Reuse NodeID field for relationship ID
	}

	if err := validateNodeRetrieval(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := buildGetRelationshipQuery()
	params := prepareRelationshipParameters(request)

	results, err := executeReadQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("relationship not found")
	}

	return processRelationshipResult(results[0])
}

// deleteRelationship deletes a relationship using Pure Core + Impure Shell pattern
func deleteRelationship(ctx context.Context, conn GraphConnection, relationshipID string) error {
	request := GraphOperationRequest{
		Operation: "delete_relationship",
		NodeID:    relationshipID, // Reuse NodeID field for relationship ID
	}

	if err := validateNodeRetrieval(request); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	query := buildDeleteRelationshipQuery()
	params := prepareRelationshipParameters(request)

	_, err := executeWriteQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	return nil
}

// executeCustomQuery executes a custom query using Pure Core + Impure Shell pattern
func executeCustomQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if query == "" {
		panic("Query cannot be empty")
	}
	if params == nil {
		params = make(map[string]interface{})
	}

	return executeQuery(ctx, conn, query, params)
}

// executeCustomReadQuery executes a custom read query using Pure Core + Impure Shell pattern
func executeCustomReadQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if query == "" {
		panic("Query cannot be empty")
	}
	if params == nil {
		params = make(map[string]interface{})
	}

	return executeReadQuery(ctx, conn, query, params)
}

// executeCustomWriteQuery executes a custom write query using Pure Core + Impure Shell pattern
func executeCustomWriteQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if query == "" {
		panic("Query cannot be empty")
	}
	if params == nil {
		params = make(map[string]interface{})
	}

	return executeWriteQuery(ctx, conn, query, params)
}

// executeBatch executes batch operations using Pure Core + Impure Shell pattern
func executeBatch(ctx context.Context, conn GraphConnection, operations []BatchOperation) error {
	if len(operations) == 0 {
		return fmt.Errorf("no operations provided")
	}

	for _, op := range operations {
		if op.Query == "" {
			panic("Operation query cannot be empty")
		}
	}

	return executeBatchOperations(ctx, conn, operations)
}

// clearAll clears all data using Pure Core + Impure Shell pattern
func clearAll(ctx context.Context, conn GraphConnection) error {
	return clearAllData(ctx, conn)
}

// healthCheck checks connection health using Pure Core + Impure Shell pattern
func healthCheck(ctx context.Context, conn GraphConnection) error {
	return verifyConnection(ctx, conn)
}