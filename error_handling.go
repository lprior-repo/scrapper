package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
)

// ErrorType represents different types of errors in the system
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeNetwork        ErrorType = "network"
	ErrorTypeDatabase       ErrorType = "database"
	ErrorTypeRateLimit      ErrorType = "rate_limit"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeTimeout        ErrorType = "timeout"
	ErrorTypeInternal       ErrorType = "internal"
	ErrorTypeExternal       ErrorType = "external"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	SeverityCritical ErrorSeverity = "critical"
	SeverityHigh     ErrorSeverity = "high"
	SeverityMedium   ErrorSeverity = "medium"
	SeverityLow      ErrorSeverity = "low"
	SeverityInfo     ErrorSeverity = "info"
)

// AppError represents a structured application error
type AppError struct {
	Type        ErrorType     `json:"type"`
	Severity    ErrorSeverity `json:"severity"`
	Message     string        `json:"message"`
	Details     string        `json:"details,omitempty"`
	Code        string        `json:"code"`
	Context     string        `json:"context,omitempty"`
	Timestamp   time.Time     `json:"timestamp"`
	Recoverable bool          `json:"recoverable"`
	RetryAfter  time.Duration `json:"retry_after,omitempty"`
	Cause       error         `json:"-"`
}

// Error implements the error interface
func (e AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e AppError) Unwrap() error {
	return e.Cause
}

// createAppError creates a new application error (Pure Core)
func createAppError(errorType ErrorType, severity ErrorSeverity, code, message string) AppError {
	if code == "" {
		panic("Error code cannot be empty")
	}
	if message == "" {
		panic("Error message cannot be empty")
	}

	return AppError{
		Type:        errorType,
		Severity:    severity,
		Message:     message,
		Code:        code,
		Timestamp:   time.Now(),
		Recoverable: determineRecoverability(errorType),
		RetryAfter:  determineRetryDelay(errorType, severity),
	}
}

// determineRecoverability determines if an error type is recoverable (Pure Core)
func determineRecoverability(errorType ErrorType) bool {
	recoverableTypes := []ErrorType{
		ErrorTypeNetwork,
		ErrorTypeRateLimit,
		ErrorTypeTimeout,
		ErrorTypeExternal,
	}
	return lo.Contains(recoverableTypes, errorType)
}

// determineRetryDelay determines appropriate retry delay for error types (Pure Core)
func determineRetryDelay(errorType ErrorType, severity ErrorSeverity) time.Duration {
	baseDelays := map[ErrorType]time.Duration{
		ErrorTypeNetwork:   5 * time.Second,
		ErrorTypeRateLimit: 60 * time.Second,
		ErrorTypeTimeout:   10 * time.Second,
		ErrorTypeExternal:  30 * time.Second,
		ErrorTypeDatabase:  2 * time.Second,
	}

	severityMultipliers := map[ErrorSeverity]float64{
		SeverityCritical: 2.0,
		SeverityHigh:     1.5,
		SeverityMedium:   1.0,
		SeverityLow:      0.5,
		SeverityInfo:     0.1,
	}

	baseDelay, exists := baseDelays[errorType]
	if !exists {
		return 0
	}

	multiplier := severityMultipliers[severity]
	return time.Duration(float64(baseDelay) * multiplier)
}

// Validation Errors

// createValidationError creates a validation error (Pure Core)
func createValidationError(field, message string) AppError {
	return createAppError(
		ErrorTypeValidation,
		SeverityMedium,
		"VALIDATION_FAILED",
		fmt.Sprintf("Validation failed for field '%s': %s", field, message),
	)
}

// createRequiredFieldError creates a required field error (Pure Core)
func createRequiredFieldError(field string) AppError {
	return createAppError(
		ErrorTypeValidation,
		SeverityMedium,
		"REQUIRED_FIELD_MISSING",
		fmt.Sprintf("Required field '%s' is missing or empty", field),
	)
}

// createInvalidFormatError creates an invalid format error (Pure Core)
func createInvalidFormatError(field, expectedFormat string) AppError {
	return createAppError(
		ErrorTypeValidation,
		SeverityMedium,
		"INVALID_FORMAT",
		fmt.Sprintf("Field '%s' has invalid format. Expected: %s", field, expectedFormat),
	)
}

// Network Errors

// createNetworkError creates a network error (Pure Core)
func createNetworkError(operation string, err error) AppError {
	appErr := createAppError(
		ErrorTypeNetwork,
		SeverityHigh,
		"NETWORK_ERROR",
		fmt.Sprintf("Network error during %s", operation),
	)
	appErr.Cause = err
	if err != nil {
		appErr.Details = err.Error()
	}
	return appErr
}

// createTimeoutError creates a timeout error (Pure Core)
func createTimeoutError(operation string, duration time.Duration) AppError {
	return createAppError(
		ErrorTypeTimeout,
		SeverityHigh,
		"OPERATION_TIMEOUT",
		fmt.Sprintf("Operation '%s' timed out after %v", operation, duration),
	)
}

// createRateLimitError creates a rate limit error (Pure Core)
func createRateLimitError(resetTime time.Time, remaining int) AppError {
	appErr := createAppError(
		ErrorTypeRateLimit,
		SeverityMedium,
		"RATE_LIMIT_EXCEEDED",
		fmt.Sprintf("Rate limit exceeded. %d requests remaining. Reset at %s", remaining, resetTime.Format(time.RFC3339)),
	)
	appErr.RetryAfter = time.Until(resetTime)
	return appErr
}

// Database Errors

// createDatabaseError creates a database error (Pure Core)
func createDatabaseError(operation string, err error) AppError {
	appErr := createAppError(
		ErrorTypeDatabase,
		SeverityCritical,
		"DATABASE_ERROR",
		fmt.Sprintf("Database error during %s", operation),
	)
	appErr.Cause = err
	if err != nil {
		appErr.Details = err.Error()
	}
	return appErr
}

// createConnectionError creates a database connection error (Pure Core)
func createConnectionError(err error) AppError {
	appErr := createAppError(
		ErrorTypeDatabase,
		SeverityCritical,
		"DATABASE_CONNECTION_FAILED",
		"Failed to connect to database",
	)
	appErr.Cause = err
	if err != nil {
		appErr.Details = err.Error()
	}
	return appErr
}

// Authentication Errors

// createAuthenticationError creates an authentication error (Pure Core)
func createAuthenticationError(message string) AppError {
	return createAppError(
		ErrorTypeAuthentication,
		SeverityHigh,
		"AUTHENTICATION_FAILED",
		message,
	)
}

// createUnauthorizedError creates an unauthorized error (Pure Core)
func createUnauthorizedError(resource string) AppError {
	return createAppError(
		ErrorTypeAuthentication,
		SeverityHigh,
		"UNAUTHORIZED_ACCESS",
		fmt.Sprintf("Unauthorized access to resource: %s", resource),
	)
}

// Not Found Errors

// createNotFoundError creates a not found error (Pure Core)
func createNotFoundError(resource, identifier string) AppError {
	return createAppError(
		ErrorTypeNotFound,
		SeverityMedium,
		"RESOURCE_NOT_FOUND",
		fmt.Sprintf("%s with identifier '%s' not found", resource, identifier),
	)
}

// Internal Errors

// createInternalError creates an internal error (Pure Core)
func createInternalError(operation string, err error) AppError {
	appErr := createAppError(
		ErrorTypeInternal,
		SeverityCritical,
		"INTERNAL_ERROR",
		fmt.Sprintf("Internal error during %s", operation),
	)
	appErr.Cause = err
	if err != nil {
		appErr.Details = err.Error()
	}
	return appErr
}

// createConfigurationError creates a configuration error (Pure Core)
func createConfigurationError(setting, issue string) AppError {
	return createAppError(
		ErrorTypeInternal,
		SeverityCritical,
		"CONFIGURATION_ERROR",
		fmt.Sprintf("Configuration error for setting '%s': %s", setting, issue),
	)
}

// External Service Errors

// createExternalServiceError creates an external service error (Pure Core)
func createExternalServiceError(service, operation string, err error) AppError {
	appErr := createAppError(
		ErrorTypeExternal,
		SeverityHigh,
		"EXTERNAL_SERVICE_ERROR",
		fmt.Sprintf("External service '%s' error during %s", service, operation),
	)
	appErr.Cause = err
	if err != nil {
		appErr.Details = err.Error()
	}
	return appErr
}

// createGitHubAPIError creates a GitHub API error (Pure Core)
func createGitHubAPIError(operation string, statusCode int, message string) AppError {
	return createAppError(
		ErrorTypeExternal,
		SeverityHigh,
		"GITHUB_API_ERROR",
		fmt.Sprintf("GitHub API error during %s (HTTP %d): %s", operation, statusCode, message),
	)
}

// Error Recovery and Context

// RecoveryStrategy represents a strategy for error recovery
type RecoveryStrategy struct {
	MaxAttempts    int           `json:"max_attempts"`
	BaseDelay      time.Duration `json:"base_delay"`
	MaxDelay       time.Duration `json:"max_delay"`
	BackoffFactor  float64       `json:"backoff_factor"`
	RetryableTypes []ErrorType   `json:"retryable_types"`
}

// getDefaultRecoveryStrategy returns the default recovery strategy (Pure Core)
func getDefaultRecoveryStrategy() RecoveryStrategy {
	return RecoveryStrategy{
		MaxAttempts:   3,
		BaseDelay:     1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		RetryableTypes: []ErrorType{
			ErrorTypeNetwork,
			ErrorTypeTimeout,
			ErrorTypeRateLimit,
			ErrorTypeExternal,
		},
	}
}

// calculateBackoffDelay calculates the delay for exponential backoff (Pure Core)
func calculateBackoffDelay(attempt int, strategy RecoveryStrategy) time.Duration {
	if attempt <= 0 {
		return strategy.BaseDelay
	}

	// Exponential backoff: baseDelay * (backoffFactor ^ attempt)
	delay := float64(strategy.BaseDelay) * pow(strategy.BackoffFactor, float64(attempt))
	delayDuration := time.Duration(delay)

	// Cap at maximum delay
	if delayDuration > strategy.MaxDelay {
		return strategy.MaxDelay
	}

	return delayDuration
}

// isRetryable determines if an error is retryable based on strategy (Pure Core)
func isRetryable(err AppError, strategy RecoveryStrategy) bool {
	return err.Recoverable && lo.Contains(strategy.RetryableTypes, err.Type)
}

// shouldRetry determines if operation should be retried (Pure Core)
func shouldRetry(attempt int, err AppError, strategy RecoveryStrategy) bool {
	if attempt >= strategy.MaxAttempts {
		return false
	}
	return isRetryable(err, strategy)
}

// Helper function for power calculation (Pure Core)
func pow(base, exponent float64) float64 {
	if exponent == 0 {
		return 1
	}
	if exponent == 1 {
		return base
	}

	result := base
	for i := 1; i < int(exponent); i++ {
		result *= base
	}
	return result
}

// Error Context and Aggregation

// ErrorContext represents the context in which an error occurred
type ErrorContext struct {
	Operation   string                 `json:"operation"`
	Component   string                 `json:"component"`
	User        string                 `json:"user,omitempty"`
	Request     string                 `json:"request,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	StackTrace  []string               `json:"stack_trace,omitempty"`
	Correlation string                 `json:"correlation_id,omitempty"`
}

// enrichErrorWithContext adds context to an error (Pure Core)
func enrichErrorWithContext(err AppError, context ErrorContext) AppError {
	err.Context = fmt.Sprintf("%s.%s", context.Component, context.Operation)

	if context.User != "" {
		err.Details = fmt.Sprintf("%s [User: %s]", err.Details, context.User)
	}

	if context.Correlation != "" {
		err.Details = fmt.Sprintf("%s [Correlation: %s]", err.Details, context.Correlation)
	}

	return err
}

// AggregatedError represents multiple related errors
type AggregatedError struct {
	Operation string     `json:"operation"`
	Errors    []AppError `json:"errors"`
	Timestamp time.Time  `json:"timestamp"`
	Summary   string     `json:"summary"`
}

// Error implements the error interface
func (e AggregatedError) Error() string {
	if len(e.Errors) == 0 {
		return "No errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%s: %d errors occurred", e.Operation, len(e.Errors))
}

// createAggregatedError creates an aggregated error from multiple errors (Pure Core)
func createAggregatedError(operation string, errors []AppError) AggregatedError {
	if operation == "" {
		panic("Operation cannot be empty")
	}

	summary := generateErrorSummary(errors)
	return AggregatedError{
		Operation: operation,
		Errors:    errors,
		Timestamp: time.Now(),
		Summary:   summary,
	}
}

// generateErrorSummary generates a summary of multiple errors (Pure Core)
func generateErrorSummary(errors []AppError) string {
	if len(errors) == 0 {
		return "No errors"
	}

	if len(errors) == 1 {
		return errors[0].Message
	}

	// Group errors by type
	typeCount := make(map[ErrorType]int)
	for _, err := range errors {
		typeCount[err.Type]++
	}

	// Build summary
	parts := []string{}
	for errorType, count := range typeCount {
		if count == 1 {
			parts = append(parts, fmt.Sprintf("1 %s error", errorType))
		} else {
			parts = append(parts, fmt.Sprintf("%d %s errors", count, errorType))
		}
	}

	return strings.Join(parts, ", ")
}

// hasCriticalErrors checks if any error is critical (Pure Core)
func hasCriticalErrors(errors []AppError) bool {
	return lo.SomeBy(errors, func(err AppError) bool {
		return err.Severity == SeverityCritical
	})
}

// filterRetryableErrors filters errors that are retryable (Pure Core)
func filterRetryableErrors(errors []AppError) []AppError {
	return lo.Filter(errors, func(err AppError, _ int) bool {
		return err.Recoverable
	})
}

// getHighestSeverity returns the highest severity among errors (Pure Core)
func getHighestSeverity(errors []AppError) ErrorSeverity {
	if len(errors) == 0 {
		return SeverityInfo
	}

	severityOrder := []ErrorSeverity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		SeverityInfo,
	}

	for _, severity := range severityOrder {
		for _, err := range errors {
			if err.Severity == severity {
				return severity
			}
		}
	}

	return SeverityInfo
}
