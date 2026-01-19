# Backend - CLAUDE.md

## Overview

Go backend services for the Personal Music Search Engine. Implements a serverless API using AWS Lambda with Echo framework, DynamoDB for data storage, and S3 for media files.

## Directory Structure

```
backend/
├── api/                    # OpenAPI specification
│   └── openapi.yaml        # API contract definition
├── cmd/                    # Lambda entrypoints
│   ├── api/                # Main API Lambda
│   ├── indexer/            # Search indexer Lambda
│   └── processor/          # Upload processor Step Functions Lambdas
└── internal/               # Internal packages (not exported)
    ├── handlers/           # HTTP request handlers
    ├── metadata/           # Audio metadata extraction
    ├── models/             # Domain models and DTOs
    ├── repository/         # Data access layer
    ├── search/             # Nixiesearch client
    └── service/            # Business logic layer
```

## File Descriptions

| File | Purpose |
|------|---------|
| `go.mod` | Go module definition with dependencies |
| `go.sum` | Dependency checksums |

## Dependencies

### External Dependencies
| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/aws/aws-lambda-go` | v1.47.0 | Lambda runtime |
| `github.com/aws/aws-sdk-go-v2` | v1.30.3 | AWS SDK |
| `github.com/labstack/echo/v4` | v4.12.0 | HTTP framework |
| `github.com/awslabs/aws-lambda-go-api-proxy` | v0.16.2 | Echo-Lambda adapter |
| `github.com/dhowden/tag` | latest | Audio metadata extraction |
| `github.com/google/uuid` | v1.6.0 | UUID generation |
| `github.com/stretchr/testify` | v1.9.0 | Testing assertions |

### Internal Dependencies
- `internal/models` - Domain models
- `internal/repository` - Data access
- `internal/service` - Business logic
- `internal/handlers` - HTTP handlers

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Lambda Handler                          │
│                    (cmd/api/main.go)                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Echo Router                                │
│              (internal/handlers/*)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Service Layer                               │
│              (internal/service/*)                            │
│           Business logic, validation                         │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                Repository Layer                              │
│             (internal/repository/*)                          │
│            DynamoDB, S3 operations                           │
└─────────────────────────────────────────────────────────────┘
```

## Build Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build API Lambda
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/api/main.go

# Build all processors
for dir in cmd/processor/*/; do
  GOOS=linux GOARCH=amd64 go build -o "${dir}bootstrap" "${dir}main.go"
done
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DYNAMODB_TABLE_NAME` | DynamoDB table name | `MusicLibrary` |
| `MEDIA_BUCKET` | S3 media bucket | `music-library-media` |
| `CLOUDFRONT_DOMAIN` | CloudFront domain | - |
| `STEP_FUNCTIONS_ARN` | Upload processor ARN | - |
| `SEARCH_INDEX_BUCKET` | Nixiesearch index bucket | - |

## Testing Strategy

All code must follow TDD:
1. Write failing tests first
2. Implement minimal code to pass
3. Refactor while keeping tests green

### Test File Naming
- Unit tests: `*_test.go` in same package
- Integration tests: `*_integration_test.go`

### Running Tests
```bash
# Unit tests only
go test -short ./...

# All tests including integration
go test ./...

# Specific package
go test ./internal/models/...
```
