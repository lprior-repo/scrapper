# GitHub Codeowners API Test Suite
# This file contains comprehensive tests for all API endpoints

# Test 1: Health Check Endpoint
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
jsonpath "$.data.version" == "1.0.0"

# Test 2: API Documentation Endpoint
GET http://localhost:8081/api/docs
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"

# Test 3: OpenAPI Specification Endpoint
GET http://localhost:8081/api/openapi.yaml
HTTP 200
[Asserts]
body contains "openapi: 3.0.3"
body contains "GitHub Codeowners Visualization API"

# Test 4: Scan Organization Endpoint - Valid Organization
POST http://localhost:8081/api/scan/microsoft
HTTP 200
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.organization" == "microsoft"
jsonpath "$.data.summary" exists
jsonpath "$.data.summary.total_repos" >= 0
jsonpath "$.data.summary.total_teams" >= 0

# Test 5: Scan Organization Endpoint - With Query Parameters
POST http://localhost:8081/api/scan/microsoft?max_repos=5&max_teams=5
HTTP 200
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.organization" == "microsoft"
jsonpath "$.data.summary.total_repos" <= 5

# Test 6: Scan Organization Endpoint - Invalid Organization
POST http://localhost:8081/api/scan/nonexistent-org-12345
HTTP 400

# Test 7: Get Organization Graph - Valid Organization (after scan)
GET http://localhost:8081/api/graph/microsoft
HTTP 200
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

# Test 8: Get Organization Graph - Non-existent Organization
GET http://localhost:8081/api/graph/nonexistent-org-12345
HTTP 400

# Test 9: Get Organization Stats - Valid Organization (after scan)
GET http://localhost:8081/api/stats/microsoft
HTTP 200
[Asserts]
jsonpath "$.data.organization" == "microsoft"
jsonpath "$.data.total_repositories" >= 0
jsonpath "$.data.total_teams" >= 0
jsonpath "$.data.total_users" >= 0

# Test 10: Get Organization Stats - Non-existent Organization
GET http://localhost:8081/api/stats/nonexistent-org-12345
HTTP 400

# Test 11: Missing Path Parameter Test
POST http://localhost:8081/api/scan/
HTTP 404

# Test 12: Invalid HTTP Method Test
DELETE http://localhost:8081/api/health
HTTP 405

# Test 13: CORS Headers Test (if applicable)
OPTIONS http://localhost:8081/api/health
HTTP 200

# Test 14: Large Organization Test (with limits)
POST http://localhost:8081/api/scan/google?max_repos=3&max_teams=3
HTTP 200
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.organization" == "google"

# Test 15: Performance Test - Multiple Rapid Requests
GET http://localhost:8081/api/health
HTTP 200

GET http://localhost:8081/api/health
HTTP 200

GET http://localhost:8081/api/health
HTTP 200