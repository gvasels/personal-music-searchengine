package service

import (
	"context"
	"os"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// AlbumService handles album-related operations
type AlbumService struct {
	repo        repository.Repository
	mediaBucket string
}

// NewAlbumService creates a new AlbumService
func NewAlbumService(repo repository.Repository) *AlbumService {
	mediaBucket := os.Getenv("MEDIA_BUCKET")
	if mediaBucket == "" {
		mediaBucket = "music-library-media"
	}

	return &AlbumService{
		repo:        repo,
		mediaBucket: mediaBucket,
	}
}

// ListAlbums lists albums with filtering
func (s *AlbumService) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResponse[models.AlbumResponse], error) {
	result, err := s.repo.ListAlbums(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.AlbumResponse, len(result.Items))
	for i, album := range result.Items {
		coverArtURL := s.getCoverArtURL(album.CoverArtKey)
		responses[i] = album.ToResponse(coverArtURL)
	}

	return &models.PaginatedResponse[models.AlbumResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}, nil
}

// GetAlbumWithTracks retrieves an album with its tracks
func (s *AlbumService) GetAlbumWithTracks(ctx context.Context, userID, albumID string) (*models.AlbumWithTracks, error) {
	album, err := s.repo.GetAlbum(ctx, userID, albumID)
	if err != nil {
		return nil, err
	}

	tracks, err := s.repo.ListTracksByAlbum(ctx, userID, albumID)
	if err != nil {
		return nil, err
	}

	coverArtURL := s.getCoverArtURL(album.CoverArtKey)
	albumResponse := album.ToResponse(coverArtURL)

	trackResponses := make([]models.TrackResponse, len(tracks))
	for i, track := range tracks {
		trackCoverURL := s.getCoverArtURL(track.CoverArtKey)
		trackResponses[i] = track.ToResponse(trackCoverURL)
	}

	return &models.AlbumWithTracks{
		Album:  albumResponse,
		Tracks: trackResponses,
	}, nil
}

// ListArtists lists artists with their stats
func (s *AlbumService) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*models.PaginatedResponse[models.ArtistSummary], error) {
	// Get all albums to build artist list
	albumFilter := models.AlbumFilter{
		Limit:     1000,
		SortBy:    "artist",
		SortOrder: "asc",
	}
	albums, err := s.repo.ListAlbums(ctx, userID, albumFilter)
	if err != nil {
		return nil, err
	}

	// Get all tracks to count per artist
	trackFilter := models.TrackFilter{
		Limit: 10000,
	}
	tracks, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err != nil {
		return nil, err
	}

	// Build artist summaries
	artistMap := make(map[string]*models.ArtistSummary)

	for _, album := range albums.Items {
		if _, exists := artistMap[album.Artist]; !exists {
			artistMap[album.Artist] = &models.ArtistSummary{
				Name:       album.Artist,
				CoverArtURL: s.getCoverArtURL(album.CoverArtKey),
			}
		}
		artistMap[album.Artist].AlbumCount++
	}

	for _, track := range tracks.Items {
		if _, exists := artistMap[track.Artist]; !exists {
			artistMap[track.Artist] = &models.ArtistSummary{
				Name: track.Artist,
			}
		}
		artistMap[track.Artist].TrackCount++
	}

	artists := make([]models.ArtistSummary, 0, len(artistMap))
	for _, artist := range artistMap {
		artists = append(artists, *artist)
	}

	// Apply pagination
	start := 0
	end := len(artists)
	if filter.Limit > 0 && filter.Limit < end {
		end = filter.Limit
	}

	return &models.PaginatedResponse[models.ArtistSummary]{
		Items: artists[start:end],
		Pagination: models.Pagination{
			Limit: filter.Limit,
		},
	}, nil
}

// GetArtist retrieves an artist with their albums and tracks
func (s *AlbumService) GetArtist(ctx context.Context, userID, artistName string) (map[string]interface{}, error) {
	albums, err := s.repo.ListAlbumsByArtist(ctx, userID, artistName)
	if err != nil {
		return nil, err
	}

	// Get all tracks by artist
	trackFilter := models.TrackFilter{
		Artist: artistName,
		Limit:  1000,
	}
	tracks, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err != nil {
		return nil, err
	}

	albumResponses := make([]models.AlbumResponse, len(albums))
	for i, album := range albums {
		coverArtURL := s.getCoverArtURL(album.CoverArtKey)
		albumResponses[i] = album.ToResponse(coverArtURL)
	}

	trackResponses := make([]models.TrackResponse, len(tracks.Items))
	for i, track := range tracks.Items {
		coverURL := s.getCoverArtURL(track.CoverArtKey)
		trackResponses[i] = track.ToResponse(coverURL)
	}

	return map[string]interface{}{
		"name":       artistName,
		"albums":     albumResponses,
		"tracks":     trackResponses,
		"albumCount": len(albums),
		"trackCount": len(tracks.Items),
	}, nil
}

func (s *AlbumService) getCoverArtURL(coverArtKey string) string {
	if coverArtKey == "" {
		return ""
	}
	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
	if cloudFrontDomain == "" {
		return ""
	}
	return "https://" + cloudFrontDomain + "/" + coverArtKey
}
