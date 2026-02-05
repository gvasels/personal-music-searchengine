# Requirements Document: Hello World Local Development Validation

## Introduction

This feature creates a standalone hello-world music search page to validate the full local development stack works end-to-end. It seeds mock track data into LocalStack DynamoDB and provides unauthenticated REST endpoints plus a React search UI, confirming that LocalStack, the Go backend, and the React frontend all work together correctly via `make local`.

## Alignment with Product Vision

This feature validates the local development infrastructure established by the LocalStack Dev Environment spec (Epic 8). By exercising the full stack with a simple, isolated feature, developers gain confidence that:
- DynamoDB read/write operations work via LocalStack
- The Go backend serves API requests correctly
- The React frontend fetches and renders data from the backend
- The Vite proxy correctly forwards `/api` requests to the backend

## Requirements

### Requirement 1: Seed Data

**User Story:** As a developer, I want 20 mock music tracks pre-loaded into LocalStack DynamoDB, so that I have realistic data to validate the local stack without manual setup.

#### Acceptance Criteria

1. WHEN the seed script runs THEN 20 tracks across 5 artists SHALL be inserted into the `MusicLibrary` DynamoDB table
2. WHEN a track is seeded THEN it SHALL have PK=`USER#seed-user`, SK=`TRACK#{trackId}` following the existing single-table key pattern
3. WHEN a track is seeded THEN it SHALL include: title, artist, album, genre, year (int), duration (int, seconds)
4. WHEN the seed script runs multiple times THEN it SHALL be idempotent (no duplicates, no errors) using DynamoDB condition expressions
5. WHEN the seed script runs THEN it SHALL use the AWS CLI pointing at LocalStack (`--endpoint-url http://localhost:4566`)

### Requirement 2: Backend Search API

**User Story:** As a developer, I want a search endpoint that filters seed tracks by query, so that I can validate the backend reads from DynamoDB and returns filtered results.

#### Acceptance Criteria

1. WHEN a GET request is made to `/api/v1/hello/search?q={query}` THEN the response SHALL contain tracks matching the query
2. WHEN searching THEN matching SHALL be case-insensitive across title, artist, album, and genre fields
3. WHEN the query is empty or missing THEN the response SHALL return an empty items array
4. WHEN matches are found THEN the response format SHALL be `{ "items": [...], "total": N }`
5. WHEN no matches are found THEN the response SHALL return `{ "items": [], "total": 0 }`

### Requirement 3: Backend Featured API

**User Story:** As a developer, I want a featured tracks endpoint that returns all seed tracks, so that I can validate the full DynamoDB scan path works locally.

#### Acceptance Criteria

1. WHEN a GET request is made to `/api/v1/hello/featured` THEN the response SHALL return all seed tracks (up to limit)
2. WHEN an optional `limit` query parameter is provided THEN the response SHALL return at most that many tracks
3. WHEN no limit is specified THEN the response SHALL default to 20 tracks
4. WHEN the response is returned THEN it SHALL use the format `{ "items": [...], "total": N }`

### Requirement 4: Backend Health Endpoint

**User Story:** As a developer, I want a hello-specific health endpoint, so that I can verify the hello module is wired and responding.

#### Acceptance Criteria

1. WHEN a GET request is made to `/api/v1/hello/health` THEN the response SHALL be `{ "status": "ok", "service": "hello" }`
2. WHEN the endpoint is called THEN it SHALL return HTTP 200

### Requirement 5: Frontend Search Page

**User Story:** As a developer, I want a search page at `/hello-search` with a hero section, search input, featured tracks grid, and search results, so that I can visually confirm the full stack is working.

#### Acceptance Criteria

1. WHEN navigating to `/hello-search` THEN a hero section with title and description SHALL render
2. WHEN the page loads THEN featured tracks SHALL be fetched and displayed in a responsive grid
3. WHEN the user types in the search input THEN search results SHALL appear after a debounce period (300ms)
4. WHEN search results are loading THEN skeleton loading cards SHALL be displayed
5. WHEN an error occurs THEN an error message SHALL be displayed with a retry option
6. WHEN a track card renders THEN it SHALL show title, artist, album, genre, year, and formatted duration
7. WHEN the page is styled THEN it SHALL use DaisyUI 5 semantic classes consistent with the existing theme

### Requirement 6: No Authentication Required

**User Story:** As a developer, I want the hello-world endpoints and page to be fully public (no auth), so that I can test without configuring Cognito locally.

#### Acceptance Criteria

1. WHEN any `/api/v1/hello/*` endpoint is called THEN it SHALL NOT require authentication headers
2. WHEN navigating to `/hello-search` THEN it SHALL NOT redirect to login
3. WHEN the hello API is called without an X-User-ID header THEN it SHALL still return data (using hardcoded `seed-user`)

## Non-Functional Requirements

### Code Architecture
- **Existing Pattern Compliance**: Backend code follows handler -> service -> repository layering
- **Frontend Pattern Compliance**: Uses TanStack Query hooks, apiClient from client.ts, DaisyUI 5 components
- **Isolation**: Hello-world code is self-contained and does not modify existing services/routes

### Performance
- Seed script SHALL complete within 10 seconds for 20 tracks
- Search endpoint SHALL respond within 500ms for DynamoDB query
- Frontend SHALL debounce search input by 300ms to avoid excessive API calls

### Reliability
- Seed script SHALL be idempotent (safe to run multiple times)
- Backend SHALL handle DynamoDB connection errors gracefully
- Frontend SHALL show error states with retry capability

## Out of Scope

- Full-text search indexing (Nixiesearch) -- uses simple DynamoDB scan + filter
- Audio playback or streaming
- User authentication or role-based access
- Track creation/update/delete mutations
- Pagination (all 20 tracks fit in a single response)
- Production deployment
