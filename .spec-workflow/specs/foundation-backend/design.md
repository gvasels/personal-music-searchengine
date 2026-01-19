# Design - Foundation Backend

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway                               │
│              (Cognito JWT Authorizer)                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Lambda Handler                             │
│                  (cmd/api/main.go)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Echo Router                                │
│              (internal/handlers/*)                           │
│   - TrackHandler, AlbumHandler, UserHandler                  │
│   - PlaylistHandler, TagHandler, UploadHandler               │
│   - StreamHandler, SearchHandler                             │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Service Layer                               │
│              (internal/service/*)                            │
│   - TrackService, AlbumService, UserService                  │
│   - PlaylistService, TagService, UploadService               │
│   - StreamService, SearchService                             │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                Repository Layer                              │
│             (internal/repository/*)                          │
│   - Repository interface                                     │
│   - DynamoDBRepository, S3Repository                         │
└─────────────────────────┬───────────────────────────────────┘
                          │
              ┌───────────┴───────────┐
              ▼                       ▼
┌─────────────────────┐   ┌─────────────────────┐
│     DynamoDB        │   │        S3           │
│   (MusicLibrary)    │   │   (media bucket)    │
└─────────────────────┘   └─────────────────────┘
```

## DynamoDB Single-Table Design

**Table**: `MusicLibrary`

| Entity | PK | SK | GSI1PK | GSI1SK |
|--------|----|----|--------|--------|
| User | `USER#{userId}` | `PROFILE` | - | - |
| Track | `USER#{userId}` | `TRACK#{trackId}` | `USER#{userId}#ARTIST#{artist}` | `TRACK#{trackId}` |
| Album | `USER#{userId}` | `ALBUM#{albumId}` | `USER#{userId}#ARTIST#{artist}` | `ALBUM#{year}` |
| Playlist | `USER#{userId}` | `PLAYLIST#{playlistId}` | - | - |
| PlaylistTrack | `PLAYLIST#{playlistId}` | `POSITION#{pos:08d}` | - | - |
| Upload | `USER#{userId}` | `UPLOAD#{uploadId}` | `UPLOAD#STATUS#{status}` | `{timestamp}` |
| Tag | `USER#{userId}` | `TAG#{tagName}` | - | - |
| TrackTag | `USER#{userId}#TRACK#{trackId}` | `TAG#{tagName}` | `USER#{userId}#TAG#{tagName}` | `TRACK#{trackId}` |

## Repository Interface

```go
// Repository defines the data access interface
type Repository interface {
    // Track operations
    CreateTrack(ctx context.Context, track models.Track) error
    GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error)
    UpdateTrack(ctx context.Context, track models.Track) error
    DeleteTrack(ctx context.Context, userID, trackID string) error
    ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResult[models.Track], error)
    ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error)

    // Album operations
    GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error)
    GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error)
    ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResult[models.Album], error)
    ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error)
    UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error

    // User operations
    CreateUser(ctx context.Context, user models.User) error
    GetUser(ctx context.Context, userID string) (*models.User, error)
    UpdateUser(ctx context.Context, user models.User) error
    UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error

    // Playlist operations
    CreatePlaylist(ctx context.Context, playlist models.Playlist) error
    GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error)
    UpdatePlaylist(ctx context.Context, playlist models.Playlist) error
    DeletePlaylist(ctx context.Context, userID, playlistID string) error
    ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResult[models.Playlist], error)
    AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error
    RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error
    GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error)

    // Tag operations
    CreateTag(ctx context.Context, tag models.Tag) error
    GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error)
    UpdateTag(ctx context.Context, tag models.Tag) error
    DeleteTag(ctx context.Context, userID, tagName string) error
    ListTags(ctx context.Context, userID string) ([]models.Tag, error)
    AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error
    RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error
    GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error)
    GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error)

    // Upload operations
    CreateUpload(ctx context.Context, upload models.Upload) error
    GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error)
    UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error
    ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResult[models.Upload], error)
    ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error)
}
```

## S3 Repository Interface

```go
// S3Repository defines media storage operations
type S3Repository interface {
    GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
    GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error)
    DeleteObject(ctx context.Context, key string) error
    CopyObject(ctx context.Context, sourceKey, destKey string) error
    GetObjectMetadata(ctx context.Context, key string) (map[string]string, error)
}
```

## Service Layer Design

### TrackService
```go
type TrackService interface {
    GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error)
    UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error)
    DeleteTrack(ctx context.Context, userID, trackID string) error
    ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResult[models.TrackResponse], error)
}
```

### UploadService
```go
type UploadService interface {
    CreatePresignedUpload(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error)
    ConfirmUpload(ctx context.Context, userID string, req models.ConfirmUploadRequest) (*models.ConfirmUploadResponse, error)
    GetUploadStatus(ctx context.Context, userID, uploadID string) (*models.UploadResponse, error)
    ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResult[models.UploadResponse], error)
}
```

### StreamService
```go
type StreamService interface {
    GetStreamURL(ctx context.Context, userID, trackID string) (*models.StreamResponse, error)
    GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadResponse, error)
}
```

## Handler Design

### No Custom Auth Middleware

**IMPORTANT**: API Gateway provides native Cognito JWT authorization. The Lambda receives the user ID from the API Gateway context, NOT from custom middleware.

```go
// Extract user ID from API Gateway context (Cognito authorizer claims)
func getUserIDFromContext(c echo.Context) string {
    // API Gateway passes claims in request context
    claims := c.Request().Context().Value("claims").(map[string]interface{})
    return claims["sub"].(string)
}
```

### Handler Structure

```go
type Handlers struct {
    trackService    TrackService
    albumService    AlbumService
    userService     UserService
    playlistService PlaylistService
    tagService      TagService
    uploadService   UploadService
    streamService   StreamService
    searchService   SearchService
}

func NewHandlers(
    track TrackService,
    album AlbumService,
    user UserService,
    playlist PlaylistService,
    tag TagService,
    upload UploadService,
    stream StreamService,
    search SearchService,
) *Handlers
```

## API Endpoints

### Track Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/tracks` | ListTracks | List user's tracks with filters |
| GET | `/api/v1/tracks/:id` | GetTrack | Get track details |
| PUT | `/api/v1/tracks/:id` | UpdateTrack | Update track metadata |
| DELETE | `/api/v1/tracks/:id` | DeleteTrack | Delete track |

### Album Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/albums` | ListAlbums | List user's albums |
| GET | `/api/v1/albums/:id` | GetAlbum | Get album details |
| GET | `/api/v1/artists` | ListArtists | List artists with aggregation |

### Upload Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/v1/upload/presigned` | CreatePresignedUpload | Get presigned S3 URL |
| POST | `/api/v1/upload/confirm` | ConfirmUpload | Trigger processing |
| GET | `/api/v1/uploads` | ListUploads | List upload history |
| GET | `/api/v1/uploads/:id` | GetUpload | Get upload status |

### Streaming Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/stream/:trackId` | GetStreamURL | Get CloudFront signed URL |
| GET | `/api/v1/download/:trackId` | GetDownloadURL | Get download URL |

### Playlist Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/playlists` | ListPlaylists | List user's playlists |
| POST | `/api/v1/playlists` | CreatePlaylist | Create playlist |
| GET | `/api/v1/playlists/:id` | GetPlaylist | Get playlist details |
| PUT | `/api/v1/playlists/:id` | UpdatePlaylist | Update playlist |
| DELETE | `/api/v1/playlists/:id` | DeletePlaylist | Delete playlist |
| POST | `/api/v1/playlists/:id/tracks` | AddTracksToPlaylist | Add tracks |
| DELETE | `/api/v1/playlists/:id/tracks` | RemoveTracksFromPlaylist | Remove tracks |

### Tag Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/tags` | ListTags | List user's tags |
| POST | `/api/v1/tags` | CreateTag | Create tag |
| PUT | `/api/v1/tags/:name` | UpdateTag | Update tag |
| DELETE | `/api/v1/tags/:name` | DeleteTag | Delete tag |
| POST | `/api/v1/tracks/:id/tags` | AddTagsToTrack | Add tags to track |
| DELETE | `/api/v1/tracks/:id/tags/:tag` | RemoveTagFromTrack | Remove tag |

### Search Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/search?q=` | SimpleSearch | Basic search |
| POST | `/api/v1/search` | AdvancedSearch | Filtered search |

### User Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/me` | GetProfile | Get user profile |
| PUT | `/api/v1/me` | UpdateProfile | Update profile |

## Error Handling

```go
// APIError represents a structured API error
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

// Common error codes
const (
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeBadRequest     = "BAD_REQUEST"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeForbidden      = "FORBIDDEN"
    ErrCodeConflict       = "CONFLICT"
    ErrCodeInternalError  = "INTERNAL_ERROR"
    ErrCodeStorageLimit   = "STORAGE_LIMIT_EXCEEDED"
)
```

## Testing Strategy

### Unit Tests
- Repository: Mock DynamoDB client
- Service: Mock Repository interface
- Handler: Mock Service interfaces

### Integration Tests
- Repository: Use DynamoDB Local
- End-to-end: Use testcontainers-go

### Test Coverage Target
- Minimum 80% per package
- Critical paths require 100% coverage
