# Stats Endpoint Test Suite
# Tests for /api/stats/{org} endpoint to verify organization statistics retrieval

# Test 1: Basic stats retrieval for valid organization
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists
jsonpath "$.total_teams" exists
jsonpath "$.total_topics" exists
jsonpath "$.total_users" exists
jsonpath "$.total_codeowners" exists
jsonpath "$.codeowner_coverage" exists
jsonpath "$.last_scan_time" exists
duration < {{response_time_threshold}}

# Test 2: Stats structure validation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" isString
jsonpath "$.total_repositories" isInteger
jsonpath "$.total_teams" isInteger
jsonpath "$.total_topics" isInteger
jsonpath "$.total_users" isInteger
jsonpath "$.total_codeowners" isInteger
jsonpath "$.codeowner_coverage" isString
jsonpath "$.last_scan_time" isString

# Test 3: Stats values validation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.total_repositories" >= 0
jsonpath "$.total_teams" >= 0
jsonpath "$.total_topics" >= 0
jsonpath "$.total_users" >= 0
jsonpath "$.total_codeowners" >= 0

# Test 4: Stats with invalid organization
GET {{base_url}}/api/stats/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 5: Stats with empty organization parameter
GET {{base_url}}/api/stats/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 6: Stats coverage format validation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.codeowner_coverage" matches "^\\d+(\\.\\d+)?%$"

# Test 7: Stats timestamp format validation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.last_scan_time" matches "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*$"

# Test 8: Stats response timing
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
duration < 5000  # Should respond within 5 seconds

# Test 9: Stats with special characters in org name
GET {{base_url}}/api/stats/test-org-123
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 10: Stats logical consistency
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.total_codeowners" <= jsonpath "$.total_repositories"

# Test 11: Stats with different Accept headers
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists

# Test 12: Stats response completeness
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$" count == 8  # Should contain exactly 8 fields
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists
jsonpath "$.total_teams" exists
jsonpath "$.total_topics" exists
jsonpath "$.total_users" exists
jsonpath "$.total_codeowners" exists
jsonpath "$.codeowner_coverage" exists
jsonpath "$.last_scan_time" exists

# Test 13: Stats case sensitivity
GET {{base_url}}/api/stats/GOLANG
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 14: Stats with minimal headers
GET {{base_url}}/api/stats/{{test_org_small}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists

# Test 15: Stats numerical bounds
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.total_repositories" >= 0
jsonpath "$.total_teams" >= 0
jsonpath "$.total_topics" >= 0
jsonpath "$.total_users" >= 0
jsonpath "$.total_codeowners" >= 0

# Test 16: Stats coverage calculation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.codeowner_coverage" isString
jsonpath "$.codeowner_coverage" not == ""

# Test 17: Stats timestamp recency
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.last_scan_time" isString
jsonpath "$.last_scan_time" not == ""

# Test 18: Stats organization name matching
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

# Test 19: Stats response size validation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.organization" exists
bytes count > 0
bytes count < 2048  # Should be reasonably sized

# Test 20: Stats consistency across calls
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

# Test 21: Stats with custom User-Agent
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: Custom-Stats-Client/1.0
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

# Test 22: Stats field data types
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.total_repositories" isInteger
jsonpath "$.total_teams" isInteger
jsonpath "$.total_topics" isInteger
jsonpath "$.total_users" isInteger
jsonpath "$.total_codeowners" isInteger

# Test 23: Stats coverage percentage format
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.codeowner_coverage" matches "^(\\d{1,2}(\\.\\d{1,2})?|100(\\.00?)?)%$"

# Test 24: Stats organization field validation
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" isString
jsonpath "$.organization" count > 0
jsonpath "$.organization" matches "^[a-zA-Z0-9]([a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$"

# Test 25: Stats response without extra fields
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$" isObject
jsonpath "$" not contains "extra"
jsonpath "$" not contains "metadata"