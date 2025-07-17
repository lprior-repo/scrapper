package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGitHubTokenValidation tests the token validation logic
func TestGitHubTokenValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		token     string
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "valid classic token",
			token:     "ghp_1234567890123456789012345678901234567890",
			shouldErr: false,
		},
		{
			name:      "valid oauth token",
			token:     "gho_1234567890123456789012345678901234567890",
			shouldErr: false,
		},
		{
			name:      "empty token",
			token:     "",
			shouldErr: true,
			errMsg:    "token cannot be empty",
		},
		{
			name:      "token too short",
			token:     "ghp_123",
			shouldErr: true,
			errMsg:    "too short",
		},
		{
			name:      "invalid prefix",
			token:     "invalid_1234567890123456789012345678901234567890",
			shouldErr: true,
			errMsg:    "valid prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateGitHubToken(tt.token)

			if tt.shouldErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGitHubOrgValidation tests the organization name validation
func TestGitHubOrgValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		org       string
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "valid org name",
			org:       "microsoft",
			shouldErr: false,
		},
		{
			name:      "valid org with hyphen",
			org:       "my-org",
			shouldErr: false,
		},
		{
			name:      "valid org with numbers",
			org:       "org123",
			shouldErr: false,
		},
		{
			name:      "empty org",
			org:       "",
			shouldErr: true,
			errMsg:    "cannot be empty",
		},
		{
			name:      "whitespace only",
			org:       "   ",
			shouldErr: true,
			errMsg:    "whitespace only",
		},
		{
			name:      "starts with hyphen",
			org:       "-invalid",
			shouldErr: true,
			errMsg:    "start or end with hyphen",
		},
		{
			name:      "ends with hyphen",
			org:       "invalid-",
			shouldErr: true,
			errMsg:    "start or end with hyphen",
		},
		{
			name:      "too long",
			org:       "this-organization-name-is-way-too-long-to-be-valid",
			shouldErr: true,
			errMsg:    "too long",
		},
		{
			name:      "invalid characters",
			org:       "org@invalid",
			shouldErr: true,
			errMsg:    "invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateGitHubOrgName(tt.org)

			if tt.shouldErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTokenMasking tests the token masking functionality
func TestTokenMasking(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "long token",
			token:    "ghp_1234567890123456789012345678901234567890",
			expected: "ghp_****7890",
		},
		{
			name:     "short token",
			token:    "short",
			expected: "****",
		},
		{
			name:     "empty token",
			token:    "",
			expected: "",
		},
		{
			name:     "exactly 8 chars",
			token:    "12345678",
			expected: "****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := maskToken(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGitHubScannerIntegration tests the full scanner integration (requires environment variables)
func TestGitHubScannerIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if no GitHub token available
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("Skipping integration test - GITHUB_TOKEN not set")
	}

	if os.Getenv("GITHUB_ORG") == "" {
		t.Skip("Skipping integration test - GITHUB_ORG not set")
	}

	t.Run("full_scan_integration", func(t *testing.T) {
		// Set limits for testing
		os.Setenv("GITHUB_MAX_REPOS", "5")
		os.Setenv("GITHUB_MAX_TEAMS", "3")
		defer func() {
			os.Unsetenv("GITHUB_MAX_REPOS")
			os.Unsetenv("GITHUB_MAX_TEAMS")
		}()

		// Create configuration
		config, err := parseGitHubScanConfig()
		require.NoError(t, err)

		// Validate configuration
		err = validateScanConfig(config)
		require.NoError(t, err)

		// Run scan with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result := scanGitHubOrganization(ctx, config)

		// Verify results
		assert.True(t, result.Success, "Scan should succeed, error: %s", result.Error)
		assert.NotEmpty(t, result.Data.Organization, "Organization should be populated")
		assert.GreaterOrEqual(t, len(result.Data.Repos), 1, "Should find at least one repository")
		assert.LessOrEqual(t, result.Summary.APICallsUsed, 50, "Should use ≤50 API calls")
		assert.GreaterOrEqual(t, result.Summary.TotalRepos, 1, "Should count repositories")

		// Verify API call efficiency
		expectedCalls := calculateAPICallsNeeded(len(result.Data.Repos), len(result.Data.Teams))
		assert.LessOrEqual(t, result.Summary.APICallsUsed, expectedCalls+1, "API calls should be within expected range")

		t.Logf("✅ Integration test passed:")
		t.Logf("   Organization: %s", result.Data.Organization)
		t.Logf("   Repositories: %d", result.Summary.TotalRepos)
		t.Logf("   Teams: %d", result.Summary.TotalTeams)
		t.Logf("   API calls used: %d/50", result.Summary.APICallsUsed)
		t.Logf("   Repos with CODEOWNERS: %d", result.Summary.ReposWithCodeowners)
	})
}

// TestGitHubClientCreation tests the GitHub client creation
func TestGitHubClientCreation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("valid_client_creation", func(t *testing.T) {
		t.Parallel()

		client, err := createGitHubClient(ctx, "ghp_validtoken1234567890123456789012345678901234567890", "testorg")

		assert.NoError(t, err)
		assert.NotNil(t, client.GraphQL)
		assert.Equal(t, "testorg", client.Org)
	})

	t.Run("empty_token", func(t *testing.T) {
		t.Parallel()

		_, err := createGitHubClient(ctx, "", "testorg")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is required")
	})

	t.Run("empty_org", func(t *testing.T) {
		t.Parallel()

		_, err := createGitHubClient(ctx, "ghp_validtoken1234567890123456789012345678901234567890", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "organization name is required")
	})
}

// BenchmarkCodeownersProcessing benchmarks the CODEOWNERS processing
func BenchmarkCodeownersProcessing(b *testing.B) {
	codeownersContent := `
# Global owners
* @global-team

# Frontend
*.js @frontend-team
*.jsx @frontend-team
*.ts @typescript-team
*.tsx @typescript-team

# Backend
*.go @backend-team
*.py @python-team
*.java @java-team

# DevOps
Dockerfile @devops-team
*.yml @devops-team
*.yaml @devops-team
.github/ @devops-team

# Documentation
*.md @docs-team
docs/ @docs-team

# Specific directories
/src/api/ @api-team
/src/web/ @web-team
/tests/ @qa-team
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entries := parseCodeownersContent(codeownersContent)
		_ = extractUniqueOwners(entries)
	}
}
