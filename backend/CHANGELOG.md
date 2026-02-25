# Changelog

All notable changes to the backend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Artist Events & Shows (Phase 1 — Mock)**
  - `EventsProvider` interface with `MockProvider` for development/testing (`internal/events/`)
  - `ArtistWatch` model and DynamoDB CRUD for watch/unwatch (`internal/models/`, `internal/repository/`)
  - `ArtistWatchService` for watch operations (`internal/service/artist_watch.go`)
  - `EventsService` for event retrieval and search (`internal/service/events.go`)
  - `ArtistWatchHandler` with 4 endpoints: watch, unwatch, status, list (`internal/handlers/artist_watch.go`)
  - `EventsHandler` with 2 endpoints: get artist events, search events (`internal/handlers/events.go`)
  - 85+ new backend tests covering handlers, services, and mock provider
- **Audio Pipeline Improvements**
  - Transcode-complete wiring to audio analysis Lambda
  - Lambda architecture fix (x86_64 → arm64)
  - MediaConvert UserMetadata for pipeline reliability
- **Embedding Enhancements**
  - Musical key incorporated into track embeddings
  - Camelot wheel neighbors encoded for harmonic mixing similarity
- **Duplicate Upload Guard**
  - "Already in library" check with edge case handling

### Changed
- Events provider gated behind `EVENTS_PROVIDER=mock` env var for production safety

### Fixed
- **Auth middleware on events routes**: Applied `RequireAuth()` middleware to events route group (was missing)
- **XSS in EventCard**: Validated `ticketUrl` to only allow `http://`/`https://` protocols
- **URL encoding**: Applied `encodeURIComponent()` to artist name URL paths in frontend API clients
- **Limit cap**: Capped `limit` query parameter at 100 in artist watch and events handlers
- **Input validation**: Added max length (256 chars) on `artistName` and `q` parameters
- **Typed error wrapping**: `ArtistWatchService` now returns proper `*models.APIError` types (409 Conflict, 404 Not Found) instead of plain `fmt.Errorf`
- **Dead code removal**: Removed unreachable `errors.Is` branch in `UnwatchArtist` handler
- **Search debounce**: Added 300ms debounce to shows page search input
- **Auth guard**: Added `enabled: isAuthenticated` guard to `useWatchedArtists` hook
- **Error feedback**: Added `onError` toast handlers to watch/unwatch mutations

### Added (continued)
- **Admin Panel & Track Visibility Feature**
  - Admin handlers for user management (`internal/handlers/admin.go`)
  - Admin service with Cognito integration (`internal/service/admin.go`)
  - Track visibility service for private/unlisted/public tracks (`internal/service/track_visibility.go`)
  - Admin handler integration tests (`internal/handlers/admin_test.go`)
- **Global User Type Feature** (Role-Based Access Control)
  - User roles: `guest`, `subscriber`, `artist`, `admin` with Cognito Groups integration
  - Permission system with 12 granular permissions (browse, listen, upload, publish, etc.)
  - Artist profile service with CRUD operations and catalog linking
  - Follow service for user-to-artist-profile relationships
  - Authorization middleware (`internal/handlers/middleware/auth.go`) with role extraction
  - Playlist visibility: `private`, `unlisted`, `public` visibility levels
  - Public playlist discovery endpoint (`GET /playlists/public`)
  - New models: `ArtistProfile`, `Follow`, `UserRole`, `Permission`, `PlaylistVisibility`
  - New handlers: artist profile CRUD, follow/unfollow, role management
  - New services: `ArtistProfileService`, `FollowService`, `RoleService`
  - Comprehensive test coverage for new services (role, artist profile, follow)
- Audio analysis package (`internal/analysis/`) with BPM detection using multi-segment autocorrelation algorithm
- Camelot wheel mapping for harmonic mixing support (24 key mappings including enharmonic equivalents)
- Analyzer Lambda processor (`cmd/processor/analyzer/`) for Step Functions integration
- Migration service (`internal/service/migration.go`) for string-to-entity artist migration
- Playlist reorder endpoint for track position management
- Matching service for DJ-style track compatibility scoring
- FFmpeg input validation to prevent command injection

### Changed
- Updated CI coverage threshold from 19% to 24%
- Added golangci-lint job to CI workflow

### Fixed
- CORS handling for playlist reorder endpoint
- 404 error on playlist reorder route
