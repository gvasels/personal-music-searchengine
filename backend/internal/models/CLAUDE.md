# Domain Models - CLAUDE.md

## Overview

Domain models, DTOs (Data Transfer Objects), and constants for the Personal Music Search Engine. These models define the core data structures used throughout the application.

## File Descriptions

| File | Purpose |
|------|---------|
| `common.go` | Shared types: `EntityType`, `UploadStatus`, `AudioFormat`, `Timestamps`, `DynamoDBItem`, `Pagination` |
| `user.go` | User model and profile DTOs |
| `track.go` | Track model, create/update requests, filter options |
| `album.go` | Album model, artist aggregation |
| `playlist.go` | Playlist and PlaylistTrack models |
| `tag.go` | Tag and TrackTag models |
| `upload.go` | Upload tracking, presigned URL requests/responses |
| `search.go` | Search request/response, Nixiesearch types |
| `streaming.go` | Stream/download URLs, playback queue |
| `errors.go` | API error types and formatting |

## Key Types

### Entity Types (`common.go`)
```go
type EntityType string
const (
    EntityUser          EntityType = "USER"
    EntityTrack         EntityType = "TRACK"
    EntityAlbum         EntityType = "ALBUM"
    EntityPlaylist      EntityType = "PLAYLIST"
    EntityPlaylistTrack EntityType = "PLAYLIST_TRACK"
    EntityUpload        EntityType = "UPLOAD"
    EntityTag           EntityType = "TAG"
    EntityTrackTag      EntityType = "TRACK_TAG"
)
```

### Upload Status (`common.go`)
```go
type UploadStatus string
const (
    UploadStatusPending    UploadStatus = "PENDING"
    UploadStatusProcessing UploadStatus = "PROCESSING"
    UploadStatusCompleted  UploadStatus = "COMPLETED"
    UploadStatusFailed     UploadStatus = "FAILED"
)
```

### Audio Formats (`common.go`)
```go
type AudioFormat string
const (
    AudioFormatMP3  AudioFormat = "MP3"
    AudioFormatFLAC AudioFormat = "FLAC"
    AudioFormatWAV  AudioFormat = "WAV"
    AudioFormatAAC  AudioFormat = "AAC"
    AudioFormatOGG  AudioFormat = "OGG"
)
```

### DynamoDB Base Item (`common.go`)
```go
type DynamoDBItem struct {
    PK     string `dynamodbav:"PK"`
    SK     string `dynamodbav:"SK"`
    GSI1PK string `dynamodbav:"GSI1PK,omitempty"`
    GSI1SK string `dynamodbav:"GSI1SK,omitempty"`
    Type   string `dynamodbav:"Type"`
}
```

### Pagination Cursor (`common.go`)
```go
// PaginationCursor represents the internal structure of a pagination cursor
// Encoded to base64 for opaque client-side handling
type PaginationCursor struct {
    PK     string `json:"pk"`
    SK     string `json:"sk"`
    GSI1PK string `json:"gsi1pk,omitempty"`
    GSI1SK string `json:"gsi1sk,omitempty"`
}
```

### Track (`track.go`)
```go
type Track struct {
    ID          string      `json:"id"`
    UserID      string      `json:"userId"`
    Title       string      `json:"title"`
    Artist      string      `json:"artist"`
    Album       string      `json:"album,omitempty"`
    Genre       string      `json:"genre,omitempty"`
    Year        int         `json:"year,omitempty"`
    Duration    int         `json:"duration"` // seconds
    Format      AudioFormat `json:"format"`
    FileSize    int64       `json:"fileSize"`
    S3Key       string      `json:"s3Key"`
    CoverArtKey string      `json:"coverArtKey,omitempty"`
    PlayCount   int         `json:"playCount"`
    Tags        []string    `json:"tags,omitempty"`
    Timestamps
}
```

### Upload (`upload.go`)
```go
type Upload struct {
    ID          string       `json:"id"`
    UserID      string       `json:"userId"`
    FileName    string       `json:"fileName"`
    FileSize    int64        `json:"fileSize"`         // Max 1GB
    ContentType string       `json:"contentType"`
    S3Key       string       `json:"s3Key"`
    Status      UploadStatus `json:"status"`
    ErrorMsg    string       `json:"errorMsg,omitempty"`
    TrackID     string       `json:"trackId,omitempty"`
    Timestamps
    CompletedAt *time.Time `json:"completedAt,omitempty"`

    // Step tracking for partial success recovery
    MetadataExtracted bool
    CoverArtExtracted bool
    TrackCreated      bool
    Indexed           bool
    FileMoved         bool

    // Multipart upload tracking
    IsMultipart   bool
    MultipartID   string
    PartsUploaded int
    TotalParts    int
}
```

### Processing Steps (`upload.go`)
```go
type ProcessingStep string
const (
    StepExtractMetadata ProcessingStep = "extract_metadata"
    StepExtractCover    ProcessingStep = "extract_cover"
    StepCreateTrack     ProcessingStep = "create_track"
    StepIndex           ProcessingStep = "index"
    StepMoveFile        ProcessingStep = "move_file"
)
```

## Functions

### Common Functions (`common.go`)
| Function | Signature | Description |
|----------|-----------|-------------|
| `EncodeCursor` | `(cursor PaginationCursor) string` | Encodes cursor to base64 string |
| `DecodeCursor` | `(encoded string) (PaginationCursor, error)` | Decodes base64 string to cursor |
| `NewPaginationCursor` | `(pk, sk string) PaginationCursor` | Creates cursor with PK/SK |
| `NewPaginationCursorWithGSI` | `(pk, sk, gsi1pk, gsi1sk string) PaginationCursor` | Creates cursor with GSI keys |

### Track Functions (`track.go`)
| Function | Signature | Description |
|----------|-----------|-------------|
| `NewTrackItem` | `(track Track) TrackItem` | Creates DynamoDB item with PK/SK |
| `ToResponse` | `(t *Track) TrackResponse` | Converts to API response |
| `formatDuration` | `(seconds int) string` | Formats duration as "M:SS" |
| `formatFileSize` | `(bytes int64) string` | Formats size as "X.XX MB" |

### Upload Functions (`upload.go`)
| Function | Signature | Description |
|----------|-----------|-------------|
| `NewUploadItem` | `(upload Upload) UploadItem` | Creates DynamoDB item with status GSI |
| `ToResponse` | `(u *Upload) UploadResponse` | Converts to API response with step tracking |

### Error Functions (`errors.go`)
| Function | Signature | Description |
|----------|-----------|-------------|
| `NewAPIError` | `(code, message string, statusCode int) *APIError` | Creates custom API error |
| `NewValidationError` | `(details any) *APIError` | Creates validation error with details |
| `NewNotFoundError` | `(resource, id string) *APIError` | Creates not found error for resource |
| `NewConflictError` | `(message string) *APIError` | Creates conflict error |

## Dependencies

- `time` - Standard library timestamps
- `fmt` - String formatting

No external dependencies - these are pure domain models.

## Usage Examples

### Creating a Track
```go
track := models.Track{
    ID:       uuid.New().String(),
    UserID:   userID,
    Title:    "Song Title",
    Artist:   "Artist Name",
    Duration: 180,
    Format:   models.AudioFormatMP3,
    FileSize: 5242880,
}
item := models.NewTrackItem(track)
```

### Creating an Upload
```go
upload := models.Upload{
    ID:          uuid.New().String(),
    UserID:      userID,
    FileName:    "song.mp3",
    FileSize:    5242880,
    ContentType: "audio/mpeg",
    Status:      models.UploadStatusPending,
}
```

## Testing

Models should have comprehensive unit tests for:
- Struct field validation
- Helper function outputs (formatDuration, formatFileSize)
- DynamoDB item creation (key patterns)
- Response conversion methods
