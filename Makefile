# GoOSQuery Makefile

# Variables
BINARY_NAME=goosquery
OUTPUT_DIR=goosquery_output
VERSION=1.0.0
# Use a hardcoded date for Windows compatibility
BUILD_TIME="$(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH-mm-ss'")"
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Directories
SRC_DIR=.
BUILD_DIR=build

# Default target - for Windows
.PHONY: all
all: clean build-windows

# Build the application for Windows (default)
.PHONY: build
build: build-windows

# Build for Windows specifically
.PHONY: build-windows
build-windows:
	@echo "Building GoOSQuery for Windows..."
	@powershell "if (-not (Test-Path $(BUILD_DIR))) { New-Item -ItemType Directory -Path $(BUILD_DIR) -Force }"
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(SRC_DIR)
	@echo "Windows build complete: $(BUILD_DIR)/$(BINARY_NAME).exe"

# Run the application
.PHONY: run
run: build-windows
	@echo "Running GoOSQuery..."
	@powershell "if (-not (Test-Path $(OUTPUT_DIR))) { New-Item -ItemType Directory -Path $(OUTPUT_DIR) -Force }"
	$(BUILD_DIR)/$(BINARY_NAME).exe $(OUTPUT_DIR)

# Run with debug mode
.PHONY: debug
debug: build-windows
	@echo "Running GoOSQuery in debug mode..."
	@powershell "if (-not (Test-Path $(OUTPUT_DIR))) { New-Item -ItemType Directory -Path $(OUTPUT_DIR) -Force }"
	$(BUILD_DIR)/$(BINARY_NAME).exe $(OUTPUT_DIR) --debug

# Run benchmarks
.PHONY: benchmark
benchmark: build-windows
	@echo "Running GoOSQuery benchmarks..."
	@powershell "if (-not (Test-Path $(OUTPUT_DIR))) { New-Item -ItemType Directory -Path $(OUTPUT_DIR) -Force }"
	$(BUILD_DIR)/$(BINARY_NAME).exe $(OUTPUT_DIR) --benchmark 5

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@powershell "if (Test-Path $(BUILD_DIR)) { Remove-Item -Recurse -Force $(BUILD_DIR) }"
	@echo "Clean complete"

# Clean output directory
.PHONY: clean-output
clean-output:
	@echo "Cleaning output directory..."
	@powershell "if (Test-Path $(OUTPUT_DIR)) { Remove-Item -Recurse -Force $(OUTPUT_DIR) }"
	@echo "Output directory cleaned"

# Clean everything
.PHONY: clean-all
clean-all: clean clean-output
	@echo "All cleaned"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...
	@echo "Dependencies installed"

# Help target
.PHONY: help
help:
	@echo "GoOSQuery Makefile Help (Windows-focused)"
	@echo "---------------------------------------"
	@echo "Available targets:"
	@echo "  all          : Clean and build the application for Windows (default)"
	@echo "  build        : Build the application for Windows"
	@echo "  build-windows: Build for Windows"
	@echo "  run          : Build and run the application"
	@echo "  debug        : Build and run the application in debug mode"
	@echo "  benchmark    : Build and run benchmarks"
	@echo "  clean        : Clean build artifacts"
	@echo "  clean-output : Clean output directory"
	@echo "  clean-all    : Clean everything (build artifacts and output)"
	@echo "  test         : Run tests"
	@echo "  deps         : Install dependencies"
	@echo "  help         : Show this help message"
