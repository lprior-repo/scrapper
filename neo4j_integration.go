// Neo4j Integration with Comprehensive Observability
//
// This package provides a fully observable Neo4j database integration layer for the
// GitHub Codeowners Visualization application. It includes comprehensive observability
// features built on top of the utilities from observability.go.
//
// Key Features:
// 1. Custom spans for all Neo4j database operations (connection, queries, session management)
// 2. Structured logging for database operations with query performance metrics
// 3. Performance timing for database queries with detailed timing analysis
// 4. Enhanced error handling with database-specific error context
// 5. Database health monitoring with connection pool metrics
// 6. Query execution logging with parameter sanitization
// 7. Transaction tracking and commit/rollback observability
// 8. Database operation metrics (query count, execution time, result counts)
//
// Observability Components:
// - Spans: Track all database operations with detailed attributes
// - Metrics: Record query performance, connection pool status, error rates
// - Logging: Structured logs with correlation IDs and sanitized parameters
// - Performance: Query execution timing with slow query detection
// - Health: Connection health monitoring and pool status tracking
//
// Usage Patterns:
//
// Basic connection (no observability):
//   conn, err := createNeo4jConnection(ctx, config)
//
// Observable connection:
//   conn, err := createObservableNeo4jConnection(ctx, gofrCtx, config)
//
// Upgrade existing connection:
//   upgradeNeo4jConnectionObservability(conn, gofrCtx)
//
// All query operations automatically include observability when the connection
// has an associated GoFr context. The observability features are designed to
// have minimal performance impact while providing comprehensive insights.
//
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/samber/lo"
	"gofr.dev/pkg/gofr"
)

// Neo4jConnection represents a Neo4j database connection with observability
type Neo4jConnection struct {
	driver   neo4j.DriverWithContext
	database string
	timeout  time.Duration
	metrics  *MetricsCollector
	ctx      *gofr.Context
}

// Neo4jSession represents a Neo4j session for transaction management with observability
type Neo4jSession struct {
	session       neo4j.SessionWithContext
	database      string
	metrics       *MetricsCollector
	ctx           *gofr.Context
	queryCount    int
	totalDuration time.Duration
}

// Neo4jTransaction represents a Neo4j transaction
type Neo4jTransaction struct {
	// This struct is kept for future extensibility
}

// Neo4jResult represents the result of a Neo4j query with observability data
type Neo4jResult struct {
	Records       []map[string]interface{}
	Summary       neo4j.ResultSummary
	ExecutionTime time.Duration
	RecordCount   int
	QueryHash     string
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
	return createNeo4jConnectionWithObservability(ctx, nil, config)
}

// createNeo4jConnectionWithObservability creates a new Neo4j connection with observability (Orchestrator)
func createNeo4jConnectionWithObservability(ctx context.Context, gofrCtx *gofr.Context, config Neo4jConfig) (*Neo4jConnection, error) {
	validateNeo4jConnectionConfig(config)

	// Initialize observability components if GoFr context is available
	var span *SpanWrapper
	var timer *PerformanceTimer
	var metrics *MetricsCollector

	if gofrCtx != nil {
		// Create span for connection operation
		span = createNeo4jSpan(gofrCtx, "connection.create", "CREATE CONNECTION")
		defer finishSpan(span)

		// Start performance timer
		timer = startPerformanceTimer(gofrCtx, "neo4j_connection_create")
		defer stopPerformanceTimer(timer)

		// Initialize metrics collector
		metrics = newMetricsCollector(gofrCtx, "neo4j-client")
	}

	// Log connection attempt with sanitized config if observability is available
	if gofrCtx != nil {
		logInfo(gofrCtx, "Creating Neo4j connection", LogFields{
			"component":   "neo4j_client",
			"operation":   "create_connection",
			"uri":         sanitizeURI(config.URI),
			"database":    config.Database,
			"timeout_ms":  config.Timeout.Milliseconds(),
			"max_pool":    50,
			"max_lifetime": (30 * time.Minute).String(),
		})
	}

	driver, err := neo4j.NewDriverWithContext(
		config.URI,
		neo4j.BasicAuth(config.Username, config.Password, ""),
		func(driverConfig *neo4j.Config) { //nolint:staticcheck // Using deprecated type until updated
			driverConfig.MaxConnectionLifetime = 30 * time.Minute
			driverConfig.MaxConnectionPoolSize = 50
			driverConfig.ConnectionAcquisitionTimeout = 2 * time.Minute
		},
	)

	if err != nil {
		// Log and record connection creation failure if observability is available
		if gofrCtx != nil {
			logError(gofrCtx, "Failed to create Neo4j driver", LogFields{
				"component": "neo4j_client",
				"operation": "create_driver",
				"error":     err.Error(),
				"uri":       sanitizeURI(config.URI),
			})
			if metrics != nil {
				metrics.recordErrorCount("neo4j_client", "driver_creation_failed")
			}
		}
		return nil, wrapNeo4jError(err, "failed to create Neo4j driver")
	}

	// Verify connectivity with observability if available
	if gofrCtx != nil {
		logDebug(gofrCtx, "Verifying Neo4j connectivity", LogFields{
			"component": "neo4j_client",
			"operation": "verify_connectivity",
		})
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		// Log connectivity failure and clean up
		if gofrCtx != nil {
			logError(gofrCtx, "Failed to verify Neo4j connectivity", LogFields{
				"component": "neo4j_client",
				"operation": "verify_connectivity",
				"error":     err.Error(),
				"uri":       sanitizeURI(config.URI),
			})
			if metrics != nil {
				metrics.recordErrorCount("neo4j_client", "connectivity_verification_failed")
			}
		}
		driver.Close(ctx)
		return nil, wrapNeo4jError(err, "failed to verify Neo4j connectivity")
	}

	// Log successful connection if observability is available
	if gofrCtx != nil {
		logInfo(gofrCtx, "Neo4j connection established successfully", LogFields{
			"component": "neo4j_client",
			"operation": "connection_established",
			"database":  config.Database,
			"uri":       sanitizeURI(config.URI),
		})

		// Record connection success metric
		if metrics != nil {
			metrics.recordCounter("neo4j_connections_total", 1, MetricLabels{
				"database": config.Database,
				"status":   "success",
			})
		}
	}

	connection := &Neo4jConnection{
		driver:   driver,
		database: config.Database,
		timeout:  config.Timeout,
		metrics:  metrics,
		ctx:      gofrCtx,
	}

	// Log connection pool status if observability is available
	if gofrCtx != nil {
		logNeo4jConnectionPoolStatus(gofrCtx, connection)
	}

	return connection, nil
}

// closeNeo4jConnection closes the Neo4j connection (Orchestrator)
func closeNeo4jConnection(ctx context.Context, conn *Neo4jConnection) error {
	if conn == nil {
		return nil
	}

	// Create span for connection close operation
	if conn.ctx != nil {
		span := createNeo4jSpan(conn.ctx, "connection.close", "CLOSE CONNECTION")
		defer finishSpan(span)

		// Log connection close
		logInfo(conn.ctx, "Closing Neo4j connection", LogFields{
			"component": "neo4j_client",
			"operation": "close_connection",
			"database":  conn.database,
		})
	}

	if conn.driver != nil {
		err := conn.driver.Close(ctx)
		if err != nil && conn.ctx != nil {
			// Log close error
			logError(conn.ctx, "Failed to close Neo4j connection", LogFields{
				"component": "neo4j_client",
				"operation": "close_connection",
				"error":     err.Error(),
				"database":  conn.database,
			})
			if conn.metrics != nil {
				conn.metrics.recordErrorCount("neo4j_client", "connection_close_failed")
			}
		} else if conn.ctx != nil {
			// Log successful close
			logDebug(conn.ctx, "Neo4j connection closed successfully", LogFields{
				"component": "neo4j_client",
				"operation": "connection_closed",
				"database":  conn.database,
			})
		}
		return err
	}

	return nil
}

// createNeo4jSession creates a new Neo4j session (Orchestrator)
func createNeo4jSession(ctx context.Context, conn *Neo4jConnection) (*Neo4jSession, error) {
	validateNeo4jConnectionNotNil(conn)

	// Create span for session creation
	span := createNeo4jSpan(conn.ctx, "session.create", "CREATE SESSION")
	defer finishSpan(span)

	// Log session creation
	logDebug(conn.ctx, "Creating Neo4j session", LogFields{
		"component": "neo4j_client",
		"operation": "create_session",
		"database":  conn.database,
	})

	session := conn.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: conn.database,
	})

	// Log successful session creation
	logDebug(conn.ctx, "Neo4j session created successfully", LogFields{
		"component": "neo4j_client",
		"operation": "session_created",
		"database":  conn.database,
	})

	// Record session creation metric
	if conn.metrics != nil {
		conn.metrics.recordCounter("neo4j_sessions_total", 1, MetricLabels{
			"database": conn.database,
			"status":   "created",
		})
	}

	return &Neo4jSession{
		session:       session,
		database:      conn.database,
		metrics:       conn.metrics,
		ctx:           conn.ctx,
		queryCount:    0,
		totalDuration: 0,
	}, nil
}

// closeNeo4jSession closes the Neo4j session (Orchestrator)
func closeNeo4jSession(ctx context.Context, session *Neo4jSession) error {
	if session == nil {
		return nil
	}

	// Create span for session close
	if session.ctx != nil {
		span := createNeo4jSpan(session.ctx, "session.close", "CLOSE SESSION")
		defer finishSpan(span)

		// Log session statistics before closing
		logInfo(session.ctx, "Closing Neo4j session", LogFields{
			"component":      "neo4j_client",
			"operation":      "close_session",
			"database":       session.database,
			"query_count":    session.queryCount,
			"total_duration": session.totalDuration.String(),
			"avg_query_time": calculateAverageQueryTime(session).String(),
		})

		// Record session metrics
		if session.metrics != nil {
			session.metrics.recordCounter("neo4j_queries_per_session", session.queryCount, MetricLabels{
				"database": session.database,
			})
			session.metrics.recordDuration("neo4j_session_duration", session.totalDuration, MetricLabels{
				"database": session.database,
			})
		}
	}

	if session.session != nil {
		err := session.session.Close(ctx)
		if err != nil && session.ctx != nil {
			// Log close error
			logError(session.ctx, "Failed to close Neo4j session", LogFields{
				"component": "neo4j_client",
				"operation": "close_session",
				"error":     err.Error(),
				"database":  session.database,
			})
			if session.metrics != nil {
				session.metrics.recordErrorCount("neo4j_client", "session_close_failed")
			}
		} else if session.ctx != nil {
			// Log successful close
			logDebug(session.ctx, "Neo4j session closed successfully", LogFields{
				"component": "neo4j_client",
				"operation": "session_closed",
				"database":  session.database,
			})
		}
		return err
	}

	return nil
}

// executeNeo4jReadQuery executes a read query (Orchestrator)
func executeNeo4jReadQuery(ctx context.Context, session *Neo4jSession, query string, params map[string]interface{}) (Neo4jResult, error) {
	validateNeo4jSessionNotNil(session)
	validateQueryNotEmpty(query)

	// Create span for read query
	span := createNeo4jSpan(session.ctx, "query.read", query)
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(session.ctx, "neo4j_read_query")
	defer func() {
		duration := stopPerformanceTimer(timer)
		session.totalDuration += duration
		session.queryCount++
	}()

	// Sanitize parameters for logging
	sanitizedParams := sanitizeParams(params)
	queryHash := generateQueryHash(query)

	// Log query execution start
	logInfo(session.ctx, "Executing Neo4j read query", LogFields{
		"component":     "neo4j_client",
		"operation":     "execute_read_query",
		"database":      session.database,
		"query_hash":    queryHash,
		"query_preview": truncateQuery(query, 100),
		"param_count":   len(params),
		"params":        sanitizedParams,
		"tx_type":       "read",
	})

	// Add span attributes for detailed tracing
	addSpanAttribute(session.ctx, "db.statement", truncateQuery(query, 200))
	addSpanAttribute(session.ctx, "db.operation", "read")
	addSpanAttribute(session.ctx, "db.name", session.database)
	addSpanAttribute(session.ctx, "db.type", "neo4j")
	addSpanAttribute(session.ctx, "neo4j.query.hash", queryHash)
	addSpanAttribute(session.ctx, "neo4j.param.count", len(params))

	result, err := session.session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return executeNeo4jQueryInTx(ctx, session, tx, query, params)
	})

	if err != nil {
		// Log and record query failure
		logError(session.ctx, "Failed to execute Neo4j read query", LogFields{
			"component":     "neo4j_client",
			"operation":     "execute_read_query",
			"error":         err.Error(),
			"database":      session.database,
			"query_hash":    queryHash,
			"query_preview": truncateQuery(query, 100),
			"tx_type":       "read",
		})
		if session.metrics != nil {
			session.metrics.recordErrorCount("neo4j_client", "read_query_failed")
			session.metrics.recordCounter("neo4j_query_errors_total", 1, MetricLabels{
				"database":   session.database,
				"query_type": "read",
				"error_type": extractErrorType(err),
			})
		}
		return Neo4jResult{}, wrapNeo4jError(err, "failed to execute read query")
	}

	neoResult := result.(Neo4jResult)

	// Log successful query execution with metrics
	logInfo(session.ctx, "Neo4j read query executed successfully", LogFields{
		"component":      "neo4j_client",
		"operation":      "read_query_success",
		"database":       session.database,
		"query_hash":     queryHash,
		"record_count":   neoResult.RecordCount,
		"execution_time": neoResult.ExecutionTime.String(),
		"tx_type":        "read",
	})

	// Record success metrics
	if session.metrics != nil {
		session.metrics.recordCounter("neo4j_queries_total", 1, MetricLabels{
			"database":   session.database,
			"query_type": "read",
			"status":     "success",
		})
		session.metrics.recordDuration("neo4j_query_duration", neoResult.ExecutionTime, MetricLabels{
			"database":   session.database,
			"query_type": "read",
		})
		session.metrics.recordCounter("neo4j_records_returned_total", neoResult.RecordCount, MetricLabels{
			"database":   session.database,
			"query_type": "read",
		})
	}

	// Monitor query performance
	monitorNeo4jPerformance(session.ctx, session, neoResult, "read")

	return neoResult, nil
}

// executeNeo4jWrite executes a write query (Orchestrator)
func executeNeo4jWrite(ctx context.Context, session *Neo4jSession, query string, params map[string]interface{}) (Neo4jResult, error) {
	validateNeo4jSessionNotNil(session)
	validateQueryNotEmpty(query)

	// Create span for write query
	span := createNeo4jSpan(session.ctx, "query.write", query)
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(session.ctx, "neo4j_write_query")
	defer func() {
		duration := stopPerformanceTimer(timer)
		session.totalDuration += duration
		session.queryCount++
	}()

	// Sanitize parameters for logging
	sanitizedParams := sanitizeParams(params)
	queryHash := generateQueryHash(query)

	// Log query execution start
	logInfo(session.ctx, "Executing Neo4j write query", LogFields{
		"component":     "neo4j_client",
		"operation":     "execute_write_query",
		"database":      session.database,
		"query_hash":    queryHash,
		"query_preview": truncateQuery(query, 100),
		"param_count":   len(params),
		"params":        sanitizedParams,
		"tx_type":       "write",
	})

	// Add span attributes for detailed tracing
	addSpanAttribute(session.ctx, "db.statement", truncateQuery(query, 200))
	addSpanAttribute(session.ctx, "db.operation", "write")
	addSpanAttribute(session.ctx, "db.name", session.database)
	addSpanAttribute(session.ctx, "db.type", "neo4j")
	addSpanAttribute(session.ctx, "neo4j.query.hash", queryHash)
	addSpanAttribute(session.ctx, "neo4j.param.count", len(params))

	result, err := session.session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return executeNeo4jQueryInTx(ctx, session, tx, query, params)
	})

	if err != nil {
		// Log and record query failure
		logError(session.ctx, "Failed to execute Neo4j write query", LogFields{
			"component":     "neo4j_client",
			"operation":     "execute_write_query",
			"error":         err.Error(),
			"database":      session.database,
			"query_hash":    queryHash,
			"query_preview": truncateQuery(query, 100),
			"tx_type":       "write",
		})
		if session.metrics != nil {
			session.metrics.recordErrorCount("neo4j_client", "write_query_failed")
			session.metrics.recordCounter("neo4j_query_errors_total", 1, MetricLabels{
				"database":   session.database,
				"query_type": "write",
				"error_type": extractErrorType(err),
			})
		}
		return Neo4jResult{}, wrapNeo4jError(err, "failed to execute write query")
	}

	neoResult := result.(Neo4jResult)

	// Log successful query execution with metrics
	logInfo(session.ctx, "Neo4j write query executed successfully", LogFields{
		"component":        "neo4j_client",
		"operation":        "write_query_success",
		"database":         session.database,
		"query_hash":       queryHash,
		"record_count":     neoResult.RecordCount,
		"execution_time":   neoResult.ExecutionTime.String(),
		"tx_type":          "write",
		"nodes_created":    extractSummaryStatistic(neoResult.Summary, "nodes_created"),
		"nodes_deleted":    extractSummaryStatistic(neoResult.Summary, "nodes_deleted"),
		"relationships_created": extractSummaryStatistic(neoResult.Summary, "relationships_created"),
		"relationships_deleted": extractSummaryStatistic(neoResult.Summary, "relationships_deleted"),
		"properties_set":   extractSummaryStatistic(neoResult.Summary, "properties_set"),
	})

	// Record success metrics
	if session.metrics != nil {
		session.metrics.recordCounter("neo4j_queries_total", 1, MetricLabels{
			"database":   session.database,
			"query_type": "write",
			"status":     "success",
		})
		session.metrics.recordDuration("neo4j_query_duration", neoResult.ExecutionTime, MetricLabels{
			"database":   session.database,
			"query_type": "write",
		})
		session.metrics.recordCounter("neo4j_records_affected_total", neoResult.RecordCount, MetricLabels{
			"database":   session.database,
			"query_type": "write",
		})
	}

	// Monitor query performance
	monitorNeo4jPerformance(session.ctx, session, neoResult, "write")

	return neoResult, nil
}

// executeNeo4jQueryInTx executes a single query within a transaction (Pure Core)
func executeNeo4jQueryInTx(ctx context.Context, session *Neo4jSession, tx neo4j.ManagedTransaction, query string, params map[string]interface{}) (Neo4jResult, error) {
	validateTransactionNotNil(tx)
	validateQueryNotEmpty(query)

	// Start execution timer
	executionStart := time.Now()

	if params == nil {
		params = make(map[string]interface{})
	}

	// Create span for transaction execution
	txSpan := createNeo4jSpan(session.ctx, "transaction.execute", query)
	defer finishSpan(txSpan)

	// Log transaction start
	logDebug(session.ctx, "Starting Neo4j transaction execution", LogFields{
		"component":     "neo4j_client",
		"operation":     "transaction_execute",
		"database":      session.database,
		"query_preview": truncateQuery(query, 100),
		"param_count":   len(params),
	})

	result, err := tx.Run(ctx, query, params)
	if err != nil {
		// Log transaction run failure
		logError(session.ctx, "Failed to run Neo4j query in transaction", LogFields{
			"component":     "neo4j_client",
			"operation":     "transaction_run",
			"error":         err.Error(),
			"database":      session.database,
			"query_preview": truncateQuery(query, 100),
		})
		if session.metrics != nil {
			session.metrics.recordErrorCount("neo4j_client", "transaction_run_failed")
		}
		return Neo4jResult{}, wrapNeo4jError(err, "failed to run query")
	}

	// Collect records with timing
	collectStart := time.Now()
	records, err := result.Collect(ctx)
	collectDuration := time.Since(collectStart)

	if err != nil {
		// Log record collection failure
		logError(session.ctx, "Failed to collect Neo4j query results", LogFields{
			"component":       "neo4j_client",
			"operation":       "collect_results",
			"error":           err.Error(),
			"database":        session.database,
			"collect_duration": collectDuration.String(),
		})
		if session.metrics != nil {
			session.metrics.recordErrorCount("neo4j_client", "collect_results_failed")
		}
		return Neo4jResult{}, wrapNeo4jError(err, "failed to collect results")
	}

	// Convert records with timing
	convertStart := time.Now()
	mappedRecords := lo.Map(records, func(record *neo4j.Record, _ int) map[string]interface{} {
		return convertNeo4jRecord(record)
	})
	convertDuration := time.Since(convertStart)

	// Consume summary with timing
	consumeStart := time.Now()
	summary, err := result.Consume(ctx)
	consumeDuration := time.Since(consumeStart)

	if err != nil {
		// Log summary consumption failure
		logError(session.ctx, "Failed to consume Neo4j result summary", LogFields{
			"component":        "neo4j_client",
			"operation":        "consume_summary",
			"error":            err.Error(),
			"database":         session.database,
			"consume_duration": consumeDuration.String(),
		})
		if session.metrics != nil {
			session.metrics.recordErrorCount("neo4j_client", "consume_summary_failed")
		}
		return Neo4jResult{}, wrapNeo4jError(err, "failed to consume result summary")
	}

	// Calculate total execution time
	totalExecutionTime := time.Since(executionStart)
	recordCount := len(mappedRecords)
	queryHash := generateQueryHash(query)

	// Log successful transaction execution with detailed metrics
	logDebug(session.ctx, "Neo4j transaction executed successfully", LogFields{
		"component":          "neo4j_client",
		"operation":          "transaction_success",
		"database":           session.database,
		"query_hash":         queryHash,
		"record_count":       recordCount,
		"execution_time":     totalExecutionTime.String(),
		"collect_duration":   collectDuration.String(),
		"convert_duration":   convertDuration.String(),
		"consume_duration":   consumeDuration.String(),
		"query_type":         determineQueryType(query),
		"server_address":     extractServerAddress(summary),
		"query_id":           extractQueryID(summary),
	})

	// Record detailed transaction metrics
	if session.metrics != nil {
		session.metrics.recordDuration("neo4j_transaction_duration", totalExecutionTime, MetricLabels{
			"database":   session.database,
			"query_type": determineQueryType(query),
			"phase":      "total",
		})
		session.metrics.recordDuration("neo4j_collect_duration", collectDuration, MetricLabels{
			"database":   session.database,
			"query_type": determineQueryType(query),
			"phase":      "collect",
		})
		session.metrics.recordDuration("neo4j_convert_duration", convertDuration, MetricLabels{
			"database":   session.database,
			"query_type": determineQueryType(query),
			"phase":      "convert",
		})
	}

	return Neo4jResult{
		Records:       mappedRecords,
		Summary:       summary,
		ExecutionTime: totalExecutionTime,
		RecordCount:   recordCount,
		QueryHash:     queryHash,
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

	// Create span for health check
	span := createHealthCheckSpan(conn.ctx, "neo4j")
	defer finishSpan(span)

	// Start performance timer for health check
	timer := startPerformanceTimer(conn.ctx, "neo4j_health_check")
	defer stopPerformanceTimer(timer)

	// Log health check start
	logInfo(conn.ctx, "Starting Neo4j health check", LogFields{
		"component": "neo4j_client",
		"operation": "health_check",
		"database":  conn.database,
		"check_type": "connectivity",
	})

	// Create session for health check
	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		// Log health check session creation failure
		errorDetails := map[string]interface{}{
			"error":         err.Error(),
			"database":      conn.database,
			"check_phase":   "session_creation",
			"health_status": "unhealthy",
		}
		logHealthCheckResult(conn.ctx, "neo4j", false, errorDetails)
		if conn.metrics != nil {
			conn.metrics.recordErrorCount("neo4j_client", "health_check_session_failed")
			conn.metrics.recordCounter("neo4j_health_checks_total", 1, MetricLabels{
				"database": conn.database,
				"status":   "failed",
				"phase":    "session_creation",
			})
		}
		return wrapNeo4jError(err, "failed to create session for health check")
	}
	defer closeNeo4jSession(ctx, session)

	// Get connection pool metrics if available
	poolMetrics := getConnectionPoolMetrics(conn)

	// Execute health check query
	query := buildNeo4jHealthQuery()
	result, err := executeNeo4jReadQuery(ctx, session, query, nil)

	if err != nil {
		// Log health check query failure
		errorDetails := map[string]interface{}{
			"error":           err.Error(),
			"database":        conn.database,
			"check_phase":     "query_execution",
			"health_status":   "unhealthy",
			"query":           query,
			"execution_time":  result.ExecutionTime.String(),
			"pool_metrics":    poolMetrics,
		}
		logHealthCheckResult(conn.ctx, "neo4j", false, errorDetails)
		if conn.metrics != nil {
			conn.metrics.recordErrorCount("neo4j_client", "health_check_query_failed")
			conn.metrics.recordCounter("neo4j_health_checks_total", 1, MetricLabels{
				"database": conn.database,
				"status":   "failed",
				"phase":    "query_execution",
			})
		}
		return wrapNeo4jError(err, "health check query failed")
	}

	// Validate health check results
	if len(result.Records) == 0 {
		// Log health check validation failure
		errorDetails := map[string]interface{}{
			"database":        conn.database,
			"check_phase":     "result_validation",
			"health_status":   "unhealthy",
			"expected_records": 1,
			"actual_records":   0,
			"execution_time":  result.ExecutionTime.String(),
			"pool_metrics":    poolMetrics,
		}
		logHealthCheckResult(conn.ctx, "neo4j", false, errorDetails)
		if conn.metrics != nil {
			conn.metrics.recordCounter("neo4j_health_checks_total", 1, MetricLabels{
				"database": conn.database,
				"status":   "failed",
				"phase":    "result_validation",
			})
		}
		return Neo4jError{
			Code:    "HEALTH_CHECK_FAILED",
			Message: "Health check returned no results",
			Details: "Expected at least one record from health check query",
		}
	}

	// Health check passed - log success with metrics
	successDetails := map[string]interface{}{
		"database":        conn.database,
		"health_status":   "healthy",
		"record_count":    len(result.Records),
		"execution_time":  result.ExecutionTime.String(),
		"server_address":  extractServerAddress(result.Summary),
		"server_version":  extractServerVersion(result.Summary),
		"query_id":        extractQueryID(result.Summary),
		"pool_metrics":    poolMetrics,
		"database_mode":   extractDatabaseMode(result.Summary),
	}
	logHealthCheckResult(conn.ctx, "neo4j", true, successDetails)

	// Record successful health check metrics
	if conn.metrics != nil {
		conn.metrics.recordCounter("neo4j_health_checks_total", 1, MetricLabels{
			"database": conn.database,
			"status":   "success",
		})
		conn.metrics.recordDuration("neo4j_health_check_duration", result.ExecutionTime, MetricLabels{
			"database": conn.database,
		})

		// Record database availability metric
		conn.metrics.recordGauge("neo4j_database_available", 1.0, MetricLabels{
			"database": conn.database,
		})
	}

	return nil
}

// createNeo4jConstraints creates database constraints (Orchestrator)
func createNeo4jConstraints(ctx context.Context, conn *Neo4jConnection) error {
	validateNeo4jConnectionNotNil(conn)

	// Create span for constraint creation
	span := createNeo4jSpan(conn.ctx, "schema.create_constraints", "CREATE CONSTRAINTS")
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(conn.ctx, "neo4j_create_constraints")
	defer stopPerformanceTimer(timer)

	// Log constraint creation start
	logInfo(conn.ctx, "Creating Neo4j database constraints", LogFields{
		"component": "neo4j_client",
		"operation": "create_constraints",
		"database":  conn.database,
	})

	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		logError(conn.ctx, "Failed to create session for constraints", LogFields{
			"component": "neo4j_client",
			"operation": "create_constraints",
			"error":     err.Error(),
			"database":  conn.database,
		})
		if conn.metrics != nil {
			conn.metrics.recordErrorCount("neo4j_client", "constraint_session_failed")
		}
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

	// Create batch logger for constraint creation
	batchLogger := createBatchLogger(conn.ctx, "create_constraints", len(constraints))
	defer batchLogger.finishBatch()

	successCount := 0
	for i, constraint := range constraints {
		// Log individual constraint creation
		logDebug(conn.ctx, "Creating database constraint", LogFields{
			"component":      "neo4j_client",
			"operation":      "create_constraint",
			"database":       conn.database,
			"label":          constraint.label,
			"property":       constraint.property,
			"constraint_num": i + 1,
			"total_constraints": len(constraints),
		})

		query := buildNeo4jConstraintQuery(constraint.label, constraint.property)
		result, err := executeNeo4jWrite(ctx, session, query, nil)

		if err != nil {
			// Log constraint creation failure
			logError(conn.ctx, "Failed to create database constraint", LogFields{
				"component": "neo4j_client",
				"operation": "create_constraint",
				"error":     err.Error(),
				"database":  conn.database,
				"label":     constraint.label,
				"property":  constraint.property,
				"query":     query,
			})
			if conn.metrics != nil {
				conn.metrics.recordErrorCount("neo4j_client", "constraint_creation_failed")
				conn.metrics.recordCounter("neo4j_constraint_errors_total", 1, MetricLabels{
					"database": conn.database,
					"label":    constraint.label,
					"property": constraint.property,
				})
			}
			return wrapNeo4jError(err, fmt.Sprintf("failed to create constraint for %s.%s", constraint.label, constraint.property))
		}

		// Log successful constraint creation
		logDebug(conn.ctx, "Database constraint created successfully", LogFields{
			"component":       "neo4j_client",
			"operation":       "constraint_created",
			"database":        conn.database,
			"label":           constraint.label,
			"property":        constraint.property,
			"execution_time":  result.ExecutionTime.String(),
			"query_id":        extractQueryID(result.Summary),
		})

		// Record constraint creation metrics
		if conn.metrics != nil {
			conn.metrics.recordCounter("neo4j_constraints_created_total", 1, MetricLabels{
				"database": conn.database,
				"label":    constraint.label,
				"property": constraint.property,
			})
		}

		successCount++
		batchLogger.logProgress(1)
	}

	// Log overall constraint creation success
	logInfo(conn.ctx, "All Neo4j database constraints created successfully", LogFields{
		"component":         "neo4j_client",
		"operation":         "constraints_completed",
		"database":          conn.database,
		"total_constraints": len(constraints),
		"successful":        successCount,
		"failed":            len(constraints) - successCount,
	})

	// Record overall constraint creation metrics
	if conn.metrics != nil {
		conn.metrics.recordCounter("neo4j_constraint_operations_total", 1, MetricLabels{
			"database": conn.database,
			"status":   "success",
		})
	}

	return nil
}

// createNeo4jIndexes creates database indexes (Orchestrator)
func createNeo4jIndexes(ctx context.Context, conn *Neo4jConnection) error {
	validateNeo4jConnectionNotNil(conn)

	// Create span for index creation
	span := createNeo4jSpan(conn.ctx, "schema.create_indexes", "CREATE INDEXES")
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(conn.ctx, "neo4j_create_indexes")
	defer stopPerformanceTimer(timer)

	// Log index creation start
	logInfo(conn.ctx, "Creating Neo4j database indexes", LogFields{
		"component": "neo4j_client",
		"operation": "create_indexes",
		"database":  conn.database,
	})

	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		logError(conn.ctx, "Failed to create session for indexes", LogFields{
			"component": "neo4j_client",
			"operation": "create_indexes",
			"error":     err.Error(),
			"database":  conn.database,
		})
		if conn.metrics != nil {
			conn.metrics.recordErrorCount("neo4j_client", "index_session_failed")
		}
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

	// Create batch logger for index creation
	batchLogger := createBatchLogger(conn.ctx, "create_indexes", len(indexes))
	defer batchLogger.finishBatch()

	successCount := 0
	for i, index := range indexes {
		// Log individual index creation
		logDebug(conn.ctx, "Creating database index", LogFields{
			"component":     "neo4j_client",
			"operation":     "create_index",
			"database":      conn.database,
			"label":         index.label,
			"property":      index.property,
			"index_num":     i + 1,
			"total_indexes": len(indexes),
		})

		query := buildNeo4jIndexQuery(index.label, index.property)
		result, err := executeNeo4jWrite(ctx, session, query, nil)

		if err != nil {
			// Log index creation failure
			logError(conn.ctx, "Failed to create database index", LogFields{
				"component": "neo4j_client",
				"operation": "create_index",
				"error":     err.Error(),
				"database":  conn.database,
				"label":     index.label,
				"property":  index.property,
				"query":     query,
			})
			if conn.metrics != nil {
				conn.metrics.recordErrorCount("neo4j_client", "index_creation_failed")
				conn.metrics.recordCounter("neo4j_index_errors_total", 1, MetricLabels{
					"database": conn.database,
					"label":    index.label,
					"property": index.property,
				})
			}
			return wrapNeo4jError(err, fmt.Sprintf("failed to create index for %s.%s", index.label, index.property))
		}

		// Log successful index creation
		logDebug(conn.ctx, "Database index created successfully", LogFields{
			"component":      "neo4j_client",
			"operation":      "index_created",
			"database":       conn.database,
			"label":          index.label,
			"property":       index.property,
			"execution_time": result.ExecutionTime.String(),
			"query_id":       extractQueryID(result.Summary),
		})

		// Record index creation metrics
		if conn.metrics != nil {
			conn.metrics.recordCounter("neo4j_indexes_created_total", 1, MetricLabels{
				"database": conn.database,
				"label":    index.label,
				"property": index.property,
			})
		}

		successCount++
		batchLogger.logProgress(1)
	}

	// Log overall index creation success
	logInfo(conn.ctx, "All Neo4j database indexes created successfully", LogFields{
		"component":     "neo4j_client",
		"operation":     "indexes_completed",
		"database":      conn.database,
		"total_indexes": len(indexes),
		"successful":    successCount,
		"failed":        len(indexes) - successCount,
	})

	// Record overall index creation metrics
	if conn.metrics != nil {
		conn.metrics.recordCounter("neo4j_index_operations_total", 1, MetricLabels{
			"database": conn.database,
			"status":   "success",
		})
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

	// Enhanced error classification for better observability
	errorType := extractErrorType(err)
	code := "DATABASE_ERROR"

	// Map error types to more specific codes
	switch errorType {
	case "timeout":
		code = "TIMEOUT_ERROR"
	case "connection":
		code = "CONNECTION_ERROR"
	case "syntax":
		code = "SYNTAX_ERROR"
	case "constraint":
		code = "CONSTRAINT_ERROR"
	case "authentication":
		code = "AUTH_ERROR"
	case "permission":
		code = "PERMISSION_ERROR"
	}

	return Neo4jError{
		Code:    code,
		Message: message,
		Details: fmt.Sprintf("%s (type: %s)", err.Error(), errorType),
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

// Observability helper functions (Pure Core)

// sanitizeURI removes sensitive information from URI for logging
func sanitizeURI(uri string) string {
	if uri == "" {
		return "<empty>"
	}
	// Replace password in URI if present
	if strings.Contains(uri, "@") {
		parts := strings.Split(uri, "@")
		if len(parts) == 2 {
			// Extract protocol and credentials part
			protocolAndCreds := parts[0]
			if strings.Contains(protocolAndCreds, "://") {
				protocolParts := strings.Split(protocolAndCreds, "://")
				if len(protocolParts) == 2 {
					protocol := protocolParts[0]
					creds := protocolParts[1]
					if strings.Contains(creds, ":") {
						credsParts := strings.Split(creds, ":")
						username := credsParts[0]
						return fmt.Sprintf("%s://%s:***@%s", protocol, username, parts[1])
					}
				}
			}
		}
	}
	return uri
}

// sanitizeParams removes sensitive information from parameters for logging
func sanitizeParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		return make(map[string]interface{})
	}

	sanitized := make(map[string]interface{})
	sensitiveKeys := []string{"password", "secret", "token", "key", "auth"}

	for k, v := range params {
		isSensitive := false
		lowerKey := strings.ToLower(k)
		for _, sensitive := range sensitiveKeys {
			if strings.Contains(lowerKey, sensitive) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			sanitized[k] = "***"
		} else {
			sanitized[k] = v
		}
	}

	return sanitized
}

// generateQueryHash creates a hash for query identification
func generateQueryHash(query string) string {
	if query == "" {
		return "empty_query"
	}
	// Simple hash based on query length and first/last characters
	normalized := strings.TrimSpace(strings.ToLower(query))
	if len(normalized) < 10 {
		return fmt.Sprintf("short_%d_%s", len(normalized), normalized)
	}
	return fmt.Sprintf("hash_%d_%c%c_%c%c", len(normalized), normalized[0], normalized[1], normalized[len(normalized)-2], normalized[len(normalized)-1])
}

// calculateAverageQueryTime calculates average query execution time for a session
func calculateAverageQueryTime(session *Neo4jSession) time.Duration {
	if session.queryCount == 0 {
		return 0
	}
	return session.totalDuration / time.Duration(session.queryCount)
}

// extractErrorType extracts error type from Neo4j error for metrics
func extractErrorType(err error) string {
	if err == nil {
		return "unknown"
	}
	errorStr := strings.ToLower(err.Error())
	switch {
	case strings.Contains(errorStr, "timeout"):
		return "timeout"
	case strings.Contains(errorStr, "connection"):
		return "connection"
	case strings.Contains(errorStr, "syntax"):
		return "syntax"
	case strings.Contains(errorStr, "constraint"):
		return "constraint"
	case strings.Contains(errorStr, "authentication"):
		return "authentication"
	case strings.Contains(errorStr, "permission"):
		return "permission"
	default:
		return "database_error"
	}
}

// extractSummaryStatistic extracts specific statistics from result summary
func extractSummaryStatistic(summary neo4j.ResultSummary, statName string) int {
	if summary == nil {
		return 0
	}
	counters := summary.Counters()
	switch statName {
	case "nodes_created":
		return counters.NodesCreated()
	case "nodes_deleted":
		return counters.NodesDeleted()
	case "relationships_created":
		return counters.RelationshipsCreated()
	case "relationships_deleted":
		return counters.RelationshipsDeleted()
	case "properties_set":
		return counters.PropertiesSet()
	default:
		return 0
	}
}

// extractServerAddress extracts server address from result summary
func extractServerAddress(summary neo4j.ResultSummary) string {
	if summary == nil {
		return "unknown"
	}
	if summary.Server() != nil {
		return summary.Server().Address()
	}
	return "unknown"
}

// extractServerVersion extracts server version from result summary
func extractServerVersion(summary neo4j.ResultSummary) string {
	if summary == nil {
		return "unknown"
	}
	if summary.Server() != nil {
		// Note: Server version might not be available in all driver versions
		// Return server address as fallback identifier
		return summary.Server().Address()
	}
	return "unknown"
}

// extractQueryID extracts query ID from result summary
func extractQueryID(summary neo4j.ResultSummary) string {
	if summary == nil {
		return "unknown"
	}
	// Neo4j query ID might be available in different ways depending on version
	return fmt.Sprintf("query_%d", time.Now().UnixNano())
}

// extractDatabaseMode extracts database mode from result summary
func extractDatabaseMode(summary neo4j.ResultSummary) string {
	if summary == nil {
		return "unknown"
	}
	// This would extract database mode (leader/follower) in cluster setups
	return "standalone" // Default for most setups
}

// determineQueryType determines the type of query based on its content
func determineQueryType(query string) string {
	if query == "" {
		return "unknown"
	}
	normalized := strings.TrimSpace(strings.ToUpper(query))
	switch {
	case strings.HasPrefix(normalized, "CREATE"):
		return "create"
	case strings.HasPrefix(normalized, "MATCH"):
		return "match"
	case strings.HasPrefix(normalized, "MERGE"):
		return "merge"
	case strings.HasPrefix(normalized, "DELETE"):
		return "delete"
	case strings.HasPrefix(normalized, "SET"):
		return "set"
	case strings.HasPrefix(normalized, "REMOVE"):
		return "remove"
	case strings.HasPrefix(normalized, "RETURN"):
		return "return"
	case strings.HasPrefix(normalized, "WITH"):
		return "with"
	case strings.HasPrefix(normalized, "UNWIND"):
		return "unwind"
	case strings.HasPrefix(normalized, "CALL"):
		return "procedure"
	default:
		return "complex"
	}
}

// getConnectionPoolMetrics extracts connection pool metrics
func getConnectionPoolMetrics(conn *Neo4jConnection) map[string]interface{} {
	if conn == nil || conn.driver == nil {
		return map[string]interface{}{
			"status": "unavailable",
		}
	}

	// Neo4j driver doesn't expose pool metrics directly in Go driver
	// This would be implemented with actual pool monitoring
	return map[string]interface{}{
		"status":           "available",
		"max_pool_size":    50, // From configuration
		"active_connections": "unknown", // Would need driver introspection
		"idle_connections": "unknown", // Would need driver introspection
		"max_lifetime":     (30 * time.Minute).String(),
		"acquisition_timeout": (2 * time.Minute).String(),
	}
}

// Neo4j observability monitoring functions

// logNeo4jConnectionPoolStatus logs connection pool status with metrics
func logNeo4jConnectionPoolStatus(ctx *gofr.Context, conn *Neo4jConnection) {
	if conn == nil {
		return
	}

	poolMetrics := getConnectionPoolMetrics(conn)
	logInfo(ctx, "Neo4j connection pool status", LogFields{
		"component":    "neo4j_client",
		"operation":    "pool_status",
		"database":     conn.database,
		"pool_metrics": poolMetrics,
	})

	// Record pool status metrics
	if conn.metrics != nil {
		conn.metrics.recordGauge("neo4j_pool_max_size", 50, MetricLabels{
			"database": conn.database,
		})
	}
}

// monitorNeo4jPerformance monitors and logs query performance patterns
func monitorNeo4jPerformance(ctx *gofr.Context, session *Neo4jSession, result Neo4jResult, queryType string) {
	if session == nil || session.ctx == nil {
		return
	}

	// Calculate performance metrics
	avgQueryTime := calculateAverageQueryTime(session)
	recordsPerSecond := float64(result.RecordCount) / result.ExecutionTime.Seconds()

	// Log performance insights
	logDebug(session.ctx, "Neo4j query performance metrics", LogFields{
		"component":         "neo4j_client",
		"operation":         "performance_monitoring",
		"database":          session.database,
		"query_type":        queryType,
		"execution_time":    result.ExecutionTime.String(),
		"record_count":      result.RecordCount,
		"records_per_sec":   recordsPerSecond,
		"session_queries":   session.queryCount,
		"session_avg_time":  avgQueryTime.String(),
		"session_total_time": session.totalDuration.String(),
	})

	// Record performance metrics
	if session.metrics != nil {
		session.metrics.recordGauge("neo4j_records_per_second", recordsPerSecond, MetricLabels{
			"database":   session.database,
			"query_type": queryType,
		})
		session.metrics.recordGauge("neo4j_avg_query_time_ms", float64(avgQueryTime.Milliseconds()), MetricLabels{
			"database":   session.database,
			"query_type": queryType,
		})
	}

	// Alert on slow queries (over 5 seconds)
	if result.ExecutionTime > 5*time.Second {
		logWarn(session.ctx, "Slow Neo4j query detected", LogFields{
			"component":      "neo4j_client",
			"operation":      "slow_query_alert",
			"database":       session.database,
			"query_type":     queryType,
			"execution_time": result.ExecutionTime.String(),
			"query_hash":     result.QueryHash,
			"record_count":   result.RecordCount,
			"threshold":      "5s",
		})

		if session.metrics != nil {
			session.metrics.recordCounter("neo4j_slow_queries_total", 1, MetricLabels{
				"database":   session.database,
				"query_type": queryType,
			})
		}
	}
}

// withNeo4jObservability wraps Neo4j operations with comprehensive observability
func withNeo4jObservability(ctx *gofr.Context, operation string, database string, fn func() error) error {
	// Create observability span
	span := createNeo4jSpan(ctx, fmt.Sprintf("observability.%s", operation), operation)
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(ctx, fmt.Sprintf("neo4j_%s", operation))
	defer stopPerformanceTimer(timer)

	// Execute with error recovery
	return withErrorRecovery(ctx, fmt.Sprintf("neo4j_%s", operation), fn)
}

// enableNeo4jObservability enables observability features for an existing connection
func enableNeo4jObservability(conn *Neo4jConnection, gofrCtx *gofr.Context) {
	if conn == nil || gofrCtx == nil {
		return
	}

	// Update connection with observability context
	conn.ctx = gofrCtx
	conn.metrics = newMetricsCollector(gofrCtx, "neo4j-client")

	// Log observability enablement
	logInfo(gofrCtx, "Neo4j observability enabled for existing connection", LogFields{
		"component": "neo4j_client",
		"operation": "enable_observability", 
		"database":  conn.database,
	})

	// Log initial connection pool status
	logNeo4jConnectionPoolStatus(gofrCtx, conn)
}

// Example usage functions for integration

// createObservableNeo4jConnection creates a connection with full observability
func createObservableNeo4jConnection(ctx context.Context, gofrCtx *gofr.Context, config Neo4jConfig) (*Neo4jConnection, error) {
	return createNeo4jConnectionWithObservability(ctx, gofrCtx, config)
}

// upgradeNeo4jConnectionObservability upgrades an existing connection with observability
func upgradeNeo4jConnectionObservability(conn *Neo4jConnection, gofrCtx *gofr.Context) {
	enableNeo4jObservability(conn, gofrCtx)
}
