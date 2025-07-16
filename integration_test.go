package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Integration tests for the interaction between orchestrator, pure core, and impure shell

type IntegrationTestSuite struct {
	suite.Suite
	service GraphService
	config  *Config
	ctx     context.Context
	cancel  context.CancelFunc
}

func (suite *IntegrationTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip("Skipping integration tests in short mode")
	}

	suite.ctx, suite.cancel = context.WithCancel(context.Background())

	// Load test configuration
	config, err := LoadConfig()
	suite.Require().NoError(err)
	suite.config = config

	// Create and connect to graph service
	service, err := NewGraphService(config.GraphDB)
	suite.Require().NoError(err)
	suite.service = service

	err = suite.service.Connect(suite.ctx)
	suite.Require().NoError(err)

	// Clear any existing data
	err = suite.service.ClearAll(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.service != nil {
		if err := suite.service.Close(suite.ctx); err != nil {
			suite.T().Logf("Failed to close service: %v", err)
		}
	}
	suite.cancel()
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.service.ClearAll(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TestServiceLayerIntegration() {
	t := suite.T()

	// Test the full flow from configuration to service operations
	assert.NotNil(t, suite.config)
	assert.NotNil(t, suite.service)

	// Health check
	err := suite.service.Health(suite.ctx)
	require.NoError(t, err)

	// Create a node through the service layer
	node, err := suite.service.CreateNode(suite.ctx, "IntegrationTest", map[string]interface{}{
		"name":      "test-node",
		"timestamp": time.Now().Unix(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, node.ID)
	assert.Contains(t, node.Labels, "IntegrationTest")

	// Verify the node was created
	retrievedNode, err := suite.service.GetNode(suite.ctx, node.ID)
	require.NoError(t, err)
	assert.Equal(t, node.ID, retrievedNode.ID)
	assert.Equal(t, node.Properties["name"], retrievedNode.Properties["name"])
}

func (suite *IntegrationTestSuite) TestConfigurationToServiceFlow() {
	t := suite.T()

	// Test that configuration properly flows through to service creation

	// Verify configuration is loaded correctly
	assert.NotEmpty(t, suite.config.GraphDB.Provider)

	// Verify service responds to configuration
	switch suite.config.GraphDB.Provider {
	case "neo4j":
		_, ok := suite.service.(*Neo4jService)
		assert.True(t, ok, "Should create Neo4j service when configured")
	case "neptune":
		_, ok := suite.service.(*NeptuneService)
		assert.True(t, ok, "Should create Neptune service when configured")
	}
}

func (suite *IntegrationTestSuite) TestConcurrentOperations() {
	t := suite.T()

	// Test concurrent operations to ensure thread safety
	const numGoroutines = 10
	resultChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			node, err := suite.service.CreateNode(suite.ctx, "ConcurrentTest", map[string]interface{}{
				"id":   id,
				"name": "concurrent-node",
			})
			if err != nil {
				resultChan <- err
				return
			}

			// Try to read the node back
			_, err = suite.service.GetNode(suite.ctx, node.ID)
			resultChan <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-resultChan
		require.NoError(t, err, "Concurrent operation %d should succeed", i)
	}

	// Verify all nodes were created
	results, err := suite.service.ExecuteReadQuery(suite.ctx,
		"MATCH (n:ConcurrentTest) RETURN count(n) as count", nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, int64(numGoroutines), results[0]["count"])
}

func (suite *IntegrationTestSuite) TestTransactionIntegrity() {
	t := suite.T()

	// Test batch operations to ensure transaction integrity
	operations := []BatchOperation{
		{
			Type:  "create_node",
			Query: "CREATE (n:TransactionTest {name: $name1, order: $order1})",
			Parameters: map[string]interface{}{
				"name1":  "node-1",
				"order1": 1,
			},
		},
		{
			Type:  "create_node",
			Query: "CREATE (n:TransactionTest {name: $name2, order: $order2})",
			Parameters: map[string]interface{}{
				"name2":  "node-2",
				"order2": 2,
			},
		},
		{
			Type: "create_relationship",
			Query: `MATCH (a:TransactionTest {order: 1}), (b:TransactionTest {order: 2}) 
					CREATE (a)-[r:FOLLOWS]->(b)`,
			Parameters: map[string]interface{}{},
		},
	}

	err := suite.service.ExecuteBatch(suite.ctx, operations)
	require.NoError(t, err)

	// Verify all operations succeeded
	results, err := suite.service.ExecuteReadQuery(suite.ctx,
		"MATCH (n:TransactionTest) RETURN count(n) as node_count", nil)
	require.NoError(t, err)
	assert.Equal(t, int64(2), results[0]["node_count"])

	results, err = suite.service.ExecuteReadQuery(suite.ctx,
		"MATCH ()-[r:FOLLOWS]->() RETURN count(r) as rel_count", nil)
	require.NoError(t, err)
	assert.Equal(t, int64(1), results[0]["rel_count"])
}

func (suite *IntegrationTestSuite) TestErrorHandlingIntegration() {
	t := suite.T()

	// Test error handling across the service layer

	// Test getting non-existent node
	_, err := suite.service.GetNode(suite.ctx, "999999")
	require.Error(t, err, "Should return error for non-existent node")

	// Test invalid query
	_, err = suite.service.ExecuteReadQuery(suite.ctx, "INVALID CYPHER QUERY", nil)
	require.Error(t, err, "Should return error for invalid query")

	// Test creating relationship with non-existent nodes
	_, err = suite.service.CreateRelationship(suite.ctx, "999999", "888888", "INVALID", nil)
	require.Error(t, err, "Should return error for relationship between non-existent nodes")
}

func (suite *IntegrationTestSuite) TestDataConsistency() {
	t := suite.T()

	// Test data consistency across operations

	// Create initial data
	node1, err := suite.service.CreateNode(suite.ctx, "ConsistencyTest", map[string]interface{}{
		"name":  "node1",
		"value": 100,
	})
	require.NoError(t, err)

	node2, err := suite.service.CreateNode(suite.ctx, "ConsistencyTest", map[string]interface{}{
		"name":  "node2",
		"value": 200,
	})
	require.NoError(t, err)

	// Create relationship
	rel, err := suite.service.CreateRelationship(suite.ctx, node1.ID, node2.ID, "RELATES_TO", map[string]interface{}{
		"strength": 0.8,
	})
	require.NoError(t, err)

	// Update node1
	err = suite.service.UpdateNode(suite.ctx, node1.ID, map[string]interface{}{
		"value":   150,
		"updated": true,
	})
	require.NoError(t, err)

	// Verify consistency with complex query
	results, err := suite.service.ExecuteReadQuery(suite.ctx, `
		MATCH (n1:ConsistencyTest)-[r:RELATES_TO]->(n2:ConsistencyTest)
		WHERE n1.updated = true
		RETURN n1.name as name1, n1.value as value1, n2.name as name2, n2.value as value2, r.strength as strength
	`, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)

	result := results[0]
	assert.Equal(t, "node1", result["name1"])
	assert.Equal(t, int64(150), result["value1"])
	assert.Equal(t, "node2", result["name2"])
	assert.Equal(t, int64(200), result["value2"])
	assert.InEpsilon(t, 0.8, result["strength"], 0.01)

	// Clean up and verify deletion
	err = suite.service.DeleteRelationship(suite.ctx, rel.ID)
	require.NoError(t, err)

	err = suite.service.DeleteNode(suite.ctx, node1.ID)
	require.NoError(t, err)

	err = suite.service.DeleteNode(suite.ctx, node2.ID)
	require.NoError(t, err)

	// Verify all data is gone
	results, err = suite.service.ExecuteReadQuery(suite.ctx,
		"MATCH (n:ConsistencyTest) RETURN count(n) as count", nil)
	require.NoError(t, err)
	assert.Equal(t, int64(0), results[0]["count"])
}

func (suite *IntegrationTestSuite) TestServiceRecovery() {
	t := suite.T()

	// Test service recovery after connection issues

	// Create some data
	node, err := suite.service.CreateNode(suite.ctx, "RecoveryTest", map[string]interface{}{
		"name": "recovery-node",
	})
	require.NoError(t, err)

	// Verify health before
	err = suite.service.Health(suite.ctx)
	require.NoError(t, err)

	// Simulate recovery by reconnecting
	err = suite.service.Close(suite.ctx)
	require.NoError(t, err)

	err = suite.service.Connect(suite.ctx)
	require.NoError(t, err)

	// Verify health after reconnection
	err = suite.service.Health(suite.ctx)
	require.NoError(t, err)

	// Verify data is still accessible
	retrievedNode, err := suite.service.GetNode(suite.ctx, node.ID)
	require.NoError(t, err)
	assert.Equal(t, node.ID, retrievedNode.ID)
	assert.Equal(t, node.Properties["name"], retrievedNode.Properties["name"])
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
