# Kubernetes-style Health Check Tests
# These tests simulate the checks that Kubernetes would perform

# Readiness Probe - Check if API is ready to serve traffic
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"
jsonpath "$.data.database" == "connected"
duration < 5000

# Liveness Probe - Check if API is alive
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
duration < 3000
header "Content-Type" contains "application/json"

# Startup Probe - Check if API has started successfully
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
jsonpath "$.data.version" exists
jsonpath "$.data.timestamp" exists

# Service Discovery Test - Verify OpenAPI spec is accessible
GET http://localhost:8081/api/openapi.yaml
HTTP 200
[Asserts]
duration < 2000
body contains "paths:"

# Network Connectivity Test - Multiple endpoints
GET http://localhost:8081/api/docs
HTTP 200

GET http://localhost:8081/api/health
HTTP 200

# Resource Availability Test - Check database connectivity
GET http://localhost:8081/api/health
HTTP 200
[Asserts]
jsonpath "$.data.database" == "connected"

# Load Balancer Health Check Simulation
GET http://localhost:8081/api/health
User-Agent: kube-probe/1.0
HTTP 200
[Asserts]
jsonpath "$.data.status" == "healthy"