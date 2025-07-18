#!/bin/bash

# Hurl Test Runner Script
# This script runs different types of tests with proper environment setup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored messages
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [TEST_TYPE]"
    echo ""
    echo "Test Types:"
    echo "  all         - Run all tests"
    echo "  functional  - Run functional tests"
    echo "  security    - Run security tests"
    echo "  performance - Run performance tests"
    echo "  integration - Run integration tests"
    echo "  api         - Run API tests"
    echo ""
    echo "Options:"
    echo "  -e, --env ENV       Set environment (development, testing, ci, docker)"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -r, --report        Generate HTML report"
    echo "  -h, --help          Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 functional"
    echo "  $0 --env ci --report all"
    echo "  $0 -v security"
}

# Parse command line arguments
parse_args() {
    ENVIRONMENT="development"
    VERBOSE=false
    GENERATE_REPORT=false
    TEST_TYPE=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--env)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -r|--report)
                GENERATE_REPORT=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            all|functional|security|performance|integration|api)
                TEST_TYPE="$1"
                shift
                ;;
            *)
                print_message $RED "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    if [[ -z "$TEST_TYPE" ]]; then
        TEST_TYPE="all"
    fi
}

# Run functional tests
run_functional_tests() {
    print_message $BLUE "Running functional tests..."
    
    local tests=(
        "tests/functional/health.hurl"
        "tests/functional/openapi.hurl"
        "tests/functional/docs.hurl"
        "tests/functional/scan.hurl"
        "tests/functional/graph.hurl"
        "tests/functional/stats.hurl"
    )
    
    for test in "${tests[@]}"; do
        if [[ -f "$test" ]]; then
            print_message $YELLOW "Running $test..."
            run_single_test "$test"
        fi
    done
}

# Run security tests
run_security_tests() {
    print_message $BLUE "Running security tests..."
    
    local tests=(
        "tests/security/security-health.hurl"
        "tests/security/security-openapi.hurl"
        "tests/security/security-docs.hurl"
        "tests/security/security-scan.hurl"
        "tests/security/security-graph.hurl"
        "tests/security/security-stats.hurl"
    )
    
    for test in "${tests[@]}"; do
        if [[ -f "$test" ]]; then
            print_message $YELLOW "Running $test..."
            run_single_test "$test"
        fi
    done
}

# Run performance tests
run_performance_tests() {
    print_message $BLUE "Running performance tests..."
    
    local tests=(
        "tests/performance/performance-health.hurl"
        "tests/performance/performance-scan.hurl"
        "tests/performance/performance-graph.hurl"
    )
    
    for test in "${tests[@]}"; do
        if [[ -f "$test" ]]; then
            print_message $YELLOW "Running $test..."
            run_single_test "$test"
        fi
    done
}

# Run integration tests
run_integration_tests() {
    print_message $BLUE "Running integration tests..."
    
    # Find all .hurl files in integration directory
    local integration_tests=$(find tests/integration -name "*.hurl" 2>/dev/null || true)
    
    if [[ -n "$integration_tests" ]]; then
        while IFS= read -r test; do
            print_message $YELLOW "Running $test..."
            run_single_test "$test"
        done <<< "$integration_tests"
    else
        print_message $YELLOW "No integration tests found"
    fi
}

# Run API tests
run_api_tests() {
    print_message $BLUE "Running API tests..."
    
    # Find all .hurl files in api directory
    local api_tests=$(find tests/api -name "*.hurl" 2>/dev/null || true)
    
    if [[ -n "$api_tests" ]]; then
        while IFS= read -r test; do
            print_message $YELLOW "Running $test..."
            run_single_test "$test"
        done <<< "$api_tests"
    else
        print_message $YELLOW "No API tests found"
    fi
}

# Run a single test with proper configuration
run_single_test() {
    local test_file="$1"
    local hurl_opts=()
    
    # Add environment configuration
    hurl_opts+=("--variables-file" "tests/config/global.hurl")
    hurl_opts+=("--variables-file" "tests/config/environments.hurl")
    
    # Add verbose output if requested
    if [[ "$VERBOSE" == "true" ]]; then
        hurl_opts+=("--verbose")
    fi
    
    # Add report generation if requested
    if [[ "$GENERATE_REPORT" == "true" ]]; then
        local report_file="tests/reports/$(basename "$test_file" .hurl).html"
        hurl_opts+=("--report-html" "$report_file")
    fi
    
    # Set environment
    export HURL_TEST_ENV="$ENVIRONMENT"
    
    # Run the test
    if hurl "${hurl_opts[@]}" "$test_file"; then
        print_message $GREEN "✓ $test_file passed"
    else
        print_message $RED "✗ $test_file failed"
        return 1
    fi
}

# Main function
main() {
    parse_args "$@"
    
    print_message $GREEN "Starting Hurl tests..."
    print_message $YELLOW "Environment: $ENVIRONMENT"
    print_message $YELLOW "Test type: $TEST_TYPE"
    print_message $YELLOW "Verbose: $VERBOSE"
    print_message $YELLOW "Generate report: $GENERATE_REPORT"
    
    # Ensure reports directory exists
    mkdir -p tests/reports
    
    # Run tests based on type
    case "$TEST_TYPE" in
        functional)
            run_functional_tests
            ;;
        security)
            run_security_tests
            ;;
        performance)
            run_performance_tests
            ;;
        integration)
            run_integration_tests
            ;;
        api)
            run_api_tests
            ;;
        all)
            run_functional_tests
            run_security_tests
            run_performance_tests
            run_integration_tests
            run_api_tests
            ;;
        *)
            print_message $RED "Unknown test type: $TEST_TYPE"
            show_usage
            exit 1
            ;;
    esac
    
    print_message $GREEN "✅ Test run complete!"
    
    if [[ "$GENERATE_REPORT" == "true" ]]; then
        print_message $YELLOW "Reports generated in tests/reports/"
    fi
}

# Run main function
main "$@"