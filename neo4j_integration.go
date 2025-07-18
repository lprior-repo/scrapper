package main

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/samber/lo"
)

// Neo4jConnection represents a Neo4j database connection
type Neo4jConnection struct {
	driver   neo4j.DriverWithContext
	database string
	timeout  time.Duration
}

// Neo4jSession represents a Neo4j session for transaction management
type Neo4jSession struct {
	session  neo4j.SessionWithContext
	database string
}

// Neo4jTransaction represents a Neo4j transaction
type Neo4jTransaction struct {
	// This struct is kept for future extensibility
}

// Neo4jResult represents the result of a Neo4j query
type Neo4jResult struct {
	Records []map[string]interface{}
	Summary neo4j.ResultSummary
}

// Neo4jError represents Neo4j-specific errors
type Neo4jError struct {
	Code    string
	Message string
	Details string
}

// Error implements the error interface for Neo4jError
func (e Neo4jError) Error() string {
	return fmt.Sprintf("Neo4j error [%s]: %s - %s", e.Code, e.Message, e.Details)
}

// createNeo4jConnection creates a new Neo4j connection (Orchestrator)
func createNeo4jConnection(ctx context.Context, config Neo4jConfig) (*Neo4jConnection, error) {
	validateNeo4jConnectionConfig(config)

	driver, err := neo4j.NewDriverWithContext(
		config.URI,
		neo4j.BasicAuth(config.Username, config.Password, ""),
		func(config *neo4j.Config) { //nolint:staticcheck // Using deprecated type until updated
			config.MaxConnectionLifetime = 30 * time.Minute
			config.MaxConnectionPoolSize = 50
			config.ConnectionAcquisitionTimeout = 2 * time.Minute
		},
	)

	if err != nil {
		return nil, wrapNeo4jError(err, "failed to create Neo4j driver")
	}

	// Verify connectivity
	if err := driver.VerifyConnectivity(ctx); err != nil {
		driver.Close(ctx)
		return nil, wrapNeo4jError(err, "failed to verify Neo4j connectivity")
	}

	return &Neo4jConnection{
		driver:   driver,
		database: config.Database,
		timeout:  config.Timeout,
	}, nil
}

// closeNeo4jConnection closes the Neo4j connection (Orchestrator)
func closeNeo4jConnection(ctx context.Context, conn *Neo4jConnection) error {
	if conn == nil {
		return nil
	}

	if conn.driver != nil {
		return conn.driver.Close(ctx)
	}

	return nil
}

// createNeo4jSession creates a new Neo4j session (Orchestrator)
func createNeo4jSession(ctx context.Context, conn *Neo4jConnection) (*Neo4jSession, error) {
	validateNeo4jConnectionNotNil(conn)

	session := conn.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: conn.database,
	})

	return &Neo4jSession{
		session:  session,
		database: conn.database,
	}, nil
}

// closeNeo4jSession closes the Neo4j session (Orchestrator)
func closeNeo4jSession(ctx context.Context, session *Neo4jSession) error {
	if session == nil {
		return nil
	}

	if session.session != nil {
		return session.session.Close(ctx)
	}

	return nil
}

// executeNeo4jReadQuery executes a read query (Orchestrator)
func executeNeo4jReadQuery(ctx context.Context, session *Neo4jSession, query string, params map[string]interface{}) (Neo4jResult, error) {
	validateNeo4jSessionNotNil(session)
	validateQueryNotEmpty(query)

	result, err := session.session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return executeNeo4jQueryInTx(ctx, tx, query, params)
	})

	if err != nil {
		return Neo4jResult{}, wrapNeo4jError(err, "failed to execute read query")
	}

	return result.(Neo4jResult), nil
}

// executeNeo4jWrite executes a write query (Orchestrator)
func executeNeo4jWrite(ctx context.Context, session *Neo4jSession, query string, params map[string]interface{}) (Neo4jResult, error) {
	validateNeo4jSessionNotNil(session)
	validateQueryNotEmpty(query)

	result, err := session.session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return executeNeo4jQueryInTx(ctx, tx, query, params)
	})

	if err != nil {
		return Neo4jResult{}, wrapNeo4jError(err, "failed to execute write query")
	}

	return result.(Neo4jResult), nil
}

// executeNeo4jQueryInTx executes a single query within a transaction (Pure Core)
func executeNeo4jQueryInTx(ctx context.Context, tx neo4j.ManagedTransaction, query string, params map[string]interface{}) (Neo4jResult, error) {
	validateTransactionNotNil(tx)
	validateQueryNotEmpty(query)

	if params == nil {
		params = make(map[string]interface{})
	}

	result, err := tx.Run(ctx, query, params)
	if err != nil {
		return Neo4jResult{}, wrapNeo4jError(err, "failed to run query")
	}

	records, err := result.Collect(ctx)
	if err != nil {
		return Neo4jResult{}, wrapNeo4jError(err, "failed to collect results")
	}

	mappedRecords := lo.Map(records, func(record *neo4j.Record, _ int) map[string]interface{} {
		return convertNeo4jRecord(record)
	})

	summary, err := result.Consume(ctx)
	if err != nil {
		return Neo4jResult{}, wrapNeo4jError(err, "failed to consume result summary")
	}

	return Neo4jResult{
		Records: mappedRecords,
		Summary: summary,
	}, nil
}

// convertNeo4jRecord converts a Neo4j record to a map (Pure Core)
func convertNeo4jRecord(record *neo4j.Record) map[string]interface{} {
	if record == nil {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})

	for _, key := range record.Keys {
		value, found := record.Get(key)
		if found {
			result[key] = value
		}
	}

	return result
}

// buildNeo4jHealthQuery builds a health check query (Pure Core)
func buildNeo4jHealthQuery() string {
	return "RETURN 1 as health_check"
}

// buildNeo4jConstraintQuery builds a query to create constraints (Pure Core)
func buildNeo4jConstraintQuery(label string, property string) string {
	validateLabelNotEmpty(label)
	validatePropertyNotEmpty(property)

	return fmt.Sprintf("CREATE CONSTRAINT IF NOT EXISTS FOR (n:%s) REQUIRE n.%s IS UNIQUE", label, property)
}

// buildNeo4jIndexQuery builds a query to create an index (Pure Core)
func buildNeo4jIndexQuery(label string, property string) string {
	validateLabelNotEmpty(label)
	validatePropertyNotEmpty(property)

	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS FOR (n:%s) ON (n.%s)", label, property)
}

// checkNeo4jHealth checks database health (Orchestrator)
func checkNeo4jHealth(ctx context.Context, conn *Neo4jConnection) error {
	validateNeo4jConnectionNotNil(conn)

	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		return wrapNeo4jError(err, "failed to create session for health check")
	}
	defer closeNeo4jSession(ctx, session)

	query := buildNeo4jHealthQuery()
	result, err := executeNeo4jReadQuery(ctx, session, query, nil)
	if err != nil {
		return wrapNeo4jError(err, "health check query failed")
	}

	if len(result.Records) == 0 {
		return Neo4jError{
			Code:    "HEALTH_CHECK_FAILED",
			Message: "Health check returned no results",
			Details: "Expected at least one record from health check query",
		}
	}

	return nil
}

// createNeo4jConstraints creates database constraints (Orchestrator)
func createNeo4jConstraints(ctx context.Context, conn *Neo4jConnection) error {
	validateNeo4jConnectionNotNil(conn)

	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		return wrapNeo4jError(err, "failed to create session for constraints")
	}
	defer closeNeo4jSession(ctx, session)

	constraints := []struct {
		label    string
		property string
	}{
		{"Organization", "login"},
		{"Repository", "full_name"},
		{"User", "login"},
		{"Team", "slug"},
	}

	for _, constraint := range constraints {
		query := buildNeo4jConstraintQuery(constraint.label, constraint.property)
		_, err := executeNeo4jWrite(ctx, session, query, nil)
		if err != nil {
			return wrapNeo4jError(err, fmt.Sprintf("failed to create constraint for %s.%s", constraint.label, constraint.property))
		}
	}

	return nil
}

// createNeo4jIndexes creates database indexes (Orchestrator)
func createNeo4jIndexes(ctx context.Context, conn *Neo4jConnection) error {
	validateNeo4jConnectionNotNil(conn)

	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		return wrapNeo4jError(err, "failed to create session for indexes")
	}
	defer closeNeo4jSession(ctx, session)

	indexes := []struct {
		label    string
		property string
	}{
		{"Repository", "name"},
		{"Repository", "updated_at"},
		{"User", "name"},
		{"Team", "name"},
	}

	for _, index := range indexes {
		query := buildNeo4jIndexQuery(index.label, index.property)
		_, err := executeNeo4jWrite(ctx, session, query, nil)
		if err != nil {
			return wrapNeo4jError(err, fmt.Sprintf("failed to create index for %s.%s", index.label, index.property))
		}
	}

	return nil
}

// wrapNeo4jError wraps an error with Neo4j-specific context (Pure Core)
func wrapNeo4jError(err error, message string) Neo4jError {
	if err == nil {
		return Neo4jError{
			Code:    "INTERNAL_ERROR",
			Message: message,
			Details: "nil error wrapped",
		}
	}

	return Neo4jError{
		Code:    "DATABASE_ERROR",
		Message: message,
		Details: err.Error(),
	}
}

// Validation helper functions (Pure Core)
func validateNeo4jConnectionNotNil(conn *Neo4jConnection) {
	if conn == nil {
		panic("Neo4j connection cannot be nil")
	}
}

func validateNeo4jSessionNotNil(session *Neo4jSession) {
	if session == nil {
		panic("Neo4j session cannot be nil")
	}
}

func validateTransactionNotNil(tx neo4j.ManagedTransaction) {
	if tx == nil {
		panic("Neo4j transaction cannot be nil")
	}
}

func validateQueryNotEmpty(query string) {
	if query == "" {
		panic("Query cannot be empty")
	}
}

func validateLabelNotEmpty(label string) {
	if label == "" {
		panic("Label cannot be empty")
	}
}

func validatePropertyNotEmpty(property string) {
	if property == "" {
		panic("Property cannot be empty")
	}
}

func validateNeo4jConnectionConfig(config Neo4jConfig) {
	if config.URI == "" {
		panic("Neo4j URI cannot be empty")
	}
	if config.Username == "" {
		panic("Neo4j username cannot be empty")
	}
	if config.Password == "" {
		panic("Neo4j password cannot be empty")
	}
	if config.Database == "" {
		panic("Neo4j database cannot be empty")
	}
	if config.Timeout <= 0 {
		panic("Neo4j timeout must be positive")
	}
}
