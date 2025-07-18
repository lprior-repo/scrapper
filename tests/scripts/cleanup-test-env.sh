#!/bin/bash

# Test Environment Cleanup Script
# This script cleans up test artifacts and temporary files

set -e

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

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --reports      Clean test reports only"
    echo "  --logs         Clean log files only"
    echo "  --all          Clean everything (default)"
    echo "  --dry-run      Show what would be cleaned without actually cleaning"
    echo "  -h, --help     Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 --reports"
    echo "  $0 --dry-run"
    echo "  $0 --all"
}

# Parse command line arguments
parse_args() {
    CLEAN_REPORTS=false
    CLEAN_LOGS=false
    CLEAN_ALL=true
    DRY_RUN=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --reports)
                CLEAN_REPORTS=true
                CLEAN_ALL=false
                shift
                ;;
            --logs)
                CLEAN_LOGS=true
                CLEAN_ALL=false
                shift
                ;;
            --all)
                CLEAN_ALL=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_message $RED "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Clean test reports
clean_reports() {
    print_message $YELLOW "Cleaning test reports..."
    
    local report_files=(
        "tests/reports/*.html"
        "tests/reports/*.json"
        "tests/reports/*.xml"
        "tests/reports/*.log"
    )
    
    for pattern in "${report_files[@]}"; do
        if [[ "$DRY_RUN" == "true" ]]; then
            if ls $pattern 1> /dev/null 2>&1; then
                print_message $BLUE "Would remove: $pattern"
            fi
        else
            if ls $pattern 1> /dev/null 2>&1; then
                rm -f $pattern
                print_message $GREEN "Removed: $pattern"
            fi
        fi
    done
}

# Clean log files
clean_logs() {
    print_message $YELLOW "Cleaning log files..."
    
    local log_files=(
        "*.log"
        "tests/**/*.log"
        "logs/*.log"
        "tmp/*.log"
    )
    
    for pattern in "${log_files[@]}"; do
        if [[ "$DRY_RUN" == "true" ]]; then
            if ls $pattern 1> /dev/null 2>&1; then
                print_message $BLUE "Would remove: $pattern"
            fi
        else
            if ls $pattern 1> /dev/null 2>&1; then
                rm -f $pattern
                print_message $GREEN "Removed: $pattern"
            fi
        fi
    done
}

# Clean temporary files
clean_temp_files() {
    print_message $YELLOW "Cleaning temporary files..."
    
    local temp_files=(
        "*.tmp"
        "tests/**/*.tmp"
        "tmp/*"
        ".hurl_*"
    )
    
    for pattern in "${temp_files[@]}"; do
        if [[ "$DRY_RUN" == "true" ]]; then
            if ls $pattern 1> /dev/null 2>&1; then
                print_message $BLUE "Would remove: $pattern"
            fi
        else
            if ls $pattern 1> /dev/null 2>&1; then
                rm -f $pattern
                print_message $GREEN "Removed: $pattern"
            fi
        fi
    done
}

# Clean cache files
clean_cache() {
    print_message $YELLOW "Cleaning cache files..."
    
    local cache_dirs=(
        ".cache"
        "tests/.cache"
        "tmp/cache"
    )
    
    for dir in "${cache_dirs[@]}"; do
        if [[ -d "$dir" ]]; then
            if [[ "$DRY_RUN" == "true" ]]; then
                print_message $BLUE "Would remove directory: $dir"
            else
                rm -rf "$dir"
                print_message $GREEN "Removed directory: $dir"
            fi
        fi
    done
}

# Main cleanup function
main() {
    parse_args "$@"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message $YELLOW "DRY RUN - No files will be actually removed"
    fi
    
    print_message $GREEN "Starting cleanup..."
    
    if [[ "$CLEAN_ALL" == "true" ]]; then
        clean_reports
        clean_logs
        clean_temp_files
        clean_cache
    else
        if [[ "$CLEAN_REPORTS" == "true" ]]; then
            clean_reports
        fi
        
        if [[ "$CLEAN_LOGS" == "true" ]]; then
            clean_logs
        fi
    fi
    
    print_message $GREEN "✅ Cleanup complete!"
}

# Run main function
main "$@"