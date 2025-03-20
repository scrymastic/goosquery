.PHONY: build test run clean lint help

BINARY_NAME=goosquery
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOLINT=golint

# Main targets
all: test build

build:
	@echo "Building goosquery..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/goosquery

# Run the application
run:
	@echo "Running goosquery in interactive mode..."
	$(GORUN) ./cmd/goosquery/main.go -i

# Run with JSON output
run-json:
	@echo "Running goosquery in interactive mode with JSON output..."
	$(GORUN) ./cmd/goosquery/main.go -i -json

# Run a specific query
query:
	@echo "Running query: SELECT * FROM processes"
	$(GORUN) ./cmd/goosquery/main.go -q "SELECT * FROM processes"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -cover ./...

# Clean build files
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN) ./...

# Lint code
lint:
	@echo "Linting code..."
	$(GOLINT) ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOCMD) mod download

# Generate a build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/goosquery

# Help command
help:
	@echo "goosquery Makefile"
	@echo "Usage:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run in interactive mode"
	@echo "  make run-json      - Run in interactive mode with JSON output"
	@echo "  make query         - Run a sample query"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build files"
	@echo "  make lint          - Lint code"
	@echo "  make deps          - Install dependencies"
	@echo "  make build-windows - Build for Windows"
	@echo "  make help          - Show this help"

# Default target
default: help
