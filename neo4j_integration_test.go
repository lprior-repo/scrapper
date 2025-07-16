package main

import (
	"context"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeo4jIntegration(t *testing.T) {
	t.Parallel()
	driver, cleanup := setupNeo4jForTests(t)
	defer cleanup()

	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			// Log error but continue
			_ = err
		}
	}()

	// Test basic connectivity
	result, err := session.Run(ctx, "RETURN 'Hello Neo4j' AS message", nil)
	require.NoError(t, err)

	record, err := result.Single(ctx)
	require.NoError(t, err)

	message, found := record.Get("message")
	require.True(t, found)
	assert.Equal(t, "Hello Neo4j", message)
}

func TestNeo4jConnectionUtilities(t *testing.T) {
	t.Parallel()
	t.Run("createNeo4jDriver", func(t *testing.T) {
		t.Parallel()
		config := getTestNeo4jConfig()
		driver, err := createNeo4jDriver(config)
		require.NoError(t, err)
		require.NotNil(t, driver)

		ctx := context.Background()
		defer func() {
			if err := driver.Close(ctx); err != nil {
				// Log error but continue
				_ = err
			}
		}()

		err = verifyNeo4jConnection(ctx, driver)
		assert.NoError(t, err)
	})

	t.Run("waitForNeo4jReady", func(t *testing.T) {
		t.Parallel()
		if testing.Short() {
			t.Skip("Skipping Neo4j integration test in short mode")
		}

		ctx := context.Background()
		config := getTestNeo4jConfig()

		err := waitForNeo4jReady(ctx, config, 10*time.Second)
		assert.NoError(t, err)
	})
}

func TestNeo4jTestSetup(t *testing.T) {
	t.Parallel()
	t.Run("setupNeo4jForTests", func(t *testing.T) {
		t.Parallel()
		driver, cleanup := setupNeo4jForTests(t)
		defer cleanup()

		assert.NotNil(t, driver)

		// Verify we can use the driver
		ctx := context.Background()
		err := verifyNeo4jConnection(ctx, driver)
		assert.NoError(t, err)
	})

	t.Run("cleanupNeo4jTestData", func(t *testing.T) {
		t.Parallel()
		driver, cleanup := setupNeo4jForTests(t)
		defer cleanup()

		ctx := context.Background()
		session := driver.NewSession(ctx, neo4j.SessionConfig{})
		defer func() {
			if err := session.Close(ctx); err != nil {
				// Log error but continue
				_ = err
			}
		}()

		// Create some test data with unique identifier
		uniqueID := "test-cleanup-" + t.Name()
		_, err := session.Run(ctx, "CREATE (n:TestNode {name: $name})", map[string]interface{}{"name": uniqueID})
		require.NoError(t, err)

		// Verify data exists
		result, err := session.Run(ctx, "MATCH (n:TestNode {name: $name}) RETURN count(n) AS count", map[string]interface{}{"name": uniqueID})
		require.NoError(t, err)

		record, err := result.Single(ctx)
		require.NoError(t, err)

		count, found := record.Get("count")
		require.True(t, found)
		assert.Equal(t, int64(1), count)

		// Clean up
		cleanupNeo4jTestData(ctx, t, driver)

		// Verify our specific test data is gone
		result, err = session.Run(ctx, "MATCH (n:TestNode {name: $name}) RETURN count(n) AS count", map[string]interface{}{"name": uniqueID})
		require.NoError(t, err)

		record, err = result.Single(ctx)
		require.NoError(t, err)

		count, found = record.Get("count")
		require.True(t, found)
		assert.Equal(t, int64(0), count)
	})
}
