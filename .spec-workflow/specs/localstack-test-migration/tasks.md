# Tasks: LocalStack Integration Test Migration

## Phase 1: Testutil Infrastructure

- [x] 1.1. Create test server helper (`testutil/server.go`)
  - File: `backend/internal/testutil/server.go`
  - Create `SetupTestServer(t)` that calls `setupEcho()` with `AWS_ENDPOINT` pointing to LocalStack, wraps in `httptest.NewServer`, returns `TestServerContext` with embedded `TestContext` + `Server` + `BaseURL`
  - Must set env vars (`AWS_ENDPOINT`, `DYNAMODB_TABLE_NAME`, `MEDIA_BUCKET`, `AWS_REGION`) before calling `setupEcho()`, restore after
  - Returns cleanup function that stops server and cleans up test data
  - Skips if LocalStack not running
  - Purpose: One-line setup for full HTTP API integration tests
  - _Leverage: `backend/cmd/api/main.go` (`setupEcho()`), `backend/internal/testutil/localstack.go` (`SetupLocalStack`, `IsLocalStackRunning`)_
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer specializing in test infrastructure | Task: Create `testutil/server.go` with `SetupTestServer(t)` function. Read `backend/cmd/api/main.go` to understand how `setupEcho()` works — it already detects `AWS_ENDPOINT` env var and routes all AWS clients to that endpoint. The helper should: (1) check `IsLocalStackRunning()`, skip if not, (2) set env vars for LocalStack, (3) call `setupEcho()` to get a fully-wired Echo instance, (4) wrap in `httptest.NewServer`, (5) also create a `TestContext` for direct DB verification, (6) return `TestServerContext` struct and cleanup func. Use `//go:build integration` tag. | Restrictions: Do not duplicate `setupEcho()` logic — call it directly. Do not modify `cmd/api/main.go`. Env vars must be restored after test. | _Leverage: `backend/cmd/api/main.go`, `backend/internal/testutil/localstack.go` | _Requirements: Req 7 (Test Server Helper) | Success: `SetupTestServer(t)` returns a running HTTP server backed by LocalStack. Tests can make requests to `tsc.BaseURL + "/api/v1/tracks"` and get real responses. Skips gracefully when LocalStack is down. Mark task in-progress in tasks.md, log implementation, mark complete._

- [x] 1.2. Create HTTP test helpers (`testutil/http_helpers.go`)
  - File: `backend/internal/testutil/http_helpers.go`
  - Create `RequestOption` type and helpers: `AsUser(userID, role)` sets `X-User-ID`/`X-User-Role` headers, `WithJSON(body)` marshals and sets Content-Type, `WithHeader(k,v)`
  - Create `DoRequest(t, method, path, opts...)` on `TestServerContext` — builds URL from `BaseURL + path`, applies options, executes request
  - Create `DecodeJSON[T](t, resp)` generic helper for response unmarshaling
  - Create `AssertStatus(t, resp, expected)` with body dump on mismatch for debugging
  - Purpose: Reduce HTTP test boilerplate
  - _Leverage: Auth middleware header fallback in `backend/internal/handlers/middleware/auth.go` (X-User-ID, X-User-Role)_
  - _Requirements: 6.2_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in test utilities | Task: Create `testutil/http_helpers.go` with HTTP test helpers. Read `backend/internal/handlers/middleware/auth.go` to see the header fallback — when no API Gateway JWT is present, it reads `X-User-ID` and `X-User-Role` headers. Build helpers that leverage this: `AsUser(userID string, role models.UserRole) RequestOption` sets those headers. `WithJSON(body interface{}) RequestOption` marshals to JSON and sets Content-Type. `DoRequest` builds full URL, applies all options, returns `*http.Response`. `DecodeJSON` uses generics. `AssertStatus` dumps response body on failure. Use `//go:build integration` tag. | Restrictions: Do not modify auth middleware. Keep helpers stateless. | _Leverage: `backend/internal/handlers/middleware/auth.go` | _Requirements: Req 6 (Full HTTP API Integration Tests) | Success: Tests can write `tsc.DoRequest(t, "GET", "/api/v1/tracks", AsUser(id, models.RoleSubscriber))` and get back `*http.Response`. Mark task in-progress in tasks.md, log implementation, mark complete._

- [x] 1.3. Extend fixtures with new entity helpers (`testutil/fixtures.go`)
  - File: `backend/internal/testutil/fixtures.go` (extend existing)
  - Add `CreateTestArtistProfile(t, userID, opts...)` — creates `ARTIST_PROFILE#{id}` / `PROFILE` item + GSI1 entry (`USER#{userId}` / `ARTIST_PROFILE`)
  - Add `CreateTestFollow(t, followerID, followedID)` — creates `FOLLOW#{followerId}` / `FOLLOWING#{followedId}` item + GSI1 entry
  - Add `CreateTestTag(t, userID, tagName, opts...)` — creates `USER#{userId}` / `TAG#{tagName}` item
  - Add `CreateTestAlbum(t, userID, opts...)` — creates `USER#{userId}` / `ALBUM#{albumId}` item
  - Add `CreateTestS3Object(t, key, content)` — puts object in S3 bucket, registers for cleanup
  - All helpers register items for cleanup and return entity IDs
  - Purpose: Enable setup for all integration test scenarios
  - _Leverage: Existing `CreateTestTrack`, `CreateTestUser`, `CreateTestPlaylist` patterns in `testutil/fixtures.go`. DynamoDB schema from CLAUDE.md._
  - _Requirements: 5.1, 5.2, 5.3, 5.4_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in test infrastructure | Task: Extend `testutil/fixtures.go` with new entity creation helpers. Read the existing file to follow the exact pattern (DynamoDB PutItem with attribute map, functional options, RegisterCleanup). Read `backend/internal/models/artist_profile.go`, `follow.go`, `tag.go` for the model fields. Read CLAUDE.md for the DynamoDB schema (PK/SK/GSI patterns). Each helper creates items matching the single-table schema. `CreateTestS3Object` also needs cleanup (delete object on teardown). | Restrictions: Follow the exact pattern of `CreateTestTrack` — DynamoDB PutItem, uuid for IDs, RegisterCleanup. Do not modify existing helpers. Use `//go:build integration` tag. | _Leverage: `backend/internal/testutil/fixtures.go`, `backend/internal/models/` | _Requirements: Req 5 (Testutil Enhancements) | Success: All entity types can be created in tests with one-line helpers. Each registers for cleanup. Mark task in-progress in tasks.md, log implementation, mark complete._

- [x] 1.4. Add S3 object cleanup to cleanup utilities (`testutil/cleanup.go`)
  - File: `backend/internal/testutil/cleanup.go` (extend existing)
  - Add S3 cleanup item type to `cleanupItem` struct (track S3 keys alongside DynamoDB items)
  - Add `CleanupS3Object(t, key)` method
  - Update `runCleanup` to also delete registered S3 objects
  - Purpose: Ensure S3 test data doesn't accumulate
  - _Leverage: Existing cleanup pattern in `testutil/cleanup.go`, S3 client from `TestContext`_
  - _Requirements: 5.2_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer | Task: Extend `testutil/cleanup.go` to support S3 cleanup. Read the existing file to understand the `cleanupItem` struct and `runCleanup` flow. Add an `itemType` value for S3 objects. In `runCleanup`, when type is "s3", call `tc.S3.DeleteObject` instead of DynamoDB delete. Add `CleanupS3Object(t, key)` public method. Add `RegisterS3Cleanup(key)` method. | Restrictions: Do not break existing DynamoDB cleanup. Keep the safety check for localhost. Use `//go:build integration` tag. | _Leverage: `backend/internal/testutil/cleanup.go` | _Requirements: Req 5 | Success: S3 objects created in tests are automatically cleaned up. Mark task in-progress in tasks.md, log implementation, mark complete._

## Phase 2: Repository Integration Tests

- [x] 2.1. DynamoDB repository integration tests
  - File: `backend/internal/repository/dynamodb_integration_test.go`
  - Test Track CRUD: Create, Get, GetByID, Update, Delete, List with pagination
  - Test User CRUD: CreateIfNotExists, Get, Update, UpdateRole
  - Test Playlist CRUD: Create, Get, Update, Delete, List, AddTracks, RemoveTracks
  - Test Tag CRUD: Create, Get, Update, Delete, List, GetTracksByTag
  - Test pagination: verify `LastEvaluatedKey` cursor behavior with 20+ items
  - Purpose: Verify repository methods work correctly against real DynamoDB
  - _Leverage: `testutil.SetupLocalStack()`, `repository.NewDynamoDBRepository()`_
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer specializing in DynamoDB and integration testing | Task: Create `repository/dynamodb_integration_test.go`. Read `backend/internal/repository/repository.go` for the `Repository` interface, then `dynamodb.go` and `tracks.go` for implementations. For each major entity (Track, User, Playlist, Tag), write integration tests that: (1) setup LocalStack, (2) create a real `DynamoDBRepository`, (3) test CRUD operations, (4) verify data in DynamoDB via `TestContext.GetItem`. Test pagination by creating 20+ items and verifying cursor works. Use `//go:build integration` tag and `t.Parallel()` for independent subtests. | Restrictions: Do not mock anything — use real LocalStack. Do not modify repository code. Each test must clean up via `defer cleanup()`. | _Leverage: `testutil.SetupLocalStack()`, `repository.NewDynamoDBRepository()`, `repository/repository.go` | _Requirements: Req 1 (Repository Integration Tests) | Success: All major repository CRUD operations verified against real DynamoDB. Pagination produces correct pages. Mark task in-progress in tasks.md, log implementation, mark complete._

- [x] 2.2. Artist profile and follow repository integration tests
  - File: `backend/internal/repository/artist_profile_integration_test.go`
  - Test ArtistProfile CRUD with GSI1 queries (lookup by userID via GSI1PK=`USER#{userId}`, GSI1SK=`ARTIST_PROFILE`)
  - Test Follow CRUD with GSI1 queries (lookup followers via GSI1PK=`ARTIST_PROFILE#{followedId}`)
  - Test follower/following count updates (atomic increments)
  - Purpose: Verify GSI-based queries work correctly against real DynamoDB
  - _Leverage: `testutil.SetupLocalStack()`, `repository.NewDynamoDBRepository()`, `repository/artist_profile.go`, `repository/follow.go`_
  - _Requirements: 1.4_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in DynamoDB GSI queries | Task: Create `repository/artist_profile_integration_test.go`. Read `repository/artist_profile.go` and `repository/follow.go` to understand the GSI query patterns. Write integration tests that: (1) create artist profiles and verify GSI1 lookup by userID, (2) create follows and verify GSI1 lookup of followers for an artist, (3) test atomic counter increments for follower/following counts. These are the most important tests because GSI behavior is where mocks diverge most from real DynamoDB. | Restrictions: Must verify actual GSI query results (not just PK/SK access). Use `//go:build integration` tag. | _Leverage: `repository/artist_profile.go`, `repository/follow.go` | _Requirements: Req 1 | Success: GSI1 queries return correct results. Atomic counter increments work. Mark task in-progress in tasks.md, log implementation, mark complete._

- [x] 2.3. S3 repository integration tests
  - File: `backend/internal/repository/s3_integration_test.go`
  - Test PutObject, GetObject, DeleteObject
  - Test GeneratePresignedUploadURL and GeneratePresignedDownloadURL
  - Test DeleteByPrefix (HLS cleanup): create `hls/user1/track1/seg1.ts`, `hls/user1/track1/seg2.ts`, delete by prefix, verify all gone
  - Purpose: Verify S3 operations including the batch prefix delete
  - _Leverage: `testutil.SetupLocalStack()`, `repository.NewS3Repository()`, `repository/s3.go`_
  - _Requirements: 3.1, 3.2, 3.3, 3.4_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in S3 operations | Task: Create `repository/s3_integration_test.go`. Read `repository/s3.go` for the `S3Repository` implementation. Write integration tests for: (1) PutObject + GetObject roundtrip, (2) presigned URL generation — upload via presigned URL, then download and verify content matches, (3) DeleteObject, (4) DeleteByPrefix — create multiple objects with shared prefix `hls/userX/trackY/`, call DeleteByPrefix, verify all are gone. Use `testutil.CreateTestS3Object` for setup. | Restrictions: Use `//go:build integration` tag. Clean up all S3 objects. | _Leverage: `repository/s3.go`, `testutil.CreateTestS3Object` | _Requirements: Req 3 (S3 Integration Tests) | Success: All S3 operations verified. DeleteByPrefix removes all matching objects. Presigned URLs work. Mark task in-progress in tasks.md, log implementation, mark complete._

## Phase 3: Service Integration Tests

- [ ] 3.1. Track service integration tests
  - File: `backend/internal/service/track_integration_test.go` (extend existing)
  - Test visibility enforcement: private track → owner 200, other user 403, admin 200
  - Test `ListTracks` with `hasGlobal=true` (admin) returns all tracks vs `hasGlobal=false` returns own + public
  - Test `DeleteTrack` with `hasGlobal=true` (admin can delete any track) including S3 cleanup
  - Test `UpdateVisibility` and verify track appears/disappears from public queries
  - Purpose: Verify track business logic with real DynamoDB/S3
  - _Leverage: `testutil.SetupLocalStack()`, `repository.NewDynamoDBRepository()`, `service.NewTrackService()`_
  - _Requirements: 2.1, 2.2, 2.3_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Extend `service/track_integration_test.go` with service-level integration tests. Read `service/track.go` for `TrackService` implementation — focus on the `hasGlobal` parameter in `GetTrack`, `ListTracks`, `DeleteTrack`. Create real `DynamoDBRepository` + `S3Repository` + `TrackService`. Test: (1) create users with different roles, (2) create private/public tracks, (3) verify GetTrack visibility enforcement (owner OK, other 403, admin OK), (4) verify ListTracks filtering, (5) verify admin DeleteTrack with S3 cleanup. | Restrictions: Use real services and repos, not mocks. Do not modify service code. | _Leverage: `service/track.go`, `testutil.SetupLocalStack()` | _Requirements: Req 2 | Success: Visibility enforcement verified end-to-end. Admin access works. Mark task in-progress in tasks.md, log implementation, mark complete._

- [ ] 3.2. Playlist service integration tests
  - File: `backend/internal/service/playlist_integration_test.go`
  - Test Playlist CRUD via service layer with real DynamoDB
  - Test visibility: set public → verify appears in `ListPublicPlaylists` (GSI2 query), set private → verify disappears
  - Test AddTracks/RemoveTracks with real track data
  - Purpose: Verify playlist business logic including public discovery
  - _Leverage: `testutil.SetupLocalStack()`, `service.NewPlaylistService()`_
  - _Requirements: 2.4_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create `service/playlist_integration_test.go`. Read `service/playlist.go` for implementation. Test: (1) create playlist, (2) add tracks, (3) list playlists, (4) set visibility to public, (5) verify ListPublicPlaylists returns it, (6) set private, verify it disappears from public list. Also test RemoveTracks and DeletePlaylist. | Restrictions: Use real repos, not mocks. | _Leverage: `service/playlist.go` | _Requirements: Req 2 | Success: Playlist visibility toggle verified against real DynamoDB GSI. Mark task in-progress in tasks.md, log implementation, mark complete._

- [ ] 3.3. Tag, Follow, ArtistProfile, Role, User service integration tests
  - Files: `backend/internal/service/tag_integration_test.go`, `follow_integration_test.go`, `artist_profile_integration_test.go`, `role_integration_test.go`, `user_integration_test.go`
  - Tag: create tag, add to track, list tracks by tag, case-insensitive lookup, rename, delete
  - Follow: follow, verify count increment, unfollow, verify count decrement, prevent duplicate follow
  - ArtistProfile: create, get, update, link to user via GSI
  - Role: get/update role via DynamoDB
  - User: CreateUserIfNotExists idempotency, profile updates
  - Purpose: Verify remaining service business logic with real DynamoDB
  - _Leverage: `testutil.SetupLocalStack()`, respective service constructors_
  - _Requirements: 2.5_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create integration tests for Tag, Follow, ArtistProfile, Role, and User services. Read each service file for the implementations. For each service: create real repo + service, test all major operations. Key scenarios: TagService case-insensitive lookup, FollowService atomic count updates and duplicate prevention, ArtistProfileService GSI-based user lookup, UserService CreateUserIfNotExists idempotency. Each file uses `//go:build integration` tag. | Restrictions: Separate file per service. Use real repos. | _Leverage: `service/tag.go`, `service/follow.go`, `service/artist_profile.go`, `service/role.go`, `service/user.go` | _Requirements: Req 2 | Success: All service business logic verified against real DynamoDB. Mark task in-progress in tasks.md, log implementation, mark complete._

## Phase 4: Cognito Integration Tests

- [ ] 4.1. Cognito client integration tests
  - File: `backend/internal/service/cognito_integration_test.go`
  - Test authentication: authenticate pre-created test users (admin, subscriber, artist), verify tokens
  - Test group operations: add user to group, remove from group, list groups for user
  - Test user management: search users, disable/enable user
  - Skip gracefully if Cognito not initialized
  - Purpose: Verify CognitoClient operations against LocalStack Cognito
  - _Leverage: `testutil.SetupLocalStack()`, `service.NewCognitoClient()`, pre-created users from `docker/localstack-init/init-cognito.sh`_
  - _Requirements: 4.1, 4.2, 4.3, 4.4_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in AWS Cognito | Task: Create `service/cognito_integration_test.go`. Read `service/cognito_client.go` for the CognitoClient implementation. Read `docker/localstack-init/init-cognito.sh` for pre-created test users. Test: (1) authenticate admin/subscriber/artist users, (2) list groups for admin user (verify "admin" group), (3) add/remove user from group, (4) search for user by email. Skip with `t.Skip()` if `tc.UserPoolID == ""`. | Restrictions: Use `//go:build integration` tag. Do not create new Cognito users — use pre-created ones. | _Leverage: `service/cognito_client.go`, `testutil.TestUsers`, `docker/localstack-init/init-cognito.sh` | _Requirements: Req 4 (Cognito Integration Tests) | Success: All CognitoClient operations verified against LocalStack. Mark task in-progress in tasks.md, log implementation, mark complete._

## Phase 5: Full API Integration Tests

- [ ] 5.1. API auth and health integration tests
  - File: `backend/test/api_auth_integration_test.go`
  - Test `/health` returns 200 without auth
  - Test protected endpoints return 401 without auth headers
  - Test role-based access: subscriber can't access admin routes (403)
  - Test guest vs subscriber vs artist vs admin access levels
  - Purpose: Verify auth middleware chain end-to-end
  - _Leverage: `testutil.SetupTestServer()`, `testutil.AsUser()`, `testutil.AssertStatus()`_
  - _Requirements: 6.9_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in API testing | Task: Create `test/api_auth_integration_test.go`. Use `testutil.SetupTestServer(t)` to spin up the full Echo server. Test: (1) GET `/health` → 200 without auth, (2) GET `/api/v1/tracks` without headers → 401, (3) GET `/api/v1/admin/users` as subscriber → 403, (4) GET `/api/v1/admin/users` as admin → 200. Use `AsUser()` to set auth headers. | Restrictions: Use `//go:build integration` tag. Create test users in DynamoDB for DB role checks. | _Leverage: `testutil/server.go`, `testutil/http_helpers.go` | _Requirements: Req 6 | Success: Full auth chain verified via HTTP. Mark task in-progress in tasks.md, log implementation, mark complete._

- [ ] 5.2. API tracks CRUD and visibility integration tests
  - File: `backend/test/api_tracks_integration_test.go`
  - Test full CRUD cycle via HTTP: list tracks (empty), create track data via fixture, GET track, list tracks (1 result), DELETE track, verify 404
  - Test visibility enforcement via HTTP: user A private track → user B GET returns 403 → admin GET returns 200
  - Test admin delete via HTTP: admin DELETE another user's track → 200, verify gone
  - Test list tracks as admin (sees all) vs subscriber (sees own + public)
  - Purpose: Verify track endpoints end-to-end through the full stack
  - _Leverage: `testutil.SetupTestServer()`, `testutil.AsUser()`, `testutil.WithJSON()`, `testutil.DecodeJSON()`_
  - _Requirements: 6.1, 6.3, 6.4_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go API Developer | Task: Create `test/api_tracks_integration_test.go`. Use `SetupTestServer` for full Echo + LocalStack. Test the complete track lifecycle via HTTP requests: (1) create users with different roles, (2) create tracks via fixtures, (3) GET /api/v1/tracks as different users — verify filtering, (4) GET /api/v1/tracks/:id — verify visibility (owner OK, other 403, admin OK), (5) DELETE /api/v1/tracks/:id as admin — verify deleted, (6) verify DynamoDB directly via `tsc.ItemExists`. | Restrictions: Use HTTP requests only for the operations being tested. Use fixtures for setup data. | _Leverage: `testutil/server.go`, `testutil/http_helpers.go` | _Requirements: Req 6 | Success: Full track CRUD and visibility verified through HTTP. Mark task in-progress in tasks.md, log implementation, mark complete._

- [ ] 5.3. API playlists and tags integration tests
  - File: `backend/test/api_playlists_integration_test.go`, `backend/test/api_tags_integration_test.go`
  - Playlists: create via POST, list, update visibility, verify public playlist discovery via GET `/playlists/public`, add/remove tracks
  - Tags: create via POST, add to track, list tracks by tag, remove from track, verify case-insensitive behavior
  - Purpose: Verify playlist and tag endpoints end-to-end
  - _Leverage: `testutil.SetupTestServer()`, `testutil.AsUser()`, `testutil.WithJSON()`_
  - _Requirements: 6.6, 6.7_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go API Developer | Task: Create `test/api_playlists_integration_test.go` and `test/api_tags_integration_test.go`. For playlists: POST create → GET list → PUT visibility to public → GET /playlists/public (verify appears) → PUT visibility to private → GET /playlists/public (verify gone). For tags: POST create → POST add-to-track → GET tracks-by-tag → DELETE remove-from-track. All via HTTP to test server. | Restrictions: Use `//go:build integration` tag. Both files. | _Leverage: `testutil/server.go`, `testutil/http_helpers.go` | _Requirements: Req 6 | Success: Playlist visibility and tag operations verified through HTTP. Mark task in-progress in tasks.md, log implementation, mark complete._

- [ ] 5.4. API follows and artist profiles integration tests
  - File: `backend/test/api_follows_integration_test.go`, `backend/test/api_artists_integration_test.go`
  - Artist profiles: create via POST, GET, update via PUT, list, search
  - Follows: POST follow → GET followers (verify appears) → verify count incremented → DELETE unfollow → verify count decremented
  - Purpose: Verify follow system and artist profile endpoints end-to-end
  - _Leverage: `testutil.SetupTestServer()`, `testutil.AsUser()`_
  - _Requirements: 6.8_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go API Developer | Task: Create `test/api_follows_integration_test.go` and `test/api_artists_integration_test.go`. For artist profiles: create artist user + profile via fixtures, GET /artists/entity/:id, PUT update, GET /artists/entity (list). For follows: subscriber follows artist via POST /artists/entity/:id/follow, GET /artists/entity/:id/followers (verify listed), check follower count incremented, DELETE /artists/entity/:id/follow, verify count decremented. | Restrictions: Create prerequisite data (users, artist profiles) via fixtures, test follow operations via HTTP. | _Leverage: `testutil/server.go`, `testutil/http_helpers.go`, `testutil/fixtures.go` | _Requirements: Req 6 | Success: Follow system verified through HTTP with count updates. Mark task in-progress in tasks.md, log implementation, mark complete._

- [ ] 5.5. API admin routes integration tests
  - File: `backend/test/api_admin_integration_test.go`
  - Test admin user search via GET `/api/v1/admin/users?q=...`
  - Test admin role update via PUT `/api/v1/admin/users/:id/role`
  - Test `RequireRoleWithDBCheck` middleware: change user role in DB → verify access changes on next request
  - Test non-admin gets 403 on all admin routes
  - Purpose: Verify admin endpoints with real DB role resolution
  - _Leverage: `testutil.SetupTestServer()`, `testutil.AsUser()`, LocalStack Cognito_
  - _Requirements: 6.5_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go API Developer | Task: Create `test/api_admin_integration_test.go`. The key test: `RequireRoleWithDBCheck` middleware queries DynamoDB for the real role on every request. Test: (1) create admin user in DB, (2) GET /api/v1/admin/users as admin → 200, (3) GET as subscriber → 403, (4) key scenario: user has admin in headers but subscriber in DB → middleware should use DB role → 403. Also test user search and role update endpoints if Cognito is available. | Restrictions: Admin tests may need Cognito — skip Cognito-dependent tests if not available. | _Leverage: `testutil/server.go`, `service/cognito_client.go` | _Requirements: Req 6 | Success: Admin routes verified with real DB role resolution. Mark task in-progress in tasks.md, log implementation, mark complete._

## Phase 6: CI Integration

- [ ] 6.1. GitHub Actions integration test workflow
  - File: `.github/workflows/integration.yml` (new) or extend `.github/workflows/ci.yml`
  - Add job that starts LocalStack via Docker Compose (or `services:` block)
  - Run `scripts/wait-for-localstack.sh`, then init scripts, then `go test -tags=integration ./...`
  - Set `AWS_ENDPOINT`, `DYNAMODB_TABLE_NAME`, `MEDIA_BUCKET` env vars
  - Collect coverage data from integration tests
  - Purpose: Run integration tests on every PR
  - _Leverage: `docker/docker-compose.yml`, `scripts/wait-for-localstack.sh`, `docker/localstack-init/`_
  - _Requirements: 8.1, 8.2, 8.3_
  - _Prompt: Implement the task for spec localstack-test-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer specializing in GitHub Actions and CI/CD | Task: Create `.github/workflows/integration.yml`. Read existing `.github/workflows/ci.yml` for patterns. The workflow should: (1) use `services:` block with `localstack/localstack:3.4` image on port 4566, (2) wait for LocalStack health, (3) run `docker/localstack-init/init-aws.sh` and `init-cognito.sh`, (4) run `cd backend && go test -tags=integration -coverprofile=integration-coverage.out ./...`, (5) upload coverage artifact. Set env vars: `AWS_ENDPOINT=http://localhost:4566`, `DYNAMODB_TABLE_NAME=MusicLibrary`, `MEDIA_BUCKET=music-library-local-media`. | Restrictions: Do not modify existing CI workflow. Retry LocalStack health check up to 3 times. | _Leverage: `.github/workflows/ci.yml`, `docker/docker-compose.yml` | _Requirements: Req 8 (CI Integration) | Success: Integration tests run and pass in GitHub Actions. Mark task in-progress in tasks.md, log implementation, mark complete._

## Task Dependencies

```
Phase 1 (Infrastructure)
  1.1 SetupTestServer ──────────────────────────────┐
  1.2 HTTP Helpers ─────────────────────────────────┤
  1.3 Extended Fixtures ────┐                       │
  1.4 S3 Cleanup ───────────┤                       │
                            ↓                       ↓
Phase 2 (Repository)     Phase 5 (API Tests)
  2.1 DynamoDB CRUD         5.1 Auth ──────── needs 1.1, 1.2
  2.2 GSI queries           5.2 Tracks ────── needs 1.1, 1.2, 1.3
  2.3 S3 ops                5.3 Playlists/Tags ── needs 1.1, 1.2, 1.3
         ↓                  5.4 Follows/Artists ── needs 1.1, 1.2, 1.3
Phase 3 (Service)           5.5 Admin ──────── needs 1.1, 1.2
  3.1 Track svc
  3.2 Playlist svc       Phase 6 (CI)
  3.3 Other svcs           6.1 GitHub Actions ── needs all tests written
         ↓
Phase 4 (Cognito)
  4.1 Cognito client
```

**Critical path**: 1.1 → 1.2 → 5.2 (test server + helpers → tracks API tests — highest value)

**Parallelizable**: Phase 2 (repo tests) can run in parallel with Phase 3 (service tests) once Phase 1 fixtures are done. Phase 5 API tests depend on Phase 1 only.
