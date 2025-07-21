package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
)

// parseIntFromQuery extracts integer from query parameters
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

// parseBoolFromQuery extracts boolean from query parameters
func parseBoolFromQuery(ctx *gofr.Context, key string, defaultValue bool) bool {
	value := ctx.Param(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// parseRepositoryFullName splits repository full name into owner and name
func parseRepositoryFullName(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// calculateScanSummary calculates summary statistics from scan results
func calculateScanSummary(repos []GitHubRepository, codeowners []GitHubCodeowners, teams []GitHubTeam, topics []GitHubTopic, duration time.Duration) ScanSummary {
	uniqueOwners := extractUniqueOwners(codeowners)
	ownersList := lo.Keys(uniqueOwners)

	return ScanSummary{
		TotalRepos:          len(repos),
		ReposWithCodeowners: len(codeowners),
		TotalTeams:          len(teams),
		TotalTopics:         len(topics),
		UniqueOwners:        ownersList,
		APICallsUsed:        estimateAPICallsUsed(repos, teams, codeowners),
		ProcessingTimeMs:    duration.Milliseconds(),
	}
}

// fetchCodeownersForSingleRepo fetches CODEOWNERS for a single repository
func fetchCodeownersForSingleRepo(ctx *gofr.Context, repo GitHubRepository) *GitHubCodeowners {
	owner, name := parseRepositoryFullName(repo.FullName)
	if owner == "" || name == "" {
		return nil
	}

	codeowner, err := fetchGitHubCodeownersWithService(ctx, owner, name)
	if err != nil {
		return nil
	}

	return &codeowner
}

// convertNeo4jErrorToGoFr converts Neo4j errors to appropriate GoFr error types
func convertNeo4jErrorToGoFr(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already a GoFr error
	switch err.(type) {
	case *gofrhttp.ErrorMissingParam,
		*gofrhttp.ErrorEntityNotFound,
		*gofrhttp.ErrorInvalidParam,
		*gofrhttp.ErrorRequestTimeout:
		return err
	}

	return convertNeo4jErrorByMessage(err)
}

// fetchTeamsOrTopics fetches teams or topics based on configuration
func fetchTeamsOrTopics(ctx *gofr.Context, request ScanRequest, repos []GitHubRepository) ([]GitHubTeam, []GitHubTopic, error) {
	var teams []GitHubTeam
	var topics []GitHubTopic

	if request.UseTopics {
		topics = collectTopicsFromRepositories(repos)
		ctx.Logger.Infof("Collected %d unique topics from repositories", len(topics))
	} else {
		teamsResult, err := fetchGitHubTeamsWithService(ctx, request.Organization, request.MaxTeams)
		if err != nil {
			ctx.Logger.Warnf("Failed to fetch teams for organization %s (likely due to permissions): %v", request.Organization, err)
			teams = []GitHubTeam{}
		} else {
			teams = teamsResult
		}
	}

	return teams, topics, nil
}

// performNeo4jHealthCheck performs health check on Neo4j connection
func performNeo4jHealthCheck(ctx context.Context, neo4jConn *Neo4jConnection) error {
	if err := checkNeo4jHealth(ctx, neo4jConn); err != nil {
		return err
	}
	return nil
}

// initializeNeo4jSchema initializes Neo4j database schema
func initializeNeo4jSchema(ctx context.Context, neo4jConn *Neo4jConnection) error {
	if err := createNeo4jConstraints(ctx, neo4jConn); err != nil {
		return fmt.Errorf("failed to create Neo4j constraints: %w", err)
	}

	if err := createNeo4jIndexes(ctx, neo4jConn); err != nil {
		return fmt.Errorf("failed to create Neo4j indexes: %w", err)
	}

	return nil
}

// fetchGraphNodes fetches graph nodes from Neo4j
func fetchGraphNodes(ctx *gofr.Context, session *Neo4jSession, orgName string, useTopics bool) ([]GraphNode, error) {
	nodesQuery := buildGraphNodesQuery(orgName, useTopics)
	nodesResult, err := executeNeo4jReadQuery(ctx, session, nodesQuery, map[string]interface{}{
		"orgName": orgName,
	})
	if err != nil {
		return nil, err
	}

	return convertToGraphNodes(nodesResult.Records), nil
}

// fetchGraphEdges fetches graph edges from Neo4j
func fetchGraphEdges(ctx *gofr.Context, session *Neo4jSession, orgName string, useTopics bool) ([]GraphEdge, error) {
	edgesQuery := buildGraphEdgesQuery(orgName, useTopics)
	edgesResult, err := executeNeo4jReadQuery(ctx, session, edgesQuery, map[string]interface{}{
		"orgName": orgName,
	})
	if err != nil {
		return nil, err
	}

	return convertToGraphEdges(edgesResult.Records), nil
}

// buildScanResponse builds scan response from components
func buildScanResponse(organization string, summary ScanSummary, org GitHubOrganization, repos []GitHubRepository, teams []GitHubTeam, topics []GitHubTopic, codeowners []GitHubCodeowners) ScanResponse {
	return ScanResponse{
		Success:      true,
		Organization: organization,
		Summary:      summary,
		Errors:       []string{},
		Data: map[string]interface{}{
			"organization": org,
			"repositories": repos,
			"teams":        teams,
			"topics":       topics,
			"codeowners":   codeowners,
		},
	}
}

// storeRepositories stores multiple repositories in Neo4j
func storeRepositories(ctx *gofr.Context, session *Neo4jSession, repos []GitHubRepository, orgLogin string) error {
	for _, repo := range repos {
		if err := storeRepository(ctx, session, repo, orgLogin); err != nil {
			return fmt.Errorf("failed to store repository %s: %w", repo.Name, err)
		}
	}
	return nil
}

// storeTeamsAndTopics stores teams and topics in Neo4j
func storeTeamsAndTopics(ctx *gofr.Context, session *Neo4jSession, teams []GitHubTeam, topics []GitHubTopic, orgLogin string) error {
	for _, team := range teams {
		if err := storeTeam(ctx, session, team, orgLogin); err != nil {
			return fmt.Errorf("failed to store team %s: %w", team.Name, err)
		}
	}

	for _, topic := range topics {
		if err := storeTopic(ctx, session, topic, orgLogin); err != nil {
			return fmt.Errorf("failed to store topic %s: %w", topic.Name, err)
		}
	}

	return nil
}

// storeCodeownersData stores codeowners data in Neo4j
func storeCodeownersData(ctx *gofr.Context, session *Neo4jSession, codeowners []GitHubCodeowners, orgLogin string) error {
	for _, codeowner := range codeowners {
		if err := storeCodeowners(ctx, session, codeowner, orgLogin); err != nil {
			return fmt.Errorf("failed to store CODEOWNERS for %s: %w", codeowner.Repository, err)
		}
	}
	return nil
}

// extractUniqueOwners extracts unique owners from codeowners data
func extractUniqueOwners(codeowners []GitHubCodeowners) map[string]bool {
	uniqueOwners := make(map[string]bool)

	for _, codeowner := range codeowners {
		for _, rule := range codeowner.Rules {
			for _, owner := range rule.Owners {
				uniqueOwners[owner] = true
			}
		}
	}

	return uniqueOwners
}

// estimateAPICallsUsed estimates the number of API calls used
func estimateAPICallsUsed(repos []GitHubRepository, teams []GitHubTeam, codeowners []GitHubCodeowners) int {
	return len(repos) + len(teams) + len(codeowners) + 1
}

// convertNeo4jErrorByMessage converts Neo4j errors based on message content
func convertNeo4jErrorByMessage(err error) error {
	errStr := err.Error()
	switch {
	case containsErrorKeywords(errStr, []string{"no data found", "not found", "record not found"}):
		return &gofrhttp.ErrorEntityNotFound{
			Name:  "data",
			Value: "requested resource",
		}
	case containsErrorKeywords(errStr, []string{"invalid parameter", "validation failed", "constraint violation"}):
		return &gofrhttp.ErrorInvalidParam{
			Params: []string{"database_constraint"},
		}
	case containsErrorKeywords(errStr, []string{"timeout", "connection timeout"}):
		return &gofrhttp.ErrorRequestTimeout{}
	default:
		return err
	}
}

// containsErrorKeywords checks if error string contains any of the keywords
func containsErrorKeywords(errStr string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}