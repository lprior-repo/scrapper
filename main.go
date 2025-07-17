package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// getDefaultConfig returns default configuration for the application
func getDefaultConfig() Config {
	config := Config{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}
	
	// Configure GraphDB based on environment
	if config.Environment == "production" {
		config.GraphDB = getDefaultGraphServiceConfig()
	} else {
		config.GraphDB = GraphServiceConfig{
			Provider: "neo4j",
			Neo4j: struct {
				URI      string `json:"uri"`
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				URI:      getEnvOrDefault("NEO4J_URI", "bolt://localhost:7687"),
				Username: getEnvOrDefault("NEO4J_USERNAME", "neo4j"),
				Password: getEnvOrDefault("NEO4J_PASSWORD", "password"),
			},
		}
	}
	
	return config
}

func main() {
	// Check command line arguments
	args := os.Args[1:]
	
	if len(args) == 0 {
		fmt.Println("GitHub Codeowners Visualization Tool")
		fmt.Println("Usage:")
		fmt.Println("  overseer api          - Start the HTTP API server")
		fmt.Println("  overseer scan <org>   - Scan a GitHub organization")
		fmt.Println("  overseer help         - Show this help message")
		return
	}

	command := args[0]
	
	switch command {
	case "api":
		startAPIServer()
	case "scan":
		if len(args) < 2 {
			fmt.Println("Error: Organization name required for scan command")
			fmt.Println("Usage: overseer scan <organization>")
			os.Exit(1)
		}
		scanOrganization(args[1])
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Use 'overseer help' for usage information")
		os.Exit(1)
	}
}

// startAPIServer starts the HTTP API server
func startAPIServer() {
	fmt.Println("🚀 Starting GitHub Codeowners Visualization API Server")
	
	// Load configuration
	config := getDefaultConfig()
	
	// Create database connection
	ctx := context.Background()
	conn, err := createConnection(ctx, config)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer closeConnection(ctx, conn)
	
	// Verify database connection
	if err := verifyConnection(ctx, conn); err != nil {
		fmt.Printf("❌ Database connection verification failed: %v\n", err)
		os.Exit(1)
	}
	
	// Run migrations
	if err := runMigrationsUp(ctx, conn); err != nil {
		fmt.Printf("⚠️ Database migrations failed: %v\n", err)
		// Continue anyway - might be already migrated
	}
	
	// Create HTTP server
	server := createHTTPServer(8081, conn)
	
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		fmt.Println("\n🛑 Received shutdown signal")
		cancel()
	}()
	
	// Start server
	if err := server.startServer(ctx); err != nil {
		fmt.Printf("❌ Server error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("✅ Server shutdown complete")
}

// scanOrganization scans a GitHub organization
func scanOrganization(orgName string) {
	fmt.Printf("🔍 Scanning GitHub organization: %s\n", orgName)
	
	// Load configuration
	config := getDefaultConfig()
	
	// Create database connection
	ctx := context.Background()
	conn, err := createConnection(ctx, config)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer closeConnection(ctx, conn)
	
	// TODO: Implement GitHub scanning
	// This would use the GitHub orchestrator to scan the organization
	batchReq := BatchRequest{
		Organization: orgName,
		MaxRepos:     100,
		MaxTeams:     50,
	}
	
	fmt.Printf("📊 Scan configuration: %+v\n", batchReq)
	fmt.Println("🚧 GitHub scanning implementation coming soon...")
	fmt.Println("✅ Scan request processed")
}

// showHelp displays help information
func showHelp() {
	fmt.Println("GitHub Codeowners Visualization Tool")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  api                   Start the HTTP API server on port 8081")
	fmt.Println("  scan <organization>   Scan a GitHub organization and store data")
	fmt.Println("  help                  Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  overseer api                  # Start API server")
	fmt.Println("  overseer scan microsoft       # Scan Microsoft organization")
	fmt.Println("  overseer scan google          # Scan Google organization")
	fmt.Println()
	fmt.Println("ENVIRONMENT VARIABLES:")
	fmt.Println("  GITHUB_TOKEN          GitHub personal access token (required for scanning)")
	fmt.Println("  NEO4J_URI            Neo4j database URI (default: bolt://localhost:7687)")
	fmt.Println("  NEO4J_USERNAME       Neo4j username (default: neo4j)")
	fmt.Println("  NEO4J_PASSWORD       Neo4j password (default: password)")
	fmt.Println()
	fmt.Println("WEB INTERFACE:")
	fmt.Println("  After starting the API server, open http://localhost:3000 in your browser")
	fmt.Println("  to access the interactive visualization interface.")
}