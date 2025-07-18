# Scan Endpoint Test Suite
# Tests for /api/scan/{org} endpoint to verify organization scanning functionality

# Test 1: Basic organization scan with valid org
POST {{base_url}}/api/scan/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.summary" exists
jsonpath "$.summary.total_repos" exists
jsonpath "$.summary.repos_with_codeowners" exists
jsonpath "$.summary.total_teams" exists
jsonpath "$.summary.total_topics" exists
jsonpath "$.summary.unique_owners" exists
jsonpath "$.summary.api_calls_used" exists
jsonpath "$.summary.processing_time_ms" exists
jsonpath "$.errors" exists
jsonpath "$.data" exists
duration < {{response_time_threshold}}

# Test 2: Scan with query parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos={{max_repos}}&max_teams={{max_teams}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.summary.total_repos" <= {{max_repos}}
jsonpath "$.summary.total_teams" <= {{max_teams}}

# Test 3: Scan with use_topics parameter
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.summary.total_topics" >= 0

# Test 4: Scan with minimal parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_repos" <= 1
jsonpath "$.summary.total_teams" <= 1

# Test 5: Scan response structure validation
POST {{base_url}}/api/scan/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" exists
jsonpath "$.organization" exists
jsonpath "$.summary" exists
jsonpath "$.summary.total_repos" exists
jsonpath "$.summary.repos_with_codeowners" exists
jsonpath "$.summary.total_teams" exists
jsonpath "$.summary.total_topics" exists
jsonpath "$.summary.unique_owners" exists
jsonpath "$.summary.api_calls_used" exists
jsonpath "$.summary.processing_time_ms" exists
jsonpath "$.errors" exists
jsonpath "$.data" exists

# Test 6: Scan with invalid organization name
POST {{base_url}}/api/scan/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == false
jsonpath "$.errors" exists
jsonpath "$.errors" count > 0

# Test 7: Scan with empty organization name
POST {{base_url}}/api/scan/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 8: Scan with special characters in org name
POST {{base_url}}/api/scan/test-org-123
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "test-org-123"

# Test 9: Scan with numeric limits validation
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=5&max_teams=3
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_repos" <= 5
jsonpath "$.summary.total_teams" <= 3

# Test 10: Scan with boolean parameter validation
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true

# Test 11: Scan response timing validation
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.summary.processing_time_ms" > 0
jsonpath "$.summary.processing_time_ms" exists

# Test 12: Scan with maximum parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1000&max_teams=500&use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true

# Test 13: Scan data structure validation
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.organization" exists
jsonpath "$.data.repositories" exists
jsonpath "$.data.repositories" exists
jsonpath "$.data.teams" exists
jsonpath "$.data.teams" exists
jsonpath "$.data.topics" exists
jsonpath "$.data.topics" exists
jsonpath "$.data.codeowners" exists
jsonpath "$.data.codeowners" exists

# Test 14: Scan summary calculations
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=3
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.summary.total_repos" >= 0
jsonpath "$.summary.repos_with_codeowners" >= 0
jsonpath "$.summary.repos_with_codeowners" exists
jsonpath "$.summary.api_calls_used" > 0

# Test 15: Scan with case-insensitive org name
POST {{base_url}}/api/scan/GOLANG
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "GOLANG"

# Test 16: Scan unique owners validation
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.summary.unique_owners" exists
jsonpath "$.summary.unique_owners" count >= 0

# Test 17: Scan error handling
POST {{base_url}}/api/scan/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.errors" exists
# Errors array should exist even if empty

# Test 18: Scan with various max_repos values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.summary.total_repos" <= 1

# Test 19: Scan with zero limits (should use defaults)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=0&max_teams=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true

# Test 20: Scan response completeness
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" exists
jsonpath "$.organization" exists
jsonpath "$.summary" exists
jsonpath "$.errors" exists
jsonpath "$.data" exists