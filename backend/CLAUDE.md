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

# Build API Lambda (ARM64 - default for AWS Lambda)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/api/main.go

# Build all processors (ARM64)
for dir in cmd/processor/*/; do
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o "${dir}bootstrap" "${dir}main.go"
done

# Deploy to Lambda
zip -j function.zip bootstrap
aws lambda update-function-code --function-name <function-name> --zip-file fileb://function.zip
rm bootstrap function.zip
```

### Build Flags Explained
| Flag | Purpose |
|------|---------|
| `GOOS=linux` | Target Linux OS (Lambda runtime) |
| `GOARCH=arm64` | Target ARM64 architecture (Graviton2) - **default** |
| `CGO_ENABLED=0` | Disable CGO for static linking |
| `-ldflags="-s -w"` | Strip debug info for smaller binary |

> **Note**: All Lambda functions use ARM64 (Graviton2) unless explicitly configured otherwise. Use `GOARCH=amd64` only if the Lambda is configured for x86_64.

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
- Integration tests: `*_integration_test.go` with `//go:build integration` tag

### Running Tests
```bash
# Unit tests only
go test ./...

# Integration tests (requires LocalStack running on port 4566)
go test -tags=integration ./internal/repository/ ./internal/service/ ./test/

# All tests
go test -tags=integration ./...

# Specific package
go test ./internal/models/...
```

### Integration Test Infrastructure (`internal/testutil/`)

| File | Purpose |
|------|---------|
| `localstack.go` | `SetupLocalStack(t)` — creates `TestContext` with DynamoDB, S3, Cognito clients pointing to LocalStack |
| `server.go` | `SetupTestServer(t)` — full Echo HTTP server backed by LocalStack for API testing |
| `http_helpers.go` | `AsUser()`, `WithJSON()`, `DoRequest()`, `AssertStatus()`, `DecodeJSON[T]()` request helpers |
| `fixtures.go` | Entity creation helpers: `CreateTestTrack`, `CreateTestUser`, `CreateTestPlaylist`, `CreateTestArtistProfile`, `CreateTestFollow`, `CreateTestTag`, `CreateTestAlbum`, `CreateTestS3Object` |
| `cleanup.go` | Automatic cleanup of DynamoDB items and S3 objects after tests |

### Integration Test Locations

| Package | File | Tests |
|---------|------|-------|
| `repository` | `dynamodb_integration_test.go` | DynamoDB CRUD + pagination |
| `repository` | `artist_follow_integration_test.go` | Artist profile and follow GSI queries |
| `repository` | `s3_integration_test.go` | S3 ops, presigned URLs, prefix delete |
| `service` | `track_service_integration_test.go` | Visibility enforcement, admin access |
| `service` | `playlist_service_integration_test.go` | Playlist CRUD, public discovery |
| `service` | `other_services_integration_test.go` | Tag, User, Role, Follow, ArtistProfile |
| `service` | `cognito_integration_test.go` | Cognito auth, groups, user management |
| `test` | `api_auth_integration_test.go` | Auth middleware, role-based access |
| `test` | `api_tracks_integration_test.go` | Track CRUD and visibility via HTTP |
| `test` | `api_playlists_tags_integration_test.go` | Playlist and tag endpoints |
| `test` | `api_follows_artists_integration_test.go` | Follow system and artist profiles |
| `test` | `api_admin_integration_test.go` | Admin routes, DB role resolution |
