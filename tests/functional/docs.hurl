# Docs Endpoint Test Suite
# Tests for /api/docs endpoint to verify OpenAPI documentation UI

# Test 1: Basic docs page retrieval
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"
body contains "<title>GitHub Codeowners API Documentation</title>"
body contains "swagger-ui"
duration < {{response_time_threshold}}

# Test 2: HTML structure validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<html>"
body contains "<head>"
body contains "<body>"
body contains "</html>"
body contains "</head>"
body contains "</body>"

# Test 3: Swagger UI resources validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "swagger-ui-dist"
body contains "swagger-ui-bundle.js"
body contains "swagger-ui-standalone-preset.js"
body contains "swagger-ui.css"

# Test 4: OpenAPI spec reference validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "/api/openapi.yaml"
body contains "url: '/api/openapi.yaml'"

# Test 5: Swagger UI configuration validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "SwaggerUIBundle"
body contains "dom_id: '#swagger-ui'"
body contains "deepLinking: true"
body contains "SwaggerUIStandalonePreset"

# Test 6: CSS styling validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<style>"
body contains "box-sizing: border-box"
body contains "overflow-y: scroll"
body contains "background: #fafafa"

# Test 7: Docs with different Accept headers
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"

# Test 8: Docs with minimal headers
GET {{base_url}}/api/docs

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"

# Test 9: Docs response timing
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
[Asserts]
duration < 2000  # Should respond within 2 seconds

# Test 10: Docs page title validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<title>GitHub Codeowners API Documentation</title>"

# Test 11: Swagger UI container validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<div id=\"swagger-ui\"></div>"

# Test 12: JavaScript initialization validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "window.onload = function()"
body contains "const ui = SwaggerUIBundle"

# Test 13: Swagger UI plugins validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "plugins: ["
body contains "SwaggerUIBundle.plugins.DownloadUrl"
body contains "layout: \"StandaloneLayout\""

# Test 14: Swagger UI presets validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "presets: ["
body contains "SwaggerUIBundle.presets.apis"
body contains "SwaggerUIStandalonePreset"

# Test 15: HTML document structure
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"
body matches "(?s).*<html>.*<head>.*<body>.*</body>.*</html>.*"

# Test 16: External CDN resources
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "https://unpkg.com/swagger-ui-dist@5.11.0"
body contains "swagger-ui.css"
body contains "swagger-ui-bundle.js"
body contains "swagger-ui-standalone-preset.js"

# Test 17: Docs with custom User-Agent
GET {{base_url}}/api/docs
User-Agent: Custom-Browser/1.0
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"

# Test 18: Swagger UI version validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "swagger-ui-dist@5.11.0"

# Test 19: Response size validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
[Asserts]
body contains "<!DOCTYPE html>"
bytes count > 1000  # Should be substantial HTML
bytes count < 10000  # But not too large

# Test 20: HTML meta tags validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<head>"
body contains "<title>"

# Test 21: Swagger UI deep linking
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "deepLinking: true"

# Test 22: Swagger UI DOM configuration
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "dom_id: '#swagger-ui'"

# Test 23: JavaScript error handling
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "window.onload"
body not contains "console.error"
body not contains "throw new Error"

# Test 24: CSS box model configuration
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "box-sizing: border-box"
body contains "box-sizing: inherit"

# Test 25: Complete HTML document validation
GET {{base_url}}/api/docs
User-Agent: {{user_agent}}
Accept: {{accept_html}}

HTTP 200
Content-Type: text/html
[Asserts]
body contains "<!DOCTYPE html>"
body contains "<html>"
body contains "<head>"
body contains "<title>GitHub Codeowners API Documentation</title>"
body contains "<body>"
body contains "<div id=\"swagger-ui\"></div>"
body contains "<script>"
body contains "window.onload"
body contains "</script>"
body contains "</body>"
body contains "</html>"