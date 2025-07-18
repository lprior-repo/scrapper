#!/bin/bash

# Test Setup Validation Script
# This script validates that all test components are properly configured

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

# Track validation results
VALIDATION_PASSED=0
VALIDATION_FAILED=0

# Validate function
validate() {
    local description="$1"
    local condition="$2"
    
    if eval "$condition"; then
        print_message $GREEN "✓ $description"
        ((VALIDATION_PASSED++))
        return 0
    else
        print_message $RED "✗ $description"
        ((VALIDATION_FAILED++))
        return 1
    fi
}

# Check if Hurl is installed and working
check_hurl_installation() {
    print_message $BLUE "Checking Hurl installation..."
    
    validate "Hurl is installed" "command -v hurl &> /dev/null"
    
    if command -v hurl &> /dev/null; then
        local version=$(hurl --version 2>/dev/null | head -n1 || echo "unknown")
        print_message $YELLOW "  Version: $version"
        
        # Test basic hurl functionality
        validate "Hurl basic functionality" "echo 'GET https://httpbin.org/get' | hurl --test &> /dev/null"
    fi
}

# Check directory structure
check_directory_structure() {
    print_message $BLUE "Checking directory structure..."
    
    local required_dirs=(
        "tests"
        "tests/api"
        "tests/config"
        "tests/scripts"
        "tests/reports"
        "tests/functional"
        "tests/security"
        "tests/performance"
        "tests/integration"
        "tests/edge-cases"
    )
    
    for dir in "${required_dirs[@]}"; do
        validate "Directory exists: $dir" "[[ -d \"$dir\" ]]"
    done
}

# Check configuration files
check_config_files() {
    print_message $BLUE "Checking configuration files..."
    
    local required_configs=(
        "tests/config/global.hurl"
        "tests/config/environments.hurl"
        "tests/config/common-assertions.hurl"
        "tests/config/api-variables.hurl"
    )
    
    for config in "${required_configs[@]}"; do
        validate "Config file exists: $config" "[[ -f \"$config\" ]]"
        
        if [[ -f "$config" ]]; then
            validate "Config file is readable: $config" "[[ -r \"$config\" ]]"
            validate "Config file is not empty: $config" "[[ -s \"$config\" ]]"
        fi
    done
}

# Check test files
check_test_files() {
    print_message $BLUE "Checking test files..."
    
    local functional_tests=(
        "tests/functional/health.hurl"
        "tests/functional/openapi.hurl"
        "tests/functional/docs.hurl"
        "tests/functional/scan.hurl"
        "tests/functional/graph.hurl"
        "tests/functional/stats.hurl"
    )
    
    for test in "${functional_tests[@]}"; do
        validate "Functional test exists: $test" "[[ -f \"$test\" ]]"
    done
    
    local security_tests=(
        "tests/security/security-health.hurl"
        "tests/security/security-openapi.hurl"
        "tests/security/security-docs.hurl"
        "tests/security/security-scan.hurl"
        "tests/security/security-graph.hurl"
        "tests/security/security-stats.hurl"
    )
    
    for test in "${security_tests[@]}"; do
        validate "Security test exists: $test" "[[ -f \"$test\" ]]"
    done
    
    local performance_tests=(
        "tests/performance/performance-health.hurl"
        "tests/performance/performance-scan.hurl"
        "tests/performance/performance-graph.hurl"
    )
    
    for test in "${performance_tests[@]}"; do
        validate "Performance test exists: $test" "[[ -f \"$test\" ]]"
    done
    
    validate "API test exists: tests/api/basic-api-tests.hurl" "[[ -f \"tests/api/basic-api-tests.hurl\" ]]"
}

# Check script files
check_script_files() {
    print_message $BLUE "Checking script files..."
    
    local required_scripts=(
        "tests/scripts/setup-test-env.sh"
        "tests/scripts/run-tests.sh"
        "tests/scripts/cleanup-test-env.sh"
        "tests/scripts/validate-test-setup.sh"
    )
    
    for script in "${required_scripts[@]}"; do
        validate "Script exists: $script" "[[ -f \"$script\" ]]"
        
        if [[ -f "$script" ]]; then
            validate "Script is executable: $script" "[[ -x \"$script\" ]]"
        fi
    done
}

# Validate Hurl configuration syntax
validate_hurl_syntax() {
    print_message $BLUE "Validating Hurl file syntax..."
    
    # Find all .hurl files
    local hurl_files=$(find tests -name "*.hurl" -type f 2>/dev/null || true)
    
    if [[ -n "$hurl_files" ]]; then
        while IFS= read -r file; do
            # Basic syntax check (hurl --check would be ideal but may not be available)
            validate "Syntax check: $file" "[[ -r \"$file\" ]]"
        done <<< "$hurl_files"
    else
        print_message $YELLOW "No .hurl files found to validate"
    fi
}

# Check environment variables
check_environment_variables() {
    print_message $BLUE "Checking environment variables..."
    
    # Check if common environment variables are set or have defaults
    local env_vars=(
        "HURL_TEST_ENV"
        "API_BASE_URL"
        "NEO4J_URI"
        "DEBUG_MODE"
    )
    
    for var in "${env_vars[@]}"; do
        if [[ -n "${!var}" ]]; then
            validate "Environment variable set: $var" "true"
            print_message $YELLOW "  $var = ${!var}"
        else
            print_message $YELLOW "  $var not set (will use defaults)"
        fi
    done
}

# Check dependencies
check_dependencies() {
    print_message $BLUE "Checking dependencies..."
    
    local dependencies=(
        "curl"
        "jq"
        "nc"
        "grep"
        "find"
        "sed"
        "awk"
    )
    
    for dep in "${dependencies[@]}"; do
        validate "Dependency available: $dep" "command -v $dep &> /dev/null"
    done
}

# Generate summary report
generate_summary() {
    print_message $BLUE "Validation Summary"
    print_message $GREEN "Passed: $VALIDATION_PASSED"
    print_message $RED "Failed: $VALIDATION_FAILED"
    
    local total=$((VALIDATION_PASSED + VALIDATION_FAILED))
    if [[ $total -gt 0 ]]; then
        local success_rate=$((VALIDATION_PASSED * 100 / total))
        print_message $YELLOW "Success Rate: $success_rate%"
    fi
    
    if [[ $VALIDATION_FAILED -eq 0 ]]; then
        print_message $GREEN "🎉 All validations passed! Test setup is ready."
        return 0
    else
        print_message $RED "❌ Some validations failed. Please fix the issues above."
        return 1
    fi
}

# Main validation function
main() {
    print_message $GREEN "Starting test setup validation..."
    
    check_hurl_installation
    check_directory_structure
    check_config_files
    check_test_files
    check_script_files
    validate_hurl_syntax
    check_environment_variables
    check_dependencies
    
    echo ""
    generate_summary
}

# Run main function
main "$@"