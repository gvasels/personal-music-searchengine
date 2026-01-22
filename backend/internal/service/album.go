package service

import (
	"context"
	"sort"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// albumService implements AlbumService
type albumService struct {
	repo   repository.Repository
	s3Repo repository.S3Repository
}

// NewAlbumService creates a new album service
func NewAlbumService(repo repository.Repository, s3Repo repository.S3Repository) AlbumService {
	return &albumService{
		repo:   repo,
		s3Repo: s3Repo,
	}
}

func (s *albumService) GetAlbum(ctx context.Context, userID, albumID string) (*models.AlbumWithTracks, error) {
	album, err := s.repo.GetAlbum(ctx, userID, albumID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("Album", albumID)
		}
		return nil, err
	}

	coverArtURL := ""
	if album.CoverArtKey != "" {
		url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, album.CoverArtKey, 24*time.Hour)
		if err == nil {
			coverArtURL = url
		}
	}

	// Get tracks for this album
	filter := models.TrackFilter{
		Album: album.Title,
		Limit: 100,
	}
	trackResult, err := s.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Filter to only tracks matching this album and convert to responses
	var tracks []models.TrackResponse
	for _, track := range trackResult.Items {
		if track.Album == album.Title && track.Artist == album.Artist {
			trackCoverURL := ""
			if track.CoverArtKey != "" {
				url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, track.CoverArtKey, 24*time.Hour)
				if err == nil {
					trackCoverURL = url
				}
			}
			tracks = append(tracks, track.ToResponse(trackCoverURL))
		}
	}

	// Sort tracks by disc number, then track number
	sort.Slice(tracks, func(i, j int) bool {
		if tracks[i].DiscNumber != tracks[j].DiscNumber {
			return tracks[i].DiscNumber < tracks[j].DiscNumber
		}
		return tracks[i].TrackNumber < tracks[j].TrackNumber
	})

	// Calculate total duration
	totalDuration := 0
	for _, t := range tracks {
		totalDuration += t.Duration
	}

	albumResp := album.ToResponse(coverArtURL)
	// Override with calculated values
	albumResp.TrackCount = len(tracks)
	albumResp.TotalDuration = totalDuration

	return &models.AlbumWithTracks{
		Album:  albumResp,
		Tracks: tracks,
	}, nil
}

func (s *albumService) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.AlbumResponse], error) {
	result, err := s.repo.ListAlbums(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Get all tracks to calculate actual track counts per album
	trackFilter := models.TrackFilter{
		Limit: 1000, // Get all tracks for aggregation
	}
	trackResult, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err != nil {
		return nil, err
	}

	// Build album -> track count and duration map
	// Key format: "albumTitle|artist" for uniqueness
	type albumStats struct {
		trackCount    int
		totalDuration int
	}
	albumStatsMap := make(map[string]*albumStats)
	for _, track := range trackResult.Items {
		if track.Album == "" {
			continue
		}
		key := track.Album + "|" + track.Artist
		if _, exists := albumStatsMap[key]; !exists {
			albumStatsMap[key] = &albumStats{}
		}
		albumStatsMap[key].trackCount++
		albumStatsMap[key].totalDuration += track.Duration
	}

	responses := make([]models.AlbumResponse, 0, len(result.Items))
	for _, album := range result.Items {
		// Look up actual track count
		key := album.Title + "|" + album.Artist
		stats := albumStatsMap[key]

		// Skip albums with 0 tracks (orphaned albums)
		if stats == nil || stats.trackCount == 0 {
			continue
		}

		coverArtURL := ""
		if album.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, album.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}

		resp := album.ToResponse(coverArtURL)
		// Override with calculated values
		resp.TrackCount = stats.trackCount
		resp.TotalDuration = stats.totalDuration
		responses = append(responses, resp)
	}

	return &repository.PaginatedResult[models.AlbumResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (s *albumService) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.AlbumResponse, error) {
	albums, err := s.repo.ListAlbumsByArtist(ctx, userID, artist)
	if err != nil {
		return nil, err
	}

	// Get tracks by this artist to calculate actual track counts
	trackFilter := models.TrackFilter{
		Artist: artist,
		Limit:  1000,
	}
	trackResult, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err != nil {
		return nil, err
	}

	// Build album -> track count and duration map
	type albumStats struct {
		trackCount    int
		totalDuration int
	}
	albumStatsMap := make(map[string]*albumStats)
	for _, track := range trackResult.Items {
		if track.Album == "" || track.Artist != artist {
			continue
		}
		if _, exists := albumStatsMap[track.Album]; !exists {
			albumStatsMap[track.Album] = &albumStats{}
		}
		albumStatsMap[track.Album].trackCount++
		albumStatsMap[track.Album].totalDuration += track.Duration
	}

	responses := make([]models.AlbumResponse, 0, len(albums))
	for _, album := range albums {
		// Look up actual track count
		stats := albumStatsMap[album.Title]

		// Skip albums with 0 tracks
		if stats == nil || stats.trackCount == 0 {
			continue
		}

		coverArtURL := ""
		if album.CoverArtKey != "" {
			url, err := s.s3Repo.GeneratePresignedDownloadURL(ctx, album.CoverArtKey, 24*time.Hour)
			if err == nil {
				coverArtURL = url
			}
		}

		resp := album.ToResponse(coverArtURL)
		resp.TrackCount = stats.trackCount
		resp.TotalDuration = stats.totalDuration
		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *albumService) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) ([]models.ArtistSummary, error) {
	// Get all tracks to aggregate artists
	trackFilter := models.TrackFilter{
		Limit: 1000, // Get all tracks for aggregation
	}
	result, err := s.repo.ListTracks(ctx, userID, trackFilter)
	if err != nil {
		return nil, err
	}

	// Aggregate by artist
	artistMap := make(map[string]*models.ArtistSummary)
	albumMap := make(map[string]map[string]bool) // artist -> set of albums

	for _, track := range result.Items {
		artist := track.Artist
		if artist == "" {
			continue
		}

		if _, exists := artistMap[artist]; !exists {
			artistMap[artist] = &models.ArtistSummary{
				Name:       artist,
				TrackCount: 0,
				AlbumCount: 0,
			}
			albumMap[artist] = make(map[string]bool)
		}

		artistMap[artist].TrackCount++
		if track.Album != "" {
			albumMap[artist][track.Album] = true
		}
	}

	// Count unique albums per artist
	for artist, albums := range albumMap {
		if summary, exists := artistMap[artist]; exists {
			summary.AlbumCount = len(albums)
		}
	}

	// Convert map to slice
	artists := make([]models.ArtistSummary, 0, len(artistMap))
	for _, summary := range artistMap {
		artists = append(artists, *summary)
	}

	// Sort based on filter
	switch filter.SortBy {
	case "trackCount":
		sort.Slice(artists, func(i, j int) bool {
			if filter.SortOrder == "desc" {
				return artists[i].TrackCount > artists[j].TrackCount
			}
			return artists[i].TrackCount < artists[j].TrackCount
		})
	case "albumCount":
		sort.Slice(artists, func(i, j int) bool {
			if filter.SortOrder == "desc" {
				return artists[i].AlbumCount > artists[j].AlbumCount
			}
			return artists[i].AlbumCount < artists[j].AlbumCount
		})
	default: // name
		sort.Slice(artists, func(i, j int) bool {
			if filter.SortOrder == "desc" {
				return artists[i].Name > artists[j].Name
			}
			return artists[i].Name < artists[j].Name
		})
	}

	// Apply limit
	if filter.Limit > 0 && filter.Limit < len(artists) {
		artists = artists[:filter.Limit]
	}

	return artists, nil
}
