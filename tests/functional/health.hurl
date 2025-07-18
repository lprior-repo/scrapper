# Health Endpoint Test Suite
# Tests for /api/health endpoint to verify system health and readiness

# Test 1: Basic health check - should return 200 with healthy status
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
jsonpath "$.data.version" exists
jsonpath "$.data.timestamp" exists
jsonpath "$.data.timestamp" matches "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*$"
duration < {{health_check_threshold}}

# Test 2: Health check response structure validation
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data.status" exists
jsonpath "$.data.database" exists
jsonpath "$.data.version" exists
jsonpath "$.data.timestamp" exists
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
jsonpath "$.data.version" == "1.0.0"

# Test 3: Health check with different Accept headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 4: Health check with custom User-Agent
GET {{base_url}}/api/health
User-Agent: Custom-Health-Monitor/1.0
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"

# Test 5: Health check response timing
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
duration < 1000  # Should respond within 1 second

# Test 6: Health check with no Accept header
GET {{base_url}}/api/health
User-Agent: {{user_agent}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 7: Health check with multiple concurrent requests
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 8: Health check timestamp validation
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.timestamp" exists
jsonpath "$.data.timestamp" isString
# Timestamp should be recent (within last 5 minutes)
jsonpath "$.data.timestamp" matches "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*$"

# Test 9: Health check with HTTP/1.1
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
header "Content-Type" == "application/json"

# Test 10: Health check response size validation
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
# Response should be reasonably sized (not too large)
bytes count < 1024

# Test 11: Health check with minimal headers
GET {{base_url}}/api/health

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 12: Health check database connection validation
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.database" == "connected"
jsonpath "$.data.status" == "healthy"

# Test 13: Health check version information
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.version" == "1.0.0"
jsonpath "$.data.version" matches "^\\d+\\.\\d+\\.\\d+$"

# Test 14: Health check consistency across multiple calls
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"

# Test 15: Health check with case variations in headers
GET {{base_url}}/api/health
user-agent: {{user_agent}}
accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"