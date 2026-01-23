package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// ArtistService defines artist management operations
type ArtistService interface {
	CreateArtist(ctx context.Context, userID string, req models.CreateArtistRequest) (*models.ArtistResponse, error)
	GetArtist(ctx context.Context, userID, artistID string) (*models.ArtistWithStatsResponse, error)
	UpdateArtist(ctx context.Context, userID, artistID string, req models.UpdateArtistRequest) (*models.ArtistResponse, error)
	DeleteArtist(ctx context.Context, userID, artistID string) error
	ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.ArtistResponse], error)
	SearchArtists(ctx context.Context, userID, query string, limit int) ([]models.ArtistResponse, error)
	GetArtistTracks(ctx context.Context, userID, artistID string) ([]models.TrackResponse, error)
}

// ArtistRepository defines the repository interface for artist operations
type ArtistRepository interface {
	CreateArtist(ctx context.Context, artist models.Artist) error
	GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error)
	GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error)
	ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error)
	UpdateArtist(ctx context.Context, artist models.Artist) error
	DeleteArtist(ctx context.Context, userID, artistID string) error
	BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error)
	SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error)
	GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error)
	GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error)
	GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error)
	ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error)
}

// artistService implements ArtistService
type artistService struct {
	artistRepo ArtistRepository
	s3Repo     repository.S3Repository
}

// NewArtistService creates a new ArtistService
func NewArtistService(artistRepo ArtistRepository, s3Repo repository.S3Repository) ArtistService {
	return &artistService{
		artistRepo: artistRepo,
		s3Repo:     s3Repo,
	}
}

// CreateArtist creates a new artist
func (s *artistService) CreateArtist(ctx context.Context, userID string, req models.CreateArtistRequest) (*models.ArtistResponse, error) {
	// Generate sort name if not provided
	sortName := req.SortName
	if sortName == "" {
		sortName = models.GenerateSortName(req.Name)
	}

	artist := models.Artist{
		ID:            uuid.New().String(),
		UserID:        userID,
		Name:          req.Name,
		SortName:      sortName,
		Bio:           req.Bio,
		ImageURL:      req.ImageURL,
		ExternalLinks: req.ExternalLinks,
		IsActive:      true,
	}
	artist.CreatedAt = time.Now()
	artist.UpdatedAt = artist.CreatedAt

	if err := s.artistRepo.CreateArtist(ctx, artist); err != nil {
		return nil, fmt.Errorf("failed to create artist: %w", err)
	}

	response := artist.ToResponse()
	return &response, nil
}

// GetArtist retrieves an artist with stats
func (s *artistService) GetArtist(ctx context.Context, userID, artistID string) (*models.ArtistWithStatsResponse, error) {
	artist, err := s.artistRepo.GetArtist(ctx, userID, artistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("artist", artistID)
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}

	// Get stats concurrently
	trackCount, err := s.artistRepo.GetArtistTrackCount(ctx, userID, artistID)
	if err != nil {
		trackCount = 0 // Non-critical error, just use 0
	}

	albumCount, err := s.artistRepo.GetArtistAlbumCount(ctx, userID, artistID)
	if err != nil {
		albumCount = 0
	}

	totalPlays, err := s.artistRepo.GetArtistTotalPlays(ctx, userID, artistID)
	if err != nil {
		totalPlays = 0
	}

	artistWithStats := models.ArtistWithStats{
		Artist:     *artist,
		TrackCount: trackCount,
		AlbumCount: albumCount,
		TotalPlays: totalPlays,
	}

	response := artistWithStats.ToResponseWithStats()
	return &response, nil
}

// UpdateArtist updates an existing artist
func (s *artistService) UpdateArtist(ctx context.Context, userID, artistID string, req models.UpdateArtistRequest) (*models.ArtistResponse, error) {
	artist, err := s.artistRepo.GetArtist(ctx, userID, artistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("artist", artistID)
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		artist.Name = *req.Name
		// Regenerate sort name if name changed and sort name not provided
		if req.SortName == nil {
			artist.SortName = models.GenerateSortName(*req.Name)
		}
	}
	if req.SortName != nil {
		artist.SortName = *req.SortName
	}
	if req.Bio != nil {
		artist.Bio = *req.Bio
	}
	if req.ImageURL != nil {
		artist.ImageURL = *req.ImageURL
	}
	if req.ExternalLinks != nil {
		artist.ExternalLinks = req.ExternalLinks
	}

	artist.UpdatedAt = time.Now()

	if err := s.artistRepo.UpdateArtist(ctx, *artist); err != nil {
		return nil, fmt.Errorf("failed to update artist: %w", err)
	}

	response := artist.ToResponse()
	return &response, nil
}

// DeleteArtist soft-deletes an artist
func (s *artistService) DeleteArtist(ctx context.Context, userID, artistID string) error {
	// Verify artist exists
	_, err := s.artistRepo.GetArtist(ctx, userID, artistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("artist", artistID)
		}
		return fmt.Errorf("failed to get artist: %w", err)
	}

	if err := s.artistRepo.DeleteArtist(ctx, userID, artistID); err != nil {
		return fmt.Errorf("failed to delete artist: %w", err)
	}

	return nil
}

// ListArtists returns paginated list of artists
func (s *artistService) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.ArtistResponse], error) {
	result, err := s.artistRepo.ListArtists(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list artists: %w", err)
	}

	responses := make([]models.ArtistResponse, 0, len(result.Items))
	for _, artist := range result.Items {
		responses = append(responses, artist.ToResponse())
	}

	return &repository.PaginatedResult[models.ArtistResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// SearchArtists searches for artists by name
func (s *artistService) SearchArtists(ctx context.Context, userID, query string, limit int) ([]models.ArtistResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	artists, err := s.artistRepo.SearchArtists(ctx, userID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search artists: %w", err)
	}

	responses := make([]models.ArtistResponse, 0, len(artists))
	for _, artist := range artists {
		responses = append(responses, artist.ToResponse())
	}

	return responses, nil
}

// GetArtistTracks retrieves all tracks by an artist
func (s *artistService) GetArtistTracks(ctx context.Context, userID, artistID string) ([]models.TrackResponse, error) {
	// Verify artist exists
	artist, err := s.artistRepo.GetArtist(ctx, userID, artistID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("artist", artistID)
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}

	// Get tracks by artist name (for now, until artistId field is populated on tracks)
	tracks, err := s.artistRepo.ListTracksByArtist(ctx, userID, artist.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracks: %w", err)
	}

	responses := make([]models.TrackResponse, 0, len(tracks))
	for _, track := range tracks {
		coverArtURL := ""
		if track.CoverArtKey != "" && s.s3Repo != nil {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}
		responses = append(responses, track.ToResponse(coverArtURL))
	}

	return responses, nil
}

// GenerateDeterministicArtistID generates a deterministic UUID for an artist
// based on the user ID and artist name. Used during migration to ensure
// idempotent artist creation.
func GenerateDeterministicArtistID(userID, artistName string) string {
	// Use SHA-1 to create a namespace-based UUID (v5-like)
	h := sha1.New()
	h.Write([]byte(userID))
	h.Write([]byte(artistName))
	hash := h.Sum(nil)

	// Format as UUID (set version and variant bits)
	hash[6] = (hash[6] & 0x0f) | 0x50 // Version 5
	hash[8] = (hash[8] & 0x3f) | 0x80 // Variant

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		hash[0:4],
		hash[4:6],
		hash[6:8],
		hash[8:10],
		hash[10:16])
}
