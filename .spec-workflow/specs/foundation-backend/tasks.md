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
| Group 1: Domain Models | 9 | **Complete** |
| Group 2: Repository Layer | 4 | Not Started |
| Group 3: Service Layer | 1 | Not Started |
| Group 4: HTTP Handlers | 2 | Not Started |
| **Total** | **16** | **1 Complete** |

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
