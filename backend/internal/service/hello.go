package service

import (
	"context"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

const (
	helloDefaultLimit = 20
	helloSeedUserID   = "seed-user"
)

// HelloTrack represents a simplified track for the hello-world feature
type HelloTrack struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	Year     int    `json:"year"`
	Duration int    `json:"duration"`
}

// HelloRepository defines the data access interface for hello tracks
type HelloRepository interface {
	GetTracksByUser(ctx context.Context, userID string) ([]HelloTrack, error)
}

// HelloService provides search and featured track operations
type HelloService struct {
	repo HelloRepository
}

// NewHelloService creates a new HelloService with the given repository
func NewHelloService(repo HelloRepository) *HelloService {
	return &HelloService{repo: repo}
}

// SearchTracks filters tracks matching query (case-insensitive) against title, artist, album, genre.
// Returns empty slice for empty query or no matches (never nil).
func (s *HelloService) SearchTracks(ctx context.Context, query string) ([]HelloTrack, error) {
	if query == "" {
		return []HelloTrack{}, nil
	}

	tracks, err := s.repo.GetTracksByUser(ctx, helloSeedUserID)
	if err != nil {
		return nil, err
	}

	lowerQuery := strings.ToLower(query)
	var results []HelloTrack
	for _, t := range tracks {
		if strings.Contains(strings.ToLower(t.Title), lowerQuery) ||
			strings.Contains(strings.ToLower(t.Artist), lowerQuery) ||
			strings.Contains(strings.ToLower(t.Album), lowerQuery) ||
			strings.Contains(strings.ToLower(t.Genre), lowerQuery) {
			results = append(results, t)
		}
	}

	if results == nil {
		results = []HelloTrack{}
	}

	return results, nil
}

// GetFeaturedTracks returns tracks up to limit.
// When limit == 0, uses the default limit of 20.
func (s *HelloService) GetFeaturedTracks(ctx context.Context, limit int) ([]HelloTrack, error) {
	if limit <= 0 {
		limit = helloDefaultLimit
	}

	tracks, err := s.repo.GetTracksByUser(ctx, helloSeedUserID)
	if err != nil {
		return nil, err
	}

	if len(tracks) > limit {
		tracks = tracks[:limit]
	}

	return tracks, nil
}

// HelloRepoAdapter bridges repository.Repository to HelloRepository
type HelloRepoAdapter struct {
	repo repository.Repository
}

// NewHelloRepoAdapter creates a new HelloRepoAdapter
func NewHelloRepoAdapter(repo repository.Repository) *HelloRepoAdapter {
	return &HelloRepoAdapter{repo: repo}
}

// GetTracksByUser retrieves tracks for a user and converts them to HelloTrack
func (a *HelloRepoAdapter) GetTracksByUser(ctx context.Context, userID string) ([]HelloTrack, error) {
	result, err := a.repo.ListTracks(ctx, userID, models.TrackFilter{Limit: 100})
	if err != nil {
		return nil, err
	}

	tracks := make([]HelloTrack, 0, len(result.Items))
	for _, t := range result.Items {
		tracks = append(tracks, HelloTrack{
			ID:       t.ID,
			Title:    t.Title,
			Artist:   t.Artist,
			Album:    t.Album,
			Genre:    t.Genre,
			Year:     t.Year,
			Duration: t.Duration,
		})
	}

	return tracks, nil
}
