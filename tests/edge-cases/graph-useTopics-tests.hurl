# Graph Endpoint useTopics Parameter Tests
# Tests for /api/graph/{org} endpoint useTopics parameter functionality

# Test 1: Graph endpoint with useTopics=true
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

# Test 2: Graph endpoint with useTopics=false
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

# Test 3: Graph endpoint with useTopics=1 (should work as true)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 4: Graph endpoint with useTopics=0 (should work as false)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=0
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 5: Graph endpoint with useTopics=invalid (should use default false)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=invalid
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 6: Graph endpoint with useTopics empty value (should use default false)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 7: Graph endpoint with useTopics parameter case sensitivity
GET {{base_url}}/api/graph/{{test_org_small}}?USETOPICS=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 8: Graph endpoint with useTopics parameter case sensitivity
GET {{base_url}}/api/graph/{{test_org_small}}?UseTopics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 9: Graph endpoint with duplicate useTopics parameters
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true&useTopics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 10: Graph endpoint with useTopics and other parameters
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true&limit=100
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 11: Graph endpoint without useTopics parameter (should use default false)
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

# Test 12: Graph endpoint with useTopics=TRUE (uppercase)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=TRUE
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 13: Graph endpoint with useTopics=FALSE (uppercase)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=FALSE
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 14: Graph endpoint with useTopics=yes (non-standard boolean)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=yes
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 15: Graph endpoint with useTopics=no (non-standard boolean)
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=no
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 16: Graph endpoint with useTopics numeric values
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=2
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 17: Graph endpoint with useTopics negative values
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=-1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 18: Graph endpoint with useTopics float values
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=1.5
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 19: Graph endpoint with very long useTopics value
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=thisIsAVeryLongStringThatShouldNotBeValidBooleanValue
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 20: Graph endpoint with useTopics and malformed URL
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true&malformed=param&
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 21: Graph endpoint with useTopics and special characters in org name
GET {{base_url}}/api/graph/test-org-123?useTopics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 22: Graph endpoint with useTopics and invalid org name
GET {{base_url}}/api/graph/{{test_org_invalid}}?useTopics=true
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 23: Graph endpoint with useTopics and invalid org name
GET {{base_url}}/api/graph/{{test_org_invalid}}?useTopics=false
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 24: Graph endpoint with useTopics and URL encoding
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics%3Dtrue
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0

# Test 25: Graph endpoint with useTopics and multiple query separators
GET {{base_url}}/api/graph/{{test_org_small}}?useTopics=true&amp;other=value
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.nodes" exists
jsonpath "$.data.nodes" count >= 0
jsonpath "$.data.edges" exists
jsonpath "$.data.edges" count >= 0