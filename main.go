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

	graphService, err := setupGraphService(ctx, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := graphService.Close(ctx); err != nil {
			log.Printf("Failed to close graph service: %v", err)
		}
	}()

	if config.IsDevelopment() {
		runDemoOperations(ctx, graphService)
	}

	return waitForShutdown(ctx)
}

// initializeConfiguration loads and validates configuration
func initializeConfiguration() (*Config, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Using %s graph database", config.GraphDB.Provider)
	return config, nil
}

// setupGraphService initializes and connects to graph service
func setupGraphService(ctx context.Context, config *Config) (GraphService, error) {
	graphService, err := NewGraphService(config.GraphDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph service: %w", err)
	}

	log.Println("Connecting to graph database...")
	if err := graphService.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to graph database: %w", err)
	}

	if err := graphService.Health(ctx); err != nil {
		return nil, fmt.Errorf("graph database health check failed: %w", err)
	}

	log.Println("Graph database connection established successfully")
	return graphService, nil
}

// runDemoOperations runs demonstration operations in development mode
func runDemoOperations(ctx context.Context, service GraphService) {
	err := demonstrateGraphOperations(ctx, service)
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

func demonstrateGraphOperations(ctx context.Context, service GraphService) error {
	node, err := service.CreateNode(ctx, "TestNode", map[string]interface{}{
		"name":        "demo-node",
		"created_at":  "2024-01-01T00:00:00Z",
		"description": "A demonstration node created by overseer",
	})
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	log.Printf("Created test node: %+v", node)

	node2, err := service.CreateNode(ctx, "TestNode", map[string]interface{}{
		"name":        "demo-node-2",
		"created_at":  "2024-01-01T00:00:00Z",
		"description": "Another demonstration node",
	})
	if err != nil {
		return fmt.Errorf("failed to create second test node: %w", err)
	}

	log.Printf("Created second test node: %+v", node2)

	rel, err := service.CreateRelationship(ctx, node.ID, node2.ID, "CONNECTED_TO", map[string]interface{}{
		"relationship_type": "demo",
		"created_at":        "2024-01-01T00:00:00Z",
	})
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	log.Printf("Created relationship: %+v", rel)

	results, err := service.ExecuteReadQuery(ctx, "MATCH (n:TestNode) RETURN n.name as name, n.description as description", nil)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	log.Printf("Query results: %+v", results)

	return nil
}
