# Requirements Document: Hello World Local Development Validation

## Overview

This feature adds a standalone hello-world music search page to validate the local development stack (LocalStack + Go API + React frontend) end-to-end. It is isolated from the main application and requires no authentication. This is iteration 3 of the SDLC validation, incorporating fixes for API contract mismatches and BUILD phase gaps found in iterations 1 and 2.

## User Stories

### US-1: Seed Data
As a developer, I want 20 pre-seeded tracks in the local DynamoDB so that I can test search and listing without manual data entry.

**Acceptance Criteria:**
- 20 tracks exist in MusicLibrary table under `USER#seed-user`
- Tracks span 5 artists with 4 tracks each
- Each track has: Title, Artist, Album, Genre, Year, Duration
- Seed script is idempotent (safe to run multiple times)
- Seed script runs automatically as part of `make local`

### US-2: Backend Search API
As a developer, I want a `GET /api/v1/hello/search?q={query}` endpoint that returns matching seed tracks.

**Acceptance Criteria:**
- Returns JSON: `{"items": [...], "total": N}` (MUST use `items` key, NOT `tracks`)
- Query parameter is `q` (NOT `query`)
- Search is case-insensitive across title, artist, album, genre
- Empty query returns empty items array
- No authentication required

### US-3: Backend Featured API
As a developer, I want a `GET /api/v1/hello/featured?limit={N}` endpoint that returns all seed tracks.

**Acceptance Criteria:**
- Returns JSON: `{"items": [...], "total": N}` (MUST use `items` key)
- Default limit is 20 (returns all seed tracks)
- Optional `limit` parameter to restrict count
- No authentication required

### US-4: Backend Health
As a developer, I want a `GET /api/v1/hello/health` endpoint for health checks.

**Acceptance Criteria:**
- Returns `{"status": "ok", "service": "hello"}`
- No authentication required

### US-5: Frontend Search Page
As a developer, I want a `/hello-search` page to search and browse seed tracks.

**Acceptance Criteria:**
- Page renders at `/hello-search` (public route, no auth)
- Shows featured tracks on initial load
- Search input filters results as user types
- Shows loading skeletons during fetch
- Shows error message on API failure
- Track cards show title, artist, album, genre, year, duration
- Duration formatted as MM:SS

### US-6: No Authentication
As a developer, I want all hello endpoints and the search page accessible without login.

**Acceptance Criteria:**
- Backend routes have no auth middleware
- Frontend route is in PUBLIC_ROUTES array in `__root.tsx`
- Works with `make local` without logging in

## Non-Functional Requirements

- R-NF1: Isolated from main application logic (separate service, handler, route)
- R-NF2: All existing tests must continue to pass (437 frontend, all backend)
- R-NF3: Code follows existing codebase patterns (handler/service/repository layering)
- R-NF4: API contract between frontend and backend must be verified after implementation (contract alignment check)

## Iteration 3 Fixes (from iterations 1 and 2)

| # | Gap Found | Fix Applied |
|---|-----------|-------------|
| 1 | Frontend sent `query` param, backend expected `q` | API contract explicitly specifies `q` in requirements, design, and tasks |
| 2 | Frontend expected `tracks` key, backend returned `items` | API contract explicitly specifies `items` key everywhere |
| 3 | No BUILD phase validation in tasks | Explicit BUILD validation task with `go vet`, `tsc --noEmit`, `eslint`, `npm run build` |
| 4 | Route tree not regenerated after adding new route file | Task includes `npx @tanstack/router-cli generate` step |
| 5 | No contract alignment verification step | Task 7.1 explicitly checks FE/BE param and key alignment |
