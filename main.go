package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/http/response"
)

// ScanRequest represents a request to scan a GitHub organization
type ScanRequest struct {
	Organization string `json:"organization"`
	MaxRepos     int    `json:"max_repos"`
	MaxTeams     int    `json:"max_teams"`
}

// ScanResponse represents the response from scanning an organization
type ScanResponse struct {
	Success      bool                   `json:"success"`
	Organization string                 `json:"organization"`
	Summary      ScanSummary            `json:"summary"`
	Errors       []string               `json:"errors"`
	Data         map[string]interface{} `json:"data"`
}

// ScanSummary represents scan statistics
type ScanSummary struct {
	TotalRepos          int      `json:"total_repos"`
	ReposWithCodeowners int      `json:"repos_with_codeowners"`
	TotalTeams          int      `json:"total_teams"`
	UniqueOwners        []string `json:"unique_owners"`
	APICallsUsed        int      `json:"api_calls_used"`
	ProcessingTimeMs    int64    `json:"processing_time_ms"`
}

// GraphResponse represents graph visualization data
type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode represents a node in the graph
type GraphNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Label    string                 `json:"label"`
	Data     map[string]interface{} `json:"data"`
	Position GraphPosition          `json:"position"`
}

// GraphEdge represents an edge in the graph
type GraphEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Label  string `json:"label"`
}

// GraphPosition represents node position in the graph
type GraphPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// StatsResponse represents organization statistics
type StatsResponse struct {
	Organization      string `json:"organization"`
	TotalRepositories int    `json:"total_repositories"`
	TotalTeams        int    `json:"total_teams"`
	TotalUsers        int    `json:"total_users"`
	TotalCodeowners   int    `json:"total_codeowners"`
	CodeownerCoverage string `json:"codeowner_coverage"`
	LastScanTime      string `json:"last_scan_time"`
}

// AppDependencies represents application dependencies
type AppDependencies struct {
	Config    AppConfig
	Neo4jConn *Neo4jConnection
}

// createAppDependencies creates application dependencies (Orchestrator)
func createAppDependencies(ctx context.Context) (*AppDependencies, error) {
	// Load configuration using environment variables
	config := loadConfigFromEnv()

	// Validate configuration
	if err := validateConfiguration(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	neo4jConn, err := createNeo4jConnection(ctx, config.Neo4j)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j connection: %w", err)
	}

	if err := checkNeo4jHealth(ctx, neo4jConn); err != nil {
		return nil, fmt.Errorf("Neo4j health check failed: %w", err)
	}

	if err := createNeo4jConstraints(ctx, neo4jConn); err != nil {
		return nil, fmt.Errorf("failed to create Neo4j constraints: %w", err)
	}

	if err := createNeo4jIndexes(ctx, neo4jConn); err != nil {
		return nil, fmt.Errorf("failed to create Neo4j indexes: %w", err)
	}

	return &AppDependencies{
		Config:    config,
		Neo4jConn: neo4jConn,
	}, nil
}

// cleanupAppDependencies cleans up application dependencies (Orchestrator)
func cleanupAppDependencies(ctx context.Context, deps *AppDependencies) error {
	if deps == nil {
		return nil
	}

	if deps.Neo4jConn != nil {
		if err := closeNeo4jConnection(ctx, deps.Neo4jConn); err != nil {
			return fmt.Errorf("failed to close Neo4j connection: %w", err)
		}
	}

	return nil
}

// handleScanOrganization handles organization scanning (Orchestrator)
func (h *AppHandler) handleScanOrganization(ctx *gofr.Context) (interface{}, error) {
	orgName := ctx.PathParam("org")
	if orgName == "" {
		return nil, &gofrhttp.ErrorMissingParam{
			Params: []string{"org"},
		}
	}

	maxRepos := parseIntFromQuery(ctx, "max_repos", 100)
	maxTeams := parseIntFromQuery(ctx, "max_teams", 50)

	scanRequest := ScanRequest{
		Organization: orgName,
		MaxRepos:     maxRepos,
		MaxTeams:     maxTeams,
	}

	response, err := scanOrganization(ctx, h.deps, scanRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to scan organization: %w", err)
	}

	return response, nil
}

// handleGetGraph handles graph data retrieval (Orchestrator)
func (h *AppHandler) handleGetGraph(ctx *gofr.Context) (interface{}, error) {
	orgName := ctx.PathParam("org")
	if orgName == "" {
		return nil, &gofrhttp.ErrorMissingParam{
			Params: []string{"org"},
		}
	}

	response, err := getOrganizationGraph(ctx, h.deps, orgName)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization graph: %w", err)
	}

	return response, nil
}

// handleGetStats handles statistics retrieval (Orchestrator)
func (h *AppHandler) handleGetStats(ctx *gofr.Context) (interface{}, error) {
	orgName := ctx.PathParam("org")
	if orgName == "" {
		return nil, &gofrhttp.ErrorMissingParam{
			Params: []string{"org"},
		}
	}

	response, err := getOrganizationStats(ctx, h.deps, orgName)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization stats: %w", err)
	}

	return response, nil
}

// handleHealth handles health check (Orchestrator)
func (h *AppHandler) handleHealth(ctx *gofr.Context) (interface{}, error) {
	if err := checkNeo4jHealth(ctx, h.deps.Neo4jConn); err != nil {
		return nil, fmt.Errorf("database health check failed: %w", err)
	}

	return map[string]interface{}{
		"status":    "healthy",
		"database":  "connected",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil
}

// handleOpenAPI serves the OpenAPI documentation UI (Orchestrator)
func (h *AppHandler) handleOpenAPI(ctx *gofr.Context) (interface{}, error) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>GitHub Codeowners API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api/openapi.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	return response.Raw{Data: html}, nil
}

// handleOpenAPISpec serves the OpenAPI specification file (Orchestrator)
func (h *AppHandler) handleOpenAPISpec(ctx *gofr.Context) (interface{}, error) {
	specContent := `openapi: 3.0.3
info:
  title: GitHub Codeowners Visualization API
  description: API for scanning GitHub organizations and visualizing codeowners relationships
  version: 1.0.0
  contact:
    name: API Support
    url: https://github.com/your-org/scrapper
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: http://localhost:8080
    description: Development server

paths:
  /api/health:
    get:
      summary: Health check endpoint
      description: Returns the health status of the API and its dependencies
      operationId: healthCheck
      tags:
        - System
      responses:
        '200':
          description: System is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: object
                    properties:
                      status:
                        type: string
                        example: "healthy"
                      database:
                        type: string
                        example: "connected"
                      version:
                        type: string
                        example: "1.0.0"
                      timestamp:
                        type: string
                        format: date-time
                        example: "2025-07-17T21:08:23-05:00"

  /api/scan/{org}:
    post:
      summary: Scan GitHub organization
      description: Scans a GitHub organization for repositories, teams, and CODEOWNERS files
      operationId: scanOrganization
      tags:
        - Scanning
      parameters:
        - name: org
          in: path
          required: true
          description: GitHub organization name
          schema:
            type: string
            example: "microsoft"
        - name: max_repos
          in: query
          required: false
          description: Maximum number of repositories to scan
          schema:
            type: integer
            default: 100
            minimum: 1
            maximum: 1000
        - name: max_teams
          in: query
          required: false
          description: Maximum number of teams to scan
          schema:
            type: integer
            default: 50
            minimum: 1
            maximum: 500

components:
  schemas:
    ScanResponse:
      type: object
      properties:
        success:
          type: boolean
          description: Whether the scan was successful
        organization:
          type: string
          description: Name of the scanned organization
        summary:
          type: object
          properties:
            total_repos:
              type: integer
            repos_with_codeowners:
              type: integer
            total_teams:
              type: integer
            unique_owners:
              type: array
              items:
                type: string

tags:
  - name: System
    description: System health and status operations
  - name: Scanning
    description: GitHub organization scanning operations`

	return response.Raw{Data: specContent}, nil
}

// scanOrganization scans a GitHub organization (Orchestrator)
func scanOrganization(ctx *gofr.Context, deps *AppDependencies, request ScanRequest) (ScanResponse, error) {
	startTime := time.Now()

	// Fetch organization data
	org, err := fetchGitHubOrganizationWithService(ctx, request.Organization)
	if err != nil {
		return ScanResponse{}, fmt.Errorf("failed to fetch organization: %w", err)
	}

	// Fetch repositories
	repos, err := fetchGitHubRepositoriesWithService(ctx, request.Organization, request.MaxRepos)
	if err != nil {
		return ScanResponse{}, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	// Fetch teams
	teams, err := fetchGitHubTeamsWithService(ctx, request.Organization, request.MaxTeams)
	if err != nil {
		return ScanResponse{}, fmt.Errorf("failed to fetch teams: %w", err)
	}

	// Fetch CODEOWNERS files
	codeowners, err := fetchCodeownersForReposWithService(ctx, repos)
	if err != nil {
		return ScanResponse{}, fmt.Errorf("failed to fetch CODEOWNERS: %w", err)
	}

	// Store data in Neo4j
	if err := storeOrganizationData(ctx, deps.Neo4jConn, org, repos, teams, codeowners); err != nil {
		return ScanResponse{}, fmt.Errorf("failed to store data: %w", err)
	}

	// Calculate summary
	summary := calculateScanSummary(repos, codeowners, teams, time.Since(startTime))

	return ScanResponse{
		Success:      true,
		Organization: request.Organization,
		Summary:      summary,
		Errors:       []string{},
		Data: map[string]interface{}{
			"organization": org,
			"repositories": repos,
			"teams":        teams,
			"codeowners":   codeowners,
		},
	}, nil
}

// getOrganizationGraph retrieves graph data for an organization (Orchestrator)
func getOrganizationGraph(ctx *gofr.Context, deps *AppDependencies, orgName string) (GraphResponse, error) {
	session, err := createNeo4jSession(ctx, deps.Neo4jConn)
	if err != nil {
		return GraphResponse{}, fmt.Errorf("failed to create Neo4j session: %w", err)
	}
	defer closeNeo4jSession(ctx, session)

	// Fetch nodes
	nodesQuery := buildGraphNodesQuery(orgName)
	nodesResult, err := executeNeo4jReadQuery(ctx, session, nodesQuery, map[string]interface{}{
		"orgName": orgName,
	})
	if err != nil {
		return GraphResponse{}, fmt.Errorf("failed to fetch graph nodes: %w", err)
	}

	// Fetch edges
	edgesQuery := buildGraphEdgesQuery(orgName)
	edgesResult, err := executeNeo4jReadQuery(ctx, session, edgesQuery, map[string]interface{}{
		"orgName": orgName,
	})
	if err != nil {
		return GraphResponse{}, fmt.Errorf("failed to fetch graph edges: %w", err)
	}

	nodes := convertToGraphNodes(nodesResult.Records)
	edges := convertToGraphEdges(edgesResult.Records)

	return GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// getOrganizationStats retrieves statistics for an organization (Orchestrator)
func getOrganizationStats(ctx *gofr.Context, deps *AppDependencies, orgName string) (StatsResponse, error) {
	session, err := createNeo4jSession(ctx, deps.Neo4jConn)
	if err != nil {
		return StatsResponse{}, fmt.Errorf("failed to create Neo4j session: %w", err)
	}
	defer closeNeo4jSession(ctx, session)

	query := buildStatsQuery(orgName)
	result, err := executeNeo4jReadQuery(ctx, session, query, map[string]interface{}{
		"orgName": orgName,
	})
	if err != nil {
		return StatsResponse{}, fmt.Errorf("failed to fetch stats: %w", err)
	}

	if len(result.Records) == 0 {
		return StatsResponse{}, fmt.Errorf("no data found for organization: %s", orgName)
	}

	return convertToStatsResponse(result.Records[0], orgName), nil
}

// fetchCodeownersForReposWithService fetches CODEOWNERS files for repositories (Orchestrator)
func fetchCodeownersForReposWithService(ctx *gofr.Context, repos []GitHubRepository) ([]GitHubCodeowners, error) {
	codeowners := make([]GitHubCodeowners, 0, len(repos))

	for _, repo := range repos {
		owner, name := parseRepositoryFullName(repo.FullName)
		if owner == "" || name == "" {
			continue
		}

		codeowner, err := fetchGitHubCodeownersWithService(ctx, owner, name)
		if err != nil {
			// Continue on error - not all repos have CODEOWNERS
			continue
		}

		if len(codeowner.Rules) > 0 {
			codeowners = append(codeowners, codeowner)
		}
	}

	return codeowners, nil
}

// storeOrganizationData stores organization data in Neo4j (Orchestrator)
func storeOrganizationData(ctx *gofr.Context, conn *Neo4jConnection, org GitHubOrganization, repos []GitHubRepository, teams []GitHubTeam, codeowners []GitHubCodeowners) error {
	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j session: %w", err)
	}
	defer closeNeo4jSession(ctx, session)

	// Store organization
	if err := storeOrganization(ctx, session, org); err != nil {
		return fmt.Errorf("failed to store organization: %w", err)
	}

	// Store repositories
	for _, repo := range repos {
		if err := storeRepository(ctx, session, repo, org.Login); err != nil {
			return fmt.Errorf("failed to store repository %s: %w", repo.Name, err)
		}
	}

	// Store teams
	for _, team := range teams {
		if err := storeTeam(ctx, session, team, org.Login); err != nil {
			return fmt.Errorf("failed to store team %s: %w", team.Name, err)
		}
	}

	// Store CODEOWNERS
	for _, codeowner := range codeowners {
		if err := storeCodeowners(ctx, session, codeowner, org.Login); err != nil {
			return fmt.Errorf("failed to store CODEOWNERS for %s: %w", codeowner.Repository, err)
		}
	}

	return nil
}

// parseIntFromQuery extracts integer from query parameters (Pure Core)
func parseIntFromQuery(ctx *gofr.Context, key string, defaultValue int) int {
	value := ctx.Param(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// parseRepositoryFullName splits repository full name into owner and name (Pure Core)
func parseRepositoryFullName(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// calculateScanSummary calculates summary statistics from scan results (Pure Core)
func calculateScanSummary(repos []GitHubRepository, codeowners []GitHubCodeowners, teams []GitHubTeam, duration time.Duration) ScanSummary {
	uniqueOwners := make(map[string]bool)

	for _, codeowner := range codeowners {
		for _, rule := range codeowner.Rules {
			for _, owner := range rule.Owners {
				uniqueOwners[owner] = true
			}
		}
	}

	ownersList := lo.Keys(uniqueOwners)

	return ScanSummary{
		TotalRepos:          len(repos),
		ReposWithCodeowners: len(codeowners),
		TotalTeams:          len(teams),
		UniqueOwners:        ownersList,
		APICallsUsed:        len(repos) + len(teams) + len(codeowners) + 1, // Estimated
		ProcessingTimeMs:    duration.Milliseconds(),
	}
}

// AppHandler contains the application dependencies
type AppHandler struct {
	deps *AppDependencies
}

// loadConfigFromEnv loads configuration from environment variables
func loadConfigFromEnv() AppConfig {
	return AppConfig{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Port:        getIntEnvOrDefault("HTTP_PORT", 8081),
		GitHub: GitHubConfig{
			Token:        os.Getenv("GITHUB_TOKEN"),
			BaseURL:      getEnvOrDefault("GITHUB_BASE_URL", "https://api.github.com"),
			UserAgent:    getEnvOrDefault("GITHUB_USER_AGENT", "overseer-codeowners-scanner/1.0"),
			Timeout:      getDurationEnvOrDefault("GITHUB_TIMEOUT", 30*time.Second),
			MaxRetries:   getIntEnvOrDefault("GITHUB_MAX_RETRIES", 3),
			RateLimitMin: getIntEnvOrDefault("GITHUB_RATE_LIMIT_MIN", 100),
		},
		Neo4j: Neo4jConfig{
			URI:      getEnvOrDefault("NEO4J_URI", "bolt://localhost:7687"),
			Username: getEnvOrDefault("NEO4J_USERNAME", "neo4j"),
			Password: getEnvOrDefault("NEO4J_PASSWORD", "password"),
			Database: getEnvOrDefault("NEO4J_DATABASE", "neo4j"),
			Timeout:  getDurationEnvOrDefault("NEO4J_TIMEOUT", 30*time.Second),
		},
		Server: ServerConfig{
			ReadTimeout:    getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:   getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:    getDurationEnvOrDefault("SERVER_IDLE_TIMEOUT", 60*time.Second),
			MaxHeaderBytes: getIntEnvOrDefault("SERVER_MAX_HEADER_BYTES", 1<<20),
		},
	}
}

// getEnvOrDefault gets environment variable or returns default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnvOrDefault gets int environment variable or returns default
func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getDurationEnvOrDefault gets duration environment variable or returns default
func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// validateConfiguration validates the loaded configuration
func validateConfiguration(config AppConfig) error {
	validationErrors := validateAppConfig(config)
	if len(validationErrors) > 0 {
		return fmt.Errorf("configuration validation failed: %d errors found", len(validationErrors))
	}
	return nil
}

// NewAppHandler creates a new app handler with dependencies
func NewAppHandler(deps *AppDependencies) *AppHandler {
	return &AppHandler{deps: deps}
}

func main() {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--cleanup", "cleanup":
			emergencyCleanup()
			return
		case "api":
			// Continue with API server startup
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: api, --cleanup, cleanup")
			return
		}
	}

	// Create GoFr app
	app := gofr.New()

	// Create application dependencies using GoFr context
	ctx := context.Background()
	deps, err := createAppDependencies(ctx)
	if err != nil {
		app.Logger().Fatalf("Failed to create app dependencies: %v", err)
	}

	// Register GitHub as an HTTP service
	RegisterGitHubService(app, GitHubServiceConfig{
		Token:        deps.Config.GitHub.Token,
		BaseURL:      deps.Config.GitHub.BaseURL,
		UserAgent:    deps.Config.GitHub.UserAgent,
		Timeout:      deps.Config.GitHub.Timeout,
		MaxRetries:   deps.Config.GitHub.MaxRetries,
		RateLimitMin: deps.Config.GitHub.RateLimitMin,
	})

	// Create handler with dependencies
	handler := NewAppHandler(deps)

	// Set up graceful shutdown
	defer func() {
		if err := cleanupAppDependencies(ctx, deps); err != nil {
			app.Logger().Errorf("Failed to cleanup dependencies: %v", err)
		}
	}()

	// API routes
	app.POST("/api/scan/{org}", handler.handleScanOrganization)
	app.GET("/api/graph/{org}", handler.handleGetGraph)
	app.GET("/api/stats/{org}", handler.handleGetStats)
	app.GET("/api/health", handler.handleHealth)
	
	// OpenAPI documentation
	app.GET("/api/docs", handler.handleOpenAPI)
	app.GET("/api/openapi.yaml", handler.handleOpenAPISpec)

	// Start server
	app.Logger().Infof("Starting GitHub Codeowners Visualization API")
	app.Run()
}
