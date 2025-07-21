package main

// ScanRequest represents a request to scan a GitHub organization
type ScanRequest struct {
	Organization string `json:"organization"`
	MaxRepos     int    `json:"max_repos"`
	MaxTeams     int    `json:"max_teams"`
	UseTopics    bool   `json:"use_topics"`
}

// ScanResponse represents the response from scanning an organization
type ScanResponse struct {
	Success      bool                   `json:"success"`
	Organization string                 `json:"organization"`
	Summary      ScanSummary            `json:"summary"`
	Errors       []string               `json:"errors"`
	Data         map[string]interface{} `json:"data"`
}

// ScanSummary represents scan statistics
type ScanSummary struct {
	TotalRepos          int      `json:"total_repos"`
	ReposWithCodeowners int      `json:"repos_with_codeowners"`
	TotalTeams          int      `json:"total_teams"`
	TotalTopics         int      `json:"total_topics"`
	UniqueOwners        []string `json:"unique_owners"`
	APICallsUsed        int      `json:"api_calls_used"`
	ProcessingTimeMs    int64    `json:"processing_time_ms"`
}

// GraphResponse represents graph visualization data
type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode represents a node in the graph
type GraphNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Label    string                 `json:"label"`
	Data     map[string]interface{} `json:"data"`
	Position GraphPosition          `json:"position"`
}

// GraphEdge represents an edge in the graph
type GraphEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Label  string `json:"label"`
}

// GraphPosition represents node position in the graph
type GraphPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// StatsResponse represents organization statistics
type StatsResponse struct {
	Organization      string `json:"organization"`
	TotalRepositories int    `json:"total_repositories"`
	TotalTeams        int    `json:"total_teams"`
	TotalTopics       int    `json:"total_topics"`
	TotalUsers        int    `json:"total_users"`
	TotalCodeowners   int    `json:"total_codeowners"`
	CodeownerCoverage string `json:"codeowner_coverage"`
	LastScanTime      string `json:"last_scan_time"`
}

// AppDependencies represents application dependencies
type AppDependencies struct {
	Config    AppConfig
	Neo4jConn *Neo4jConnection
}

// AppHandler contains the application dependencies
type AppHandler struct {
	deps *AppDependencies
}