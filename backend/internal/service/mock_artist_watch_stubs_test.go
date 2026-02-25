package service

// This file adds ArtistWatch repository method stubs to all mock repository types
// that implement repository.Repository. These stubs are needed because the
// ArtistWatch methods were added to the Repository interface and existing mocks
// in other test files don't have them.

import (
	"context"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// MockPlaylistRepository stubs

func (m *MockPlaylistRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockPlaylistRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockPlaylistRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockPlaylistRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}

// MockRepository (search) stubs

func (m *MockRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}

// MockFilterTagsRepository stubs

func (m *MockFilterTagsRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockFilterTagsRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockFilterTagsRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockFilterTagsRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}

// MockSimilarityRepository stubs

func (m *MockSimilarityRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockSimilarityRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockSimilarityRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockSimilarityRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}

// MockStatsRepository stubs

func (m *MockStatsRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockStatsRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockStatsRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockStatsRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}

// MockTagRepository stubs

func (m *MockTagRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockTagRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockTagRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockTagRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}

// MockTrackServiceRepository stubs

func (m *MockTrackServiceRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	return nil
}

func (m *MockTrackServiceRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	return nil
}

func (m *MockTrackServiceRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	return nil, repository.ErrNotFound
}

func (m *MockTrackServiceRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	return &repository.PaginatedResult[models.ArtistWatch]{}, nil
}
