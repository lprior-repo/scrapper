# Edge Case Tests - Boundary Value Testing (Ultra Fast)
# Tests API endpoints with minimal boundary values

### Test 1: Minimum valid scan (ultra fast)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1&use_topics=false

HTTP 201
[Asserts]
jsonpath "$.data.success" == true

### Test 2: Empty org name (should fail)
POST {{base_url}}/api/scan/

HTTP 404

### Test 3: Invalid org name
POST {{base_url}}/api/scan/invalid-org-xyz-123

HTTP 404
[Asserts]
jsonpath "$.error" exists

GET {{base_url}}/api/health

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"

### Test 8: Test graph endpoint with boolean parameters
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true

HTTP 200
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=false

HTTP 200
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

### Test 9: Test scan with zero max_repos (boundary case)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=0&max_teams=10

HTTP 400
[Asserts]
jsonpath "$.error" exists

### Test 10: Test scan with zero max_teams (boundary case)
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=10&max_teams=0

HTTP 201
[Asserts]
jsonpath "$.data.success" == true