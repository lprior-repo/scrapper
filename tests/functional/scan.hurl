# Scan Endpoint Test Suite - Core Functionality Only
# Simplified tests for /api/scan/{org} endpoint

# Test 1: Basic organization scan with valid org
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 201
Content-Type: application/json
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.organization" == "{{test_org_small}}"
jsonpath "$.data.summary" exists
jsonpath "$.data.summary.total_repos" exists
jsonpath "$.data.data" exists

# Test 2: Scan with invalid organization name
POST {{base_url}}/api/scan/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 3: Scan with minimal parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 201
Content-Type: application/json
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.summary.total_repos" <= 1