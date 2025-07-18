# Performance Tests for Health Endpoint
# Tests for /api/health endpoint to verify response times and performance under load

# Test 1: Health endpoint response time baseline
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 2: Health endpoint response time under multiple rapid requests
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 3: Health endpoint response time with large headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Large-Header: ${"A" * 1000}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 4: Health endpoint response time with minimal headers
GET {{base_url}}/api/health

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 5: Health endpoint response time with various Accept headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 6: Health endpoint response time consistency check
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 7: Health endpoint response time with different User-Agent strings
GET {{base_url}}/api/health
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 8: Health endpoint cold start performance
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < 2000  # Allow slightly longer for cold start

# Test 9: Health endpoint response size validation
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}
bytes count < 1024  # Should be small response

# Test 10: Health endpoint repeated calls performance
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 11: Health endpoint performance with Keep-Alive
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: keep-alive

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 12: Health endpoint performance with gzip compression
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Accept-Encoding: gzip, deflate

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 13: Health endpoint performance with cache headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Cache-Control: no-cache

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 14: Health endpoint performance with conditional requests
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
If-Modified-Since: Thu, 01 Jan 1970 00:00:00 GMT

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

# Test 15: Health endpoint performance stress test
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < {{health_check_threshold}}