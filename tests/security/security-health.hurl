# Security Tests for Health Endpoint
# Tests for /api/health endpoint to verify security hardening and attack resistance

# Test 1: Health endpoint with XSS attempt
GET {{base_url}}/api/health
User-Agent: {{xss_pattern}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
# Response should not contain XSS payload
body not contains "<script>"

# Test 2: Health endpoint with SQL injection attempt in headers
GET {{base_url}}/api/health
User-Agent: {{sql_injection_pattern}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
# Response should not contain SQL injection payload
body not contains "DROP TABLE"

# Test 3: Health endpoint with path traversal attempt
GET {{base_url}}/api/health
User-Agent: {{path_traversal_pattern}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
# Response should not contain path traversal payload
body not contains "../../../"

# Test 4: Health endpoint with malicious Accept header
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{xss_pattern}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 5: Health endpoint with oversized headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-Large-Header: ${"A" * 10000}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 6: Health endpoint with null bytes
GET {{base_url}}/api/health
User-Agent: test\x00user
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 7: Health endpoint with Unicode attack
GET {{base_url}}/api/health
User-Agent: test\u0000\u0001\u0002
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 8: Health endpoint with CRLF injection
GET {{base_url}}/api/health
User-Agent: test\r\nX-Injected-Header: malicious
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 9: Health endpoint with multiple malicious headers
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-XSS-Test: {{xss_pattern}}
X-SQL-Test: {{sql_injection_pattern}}
X-Path-Test: {{path_traversal_pattern}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 10: Health endpoint with HTTP method override attempt
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}
X-HTTP-Method-Override: DELETE

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 11: Health endpoint security headers check
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"
# Check for security headers if implemented
header "X-Content-Type-Options" exists
header "X-Frame-Options" exists
header "X-XSS-Protection" exists

# Test 12: Health endpoint with malformed JSON in Accept
GET {{base_url}}/api/health
User-Agent: {{user_agent}}
Accept: application/json; charset=utf-8; malformed="

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 13: Health endpoint with command injection attempt
GET {{base_url}}/api/health
User-Agent: $(whoami)
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 14: Health endpoint with LDAP injection attempt
GET {{base_url}}/api/health
User-Agent: *)(uid=*))(|(uid=*
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"

# Test 15: Health endpoint with NoSQL injection attempt
GET {{base_url}}/api/health
User-Agent: {"$ne": null}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.data.status" == "healthy"