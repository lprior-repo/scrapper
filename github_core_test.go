package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateBatchRequest(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name      string
		request   BatchRequest
		shouldErr bool
		errMsg    string
	}{
		{
			name: "valid request",
			request: BatchRequest{
				Organization: "testorg",
				MaxRepos:     100,
				MaxTeams:     50,
			},
			shouldErr: false,
		},
		{
			name: "empty organization",
			request: BatchRequest{
				Organization: "",
				MaxRepos:     100,
				MaxTeams:     50,
			},
			shouldErr: true,
			errMsg:    "organization name is required",
		},
		{
			name: "whitespace organization",
			request: BatchRequest{
				Organization: "   ",
				MaxRepos:     100,
				MaxTeams:     50,
			},
			shouldErr: true,
			errMsg:    "organization name cannot be empty or whitespace",
		},
		{
			name: "negative max repos",
			request: BatchRequest{
				Organization: "testorg",
				MaxRepos:     -1,
				MaxTeams:     50,
			},
			shouldErr: true,
			errMsg:    "max repos cannot be negative",
		},
		{
			name: "negative max teams",
			request: BatchRequest{
				Organization: "testorg",
				MaxRepos:     100,
				MaxTeams:     -1,
			},
			shouldErr: true,
			errMsg:    "max teams cannot be negative",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			err := validateBatchRequest(tt.request)
			
			if tt.shouldErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseCodeownersContent(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name     string
		content  string
		expected []CodeownersEntry
	}{
		{
			name:     "empty content",
			content:  "",
			expected: []CodeownersEntry{},
		},
		{
			name:     "single entry",
			content:  "* @team1",
			expected: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1"}},
			},
		},
		{
			name:     "multiple entries",
			content:  "* @team1\n*.js @frontend\n*.go @backend",
			expected: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1"}},
				{Pattern: "*.js", Owners: []string{"@frontend"}},
				{Pattern: "*.go", Owners: []string{"@backend"}},
			},
		},
		{
			name:     "with comments",
			content:  "# This is a comment\n* @team1\n# Another comment\n*.js @frontend",
			expected: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1"}},
				{Pattern: "*.js", Owners: []string{"@frontend"}},
			},
		},
		{
			name:     "multiple owners",
			content:  "* @team1 @team2 @user1",
			expected: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1", "@team2", "@user1"}},
			},
		},
		{
			name:     "empty lines",
			content:  "* @team1\n\n*.js @frontend\n\n",
			expected: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1"}},
				{Pattern: "*.js", Owners: []string{"@frontend"}},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			result := parseCodeownersContent(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractUniqueOwners(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name     string
		entries  []CodeownersEntry
		expected []string
	}{
		{
			name:     "empty entries",
			entries:  []CodeownersEntry{},
			expected: []string{},
		},
		{
			name: "single entry",
			entries: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1"}},
			},
			expected: []string{"@team1"},
		},
		{
			name: "multiple entries with duplicates",
			entries: []CodeownersEntry{
				{Pattern: "*", Owners: []string{"@team1", "@team2"}},
				{Pattern: "*.js", Owners: []string{"@team1", "@frontend"}},
			},
			expected: []string{"@team1", "@team2", "@frontend"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			result := extractUniqueOwners(tt.entries)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestCalculateAPICallsNeeded(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name      string
		repoCount int
		teamCount int
		expected  int
	}{
		{
			name:      "small organization",
			repoCount: 10,
			teamCount: 5,
			expected:  1,
		},
		{
			name:      "medium organization",
			repoCount: 150,
			teamCount: 25,
			expected:  2,
		},
		{
			name:      "large organization",
			repoCount: 500,
			teamCount: 200,
			expected:  5,
		},
		{
			name:      "zero counts",
			repoCount: 0,
			teamCount: 0,
			expected:  1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			result := calculateAPICallsNeeded(tt.repoCount, tt.teamCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildGraphQLOrgQuery(t *testing.T) {
	t.Parallel()
	
	request := BatchRequest{
		Organization: "testorg",
		MaxRepos:     50,
		MaxTeams:     25,
	}
	
	query := buildGraphQLOrgQuery(request)
	
	// Verify query contains expected elements
	assert.Contains(t, query, "organization(login: $org)")
	assert.Contains(t, query, "repositories(first: 50")
	assert.Contains(t, query, "teams(first: 25")
	assert.Contains(t, query, "codeowners: object(expression: \"HEAD:CODEOWNERS\")")
	assert.Contains(t, query, "docsCodeowners: object(expression: \"HEAD:.github/CODEOWNERS\")")
	assert.Contains(t, query, "rootCodeowners: object(expression: \"HEAD:docs/CODEOWNERS\")")
}