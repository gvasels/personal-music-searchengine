# Requirements: LocalStack Integration Test Migration

## Introduction

Migrate all backend service tests that currently mock AWS services (DynamoDB, S3, Cognito) to also have LocalStack-based integration test counterparts. The project has mature LocalStack test utilities (`backend/internal/testutil/`) and Docker infrastructure (`docker/docker-compose.yml`) but only 1 of 20+ service test files uses them. This spec covers creating integration tests that exercise real AWS API contracts via LocalStack, catching bugs that interface mocks cannot.

**Scope**: Backend Go service and repository integration tests only. Frontend tests and model tests (pure unit tests with no AWS interaction) are out of scope.

## Alignment with Product Vision

Integration tests against LocalStack validate real DynamoDB query behavior, S3 operations, and Cognito authentication flows — catching contract mismatches, pagination bugs, and GSI query issues that mock-based tests miss. This directly supports platform reliability as the multi-user role-based system grows in complexity.

## Requirements

### Requirement 1: Repository Integration Tests

**User Story:** As a developer, I want integration tests for the DynamoDB repository layer, so that I can verify actual DynamoDB operations (queries, scans, GSI lookups, pagination) work correctly against the single-table schema.

#### Acceptance Criteria

1. WHEN a repository integration test runs against LocalStack THEN the test SHALL use `testutil.SetupLocalStack(t)` for setup and cleanup
2. WHEN LocalStack is not running THEN integration tests SHALL skip gracefully via `t.Skip()`
3. WHEN repository tests create DynamoDB items THEN they SHALL use the `testutil.TestContext` fixtures and register items for cleanup
4. IF the repository method uses a GSI (GSI1PK/GSI1SK) THEN the integration test SHALL verify the GSI query returns correct results
5. WHEN repository pagination is tested THEN the test SHALL verify `LastEvaluatedKey` cursor behavior with real DynamoDB pagination

**In-scope repository operations:**
- Track CRUD (PK=`USER#{userId}`, SK=`TRACK#{trackId}`)
- User CRUD (PK=`USER#{userId}`, SK=`PROFILE`)
- Playlist CRUD (PK=`USER#{userId}`, SK=`PLAYLIST#{playlistId}`)
- Tag CRUD (PK=`USER#{userId}`, SK=`TAG#{tagName}`)
- Artist Profile (PK=`ARTIST_PROFILE#{id}`, SK=`PROFILE`, GSI1)
- Follow relationships (PK=`FOLLOW#{followerId}`, SK=`FOLLOWING#{followedId}`, GSI1)
- S3 operations (upload, get presigned URL, delete, delete by prefix)

### Requirement 2: Service Layer Integration Tests

**User Story:** As a developer, I want integration tests for the service layer using real repository instances against LocalStack, so that I can validate end-to-end business logic including visibility enforcement, role checks, and cross-entity operations.

#### Acceptance Criteria

1. WHEN a service integration test runs THEN it SHALL instantiate real `DynamoDBRepository` and real service instances (not mocks)
2. WHEN testing track visibility enforcement THEN the test SHALL verify that private tracks return 403 for non-owners and public tracks are accessible to all
3. WHEN testing admin operations (hasGlobal=true) THEN the test SHALL verify admins can access all tracks regardless of ownership
4. WHEN testing playlist visibility THEN the test SHALL verify public playlists appear in GSI2 queries and private playlists do not
5. WHEN testing follow/unfollow THEN the test SHALL verify follower counts increment/decrement correctly in DynamoDB

**In-scope services:**
- TrackService (visibility enforcement, admin global access, CRUD)
- PlaylistService (visibility, public playlist discovery, track management)
- TagService (tag CRUD, track-tag associations, case-insensitive lookup)
- FollowService (follow/unfollow, count updates, duplicate prevention)
- ArtistProfileService (CRUD, user-to-profile linking)
- RoleService (role get/update via DynamoDB)
- UserService (create-if-not-exists, profile updates)

### Requirement 3: S3 Integration Tests

**User Story:** As a developer, I want integration tests for S3 operations, so that I can verify file uploads, presigned URLs, and deletion (including HLS prefix cleanup) work correctly.

#### Acceptance Criteria

1. WHEN an S3 integration test uploads a file THEN it SHALL verify the object exists in the LocalStack S3 bucket
2. WHEN testing presigned URL generation THEN the test SHALL verify the URL is accessible and returns the uploaded content
3. WHEN testing delete-by-prefix (HLS cleanup) THEN the test SHALL create multiple objects with a shared prefix and verify all are deleted
4. WHEN S3 operations fail THEN the test SHALL verify appropriate error types are returned

### Requirement 4: Cognito Integration Tests

**User Story:** As a developer, I want integration tests for Cognito operations, so that I can verify authentication flows and group management against LocalStack Cognito.

#### Acceptance Criteria

1. WHEN testing authentication THEN the test SHALL authenticate pre-created test users (admin, subscriber, artist) and verify JWT tokens are returned
2. WHEN testing group operations THEN the test SHALL verify adding/removing users from Cognito groups works correctly
3. WHEN testing the admin CognitoClient THEN the test SHALL verify user search, disable, and group listing operations
4. IF Cognito is not initialized in LocalStack THEN the test SHALL skip gracefully

### Requirement 5: Testutil Enhancements

**User Story:** As a developer, I want the testutil package to support all entity types needed for integration tests, so that test setup is consistent and cleanup is automatic.

#### Acceptance Criteria

1. WHEN creating test fixtures THEN `testutil` SHALL provide helpers for: artist profiles, follows, tags, albums, and S3 objects
2. WHEN integration tests complete THEN cleanup SHALL remove all created entities (DynamoDB items + S3 objects)
3. WHEN creating fixtures THEN each helper SHALL return the entity ID and accept functional options for customization
4. IF a fixture depends on another entity (e.g., follow requires artist profile) THEN the helper SHALL document this dependency

### Requirement 6: Full HTTP API Integration Tests (End-to-End)

**User Story:** As a developer, I want integration tests that spin up the real Echo HTTP server backed by LocalStack, so that I can verify the complete request lifecycle: HTTP request → Echo router → auth middleware → handler → service → repository → DynamoDB/S3/Cognito.

#### Acceptance Criteria

1. WHEN an API integration test starts THEN it SHALL call `setupEcho()` with `AWS_ENDPOINT` pointing at LocalStack and wrap it in `httptest.NewServer` to get a test HTTP server
2. WHEN testing authenticated endpoints THEN the test SHALL use `X-User-ID` and `X-User-Role` headers (the existing auth middleware fallback for non-API-Gateway environments)
3. WHEN testing CRUD flows THEN the test SHALL verify the full cycle through HTTP: create via POST → verify via GET → update via PUT → list via GET → delete via DELETE → verify 404
4. WHEN testing track visibility enforcement end-to-end THEN the test SHALL: create a private track as user A via HTTP, attempt GET as user B via HTTP and verify 403, attempt GET as admin via HTTP and verify 200
5. WHEN testing admin routes (`/api/v1/admin/*`) THEN the test SHALL verify the `RequireRoleWithDBCheck` middleware queries the real DynamoDB for role resolution
6. WHEN testing playlist visibility THEN the test SHALL: create a public playlist via PUT visibility, verify it appears in GET `/playlists/public`, create a private playlist and verify it does not appear
7. WHEN testing tag operations THEN the test SHALL verify the full HTTP flow: create tag → add to track → list tracks by tag → remove tag → verify removed
8. WHEN testing follow operations THEN the test SHALL verify HTTP: POST follow → GET followers list → verify count incremented → DELETE unfollow → verify count decremented
9. WHEN testing unauthenticated requests THEN the test SHALL verify protected endpoints return 401 and public endpoints (health, public playlists) return 200

**In-scope API routes:**
- `/health` — health check (public)
- `/api/v1/me` — user profile (auth required)
- `/api/v1/tracks` — track CRUD (auth + visibility)
- `/api/v1/albums` — album listing (auth)
- `/api/v1/playlists` — playlist CRUD + `/playlists/public` (mixed auth)
- `/api/v1/tags` — tag CRUD (auth)
- `/api/v1/search` — search (auth)
- `/api/v1/uploads` — upload management (auth)
- `/api/v1/artists/entity` — artist profiles (mixed auth)
- `/api/v1/artists/entity/:id/follow` — follow system (auth)
- `/api/v1/admin/*` — admin routes (admin role + DB check)

### Requirement 7: Test Server Helper

**User Story:** As a developer, I want a `testutil.SetupTestServer(t)` helper that returns a running HTTP test server backed by LocalStack, so that API integration tests have a one-line setup.

#### Acceptance Criteria

1. WHEN `SetupTestServer(t)` is called THEN it SHALL start LocalStack clients, construct repositories, services, handlers, and return an `*httptest.Server` with the full Echo app
2. WHEN `SetupTestServer(t)` is called THEN it SHALL also return a `*TestContext` for direct DynamoDB/S3 verification alongside HTTP tests
3. WHEN the test completes THEN the cleanup function SHALL stop the HTTP server and clean up all test data
4. IF LocalStack is not running THEN `SetupTestServer(t)` SHALL call `t.Skip()`
5. WHEN constructing the server THEN optional services (Search, Admin) SHALL be wired if LocalStack Cognito is available

### Requirement 8: CI Integration

**User Story:** As a developer, I want integration tests to run in CI via GitHub Actions, so that LocalStack tests validate changes on every PR.

#### Acceptance Criteria

1. WHEN a PR is opened THEN GitHub Actions SHALL start LocalStack via Docker Compose and run `go test -tags=integration ./...`
2. WHEN LocalStack health check fails THEN the CI step SHALL retry up to 3 times before failing
3. WHEN integration tests complete THEN coverage data SHALL be collected alongside unit test coverage

## Non-Functional Requirements

### Code Architecture and Modularity
- Integration test files use `//go:build integration` build tag and `_integration_test.go` suffix
- Tests live alongside the code they test (e.g., `service/track_integration_test.go`)
- All integration tests use `testutil.SetupLocalStack(t)` — no direct AWS client construction in test files
- Existing unit tests with mocks remain unchanged — integration tests are additive

### Performance
- Individual integration tests complete within 5 seconds
- Full integration test suite completes within 60 seconds
- Tests run in parallel where possible (`t.Parallel()`) for non-conflicting data

### Security
- Integration tests only connect to LocalStack (localhost:4566), never real AWS
- Static credentials (`test/test/test`) used for all LocalStack connections
- `CleanupAll` safety check prevents running against non-localhost endpoints

### Reliability
- Tests skip gracefully when LocalStack is unavailable
- Cleanup runs even when tests fail (via `defer cleanup()`)
- Each test creates its own isolated data (unique UUIDs) to avoid cross-test interference
