#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Helper Functions ---
print_header() {
    echo ""
    echo "================================================="
    echo "  $1"
    echo "================================================="
}

print_success() {
    echo "✅ SUCCESS: $1"
}

print_fail() {
    echo "❌ FAILURE: $1" >&2
    exit 1
}

# --- Main Test Runner ---
print_header "DSL BYPASS ULTRA - UNIT TEST SUITE"

# Use python's unittest discovery to find and run all tests
# The -s flag specifies the start directory for discovery.
# The -p flag specifies the pattern for test files.
if python3 -m unittest discover -s tests -p "test_*.py"; then
    print_success "All Python unit tests passed."
else
    print_fail "Python unit tests failed."
fi

# You can still run non-python tests separately if needed
# For example, the performance shell script
print_header "RUNNING PERFORMANCE VALIDATION"
if bash "tests/suite/test_performance.sh"; then
    print_success "Performance validation passed."
else
    print_fail "Performance validation failed."
fi

print_header "ALL TESTS PASSED"
echo "✅ Deployment validation successful."
exit 0