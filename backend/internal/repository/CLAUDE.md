# Repository Layer - CLAUDE.md

## Overview

Data access layer implementing DynamoDB single-table design and S3 media storage operations for the Personal Music Search Engine.

## File Descriptions

| File | Purpose |
|------|---------|
| `repository.go` | Interface definitions for Repository, S3Repository, CloudFrontSigner |
| `dynamodb.go` | DynamoDB implementation of Repository interface |
| `s3.go` | S3 implementation of S3Repository interface |

## Key Interfaces

### Repository Interface (`repository.go`)
Primary data access interface for DynamoDB operations:
- Track CRUD operations
- Album operations with artist queries
- User profile management
- Playlist operations with track ordering
- Tag operations with track associations
- Upload status tracking and step management

### S3Repository Interface (`repository.go`)
Media storage operations:
- Presigned URL generation (upload/download)
- Multipart upload support for files > 100MB
- Object operations (delete, copy, metadata)

### CloudFrontSigner Interface (`repository.go`)
Signed URL generation for streaming via CloudFront.

## DynamoDB Key Patterns

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

## Functions

### DynamoDB Repository
| Function | Description |
|----------|-------------|
| `NewDynamoDBRepository` | Creates new DynamoDB repository with client and table name |
| `CreateTrack`, `GetTrack`, `UpdateTrack`, `DeleteTrack` | Track CRUD |
| `ListTracks` | Paginated track listing with cursor |
| `ListTracksByArtist` | Query tracks by artist using GSI1 |
| `GetOrCreateAlbum` | Idempotent album creation |
| `CreateUser`, `GetUser`, `UpdateUser` | User profile operations |
| `UpdateUserStats`, `UpdateAlbumStats` | Stat update operations |
| `CreatePlaylist`, `GetPlaylist`, etc. | Playlist CRUD |
| `AddTracksToPlaylist`, `RemoveTracksFromPlaylist` | Playlist track management |
| `CreateTag`, `AddTagsToTrack`, `GetTracksByTag` | Tag operations |
| `CreateUpload`, `UpdateUploadStatus`, `UpdateUploadStep` | Upload tracking |
| `ListUploadsByStatus` | Query uploads by status using GSI1 |

### S3 Repository
| Function | Description |
|----------|-------------|
| `NewS3Repository` | Creates new S3 repository with client and bucket name |
| `GeneratePresignedUploadURL` | Generate presigned PUT URL |
| `GeneratePresignedDownloadURL` | Generate presigned GET URL |
| `InitiateMultipartUpload` | Start multipart upload |
| `GenerateMultipartUploadURLs` | Generate presigned URLs for all parts |
| `CompleteMultipartUpload` | Complete multipart upload |
| `AbortMultipartUpload` | Abort multipart upload |
| `DeleteObject`, `CopyObject` | Object operations |
| `DeleteByPrefix(ctx, prefix)` | Batch delete all objects with given prefix (used for HLS cleanup) |
| `GetObjectMetadata`, `ObjectExists` | Metadata operations |

### S3Client Interface Methods
```go
// Object operations interface
DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, ...) (*s3.DeleteObjectsOutput, error)
ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, ...) (*s3.ListObjectsV2Output, error)
```

## Error Handling

Common errors defined in `repository.go`:
- `ErrNotFound` - Item not found
- `ErrAlreadyExists` - Item already exists
- `ErrInvalidCursor` - Invalid pagination cursor
- `ErrInvalidInput` - Invalid input

## Pagination

Uses opaque base64-encoded cursors for pagination:
```go
type PaginatedResult[T any] struct {
    Items      []T
    NextCursor string
    HasMore    bool
}
```

Cursors encode `PK`, `SK`, `GSI1PK`, `GSI1SK` from DynamoDB's `LastEvaluatedKey`.

## Dependencies

### External
- `github.com/aws/aws-sdk-go-v2/service/dynamodb`
- `github.com/aws/aws-sdk-go-v2/service/s3`
- `github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue`
- `github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression`

### Internal
- `github.com/gvasels/personal-music-searchengine/internal/models`

## Testing

Repository implementations use client interfaces for testability:
- `DynamoDBClient` - Interface for DynamoDB operations
- `S3Client` - Interface for S3 operations
- `S3PresignClient` - Interface for presigned URL operations

Integration tests should use DynamoDB Local and LocalStack for S3.

## Usage Examples

### Creating a Repository
```go
cfg, _ := config.LoadDefaultConfig(ctx)
dynamoClient := dynamodb.NewFromConfig(cfg)
s3Client := s3.NewFromConfig(cfg)
presignClient := s3.NewPresignClient(s3Client)

repo := repository.NewDynamoDBRepository(dynamoClient, "MusicLibrary")
s3Repo := repository.NewS3Repository(s3Client, presignClient, "music-library-media")
```

### Paginated Queries
```go
filter := models.TrackFilter{Limit: 20, Cursor: ""}
result, err := repo.ListTracks(ctx, userID, filter)
// result.Items, result.NextCursor, result.HasMore
```
