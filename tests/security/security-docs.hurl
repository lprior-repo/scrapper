# Security Tests for Docs Endpoint
# Tests for /api/docs endpoint to verify security hardening and attack resistance

# Test 1: Docs endpoint with XSS attempt in headers
GET {{base_url}}/api/docs
User-Agent: {{xss_pattern}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
# Response should not contain XSS payload
body not contains "<script>alert('xss')</script>"

# Test 2: Docs endpoint with SQL injection attempt in headers
GET {{base_url}}/api/docs
User-Agent: {{sql_injection_pattern}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
# Response should not contain SQL injection payload
body not contains "DROP TABLE"

# Test 3: Docs endpoint with path traversal attempt in headers
GET {{base_url}}/api/docs
User-Agent: {{path_traversal_pattern}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
# Response should not contain path traversal payload
body not contains "../../../"

# Test 4: Docs endpoint with oversized headers
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
X-Large-Header: ${"A" * 10000}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 5: Docs endpoint with null bytes in headers
GET {{base_url}}/api/docs
User-Agent: test\x00user
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 6: Docs endpoint with Unicode attack in headers
GET {{base_url}}/api/docs
User-Agent: test\u0000\u0001\u0002
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 7: Docs endpoint with CRLF injection in headers
GET {{base_url}}/api/docs
User-Agent: test\r\nX-Injected-Header: malicious
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 8: Docs endpoint with command injection in headers
GET {{base_url}}/api/docs
User-Agent: $(whoami)
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 9: Docs endpoint with LDAP injection in headers
GET {{base_url}}/api/docs
User-Agent: *)(uid=*))(|(uid=*
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 10: Docs endpoint with NoSQL injection in headers
GET {{base_url}}/api/docs
User-Agent: {"$ne": null}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 11: Docs endpoint with malicious Accept header
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{xss_pattern}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 12: Docs endpoint with HTTP method override attempt
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
X-HTTP-Method-Override: DELETE

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 13: Docs endpoint with multiple malicious headers
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
X-XSS-Test: {{xss_pattern}}
X-SQL-Test: {{sql_injection_pattern}}
X-Path-Test: {{path_traversal_pattern}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 14: Docs endpoint with malformed Accept header
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: text/html; charset=utf-8; malformed="

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 15: Docs endpoint with case variations in headers
GET {{base_url}}/api/docs
user-agent: {{user_agent}}
accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 16: Docs endpoint security headers check
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"
# Check for security headers if implemented
header "X-Content-Type-Options" exists
header "X-Frame-Options" exists
header "X-XSS-Protection" exists

# Test 17: Docs endpoint with malicious Referer header
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
Referer: {{xss_pattern}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 18: Docs endpoint with malicious Origin header
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
Origin: {{xss_pattern}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 19: Docs endpoint with cookie injection attempt
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
Cookie: malicious={{xss_pattern}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"

# Test 20: Docs endpoint with authorization bypass attempt
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}
Authorization: Bearer {{xss_pattern}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "API Documentation"