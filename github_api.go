package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/samber/lo"
)


// GitHubAPIRequest represents a GitHub API request
type GitHubAPIRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// GitHubAPIResponse represents a GitHub API response
type GitHubAPIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// GitHubGraphQLRequest represents a GraphQL request
type GitHubGraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GitHubGraphQLResponse represents a GraphQL response
type GitHubGraphQLResponse struct {
	Data   interface{}          `json:"data"`
	Errors []GitHubGraphQLError `json:"errors,omitempty"`
}

// GitHubGraphQLError represents a GraphQL error
type GitHubGraphQLError struct {
	Message string   `json:"message"`
	Type    string   `json:"type"`
	Path    []string `json:"path,omitempty"`
}

// GitHubOrganization represents a GitHub organization
type GitHubOrganization struct {
	ID          int       `json:"id"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Private     bool      `json:"private"`
	Topics      []string  `json:"topics"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GitHubUser represents a GitHub user
type GitHubUser struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GitHubTeam represents a GitHub team
type GitHubTeam struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// GitHubTopic represents a GitHub repository topic
type GitHubTopic struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GitHubCodeowners represents CODEOWNERS file content
type GitHubCodeowners struct {
	Repository string                  `json:"repository"`
	Rules      []GitHubCodeownersRule  `json:"rules"`
	Errors     []GitHubCodeownersError `json:"errors"`
}

// GitHubCodeownersRule represents a CODEOWNERS rule
type GitHubCodeownersRule struct {
	Pattern string   `json:"pattern"`
	Owners  []string `json:"owners"`
	Line    int      `json:"line"`
}

// GitHubCodeownersError represents a CODEOWNERS parsing error
type GitHubCodeownersError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// GitHubAPIError represents GitHub API errors that implement GoFr error patterns
type GitHubAPIError struct {
	Code       string
	Message    string
	Details    string
	HTTPStatus int
}

// Error implements the error interface for GitHubAPIError
func (e GitHubAPIError) Error() string {
	return fmt.Sprintf("GitHub API error [%s]: %s - %s", e.Code, e.Message, e.Details)
}

// StatusCode returns the HTTP status code for the error
func (e GitHubAPIError) StatusCode() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}
	return http.StatusInternalServerError
}

// collectTopicsFromRepositories collects all unique topics from repositories with their counts (Pure Core)
func collectTopicsFromRepositories(repos []GitHubRepository) []GitHubTopic {
	topicCounts := make(map[string]int)

	for _, repo := range repos {
		for _, topic := range repo.Topics {
			if topic != "" {
				topicCounts[topic]++
			}
		}
	}

	topics := lo.Map(lo.Keys(topicCounts), func(topicName string, _ int) GitHubTopic {
		return GitHubTopic{
			Name:  topicName,
			Count: topicCounts[topicName],
		}
	})

	// Sort topics by count (descending) for better organization
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Count > topics[j].Count
	})
	return topics
}

// Helper functions (Pure Core)
func getStringFromMap(m map[string]interface{}, key string) string {
	value, exists := m[key]
	if !exists {
		return ""
	}
	return convertToString(value)
}

// convertToString converts various types to string representation
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

func getIntFromMap(m map[string]interface{}, key string) int {
	if value, exists := m[key]; exists {
		if i, ok := value.(int); ok {
			return i
		}
		if f, ok := value.(float64); ok {
			return int(f)
		}
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Validation helper functions (Pure Core)
func validateOrgLoginNotEmpty(orgLogin string) {
	if orgLogin == "" {
		panic("Organization login cannot be empty")
	}
}

func validateOwnerNotEmpty(owner string) {
	if owner == "" {
		panic("Owner cannot be empty")
	}
}
