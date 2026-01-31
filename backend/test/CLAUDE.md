# Backend Test - CLAUDE.md

## Overview

Full-stack API integration tests that exercise HTTP endpoints through a real Echo server backed by LocalStack. These tests verify the complete request flow: HTTP request -> middleware (auth, role checking) -> handler -> service -> repository -> DynamoDB/S3.

## Directory Structure

```
backend/test/
├── integration_test.go                    # Original integration test (track upload flow)
├── api_auth_integration_test.go           # Auth middleware and health endpoint tests
├── api_tracks_integration_test.go         # Track CRUD and visibility via HTTP
├── api_playlists_tags_integration_test.go # Playlist and tag endpoints
├── api_follows_artists_integration_test.go # Follow system and artist profiles
├── api_admin_integration_test.go          # Admin routes with DB role resolution
└── CLAUDE.md                              # This file
```

## File Descriptions

| File | Test Functions | What It Tests |
|------|---------------|---------------|
| `integration_test.go` | `TestEpic1_*`, `TestEpic2_*`, `TestEpic3_*`, `TestEpic4_*`, `TestAuthenticated_*` | **Deployed-environment tests** — hit real AWS (API Gateway, DynamoDB, Cognito, CloudFront). Skipped when `LOCALSTACK_ENDPOINT` is set. |
| `api_auth_integration_test.go` | `TestIntegration_API_Health`, `TestIntegration_API_AuthMiddleware` | Health endpoint (200 without auth), 401 without headers, role-based 403, admin route access |
| `api_tracks_integration_test.go` | `TestIntegration_API_TracksCRUD`, `TestIntegration_API_TrackVisibility`, `TestIntegration_API_AdminTrackDelete` | Full CRUD lifecycle, visibility enforcement (owner/other/admin), admin cross-user delete |
| `api_playlists_tags_integration_test.go` | `TestIntegration_API_PlaylistCRUD`, `TestIntegration_API_TagsCRUD` | Playlist create/add tracks/visibility/public discovery/delete, tag create/add to track/list/remove/delete |
| `api_follows_artists_integration_test.go` | `TestIntegration_API_ArtistProfileCRUD`, `TestIntegration_API_FollowSystem` | Artist profile create/get/list/update, follow/unfollow/is-following/followers list |
| `api_admin_integration_test.go` | `TestIntegration_API_AdminRoutes` | Non-admin 403, admin user search, role update, DB role check overrides header role |

## Build Tag

All files use `//go:build integration` and are excluded from normal `go test ./...` runs.

## Package

```go
package integration
```

Uses an external test package to avoid import cycles with `testutil`.

## Key Test Patterns

### Test Server Setup
```go
tsc, cleanup := testutil.SetupTestServer(t)
defer cleanup()
```
Creates a full Echo HTTP server connected to LocalStack (DynamoDB + S3 + Cognito).

### Making Authenticated Requests
```go
resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/tracks",
    testutil.AsUser(userID, models.RoleSubscriber),
)
testutil.AssertStatus(t, resp, http.StatusOK)
body := testutil.DecodeJSONBody(t, resp)
```

### Creating Test Data
```go
userID := tsc.CreateTestUser(t, "email@test.com", "subscriber")
trackID := tsc.CreateTestTrack(t, userID, testutil.WithTrackTitle("My Track"))
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `testutil` | `SetupTestServer`, `AsUser`, `WithJSON`, `DoRequest`, `AssertStatus`, `DecodeJSONBody` |
| `models` | `RoleSubscriber`, `RoleAdmin`, `RoleArtist` for role constants |
| `testify` | `assert`, `require` for test assertions |

## Running Tests

```bash
# Run all API integration tests (requires LocalStack on port 4566)
cd backend && go test -tags=integration -v ./test/

# Run a specific test
cd backend && go test -tags=integration -v -run TestIntegration_API_TrackVisibility ./test/
```

## Note on .gitignore

The `backend/test/` directory is listed in `.gitignore`. Files must be force-added with `git add -f backend/test/`.
