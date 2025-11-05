#!/bin/bash

# Validation Middleware Integration Tests
# Tests gRPC request validation using grpcurl against running services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}Error: grpcurl is not installed${NC}"
    echo "Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Check if services are running
if ! docker ps | grep -q "go-chat-auth"; then
    echo -e "${RED}Error: Services are not running${NC}"
    echo "Start services with: make docker-up"
    exit 1
fi

echo "======================================"
echo "Validation Middleware Integration Tests"
echo "======================================"
echo ""

# Test helper function
run_test() {
    local test_name="$1"
    local service_port="$2"
    local service_method="$3"
    local request_data="$4"
    local should_fail="$5"
    local expected_error="$6"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -n "Test $TOTAL_TESTS: $test_name ... "
    
    if [ "$should_fail" = "true" ]; then
        # Test should fail validation
        if grpcurl -plaintext -d "$request_data" "localhost:$service_port" "$service_method" 2>&1 | grep -q "$expected_error"; then
            echo -e "${GREEN}PASS${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}FAIL${NC}"
            echo "  Expected error: $expected_error"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    else
        # Test should pass validation
        if grpcurl -plaintext -d "$request_data" "localhost:$service_port" "$service_method" > /dev/null 2>&1; then
            echo -e "${GREEN}PASS${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}FAIL${NC}"
            echo "  Request should have passed validation"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    fi
}

echo "=== Auth Service Validation Tests ==="
echo ""

# Test 1: Invalid email format
run_test \
    "Auth Register - Invalid email format" \
    "9001" \
    "api.auth.v1.AuthService/Register" \
    '{"email": "not-an-email", "password": "SecurePass123!"}' \
    "true" \
    "email.*valid email"

# Test 2: Password too short
run_test \
    "Auth Register - Password too short" \
    "9001" \
    "api.auth.v1.AuthService/Register" \
    '{"email": "user@example.com", "password": "short"}' \
    "true" \
    "password.*at least 8"

# Test 3: Empty email
run_test \
    "Auth Register - Empty email" \
    "9001" \
    "api.auth.v1.AuthService/Register" \
    '{"email": "", "password": "SecurePass123!"}' \
    "true" \
    "email"

# Test 4: Valid registration request
run_test \
    "Auth Register - Valid request" \
    "9001" \
    "api.auth.v1.AuthService/Register" \
    '{"email": "user@example.com", "password": "SecurePass123!"}' \
    "false" \
    ""

# Test 5: Login with invalid email
run_test \
    "Auth Login - Invalid email" \
    "9001" \
    "api.auth.v1.AuthService/Login" \
    '{"email": "invalid", "password": "SecurePass123!"}' \
    "true" \
    "email.*valid email"

echo ""
echo "=== Users Service Validation Tests ==="
echo ""

# Test 6: Invalid UUID format
run_test \
    "Users CreateProfile - Invalid UUID" \
    "9002" \
    "api.users.v1.UserService/CreateProfile" \
    '{"user_id": "not-a-uuid", "nickname": "john_doe"}' \
    "true" \
    "user_id.*UUID"

# Test 7: Empty nickname
run_test \
    "Users CreateProfile - Empty nickname" \
    "9002" \
    "api.users.v1.UserService/CreateProfile" \
    '{"user_id": "550e8400-e29b-41d4-a716-446655440000", "nickname": ""}' \
    "true" \
    "nickname.*at least 1"

# Test 8: Valid profile creation
run_test \
    "Users CreateProfile - Valid request" \
    "9002" \
    "api.users.v1.UserService/CreateProfile" \
    '{"user_id": "550e8400-e29b-41d4-a716-446655440000", "nickname": "john_doe"}' \
    "false" \
    ""

# Test 9: Valid profile with avatar URL
run_test \
    "Users CreateProfile - Valid with avatar URL" \
    "9002" \
    "api.users.v1.UserService/CreateProfile" \
    '{"user_id": "550e8400-e29b-41d4-a716-446655440000", "nickname": "jane_doe", "avatar_url": "https://example.com/avatar.jpg"}' \
    "false" \
    ""

# Test 10: Invalid avatar URL
run_test \
    "Users CreateProfile - Invalid avatar URL" \
    "9002" \
    "api.users.v1.UserService/CreateProfile" \
    '{"user_id": "550e8400-e29b-41d4-a716-446655440000", "nickname": "john", "avatar_url": "not-a-url"}' \
    "true" \
    "avatar_url.*URI"

# Test 11: Search query too long
run_test \
    "Users Search - Query too long" \
    "9002" \
    "api.users.v1.UserService/SearchByNickname" \
    "{\"query\": \"$(printf 'a%.0s' {1..150})\", \"limit\": 20}" \
    "true" \
    "query.*at most 100"

# Test 12: Search limit out of range
run_test \
    "Users Search - Limit too high" \
    "9002" \
    "api.users.v1.UserService/SearchByNickname" \
    '{"query": "john", "limit": 150}' \
    "true" \
    "limit.*less than or equal to 100"

# Test 13: Valid search request
run_test \
    "Users Search - Valid request" \
    "9002" \
    "api.users.v1.UserService/SearchByNickname" \
    '{"query": "john", "limit": 20}' \
    "false" \
    ""

echo ""
echo "======================================"
echo "Test Results"
echo "======================================"
echo "Total:  $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
echo "======================================"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi

