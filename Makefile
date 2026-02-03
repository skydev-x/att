# Makefile for your TUI app

APP_NAME := at
VERSION := 1.0.0
BUILD_DIR := build

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

# Build flags
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: all build clean test install uninstall release

all: test build

# Build for current platform
build:
	@echo "Building $(APP_NAME)..."
	@$(GOBUILD) $(LDFLAGS) -o $(APP_NAME) .
	@echo "Build complete: ./$(APP_NAME)"

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(APP_NAME)
	@rm -rf $(BUILD_DIR)

# Install to /usr/local/bin
install: build
	@echo "Installing to /usr/local/bin..."
	@sudo cp $(APP_NAME) /usr/local/bin/
	@echo "Installed! Run '$(APP_NAME)' to start"

# Uninstall from /usr/local/bin
uninstall:
	@echo "Uninstalling..."
	@sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "Uninstalled"

# Build for multiple platforms (for releases)
release:
	@echo "Building releases..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	
	# Linux ARM64
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .
	
	# macOS AMD64 (Intel)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .
	
	@echo "Release builds complete in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)

# Show version
version:
	@echo "$(APP_NAME) version $(VERSION)"

# Run the app
run: build
	@./$(APP_NAME)
