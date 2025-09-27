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
print_header "K2Ray DEPLOYMENT VALIDATION SUITE"

# Get the root directory of the project
ROOT_DIR=$(git rev-parse --show-toplevel)
TEST_SUITE_DIR="$ROOT_DIR/tests/suite"

# Run Integration Tests
print_header "RUNNING INTEGRATION TESTS"
if python3 "$TEST_SUITE_DIR/test_integration.py"; then
    print_success "Integration tests passed."
else
    print_fail "Integration tests failed."
fi

# Run Safety System Tests
print_header "RUNNING SAFETY SYSTEM TESTS"
if python3 "$TEST_SUITE_DIR/test_safety_checks.py"; then
    print_success "Safety system tests passed."
else
    print_fail "Safety system tests failed."
fi

# Run Performance Validation
print_header "RUNNING PERFORMANCE VALIDATION"
if bash "$TEST_SUITE_DIR/test_performance.sh"; then
    print_success "Performance validation passed."
else
    print_fail "Performance validation failed."
fi


print_header "ALL TESTS PASSED"
echo "✅ Deployment validation successful."
exit 0