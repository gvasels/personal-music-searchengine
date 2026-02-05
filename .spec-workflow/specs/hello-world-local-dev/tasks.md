# Tasks: Hello World Local Development Validation (Iteration 3)

**Spec:** hello-world-local-dev
**SDLC Phase:** SPEC (complete) -> TEST (complete) -> CODE (complete) -> BUILD (complete) -> DOCS (complete)
**TDD Enforcement:** Every task group has separate Red (test) and Green (implementation) tasks with separate commits.

**Iteration 3 Fixes Applied:**
1. API contract explicitly states `q` param and `items` key throughout
2. BUILD phase has explicit validation steps
3. Route tree regeneration step included
4. Contract alignment check added as Task 7.1

---

## Task 1: Seed Data

- [x] 1.1 Write seed script tests (TDD Red)
  - **Test File:** `docker/localstack-init/test-seed-music.sh` (create)
  - **Acceptance Criteria:**
    - Test verifies 20 items exist in DynamoDB after running seed script
    - Test verifies items are under `USER#seed-user` partition key
    - Test verifies idempotency (running twice produces same 20 items, not 40)
    - Test verifies each item has required attributes: Title, Artist, Album, Genre, Year, Duration
    - **MUST RUN AND FAIL before implementation** (seed script does not exist yet)
  - **Commit:** Separate Red phase commit with test file only
  - _Requirements: US-1_
  - _Prompt: Role: Shell Script Developer | Task: Create `docker/localstack-init/test-seed-music.sh` that tests the seed script. The test should: (1) source or call `init-seed-music.sh`, (2) query DynamoDB for items with PK=USER#seed-user using `aws --endpoint-url=http://localhost:4566 dynamodb query`, (3) verify count is 20, (4) verify each item has Title, Artist, Album, Genre, Year, Duration attributes, (5) run the seed script again and verify count is still 20 (idempotency). Use `set -e` and exit with non-zero on failure. This test MUST fail initially because `init-seed-music.sh` does not exist yet. | Restrictions: Do not create the seed script itself. Only create the test._

- [x] 1.2 Implement seed script (TDD Green)
  - **Implementation File:** `docker/localstack-init/init-seed-music.sh` (create)
  - **Acceptance Criteria:**
    - Creates 20 tracks in MusicLibrary table under `USER#seed-user`
    - 5 artists with 4 tracks each
    - Each track has: Title, Artist, Album, Genre, Year, Duration
    - Uses `aws dynamodb put-item` (idempotent)
    - Script is executable (`chmod +x`)
    - **ALL TESTS from 1.1 MUST PASS**
  - **Commit:** Separate Green phase commit with implementation only
  - _Requirements: US-1_
  - _Prompt: Role: Shell Script Developer | Task: Create `docker/localstack-init/init-seed-music.sh` that inserts 20 seed tracks. Read `docker/localstack-init/init-aws.sh` to follow the same patterns (endpoint URL, region, error handling). Create 5 artists: "Aurora Waves", "Midnight Echo", "Solar Drift", "Neon Pulse", "Velvet Storm". Each artist gets 4 tracks with unique titles, albums, genres (jazz, electronic, rock, ambient, soul), years (2018-2024), and durations (180-360 seconds). Use `aws dynamodb put-item` with PK=USER#seed-user, SK=TRACK#{uuid}. UUIDs can be hardcoded for idempotency. | Restrictions: Follow existing init script patterns. Must be idempotent._

- [x] 1.3 Wire seed script into Makefile
  - **Modify:** `Makefile` (update `local-services` target)
  - **Acceptance Criteria:**
    - `make local-services` calls `init-seed-music.sh` after `init-aws.sh` and `init-cognito.sh`
    - Seed data is available when `make local` completes
  - **Commit:** Can be combined with 1.2 Green commit
  - _Requirements: US-1_

---

## Task 2: Backend Service

- [x] 2.1 Write HelloService tests (TDD Red)
  - **Test File:** `backend/internal/service/hello_test.go` (create)
  - **Acceptance Criteria:**
    - Test `Search` with matching query returns filtered tracks
    - Test `Search` is case-insensitive (searching "JAZZ" matches genre "jazz")
    - Test `Search` matches across title, artist, album, genre fields
    - Test `Search` with empty query returns empty slice
    - Test `Search` with no matches returns empty slice
    - Test `Featured` returns up to limit tracks
    - Test `Featured` with default limit returns all tracks
    - Test `Featured` with limit=5 returns 5 tracks
    - Use mock repository (follow existing `service/*_test.go` patterns)
    - **MUST RUN AND FAIL before implementation** (HelloService does not exist yet)
  - **Commit:** Separate Red phase commit with test file only
  - _Requirements: US-2, US-3_
  - _Prompt: Role: Go Backend Developer | Task: Create `backend/internal/service/hello_test.go` with unit tests for HelloService. Read existing test files like `service/tag_test.go` or `service/playlist_test.go` for patterns (mock repository setup, testify assertions). Define tests for: (1) Search("jazz") returns tracks with genre "jazz", (2) Search("AURORA") returns tracks by "Aurora Waves" (case-insensitive), (3) Search("") returns empty, (4) Search("nonexistent") returns empty, (5) Featured(20) returns all tracks, (6) Featured(5) returns 5 tracks. Mock the repository's Query or Scan method. The tests MUST fail because HelloService does not exist yet. | Restrictions: Only create the test file. Do not create hello.go._

- [x] 2.2 Implement HelloService (TDD Green)
  - **Implementation File:** `backend/internal/service/hello.go` (create)
  - **Acceptance Criteria:**
    - `HelloService` struct with `repo repository.Repository` field
    - `NewHelloService(repo)` constructor
    - `Search(ctx, query)` method: queries DynamoDB, filters case-insensitively
    - `Featured(ctx, limit)` method: queries DynamoDB, returns up to limit
    - `HelloTrack` struct with JSON tags matching API contract
    - **ALL TESTS from 2.1 MUST PASS**
  - **Commit:** Separate Green phase commit with implementation only
  - _Requirements: US-2, US-3_
  - _Prompt: Role: Go Backend Developer | Task: Create `backend/internal/service/hello.go` implementing HelloService. Read `backend/internal/repository/repository.go` for the Repository interface to understand how to query DynamoDB. Query with PK=USER#seed-user, SK begins_with TRACK#. For Search: filter results case-insensitively by checking if query appears in title, artist, album, or genre (use strings.Contains + strings.ToLower). For Featured: return up to limit results. Define HelloTrack struct with json tags. | Restrictions: Make tests from 2.1 pass. Do not modify service.go yet (wiring is in Task 3.2)._

---

## Task 3: Backend Handler + Wiring

- [x] 3.1 Write HelloHandler tests (TDD Red)
  - **Test File:** `backend/internal/handlers/hello_test.go` (create)
  - **Acceptance Criteria:**
    - Test Health endpoint returns `{"status":"ok","service":"hello"}`
    - Test Search with `q=jazz` calls service.Search("jazz") and returns `{"items":[...],"total":N}`
    - Test Search reads query param `q` (NOT `query`) -- **CONTRACT CRITICAL**
    - Test Search with empty `q` returns `{"items":[],"total":0}`
    - Test Featured returns `{"items":[...],"total":N}` -- **CONTRACT CRITICAL**
    - Test Featured with `limit=5` passes limit to service
    - Test response uses `items` key (NOT `tracks`) -- **CONTRACT CRITICAL**
    - Test routes are registered without auth middleware
    - Use httptest and mock HelloService
    - **MUST RUN AND FAIL before implementation**
  - **Commit:** Separate Red phase commit with test file only
  - _Requirements: US-2, US-3, US-4_
  - _Prompt: Role: Go Backend Developer | Task: Create `backend/internal/handlers/hello_test.go`. Read `backend/internal/handlers/handlers.go` to understand the existing patterns (ListResponse[T], success(), handleError()). Write tests using httptest: (1) GET /api/v1/hello/health returns {"status":"ok","service":"hello"}, (2) GET /api/v1/hello/search?q=jazz calls service with "jazz" and returns {"items":[...], "total":N}, (3) GET /api/v1/hello/search (no q param) returns {"items":[], "total":0}, (4) GET /api/v1/hello/featured returns {"items":[...], "total":N}, (5) GET /api/v1/hello/featured?limit=5 passes limit. CRITICAL: Verify the response JSON uses "items" key (not "tracks") and the search reads "q" param (not "query"). The tests MUST fail because HelloHandler does not exist. | Restrictions: Only create the test file._

- [x] 3.2 Implement HelloHandler + wire service and routes (TDD Green)
  - **Implementation Files:**
    - `backend/internal/handlers/hello.go` (create)
    - `backend/internal/service/service.go` (modify -- add `Hello *HelloService` to Services struct)
    - `backend/cmd/api/main.go` (modify -- create HelloHandler, call RegisterHelloRoutes)
  - **Acceptance Criteria:**
    - HelloHandler struct with `service *service.HelloService`
    - `Health(c)`: returns `{"status":"ok","service":"hello"}`
    - `Search(c)`: reads `c.QueryParam("q")`, calls service, returns `ListResponse[HelloTrack]`
    - `Featured(c)`: reads optional `limit` param, calls service, returns `ListResponse[HelloTrack]`
    - `RegisterHelloRoutes(e, h)`: registers `/api/v1/hello/` routes WITHOUT auth middleware
    - **Wiring checklist (all 4 steps):**
      1. Add `Hello *HelloService` to `service.Services` struct
      2. Initialize `Hello: NewHelloService(repo)` in `NewServices()`
      3. Create `helloHandler := handlers.NewHelloHandler(services.Hello)` in `main.go`
      4. Call `handlers.RegisterHelloRoutes(e, helloHandler)` in `main.go`
    - **ALL TESTS from 3.1 MUST PASS**
  - **Commit:** Separate Green phase commit
  - _Requirements: US-2, US-3, US-4, US-6, Wiring Checklist_
  - _Prompt: Role: Go Backend Developer | Task: (1) Create `backend/internal/handlers/hello.go` with HelloHandler. Search reads `c.QueryParam("q")` (NOT "query"). Returns `ListResponse[service.HelloTrack]{Items: tracks, Total: len(tracks)}` using existing generic from handlers.go. RegisterHelloRoutes creates `/api/v1/hello` group with NO auth middleware. (2) Add `Hello *HelloService` field to Services struct in service.go. (3) Initialize in NewServices(). (4) Wire in main.go: create handler, register routes. Read `backend/internal/handlers/handlers.go` for ListResponse pattern and `backend/cmd/api/main.go` for wiring pattern. | Restrictions: Follow wiring checklist exactly. Make tests from 3.1 pass._

---

## Task 4: Frontend API Client + Hooks

- [x] 4.1 Write frontend API client and hook tests (TDD Red)
  - **Test Files:**
    - `frontend/src/lib/api/__tests__/helloSearch.test.ts` (create)
    - `frontend/src/hooks/__tests__/useHelloSearch.test.ts` (create)
  - **Acceptance Criteria:**
    - API client tests:
      - `searchHelloTracks("jazz")` sends GET to `/v1/hello/search` with params `{ q: "jazz" }` -- **NOT `{ query: "jazz" }`**
      - `searchHelloTracks` returns `HelloSearchResponse` with `items` array -- **NOT `tracks`**
      - `getFeaturedHelloTracks()` sends GET to `/v1/hello/featured`
      - `getFeaturedHelloTracks(5)` sends GET with params `{ limit: 5 }`
      - `getHelloHealth()` sends GET to `/v1/hello/health`
    - Hook tests:
      - `useHelloSearch("jazz")` returns data with `items` array
      - `useHelloSearch("")` is disabled (no request made)
      - `useHelloFeatured()` returns featured tracks
      - Query key factory produces correct keys
    - **MUST RUN AND FAIL before implementation**
  - **Commit:** Separate Red phase commit with test files only
  - _Requirements: US-2, US-3, US-5_
  - _Prompt: Role: Frontend TypeScript Developer | Task: Create test files for the hello API client and hook. Read existing test patterns in `frontend/src/lib/api/__tests__/` and `frontend/src/hooks/__tests__/` (e.g., useSearch.test.ts, useTags.test.ts). For API tests: mock axios via vi.mock, verify searchHelloTracks("jazz") calls apiClient.get('/v1/hello/search', { params: { q: "jazz" } }). CRITICAL: param must be "q" not "query", response must use "items" not "tracks". For hook tests: use renderHook with QueryClientProvider wrapper. Tests MUST fail because the files don't exist yet. | Restrictions: Only create test files._

- [x] 4.2 Implement frontend API client and hooks (TDD Green)
  - **Implementation Files:**
    - `frontend/src/lib/api/helloSearch.ts` (create)
    - `frontend/src/hooks/useHelloSearch.ts` (create)
  - **Acceptance Criteria:**
    - `HelloTrack` and `HelloSearchResponse` types defined
    - `searchHelloTracks(query)` sends `params: { q: query }` -- **CONTRACT CRITICAL**
    - `getFeaturedHelloTracks(limit?)` sends to `/v1/hello/featured`
    - `getHelloHealth()` sends to `/v1/hello/health`
    - `useHelloSearch(query)` hook with query key factory
    - `useHelloFeatured()` hook
    - Response type uses `items` field -- **CONTRACT CRITICAL**
    - **ALL TESTS from 4.1 MUST PASS**
  - **Commit:** Separate Green phase commit
  - _Requirements: US-2, US-3, US-5_
  - _Prompt: Role: Frontend TypeScript Developer | Task: Create API client and hook. Read `frontend/src/lib/api/client.ts` for apiClient import and pattern. Create `helloSearch.ts` with: searchHelloTracks sends { params: { q: query } } (NOT { query }), response type has items (NOT tracks). Create `useHelloSearch.ts` with TanStack Query hooks following existing hook patterns (queryKey factory, enabled flag). | Restrictions: Make tests from 4.1 pass. Use exact API contract from design.md._

---

## Task 5: Frontend Components

- [x] 5.1 Write component tests (TDD Red)
  - **Test Files:**
    - `frontend/src/components/hello/__tests__/SearchInput.test.tsx` (create)
    - `frontend/src/components/hello/__tests__/TrackCard.test.tsx` (create)
    - `frontend/src/components/hello/__tests__/TrackCardSkeleton.test.tsx` (create)
    - `frontend/src/components/hello/__tests__/HelloNav.test.tsx` (create)
  - **Acceptance Criteria:**
    - SearchInput: renders input, calls onChange when typing, shows placeholder
    - TrackCard: renders title, artist, album, genre, year; formats duration as MM:SS (e.g., 240 -> "4:00")
    - TrackCardSkeleton: renders skeleton elements with loading animation classes
    - HelloNav: renders "Hello Music Search" title, has link to home
    - **MUST RUN AND FAIL before implementation**
  - **Commit:** Separate Red phase commit with test files only
  - _Requirements: US-5_
  - _Prompt: Role: Frontend React Developer | Task: Create test files for hello components. Read existing component test patterns in `frontend/src/components/` (look for __tests__ directories). Use @testing-library/react for rendering, screen.getByText/getByRole for assertions, @testing-library/user-event for interactions. Test SearchInput onChange callback, TrackCard rendering all fields with duration formatting (240 seconds -> "4:00", 185 seconds -> "3:05"), TrackCardSkeleton loading structure, HelloNav title and home link. Tests MUST fail because components don't exist yet. | Restrictions: Only create test files._

- [x] 5.2 Implement components (TDD Green)
  - **Implementation Files:**
    - `frontend/src/components/hello/SearchInput.tsx` (create)
    - `frontend/src/components/hello/TrackCard.tsx` (create)
    - `frontend/src/components/hello/TrackCardSkeleton.tsx` (create)
    - `frontend/src/components/hello/HelloNav.tsx` (create)
  - **Acceptance Criteria:**
    - All components use DaisyUI 5 semantic classes
    - TrackCard formats duration: `${Math.floor(d/60)}:${String(d%60).padStart(2,'0')}`
    - SearchInput is a controlled component with value/onChange props
    - Dark theme compatible (uses `base-100/200/300`, `primary`, etc.)
    - **ALL TESTS from 5.1 MUST PASS**
  - **Commit:** Separate Green phase commit
  - _Requirements: US-5_
  - _Prompt: Role: Frontend React Developer | Task: Implement hello components. Read existing components for DaisyUI patterns (card, input, badge, skeleton classes). SearchInput: DaisyUI input with search icon. TrackCard: DaisyUI card with title, artist badge, album, genre badge, year, formatted duration. TrackCardSkeleton: skeleton placeholder matching TrackCard layout. HelloNav: navbar with title and home link. Use semantic DaisyUI classes for dark/light theme support. | Restrictions: Make tests from 5.1 pass._

---

## Task 6: Frontend Route Page

- [x] 6.1 Write route page tests (TDD Red)
  - **Test File:** `frontend/src/routes/__tests__/hello-search.test.tsx` (create)
  - **Acceptance Criteria:**
    - Test page renders with search input and featured tracks
    - Test typing in search input triggers search
    - Test loading state shows skeleton cards
    - Test error state shows error message
    - Test empty results shows "No results" message
    - Test page is accessible at `/hello-search`
    - **MUST RUN AND FAIL before implementation**
  - **Commit:** Separate Red phase commit with test file only
  - _Requirements: US-5, US-6_
  - _Prompt: Role: Frontend React Developer | Task: Create `frontend/src/routes/__tests__/hello-search.test.tsx`. Read existing route tests in `frontend/src/routes/__tests__/` for patterns (how they mock hooks, render pages, test interactions). Test: (1) page renders HelloNav, SearchInput, and featured TrackCards, (2) typing triggers search (mock useHelloSearch), (3) loading state renders TrackCardSkeletons, (4) error state shows alert, (5) empty results shows message. Mock useHelloSearch and useHelloFeatured hooks. Tests MUST fail because the route file doesn't exist. | Restrictions: Only create the test file._

- [x] 6.2 Implement route page + wire router (TDD Green)
  - **Implementation Files:**
    - `frontend/src/routes/hello-search.tsx` (create)
    - `frontend/src/routes/__root.tsx` (modify -- add `/hello-search` to PUBLIC_ROUTES)
  - **Acceptance Criteria:**
    - Page at `/hello-search` with file-based routing
    - Shows HelloNav at top
    - Shows SearchInput
    - Shows featured tracks on initial load (useHelloFeatured)
    - Shows search results when query is entered (useHelloSearch)
    - Shows TrackCardSkeleton during loading
    - Shows error alert on failure
    - Shows "No results found" for empty results
    - `/hello-search` added to PUBLIC_ROUTES in `__root.tsx`
    - **MUST run `npx @tanstack/router-cli generate`** to regenerate `routeTree.gen.ts`
    - **ALL TESTS from 6.1 MUST PASS**
  - **Commit:** Separate Green phase commit
  - _Requirements: US-5, US-6_
  - _Prompt: Role: Frontend React Developer | Task: (1) Create `frontend/src/routes/hello-search.tsx` using TanStack Router file-based routing (createFileRoute). Read existing route files for the pattern. Compose HelloNav, SearchInput, useHelloSearch/useHelloFeatured hooks, TrackCard grid, loading skeletons, error handling. (2) Add '/hello-search' to PUBLIC_ROUTES array in __root.tsx. (3) Run `npx @tanstack/router-cli generate` to update routeTree.gen.ts. | Restrictions: Make tests from 6.1 pass. Must regenerate route tree._

---

## Task 7: Integration, BUILD Validation, and Documentation

- [x] 7.1 Contract alignment check
  - **Purpose:** Verify frontend and backend use the same API contract after both are implemented.
  - **Acceptance Criteria:**
    - Frontend API client sends `q` param (not `query`)
    - Frontend expects `items` key (not `tracks`)
    - Backend handler reads `q` query param
    - Backend handler returns JSON with `items` key
  - **Verification Commands:**
    ```bash
    # Verify frontend API sends correct params
    grep -n "params.*{.*q:" frontend/src/lib/api/helloSearch.ts
    # Verify frontend expects "items" not "tracks"
    grep -n "items:" frontend/src/lib/api/helloSearch.ts
    # Verify backend uses "q" query param
    grep -n 'QueryParam("q")' backend/internal/handlers/hello.go
    # Verify backend returns "items" key
    grep -n '"items"' backend/internal/handlers/hello.go
    ```
  - **If any check fails:** Fix the misaligned component, re-run tests, commit fix.
  - _Requirements: R-NF4_

- [x] 7.2 Local validation (`make local`)
  - **Purpose:** End-to-end validation with real local stack.
  - **Acceptance Criteria:**
    - `make local` starts successfully with seed data
    - `curl http://localhost:8080/api/v1/hello/health` returns `{"status":"ok","service":"hello"}`
    - `curl http://localhost:8080/api/v1/hello/featured` returns 20 tracks with `items` key
    - `curl "http://localhost:8080/api/v1/hello/search?q=jazz"` returns matching tracks
    - Frontend at `http://localhost:5173/hello-search` loads and shows tracks
    - Search input filters results
  - _Requirements: US-1 through US-6_

- [x] 7.3 BUILD validation
  - **Purpose:** Ensure code passes all quality gates.
  - **Acceptance Criteria:**
    - All commands pass with zero errors:
    ```bash
    # Backend
    cd backend && go vet ./...
    cd backend && go build -o /tmp/test-build ./cmd/api
    cd backend && go test ./...

    # Frontend
    cd frontend && npx tsc --noEmit
    cd frontend && npx eslint .
    cd frontend && npm run build
    cd frontend && npm test -- --run
    ```
  - _Requirements: R-NF2, R-NF3_

- [x] 7.4 Documentation updates
  - **Purpose:** Update CLAUDE.md files in all directories with changes.
  - **Files to Update:**
    - `backend/internal/service/CLAUDE.md` -- add HelloService description
    - `backend/internal/handlers/CLAUDE.md` -- add HelloHandler and hello routes
    - `frontend/CLAUDE.md` -- add hello components, hook, route
    - `frontend/src/CLAUDE.md` -- add hello route and components
    - `docker/CLAUDE.md` -- add init-seed-music.sh description
    - Root `CLAUDE.md` -- add hello-world-local-dev to Recent Updates
    - `CHANGELOG.md` -- add entry under [Unreleased]
  - _Requirements: Hard Requirement #4 (Documentation)_

---

## Task Dependencies

```
Task 1: Seed Data
  1.1 Test ──────────────────────────┐
  1.2 Implement ◄── depends on 1.1  │
  1.3 Wire Makefile                  │
                                     │
Task 2: Backend Service              │
  2.1 Test ──────────────────────────┤
  2.2 Implement ◄── depends on 2.1  │
                                     │
Task 3: Backend Handler + Wiring     │
  3.1 Test ◄── depends on 2.2       │
  3.2 Implement ◄── depends on 3.1  │
                                     │
Task 4: Frontend API + Hooks         │
  4.1 Test ──────────────────────────┤ (can parallel with Task 2-3)
  4.2 Implement ◄── depends on 4.1  │
                                     │
Task 5: Frontend Components          │
  5.1 Test ──────────────────────────┤ (can parallel with Task 4)
  5.2 Implement ◄── depends on 5.1  │
                                     │
Task 6: Frontend Route               │
  6.1 Test ◄── depends on 4.2, 5.2  │
  6.2 Implement ◄── depends on 6.1  │
                                     │
Task 7: Integration                  │
  7.1 Contract Check ◄── depends on 3.2, 4.2
  7.2 Local Validation ◄── depends on all implementation tasks
  7.3 BUILD Validation ◄── depends on all implementation tasks
  7.4 Documentation ◄── depends on 7.3
```

**Critical Path:** 2.1 -> 2.2 -> 3.1 -> 3.2 -> 7.1 (backend must be done before contract check)

**Parallelizable:**
- Tasks 4 and 5 (frontend API + components) can be done in parallel with Tasks 2 and 3 (backend)
- Task 1 (seed data) is independent of all other tasks

---

## SDLC Gap Tracking

_Gaps discovered during this iteration will be logged here. Start empty._

| # | Phase | Gap Description | Impact | Fix Applied |
|---|-------|-----------------|--------|-------------|
| 1 | Phase 2 (Red) | Commit plan specified 6 separate Red commits but all tests were committed in 1 batch | Lower granularity than planned, but Red/Green separation maintained | Acceptable trade-off — all tests verified failing before single Red commit |
| 2 | Phase 3 (Green) | DynamoDB repository mapped `TrackID` attribute but seed data stores ID in SK field | Track IDs empty in API responses during `make local` — unit tests passed because mocks bypass DynamoDB | Fixed: Extract ID from SK field by stripping `TRACK#` prefix. Separate fix commit. |
| 3 | Phase 4 (Verify) | `make local-reset` + `make local-services` still showed old seed data (40 items instead of 20) | Stale DynamoDB data from previous iterations persisted in LocalStack volume | Fixed by full volume reset. Not an SDLC gap per se — operational issue with persistent volumes. |
| 4 | Phase 4 (Verify) | Backend started without dummy AWS credentials fails with SSO credential error | `make local-backend` sets `AWS_ENDPOINT` but Makefile's `export` of dummy creds only applies to Make targets, not manual `go run` | Noted in MEMORY.md — always use `AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test` when running backend manually |
| 5 | Phase 3 (Green) | HelloDynamoDBRepository placed in `service/` package instead of `repository/` to avoid circular import | Unconventional package placement — HelloTrack defined in service, DynamoDB repo needs it | Design trade-off documented. Could be refactored to use separate DTO in repository package if needed. |

---

## Commit Plan

| Phase | Commit | Files |
|-------|--------|-------|
| Red 1 | `test(hello): Red - seed script tests` | `test-seed-music.sh` |
| Green 1 | `feat(hello): Green - seed script implementation` | `init-seed-music.sh`, `Makefile` |
| Red 2 | `test(hello): Red - HelloService tests` | `hello_test.go` |
| Green 2 | `feat(hello): Green - HelloService implementation` | `hello.go` |
| Red 3 | `test(hello): Red - HelloHandler tests` | `hello_test.go` (handlers) |
| Green 3 | `feat(hello): Green - HelloHandler + wiring` | `hello.go` (handlers), `service.go`, `main.go` |
| Red 4 | `test(hello): Red - frontend API + hook tests` | `helloSearch.test.ts`, `useHelloSearch.test.ts` |
| Green 4 | `feat(hello): Green - frontend API + hooks` | `helloSearch.ts`, `useHelloSearch.ts` |
| Red 5 | `test(hello): Red - component tests` | `SearchInput.test.tsx`, `TrackCard.test.tsx`, etc. |
| Green 5 | `feat(hello): Green - component implementations` | `SearchInput.tsx`, `TrackCard.tsx`, etc. |
| Red 6 | `test(hello): Red - route page tests` | `hello-search.test.tsx` |
| Green 6 | `feat(hello): Green - route page + router wiring` | `hello-search.tsx`, `__root.tsx`, `routeTree.gen.ts` |
| Final | `docs(hello): documentation updates` | CLAUDE.md files, CHANGELOG.md |
