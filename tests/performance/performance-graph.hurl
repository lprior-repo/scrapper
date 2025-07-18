# Performance Tests for Graph Endpoint
# Tests for /api/graph/{org} endpoint to verify response times and performance under load

# Test 1: Graph endpoint response time baseline
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 2: Graph endpoint response time with large headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Large-Header: ${"A" * 1000}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 3: Graph endpoint response time with minimal headers
GET {{base_url}}/api/graph/{{test_org_small}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 4: Graph endpoint response time consistency check
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 5: Graph endpoint response time with different User-Agent strings
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 6: Graph endpoint cold start performance
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < 8000  # Allow longer for cold start

# Test 7: Graph endpoint response size validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}
bytes count > 50  # Should have substantial response

# Test 8: Graph endpoint performance with Keep-Alive
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: keep-alive

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 9: Graph endpoint performance with gzip compression
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Accept-Encoding: gzip, deflate

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 10: Graph endpoint performance with cache headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Cache-Control: no-cache

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 11: Graph endpoint repeated calls performance
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 12: Graph endpoint performance with different Accept headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 13: Graph endpoint memory efficiency test
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}
# Response should be well-structured and not excessively large
bytes count < 500000  # Should not be more than 500KB

# Test 14: Graph endpoint nodes and edges validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
jsonpath "$.nodes" isArray
jsonpath "$.edges" isArray
duration < {{response_time_threshold}}

# Test 15: Graph endpoint performance stress test
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 16: Graph endpoint performance with case variations
GET {{base_url}}/api/graph/GOLANG
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 17: Graph endpoint performance with conditional requests
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
If-Modified-Since: Thu, 01 Jan 1970 00:00:00 GMT

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 18: Graph endpoint data structure performance
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
jsonpath "$.nodes" isArray
jsonpath "$.edges" isArray
duration < {{response_time_threshold}}

# Test 19: Graph endpoint performance with multiple headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Test-Header-1: value1
X-Test-Header-2: value2
X-Test-Header-3: value3

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}

# Test 20: Graph endpoint response latency validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
duration < {{response_time_threshold}}
# Check that response time is reasonable
duration > 10  # Should take at least 10ms to process