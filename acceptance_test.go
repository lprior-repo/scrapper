package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AcceptanceTestSuite struct {
	suite.Suite
	service GraphService
	config  *Config
	ctx     context.Context
	cancel  context.CancelFunc
}

func (suite *AcceptanceTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip("Skipping acceptance tests in short mode")
	}

	suite.ctx, suite.cancel = context.WithCancel(context.Background())

	config, err := LoadConfig()
	suite.Require().NoError(err)
	suite.config = config

	service, err := NewGraphService(config.GraphDB)
	suite.Require().NoError(err)
	suite.service = service

	err = suite.service.Connect(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *AcceptanceTestSuite) TearDownSuite() {
	if suite.service != nil {
		if err := suite.service.Close(suite.ctx); err != nil {
			suite.T().Logf("Failed to close service: %v", err)
		}
	}
	suite.cancel()
}

func (suite *AcceptanceTestSuite) SetupTest() {
	err := suite.service.ClearAll(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *AcceptanceTestSuite) TestUserCanManageGraphNodes() {
	t := suite.T()

	assert.NotNil(t, suite.service)
	err := suite.service.Health(suite.ctx)
	require.NoError(t, err, "Graph service should be healthy")

	createdNode, err := suite.service.CreateNode(suite.ctx, "Person", map[string]interface{}{
		"name":       "John Doe",
		"age":        30,
		"occupation": "Software Engineer",
		"created_at": time.Now().Unix(),
	})

	require.NoError(t, err, "Node creation should succeed")
	assert.NotEmpty(t, createdNode.ID, "Created node should have an ID")
	assert.Contains(t, createdNode.Labels, "Person", "Created node should have Person label")
	assert.Equal(t, "John Doe", createdNode.Properties["name"], "Node should have correct name")
	assert.Equal(t, int64(30), createdNode.Properties["age"], "Node should have correct age")

	// WHEN: User retrieves the node
	retrievedNode, err := suite.service.GetNode(suite.ctx, createdNode.ID)

	// THEN: The node should be retrieved with all properties intact
	require.NoError(t, err, "Node retrieval should succeed")
	assert.Equal(t, createdNode.ID, retrievedNode.ID, "Retrieved node should have same ID")
	assert.Equal(t, createdNode.Properties["name"], retrievedNode.Properties["name"], "Retrieved node should have same name")
	assert.Equal(t, createdNode.Properties["age"], retrievedNode.Properties["age"], "Retrieved node should have same age")

	// WHEN: User updates the node
	err = suite.service.UpdateNode(suite.ctx, createdNode.ID, map[string]interface{}{
		"age":        31,
		"updated_at": time.Now().Unix(),
	})

	// THEN: The node should be updated successfully
	require.NoError(t, err, "Node update should succeed")

	// AND: The updated values should be persisted
	updatedNode, err := suite.service.GetNode(suite.ctx, createdNode.ID)
	require.NoError(t, err, "Updated node retrieval should succeed")
	assert.Equal(t, int64(31), updatedNode.Properties["age"], "Node age should be updated")
	assert.NotNil(t, updatedNode.Properties["updated_at"], "Node should have updated_at timestamp")

	// WHEN: User deletes the node
	err = suite.service.DeleteNode(suite.ctx, createdNode.ID)

	// THEN: The node should be deleted successfully
	require.NoError(t, err, "Node deletion should succeed")

	// AND: The node should no longer exist
	_, err = suite.service.GetNode(suite.ctx, createdNode.ID)
	require.Error(t, err, "Deleted node should not be retrievable")
}

func (suite *AcceptanceTestSuite) TestUserCanManageGraphRelationships() {
	t := suite.T()

	// GIVEN: A user wants to manage relationships between nodes
	// AND: Two nodes exist in the graph

	// Create prerequisite nodes
	person1, err := suite.service.CreateNode(suite.ctx, "Person", map[string]interface{}{
		"name": "Alice",
		"age":  28,
	})
	require.NoError(t, err, "Person1 creation should succeed")

	person2, err := suite.service.CreateNode(suite.ctx, "Person", map[string]interface{}{
		"name": "Bob",
		"age":  32,
	})
	require.NoError(t, err, "Person2 creation should succeed")

	// WHEN: User creates a relationship between the nodes
	relationship, err := suite.service.CreateRelationship(suite.ctx, person1.ID, person2.ID, "KNOWS", map[string]interface{}{
		"since":    "2020-01-01",
		"strength": 0.8,
		"context":  "work",
	})

	// THEN: The relationship should be created successfully
	require.NoError(t, err, "Relationship creation should succeed")
	assert.NotEmpty(t, relationship.ID, "Relationship should have an ID")
	assert.Equal(t, "KNOWS", relationship.Type, "Relationship should have correct type")
	assert.Equal(t, person1.ID, relationship.FromID, "Relationship should have correct from ID")
	assert.Equal(t, person2.ID, relationship.ToID, "Relationship should have correct to ID")
	assert.Equal(t, "2020-01-01", relationship.Properties["since"], "Relationship should have correct since date")
	assert.InEpsilon(t, 0.8, relationship.Properties["strength"], 0.01, "Relationship should have correct strength")

	// WHEN: User retrieves the relationship
	retrievedRel, err := suite.service.GetRelationship(suite.ctx, relationship.ID)

	// THEN: The relationship should be retrieved with all properties intact
	require.NoError(t, err, "Relationship retrieval should succeed")
	assert.Equal(t, relationship.ID, retrievedRel.ID, "Retrieved relationship should have same ID")
	assert.Equal(t, relationship.Type, retrievedRel.Type, "Retrieved relationship should have same type")
	assert.Equal(t, relationship.Properties["since"], retrievedRel.Properties["since"], "Retrieved relationship should have same since date")

	// WHEN: User queries for connected nodes
	results, err := suite.service.ExecuteReadQuery(suite.ctx, `
		MATCH (p1:Person)-[r:KNOWS]->(p2:Person)
		WHERE p1.name = $name
		RETURN p1.name as from_name, p2.name as to_name, r.since as since, r.strength as strength
	`, map[string]interface{}{
		"name": "Alice",
	})

	// THEN: The query should return the connected nodes
	require.NoError(t, err, "Query should succeed")
	require.Len(t, results, 1, "Query should return one result")
	result := results[0]
	assert.Equal(t, "Alice", result["from_name"], "Query should return correct from name")
	assert.Equal(t, "Bob", result["to_name"], "Query should return correct to name")
	assert.Equal(t, "2020-01-01", result["since"], "Query should return correct since date")
	assert.InEpsilon(t, 0.8, result["strength"], 0.01, "Query should return correct strength")

	// WHEN: User deletes the relationship
	err = suite.service.DeleteRelationship(suite.ctx, relationship.ID)

	// THEN: The relationship should be deleted successfully
	require.NoError(t, err, "Relationship deletion should succeed")

	// AND: The relationship should no longer exist
	_, err = suite.service.GetRelationship(suite.ctx, relationship.ID)
	require.Error(t, err, "Deleted relationship should not be retrievable")

	// AND: The nodes should still exist
	_, err = suite.service.GetNode(suite.ctx, person1.ID)
	require.NoError(t, err, "Person1 should still exist after relationship deletion")
	_, err = suite.service.GetNode(suite.ctx, person2.ID)
	require.NoError(t, err, "Person2 should still exist after relationship deletion")
}

func (suite *AcceptanceTestSuite) TestUserCanExecuteComplexQueries() {
	t := suite.T()

	// GIVEN: A user wants to execute complex queries on graph data
	// AND: A complex graph structure exists

	// Create a small network of people and relationships
	alice, err := suite.service.CreateNode(suite.ctx, "Person", map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"city": "New York",
	})
	require.NoError(t, err)

	bob, err := suite.service.CreateNode(suite.ctx, "Person", map[string]interface{}{
		"name": "Bob",
		"age":  25,
		"city": "San Francisco",
	})
	require.NoError(t, err)

	charlie, err := suite.service.CreateNode(suite.ctx, "Person", map[string]interface{}{
		"name": "Charlie",
		"age":  35,
		"city": "New York",
	})
	require.NoError(t, err)

	// Create relationships
	_, err = suite.service.CreateRelationship(suite.ctx, alice.ID, bob.ID, "KNOWS", map[string]interface{}{
		"since": "2020-01-01",
	})
	require.NoError(t, err)

	_, err = suite.service.CreateRelationship(suite.ctx, alice.ID, charlie.ID, "WORKS_WITH", map[string]interface{}{
		"since": "2019-06-01",
	})
	require.NoError(t, err)

	_, err = suite.service.CreateRelationship(suite.ctx, bob.ID, charlie.ID, "KNOWS", map[string]interface{}{
		"since": "2021-03-15",
	})
	require.NoError(t, err)

	// WHEN: User executes a complex query to find mutual connections
	results, err := suite.service.ExecuteReadQuery(suite.ctx, `
		MATCH (p1:Person)-[r1]->(p2:Person)-[r2]->(p3:Person)
		WHERE p1.name = $name AND p1 <> p3
		RETURN p1.name as person1, p2.name as person2, p3.name as person3, 
			   type(r1) as relationship1, type(r2) as relationship2
		ORDER BY p2.name, p3.name
	`, map[string]interface{}{
		"name": "Alice",
	})

	// THEN: The query should return the expected paths
	require.NoError(t, err, "Complex query should succeed")
	require.GreaterOrEqual(t, len(results), 1, "Query should return at least one path")

	// Verify paths contain expected elements
	foundAliceToCharlie := false
	for _, result := range results {
		if result["person1"] == "Alice" && result["person3"] == "Charlie" {
			foundAliceToCharlie = true
			assert.Equal(t, "Alice", result["person1"], "Path should start with Alice")
			assert.Equal(t, "Charlie", result["person3"], "Path should end with Charlie")
			break
		}
	}
	assert.True(t, foundAliceToCharlie, "Should find a path from Alice to Charlie")

	// WHEN: User executes an aggregation query
	aggregationResults, err := suite.service.ExecuteReadQuery(suite.ctx, `
		MATCH (p:Person)
		RETURN p.city as city, count(p) as person_count, avg(p.age) as avg_age
		ORDER BY person_count DESC
	`, nil)

	// THEN: The aggregation should return correct results
	require.NoError(t, err, "Aggregation query should succeed")
	require.Len(t, aggregationResults, 2, "Should return results for both cities")

	// New York should have 2 people (Alice and Charlie)
	nyResult := aggregationResults[0]
	assert.Equal(t, "New York", nyResult["city"], "First result should be New York")
	assert.Equal(t, int64(2), nyResult["person_count"], "New York should have 2 people")
	assert.InEpsilon(t, float64(32.5), nyResult["avg_age"], 0.01, "New York average age should be 32.5")

	// San Francisco should have 1 person (Bob)
	sfResult := aggregationResults[1]
	assert.Equal(t, "San Francisco", sfResult["city"], "Second result should be San Francisco")
	assert.Equal(t, int64(1), sfResult["person_count"], "San Francisco should have 1 person")
	assert.InEpsilon(t, float64(25), sfResult["avg_age"], 0.01, "San Francisco average age should be 25")
}

func (suite *AcceptanceTestSuite) TestUserCanPerformBatchOperations() {
	t := suite.T()

	// GIVEN: A user wants to perform multiple operations atomically
	// AND: The system supports batch operations

	// Verify prerequisites
	assert.NotNil(t, suite.service)
	err := suite.service.Health(suite.ctx)
	require.NoError(t, err, "Graph service should be healthy")

	// WHEN: User executes a batch of operations
	operations := []BatchOperation{
		{
			Type:  "create_node",
			Query: "CREATE (p:Person {name: $name1, age: $age1, team: $team})",
			Parameters: map[string]interface{}{
				"name1": "Alice",
				"age1":  30,
				"team":  "Engineering",
			},
		},
		{
			Type:  "create_node",
			Query: "CREATE (p:Person {name: $name2, age: $age2, team: $team})",
			Parameters: map[string]interface{}{
				"name2": "Bob",
				"age2":  25,
				"team":  "Engineering",
			},
		},
		{
			Type:  "create_node",
			Query: "CREATE (p:Person {name: $name3, age: $age3, team: $team})",
			Parameters: map[string]interface{}{
				"name3": "Charlie",
				"age3":  35,
				"team":  "Design",
			},
		},
		{
			Type: "create_relationships",
			Query: `MATCH (a:Person {name: "Alice"}), (b:Person {name: "Bob"})
					CREATE (a)-[:WORKS_WITH {since: $since}]->(b)`,
			Parameters: map[string]interface{}{
				"since": "2020-01-01",
			},
		},
		{
			Type: "create_relationships",
			Query: `MATCH (a:Person {name: "Alice"}), (c:Person {name: "Charlie"})
					CREATE (a)-[:COLLABORATES_WITH {since: $since}]->(c)`,
			Parameters: map[string]interface{}{
				"since": "2019-06-01",
			},
		},
	}

	err = suite.service.ExecuteBatch(suite.ctx, operations)

	// THEN: All operations should succeed atomically
	require.NoError(t, err, "Batch operations should succeed")

	// AND: All nodes should be created
	nodeResults, err := suite.service.ExecuteReadQuery(suite.ctx, `
		MATCH (p:Person)
		RETURN p.name as name, p.age as age, p.team as team
		ORDER BY p.name
	`, nil)
	require.NoError(t, err, "Node query should succeed")
	require.Len(t, nodeResults, 3, "Should have 3 nodes")

	// Verify node data
	assert.Equal(t, "Alice", nodeResults[0]["name"])
	assert.Equal(t, int64(30), nodeResults[0]["age"])
	assert.Equal(t, "Engineering", nodeResults[0]["team"])

	assert.Equal(t, "Bob", nodeResults[1]["name"])
	assert.Equal(t, int64(25), nodeResults[1]["age"])
	assert.Equal(t, "Engineering", nodeResults[1]["team"])

	assert.Equal(t, "Charlie", nodeResults[2]["name"])
	assert.Equal(t, int64(35), nodeResults[2]["age"])
	assert.Equal(t, "Design", nodeResults[2]["team"])

	// AND: All relationships should be created
	relResults, err := suite.service.ExecuteReadQuery(suite.ctx, `
		MATCH (a:Person)-[r]->(b:Person)
		RETURN a.name as from_name, b.name as to_name, type(r) as rel_type, r.since as since
		ORDER BY a.name, b.name
	`, nil)
	require.NoError(t, err, "Relationship query should succeed")
	require.Len(t, relResults, 2, "Should have 2 relationships")

	// Verify relationship data
	assert.Equal(t, "Alice", relResults[0]["from_name"])
	assert.Equal(t, "Bob", relResults[0]["to_name"])
	assert.Equal(t, "WORKS_WITH", relResults[0]["rel_type"])
	assert.Equal(t, "2020-01-01", relResults[0]["since"])

	assert.Equal(t, "Alice", relResults[1]["from_name"])
	assert.Equal(t, "Charlie", relResults[1]["to_name"])
	assert.Equal(t, "COLLABORATES_WITH", relResults[1]["rel_type"])
	assert.Equal(t, "2019-06-01", relResults[1]["since"])
}

func (suite *AcceptanceTestSuite) TestSystemHandlesErrorsGracefully() {
	t := suite.T()

	// GIVEN: A user interacts with the system
	// AND: Various error conditions may occur

	// Verify prerequisites
	assert.NotNil(t, suite.service)
	err := suite.service.Health(suite.ctx)
	require.NoError(t, err, "Graph service should be healthy")

	// WHEN: User tries to retrieve a non-existent node
	_, err = suite.service.GetNode(suite.ctx, "999999")

	// THEN: The system should return a meaningful error
	require.Error(t, err, "Should return error for non-existent node")
	assert.Contains(t, err.Error(), "not found", "Error should indicate node not found")

	// WHEN: User tries to create a relationship with invalid nodes
	_, err = suite.service.CreateRelationship(suite.ctx, "999999", "888888", "INVALID", nil)

	// THEN: The system should return a meaningful error
	require.Error(t, err, "Should return error for invalid relationship")

	// WHEN: User executes an invalid query
	_, err = suite.service.ExecuteReadQuery(suite.ctx, "INVALID CYPHER SYNTAX", nil)

	// THEN: The system should return a meaningful error
	require.Error(t, err, "Should return error for invalid query")

	// WHEN: User tries to delete a non-existent node
	err = suite.service.DeleteNode(suite.ctx, "999999")

	// THEN: The system should handle it gracefully (no error for idempotent operations)
	require.NoError(t, err, "Deleting non-existent node should be idempotent")

	// AND: The system should remain healthy after errors
	err = suite.service.Health(suite.ctx)
	require.NoError(t, err, "System should remain healthy after error conditions")
}

func TestAcceptanceTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptanceTestSuite))
}
