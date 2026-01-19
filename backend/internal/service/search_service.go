package service

import (
	"context"
	"os"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// SearchService handles search-related operations
type SearchService struct {
	repo          repository.Repository
	nixieEndpoint string
}

// NewSearchService creates a new SearchService
func NewSearchService(repo repository.Repository) *SearchService {
	nixieEndpoint := os.Getenv("NIXIESEARCH_ENDPOINT")
	if nixieEndpoint == "" {
		nixieEndpoint = "http://localhost:8080"
	}

	return &SearchService{
		repo:          repo,
		nixieEndpoint: nixieEndpoint,
	}
}

// Search performs a search across tracks, albums, and artists
func (s *SearchService) Search(ctx context.Context, userID string, req models.SearchRequest) (*models.SearchResponse, error) {
	// For now, implement simple DynamoDB-based search
	// In production, this would use Nixiesearch
	query := strings.ToLower(req.Query)

	// Search tracks
	trackFilter := models.TrackFilter{
		Limit: req.Limit,
	}
	tracks, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err != nil {
		return nil, err
	}

	// Filter tracks that match the query
	matchedTracks := make([]models.TrackResponse, 0)
	for _, track := range tracks.Items {
		if matchesQuery(track, query) {
			coverArtURL := s.getCoverArtURL(track.CoverArtKey)
			matchedTracks = append(matchedTracks, track.ToResponse(coverArtURL))
		}
	}

	// Search albums
	albumFilter := models.AlbumFilter{
		Limit: req.Limit,
	}
	albums, err := s.repo.ListAlbums(ctx, userID, albumFilter)
	if err != nil {
		return nil, err
	}

	matchedAlbums := make([]models.AlbumResponse, 0)
	for _, album := range albums.Items {
		if matchesAlbumQuery(album, query) {
			coverArtURL := s.getCoverArtURL(album.CoverArtKey)
			matchedAlbums = append(matchedAlbums, album.ToResponse(coverArtURL))
		}
	}

	// Build artist summaries from matched tracks
	artistMap := make(map[string]*models.ArtistSummary)
	for _, track := range matchedTracks {
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

	return &models.SearchResponse{
		Query:        req.Query,
		TotalResults: len(matchedTracks) + len(matchedAlbums),
		Tracks:       matchedTracks,
		Albums:       matchedAlbums,
		Artists:      artists,
		Limit:        req.Limit,
		Offset:       req.Offset,
	}, nil
}

// GetSuggestions returns autocomplete suggestions
func (s *SearchService) GetSuggestions(ctx context.Context, userID, query string) (*models.AutocompleteResponse, error) {
	query = strings.ToLower(query)
	suggestions := make([]models.SearchSuggestion, 0)

	// Get track suggestions
	trackFilter := models.TrackFilter{Limit: 100}
	tracks, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err == nil {
		for _, track := range tracks.Items {
			if strings.Contains(strings.ToLower(track.Title), query) {
				suggestions = append(suggestions, models.SearchSuggestion{
					Text:     track.Title,
					Type:     "track",
					ID:       track.ID,
					ImageURL: s.getCoverArtURL(track.CoverArtKey),
				})
			}
			if strings.Contains(strings.ToLower(track.Artist), query) {
				// Add artist if not already in suggestions
				found := false
				for _, s := range suggestions {
					if s.Type == "artist" && s.Text == track.Artist {
						found = true
						break
					}
				}
				if !found {
					suggestions = append(suggestions, models.SearchSuggestion{
						Text: track.Artist,
						Type: "artist",
					})
				}
			}
		}
	}

	// Limit suggestions
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return &models.AutocompleteResponse{
		Query:       query,
		Suggestions: suggestions,
	}, nil
}

func matchesQuery(track models.Track, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(track.Title), query) ||
		strings.Contains(strings.ToLower(track.Artist), query) ||
		strings.Contains(strings.ToLower(track.Album), query) ||
		strings.Contains(strings.ToLower(track.Genre), query)
}

func matchesAlbumQuery(album models.Album, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(album.Title), query) ||
		strings.Contains(strings.ToLower(album.Artist), query)
}

func (s *SearchService) getCoverArtURL(coverArtKey string) string {
	if coverArtKey == "" {
		return ""
	}
	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
	if cloudFrontDomain == "" {
		return ""
	}
	return "https://" + cloudFrontDomain + "/" + coverArtKey
}
