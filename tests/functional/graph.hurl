# Graph Endpoint Test Suite
# Tests for /api/graph/{org} endpoint to verify graph data retrieval functionality

# Test 1: Basic graph retrieval for valid organization
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.nodes" isArray
jsonpath "$.edges" exists
jsonpath "$.edges" isArray
duration < {{response_time_threshold}}

# Test 2: Graph structure validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" isArray
jsonpath "$.edges" isArray

# Test 3: Graph nodes structure validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" isArray
jsonpath "$.nodes[0].id" exists
jsonpath "$.nodes[0].type" exists
jsonpath "$.nodes[0].label" exists
jsonpath "$.nodes[0].data" exists
jsonpath "$.nodes[0].position" exists
jsonpath "$.nodes[0].position.x" exists
jsonpath "$.nodes[0].position.y" exists

# Test 4: Graph edges structure validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.edges" isArray
jsonpath "$.edges[0].id" exists
jsonpath "$.edges[0].source" exists
jsonpath "$.edges[0].target" exists
jsonpath "$.edges[0].type" exists
jsonpath "$.edges[0].label" exists

# Test 5: Graph with invalid organization
GET {{base_url}}/api/graph/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 6: Graph with empty organization parameter
GET {{base_url}}/api/graph/
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 7: Graph nodes data types validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes[0].id" isString
jsonpath "$.nodes[0].type" isString
jsonpath "$.nodes[0].label" isString
jsonpath "$.nodes[0].data" isObject
jsonpath "$.nodes[0].position" isObject
jsonpath "$.nodes[0].position.x" isNumber
jsonpath "$.nodes[0].position.y" isNumber

# Test 8: Graph edges data types validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.edges[0].id" isString
jsonpath "$.edges[0].source" isString
jsonpath "$.edges[0].target" isString
jsonpath "$.edges[0].type" isString
jsonpath "$.edges[0].label" isString

# Test 9: Graph response consistency
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" count >= 0
jsonpath "$.edges" count >= 0

# Test 10: Graph with special characters in org name
GET {{base_url}}/api/graph/test-org-123
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 11: Graph response timing
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
duration < 10000  # Should respond within 10 seconds

# Test 12: Graph with case-sensitive org name
GET {{base_url}}/api/graph/GOLANG
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 13: Graph node position validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes[0].position.x" isNumber
jsonpath "$.nodes[0].position.y" isNumber
jsonpath "$.nodes[0].position.x" >= 0
jsonpath "$.nodes[0].position.y" >= 0

# Test 14: Graph edge relationships validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.edges[0].source" isString
jsonpath "$.edges[0].target" isString
jsonpath "$.edges[0].source" != jsonpath "$.edges[0].target"

# Test 15: Graph with different Accept headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 16: Graph response size validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
bytes count > 0

# Test 17: Graph node types validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes[0].type" matches "^(organization|repository|team|topic|user|codeowner)$"

# Test 18: Graph edge types validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.edges[0].type" matches "^(owns|belongs_to|manages|contributes_to|depends_on)$"

# Test 19: Graph nodes uniqueness validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes[0].id" isString
# Note: Uniqueness validation would require custom logic

# Test 20: Graph without data (empty organization)
GET {{base_url}}/api/graph/empty-org-test
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 21: Graph with minimal headers
GET {{base_url}}/api/graph/{{test_org_small}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 22: Graph data structure completeness
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes[0].data" isObject
jsonpath "$.nodes[0].data" count >= 0

# Test 23: Graph edge labels validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.edges[0].label" isString
jsonpath "$.edges[0].label" count > 0

# Test 24: Graph position coordinates validation
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes[0].position.x" isNumber
jsonpath "$.nodes[0].position.y" isNumber

# Test 25: Graph response format consistency
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$" isObject
jsonpath "$.nodes" isArray
jsonpath "$.edges" isArray
jsonpath "$" count == 2  # Should only contain nodes and edges