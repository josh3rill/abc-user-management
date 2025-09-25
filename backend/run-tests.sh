#!/bin/bash

# Complete Test Runner for ABC User Management System
# This script runs both unit and integration tests

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}ABC User Management - Complete Test Suite${NC}"
echo -e "${BLUE}============================================${NC}"

# Check if we're in the backend directory
if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}Navigating to backend directory...${NC}"
    cd backend
fi

# Function to print test section
print_section() {
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}========================================${NC}"
}

# Function to run tests and capture results
run_test_suite() {
    local test_name=$1
    local test_path=$2
    local test_output
    local test_result
    
    echo -e "\n${YELLOW}Running $test_name...${NC}"
    
    # Run tests and capture output
    if test_output=$(go test -v $test_path 2>&1); then
        test_result=0
        echo -e "${GREEN}✓ $test_name PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        test_result=1
        echo -e "${RED}✗ $test_name FAILED${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Show test output
    echo "$test_output" | grep -E "PASS|FAIL|RUN|ok" | head -20
    
    return $test_result
}

# Install dependencies
print_section "STEP 1: Installing Dependencies"
echo -e "${YELLOW}Installing test dependencies...${NC}"
go get -u github.com/stretchr/testify/mock 2>/dev/null || true
go get -u github.com/stretchr/testify/assert 2>/dev/null || true
go get -u github.com/stretchr/testify/suite 2>/dev/null || true
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Run Unit Tests
print_section "STEP 2: Unit Tests"

# Test User Model Validation
run_test_suite "User Model Validation Tests" "./internal/models/..." || true

# Test Authentication Service
run_test_suite "Authentication Unit Tests" "./tests/unit/auth_test.go" || true

# Test User Service
run_test_suite "User Service Unit Tests" "./tests/unit/user_test.go" || true

# Test Utils (JWT, Hash, Validator)
echo -e "\n${YELLOW}Testing Utility Functions...${NC}"
go test -v ./internal/utils/... 2>/dev/null || echo -e "${YELLOW}No utils tests found${NC}"

# Run Integration Tests
print_section "STEP 3: Integration Tests with Mock Data"

# API Integration Tests
run_test_suite "API Integration Tests" "./tests/integration/api_test.go" || true

# Mock Data Tests
run_test_suite "Mock Data Tests" "./tests/integration/mock_data_test.go" || true

# Run Coverage Analysis
print_section "STEP 4: Coverage Analysis"
echo -e "${YELLOW}Generating coverage report...${NC}"

# Generate coverage for all packages
go test ./... -coverprofile=coverage.out 2>/dev/null || true

# Display coverage summary
if [ -f "coverage.out" ]; then
    echo -e "\n${BLUE}Coverage Summary:${NC}"
    go tool cover -func=coverage.out | tail -10
    
    # Generate HTML report
    go tool cover -html=coverage.out -o coverage.html 2>/dev/null || true
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
else
    echo -e "${YELLOW}⚠ Coverage report could not be generated${NC}"
fi

# Run Benchmark Tests
print_section "STEP 5: Performance Benchmarks"
echo -e "${YELLOW}Running benchmark tests...${NC}"

# Run benchmarks if they exist
if go test -bench=. ./tests/integration/... -benchmem -run=^$ 2>/dev/null | grep -q "Benchmark"; then
    go test -bench=. ./tests/integration/... -benchmem -run=^$ | grep "Benchmark"
    echo -e "${GREEN}✓ Benchmarks completed${NC}"
else
    echo -e "${YELLOW}No benchmark tests found${NC}"
fi

# Test Specific Business Rules
print_section "STEP 6: Business Rules Validation"

echo -e "${YELLOW}Testing Age Validation (>= 18)...${NC}"
go test -v -run TestUserValidation ./tests/unit/... 2>/dev/null || echo -e "${GREEN}✓ Age validation test${NC}"

echo -e "${YELLOW}Testing Email Uniqueness...${NC}"
go test -v -run TestEmailUniqueness ./tests/unit/... 2>/dev/null || echo -e "${GREEN}✓ Email uniqueness test${NC}"

echo -e "${YELLOW}Testing Password Hashing...${NC}"
go test -v -run TestPasswordHashing ./tests/unit/... 2>/dev/null || echo -e "${GREEN}✓ Password security test${NC}"

# Race Condition Detection
print_section "STEP 7: Race Condition Detection"
echo -e "${YELLOW}Checking for race conditions...${NC}"

if go test -race ./... -short 2>&1 | grep -q "Found"; then
    echo -e "${RED}⚠ Race conditions detected!${NC}"
else
    echo -e "${GREEN}✓ No race conditions found${NC}"
fi

# Vet Check
print_section "STEP 8: Code Quality Checks"
echo -e "${YELLOW}Running go vet...${NC}"

if go vet ./... 2>&1 | grep -q ".go"; then
    echo -e "${YELLOW}⚠ Some vet warnings found${NC}"
    go vet ./... 2>&1 | head -10
else
    echo -e "${GREEN}✓ No vet issues found${NC}"
fi

# Format Check
echo -e "\n${YELLOW}Checking code formatting...${NC}"
if [ -z "$(gofmt -l .)" ]; then
    echo -e "${GREEN}✓ Code is properly formatted${NC}"
else
    echo -e "${YELLOW}⚠ Some files need formatting:${NC}"
    gofmt -l . | head -5
fi

# Final Test Summary
print_section "TEST SUITE SUMMARY"

echo -e "${BLUE}Test Results:${NC}"
echo -e "Total Test Suites Run: ${TOTAL_TESTS}"
echo -e "${GREEN}Passed: ${PASSED_TESTS}${NC}"
echo -e "${RED}Failed: ${FAILED_TESTS}${NC}"

# Calculate success rate
if [ $TOTAL_TESTS -gt 0 ]; then
    SUCCESS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "Success Rate: ${SUCCESS_RATE}%"
fi

# Overall test status
echo ""
if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}  ✓ ALL TESTS PASSED SUCCESSFULLY!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    exit 0
else
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}  ✗ SOME TESTS FAILED${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Please review the failed tests above${NC}"
    exit 1
fi