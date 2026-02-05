# Service Layer - CLAUDE.md

## Overview

Business logic layer implementing domain operations for the Personal Music Search Engine. Services orchestrate repository calls, handle validation, and implement business rules.

## File Descriptions

| File | Purpose |
|------|---------|
| `service.go` | Service interfaces and Services container |
| `track.go` | TrackService - track management operations |
| `album.go` | AlbumService - album operations and artist aggregation |
| `user.go` | UserService - user profile management |
| `playlist.go` | PlaylistService - playlist CRUD and track management |
| `playlist_test.go` | Unit tests for PlaylistService (16 tests) |
| `tag.go` | TagService - tag management and track associations |
| `tag_test.go` | Unit tests for TagService (24 tests) |
| `upload.go` | UploadService - upload workflow and presigned URLs |
| `stream.go` | StreamService - streaming and download URL generation |
| `search.go` | SearchService - Nixiesearch integration for full-text search |
| `search_test.go` | Unit tests for SearchService including filterByTags (8 tests) |
| `transcode.go` | TranscodeService - MediaConvert HLS transcoding |
| `transcode_test.go` | Unit tests for TranscodeService |
| `migration.go` | MigrationService - artist migration from string to entity model |
| `migration_test.go` | Unit tests for MigrationService |
| `embedding.go` | EmbeddingService - Bedrock Titan text embeddings |
| `embedding_test.go` | Unit tests for EmbeddingService (20 tests) |
| `camelot.go` | Camelot key compatibility utilities for DJ mixing |
| `camelot_test.go` | Unit tests for Camelot utilities |
| `similarity.go` | SimilarityService - similar/mixable tracks for DJs |

## Service Interfaces

### TrackService
- `GetTrack(ctx, userID, trackID, hasGlobal)` - Get track with visibility enforcement
  - **hasGlobal=true**: Admin/global access - can access any track
  - **hasGlobal=false**: Regular user - must own track OR track must be public/unlisted
  - Returns **403 Forbidden** for unauthorized access to private tracks
  - Returns **404 Not Found** only for truly non-existent tracks
- `UpdateTrack` - Update track metadata
- `DeleteTrack(ctx, userID, trackID, hasGlobal)` - Delete track with admin support
  - **hasGlobal=false**: Can only delete own tracks
  - **hasGlobal=true**: Admin can delete ANY track
  - Uses `GetTrackByID` for admin to find track regardless of owner
  - Deletes from DynamoDB using actual owner's ID
  - Cleans up S3 files: audio, cover art, HLS transcoded files
  - HLS cleanup uses `DeleteByPrefix("hls/{ownerID}/{trackID}/")`
- `ListTracks(ctx, userID, hasGlobal, filter)` - Paginated track listing with visibility filtering
  - **hasGlobal=true**: Returns ALL tracks (admin view) - scans in batches of 100
  - **hasGlobal=false**: Returns only user's own tracks + public tracks from others
- `ListTracksByArtist` - Query tracks by artist
- `IncrementPlayCount` - Update play count and last played

### TrackVisibilityService
- `ListTracksWithVisibility` - List tracks with proper visibility enforcement and owner display names
  - Admins see all tracks globally
  - Regular users see own tracks + public tracks
  - Deduplicates public tracks that appear in both queries
  - Sets `OwnerDisplayName` to "You" for own tracks

### AlbumService
- `GetAlbum` - Get album with tracks
- `ListAlbums` - Paginated album listing
- `ListAlbumsByArtist` - Query albums by artist
- `ListArtists` - Aggregate artists with track/album counts

### UserService
- `GetProfile` - Get user profile
- `UpdateProfile` - Update user profile
- `CreateUserIfNotExists` - Idempotent user creation

### PlaylistService
- `CreatePlaylist` - Create new playlist
- `GetPlaylist` - Get playlist with tracks
- `UpdatePlaylist` - Update playlist details
- `DeletePlaylist` - Delete playlist and track associations
- `ListPlaylists` - Paginated playlist listing
- `AddTracks` - Add tracks to playlist at position
- `RemoveTracks` - Remove tracks from playlist

### TagService
- `CreateTag` - Create new tag (normalized to lowercase)
- `GetTag` - Get tag details (case-insensitive lookup)
- `UpdateTag` - Update tag (including rename, normalized)
- `DeleteTag` - Delete tag (case-insensitive)
- `ListTags` - List all user tags
- `AddTagsToTrack` - Add tags to track (creates tags if needed, normalized)
- `RemoveTagFromTrack` - Remove tag from track (case-insensitive)
- `GetTracksByTag` - Query tracks by tag (case-insensitive)
- `normalizeTagName` - Helper that converts tag name to lowercase

### UploadService
- `CreatePresignedUpload` - Generate presigned URL for upload
- `ConfirmUpload` - Confirm upload and trigger processing
- `CompleteMultipartUpload` - Complete multipart upload
- `GetUploadStatus` - Get upload status
- `ListUploads` - List upload history
- `ReprocessUpload` - Retry failed upload from step
- `UploadCoverArt` - Generate presigned URL for cover art

### StreamService
- `GetStreamURL` - Get CloudFront signed URL for streaming
- `GetDownloadURL` - Get signed URL for download
- `GetCoverArtURL` - Get signed URL for cover art

### SearchService
- `Search` - Execute full-text search with filters and pagination
  - Validates query is not empty
  - Validates query length (max 500 characters via `MaxQueryLength`)
  - Applies tag filtering via `filterByTags` when tags specified
- `Autocomplete` - Provide search suggestions
- `IndexTrack` - Index a track in the search engine
- `RemoveTrack` - Remove a track from the search index
- `RebuildIndex` - Rebuild the entire search index for a user
- `filterByTags` - Post-filter search results by tags
  - Validates all tags exist (returns NotFoundError if not)
  - Uses AND logic (tracks must have ALL specified tags)
  - Deduplicates and normalizes tag names to lowercase

### TranscodeService
- `StartTranscode` - Create MediaConvert job for HLS transcoding
- `GetTranscodeStatus` - Get status of a MediaConvert job
- `buildJobSettings` - Build MediaConvert job settings for HLS output

### MigrationService
- `MigrateArtists` - Migrate string-based artists to entity model (idempotent)
- `GetMigrationStatus` - Get migration status (not_started, partial, completed)

### EmbeddingService
- `ComposeEmbedText` - Create text representation of track for embedding
- `GenerateTrackEmbedding` - Generate 1024-dim embedding for track metadata
- `GenerateQueryEmbedding` - Generate embedding for search query
- `BatchGenerateEmbeddings` - Batch embed multiple tracks with partial failure handling

### SimilarityService
- `FindSimilarTracks` - Find tracks similar by semantic/features
- `FindMixableTracks` - Find DJ-compatible tracks (BPM + key)
- `CosineSimilarity` - Calculate vector similarity

### Camelot Key Utilities
- `IsKeyCompatible` - Check if two keys can be mixed harmonically
- `GetCompatibleKeys` - Get all compatible keys for a key
- `GetKeyTransition` - Describe the mixing transition type
- `GetBPMCompatibility` - Check BPM compatibility with half/double time

## Dependencies

### Internal
- `github.com/gvasels/personal-music-searchengine/internal/models`
- `github.com/gvasels/personal-music-searchengine/internal/repository`

### External
- `github.com/google/uuid` - UUID generation for new entities

## Key Design Patterns

### Cover Art URLs
All services that return tracks, albums, or playlists generate presigned URLs for cover art. The URL is generated on each request with a 24-hour expiry.

### Pagination
Services use the repository's `PaginatedResult[T]` type with opaque cursors. The cursor is passed through from repository to handler.

### Error Handling
Services convert repository `ErrNotFound` to `models.NewNotFoundError()` for proper API error responses.

### Access Control & Visibility
Track visibility is enforced at the service layer (not just handlers):
- **403 Forbidden**: Returned when user lacks permission to access a track they aren't authorized to view
- **404 Not Found**: Returned only when the resource truly doesn't exist
- The `hasGlobal` parameter determines if the user has admin/global read permissions
- Visibility levels: `private` (owner only), `unlisted` (anyone with link), `public` (discoverable)

### Async Play Count
Play count is incremented asynchronously in a goroutine to avoid blocking the stream URL response.

## Usage Example

```go
// Create services
services := service.NewServices(
    repo,
    s3Repo,
    cloudfront,
    "music-library-media",
    "arn:aws:states:us-east-1:123:stateMachine:upload-processor",
)

// Use a service
track, err := services.Track.GetTrack(ctx, userID, trackID)
if err != nil {
    // Handle error
}
```

## Testing

Services should be tested with mocked repository interfaces. Each service method should have:
- Happy path test
- Not found error test
- Validation error tests
- Edge case tests (empty lists, etc.)
