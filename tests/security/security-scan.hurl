# Security Tests for Scan Endpoint
# Tests for /api/scan/{org} endpoint to verify security hardening and attack resistance

# Test 1: Scan endpoint with XSS attempt in org parameter
POST {{base_url}}/api/scan/{{xss_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain XSS payload
body not contains "<script>"

# Test 2: Scan endpoint with SQL injection attempt in org parameter
POST {{base_url}}/api/scan/{{sql_injection_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain SQL injection payload
body not contains "DROP TABLE"

# Test 3: Scan endpoint with path traversal attempt in org parameter
POST {{base_url}}/api/scan/{{path_traversal_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists
# Response should not contain path traversal payload
body not contains "../../../"

# Test 4: Scan endpoint with malicious query parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos={{xss_pattern}}&max_teams={{sql_injection_pattern}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 5: Scan endpoint with oversized org parameter
POST {{base_url}}/api/scan/${"A" * 1000}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 6: Scan endpoint with null bytes in org parameter
POST {{base_url}}/api/scan/test\x00org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 7: Scan endpoint with Unicode attack in org parameter
POST {{base_url}}/api/scan/test\u0000\u0001\u0002
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 8: Scan endpoint with CRLF injection in headers
POST {{base_url}}/api/scan/{{test_org_small}}
User-Agent: test\r\nX-Injected-Header: malicious
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" exists

# Test 9: Scan endpoint with command injection in org parameter
POST {{base_url}}/api/scan/$(whoami)
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 10: Scan endpoint with LDAP injection in org parameter
POST {{base_url}}/api/scan/*)(uid=*))(|(uid=*
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 11: Scan endpoint with NoSQL injection in org parameter
POST {{base_url}}/api/scan/{"$ne": null}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 12: Scan endpoint with malicious User-Agent
POST {{base_url}}/api/scan/{{test_org_small}}
User-Agent: {{xss_pattern}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" exists
# Response should not contain XSS payload
body not contains "<script>"

# Test 13: Scan endpoint with negative parameter values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=-1&max_teams=-1
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 14: Scan endpoint with extremely large parameter values
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=999999999&max_teams=999999999
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 15: Scan endpoint with malformed query parameters
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=abc&max_teams=xyz&use_topics=maybe
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 16: Scan endpoint with special characters in org parameter
POST {{base_url}}/api/scan/test!@#$%^&*()
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 17: Scan endpoint with multiple slashes in org parameter
POST {{base_url}}/api/scan/test//org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 18: Scan endpoint with encoded characters in org parameter
POST {{base_url}}/api/scan/test%20org
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 19: Scan endpoint with JSON injection in query params
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos={"$gt": 0}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 400
Content-Type: application/json
[Asserts]
jsonpath "$.error" exists

# Test 20: Scan endpoint with HTTP parameter pollution
POST {{base_url}}/api/scan/{{test_org_small}}?max_repos=1&max_repos=999
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.success" exists