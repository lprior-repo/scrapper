# Global Test Configuration for GitHub Codeowners API
# This file contains shared variables and base configurations for all test suites

# Base URL for the API
base_url=http://localhost:8081

# Test organization names
test_org=microsoft
test_org_small=golang
test_org_invalid=thisorgdoesnotexist12345

# Test limits
max_repos=10
max_teams=5

# Performance thresholds (milliseconds)
response_time_threshold=60000
health_check_threshold=1000

# GitHub API rate limits
github_rate_limit_min=100

# Test data patterns
valid_org_pattern=^[a-zA-Z0-9]([a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$

# Content type expectations
json_content_type=application/json
yaml_content_type=application/yaml
html_content_type=text/html

# Test environment detection
environment=development

# Common HTTP headers
user_agent=Hurl-Test-Suite/1.0
accept_json=application/json
accept_html=text/html
accept_yaml=application/yaml

# Security test patterns
sql_injection_pattern='; DROP TABLE users; --
xss_pattern=<script>alert('xss')</script>
path_traversal_pattern=../../../etc/passwd

# Load test configurations
concurrent_requests=10
test_duration_seconds=30

# Error message expectations
missing_param_error=missing parameter
not_found_error=not found
invalid_input_error=invalid input

# Database connection test
neo4j_health_check=connected

# API version
api_version=1.0.0

# Test timeouts
default_timeout=30
long_timeout=120
short_timeout=5