.PHONY: test test-unit test-integration test-all test-coverage test-benchmark clean test-verbose test-short

# Default test command - runs everything
test: test-all

# Run all tests (unit + integration)
test-all:
	@echo "Running complete test suite..."
	@./run-tests.sh

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test ./tests/unit/... -v
	@go test ./internal/models/... -v
	@go test ./internal/utils/... -v

# Run integration tests with mock data
test-integration:
	@echo "Running integration tests with mock data..."
	@go test ./tests/integration/... -v -count=1

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | grep total
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
test-benchmark:
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem -run=^$ ./tests/integration/...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run short tests only (skip integration)
test-short:
	@echo "Running short tests..."
	@go test -short ./...

# Run specific test
test-specific:
	@echo "Running specific test: $(TEST)"
	@go test -v -run $(TEST) ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -race ./...

# Clean test artifacts
clean:
	@rm -f coverage.out coverage.html test-results.json coverage.txt
	@echo "Test artifacts cleaned"

# Install test dependencies
deps:
	@echo "Installing test dependencies..."
	@go get -u github.com/stretchr/testify/mock
	@go get -u github.com/stretchr/testify/assert
	@go get -u github.com/stretchr/testify/suite
	@go get -u github.com/sirupsen/logrus
	@go mod tidy

# Run tests in Docker
test-docker:
	@echo "Running tests in Docker..."
	@docker run --rm -v $(PWD):/app -w /app golang:1.21-alpine sh -c "go test -v ./..."

# Run CI tests
test-ci:
	@echo "Running CI test suite..."
	@./test-ci.sh

# Generate test report
test-report:
	@echo "Generating test report..."
	@go test -v ./... -json > test-report.json
	@echo "Test report saved to test-report.json"

# Check test coverage threshold
test-threshold:
	@echo "Checking coverage threshold..."
	@go test -coverprofile=coverage.out ./...
	@COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//'); \
	if [ "$COVERAGE" -lt "70" ]; then \
		echo "Coverage $COVERAGE% is below 70% threshold"; \
		exit 1; \
	else \
		echo "Coverage $COVERAGE% meets threshold"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Running go vet..."
	@go vet ./...

# Help
help:
	@echo "Available test commands:"
	@echo "  make test              - Run complete test suite (default)"
	@echo "  make test-all          - Run all tests (unit + integration)"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests only"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-benchmark    - Run benchmark tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-short        - Run short tests only"
	@echo "  make test-race         - Run tests with race detection"
	@echo "  make test-specific TEST=TestName - Run specific test"
	@echo "  make test-docker       - Run tests in Docker"
	@echo "  make test-ci           - Run CI test suite"
	@echo "  make test-report       - Generate JSON test report"
	@echo "  make test-threshold    - Check coverage threshold (70%)"
	@echo "  make deps              - Install test dependencies"
	@echo "  make clean             - Clean test artifacts"
	@echo "  make fmt               - Format code"
	@echo "  make vet               - Run go vet"
	@echo "  make help              - Show this help"