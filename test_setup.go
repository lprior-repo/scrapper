package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
)

func setupNeo4jForTests(t *testing.T) (neo4j.DriverWithContext, func()) {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping Neo4j integration test in short mode")
	}

	ctx := context.Background()
	config := getTestNeo4jConfig()

	err := waitForNeo4jReady(ctx, config, 30*time.Second)
	require.NoError(t, err, "Neo4j should be ready for testing")

	driver, err := createNeo4jDriver(config)
	require.NoError(t, err, "Should create Neo4j driver")

	err = verifyNeo4jConnection(ctx, driver)
	require.NoError(t, err, "Should connect to Neo4j")

	cleanup := func() {
		cleanupNeo4jTestData(ctx, t, driver)
		if err := driver.Close(ctx); err != nil {
			t.Logf("Failed to close Neo4j driver: %v", err)
		}
	}

	return driver, cleanup
}

func cleanupNeo4jTestData(ctx context.Context, t *testing.T, driver neo4j.DriverWithContext) {
	t.Helper()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			t.Logf("Failed to close session: %v", err)
		}
	}()

	_, err := session.Run(ctx, "MATCH (n) DETACH DELETE n", nil)
	require.NoError(t, err, "Should clean up test data")
}

func getTestNeo4jConfig() Neo4jConfig {
	uri := getEnvOrDefault("NEO4J_TEST_URI", "bolt://localhost:7687")
	username := getEnvOrDefault("NEO4J_TEST_USERNAME", "neo4j")
	password := getEnvOrDefault("NEO4J_TEST_PASSWORD", "password")

	return Neo4jConfig{
		URI:      uri,
		Username: username,
		Password: password,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
