# Tasks: Hello World Local Development (Iteration 4)

**Spec:** hello-world-local-dev
**SDLC Phase:** SPEC → TEST → CODE → BUILD → DOCS
**TDD Enforcement:** Every task group has separate Red (test) and Green (implementation) sub-tasks.

**Iteration 4 Goals:**
- Validate updated SDLC process docs from iterations 1-3
- Follow Contract Alignment Check (Task 7.1)
- Follow BUILD validation (Task 7.3)
- Follow updated wiring checklist with route tree regeneration

---

## Task 1: Seed Data

- [ ] 1.1 Write seed script tests (TDD Red)
  - **Test File:** `docker/localstack-init/test-seed-music.sh` (create)
  - **Acceptance Criteria:**
    - Test verifies 20 items exist in DynamoDB after running seed script
    - Test verifies items are under `USER#seed-user` partition key
    - Test verifies idempotency (running twice produces same 20 items)
    - Test verifies each item has: Title, Artist, Album, Genre, Year, Duration
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: US-1_

- [ ] 1.2 Implement seed script (TDD Green)
  - **Implementation File:** `docker/localstack-init/init-seed-music.sh` (create)
  - **Acceptance Criteria:**
    - Creates 20 tracks with 5 artists × 4 tracks each
    - Uses `aws dynamodb put-item` (idempotent)
    - Script is executable (`chmod +x`)
    - **ALL TESTS from 1.1 MUST PASS**
  - _Requirements: US-1_

- [ ] 1.3 Wire seed script into Makefile
  - **Modify:** `Makefile` (update `local-services` target)
  - **Acceptance Criteria:**
    - `make local-services` calls `init-seed-music.sh`
    - Seed data available when `make local` completes
  - _Requirements: US-1_

---

## Task 2: Backend Service

- [ ] 2.1 Write HelloService tests (TDD Red)
  - **Test File:** `backend/internal/service/hello_test.go` (create)
  - **Acceptance Criteria:**
    - Test `Search` matches title, artist, album, genre (case-insensitive)
    - Test `Search("")` returns empty slice
    - Test `Search` with no matches returns empty slice
    - Test `Featured(limit)` returns up to limit tracks
    - Test `Featured(0)` returns all tracks
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: US-2, US-3_

- [ ] 2.2 Implement HelloService (TDD Green)
  - **Implementation Files:**
    - `backend/internal/service/hello.go` (create)
    - `backend/internal/service/hello_dynamo.go` (create)
  - **Acceptance Criteria:**
    - `HelloTrack` struct with JSON tags
    - `HelloRepository` interface with `GetSeedTracks`
    - `HelloDynamoDBRepository` queries `PK=USER#seed-user, SK begins_with TRACK#`
    - **ID extracted from SK field** by stripping `TRACK#` prefix
    - **ALL TESTS from 2.1 MUST PASS**
  - _Requirements: US-2, US-3_

---

## Task 3: Backend Handler + Wiring

- [ ] 3.1 Write HelloHandler tests (TDD Red)
  - **Test File:** `backend/internal/handlers/hello_test.go` (create)
  - **Acceptance Criteria:**
    - Test Health returns `{"status":"ok","service":"hello"}`
    - Test Search reads `c.QueryParam("q")` — **NOT "query"**
    - Test Search returns `{"items":[...],"total":N}` — **NOT "tracks"**
    - Test Featured returns `{"items":[...],"total":N}`
    - Test Featured default limit is 10
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: US-2, US-3, US-4_

- [ ] 3.2 Implement HelloHandler + wire service and routes (TDD Green)
  - **Implementation Files:**
    - `backend/internal/handlers/hello.go` (create)
    - `backend/internal/service/service.go` (modify — add `Hello HelloService`)
    - `backend/cmd/api/main.go` (modify — create handler, register routes)
  - **Acceptance Criteria:**
    - Search reads `c.QueryParam("q")` — **CONTRACT CRITICAL**
    - Returns `ListResponse[HelloTrack]` with `Items` field → JSON `items`
    - `RegisterHelloRoutes(e, h)` called in `main.go`
    - Routes under `/api/v1/hello/` without auth middleware
    - **ALL TESTS from 3.1 MUST PASS**
  - **Wiring Checklist:** (1) Add `Hello HelloService` to `Services` struct, (2) Create `helloHandler := NewHelloHandler(...)` in main.go, (3) Call `RegisterHelloRoutes(e, helloHandler)` in main.go
  - _Requirements: US-2, US-3, US-4, US-6_

---

## Task 4: Frontend API Client + Hooks

- [ ] 4.1 Write frontend API and hook tests (TDD Red)
  - **Test Files:**
    - `frontend/src/lib/api/__tests__/helloSearch.test.ts` (create)
    - `frontend/src/hooks/__tests__/useHelloSearch.test.tsx` (create)
  - **Acceptance Criteria:**
    - API test: `searchHelloTracks("jazz")` sends `params: { q: "jazz" }` — **NOT { query }**
    - API test: Response uses `items` field — **NOT "tracks"**
    - Hook test: `useHelloSearch("")` is disabled
    - Hook test: `helloKeys` factory produces correct keys
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: US-2, US-3, US-5_

- [ ] 4.2 Implement frontend API client and hooks (TDD Green)
  - **Implementation Files:**
    - `frontend/src/lib/api/helloSearch.ts` (create)
    - `frontend/src/hooks/useHelloSearch.ts` (create)
  - **Acceptance Criteria:**
    - `searchHelloTracks` sends `params: { q: query }` — **CONTRACT CRITICAL**
    - Response type uses `items` field — **CONTRACT CRITICAL**
    - `useHelloSearch(query)` disabled when empty
    - **ALL TESTS from 4.1 MUST PASS**
  - _Requirements: US-2, US-3, US-5_

---

## Task 5: Frontend Components

- [ ] 5.1 Write component tests (TDD Red)
  - **Test Files:**
    - `frontend/src/components/hello/__tests__/SearchInput.test.tsx`
    - `frontend/src/components/hello/__tests__/TrackCard.test.tsx`
    - `frontend/src/components/hello/__tests__/TrackCardSkeleton.test.tsx`
    - `frontend/src/components/hello/__tests__/HelloNav.test.tsx`
  - **Acceptance Criteria:**
    - SearchInput: renders input, calls onChange, shows placeholder
    - TrackCard: formats duration (240→"4:00", 185→"3:05"), renders all fields
    - TrackCardSkeleton: has `.skeleton` class elements
    - HelloNav: shows "Hello Music Search", has home link
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: US-5_

- [ ] 5.2 Implement components (TDD Green)
  - **Implementation Files:**
    - `frontend/src/components/hello/SearchInput.tsx`
    - `frontend/src/components/hello/TrackCard.tsx`
    - `frontend/src/components/hello/TrackCardSkeleton.tsx`
    - `frontend/src/components/hello/HelloNav.tsx`
  - **Acceptance Criteria:**
    - All use DaisyUI 5 semantic classes
    - TrackCard duration format: `${Math.floor(d/60)}:${String(d%60).padStart(2,'0')}`
    - **ALL TESTS from 5.1 MUST PASS**
  - _Requirements: US-5_

---

## Task 6: Frontend Route Page

- [ ] 6.1 Write route page tests (TDD Red)
  - **Test File:** `frontend/src/routes/__tests__/hello-search.test.tsx`
  - **Acceptance Criteria:**
    - Test page renders search input and featured tracks
    - Test loading state shows skeletons
    - Test error state shows alert
    - Test empty results shows "No results found"
    - Test typing triggers search
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: US-5, US-6_

- [ ] 6.2 Implement route page + wire router (TDD Green)
  - **Implementation Files:**
    - `frontend/src/routes/hello-search.tsx` (create)
    - `frontend/src/routes/__root.tsx` (modify — add to PUBLIC_ROUTES)
  - **Acceptance Criteria:**
    - Uses `createFileRoute('/hello-search')`
    - Shows HelloNav, SearchInput, TrackCard grid
    - Shows TrackCardSkeleton during loading
    - Shows error alert on failure
    - Shows "No results found" for empty results
    - `/hello-search` added to PUBLIC_ROUTES
    - **MUST run `npx @tanstack/router-cli generate`** to regenerate routeTree.gen.ts
    - **ALL TESTS from 6.1 MUST PASS**
  - **Wiring Checklist (Frontend):** (1) Create route file with `createFileRoute('/hello-search')`, (2) Run `npx @tanstack/router-cli generate`, (3) Verify `routeTree.gen.ts` includes `/hello-search`, (4) Add `/hello-search` to `PUBLIC_ROUTES` in `__root.tsx`
  - _Requirements: US-5, US-6_

---

## Task 7: Integration, Validation, and Documentation

- [ ] 7.1 Contract alignment check
  - **Purpose:** Verify frontend and backend use the same API contract
  - **Acceptance Criteria:**
    - Frontend sends `q` param: `grep -n "params.*{.*q:" frontend/src/lib/api/helloSearch.ts`
    - Frontend expects `items`: `grep -n "items" frontend/src/lib/api/helloSearch.ts`
    - Backend reads `q`: `grep -n 'QueryParam("q")' backend/internal/handlers/hello.go`
    - Backend returns `items`: `grep -n '"items"\|Items' backend/internal/handlers/hello.go`
  - **If any check fails:** Fix the misaligned component, re-run tests, commit fix
  - _Requirements: R-NF4_

- [ ] 7.2 Local stack validation (`make local`)
  - **Purpose:** End-to-end validation with real local stack
  - **Acceptance Criteria:**
    - `make local` starts successfully with seed data
    - `curl http://localhost:8080/api/v1/hello/health` returns `{"status":"ok","service":"hello"}`
    - `curl http://localhost:8080/api/v1/hello/featured` returns 20 tracks with `items` key
    - `curl "http://localhost:8080/api/v1/hello/search?q=jazz"` returns matching tracks
    - Frontend at `http://localhost:5173/hello-search` loads and shows tracks
  - _Requirements: US-1 through US-6_

- [ ] 7.3 BUILD validation
  - **Purpose:** Ensure code passes all quality gates
  - **Acceptance Criteria:**
    - All commands pass with zero errors:
    ```bash
    # Backend
    cd backend && go vet ./...
    cd backend && go build -o /tmp/test-build ./cmd/api
    cd backend && go test ./...

    # Frontend
    cd frontend && npx --package=typescript tsc --noEmit
    cd frontend && npx eslint .
    cd frontend && npm run build
    cd frontend && npm test -- --run
    ```
  - _Requirements: R-NF2, R-NF3_

- [ ] 7.4 Update CLAUDE.md for all changed directories
  - **Purpose:** Document all new/modified code
  - **Directories requiring CLAUDE.md:**
    - `backend/cmd/api/`
    - `backend/internal/handlers/` (update)
    - `backend/internal/service/` (update)
    - `docker/localstack-init/`
    - `frontend/src/components/hello/`
    - `frontend/src/hooks/` (update)
    - `frontend/src/lib/api/` (update)
    - `frontend/src/routes/` (update)
  - **Acceptance Criteria:**
    - Each CLAUDE.md has: Overview, Files, Key Functions, Dependencies
    - Zero missing CLAUDE.md for changed directories
  - _Requirements: Hard Requirement #4 (Documentation)_

---

## Task Dependencies

```
Task 1 (Seed): Independent
Task 2 (Service): Independent
Task 3 (Handler): Depends on Task 2
Task 4 (FE API): Independent (can parallel with 2-3)
Task 5 (FE Components): Independent (can parallel with 4)
Task 6 (FE Route): Depends on Tasks 4, 5
Task 7 (Integration): Depends on all implementation tasks
```

**Parallelizable:**
- Tasks 1, 2, 4, 5 can run in parallel
- Task 6 after 4, 5 complete
- Task 7 after all implementation complete

---

## SDLC Gap Tracking

_Gaps discovered during iteration 4 will be logged here._

| # | Phase | Gap Description | Impact | Fix Applied |
|---|-------|-----------------|--------|-------------|
| - | - | - | - | - |
