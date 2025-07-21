package main

// buildOpenAPIHTML returns the HTML for OpenAPI documentation UI
func buildOpenAPIHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>GitHub Codeowners API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api/openapi.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
}

// buildOpenAPISpec returns the OpenAPI specification content
func buildOpenAPISpec() string {
	info := buildOpenAPIInfo()
	servers := buildOpenAPIServers()
	paths := buildOpenAPIPaths()
	components := buildOpenAPIComponents()
	tags := buildOpenAPITags()

	return info + servers + paths + components + tags
}

// buildOpenAPIInfo returns the info section of OpenAPI spec
func buildOpenAPIInfo() string {
	return `openapi: 3.0.3
info:
  title: GitHub Codeowners Visualization API
  description: API for scanning GitHub organizations and visualizing codeowners relationships
  version: 1.0.0
  contact:
    name: API Support
    url: https://github.com/your-org/scrapper
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

`
}

// buildOpenAPIServers returns the servers section of OpenAPI spec
func buildOpenAPIServers() string {
	return `servers:
  - url: http://localhost:8080
    description: Development server

`
}

// buildOpenAPIPaths returns the paths section of OpenAPI spec
func buildOpenAPIPaths() string {
	healthPath := buildHealthPathSpec()
	scanPath := buildScanPathSpec()
	return "paths:\n" + healthPath + scanPath
}

// buildHealthPathSpec returns the health endpoint specification
func buildHealthPathSpec() string {
	return `  /api/health:
    get:
      summary: Health check endpoint
      description: Returns the health status of the API and its dependencies
      operationId: healthCheck
      tags:
        - System
      responses:
        '200':
          description: System is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: object
                    properties:
                      status:
                        type: string
                        example: "healthy"
                      database:
                        type: string
                        example: "connected"
                      version:
                        type: string
                        example: "1.0.0"
                      timestamp:
                        type: string
                        format: date-time
                        example: "2025-07-17T21:08:23-05:00"

`
}

// buildScanPathSpec returns the scan endpoint specification
func buildScanPathSpec() string {
	return `  /api/scan/{org}:
    post:
      summary: Scan GitHub organization
      description: Scans a GitHub organization for repositories, teams, and CODEOWNERS files
      operationId: scanOrganization
      tags:
        - Scanning
      parameters:
        - name: org
          in: path
          required: true
          description: GitHub organization name
          schema:
            type: string
            example: "microsoft"
        - name: max_repos
          in: query
          required: false
          description: Maximum number of repositories to scan
          schema:
            type: integer
            default: 100
            minimum: 1
            maximum: 1000
        - name: max_teams
          in: query
          required: false
          description: Maximum number of teams to scan
          schema:
            type: integer
            default: 50
            minimum: 1
            maximum: 500
        - name: use_topics
          in: query
          required: false
          description: Use repository topics instead of teams for organization
          schema:
            type: boolean
            default: false

`
}

// buildOpenAPIComponents returns the components section of OpenAPI spec
func buildOpenAPIComponents() string {
	return `components:
  schemas:
    ScanResponse:
      type: object
      properties:
        success:
          type: boolean
          description: Whether the scan was successful
        organization:
          type: string
          description: Name of the scanned organization
        summary:
          type: object
          properties:
            total_repos:
              type: integer
            repos_with_codeowners:
              type: integer
            total_teams:
              type: integer
            total_topics:
              type: integer
            unique_owners:
              type: array
              items:
                type: string

`
}

// buildOpenAPITags returns the tags section of OpenAPI spec
func buildOpenAPITags() string {
	return `tags:
  - name: System
    description: System health and status operations
  - name: Scanning
    description: GitHub organization scanning operations`
}