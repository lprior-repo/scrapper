package main

import (
	"fmt"
	"time"
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


