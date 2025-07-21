package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"gofr.dev/pkg/gofr"
)

func main() {
	if shouldHandleCommand() {
		return
	}

	app := gofr.New()
	ctx := context.Background()

	deps, err := createAppDependencies(ctx)
	if err != nil {
		app.Logger().Fatalf("Failed to create app dependencies: %v", err)
	}

	logApplicationStartup(app, deps)
	registerGitHubService(app, deps.Config.GitHub)

	handler := NewAppHandler(deps)
	setupGracefulShutdown(app, ctx, deps)
	registerAPIRoutes(app, handler)
	logServerReady(app, deps)

	app.Run()
}

// shouldHandleCommand handles command line arguments
func shouldHandleCommand() bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--cleanup", "cleanup":
			emergencyCleanup()
			return true
		case "api":
			return false
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: api, --cleanup, cleanup")
			return true
		}
	}
	return false
}

// logApplicationStartup logs application startup information
func logApplicationStartup(app *gofr.App, deps *AppDependencies) {
	app.Logger().Infof("Starting GitHub Codeowners Visualization API - service_name=codeowners-scanner service_version=1.0.0 environment=%s port=%d github_base_url=%s neo4j_uri=%s startup_time=%s",
		deps.Config.Environment,
		deps.Config.Port,
		deps.Config.GitHub.BaseURL,
		deps.Config.Neo4j.URI,
		time.Now().UTC().Format(time.RFC3339),
	)
}

// registerGitHubService registers GitHub as an HTTP service
func registerGitHubService(app *gofr.App, config GitHubConfig) {
	RegisterGitHubService(app, GitHubServiceConfig{
		Token:        config.Token,
		BaseURL:      config.BaseURL,
		UserAgent:    config.UserAgent,
		Timeout:      config.Timeout,
		MaxRetries:   config.MaxRetries,
		RateLimitMin: config.RateLimitMin,
	})
}

// setupGracefulShutdown sets up graceful shutdown handling
func setupGracefulShutdown(app *gofr.App, ctx context.Context, deps *AppDependencies) {
	defer func() {
		startTime := time.Now()
		app.Logger().Infof("Starting graceful shutdown - component=main operation=shutdown")
		
		if err := cleanupAppDependencies(ctx, deps); err != nil {
			app.Logger().Errorf("Failed to cleanup dependencies: %v - component=main operation=cleanup_dependencies severity=medium user_impact=cleanup_incomplete", err)
		} else {
			app.Logger().Infof("Dependencies cleaned up successfully - component=main operation=cleanup_dependencies")
		}
		
		duration := time.Since(startTime)
		app.Logger().Infof("Graceful shutdown completed in %v - component=main operation=graceful_shutdown", duration)
	}()
}

// registerAPIRoutes registers all API routes
func registerAPIRoutes(app *gofr.App, handler *AppHandler) {
	app.POST("/api/scan/{org}", handler.handleScanOrganization)
	app.GET("/api/graph/{org}", handler.handleGetGraph)
	app.GET("/api/stats/{org}", handler.handleGetStats)
	app.GET("/api/health", handler.handleHealth)
	app.GET("/api/docs", handler.handleOpenAPI)
	app.GET("/api/openapi.yaml", handler.handleOpenAPISpec)
}

// logServerReady logs server ready information
func logServerReady(app *gofr.App, deps *AppDependencies) {
	app.Logger().Infof("API server routes registered successfully - component=main operation=register_routes routes_count=6 api_endpoints=[/api/scan/{org},/api/graph/{org},/api/stats/{org},/api/health] docs_endpoints=[/api/docs,/api/openapi.yaml]")
	app.Logger().Infof("GitHub Codeowners Visualization API starting on port %d - component=main operation=start_server ready=true", deps.Config.Port)
}

