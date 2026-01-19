# Repository Layer - CLAUDE.md

## Overview

Data access layer implementing the repository pattern for DynamoDB and S3 operations. Provides an abstraction over AWS services to enable testing and separation of concerns.

## File Descriptions

| File | Purpose |
|------|---------|
| `repository.go` | Interface definition and DynamoDBRepository factory |
| `dynamodb_user.go` | User CRUD operations |
| `dynamodb_track.go` | Track CRUD and query operations |
| `dynamodb_album.go` | Album CRUD and artist queries |
| `dynamodb_playlist.go` | Playlist and PlaylistTrack operations |
| `dynamodb_tag.go` | Tag and TrackTag operations |
| `dynamodb_upload.go` | Upload tracking operations |

## Key Types

### Repository Interface (`repository.go`)
```go
type Repository interface {
    // User operations
    GetUser(ctx context.Context, userID string) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    UpdateUser(ctx context.Context, user *models.User) error

    // Track operations
    GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error)
    CreateTrack(ctx context.Context, track *models.Track) error
    UpdateTrack(ctx context.Context, track *models.Track) error
    DeleteTrack(ctx context.Context, userID, trackID string) error
    ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.Track], error)
    ListTracksByAlbum(ctx context.Context, userID, albumID string) ([]models.Track, error)
    ListTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error)

    // Album operations
    GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error)
    CreateAlbum(ctx context.Context, album *models.Album) error
    UpdateAlbum(ctx context.Context, album *models.Album) error
    DeleteAlbum(ctx context.Context, userID, albumID string) error
    ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResponse[models.Album], error)
    ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error)
    GetOrCreateAlbum(ctx context.Context, userID, title, artist string, year int) (*models.Album, error)

    // Playlist operations
    GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error)
    CreatePlaylist(ctx context.Context, playlist *models.Playlist) error
    UpdatePlaylist(ctx context.Context, playlist *models.Playlist) error
    DeletePlaylist(ctx context.Context, userID, playlistID string) error
    ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResponse[models.Playlist], error)
    GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error)
    AddPlaylistTrack(ctx context.Context, pt *models.PlaylistTrack) error
    RemovePlaylistTrack(ctx context.Context, playlistID, trackID string) error
    ReorderPlaylistTracks(ctx context.Context, playlistID string, trackIDs []string) error

    // Tag operations
    GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error)
    CreateTag(ctx context.Context, tag *models.Tag) error
    UpdateTag(ctx context.Context, tag *models.Tag) error
    DeleteTag(ctx context.Context, userID, tagName string) error
    ListTags(ctx context.Context, userID string, filter models.TagFilter) (*models.PaginatedResponse[models.Tag], error)
    AddTrackTag(ctx context.Context, tt *models.TrackTag) error
    RemoveTrackTag(ctx context.Context, userID, trackID, tagName string) error

    // Upload operations
    GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error)
    CreateUpload(ctx context.Context, upload *models.Upload) error
    UpdateUpload(ctx context.Context, upload *models.Upload) error
    ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.Upload], error)
}
```

### DynamoDBRepository (`repository.go`)
```go
type DynamoDBRepository struct {
    client    *dynamodb.Client
    s3Client  *s3.Client
    tableName string
}

func NewDynamoDBRepository() *DynamoDBRepository
func (r *DynamoDBRepository) GetTableName() string
func (r *DynamoDBRepository) GetDynamoDBClient() *dynamodb.Client
func (r *DynamoDBRepository) GetS3Client() *s3.Client
```

## DynamoDB Key Patterns

| Entity | PK | SK | GSI1PK | GSI1SK |
|--------|----|----|--------|--------|
| User | `USER#{userId}` | `PROFILE` | - | - |
| Track | `USER#{userId}` | `TRACK#{trackId}` | `USER#{userId}#ARTIST#{artist}` | `TRACK#{trackId}` |
| Album | `USER#{userId}` | `ALBUM#{albumId}` | `USER#{userId}#ARTIST#{artist}` | `ALBUM#{year}` |
| Playlist | `USER#{userId}` | `PLAYLIST#{playlistId}` | - | - |
| PlaylistTrack | `PLAYLIST#{playlistId}` | `POSITION#{pos}` | - | - |
| Upload | `USER#{userId}` | `UPLOAD#{uploadId}` | `UPLOAD#STATUS#{status}` | `{timestamp}` |
| Tag | `USER#{userId}` | `TAG#{tagName}` | - | - |
| TrackTag | `USER#{userId}#TRACK#{trackId}` | `TAG#{tagName}` | `USER#{userId}#TAG#{tagName}` | `TRACK#{trackId}` |

## Dependencies

### External
- `github.com/aws/aws-sdk-go-v2/service/dynamodb` - DynamoDB client
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 client
- `github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue` - Marshaling
- `github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression` - Query building

### Internal
- `internal/models` - Domain models

## Usage Examples

### Initialize Repository
```go
repo := repository.NewDynamoDBRepository()
```

### Query Tracks
```go
filter := models.TrackFilter{
    Artist: "Artist Name",
    Limit:  20,
}
result, err := repo.ListTracks(ctx, userID, filter)
```

### Create Upload
```go
upload := &models.Upload{
    ID:       uuid.New().String(),
    UserID:   userID,
    Status:   models.UploadStatusPending,
}
err := repo.CreateUpload(ctx, upload)
```

## Testing

Repository tests require DynamoDB Local:

```bash
# Start DynamoDB Local
docker run -p 8000:8000 amazon/dynamodb-local

# Run integration tests
AWS_REGION=us-east-1 go test -v ./internal/repository/...
```

Test files should mock the AWS SDK clients or use DynamoDB Local for integration tests.
