package main

import (
	"context"
	"fmt"
)

// findRepositoriesByCodeowner finds all repositories that have a specific codeowner
func findRepositoriesByCodeowner(ctx context.Context, conn GraphConnection, ownerName string) ([]*Node, error) {
	if ownerName == "" {
		return nil, fmt.Errorf("owner name cannot be empty")
	}

	query := `
		MATCH (owner {name: $ownerName})<-[:HAS_CODEOWNER]-(repo:Repository)
		RETURN id(repo) as id, labels(repo) as labels, properties(repo) as properties
		ORDER BY repo.full_name
	`

	params := map[string]interface{}{
		"ownerName": ownerName,
	}

	results, err := executeReadQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find repositories by codeowner: %w", err)
	}

	nodes := make([]*Node, len(results))
	for i, result := range results {
		node, err := processNodeResult(result)
		if err != nil {
			return nil, fmt.Errorf("failed to process node result: %w", err)
		}
		nodes[i] = node
	}

	return nodes, nil
}

// findCodeownersByRepository finds all codeowners for a specific repository
func findCodeownersByRepository(ctx context.Context, conn GraphConnection, repoName string) ([]*Node, error) {
	if repoName == "" {
		return nil, fmt.Errorf("repository name cannot be empty")
	}

	query := `
		MATCH (repo:Repository {full_name: $repoName})-[:HAS_CODEOWNER]->(owner)
		RETURN id(owner) as id, labels(owner) as labels, properties(owner) as properties
		ORDER BY owner.name
	`

	params := map[string]interface{}{
		"repoName": repoName,
	}

	results, err := executeReadQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find codeowners by repository: %w", err)
	}

	nodes := make([]*Node, len(results))
	for i, result := range results {
		node, err := processNodeResult(result)
		if err != nil {
			return nil, fmt.Errorf("failed to process node result: %w", err)
		}
		nodes[i] = node
	}

	return nodes, nil
}

// GetCodeownershipStats gets comprehensive codeownership statistics for an organization
func GetCodeownershipStats(ctx context.Context, conn GraphConnection, orgName string) (map[string]interface{}, error) {
	if orgName == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}

	query := `
		MATCH (org:Organization {name: $orgName})-[:OWNS]->(repo:Repository)
		OPTIONAL MATCH (repo)-[:HAS_CODEOWNER]->(owner)
		RETURN 
			COUNT(DISTINCT repo) as total_repos,
			COUNT(DISTINCT CASE WHEN owner IS NOT NULL THEN repo END) as repos_with_codeowners,
			COUNT(DISTINCT owner) as unique_owners,
			COUNT(DISTINCT CASE WHEN owner.type = 'team' THEN owner END) as team_owners,
			COUNT(DISTINCT CASE WHEN owner.type = 'user' THEN owner END) as user_owners
	`

	params := map[string]interface{}{
		"orgName": orgName,
	}

	results, err := executeReadQuery(ctx, conn, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get codeownership stats: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no data found for organization: %s", orgName)
	}

	stats := results[0]
	stats["organization"] = orgName

	return stats, nil
}
