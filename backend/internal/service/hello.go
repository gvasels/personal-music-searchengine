package service

import (
	"context"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

const seedUserID = "seed-user"

// HelloRepository defines the data access methods needed by HelloService
type HelloRepository interface {
	ListTracksByUser(ctx context.Context, userID string) ([]models.Track, error)
}

// HelloService provides hello-world search functionality
type HelloService struct {
	repo HelloRepository
}

// NewHelloService creates a new HelloService
func NewHelloService(repo HelloRepository) *HelloService {
	return &HelloService{repo: repo}
}

// SearchTracks searches seed tracks by title, artist, or genre (case-insensitive)
func (s *HelloService) SearchTracks(ctx context.Context, query string, limit int) ([]models.TrackResponse, error) {
	if query == "" {
		return []models.TrackResponse{}, nil
	}

	tracks, err := s.repo.ListTracksByUser(ctx, seedUserID)
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	var results []models.TrackResponse

	for _, t := range tracks {
		if strings.Contains(strings.ToLower(t.Title), queryLower) ||
			strings.Contains(strings.ToLower(t.Artist), queryLower) ||
			strings.Contains(strings.ToLower(t.Genre), queryLower) {
			results = append(results, t.ToResponse(""))
		}
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	if results == nil {
		results = []models.TrackResponse{}
	}

	return results, nil
}

// GetFeaturedTracks returns all seed tracks
func (s *HelloService) GetFeaturedTracks(ctx context.Context, limit int) ([]models.TrackResponse, error) {
	tracks, err := s.repo.ListTracksByUser(ctx, seedUserID)
	if err != nil {
		return nil, err
	}

	var results []models.TrackResponse
	for i, t := range tracks {
		if limit > 0 && i >= limit {
			break
		}
		results = append(results, t.ToResponse(""))
	}

	if results == nil {
		results = []models.TrackResponse{}
	}

	return results, nil
}

// HelloRepoAdapter adapts the existing Repository interface for HelloService
type HelloRepoAdapter struct {
	repo repository.Repository
}

// NewHelloRepoAdapter creates a new HelloRepoAdapter
func NewHelloRepoAdapter(repo repository.Repository) *HelloRepoAdapter {
	return &HelloRepoAdapter{repo: repo}
}

// ListTracksByUser lists all tracks for a user using the existing repository
func (a *HelloRepoAdapter) ListTracksByUser(ctx context.Context, userID string) ([]models.Track, error) {
	filter := models.TrackFilter{Limit: 100}
	result, err := a.repo.ListTracks(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
