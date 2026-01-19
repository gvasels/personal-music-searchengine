# Backend Makefile
# Build and run commands for the music library backend

.PHONY: all build build-api build-processors test clean run-local deps lint fmt help

# Variables
GOOS ?= linux
GOARCH ?= arm64
CGO_ENABLED ?= 0
GO_BUILD_FLAGS := -ldflags="-s -w" -trimpath

# Processor directories
PROCESSOR_DIRS := metadata coverart track mover indexer status

# Default target
all: build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint code
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Build all binaries
build: build-api build-processors

# Build API Lambda
build-api:
	@echo "Building API Lambda..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		go build $(GO_BUILD_FLAGS) -o bin/api/bootstrap ./cmd/api

# Build all processor Lambdas
build-processors: $(addprefix build-processor-,$(PROCESSOR_DIRS))

# Build individual processor Lambda
build-processor-%:
	@echo "Building processor: $*..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		go build $(GO_BUILD_FLAGS) -o bin/processor/$*/bootstrap ./cmd/processor/$*

# Build for local development (native architecture)
build-local:
	@echo "Building for local development..."
	CGO_ENABLED=0 go build -o bin/api-local ./cmd/api

# Run API locally (requires LocalStack)
run-local: build-local
	@echo "Starting API locally..."
	@echo "Make sure LocalStack is running: cd ../docker && docker-compose up -d"
	DYNAMODB_TABLE_NAME=MusicLibrary \
	MEDIA_BUCKET=music-library-local-media \
	AWS_REGION=us-east-1 \
	LOCAL_MODE=true \
	LOCALSTACK_ENDPOINT=http://localhost:4566 \
	./bin/api-local

# Start LocalStack
localstack-up:
	@echo "Starting LocalStack..."
	cd ../docker && docker-compose up -d

# Stop LocalStack
localstack-down:
	@echo "Stopping LocalStack..."
	cd ../docker && docker-compose down

# View LocalStack logs
localstack-logs:
	cd ../docker && docker-compose logs -f localstack

# Initialize LocalStack resources
localstack-init:
	@echo "Initializing LocalStack resources..."
	cd ../docker && bash localstack-init/init-aws.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Package Lambda functions for deployment
package: build
	@echo "Packaging Lambda functions..."
	@mkdir -p dist
	cd bin/api && zip ../../dist/api.zip bootstrap
	@for dir in $(PROCESSOR_DIRS); do \
		cd bin/processor/$$dir && zip ../../../dist/processor-$$dir.zip bootstrap && cd ../../..; \
	done
	@echo "Packages created in dist/"

# Show help
help:
	@echo "Available targets:"
	@echo "  make deps           - Install dependencies"
	@echo "  make build          - Build all Lambda binaries (Linux ARM64)"
	@echo "  make build-api      - Build API Lambda only"
	@echo "  make build-processors - Build all processor Lambdas"
	@echo "  make build-local    - Build for local development (native arch)"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make run-local      - Run API locally (requires LocalStack)"
	@echo "  make localstack-up  - Start LocalStack"
	@echo "  make localstack-down - Stop LocalStack"
	@echo "  make localstack-init - Initialize LocalStack resources"
	@echo "  make package        - Package Lambdas for deployment"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make help           - Show this help"
