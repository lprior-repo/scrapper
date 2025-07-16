package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down overseer...")
		cancel()
	}()

	if err := runOverseer(ctx); err != nil {
		log.Printf("Overseer failed: %v", err)
		cancel()
		os.Exit(1)
	}

	cancel()
}

func runOverseer(ctx context.Context) error {
	log.Println("Starting overseer...")

	config, err := initializeConfiguration()
	if err != nil {
		return err
	}

	graphConn, err := setupGraphConnection(ctx, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := closeConnection(ctx, graphConn); err != nil {
			log.Printf("Failed to close graph connection: %v", err)
		}
	}()

	if checkIsDevelopment(config) {
		runDemoOperations(ctx, graphConn)
	}

	return waitForShutdown(ctx)
}

// initializeConfiguration loads and validates configuration
func initializeConfiguration() (Config, error) {
	config, err := LoadConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := validateConfig(*config); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Using %s graph database", config.GraphDB.Provider)
	return *config, nil
}

// setupGraphConnection initializes and connects to graph database
func setupGraphConnection(ctx context.Context, config Config) (GraphConnection, error) {
	log.Println("Connecting to graph database...")
	
	graphConn, err := createConnection(ctx, config)
	if err != nil {
		return GraphConnection{}, fmt.Errorf("failed to create graph connection: %w", err)
	}

	if err := healthCheck(ctx, graphConn); err != nil {
		_ = closeConnection(ctx, graphConn)
		return GraphConnection{}, fmt.Errorf("graph database health check failed: %w", err)
	}

	log.Println("Graph database connection established successfully")
	return graphConn, nil
}

// runDemoOperations runs demonstration operations in development mode
func runDemoOperations(ctx context.Context, conn GraphConnection) {
	err := demonstrateGraphOperations(ctx, conn)
	if err != nil {
		log.Printf("Demo operations failed: %v", err)
	}
}

// waitForShutdown waits for context cancellation
func waitForShutdown(ctx context.Context) error {
	log.Println("Overseer is running... Press Ctrl+C to stop")
	<-ctx.Done()
	log.Println("Overseer stopped")
	return nil
}

func demonstrateGraphOperations(ctx context.Context, conn GraphConnection) error {
	nodeRequest := GraphOperationRequest{
		Operation: "create_node",
		Label:     "TestNode",
		Properties: map[string]interface{}{
			"name":        "demo-node",
			"created_at":  "2024-01-01T00:00:00Z",
			"description": "A demonstration node created by overseer",
		},
	}

	node, err := createNode(ctx, conn, nodeRequest)
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	log.Printf("Created test node: %+v", node)

	nodeRequest2 := GraphOperationRequest{
		Operation: "create_node",
		Label:     "TestNode",
		Properties: map[string]interface{}{
			"name":        "demo-node-2",
			"created_at":  "2024-01-01T00:00:00Z",
			"description": "Another demonstration node",
		},
	}

	node2, err := createNode(ctx, conn, nodeRequest2)
	if err != nil {
		return fmt.Errorf("failed to create second test node: %w", err)
	}

	log.Printf("Created second test node: %+v", node2)

	relRequest := GraphOperationRequest{
		Operation: "create_relationship",
		FromID:    node.ID,
		ToID:      node2.ID,
		RelType:   "CONNECTED_TO",
		Properties: map[string]interface{}{
			"relationship_type": "demo",
			"created_at":        "2024-01-01T00:00:00Z",
		},
	}

	rel, err := createRelationship(ctx, conn, relRequest)
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	log.Printf("Created relationship: %+v", rel)

	results, err := executeCustomReadQuery(ctx, conn, "MATCH (n:TestNode) RETURN n.name as name, n.description as description", nil)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	log.Printf("Query results: %+v", results)

	return nil
}
