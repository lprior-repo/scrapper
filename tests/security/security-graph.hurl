# Security Tests for Graph Endpoint
# Tests for /api/graph/{org} endpoint to verify security hardening and attack resistance

# Test 1: Graph endpoint with XSS attempt in org parameter
GET {{base_url}}/api/graph/{{xss_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain XSS payload
body not contains "<script>"

# Test 2: Graph endpoint with SQL injection attempt in org parameter
GET {{base_url}}/api/graph/{{sql_injection_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain SQL injection payload
body not contains "DROP TABLE"

# Test 3: Graph endpoint with path traversal attempt in org parameter
GET {{base_url}}/api/graph/{{path_traversal_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain path traversal payload
body not contains "../../../"

# Test 4: Graph endpoint with oversized org parameter
GET {{base_url}}/api/graph/${"A" * 1000}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 5: Graph endpoint with null bytes in org parameter
GET {{base_url}}/api/graph/test\x00org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 6: Graph endpoint with Unicode attack in org parameter
GET {{base_url}}/api/graph/test\u0000\u0001\u0002
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 7: Graph endpoint with CRLF injection in headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: test\r\nX-Injected-Header: malicious
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 8: Graph endpoint with command injection in org parameter
GET {{base_url}}/api/graph/$(whoami)
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 9: Graph endpoint with LDAP injection in org parameter
GET {{base_url}}/api/graph/*)(uid=*))(|(uid=*
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 10: Graph endpoint with NoSQL injection in org parameter
GET {{base_url}}/api/graph/{"$ne": null}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 11: Graph endpoint with malicious User-Agent
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{xss_pattern}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists
# Response should not contain XSS payload
body not contains "<script>"

# Test 12: Graph endpoint with special characters in org parameter
GET {{base_url}}/api/graph/test!@#$%^&*()
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 13: Graph endpoint with multiple slashes in org parameter
GET {{base_url}}/api/graph/test//org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 14: Graph endpoint with encoded characters in org parameter
GET {{base_url}}/api/graph/test%20org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 15: Graph endpoint with malformed Accept header
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: application/json; charset=utf-8; malformed="

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 16: Graph endpoint with oversized headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Large-Header: ${"A" * 10000}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 17: Graph endpoint with HTTP method override attempt
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-HTTP-Method-Override: DELETE

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 18: Graph endpoint with malicious Accept header
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{xss_pattern}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 19: Graph endpoint with multiple malicious headers
GET {{base_url}}/api/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-XSS-Test: {{xss_pattern}}
X-SQL-Test: {{sql_injection_pattern}}
X-Path-Test: {{path_traversal_pattern}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists

# Test 20: Graph endpoint with case variations in org parameter
GET {{base_url}}/api/graph/TEST-ORG
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.nodes" exists
jsonpath "$.edges" exists