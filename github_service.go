package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
)

// GitHubServiceConfig represents GitHub service configuration
type GitHubServiceConfig struct {
	Token        string
	BaseURL      string
	UserAgent    string
	Timeout      time.Duration
	MaxRetries   int
	RateLimitMin int
}

// RegisterGitHubService registers GitHub as an HTTP service in GoFr
func RegisterGitHubService(app *gofr.App, config GitHubServiceConfig) {
	// Register GitHub API as an HTTP service
	app.AddHTTPService("github", config.BaseURL)
}

// fetchGitHubOrganizationWithService fetches organization data using GoFr HTTP service
func fetchGitHubOrganizationWithService(ctx *gofr.Context, orgName string) (GitHubOrganization, error) {
	// Create span for tracking organization fetch
	span := createGitHubScanSpan(ctx, orgName, "fetch_organization")
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(ctx, "github_fetch_organization")
	defer stopPerformanceTimer(timer)

	// Initialize metrics collector
	metrics := newMetricsCollector(ctx, "codeowners-scanner")

	// Structured logging with context
	logInfo(ctx, "Starting GitHub organization fetch", LogFields{
		"component":    "github_client",
		"operation":    "fetch_organization",
		"organization": orgName,
		"service":      "github",
	})

	if orgName == "" {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("organization name is required but was empty"),
			Operation:   "fetch_organization",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "request_failed",
			Context: map[string]interface{}{
				"organization": orgName,
				"validation":   "missing_param",
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "validation_error")
		return GitHubOrganization{}, &gofrhttp.ErrorMissingParam{
			Params: []string{"organization_name"},
		}
	}

	resp, err := makeGitHubOrgAPIRequest(ctx, orgName)
	if err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "fetch_organization",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "api_request_failed",
			Context: map[string]interface{}{
				"organization": orgName,
				"api_endpoint": fmt.Sprintf("orgs/%s", orgName),
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "api_request_error")
		return GitHubOrganization{}, err
	}
	defer resp.Body.Close()

	// Record API call metrics
	metrics.recordAPICallCount("github", "organization", resp.StatusCode)

	org, err := parseGitHubOrgResponse(ctx, resp, orgName)
	if err != nil {
		metrics.recordErrorCount("github_client", "response_parse_error")
		return GitHubOrganization{}, err
	}

	logInfo(ctx, "Successfully fetched GitHub organization", LogFields{
		"component":       "github_client",
		"operation":       "fetch_organization",
		"organization":     orgName,
		"organization_id":  org.ID,
		"public_repos":    org.PublicRepos,
		"followers":       org.Followers,
		"response_status": resp.StatusCode,
	})

	return org, nil
}

// makeGitHubOrgAPIRequest makes API request to GitHub for organization data (Pure Core)
func makeGitHubOrgAPIRequest(ctx *gofr.Context, orgName string) (*http.Response, error) {
	githubSvc := ctx.GetHTTPService("github")
	headers := buildGitHubRequestHeaders()

	// Log API request with structured context
	logDebug(ctx, "Making GitHub API request", LogFields{
		"component":    "github_client",
		"operation":    "api_request",
		"organization": orgName,
		"endpoint":     fmt.Sprintf("orgs/%s", orgName),
		"headers_count": len(headers),
	})

	// Start timing for API call
	apiTimer := startPerformanceTimer(ctx, "github_api_call")
	defer stopPerformanceTimer(apiTimer)

	resp, err := githubSvc.GetWithHeaders(ctx, fmt.Sprintf("orgs/%s", orgName), nil, headers)
	if err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "api_request",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "api_unavailable",
			Context: map[string]interface{}{
				"organization": orgName,
				"endpoint":     fmt.Sprintf("orgs/%s", orgName),
				"error_type":   fmt.Sprintf("%T", err),
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		return nil, &gofrhttp.ErrorRequestTimeout{}
	}

	// Log rate limit information if available
	logRateLimitInfo(ctx, resp)

	return resp, nil
}

// parseGitHubOrgResponse parses GitHub organization API response (Pure Core)
func parseGitHubOrgResponse(ctx *gofr.Context, resp *http.Response, orgName string) (GitHubOrganization, error) {
	// Log response details with structured context
	logInfo(ctx, "Processing GitHub API response", LogFields{
		"component":    "github_client",
		"operation":    "parse_response",
		"organization": orgName,
		"status_code":  resp.StatusCode,
		"content_type": resp.Header.Get("Content-Type"),
	})

	if resp.StatusCode == http.StatusNotFound {
		logWarn(ctx, "Organization not found", LogFields{
			"component":    "github_client",
			"organization": orgName,
			"status_code":  resp.StatusCode,
			"error_type":   "not_found",
		})
		return GitHubOrganization{}, &gofrhttp.ErrorEntityNotFound{
			Name:  "organization",
			Value: orgName,
		}
	}

	if resp.StatusCode != http.StatusOK {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("GitHub API returned error status %d", resp.StatusCode),
			Operation:   "parse_response",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "api_error",
			Context: map[string]interface{}{
				"organization": orgName,
				"status_code":  resp.StatusCode,
				"response_headers": extractResponseHeaders(resp),
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		return GitHubOrganization{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"github_api_status", fmt.Sprintf("status_code_%d", resp.StatusCode)},
		}
	}

	var org GitHubOrganization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "json_decode",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: false,
			UserImpact:  "data_corruption",
			Context: map[string]interface{}{
				"organization": orgName,
				"content_type": resp.Header.Get("Content-Type"),
				"content_length": resp.Header.Get("Content-Length"),
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		return GitHubOrganization{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_format", err.Error()},
		}
	}

	logInfo(ctx, "Successfully parsed organization response", LogFields{
		"component":      "github_client",
		"operation":      "parse_response",
		"organization":   org.Login,
		"organization_id": org.ID,
		"public_repos":   org.PublicRepos,
		"followers":      org.Followers,
		"created_at":     org.CreatedAt,
	})
	return org, nil
}

// fetchGitHubRepositoriesWithService fetches repositories using GoFr HTTP service
func fetchGitHubRepositoriesWithService(ctx *gofr.Context, orgName string, maxRepos int) ([]GitHubRepository, error) {
	// Create span for tracking repository fetch
	span := createGitHubScanSpan(ctx, orgName, "fetch_repositories")
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(ctx, "github_fetch_repositories")
	defer stopPerformanceTimer(timer)

	// Initialize metrics collector
	metrics := newMetricsCollector(ctx, "codeowners-scanner")

	// Structured logging with context
	logInfo(ctx, "Starting GitHub repositories fetch", LogFields{
		"component":    "github_client",
		"operation":    "fetch_repositories",
		"organization": orgName,
		"max_repos":    maxRepos,
		"service":      "github",
	})

	if err := validateRepositoryParams(orgName, maxRepos); err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "validate_params",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "request_failed",
			Context: map[string]interface{}{
				"organization": orgName,
				"max_repos":    maxRepos,
				"validation":   "param_error",
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "validation_error")
		return nil, err
	}

	githubSvc := ctx.GetHTTPService("github")
	allRepos, err := fetchAllRepositoryPages(ctx, githubSvc, orgName, maxRepos)
	if err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "fetch_repositories",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "api_request_failed",
			Context: map[string]interface{}{
				"organization": orgName,
				"max_repos":    maxRepos,
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "fetch_error")
		return nil, err
	}

	// Record repository metrics
	metrics.recordRepositoryCount(orgName, len(allRepos))

	logInfo(ctx, "Successfully fetched GitHub repositories", LogFields{
		"component":      "github_client",
		"operation":      "fetch_repositories",
		"organization":   orgName,
		"max_repos":      maxRepos,
		"fetched_repos":  len(allRepos),
		"fetch_complete": len(allRepos) < maxRepos,
	})
	return allRepos, nil
}

// validateRepositoryParams validates input parameters for repository fetching
func validateRepositoryParams(orgName string, maxRepos int) error {
	if orgName == "" {
		return &gofrhttp.ErrorMissingParam{
			Params: []string{"organization_name"},
		}
	}

	if maxRepos <= 0 {
		return &gofrhttp.ErrorInvalidParam{
			Params: []string{"max_repos", fmt.Sprintf("%d", maxRepos)},
		}
	}

	return nil
}

// fetchAllRepositoryPages fetches all repository pages up to maxRepos
func fetchAllRepositoryPages(ctx *gofr.Context, githubSvc any, orgName string, maxRepos int) ([]GitHubRepository, error) {
	// Create batch logger for pagination progress
	batchLogger := createBatchLogger(ctx, "repository_pagination", maxRepos)
	defer batchLogger.finishBatch()

	var allRepos []GitHubRepository
	page := 1
	perPage := 100

	logInfo(ctx, "Starting repository pagination", LogFields{
		"component":    "github_client",
		"operation":    "paginate_repositories",
		"organization": orgName,
		"max_repos":    maxRepos,
		"per_page":     perPage,
	})

	for len(allRepos) < maxRepos {
		repos, shouldContinue, err := fetchRepositoryPage(ctx, githubSvc, orgName, page, perPage)
		if err != nil {
			errCtx := ErrorContext{
				Error:       err,
				Operation:   "fetch_page",
				Component:   "github_client",
				Severity:    "error",
				Recoverable: true,
				UserImpact:  "partial_data",
				Context: map[string]interface{}{
					"organization": orgName,
					"page":         page,
					"per_page":     perPage,
					"repos_so_far": len(allRepos),
				},
			}
			logErrorWithStackTrace(ctx, errCtx)
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		// Log pagination progress
		batchLogger.logProgress(len(repos))
		logDebug(ctx, "Repository page processed", LogFields{
			"component":      "github_client",
			"operation":      "paginate_repositories",
			"organization":   orgName,
			"page":           page,
			"repos_in_page":  len(repos),
			"total_repos":    len(allRepos),
			"should_continue": shouldContinue,
		})

		page++

		if !shouldContinue {
			logInfo(ctx, "Repository pagination completed - no more pages", LogFields{
				"component":    "github_client",
				"organization": orgName,
				"total_pages":  page - 1,
				"total_repos":  len(allRepos),
			})
			break
		}
	}

	return limitRepositories(ctx, allRepos, maxRepos, orgName), nil
}

// fetchRepositoryPage fetches a single page of repositories
func fetchRepositoryPage(ctx *gofr.Context, githubSvc any, orgName string, page, perPage int) ([]GitHubRepository, bool, error) {
	// Start timer for page fetch
	pageTimer := startPerformanceTimer(ctx, fmt.Sprintf("github_fetch_page_%d", page))
	defer stopPerformanceTimer(pageTimer)

	logDebug(ctx, "Fetching repository page", LogFields{
		"component":    "github_client",
		"operation":    "fetch_page",
		"organization": orgName,
		"page":         page,
		"per_page":     perPage,
	})

	resp, err := executeRepositoryRequest(ctx, githubSvc, orgName, page, perPage)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	// Log rate limit information
	logRateLimitInfo(ctx, resp)

	if err := validateRepositoryResponse(ctx, resp, orgName, page); err != nil {
		return nil, false, err
	}

	repos, err := decodeRepositoryResponse(ctx, resp, orgName)
	if err != nil {
		return nil, false, err
	}

	shouldContinue := len(repos) > 0 && len(repos) == perPage

	logInfo(ctx, "Repository page fetched successfully", LogFields{
		"component":       "github_client",
		"operation":       "fetch_page",
		"organization":    orgName,
		"page":            page,
		"repos_in_page":   len(repos),
		"should_continue": shouldContinue,
		"response_status": resp.StatusCode,
	})

	return repos, shouldContinue, nil
}

// executeRepositoryRequest executes a repository API request
func executeRepositoryRequest(ctx *gofr.Context, githubSvc any, orgName string, page, perPage int) (*http.Response, error) {
	query := map[string]any{
		"page":     fmt.Sprintf("%d", page),
		"per_page": fmt.Sprintf("%d", perPage),
		"sort":     "updated",
	}

	headers := buildGitHubRequestHeaders()
	endpoint := fmt.Sprintf("orgs/%s/repos", orgName)

	// Log API request details
	logDebug(ctx, "Executing repository API request", LogFields{
		"component":    "github_client",
		"operation":    "api_request",
		"organization": orgName,
		"endpoint":     endpoint,
		"page":         page,
		"per_page":     perPage,
		"query_params": len(query),
		"headers_count": len(headers),
	})

	// Start API call timer
	apiTimer := startPerformanceTimer(ctx, "github_api_call_repos")
	defer stopPerformanceTimer(apiTimer)
	
	// Get the GitHub service from context (same pattern as working organization request)
	githubHttpSvc := ctx.GetHTTPService("github")
	resp, err := githubHttpSvc.GetWithHeaders(ctx, endpoint, query, headers)
	if err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "api_request",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "api_unavailable",
			Context: map[string]interface{}{
				"organization": orgName,
				"endpoint":     endpoint,
				"page":         page,
				"per_page":     perPage,
				"error_type":   fmt.Sprintf("%T", err),
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		return nil, &gofrhttp.ErrorRequestTimeout{}
	}

	return resp, nil
}

// validateRepositoryResponse validates the HTTP response from GitHub API
func validateRepositoryResponse(ctx *gofr.Context, resp *http.Response, orgName string, page int) error {
	logDebug(ctx, "Validating repository API response", LogFields{
		"component":    "github_client",
		"operation":    "validate_response",
		"organization": orgName,
		"page":         page,
		"status_code":  resp.StatusCode,
		"content_type": resp.Header.Get("Content-Type"),
	})

	if resp.StatusCode == http.StatusNotFound {
		logWarn(ctx, "No repositories found for organization", LogFields{
			"component":    "github_client",
			"organization": orgName,
			"page":         page,
			"status_code":  resp.StatusCode,
			"error_type":   "not_found",
		})
		return &gofrhttp.ErrorEntityNotFound{
			Name:  "organization_repositories",
			Value: orgName,
		}
	}

	if resp.StatusCode != http.StatusOK {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("GitHub API returned error status %d for repositories", resp.StatusCode),
			Operation:   "validate_response",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "api_error",
			Context: map[string]interface{}{
				"organization": orgName,
				"page":         page,
				"status_code":  resp.StatusCode,
				"response_headers": extractResponseHeaders(resp),
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		return &gofrhttp.ErrorInvalidParam{
			Params: []string{"github_api_status", fmt.Sprintf("status_code_%d", resp.StatusCode)},
		}
	}

	return nil
}

// decodeRepositoryResponse decodes the JSON response into GitHubRepository slice
func decodeRepositoryResponse(ctx *gofr.Context, resp *http.Response, orgName string) ([]GitHubRepository, error) {
	var repos []GitHubRepository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		errCtx := ErrorContext{
			Error:       err,
			Operation:   "json_decode",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: false,
			UserImpact:  "data_corruption",
			Context: map[string]interface{}{
				"organization": orgName,
				"content_type": resp.Header.Get("Content-Type"),
				"content_length": resp.Header.Get("Content-Length"),
				"response_size": resp.ContentLength,
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		return nil, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_format", err.Error()},
		}
	}

	logDebug(ctx, "Successfully decoded repositories response", LogFields{
		"component":    "github_client",
		"operation":    "json_decode",
		"organization": orgName,
		"repo_count":   len(repos),
		"content_type": resp.Header.Get("Content-Type"),
	})

	return repos, nil
}

// limitRepositories limits the number of repositories to maxRepos
func limitRepositories(ctx *gofr.Context, allRepos []GitHubRepository, maxRepos int, orgName string) []GitHubRepository {
	if len(allRepos) > maxRepos {
		logInfo(ctx, "Limiting repositories to maximum requested", LogFields{
			"component":      "github_client",
			"operation":      "limit_repositories",
			"organization":   orgName,
			"total_fetched":  len(allRepos),
			"max_requested": maxRepos,
			"limited":        true,
		})
		return allRepos[:maxRepos]
	}

	logDebug(ctx, "No repository limiting needed", LogFields{
		"component":      "github_client",
		"operation":      "limit_repositories",
		"organization":   orgName,
		"total_fetched":  len(allRepos),
		"max_requested": maxRepos,
		"limited":        false,
	})
	return allRepos
}

// fetchGitHubTeamsWithService fetches teams using GoFr HTTP service
func fetchGitHubTeamsWithService(ctx *gofr.Context, orgName string, maxTeams int) ([]GitHubTeam, error) {
	// Create span for tracking team fetch
	span := createGitHubScanSpan(ctx, orgName, "fetch_teams")
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(ctx, "github_fetch_teams")
	defer stopPerformanceTimer(timer)

	// Initialize metrics collector
	metrics := newMetricsCollector(ctx, "codeowners-scanner")

	// Structured logging with context
	logInfo(ctx, "Starting GitHub teams fetch", LogFields{
		"component":    "github_client",
		"operation":    "fetch_teams",
		"organization": orgName,
		"max_teams":    maxTeams,
		"service":      "github",
	})

	if orgName == "" {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("organization name is required but was empty"),
			Operation:   "validate_params",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "request_failed",
			Context: map[string]interface{}{
				"organization": orgName,
				"validation":   "missing_param",
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "validation_error")
		return nil, &gofrhttp.ErrorMissingParam{
			Params: []string{"organization_name"},
		}
	}

	if maxTeams <= 0 {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("max_teams must be positive, got %d", maxTeams),
			Operation:   "validate_params",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "request_failed",
			Context: map[string]interface{}{
				"max_teams":    maxTeams,
				"validation":   "invalid_param",
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "validation_error")
		return nil, &gofrhttp.ErrorInvalidParam{
			Params: []string{"max_teams", fmt.Sprintf("%d", maxTeams)},
		}
	}

	githubSvc := ctx.GetHTTPService("github")

	// Create batch logger for team pagination
	batchLogger := createBatchLogger(ctx, "team_pagination", maxTeams)
	defer batchLogger.finishBatch()

	var allTeams []GitHubTeam
	page := 1
	perPage := 100

	logInfo(ctx, "Starting team pagination", LogFields{
		"component":    "github_client",
		"operation":    "paginate_teams",
		"organization": orgName,
		"max_teams":    maxTeams,
		"per_page":     perPage,
	})

	for len(allTeams) < maxTeams {
		// Start page timer
		pageTimer := startPerformanceTimer(ctx, fmt.Sprintf("github_fetch_teams_page_%d", page))

		query := map[string]any{
			"page":     fmt.Sprintf("%d", page),
			"per_page": fmt.Sprintf("%d", perPage),
		}

		headers := buildGitHubRequestHeaders()
		endpoint := fmt.Sprintf("orgs/%s/teams", orgName)

		logDebug(ctx, "Fetching teams page", LogFields{
			"component":    "github_client",
			"operation":    "fetch_teams_page",
			"organization": orgName,
			"page":         page,
			"endpoint":     endpoint,
		})

		resp, err := githubSvc.GetWithHeaders(ctx, endpoint, query, headers)
		if err != nil {
			stopPerformanceTimer(pageTimer)
			errCtx := ErrorContext{
				Error:       err,
				Operation:   "api_request",
				Component:   "github_client",
				Severity:    "error",
				Recoverable: true,
				UserImpact:  "api_unavailable",
				Context: map[string]interface{}{
					"organization": orgName,
					"endpoint":     endpoint,
					"page":         page,
					"error_type":   fmt.Sprintf("%T", err),
				},
			}
			logErrorWithStackTrace(ctx, errCtx)
			metrics.recordErrorCount("github_client", "api_request_error")
			return nil, &gofrhttp.ErrorRequestTimeout{}
		}
		defer resp.Body.Close()

		// Log rate limit information
		logRateLimitInfo(ctx, resp)

		// Record API call metrics
		metrics.recordAPICallCount("github", "teams", resp.StatusCode)

		if resp.StatusCode == http.StatusNotFound {
			stopPerformanceTimer(pageTimer)
			logWarn(ctx, "No teams found for organization", LogFields{
				"component":    "github_client",
				"organization": orgName,
				"page":         page,
				"status_code":  resp.StatusCode,
				"error_type":   "not_found",
			})
			return nil, &gofrhttp.ErrorEntityNotFound{
				Name:  "organization_teams",
				Value: orgName,
			}
		}

		if resp.StatusCode != http.StatusOK {
			stopPerformanceTimer(pageTimer)
			errCtx := ErrorContext{
				Error:       fmt.Errorf("GitHub API returned error status %d for teams", resp.StatusCode),
				Operation:   "validate_response",
				Component:   "github_client",
				Severity:    "error",
				Recoverable: true,
				UserImpact:  "api_error",
				Context: map[string]interface{}{
					"organization": orgName,
					"page":         page,
					"status_code":  resp.StatusCode,
					"response_headers": extractResponseHeaders(resp),
				},
			}
			logErrorWithStackTrace(ctx, errCtx)
			metrics.recordErrorCount("github_client", "api_error")
			return nil, &gofrhttp.ErrorInvalidParam{
				Params: []string{"github_api_status", fmt.Sprintf("status_code_%d", resp.StatusCode)},
			}
		}

		var teams []GitHubTeam
		if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
			stopPerformanceTimer(pageTimer)
			errCtx := ErrorContext{
				Error:       err,
				Operation:   "json_decode",
				Component:   "github_client",
				Severity:    "error",
				Recoverable: false,
				UserImpact:  "data_corruption",
				Context: map[string]interface{}{
					"organization": orgName,
					"page":         page,
					"content_type": resp.Header.Get("Content-Type"),
					"content_length": resp.Header.Get("Content-Length"),
				},
			}
			logErrorWithStackTrace(ctx, errCtx)
			metrics.recordErrorCount("github_client", "decode_error")
			return nil, &gofrhttp.ErrorInvalidParam{
				Params: []string{"response_format", err.Error()},
			}
		}

		stopPerformanceTimer(pageTimer)

		if len(teams) == 0 {
			logInfo(ctx, "Team pagination completed - no more teams", LogFields{
				"component":    "github_client",
				"organization": orgName,
				"total_pages":  page,
				"total_teams":  len(allTeams),
			})
			break
		}

		allTeams = append(allTeams, teams...)

		// Log pagination progress
		batchLogger.logProgress(len(teams))
		logDebug(ctx, "Teams page processed", LogFields{
			"component":     "github_client",
			"operation":     "paginate_teams",
			"organization":  orgName,
			"page":          page,
			"teams_in_page": len(teams),
			"total_teams":   len(allTeams),
		})

		page++

		if len(teams) < perPage {
			logInfo(ctx, "Team pagination completed - partial page", LogFields{
				"component":     "github_client",
				"organization":  orgName,
				"total_pages":   page - 1,
				"total_teams":   len(allTeams),
				"teams_in_page": len(teams),
			})
			break
		}
	}

	if len(allTeams) > maxTeams {
		logInfo(ctx, "Limiting teams to maximum requested", LogFields{
			"component":     "github_client",
			"organization":  orgName,
			"fetched_teams": len(allTeams),
			"max_teams":     maxTeams,
			"limited":       true,
		})
		allTeams = allTeams[:maxTeams]
	}

	// Record team metrics
	metrics.recordCounter("teams_processed", len(allTeams), MetricLabels{
		"organization": orgName,
		"service":      "codeowners-scanner",
	})

	logInfo(ctx, "Successfully fetched GitHub teams", LogFields{
		"component":     "github_client",
		"operation":     "fetch_teams",
		"organization":  orgName,
		"max_teams":     maxTeams,
		"fetched_teams": len(allTeams),
		"total_pages":   page - 1,
	})

	return allTeams, nil
}

// fetchGitHubCodeownersWithService fetches CODEOWNERS file using GoFr HTTP service
func fetchGitHubCodeownersWithService(ctx *gofr.Context, owner, repo string) (GitHubCodeowners, error) {
	// Create span for tracking CODEOWNERS fetch
	span := createGitHubScanSpan(ctx, owner, "fetch_codeowners")
	defer finishSpan(span)

	// Start performance timer
	timer := startPerformanceTimer(ctx, "github_fetch_codeowners")
	defer stopPerformanceTimer(timer)

	// Initialize metrics collector
	metrics := newMetricsCollector(ctx, "codeowners-scanner")

	// Structured logging with context
	logInfo(ctx, "Starting CODEOWNERS fetch", LogFields{
		"component":  "github_client",
		"operation":  "fetch_codeowners",
		"owner":      owner,
		"repository": repo,
		"service":    "github",
	})

	if owner == "" {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("owner is required but was empty"),
			Operation:   "validate_params",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "request_failed",
			Context: map[string]interface{}{
				"owner":      owner,
				"repository": repo,
				"validation": "missing_param",
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "validation_error")
		return GitHubCodeowners{}, &gofrhttp.ErrorMissingParam{
			Params: []string{"owner"},
		}
	}

	if repo == "" {
		errCtx := ErrorContext{
			Error:       fmt.Errorf("repository is required but was empty"),
			Operation:   "validate_params",
			Component:   "github_client",
			Severity:    "error",
			Recoverable: true,
			UserImpact:  "request_failed",
			Context: map[string]interface{}{
				"owner":      owner,
				"repository": repo,
				"validation": "missing_param",
			},
		}
		logErrorWithStackTrace(ctx, errCtx)
		metrics.recordErrorCount("github_client", "validation_error")
		return GitHubCodeowners{}, &gofrhttp.ErrorMissingParam{
			Params: []string{"repository"},
		}
	}

	githubSvc := ctx.GetHTTPService("github")

	// Try different CODEOWNERS locations
	locations := []string{
		fmt.Sprintf("repos/%s/%s/contents/CODEOWNERS", owner, repo),
		fmt.Sprintf("repos/%s/%s/contents/.github/CODEOWNERS", owner, repo),
		fmt.Sprintf("repos/%s/%s/contents/docs/CODEOWNERS", owner, repo),
	}

	logInfo(ctx, "Searching for CODEOWNERS file in multiple locations", LogFields{
		"component":       "github_client",
		"operation":       "search_codeowners",
		"owner":           owner,
		"repository":      repo,
		"search_locations": len(locations),
	})

	for i, location := range locations {
		// Start location timer
		locationTimer := startPerformanceTimer(ctx, fmt.Sprintf("codeowners_location_%d", i+1))

		logDebug(ctx, "Checking CODEOWNERS location", LogFields{
			"component":  "github_client",
			"operation":  "check_location",
			"owner":      owner,
			"repository": repo,
			"location":   location,
			"attempt":    i + 1,
		})

		headers := buildGitHubRequestHeaders()
		resp, err := githubSvc.GetWithHeaders(ctx, location, nil, headers)
		if err != nil {
			stopPerformanceTimer(locationTimer)
			logDebug(ctx, "CODEOWNERS location request failed", LogFields{
				"component":  "github_client",
				"operation":  "check_location",
				"location":   location,
				"error":      err.Error(),
				"attempt":    i + 1,
			})
			metrics.recordErrorCount("github_client", "location_request_error")
			continue
		}
		defer resp.Body.Close()

		// Log rate limit information
		logRateLimitInfo(ctx, resp)

		// Record API call metrics
		metrics.recordAPICallCount("github", "codeowners", resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			logInfo(ctx, "CODEOWNERS file found", LogFields{
				"component":  "github_client",
				"operation":  "file_found",
				"owner":      owner,
				"repository": repo,
				"location":   location,
				"attempt":    i + 1,
			})

			var fileContent struct {
				Content string `json:"content"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
				stopPerformanceTimer(locationTimer)
				errCtx := ErrorContext{
					Error:       err,
					Operation:   "json_decode",
					Component:   "github_client",
					Severity:    "error",
					Recoverable: true,
					UserImpact:  "data_corruption",
					Context: map[string]interface{}{
						"owner":        owner,
						"repository":   repo,
						"location":     location,
						"content_type": resp.Header.Get("Content-Type"),
					},
				}
				logErrorWithStackTrace(ctx, errCtx)
				metrics.recordErrorCount("github_client", "decode_error")
				continue
			}

			stopPerformanceTimer(locationTimer)

			// Parse CODEOWNERS content
			rules := parseCodeownersContent(fileContent.Content)

			logInfo(ctx, "CODEOWNERS file parsed successfully", LogFields{
				"component":    "github_client",
				"operation":    "parse_success",
				"owner":        owner,
				"repository":   repo,
				"location":     location,
				"rules_count":  len(rules),
				"content_size": len(fileContent.Content),
			})

			// Record CODEOWNERS metrics
			metrics.recordGauge("codeowners_rules_count", float64(len(rules)), MetricLabels{
				"owner":      owner,
				"repository": repo,
				"service":    "codeowners-scanner",
			})

			return GitHubCodeowners{
				Repository: fmt.Sprintf("%s/%s", owner, repo),
				Rules:      rules,
				Errors:     []GitHubCodeownersError{},
			}, nil
		} else {
			stopPerformanceTimer(locationTimer)
			logDebug(ctx, "CODEOWNERS not found at location", LogFields{
				"component":   "github_client",
				"operation":   "location_not_found",
				"owner":       owner,
				"repository":  repo,
				"location":    location,
				"status_code": resp.StatusCode,
				"attempt":     i + 1,
			})
		}
	}

	// Return empty CODEOWNERS (not an error - many repos don't have CODEOWNERS)
	logInfo(ctx, "No CODEOWNERS file found in any location", LogFields{
		"component":         "github_client",
		"operation":         "codeowners_not_found",
		"owner":             owner,
		"repository":        repo,
		"locations_checked": len(locations),
		"result":            "empty_codeowners",
	})

	// Record metric for repositories without CODEOWNERS
	metrics.recordCounter("codeowners_not_found", 1, MetricLabels{
		"owner":      owner,
		"repository": repo,
		"service":    "codeowners-scanner",
	})

	return GitHubCodeowners{
		Repository: fmt.Sprintf("%s/%s", owner, repo),
		Rules:      []GitHubCodeownersRule{},
		Errors:     []GitHubCodeownersError{},
	}, nil
}

// parseCodeownersContent parses base64-encoded CODEOWNERS content (Pure Core)
func parseCodeownersContent(base64Content string) []GitHubCodeownersRule {
	validateBase64ContentNotEmpty(base64Content)

	// Decode base64 content
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return []GitHubCodeownersRule{}
	}

	content := string(decodedBytes)
	if content == "" {
		return []GitHubCodeownersRule{}
	}

	// Split into lines and filter out comments and empty lines
	lines := lo.Filter(lo.Map(strings.Split(content, "\n"), func(line string, _ int) string {
		return strings.TrimSpace(line)
	}), func(line string, _ int) bool {
		return line != "" && !strings.HasPrefix(line, "#")
	})

	// Parse each line into a rule
	rules := lo.Map(lines, func(line string, index int) GitHubCodeownersRule {
		parts := lo.Filter(strings.Fields(line), func(part string, _ int) bool {
			return part != ""
		})

		if len(parts) < 2 {
			return GitHubCodeownersRule{
				Pattern: line,
				Owners:  []string{},
				Line:    index + 1,
			}
		}

		return GitHubCodeownersRule{
			Pattern: parts[0],
			Owners:  parts[1:],
			Line:    index + 1,
		}
	})

	return rules
}

// validateBase64ContentNotEmpty validates base64 content is not empty (Pure Core)
func validateBase64ContentNotEmpty(content string) {
	if content == "" {
		panic("Base64 content cannot be empty")
	}
}

// Rate limit monitoring and logging utilities

// logRateLimitInfo logs GitHub API rate limit information from response headers
func logRateLimitInfo(ctx *gofr.Context, resp *http.Response) {
	if resp == nil {
		return
	}

	// Extract rate limit headers
	rateLimit := resp.Header.Get("X-RateLimit-Limit")
	rateRemaining := resp.Header.Get("X-RateLimit-Remaining")
	rateReset := resp.Header.Get("X-RateLimit-Reset")
	rateUsed := resp.Header.Get("X-RateLimit-Used")
	rateResource := resp.Header.Get("X-RateLimit-Resource")

	if rateLimit != "" || rateRemaining != "" {
		// Parse reset time
		resetTime := ""
		if rateReset != "" {
			if resetTimestamp, err := strconv.ParseInt(rateReset, 10, 64); err == nil {
				resetTime = time.Unix(resetTimestamp, 0).UTC().Format(time.RFC3339)
			}
		}

		// Calculate remaining percentage
		remainingPct := 0.0
		if rateLimit != "" && rateRemaining != "" {
			if limit, err := strconv.Atoi(rateLimit); err == nil && limit > 0 {
				if remaining, err := strconv.Atoi(rateRemaining); err == nil {
					remainingPct = float64(remaining) / float64(limit) * 100
				}
			}
		}

		// Determine log level based on remaining rate limit
		logLevel := "debug"
		if remainingPct < 10 {
			logLevel = "warn"
		} else if remainingPct < 25 {
			logLevel = "info"
		}

		logWithContext(ctx, logLevel, "GitHub API rate limit status", LogFields{
			"component":         "github_client",
			"operation":         "rate_limit_check",
			"rate_limit":        rateLimit,
			"rate_remaining":    rateRemaining,
			"rate_used":         rateUsed,
			"rate_reset":        rateReset,
			"rate_reset_time":   resetTime,
			"rate_resource":     rateResource,
			"remaining_percent": fmt.Sprintf("%.1f%%", remainingPct),
			"status_code":       resp.StatusCode,
		})

		// Record rate limit metrics
		metrics := newMetricsCollector(ctx, "codeowners-scanner")
		if rateRemaining != "" {
			if remaining, err := strconv.Atoi(rateRemaining); err == nil {
				metrics.recordGauge("github_rate_limit_remaining", float64(remaining), MetricLabels{
					"resource": rateResource,
					"service":  "codeowners-scanner",
				})
			}
		}

		if rateLimit != "" {
			if limit, err := strconv.Atoi(rateLimit); err == nil {
				metrics.recordGauge("github_rate_limit_total", float64(limit), MetricLabels{
					"resource": rateResource,
					"service":  "codeowners-scanner",
				})
			}
		}

		// Log warning if rate limit is low
		if remainingPct < 10 && remainingPct > 0 {
			logWarn(ctx, "GitHub API rate limit critically low", LogFields{
				"component":         "github_client",
				"rate_remaining":    rateRemaining,
				"remaining_percent": fmt.Sprintf("%.1f%%", remainingPct),
				"rate_reset_time":   resetTime,
				"recommendation":    "consider_throttling",
			})
		}
	}
}

// extractResponseHeaders extracts important response headers for logging
func extractResponseHeaders(resp *http.Response) map[string]string {
	if resp == nil {
		return map[string]string{}
	}

	headers := make(map[string]string)
	importantHeaders := []string{
		"Content-Type",
		"Content-Length",
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
		"X-RateLimit-Resource",
		"X-GitHub-Request-Id",
		"X-GitHub-Media-Type",
		"ETag",
		"Last-Modified",
	}

	for _, header := range importantHeaders {
		if value := resp.Header.Get(header); value != "" {
			headers[strings.ToLower(strings.ReplaceAll(header, "-", "_"))] = value
		}
	}

	return headers
}

// buildGitHubRequestHeaders builds headers for GitHub API requests (Pure Core)
func buildGitHubRequestHeaders() map[string]string {
	// Note: In a real implementation, we would get the token from configuration
	// For now, we'll use a placeholder that expects GITHUB_TOKEN environment variable
	token := os.Getenv("GITHUB_TOKEN")

	headers := map[string]string{
		"Accept":     "application/vnd.github.v3+json",
		"User-Agent": "overseer-codeowners-scanner/1.0",
	}

	if token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", token)
	}

	return headers
}
