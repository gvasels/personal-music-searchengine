package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// Common repository errors
var (
	ErrNotFound      = errors.New("item not found")
	ErrAlreadyExists = errors.New("item already exists")
	ErrInvalidCursor = errors.New("invalid pagination cursor")
	ErrInvalidInput  = errors.New("invalid input")
)

// PaginatedResult represents a paginated query result
type PaginatedResult[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// Repository defines the data access interface for DynamoDB operations
type Repository interface {
	// Track operations
	CreateTrack(ctx context.Context, track models.Track) error
	GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error)
	UpdateTrack(ctx context.Context, track models.Track) error
	DeleteTrack(ctx context.Context, userID, trackID string) error
	ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*PaginatedResult[models.Track], error)
	ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error)

	// Album operations
	GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error)
	GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error)
	ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*PaginatedResult[models.Album], error)
	ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error)
	UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error

	// Artist operations
	CreateArtist(ctx context.Context, artist models.Artist) error
	GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error)
	GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error)
	ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*PaginatedResult[models.Artist], error)
	UpdateArtist(ctx context.Context, artist models.Artist) error
	DeleteArtist(ctx context.Context, userID, artistID string) error
	BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error)
	SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error)
	GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error)
	GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error)
	GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error)

	// User operations
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error
	UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error
	ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*PaginatedResult[models.User], error)

	// Playlist operations
	CreatePlaylist(ctx context.Context, playlist models.Playlist) error
	GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlist models.Playlist) error
	DeletePlaylist(ctx context.Context, userID, playlistID string) error
	ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*PaginatedResult[models.Playlist], error)
	SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error)
	AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error
	RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error
	GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error)
	ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error
	UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error
	ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*PaginatedResult[models.Playlist], error)

	// ArtistProfile operations
	CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error
	GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error)
	UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error
	DeleteArtistProfile(ctx context.Context, userID string) error
	ListArtistProfiles(ctx context.Context, limit int, cursor string) (*PaginatedResult[models.ArtistProfile], error)
	IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error

	// Follow operations
	CreateFollow(ctx context.Context, follow models.Follow) error
	DeleteFollow(ctx context.Context, followerID, followedID string) error
	GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error)
	ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*PaginatedResult[models.Follow], error)
	ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*PaginatedResult[models.Follow], error)
	IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error

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
	UpdateUpload(ctx context.Context, upload models.Upload) error
	UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error
	UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error
	ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*PaginatedResult[models.Upload], error)
	ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error)
}

// S3Repository defines media storage operations
type S3Repository interface {
	// Presigned URL operations
	GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, expiry time.Duration, filename string) (string, error)

	// Multipart upload operations
	InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error)
	GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error)
	CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error
	AbortMultipartUpload(ctx context.Context, key, uploadID string) error

	// Object operations
	DeleteObject(ctx context.Context, key string) error
	CopyObject(ctx context.Context, sourceKey, destKey string) error
	GetObjectMetadata(ctx context.Context, key string) (map[string]string, error)
	ObjectExists(ctx context.Context, key string) (bool, error)
}

// CloudFrontSigner defines signed URL operations for streaming
type CloudFrontSigner interface {
	GenerateSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	GenerateSignedDownloadURL(ctx context.Context, key string, expiry time.Duration, filename string) (string, error)
}
