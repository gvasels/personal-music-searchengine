#!/bin/bash

# Local Development Environment Script
# Alternative to Makefile for shell users
#
# Usage:
#   ./scripts/local-dev.sh start   - Start full environment
#   ./scripts/local-dev.sh stop    - Stop all services
#   ./scripts/local-dev.sh test    - Run integration tests
#   ./scripts/local-dev.sh status  - Check service status
#   ./scripts/local-dev.sh reset   - Reset LocalStack data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to project root
cd "$PROJECT_ROOT"

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

show_help() {
    echo "Personal Music Search Engine - Local Development"
    echo ""
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  start    Start full local development environment"
    echo "  stop     Stop all local services"
    echo "  test     Run integration tests against LocalStack"
    echo "  status   Check service status"
    echo "  reset    Reset LocalStack data (destructive)"
    echo "  help     Show this help message"
    echo ""
    echo "Test Users (LocalStack):"
    echo "  admin@local.test      Password: LocalTest123!"
    echo "  subscriber@local.test Password: LocalTest123!"
    echo "  artist@local.test     Password: LocalTest123!"
}

check_dependencies() {
    local missing=0

    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        missing=1
    fi

    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        missing=1
    fi

    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed"
        missing=1
    fi

    if ! command -v aws &> /dev/null; then
        print_warning "AWS CLI not installed (needed for init scripts)"
    fi

    if [ $missing -eq 1 ]; then
        exit 1
    fi
}

start_localstack() {
    print_header "Starting LocalStack"

    docker-compose -f docker/docker-compose.yml up -d
    print_success "LocalStack container started"

    echo "Waiting for services to be healthy..."
    ./scripts/wait-for-localstack.sh 60

    echo "Running initialization scripts..."
    ./docker/localstack-init/init-aws.sh
    ./docker/localstack-init/init-cognito.sh

    print_success "LocalStack is ready"
}

start_backend() {
    print_header "Starting Backend"

    echo "Backend will be available at http://localhost:8080"

    cd backend
    AWS_ENDPOINT=http://localhost:4566 \
    DYNAMODB_TABLE_NAME=MusicLibrary \
    MEDIA_BUCKET=music-library-local-media \
    go run ./cmd/api &

    cd "$PROJECT_ROOT"
    print_success "Backend starting in background"
}

start_frontend() {
    print_header "Starting Frontend"

    if [ ! -f frontend/.env.local ]; then
        print_warning "frontend/.env.local not found"
        echo "Copy frontend/.env.local.example and fill in Cognito values"
    fi

    echo "Frontend will be available at http://localhost:5173"

    cd frontend
    npm run dev:local &

    cd "$PROJECT_ROOT"
    print_success "Frontend starting in background"
}

cmd_start() {
    print_header "Starting Local Development Environment"
    check_dependencies

    start_localstack
    echo ""
    start_backend
    echo ""
    start_frontend

    echo ""
    print_success "Local environment is starting"
    echo ""
    echo "Services:"
    echo "  LocalStack:  http://localhost:4566"
    echo "  Backend:     http://localhost:8080"
    echo "  Frontend:    http://localhost:5173"
    echo ""
    echo "Use '$0 stop' to stop all services"
}

cmd_stop() {
    print_header "Stopping Local Services"

    echo "Stopping backend..."
    pkill -f "go run ./cmd/api" 2>/dev/null || true

    echo "Stopping frontend..."
    pkill -f "vite" 2>/dev/null || true

    echo "Stopping LocalStack..."
    docker-compose -f docker/docker-compose.yml down

    print_success "All services stopped"
}

cmd_test() {
    print_header "Running Integration Tests"
    check_dependencies

    echo "Starting LocalStack..."
    docker-compose -f docker/docker-compose.yml up -d
    ./scripts/wait-for-localstack.sh 60
    ./docker/localstack-init/init-aws.sh

    echo ""
    echo "Running tests..."
    cd backend
    go test -tags=integration -v ./...

    print_success "Integration tests complete"
}

cmd_status() {
    print_header "Service Status"

    # Check LocalStack
    if curl -s http://localhost:4566/_localstack/health > /dev/null 2>&1; then
        print_success "LocalStack: Running"
        curl -s http://localhost:4566/_localstack/health | grep -o '"[^"]*": "[^"]*"' | head -5
    else
        print_error "LocalStack: Not running"
    fi

    echo ""

    # Check Backend
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "Backend: Running at http://localhost:8080"
    else
        print_error "Backend: Not running"
    fi

    echo ""

    # Check Frontend
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        print_success "Frontend: Running at http://localhost:5173"
    else
        print_error "Frontend: Not running"
    fi
}

cmd_reset() {
    print_header "Resetting LocalStack Data"

    echo "This will delete all LocalStack data. Continue? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Cancelled."
        exit 0
    fi

    docker-compose -f docker/docker-compose.yml down -v
    print_success "LocalStack data cleared"
    echo "Run '$0 start' to restart with fresh data"
}

# Main command handler
case "${1:-help}" in
    start)
        cmd_start
        ;;
    stop)
        cmd_stop
        ;;
    test)
        cmd_test
        ;;
    status)
        cmd_status
        ;;
    reset)
        cmd_reset
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
