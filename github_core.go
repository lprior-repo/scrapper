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

// buildGraphQLOrgQuery builds a GraphQL query for organization data (Pure Core)
func buildGraphQLOrgQuery(request BatchRequest) string {
	if request.Organization == "" {
		panic("Organization cannot be empty")
	}

	repoLimit := 100
	if request.MaxRepos > 0 && request.MaxRepos < 100 {
		repoLimit = request.MaxRepos
	}

	teamLimit := 100
	if request.MaxTeams > 0 && request.MaxTeams < 100 {
		teamLimit = request.MaxTeams
	}

	return fmt.Sprintf(`
		query($org: String!, $reposCursor: String, $teamsCursor: String) {
			organization(login: $org) {
				repositories(first: %d, after: $reposCursor) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						name
						nameWithOwner
						defaultBranchRef {
							name
						}
						codeowners: object(expression: "HEAD:CODEOWNERS") {
							... on Blob {
								text
							}
						}
						docsCodeowners: object(expression: "HEAD:.github/CODEOWNERS") {
							... on Blob {
								text
							}
						}
						rootCodeowners: object(expression: "HEAD:docs/CODEOWNERS") {
							... on Blob {
								text
							}
						}
					}
				}
				teams(first: %d, after: $teamsCursor) {
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
						members {
							totalCount
						}
					}
				}
			}
		}
	`, repoLimit, teamLimit)
}

// parseCodeownersContent parses CODEOWNERS file content (Pure Core)
func parseCodeownersContent(content string) []CodeownersEntry {
	if content == "" {
		return []CodeownersEntry{}
	}

	lines := strings.Split(content, "\n")
	entries := []CodeownersEntry{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}

		pattern := parts[0]
		owners := parts[1:]

		entries = append(entries, CodeownersEntry{
			Pattern: pattern,
			Owners:  owners,
		})
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
