# Tasks Document: Hello World Local Development Validation

**Spec:** `hello-world-local-dev`
**TDD Enforcement:** Every task group separates tests (Red) from implementation (Green). Tests MUST run and FAIL before implementation begins. Separate commits required for Red and Green phases.

---

## Task 1: Seed Data

### 1.1 Write seed script validation test (TDD Red)

- **Test File:** `docker/localstack-init/test-seed-music.sh` (create)
- **Acceptance Criteria:**
  - Script verifies that after running `init-seed-music.sh`, exactly 20 items exist in DynamoDB with PK=`USER#seed-user`
  - Script verifies idempotency: running seed twice produces no errors and still has exactly 20 items
  - Script verifies at least one track has expected fields (Title, Artist, Album, Genre, Year, Duration)
  - **MUST RUN AND FAIL** before seed script is implemented (because `init-seed-music.sh` does not exist yet)
- **_Leverage:_** Existing `init-aws.sh` and `init-cognito.sh` patterns for AWS CLI usage
- **_Requirements:_** R1 (Seed Data)

### 1.2 Implement seed script (TDD Green)

- **Implementation:** `docker/localstack-init/init-seed-music.sh` (create)
- **Acceptance Criteria:**
  - Inserts 20 tracks across 5 artists into `MusicLibrary` table under `USER#seed-user`
  - Each track has: PK, SK, TrackID, UserID, Title, Artist, Album, Genre, Year, Duration, CreatedAt, UpdatedAt, Visibility
  - Uses `--condition-expression "attribute_not_exists(PK)"` for idempotency
  - Uses `AWS_ENDPOINT` environment variable (default `http://localhost:4566`)
  - **ALL TESTS FROM 1.1 MUST PASS**
- **_Leverage:_** `docker/localstack-init/init-aws.sh` script pattern
- **_Requirements:_** R1 (Seed Data)

### 1.3 Wire seed script into Makefile

- **Implementation:** `Makefile` (modify)
- **Acceptance Criteria:**
  - `make local-services` calls `init-seed-music.sh` after `init-cognito.sh`
  - Seed script runs as part of normal `make local` flow
  - Script has executable permissions (`chmod +x`)
- **_Leverage:_** Existing Makefile `local-services` target
- **_Requirements:_** R1 (Seed Data)

---

## Task 2: Backend Service

### 2.1 Write HelloService unit tests (TDD Red)

- **Test File:** `backend/internal/service/hello_test.go` (create)
- **Acceptance Criteria:**
  - Tests use a mock `HelloRepository` (defined in test file)
  - Test: `TestHelloService_SearchTracks_MatchesTitle` -- searches by title, expects matching tracks
  - Test: `TestHelloService_SearchTracks_CaseInsensitive` -- verifies case-insensitive matching
  - Test: `TestHelloService_SearchTracks_MatchesArtist` -- searches by artist name
  - Test: `TestHelloService_SearchTracks_MatchesAlbum` -- searches by album name
  - Test: `TestHelloService_SearchTracks_MatchesGenre` -- searches by genre
  - Test: `TestHelloService_SearchTracks_EmptyQuery` -- empty query returns empty slice
  - Test: `TestHelloService_SearchTracks_NoMatch` -- non-matching query returns empty slice
  - Test: `TestHelloService_GetFeaturedTracks_ReturnsAll` -- returns all tracks up to limit
  - Test: `TestHelloService_GetFeaturedTracks_WithLimit` -- respects limit parameter
  - Test: `TestHelloService_GetFeaturedTracks_DefaultLimit` -- default limit is 20
  - **MUST RUN AND FAIL** because `hello.go` does not exist yet
  - **MUST commit test file separately** before implementation
- **_Leverage:_** Mock pattern from `backend/internal/service/search_test.go`, `testify/assert`
- **_Requirements:_** R2 (Backend Search API), R3 (Backend Featured API)

### 2.2 Implement HelloService (TDD Green)

- **Implementation:** `backend/internal/service/hello.go` (create)
- **Acceptance Criteria:**
  - `HelloTrack` struct with JSON tags
  - `HelloRepository` interface with `GetTracksByUser` method
  - `HelloService` struct with `repo` and `userID` fields
  - `NewHelloService(repo HelloRepository) *HelloService` constructor
  - `SearchTracks(ctx, query)` -- fetches all tracks, filters in-memory (case-insensitive across title, artist, album, genre)
  - `GetFeaturedTracks(ctx, limit)` -- returns all tracks up to limit (default 20)
  - `HelloRepoAdapter` struct that wraps `repository.Repository`
  - `NewHelloRepoAdapter(repo repository.Repository) *HelloRepoAdapter`
  - `GetTracksByUser(ctx, userID)` -- calls `repo.ListTracks` and maps `models.Track` to `HelloTrack`
  - **ALL TESTS FROM 2.1 MUST PASS**
  - **MUST commit implementation separately** after Red phase commit
- **_Leverage:_** `service/service.go` interface pattern, `repository.Repository` ListTracks method
- **_Requirements:_** R2 (Backend Search API), R3 (Backend Featured API)

---

## Task 3: Backend Handler

### 3.1 Write HelloHandler unit tests (TDD Red)

- **Test File:** `backend/internal/handlers/hello_test.go` (create)
- **Acceptance Criteria:**
  - Tests use `httptest.NewRecorder` and `echo.New()` for handler testing
  - Tests mock `HelloService` by providing a mock `HelloRepository`
  - Test: `TestHelloHealth_ReturnsOK` -- GET /api/v1/hello/health returns 200 with `{"status":"ok","service":"hello"}`
  - Test: `TestHelloSearch_WithQuery` -- GET /api/v1/hello/search?q=neon returns matching tracks
  - Test: `TestHelloSearch_EmptyQuery` -- GET /api/v1/hello/search (no q) returns empty items
  - Test: `TestHelloFeatured_ReturnsAll` -- GET /api/v1/hello/featured returns all tracks
  - Test: `TestHelloFeatured_WithLimit` -- GET /api/v1/hello/featured?limit=5 returns at most 5
  - Test: `TestRegisterHelloRoutes_RoutesExist` -- verifies all 3 routes are registered on Echo
  - **MUST RUN AND FAIL** because `handlers/hello.go` does not exist yet
  - **MUST commit test file separately** before implementation
- **_Leverage:_** Handler test patterns from existing `handlers/*_test.go` files, `echo` test utilities
- **_Requirements:_** R2 (Backend Search API), R3 (Backend Featured API), R4 (Backend Health)

### 3.2 Implement HelloHandler + WIRE into application (TDD Green)

- **Implementation Files:**
  - `backend/internal/handlers/hello.go` (create)
  - `backend/internal/service/service.go` (modify)
  - `backend/cmd/api/main.go` (modify)
- **Acceptance Criteria:**
  - **Handler Implementation:**
    - `HelloHandler` struct with `service *service.HelloService`
    - `NewHelloHandler(svc *service.HelloService) *HelloHandler`
    - `HelloHealth(c echo.Context) error` -- returns `{"status":"ok","service":"hello"}`
    - `HelloSearch(c echo.Context) error` -- reads `q` query param, calls service, returns `{"items":[...],"total":N}`
    - `HelloFeatured(c echo.Context) error` -- reads optional `limit` query param, calls service, returns `{"items":[...],"total":N}`
    - `RegisterHelloRoutes(e *echo.Echo, h *HelloHandler)` -- registers GET routes under `/api/v1/hello`
  - **WIRING (CRITICAL -- unit tests CANNOT catch missing wiring):**
    - Add `Hello *HelloService` field to `Services` struct in `service/service.go`
    - Initialize `Hello: NewHelloService(NewHelloRepoAdapter(repo))` in `NewServices()` in `service/service.go`
    - Create `helloHandler := handlers.NewHelloHandler(services.Hello)` in `setupEcho()` in `main.go`
    - Call `handlers.RegisterHelloRoutes(e, helloHandler)` in `setupEcho()` in `main.go` (AFTER `h.RegisterRoutes(e)`)
  - **ALL TESTS FROM 3.1 MUST PASS**
  - **Application MUST compile:** `cd backend && go build ./cmd/api`
  - **MUST commit implementation separately** after Red phase commit
- **_Leverage:_** `handlers/handlers.go` RegisterRoutes pattern, `service/service.go` Services struct, `cmd/api/main.go` setupEcho() wiring pattern, `.claude/docs/wiring-checklist.md`
- **_Requirements:_** R2 (Backend Search API), R3 (Backend Featured API), R4 (Backend Health), R6 (No Auth)

---

## Task 4: Frontend API Client + Hooks

### 4.1 Write API client and hooks tests (TDD Red)

- **Test Files:**
  - `frontend/src/lib/api/__tests__/helloSearch.test.ts` (create)
  - `frontend/src/hooks/__tests__/useHelloSearch.test.ts` (create)
- **Acceptance Criteria:**
  - **API Client Tests (`helloSearch.test.ts`):**
    - Test: `searchHelloTracks` calls GET `/v1/hello/search` with query param
    - Test: `getFeaturedTracks` calls GET `/v1/hello/featured`
    - Test: `getFeaturedTracks` passes limit parameter when provided
    - Test: `getHelloHealth` calls GET `/v1/hello/health`
    - Tests mock `apiClient` using `vi.mock`
  - **Hook Tests (`useHelloSearch.test.ts`):**
    - Test: `useHelloSearch` is disabled when query is empty
    - Test: `useHelloSearch` fetches when query is non-empty
    - Test: `useFeaturedTracks` fetches on mount
    - Tests use `@tanstack/react-query` test utilities with `QueryClientProvider` wrapper
  - **MUST RUN AND FAIL** because source files do not exist yet
  - **MUST commit test files separately** before implementation
- **_Leverage:_** `frontend/src/lib/api/__tests__/` existing test patterns, `vi.mock` for axios, `@testing-library/react` renderHook
- **_Requirements:_** R2, R3, R5 (Frontend Search Page)

### 4.2 Implement API client and hooks (TDD Green)

- **Implementation Files:**
  - `frontend/src/lib/api/helloSearch.ts` (create)
  - `frontend/src/hooks/useHelloSearch.ts` (create)
- **Acceptance Criteria:**
  - **API Client (`helloSearch.ts`):**
    - `HelloTrack` interface with id, title, artist, album, genre, year, duration
    - `HelloSearchResponse` interface with items and total
    - `searchHelloTracks(query: string): Promise<HelloSearchResponse>`
    - `getFeaturedTracks(limit?: number): Promise<HelloSearchResponse>`
    - `getHelloHealth(): Promise<{ status: string; service: string }>`
    - Uses `apiClient` from `./client`
  - **Hooks (`useHelloSearch.ts`):**
    - `helloKeys` query key factory
    - `useHelloSearch(query: string)` -- enabled only when query.length > 0
    - `useFeaturedTracks(limit?: number)` -- always enabled
  - **ALL TESTS FROM 4.1 MUST PASS**
  - **MUST commit implementation separately** after Red phase commit
- **_Leverage:_** `frontend/src/lib/api/client.ts` apiClient, `frontend/src/hooks/useTracks.ts` queryKey factory pattern
- **_Requirements:_** R2, R3, R5

---

## Task 5: Frontend Components

### 5.1 Write component tests (TDD Red)

- **Test Files:**
  - `frontend/src/components/hello/__tests__/SearchInput.test.tsx` (create)
  - `frontend/src/components/hello/__tests__/TrackCard.test.tsx` (create)
  - `frontend/src/components/hello/__tests__/TrackCardSkeleton.test.tsx` (create)
- **Acceptance Criteria:**
  - **SearchInput Tests:**
    - Test: renders input with placeholder text
    - Test: calls onChange when user types
    - Test: displays current value
  - **TrackCard Tests:**
    - Test: renders track title
    - Test: renders artist name
    - Test: renders album name
    - Test: renders genre as badge
    - Test: renders year
    - Test: formats duration as MM:SS (e.g., 234 seconds -> "3:54")
  - **TrackCardSkeleton Tests:**
    - Test: renders skeleton loading placeholders
    - Test: has correct number of skeleton elements
  - **MUST RUN AND FAIL** because component files do not exist yet
  - **MUST commit test files separately** before implementation
- **_Leverage:_** `@testing-library/react`, existing component test patterns in `frontend/src/components/`
- **_Requirements:_** R5 (Frontend Search Page)

### 5.2 Implement components (TDD Green)

- **Implementation Files:**
  - `frontend/src/components/hello/SearchInput.tsx` (create)
  - `frontend/src/components/hello/TrackCard.tsx` (create)
  - `frontend/src/components/hello/TrackCardSkeleton.tsx` (create)
  - `frontend/src/components/hello/HelloNav.tsx` (create)
- **Acceptance Criteria:**
  - **SearchInput:** Controlled input with DaisyUI `input input-bordered` classes, magnifying glass icon
  - **TrackCard:** DaisyUI `card bg-base-200`, displays title/artist/album/genre/year/duration
  - **TrackCardSkeleton:** DaisyUI `skeleton` classes matching TrackCard dimensions
  - **HelloNav:** Simple link component for navigation
  - All components use DaisyUI 5 semantic classes (`base-100`, `base-200`, `primary`, etc.)
  - **ALL TESTS FROM 5.1 MUST PASS**
  - **MUST commit implementation separately** after Red phase commit
- **_Leverage:_** DaisyUI 5 card patterns, existing `frontend/src/components/library/` styling
- **_Requirements:_** R5 (Frontend Search Page)

---

## Task 6: Frontend Route

### 6.1 Write route page tests (TDD Red)

- **Test File:** `frontend/src/routes/__tests__/hello-search.test.tsx` (create)
- **Acceptance Criteria:**
  - Test: renders hero section with title text
  - Test: shows featured tracks when loaded (mock `useFeaturedTracks`)
  - Test: shows loading skeletons while fetching
  - Test: shows error message when API fails
  - Test: performs search when query is entered (mock `useHelloSearch`)
  - Test: shows search results instead of featured tracks when searching
  - Tests mock hooks (`useHelloSearch`, `useFeaturedTracks`) using `vi.mock`
  - **MUST RUN AND FAIL** because route file does not exist yet
  - **MUST commit test file separately** before implementation
- **_Leverage:_** `frontend/src/routes/__tests__/` existing route test patterns, `@testing-library/react`
- **_Requirements:_** R5 (Frontend Search Page)

### 6.2 Implement route page + wire into router (TDD Green)

- **Implementation Files:**
  - `frontend/src/routes/hello-search.tsx` (create)
  - `frontend/src/routes/__root.tsx` (modify)
- **Acceptance Criteria:**
  - **Route Page (`hello-search.tsx`):**
    - Uses `createFileRoute('/hello-search')` for TanStack Router file-based routing
    - Hero section with title ("Hello Music Search") and subtitle
    - SearchInput component with 300ms debounce
    - Featured tracks grid (responsive: 1/2/3 columns)
    - Search results grid when query is active
    - TrackCardSkeleton grid during loading
    - Error state with retry button
    - DaisyUI 5 styling consistent with app theme
  - **Router Wiring (`__root.tsx`):**
    - Add `/hello-search` to `PUBLIC_ROUTES` array so unauthenticated access works
  - **ALL TESTS FROM 6.1 MUST PASS**
  - **MUST commit implementation separately** after Red phase commit
- **_Leverage:_** `frontend/src/routes/search.tsx` page structure, `__root.tsx` PUBLIC_ROUTES pattern, TanStack Router `createFileRoute`
- **_Requirements:_** R5 (Frontend Search Page), R6 (No Auth Required)

---

## Task 7: Integration and Documentation

### 7.1 End-to-end validation with `make local`

- **Validation Steps:**
  1. Run `make local` -- confirm LocalStack + backend + frontend start
  2. Verify seed data: `aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name MusicLibrary --filter-expression "begins_with(PK, :pk)" --expression-attribute-values '{":pk":{"S":"USER#seed-user"}}' --select COUNT` returns 20
  3. Verify health: `curl http://localhost:8080/api/v1/hello/health` returns `{"status":"ok","service":"hello"}`
  4. Verify featured: `curl http://localhost:8080/api/v1/hello/featured` returns 20 tracks
  5. Verify search: `curl "http://localhost:8080/api/v1/hello/search?q=aurora"` returns matching tracks
  6. Verify frontend: Open `http://localhost:5173/hello-search` in browser, confirm hero + featured tracks render
  7. Verify search UI: Type a query, confirm search results appear
  8. Run `make local-stop` -- confirm clean shutdown
- **Acceptance Criteria:**
  - All 8 validation steps pass without errors
- **_Requirements:_** All (R1-R6)

### 7.2 Update CLAUDE.md files

- **Files to Update:**
  - `backend/internal/service/CLAUDE.md` -- add HelloService section
  - `backend/internal/handlers/CLAUDE.md` -- add HelloHandler routes
  - `frontend/src/CLAUDE.md` -- add hello components, hooks, route
  - `docker/CLAUDE.md` -- add seed script documentation
- **Acceptance Criteria:**
  - Each CLAUDE.md documents the new files, functions, and patterns
  - Follows existing CLAUDE.md format in each directory
- **_Leverage:_** Existing CLAUDE.md files in each directory
- **_Requirements:_** Documentation standard from root CLAUDE.md

---

## SDLC Gap Tracking

| # | Phase | Gap Description | Impact | Proposed Fix | Status |
|---|-------|----------------|--------|--------------|--------|
| | | | | | |

_No gaps identified yet. Update this table if SDLC process issues are discovered during implementation._
