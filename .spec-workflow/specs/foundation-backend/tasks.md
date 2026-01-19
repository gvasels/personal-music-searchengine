# Tasks - Foundation Backend

## Epic: Foundation Backend
**Status**: Not Started
**Wave**: 0-2

---

## Group 1: Domain Models (Wave 0)

### Task 1.1: Common Types and Constants
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/common.go`
- `backend/internal/models/common_test.go`

**Acceptance Criteria**:
- [ ] EntityType constants defined
- [ ] UploadStatus constants defined
- [ ] AudioFormat constants defined
- [ ] Timestamps embed struct
- [ ] DynamoDBItem base struct
- [ ] Pagination helpers
- [ ] Duration and file size formatters
- [ ] Unit tests passing

### Task 1.2: Track Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/track.go`
- `backend/internal/models/track_test.go`

**Acceptance Criteria**:
- [ ] Track struct with JSON tags
- [ ] TrackItem with DynamoDB keys
- [ ] NewTrackItem() function with GSI1 for artist queries
- [ ] ToResponse() conversion
- [ ] TrackFilter for queries
- [ ] UpdateTrackRequest
- [ ] Unit tests passing

### Task 1.3: Album Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/album.go`
- `backend/internal/models/album_test.go`

**Acceptance Criteria**:
- [ ] Album struct with JSON tags
- [ ] AlbumItem with DynamoDB keys
- [ ] NewAlbumItem() with GSI1 for artist queries
- [ ] ToResponse() conversion
- [ ] ArtistSummary for aggregation
- [ ] AlbumFilter for queries
- [ ] Unit tests passing

### Task 1.4: User Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/user.go`
- `backend/internal/models/user_test.go`

**Acceptance Criteria**:
- [ ] User struct with storage tracking
- [ ] UserItem with DynamoDB keys
- [ ] NewUserItem() function
- [ ] ToResponse() conversion
- [ ] UpdateUserRequest
- [ ] Unit tests passing

### Task 1.5: Playlist Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/playlist.go`
- `backend/internal/models/playlist_test.go`

**Acceptance Criteria**:
- [ ] Playlist struct
- [ ] PlaylistTrack for ordering
- [ ] PlaylistItem and PlaylistTrackItem
- [ ] Position zero-padding for sort order
- [ ] ToResponse() conversion
- [ ] CRUD request types
- [ ] Unit tests passing

### Task 1.6: Tag Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/tag.go`
- `backend/internal/models/tag_test.go`

**Acceptance Criteria**:
- [ ] Tag struct
- [ ] TrackTag for associations
- [ ] TagItem and TrackTagItem
- [ ] GSI1 for tag-based track lookup
- [ ] ToResponse() conversion
- [ ] Unit tests passing

### Task 1.7: Upload Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/upload.go`
- `backend/internal/models/upload_test.go`

**Acceptance Criteria**:
- [ ] Upload struct with status tracking
- [ ] UploadItem with status GSI
- [ ] PresignedUploadRequest/Response
- [ ] ConfirmUploadRequest/Response
- [ ] UploadMetadata for extracted data
- [ ] ToResponse() conversion
- [ ] Unit tests passing

---

## Group 2: Repository Layer (Wave 1)

### Task 2.1: Repository Interface
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/repository.go`

**Acceptance Criteria**:
- [ ] Repository interface defined (per design.md)
- [ ] S3Repository interface defined
- [ ] Common error types

### Task 2.2: DynamoDB Repository - Core
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/dynamodb.go`
- `backend/internal/repository/dynamodb_test.go`

**Acceptance Criteria**:
- [ ] DynamoDBRepository struct
- [ ] Connection and configuration
- [ ] Put, Get, Update, Delete, Query helpers
- [ ] Error handling and mapping
- [ ] Integration tests with DynamoDB Local

### Task 2.3: DynamoDB Repository - Entity Operations
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/dynamodb_*.go`

**Acceptance Criteria**:
- [ ] Track CRUD operations
- [ ] Album operations with artist queries
- [ ] User profile operations
- [ ] Playlist operations with track ordering
- [ ] Tag operations with track associations
- [ ] Upload status tracking
- [ ] Integration tests

### Task 2.4: S3 Repository
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/s3.go`
- `backend/internal/repository/s3_test.go`

**Acceptance Criteria**:
- [ ] S3Repository implementation
- [ ] GeneratePresignedUploadURL
- [ ] GeneratePresignedDownloadURL
- [ ] DeleteObject, CopyObject
- [ ] GetObjectMetadata
- [ ] Unit tests with mocked S3 client

---

## Group 3: Service Layer (Wave 1)

### Task 3.1: Services Implementation
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/*.go`
- `backend/internal/service/*_test.go`

**Acceptance Criteria**:
- [ ] TrackService for track management
- [ ] AlbumService for album aggregation
- [ ] UserService for profile management
- [ ] PlaylistService for playlist operations
- [ ] TagService for tagging operations
- [ ] UploadService for presigned URLs
- [ ] StreamService for CloudFront signed URLs
- [ ] SearchService for Nixiesearch integration
- [ ] Unit tests with mocked repository

---

## Group 4: HTTP Handlers (Wave 2)

### Task 4.1: Handler Setup
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/handlers.go`
- `backend/cmd/api/main.go`

**Acceptance Criteria**:
- [ ] Handlers struct with service dependencies
- [ ] getUserIDFromContext (from API Gateway claims)
- [ ] Common error response helpers
- [ ] Echo router configuration
- [ ] Lambda adapter setup

### Task 4.2: Entity Handlers
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/*.go`
- `backend/internal/handlers/*_test.go`

**Acceptance Criteria**:
- [ ] Track handlers (CRUD)
- [ ] Album handlers (list, get)
- [ ] User handlers (profile)
- [ ] Playlist handlers (CRUD + tracks)
- [ ] Tag handlers (CRUD + associations)
- [ ] Upload handlers (presigned, confirm, status)
- [ ] Stream handlers (stream/download URLs)
- [ ] Search handlers (query, advanced)
- [ ] HTTP tests with mocked services

---

## Summary

| Group | Tasks | Status |
|-------|-------|--------|
| Group 1: Domain Models | 7 | Not Started |
| Group 2: Repository Layer | 4 | Not Started |
| Group 3: Service Layer | 1 | Not Started |
| Group 4: HTTP Handlers | 2 | Not Started |
| **Total** | **14** | **Not Started** |

---

## Design Questions to Resolve

Before implementation, clarify:
1. S3 storage class (Intelligent-Tiering vs Standard?)
2. Caching strategy (CloudFront? DynamoDB DAX? In-memory?)
3. Search indexing approach (real-time vs batch?)
4. Error handling patterns
5. Pagination cursor format
