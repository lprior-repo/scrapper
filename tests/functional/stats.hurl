# Stats Endpoint Test Suite - Core Functionality Only
# Simplified tests for /api/stats/{org} endpoint

# Test 1: Basic stats retrieval for valid organization
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.total_repositories" exists
jsonpath "$.data.total_codeowners" exists

# Test 2: Stats with invalid organization
GET {{base_url}}/api/stats/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists