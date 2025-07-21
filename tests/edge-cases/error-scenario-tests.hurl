# Edge Case Tests - Error Scenario Testing
# Tests API endpoints with various error conditions and malformed inputs

### Test 1: Test with invalid HTTP method on scan endpoint
GET {{base_url}}/api/scan/{{test_org_small}}

HTTP 404

### Test 2: Test with invalid query parameter values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=invalid&max_teams=notanumber

HTTP 201
[Asserts]
# Should use default values when invalid parameters provided
jsonpath "$.data.success" == true

### Test 3: Test with negative values for max parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=-1&max_teams=-5

HTTP 400
[Asserts]
# Should reject negative values
jsonpath "$.error" exists

### Test 4: Test with extremely large values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=999999999&max_teams=999999999

HTTP 201
[Asserts]
jsonpath "$.data.success" == true

### Test 5: Test with malformed boolean parameter
POST {{base_url}}/api/scan/{{test_org_small}}?use_topics=notabool

HTTP 201
[Asserts]
# Should handle malformed boolean gracefully
jsonpath "$.data.success" == true

### Test 6: Test graph endpoint with malformed useTopics parameter
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=invalid

HTTP 200
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.edges" exists

### Test 7: Test with organization name containing special characters
POST {{base_url}}/api/scan/test-org-with-dashes

HTTP 404
[Asserts]
jsonpath "$.error" exists

### Test 8: Test with organization name containing URL-encoded characters
POST {{base_url}}/api/scan/test%20org

HTTP 404
[Asserts]
jsonpath "$.error" exists

### Test 9: Test non-existent endpoints
GET {{base_url}}/api/nonexistent

HTTP 404

POST {{base_url}}/api/invalid/endpoint

HTTP 404

### Test 10: Test with missing required path parameters
GET {{base_url}}/api/graph/

HTTP 404

GET {{base_url}}/api/stats/

HTTP 404

### Test 11: Test health endpoint with invalid method
POST {{base_url}}/api/health

HTTP 404

DELETE {{base_url}}/api/health

HTTP 404

### Test 12: Test endpoints with extra path segments
GET {{base_url}}/api/health/extra/path

HTTP 404

GET {{base_url}}/api/graph/{{test_org_small}}/extra

HTTP 404

### Test 13: Test with organization name containing only numbers
POST {{base_url}}/api/scan/12345

HTTP 404
[Asserts]
jsonpath "$.error" exists

### Test 14: Test concurrent requests to same endpoint
GET {{base_url}}/api/health

HTTP 200

GET {{base_url}}/api/health

HTTP 200

### Test 15: Test with very long query parameter values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1000000000000000000000000000000000000000000000000000

HTTP 201
[Asserts]
jsonpath "$.data.success" == true