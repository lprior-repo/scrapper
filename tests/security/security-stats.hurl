# Security Tests for Stats Endpoint
# Tests for /api/stats/{org} endpoint to verify security hardening and attack resistance

# Test 1: Stats endpoint with XSS attempt in org parameter
GET {{base_url}}/api/stats/{{xss_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain XSS payload
body not contains "<script>"

# Test 2: Stats endpoint with SQL injection attempt in org parameter
GET {{base_url}}/api/stats/{{sql_injection_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain SQL injection payload
body not contains "DROP TABLE"

# Test 3: Stats endpoint with path traversal attempt in org parameter
GET {{base_url}}/api/stats/{{path_traversal_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain path traversal payload
body not contains "../../../"

# Test 4: Stats endpoint with oversized org parameter
GET {{base_url}}/api/stats/${"A" * 1000}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 5: Stats endpoint with null bytes in org parameter
GET {{base_url}}/api/stats/test\x00org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 6: Stats endpoint with Unicode attack in org parameter
GET {{base_url}}/api/stats/test\u0000\u0001\u0002
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 7: Stats endpoint with CRLF injection in headers
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: test\r\nX-Injected-Header: malicious
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists

# Test 8: Stats endpoint with command injection in org parameter
GET {{base_url}}/api/stats/$(whoami)
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 9: Stats endpoint with LDAP injection in org parameter
GET {{base_url}}/api/stats/*)(uid=*))(|(uid=*
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 10: Stats endpoint with NoSQL injection in org parameter
GET {{base_url}}/api/stats/{"$ne": null}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 11: Stats endpoint with malicious User-Agent
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{xss_pattern}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists
# Response should not contain XSS payload
body not contains "<script>"

# Test 12: Stats endpoint with special characters in org parameter
GET {{base_url}}/api/stats/test!@#$%^&*()
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 13: Stats endpoint with multiple slashes in org parameter
GET {{base_url}}/api/stats/test//org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 14: Stats endpoint with encoded characters in org parameter
GET {{base_url}}/api/stats/test%20org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 15: Stats endpoint with malformed Accept header
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: application/json; charset=utf-8; malformed="

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists

# Test 16: Stats endpoint with oversized headers
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Large-Header: ${"A" * 10000}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists

# Test 17: Stats endpoint with HTTP method override attempt
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-HTTP-Method-Override: DELETE

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists

# Test 18: Stats endpoint with malicious Accept header
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{xss_pattern}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists

# Test 19: Stats endpoint with multiple malicious headers
GET {{base_url}}/api/stats/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-XSS-Test: {{xss_pattern}}
X-SQL-Test: {{sql_injection_pattern}}
X-Path-Test: {{path_traversal_pattern}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists

# Test 20: Stats endpoint with case variations in org parameter
GET {{base_url}}/api/stats/TEST-ORG
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists