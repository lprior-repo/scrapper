# OpenAPI Spec Endpoint Test Suite
# Tests for /api/openapi.yaml endpoint to verify OpenAPI specification file

# Test 1: Basic OpenAPI spec retrieval
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"
body contains "info:"
body contains "title: GitHub Codeowners Visualization API"
body contains "version: 1.0.0"
duration < {{response_time_threshold}}

# Test 2: OpenAPI version validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"
body matches "(?m)^openapi: 3\\.0\\.3$"

# Test 3: API info section validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "info:"
body contains "title: GitHub Codeowners Visualization API"
body contains "description: API for scanning GitHub organizations"
body contains "version: 1.0.0"
body contains "contact:"
body contains "license:"

# Test 4: Server configuration validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "servers:"
body contains "url: http://localhost:8080"
body contains "description: Development server"

# Test 5: Paths section validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "paths:"
body contains "/api/health:"
body contains "/api/scan/{org}:"

# Test 6: Health endpoint specification
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "/api/health:"
body contains "get:"
body contains "summary: Health check endpoint"
body contains "operationId: healthCheck"
body contains "tags:"
body contains "- System"

# Test 7: Scan endpoint specification
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "/api/scan/{org}:"
body contains "post:"
body contains "summary: Scan GitHub organization"
body contains "operationId: scanOrganization"
body contains "tags:"
body contains "- Scanning"

# Test 8: Parameters validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "parameters:"
body contains "name: org"
body contains "in: path"
body contains "required: true"
body contains "name: max_repos"
body contains "name: max_teams"
body contains "name: use_topics"

# Test 9: Response schemas validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "responses:"
body contains "'200':"
body contains "description:"
body contains "content:"
body contains "application/json:"
body contains "schema:"

# Test 10: Components section validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "components:"
body contains "schemas:"
body contains "ScanResponse:"

# Test 11: Schema definitions validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "ScanResponse:"
body contains "type: object"
body contains "properties:"
body contains "success:"
body contains "organization:"
body contains "summary:"

# Test 12: Tags validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "tags:"
body contains "- name: System"
body contains "description: System health and status operations"
body contains "- name: Scanning"
body contains "description: GitHub organization scanning operations"

# Test 13: Contact information validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "contact:"
body contains "name: API Support"
body contains "url: https://github.com/your-org/scrapper"

# Test 14: License information validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "license:"
body contains "name: MIT"
body contains "url: https://opensource.org/licenses/MIT"

# Test 15: Response timing validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
[Asserts]
duration < 1000  # Should respond quickly

# Test 16: YAML format validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi:"
body contains "info:"
body contains "paths:"
body matches "(?m)^\\s*[a-zA-Z_][a-zA-Z0-9_]*:\\s*"

# Test 17: Property data types validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "type: string"
body contains "type: integer"
body contains "type: boolean"
body contains "type: object"
body contains "type: array"

# Test 18: Parameter constraints validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "minimum: 1"
body contains "maximum: 1000"
body contains "default: 100"
body contains "default: 50"
body contains "default: false"

# Test 19: Example values validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "example:"
body contains "example: \"microsoft\""
body contains "example: \"healthy\""
body contains "example: \"connected\""
body contains "example: \"1.0.0\""

# Test 20: OpenAPI with different Accept headers
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: application/yaml, text/yaml, text/plain, */*

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"

# Test 21: OpenAPI with minimal headers
GET {{base_url}}/api/openapi.yaml

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"

# Test 22: Response size validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
[Asserts]
body contains "openapi: 3.0.3"
bytes count > 1000  # Should be substantial
bytes count < 20000  # But not too large

# Test 23: YAML indentation validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"
body contains "info:"
body contains "  title:"
body contains "  version:"

# Test 24: Required fields validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "required: true"
body contains "required: false"

# Test 25: Complete OpenAPI structure
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"
body contains "info:"
body contains "servers:"
body contains "paths:"
body contains "components:"
body contains "tags:"

# Test 26: Custom User-Agent handling
GET {{base_url}}/api/openapi.yaml
User-Agent: OpenAPI-Client/1.0
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "openapi: 3.0.3"

# Test 27: API title and description
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "title: GitHub Codeowners Visualization API"
body contains "description: API for scanning GitHub organizations and visualizing codeowners relationships"

# Test 28: Operation IDs validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "operationId: healthCheck"
body contains "operationId: scanOrganization"

# Test 29: Date-time format validation
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "format: date-time"

# Test 30: Response content types
GET {{base_url}}/api/openapi.yaml
User-Agent: {{user_agent}}
Accept: {{accept_yaml}}

HTTP 200
Content-Type: application/yaml
[Asserts]
body contains "application/json:"
body contains "schema:"