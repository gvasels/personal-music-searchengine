package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/search"
)

// searchServiceImpl implements SearchService using Nixiesearch.
type searchServiceImpl struct {
	client *search.Client
	repo   repository.Repository
	s3Repo repository.S3Repository
}

// NewSearchService creates a new search service.
func NewSearchService(client *search.Client, repo repository.Repository, s3Repo repository.S3Repository) SearchService {
	return &searchServiceImpl{
		client: client,
		repo:   repo,
		s3Repo: s3Repo,
	}
}

// Search executes a search query scoped to the user.
func (s *searchServiceImpl) Search(ctx context.Context, userID string, req models.SearchRequest) (*models.SearchResponse, error) {
	if req.Query == "" {
		return nil, models.NewValidationError("search query cannot be empty")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Convert models.SearchRequest to search.SearchQuery
	searchQuery := search.SearchQuery{
		Query:  req.Query,
		Limit:  limit,
		Cursor: req.Cursor,
	}

	// Convert filters
	searchQuery.Filters = s.convertFilters(req.Filters)

	// Convert sort
	if req.Sort.Field != "" {
		searchQuery.Sort = &search.SortOption{
			Field: req.Sort.Field,
			Order: req.Sort.Order,
		}
	}

	// Execute search
	resp, err := s.client.Search(ctx, userID, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert to API response
	tracks := make([]models.TrackResponse, 0, len(resp.Results))
	for _, result := range resp.Results {
		track := s.searchResultToTrackResponse(result)
		tracks = append(tracks, track)
	}

	// Enrich with cover art URLs
	s.enrichTracksWithCoverArt(ctx, userID, tracks)

	return &models.SearchResponse{
		Query:        req.Query,
		TotalResults: resp.Total,
		Tracks:       tracks,
		Limit:        limit,
		NextCursor:   resp.NextCursor,
		HasMore:      resp.NextCursor != "",
	}, nil
}

// Autocomplete provides search suggestions.
func (s *searchServiceImpl) Autocomplete(ctx context.Context, userID, query string) (*models.AutocompleteResponse, error) {
	if query == "" {
		return &models.AutocompleteResponse{
			Query:       query,
			Suggestions: []models.SearchSuggestion{},
		}, nil
	}

	// Execute a limited search for autocomplete
	searchQuery := search.SearchQuery{
		Query: query,
		Limit: 10,
	}

	resp, err := s.client.Search(ctx, userID, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("autocomplete failed: %w", err)
	}

	// Convert results to suggestions
	suggestions := make([]models.SearchSuggestion, 0)
	seenArtists := make(map[string]bool)
	seenAlbums := make(map[string]bool)

	for _, result := range resp.Results {
		// Add track suggestion
		suggestions = append(suggestions, models.SearchSuggestion{
			Text: result.Title,
			Type: "track",
			ID:   result.ID,
		})

		// Add artist suggestion (deduplicated)
		if result.Artist != "" && !seenArtists[result.Artist] {
			seenArtists[result.Artist] = true
			suggestions = append(suggestions, models.SearchSuggestion{
				Text: result.Artist,
				Type: "artist",
			})
		}

		// Add album suggestion (deduplicated)
		if result.Album != "" && !seenAlbums[result.Album] {
			seenAlbums[result.Album] = true
			suggestions = append(suggestions, models.SearchSuggestion{
				Text: result.Album,
				Type: "album",
			})
		}

		// Limit total suggestions
		if len(suggestions) >= 10 {
			break
		}
	}

	return &models.AutocompleteResponse{
		Query:       query,
		Suggestions: suggestions,
	}, nil
}

// IndexTrack indexes a track in the search engine.
func (s *searchServiceImpl) IndexTrack(ctx context.Context, track models.Track) error {
	doc := search.Document{
		ID:        track.ID,
		UserID:    track.UserID,
		Title:     track.Title,
		Artist:    track.Artist,
		Album:     track.Album,
		Genre:     track.Genre,
		Year:      track.Year,
		Duration:  track.Duration,
		Filename:  track.S3Key,
		IndexedAt: time.Now(),
	}

	resp, err := s.client.Index(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to index track %s: %w", track.ID, err)
	}

	if !resp.Indexed {
		return fmt.Errorf("track %s was not indexed", track.ID)
	}

	return nil
}

// RemoveTrack removes a track from the search index.
func (s *searchServiceImpl) RemoveTrack(ctx context.Context, trackID string) error {
	resp, err := s.client.Delete(ctx, trackID)
	if err != nil {
		return fmt.Errorf("failed to remove track %s from index: %w", trackID, err)
	}

	if !resp.Deleted {
		// Track might not exist in index - log but don't fail
		fmt.Printf("Warning: track %s was not found in search index\n", trackID)
	}

	return nil
}

// RebuildIndex rebuilds the entire search index for a user.
func (s *searchServiceImpl) RebuildIndex(ctx context.Context, userID string) error {
	// Collect all tracks for the user using pagination
	var allTracks []models.Track
	cursor := ""

	for {
		filter := models.TrackFilter{
			Limit:   100,
			LastKey: cursor,
		}

		result, err := s.repo.ListTracks(ctx, userID, filter)
		if err != nil {
			return fmt.Errorf("failed to list tracks for rebuild: %w", err)
		}

		allTracks = append(allTracks, result.Items...)

		if !result.HasMore || result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}

	if len(allTracks) == 0 {
		return nil
	}

	// Convert tracks to documents
	docs := make([]search.Document, len(allTracks))
	for i, track := range allTracks {
		docs[i] = search.Document{
			ID:        track.ID,
			UserID:    track.UserID,
			Title:     track.Title,
			Artist:    track.Artist,
			Album:     track.Album,
			Genre:     track.Genre,
			Year:      track.Year,
			Duration:  track.Duration,
			Filename:  track.S3Key,
			IndexedAt: time.Now(),
		}
	}

	// Bulk index in batches of 100
	batchSize := 100
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}

		batch := docs[i:end]
		resp, err := s.client.BulkIndex(ctx, batch)
		if err != nil {
			return fmt.Errorf("bulk index failed at batch %d: %w", i/batchSize, err)
		}

		if resp.Failed > 0 {
			fmt.Printf("Warning: %d documents failed to index in batch %d\n", resp.Failed, i/batchSize)
		}
	}

	return nil
}

// convertFilters converts models.SearchFilters to search.SearchFilters.
func (s *searchServiceImpl) convertFilters(filters models.SearchFilters) search.SearchFilters {
	result := search.SearchFilters{}

	// Use first artist if provided
	if len(filters.Artists) > 0 {
		result.Artist = filters.Artists[0]
	}

	// Use first album if provided
	if len(filters.Albums) > 0 {
		result.Album = filters.Albums[0]
	}

	// Use first genre if provided
	if len(filters.Genres) > 0 {
		result.Genre = filters.Genres[0]
	}

	// Convert year range
	if len(filters.Years) > 0 {
		minYear, maxYear := filters.Years[0], filters.Years[0]
		for _, year := range filters.Years {
			if year < minYear {
				minYear = year
			}
			if year > maxYear {
				maxYear = year
			}
		}
		result.YearFrom = minYear
		result.YearTo = maxYear
	}

	return result
}

// searchResultToTrackResponse converts a search result to a track response.
func (s *searchServiceImpl) searchResultToTrackResponse(result search.SearchResult) models.TrackResponse {
	return models.TrackResponse{
		ID:          result.ID,
		Title:       result.Title,
		Artist:      result.Artist,
		Album:       result.Album,
		Genre:       result.Genre,
		Year:        result.Year,
		Duration:    result.Duration,
		DurationStr: formatDuration(result.Duration),
	}
}

// enrichTracksWithCoverArt adds cover art URLs to track responses.
func (s *searchServiceImpl) enrichTracksWithCoverArt(ctx context.Context, userID string, tracks []models.TrackResponse) {
	for i := range tracks {
		track, err := s.repo.GetTrack(ctx, userID, tracks[i].ID)
		if err != nil {
			continue
		}
		if track.CoverArtKey != "" && s.s3Repo != nil {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
			if err == nil {
				tracks[i].CoverArtURL = url
			}
		}
	}
}

// formatDuration formats seconds as "M:SS".
func formatDuration(seconds int) string {
	if seconds <= 0 {
		return "0:00"
	}
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}
