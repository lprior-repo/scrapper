package main

import (
	"fmt"
	"time"

	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"
)

// NewAppHandler creates a new app handler with dependencies
func NewAppHandler(deps *AppDependencies) *AppHandler {
	return &AppHandler{deps: deps}
}

// handleScanOrganization handles organization scanning
func (h *AppHandler) handleScanOrganization(ctx *gofr.Context) (interface{}, error) {
	orgName := extractOrgParam(ctx)
	if orgName == "" {
		return nil, createMissingParamError("org")
	}

	scanRequest := buildScanRequest(ctx, h.deps.Config, orgName)
	response, err := scanOrganization(ctx, h.deps, scanRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// handleGetGraph handles graph data retrieval
func (h *AppHandler) handleGetGraph(ctx *gofr.Context) (interface{}, error) {
	orgName := extractOrgParam(ctx)
	if orgName == "" {
		return nil, createMissingParamError("org")
	}

	useTopics := parseBoolFromQuery(ctx, "useTopics", false)
	response, err := getOrganizationGraph(ctx, h.deps, orgName, useTopics)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// handleGetStats handles statistics retrieval
func (h *AppHandler) handleGetStats(ctx *gofr.Context) (interface{}, error) {
	orgName := extractOrgParam(ctx)
	if orgName == "" {
		return nil, createMissingParamError("org")
	}

	response, err := getOrganizationStats(ctx, h.deps, orgName)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// handleHealth handles health check
func (h *AppHandler) handleHealth(ctx *gofr.Context) (interface{}, error) {
	if err := checkNeo4jHealth(ctx, h.deps.Neo4jConn); err != nil {
		return nil, fmt.Errorf("database health check failed: %w", err)
	}

	return buildHealthResponse(), nil
}

// handleOpenAPI serves the OpenAPI documentation UI
func (*AppHandler) handleOpenAPI(_ *gofr.Context) (interface{}, error) {
	html := buildOpenAPIHTML()
	return response.Raw{Data: html}, nil
}

// handleOpenAPISpec serves the OpenAPI specification file
func (*AppHandler) handleOpenAPISpec(_ *gofr.Context) (interface{}, error) {
	specContent := buildOpenAPISpec()
	return response.Raw{Data: specContent}, nil
}

// extractOrgParam extracts organization parameter from path
func extractOrgParam(ctx *gofr.Context) string {
	return ctx.PathParam("org")
}

// createMissingParamError creates missing parameter error
func createMissingParamError(param string) error {
	return &gofrhttp.ErrorMissingParam{
		Params: []string{param},
	}
}

// buildScanRequest constructs scan request from context and config
func buildScanRequest(ctx *gofr.Context, config AppConfig, orgName string) ScanRequest {
	maxRepos := parseIntFromQuery(ctx, "max_repos", 100)
	maxTeams := parseIntFromQuery(ctx, "max_teams", 50)
	useTopics := parseBoolFromQuery(ctx, "use_topics", config.GitHub.UseTopics)

	return ScanRequest{
		Organization: orgName,
		MaxRepos:     maxRepos,
		MaxTeams:     maxTeams,
		UseTopics:    useTopics,
	}
}

// buildHealthResponse constructs health check response
func buildHealthResponse() map[string]interface{} {
	return map[string]interface{}{
		"status":    "healthy",
		"database":  "connected",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	}
}