# Edge Case Tests - Boundary Value Testing
# Tests API endpoints with boundary values and edge cases

### Test 1: Scan organization with minimum valid parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_teams=1&use_topics=false

HTTP 201
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.organization" == "{{test_org_small}}"
jsonpath "$.data.summary.total_repos" >= 0
jsonpath "$.data.summary.total_teams" >= 0

### Test 2: Scan organization with maximum valid parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1000&max_teams=500&use_topics=true

HTTP 201
[Asserts]
jsonpath "$.data.success" == true
jsonpath "$.data.organization" == "{{test_org_small}}"

### Test 3: Test with empty organization name (should fail)
POST {{base_url}}/api/scan/

HTTP 404

### Test 4: Test with very long organization name
POST {{base_url}}/api/scan/thisIsAVeryLongOrganizationNameThatMightCauseIssuesInSomeSystemsAndShouldBeTestedForProperHandling

HTTP 404
[Asserts]
jsonpath "$.error" exists

### Test 5: Test graph endpoint with invalid org
GET {{base_url}}/api/graph/{{test_org_invalid}}

HTTP 200
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

### Test 6: Test stats endpoint with special characters in org name
GET {{base_url}}/api/stats/{{test_org_invalid}}

HTTP 404
[Asserts]
jsonpath "$.error" exists

### Test 7: Test health endpoint multiple times (should be consistent)
GET {{base_url}}/api/health

HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"

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