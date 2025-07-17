package main

import "fmt"

// buildCreateNodeQuery builds a Cypher query for node creation (Pure Core)
func buildCreateNodeQuery(request GraphOperationRequest) string {
	if request.Label == "" {
		panic("Label cannot be empty")
	}
	return fmt.Sprintf("CREATE (n:%s) SET n = $properties RETURN id(n) as id, labels(n) as labels, properties(n) as properties", request.Label)
}

// buildGetNodeQuery builds a Cypher query for node retrieval (Pure Core)
func buildGetNodeQuery() string {
	return "MATCH (n) WHERE id(n) = $id RETURN id(n) as id, labels(n) as labels, properties(n) as properties"
}

// buildUpdateNodeQuery builds a Cypher query for node update (Pure Core)
func buildUpdateNodeQuery() string {
	return "MATCH (n) WHERE id(n) = $id SET n += $properties"
}

// buildDeleteNodeQuery builds a Cypher query for node deletion (Pure Core)
func buildDeleteNodeQuery() string {
	return "MATCH (n) WHERE id(n) = $id DETACH DELETE n"
}

// buildCreateRelationshipQuery builds a Cypher query for relationship creation (Pure Core)
func buildCreateRelationshipQuery(request GraphOperationRequest) string {
	if request.RelType == "" {
		panic("Relationship type cannot be empty")
	}
	return fmt.Sprintf(`
		MATCH (from) WHERE id(from) = $fromID
		MATCH (to) WHERE id(to) = $toID
		CREATE (from)-[r:%s]->(to)
		SET r = $properties
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`, request.RelType)
}

// buildGetRelationshipQuery builds a Cypher query for relationship retrieval (Pure Core)
func buildGetRelationshipQuery() string {
	return `
		MATCH (from)-[r]->(to) WHERE id(r) = $id
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`
}

// buildUpdateRelationshipQuery builds a Cypher query for relationship update (Pure Core)
func buildUpdateRelationshipQuery() string {
	return `
		MATCH (from)-[r]->(to) WHERE id(r) = $id
		SET r += $properties
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
	`
}

// buildDeleteRelationshipQuery builds a Cypher query for relationship deletion (Pure Core)
func buildDeleteRelationshipQuery() string {
	return "MATCH ()-[r]-() WHERE id(r) = $id DELETE r"
}

// buildFindRelationshipsByTypeQuery builds a Cypher query to find relationships by type (Pure Core)
func buildFindRelationshipsByTypeQuery(relType string) string {
	if relType == "" {
		panic("Relationship type cannot be empty")
	}
	return fmt.Sprintf(`
		MATCH (from)-[r:%s]->(to)
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
		ORDER BY id(r)
	`, relType)
}

// buildFindRelationshipsByNodeQuery builds a Cypher query to find relationships for a specific node (Pure Core)
func buildFindRelationshipsByNodeQuery(nodeID string, direction string) string {
	if nodeID == "" {
		panic("Node ID cannot be empty")
	}

	switch direction {
	case "outgoing":
		return `
			MATCH (n)-[r]->(to) WHERE id(n) = $nodeID
			RETURN id(r) as id, type(r) as type, id(n) as from_id, id(to) as to_id, properties(r) as properties
			ORDER BY id(r)
		`
	case "incoming":
		return `
			MATCH (from)-[r]->(n) WHERE id(n) = $nodeID
			RETURN id(r) as id, type(r) as type, id(from) as from_id, id(n) as to_id, properties(r) as properties
			ORDER BY id(r)
		`
	default: // both directions
		return `
			MATCH (n)-[r]-(other) WHERE id(n) = $nodeID
			RETURN id(r) as id, type(r) as type, 
				   CASE WHEN id(startNode(r)) = $nodeID 
					   THEN id(startNode(r)) 
					   ELSE id(endNode(r)) 
				   END as from_id,
				   CASE WHEN id(startNode(r)) = $nodeID 
					   THEN id(endNode(r)) 
					   ELSE id(startNode(r)) 
				   END as to_id,
				   properties(r) as properties
			ORDER BY id(r)
		`
	}
}

// buildFindRelationshipsBetweenNodesQuery builds a Cypher query to find relationships between two nodes (Pure Core)
func buildFindRelationshipsBetweenNodesQuery() string {
	return `
		MATCH (from)-[r]->(to) 
		WHERE id(from) = $fromID AND id(to) = $toID
		RETURN id(r) as id, type(r) as type, id(from) as from_id, id(to) as to_id, properties(r) as properties
		ORDER BY id(r)
	`
}

// buildFindShortestPathQuery builds a Cypher query to find the shortest path between two nodes (Pure Core)
func buildFindShortestPathQuery() string {
	return `
		MATCH (from), (to) 
		WHERE id(from) = $fromID AND id(to) = $toID
		MATCH path = shortestPath((from)-[*]-(to))
		RETURN 
			[node in nodes(path) | {id: id(node), labels: labels(node), properties: properties(node)}] as nodes,
			[rel in relationships(path) | {
				id: id(rel), 
				type: type(rel), 
				from_id: id(startNode(rel)), 
				to_id: id(endNode(rel)), 
				properties: properties(rel)
			}] as relationships,
			length(path) as length
	`
}

// buildFindAllPathsQuery builds a Cypher query to find all paths between two nodes up to a maximum depth (Pure Core)
func buildFindAllPathsQuery(maxDepth int) string {
	if maxDepth <= 0 {
		maxDepth = 5 // Default maximum depth
	}
	return fmt.Sprintf(`
		MATCH (from), (to) 
		WHERE id(from) = $fromID AND id(to) = $toID
		MATCH path = (from)-[*1..%d]-(to)
		RETURN 
			[node in nodes(path) | {id: id(node), labels: labels(node), properties: properties(node)}] as nodes,
			[rel in relationships(path) | {
				id: id(rel), 
				type: type(rel), 
				from_id: id(startNode(rel)), 
				to_id: id(endNode(rel)), 
				properties: properties(rel)
			}] as relationships,
			length(path) as length
		ORDER BY length(path)
		LIMIT 100
	`, maxDepth)
}