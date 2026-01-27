# Personal Music Search Engine - Development Makefile
#
# Usage:
#   make local           - Start full local development environment
#   make local-stop      - Stop all local services
#   make test-integration - Run integration tests against LocalStack
#
# Prerequisites:
#   - Docker and docker-compose
#   - Go 1.22+
#   - Node.js 18+
#   - AWS CLI (for LocalStack init scripts)

.PHONY: local local-services local-backend local-frontend local-stop test-integration clean help

# Default target
help:
	@echo "Personal Music Search Engine - Development Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make local            Start full local development environment"
	@echo "  make local-services   Start LocalStack only"
	@echo "  make local-backend    Start backend server (requires LocalStack)"
	@echo "  make local-frontend   Start frontend dev server"
	@echo "  make local-stop       Stop all local services"
	@echo "  make test-integration Run integration tests against LocalStack"
	@echo "  make clean            Clean build artifacts"
	@echo ""
	@echo "Test Users (LocalStack):"
	@echo "  admin@local.test      Password: LocalTest123!"
	@echo "  subscriber@local.test Password: LocalTest123!"
	@echo "  artist@local.test     Password: LocalTest123!"

# Start full local environment
local: local-services
	@echo ""
	@echo "Starting backend and frontend..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Frontend will be available at http://localhost:5173"
	@echo ""
	@$(MAKE) -j2 local-backend local-frontend

# Start LocalStack services
local-services:
	@echo "Starting LocalStack..."
	docker-compose -f docker/docker-compose.yml up -d
	@echo "Waiting for LocalStack to be healthy..."
	./scripts/wait-for-localstack.sh 60
	@echo ""
	@echo "Running initialization scripts..."
	./docker/localstack-init/init-aws.sh
	./docker/localstack-init/init-cognito.sh
	@echo ""
	@echo "LocalStack is ready!"

# Start backend against LocalStack
local-backend:
	@echo "Starting backend server..."
	cd backend && \
		AWS_ENDPOINT=http://localhost:4566 \
		DYNAMODB_TABLE_NAME=MusicLibrary \
		MEDIA_BUCKET=music-library-local-media \
		go run ./cmd/api

# Start frontend in local mode
local-frontend:
	@echo "Starting frontend dev server..."
	@if [ ! -f frontend/.env.local ]; then \
		echo "WARNING: frontend/.env.local not found"; \
		echo "Copy frontend/.env.local.example and fill in Cognito values from init-cognito.sh output"; \
	fi
	cd frontend && npm run dev:local

# Stop all local services
local-stop:
	@echo "Stopping local services..."
	-@pkill -f "go run ./cmd/api" 2>/dev/null || true
	-@pkill -f "vite" 2>/dev/null || true
	docker-compose -f docker/docker-compose.yml down
	@echo "Local services stopped."

# Run integration tests
test-integration:
	@echo "Starting LocalStack for integration tests..."
	docker-compose -f docker/docker-compose.yml up -d
	./scripts/wait-for-localstack.sh 60
	./docker/localstack-init/init-aws.sh
	@echo ""
	@echo "Running integration tests..."
	cd backend && go test -tags=integration -v ./...
	@echo ""
	@echo "Integration tests complete."

# Run unit tests only
test:
	@echo "Running unit tests..."
	cd backend && go test -short ./...
	cd frontend && npm test -- --run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bootstrap backend/lambda.zip
	rm -rf frontend/dist
	@echo "Clean complete."

# Reset LocalStack data
local-reset:
	@echo "Resetting LocalStack data..."
	docker-compose -f docker/docker-compose.yml down -v
	@echo "LocalStack data cleared. Run 'make local' to restart."
