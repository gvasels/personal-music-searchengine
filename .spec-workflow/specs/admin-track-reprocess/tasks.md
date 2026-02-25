# Tasks Document: Admin Track Reprocess Button

## TDD Task Structure

Tasks are organized following strict TDD workflow:
- **N.1**: Backend unit tests (Red phase)
- **N.2**: Frontend unit tests (Red phase)
- **N.3**: Backend implementation (Green phase)
- **N.4**: Frontend implementation (Green phase)
- **N.5**: Integration tests (if applicable)

---

## Task 1: Backend Model & Service

- [ ] 1.1 Backend Unit Tests - Track Reprocess Service: Create failing tests for `TrackService.ReprocessTrack`
  - Files: `backend/internal/service/track_test.go` (extend), `backend/internal/models/track_test.go` (extend)
  - Tests: Test successful reprocess, track not found, partial success, status transitions
  - _Leverage: Existing mock patterns in `service/track_test.go`, `models/track_test.go`_
  - _Requirements: 2.1-2.7, 3.1-3.3_
  - _Prompt: Role: Go Test Engineer specializing in TDD and testify mocks | Task: Write failing unit tests for TrackService.ReprocessTrack method and Track model extensions including TestReprocessTrack_Success, TestReprocessTrack_TrackNotFound, TestReprocessTrack_PartialSuccess scenarios | Restrictions: Tests MUST fail initially, use testify/assert and testify/mock, follow existing test patterns | Success: Tests compile but fail with method not found_

- [ ] 1.2 Backend Unit Tests - Track Handler: Create failing tests for `Handlers.ReprocessTrack`
  - Files: `backend/internal/handlers/track_test.go` (extend or create)
  - Tests: Test admin can reprocess (202), non-admin gets 403, missing track returns 404
  - _Leverage: Existing handler test patterns, mock service_
  - _Requirements: 2.1-2.3_
  - _Prompt: Role: Go Test Engineer specializing in HTTP handler testing | Task: Write failing unit tests for Handlers.ReprocessTrack endpoint including TestReprocessTrack_AdminSuccess, TestReprocessTrack_NonAdmin_Forbidden, TestReprocessTrack_NotFound | Restrictions: Tests MUST fail initially, use httptest and Echo test helpers, mock TrackService | Success: Tests compile but fail with handler not found_

- [ ] 1.3 Backend Implementation - Model & Service: Implement Track model extensions and ReprocessTrack service method
  - Files: `backend/internal/models/track.go`, `backend/internal/service/track.go`, `backend/internal/service/service.go`
  - Implementation: Add ReprocessStatus, ReprocessedAt, ReprocessError to Track, add ReprocessResult struct, implement ReprocessTrack
  - _Leverage: Existing TrackService patterns, EmbeddingService.GenerateTrackEmbedding_
  - _Requirements: 2.4-2.7, 3.1-3.3_
  - _Prompt: Role: Go Backend Developer with DynamoDB and service layer expertise | Task: Implement Track model extensions and ReprocessTrack service method to make tests pass including status transitions and partial failure handling | Restrictions: ONLY implement what is needed to pass tests, follow existing patterns | Success: All 1.1 tests pass_

- [ ] 1.4 Backend Implementation - Handler & Route: Implement ReprocessTrack handler and register route
  - Files: `backend/internal/handlers/track.go`, `backend/internal/handlers/handlers.go`
  - Implementation: Add ReprocessTrack handler, register POST /tracks/:id/reprocess with admin middleware
  - _Leverage: Existing handler patterns, getAuthContextWithDBRole_
  - _Requirements: 2.1-2.3_
  - _Prompt: Role: Go Backend Developer with Echo HTTP framework expertise | Task: Implement ReprocessTrack handler using getAuthContextWithDBRole for real-time admin check and register route with RequireRole middleware | Restrictions: Follow existing handler patterns, real-time DB role check required | Success: All 1.2 tests pass_

---

## Task 2: Frontend Components

- [ ] 2.1 Frontend Unit Tests - ReprocessButton Component: Create failing tests for ReprocessButton component
  - Files: `frontend/src/components/library/__tests__/ReprocessButton.test.tsx` (new)
  - Tests: Button renders, loading state, success toast, error toast
  - _Leverage: Existing component test patterns, vitest, testing-library_
  - _Requirements: 1.3-1.5_
  - _Prompt: Role: React Test Engineer with Vitest and Testing Library expertise | Task: Write failing unit tests for ReprocessButton component including renders button with icon, shows loading spinner, disables during reprocessing, calls onSuccess callback | Restrictions: Tests MUST fail initially, mock useMutation and react-hot-toast | Success: Tests compile but fail with component not found_

- [ ] 2.2 Frontend Unit Tests - TrackList with Reprocess: Create failing tests for TrackList reprocess button integration
  - Files: `frontend/src/components/library/__tests__/TrackList.test.tsx` (extend)
  - Tests: Button visible for admin, hidden for non-admin, hidden when showUploadedBy false
  - _Leverage: Existing TrackList tests_
  - _Requirements: 1.1, 1.2_
  - _Prompt: Role: React Test Engineer | Task: Write failing tests for TrackList reprocess button integration including shows button for admin, hides for non-admin, respects showUploadedBy setting | Restrictions: Tests MUST fail initially, mock useAuth hook, extend existing tests | Success: Tests compile but fail_

- [ ] 2.3 Frontend Unit Tests - API Client: Create failing tests for reprocessTrack API function
  - Files: `frontend/src/lib/api/__tests__/tracks.test.ts` (extend)
  - Tests: Sends POST to correct endpoint, returns ReprocessResult
  - _Leverage: Existing API test patterns_
  - _Requirements: 2.1_
  - _Prompt: Role: TypeScript Test Engineer | Task: Write failing tests for reprocessTrack API function including sends POST request to /tracks/:id/reprocess, returns ReprocessResult on success | Restrictions: Tests MUST fail initially, mock axios/apiClient | Success: Tests compile but fail with function not found_

- [ ] 2.4 Frontend Implementation - API Client & Types: Implement reprocessTrack API function and types
  - Files: `frontend/src/lib/api/tracks.ts`, `frontend/src/types/index.ts`
  - Implementation: Add ReprocessResult interface, add reprocessTrack function
  - _Leverage: Existing API function patterns_
  - _Requirements: 2.1_
  - _Prompt: Role: TypeScript Developer | Task: Implement reprocessTrack API function that POSTs to /api/v1/tracks/:id/reprocess and ReprocessResult interface | Restrictions: Follow existing API patterns, use apiClient from client.ts | Success: All 2.3 tests pass_

- [ ] 2.5 Frontend Implementation - ReprocessButton Component: Implement ReprocessButton component
  - Files: `frontend/src/components/library/ReprocessButton.tsx` (new), `frontend/src/components/library/index.ts`
  - Implementation: Create ReprocessButton with useMutation, loading state, toast notifications
  - _Leverage: Existing button patterns, useMutation patterns_
  - _Requirements: 1.3-1.5_
  - _Prompt: Role: React Developer with TanStack Query expertise | Task: Implement ReprocessButton component with useMutation hook, loading spinner, disabled state during loading, toast on success/error, onSuccess callback | Restrictions: Follow existing component patterns, use react-hot-toast, style with DaisyUI btn-ghost btn-xs | Success: All 2.1 tests pass_

- [ ] 2.6 Frontend Implementation - TrackList Integration: Integrate ReprocessButton into TrackList
  - Files: `frontend/src/components/library/TrackList.tsx`
  - Implementation: Import ReprocessButton, add to actions column when admin + showUploadedBy
  - _Leverage: Existing TrackList patterns_
  - _Requirements: 1.1, 1.2_
  - _Prompt: Role: React Developer | Task: Integrate ReprocessButton into TrackList actions column, showing only when isAdmin and showUploadedByColumn are true, pass onSuccess to invalidate queries | Restrictions: Minimal changes, follow existing action button patterns | Success: All 2.2 tests pass_

---

## Task 3: Integration & Documentation

- [ ] 3.1 Integration Tests: Create integration tests for full reprocess flow
  - Files: `backend/test/api_tracks_integration_test.go` (extend)
  - Tests: Full reprocess flow with LocalStack, admin authorization
  - _Leverage: Existing integration test patterns_
  - _Requirements: All_
  - _Prompt: Role: Integration Test Engineer | Task: Create integration tests for track reprocess endpoint including TestReprocessTrack_Integration_AdminSuccess with LocalStack, TestReprocessTrack_Integration_NonAdminForbidden | Restrictions: Use LocalStack S3 + DynamoDB, tag with go:build integration | Success: Integration tests pass_

- [ ] 3.2 Documentation Updates: Update CLAUDE.md files for changed directories
  - Files: `backend/internal/handlers/CLAUDE.md`, `backend/internal/service/CLAUDE.md`, `frontend/src/components/library/CLAUDE.md`
  - Updates: Add ReprocessTrack handler, method, and ReprocessButton component
  - _Leverage: Existing CLAUDE.md templates_
  - _Requirements: Documentation standard_
  - _Prompt: Role: Technical Writer | Task: Update CLAUDE.md files to document new ReprocessTrack handler, ReprocessTrack service method, and ReprocessButton component | Restrictions: Follow existing CLAUDE.md format, keep updates concise | Success: CLAUDE.md files updated with new components_

---

## Summary

| Task | Type | Dependencies | Files |
|------|------|--------------|-------|
| 1.1 | BE Tests (Red) | - | service/track_test.go, models/track_test.go |
| 1.2 | BE Tests (Red) | - | handlers/track_test.go |
| 1.3 | BE Impl (Green) | 1.1 | models/track.go, service/track.go |
| 1.4 | BE Impl (Green) | 1.2, 1.3 | handlers/track.go, handlers.go |
| 2.1 | FE Tests (Red) | - | ReprocessButton.test.tsx |
| 2.2 | FE Tests (Red) | - | TrackList.test.tsx |
| 2.3 | FE Tests (Red) | - | api/tracks.test.ts |
| 2.4 | FE Impl (Green) | 2.3 | api/tracks.ts, types/index.ts |
| 2.5 | FE Impl (Green) | 2.1, 2.4 | ReprocessButton.tsx |
| 2.6 | FE Impl (Green) | 2.2, 2.5 | TrackList.tsx |
| 3.1 | Integration | 1.*, 2.* | api_tracks_integration_test.go |
| 3.2 | Docs | All | CLAUDE.md files |
