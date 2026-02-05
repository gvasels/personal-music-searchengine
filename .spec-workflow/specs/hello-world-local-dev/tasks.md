# Tasks Document

**Spec**: hello-world-local-dev
**Design**: [design.md](./design.md)
**Requirements**: [requirements.md](./requirements.md)

---

## CRITICAL: TDD Task Structure

**Every feature MUST follow this ordering:**
1. **Test tasks FIRST** (TDD Red) - Write failing tests
2. **Implementation tasks SECOND** (TDD Green) - Write code to make tests pass

**DO NOT** put implementation tasks before their corresponding test tasks.
**DO NOT** combine test and implementation in a single task.

---

## Task 1: Seed Data Scripts

- [ ] 1.1 Write seed data validation script (TDD Red)
  - **Test Script**: `docker/localstack-init/test-seed-data.sh` (create)
  - **Acceptance Criteria**:
    - Validates 20 items exist in DynamoDB under `USER#seed-user` partition
    - Validates each item has required fields: title, artist, album, genre, year, duration
    - Validates 5 distinct artists exist
    - Validates idempotency (running seed twice produces same count)
    - **MUST RUN AND FAIL** before seed script implementation (no data exists yet)
  - _Leverage: existing init-aws.sh for AWS CLI patterns_
  - _Requirements: R1_

- [ ] 1.2 Implement seed data script (TDD Green)
  - **Implementation**: `docker/localstack-init/init-seed-music.sh` (create)
  - **Acceptance Criteria**:
    - Seeds 20 mock tracks across 5 artists into MusicLibrary table
    - Each track has: PK, SK, id, title, artist, album, genre, year, duration, Type, createdAt
    - Script is idempotent (checks for existing data before inserting)
    - **test-seed-data.sh MUST PASS after implementation**
  - _Leverage: init-aws.sh for AWS CLI and DynamoDB put-item patterns_
  - _Requirements: R1_

- [ ] 1.3 Wire seed script into Makefile
  - **Implementation**: `Makefile` (modify)
  - **Acceptance Criteria**:
    - `local-services` target calls `./docker/localstack-init/init-seed-music.sh` after `init-cognito.sh`
  - _Requirements: R1_

---

## Task 2: Backend Hello Service

- [ ] 2.1 Write backend unit tests (TDD Red)
  - **Unit Test**: `backend/internal/service/hello_test.go` (create)
  - **Acceptance Criteria**:
    - Test SearchTracks with matching query returns correct tracks
    - Test SearchTracks case-insensitivity
    - Test SearchTracks with no matches returns empty slice
    - Test SearchTracks respects limit parameter
    - Test GetFeaturedTracks returns all tracks up to limit
    - Test GetFeaturedTracks with custom limit
    - Test GetFeaturedTracks with zero tracks
    - **MUST RUN AND FAIL** before implementation
  - _Leverage: testify/assert, testify/mock for HelloRepository_
  - _Requirements: R2, R3_

- [ ] 2.2 Implement HelloService (TDD Green)
  - **Implementation**: `backend/internal/service/hello.go` (create)
  - **Acceptance Criteria**:
    - HelloService interface with SearchTracks and GetFeaturedTracks
    - HelloRepository interface with ListTracksByUser method
    - HelloRepoAdapter adapts Repository to HelloRepository
    - Case-insensitive search across title, artist, album, genre
    - **ALL TESTS FROM 2.1 MUST PASS**
  - _Leverage: service.go Services struct pattern_
  - _Requirements: R2, R3_

---

## Task 3: Backend Hello Handler

- [ ] 3.1 Write backend handler tests (TDD Red)
  - **Unit Test**: `backend/internal/handlers/hello_test.go` (create)
  - **Acceptance Criteria**:
    - Test HelloHealth returns {"status":"ok","service":"hello"}
    - Test HelloSearch with valid query returns 200 with items
    - Test HelloSearch without query returns 400
    - Test HelloFeatured returns 200 with items
    - Test HelloSearch error handling
    - **MUST RUN AND FAIL** before implementation
  - _Leverage: echo test context, httptest, testify/mock_
  - _Requirements: R2, R3, R4_

- [ ] 3.2 Implement HelloHandler and wire routes (TDD Green)
  - **Implementation**: `backend/internal/handlers/hello.go` (create), `handlers.go` (modify), `service.go` (modify)
  - **Acceptance Criteria**:
    - HelloHealth, HelloSearch, HelloFeatured handler methods
    - Routes registered: GET /api/v1/hello/health, /search, /featured
    - HelloService added to Services struct and wired in NewServices
    - **ALL TESTS FROM 3.1 MUST PASS**
  - _Leverage: handlers.go RegisterRoutes pattern, handleError/success helpers_
  - _Requirements: R2, R3, R4_

---

## Task 4: Frontend API Client and Hook

- [ ] 4.1 Write frontend API and hook tests (TDD Red)
  - **Unit Tests**:
    - `frontend/src/lib/api/__tests__/helloSearch.test.ts` (create)
    - `frontend/src/hooks/__tests__/useHelloSearch.test.ts` (create)
  - **Acceptance Criteria**:
    - Test searchHello calls correct URL with query param
    - Test getFeaturedTracks calls correct URL
    - Test useHelloSearch returns data and respects enabled flag
    - Test useFeaturedTracks returns featured tracks
    - Test helloSearchKeys factory produces correct keys
    - **MUST RUN AND FAIL** before implementation
  - _Leverage: vitest mock for apiClient, test-utils.tsx renderHook wrapper_
  - _Requirements: R2, R3, R5_

- [ ] 4.2 Implement API client and hook (TDD Green)
  - **Implementation**:
    - `frontend/src/lib/api/helloSearch.ts` (create)
    - `frontend/src/hooks/useHelloSearch.ts` (create)
  - **Acceptance Criteria**:
    - searchHello and getFeaturedTracks API functions
    - helloSearchKeys query key factory
    - useHelloSearch and useFeaturedTracks hooks
    - **ALL TESTS FROM 4.1 MUST PASS**
  - _Leverage: apiClient from client.ts, queryKey pattern from other hooks_
  - _Requirements: R2, R3, R5_

---

## Task 5: Frontend Hello Components

- [ ] 5.1 Write component tests (TDD Red)
  - **Unit Tests**:
    - `frontend/src/components/hello/__tests__/SearchInput.test.tsx` (create)
    - `frontend/src/components/hello/__tests__/TrackCard.test.tsx` (create)
    - `frontend/src/components/hello/__tests__/TrackCardSkeleton.test.tsx` (create)
    - `frontend/src/components/hello/__tests__/HelloNav.test.tsx` (create)
  - **Acceptance Criteria**:
    - SearchInput: test Enter fires onSearch, Escape clears, value updates
    - TrackCard: test renders title, artist, album, genre badge, duration
    - TrackCardSkeleton: test renders skeleton elements
    - HelloNav: test search callback integration, theme toggle presence
    - **MUST RUN AND FAIL** before implementation
  - _Leverage: @testing-library/react, @testing-library/user-event, vitest_
  - _Requirements: R5_

- [ ] 5.2 Implement hello components (TDD Green)
  - **Implementation**:
    - `frontend/src/components/hello/SearchInput.tsx` (create)
    - `frontend/src/components/hello/TrackCard.tsx` (create)
    - `frontend/src/components/hello/TrackCardSkeleton.tsx` (create)
    - `frontend/src/components/hello/HelloNav.tsx` (create)
  - **Acceptance Criteria**:
    - SearchInput with Enter/Escape key handling
    - TrackCard with DaisyUI card styling, genre badge, formatted duration
    - TrackCardSkeleton matching TrackCard layout
    - HelloNav with search, theme toggle (swap component), mobile dock
    - **ALL TESTS FROM 5.1 MUST PASS**
  - _Leverage: DaisyUI 5 component classes, themeStore for toggle_
  - _Requirements: R5_

---

## Task 6: Frontend Hello Search Route

- [ ] 6.1 Write route tests (TDD Red)
  - **Unit Test**: `frontend/src/routes/__tests__/hello-search.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test hero section renders with title and description
    - Test featured tracks display when loaded
    - Test search triggers API call and displays results
    - Test loading state shows skeletons
    - Test empty search results show message
    - Test error state renders error message
    - **MUST RUN AND FAIL** before implementation
  - _Leverage: test-utils.tsx custom render, vitest mock for hooks_
  - _Requirements: R5_

- [ ] 6.2 Implement hello-search route (TDD Green)
  - **Implementation**: `frontend/src/routes/hello-search.tsx` (create)
  - **Acceptance Criteria**:
    - Route at `/hello-search` with hero section
    - Featured tracks grid using TrackCard components
    - Search functionality via SearchInput
    - Loading skeletons via TrackCardSkeleton
    - Empty state and error handling
    - Add `/hello-search` to PUBLIC_ROUTES in `__root.tsx`
    - **ALL TESTS FROM 6.1 MUST PASS**
  - _Leverage: existing route patterns, hello components from Task 5_
  - _Requirements: R5_

---

## Task 7: Integration and Documentation

- [ ] 7.1 Validate full stack locally
  - **Validation**: `make local` + manual verification
  - **Acceptance Criteria**:
    - `make local` starts successfully with seed data
    - `curl http://localhost:8080/api/v1/hello/health` returns ok
    - `curl http://localhost:8080/api/v1/hello/featured` returns 20 tracks
    - `curl http://localhost:8080/api/v1/hello/search?q=luna` returns matching tracks
    - Frontend at `http://localhost:5173/hello-search` renders search page

- [ ] 7.2 Update CLAUDE.md for all changed directories
  - **Documentation**: Create/update CLAUDE.md files
  - **Directories requiring CLAUDE.md updates**:
    - `backend/internal/handlers/CLAUDE.md`
    - `backend/internal/service/CLAUDE.md`
    - `docker/CLAUDE.md`
    - `frontend/CLAUDE.md`
    - `frontend/src/CLAUDE.md`
  - **Acceptance Criteria**:
    - Each CLAUDE.md has: Overview, Files, Key Functions, Dependencies
    - Zero missing CLAUDE.md for changed directories

---

## Task Structure Rules

| Rule | Requirement |
|------|-------------|
| **Test before code** | Unit test tasks (N.1) MUST come before implementation tasks (N.2) |
| **Separate test tasks** | Unit tests and implementation are separate tasks |
| **MUST RUN AND FAIL** | Every test task must include this requirement |
| **ALL TESTS MUST PASS** | Every implementation task must include this requirement |
| **Leverage existing code** | Each task notes what existing code it reuses |
| **Requirement traceability** | Each task links to requirement IDs |
| **CLAUDE.md final task** | Last task group includes CLAUDE.md verification |
| **Separate commits** | Red phase gets its own commit, Green phase gets its own commit |

---

## SDLC Gap Tracking

[Document any gaps found during this spec's lifecycle. Include: what was expected, what actually happened, and proposed fix.]

| # | Gap | Where | Impact | Proposed Fix |
|---|-----|-------|--------|--------------|
| 1 | `test-engineer` and `implementation-agent` subagent_types hallucinate tool calls (tool_uses: 0) | Task tool agent spawning | Files never written to disk; entire TDD phase silently fails | Use `general-purpose` subagent_type for ALL file-writing agents |
| 2 | Dashboard not available by default for spec approvals | Phase 1 Spec Approval | User must manually run `spec-workflow-mcp --dashboard` | Document dashboard setup in Phase 1 prerequisites |
| 3 | Combined test+implementation in single agents (prior session) | Phase 2/3 agent prompts | TDD Red phase skipped entirely | Separate agent spawns + separate commits enforced |
