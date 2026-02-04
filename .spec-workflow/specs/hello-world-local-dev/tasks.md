# Tasks: Hello World Local Dev Validation

**Spec**: hello-world-local-dev
**Design**: [design.md](./design.md)
**Requirements**: [requirements.md](./requirements.md)

---

## Task 1: Mock Data Seed Script

- [ ] 1.1 Write unit tests for seed data validation (TDD Red)
  - **Unit Test**: `docker/localstack-init/test-seed-data.sh` (create)
  - **Acceptance Criteria**:
    - Test verifies DynamoDB table has 20 items after seeding
    - Test verifies each item has required fields (title, artist, album, genre, year, duration)
    - Test verifies 5 distinct artists exist
    - Test verifies seed script is idempotent (running twice doesn't duplicate)
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 2.1, 2.2, 2.3_

- [ ] 1.2 Implement seed script (TDD Green)
  - **Implementation**: `docker/localstack-init/init-seed-music.sh` (create)
  - **Acceptance Criteria**:
    - Seeds 20 mock tracks across 5 artists into DynamoDB MusicLibrary table
    - Uses `aws dynamodb put-item` with proper PK/SK pattern (`USER#seed-user` / `TRACK#<uuid>`)
    - Includes all fields: id, userId, title, artist, album, genre, year, duration, format, fileSize, s3Key, Visibility=public
    - Idempotent: checks if seed data exists before inserting
    - Integrates with existing docker-compose startup via Makefile
    - **ALL TESTS MUST PASS**
  - _Requirements: 2.1, 2.2, 2.3_

---

## Task 2: Backend Search Service

- [ ] 2.1 Write unit tests for HelloService (TDD Red)
  - **Unit Test**: `backend/internal/service/hello_test.go` (create)
  - **Acceptance Criteria**:
    - Test `SearchTracks` returns matching tracks for valid query
    - Test `SearchTracks` returns empty slice for no-match query
    - Test `SearchTracks` is case-insensitive
    - Test `SearchTracks` matches on title, artist, and genre
    - Test `GetFeaturedTracks` returns all seeded tracks
    - Mock the repository interface
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 1.2, 1.3_

- [ ] 2.2 Write unit tests for HelloHandler (TDD Red)
  - **Unit Test**: `backend/internal/handlers/hello_test.go` (create)
  - **Acceptance Criteria**:
    - Test `GET /api/v1/health` returns `200` with `{"status":"ok"}`
    - Test `GET /api/v1/hello/search?q=luna` returns `200` with tracks
    - Test `GET /api/v1/hello/search` (no query) returns `400`
    - Test `GET /api/v1/hello/search?q=` (empty query) returns `400`
    - Test `GET /api/v1/hello/featured` returns `200` with all tracks
    - Mock the service interface
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [ ] 2.3 Implement HelloService and HelloHandler (TDD Green)
  - **Implementation**:
    - `backend/internal/service/hello.go` (create)
    - `backend/internal/handlers/hello.go` (create)
    - `backend/internal/handlers/handlers.go` (modify - add route registration)
  - **Acceptance Criteria**:
    - `HelloService.SearchTracks(ctx, query, limit)` scans DynamoDB for seed user tracks, filters by title/artist/genre containing query (case-insensitive)
    - `HelloService.GetFeaturedTracks(ctx, limit)` returns all seed user tracks
    - `HelloHandler.Health` returns health check response
    - `HelloHandler.Search` validates query param, calls service, returns `HelloSearchResponse`
    - `HelloHandler.Featured` calls service, returns featured tracks
    - Routes registered at `/api/v1/health`, `/api/v1/hello/search`, `/api/v1/hello/featured`
    - **ALL TESTS MUST PASS**
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

---

## Task 3: Frontend API Client and Hook

- [ ] 3.1 Write unit tests for API client (TDD Red)
  - **Unit Test**: `frontend/src/lib/api/__tests__/helloSearch.test.ts` (create)
  - **Acceptance Criteria**:
    - Test `searchHello(query)` calls `GET /api/v1/hello/search?q={query}`
    - Test `getFeaturedTracks()` calls `GET /api/v1/hello/featured`
    - Test response types match `HelloSearchResponse` interface
    - Mock axios
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 3.2_

- [ ] 3.2 Write unit tests for useHelloSearch hook (TDD Red)
  - **Unit Test**: `frontend/src/hooks/__tests__/useHelloSearch.test.ts` (create)
  - **Acceptance Criteria**:
    - Test hook returns `isLoading: true` initially
    - Test hook returns track data on success
    - Test hook debounces queries (300ms)
    - Test hook returns `isError: true` on API failure
    - Test hook skips query when search term is empty
    - Use TanStack Query test utilities
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 3.2, 3.6_

- [ ] 3.3 Implement API client and hook (TDD Green)
  - **Implementation**:
    - `frontend/src/lib/api/helloSearch.ts` (create)
    - `frontend/src/hooks/useHelloSearch.ts` (create)
  - **Acceptance Criteria**:
    - `searchHello(query)` calls the hello search API endpoint
    - `getFeaturedTracks()` calls the hello featured endpoint
    - `useHelloSearch(query)` wraps searchHello with TanStack Query
    - `useFeaturedTracks()` wraps getFeaturedTracks with TanStack Query
    - Debounce implemented (300ms delay before API call)
    - **ALL TESTS MUST PASS**
  - _Requirements: 3.2, 3.6_

---

## Task 4: Frontend UI Components

- [ ] 4.1 Write unit tests for SearchInput component (TDD Red)
  - **Unit Test**: `frontend/src/components/hello/__tests__/SearchInput.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test renders input element with placeholder text
    - Test calls onChange when user types
    - Test calls onSubmit on Enter key press
    - Test clears input on Escape key press
    - Test autoFocus prop focuses input on mount
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 3.1_

- [ ] 4.2 Write unit tests for TrackCard component (TDD Red)
  - **Unit Test**: `frontend/src/components/hello/__tests__/TrackCard.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test renders track title, artist name, album name
    - Test renders genre as DaisyUI badge
    - Test renders formatted duration string
    - Test renders cover art image with alt text
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 3.4_

- [ ] 4.3 Write unit tests for TrackCardSkeleton component (TDD Red)
  - **Unit Test**: `frontend/src/components/hello/__tests__/TrackCardSkeleton.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test renders skeleton elements (DaisyUI skeleton class)
    - Test renders correct number of skeleton lines (image, title, artist, badge)
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 3.3_

- [ ] 4.4 Implement SearchInput, TrackCard, TrackCardSkeleton (TDD Green)
  - **Implementation**:
    - `frontend/src/components/hello/SearchInput.tsx` (create)
    - `frontend/src/components/hello/TrackCard.tsx` (create)
    - `frontend/src/components/hello/TrackCardSkeleton.tsx` (create)
  - **Acceptance Criteria**:
    - SearchInput: DaisyUI `input` + `join` + search icon, keyboard handling
    - TrackCard: DaisyUI `card` with image, title, artist, genre badge, duration
    - TrackCardSkeleton: DaisyUI `skeleton` placeholders matching card layout
    - All components use DaisyUI 5 semantic color classes
    - **ALL TESTS MUST PASS**
  - _Requirements: 3.1, 3.3, 3.4_

---

## Task 5: Frontend Route and Page Assembly

- [ ] 5.1 Write unit tests for SearchHeroPage route (TDD Red)
  - **Unit Test**: `frontend/src/routes/__tests__/hello-search.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test page renders hero section with "Music Search" heading
    - Test page renders search input
    - Test typing in search shows loading skeletons
    - Test search results display as track cards
    - Test empty search shows "Featured Tracks" section
    - Test no results shows "No tracks found" message
    - Test result count "Found N tracks" is displayed
    - Mock useHelloSearch and useFeaturedTracks hooks
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

- [ ] 5.2 Write unit tests for HelloNav component (TDD Red)
  - **Unit Test**: `frontend/src/components/hello/__tests__/HelloNav.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test renders navbar with app name
    - Test renders search input in navbar
    - Test renders theme toggle button
    - Test renders dock on mobile viewport (dock element present)
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 4.1, 4.2_

- [ ] 5.3 Write E2E tests for hello-search page (TDD Red)
  - **E2E Test**: `frontend/e2e/hello-search.spec.ts` (create)
  - **Acceptance Criteria**:
    - Test navigates to `/hello-search` route
    - Test page has hero with search input
    - Test search input is auto-focused
    - Test featured tracks are displayed on load
    - Test typing query shows search results after debounce
    - Test empty results show "No tracks found"
    - Test responsive: mobile viewport shows dock
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: All Requirement 3, 4_

- [ ] 5.4 Implement SearchHeroPage and HelloNav (TDD Green)
  - **Implementation**:
    - `frontend/src/routes/hello-search.tsx` (create)
    - `frontend/src/components/hello/HelloNav.tsx` (create)
  - **Acceptance Criteria**:
    - Route registered at `/hello-search` in TanStack Router
    - Hero section: DaisyUI `hero` with gradient background, heading, subtext, search input
    - Results grid: responsive 1-4 column grid with TrackCard components
    - Loading state: grid of TrackCardSkeleton components
    - Empty state: centered message "No tracks found for '{query}'"
    - Featured state: shows all tracks with "Featured Tracks" heading
    - Result count: "Found N tracks for '{query}'"
    - HelloNav: DaisyUI `navbar` with logo, search, theme toggle
    - Bottom dock: DaisyUI `dock` for mobile with Home/Search/About icons
    - **ALL TESTS MUST PASS**
  - _Requirements: All Requirement 3, 4, 5_

---

## Task 6: Theme and Styling Polish

- [ ] 6.1 Write unit tests for theme integration (TDD Red)
  - **Unit Test**: `frontend/src/components/hello/__tests__/ThemeToggle.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test dark theme is default
    - Test toggle switches between dark and light themes
    - Test theme persists after toggle (calls themeStore)
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: 5.1, 5.2, 5.3_

- [ ] 6.2 Implement theme integration (TDD Green)
  - **Implementation**:
    - Update `frontend/src/routes/hello-search.tsx` (modify - add theme classes)
    - Verify existing `themeStore` integration works with hello components
  - **Acceptance Criteria**:
    - Dark theme active by default on `/hello-search`
    - All components use semantic DaisyUI color classes
    - Theme toggle in navbar works
    - **ALL TESTS MUST PASS**
  - _Requirements: 5.1, 5.2, 5.3_

---

## Task 7: Integration and End-to-End Validation

- [ ] 7.1 Validate full stack with `make local`
  - **Validation**: Manual + automated
  - **Acceptance Criteria**:
    - `make local` starts LocalStack, seeds data, starts backend, starts frontend
    - `curl http://localhost:8080/api/v1/health` returns `{"status":"ok"}`
    - `curl http://localhost:8080/api/v1/hello/search?q=luna` returns matching tracks
    - `http://localhost:5173/hello-search` loads with featured tracks
    - Searching "electronic" shows Luna Waves and DJ Phantom tracks
    - Mobile viewport shows bottom dock
    - Theme toggle works

- [ ] 7.2 Update CLAUDE.md for all changed directories
  - **Documentation**: Create/update CLAUDE.md files
  - **Directories requiring CLAUDE.md**:
    - `frontend/src/components/hello/` (create)
    - `frontend/src/hooks/` (update with useHelloSearch)
    - `frontend/src/lib/api/` (update with helloSearch.ts)
    - `docker/localstack-init/` (update with seed script)
    - `backend/internal/handlers/` (update with hello.go)
    - `backend/internal/service/` (update with hello.go)
  - **Acceptance Criteria**:
    - Each CLAUDE.md has: Overview, Files, Key Functions, Dependencies
    - Zero missing CLAUDE.md for changed directories

---

## SDLC Gap Tracking

### Gaps Found During This Spec

| # | Gap | Where | Impact | Fix Applied |
|---|-----|-------|--------|-------------|
| 1 | SDLC Phase 1 doesn't mention starting spec-workflow dashboard | `sdlc.md` Phase 1 | Approvals can't be completed without dashboard | ✅ Added Prerequisites section with VS Code extension (preferred) and dashboard (alternative) |
| 2 | No clear instructions for approval dashboard vs VS Code extension | `sdlc.md` Phase 1 | Users don't know which approval method to use | ✅ Added "Verbal approval is NOT accepted" and preference guidance |
| 3 | No steering documents exist | `.spec-workflow/steering/` | Design template references product.md, tech.md, structure.md that don't exist | ✅ Updated requirements + design templates to reference root CLAUDE.md as fallback |
| 4 | Spec workflow requires sequential approval of each doc | `sdlc.md` Phase 1 | Blocks progress when reviewer isn't immediately available | ✅ Added step 7: "Create all 3 docs first, submit all for review, batch-approve" |
| 5 | `tasks.md` template doesn't follow TDD structure | `.spec-workflow/templates/tasks-template.md` | Template has code-first tasks, violates TDD rule | ✅ Rewrote template with TDD structure, MUST RUN AND FAIL, Task Structure Rules table |
| 6 | No `make local` validation step in SDLC | `sdlc.md` Phase 4 | No requirement to verify feature works end-to-end locally before PR | ✅ Added "Local Stack Validation" section + status tracking checklist item |
| 7 | SDLC doesn't specify when to use existing vs new endpoints | `sdlc.md` Phase 1 | Could lead to duplicate endpoints or unnecessary new code | ✅ Added step 3 "Code Reuse Analysis" + MUST include in design.md |
