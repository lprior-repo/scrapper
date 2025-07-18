# Concurrent Access Tests
# Tests for multiple simultaneous operations and resource contention

# Test 1: Multiple health checks in sequence
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 2: Multiple graph requests for same organization
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

# Test 3: Multiple stats requests for same organization
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

# Test 4: Mixed endpoint requests
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists

GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"

# Test 5: Graph requests with different parameters
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists

GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists

# Test 6: Multiple OpenAPI doc requests
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"

GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"

# Test 7: Rapid successive requests with timing
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

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
duration < 1000

# Test 8: Different User-Agent concurrent requests
GET {{base_url}}/api/health
User-Agent: Browser-Client/1.0
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: API-Client/2.0
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: Mobile-App/1.5
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 9: Resource intensive requests in sequence
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
duration < 10000

GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" == "{{test_org_small}}"
duration < 5000

# Test 10: Same endpoint with different accept headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: application/json

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: application/json, text/plain, */*

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 11: Connection persistence test
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: keep-alive

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Connection: keep-alive

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 12: Large response handling
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists
bytes count > 100

# Test 13: Multiple documentation requests
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"

GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"

# Test 14: Error handling under load
GET {{base_url}}/api/graph/non-existent-org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

GET {{base_url}}/api/stats/non-existent-org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 15: Mixed success and failure requests
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

GET {{base_url}}/api/graph/non-existent-org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists