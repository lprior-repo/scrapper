#!/bin/bash

# Test Environment Setup Script
# This script prepares the test environment for running Hurl tests

set -e

echo "Setting up test environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if hurl is installed
check_hurl() {
    print_message $YELLOW "Checking Hurl installation..."
    if command -v hurl &> /dev/null; then
        local version=$(hurl --version | head -n1)
        print_message $GREEN "✓ Hurl is installed: $version"
        return 0
    else
        print_message $RED "✗ Hurl is not installed"
        return 1
    fi
}

# Wait for API server to be ready
wait_for_api() {
    local url=${1:-"http://localhost:8080"}
    local timeout=${2:-30}
    
    print_message $YELLOW "Waiting for API server at $url..."
    
    for i in $(seq 1 $timeout); do
        if curl -s --fail "$url/health" > /dev/null 2>&1; then
            print_message $GREEN "✓ API server is ready"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    
    print_message $RED "✗ API server not ready after ${timeout}s"
    return 1
}

# Wait for Neo4j to be ready
wait_for_neo4j() {
    local uri=${1:-"bolt://localhost:7687"}
    local timeout=${2:-30}
    
    print_message $YELLOW "Waiting for Neo4j at $uri..."
    
    for i in $(seq 1 $timeout); do
        if nc -z localhost 7687 > /dev/null 2>&1; then
            print_message $GREEN "✓ Neo4j is ready"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    
    print_message $RED "✗ Neo4j not ready after ${timeout}s"
    return 1
}

# Create test reports directory
setup_reports() {
    print_message $YELLOW "Setting up test reports directory..."
    mkdir -p tests/reports
    chmod 755 tests/reports
    print_message $GREEN "✓ Reports directory ready"
}

# Export environment variables for tests
setup_env_vars() {
    print_message $YELLOW "Setting up environment variables..."
    
    # Default values
    export HURL_TEST_ENV=${HURL_TEST_ENV:-development}
    export API_BASE_URL=${API_BASE_URL:-http://localhost:8080}
    export NEO4J_URI=${NEO4J_URI:-bolt://localhost:7687}
    export GITHUB_TOKEN=${GITHUB_TOKEN:-}
    export DEBUG_MODE=${DEBUG_MODE:-true}
    
    print_message $GREEN "✓ Environment variables set"
    print_message $YELLOW "  HURL_TEST_ENV: $HURL_TEST_ENV"
    print_message $YELLOW "  API_BASE_URL: $API_BASE_URL"
    print_message $YELLOW "  NEO4J_URI: $NEO4J_URI"
    print_message $YELLOW "  DEBUG_MODE: $DEBUG_MODE"
}

# Main setup function
main() {
    print_message $GREEN "Starting test environment setup..."
    
    # Check dependencies
    if ! check_hurl; then
        print_message $RED "Please install Hurl first"
        exit 1
    fi
    
    # Setup environment
    setup_env_vars
    setup_reports
    
    # Wait for services (optional, controlled by flags)
    if [[ "${1:-}" == "--wait-services" ]]; then
        wait_for_api "$API_BASE_URL"
        wait_for_neo4j "$NEO4J_URI"
    fi
    
    print_message $GREEN "✅ Test environment setup complete!"
    print_message $YELLOW "Ready to run Hurl tests"
}

# Run main function
main "$@"