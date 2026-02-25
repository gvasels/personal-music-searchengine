package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// ArtistWatchRepository defines the repository methods needed by ArtistWatchService.
type ArtistWatchRepository interface {
	CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error
	DeleteArtistWatch(ctx context.Context, userID, artistName string) error
	GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error)
	ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error)
}

// ArtistWatchService handles artist watch operations.
type ArtistWatchService struct {
	repo ArtistWatchRepository
}

// NewArtistWatchService creates a new ArtistWatchService.
func NewArtistWatchService(repo ArtistWatchRepository) *ArtistWatchService {
	return &ArtistWatchService{repo: repo}
}

// WatchArtist creates a watch for the given artist.
func (s *ArtistWatchService) WatchArtist(ctx context.Context, userID, artistName string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if strings.TrimSpace(artistName) == "" {
		return fmt.Errorf("artist name is required")
	}

	watch := models.NewArtistWatch(userID, artistName)
	err := s.repo.CreateArtistWatch(ctx, *watch)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return fmt.Errorf("already watching this artist")
		}
		return err
	}
	return nil
}

// UnwatchArtist removes a watch for the given artist.
func (s *ArtistWatchService) UnwatchArtist(ctx context.Context, userID, artistName string) error {
	err := s.repo.DeleteArtistWatch(ctx, userID, artistName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("not watching this artist")
		}
		return err
	}
	return nil
}

// IsWatching checks if the user is watching the given artist.
func (s *ArtistWatchService) IsWatching(ctx context.Context, userID, artistName string) (bool, error) {
	_, err := s.repo.GetArtistWatch(ctx, userID, artistName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListWatchedArtists returns a paginated list of the user's watched artists.
func (s *ArtistWatchService) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return s.repo.ListWatchedArtists(ctx, userID, limit, cursor)
}
