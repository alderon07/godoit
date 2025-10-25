.PHONY: help build run test clean install fmt vet lint all

# Binary name
BINARY_NAME=godo
OUTPUT_DIR=bin

# Version info (can be overridden)
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
all: fmt vet test build

## help: Display this help message
help:
	@echo "Available targets:"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build the binary"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make dev            - Run with auto-reload (requires entr)"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make all            - Format, vet, test, and build"
	@echo ""
	@echo "Cross-platform Builds:"
	@echo "  make build-all      - Build for all platforms"
	@echo "  make build-linux    - Build for Linux (amd64, arm64)"
	@echo "  make build-darwin   - Build for macOS (amd64, arm64)"
	@echo "  make build-windows  - Build for Windows (amd64)"
	@echo "  make build-platform - Build for specific platform (GOOS=... GOARCH=...)"
	@echo "  make package        - Create compressed archives"
	@echo "  make checksums      - Generate checksums"
	@echo ""
	@echo "Other:"
	@echo "  make clean          - Remove built binaries"
	@echo "  make install        - Install binary to GOPATH/bin"
	@echo "  make deps           - Download dependencies"
	@echo ""
	@echo "Examples:"
	@echo "  make build VERSION=1.0.0"
	@echo "  make run ARGS='list -all'"
	@echo "  make build-platform GOOS=linux GOARCH=arm64"

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/todo
	@echo "Binary created at $(OUTPUT_DIR)/$(BINARY_NAME)"

## run: Run the application
run:
	@go run ./cmd/todo/main.go $(ARGS)

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(OUTPUT_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

## fmt: Format all Go files
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

## dev: Run in development mode with auto-reload (requires entr)
dev:
	@if command -v entr > /dev/null; then \
		find . -name "*.go" | entr -r go run ./cmd/todo/main.go; \
	else \
		echo "entr not installed. Install it for auto-reload support"; \
		echo "On Ubuntu/Debian: sudo apt-get install entr"; \
	fi

## build-all: Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "All platform builds complete"

## build-linux: Build for Linux
build-linux:
	@echo "Building for Linux amd64..."
	@mkdir -p $(OUTPUT_DIR)/linux-amd64
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/linux-amd64/$(BINARY_NAME) ./cmd/todo
	@echo "Building for Linux arm64..."
	@mkdir -p $(OUTPUT_DIR)/linux-arm64
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/linux-arm64/$(BINARY_NAME) ./cmd/todo

## build-darwin: Build for macOS
build-darwin:
	@echo "Building for macOS amd64..."
	@mkdir -p $(OUTPUT_DIR)/darwin-amd64
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/darwin-amd64/$(BINARY_NAME) ./cmd/todo
	@echo "Building for macOS arm64..."
	@mkdir -p $(OUTPUT_DIR)/darwin-arm64
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/darwin-arm64/$(BINARY_NAME) ./cmd/todo

## build-windows: Build for Windows
build-windows:
	@echo "Building for Windows amd64..."
	@mkdir -p $(OUTPUT_DIR)/windows-amd64
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/windows-amd64/$(BINARY_NAME).exe ./cmd/todo

## build-platform: Build for a specific platform (use: make build-platform GOOS=linux GOARCH=amd64)
build-platform:
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Error: GOOS and GOARCH must be set"; \
		echo "Example: make build-platform GOOS=linux GOARCH=amd64"; \
		exit 1; \
	fi
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)/$(GOOS)-$(GOARCH)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(GOOS)-$(GOARCH)/$(BINARY_NAME)$(if $(filter windows,$(GOOS)),.exe,) ./cmd/todo
	@echo "Binary created at $(OUTPUT_DIR)/$(GOOS)-$(GOARCH)/"

## package: Create compressed archives for all platform builds
package: build-all
	@echo "Creating release packages..."
	@cd $(OUTPUT_DIR) && \
	for dir in */; do \
		platform=$${dir%/}; \
		if [ "$$platform" != "$(BINARY_NAME)" ]; then \
			echo "Packaging $$platform..."; \
			if echo "$$platform" | grep -q "windows"; then \
				zip -q -r $(BINARY_NAME)-$(VERSION)-$$platform.zip $$platform; \
			else \
				tar -czf $(BINARY_NAME)-$(VERSION)-$$platform.tar.gz $$platform; \
			fi; \
		fi; \
	done
	@echo "Packages created in $(OUTPUT_DIR)/"

## checksums: Generate checksums for release packages
checksums:
	@echo "Generating checksums..."
	@cd $(OUTPUT_DIR) && \
	if ls *.tar.gz >/dev/null 2>&1 || ls *.zip >/dev/null 2>&1; then \
		shasum -a 256 *.tar.gz *.zip 2>/dev/null > checksums.txt || \
		sha256sum *.tar.gz *.zip 2>/dev/null > checksums.txt; \
		echo "Checksums saved to $(OUTPUT_DIR)/checksums.txt"; \
	else \
		echo "No packages found. Run 'make package' first."; \
	fi
