# Requirements: Hello World Local Dev Validation

## Introduction

This spec validates the full local development stack by building a minimal but visually polished media search experience. It proves the Go API (Echo v4) serves data locally, the React frontend (TanStack Router/Query, Zustand, DaisyUI 5) consumes it, and the existing `make local` tooling works end-to-end. The UI is designed as a **future-facing foundation** for the music search engine - not throwaway scaffolding.

## Alignment with Product Vision

This feature directly supports the Personal Music Search Engine by:
- Validating the local development loop (LocalStack + Go API + React frontend)
- Establishing DaisyUI 5 design patterns for media search/discovery
- Proving TDD workflow compliance end-to-end (spec -> test -> code -> verify)
- Creating reusable search UI components that will be used in production

## Requirements

### Requirement 1: Local Go API Health + Search Endpoint

**User Story:** As a developer, I want a Go API running locally that serves mock music search results, so that I can validate the backend stack and API contract.

#### Acceptance Criteria

1. WHEN the API server starts locally THEN the system SHALL respond to `GET /api/v1/health` with `200 OK` and `{"status": "ok"}`
2. WHEN a client sends `GET /api/v1/search?q={term}` THEN the system SHALL return a JSON array of track results matching the query
3. WHEN a client sends `GET /api/v1/search?q={term}` with no matches THEN the system SHALL return an empty array `[]` with `200 OK`
4. WHEN the search query parameter is empty THEN the system SHALL return `400 Bad Request` with `{"error": "query parameter 'q' is required"}`
5. WHEN the API starts THEN it SHALL use the existing Echo v4 server pattern from `backend/cmd/api/main.go` (HTTP mode, not Lambda)

### Requirement 2: Mock Data Seeding

**User Story:** As a developer, I want the local environment to have realistic mock music data, so that search results look meaningful during development.

#### Acceptance Criteria

1. WHEN LocalStack starts via `make local-services` THEN the system SHALL seed DynamoDB with at least 20 mock tracks across 5 artists
2. IF mock data already exists in DynamoDB THEN the seed script SHALL skip insertion (idempotent)
3. WHEN mock tracks are seeded THEN each track SHALL have: title, artist name, album name, duration (seconds), genre, year, and cover art URL (placeholder)

### Requirement 3: React Frontend with Media Search UI

**User Story:** As a user, I want a visually polished search interface to find music, so that I can validate the frontend stack and see a future-facing design.

#### Acceptance Criteria

1. WHEN the React app loads at `http://localhost:5173` THEN the system SHALL display a hero section with a prominent search bar (DaisyUI hero + input components)
2. WHEN a user types in the search bar THEN the system SHALL debounce input (300ms) and call `GET /api/v1/search?q={term}` via TanStack Query
3. WHEN search results are loading THEN the system SHALL display skeleton placeholders (DaisyUI skeleton component)
4. WHEN search results return THEN the system SHALL display results as media cards in a responsive grid (DaisyUI card components with album art, track title, artist, duration, genre badge)
5. WHEN there are no search results THEN the system SHALL display an empty state with "No tracks found" message
6. WHEN the search bar is empty THEN the system SHALL display a curated "Featured Tracks" section showing recent/popular tracks

### Requirement 4: Frontend Navigation Shell

**User Story:** As a user, I want consistent navigation that feels like a real music app, so that the UI is future-ready for additional pages.

#### Acceptance Criteria

1. WHEN the app loads THEN the system SHALL display a top navbar with app logo/name and search input (DaisyUI navbar with search)
2. WHEN on mobile viewport THEN the system SHALL display a bottom dock with navigation icons (DaisyUI dock component)
3. WHEN on desktop viewport THEN the system SHALL display a sidebar with navigation links (DaisyUI drawer component)
4. WHEN the user navigates THEN the system SHALL use TanStack Router with the following routes:
   - `/` - Home/search page (hero + search)
   - `/search?q={term}` - Search results page
   - `/about` - About page (placeholder)

### Requirement 5: Theme and Styling

**User Story:** As a developer, I want a cohesive dark-themed music app aesthetic, so that the UI looks professional from day one.

#### Acceptance Criteria

1. WHEN the app loads THEN the system SHALL use a dark theme by default (DaisyUI `dark` or custom theme)
2. WHEN displaying the UI THEN all components SHALL use DaisyUI 5 semantic color classes (`primary`, `secondary`, `accent`, `base-100/200/300`)
3. WHEN the user toggles theme THEN the system SHALL persist the preference in localStorage (using existing `themeStore`)

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Each file should have a single, well-defined purpose
- **Modular Design**: Search components, API clients, and hooks should be isolated and reusable
- **Dependency Management**: Frontend search components must not import backend code
- **Clear Interfaces**: API contract defined in design.md with TypeScript types matching Go structs
- **Leverage Existing Code**: MUST use existing patterns from `backend/cmd/api/main.go`, `frontend/src/lib/api/client.ts`, `frontend/src/hooks/`, and `frontend/src/components/layout/`

### Performance
- Search debounce: 300ms minimum to avoid excessive API calls
- Search response: < 200ms for local DynamoDB queries
- Frontend bundle: No new dependencies beyond what's already in `package.json`
- Skeleton loading states must appear within 100ms of search initiation

### Security
- No authentication required for hello-world endpoints (search is public)
- API input validation on search query (sanitize, length limits)
- No sensitive data in mock seed data

### Reliability
- Health endpoint must always return 200 if server is running
- Search must gracefully handle DynamoDB connection failures with user-friendly error message
- Frontend must handle API unavailability with retry logic (TanStack Query built-in)

### Usability
- Search bar auto-focused on page load
- Responsive layout works from 320px to 2560px viewport
- Keyboard navigation: Enter to search, Escape to clear
- Results show track count ("Found N tracks")
