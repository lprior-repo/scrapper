// Package main provides comprehensive observability utilities for the GoFr-based
// GitHub Codeowners Visualization application. This module includes:
//
// 1. Custom span creation utilities that wrap GoFr's ctx.Trace() functionality
// 2. Structured logging helpers with correlation IDs and consistent field formatting
// 3. Custom metrics helpers for business operations (scan duration, repository counts, etc.)
// 4. Error logging utilities with stack traces and context
// 5. Performance timing utilities for function execution tracking
//
// The utilities are designed to integrate seamlessly with GoFr's observability features
// and support high-cardinality logging for detailed debugging while maintaining
// performance in production environments.
//
// Key Features:
// - Pure function design following codebase patterns
// - Comprehensive error handling with structured context
// - Business-specific metrics for GitHub scanning operations
// - Correlation ID tracking across distributed operations
// - Stack trace capture for debugging
// - Performance monitoring with automatic metrics collection
// - High-cardinality logging with configurable enablement
//
// Usage Example:
//
//	// Create a span for tracking GitHub operations
//	span := createGitHubScanSpan(ctx, "microsoft", "fetch_repositories")
//	defer finishSpan(span)
//
//	// Log with structured context
//	logInfo(ctx, "Starting repository scan", LogFields{
//	    "organization": "microsoft",
//	    "max_repos": 100,
//	})
//
//	// Record business metrics
//	metrics := newMetricsCollector(ctx, "codeowners-scanner")
//	metrics.recordRepositoryCount("microsoft", 50)
//
//	// Track performance
//	timer := startPerformanceTimer(ctx, "github_scan")
//	defer stopPerformanceTimer(timer)
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"gofr.dev/pkg/gofr"
)

// ObservabilityConfig represents observability configuration
type ObservabilityConfig struct {
	EnableTracing      bool
	EnableMetrics      bool
	EnableHighCardinal bool
	ServiceName        string
	ServiceVersion     string
}

// SpanConfig represents custom span configuration
type SpanConfig struct {
	OperationName string
	Tags          map[string]interface{}
	Component     string
	Kind          string
}

// MetricLabels represents key-value pairs for metrics
type MetricLabels map[string]string

// LogFields represents structured log fields
type LogFields map[string]interface{}

// PerformanceTimer tracks function execution time
type PerformanceTimer struct {
	StartTime     time.Time
	OperationName string
	Context       *gofr.Context
	Tags          map[string]interface{}
}

// SpanWrapper wraps GoFr context with custom span functionality
type SpanWrapper struct {
	ctx       *gofr.Context
	spanName  string
	tags      map[string]interface{}
	startTime time.Time
}

// LogContext provides correlation IDs and consistent field formatting
type LogContext struct {
	CorrelationID string
	SessionID     string
	UserID        string
	RequestID     string
	TraceID       string
	Component     string
}

// MetricsCollector handles business metrics collection
type MetricsCollector struct {
	ctx         *gofr.Context
	serviceName string
}

// ErrorContext provides enhanced error logging with stack traces
type ErrorContext struct {
	Error       error
	Operation   string
	Component   string
	StackTrace  string
	Context     map[string]interface{}
	Severity    string
	Recoverable bool
	UserImpact  string
}

// createSpan creates a custom span wrapping GoFr's trace functionality
func createSpan(ctx *gofr.Context, config SpanConfig) *SpanWrapper {
	// Handle nil context gracefully
	if ctx == nil {
		return &SpanWrapper{
			ctx:       nil,
			spanName:  config.OperationName,
			tags:      config.Tags,
			startTime: time.Now(),
		}
	}

	// Create the span using GoFr's tracing
	_ = ctx.Trace(config.OperationName)

	// Log span start with attributes
	if ctx.Logger != nil {
		ctx.Logger.Debugf("Starting span: %s", config.OperationName)
		for key, value := range config.Tags {
			ctx.Logger.Debugf("Span attribute: %s = %v", key, value)
		}
	}

	return &SpanWrapper{
		ctx:       ctx, // Use original context
		spanName:  config.OperationName,
		tags:      config.Tags,
		startTime: time.Now(),
	}
}

// addSpanAttribute adds an attribute to the span (GoFr implementation dependent)
func addSpanAttribute(ctx *gofr.Context, key string, value interface{}) {
	// GoFr's context logger for span attributes
	// This logs span attributes as structured log entries
	if ctx != nil && ctx.Logger != nil {
		ctx.Logger.Debugf("Span attribute: %s = %v", key, value)
	}
}

// finishSpan completes the span and logs duration
func finishSpan(span *SpanWrapper) {
	if span == nil {
		return
	}

	duration := time.Since(span.startTime)
	
	// Handle nil context gracefully
	if span.ctx != nil && span.ctx.Logger != nil {
		span.ctx.Logger.Debugf("Span '%s' completed in %v", span.spanName, duration)

		// Log span completion with tags
		for key, value := range span.tags {
			span.ctx.Logger.Debugf("Span tag: %s = %v", key, value)
		}
	}
}

// createGitHubScanSpan creates a span for GitHub scanning operations
func createGitHubScanSpan(ctx *gofr.Context, organization string, operation string) *SpanWrapper {
	return createSpan(ctx, SpanConfig{
		OperationName: fmt.Sprintf("github.scan.%s", operation),
		Tags: map[string]interface{}{
			"organization":    organization,
			"operation":       operation,
			"component":       "github_scanner",
			"span.kind":       "client",
			"service.name":    "codeowners-scanner",
			"service.version": "1.0.0",
		},
		Component: "github_scanner",
		Kind:      "client",
	})
}

// createNeo4jSpan creates a span for Neo4j database operations
func createNeo4jSpan(ctx *gofr.Context, operation string, query string) *SpanWrapper {
	return createSpan(ctx, SpanConfig{
		OperationName: fmt.Sprintf("neo4j.%s", operation),
		Tags: map[string]interface{}{
			"operation":       operation,
			"db.type":         "neo4j",
			"db.statement":    truncateQuery(query, 100),
			"component":       "neo4j_client",
			"span.kind":       "client",
			"service.name":    "codeowners-scanner",
			"service.version": "1.0.0",
		},
		Component: "neo4j_client",
		Kind:      "client",
	})
}

// createAPISpan creates a span for API endpoint operations
func createAPISpan(ctx *gofr.Context, endpoint string, method string) *SpanWrapper {
	return createSpan(ctx, SpanConfig{
		OperationName: fmt.Sprintf("api.%s", endpoint),
		Tags: map[string]interface{}{
			"http.method":     method,
			"http.endpoint":   endpoint,
			"component":       "api_handler",
			"span.kind":       "server",
			"service.name":    "codeowners-scanner",
			"service.version": "1.0.0",
		},
		Component: "api_handler",
		Kind:      "server",
	})
}

// createLogContext creates a structured log context with correlation IDs
func createLogContext(ctx *gofr.Context, component string) LogContext {
	// Handle nil context gracefully
	if ctx == nil {
		return LogContext{
			CorrelationID: "no-correlation",
			SessionID:     "no-session",
			UserID:        "no-user",
			RequestID:     "no-request",
			TraceID:       "no-trace",
			Component:     component,
		}
	}
	
	return LogContext{
		CorrelationID: generateCorrelationID(ctx),
		SessionID:     extractSessionID(ctx),
		UserID:        extractUserID(ctx),
		RequestID:     extractRequestID(ctx),
		TraceID:       extractTraceID(ctx),
		Component:     component,
	}
}

// logWithContext logs a message with structured context and correlation IDs
func logWithContext(ctx *gofr.Context, level string, message string, fields LogFields) {
	// Return early if gofr context is nil to prevent panic
	if ctx == nil || ctx.Logger == nil {
		return
	}
	
	logCtx := createLogContext(ctx, extractComponent(fields))

	// Enhance fields with context
	enhancedFields := make(LogFields)
	for k, v := range fields {
		enhancedFields[k] = v
	}

	enhancedFields["correlation_id"] = logCtx.CorrelationID
	enhancedFields["trace_id"] = logCtx.TraceID
	enhancedFields["component"] = logCtx.Component
	enhancedFields["timestamp"] = time.Now().UTC().Format(time.RFC3339Nano)

	// Format the log message
	logMessage := formatLogMessage(message, enhancedFields)

	// Log based on level
	switch strings.ToLower(level) {
	case "debug":
		ctx.Logger.Debugf(logMessage)
	case "info":
		ctx.Logger.Infof(logMessage)
	case "warn", "warning":
		ctx.Logger.Warnf(logMessage)
	case "error":
		ctx.Logger.Errorf(logMessage)
	default:
		ctx.Logger.Infof(logMessage)
	}
}

// logInfo logs an info message with structured context
func logInfo(ctx *gofr.Context, message string, fields LogFields) {
	logWithContext(ctx, "info", message, fields)
}

// logError logs an error message with structured context
func logError(ctx *gofr.Context, message string, fields LogFields) {
	logWithContext(ctx, "error", message, fields)
}

// logWarn logs a warning message with structured context
func logWarn(ctx *gofr.Context, message string, fields LogFields) {
	logWithContext(ctx, "warn", message, fields)
}

// logDebug logs a debug message with structured context
func logDebug(ctx *gofr.Context, message string, fields LogFields) {
	logWithContext(ctx, "debug", message, fields)
}

// newMetricsCollector creates a new metrics collector
func newMetricsCollector(ctx *gofr.Context, serviceName string) *MetricsCollector {
	return &MetricsCollector{
		ctx:         ctx,
		serviceName: serviceName,
	}
}

// recordScanDuration records the duration of a scan operation
func (mc *MetricsCollector) recordScanDuration(organization string, duration time.Duration) {
	labels := MetricLabels{
		"organization": organization,
		"service":      mc.serviceName,
		"operation":    "scan",
	}

	mc.recordDuration("scan_duration_ms", duration, labels)
}

// recordRepositoryCount records the number of repositories processed
func (mc *MetricsCollector) recordRepositoryCount(organization string, count int) {
	labels := MetricLabels{
		"organization": organization,
		"service":      mc.serviceName,
	}

	mc.recordCounter("repositories_processed", count, labels)
}

// recordAPICallCount records API call metrics
func (mc *MetricsCollector) recordAPICallCount(service string, endpoint string, status int) {
	labels := MetricLabels{
		"service":     service,
		"endpoint":    endpoint,
		"status_code": fmt.Sprintf("%d", status),
		"app_service": mc.serviceName,
	}

	mc.recordCounter("api_calls_total", 1, labels)
}

// recordErrorCount records error metrics
func (mc *MetricsCollector) recordErrorCount(component string, errorType string) {
	labels := MetricLabels{
		"component":  component,
		"error_type": errorType,
		"service":    mc.serviceName,
	}

	mc.recordCounter("errors_total", 1, labels)
}

// recordDuration records a duration metric (placeholder for GoFr implementation)
func (mc *MetricsCollector) recordDuration(metricName string, duration time.Duration, labels MetricLabels) {
	// GoFr should provide metrics recording capabilities
	if mc.ctx != nil && mc.ctx.Logger != nil {
		mc.ctx.Logger.Infof("Metric [%s]: %v ms (labels: %v)", metricName, duration.Milliseconds(), labels)
	}
}

// recordCounter records a counter metric (placeholder for GoFr implementation)
func (mc *MetricsCollector) recordCounter(metricName string, value int, labels MetricLabels) {
	// GoFr should provide metrics recording capabilities
	if mc.ctx != nil && mc.ctx.Logger != nil {
		mc.ctx.Logger.Infof("Metric [%s]: %d (labels: %v)", metricName, value, labels)
	}
}

// recordGauge records a gauge metric (placeholder for GoFr implementation)
func (mc *MetricsCollector) recordGauge(metricName string, value float64, labels MetricLabels) {
	// GoFr should provide metrics recording capabilities
	if mc.ctx != nil && mc.ctx.Logger != nil {
		mc.ctx.Logger.Infof("Metric [%s]: %.2f (labels: %v)", metricName, value, labels)
	}
}

// logErrorWithStackTrace logs an error with enhanced context and stack trace
func logErrorWithStackTrace(ctx *gofr.Context, errCtx ErrorContext) {
	// Generate stack trace if not provided
	if errCtx.StackTrace == "" {
		errCtx.StackTrace = captureStackTrace(3) // Skip this function and 2 callers
	}

	fields := LogFields{
		"error":       errCtx.Error.Error(),
		"operation":   errCtx.Operation,
		"component":   errCtx.Component,
		"stack_trace": errCtx.StackTrace,
		"severity":    errCtx.Severity,
		"recoverable": errCtx.Recoverable,
		"user_impact": errCtx.UserImpact,
		"error_type":  fmt.Sprintf("%T", errCtx.Error),
	}

	// Add context fields
	for key, value := range errCtx.Context {
		fields[key] = value
	}

	logError(ctx, "Error occurred", fields)
}

// startPerformanceTimer starts a timer for performance tracking
func startPerformanceTimer(ctx *gofr.Context, operationName string) *PerformanceTimer {
	return &PerformanceTimer{
		StartTime:     time.Now(),
		OperationName: operationName,
		Context:       ctx,
		Tags: map[string]interface{}{
			"operation":  operationName,
			"start_time": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}
}

// stopPerformanceTimer stops the timer and logs the execution time
func stopPerformanceTimer(timer *PerformanceTimer) time.Duration {
	if timer == nil {
		return 0
	}

	duration := time.Since(timer.StartTime)

	// Handle nil context gracefully
	if timer.Context != nil {
		timer.Context.Logger.Infof("Performance [%s]: %v", timer.OperationName, duration)

		// Create metrics collector and record the duration
		metrics := newMetricsCollector(timer.Context, "codeowners-scanner")
		labels := MetricLabels{
			"operation": timer.OperationName,
			"service":   "codeowners-scanner",
		}
		metrics.recordDuration(fmt.Sprintf("%s_duration_ms", timer.OperationName), duration, labels)
	}

	return duration
}

// withPerformanceTracking wraps a function with performance tracking
func withPerformanceTracking(ctx *gofr.Context, operationName string, fn func() error) error {
	timer := startPerformanceTimer(ctx, operationName)
	defer stopPerformanceTimer(timer)

	return fn()
}

// withSpanTracking wraps a function with span tracking
func withSpanTracking(ctx *gofr.Context, config SpanConfig, fn func(*gofr.Context) error) error {
	span := createSpan(ctx, config)
	defer finishSpan(span)

	return fn(span.ctx)
}

// Helper functions

// generateCorrelationID generates a unique correlation ID
func generateCorrelationID(ctx *gofr.Context) string {
	// In a real implementation, this could use request headers or generate a UUID
	return fmt.Sprintf("corr_%d", time.Now().UnixNano())
}

// extractSessionID extracts session ID from context
func extractSessionID(ctx *gofr.Context) string {
	// Extract from request headers using Param method
	return ctx.Param("session_id")
}

// extractUserID extracts user ID from context
func extractUserID(ctx *gofr.Context) string {
	// Extract from authentication context using Param method
	return ctx.Param("user_id")
}

// extractRequestID extracts request ID from context
func extractRequestID(ctx *gofr.Context) string {
	// Extract from request headers using Param method
	return ctx.Param("request_id")
}

// extractTraceID extracts trace ID from context
func extractTraceID(ctx *gofr.Context) string {
	// Extract from tracing context using Param method
	return ctx.Param("trace_id")
}

// extractComponent extracts component name from log fields
func extractComponent(fields LogFields) string {
	if comp, exists := fields["component"]; exists {
		if compStr, ok := comp.(string); ok {
			return compStr
		}
	}
	return "unknown"
}

// formatLogMessage formats a log message with structured fields
func formatLogMessage(message string, fields LogFields) string {
	var parts []string
	parts = append(parts, message)

	for key, value := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	return strings.Join(parts, " ")
}

// truncateQuery truncates a long query for logging
func truncateQuery(query string, maxLength int) string {
	if len(query) <= maxLength {
		return query
	}
	return query[:maxLength] + "..."
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) string {
	const maxDepth = 32
	pcs := make([]uintptr, maxDepth)
	depth := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	var trace []string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		trace = append(trace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
	}

	return strings.Join(trace, "\n")
}

// Business-specific logging helpers

// logGitHubOperation logs GitHub API operations with consistent formatting
func logGitHubOperation(ctx *gofr.Context, operation string, organization string, details map[string]interface{}) {
	fields := LogFields{
		"component":    "github_client",
		"operation":    operation,
		"organization": organization,
		"service":      "github",
	}

	for key, value := range details {
		fields[key] = value
	}

	logInfo(ctx, fmt.Sprintf("GitHub %s operation", operation), fields)
}

// logNeo4jOperation logs Neo4j operations with consistent formatting
func logNeo4jOperation(ctx *gofr.Context, operation string, query string, params map[string]interface{}) {
	fields := LogFields{
		"component":   "neo4j_client",
		"operation":   operation,
		"query":       truncateQuery(query, 200),
		"db_type":     "neo4j",
		"param_count": len(params),
	}

	logInfo(ctx, fmt.Sprintf("Neo4j %s operation", operation), fields)
}

// logScanProgress logs scan progress with metrics
func logScanProgress(ctx *gofr.Context, organization string, progress map[string]interface{}) {
	fields := LogFields{
		"component":    "scanner",
		"operation":    "scan_progress",
		"organization": organization,
	}

	for key, value := range progress {
		fields[key] = value
	}

	logInfo(ctx, "Scan progress update", fields)
}

// High-cardinality logging for detailed debugging

// logHighCardinalityEvent logs events with high-cardinality data for debugging
func logHighCardinalityEvent(ctx *gofr.Context, eventType string, data map[string]interface{}) {
	if !isHighCardinalityEnabled() {
		return // Skip if high-cardinality logging is disabled
	}

	fields := LogFields{
		"event_type":   eventType,
		"cardinality":  "high",
		"debug_level":  "detailed",
		"timestamp_ns": time.Now().UnixNano(),
	}

	for key, value := range data {
		fields[key] = value
	}

	logDebug(ctx, "High-cardinality debug event", fields)
}

// isHighCardinalityEnabled checks if high-cardinality logging is enabled
func isHighCardinalityEnabled() bool {
	// This would typically check configuration or environment variables
	return true // For now, always enabled
}

// Metrics recording helpers for business operations

// recordBusinessMetrics records comprehensive business metrics
func recordBusinessMetrics(ctx *gofr.Context, organization string, scanResults map[string]interface{}) {
	metrics := newMetricsCollector(ctx, "codeowners-scanner")

	// Record repository metrics
	if repoCount, ok := scanResults["repository_count"].(int); ok {
		metrics.recordRepositoryCount(organization, repoCount)
	}

	// Record scan duration
	if duration, ok := scanResults["scan_duration"].(time.Duration); ok {
		metrics.recordScanDuration(organization, duration)
	}

	// Record coverage metrics
	if coverage, ok := scanResults["coverage_percentage"].(float64); ok {
		labels := MetricLabels{
			"organization": organization,
			"service":      "codeowners-scanner",
		}
		metrics.recordGauge("codeowners_coverage_percentage", coverage, labels)
	}
}

// Advanced observability utilities

// createObservabilityMiddleware creates middleware for automatic request tracking
func createObservabilityMiddleware(serviceName string) func(*gofr.Context, func()) {
	return func(ctx *gofr.Context, next func()) {
		// Extract request path and method from context
		requestPath := extractRequestPath(ctx)
		requestMethod := extractRequestMethod(ctx)

		// Start request span
		span := createAPISpan(ctx, requestPath, requestMethod)
		defer finishSpan(span)

		// Start performance timer
		timer := startPerformanceTimer(ctx, "api_request")
		defer func() {
			duration := stopPerformanceTimer(timer)

			// Record API metrics
			metrics := newMetricsCollector(ctx, serviceName)
			statusCode := extractStatusCode(ctx)
			metrics.recordAPICallCount("api", requestPath, statusCode)

			// Log request completion
			logInfo(ctx, "API request completed", LogFields{
				"method":      requestMethod,
				"path":        requestPath,
				"status_code": statusCode,
				"duration_ms": duration.Milliseconds(),
				"component":   "api_middleware",
			})
		}()

		// Execute next handler
		next()
	}
}

// withErrorRecovery wraps a function with error recovery and logging
func withErrorRecovery(ctx *gofr.Context, operation string, fn func() error) error {
	defer func() {
		if r := recover(); r != nil {
			// Log panic with stack trace
			logErrorWithStackTrace(ctx, ErrorContext{
				Error:       fmt.Errorf("panic recovered: %v", r),
				Operation:   operation,
				Component:   "error_recovery",
				StackTrace:  captureStackTrace(2),
				Severity:    "critical",
				Recoverable: false,
				UserImpact:  "service_disruption",
				Context: map[string]interface{}{
					"panic_value": r,
					"operation":   operation,
				},
			})

			// Record error metric
			metrics := newMetricsCollector(ctx, "codeowners-scanner")
			metrics.recordErrorCount("error_recovery", "panic")
		}
	}()

	return fn()
}

// createBatchLogger creates a logger for batch operations
func createBatchLogger(ctx *gofr.Context, batchName string, totalItems int) *BatchLogger {
	return &BatchLogger{
		ctx:         ctx,
		batchName:   batchName,
		totalItems:  totalItems,
		processed:   0,
		startTime:   time.Now(),
		lastLogTime: time.Now(),
	}
}

// BatchLogger tracks progress of batch operations
type BatchLogger struct {
	ctx         *gofr.Context
	batchName   string
	totalItems  int
	processed   int
	startTime   time.Time
	lastLogTime time.Time
}

// logProgress logs batch processing progress
func (bl *BatchLogger) logProgress(increment int) {
	bl.processed += increment

	// Log every 10% or every 30 seconds
	percentComplete := float64(bl.processed) / float64(bl.totalItems) * 100
	timeSinceLastLog := time.Since(bl.lastLogTime)

	if timeSinceLastLog > 30*time.Second || bl.processed%max(1, bl.totalItems/10) == 0 {
		logInfo(bl.ctx, "Batch processing progress", LogFields{
			"batch_name":          bl.batchName,
			"processed":           bl.processed,
			"total":               bl.totalItems,
			"percent_complete":    fmt.Sprintf("%.1f%%", percentComplete),
			"elapsed_time":        time.Since(bl.startTime).String(),
			"estimated_remaining": bl.estimateRemaining().String(),
			"component":           "batch_processor",
		})
		bl.lastLogTime = time.Now()
	}
}

// estimateRemaining estimates remaining processing time
func (bl *BatchLogger) estimateRemaining() time.Duration {
	if bl.processed == 0 {
		return 0
	}

	elapsed := time.Since(bl.startTime)
	remaining := bl.totalItems - bl.processed
	avgTimePerItem := elapsed / time.Duration(bl.processed)

	return avgTimePerItem * time.Duration(remaining)
}

// finishBatch logs batch completion
func (bl *BatchLogger) finishBatch() {
	duration := time.Since(bl.startTime)
	logInfo(bl.ctx, "Batch processing completed", LogFields{
		"batch_name":    bl.batchName,
		"total_items":   bl.totalItems,
		"processed":     bl.processed,
		"duration":      duration.String(),
		"items_per_sec": float64(bl.processed) / duration.Seconds(),
		"component":     "batch_processor",
	})
}

// Helper utility functions

// extractStatusCode extracts HTTP status code from GoFr context
func extractStatusCode(ctx *gofr.Context) int {
	// GoFr should provide a way to get response status
	// This is a placeholder implementation
	return 200 // Default to 200 if not available
}

// extractRequestPath extracts request path from GoFr context
func extractRequestPath(ctx *gofr.Context) string {
	// GoFr context may store request path in context values
	// For now, extract from available parameters or use a default
	if path := ctx.Param("path"); path != "" {
		return path
	}
	return "/api/unknown"
}

// extractRequestMethod extracts request method from GoFr context
func extractRequestMethod(ctx *gofr.Context) string {
	// GoFr context may store request method in context values
	// For now, extract from available parameters or use a default
	if method := ctx.Param("method"); method != "" {
		return method
	}
	return "GET"
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// createHealthCheckSpan creates a span for health check operations
func createHealthCheckSpan(ctx *gofr.Context, checkType string) *SpanWrapper {
	// Handle nil context gracefully
	if ctx == nil {
		return &SpanWrapper{
			ctx:       nil,
			spanName:  fmt.Sprintf("health.%s", checkType),
			tags:      make(map[string]interface{}),
			startTime: time.Now(),
		}
	}
	
	return createSpan(ctx, SpanConfig{
		OperationName: fmt.Sprintf("health.%s", checkType),
		Tags: map[string]interface{}{
			"check_type":      checkType,
			"component":       "health_checker",
			"span.kind":       "internal",
			"service.name":    "codeowners-scanner",
			"service.version": "1.0.0",
		},
		Component: "health_checker",
		Kind:      "internal",
	})
}

// logHealthCheckResult logs health check results with consistent formatting
func logHealthCheckResult(ctx *gofr.Context, checkType string, healthy bool, details map[string]interface{}) {
	level := "info"
	if !healthy {
		level = "error"
	}

	fields := LogFields{
		"component":  "health_checker",
		"check_type": checkType,
		"healthy":    healthy,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	for key, value := range details {
		fields[key] = value
	}

	message := fmt.Sprintf("Health check %s: %s", checkType, map[bool]string{true: "PASS", false: "FAIL"}[healthy])
	logWithContext(ctx, level, message, fields)
}

// Assertion helpers for observability

// assertSpanExists validates that a span is properly created
func assertSpanExists(span *SpanWrapper, operation string) {
	if span == nil {
		panic(fmt.Sprintf("span should exist for operation: %s", operation))
	}
	if span.spanName == "" {
		panic(fmt.Sprintf("span name should not be empty for operation: %s", operation))
	}
}

// assertMetricsCollector validates that metrics collector is properly initialized
func assertMetricsCollector(metrics *MetricsCollector) {
	if metrics == nil {
		panic("metrics collector should not be nil")
	}
	if metrics.ctx == nil {
		panic("metrics collector context should not be nil")
	}
	if metrics.serviceName == "" {
		panic("metrics collector service name should not be empty")
	}
}

// assertLogContext validates that log context is properly formed
func assertLogContext(logCtx LogContext) {
	if logCtx.Component == "" {
		panic("log context component should not be empty")
	}
	if logCtx.CorrelationID == "" {
		panic("log context correlation ID should not be empty")
	}
}
