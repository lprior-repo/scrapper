# Data Consistency Tests  
# Tests for data consistency between endpoints and response validation

# Test 1: Health check consistency across multiple calls
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
jsonpath "$.data.version" == "1.0.0"
[Captures]
first_timestamp: jsonpath "$.data.timestamp"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
jsonpath "$.data.version" == "1.0.0"
[Captures]
second_timestamp: jsonpath "$.data.timestamp"

# Test 2: Graph endpoint data structure consistency
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0
[Captures]
nodes_count: jsonpath "$.data.nodes" count
edges_count: jsonpath "$.data.edges" count

# Test 3: Graph endpoint with topics vs without topics consistency
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0
[Captures]
topics_nodes_count: jsonpath "$.data.nodes" count
topics_edges_count: jsonpath "$.data.edges" count

GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0
[Captures]
no_topics_nodes_count: jsonpath "$.data.nodes" count
no_topics_edges_count: jsonpath "$.data.edges" count

# Test 4: OpenAPI spec consistency
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"
body contains "info:"
body contains "title: GitHub Codeowners Visualization API"
body contains "version: 1.0.0"
body contains "paths:"
body contains "/api/health:"
body contains "/api/graph/{org}:"
body contains "/api/scan/{org}:"
body contains "/api/stats/{org}:"
body contains "useTopics"

# Test 5: Documentation consistency
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"
body contains "GitHub Codeowners API Documentation"
body contains "swagger-ui"
body contains "/api/openapi.yaml"

# Test 6: Graph node structure consistency
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes[0].id" exists
jsonpath "$.data.nodes[0].type" exists
jsonpath "$.data.nodes[0].label" exists
jsonpath "$.data.nodes[0].data" exists
jsonpath "$.data.nodes[0].position" exists
jsonpath "$.data.nodes[0].position.x" exists
jsonpath "$.data.nodes[0].position.y" exists

# Test 7: Graph edge structure consistency
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.edges[0].id" exists
jsonpath "$.data.edges[0].source" exists
jsonpath "$.data.edges[0].target" exists
jsonpath "$.data.edges[0].type" exists
jsonpath "$.data.edges[0].label" exists

# Test 8: Response header consistency
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
header "Content-Type" == "application/json"
header "Access-Control-Allow-Origin" == "*"
header "Access-Control-Allow-Methods" contains "GET"
header "Access-Control-Allow-Headers" contains "Content-Type"

# Test 9: Graph response data types consistency
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes[0].position.x" isNumber
jsonpath "$.data.nodes[0].position.y" isNumber
jsonpath "$.data.nodes[0].id" isString
jsonpath "$.data.nodes[0].type" isString
jsonpath "$.data.nodes[0].label" isString

# Test 10: Error response consistency
GET {{base_url}}/api/graph/non-existent-org-12345
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error.message" exists
jsonpath "$.error.message" isString

# Test 11: Multiple graph requests return consistent data
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists
[Captures]
first_response_nodes: jsonpath "$.data.nodes" count
first_response_edges: jsonpath "$.data.edges" count

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists
[Captures]
second_response_nodes: jsonpath "$.data.nodes" count
second_response_edges: jsonpath "$.data.edges" count

# Test 12: Graph node types are valid
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes[0].type" matches "^(organization|repository|team|topic|user|codeowner)$"

# Test 13: Graph edge types are valid
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.edges[0].type" matches "^(owns|member_of|codeowner|has_topic|repo_topic|belongs_to|manages|contributes_to|depends_on)$"

# Test 14: Health check database status consistency
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.database" == "connected"
jsonpath "$.data.status" == "healthy"

# Test 15: OpenAPI spec version consistency with health check
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "version: 1.0.0"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.version" == "1.0.0"

# Test 16: Graph position coordinates are valid
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes[0].position.x" >= 0
jsonpath "$.data.nodes[0].position.y" >= 0

# Test 17: Graph edge relationships are valid
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.edges[0].source" isString
jsonpath "$.data.edges[0].target" isString
jsonpath "$.data.edges[0].source" count > 0
jsonpath "$.data.edges[0].target" count > 0

# Test 18: Response timing consistency
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < 1000

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < 1000

# Test 19: CORS headers consistency
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
header "Access-Control-Allow-Origin" == "*"
header "Access-Control-Allow-Methods" contains "GET"

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
header "Access-Control-Allow-Origin" == "*"
header "Access-Control-Allow-Methods" contains "GET"

# Test 20: Content-Type header consistency
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
header "Content-Type" == "application/json"

GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
header "Content-Type" == "application/yaml"

GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
header "Content-Type" == "text/html"