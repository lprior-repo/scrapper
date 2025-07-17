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
