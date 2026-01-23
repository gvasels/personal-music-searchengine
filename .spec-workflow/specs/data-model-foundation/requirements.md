# Requirements Document: Data Model Foundation

## Introduction

This spec establishes a proper Artist entity model with UUIDs to replace the current string-based artist references. This foundation is critical for enabling multi-artist tracks, artist profiles, rights management, and creator studio features. Currently, tracks store artist as a simple string field, making it impossible to link multiple artists to a track, maintain artist metadata, or support the streaming service model.

## Alignment with Product Vision

This directly supports:
- **Phase 2: Global User Type** - Artist profiles require proper Artist entities
- **Creator Studio** - DJ/Producer features need artist linking for collaboration
- **Rights Management** - Rights are assigned to artists, requiring proper entity relationships
- **Search & Discovery** - Artist-based recommendations need structured artist data

## Requirements

### Requirement 1: Artist Entity Model

**User Story:** As a platform developer, I want a dedicated Artist entity with UUIDs, so that I can uniquely identify artists regardless of name collisions (e.g., multiple artists named "The Band").

#### Acceptance Criteria

1. WHEN an artist is created THEN the system SHALL generate a unique UUID as the artist ID
2. WHEN an artist entity is stored THEN the system SHALL use DynamoDB single-table pattern with `PK: USER#{userId}` and `SK: ARTIST#{artistId}`
3. WHEN an artist is retrieved THEN the system SHALL return id, name, sortName, bio, imageUrl, externalLinks, createdAt, updatedAt
4. IF an artist name already exists for a user THEN the system SHALL still allow creation (names are NOT unique)
5. WHEN searching artists by name THEN the system SHALL use a GSI on the name field

### Requirement 2: Artist-Track Linking

**User Story:** As a user, I want tracks linked to artist IDs instead of artist name strings, so that I can navigate to the correct artist profile when clicking an artist name.

#### Acceptance Criteria

1. WHEN a track is created THEN the system SHALL store `artistId` (UUID) instead of `artist` (string)
2. WHEN a track is returned via API THEN the system SHALL include both `artistId` and resolved `artistName` for display
3. WHEN an artist is deleted THEN the system SHALL NOT cascade delete tracks (orphan tracks retain artistId)
4. IF a track's artistId references a non-existent artist THEN the system SHALL return "Unknown Artist" for artistName
5. WHEN migrating existing data THEN the system SHALL create Artist entities from unique artist names and update track references

### Requirement 3: Multiple Artists per Track

**User Story:** As a user uploading a collaboration track, I want to assign multiple artists (main artist, featuring artists, remixer), so that all contributors are properly credited.

#### Acceptance Criteria

1. WHEN a track is created/updated THEN the system SHALL accept an array of artist contributions: `[{artistId, role: "main"|"featuring"|"remixer"|"producer"}]`
2. WHEN a track is displayed THEN the system SHALL format artists as "Artist1 ft. Artist2 (Remixed by Artist3)"
3. IF no role is specified THEN the system SHALL default to "main"
4. WHEN searching by artist THEN the system SHALL return tracks where the artist appears in ANY role
5. WHEN viewing an artist profile THEN the system SHALL show tracks grouped by role (Main, Featured, Remixes, Productions)

### Requirement 4: Artist Metadata

**User Story:** As an artist/creator, I want to maintain a profile with bio, images, and external links, so that fans can learn about me and find me on other platforms.

#### Acceptance Criteria

1. WHEN updating artist metadata THEN the system SHALL allow: bio (text), imageUrl (S3 key), externalLinks (map of platform→URL)
2. WHEN external links are provided THEN the system SHALL support: spotify, apple_music, soundcloud, bandcamp, instagram, twitter, youtube, website
3. IF an artist has no imageUrl THEN the system SHALL return a placeholder/generated avatar
4. WHEN displaying artist metadata THEN the system SHALL show track count, album count, and total plays
5. WHEN an artist uploads a profile image THEN the system SHALL store it in S3 under `artists/{artistId}/profile.{ext}`

### Requirement 5: Artist API Endpoints

**User Story:** As a frontend developer, I want complete CRUD endpoints for artists, so that I can build artist management and profile pages.

#### Acceptance Criteria

1. WHEN `POST /api/v1/artists` is called THEN the system SHALL create a new artist and return the created entity
2. WHEN `GET /api/v1/artists` is called THEN the system SHALL return paginated list of user's artists
3. WHEN `GET /api/v1/artists/:id` is called THEN the system SHALL return artist with aggregated stats
4. WHEN `PUT /api/v1/artists/:id` is called THEN the system SHALL update artist metadata
5. WHEN `DELETE /api/v1/artists/:id` is called THEN the system SHALL soft-delete (mark inactive, don't remove)
6. WHEN `GET /api/v1/artists/:id/tracks` is called THEN the system SHALL return tracks grouped by role

### Requirement 6: Data Migration

**User Story:** As an operator, I want to migrate existing string-based artist data to the new entity model without data loss or downtime.

#### Acceptance Criteria

1. WHEN migration runs THEN the system SHALL extract unique artist names from all tracks
2. WHEN migration creates artists THEN the system SHALL generate deterministic UUIDs from (userId, artistName) to ensure idempotency
3. WHEN migration updates tracks THEN the system SHALL populate artistId and artistContributions fields
4. WHEN migration completes THEN the system SHALL retain original `artist` string field as backup for 30 days
5. IF migration fails mid-process THEN the system SHALL be resumable from last checkpoint

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Artist model in `internal/models/artist.go`, repository in `internal/repository/artist.go`, service in `internal/service/artist.go`, handlers in `internal/handlers/artist.go`
- **Modular Design**: Artist service should be independent and injectable
- **Dependency Management**: No circular dependencies between artist, track, and album packages
- **Clear Interfaces**: Define `ArtistRepository` and `ArtistService` interfaces for testability

### Performance
- Artist lookup by ID must complete in <50ms (DynamoDB point read)
- Artist search by name must complete in <200ms (GSI query)
- Track listing with artist resolution must not N+1 query (batch get artists)
- Migration must process 10,000 tracks in <5 minutes

### Security
- Artist profiles are private by default (only owner can view)
- Profile images must be validated (type, size <5MB)
- External links must be sanitized (no XSS in URLs)
- Soft-delete must be irreversible by API (admin-only restore)

### Reliability
- Migration must be idempotent (safe to re-run)
- Artist deletion must not orphan tracks
- Database operations must use transactions where atomicity required

### Usability
- Artist name autocomplete in track upload/edit forms
- "Create new artist" option when no match found
- Clear error messages for duplicate detection attempts
