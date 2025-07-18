package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
	"gofr.dev/pkg/gofr"
	gofrhttp "gofr.dev/pkg/gofr/http"
)

// GitHubServiceConfig represents GitHub service configuration
type GitHubServiceConfig struct {
	Token        string
	BaseURL      string
	UserAgent    string
	Timeout      time.Duration
	MaxRetries   int
	RateLimitMin int
}

// RegisterGitHubService registers GitHub as an HTTP service in GoFr
func RegisterGitHubService(app *gofr.App, config GitHubServiceConfig) {
	// Register GitHub API as an HTTP service
	app.AddHTTPService("github", config.BaseURL)
}

// fetchGitHubOrganizationWithService fetches organization data using GoFr HTTP service
func fetchGitHubOrganizationWithService(ctx *gofr.Context, orgName string) (GitHubOrganization, error) {
	ctx.Logger.Infof("Fetching GitHub organization: %s", orgName)

	if orgName == "" {
		ctx.Logger.Errorf("Organization name is required but was empty")
		return GitHubOrganization{}, &gofrhttp.ErrorMissingParam{
			Params: []string{"organization_name"},
		}
	}

	resp, err := makeGitHubOrgAPIRequest(ctx, orgName)
	if err != nil {
		return GitHubOrganization{}, err
	}
	defer resp.Body.Close()

	return parseGitHubOrgResponse(ctx, resp, orgName)
}

// makeGitHubOrgAPIRequest makes API request to GitHub for organization data (Pure Core)
func makeGitHubOrgAPIRequest(ctx *gofr.Context, orgName string) (*http.Response, error) {
	githubSvc := ctx.GetHTTPService("github")
	headers := buildGitHubRequestHeaders()
	ctx.Logger.Debugf("Making GitHub API request to fetch organization: %s", orgName)

	resp, err := githubSvc.GetWithHeaders(ctx, fmt.Sprintf("orgs/%s", orgName), nil, headers)
	if err != nil {
		ctx.Logger.Errorf("Failed to make GitHub API request for organization %s: %v", orgName, err)
		return nil, &gofrhttp.ErrorRequestTimeout{}
	}

	return resp, nil
}

// parseGitHubOrgResponse parses GitHub organization API response (Pure Core)
func parseGitHubOrgResponse(ctx *gofr.Context, resp *http.Response, orgName string) (GitHubOrganization, error) {
	ctx.Logger.Infof("GitHub API response status for organization %s: %d", orgName, resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		ctx.Logger.Warnf("Organization not found: %s", orgName)
		return GitHubOrganization{}, &gofrhttp.ErrorEntityNotFound{
			Name:  "organization",
			Value: orgName,
		}
	}

	if resp.StatusCode != http.StatusOK {
		ctx.Logger.Errorf("GitHub API returned error status %d for organization %s", resp.StatusCode, orgName)
		return GitHubOrganization{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"github_api_status", fmt.Sprintf("status_code_%d", resp.StatusCode)},
		}
	}

	var org GitHubOrganization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		ctx.Logger.Errorf("Failed to decode organization response for %s: %v", orgName, err)
		return GitHubOrganization{}, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_format", err.Error()},
		}
	}

	ctx.Logger.Infof("Successfully fetched organization: %s (ID: %d)", org.Login, org.ID)
	return org, nil
}

// fetchGitHubRepositoriesWithService fetches repositories using GoFr HTTP service
func fetchGitHubRepositoriesWithService(ctx *gofr.Context, orgName string, maxRepos int) ([]GitHubRepository, error) {
	ctx.Logger.Infof("Fetching repositories for organization: %s (max: %d)", orgName, maxRepos)

	if err := validateRepositoryParams(orgName, maxRepos); err != nil {
		return nil, err
	}

	githubSvc := ctx.GetHTTPService("github")
	allRepos, err := fetchAllRepositoryPages(ctx, githubSvc, orgName, maxRepos)
	if err != nil {
		return nil, err
	}

	ctx.Logger.Infof("Successfully fetched %d repositories for organization %s", len(allRepos), orgName)
	return allRepos, nil
}

// validateRepositoryParams validates input parameters for repository fetching
func validateRepositoryParams(orgName string, maxRepos int) error {
	if orgName == "" {
		return &gofrhttp.ErrorMissingParam{
			Params: []string{"organization_name"},
		}
	}

	if maxRepos <= 0 {
		return &gofrhttp.ErrorInvalidParam{
			Params: []string{"max_repos", fmt.Sprintf("%d", maxRepos)},
		}
	}

	return nil
}

// fetchAllRepositoryPages fetches all repository pages up to maxRepos
func fetchAllRepositoryPages(ctx *gofr.Context, githubSvc any, orgName string, maxRepos int) ([]GitHubRepository, error) {
	var allRepos []GitHubRepository
	page := 1
	perPage := 100

	for len(allRepos) < maxRepos {
		repos, shouldContinue, err := fetchRepositoryPage(ctx, githubSvc, orgName, page, perPage)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)
		page++

		if !shouldContinue {
			break
		}
	}

	return limitRepositories(ctx, allRepos, maxRepos, orgName), nil
}

// fetchRepositoryPage fetches a single page of repositories
func fetchRepositoryPage(ctx *gofr.Context, githubSvc any, orgName string, page, perPage int) ([]GitHubRepository, bool, error) {
	ctx.Logger.Debugf("Fetching repositories page %d for organization %s", page, orgName)

	resp, err := executeRepositoryRequest(ctx, githubSvc, orgName, page, perPage)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if err := validateRepositoryResponse(ctx, resp, orgName, page); err != nil {
		return nil, false, err
	}

	repos, err := decodeRepositoryResponse(ctx, resp, orgName)
	if err != nil {
		return nil, false, err
	}

	ctx.Logger.Infof("Fetched %d repositories from page %d for organization %s", len(repos), page, orgName)
	shouldContinue := len(repos) > 0 && len(repos) == perPage

	return repos, shouldContinue, nil
}

// executeRepositoryRequest executes a repository API request
func executeRepositoryRequest(ctx *gofr.Context, githubSvc any, orgName string, page, perPage int) (*http.Response, error) {
	query := map[string]any{
		"page":     fmt.Sprintf("%d", page),
		"per_page": fmt.Sprintf("%d", perPage),
		"sort":     "updated",
	}

	headers := buildGitHubRequestHeaders()
	
	// Get the GitHub service from context (same pattern as working organization request)
	githubHttpSvc := ctx.GetHTTPService("github")
	resp, err := githubHttpSvc.GetWithHeaders(ctx, fmt.Sprintf("orgs/%s/repos", orgName), query, headers)
	if err != nil {
		ctx.Logger.Errorf("Failed to fetch repositories for %s (page %d): %v", orgName, page, err)
		return nil, &gofrhttp.ErrorRequestTimeout{}
	}

	return resp, nil
}

// validateRepositoryResponse validates the HTTP response from GitHub API
func validateRepositoryResponse(ctx *gofr.Context, resp *http.Response, orgName string, page int) error {
	ctx.Logger.Debugf("Repository API response status for %s (page %d): %d", orgName, page, resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		ctx.Logger.Warnf("No repositories found for organization: %s", orgName)
		return &gofrhttp.ErrorEntityNotFound{
			Name:  "organization_repositories",
			Value: orgName,
		}
	}

	if resp.StatusCode != http.StatusOK {
		ctx.Logger.Errorf("GitHub API returned error status %d for repositories of %s", resp.StatusCode, orgName)
		return &gofrhttp.ErrorInvalidParam{
			Params: []string{"github_api_status", fmt.Sprintf("status_code_%d", resp.StatusCode)},
		}
	}

	return nil
}

// decodeRepositoryResponse decodes the JSON response into GitHubRepository slice
func decodeRepositoryResponse(ctx *gofr.Context, resp *http.Response, orgName string) ([]GitHubRepository, error) {
	var repos []GitHubRepository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		ctx.Logger.Errorf("Failed to decode repositories response for %s: %v", orgName, err)
		return nil, &gofrhttp.ErrorInvalidParam{
			Params: []string{"response_format", err.Error()},
		}
	}

	return repos, nil
}

// limitRepositories limits the number of repositories to maxRepos
func limitRepositories(ctx *gofr.Context, allRepos []GitHubRepository, maxRepos int, orgName string) []GitHubRepository {
	if len(allRepos) > maxRepos {
		ctx.Logger.Infof("Limiting repositories from %d to %d for organization %s", len(allRepos), maxRepos, orgName)
		return allRepos[:maxRepos]
	}
	return allRepos
}

// fetchGitHubTeamsWithService fetches teams using GoFr HTTP service
func fetchGitHubTeamsWithService(ctx *gofr.Context, orgName string, maxTeams int) ([]GitHubTeam, error) {
	if orgName == "" {
		return nil, &gofrhttp.ErrorMissingParam{
			Params: []string{"organization_name"},
		}
	}

	if maxTeams <= 0 {
		return nil, &gofrhttp.ErrorInvalidParam{
			Params: []string{"max_teams", fmt.Sprintf("%d", maxTeams)},
		}
	}

	githubSvc := ctx.GetHTTPService("github")

	var allTeams []GitHubTeam
	page := 1
	perPage := 100

	for len(allTeams) < maxTeams {
		query := map[string]any{
			"page":     fmt.Sprintf("%d", page),
			"per_page": fmt.Sprintf("%d", perPage),
		}

		headers := buildGitHubRequestHeaders()
		resp, err := githubSvc.GetWithHeaders(ctx, fmt.Sprintf("orgs/%s/teams", orgName), query, headers)
		if err != nil {
			return nil, &gofrhttp.ErrorRequestTimeout{}
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, &gofrhttp.ErrorEntityNotFound{
				Name:  "organization_teams",
				Value: orgName,
			}
		}

		if resp.StatusCode != http.StatusOK {
			return nil, &gofrhttp.ErrorInvalidParam{
				Params: []string{"github_api_status", fmt.Sprintf("status_code_%d", resp.StatusCode)},
			}
		}

		var teams []GitHubTeam
		if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
			return nil, &gofrhttp.ErrorInvalidParam{
				Params: []string{"response_format", err.Error()},
			}
		}

		if len(teams) == 0 {
			break
		}

		allTeams = append(allTeams, teams...)
		page++

		if len(teams) < perPage {
			break
		}
	}

	if len(allTeams) > maxTeams {
		allTeams = allTeams[:maxTeams]
	}

	return allTeams, nil
}

// fetchGitHubCodeownersWithService fetches CODEOWNERS file using GoFr HTTP service
func fetchGitHubCodeownersWithService(ctx *gofr.Context, owner, repo string) (GitHubCodeowners, error) {
	if owner == "" {
		return GitHubCodeowners{}, &gofrhttp.ErrorMissingParam{
			Params: []string{"owner"},
		}
	}

	if repo == "" {
		return GitHubCodeowners{}, &gofrhttp.ErrorMissingParam{
			Params: []string{"repository"},
		}
	}

	githubSvc := ctx.GetHTTPService("github")

	// Try different CODEOWNERS locations
	locations := []string{
		fmt.Sprintf("repos/%s/%s/contents/CODEOWNERS", owner, repo),
		fmt.Sprintf("repos/%s/%s/contents/.github/CODEOWNERS", owner, repo),
		fmt.Sprintf("repos/%s/%s/contents/docs/CODEOWNERS", owner, repo),
	}

	for _, location := range locations {
		headers := buildGitHubRequestHeaders()
		resp, err := githubSvc.GetWithHeaders(ctx, location, nil, headers)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var fileContent struct {
				Content string `json:"content"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
				continue
			}

			// Parse CODEOWNERS content
			rules := parseCodeownersContent(fileContent.Content)

			return GitHubCodeowners{
				Repository: fmt.Sprintf("%s/%s", owner, repo),
				Rules:      rules,
				Errors:     []GitHubCodeownersError{},
			}, nil
		}
	}

	// Return empty CODEOWNERS (not an error - many repos don't have CODEOWNERS)
	return GitHubCodeowners{
		Repository: fmt.Sprintf("%s/%s", owner, repo),
		Rules:      []GitHubCodeownersRule{},
		Errors:     []GitHubCodeownersError{},
	}, nil
}

// parseCodeownersContent parses base64-encoded CODEOWNERS content (Pure Core)
func parseCodeownersContent(base64Content string) []GitHubCodeownersRule {
	validateBase64ContentNotEmpty(base64Content)

	// Decode base64 content
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return []GitHubCodeownersRule{}
	}

	content := string(decodedBytes)
	if content == "" {
		return []GitHubCodeownersRule{}
	}

	// Split into lines and filter out comments and empty lines
	lines := lo.Filter(lo.Map(strings.Split(content, "\n"), func(line string, _ int) string {
		return strings.TrimSpace(line)
	}), func(line string, _ int) bool {
		return line != "" && !strings.HasPrefix(line, "#")
	})

	// Parse each line into a rule
	rules := lo.Map(lines, func(line string, index int) GitHubCodeownersRule {
		parts := lo.Filter(strings.Fields(line), func(part string, _ int) bool {
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

	return rules
}

// validateBase64ContentNotEmpty validates base64 content is not empty (Pure Core)
func validateBase64ContentNotEmpty(content string) {
	if content == "" {
		panic("Base64 content cannot be empty")
	}
}

// buildGitHubRequestHeaders builds headers for GitHub API requests (Pure Core)
func buildGitHubRequestHeaders() map[string]string {
	// Note: In a real implementation, we would get the token from configuration
	// For now, we'll use a placeholder that expects GITHUB_TOKEN environment variable
	token := os.Getenv("GITHUB_TOKEN")

	headers := map[string]string{
		"Accept":     "application/vnd.github.v3+json",
		"User-Agent": "overseer-codeowners-scanner/1.0",
	}

	if token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", token)
	}

	return headers
}
