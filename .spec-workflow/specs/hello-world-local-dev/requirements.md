# Requirements: Hello World Local Development

## Introduction

A minimal "Hello World" feature to validate the SDLC workflow with full-stack integration: seed data in DynamoDB, backend Go API, and React frontend. This serves as a validation exercise for the SDLC process improvements.

## Alignment with Product Vision

This feature validates the local development workflow by creating a simple music search page that demonstrates the full stack: LocalStack DynamoDB → Go Echo API → React TanStack frontend. It follows patterns from the root `CLAUDE.md` for backend layering (handler → service → repository) and frontend architecture (TanStack Router/Query + DaisyUI).

## Requirements

### US-1: Seed Data

**User Story:** As a developer, I want seed music data in LocalStack DynamoDB, so that I can test the API locally.

#### Acceptance Criteria

1. WHEN LocalStack starts THEN the system SHALL create 20 seed tracks in the MusicLibrary table
2. IF seed script runs multiple times THEN the system SHALL maintain exactly 20 tracks (idempotent)
3. WHEN querying seed data THEN each track SHALL have: Title, Artist, Album, Genre, Year, Duration

### US-2: Search API

**User Story:** As a user, I want to search tracks by query, so that I can find music I'm looking for.

#### Acceptance Criteria

1. WHEN I send GET `/api/v1/hello/search?q=jazz` THEN the system SHALL return tracks matching "jazz" in title, artist, album, or genre
2. IF query matches multiple fields THEN the system SHALL return all matching tracks
3. WHEN search query is empty THEN the system SHALL return an empty array

### US-3: Featured API

**User Story:** As a user, I want to see featured tracks, so that I can discover music without searching.

#### Acceptance Criteria

1. WHEN I send GET `/api/v1/hello/featured` THEN the system SHALL return up to 10 tracks by default
2. IF `limit` parameter is provided THEN the system SHALL return up to that many tracks
3. WHEN `limit=0` or no limit THEN the system SHALL return all available tracks

### US-4: Health API

**User Story:** As a developer, I want a health endpoint, so that I can verify the hello service is running.

#### Acceptance Criteria

1. WHEN I send GET `/api/v1/hello/health` THEN the system SHALL return `{"status":"ok","service":"hello"}`

### US-5: Search UI

**User Story:** As a user, I want a search page with featured tracks and search input, so that I can browse and find music visually.

#### Acceptance Criteria

1. WHEN I visit `/hello-search` THEN the system SHALL display featured tracks
2. WHEN I type in the search box THEN the system SHALL filter tracks matching my query
3. WHEN data is loading THEN the system SHALL display skeleton loading cards
4. WHEN an error occurs THEN the system SHALL display an error message
5. WHEN no results match THEN the system SHALL display "No results found"

### US-6: Public Route

**User Story:** As a user, I want to access the hello search page without authentication, so that it's easy to demo.

#### Acceptance Criteria

1. WHEN I visit `/hello-search` without authentication THEN the system SHALL allow access
2. IF I am not logged in THEN the system SHALL NOT redirect me to login

## Non-Functional Requirements

### R-NF1: Code Architecture
- Backend follows handler → service → repository pattern
- Frontend uses TanStack Router file-based routing with TanStack Query hooks
- DaisyUI 5 semantic classes for UI components

### R-NF2: Performance
- API responses under 200ms for local development
- Frontend renders within 100ms after data received

### R-NF3: Testing
- 80% test coverage minimum
- Unit tests for service and handler layers
- Component tests for all UI components

### R-NF4: API Contract
- Search endpoint uses query param `q` (NOT `query`)
- All list responses use `items` array (NOT `tracks`)
- Response format: `{"items": [...], "total": N}`
