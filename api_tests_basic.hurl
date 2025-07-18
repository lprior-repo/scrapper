# GitHub Codeowners API Basic Test Suite
# This file contains basic tests for working endpoints

# Test 1: Health Check Endpoint
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
jsonpath "$.data.version" == "1.0.0"

# Test 2: OpenAPI Specification Endpoint
GET http://localhost:8081/api/openapi.yaml
HTTP 200
[Asserts]
body contains "openapi: 3.0.3"
body contains "GitHub Codeowners Visualization API"

# Test 3: Health Check Response Structure
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data.status" exists
jsonpath "$.data.database" exists
jsonpath "$.data.version" exists
jsonpath "$.data.timestamp" exists

# Test 4: Health Check with different Accept header
GET http://localhost:8081/api/health
Accept: application/json
HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 5: OpenAPI Response Content
GET http://localhost:8081/api/openapi.yaml
HTTP 200
[Asserts]
body contains "paths:"
body contains "/api/health"
body contains "/api/scan"