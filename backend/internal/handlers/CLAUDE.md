# HTTP Handlers - CLAUDE.md

## Overview

HTTP request handlers using Echo framework. Handlers parse requests, call services, and format responses. User authentication is handled by API Gateway Cognito Authorizer - handlers extract user ID from JWT claims.

## File Descriptions

| File | Purpose |
|------|---------|
| `handlers.go` | Handlers struct and dependency injection |
| `user.go` | User profile endpoints |
| `track.go` | Track CRUD endpoints |
| `album.go` | Album and artist endpoints |
| `playlist.go` | Playlist management endpoints |
| `tag.go` | Tag management endpoints |
| `upload.go` | Upload presigned URL and confirmation |
| `search.go` | Search and autocomplete endpoints |
| `stream.go` | Streaming and download endpoints |

## Key Types

### Handlers (`handlers.go`)
```go
type Handlers struct {
    userService     *service.UserService
    trackService    *service.TrackService
    albumService    *service.AlbumService
    playlistService *service.PlaylistService
    tagService      *service.TagService
    uploadService   *service.UploadService
    searchService   *service.SearchService
    streamService   *service.StreamService
}

func NewHandlers() *Handlers
```

## Endpoints

### User (`user.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/me` | `GetMe` | Get current user profile |
| PUT | `/api/v1/me` | `UpdateMe` | Update user profile |

### Track (`track.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/tracks` | `ListTracks` | List tracks with filters |
| GET | `/api/v1/tracks/:id` | `GetTrack` | Get track by ID |
| PUT | `/api/v1/tracks/:id` | `UpdateTrack` | Update track metadata |
| DELETE | `/api/v1/tracks/:id` | `DeleteTrack` | Delete track |
| POST | `/api/v1/tracks/:id/tags` | `AddTagsToTrack` | Add tags to track |
| DELETE | `/api/v1/tracks/:id/tags/:tag` | `RemoveTagFromTrack` | Remove tag |

### Album (`album.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/albums` | `ListAlbums` | List albums |
| GET | `/api/v1/albums/:id` | `GetAlbum` | Get album with tracks |
| GET | `/api/v1/artists` | `ListArtists` | List artists |
| GET | `/api/v1/artists/:name` | `GetArtist` | Get artist with albums |

### Playlist (`playlist.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/playlists` | `ListPlaylists` | List playlists |
| POST | `/api/v1/playlists` | `CreatePlaylist` | Create playlist |
| GET | `/api/v1/playlists/:id` | `GetPlaylist` | Get playlist with tracks |
| PUT | `/api/v1/playlists/:id` | `UpdatePlaylist` | Update playlist |
| DELETE | `/api/v1/playlists/:id` | `DeletePlaylist` | Delete playlist |
| POST | `/api/v1/playlists/:id/tracks` | `AddPlaylistTracks` | Add tracks |
| DELETE | `/api/v1/playlists/:id/tracks` | `RemovePlaylistTracks` | Remove tracks |

### Tag (`tag.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/tags` | `ListTags` | List user tags |
| POST | `/api/v1/tags` | `CreateTag` | Create tag |
| PUT | `/api/v1/tags/:name` | `UpdateTag` | Update tag |
| DELETE | `/api/v1/tags/:name` | `DeleteTag` | Delete tag |

### Upload (`upload.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/v1/upload/presigned` | `GetPresignedURL` | Get S3 presigned URL |
| POST | `/api/v1/upload/confirm` | `ConfirmUpload` | Confirm and process |
| GET | `/api/v1/uploads` | `ListUploads` | List uploads |

### Search (`search.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/search` | `Search` | Full-text search |
| POST | `/api/v1/search` | `AdvancedSearch` | Search with filters |
| GET | `/api/v1/search/suggest` | `GetSuggestions` | Autocomplete |

### Stream (`stream.go`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/api/v1/stream/:trackId` | `GetStreamURL` | Get streaming URL |
| GET | `/api/v1/download/:trackId` | `GetDownloadURL` | Get download URL |
| POST | `/api/v1/playback/record` | `RecordPlayback` | Record play event |

## Authentication

**IMPORTANT**: Authentication is handled by API Gateway Cognito Authorizer, NOT custom middleware.

User ID extraction from JWT claims:
```go
func getUserID(c echo.Context) string {
    // JWT claims are added by API Gateway authorizer
    claims := c.Get("claims").(map[string]interface{})
    return claims["sub"].(string)
}
```

## Error Handling

All handlers use consistent error responses:
```go
func (h *Handlers) handleError(c echo.Context, err error) error {
    switch {
    case errors.Is(err, models.ErrNotFound):
        return c.JSON(http.StatusNotFound, models.ErrorResponse{
            Code:    "NOT_FOUND",
            Message: err.Error(),
        })
    case errors.Is(err, models.ErrValidation):
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Code:    "VALIDATION_ERROR",
            Message: err.Error(),
        })
    default:
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Code:    "INTERNAL_ERROR",
            Message: "An unexpected error occurred",
        })
    }
}
```

## Dependencies

### External
- `github.com/labstack/echo/v4` - HTTP framework

### Internal
- `internal/service` - Business logic
- `internal/models` - Request/response types

## Testing

Handler tests use Echo's test utilities:

```go
func TestGetTrack(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPath("/api/v1/tracks/:id")
    c.SetParamNames("id")
    c.SetParamValues("track-123")

    // Mock service
    mockService := &MockTrackService{}
    mockService.On("GetTrack", mock.Anything, "user-123", "track-123").Return(&models.TrackResponse{...}, nil)

    h := &Handlers{trackService: mockService}

    err := h.GetTrack(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```
