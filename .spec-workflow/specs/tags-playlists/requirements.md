# Requirements - Tags & Playlists (Epic 4)

## Epic Overview

**Epic**: Tags & Playlists
**Wave**: 3
**Dependencies**: Epic 3 (Search & Streaming)

Enable users to organize their music library with custom tags and playlists. Tags provide flexible categorization, while playlists allow curated track collections with ordering.

---

## Functional Requirements

### FR-1: Tag Management

**FR-1.1**: Users can create custom tags
- Tag names are unique per user
- Tags have optional color for visual distinction
- Tag names limited to 50 characters

**FR-1.2**: Users can manage their tags
- View list of all tags with track counts
- Update tag name or color
- Delete tags (removes all track associations)

**FR-1.3**: Users can tag tracks
- Add one or more tags to a track
- Tags are auto-created if they don't exist
- Remove tags from tracks
- View all tags on a track

**FR-1.4**: Users can browse tracks by tag
- Get all tracks with a specific tag
- Multiple tags use AND logic (tracks must have ALL specified tags)

### FR-2: Tag Filtering in Search

**FR-2.1**: Search results can be filtered by tags
- Pass tag names as filter parameters
- Multiple tags use AND logic
- Duplicate tag names in request are deduplicated silently
- Tags must exist or return NotFoundError

**FR-2.2**: Tag validation in search
- IF tag does not exist THEN return NotFoundError with tag name
- IF all tags valid THEN apply post-search filtering
- Return filtered results maintaining search ranking

**FR-2.3**: Tag name case handling
- All tag names normalized to lowercase
- "Rock" and "ROCK" match the same tag "rock"
- Normalization applied in service layer (not just search)

### FR-3: Playlist Management

**FR-3.1**: Users can create playlists
- Playlist name (required)
- Description (optional)
- Public/private visibility flag

**FR-3.2**: Users can manage playlists
- Update playlist name, description, visibility
- Delete playlist (removes all track associations)
- View all playlists with track counts and duration

**FR-3.3**: Users can add tracks to playlists
- Add tracks at specific position (append by default)
- Tracks maintain order within playlist
- Duplicate tracks allowed (same track multiple times)

**FR-3.4**: Users can remove tracks from playlists
- Remove specific tracks by ID
- Positions automatically reorder

**FR-3.5**: Users can view playlist contents
- Get playlist with all tracks in order
- Includes full track metadata and cover art URLs

---

## Non-Functional Requirements

### NFR-1: Performance

**NFR-1.1**: Tag operations complete < 200ms
**NFR-1.2**: Playlist operations complete < 300ms
**NFR-1.3**: Tag filtering in search adds < 100ms overhead

### NFR-2: Scalability

**NFR-2.1**: Support up to 1000 tags per user
**NFR-2.2**: Support up to 100 playlists per user
**NFR-2.3**: Support up to 500 tracks per playlist

### NFR-3: Data Integrity

**NFR-3.1**: Tag deletion cascades to all track associations
**NFR-3.2**: Playlist track positions maintain consistency
**NFR-3.3**: Track deletion removes from all playlists

### NFR-4: Test Coverage

**NFR-4.1**: Unit tests for all tag service methods (80%+ coverage)
**NFR-4.2**: Unit tests for all playlist service methods (80%+ coverage)
**NFR-4.3**: Unit tests for tag filtering in search

---

## User Stories

### US-1: Tag Track
```
As a user, I want to tag my tracks with custom labels
So that I can organize my music beyond album/artist
```

**Acceptance Criteria**:
- Can add tags like "workout", "chill", "favorites"
- Tags appear on track detail view
- Can remove tags from tracks

### US-2: Browse by Tag
```
As a user, I want to browse all tracks with a specific tag
So that I can find music matching a mood or activity
```

**Acceptance Criteria**:
- Click tag to see all tracks
- Multiple tag selection narrows results
- Shows track count for each tag

### US-3: Search with Tags
```
As a user, I want to filter search results by my tags
So that I can find specific tracks within a category
```

**Acceptance Criteria**:
- Search "rock" filtered by tag "favorites"
- Results only include tracks matching both
- Clear error if tag doesn't exist

### US-4: Create Playlist
```
As a user, I want to create playlists of my favorite tracks
So that I can have curated collections for different occasions
```

**Acceptance Criteria**:
- Create playlist with name and description
- Add tracks in specific order
- Rearrange tracks within playlist

### US-5: Play Playlist
```
As a user, I want to play an entire playlist in order
So that I can enjoy my curated music collection
```

**Acceptance Criteria**:
- Play button starts from first track
- Tracks play in playlist order
- Shows current position in playlist

---

## API Endpoints

### Tag Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tags` | List user's tags |
| POST | `/api/v1/tags` | Create new tag |
| GET | `/api/v1/tags/{name}` | Get tag details |
| PUT | `/api/v1/tags/{name}` | Update tag |
| DELETE | `/api/v1/tags/{name}` | Delete tag |
| POST | `/api/v1/tracks/{id}/tags` | Add tags to track |
| DELETE | `/api/v1/tracks/{id}/tags/{name}` | Remove tag from track |
| GET | `/api/v1/tags/{name}/tracks` | Get tracks by tag |

### Playlist Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/playlists` | List user's playlists |
| POST | `/api/v1/playlists` | Create playlist |
| GET | `/api/v1/playlists/{id}` | Get playlist with tracks |
| PUT | `/api/v1/playlists/{id}` | Update playlist |
| DELETE | `/api/v1/playlists/{id}` | Delete playlist |
| POST | `/api/v1/playlists/{id}/tracks` | Add tracks |
| DELETE | `/api/v1/playlists/{id}/tracks` | Remove tracks |

### Search with Tags

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/search` | Search with optional `filters.tags` array |

---

## Data Models

### Tag (Existing)
```go
type Tag struct {
    UserID     string    // Partition key
    Name       string    // Sort key
    Color      string    // Hex color for UI
    TrackCount int       // Tracks with this tag
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

### TrackTag (Existing - DynamoDB)
```
PK: USER#{userId}#TRACK#{trackId}
SK: TAG#{tagName}
GSI1PK: USER#{userId}#TAG#{tagName}
GSI1SK: TRACK#{trackId}
```

### Playlist (Existing)
```go
type Playlist struct {
    ID            string
    UserID        string
    Name          string
    Description   string
    CoverArtKey   string
    TrackCount    int
    TotalDuration int
    IsPublic      bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### PlaylistTrack (Existing - DynamoDB)
```
PK: PLAYLIST#{playlistId}
SK: POSITION#{pos}
TrackID: string
AddedAt: timestamp
```

---

## Existing Implementation Status

### Already Implemented
- Tag model (`models/tag.go`)
- Playlist model (`models/playlist.go`)
- Tag repository (`repository/dynamodb.go`)
- Playlist repository (`repository/dynamodb.go`)
- Tag service (`service/tag.go`)
- Playlist service (`service/playlist.go`)
- Tag handlers (`handlers/tag.go`)
- Playlist handlers (`handlers/playlist.go`)
- SearchFilters.Tags field (`models/search.go`)

### Missing (To Implement)
- **Tag filtering in search** - `filterByTags` function in `service/search.go`
- **Tag service unit tests** - `service/tag_test.go`
- **Playlist service unit tests** - `service/playlist_test.go`
- **filterByTags unit tests** - tests in `service/search_test.go`

---

## Out of Scope

- Shared/public tags between users
- Tag suggestions or autocomplete
- Smart playlists (auto-generated based on rules)
- Collaborative playlists
- Playlist folders or nesting
- Drag-and-drop reordering (frontend feature)
