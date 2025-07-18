# Performance Tests for Scan Endpoint
# Tests for /api/scan/{org} endpoint to verify response times and performance under load

# Test 1: Scan endpoint response time baseline with minimal parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}

# Test 2: Scan endpoint response time with moderate parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=5&max_teams=3
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}

# Test 3: Scan endpoint response time with various parameter combinations
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=2&use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}

# Test 4: Scan endpoint response time with large headers
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Large-Header: ${"A" * 1000}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 5: Scan endpoint response time with minimal headers
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 6: Scan endpoint response time consistency check
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 7: Scan endpoint response time with different User-Agent strings
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 8: Scan endpoint cold start performance
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < 10000  # Allow longer for cold start

# Test 9: Scan endpoint response size validation
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}
bytes count > 100  # Should have substantial response

# Test 10: Scan endpoint performance with Keep-Alive
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: keep-alive

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 11: Scan endpoint performance with gzip compression
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Accept-Encoding: gzip, deflate

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 12: Scan endpoint performance with cache headers
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Cache-Control: no-cache

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 13: Scan endpoint processing time validation
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
jsonpath "$.summary.processing_time_ms" < 30000  # Should not take more than 30 seconds

# Test 14: Scan endpoint API calls efficiency
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.api_calls_used" > 0
jsonpath "$.summary.api_calls_used" < 50  # Should be efficient with API calls

# Test 15: Scan endpoint performance with different parameter values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=3&max_teams=2&use_topics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}

# Test 16: Scan endpoint performance with topics enabled
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}

# Test 17: Scan endpoint performance with zero limits (should use defaults)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=0&max_teams=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}

# Test 18: Scan endpoint performance comparison between small and large limits
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < 3000  # Should be faster with smaller limits

# Test 19: Scan endpoint response time with different Accept headers
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
duration < {{response_time_threshold}}

# Test 20: Scan endpoint memory efficiency test
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
duration < {{response_time_threshold}}
# Response should be well-structured and not excessively large
bytes count < 100000  # Should not be more than 100KB