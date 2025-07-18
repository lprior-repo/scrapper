# Common assertions that can be reused across test suites

# Standard HTTP status code assertions
[template "assert_200"]
HTTP 200

[template "assert_201"]
HTTP 201

[template "assert_400"]
HTTP 400

[template "assert_401"]
HTTP 401

[template "assert_403"]
HTTP 403

[template "assert_404"]
HTTP 404

[template "assert_422"]
HTTP 422

[template "assert_500"]
HTTP 500

[template "assert_503"]
HTTP 503

# Content type assertions
[template "assert_json_content"]
Content-Type: application/json

[template "assert_html_content"]
Content-Type: text/html

[template "assert_yaml_content"]
Content-Type: application/yaml

[template "assert_text_content"]
Content-Type: text/plain

# Response time assertions
[template "assert_fast_response"]
duration < 1000

[template "assert_normal_response"]
duration < 5000

[template "assert_slow_response"]
duration < 30000

# JSON structure assertions
[template "assert_success_response"]
jsonpath "$.success" == true

[template "assert_error_response"]
jsonpath "$.success" == false

[template "assert_has_data"]
jsonpath "$.data" exists

[template "assert_has_errors"]
jsonpath "$.errors" exists

[template "assert_organization_field"]
jsonpath "$.organization" exists
jsonpath "$.organization" isString

[template "assert_summary_field"]
jsonpath "$.summary" exists
jsonpath "$.summary" isObject

# Health check assertions
[template "assert_healthy"]
jsonpath "$.status" == "healthy"
jsonpath "$.database" == "connected"
jsonpath "$.version" exists
jsonpath "$.timestamp" exists

# Graph response assertions
[template "assert_graph_structure"]
jsonpath "$.nodes" exists
jsonpath "$.nodes" isArray
jsonpath "$.edges" exists
jsonpath "$.edges" isArray

# Stats response assertions
[template "assert_stats_structure"]
jsonpath "$.organization" exists
jsonpath "$.total_repositories" exists
jsonpath "$.total_teams" exists
jsonpath "$.total_topics" exists
jsonpath "$.total_users" exists
jsonpath "$.total_codeowners" exists
jsonpath "$.codeowner_coverage" exists
jsonpath "$.last_scan_time" exists

# Scan response assertions
[template "assert_scan_structure"]
jsonpath "$.success" exists
jsonpath "$.organization" exists
jsonpath "$.summary" exists
jsonpath "$.summary.total_repos" exists
jsonpath "$.summary.repos_with_codeowners" exists
jsonpath "$.summary.total_teams" exists
jsonpath "$.summary.total_topics" exists
jsonpath "$.summary.unique_owners" exists
jsonpath "$.summary.api_calls_used" exists
jsonpath "$.summary.processing_time_ms" exists
jsonpath "$.errors" exists
jsonpath "$.data" exists

# Security assertions
[template "assert_no_sql_injection"]
body not contains "SQL"
body not contains "database"
body not contains "table"

[template "assert_no_xss"]
body not contains "<script>"
body not contains "javascript:"
body not contains "alert("

[template "assert_no_path_traversal"]
body not contains "/etc/passwd"
body not contains "../"

# Rate limiting assertions
[template "assert_rate_limit_headers"]
header "X-RateLimit-Limit" exists
header "X-RateLimit-Remaining" exists
header "X-RateLimit-Reset" exists

# CORS assertions
[template "assert_cors_headers"]
header "Access-Control-Allow-Origin" exists
header "Access-Control-Allow-Methods" exists
header "Access-Control-Allow-Headers" exists

# Cache assertions
[template "assert_no_cache"]
header "Cache-Control" contains "no-cache"

[template "assert_cacheable"]
header "Cache-Control" contains "public"

# Compression assertions
[template "assert_compression"]
header "Content-Encoding" == "gzip"

# OpenAPI spec assertions
[template "assert_openapi_spec"]
jsonpath "$.openapi" == "3.0.3"
jsonpath "$.info.title" exists
jsonpath "$.info.version" exists
jsonpath "$.paths" exists
jsonpath "$.components" exists

# Error format assertions
[template "assert_error_format"]
jsonpath "$.error" exists
jsonpath "$.error.message" exists
jsonpath "$.error.code" exists

# Pagination assertions
[template "assert_pagination"]
jsonpath "$.pagination" exists
jsonpath "$.pagination.page" exists
jsonpath "$.pagination.per_page" exists
jsonpath "$.pagination.total" exists

# Common validation patterns
[template "assert_uuid_format"]
matches "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

[template "assert_timestamp_format"]
matches "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d{3})?Z?$"

[template "assert_github_org_format"]
matches "^[a-zA-Z0-9]([a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$"

[template "assert_positive_integer"]
isInteger
greaterThan 0

[template "assert_non_negative_integer"]
isInteger
greaterThanOrEqualTo 0