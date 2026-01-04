# Makefile for exsconvert

.PHONY: all build test test-cov fmt vet clean deps help

# Default target
all: deps build

# Download dependencies
deps:
	go mod download

# Build the binary
build:
	go build -v -o exsconvert .

# Run tests
test:
	go test ./...

# Run tests with coverage
test-cov:
	go test -coverprofile=coverage.out ./...

# Format code
fmt:
	go fmt ./...

# Static analysis
vet:
	go vet ./...

# Clean build artifacts
clean:
	go clean
	rm -f exsconvert coverage.out

# Clean build and rebuild from scratch
clean-build: clean deps build test

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Download dependencies and build the binary"
	@echo "  deps         - Download Go module dependencies"
	@echo "  build        - Build the exsconvert binary"
	@echo "  test         - Run all tests"
	@echo "  test-cov     - Run tests with coverage report"
	@echo "  fmt          - Format Go code"
	@echo "  vet          - Run static analysis"
	@echo "  clean        - Clean build artifacts"
	@echo "  clean-build  - Clean, download deps, build, and test"
	@echo "  help         - Show this help message"