package handlers

import (
	"context"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// UserServiceInterface defines methods for user operations
type UserServiceInterface interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, req models.UpdateUserRequest) (*models.User, error)
}

// TrackServiceInterface defines methods for track operations
type TrackServiceInterface interface {
	GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error)
	ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.TrackResponse], error)
	UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error)
	DeleteTrack(ctx context.Context, userID, trackID string) error
	AddTagsToTrack(ctx context.Context, userID, trackID string, tags []string) (*models.TrackResponse, error)
	RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) (*models.TrackResponse, error)
}

// AlbumServiceInterface defines methods for album operations
type AlbumServiceInterface interface {
	GetAlbumWithTracks(ctx context.Context, userID, albumID string) (*models.AlbumWithTracks, error)
	ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResponse[models.AlbumResponse], error)
	ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*models.PaginatedResponse[models.ArtistSummary], error)
	GetArtist(ctx context.Context, userID, artistName string) (map[string]interface{}, error)
}

// PlaylistServiceInterface defines methods for playlist operations
type PlaylistServiceInterface interface {
	GetPlaylistWithTracks(ctx context.Context, userID, playlistID string) (*models.PlaylistWithTracks, error)
	CreatePlaylist(ctx context.Context, userID string, req models.CreatePlaylistRequest) (*models.PlaylistResponse, error)
	UpdatePlaylist(ctx context.Context, userID, playlistID string, req models.UpdatePlaylistRequest) (*models.PlaylistResponse, error)
	DeletePlaylist(ctx context.Context, userID, playlistID string) error
	ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResponse[models.PlaylistResponse], error)
	AddTracks(ctx context.Context, userID, playlistID string, req models.AddTracksToPlaylistRequest) (*models.PlaylistResponse, error)
	RemoveTracks(ctx context.Context, userID, playlistID string, trackIDs []string) (*models.PlaylistResponse, error)
	ReorderTracks(ctx context.Context, userID, playlistID string, req models.ReorderPlaylistTracksRequest) (*models.PlaylistResponse, error)
}

// TagServiceInterface defines methods for tag operations
type TagServiceInterface interface {
	CreateTag(ctx context.Context, userID string, req models.CreateTagRequest) (*models.Tag, error)
	UpdateTag(ctx context.Context, userID, tagName string, req models.UpdateTagRequest) (*models.Tag, error)
	DeleteTag(ctx context.Context, userID, tagName string) error
	ListTags(ctx context.Context, userID string, filter models.TagFilter) (*models.PaginatedResponse[models.TagResponse], error)
}

// UploadServiceInterface defines methods for upload operations
type UploadServiceInterface interface {
	CreatePresignedUpload(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error)
	ConfirmUpload(ctx context.Context, userID, uploadID string) (*models.ConfirmUploadResponse, error)
	ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.UploadResponse], error)
}

// SearchServiceInterface defines methods for search operations
type SearchServiceInterface interface {
	Search(ctx context.Context, userID string, req models.SearchRequest) (*models.SearchResponse, error)
	GetSuggestions(ctx context.Context, userID, query string) (*models.AutocompleteResponse, error)
}

// StreamServiceInterface defines methods for streaming operations
type StreamServiceInterface interface {
	GetStreamURL(ctx context.Context, userID, trackID, quality string) (*models.StreamResponse, error)
	GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadResponse, error)
	RecordPlayback(ctx context.Context, userID string, req models.RecordPlayRequest) error
	GetQueue(ctx context.Context, userID string) (*models.PlayQueue, error)
	UpdateQueue(ctx context.Context, userID string, req models.UpdateQueueRequest) (*models.PlayQueue, error)
	QueueAction(ctx context.Context, userID string, req models.QueueActionRequest) (*models.PlayQueue, error)
}
