# Graph Endpoint Test Suite - Core Functionality Only
# Simplified tests for /api/graph/{org} endpoint

# Test 1: Basic graph retrieval for valid organization
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

# Test 2: Graph with invalid organization returns empty graph
GET {{base_url}}/api/graph/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists
jsonpath "$.data.nodes" count == 0
jsonpath "$.data.edges" count == 0