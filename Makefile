.PHONY: build test clean lint fmt run all help

# Output binary name
BINARY_NAME=otel-checker

# Go compiler and tools
GO=go
GOFMT=gofmt
GOLINT=golangci-lint
GOTEST=$(GO) test

# Directories
SRC_DIR=.

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) install

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@$(GO) clean

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -w $(SRC_DIR)

# Lint code
lint:
	@echo "Linting code..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run; \
	else \
		echo "golangci-lint not installed. Please install with:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Update dependencies
deps:
	@echo "Updating dependencies..."
	@$(GO) mod tidy

# Help command
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make fmt      - Format the code"
	@echo "  make lint     - Lint the code"
	@echo "  make deps     - Update dependencies"
	@echo "  make all      - Build the application"
	@echo "  make help     - Show this help message"
