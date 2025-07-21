package main

import (
	"context"
	"fmt"
	"time"

	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
)

// createAppDependencies creates application dependencies
func createAppDependencies(ctx context.Context) (*AppDependencies, error) {
	config, err := loadAndValidateConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration setup failed: %w", err)
	}

	neo4jConn, err := setupNeo4jConnection(ctx, config.Neo4j)
	if err != nil {
		return nil, fmt.Errorf("Neo4j setup failed: %w", err)
	}

	return &AppDependencies{
		Config:    config,
		Neo4jConn: neo4jConn,
	}, nil
}

// loadAndValidateConfig loads and validates the application configuration
func loadAndValidateConfig() (AppConfig, error) {
	config := loadConfigFromEnv()

	if err := validateConfiguration(config); err != nil {
		return AppConfig{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// setupNeo4jConnection creates and initializes Neo4j connection
func setupNeo4jConnection(ctx context.Context, config Neo4jConfig) (*Neo4jConnection, error) {
	neo4jConn, err := createNeo4jConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j connection: %w", err)
	}

	if err := performNeo4jHealthCheck(ctx, neo4jConn); err != nil {
		return nil, fmt.Errorf("Neo4j health check failed: %w", err)
	}

	if err := initializeNeo4jSchema(ctx, neo4jConn); err != nil {
		return nil, fmt.Errorf("failed to initialize Neo4j schema: %w", err)
	}

	return neo4jConn, nil
}

// cleanupAppDependencies cleans up application dependencies
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

// scanOrganization scans a GitHub organization
func scanOrganization(ctx *gofr.Context, deps *AppDependencies, request ScanRequest) (ScanResponse, error) {
	startTime := time.Now()

	org, err := fetchGitHubOrganizationWithService(ctx, request.Organization)
	if err != nil {
		return ScanResponse{}, err
	}

	repos, err := fetchGitHubRepositoriesWithService(ctx, request.Organization, request.MaxRepos)
	if err != nil {
		return ScanResponse{}, err
	}

	teams, topics, err := fetchTeamsOrTopics(ctx, request, repos)
	if err != nil {
		return ScanResponse{}, err
	}

	codeowners, err := fetchCodeownersForReposWithService(ctx, repos)
	if err != nil {
		return ScanResponse{}, err
	}

	if err := storeOrganizationData(ctx, deps.Neo4jConn, org, repos, teams, topics, codeowners); err != nil {
		return ScanResponse{}, convertNeo4jErrorToGoFr(err)
	}

	summary := calculateScanSummary(repos, codeowners, teams, topics, time.Since(startTime))

	return buildScanResponse(request.Organization, summary, org, repos, teams, topics, codeowners), nil
}

// getOrganizationGraph retrieves graph data for an organization
func getOrganizationGraph(ctx *gofr.Context, deps *AppDependencies, orgName string, useTopics bool) (GraphResponse, error) {
	session, err := createNeo4jSession(ctx, deps.Neo4jConn)
	if err != nil {
		return GraphResponse{}, convertNeo4jErrorToGoFr(err)
	}
	defer closeNeo4jSession(ctx, session)

	nodes, err := fetchGraphNodes(ctx, session, orgName, useTopics)
	if err != nil {
		return GraphResponse{}, convertNeo4jErrorToGoFr(err)
	}

	edges, err := fetchGraphEdges(ctx, session, orgName, useTopics)
	if err != nil {
		return GraphResponse{}, convertNeo4jErrorToGoFr(err)
	}

	return GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// getOrganizationStats retrieves statistics for an organization
func getOrganizationStats(ctx *gofr.Context, deps *AppDependencies, orgName string) (StatsResponse, error) {
	session, err := createNeo4jSession(ctx, deps.Neo4jConn)
	if err != nil {
		return StatsResponse{}, convertNeo4jErrorToGoFr(err)
	}
	defer closeNeo4jSession(ctx, session)

	query := buildStatsQuery(orgName)
	result, err := executeNeo4jReadQuery(ctx, session, query, map[string]interface{}{
		"orgName": orgName,
	})
	if err != nil {
		return StatsResponse{}, convertNeo4jErrorToGoFr(err)
	}

	if len(result.Records) == 0 {
		return StatsResponse{}, &gofrhttp.ErrorEntityNotFound{
			Name:  "organization",
			Value: orgName,
		}
	}

	return convertToStatsResponse(result.Records[0], orgName), nil
}

// fetchCodeownersForReposWithService fetches CODEOWNERS files for repositories
func fetchCodeownersForReposWithService(ctx *gofr.Context, repos []GitHubRepository) ([]GitHubCodeowners, error) {
	codeowners := make([]GitHubCodeowners, 0, len(repos))

	for _, repo := range repos {
		codeowner := fetchCodeownersForSingleRepo(ctx, repo)
		if codeowner != nil && len(codeowner.Rules) > 0 {
			codeowners = append(codeowners, *codeowner)
		}
	}

	return codeowners, nil
}

// storeOrganizationData stores organization data in Neo4j
func storeOrganizationData(ctx *gofr.Context, conn *Neo4jConnection, org GitHubOrganization, repos []GitHubRepository, teams []GitHubTeam, topics []GitHubTopic, codeowners []GitHubCodeowners) error {
	session, err := createNeo4jSession(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j session: %w", err)
	}
	defer closeNeo4jSession(ctx, session)

	if err := storeOrganization(ctx, session, org); err != nil {
		return fmt.Errorf("failed to store organization: %w", err)
	}

	if err := storeRepositories(ctx, session, repos, org.Login); err != nil {
		return fmt.Errorf("failed to store repositories: %w", err)
	}

	if err := storeTeamsAndTopics(ctx, session, teams, topics, org.Login); err != nil {
		return fmt.Errorf("failed to store teams and topics: %w", err)
	}

	if err := storeCodeownersData(ctx, session, codeowners, org.Login); err != nil {
		return fmt.Errorf("failed to store codeowners: %w", err)
	}

	return nil
}