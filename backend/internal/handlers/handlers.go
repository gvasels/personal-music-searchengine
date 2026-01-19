package handlers

import (
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
)

// Handlers holds all API handlers
type Handlers struct {
	userService     UserServiceInterface
	trackService    TrackServiceInterface
	albumService    AlbumServiceInterface
	playlistService PlaylistServiceInterface
	tagService      TagServiceInterface
	uploadService   UploadServiceInterface
	searchService   SearchServiceInterface
	streamService   StreamServiceInterface
}

// NewHandlers creates a new Handlers instance with all dependencies
func NewHandlers() *Handlers {
	// Initialize repositories
	repo := repository.NewDynamoDBRepository()

	// Get S3 client for services that need it
	s3Client := repo.GetS3Client()

	// Initialize services
	return &Handlers{
		userService:     service.NewUserService(repo),
		trackService:    service.NewTrackService(repo),
		albumService:    service.NewAlbumService(repo),
		playlistService: service.NewPlaylistService(repo),
		tagService:      service.NewTagService(repo),
		uploadService:   service.NewUploadService(repo, s3Client),
		searchService:   service.NewSearchService(repo),
		streamService:   service.NewStreamService(repo, s3Client),
	}
}

// NewHandlersWithServices creates a new Handlers instance with injected services (for testing)
func NewHandlersWithServices(
	userService UserServiceInterface,
	trackService TrackServiceInterface,
	albumService AlbumServiceInterface,
	playlistService PlaylistServiceInterface,
	tagService TagServiceInterface,
	uploadService UploadServiceInterface,
	searchService SearchServiceInterface,
	streamService StreamServiceInterface,
) *Handlers {
	return &Handlers{
		userService:     userService,
		trackService:    trackService,
		albumService:    albumService,
		playlistService: playlistService,
		tagService:      tagService,
		uploadService:   uploadService,
		searchService:   searchService,
		streamService:   streamService,
	}
}
