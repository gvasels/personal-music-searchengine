# Tasks - Foundation Backend

## Epic: Foundation Backend
**Status**: Not Started
**Wave**: 0-2

---

## Group 1: Domain Models (Wave 0)

### Task 1.1: Common Types and Constants
**Status**: [x] Complete
**Files**:
- `backend/internal/models/common.go`
- `backend/internal/models/common_test.go`

**Acceptance Criteria**:
- [x] EntityType constants defined
- [x] UploadStatus constants defined
- [x] AudioFormat constants defined
- [x] Timestamps embed struct
- [x] DynamoDBItem base struct
- [x] Pagination helpers
- [x] PaginationCursor with encode/decode (base64 opaque cursor per design)
- [x] Duration and file size formatters
- [x] Unit tests passing

### Task 1.2: Track Model
**Status**: [x] Complete
**Files**:
- `backend/internal/models/track.go`
- `backend/internal/models/track_test.go`

**Acceptance Criteria**:
- [x] Track struct with JSON tags
- [x] TrackItem with DynamoDB keys
- [x] NewTrackItem() function with GSI1 for artist queries
- [x] ToResponse() conversion
- [x] TrackFilter for queries
- [x] UpdateTrackRequest
- [x] Unit tests passing

### Task 1.3: Album Model
**Status**: [x] Complete
**Files**:
- `backend/internal/models/album.go`
- `backend/internal/models/album_test.go`

**Acceptance Criteria**:
- [x] Album struct with JSON tags
- [x] AlbumItem with DynamoDB keys
- [x] NewAlbumItem() with GSI1 for artist queries
- [x] ToResponse() conversion
- [x] ArtistSummary for aggregation
- [x] AlbumFilter for queries
- [x] Unit tests passing

### Task 1.4: User Model
**Status**: [x] Complete
**Files**:
- `backend/internal/models/user.go`
- `backend/internal/models/user_test.go`

**Acceptance Criteria**:
- [x] User struct with storage tracking
- [x] UserItem with DynamoDB keys
- [x] NewUserItem() function
- [x] ToResponse() conversion
- [x] UpdateUserRequest
- [x] Unit tests passing

### Task 1.5: Playlist Model
**Status**: [x] Complete
**Files**:
- `backend/internal/models/playlist.go`
- `backend/internal/models/playlist_test.go`

**Acceptance Criteria**:
- [x] Playlist struct
- [x] PlaylistTrack for ordering
- [x] PlaylistItem and PlaylistTrackItem
- [x] Position zero-padding for sort order
- [x] ToResponse() conversion
- [x] CRUD request types
- [x] Unit tests passing

### Task 1.6: Tag Model
**Status**: [x] Complete
**Files**:
- `backend/internal/models/tag.go`
- `backend/internal/models/tag_test.go`

**Acceptance Criteria**:
- [x] Tag struct
- [x] TrackTag for associations
- [x] TagItem and TrackTagItem
- [x] GSI1 for tag-based track lookup
- [x] ToResponse() conversion
- [x] Unit tests passing

### Task 1.7: Upload Model
**Status**: [x] Complete
**Files**:
- `backend/internal/models/upload.go`
- `backend/internal/models/upload_test.go`

**Acceptance Criteria**:
- [x] Upload struct with status tracking
- [x] UploadItem with status GSI
- [x] PresignedUploadRequest/Response (1GB max, multipart support)
- [x] ConfirmUploadRequest/Response
- [x] UploadMetadata for extracted data
- [x] ToResponse() conversion
- [x] Step tracking for partial success recovery (per design)
- [x] Multipart upload tracking fields
- [x] ReprocessUploadRequest for retry from specific step
- [x] CoverArtUploadRequest/Response
- [x] CompleteMultipartUploadRequest
- [x] Unit tests passing

### Task 1.8: Search Models
**Status**: [x] Complete
**Files**:
- `backend/internal/models/search.go`

**Acceptance Criteria**:
- [x] SearchRequest with cursor-based pagination (per design)
- [x] SearchResponse with nextCursor and hasMore
- [x] NixieSearchQuery with search_after for pagination
- [x] SearchFilters and SearchSort
- [x] SearchFacets for filtering UI

### Task 1.9: Error Models
**Status**: [x] Complete
**Files**:
- `backend/internal/models/errors.go`

**Acceptance Criteria**:
- [x] APIError base struct
- [x] Common errors (NotFound, Unauthorized, etc.)
- [x] Upload-specific errors (UploadExpired, UploadNotFound, etc.)
- [x] ErrInvalidCursor for pagination errors
- [x] NewValidationError, NewNotFoundError, NewConflictError helpers

---

## Group 2: Repository Layer (Wave 1)

### Task 2.1: Repository Interface
**Status**: [x] Complete
**Files**:
- `backend/internal/repository/repository.go`

**Acceptance Criteria**:
- [x] Repository interface defined (per design.md)
- [x] S3Repository interface defined
- [x] CloudFrontSigner interface defined
- [x] Common error types (ErrNotFound, ErrAlreadyExists, ErrInvalidCursor)
- [x] PaginatedResult generic type

### Task 2.2: DynamoDB Repository - Core
**Status**: [x] Complete
**Files**:
- `backend/internal/repository/dynamodb.go`

**Acceptance Criteria**:
- [x] DynamoDBRepository struct with client interface
- [x] DynamoDBClient interface for testability
- [x] NewDynamoDBRepository constructor
- [x] Cursor encoding/decoding with models.PaginationCursor
- [x] Error handling and mapping

### Task 2.3: DynamoDB Repository - Entity Operations
**Status**: [x] Complete
**Files**:
- `backend/internal/repository/dynamodb.go`

**Acceptance Criteria**:
- [x] Track CRUD operations (Create, Get, Update, Delete, List, ListByArtist)
- [x] Album operations (GetOrCreate, Get, List, ListByArtist, UpdateStats)
- [x] User profile operations (Create, Get, Update, UpdateStats)
- [x] Playlist operations (CRUD, AddTracks, RemoveTracks, GetTracks)
- [x] Tag operations (CRUD, AddToTrack, RemoveFromTrack, GetTrackTags, GetTracksByTag)
- [x] Upload tracking (Create, Get, Update, UpdateStatus, UpdateStep, List, ListByStatus)
- [x] Pagination support with opaque cursors
- [x] Batch operations for playlist tracks and tags

### Task 2.4: S3 Repository
**Status**: [x] Complete
**Files**:
- `backend/internal/repository/s3.go`

**Acceptance Criteria**:
- [x] S3RepositoryImpl implementation
- [x] S3Client and S3PresignClient interfaces for testability
- [x] GeneratePresignedUploadURL (with Intelligent-Tiering storage class)
- [x] GeneratePresignedDownloadURL
- [x] Multipart upload support (Initiate, GeneratePartURLs, Complete, Abort)
- [x] DeleteObject, CopyObject
- [x] GetObjectMetadata, ObjectExists

---

## Group 3: Service Layer (Wave 1)

### Task 3.1: Services Implementation
**Status**: [x] Complete
**Files**:
- `backend/internal/service/service.go` - Interfaces and Services container
- `backend/internal/service/track.go` - TrackService implementation
- `backend/internal/service/album.go` - AlbumService implementation
- `backend/internal/service/user.go` - UserService implementation
- `backend/internal/service/playlist.go` - PlaylistService implementation
- `backend/internal/service/tag.go` - TagService implementation
- `backend/internal/service/upload.go` - UploadService implementation
- `backend/internal/service/stream.go` - StreamService implementation

**Acceptance Criteria**:
- [x] TrackService for track management (CRUD, play count)
- [x] AlbumService for album aggregation and artist listing
- [x] UserService for profile management
- [x] PlaylistService for playlist operations (CRUD, tracks)
- [x] TagService for tagging operations (CRUD, track associations)
- [x] UploadService for presigned URLs (single and multipart)
- [x] StreamService for CloudFront signed URLs (stream/download)
- [ ] SearchService for Nixiesearch integration (separate implementation)
- [x] Cover art URL generation for all entity responses
- [x] Pagination support with opaque cursors

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
| Group 1: Domain Models | 9 | **Complete** |
| Group 2: Repository Layer | 4 | **Complete** |
| Group 3: Service Layer | 1 | **Complete** |
| Group 4: HTTP Handlers | 2 | Not Started |
| **Total** | **16** | **3 Complete** |

---

## Design Questions Resolved

All design questions have been resolved in `design.md`:
1. S3 storage class → **Intelligent-Tiering** (auto-moves between tiers based on access)
2. Caching strategy → **CloudFront only** (cache static assets and API responses at edge)
3. Search indexing approach → **Hybrid** (real-time for new tracks, weekly batch re-index)
4. Pagination cursor format → **Opaque base64** (encode lastEvaluatedKey as base64)
5. Retry strategy → **AWS SDK defaults** (3 retries with exponential backoff)
6. Upload failures → **Partial success** (continue even if some steps fail, allow re-processing)
7. Max file size → **1 GB** (requires multipart upload handling)
