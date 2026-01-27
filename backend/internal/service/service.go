package service

import (
	"context"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// TrackService defines track management operations
type TrackService interface {
	GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error)
	UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error)
	DeleteTrack(ctx context.Context, userID, trackID string) error
	ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error)
	ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.TrackResponse, error)
	IncrementPlayCount(ctx context.Context, userID, trackID string) error
	// Visibility operations
	UpdateVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error
}

// AlbumService defines album management operations
type AlbumService interface {
	GetAlbum(ctx context.Context, userID, albumID string) (*models.AlbumWithTracks, error)
	ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.AlbumResponse], error)
	ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.AlbumResponse, error)
	ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) ([]models.ArtistSummary, error)
}

// UserService defines user profile operations
type UserService interface {
	GetProfile(ctx context.Context, userID string) (*models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateUserRequest) (*models.UserResponse, error)
	CreateUserIfNotExists(ctx context.Context, userID, email, displayName string) (*models.User, error)
	// User settings operations
	GetSettings(ctx context.Context, userID string) (*models.UserSettings, error)
	UpdateSettings(ctx context.Context, userID string, input *UserSettingsUpdateInput) (*models.UserSettings, error)
	// Cognito user creation
	CreateUserFromCognito(ctx context.Context, cognitoSub, email, displayName string) (*models.User, error)
}

// PlaylistService defines playlist operations
type PlaylistService interface {
	CreatePlaylist(ctx context.Context, userID string, req models.CreatePlaylistRequest) (*models.PlaylistResponse, error)
	GetPlaylist(ctx context.Context, userID, playlistID string) (*models.PlaylistWithTracks, error)
	UpdatePlaylist(ctx context.Context, userID, playlistID string, req models.UpdatePlaylistRequest) (*models.PlaylistResponse, error)
	DeletePlaylist(ctx context.Context, userID, playlistID string) error
	ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.PlaylistResponse], error)
	AddTracks(ctx context.Context, userID, playlistID string, req models.AddTracksToPlaylistRequest) (*models.PlaylistResponse, error)
	RemoveTracks(ctx context.Context, userID, playlistID string, req models.RemoveTracksFromPlaylistRequest) (*models.PlaylistResponse, error)
	ReorderTracks(ctx context.Context, userID, playlistID string, req models.ReorderPlaylistTracksRequest) (*models.PlaylistResponse, error)
	// Visibility operations
	UpdateVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error
	ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.PlaylistResponse], error)
}

// TagService defines tag management operations
type TagService interface {
	CreateTag(ctx context.Context, userID string, req models.CreateTagRequest) (*models.TagResponse, error)
	GetTag(ctx context.Context, userID, tagName string) (*models.TagResponse, error)
	UpdateTag(ctx context.Context, userID, tagName string, req models.UpdateTagRequest) (*models.TagResponse, error)
	DeleteTag(ctx context.Context, userID, tagName string) error
	ListTags(ctx context.Context, userID string) ([]models.TagResponse, error)
	AddTagsToTrack(ctx context.Context, userID, trackID string, req models.AddTagsToTrackRequest) ([]string, error)
	RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error
	GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.TrackResponse, error)
}

// UploadService defines upload and processing operations
type UploadService interface {
	CreatePresignedUpload(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error)
	ConfirmUpload(ctx context.Context, userID string, req models.ConfirmUploadRequest) (*models.ConfirmUploadResponse, error)
	CompleteMultipartUpload(ctx context.Context, userID string, req models.CompleteMultipartUploadRequest) (*models.ConfirmUploadResponse, error)
	GetUploadStatus(ctx context.Context, userID, uploadID string) (*models.UploadResponse, error)
	ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.UploadResponse], error)
	ReprocessUpload(ctx context.Context, userID, uploadID string, req models.ReprocessUploadRequest) (*models.UploadResponse, error)
	UploadCoverArt(ctx context.Context, userID, trackID string, req models.CoverArtUploadRequest) (*models.CoverArtUploadResponse, error)
}

// StreamService defines streaming and download operations
type StreamService interface {
	GetStreamURL(ctx context.Context, userID, trackID string) (*models.StreamResponse, error)
	GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadResponse, error)
	GetCoverArtURL(ctx context.Context, userID, trackID string) (string, error)
}

// SearchService defines search operations
type SearchService interface {
	Search(ctx context.Context, userID string, req models.SearchRequest) (*models.SearchResponse, error)
	Autocomplete(ctx context.Context, userID, query string) (*models.AutocompleteResponse, error)
	RemoveTrack(ctx context.Context, trackID string) error
	IndexTrack(ctx context.Context, track models.Track) error
}

// Services holds all service implementations
type Services struct {
	Track    TrackService
	Album    AlbumService
	Artist   ArtistService
	User     UserService
	Playlist PlaylistService
	Tag      TagService
	Upload   UploadService
	Stream   StreamService
	Search   SearchService
	Admin    AdminService
}

// NewServices creates a new Services instance with all dependencies
func NewServices(
	repo repository.Repository,
	s3Repo repository.S3Repository,
	cloudfront repository.CloudFrontSigner,
	mediaBucket string,
	stepFunctionsARN string,
) *Services {
	return &Services{
		Track:    NewTrackService(repo, s3Repo),
		Album:    NewAlbumService(repo, s3Repo),
		Artist:   NewArtistService(repo, s3Repo),
		User:     NewUserService(repo),
		Playlist: NewPlaylistService(repo, s3Repo),
		Tag:      NewTagService(repo),
		Upload:   NewUploadService(repo, s3Repo, mediaBucket, stepFunctionsARN),
		Stream:   NewStreamService(repo, cloudfront, s3Repo),
		// Search service requires Nixiesearch client - initialized separately
	}
}
