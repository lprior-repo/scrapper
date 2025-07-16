package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// GraphConnection represents a connection to the graph database
type GraphConnection struct {
	Driver neo4j.DriverWithContext
	Config Config
}

// executeQuery executes a query against the database (Impure Shell)
func executeQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if conn.Driver == nil {
		return nil, fmt.Errorf("driver not initialized")
	}
	if query == "" {
		panic("Query cannot be empty")
	}

	session := conn.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			// Log error but continue
			_ = err
		}
	}()

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return collectResults(ctx, result)
}

// executeReadQuery executes a read-only query (Impure Shell)
func executeReadQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if conn.Driver == nil {
		return nil, fmt.Errorf("driver not initialized")
	}
	if query == "" {
		panic("Query cannot be empty")
	}

	session := conn.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		if err := session.Close(ctx); err != nil {
			// Log error but continue
			_ = err
		}
	}()

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute read query: %w", err)
	}

	return collectResults(ctx, result)
}

// executeWriteQuery executes a write query (Impure Shell)
func executeWriteQuery(ctx context.Context, conn GraphConnection, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if conn.Driver == nil {
		return nil, fmt.Errorf("driver not initialized")
	}
	if query == "" {
		panic("Query cannot be empty")
	}

	session := conn.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		if err := session.Close(ctx); err != nil {
			// Log error but continue
			_ = err
		}
	}()

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute write query: %w", err)
	}

	return collectResults(ctx, result)
}

// collectResults collects results from a Neo4j result
func collectResults(ctx context.Context, result neo4j.ResultWithContext) ([]map[string]interface{}, error) {
	if result == nil {
		panic("Result cannot be nil")
	}

	records, err := result.Collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect results: %w", err)
	}

	results := make([]map[string]interface{}, len(records))
	for i, record := range records {
		results[i] = record.AsMap()
	}
	
	return results, nil
}

// getSingleResult gets a single result from a Neo4j result
func getSingleResult(ctx context.Context, result neo4j.ResultWithContext) (map[string]interface{}, error) {
	if result == nil {
		panic("Result cannot be nil")
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get single result: %w", err)
	}

	return record.AsMap(), nil
}

// executeBatchOperations executes multiple operations in a transaction (Impure Shell)
func executeBatchOperations(ctx context.Context, conn GraphConnection, operations []BatchOperation) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}
	if len(operations) == 0 {
		return fmt.Errorf("no operations provided")
	}

	session := conn.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			// Log error but continue
			_ = err
		}
	}()

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		for _, op := range operations {
			if op.Query == "" {
				panic("Operation query cannot be empty")
			}
			_, err := tx.Run(ctx, op.Query, op.Parameters)
			if err != nil {
				return nil, fmt.Errorf("failed to execute batch operation: %w", err)
			}
		}
		return nil, nil
	})
	
	return err
}

// clearAllData removes all nodes and relationships (Impure Shell)
func clearAllData(ctx context.Context, conn GraphConnection) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}

	session := conn.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			// Log error but continue
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

// verifyConnection verifies the database connection (Impure Shell)
func verifyConnection(ctx context.Context, conn GraphConnection) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}
	
	return verifyNeo4jConnection(ctx, conn.Driver)
}

// closeConnection closes the database connection (Impure Shell)
func closeConnection(ctx context.Context, conn GraphConnection) error {
	if conn.Driver != nil {
		return conn.Driver.Close(ctx)
	}
	return nil
}

// createConnection creates a new database connection (Impure Shell)
func createConnection(ctx context.Context, config Config) (GraphConnection, error) {
	if err := validateConfig(config); err != nil {
		return GraphConnection{}, fmt.Errorf("invalid configuration: %w", err)
	}

	var driver neo4j.DriverWithContext
	var err error

	switch config.GraphDB.Provider {
	case providerNeo4j:
		neo4jConfig := Neo4jConfig{
			URI:      config.GraphDB.Neo4j.URI,
			Username: config.GraphDB.Neo4j.Username,
			Password: config.GraphDB.Neo4j.Password,
		}
		driver, err = createNeo4jDriver(neo4jConfig)
		if err != nil {
			return GraphConnection{}, fmt.Errorf("failed to create Neo4j driver: %w", err)
		}
	case providerNeptune:
		return GraphConnection{}, fmt.Errorf("Neptune not implemented yet")
	default:
		return GraphConnection{}, fmt.Errorf("unsupported provider: %s", config.GraphDB.Provider)
	}

	conn := GraphConnection{
		Driver: driver,
		Config: config,
	}

	if err := verifyConnection(ctx, conn); err != nil {
		_ = closeConnection(ctx, conn)
		return GraphConnection{}, fmt.Errorf("failed to verify connection: %w", err)
	}

	return conn, nil
}