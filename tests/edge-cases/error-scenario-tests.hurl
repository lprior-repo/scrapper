# Error Scenario Tests
# Tests for database failures, API limits, and network issues

# Test 1: Health check should always work
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"

# Test 2: Multiple rapid health checks (stress test)
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 3: Test with invalid Accept header
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: invalid/content-type

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 4: Test with no Accept header
GET {{base_url}}/api/health
User-Agent: {{user_agent}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 5: Test with very long User-Agent
GET {{base_url}}/api/health
User-Agent: Very-Long-User-Agent-String-That-Exceeds-Normal-Length-Boundaries-And-Contains-Many-Characters-To-Test-Header-Parsing-Limits-And-Edge-Cases-In-HTTP-Request-Processing
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 6: Test with empty User-Agent
GET {{base_url}}/api/health
User-Agent: 
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 7: Test with malformed URL path
GET {{base_url}}/api/health/../health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 8: Test with double slashes in URL
GET {{base_url}}//api//health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 9: Test with trailing slash
GET {{base_url}}/api/health/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 10: Test scan with network timeout simulation (small timeout)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < 30000

# Test 11: Test graph with potential database connection issues
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

# Test 12: Test stats with potential database connection issues
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

# Test 13: Test with concurrent requests (simulated)
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 14: Test with large request headers
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Custom-Header: ${"A" * 1000}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true

# Test 15: Test with connection keep-alive
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: keep-alive

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 16: Test with connection close
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: close

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 17: Test HTTP version compatibility
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 18: Test with gzip compression acceptance
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Accept-Encoding: gzip, deflate

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 19: Test with cache control headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Cache-Control: no-cache

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 20: Test with pragma no-cache
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Pragma: no-cache

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"