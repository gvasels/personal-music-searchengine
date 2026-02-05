# Requirements Document

## Introduction

Hello World Local Dev is a validation feature that proves the full stack (LocalStack DynamoDB, Go backend, React frontend) works end-to-end in the local development environment. It seeds mock music data into DynamoDB, exposes search and featured-tracks API endpoints, and renders a search page in the frontend. This feature serves as a smoke test for `make local` and a template for future feature development.

## Alignment with Product Vision

Per root `CLAUDE.md`, this is a "Multi-user music library platform" with DynamoDB single-table design, Go/Echo backend, and React/TanStack/DaisyUI frontend. This feature validates that all layers connect correctly in the local environment before building real features on top.

## Requirements

### Requirement 1: Seed Music Data

**User Story:** As a developer, I want mock music data seeded into LocalStack DynamoDB when I run `make local`, so that I have realistic data to test against.

#### Acceptance Criteria

1. WHEN `make local` is run THEN the system SHALL seed 20 mock tracks into the MusicLibrary DynamoDB table
2. WHEN seed data is written THEN each track SHALL have: title, artist, album, genre, year, duration fields
3. WHEN seed tracks are written THEN they SHALL be stored under `PK=USER#seed-user` partition
4. WHEN seed data already exists THEN the script SHALL be idempotent (re-running does not create duplicates)
5. WHEN seed data is written THEN there SHALL be at least 5 distinct artists across the 20 tracks

### Requirement 2: Backend Search API

**User Story:** As a frontend developer, I want a search API endpoint that queries seed data by text match, so that I can build a search UI against it.

#### Acceptance Criteria

1. WHEN a GET request is made to `/api/v1/hello/search?q=<query>` THEN the system SHALL return tracks matching the query
2. WHEN searching THEN the system SHALL match case-insensitively across title, artist, album, and genre fields
3. WHEN no query parameter is provided THEN the system SHALL return a 400 error
4. WHEN the query matches no tracks THEN the system SHALL return an empty items array with 200 status
5. WHEN a `limit` query parameter is provided THEN the system SHALL return at most that many results

### Requirement 3: Backend Featured Tracks API

**User Story:** As a frontend developer, I want a featured tracks endpoint that returns all seed data, so that I can display a default track listing.

#### Acceptance Criteria

1. WHEN a GET request is made to `/api/v1/hello/featured` THEN the system SHALL return seed tracks
2. WHEN called with no parameters THEN the system SHALL return up to 20 tracks
3. WHEN a `limit` parameter is provided THEN the system SHALL return at most that many tracks

### Requirement 4: Backend Health Check

**User Story:** As a developer, I want a health check endpoint for the hello service, so that I can verify the API is running.

#### Acceptance Criteria

1. WHEN a GET request is made to `/api/v1/hello/health` THEN the system SHALL return `{"status":"ok","service":"hello"}`

### Requirement 5: Frontend Search Page

**User Story:** As a developer, I want a search page at `/hello-search` that displays featured tracks and allows searching, so that I can visually verify the full stack works.

#### Acceptance Criteria

1. WHEN navigating to `/hello-search` THEN the system SHALL display a hero section and featured tracks
2. WHEN featured tracks load THEN the system SHALL display each track as a card with title, artist, album, genre, and duration
3. WHEN a user types a query and presses Enter THEN the system SHALL call the search API and display results
4. WHEN tracks are loading THEN the system SHALL display skeleton placeholders
5. WHEN the search returns no results THEN the system SHALL display a "no results" message
6. WHEN the page loads THEN it SHALL be accessible without authentication (public route)

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Seed script, backend service, handler, and frontend components are each in separate files
- **Modular Design**: Frontend components (SearchInput, TrackCard, TrackCardSkeleton, HelloNav) are reusable
- **Dependency Management**: HelloService depends only on Repository interface, not concrete implementation

### Performance
- Search operates on at most 20 seed tracks; no pagination needed
- Frontend uses TanStack Query caching with 5-minute stale time

### Security
- Hello endpoints do not require authentication (they read seed data only)
- No user data is exposed; seed data is synthetic

### Reliability
- Seed script is idempotent
- Frontend handles loading, error, and empty states gracefully
