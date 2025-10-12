# Makefile for Server Handler

# Variables
BINARY_NAME=server_handler
MAIN_PATH=./cmd/server_handler
BUILD_DIR=./bin
CONFIG_FILE=server_handler_config.yaml

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean deps run dev install fmt vet lint setup test test-coverage help

# Help target
# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: ## Clean, install deps, and build
all: clean deps build

# Build the binary
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Built $(BINARY_NAME) in $(BUILD_DIR)/"

# Install dependencies
install: ## Install and tidy dependencies
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run the application in development mode
dev: ## Run in development mode
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_PATH) $(CONFIG_FILE)

# Run the built binary
run: ## Build and run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(CONFIG_FILE)

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet code
vet: ## Vet Go code
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run linter (requires golangci-lint)
lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

# Run tests
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

check: ## Run all code checks
	@make vet
	@make lint
	@staticcheck ./...

fix: ## Run all code quality fixes
	@make fmt
	@staticcheck ./...

# Quick development setup
setup: ## Setup development environment
setup: deps
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Development setup complete!"
	@echo "Edit $(CONFIG_FILE) with your configuration, then run 'make dev' to start development"
