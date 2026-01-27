# Test Utilities (testutil) - CLAUDE.md

## Overview

Integration test utilities for testing against LocalStack. Provides helpers for setting up LocalStack connections, creating test fixtures, and cleaning up test data.

## Directory Structure

```
testutil/
├── localstack.go    # LocalStack connection and setup
├── fixtures.go      # Test data creation helpers
├── cleanup.go       # Test data cleanup
└── CLAUDE.md        # This file
```

## File Descriptions

| File | Purpose |
|------|---------|
| `localstack.go` | LocalStack detection, client creation, and TestContext setup |
| `fixtures.go` | Test user definitions and track/user creation helpers |
| `cleanup.go` | Cleanup functions for test data |

## Core Types

### TestContext

Central struct providing LocalStack clients and test helpers:

```go
type TestContext struct {
    DynamoDB   *dynamodb.Client
    S3         *s3.Client
    Cognito    *cognitoidentityprovider.Client
    TableName  string
    BucketName string
    UserPoolID string
    ClientID   string
}
```

### TestUser

Pre-configured test user credentials:

```go
type TestUser struct {
    Email    string
    Password string
    Role     string
}
```

## Key Functions

### localstack.go

| Function | Signature | Purpose |
|----------|-----------|---------|
| `SetupLocalStack` | `(t *testing.T) (*TestContext, func())` | Creates TestContext with LocalStack clients; returns cleanup function |
| `IsLocalStackRunning` | `() bool` | Checks if LocalStack is healthy |
| `NewLocalStackConfig` | `(ctx context.Context) (aws.Config, error)` | Creates AWS config for LocalStack |

### fixtures.go

| Function | Signature | Purpose |
|----------|-----------|---------|
| `CreateTestTrack` | `(t *testing.T, userID string, opts ...TrackOption) string` | Creates a test track, returns track ID |
| `CreateTestUser` | `(t *testing.T, email, role string) string` | Creates a test user, returns user ID |
| `GetTestUserToken` | `(t *testing.T, role string) string` | Gets JWT token for test user |

### cleanup.go

| Function | Signature | Purpose |
|----------|-----------|---------|
| `CleanupUser` | `(t *testing.T, userID string)` | Deletes user and related data |
| `CleanupTrack` | `(t *testing.T, userID, trackID string)` | Deletes specific track |
| `CleanupAll` | `(t *testing.T)` | Clears all test data from table |

## Test Users

Pre-defined test users matching LocalStack Cognito init:

| Key | Email | Password | Role |
|-----|-------|----------|------|
| `admin` | `admin@local.test` | `LocalTest123!` | admin |
| `subscriber` | `subscriber@local.test` | `LocalTest123!` | subscriber |
| `artist` | `artist@local.test` | `LocalTest123!` | artist |

## Usage Examples

### Basic Integration Test

```go
//go:build integration

package mypackage_test

import (
    "testing"
    "github.com/gvasels/personal-music-searchengine/backend/internal/testutil"
)

func TestIntegration_Something(t *testing.T) {
    // Setup - connects to LocalStack or skips if unavailable
    ctx, cleanup := testutil.SetupLocalStack(t)
    defer cleanup()

    // Create test data
    userID := ctx.CreateTestUser(t, "test@example.com", "subscriber")
    trackID := ctx.CreateTestTrack(t, userID)

    // Your test assertions...
}
```

### Testing with Authentication

```go
func TestIntegration_AuthenticatedEndpoint(t *testing.T) {
    ctx, cleanup := testutil.SetupLocalStack(t)
    defer cleanup()

    // Get JWT token for admin user
    token := ctx.GetTestUserToken(t, "admin")

    // Use token in API call
    req, _ := http.NewRequest("GET", "/api/tracks", nil)
    req.Header.Set("Authorization", "Bearer " + token)

    // Make request and assert...
}
```

### Custom Track Creation

```go
func TestIntegration_CustomTrack(t *testing.T) {
    ctx, cleanup := testutil.SetupLocalStack(t)
    defer cleanup()

    trackID := ctx.CreateTestTrack(t, "user123",
        testutil.WithTitle("My Custom Track"),
        testutil.WithArtist("Test Artist"),
        testutil.WithDuration(180),
    )

    // Assertions...
}
```

### Cleanup After Test

```go
func TestIntegration_WithCleanup(t *testing.T) {
    ctx, cleanup := testutil.SetupLocalStack(t)
    defer cleanup()

    userID := ctx.CreateTestUser(t, "test@example.com", "subscriber")
    trackID := ctx.CreateTestTrack(t, userID)

    // Explicit cleanup (optional - cleanup() handles this)
    defer ctx.CleanupTrack(t, userID, trackID)
    defer ctx.CleanupUser(t, userID)

    // Test logic...
}
```

## Build Tags

Integration tests must use the `integration` build tag:

```go
//go:build integration

package mypackage_test
```

Run with:
```bash
go test -tags=integration ./...
```

Tests will skip gracefully if LocalStack is not running.

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `LOCALSTACK_ENDPOINT` | `http://localhost:4566` | LocalStack endpoint URL |
| `DYNAMODB_TABLE_NAME` | `MusicLibrary` | DynamoDB table name |
| `MEDIA_BUCKET` | `music-library-local-media` | S3 bucket name |

## Error Handling

- `SetupLocalStack` calls `t.Skip()` if LocalStack is unavailable
- All helper functions call `t.Fatal()` on errors
- Cleanup function logs warnings but doesn't fail tests

## Dependencies

```go
import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)
```

## Related Files

- `docker/docker-compose.yml` - LocalStack container definition
- `docker/localstack-init/init-aws.sh` - DynamoDB/S3 initialization
- `docker/localstack-init/init-cognito.sh` - Cognito initialization
- `scripts/wait-for-localstack.sh` - Health check script
