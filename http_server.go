package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/samber/lo"
)

// HTTPServer represents the HTTP server configuration
type HTTPServer struct {
	Port           int           `json:"port"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes"`
	GraphConn      GraphConnection
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ScanRequest represents a request to scan an organization
type ScanRequest struct {
	Organization string `json:"organization"`
	GitHubToken  string `json:"github_token"`
	MaxRepos     int    `json:"max_repos,omitempty"`
	MaxTeams     int    `json:"max_teams,omitempty"`
}

// createHTTPServer creates a new HTTP server instance (Pure Core)
func createHTTPServer(port int, graphConn GraphConnection) HTTPServer {
	if port <= 0 || port > 65535 {
		panic("Invalid port number")
	}

	return HTTPServer{
		Port:           port,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
		GraphConn:      graphConn,
	}
}

// createSuccessResponse creates a successful API response (Pure Core)
func createSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// createErrorResponse creates an error API response (Pure Core)
func createErrorResponse(code, message, details string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// validateScanRequest validates a scan request (Pure Core)
func validateScanRequest(req ScanRequest) error {
	if req.Organization == "" {
		return createRequiredFieldError("organization")
	}
	if req.GitHubToken == "" {
		return createRequiredFieldError("github_token")
	}
	if req.MaxRepos < 0 {
		return createValidationError("max_repos", "cannot be negative")
	}
	if req.MaxTeams < 0 {
		return createValidationError("max_teams", "cannot be negative")
	}
	return nil
}

// writeJSONResponse writes a JSON response to the HTTP writer (Impure Shell)
func writeJSONResponse(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback error response
		fmt.Fprintf(w, `{"success":false,"error":{"code":"ENCODING_ERROR","message":"Failed to encode response"}}`)
	}
}

// parseIntParam parses an integer parameter from URL (Pure Core)
func parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultValue
}

// HTTP Handlers

// healthHandler handles health check requests (Impure Shell)
func (s *HTTPServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check database connection
	if err := verifyConnection(ctx, s.GraphConn); err != nil {
		response := createErrorResponse("DATABASE_UNHEALTHY", "Database connection failed", err.Error())
		writeJSONResponse(w, http.StatusServiceUnavailable, response)
		return
	}

	healthData := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"database":  "connected",
		"version":   "1.0.0",
	}

	response := createSuccessResponse(healthData)
	writeJSONResponse(w, http.StatusOK, response)
}

// scanOrganizationHandler handles organization scanning requests (Impure Shell)
func (s *HTTPServer) scanOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := createErrorResponse("METHOD_NOT_ALLOWED", "Only POST method is allowed", "")
		writeJSONResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	var scanReq ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&scanReq); err != nil {
		response := createErrorResponse("INVALID_JSON", "Invalid JSON in request body", err.Error())
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := validateScanRequest(scanReq); err != nil {
		if appErr, ok := err.(AppError); ok {
			response := createErrorResponse(appErr.Code, appErr.Message, appErr.Details)
			writeJSONResponse(w, http.StatusBadRequest, response)
			return
		}
		response := createErrorResponse("VALIDATION_ERROR", "Request validation failed", err.Error())
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Create batch request
	batchReq := BatchRequest{
		Organization: scanReq.Organization,
		MaxRepos:     lo.Ternary(scanReq.MaxRepos > 0, scanReq.MaxRepos, 100),
		MaxTeams:     lo.Ternary(scanReq.MaxTeams > 0, scanReq.MaxTeams, 50),
	}

	// TODO: Implement actual GitHub scanning
	// This would call the GitHub orchestrator to scan the organization
	scanResult := map[string]interface{}{
		"organization": scanReq.Organization,
		"status":       "scanning_started",
		"batch_config": batchReq,
		"message":      "Organization scan initiated successfully",
	}

	response := createSuccessResponse(scanResult)
	writeJSONResponse(w, http.StatusAccepted, response)
}

// getOrganizationHandler handles requests for organization data (Impure Shell)
func (s *HTTPServer) getOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgName := vars["org"]

	if orgName == "" {
		response := createErrorResponse("MISSING_PARAMETER", "Organization name is required", "")
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// TODO: Implement actual data fetching from graph database
	// This would call graph queries to get organization data
	orgData := map[string]interface{}{
		"name":         orgName,
		"repositories": []map[string]interface{}{},
		"teams":        []map[string]interface{}{},
		"statistics": map[string]interface{}{
			"total_repos":            0,
			"repos_with_codeowners":  0,
			"total_teams":            0,
			"unique_owners":          0,
		},
	}

	response := createSuccessResponse(orgData)
	writeJSONResponse(w, http.StatusOK, response)
}

// getGraphDataHandler handles requests for graph visualization data (Impure Shell)
func (s *HTTPServer) getGraphDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgName := vars["org"]

	if orgName == "" {
		response := createErrorResponse("MISSING_PARAMETER", "Organization name is required", "")
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// TODO: Implement actual graph data fetching
	// This would call graph queries to get nodes and relationships
	graphData := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{
				"id":   "org-" + orgName,
				"type": "organization",
				"name": orgName,
			},
		},
		"edges": []map[string]interface{}{},
		"metadata": map[string]interface{}{
			"total_nodes": 1,
			"total_edges": 0,
			"organization": orgName,
		},
	}

	response := createSuccessResponse(graphData)
	writeJSONResponse(w, http.StatusOK, response)
}

// getRepositoryHandler handles requests for repository details (Impure Shell)
func (s *HTTPServer) getRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgName := vars["org"]
	repoName := vars["repo"]

	if orgName == "" || repoName == "" {
		response := createErrorResponse("MISSING_PARAMETER", "Organization and repository names are required", "")
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	fullName := orgName + "/" + repoName

	// TODO: Implement actual repository data fetching
	repoData := map[string]interface{}{
		"name":                orgName,
		"full_name":           fullName,
		"has_codeowners_file": false,
		"codeowners_content":  "",
		"codeowners_entries":  []CodeownersEntry{},
		"owners":              []string{},
	}

	response := createSuccessResponse(repoData)
	writeJSONResponse(w, http.StatusOK, response)
}

// corsMiddleware adds CORS headers (Impure Shell)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests (Impure Shell)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		fmt.Printf("[%s] %s %s - %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			duration)
	})
}

// setupRoutes configures HTTP routes (Impure Shell)
func (s *HTTPServer) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Health and utility endpoints
	api.HandleFunc("/health", s.healthHandler).Methods("GET")

	// Organization endpoints
	api.HandleFunc("/scan/{org}", s.scanOrganizationHandler).Methods("POST")
	api.HandleFunc("/organizations/{org}", s.getOrganizationHandler).Methods("GET")
	api.HandleFunc("/graph/{org}", s.getGraphDataHandler).Methods("GET")

	// Repository endpoints
	api.HandleFunc("/repositories/{org}/{repo}", s.getRepositoryHandler).Methods("GET")

	// Static file serving for frontend (if needed)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/dist/")))

	return r
}

// startServer starts the HTTP server (Impure Shell)
func (s *HTTPServer) startServer(ctx context.Context) error {
	router := s.setupRoutes()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", s.Port),
		Handler:        router,
		ReadTimeout:    s.ReadTimeout,
		WriteTimeout:   s.WriteTimeout,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("🔧 HTTP server starting on port %d\n", s.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("🛑 Shutting down HTTP server...")
	return server.Shutdown(shutdownCtx)
}