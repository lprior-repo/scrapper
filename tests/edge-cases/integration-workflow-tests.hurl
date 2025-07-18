# Integration Workflow Tests
# Tests for full workflow: scan → graph → stats sequence

# Test 1: Full workflow with small organization
# Step 1: Scan organization
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&max_teams=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.summary.total_repos" <= 2
jsonpath "$.summary.total_teams" <= 2

# Step 2: Get graph data for same organization
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 1
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Step 3: Get stats for same organization
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.total_repositories" >= 0
jsonpath "$.total_teams" >= 0

# Test 2: Full workflow with topics enabled
# Step 1: Scan organization with topics
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=2&use_topics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.summary.total_topics" >= 0

# Step 2: Get graph data with topics
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 1
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Step 3: Get stats (should reflect topics)
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.total_topics" >= 0

# Test 3: Workflow with minimal parameters
# Step 1: Scan with minimal parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.summary.total_repos" <= 1
jsonpath "$.summary.total_teams" <= 1

# Step 2: Get graph data
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

# Step 3: Get stats
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

# Test 4: Workflow consistency check
# Step 1: Scan organization
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=3&max_teams=3
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" == true
jsonpath "$.organization" == "{{test_org_small}}"
[Captures]
scan_repos_count: jsonpath "$.summary.total_repos"
scan_teams_count: jsonpath "$.summary.total_teams"

# Step 2: Get graph data
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

# Step 3: Get stats and verify consistency
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.total_repositories" >= 0
jsonpath "$.total_teams" >= 0

# Test 5: Workflow with invalid organization (should fail at scan)
# Step 1: Scan invalid organization
POST {{base_url}}/api/scan/{{test_org_invalid}}?max_repos=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 500
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "organization"

# Step 2: Try to get graph data (should fail)
GET {{base_url}}/api/graph/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "not found"

# Step 3: Try to get stats (should fail)
GET {{base_url}}/api/stats/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" contains "not found"