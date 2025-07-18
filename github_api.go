package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"
	gofrhttp "gofr.dev/pkg/gofr/http"
)

// GitHubAPIClient represents a GitHub API client
type GitHubAPIClient struct {
	config     GitHubConfig
	httpClient *http.Client
}

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

// createGitHubAPIClient creates a new GitHub API client (Pure Core)
func createGitHubAPIClient(config GitHubConfig) GitHubAPIClient {
	validateGitHubConfig(config)

	return GitHubAPIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// buildGitHubGraphQLRequest builds a GraphQL request (Pure Core)
func buildGitHubGraphQLRequest(query string, variables map[string]interface{}) GitHubGraphQLRequest {
	validateQueryNotEmpty(query)

	if variables == nil {
		variables = make(map[string]interface{})
	}

	return GitHubGraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}

// buildGitHubAPIRequest builds an API request (Pure Core)
func buildGitHubAPIRequest(method, url string, headers map[string]string, body []byte) GitHubAPIRequest {
	validateMethodNotEmpty(method)
	validateURLNotEmpty(url)

	if headers == nil {
		headers = make(map[string]string)
	}

	return GitHubAPIRequest{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}
}

// buildGitHubAuthHeaders builds authentication headers (Pure Core)
func buildGitHubAuthHeaders(token string) map[string]string {
	validateTokenNotEmpty(token)

	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Accept":        "application/vnd.github.v3+json",
		"User-Agent":    "overseer-github-scanner/1.0",
	}
}

// buildOrganizationQuery builds GraphQL query for organization data (Pure Core)
func buildOrganizationQuery() string {
	return `
		query($login: String!) {
			organization(login: $login) {
				id
				login
				name
				description
				url
				email
				createdAt
				updatedAt
			}
		}
	`
}

// buildRepositoriesQuery builds GraphQL query for repositories (Pure Core)
func buildRepositoriesQuery() string {
	return `
		query($login: String!, $first: Int!, $after: String) {
			organization(login: $login) {
				repositories(first: $first, after: $after) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						name
						nameWithOwner
						description
						url
						isPrivate
						createdAt
						updatedAt
					}
				}
			}
		}
	`
}

// buildTeamsQuery builds GraphQL query for teams (Pure Core)
func buildTeamsQuery() string {
	return `
		query($login: String!, $first: Int!, $after: String) {
			organization(login: $login) {
				teams(first: $first, after: $after) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						slug
						name
						description
						url
					}
				}
			}
		}
	`
}

// buildCodeownersQuery builds query for CODEOWNERS file (Pure Core)
func buildCodeownersQuery() string {
	return `
		query($owner: String!, $name: String!) {
			repository(owner: $owner, name: $name) {
				object(expression: "HEAD:.github/CODEOWNERS") {
					... on Blob {
						text
					}
				}
			}
		}
	`
}

// executeGitHubGraphQL executes a GraphQL request (Orchestrator)
func executeGitHubGraphQL(ctx context.Context, client GitHubAPIClient, request GitHubGraphQLRequest) (GitHubGraphQLResponse, error) {
	validateGitHubGraphQLRequest(request)

	body, err := json.Marshal(request)
	if err != nil {
		return GitHubGraphQLResponse{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"graphql_request", err.Error()},
		}
	}

	headers := buildGitHubAuthHeaders(client.config.Token)
	headers["Content-Type"] = "application/json"

	apiRequest := buildGitHubAPIRequest("POST", client.config.BaseURL+"/graphql", headers, body)
	response, err := executeGitHubAPIRequest(ctx, client, apiRequest)
	if err != nil {
		return GitHubGraphQLResponse{}, &gofrhttp.ErrorRequestTimeout{}
	}

	if response.StatusCode == http.StatusNotFound {
		return GitHubGraphQLResponse{}, &gofrhttp.ErrorEntityNotFound{
			Name:  "graphql_endpoint",
			Value: "github",
		}
	}

	if response.StatusCode != http.StatusOK {
		return GitHubGraphQLResponse{}, GitHubAPIError{
			Code:       "GRAPHQL_ERROR",
			Message:    fmt.Sprintf("GraphQL request failed with status %d", response.StatusCode),
			Details:    string(response.Body),
			HTTPStatus: response.StatusCode,
		}
	}

	var graphqlResponse GitHubGraphQLResponse
	if err := json.Unmarshal(response.Body, &graphqlResponse); err != nil {
		return GitHubGraphQLResponse{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_format", err.Error()},
		}
	}

	if len(graphqlResponse.Errors) > 0 {
		return GitHubGraphQLResponse{}, GitHubAPIError{
			Code:       "GRAPHQL_ERRORS",
			Message:    "GraphQL query returned errors",
			Details:    formatGraphQLErrors(graphqlResponse.Errors),
			HTTPStatus: http.StatusBadRequest,
		}
	}

	return graphqlResponse, nil
}

// executeGitHubAPIRequest executes an API request (Orchestrator)
func executeGitHubAPIRequest(ctx context.Context, client GitHubAPIClient, request GitHubAPIRequest) (GitHubAPIResponse, error) {
	validateGitHubAPIRequest(request)

	var body io.Reader
	if len(request.Body) > 0 {
		body = bytes.NewReader(request.Body)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, request.URL, body)
	if err != nil {
		return GitHubAPIResponse{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"http_request", err.Error()},
		}
	}

	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return GitHubAPIResponse{}, &gofrhttp.ErrorRequestTimeout{}
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return GitHubAPIResponse{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_body", err.Error()},
		}
	}

	responseHeaders := make(map[string]string)
	for key, values := range httpResponse.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	return GitHubAPIResponse{
		StatusCode: httpResponse.StatusCode,
		Headers:    responseHeaders,
		Body:       responseBody,
	}, nil
}

// fetchGitHubOrganization fetches organization data (Orchestrator)
func fetchGitHubOrganization(ctx context.Context, client GitHubAPIClient, orgLogin string) (GitHubOrganization, error) {
	validateOrgLoginNotEmpty(orgLogin)

	query := buildOrganizationQuery()
	variables := map[string]interface{}{
		"login": orgLogin,
	}

	request := buildGitHubGraphQLRequest(query, variables)
	response, err := executeGitHubGraphQL(ctx, client, request)
	if err != nil {
		return GitHubOrganization{}, err
	}

	return parseOrganizationResponse(response)
}

// fetchGitHubRepositories fetches repositories for an organization (Orchestrator)
func fetchGitHubRepositories(ctx context.Context, client GitHubAPIClient, orgLogin string, maxRepos int) ([]GitHubRepository, error) {
	validateOrgLoginNotEmpty(orgLogin)
	validateMaxReposPositive(maxRepos)

	var allRepos []GitHubRepository
	var cursor string
	hasNextPage := true

	for hasNextPage && len(allRepos) < maxRepos {
		query := buildRepositoriesQuery()
		variables := map[string]interface{}{
			"login": orgLogin,
			"first": min(50, maxRepos-len(allRepos)),
		}

		if cursor != "" {
			variables["after"] = cursor
		}

		request := buildGitHubGraphQLRequest(query, variables)
		response, err := executeGitHubGraphQL(ctx, client, request)
		if err != nil {
			return nil, err
		}

		repos, nextCursor, hasNext, err := parseRepositoriesResponse(response)
		if err != nil {
			return nil, &gofrhttp.ErrorInvalidParam{
				Params: []string{"response_parsing", err.Error()},
			}
		}

		allRepos = append(allRepos, repos...)
		cursor = nextCursor
		hasNextPage = hasNext
	}

	return allRepos, nil
}

// fetchGitHubTeams fetches teams for an organization (Orchestrator)
func fetchGitHubTeams(ctx context.Context, client GitHubAPIClient, orgLogin string, maxTeams int) ([]GitHubTeam, error) {
	validateOrgLoginNotEmpty(orgLogin)
	validateMaxTeamsPositive(maxTeams)

	var allTeams []GitHubTeam
	var cursor string
	hasNextPage := true

	for hasNextPage && len(allTeams) < maxTeams {
		query := buildTeamsQuery()
		variables := map[string]interface{}{
			"login": orgLogin,
			"first": min(50, maxTeams-len(allTeams)),
		}

		if cursor != "" {
			variables["after"] = cursor
		}

		request := buildGitHubGraphQLRequest(query, variables)
		response, err := executeGitHubGraphQL(ctx, client, request)
		if err != nil {
			return nil, err
		}

		teams, nextCursor, hasNext, err := parseTeamsResponse(response)
		if err != nil {
			return nil, &gofrhttp.ErrorInvalidParam{
				Params: []string{"response_parsing", err.Error()},
			}
		}

		allTeams = append(allTeams, teams...)
		cursor = nextCursor
		hasNextPage = hasNext
	}

	return allTeams, nil
}

// fetchGitHubCodeowners fetches CODEOWNERS file for a repository (Orchestrator)
func fetchGitHubCodeowners(ctx context.Context, client GitHubAPIClient, owner, repo string) (GitHubCodeowners, error) {
	validateOwnerNotEmpty(owner)
	validateRepoNotEmpty(repo)

	query := buildCodeownersQuery()
	variables := map[string]interface{}{
		"owner": owner,
		"name":  repo,
	}

	request := buildGitHubGraphQLRequest(query, variables)
	response, err := executeGitHubGraphQL(ctx, client, request)
	if err != nil {
		return GitHubCodeowners{}, err
	}

	return parseCodeownersResponse(response, fmt.Sprintf("%s/%s", owner, repo))
}

// parseOrganizationResponse parses organization GraphQL response (Pure Core)
func parseOrganizationResponse(response GitHubGraphQLResponse) (GitHubOrganization, error) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return GitHubOrganization{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_data_format", "expected_map_string_interface"},
		}
	}

	org, ok := data["organization"].(map[string]interface{})
	if !ok {
		return GitHubOrganization{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"organization_data_format", "expected_map_string_interface"},
		}
	}

	return GitHubOrganization{
		ID:          getIntFromMap(org, "id"),
		Login:       getStringFromMap(org, "login"),
		Name:        getStringFromMap(org, "name"),
		Description: getStringFromMap(org, "description"),
		URL:         getStringFromMap(org, "url"),
		Email:       getStringFromMap(org, "email"),
		CreatedAt:   parseTimeFromMap(org, "createdAt"),
		UpdatedAt:   parseTimeFromMap(org, "updatedAt"),
	}, nil
}

// parseRepositoriesResponse parses repositories GraphQL response (Pure Core)
func parseRepositoriesResponse(response GitHubGraphQLResponse) ([]GitHubRepository, string, bool, error) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid response data format",
			Details: "Expected map[string]interface{}",
		}
	}

	org, ok := data["organization"].(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid organization data format",
			Details: "Expected map[string]interface{}",
		}
	}

	repositories, ok := org["repositories"].(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid repositories data format",
			Details: "Expected map[string]interface{}",
		}
	}

	pageInfo, ok := repositories["pageInfo"].(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid pageInfo data format",
			Details: "Expected map[string]interface{}",
		}
	}

	nodes, ok := repositories["nodes"].([]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid nodes data format",
			Details: "Expected []interface{}",
		}
	}

	repos := lo.Map(nodes, func(node interface{}, _ int) GitHubRepository {
		repo := node.(map[string]interface{})
		return GitHubRepository{
			ID:          getIntFromMap(repo, "id"),
			Name:        getStringFromMap(repo, "name"),
			FullName:    getStringFromMap(repo, "nameWithOwner"),
			Description: getStringFromMap(repo, "description"),
			URL:         getStringFromMap(repo, "url"),
			Private:     getBoolFromMap(repo, "isPrivate"),
			CreatedAt:   parseTimeFromMap(repo, "createdAt"),
			UpdatedAt:   parseTimeFromMap(repo, "updatedAt"),
		}
	})

	endCursor := getStringFromMap(pageInfo, "endCursor")
	hasNextPage := getBoolFromMap(pageInfo, "hasNextPage")

	return repos, endCursor, hasNextPage, nil
}

// parseTeamsResponse parses teams GraphQL response (Pure Core)
func parseTeamsResponse(response GitHubGraphQLResponse) ([]GitHubTeam, string, bool, error) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid response data format",
			Details: "Expected map[string]interface{}",
		}
	}

	org, ok := data["organization"].(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid organization data format",
			Details: "Expected map[string]interface{}",
		}
	}

	teams, ok := org["teams"].(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid teams data format",
			Details: "Expected map[string]interface{}",
		}
	}

	pageInfo, ok := teams["pageInfo"].(map[string]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid pageInfo data format",
			Details: "Expected map[string]interface{}",
		}
	}

	nodes, ok := teams["nodes"].([]interface{})
	if !ok {
		return nil, "", false, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid nodes data format",
			Details: "Expected []interface{}",
		}
	}

	teamsList := lo.Map(nodes, func(node interface{}, _ int) GitHubTeam {
		team := node.(map[string]interface{})
		return GitHubTeam{
			ID:          getIntFromMap(team, "id"),
			Slug:        getStringFromMap(team, "slug"),
			Name:        getStringFromMap(team, "name"),
			Description: getStringFromMap(team, "description"),
			URL:         getStringFromMap(team, "url"),
		}
	})

	endCursor := getStringFromMap(pageInfo, "endCursor")
	hasNextPage := getBoolFromMap(pageInfo, "hasNextPage")

	return teamsList, endCursor, hasNextPage, nil
}

// parseCodeownersResponse parses CODEOWNERS GraphQL response (Pure Core)
func parseCodeownersResponse(response GitHubGraphQLResponse, repoName string) (GitHubCodeowners, error) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return GitHubCodeowners{}, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid response data format",
			Details: "Expected map[string]interface{}",
		}
	}

	repo, ok := data["repository"].(map[string]interface{})
	if !ok {
		return GitHubCodeowners{}, GitHubAPIError{
			Code:    "PARSE_ERROR",
			Message: "Invalid repository data format",
			Details: "Expected map[string]interface{}",
		}
	}

	object, ok := repo["object"].(map[string]interface{})
	if !ok {
		// No CODEOWNERS file found
		return GitHubCodeowners{
			Repository: repoName,
			Rules:      []GitHubCodeownersRule{},
			Errors:     []GitHubCodeownersError{},
		}, nil
	}

	text := getStringFromMap(object, "text")
	return parseCodeownersText(text, repoName), nil
}

// parseCodeownersText parses CODEOWNERS file content (Pure Core)
func parseCodeownersText(text, repoName string) GitHubCodeowners {
	if text == "" {
		return GitHubCodeowners{
			Repository: repoName,
			Rules:      []GitHubCodeownersRule{},
			Errors:     []GitHubCodeownersError{},
		}
	}

	lines := lo.Filter(lo.Map(strings.Split(text, "\n"), func(line string, _ int) string {
		return strings.TrimSpace(line)
	}), func(line string, _ int) bool {
		return line != "" && !strings.HasPrefix(line, "#")
	})

	rules := lo.Map(lines, func(line string, index int) GitHubCodeownersRule {
		parts := lo.Filter(strings.Split(line, " "), func(part string, _ int) bool {
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

	return GitHubCodeowners{
		Repository: repoName,
		Rules:      rules,
		Errors:     []GitHubCodeownersError{},
	}
}

// Helper functions (Pure Core)
func getStringFromMap(m map[string]interface{}, key string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
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

func getBoolFromMap(m map[string]interface{}, key string) bool {
	if value, exists := m[key]; exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

func parseTimeFromMap(m map[string]interface{}, key string) time.Time {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			if t, err := time.Parse(time.RFC3339, str); err == nil {
				return t
			}
		}
	}
	return time.Time{}
}

func formatGraphQLErrors(errors []GitHubGraphQLError) string {
	messages := lo.Map(errors, func(err GitHubGraphQLError, _ int) string {
		return err.Message
	})
	return strings.Join(messages, "; ")
}

func wrapGitHubAPIError(err error, message string) GitHubAPIError {
	if err == nil {
		return GitHubAPIError{
			Code:       "INTERNAL_ERROR",
			Message:    message,
			Details:    "nil error wrapped",
			HTTPStatus: http.StatusInternalServerError,
		}
	}

	return GitHubAPIError{
		Code:       "API_ERROR",
		Message:    message,
		Details:    err.Error(),
		HTTPStatus: http.StatusBadGateway,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Validation helper functions (Pure Core)

func validateGitHubGraphQLRequest(request GitHubGraphQLRequest) {
	if request.Query == "" {
		panic("GraphQL query cannot be empty")
	}
}

func validateGitHubAPIRequest(request GitHubAPIRequest) {
	if request.Method == "" {
		panic("HTTP method cannot be empty")
	}
	if request.URL == "" {
		panic("HTTP URL cannot be empty")
	}
}

func validateMethodNotEmpty(method string) {
	if method == "" {
		panic("HTTP method cannot be empty")
	}
}

func validateURLNotEmpty(url string) {
	if url == "" {
		panic("URL cannot be empty")
	}
}

func validateTokenNotEmpty(token string) {
	if token == "" {
		panic("Token cannot be empty")
	}
}

func validateOrgLoginNotEmpty(orgLogin string) {
	if orgLogin == "" {
		panic("Organization login cannot be empty")
	}
}

func validateMaxReposPositive(maxRepos int) {
	if maxRepos <= 0 {
		panic("Max repositories must be positive")
	}
}

func validateMaxTeamsPositive(maxTeams int) {
	if maxTeams <= 0 {
		panic("Max teams must be positive")
	}
}

func validateOwnerNotEmpty(owner string) {
	if owner == "" {
		panic("Owner cannot be empty")
	}
}

func validateRepoNotEmpty(repo string) {
	if repo == "" {
		panic("Repository name cannot be empty")
	}
}
