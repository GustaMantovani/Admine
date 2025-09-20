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

.PHONY: all build clean test deps run dev install help

# Default target
all: clean deps build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Built $(BINARY_NAME) in $(BUILD_DIR)/"

# Build for different platforms
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

# Build for all platforms
build-all: build-linux build-windows build-mac

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the application in development mode
dev:
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_PATH) $(CONFIG_FILE)

# Run the built binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(CONFIG_FILE)

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(MAIN_PATH)

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Check if required tools are installed
check-deps:
	@echo "Checking dependencies..."
	@which go > /dev/null || (echo "Go is not installed" && exit 1)
	@echo "All dependencies are installed"

# Quick development setup
setup: deps
	@echo "Development setup complete!"
	@echo "Edit $(CONFIG_FILE) with your configuration, then run 'make dev' to start development"

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Clean, install deps, and build"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for all platforms (Linux, Windows, macOS)"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install Go dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  dev          - Run in development mode"
	@echo "  run          - Build and run the binary"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  check-deps   - Check if required tools are installed"
	@echo "  setup        - Quick development setup"
	@echo "  help         - Show this help"