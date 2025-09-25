#!/bin/bash

# CI/CD Test Runner - Simplified version for automated pipelines
# Exit on first failure for CI/CD

set -e

echo "Starting ABC User Management Test Suite..."

cd backend 2>/dev/null || true

# Install dependencies quietly
echo "Installing dependencies..."
go mod download
go mod tidy

# Run all tests with coverage in JSON format for CI parsing
echo "Running all tests..."
go test -v -coverprofile=coverage.out -json ./... > test-results.json

# Generate coverage report
echo "Generating coverage report..."
go tool cover -func=coverage.out > coverage.txt
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')

echo "Test Coverage: $COVERAGE"

# Check if coverage meets minimum threshold (e.g., 70%)
COVERAGE_NUM=${COVERAGE%\%}
if (( $(echo "$COVERAGE_NUM < 70" | bc -l) )); then
    echo "ERROR: Coverage $COVERAGE is below 70% threshold"
    exit 1
fi

echo "All tests passed with $COVERAGE coverage"
exit 0