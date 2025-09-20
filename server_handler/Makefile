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

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run the application in development mode
dev:
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_PATH) $(CONFIG_FILE)

# Run the built binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(CONFIG_FILE)

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

# Quick development setup
setup: deps
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Development setup complete!"
	@echo "Edit $(CONFIG_FILE) with your configuration, then run 'make dev' to start development"

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Clean, install deps, and build"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install Go dependencies"
	@echo "  dev          - Run in development mode"
	@echo "  run          - Build and run the binary"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  setup        - Quick development setup"
	@echo "  help         - Show this help"