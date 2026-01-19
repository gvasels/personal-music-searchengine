# Service Layer - CLAUDE.md

## Overview

Business logic layer implementing application use cases. Services orchestrate repository operations, apply business rules, and transform data between domain models and API responses.

## File Descriptions

| File | Purpose |
|------|---------|
| `user_service.go` | User profile management |
| `track_service.go` | Track CRUD, tagging operations |
| `album_service.go` | Album queries, artist aggregation |
| `playlist_service.go` | Playlist management, track ordering |
| `tag_service.go` | Tag CRUD, track-tag associations |
| `upload_service.go` | Presigned URLs, Step Functions trigger |
| `search_service.go` | Full-text search via Nixiesearch |
| `stream_service.go` | CloudFront signed URLs for streaming |

## Key Types

### TrackService (`track_service.go`)
```go
type TrackService struct {
    repo        *repository.DynamoDBRepository
    mediaBucket string
}

func NewTrackService(repo *repository.DynamoDBRepository) *TrackService
func (s *TrackService) GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error)
func (s *TrackService) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.TrackResponse], error)
func (s *TrackService) UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error)
func (s *TrackService) DeleteTrack(ctx context.Context, userID, trackID string) error
func (s *TrackService) AddTagsToTrack(ctx context.Context, userID, trackID string, tags []string) (*models.TrackResponse, error)
func (s *TrackService) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) (*models.TrackResponse, error)
```

### UploadService (`upload_service.go`)
```go
type UploadService struct {
    repo              *repository.DynamoDBRepository
    mediaBucket       string
    stepFunctionsARN  string
}

func NewUploadService(repo *repository.DynamoDBRepository) *UploadService
func (s *UploadService) GetPresignedURL(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error)
func (s *UploadService) ConfirmUpload(ctx context.Context, userID, uploadID string) (*models.ConfirmUploadResponse, error)
func (s *UploadService) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.UploadResponse], error)
```

### StreamService (`stream_service.go`)
```go
type StreamService struct {
    repo             *repository.DynamoDBRepository
    mediaBucket      string
    cloudFrontDomain string
    keyPairID        string
    privateKey       string
}

func NewStreamService(repo *repository.DynamoDBRepository) *StreamService
func (s *StreamService) GetStreamURL(ctx context.Context, userID, trackID string) (*models.StreamURLResponse, error)
func (s *StreamService) GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadURLResponse, error)
func (s *StreamService) RecordPlayback(ctx context.Context, userID, trackID string, duration int, completed bool) error
```

### PlaylistService (`playlist_service.go`)
```go
type PlaylistService struct {
    repo *repository.DynamoDBRepository
}

func NewPlaylistService(repo *repository.DynamoDBRepository) *PlaylistService
func (s *PlaylistService) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.PlaylistResponse, error)
func (s *PlaylistService) CreatePlaylist(ctx context.Context, userID string, req models.CreatePlaylistRequest) (*models.PlaylistResponse, error)
func (s *PlaylistService) UpdatePlaylist(ctx context.Context, userID, playlistID string, req models.UpdatePlaylistRequest) (*models.PlaylistResponse, error)
func (s *PlaylistService) DeletePlaylist(ctx context.Context, userID, playlistID string) error
func (s *PlaylistService) AddTracks(ctx context.Context, userID, playlistID string, trackIDs []string, position *int) (*models.PlaylistResponse, error)
func (s *PlaylistService) RemoveTracks(ctx context.Context, userID, playlistID string, trackIDs []string) (*models.PlaylistResponse, error)
func (s *PlaylistService) ReorderTracks(ctx context.Context, userID, playlistID string, trackIDs []string) (*models.PlaylistResponse, error)
```

## Business Rules

### Track Operations
- Tracks are always scoped to a user (multi-tenant)
- Deleting a track removes it from S3 and search index
- Tags are automatically created if they don't exist

### Upload Processing
- Presigned URLs expire after 15 minutes
- Maximum file size: 500MB
- Supported formats: MP3, FLAC, WAV, AAC, OGG
- Confirm triggers Step Functions workflow

### Streaming
- CloudFront signed URLs expire after 4 hours
- Playback is recorded for analytics
- Download URLs include original filename

### Playlists
- Track positions are zero-padded for lexicographic ordering
- Reordering recreates all position entries atomically

## Dependencies

### External
- `github.com/aws/aws-sdk-go-v2/service/s3` - Presigned URLs
- `github.com/aws/aws-sdk-go-v2/service/sfn` - Step Functions
- `github.com/aws/aws-sdk-go-v2/service/cloudfront` - Signed URLs

### Internal
- `internal/repository` - Data access
- `internal/models` - Domain models

## Usage Examples

### Get Track with Cover Art URL
```go
trackService := service.NewTrackService(repo)
track, err := trackService.GetTrack(ctx, userID, trackID)
// track.CoverArtURL is populated if cover art exists
```

### Trigger Upload Processing
```go
uploadService := service.NewUploadService(repo)
response, err := uploadService.ConfirmUpload(ctx, userID, uploadID)
// Step Functions execution started
```

### Generate Stream URL
```go
streamService := service.NewStreamService(repo)
stream, err := streamService.GetStreamURL(ctx, userID, trackID)
// stream.StreamURL is CloudFront signed URL
```

## Testing

Service tests should mock the repository:

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
    args := m.Called(ctx, userID, trackID)
    return args.Get(0).(*models.Track), args.Error(1)
}
```

Test coverage should include:
- Happy path for all operations
- Error handling (not found, validation, AWS errors)
- Business rule enforcement
