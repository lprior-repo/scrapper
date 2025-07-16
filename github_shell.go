package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// GitHubClient represents a connection to GitHub
type GitHubClient struct {
	GraphQL *githubv4.Client
	Token   string
	Org     string
}

// GraphQLOrgResponse represents the GraphQL response structure
type GraphQLOrgResponse struct {
	Organization struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []map[string]interface{}
		}
		Teams struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []map[string]interface{}
		}
	}
}

// createGitHubClient creates a new GitHub client (Impure Shell)
func createGitHubClient(ctx context.Context, token string, org string) (GitHubClient, error) {
	if token == "" {
		return GitHubClient{}, fmt.Errorf("GitHub token is required")
	}
	if org == "" {
		return GitHubClient{}, fmt.Errorf("organization name is required")
	}
	
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	
	graphqlClient := githubv4.NewClient(httpClient)
	
	return GitHubClient{
		GraphQL: graphqlClient,
		Token:   token,
		Org:     org,
	}, nil
}

// executeGraphQLQuery executes a GraphQL query against GitHub (Impure Shell)
func executeGraphQLQuery(ctx context.Context, client GitHubClient, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	if client.GraphQL == nil {
		return nil, fmt.Errorf("GraphQL client not initialized")
	}
	if query == "" {
		panic("Query cannot be empty")
	}
	
	// Create HTTP request
	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+client.Token)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Check for errors
	if errors, ok := result["errors"]; ok {
		return nil, fmt.Errorf("GraphQL errors: %v", errors)
	}
	
	return result, nil
}

// fetchAllOrgData fetches all organization data with minimal API calls (Impure Shell)
func fetchAllOrgData(ctx context.Context, client GitHubClient, request BatchRequest) (GitHubOrgData, error) {
	if err := validateBatchRequest(request); err != nil {
		return GitHubOrgData{}, fmt.Errorf("invalid request: %w", err)
	}
	
	allRepos := []GitHubRepo{}
	allTeams := []GitHubTeam{}
	apiCallCount := 0
	
	reposCursor := ""
	teamsCursor := ""
	hasMoreRepos := true
	hasMoreTeams := true
	
	query := buildGraphQLOrgQuery(request)
	
	for hasMoreRepos || hasMoreTeams {
		variables := map[string]interface{}{
			"org": request.Organization,
		}
		
		if reposCursor != "" {
			variables["reposCursor"] = reposCursor
		}
		if teamsCursor != "" {
			variables["teamsCursor"] = teamsCursor
		}
		
		// Execute GraphQL query
		result, err := executeGraphQLQuery(ctx, client, query, variables)
		if err != nil {
			return GitHubOrgData{}, fmt.Errorf("failed to fetch org data: %w", err)
		}
		apiCallCount++
		
		// Process response
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			return GitHubOrgData{}, fmt.Errorf("invalid response format")
		}
		
		org, ok := data["organization"].(map[string]interface{})
		if !ok {
			return GitHubOrgData{}, fmt.Errorf("organization data not found")
		}
		
		// Process repositories
		if hasMoreRepos {
			repos, repoPageInfo := extractRepositories(org)
			for _, repo := range repos {
				allRepos = append(allRepos, processRepoCodeowners(repo))
			}
			
			hasMoreRepos = repoPageInfo["hasNextPage"].(bool)
			if cursor, ok := repoPageInfo["endCursor"].(string); ok && hasMoreRepos {
				reposCursor = cursor
			}
		}
		
		// Process teams
		if hasMoreTeams {
			teams, teamPageInfo := extractTeams(org)
			for _, team := range teams {
				allTeams = append(allTeams, processTeamData(team))
			}
			
			hasMoreTeams = teamPageInfo["hasNextPage"].(bool)
			if cursor, ok := teamPageInfo["endCursor"].(string); ok && hasMoreTeams {
				teamsCursor = cursor
			}
		}
		
		// Stop if we've reached our limits
		if request.MaxRepos > 0 && len(allRepos) >= request.MaxRepos {
			hasMoreRepos = false
		}
		if request.MaxTeams > 0 && len(allTeams) >= request.MaxTeams {
			hasMoreTeams = false
		}
		
		// Safety check to prevent infinite loops
		if apiCallCount >= 50 {
			break
		}
	}
	
	return GitHubOrgData{
		Organization: request.Organization,
		Repos:        allRepos,
		Teams:        allTeams,
		APICallCount: apiCallCount,
	}, nil
}

// extractRepositories extracts repository data from GraphQL response
func extractRepositories(org map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	repos := []map[string]interface{}{}
	pageInfo := map[string]interface{}{
		"hasNextPage": false,
		"endCursor":   "",
	}
	
	if reposData, ok := org["repositories"].(map[string]interface{}); ok {
		if nodes, ok := reposData["nodes"].([]interface{}); ok {
			for _, node := range nodes {
				if repo, ok := node.(map[string]interface{}); ok {
					repos = append(repos, repo)
				}
			}
		}
		
		if pi, ok := reposData["pageInfo"].(map[string]interface{}); ok {
			pageInfo = pi
		}
	}
	
	return repos, pageInfo
}

// extractTeams extracts team data from GraphQL response
func extractTeams(org map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	teams := []map[string]interface{}{}
	pageInfo := map[string]interface{}{
		"hasNextPage": false,
		"endCursor":   "",
	}
	
	if teamsData, ok := org["teams"].(map[string]interface{}); ok {
		if nodes, ok := teamsData["nodes"].([]interface{}); ok {
			for _, node := range nodes {
				if team, ok := node.(map[string]interface{}); ok {
					teams = append(teams, team)
				}
			}
		}
		
		if pi, ok := teamsData["pageInfo"].(map[string]interface{}); ok {
			pageInfo = pi
		}
	}
	
	return teams, pageInfo
}

// writeOrgDataToFile writes organization data to a JSON file (Impure Shell)
func writeOrgDataToFile(filename string, data GitHubOrgData) error {
	if filename == "" {
		panic("Filename cannot be empty")
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := writeFile(filename, jsonData); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// writeFile writes data to a file (Impure Shell helper)
func writeFile(filename string, data []byte) error {
	// This would normally use os.WriteFile but we're keeping it separate
	// for clarity of the impure shell pattern
	return fmt.Errorf("writeFile not implemented - use os.WriteFile")
}