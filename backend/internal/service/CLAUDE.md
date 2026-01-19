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
| `tag.go` | TagService - tag management and track associations |
| `upload.go` | UploadService - upload workflow and presigned URLs |
| `stream.go` | StreamService - streaming and download URL generation |
| `search.go` | SearchService - Nixiesearch integration for full-text search |
| `search_test.go` | Unit tests for SearchService |
| `transcode.go` | TranscodeService - MediaConvert HLS transcoding |
| `transcode_test.go` | Unit tests for TranscodeService |

## Service Interfaces

### TrackService
- `GetTrack` - Get track with cover art URL
- `UpdateTrack` - Update track metadata
- `DeleteTrack` - Delete track and S3 files
- `ListTracks` - Paginated track listing with cover art URLs
- `ListTracksByArtist` - Query tracks by artist
- `IncrementPlayCount` - Update play count and last played

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
- `CreateTag` - Create new tag
- `GetTag` - Get tag details
- `UpdateTag` - Update tag (including rename)
- `DeleteTag` - Delete tag
- `ListTags` - List all user tags
- `AddTagsToTrack` - Add tags to track (creates tags if needed)
- `RemoveTagFromTrack` - Remove tag from track
- `GetTracksByTag` - Query tracks by tag

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
- `Autocomplete` - Provide search suggestions
- `IndexTrack` - Index a track in the search engine
- `RemoveTrack` - Remove a track from the search index
- `RebuildIndex` - Rebuild the entire search index for a user

### TranscodeService
- `StartTranscode` - Create MediaConvert job for HLS transcoding
- `GetTranscodeStatus` - Get status of a MediaConvert job
- `buildJobSettings` - Build MediaConvert job settings for HLS output

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
