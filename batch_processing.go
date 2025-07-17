package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/samber/lo"
)

// BatchProcessor represents a generic batch processor
type BatchProcessor[T any] struct {
	Name           string
	BatchSize      int
	MaxConcurrency int
	Timeout        time.Duration
	RetryStrategy  RecoveryStrategy
}

// BatchResult represents the result of a batch operation
type BatchResult[T any] struct {
	Items      []T                    `json:"items"`
	Successful int                    `json:"successful"`
	Failed     int                    `json:"failed"`
	Duration   time.Duration          `json:"duration"`
	Errors     []AppError             `json:"errors,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// BatchProcessorFunc represents a function that processes a batch of items
type BatchProcessorFunc[T any] func(ctx context.Context, items []T) ([]T, error)

// createBatchProcessor creates a new batch processor (Pure Core)
func createBatchProcessor[T any](name string, batchSize, maxConcurrency int) BatchProcessor[T] {
	if name == "" {
		panic("Batch processor name cannot be empty")
	}
	if batchSize <= 0 {
		panic("Batch size must be positive")
	}
	if maxConcurrency <= 0 {
		panic("Max concurrency must be positive")
	}

	return BatchProcessor[T]{
		Name:           name,
		BatchSize:      batchSize,
		MaxConcurrency: maxConcurrency,
		Timeout:        30 * time.Second,
		RetryStrategy:  getDefaultRecoveryStrategy(),
	}
}

// splitIntoBatches splits items into batches of specified size (Pure Core)
func splitIntoBatches[T any](items []T, batchSize int) [][]T {
	if batchSize <= 0 {
		panic("Batch size must be positive")
	}
	if len(items) == 0 {
		return [][]T{}
	}

	batches := [][]T{}
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}

	return batches
}

// calculateOptimalBatches calculates optimal batch configuration (Pure Core)
func calculateOptimalBatches(totalItems, maxBatchSize, maxConcurrency int) (batchSize, numBatches int) {
	if totalItems <= 0 {
		return 0, 0
	}
	if maxBatchSize <= 0 {
		panic("Max batch size must be positive")
	}
	if maxConcurrency <= 0 {
		panic("Max concurrency must be positive")
	}

	// Start with maximum batch size
	optimalBatchSize := maxBatchSize

	// Calculate number of batches needed
	numBatches = (totalItems + optimalBatchSize - 1) / optimalBatchSize

	// If we have more batches than max concurrency, increase batch size
	if numBatches > maxConcurrency {
		optimalBatchSize = (totalItems + maxConcurrency - 1) / maxConcurrency
		numBatches = (totalItems + optimalBatchSize - 1) / optimalBatchSize
	}

	return optimalBatchSize, numBatches
}

// GitHub Organization Batch Processing

// GitHubBatchRequest represents a batch request for GitHub data
type GitHubBatchRequest struct {
	Organization string             `json:"organization"`
	Repositories []string           `json:"repositories,omitempty"`
	Teams        []string           `json:"teams,omitempty"`
	Options      GitHubBatchOptions `json:"options"`
	Context      ErrorContext       `json:"context"`
}

// GitHubBatchOptions represents options for GitHub batch processing
type GitHubBatchOptions struct {
	IncludeCodeowners    bool          `json:"include_codeowners"`
	IncludeMembers       bool          `json:"include_members"`
	IncludeCollaborators bool          `json:"include_collaborators"`
	MaxReposPerBatch     int           `json:"max_repos_per_batch"`
	MaxTeamsPerBatch     int           `json:"max_teams_per_batch"`
	Timeout              time.Duration `json:"timeout"`
	RespectRateLimit     bool          `json:"respect_rate_limit"`
}

// getDefaultGitHubBatchOptions returns default GitHub batch options (Pure Core)
func getDefaultGitHubBatchOptions() GitHubBatchOptions {
	return GitHubBatchOptions{
		IncludeCodeowners:    true,
		IncludeMembers:       true,
		IncludeCollaborators: false,
		MaxReposPerBatch:     50,
		MaxTeamsPerBatch:     25,
		Timeout:              5 * time.Minute,
		RespectRateLimit:     true,
	}
}

// validateGitHubBatchRequest validates a GitHub batch request (Pure Core)
func validateGitHubBatchRequest(request GitHubBatchRequest) error {
	if request.Organization == "" {
		return createRequiredFieldError("organization")
	}

	if request.Options.MaxReposPerBatch <= 0 {
		return createValidationError("max_repos_per_batch", "must be positive")
	}

	if request.Options.MaxTeamsPerBatch <= 0 {
		return createValidationError("max_teams_per_batch", "must be positive")
	}

	if request.Options.Timeout <= 0 {
		return createValidationError("timeout", "must be positive")
	}

	return nil
}

// splitGitHubBatchRequest splits a large GitHub request into smaller batches (Pure Core)
func splitGitHubBatchRequest(request GitHubBatchRequest) []GitHubBatchRequest {
	if len(request.Repositories) == 0 && len(request.Teams) == 0 {
		return []GitHubBatchRequest{request}
	}

	// Split repositories into batches
	repoBatches := splitIntoBatches(request.Repositories, request.Options.MaxReposPerBatch)
	teamBatches := splitIntoBatches(request.Teams, request.Options.MaxTeamsPerBatch)

	// Create batch requests - use the longer of the two lists
	maxBatches := lo.Max([]int{len(repoBatches), len(teamBatches)})
	if maxBatches == 0 {
		maxBatches = 1
	}

	batchRequests := make([]GitHubBatchRequest, maxBatches)

	for i := 0; i < maxBatches; i++ {
		batchReq := GitHubBatchRequest{
			Organization: request.Organization,
			Options:      request.Options,
			Context:      request.Context,
		}

		// Assign repositories to this batch
		if i < len(repoBatches) {
			batchReq.Repositories = repoBatches[i]
		}

		// Assign teams to this batch
		if i < len(teamBatches) {
			batchReq.Teams = teamBatches[i]
		}

		batchRequests[i] = batchReq
	}

	return batchRequests
}

// Database Batch Operations

// DatabaseBatchOperation represents a database batch operation
type DatabaseBatchOperation struct {
	Type       string                 `json:"type"`
	Query      string                 `json:"query"`
	Parameters map[string]interface{} `json:"parameters"`
	Priority   int                    `json:"priority"`
}

// DatabaseBatchRequest represents a batch of database operations
type DatabaseBatchRequest struct {
	Operations []DatabaseBatchOperation `json:"operations"`
	Options    DatabaseBatchOptions     `json:"options"`
	Context    ErrorContext             `json:"context"`
}

// DatabaseBatchOptions represents options for database batch processing
type DatabaseBatchOptions struct {
	MaxBatchSize    int           `json:"max_batch_size"`
	UseTransaction  bool          `json:"use_transaction"`
	Timeout         time.Duration `json:"timeout"`
	FailFast        bool          `json:"fail_fast"`
	OrderByPriority bool          `json:"order_by_priority"`
}

// getDefaultDatabaseBatchOptions returns default database batch options (Pure Core)
func getDefaultDatabaseBatchOptions() DatabaseBatchOptions {
	return DatabaseBatchOptions{
		MaxBatchSize:    100,
		UseTransaction:  true,
		Timeout:         30 * time.Second,
		FailFast:        false,
		OrderByPriority: true,
	}
}

// validateDatabaseBatchRequest validates a database batch request (Pure Core)
func validateDatabaseBatchRequest(request DatabaseBatchRequest) error {
	if len(request.Operations) == 0 {
		return createValidationError("operations", "at least one operation is required")
	}

	if request.Options.MaxBatchSize <= 0 {
		return createValidationError("max_batch_size", "must be positive")
	}

	if request.Options.Timeout <= 0 {
		return createValidationError("timeout", "must be positive")
	}

	// Validate each operation
	for i, op := range request.Operations {
		if op.Type == "" {
			return createValidationError(fmt.Sprintf("operations[%d].type", i), "cannot be empty")
		}
		if op.Query == "" {
			return createValidationError(fmt.Sprintf("operations[%d].query", i), "cannot be empty")
		}
	}

	return nil
}

// prepareDatabaseBatch prepares database operations for batch execution (Pure Core)
func prepareDatabaseBatch(request DatabaseBatchRequest) []DatabaseBatchRequest {
	operations := request.Operations

	// Sort by priority if requested
	if request.Options.OrderByPriority {
		sort.Slice(operations, func(i, j int) bool {
			return operations[i].Priority < operations[j].Priority
		})
	}

	// Split into batches
	batches := splitIntoBatches(operations, request.Options.MaxBatchSize)

	// Create batch requests
	batchRequests := make([]DatabaseBatchRequest, len(batches))
	for i, batch := range batches {
		batchRequests[i] = DatabaseBatchRequest{
			Operations: batch,
			Options:    request.Options,
			Context:    request.Context,
		}
	}

	return batchRequests
}

// Progress Tracking

// BatchProgress represents the progress of a batch operation
type BatchProgress struct {
	Operation           string                 `json:"operation"`
	TotalBatches        int                    `json:"total_batches"`
	CompletedBatches    int                    `json:"completed_batches"`
	FailedBatches       int                    `json:"failed_batches"`
	TotalItems          int                    `json:"total_items"`
	ProcessedItems      int                    `json:"processed_items"`
	FailedItems         int                    `json:"failed_items"`
	StartTime           time.Time              `json:"start_time"`
	LastUpdate          time.Time              `json:"last_update"`
	EstimatedCompletion time.Time              `json:"estimated_completion,omitempty"`
	Errors              []AppError             `json:"errors,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// createBatchProgress creates initial batch progress (Pure Core)
func createBatchProgress(operation string, totalBatches, totalItems int) BatchProgress {
	if operation == "" {
		panic("Operation cannot be empty")
	}
	if totalBatches < 0 || totalItems < 0 {
		panic("Totals cannot be negative")
	}

	now := time.Now()
	return BatchProgress{
		Operation:    operation,
		TotalBatches: totalBatches,
		TotalItems:   totalItems,
		StartTime:    now,
		LastUpdate:   now,
		Metadata:     make(map[string]interface{}),
	}
}

// updateBatchProgress updates batch progress with new completion info (Pure Core)
func updateBatchProgress(progress BatchProgress, completedBatches, processedItems int, errors []AppError) BatchProgress {
	if completedBatches < 0 || processedItems < 0 {
		panic("Completed counts cannot be negative")
	}

	now := time.Now()
	progress.CompletedBatches = completedBatches
	progress.ProcessedItems = processedItems
	progress.LastUpdate = now

	// Count failed items from errors
	progress.FailedItems = len(errors)
	progress.Errors = errors

	// Estimate completion time
	if completedBatches > 0 && completedBatches < progress.TotalBatches {
		elapsed := now.Sub(progress.StartTime)
		averageTimePerBatch := elapsed / time.Duration(completedBatches)
		remainingBatches := progress.TotalBatches - completedBatches
		remainingTime := averageTimePerBatch * time.Duration(remainingBatches)
		progress.EstimatedCompletion = now.Add(remainingTime)
	}

	return progress
}

// calculateProgressPercentage calculates completion percentage (Pure Core)
func calculateProgressPercentage(progress BatchProgress) float64 {
	if progress.TotalItems == 0 {
		return 0.0
	}
	return float64(progress.ProcessedItems) / float64(progress.TotalItems) * 100.0
}

// getProgressSummary generates a human-readable progress summary (Pure Core)
func getProgressSummary(progress BatchProgress) string {
	percentage := calculateProgressPercentage(progress)

	summary := fmt.Sprintf("%.1f%% complete (%d/%d items)",
		percentage, progress.ProcessedItems, progress.TotalItems)

	if progress.FailedItems > 0 {
		summary += fmt.Sprintf(", %d failed", progress.FailedItems)
	}

	if !progress.EstimatedCompletion.IsZero() && progress.EstimatedCompletion.After(time.Now()) {
		remaining := time.Until(progress.EstimatedCompletion)
		summary += fmt.Sprintf(", ~%v remaining", remaining.Round(time.Second))
	}

	return summary
}

// Batch Operation Statistics

// BatchStatistics represents statistics for batch operations
type BatchStatistics struct {
	TotalOperations  int           `json:"total_operations"`
	SuccessfulOps    int           `json:"successful_operations"`
	FailedOps        int           `json:"failed_operations"`
	TotalDuration    time.Duration `json:"total_duration"`
	AverageDuration  time.Duration `json:"average_duration"`
	ThroughputPerSec float64       `json:"throughput_per_second"`
	ErrorRate        float64       `json:"error_rate"`
	SuccessRate      float64       `json:"success_rate"`
}

// calculateBatchStatistics calculates statistics from batch results (Pure Core)
func calculateBatchStatistics[T any](results []BatchResult[T]) BatchStatistics {
	if len(results) == 0 {
		return BatchStatistics{}
	}

	totalOps := 0
	successfulOps := 0
	failedOps := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		totalOps += len(result.Items)
		successfulOps += result.Successful
		failedOps += result.Failed
		totalDuration += result.Duration
	}

	stats := BatchStatistics{
		TotalOperations: totalOps,
		SuccessfulOps:   successfulOps,
		FailedOps:       failedOps,
		TotalDuration:   totalDuration,
	}

	// Calculate averages and rates
	if len(results) > 0 {
		stats.AverageDuration = totalDuration / time.Duration(len(results))
	}

	if totalDuration > 0 {
		stats.ThroughputPerSec = float64(totalOps) / totalDuration.Seconds()
	}

	if totalOps > 0 {
		stats.ErrorRate = float64(failedOps) / float64(totalOps) * 100.0
		stats.SuccessRate = float64(successfulOps) / float64(totalOps) * 100.0
	}

	return stats
}

// formatBatchStatistics formats batch statistics for display (Pure Core)
func formatBatchStatistics(stats BatchStatistics) string {
	return fmt.Sprintf(
		"Operations: %d total, %d successful (%.1f%%), %d failed (%.1f%%), "+
			"Duration: %v total, %v average, Throughput: %.1f ops/sec",
		stats.TotalOperations,
		stats.SuccessfulOps, stats.SuccessRate,
		stats.FailedOps, stats.ErrorRate,
		stats.TotalDuration.Round(time.Millisecond),
		stats.AverageDuration.Round(time.Millisecond),
		stats.ThroughputPerSec,
	)
}

// Concurrent Batch Processing Utilities

// ConcurrentBatchConfig represents configuration for concurrent batch processing
type ConcurrentBatchConfig struct {
	MaxConcurrency   int           `json:"max_concurrency"`
	BatchSize        int           `json:"batch_size"`
	Timeout          time.Duration `json:"timeout"`
	FailFast         bool          `json:"fail_fast"`
	CollectResults   bool          `json:"collect_results"`
	ProgressCallback bool          `json:"progress_callback"`
}

// getDefaultConcurrentBatchConfig returns default concurrent batch config (Pure Core)
func getDefaultConcurrentBatchConfig() ConcurrentBatchConfig {
	return ConcurrentBatchConfig{
		MaxConcurrency:   5,
		BatchSize:        10,
		Timeout:          5 * time.Minute,
		FailFast:         false,
		CollectResults:   true,
		ProgressCallback: true,
	}
}

// validateConcurrentBatchConfig validates concurrent batch configuration (Pure Core)
func validateConcurrentBatchConfig(config ConcurrentBatchConfig) error {
	if config.MaxConcurrency <= 0 {
		return createValidationError("max_concurrency", "must be positive")
	}
	if config.BatchSize <= 0 {
		return createValidationError("batch_size", "must be positive")
	}
	if config.Timeout <= 0 {
		return createValidationError("timeout", "must be positive")
	}
	return nil
}
