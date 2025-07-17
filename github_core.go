package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

// GitHubRepo represents a GitHub repository
type GitHubRepo struct {
	Name              string   `json:"name"`
	FullName          string   `json:"full_name"`
	DefaultBranch     string   `json:"default_branch"`
	HasCodeownersFile bool     `json:"has_codeowners_file"`
	CodeownersContent string   `json:"codeowners_content"`
	CodeownersPaths   []string `json:"codeowners_paths"`
}

// GitHubTeam represents a GitHub team
type GitHubTeam struct {
	ID          int64  `json:"id"`
	NodeID      string `json:"node_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	MemberCount int    `json:"member_count"`
}

// GitHubOrgData represents all data from a GitHub organization
type GitHubOrgData struct {
	Organization string       `json:"organization"`
	Repos        []GitHubRepo `json:"repos"`
	Teams        []GitHubTeam `json:"teams"`
	APICallCount int          `json:"api_call_count"`
}

// CodeownersEntry represents a parsed CODEOWNERS entry
type CodeownersEntry struct {
	Pattern string   `json:"pattern"`
	Owners  []string `json:"owners"`
}

// BatchRequest represents a batch API request configuration
type BatchRequest struct {
	Organization string `json:"organization"`
	MaxRepos     int    `json:"max_repos"`
	MaxTeams     int    `json:"max_teams"`
}

// validateBatchRequest validates batch request parameters (Pure Core)
func validateBatchRequest(request BatchRequest) error {
	if request.Organization == "" {
		return fmt.Errorf("organization name is required")
	}
	if strings.TrimSpace(request.Organization) == "" {
		return fmt.Errorf("organization name cannot be empty or whitespace")
	}
	if request.MaxRepos < 0 {
		return fmt.Errorf("max repos cannot be negative")
	}
	if request.MaxTeams < 0 {
		return fmt.Errorf("max teams cannot be negative")
	}
	return nil
}

// getRepoLimit calculates the repository limit for queries (Pure Core)
func getRepoLimit(maxRepos int) int {
	if maxRepos > 0 && maxRepos < 100 {
		return maxRepos
	}
	return 100
}

// getTeamLimit calculates the team limit for queries (Pure Core)
func getTeamLimit(maxTeams int) int {
	if maxTeams > 0 && maxTeams < 100 {
		return maxTeams
	}
	return 100
}

// buildGraphQLOrgQuery builds an optimized GraphQL query for organization data (Pure Core)
func buildGraphQLOrgQuery(request BatchRequest) string {
	if request.Organization == "" {
		panic("Organization cannot be empty")
	}

	repoLimit := getRepoLimit(request.MaxRepos)
	teamLimit := getTeamLimit(request.MaxTeams)

	return fmt.Sprintf(`
		query($org: String!, $reposCursor: String, $teamsCursor: String) {
			rateLimit {
				cost
				remaining
				resetAt
			}
			organization(login: $org) {
				id
				login
				name
				description
				url
				createdAt
				repositories(first: %d, after: $reposCursor, orderBy: {field: UPDATED_AT, direction: DESC}) {
					totalCount
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						name
						nameWithOwner
						description
						isPrivate
						createdAt
						updatedAt
						defaultBranchRef {
							name
						}
						collaborators(first: 10) {
							totalCount
							nodes {
								login
								name
								email
							}
						}
						codeowners: object(expression: "HEAD:CODEOWNERS") {
							... on Blob {
								text
								byteSize
							}
						}
						docsCodeowners: object(expression: "HEAD:.github/CODEOWNERS") {
							... on Blob {
								text
								byteSize
							}
						}
						rootCodeowners: object(expression: "HEAD:docs/CODEOWNERS") {
							... on Blob {
								text
								byteSize
							}
						}
						languages(first: 5, orderBy: {field: SIZE, direction: DESC}) {
							nodes {
								name
								color
							}
						}
					}
				}
				teams(first: %d, after: $teamsCursor, orderBy: {field: NAME, direction: ASC}) {
					totalCount
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						databaseId
						name
						slug
						description
						privacy
						createdAt
						updatedAt
						members(first: 50) {
							totalCount
							nodes {
								login
								name
								email
								createdAt
							}
						}
						repositories(first: 20) {
							totalCount
							nodes {
								nameWithOwner
							}
						}
					}
				}
			}
		}
	`, repoLimit, teamLimit)
}

// buildOptimizedBatchQuery builds a query optimized for batch processing (Pure Core)
func buildOptimizedBatchQuery(orgName string, repoLimit, teamLimit int) string {
	if orgName == "" {
		panic("Organization name cannot be empty")
	}
	if repoLimit <= 0 || teamLimit <= 0 {
		panic("Limits must be positive")
	}

	return fmt.Sprintf(`
		query($org: String!) {
			rateLimit {
				cost
				remaining
				resetAt
				nodeCount
			}
			organization(login: $org) {
				id
				login
				name
				description
				url
				createdAt
				membersWithRole(first: 1) {
					totalCount
				}
				repositories(first: %d, orderBy: {field: UPDATED_AT, direction: DESC}) {
					totalCount
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						name
						nameWithOwner
						description
						isPrivate
						isArchived
						isFork
						createdAt
						updatedAt
						diskUsage
						defaultBranchRef {
							name
							target {
								... on Commit {
									committedDate
								}
							}
						}
						primaryLanguage {
							name
							color
						}
						collaborators(first: 10, affiliation: DIRECT) {
							totalCount
							nodes {
								login
								name
								email
								avatarUrl
							}
						}
						codeowners: object(expression: "HEAD:CODEOWNERS") {
							... on Blob {
								text
								byteSize
								oid
							}
						}
						docsCodeowners: object(expression: "HEAD:.github/CODEOWNERS") {
							... on Blob {
								text
								byteSize
								oid
							}
						}
						rootCodeowners: object(expression: "HEAD:docs/CODEOWNERS") {
							... on Blob {
								text
								byteSize
								oid
							}
						}
						issues(states: OPEN) {
							totalCount
						}
						pullRequests(states: OPEN) {
							totalCount
						}
					}
				}
				teams(first: %d, orderBy: {field: NAME, direction: ASC}) {
					totalCount
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						databaseId
						name
						slug
						description
						privacy
						createdAt
						updatedAt
						combinedSlug
						members(first: 100, membership: ALL) {
							totalCount
							nodes {
								login
								name
								email
								avatarUrl
								createdAt
								company
								location
							}
						}
						repositories(first: 50) {
							totalCount
							nodes {
								nameWithOwner
								isPrivate
							}
						}
						childTeams(first: 10) {
							totalCount
							nodes {
								name
								slug
							}
						}
					}
				}
			}
		}
	`, repoLimit, teamLimit)
}

// isValidCodeownersLine checks if a line is valid for processing (Pure Core)
func isValidCodeownersLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed != "" && !strings.HasPrefix(trimmed, "#")
}

// parseCodeownersLine parses a single CODEOWNERS line (Pure Core)
func parseCodeownersLine(line string) (CodeownersEntry, bool) {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) < 2 {
		return CodeownersEntry{}, false
	}

	return CodeownersEntry{
		Pattern: parts[0],
		Owners:  parts[1:],
	}, true
}

// parseCodeownersContent parses CODEOWNERS file content (Pure Core)
func parseCodeownersContent(content string) []CodeownersEntry {
	if content == "" {
		return []CodeownersEntry{}
	}

	lines := strings.Split(content, "\n")
	entries := []CodeownersEntry{}

	for _, line := range lines {
		if !isValidCodeownersLine(line) {
			continue
		}

		if entry, valid := parseCodeownersLine(line); valid {
			entries = append(entries, entry)
		}
	}

	return entries
}

// extractUniqueOwners extracts unique owners from all CODEOWNERS entries (Pure Core)
func extractUniqueOwners(entries []CodeownersEntry) []string {
	allOwners := []string{}

	for _, entry := range entries {
		allOwners = append(allOwners, entry.Owners...)
	}

	return lo.Uniq(allOwners)
}

// processRepoCodeowners processes codeowners data for a repository (Pure Core)
func processRepoCodeowners(repo map[string]interface{}) GitHubRepo {
	repoName := extractStringField(repo, "name")
	fullName := extractStringField(repo, "nameWithOwner")

	if repoName == "" {
		panic("Repository name cannot be empty")
	}

	defaultBranch := "main"
	if branchRef, ok := repo["defaultBranchRef"].(map[string]interface{}); ok {
		if branch := extractStringField(branchRef, "name"); branch != "" {
			defaultBranch = branch
		}
	}

	codeownersContent := ""
	codeownersPaths := []string{}

	// Check multiple possible locations for CODEOWNERS
	locations := []struct {
		field string
		path  string
	}{
		{"codeowners", "CODEOWNERS"},
		{"docsCodeowners", ".github/CODEOWNERS"},
		{"rootCodeowners", "docs/CODEOWNERS"},
	}

	for _, loc := range locations {
		if content := extractCodeownersText(repo, loc.field); content != "" {
			codeownersContent = content
			codeownersPaths = append(codeownersPaths, loc.path)
		}
	}

	return GitHubRepo{
		Name:              repoName,
		FullName:          fullName,
		DefaultBranch:     defaultBranch,
		HasCodeownersFile: codeownersContent != "",
		CodeownersContent: codeownersContent,
		CodeownersPaths:   codeownersPaths,
	}
}

// processTeamData processes team data from GraphQL response (Pure Core)
func processTeamData(team map[string]interface{}) GitHubTeam {
	teamName := extractStringField(team, "name")
	if teamName == "" {
		panic("Team name cannot be empty")
	}

	memberCount := 0
	if members, ok := team["members"].(map[string]interface{}); ok {
		if count, ok := members["totalCount"].(float64); ok {
			memberCount = int(count)
		}
	}

	return GitHubTeam{
		ID:          extractInt64Field(team, "databaseId"),
		NodeID:      extractStringField(team, "id"),
		Name:        teamName,
		Slug:        extractStringField(team, "slug"),
		Description: extractStringField(team, "description"),
		Privacy:     extractStringField(team, "privacy"),
		MemberCount: memberCount,
	}
}

// extractStringField safely extracts a string field from a map (Pure Core)
func extractStringField(data map[string]interface{}, field string) string {
	if data == nil {
		return ""
	}
	if value, ok := data[field].(string); ok {
		return value
	}
	return ""
}

// extractInt64Field safely extracts an int64 field from a map (Pure Core)
func extractInt64Field(data map[string]interface{}, field string) int64 {
	if data == nil {
		return 0
	}
	if value, ok := data[field].(float64); ok {
		return int64(value)
	}
	return 0
}

// extractCodeownersText extracts CODEOWNERS text from a blob field (Pure Core)
func extractCodeownersText(repo map[string]interface{}, field string) string {
	if repo == nil {
		return ""
	}
	if blob, ok := repo[field].(map[string]interface{}); ok {
		if text, ok := blob["text"].(string); ok {
			return text
		}
	}
	return ""
}

// calculateAPICallsNeeded estimates API calls needed for an org (Pure Core)
func calculateAPICallsNeeded(repoCount, teamCount int) int {
	if repoCount < 0 || teamCount < 0 {
		panic("Counts cannot be negative")
	}

	// GraphQL can fetch up to 100 items per call
	repoCalls := (repoCount + 99) / 100
	teamCalls := (teamCount + 99) / 100

	// We can fetch both in same query, so take the max
	return lo.Max([]int{repoCalls, teamCalls, 1})
}

// optimizeBatchSize calculates optimal batch size to minimize API calls (Pure Core)
func optimizeBatchSize(totalItems int, maxPerCall int) int {
	if totalItems <= 0 {
		return 0
	}
	if maxPerCall <= 0 {
		panic("Max per call must be positive")
	}

	if totalItems <= maxPerCall {
		return totalItems
	}
	return maxPerCall
}

// RateLimitInfo represents GitHub API rate limit information
type RateLimitInfo struct {
	Cost      int    `json:"cost"`
	Remaining int    `json:"remaining"`
	ResetAt   string `json:"resetAt"`
	NodeCount int    `json:"nodeCount"`
}

// BatchPagination represents pagination state for batch processing
type BatchPagination struct {
	ReposCursor string `json:"repos_cursor"`
	TeamsCursor string `json:"teams_cursor"`
	HasMoreData bool   `json:"has_more_data"`
}

// OptimizedQueryParams represents optimized query parameters
type OptimizedQueryParams struct {
	Organization string          `json:"organization"`
	RepoLimit    int             `json:"repo_limit"`
	TeamLimit    int             `json:"team_limit"`
	Pagination   BatchPagination `json:"pagination"`
	RateLimit    RateLimitInfo   `json:"rate_limit"`
}

// calculateOptimalBatchSize determines optimal batch size based on rate limits (Pure Core)
func calculateOptimalBatchSize(remaining int, cost int, maxItems int) int {
	if remaining <= 0 || cost <= 0 {
		return 1
	}
	if maxItems <= 0 {
		panic("Max items must be positive")
	}

	// Calculate how many full batches we can make
	possibleBatches := remaining / cost
	if possibleBatches <= 0 {
		return 1
	}

	// Don't exceed the maximum items per request
	optimalSize := lo.Min([]int{maxItems, possibleBatches * 10})
	return lo.Max([]int{optimalSize, 1})
}

// buildMinimalOrgQuery builds a minimal query for initial organization data (Pure Core)
func buildMinimalOrgQuery() string {
	return `
		query($org: String!) {
			rateLimit {
				cost
				remaining
				resetAt
			}
			organization(login: $org) {
				id
				login
				name
				description
				url
				createdAt
				repositories {
					totalCount
				}
				teams {
					totalCount
				}
				membersWithRole(first: 1) {
					totalCount
				}
			}
		}
	`
}

// buildRepositoriesOnlyQuery builds a query focused only on repositories (Pure Core)
func buildRepositoriesOnlyQuery(limit int, cursor string) string {
	if limit <= 0 {
		panic("Limit must be positive")
	}

	return fmt.Sprintf(`
		query($org: String!, $cursor: String) {
			rateLimit {
				cost
				remaining
				resetAt
			}
			organization(login: $org) {
				repositories(first: %d, after: $cursor, orderBy: {field: UPDATED_AT, direction: DESC}) {
					totalCount
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						name
						nameWithOwner
						description
						isPrivate
						isArchived
						isFork
						createdAt
						updatedAt
						defaultBranchRef {
							name
						}
						codeowners: object(expression: "HEAD:CODEOWNERS") {
							... on Blob {
								text
								byteSize
							}
						}
						docsCodeowners: object(expression: "HEAD:.github/CODEOWNERS") {
							... on Blob {
								text
								byteSize
							}
						}
						rootCodeowners: object(expression: "HEAD:docs/CODEOWNERS") {
							... on Blob {
								text
								byteSize
							}
						}
					}
				}
			}
		}
	`, limit)
}

// buildTeamsOnlyQuery builds a query focused only on teams (Pure Core)
func buildTeamsOnlyQuery(limit int, cursor string) string {
	if limit <= 0 {
		panic("Limit must be positive")
	}

	return fmt.Sprintf(`
		query($org: String!, $cursor: String) {
			rateLimit {
				cost
				remaining
				resetAt
			}
			organization(login: $org) {
				teams(first: %d, after: $cursor, orderBy: {field: NAME, direction: ASC}) {
					totalCount
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						databaseId
						name
						slug
						description
						privacy
						createdAt
						updatedAt
						members(first: 100) {
							totalCount
							nodes {
								login
								name
								email
								avatarUrl
							}
						}
					}
				}
			}
		}
	`, limit)
}

// estimateQueryCost estimates the cost of a GraphQL query (Pure Core)
func estimateQueryCost(repoCount, teamCount, memberCount int) int {
	if repoCount < 0 || teamCount < 0 || memberCount < 0 {
		panic("Counts cannot be negative")
	}

	// Base cost for organization query
	baseCost := 1

	// Repository cost (includes codeowners fetching)
	repoCost := repoCount * 3 // 3 codeowner file checks per repo

	// Team cost (includes member fetching)
	teamCost := teamCount * 2

	// Member cost
	memberCost := memberCount / 10 // Members are batched

	total := baseCost + repoCost + teamCost + memberCost
	return lo.Max([]int{total, 1})
}

// splitLargeOrganization splits large organization scan into smaller batches (Pure Core)
func splitLargeOrganization(totalRepos, totalTeams int, maxCostPerQuery int) []BatchRequest {
	if totalRepos < 0 || totalTeams < 0 {
		panic("Totals cannot be negative")
	}
	if maxCostPerQuery <= 0 {
		panic("Max cost must be positive")
	}

	batches := []BatchRequest{}

	// If small enough, do in one batch
	estimatedCost := estimateQueryCost(totalRepos, totalTeams, totalTeams*20)
	if estimatedCost <= maxCostPerQuery {
		return []BatchRequest{{
			MaxRepos: totalRepos,
			MaxTeams: totalTeams,
		}}
	}

	// Split into multiple batches
	maxReposPerBatch := 50 // Conservative batch size
	maxTeamsPerBatch := 25 // Conservative batch size

	reposBatched := 0
	teamsBatched := 0

	for reposBatched < totalRepos || teamsBatched < totalTeams {
		reposInBatch := lo.Min([]int{maxReposPerBatch, totalRepos - reposBatched})
		teamsInBatch := lo.Min([]int{maxTeamsPerBatch, totalTeams - teamsBatched})

		if reposInBatch > 0 || teamsInBatch > 0 {
			batches = append(batches, BatchRequest{
				MaxRepos: reposInBatch,
				MaxTeams: teamsInBatch,
			})
		}

		reposBatched += reposInBatch
		teamsBatched += teamsInBatch
	}

	return batches
}
