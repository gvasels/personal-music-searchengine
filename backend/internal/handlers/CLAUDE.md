# Handlers Layer - CLAUDE.md

## Overview

HTTP handlers using Echo framework for the Personal Music Search Engine API. All handlers follow a consistent pattern: extract user ID from context, validate input, call service layer, and return appropriate response.

## File Descriptions

| File | Purpose |
|------|---------|
| `handlers.go` | Main handler struct, route registration, helper functions |
| `user.go` | User profile handlers (GetProfile, UpdateProfile) |
| `track.go` | Track CRUD handlers and tag management |
| `album.go` | Album and artist handlers |
| `playlist.go` | Playlist CRUD and track management |
| `tag.go` | Tag CRUD and track associations |
| `upload.go` | Upload workflow handlers (presigned URLs, confirmation) |
| `stream.go` | Streaming and download URL handlers |
| `search.go` | Search handlers (simple and advanced) |

## Route Registration

All routes are registered under `/api/v1`:

### User Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/me` | GetProfile | Get current user's profile |
| PUT | `/me` | UpdateProfile | Update current user's profile |

### Track Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/tracks` | ListTracks | List tracks with pagination |
| GET | `/tracks/:id` | GetTrack | Get track by ID |
| PUT | `/tracks/:id` | UpdateTrack | Update track metadata |
| DELETE | `/tracks/:id` | DeleteTrack | Delete track |
| POST | `/tracks/:id/tags` | AddTagsToTrack | Add tags to track |
| DELETE | `/tracks/:id/tags/:tag` | RemoveTagFromTrack | Remove tag from track |
| PUT | `/tracks/:id/cover` | UploadCoverArt | Upload cover art |

### Album Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/albums` | ListAlbums | List albums with pagination |
| GET | `/albums/:id` | GetAlbum | Get album with tracks |

### Artist Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/artists` | ListArtists | List artists with counts |
| GET | `/artists/:name/tracks` | ListTracksByArtist | Get tracks by artist |
| GET | `/artists/:name/albums` | ListAlbumsByArtist | Get albums by artist |

### Playlist Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/playlists` | ListPlaylists | List playlists |
| POST | `/playlists` | CreatePlaylist | Create new playlist |
| GET | `/playlists/:id` | GetPlaylist | Get playlist with tracks |
| PUT | `/playlists/:id` | UpdatePlaylist | Update playlist details |
| DELETE | `/playlists/:id` | DeletePlaylist | Delete playlist |
| POST | `/playlists/:id/tracks` | AddTracksToPlaylist | Add tracks to playlist |
| DELETE | `/playlists/:id/tracks` | RemoveTracksFromPlaylist | Remove tracks |

### Tag Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/tags` | ListTags | List user's tags |
| POST | `/tags` | CreateTag | Create new tag |
| GET | `/tags/:name` | GetTag | Get tag details |
| PUT | `/tags/:name` | UpdateTag | Update tag |
| DELETE | `/tags/:name` | DeleteTag | Delete tag |
| GET | `/tags/:name/tracks` | GetTracksByTag | Get tracks with tag |

### Upload Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/upload/presigned` | CreatePresignedUpload | Get presigned URL |
| POST | `/upload/confirm` | ConfirmUpload | Confirm upload |
| POST | `/upload/complete-multipart` | CompleteMultipartUpload | Complete multipart |
| GET | `/uploads` | ListUploads | List upload history |
| GET | `/uploads/:id` | GetUploadStatus | Get upload status |
| POST | `/uploads/:id/reprocess` | ReprocessUpload | Retry failed upload |

### Streaming Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/stream/:trackId` | GetStreamURL | Get streaming URL |
| GET | `/download/:trackId` | GetDownloadURL | Get download URL |

### Search Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/search` | SimpleSearch | Simple text search |
| POST | `/search` | AdvancedSearch | Advanced search with filters |

## Helper Functions

| Function | Purpose |
|----------|---------|
| `getUserIDFromContext` | Extract user ID from API Gateway claims or X-User-ID header |
| `handleError` | Convert errors to appropriate HTTP responses |
| `bindAndValidate` | Bind and validate request body |
| `success` | Return 200 OK with JSON data |
| `created` | Return 201 Created with JSON data |
| `noContent` | Return 204 No Content |

## Dependencies

### Internal
- `github.com/gvasels/personal-music-searchengine/internal/models`
- `github.com/gvasels/personal-music-searchengine/internal/service`

### External
- `github.com/labstack/echo/v4` - HTTP framework

## Authentication

User ID is extracted from the request context in production (API Gateway Cognito JWT authorizer) or from the `X-User-ID` header for local development/testing.

```go
func getUserIDFromContext(c echo.Context) string {
    // Try API Gateway claims first
    if claims := c.Request().Context().Value("claims"); claims != nil {
        if claimsMap, ok := claims.(map[string]interface{}); ok {
            if sub, ok := claimsMap["sub"].(string); ok {
                return sub
            }
        }
    }
    // Fall back to X-User-ID header
    return c.Request().Header.Get("X-User-ID")
}
```

## Error Handling

All errors are converted to `models.APIError` and returned as JSON:

```go
func handleError(c echo.Context, err error) error {
    var apiErr *models.APIError
    if errors.As(err, &apiErr) {
        return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
    }
    return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
}
```

## Usage Example

```go
// Create handlers with services
handlers := handlers.NewHandlers(services)

// Register routes
e := echo.New()
handlers.RegisterRoutes(e)

// Start server (or use Lambda adapter)
e.Start(":8080")
```

## Testing

Handlers should be tested with mocked service interfaces:
- Mock the service methods to return expected responses
- Test successful paths
- Test error handling (not found, validation errors, etc.)
- Test authentication (missing/invalid user ID)
