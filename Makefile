.PHONY: all build test clean run

# Binary name
BINARY_NAME=goosquery.exe

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main build target
all: clean build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

# Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run the program
run: build
	@echo "Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf zdata

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...

# Build for release (with optimizations)
release:
	@echo "Building release version..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v -ldflags="-s -w"

# Show help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make run      - Build and run the application"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make release  - Build optimized release version"
	@echo "  make help     - Show this help message"
