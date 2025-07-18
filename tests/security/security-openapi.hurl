# Security Tests for OpenAPI Endpoint
# Tests for /api/openapi.yaml endpoint to verify security hardening and attack resistance

# Test 1: OpenAPI endpoint with XSS attempt in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {{xss_pattern}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
# Response should not contain XSS payload
body not contains "<script>alert('xss')</script>"

# Test 2: OpenAPI endpoint with SQL injection attempt in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {{sql_injection_pattern}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
# Response should not contain SQL injection payload
body not contains "DROP TABLE"

# Test 3: OpenAPI endpoint with path traversal attempt in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {{path_traversal_pattern}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
# Response should not contain path traversal payload
body not contains "../../../"

# Test 4: OpenAPI endpoint with oversized headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
X-Large-Header: ${"A" * 10000}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 5: OpenAPI endpoint with null bytes in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: test\x00user
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 6: OpenAPI endpoint with Unicode attack in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: test\u0000\u0001\u0002
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 7: OpenAPI endpoint with CRLF injection in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: test\r\nX-Injected-Header: malicious
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 8: OpenAPI endpoint with command injection in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: $(whoami)
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 9: OpenAPI endpoint with LDAP injection in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: *)(uid=*))(|(uid=*
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 10: OpenAPI endpoint with NoSQL injection in headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {"$ne": null}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 11: OpenAPI endpoint with malicious Accept header
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{xss_pattern}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 12: OpenAPI endpoint with HTTP method override attempt
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
X-HTTP-Method-Override: DELETE

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 13: OpenAPI endpoint with multiple malicious headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
X-XSS-Test: {{xss_pattern}}
X-SQL-Test: {{sql_injection_pattern}}
X-Path-Test: {{path_traversal_pattern}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 14: OpenAPI endpoint with malformed Accept header
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: application/yaml; charset=utf-8; malformed="

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 15: OpenAPI endpoint with case variations in headers
GET {{base_url}}/api/openapi.yaml
user-agent: {{user_agent}}
accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 16: OpenAPI endpoint security headers check
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"
# Check for security headers if implemented
header "X-Content-Type-Options" exists
header "X-Frame-Options" exists
header "X-XSS-Protection" exists

# Test 17: OpenAPI endpoint with malicious Referer header
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
Referer: {{xss_pattern}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 18: OpenAPI endpoint with malicious Origin header
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
Origin: {{xss_pattern}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 19: OpenAPI endpoint with cookie injection attempt
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
Cookie: malicious={{xss_pattern}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"

# Test 20: OpenAPI endpoint with authorization bypass attempt
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}
Authorization: Bearer {{xss_pattern}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"