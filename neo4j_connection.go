package main

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jConfig holds configuration for Neo4j connection
type Neo4jConfig struct {
	URI      string
	Username string
	Password string
}

func createNeo4jDriver(config Neo4jConfig) (neo4j.DriverWithContext, error) {
	auth := neo4j.BasicAuth(config.Username, config.Password, "")

	driver, err := neo4j.NewDriverWithContext(config.URI, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	return driver, nil
}

func verifyNeo4jConnection(ctx context.Context, driver neo4j.DriverWithContext) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if err := session.Close(ctx); err != nil {
			_ = err
		}
	}()

	_, err := session.Run(ctx, "RETURN 1", nil)
	if err != nil {
		return fmt.Errorf("failed to verify Neo4j connection: %w", err)
	}

	return nil
}

func waitForNeo4jReady(ctx context.Context, config Neo4jConfig, timeout time.Duration) error {
	driver, err := createNeo4jDriver(config)
	if err != nil {
		return err
	}
	defer func() {
		if err := driver.Close(ctx); err != nil {
			_ = err
		}
	}()

	return waitForConnectionReady(ctx, driver, timeout)
}

// waitForConnectionReady waits for Neo4j connection to be ready
func waitForConnectionReady(ctx context.Context, driver neo4j.DriverWithContext, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := verifyNeo4jConnection(ctx, driver); err == nil {
			return nil
		}

		if err := waitWithContext(ctx, time.Second); err != nil {
			return err
		}
	}

	return fmt.Errorf("Neo4j did not become ready within %v", timeout)
}

// waitWithContext waits for duration or context cancellation
func waitWithContext(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}
