# Covo Bot Makefile

# Variables
BINARY_NAME=covo-bot
BUILD_DIR=build
DOCKER_IMAGE=covo-bot
DOCKER_TAG=latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags="-s -w"
BUILD_FLAGS=-a -installsuffix cgo

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	@echo "Multi-platform build completed"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run main.go

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		$(GOGET) -u github.com/cosmtrek/air; \
		air; \
	fi

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. ./...

# Run race detection
.PHONY: race
race:
	@echo "Running race detection..."
	$(GOTEST) -race ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
		golangci-lint run; \
	fi

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) -u github.com/cosmtrek/air
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/go-delve/delve/cmd/dlv

# Docker commands
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(BINARY_NAME) --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(BINARY_NAME) || true
	docker rm $(BINARY_NAME) || true

.PHONY: docker-logs
docker-logs:
	@echo "Showing Docker logs..."
	docker logs -f $(BINARY_NAME)

# Docker Compose commands
.PHONY: compose-up
compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

.PHONY: compose-down
compose-down:
	@echo "Stopping services with Docker Compose..."
	docker-compose down

.PHONY: compose-logs
compose-logs:
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

.PHONY: compose-build
compose-build:
	@echo "Building services with Docker Compose..."
	docker-compose build

# Database commands
.PHONY: db-migrate
db-migrate:
	@echo "Running database migrations..."
	$(GOCMD) run main.go --migrate

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	$(GOCMD) run main.go --seed

# Release commands
.PHONY: release
release: clean build-all
	@echo "Creating release..."
	@mkdir -p release
	@cp $(BUILD_DIR)/* release/
	@cp README.md release/
	@cp docker-compose.yml release/
	@cp Dockerfile release/
	@echo "Release created in release/ directory"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode with hot reload"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  bench          - Run benchmarks"
	@echo "  race           - Run race detection"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  deps           - Download dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  install-tools  - Install development tools"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-stop    - Stop Docker container"
	@echo "  docker-logs    - Show Docker logs"
	@echo "  compose-up     - Start services with Docker Compose"
	@echo "  compose-down   - Stop services with Docker Compose"
	@echo "  compose-logs   - Show Docker Compose logs"
	@echo "  compose-build  - Build services with Docker Compose"
	@echo "  db-migrate     - Run database migrations"
	@echo "  db-seed        - Seed database"
	@echo "  release        - Create release package"
	@echo "  help           - Show this help message"
