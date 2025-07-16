package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupNeo4jServiceForTest creates and connects to a test Neo4j service
func setupNeo4jServiceForTest(t *testing.T) (GraphService, context.Context) {
	t.Helper()
	config := GraphServiceConfig{
		Provider: "neo4j",
		Neo4j: struct {
			URI      string `json:"uri"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			URI:      "bolt://localhost:7687",
			Username: "neo4j",
			Password: "password",
		},
	}

	service, err := NewGraphService(config)
	require.NoError(t, err)
	require.NotNil(t, service)

	ctx := context.Background()
	err = service.Connect(ctx)
	require.NoError(t, err)

	return service, ctx
}

func TestGraphServiceFactory(t *testing.T) {
	t.Parallel()
	t.Run("CreateNeo4jService", func(t *testing.T) {
		t.Parallel()
		config := GraphServiceConfig{
			Provider: "neo4j",
			Neo4j: struct {
				URI      string `json:"uri"`
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				URI:      "bolt://localhost:7687",
				Username: "neo4j",
				Password: "password",
			},
		}

		service, err := NewGraphService(config)
		require.NoError(t, err)
		require.NotNil(t, service)

		// Verify it's a Neo4j service
		_, ok := service.(*Neo4jService)
		assert.True(t, ok, "Should create Neo4j service")
	})

	t.Run("CreateNeptuneService", func(t *testing.T) {
		t.Parallel()
		config := GraphServiceConfig{
			Provider: "neptune",
			Neptune: struct {
				Endpoint string `json:"endpoint"`
				Region   string `json:"region"`
			}{
				Endpoint: "wss://test.neptune.amazonaws.com:8182/gremlin",
				Region:   "us-east-1",
			},
		}

		service, err := NewGraphService(config)
		require.NoError(t, err)
		require.NotNil(t, service)

		// Verify it's a Neptune service
		_, ok := service.(*NeptuneService)
		assert.True(t, ok, "Should create Neptune service")
	})

	t.Run("UnsupportedProvider", func(t *testing.T) {
		t.Parallel()
		config := GraphServiceConfig{
			Provider: "unsupported",
		}

		service, err := NewGraphService(config)
		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "unsupported graph service provider")
	})
}

func TestNeo4jServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Neo4j integration test in short mode")
	}

	service, ctx := setupNeo4jServiceForTest(t)
	defer func() {
		if err := service.Close(ctx); err != nil {
			t.Logf("Failed to close service: %v", err)
		}
	}()

	runIntegrationTests(ctx, t, service)
}

// runIntegrationTests runs all integration test suites
func runIntegrationTests(ctx context.Context, t *testing.T, service GraphService) {
	t.Helper()

	// Clean up any existing data
	err := service.ClearAll(ctx)
	require.NoError(t, err)

	t.Run("NodeOperations", func(t *testing.T) {
		// Clean up before test
		err = service.ClearAll(ctx)
		require.NoError(t, err)

		// Create node
		node, err := service.CreateNode(ctx, "TestLabel", map[string]interface{}{
			"name":  "test-node",
			"value": 42,
		})
		require.NoError(t, err)
		require.NotNil(t, node)
		assert.NotEmpty(t, node.ID)
		assert.Contains(t, node.Labels, "TestLabel")
		assert.Equal(t, "test-node", node.Properties["name"])
		assert.Equal(t, int64(42), node.Properties["value"])

		// Get node
		retrievedNode, err := service.GetNode(ctx, node.ID)
		require.NoError(t, err)
		assert.Equal(t, node.ID, retrievedNode.ID)
		assert.Equal(t, node.Properties["name"], retrievedNode.Properties["name"])

		// Update node
		err = service.UpdateNode(ctx, node.ID, map[string]interface{}{
			"updated": true,
			"value":   100,
		})
		require.NoError(t, err)

		// Verify update
		updatedNode, err := service.GetNode(ctx, node.ID)
		require.NoError(t, err)
		assert.Equal(t, true, updatedNode.Properties["updated"])
		assert.Equal(t, int64(100), updatedNode.Properties["value"])

		// Delete node
		err = service.DeleteNode(ctx, node.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetNode(ctx, node.ID)
		require.Error(t, err)
	})

	t.Run("RelationshipOperations", func(t *testing.T) {
		// Clean up before test
		err = service.ClearAll(ctx)
		require.NoError(t, err)

		// Create two nodes
		node1, err := service.CreateNode(ctx, "Person", map[string]interface{}{
			"name": "Alice",
		})
		require.NoError(t, err)

		node2, err := service.CreateNode(ctx, "Person", map[string]interface{}{
			"name": "Bob",
		})
		require.NoError(t, err)

		// Create relationship
		rel, err := service.CreateRelationship(ctx, node1.ID, node2.ID, "KNOWS", map[string]interface{}{
			"since": "2020-01-01",
		})
		require.NoError(t, err)
		require.NotNil(t, rel)
		assert.NotEmpty(t, rel.ID)
		assert.Equal(t, "KNOWS", rel.Type)
		assert.Equal(t, node1.ID, rel.FromID)
		assert.Equal(t, node2.ID, rel.ToID)
		assert.Equal(t, "2020-01-01", rel.Properties["since"])

		// Get relationship
		retrievedRel, err := service.GetRelationship(ctx, rel.ID)
		require.NoError(t, err)
		assert.Equal(t, rel.ID, retrievedRel.ID)
		assert.Equal(t, rel.Type, retrievedRel.Type)

		// Delete relationship
		err = service.DeleteRelationship(ctx, rel.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetRelationship(ctx, rel.ID)
		require.Error(t, err)

		// Clean up nodes
		if err := service.DeleteNode(ctx, node1.ID); err != nil {
			t.Logf("Failed to delete node1: %v", err)
		}
		if err := service.DeleteNode(ctx, node2.ID); err != nil {
			t.Logf("Failed to delete node2: %v", err)
		}
	})

	t.Run("QueryOperations", func(t *testing.T) {
		// Clean up before test
		err = service.ClearAll(ctx)
		require.NoError(t, err)

		// Create test data
		_, err := service.CreateNode(ctx, "TestQuery", map[string]interface{}{
			"name":  "query-test-1",
			"value": 1,
		})
		require.NoError(t, err)

		_, err = service.CreateNode(ctx, "TestQuery", map[string]interface{}{
			"name":  "query-test-2",
			"value": 2,
		})
		require.NoError(t, err)

		// Execute read query
		results, err := service.ExecuteReadQuery(ctx,
			"MATCH (n:TestQuery) RETURN n.name as name, n.value as value ORDER BY n.value",
			nil)
		require.NoError(t, err)
		require.Len(t, results, 2)

		assert.Equal(t, "query-test-1", results[0]["name"])
		assert.Equal(t, int64(1), results[0]["value"])
		assert.Equal(t, "query-test-2", results[1]["name"])
		assert.Equal(t, int64(2), results[1]["value"])

		// Execute write query
		writeResults, err := service.ExecuteWriteQuery(ctx,
			"CREATE (n:TestQuery {name: $name, value: $value}) RETURN n.name as name",
			map[string]interface{}{
				"name":  "query-test-3",
				"value": 3,
			})
		require.NoError(t, err)
		require.Len(t, writeResults, 1)
		assert.Equal(t, "query-test-3", writeResults[0]["name"])

		// Clean up
		if _, err := service.ExecuteWriteQuery(ctx, "MATCH (n:TestQuery) DELETE n", nil); err != nil {
			t.Logf("Failed to clean up TestQuery nodes: %v", err)
		}
	})

	t.Run("BatchOperations", func(t *testing.T) {
		// Clean up before test
		err = service.ClearAll(ctx)
		require.NoError(t, err)

		operations := []BatchOperation{
			{
				Type:  "create_node",
				Query: "CREATE (n:BatchTest {name: $name1, value: $value1})",
				Parameters: map[string]interface{}{
					"name1":  "batch-1",
					"value1": 1,
				},
			},
			{
				Type:  "create_node",
				Query: "CREATE (n:BatchTest {name: $name2, value: $value2})",
				Parameters: map[string]interface{}{
					"name2":  "batch-2",
					"value2": 2,
				},
			},
		}

		err = service.ExecuteBatch(ctx, operations)
		require.NoError(t, err)

		// Verify batch operations
		results, err := service.ExecuteReadQuery(ctx,
			"MATCH (n:BatchTest) RETURN count(n) as count",
			nil)
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, int64(2), results[0]["count"])

		// Clean up
		if _, err := service.ExecuteWriteQuery(ctx, "MATCH (n:BatchTest) DELETE n", nil); err != nil {
			t.Logf("Failed to clean up BatchTest nodes: %v", err)
		}
	})
}

func TestNeptuneServicePlaceholder(t *testing.T) {
	t.Parallel()
	service, err := NewNeptuneService(NeptuneConfig{
		Endpoint: "wss://test.neptune.amazonaws.com:8182/gremlin",
		Region:   "us-east-1",
	})
	require.NoError(t, err)
	require.NotNil(t, service)

	ctx := context.Background()

	// All Neptune operations should return not implemented errors
	err = service.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")

	err = service.Health(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")

	_, err = service.CreateNode(ctx, "Test", map[string]interface{}{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}
