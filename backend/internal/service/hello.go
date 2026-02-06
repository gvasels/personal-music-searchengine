package service

import (
	"context"
	"strings"
)

// HelloTrack represents a track from the hello world seed data
type HelloTrack struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	Year     int    `json:"year"`
	Duration int    `json:"duration"`
}

// HelloRepository defines the interface for accessing hello seed data
type HelloRepository interface {
	GetSeedTracks(ctx context.Context) ([]HelloTrack, error)
}

// HelloService provides search and featured track operations
type HelloService struct {
	repo HelloRepository
}

// NewHelloService creates a new HelloService with the given repository
func NewHelloService(repo HelloRepository) *HelloService {
	return &HelloService{
		repo: repo,
	}
}

// Search searches tracks by query, matching against title, artist, album, and genre.
// Empty query returns an empty slice.
// Search is case-insensitive.
func (s *HelloService) Search(ctx context.Context, query string) ([]HelloTrack, error) {
	// Empty query returns empty slice immediately
	if query == "" {
		return []HelloTrack{}, nil
	}

	tracks, err := s.repo.GetSeedTracks(ctx)
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	var results []HelloTrack

	for _, track := range tracks {
		if strings.Contains(strings.ToLower(track.Title), queryLower) ||
			strings.Contains(strings.ToLower(track.Artist), queryLower) ||
			strings.Contains(strings.ToLower(track.Album), queryLower) ||
			strings.Contains(strings.ToLower(track.Genre), queryLower) {
			results = append(results, track)
		}
	}

	// Return empty slice instead of nil for consistency
	if results == nil {
		return []HelloTrack{}, nil
	}

	return results, nil
}

// Featured returns featured tracks up to the specified limit.
// If limit <= 0, returns all tracks.
// If limit exceeds the number of tracks, returns all tracks.
func (s *HelloService) Featured(ctx context.Context, limit int) ([]HelloTrack, error) {
	tracks, err := s.repo.GetSeedTracks(ctx)
	if err != nil {
		return nil, err
	}

	// Handle empty tracks case
	if len(tracks) == 0 {
		return []HelloTrack{}, nil
	}

	// limit <= 0 means return all tracks
	if limit <= 0 {
		return tracks, nil
	}

	// Return up to limit tracks
	if limit >= len(tracks) {
		return tracks, nil
	}

	return tracks[:limit], nil
}
