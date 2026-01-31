# Design Document: LocalStack Integration Test Migration

## Overview

Add integration tests that exercise the full backend stack — from HTTP request through Echo router, auth middleware, handlers, services, and repositories — against LocalStack-emulated AWS services. This is additive to existing mock-based unit tests. The centerpiece is a `SetupTestServer(t)` helper that spins up the real Echo app backed by LocalStack, enabling tests to make actual HTTP requests against `/api/v1/*` routes.

## Code Reuse Analysis

### Existing Components to Leverage

- **`backend/internal/testutil/`** (localstack.go, fixtures.go, cleanup.go): All LocalStack client setup, fixture creation, and cleanup. Will be extended with new fixture helpers and the test server helper.
- **`backend/cmd/api/main.go` → `setupEcho()`**: The existing server initialization function already supports LocalStack via the `AWS_ENDPOINT` env var. The test server will reuse this directly.
- **Auth middleware header fallback**: `middleware/auth.go` already falls back to `X-User-ID` and `X-User-Role` headers when API Gateway JWT claims are unavailable — perfect for test HTTP requests without real JWT tokens.
- **`docker/docker-compose.yml`**: Existing LocalStack container with DynamoDB, S3, Cognito, STS, IAM.
- **`docker/localstack-init/`**: Existing init scripts for DynamoDB table, S3 bucket, Cognito user pool/groups/users.

### Integration Points

- **DynamoDB (LocalStack)**: Single-table `MusicLibrary` with PK/SK and GSI1. All repository operations verified against real DynamoDB behavior.
- **S3 (LocalStack)**: `music-library-local-media` bucket for upload/download/delete operations.
- **Cognito (LocalStack)**: User pool with admin/subscriber/artist groups. Used for `CognitoClient` tests and optionally for real JWT token tests.

## Architecture

### Test Execution Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Integration Test                                  │
│                                                                      │
│  1. testutil.SetupTestServer(t)                                      │
│     ├── Check LocalStack health → t.Skip() if down                   │
│     ├── Set AWS_ENDPOINT=http://localhost:4566                        │
│     ├── Call setupEcho() → real Echo app with LocalStack clients      │
│     ├── Wrap in httptest.NewServer → get base URL                     │
│     └── Return TestServerContext (server + TestContext + cleanup)     │
│                                                                      │
│  2. Create test data via TestContext fixtures                         │
│     ├── ctx.CreateTestUser(t, email, role)                           │
│     ├── ctx.CreateTestTrack(t, userID, opts...)                      │
│     └── ctx.CreateTestArtistProfile(t, userID, opts...)  [NEW]       │
│                                                                      │
│  3. Make HTTP requests to test server                                │
│     ├── GET  {baseURL}/api/v1/tracks                                 │
│     ├── POST {baseURL}/api/v1/playlists                              │
│     ├── Headers: X-User-ID, X-User-Role                             │
│     └── Assert: status code, response body, side effects in DB       │
│                                                                      │
│  4. Verify side effects directly via TestContext                     │
│     ├── ctx.ItemExists(t, pk, sk)                                    │
│     ├── ctx.GetItem(t, pk, sk)                                       │
│     └── Direct DynamoDB/S3 assertions                                │
│                                                                      │
│  5. defer cleanup()                                                  │
│     ├── Stop httptest.Server                                          │
│     └── Clean up all registered DynamoDB/S3 items                    │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                    testutil package                            │
│                                                                │
│  ┌─────────────────────┐  ┌─────────────────────────────┐    │
│  │   localstack.go     │  │   server.go [NEW]            │    │
│  │   (existing)        │  │                               │    │
│  │   SetupLocalStack() │  │   SetupTestServer(t)          │    │
│  │   IsLocalStackRunning│ │   → calls setupEcho()         │    │
│  │   TestContext        │  │   → wraps in httptest.Server  │    │
│  └─────────────────────┘  │   → returns TestServerContext  │    │
│                            └─────────────────────────────────┘    │
│  ┌─────────────────────┐  ┌─────────────────────────────┐    │
│  │   fixtures.go       │  │   http_helpers.go [NEW]      │    │
│  │   (existing+extend) │  │                               │    │
│  │   CreateTestTrack   │  │   DoRequest(method,path,body) │    │
│  │   CreateTestUser    │  │   AsUser(userID, role)        │    │
│  │   CreateTestPlaylist│  │   AssertStatus(t, resp, code) │    │
│  │   [NEW] CreateTest  │  │   DecodeJSON(t, resp, &out)   │    │
│  │     ArtistProfile   │  └─────────────────────────────┘    │
│  │   [NEW] CreateTest  │                                      │
│  │     Follow          │  ┌─────────────────────────────┐    │
│  │   [NEW] CreateTest  │  │   cleanup.go (existing)      │    │
│  │     Tag             │  │   + S3 object cleanup [NEW]  │    │
│  │   [NEW] CreateTest  │  └─────────────────────────────┘    │
│  │     Album           │                                      │
│  │   [NEW] CreateTest  │                                      │
│  │     S3Object        │                                      │
│  └─────────────────────┘                                      │
└──────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Component 1: TestServerContext (`testutil/server.go`) [NEW]

- **Purpose:** Spin up the full Echo HTTP server backed by LocalStack for end-to-end API testing.
- **Interfaces:**
  ```go
  type TestServerContext struct {
      *TestContext                    // Embedded — DynamoDB, S3, Cognito clients
      Server      *httptest.Server   // Running HTTP test server
      BaseURL     string             // e.g. "http://127.0.0.1:54321"
      Echo        *echo.Echo         // Echo instance for direct access if needed
  }

  // SetupTestServer creates a full Echo app backed by LocalStack.
  // Returns TestServerContext and cleanup function.
  // Skips if LocalStack is not running.
  func SetupTestServer(t *testing.T) (*TestServerContext, func())
  ```
- **Dependencies:** `cmd/api.setupEcho()`, `testutil.SetupLocalStack()`, `net/http/httptest`
- **Reuses:** Existing `setupEcho()` which already handles LocalStack endpoint override via `AWS_ENDPOINT` env var.

### Component 2: HTTP Test Helpers (`testutil/http_helpers.go`) [NEW]

- **Purpose:** Reduce boilerplate in HTTP integration tests — authenticated requests, JSON encoding/decoding, assertions.
- **Interfaces:**
  ```go
  // RequestOption configures an HTTP request
  type RequestOption func(*http.Request)

  // AsUser sets X-User-ID and X-User-Role headers
  func AsUser(userID string, role models.UserRole) RequestOption

  // WithJSON sets Content-Type and marshals body
  func WithJSON(body interface{}) RequestOption

  // DoRequest makes an HTTP request to the test server
  func (tsc *TestServerContext) DoRequest(t *testing.T, method, path string, opts ...RequestOption) *http.Response

  // DecodeJSON reads and unmarshals response body
  func DecodeJSON[T any](t *testing.T, resp *http.Response) T

  // AssertStatus verifies response status code with helpful error on failure
  func AssertStatus(t *testing.T, resp *http.Response, expected int)
  ```
- **Dependencies:** `net/http`, `encoding/json`, `testing`
- **Reuses:** Auth middleware's header fallback (`X-User-ID`, `X-User-Role`)

### Component 3: Extended Fixtures (`testutil/fixtures.go`) [EXTEND]

- **Purpose:** Add fixture creation helpers for all entity types needed by integration tests.
- **New Interfaces:**
  ```go
  // Artist profiles
  type ArtistProfileOption func(map[string]dynamodbtypes.AttributeValue)
  func WithArtistName(name string) ArtistProfileOption
  func WithArtistBio(bio string) ArtistProfileOption
  func (tc *TestContext) CreateTestArtistProfile(t *testing.T, userID string, opts ...ArtistProfileOption) string

  // Follow relationships
  func (tc *TestContext) CreateTestFollow(t *testing.T, followerID, followedID string) string

  // Tags
  type TagOption func(map[string]dynamodbtypes.AttributeValue)
  func WithTagColor(color string) TagOption
  func (tc *TestContext) CreateTestTag(t *testing.T, userID, tagName string, opts ...TagOption) string

  // Albums
  type AlbumOption func(map[string]dynamodbtypes.AttributeValue)
  func WithAlbumArtist(artist string) AlbumOption
  func (tc *TestContext) CreateTestAlbum(t *testing.T, userID string, opts ...AlbumOption) string

  // S3 objects
  func (tc *TestContext) CreateTestS3Object(t *testing.T, key string, content []byte) string
  ```
- **Dependencies:** Existing `TestContext` struct, DynamoDB/S3 clients
- **Reuses:** Same pattern as existing `CreateTestTrack`, `CreateTestUser`, `CreateTestPlaylist`

### Component 4: Repository Integration Tests [NEW]

- **Purpose:** Test repository methods against real DynamoDB, verifying GSI queries, pagination, and single-table key patterns.
- **Files:**
  - `backend/internal/repository/dynamodb_integration_test.go` — Track, User, Playlist, Tag CRUD
  - `backend/internal/repository/artist_profile_integration_test.go` — Artist Profile with GSI1
  - `backend/internal/repository/follow_integration_test.go` — Follow with GSI1
  - `backend/internal/repository/s3_integration_test.go` — S3 upload, download, delete, delete-by-prefix
- **Pattern:**
  ```go
  //go:build integration

  func TestIntegration_Repository_CreateAndGetTrack(t *testing.T) {
      tc, cleanup := testutil.SetupLocalStack(t)
      defer cleanup()

      repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
      // Test using real repo methods against LocalStack DynamoDB
  }
  ```
- **Reuses:** `testutil.SetupLocalStack()`, `testutil.TestContext`

### Component 5: Service Integration Tests [NEW]

- **Purpose:** Test service layer with real repositories against LocalStack, verifying business logic like visibility enforcement, role checks, and cross-entity operations.
- **Files:**
  - `backend/internal/service/track_integration_test.go` — extend existing file with service-level tests
  - `backend/internal/service/playlist_integration_test.go`
  - `backend/internal/service/tag_integration_test.go`
  - `backend/internal/service/follow_integration_test.go`
  - `backend/internal/service/artist_profile_integration_test.go`
  - `backend/internal/service/role_integration_test.go`
  - `backend/internal/service/user_integration_test.go`
- **Pattern:**
  ```go
  //go:build integration

  func TestIntegration_TrackService_VisibilityEnforcement(t *testing.T) {
      tc, cleanup := testutil.SetupLocalStack(t)
      defer cleanup()

      repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
      s3Repo := repository.NewS3Repository(tc.S3, s3.NewPresignClient(tc.S3), tc.BucketName)
      trackSvc := service.NewTrackService(repo, s3Repo)

      // Create user + private track, test visibility as different users
  }
  ```
- **Reuses:** `testutil.SetupLocalStack()`, real repository/service constructors

### Component 6: API Integration Tests [NEW]

- **Purpose:** Test the full HTTP flow through the Echo server — routes, middleware, handlers, services, repositories — all backed by LocalStack.
- **Files:**
  - `backend/test/api_tracks_integration_test.go` — Track CRUD, visibility, admin delete
  - `backend/test/api_playlists_integration_test.go` — Playlist CRUD, visibility, public listing
  - `backend/test/api_tags_integration_test.go` — Tag CRUD, track-tag associations
  - `backend/test/api_follows_integration_test.go` — Follow/unfollow, counts
  - `backend/test/api_artists_integration_test.go` — Artist profile CRUD
  - `backend/test/api_admin_integration_test.go` — Admin routes, role enforcement
  - `backend/test/api_auth_integration_test.go` — Auth requirements, role-based access
- **Pattern:**
  ```go
  //go:build integration

  func TestAPI_Tracks_CRUD(t *testing.T) {
      tsc, cleanup := testutil.SetupTestServer(t)
      defer cleanup()

      // Create user in DB
      userID := tsc.CreateTestUser(t, "api-test@example.com", "subscriber")

      // POST /api/v1/tracks — create track via HTTP
      resp := tsc.DoRequest(t, "POST", "/api/v1/tracks",
          AsUser(userID, models.RoleSubscriber),
          WithJSON(map[string]string{"title": "Test Track"}),
      )
      AssertStatus(t, resp, http.StatusCreated)
      created := DecodeJSON[models.TrackResponse](t, resp)

      // GET /api/v1/tracks/:id — read back via HTTP
      resp = tsc.DoRequest(t, "GET", "/api/v1/tracks/"+created.ID,
          AsUser(userID, models.RoleSubscriber),
      )
      AssertStatus(t, resp, http.StatusOK)

      // Verify in DynamoDB directly
      assert.True(t, tsc.ItemExists(t, "USER#"+userID, "TRACK#"+created.ID))
  }

  func TestAPI_Tracks_VisibilityEnforcement(t *testing.T) {
      tsc, cleanup := testutil.SetupTestServer(t)
      defer cleanup()

      ownerID := tsc.CreateTestUser(t, "owner@test.com", "subscriber")
      otherID := tsc.CreateTestUser(t, "other@test.com", "subscriber")
      adminID := tsc.CreateTestUser(t, "admin@test.com", "admin")

      // Create private track
      trackID := tsc.CreateTestTrack(t, ownerID, testutil.WithTrackVisibility("private"))

      // Owner can access
      resp := tsc.DoRequest(t, "GET", "/api/v1/tracks/"+trackID,
          AsUser(ownerID, models.RoleSubscriber))
      AssertStatus(t, resp, http.StatusOK)

      // Other user gets 403
      resp = tsc.DoRequest(t, "GET", "/api/v1/tracks/"+trackID,
          AsUser(otherID, models.RoleSubscriber))
      AssertStatus(t, resp, http.StatusForbidden)

      // Admin gets 200 (hasGlobal)
      resp = tsc.DoRequest(t, "GET", "/api/v1/tracks/"+trackID,
          AsUser(adminID, models.RoleAdmin))
      AssertStatus(t, resp, http.StatusOK)
  }
  ```
- **Reuses:** `testutil.SetupTestServer()`, `testutil.DoRequest()`, auth middleware header fallback

### Component 7: Cognito Integration Tests [NEW]

- **Purpose:** Test CognitoClient operations (admin user management) against LocalStack Cognito.
- **File:** `backend/internal/service/cognito_integration_test.go`
- **Pattern:**
  ```go
  //go:build integration

  func TestIntegration_CognitoClient_ListGroupsForUser(t *testing.T) {
      tc, cleanup := testutil.SetupLocalStack(t)
      defer cleanup()
      if tc.UserPoolID == "" {
          t.Skip("Cognito not configured")
      }

      cognitoClient := service.NewCognitoClient(tc.Cognito, tc.UserPoolID)
      groups, err := cognitoClient.ListGroupsForUser(ctx, "admin@local.test")
      require.NoError(t, err)
      assert.Contains(t, groups, "admin")
  }
  ```
- **Reuses:** `testutil.SetupLocalStack()`, pre-created Cognito users from `init-cognito.sh`

### Component 8: CI GitHub Action [NEW/EXTEND]

- **Purpose:** Run integration tests in CI with LocalStack via Docker Compose.
- **File:** `.github/workflows/ci.yml` (extend existing) or `.github/workflows/integration.yml` (new)
- **Additions:**
  ```yaml
  integration-tests:
    runs-on: ubuntu-latest
    services:
      localstack:
        image: localstack/localstack:3.4
        ports:
          - 4566:4566
        env:
          SERVICES: dynamodb,s3,sts,iam,cognito-idp
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - name: Wait for LocalStack
        run: scripts/wait-for-localstack.sh
      - name: Initialize LocalStack
        run: |
          docker/localstack-init/init-aws.sh
          docker/localstack-init/init-cognito.sh
      - name: Run integration tests
        env:
          AWS_ENDPOINT: http://localhost:4566
          DYNAMODB_TABLE_NAME: MusicLibrary
          MEDIA_BUCKET: music-library-local-media
        run: cd backend && go test -tags=integration -v ./...
  ```

## Data Models

No new domain models. Tests reuse existing models from `backend/internal/models/`. Test fixtures create DynamoDB items matching the existing single-table schema:

| Entity | PK | SK | Used In |
|--------|----|----|---------|
| User | `USER#{userId}` | `PROFILE` | All tests |
| Track | `USER#{userId}` | `TRACK#{trackId}` | Track/API tests |
| Playlist | `USER#{userId}` | `PLAYLIST#{playlistId}` | Playlist tests |
| Tag | `USER#{userId}` | `TAG#{tagName}` | Tag tests |
| ArtistProfile | `ARTIST_PROFILE#{id}` | `PROFILE` | Artist/Follow tests |
| Follow | `FOLLOW#{followerId}` | `FOLLOWING#{followedId}` | Follow tests |
| Album | `USER#{userId}` | `ALBUM#{albumId}` | Album tests |

## Error Handling

### Error Scenarios

1. **LocalStack not running**
   - **Handling:** `t.Skip()` with message pointing to Docker Compose command
   - **User Impact:** Tests skip gracefully, CI retries

2. **Cognito not initialized**
   - **Handling:** `t.Skip()` for Cognito-dependent tests; other tests continue
   - **User Impact:** Partial test coverage still runs

3. **Port conflict on httptest.Server**
   - **Handling:** `httptest.NewServer` uses random ports by default — no conflict
   - **User Impact:** None

4. **Test data leaking between tests**
   - **Handling:** Each test creates data with unique UUIDs; `defer cleanup()` removes all registered items
   - **User Impact:** None — tests are isolated

5. **Cleanup failure**
   - **Handling:** `t.Logf()` warning, no test failure. `CleanupAll` available as nuclear option.
   - **User Impact:** Stale data in LocalStack, resolved by container restart

## Testing Strategy

### Test Layers

```
┌─────────────────────────────────────────────────────────┐
│  Layer 3: API Integration Tests (NEW)                    │
│  backend/test/api_*_integration_test.go                  │
│  Full HTTP flow: request → router → middleware →          │
│  handler → service → repo → DynamoDB/S3                  │
│  Uses: SetupTestServer + DoRequest + X-User-ID headers   │
├─────────────────────────────────────────────────────────┤
│  Layer 2: Service Integration Tests (NEW)                │
│  backend/internal/service/*_integration_test.go          │
│  Real services + real repos against LocalStack           │
│  Uses: SetupLocalStack + real constructors               │
├─────────────────────────────────────────────────────────┤
│  Layer 1: Repository Integration Tests (NEW)             │
│  backend/internal/repository/*_integration_test.go       │
│  Real DynamoDB/S3 operations against LocalStack          │
│  Uses: SetupLocalStack + real repo constructors          │
├─────────────────────────────────────────────────────────┤
│  Layer 0: Unit Tests (EXISTING - unchanged)              │
│  backend/internal/service/*_test.go                      │
│  Mock-based tests for business logic                     │
│  Uses: testify/mock                                      │
└─────────────────────────────────────────────────────────┘
```

### Key Test Scenarios (API Layer)

| Scenario | Method | Path | Auth | Expected |
|----------|--------|------|------|----------|
| Health check | GET | `/health` | None | 200 |
| Unauthenticated access | GET | `/api/v1/tracks` | None | 401 |
| List own tracks | GET | `/api/v1/tracks` | subscriber | 200 + own tracks |
| Admin lists all tracks | GET | `/api/v1/tracks` | admin | 200 + all tracks |
| Get private track (owner) | GET | `/api/v1/tracks/:id` | owner | 200 |
| Get private track (other) | GET | `/api/v1/tracks/:id` | other | 403 |
| Get private track (admin) | GET | `/api/v1/tracks/:id` | admin | 200 |
| Admin deletes any track | DELETE | `/api/v1/tracks/:id` | admin | 200 |
| Create playlist | POST | `/api/v1/playlists` | subscriber | 201 |
| Public playlist discovery | GET | `/playlists/public` | any auth | 200 |
| Create tag | POST | `/api/v1/tags` | subscriber | 201 |
| Follow artist | POST | `/api/v1/artists/entity/:id/follow` | subscriber | 200 |
| Unfollow artist | DELETE | `/api/v1/artists/entity/:id/follow` | subscriber | 200 |
| Admin search users | GET | `/api/v1/admin/users` | admin | 200 |
| Non-admin hits admin route | GET | `/api/v1/admin/users` | subscriber | 403 |

### Running Tests

```bash
# Unit tests only (existing, fast)
cd backend && go test ./...

# Integration tests only (requires LocalStack)
cd backend && go test -tags=integration ./...

# All tests
cd backend && go test -tags=integration ./... && go test ./...

# Specific integration test
cd backend && go test -tags=integration -run TestAPI_Tracks -v ./test/
```
