package service

import (
	"context"
	"strings"
)

// HelloTrack represents a track in the hello-world seed data.
type HelloTrack struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	Year     int    `json:"year"`
	Duration int    `json:"duration"`
}

// HelloRepository defines the data access interface for hello seed tracks.
type HelloRepository interface {
	GetSeedTracks(ctx context.Context) ([]HelloTrack, error)
}

// HelloServiceImpl implements the HelloService interface.
type HelloServiceImpl struct {
	repo HelloRepository
}

// NewHelloService creates a new HelloServiceImpl.
func NewHelloService(repo HelloRepository) *HelloServiceImpl {
	return &HelloServiceImpl{repo: repo}
}

// Search filters seed tracks by a case-insensitive substring match across
// Title, Artist, Album, and Genre fields. An empty query returns an empty slice
// without calling the repository.
func (s *HelloServiceImpl) Search(ctx context.Context, query string) ([]HelloTrack, error) {
	if query == "" {
		return []HelloTrack{}, nil
	}

	tracks, err := s.repo.GetSeedTracks(ctx)
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	results := []HelloTrack{}
	for _, t := range tracks {
		if strings.Contains(strings.ToLower(t.Title), q) ||
			strings.Contains(strings.ToLower(t.Artist), q) ||
			strings.Contains(strings.ToLower(t.Album), q) ||
			strings.Contains(strings.ToLower(t.Genre), q) {
			results = append(results, t)
		}
	}

	return results, nil
}

// Featured returns the first N seed tracks. If limit <= 0, all tracks are returned.
// If limit exceeds the number of available tracks, all tracks are returned.
func (s *HelloServiceImpl) Featured(ctx context.Context, limit int) ([]HelloTrack, error) {
	tracks, err := s.repo.GetSeedTracks(ctx)
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > len(tracks) {
		return tracks, nil
	}

	return tracks[:limit], nil
}
