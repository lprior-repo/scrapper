package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// RelationshipTestSuite tests the enhanced relationship functionality
type RelationshipTestSuite struct {
	suite.Suite
	conn   GraphConnection
	config *Config
	ctx    context.Context
	cancel context.CancelFunc
}

func (suite *RelationshipTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip("Skipping relationship tests in short mode")
	}

	suite.ctx, suite.cancel = context.WithCancel(context.Background())

	config, err := LoadConfig()
	suite.Require().NoError(err)
	suite.config = config

	conn, err := createConnection(suite.ctx, *config)
	suite.Require().NoError(err)
	suite.conn = conn

	// Clear any existing data
	err = clearAll(suite.ctx, suite.conn)
	suite.Require().NoError(err)
}

func (suite *RelationshipTestSuite) TearDownSuite() {
	if suite.conn.Driver != nil {
		if err := closeConnection(suite.ctx, suite.conn); err != nil {
			suite.T().Logf("Failed to close connection: %v", err)
		}
	}
	suite.cancel()
}

func (suite *RelationshipTestSuite) SetupTest() {
	// Clean database before each test
	err := clearAll(suite.ctx, suite.conn)
	suite.Require().NoError(err)
}

func (suite *RelationshipTestSuite) TestRelationshipUpdate() {
	t := suite.T()

	// Create two nodes
	node1, err := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label: "Person",
		Properties: map[string]interface{}{
			"name": "Alice",
		},
	})
	require.NoError(t, err)

	node2, err := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label: "Person",
		Properties: map[string]interface{}{
			"name": "Bob",
		},
	})
	require.NoError(t, err)

	// Create a relationship
	rel, err := createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  node1.ID,
		ToID:    node2.ID,
		RelType: "KNOWS",
		Properties: map[string]interface{}{
			"since":    "2020",
			"strength": 0.5,
		},
	})
	require.NoError(t, err)

	// Update the relationship
	updatedRel, err := updateRelationship(suite.ctx, suite.conn, rel.ID, map[string]interface{}{
		"since":    "2021",
		"strength": 0.8,
		"context":  "work",
	})
	require.NoError(t, err)

	// Verify the update
	assert.Equal(t, rel.ID, updatedRel.ID)
	assert.Equal(t, "2021", updatedRel.Properties["since"])
	assert.InEpsilon(t, 0.8, updatedRel.Properties["strength"], 0.01)
	assert.Equal(t, "work", updatedRel.Properties["context"])
}

func (suite *RelationshipTestSuite) TestFindRelationshipsByType() {
	t := suite.T()

	// Create nodes
	alice, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Alice"},
	})

	bob, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Bob"},
	})

	charlie, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Charlie"},
	})

	// Create different types of relationships
	_, err := createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:     alice.ID,
		ToID:       bob.ID,
		RelType:    "KNOWS",
		Properties: map[string]interface{}{"since": "2020"},
	})
	require.NoError(t, err)

	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:     alice.ID,
		ToID:       charlie.ID,
		RelType:    "WORKS_WITH",
		Properties: map[string]interface{}{"department": "engineering"},
	})
	require.NoError(t, err)

	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:     bob.ID,
		ToID:       charlie.ID,
		RelType:    "KNOWS",
		Properties: map[string]interface{}{"since": "2021"},
	})
	require.NoError(t, err)

	// Find relationships by type
	knowsRels, err := findRelationshipsByType(suite.ctx, suite.conn, "KNOWS")
	require.NoError(t, err)
	assert.Len(t, knowsRels, 2)
	for _, rel := range knowsRels {
		assert.Equal(t, "KNOWS", rel.Type)
	}

	worksWithRels, err := findRelationshipsByType(suite.ctx, suite.conn, "WORKS_WITH")
	require.NoError(t, err)
	assert.Len(t, worksWithRels, 1)
	assert.Equal(t, "WORKS_WITH", worksWithRels[0].Type)
}

func (suite *RelationshipTestSuite) TestFindRelationshipsByNode() {
	t := suite.T()

	// Create nodes
	alice, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Alice"},
	})

	bob, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Bob"},
	})

	charlie, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Charlie"},
	})

	// Create relationships
	_, err := createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    bob.ID,
		RelType: "KNOWS",
	})
	require.NoError(t, err)

	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  charlie.ID,
		ToID:    alice.ID,
		RelType: "WORKS_WITH",
	})
	require.NoError(t, err)

	// Find outgoing relationships for Alice
	outgoingRels, err := findRelationshipsByNode(suite.ctx, suite.conn, alice.ID, "outgoing")
	require.NoError(t, err)
	assert.Len(t, outgoingRels, 1)
	assert.Equal(t, alice.ID, outgoingRels[0].FromID)

	// Find incoming relationships for Alice
	incomingRels, err := findRelationshipsByNode(suite.ctx, suite.conn, alice.ID, "incoming")
	require.NoError(t, err)
	assert.Len(t, incomingRels, 1)
	assert.Equal(t, alice.ID, incomingRels[0].ToID)

	// Find all relationships for Alice
	allRels, err := findRelationshipsByNode(suite.ctx, suite.conn, alice.ID, "both")
	require.NoError(t, err)
	assert.Len(t, allRels, 2)
}

func (suite *RelationshipTestSuite) TestFindRelationshipsBetweenNodes() {
	t := suite.T()

	// Create nodes
	alice, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Alice"},
	})

	bob, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Bob"},
	})

	charlie, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Charlie"},
	})

	// Create multiple relationships between Alice and Bob
	_, err := createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:     alice.ID,
		ToID:       bob.ID,
		RelType:    "KNOWS",
		Properties: map[string]interface{}{"context": "personal"},
	})
	require.NoError(t, err)

	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:     alice.ID,
		ToID:       bob.ID,
		RelType:    "WORKS_WITH",
		Properties: map[string]interface{}{"context": "professional"},
	})
	require.NoError(t, err)

	// Create a relationship between Alice and Charlie
	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    charlie.ID,
		RelType: "KNOWS",
	})
	require.NoError(t, err)

	// Find relationships between Alice and Bob
	aliceBobRels, err := findRelationshipsBetweenNodes(suite.ctx, suite.conn, alice.ID, bob.ID)
	require.NoError(t, err)
	assert.Len(t, aliceBobRels, 2)

	// Find relationships between Alice and Charlie
	aliceCharlieRels, err := findRelationshipsBetweenNodes(suite.ctx, suite.conn, alice.ID, charlie.ID)
	require.NoError(t, err)
	assert.Len(t, aliceCharlieRels, 1)

	// Find relationships between Bob and Charlie (should be empty)
	bobCharlieRels, err := findRelationshipsBetweenNodes(suite.ctx, suite.conn, bob.ID, charlie.ID)
	require.NoError(t, err)
	assert.Len(t, bobCharlieRels, 0)
}

func (suite *RelationshipTestSuite) TestAdvancedValidation() {
	t := suite.T()

	// Create nodes
	alice, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Alice"},
	})

	bob, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Bob"},
	})

	// Test invalid relationship type
	_, err := createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    bob.ID,
		RelType: "INVALID_TYPE",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid relationship type")

	// Test strength validation
	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    bob.ID,
		RelType: "KNOWS",
		Properties: map[string]interface{}{
			"strength": 1.5, // Invalid: > 1.0
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "strength must be between 0.0 and 1.0")

	// Test self-relationship restriction
	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    alice.ID,
		RelType: "MANAGES", // Self-management not allowed
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "self-relationships not allowed")

	// Test valid relationship
	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    bob.ID,
		RelType: "KNOWS",
		Properties: map[string]interface{}{
			"since":    "2020",
			"strength": 0.8,
		},
	})
	require.NoError(t, err)
}

func (suite *RelationshipTestSuite) TestPathFinding() {
	t := suite.T()

	// Create a small network
	alice, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Alice"},
	})

	bob, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Bob"},
	})

	charlie, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "Charlie"},
	})

	david, _ := createNode(suite.ctx, suite.conn, GraphOperationRequest{
		Label:      "Person",
		Properties: map[string]interface{}{"name": "David"},
	})

	// Create relationships to form a path: Alice -> Bob -> Charlie -> David
	_, err := createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    bob.ID,
		RelType: "KNOWS",
	})
	require.NoError(t, err)

	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  bob.ID,
		ToID:    charlie.ID,
		RelType: "WORKS_WITH",
	})
	require.NoError(t, err)

	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  charlie.ID,
		ToID:    david.ID,
		RelType: "KNOWS",
	})
	require.NoError(t, err)

	// Also create a direct relationship: Alice -> David
	_, err = createRelationship(suite.ctx, suite.conn, GraphOperationRequest{
		FromID:  alice.ID,
		ToID:    david.ID,
		RelType: "KNOWS",
	})
	require.NoError(t, err)

	// Find shortest path from Alice to David
	shortestPath, err := findShortestPath(suite.ctx, suite.conn, alice.ID, david.ID)
	require.NoError(t, err)
	assert.NotNil(t, shortestPath)
	assert.Equal(t, 1, shortestPath.Length) // Direct path should be shortest
	assert.Len(t, shortestPath.Nodes, 2)
	assert.Len(t, shortestPath.Relationships, 1)

	// Find all paths from Alice to David
	allPaths, err := findAllPaths(suite.ctx, suite.conn, alice.ID, david.ID, 5)
	require.NoError(t, err)
	assert.Len(t, allPaths, 2) // Direct path and indirect path through Bob and Charlie

	// Verify paths are sorted by length
	assert.LessOrEqual(t, allPaths[0].Length, allPaths[1].Length)
}

func TestRelationshipTestSuite(t *testing.T) {
	suite.Run(t, new(RelationshipTestSuite))
}
