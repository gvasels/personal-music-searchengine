package repository

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// Repository defines the interface for data access
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

// DynamoDBRepository implements Repository using DynamoDB
type DynamoDBRepository struct {
	client    *dynamodb.Client
	s3Client  *s3.Client
	tableName string
}

// NewDynamoDBRepository creates a new DynamoDB repository
func NewDynamoDBRepository() *DynamoDBRepository {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "MusicLibrary"
	}

	return &DynamoDBRepository{
		client:    dynamodb.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
		tableName: tableName,
	}
}

// GetTableName returns the DynamoDB table name
func (r *DynamoDBRepository) GetTableName() string {
	return r.tableName
}

// GetDynamoDBClient returns the DynamoDB client
func (r *DynamoDBRepository) GetDynamoDBClient() *dynamodb.Client {
	return r.client
}

// GetS3Client returns the S3 client
func (r *DynamoDBRepository) GetS3Client() *s3.Client {
	return r.s3Client
}
