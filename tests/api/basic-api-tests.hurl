# Basic API Tests
# This file contains fundamental API endpoint tests

# Test 1: Health Check
GET {{base_url}}/health
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
status == 200
header "Content-Type" contains "application/json"
jsonpath "$.status" == "healthy"
jsonpath "$.timestamp" isString
jsonpath "$.version" isString

# Test 2: OpenAPI Documentation
GET {{base_url}}/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
[Asserts]
status == 200
header "Content-Type" contains "application/yaml"
body contains "openapi:"
body contains "paths:"
body contains "components:"

# Test 3: API Documentation
GET {{base_url}}/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
[Asserts]
status == 200
header "Content-Type" contains "text/html"
body contains "<!DOCTYPE html>"
body contains "API Documentation"

# Test 4: Stats Endpoint
GET {{base_url}}/stats
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
status == 200
header "Content-Type" contains "application/json"
jsonpath "$.organizations" isNumber
jsonpath "$.repositories" isNumber
jsonpath "$.teams" isNumber
jsonpath "$.members" isNumber

# Test 5: Invalid Endpoint
GET {{base_url}}/invalid-endpoint
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
[Asserts]
status == 404
header "Content-Type" contains "application/json"
jsonpath "$.error" isString
jsonpath "$.message" contains "not found"

# Test 6: Scan with Valid Organization
POST {{base_url}}/scan
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Content-Type: {{json_content_type}}

{
  "organization": "{{test_org_small}}"
}

HTTP 200
[Asserts]
status == 200
header "Content-Type" contains "application/json"
jsonpath "$.organization" == "{{test_org_small}}"
jsonpath "$.repositories" isArray
jsonpath "$.teams" isArray
jsonpath "$.scan_timestamp" isString

# Test 7: Scan with Invalid Organization
POST {{base_url}}/scan
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Content-Type: {{json_content_type}}

{
  "organization": "{{test_org_invalid}}"
}

HTTP 404
[Asserts]
status == 404
header "Content-Type" contains "application/json"
jsonpath "$.error" isString
jsonpath "$.message" contains "not found"

# Test 8: Graph Data for Organization
GET {{base_url}}/graph/{{test_org_small}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 200
[Asserts]
status == 200
header "Content-Type" contains "application/json"
jsonpath "$.nodes" isArray
jsonpath "$.edges" isArray
jsonpath "$.metadata" isObject
jsonpath "$.metadata.organization" == "{{test_org_small}}"

# Test 9: Graph Data for Non-existent Organization
GET {{base_url}}/graph/{{test_org_invalid}}
User-Agent: {{user_agent}}
Accept: {{accept_json}}

HTTP 404
[Asserts]
status == 404
header "Content-Type" contains "application/json"
jsonpath "$.error" isString

# Test 10: Malformed JSON Request
POST {{base_url}}/scan
User-Agent: {{user_agent}}
Accept: {{accept_json}}
Content-Type: {{json_content_type}}

{
  "organization": "test"
  "invalid": "json"
}

HTTP 400
[Asserts]
status == 400
header "Content-Type" contains "application/json"
jsonpath "$.error" isString
jsonpath "$.message" contains "invalid"