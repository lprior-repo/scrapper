# Scan Endpoint Edge Cases
# Tests for /api/scan/{org} endpoint boundary conditions and edge cases

# Test 1: Scan with extremely long organization name (should return 500 as not found)
POST {{base_url}}/api/scan/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 2: Scan with organization name containing only numbers (should work)
POST {{base_url}}/api/scan/123456789
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 3: Scan with organization name starting with hyphen (should work)
POST {{base_url}}/api/scan/-invalid-org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 4: Scan with organization name ending with hyphen (should work)
POST {{base_url}}/api/scan/invalid-org-
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 5: Scan with max_repos exceeding API limit (should work but use default)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=10000
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 6: Scan with max_teams exceeding API limit (should work but use default)
POST {{base_url}}/api/scan/{{test_org_small}}?max_teams=10000
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 7: Scan with negative max_repos (should work but use default)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=-5
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 8: Scan with negative max_teams (should work but use default)
POST {{base_url}}/api/scan/{{test_org_small}}?max_teams=-5
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 9: Scan with use_topics as invalid string (should work but use default)
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=invalid
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 10: Scan with organization name containing Unicode characters (should fail)
POST {{base_url}}/api/scan/test-org-🚀
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 11: Scan with organization name containing spaces (should fail)
POST {{base_url}}/api/scan/test org space
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 12: Scan with organization name containing uppercase letters (should work)
POST {{base_url}}/api/scan/TEST-ORG-UPPER
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 13: Scan with maximum allowed parameter values (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1000&max_teams=500&use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 14: Scan with minimum allowed parameter values (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1&use_topics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_repos" <= 1
jsonpath "$.summary.total_teams" <= 1

# Test 15: Scan with zero parameter values (should use defaults)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=0&max_teams=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_repos" >= 0
jsonpath "$.summary.total_teams" >= 0

# Test 16: Scan with float parameter values (should work but be truncated)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=5.5&max_teams=3.7
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 17: Scan with empty parameter values (should use defaults)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=&max_teams=
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 18: Scan with duplicate parameter names (should use last value)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=5&max_repos=10
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 19: Scan with case-sensitive parameter names (should use defaults for unknown params)
POST {{base_url}}/api/scan/{{test_org_small}}?MAX_REPOS=5&MAX_TEAMS=3
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 20: Scan with boolean parameter variations (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 21: Scan with boolean parameter variations (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 22: Scan with boolean parameter variations (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 23: Scan with boolean parameter variations (should work)
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"

# Test 24: Scan with organization name at maximum GitHub length (39 characters)
POST {{base_url}}/api/scan/test-organization-with-maximum-length39
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Test 25: Scan with organization name at minimum length (1 character)
POST {{base_url}}/api/scan/a
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"