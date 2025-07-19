package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// buildGraphNodesQuery builds a query to fetch graph nodes (Pure Core)
func buildGraphNodesQuery(orgName string, useTopics bool) string {
	validateOrgNameNotEmpty(orgName)

	if useTopics {
		return `
			MATCH (org:Organization {login: $orgName})
			OPTIONAL MATCH (org)-[:OWNS]->(repo:Repository)
			OPTIONAL MATCH (org)-[:HAS_TOPIC]->(topic:Topic)
			OPTIONAL MATCH (repo)-[:HAS_CODEOWNER]->(user:User)
			WITH org,
				 COLLECT(DISTINCT {
					 id: repo.id,
					 type: 'repository',
					 label: repo.name,
					 data: {
						 name: repo.name,
						 fullName: repo.full_name,
						 description: repo.description,
						 private: repo.private,
						 url: repo.url,
						 createdAt: repo.created_at,
						 updatedAt: repo.updated_at
					 }
				 }) AS repos,
				 COLLECT(DISTINCT {
					 id: topic.name,
					 type: 'topic',
					 label: topic.name,
					 data: {
						 name: topic.name,
						 count: topic.count
					 }
				 }) AS topics,
				 COLLECT(DISTINCT {
					 id: user.id,
					 type: 'user',
					 label: user.login,
					 data: {
						 login: user.login,
						 name: user.name,
						 email: user.email,
						 url: user.url
					 }
				 }) AS users
			RETURN {
				id: org.id,
				type: 'organization',
				label: org.name,
				data: {
					login: org.login,
					name: org.name,
					description: org.description,
					email: org.email,
					url: org.url,
					createdAt: org.created_at,
					updatedAt: org.updated_at
				}
			} AS org_node,
			repos,
			[] AS teams,
			topics,
			users
		`
	} else {
		return `
			MATCH (org:Organization {login: $orgName})
			OPTIONAL MATCH (org)-[:OWNS]->(repo:Repository)
			OPTIONAL MATCH (org)-[:HAS_TEAM]->(team:Team)
			OPTIONAL MATCH (repo)-[:HAS_CODEOWNER]->(user:User)
			OPTIONAL MATCH (repo)-[:HAS_TEAM_OWNER]->(team)
			WITH org,
				 COLLECT(DISTINCT {
					 id: repo.id,
					 type: 'repository',
					 label: repo.name,
					 data: {
						 name: repo.name,
						 fullName: repo.full_name,
						 description: repo.description,
						 private: repo.private,
						 url: repo.url,
						 createdAt: repo.created_at,
						 updatedAt: repo.updated_at
					 }
				 }) AS repos,
				 COLLECT(DISTINCT {
					 id: team.id,
					 type: 'team',
					 label: team.name,
					 data: {
						 name: team.name,
						 slug: team.slug,
						 description: team.description,
						 url: team.url
					 }
				 }) AS teams,
				 COLLECT(DISTINCT {
					 id: user.id,
					 type: 'user',
					 label: user.login,
					 data: {
						 login: user.login,
						 name: user.name,
						 email: user.email,
						 url: user.url
					 }
				 }) AS users
			RETURN {
				id: org.id,
				type: 'organization',
				label: org.name,
				data: {
					login: org.login,
					name: org.name,
					description: org.description,
					email: org.email,
					url: org.url,
					createdAt: org.created_at,
					updatedAt: org.updated_at
				}
			} AS org_node,
			repos,
			teams,
			[] AS topics,
			users
		`
	}
}

// buildGraphEdgesQuery builds a query to fetch graph edges (Pure Core)
func buildGraphEdgesQuery(orgName string, useTopics bool) string {
	validateOrgNameNotEmpty(orgName)

	if useTopics {
		return `
			MATCH (org:Organization {login: $orgName})
			OPTIONAL MATCH (org)-[:OWNS]->(repo:Repository)
			OPTIONAL MATCH (org)-[:HAS_TOPIC]->(topic:Topic)
			OPTIONAL MATCH (repo)-[:HAS_TOPIC]->(repo_topic:Topic)
			OPTIONAL MATCH (repo)-[:HAS_CODEOWNER]->(user:User)
			WITH org,
				 COLLECT(DISTINCT {
					 id: 'owns-' + org.id + '-' + repo.id,
					 source: org.id,
					 target: repo.id,
					 type: 'owns',
					 label: 'owns'
				 }) AS owns_edges,
				 COLLECT(DISTINCT {
					 id: 'has-topic-' + org.id + '-' + topic.name,
					 source: org.id,
					 target: topic.name,
					 type: 'has_topic',
					 label: 'has topic'
				 }) AS topic_edges,
				 COLLECT(DISTINCT {
					 id: 'repo-topic-' + repo.id + '-' + repo_topic.name,
					 source: repo.id,
					 target: repo_topic.name,
					 type: 'repo_topic',
					 label: 'uses topic'
				 }) AS repo_topic_edges,
				 COLLECT(DISTINCT {
					 id: 'codeowner-' + repo.id + '-' + user.id,
					 source: repo.id,
					 target: user.id,
					 type: 'codeowner',
					 label: 'code owner'
				 }) AS codeowner_edges
			RETURN owns_edges + topic_edges + repo_topic_edges + codeowner_edges AS edges
		`
	} else {
		return `
			MATCH (org:Organization {login: $orgName})
			OPTIONAL MATCH (org)-[:OWNS]->(repo:Repository)
			OPTIONAL MATCH (org)-[:HAS_TEAM]->(team:Team)
			OPTIONAL MATCH (repo)-[:HAS_CODEOWNER]->(user:User)
			OPTIONAL MATCH (repo)-[:HAS_TEAM_OWNER]->(team)
			WITH org,
				 COLLECT(DISTINCT {
					 id: 'owns-' + org.id + '-' + repo.id,
					 source: org.id,
					 target: repo.id,
					 type: 'owns',
					 label: 'owns'
				 }) AS owns_edges,
				 COLLECT(DISTINCT {
					 id: 'has-team-' + org.id + '-' + team.id,
					 source: org.id,
					 target: team.id,
					 type: 'has_team',
					 label: 'has team'
				 }) AS team_edges,
				 COLLECT(DISTINCT {
					 id: 'codeowner-' + repo.id + '-' + user.id,
					 source: repo.id,
					 target: user.id,
					 type: 'codeowner',
					 label: 'code owner'
				 }) AS codeowner_edges,
				 COLLECT(DISTINCT {
					 id: 'team-owner-' + repo.id + '-' + team.id,
					 source: repo.id,
					 target: team.id,
					 type: 'team_owner',
					 label: 'team owner'
				 }) AS team_owner_edges
			RETURN owns_edges + team_edges + codeowner_edges + team_owner_edges AS edges
		`
	}
}

// buildStatsQuery builds a query to fetch organization statistics (Pure Core)
func buildStatsQuery(orgName string) string {
	validateOrgNameNotEmpty(orgName)

	return `
		MATCH (org:Organization {login: $orgName})
		OPTIONAL MATCH (org)-[:OWNS]->(repo:Repository)
		OPTIONAL MATCH (org)-[:HAS_TEAM]->(team:Team)
		OPTIONAL MATCH (org)-[:HAS_TOPIC]->(topic:Topic)
		OPTIONAL MATCH (repo)-[:HAS_CODEOWNER]->(user:User)
		OPTIONAL MATCH (repo)-[:HAS_TEAM_OWNER]->(team_owner:Team)
		WITH org,
			 COUNT(DISTINCT repo) AS total_repos,
			 COUNT(DISTINCT team) AS total_teams,
			 COUNT(DISTINCT topic) AS total_topics,
			 COUNT(DISTINCT user) AS total_users,
			 SIZE([r IN collect(DISTINCT repo) WHERE EXISTS((r)-[:HAS_CODEOWNER]->()) OR EXISTS((r)-[:HAS_TEAM_OWNER]->())]) AS repos_with_codeowners
		RETURN {
			organization: org.login,
			total_repositories: total_repos,
			total_teams: total_teams,
			total_topics: total_topics,
			total_users: total_users,
			total_codeowners: repos_with_codeowners,
			codeowner_coverage: CASE
				WHEN total_repos > 0 THEN toString(round(100.0 * repos_with_codeowners / total_repos)) + '%'
				ELSE '0%'
			END,
			last_scan_time: org.updated_at
		} AS stats
	`
}

// buildCreateOrganizationQuery builds a query to create/update an organization (Pure Core)
func buildCreateOrganizationQuery() string {
	return `
		MERGE (org:Organization {login: $login})
		SET org.id = $id,
			org.name = $name,
			org.description = $description,
			org.email = $email,
			org.url = $url,
			org.created_at = $created_at,
			org.updated_at = $updated_at
		RETURN org
	`
}

// buildCreateRepositoryQuery builds a query to create/update a repository (Pure Core)
func buildCreateRepositoryQuery() string {
	return `
		MERGE (repo:Repository {full_name: $full_name})
		SET repo.id = $id,
			repo.name = $name,
			repo.description = $description,
			repo.private = $private,
			repo.url = $url,
			repo.created_at = $created_at,
			repo.updated_at = $updated_at
		WITH repo
		MATCH (org:Organization {login: $org_login})
		MERGE (org)-[:OWNS]->(repo)
		RETURN repo
	`
}

// buildCreateRepositoryTopicRelationshipQuery builds a query to create repository-topic relationships (Pure Core)
func buildCreateRepositoryTopicRelationshipQuery() string {
	return `
		MATCH (repo:Repository {full_name: $repo_full_name})
		MATCH (topic:Topic {name: $topic_name})
		MERGE (repo)-[:HAS_TOPIC]->(topic)
		RETURN repo, topic
	`
}

// buildCreateTeamQuery builds a query to create/update a team (Pure Core)
func buildCreateTeamQuery() string {
	return `
		MERGE (team:Team {slug: $slug})
		SET team.id = $id,
			team.name = $name,
			team.description = $description,
			team.url = $url
		WITH team
		MATCH (org:Organization {login: $org_login})
		MERGE (org)-[:HAS_TEAM]->(team)
		RETURN team
	`
}

// buildCreateTopicQuery builds a query to create/update a topic (Pure Core)
func buildCreateTopicQuery() string {
	return `
		MERGE (topic:Topic {name: $name})
		SET topic.count = $count
		WITH topic
		MATCH (org:Organization {login: $org_login})
		MERGE (org)-[:HAS_TOPIC]->(topic)
		RETURN topic
	`
}

// buildCreateUserQuery builds a query to create/update a user (Pure Core)
func buildCreateUserQuery() string {
	return `
		MERGE (user:User {login: $login})
		SET user.id = $id,
			user.name = $name,
			user.email = CASE 
				WHEN $email = '' THEN NULL
				ELSE $email
			END,
			user.url = $url
		RETURN user
	`
}

// buildCreateCodeownerRelationshipQuery builds a query to create codeowner relationships (Pure Core)
func buildCreateCodeownerRelationshipQuery() string {
	return `
		MATCH (repo:Repository {full_name: $repo_full_name})
		MATCH (owner:User {login: $owner_login})
		MERGE (repo)-[r:HAS_CODEOWNER]->(owner)
		SET r.pattern = $pattern,
			r.line = $line
		RETURN r
	`
}

// buildCreateTeamCodeownerRelationshipQuery builds a query to create team codeowner relationships (Pure Core)
func buildCreateTeamCodeownerRelationshipQuery() string {
	return `
		MATCH (repo:Repository {full_name: $repo_full_name})
		MATCH (team:Team {slug: $team_slug})
		MERGE (repo)-[r:HAS_TEAM_OWNER]->(team)
		SET r.pattern = $pattern,
			r.line = $line
		RETURN r
	`
}


// storeOrganization stores organization data in Neo4j (Orchestrator)
func storeOrganization(ctx context.Context, session *Neo4jSession, org GitHubOrganization) error {
	validateNeo4jSessionNotNil(session)

	query := buildCreateOrganizationQuery()
	params := map[string]interface{}{
		"id":          org.ID,
		"login":       org.Login,
		"name":        org.Name,
		"description": org.Description,
		"email":       org.Email,
		"url":         org.URL,
		"created_at":  org.CreatedAt.Format(time.RFC3339),
		"updated_at":  org.UpdatedAt.Format(time.RFC3339),
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store organization: %w", err)
	}

	return nil
}

// storeRepository stores repository data in Neo4j (Orchestrator)
func storeRepository(ctx context.Context, session *Neo4jSession, repo GitHubRepository, orgLogin string) error {
	validateNeo4jSessionNotNil(session)
	validateOrgLoginNotEmpty(orgLogin)

	query := buildCreateRepositoryQuery()
	params := map[string]interface{}{
		"id":          repo.ID,
		"name":        repo.Name,
		"full_name":   repo.FullName,
		"description": repo.Description,
		"private":     repo.Private,
		"url":         repo.URL,
		"created_at":  repo.CreatedAt.Format(time.RFC3339),
		"updated_at":  repo.UpdatedAt.Format(time.RFC3339),
		"org_login":   orgLogin,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store repository: %w", err)
	}

	// Create relationships between repository and its topics
	for _, topic := range repo.Topics {
		if err := storeRepositoryTopicRelationship(ctx, session, repo.FullName, topic); err != nil {
			return fmt.Errorf("failed to store repository-topic relationship: %w", err)
		}
	}

	return nil
}

// storeRepositoryTopicRelationship stores a relationship between a repository and a topic (Orchestrator)
func storeRepositoryTopicRelationship(ctx context.Context, session *Neo4jSession, repoFullName, topicName string) error {
	validateNeo4jSessionNotNil(session)
	validateRepoFullNameNotEmpty(repoFullName)
	validateTopicNameNotEmpty(topicName)

	query := buildCreateRepositoryTopicRelationshipQuery()
	params := map[string]interface{}{
		"repo_full_name": repoFullName,
		"topic_name":     topicName,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store repository-topic relationship: %w", err)
	}

	return nil
}

// storeTeam stores team data in Neo4j (Orchestrator)
func storeTeam(ctx context.Context, session *Neo4jSession, team GitHubTeam, orgLogin string) error {
	validateNeo4jSessionNotNil(session)
	validateOrgLoginNotEmpty(orgLogin)

	query := buildCreateTeamQuery()
	params := map[string]interface{}{
		"id":          team.ID,
		"slug":        team.Slug,
		"name":        team.Name,
		"description": team.Description,
		"url":         team.URL,
		"org_login":   orgLogin,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store team: %w", err)
	}

	return nil
}

// storeTopic stores topic data in Neo4j (Orchestrator)
func storeTopic(ctx context.Context, session *Neo4jSession, topic GitHubTopic, orgLogin string) error {
	validateNeo4jSessionNotNil(session)
	validateOrgLoginNotEmpty(orgLogin)

	query := buildCreateTopicQuery()
	params := map[string]interface{}{
		"name":      topic.Name,
		"count":     topic.Count,
		"org_login": orgLogin,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store topic: %w", err)
	}

	return nil
}

// storeUser stores user data in Neo4j (Orchestrator)
func storeUser(ctx context.Context, session *Neo4jSession, user GitHubUser) error {
	validateNeo4jSessionNotNil(session)

	query := buildCreateUserQuery()
	params := map[string]interface{}{
		"id":    user.ID,
		"login": user.Login,
		"name":  user.Name,
		"email": user.Email,
		"url":   user.URL,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store user: %w", err)
	}

	return nil
}

// storeCodeowners stores CODEOWNERS data in Neo4j (Orchestrator)
func storeCodeowners(ctx context.Context, session *Neo4jSession, codeowners GitHubCodeowners, orgLogin string) error {
	validateNeo4jSessionNotNil(session)
	validateOrgLoginNotEmpty(orgLogin)

	for _, rule := range codeowners.Rules {
		for _, owner := range rule.Owners {
			if err := storeCodeownerRule(ctx, session, codeowners.Repository, owner, rule.Pattern, rule.Line); err != nil {
				return fmt.Errorf("failed to store codeowner rule: %w", err)
			}
		}
	}

	return nil
}

// storeCodeownerRule stores a single codeowner rule in Neo4j (Orchestrator)
func storeCodeownerRule(ctx context.Context, session *Neo4jSession, repoFullName, owner, pattern string, line int) error {
	validateNeo4jSessionNotNil(session)
	validateRepoFullNameNotEmpty(repoFullName)
	validateOwnerNotEmpty(owner)

	if isTeamOwner(owner) {
		return storeTeamCodeownerRule(ctx, session, repoFullName, owner, pattern, line)
	}

	return storeUserCodeownerRule(ctx, session, repoFullName, owner, pattern, line)
}

// storeUserCodeownerRule stores a user codeowner rule in Neo4j (Orchestrator)
func storeUserCodeownerRule(ctx context.Context, session *Neo4jSession, repoFullName, userLogin, pattern string, line int) error {
	validateNeo4jSessionNotNil(session)
	validateRepoFullNameNotEmpty(repoFullName)

	// Clean user login (remove @ prefix)
	cleanUserLogin := strings.TrimPrefix(userLogin, "@")

	// First, ensure the user exists
	user := GitHubUser{
		ID:    generateUserID(cleanUserLogin),
		Login: cleanUserLogin,
		Name:  cleanUserLogin,
		Email: "",
		URL:   fmt.Sprintf("https://github.com/%s", cleanUserLogin),
	}

	if err := storeUser(ctx, session, user); err != nil {
		return fmt.Errorf("failed to store user: %w", err)
	}

	// Create the codeowner relationship
	query := buildCreateCodeownerRelationshipQuery()
	params := map[string]interface{}{
		"repo_full_name": repoFullName,
		"owner_login":    cleanUserLogin,
		"pattern":        pattern,
		"line":           line,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store user codeowner relationship: %w", err)
	}

	return nil
}

// storeTeamCodeownerRule stores a team codeowner rule in Neo4j (Orchestrator)
func storeTeamCodeownerRule(ctx context.Context, session *Neo4jSession, repoFullName, teamSlug, pattern string, line int) error {
	validateNeo4jSessionNotNil(session)
	validateRepoFullNameNotEmpty(repoFullName)

	// Clean team slug (remove @org/ prefix)
	cleanTeamSlug := extractTeamSlug(teamSlug)

	query := buildCreateTeamCodeownerRelationshipQuery()
	params := map[string]interface{}{
		"repo_full_name": repoFullName,
		"team_slug":      cleanTeamSlug,
		"pattern":        pattern,
		"line":           line,
	}

	_, err := executeNeo4jWrite(ctx, session, query, params)
	if err != nil {
		return fmt.Errorf("failed to store team codeowner relationship: %w", err)
	}

	return nil
}

// convertToGraphNodes converts Neo4j records to graph nodes (Pure Core)
func convertToGraphNodes(records []map[string]interface{}) []GraphNode {
	if len(records) == 0 {
		return []GraphNode{}
	}

	record := records[0]
	var nodes []GraphNode

	nodes = append(nodes, extractOrganizationNode(record)...)
	nodes = append(nodes, extractRepositoryNodes(record)...)
	nodes = append(nodes, extractTeamNodes(record)...)
	nodes = append(nodes, extractTopicNodes(record)...)
	nodes = append(nodes, extractUserNodes(record)...)

	return nodes
}

func extractOrganizationNode(record map[string]interface{}) []GraphNode {
	orgNode, exists := record["org_node"]
	if !exists {
		return []GraphNode{}
	}

	orgMap, ok := orgNode.(map[string]interface{})
	if !ok {
		return []GraphNode{}
	}

	return []GraphNode{convertMapToGraphNode(orgMap, 0, 0)}
}

func extractRepositoryNodes(record map[string]interface{}) []GraphNode {
	repos, exists := record["repos"]
	if !exists {
		return []GraphNode{}
	}

	repoList, ok := repos.([]interface{})
	if !ok {
		return []GraphNode{}
	}

	return convertListToGraphNodes(repoList, 200, 200)
}

func extractTeamNodes(record map[string]interface{}) []GraphNode {
	teams, exists := record["teams"]
	if !exists {
		return []GraphNode{}
	}

	teamList, ok := teams.([]interface{})
	if !ok {
		return []GraphNode{}
	}

	return convertListToGraphNodes(teamList, 400, 200)
}

func extractTopicNodes(record map[string]interface{}) []GraphNode {
	topics, exists := record["topics"]
	if !exists {
		return []GraphNode{}
	}

	topicList, ok := topics.([]interface{})
	if !ok {
		return []GraphNode{}
	}

	return convertListToGraphNodes(topicList, 500, 200)
}

func extractUserNodes(record map[string]interface{}) []GraphNode {
	users, exists := record["users"]
	if !exists {
		return []GraphNode{}
	}

	userList, ok := users.([]interface{})
	if !ok {
		return []GraphNode{}
	}

	return convertListToGraphNodes(userList, 600, 200)
}

func convertListToGraphNodes(list []interface{}, yOffset, xSpacing float64) []GraphNode {
	var nodes []GraphNode

	for i, item := range list {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		node := convertMapToGraphNode(itemMap, float64(i*int(xSpacing)), yOffset)
		nodes = append(nodes, node)
	}

	return nodes
}

// convertToGraphEdges converts Neo4j records to graph edges (Pure Core)
func convertToGraphEdges(records []map[string]interface{}) []GraphEdge {
	if len(records) == 0 {
		return []GraphEdge{}
	}

	record := records[0]
	var edges []GraphEdge

	if edgeList, exists := record["edges"]; exists {
		if edgesArray, ok := edgeList.([]interface{}); ok {
			for _, edge := range edgesArray {
				if edgeMap, ok := edge.(map[string]interface{}); ok {
					edges = append(edges, convertMapToGraphEdge(edgeMap))
				}
			}
		}
	}

	return edges
}

// convertToStatsResponse converts Neo4j record to stats response (Pure Core)
func convertToStatsResponse(record map[string]interface{}, orgName string) StatsResponse {
	stats, exists := record["stats"]
	if !exists {
		return StatsResponse{Organization: orgName}
	}

	statsMap, ok := stats.(map[string]interface{})
	if !ok {
		return StatsResponse{Organization: orgName}
	}

	return StatsResponse{
		Organization:      getStringFromMap(statsMap, "organization"),
		TotalRepositories: getIntFromMap(statsMap, "total_repositories"),
		TotalTeams:        getIntFromMap(statsMap, "total_teams"),
		TotalTopics:       getIntFromMap(statsMap, "total_topics"),
		TotalUsers:        getIntFromMap(statsMap, "total_users"),
		TotalCodeowners:   getIntFromMap(statsMap, "total_codeowners"),
		CodeownerCoverage: getStringFromMap(statsMap, "codeowner_coverage"),
		LastScanTime:      getStringFromMap(statsMap, "last_scan_time"),
	}
}

// convertMapToGraphNode converts a map to a graph node (Pure Core)
func convertMapToGraphNode(nodeMap map[string]interface{}, x, y float64) GraphNode {
	return GraphNode{
		ID:    getStringFromMap(nodeMap, "id"),
		Type:  getStringFromMap(nodeMap, "type"),
		Label: getStringFromMap(nodeMap, "label"),
		Data:  getMapFromMap(nodeMap, "data"),
		Position: GraphPosition{
			X: x,
			Y: y,
		},
	}
}

// convertMapToGraphEdge converts a map to a graph edge (Pure Core)
func convertMapToGraphEdge(edgeMap map[string]interface{}) GraphEdge {
	return GraphEdge{
		ID:     getStringFromMap(edgeMap, "id"),
		Source: getStringFromMap(edgeMap, "source"),
		Target: getStringFromMap(edgeMap, "target"),
		Type:   getStringFromMap(edgeMap, "type"),
		Label:  getStringFromMap(edgeMap, "label"),
	}
}

// Helper functions (Pure Core)
func generateUserID(login string) int {
	// Simple hash-based ID generation for users
	// In a real system, this would be retrieved from GitHub API
	hash := 0
	for _, c := range login {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

func getMapFromMap(m map[string]interface{}, key string) map[string]interface{} {
	if value, exists := m[key]; exists {
		if subMap, ok := value.(map[string]interface{}); ok {
			return subMap
		}
	}
	return make(map[string]interface{})
}

func isTeamOwner(owner string) bool {
	return strings.Contains(owner, "/") && strings.HasPrefix(owner, "@")
}

func extractTeamSlug(teamOwner string) string {
	// Extract team slug from @org/team format
	cleaned := strings.TrimPrefix(teamOwner, "@")
	parts := strings.Split(cleaned, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return cleaned
}

// Validation helper functions (Pure Core)
func validateOrgNameNotEmpty(orgName string) {
	if orgName == "" {
		panic("Organization name cannot be empty")
	}
}

func validateRepoFullNameNotEmpty(repoFullName string) {
	if repoFullName == "" {
		panic("Repository full name cannot be empty")
	}
}

func validateTopicNameNotEmpty(topicName string) {
	if topicName == "" {
		panic("Topic name cannot be empty")
	}
}
