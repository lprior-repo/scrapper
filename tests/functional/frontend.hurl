# Frontend Health Check Test
# Tests the React application at localhost:3000

# Test 1: Basic health check - frontend should be accessible
GET http://localhost:3000/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "<!doctype html>"

# Test 2: Check if JavaScript bundle is loading
GET http://localhost:3000/
HTTP 200
[Asserts]
# Should contain React app root element
body contains "id=\"root\""

# Test 3: Check if application is not serving error pages
GET http://localhost:3000/
HTTP 200
[Asserts]
# Should not contain common error indicators
body not contains "Error"
body not contains "500"
body not contains "404"

# Test 4: Verify the app responds quickly (performance check)
GET http://localhost:3000/
HTTP 200
[Asserts]
duration < 1000  # Should load within 1 second