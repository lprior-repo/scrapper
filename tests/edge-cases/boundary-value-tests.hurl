# Boundary Value Tests
# Tests for parameter limits and edge conditions

# Test 1: Scan with maximum repos parameter (1000)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1000
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 2: Scan with maximum teams parameter (500)
POST {{base_url}}/api/scan/{{test_org_small}}?max_teams=500
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 3: Scan with minimum repos parameter (1)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_repos" <= 1

# Test 4: Scan with minimum teams parameter (1)
POST {{base_url}}/api/scan/{{test_org_small}}?max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_teams" <= 1

# Test 5: Scan with maximum organization name length (39 chars)
POST {{base_url}}/api/scan/test-organization-name-at-max-length-39
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 6: Scan with minimum organization name length (1 char)
POST {{base_url}}/api/scan/a
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 7: Scan with exactly valid organization name (github format)
POST {{base_url}}/api/scan/valid-org-name-123
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 8: Scan with zero repos (should use default)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 9: Scan with zero teams (should use default)
POST {{base_url}}/api/scan/{{test_org_small}}?max_teams=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 10: Scan with values just above maximum (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1001&max_teams=501
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 11: Scan with very large integer values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=999999999&max_teams=999999999
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 12: Graph endpoint with non-existent organization
GET {{base_url}}/api/graph/non-existent-org-12345
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "not found"

# Test 13: Stats endpoint with non-existent organization
GET {{base_url}}/api/stats/non-existent-org-12345
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "not found"

# Test 14: Scan with empty organization name (should fail)
POST {{base_url}}/api/scan/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 15: Graph with empty organization name (should fail)
GET {{base_url}}/api/graph/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 16: Stats with empty organization name (should fail)
GET {{base_url}}/api/stats/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 17: Scan with parameters as strings (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 18: Scan with mixed valid/invalid parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=5&invalid_param=value
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 19: Test processing time boundaries
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.processing_time_ms" > 0
jsonpath "$.summary.processing_time_ms" < 60000

# Test 20: Test API calls efficiency boundaries
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.api_calls_used" > 0
jsonpath "$.summary.api_calls_used" < 100