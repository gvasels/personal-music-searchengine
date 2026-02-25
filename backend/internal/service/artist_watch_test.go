package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// ArtistWatch Service Tests
// =============================================================================

// MockArtistWatchRepository provides mockable repository methods for artist watch tests.
type MockArtistWatchRepository struct {
	mock.Mock
}

func (m *MockArtistWatchRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	args := m.Called(ctx, watch)
	return args.Error(0)
}

func (m *MockArtistWatchRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	args := m.Called(ctx, userID, artistName)
	return args.Error(0)
}

func (m *MockArtistWatchRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	args := m.Called(ctx, userID, artistName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistWatch), args.Error(1)
}

func (m *MockArtistWatchRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	args := m.Called(ctx, userID, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.ArtistWatch]), args.Error(1)
}

// TestArtistWatchService_WatchArtist verifies that WatchArtist delegates to the
// repository's CreateArtistWatch and returns nil on success.
func TestArtistWatchService_WatchArtist(t *testing.T) {
	t.Run("creates artist watch via repository", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("CreateArtistWatch", ctx, mock.MatchedBy(func(w models.ArtistWatch) bool {
			return w.UserID == "user-123" && w.ArtistName == "Daft Punk"
		})).Return(nil)

		svc := NewArtistWatchService(mockRepo)
		err := svc.WatchArtist(ctx, "user-123", "Daft Punk")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("CreateArtistWatch", ctx, mock.Anything).Return(errors.New("dynamo error"))

		svc := NewArtistWatchService(mockRepo)
		err := svc.WatchArtist(ctx, "user-123", "Daft Punk")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "dynamo error")
	})

	t.Run("returns already exists error when already watching", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("CreateArtistWatch", ctx, mock.Anything).Return(repository.ErrAlreadyExists)

		svc := NewArtistWatchService(mockRepo)
		err := svc.WatchArtist(ctx, "user-123", "Daft Punk")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "already watching")
	})
}

// TestArtistWatchService_WatchArtist_EmptyName verifies that WatchArtist returns
// a validation error when the artist name is empty or whitespace-only.
func TestArtistWatchService_WatchArtist_EmptyName(t *testing.T) {
	t.Run("returns error for empty artist name", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		svc := NewArtistWatchService(mockRepo)
		err := svc.WatchArtist(ctx, "user-123", "")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "artist name")
	})

	t.Run("returns error for whitespace-only artist name", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		svc := NewArtistWatchService(mockRepo)
		err := svc.WatchArtist(ctx, "user-123", "   ")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "artist name")
	})

	t.Run("returns error for empty user ID", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		svc := NewArtistWatchService(mockRepo)
		err := svc.WatchArtist(ctx, "", "Daft Punk")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user ID")
	})
}

// TestArtistWatchService_UnwatchArtist verifies that UnwatchArtist delegates to
// the repository's DeleteArtistWatch.
func TestArtistWatchService_UnwatchArtist(t *testing.T) {
	t.Run("deletes artist watch via repository", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("DeleteArtistWatch", ctx, "user-123", "Daft Punk").Return(nil)

		svc := NewArtistWatchService(mockRepo)
		err := svc.UnwatchArtist(ctx, "user-123", "Daft Punk")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when not watching", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("DeleteArtistWatch", ctx, "user-123", "Daft Punk").Return(repository.ErrNotFound)

		svc := NewArtistWatchService(mockRepo)
		err := svc.UnwatchArtist(ctx, "user-123", "Daft Punk")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not watching")
	})

	t.Run("returns error from repository", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("DeleteArtistWatch", ctx, "user-123", "Daft Punk").Return(errors.New("dynamo error"))

		svc := NewArtistWatchService(mockRepo)
		err := svc.UnwatchArtist(ctx, "user-123", "Daft Punk")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "dynamo error")
	})
}

// TestArtistWatchService_IsWatching verifies that IsWatching returns true when
// a watch record exists and false when it does not.
func TestArtistWatchService_IsWatching(t *testing.T) {
	t.Run("returns true when watch exists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		watch := &models.ArtistWatch{
			UserID:     "user-123",
			ArtistName: "Daft Punk",
			WatchedAt:  time.Now(),
		}
		mockRepo.On("GetArtistWatch", ctx, "user-123", "Daft Punk").Return(watch, nil)

		svc := NewArtistWatchService(mockRepo)
		watching, err := svc.IsWatching(ctx, "user-123", "Daft Punk")

		require.NoError(t, err)
		assert.True(t, watching)
	})

	t.Run("returns false when watch does not exist", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("GetArtistWatch", ctx, "user-123", "Daft Punk").Return((*models.ArtistWatch)(nil), repository.ErrNotFound)

		svc := NewArtistWatchService(mockRepo)
		watching, err := svc.IsWatching(ctx, "user-123", "Daft Punk")

		require.NoError(t, err)
		assert.False(t, watching)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("GetArtistWatch", ctx, "user-123", "Daft Punk").Return((*models.ArtistWatch)(nil), errors.New("dynamo error"))

		svc := NewArtistWatchService(mockRepo)
		watching, err := svc.IsWatching(ctx, "user-123", "Daft Punk")

		require.Error(t, err)
		assert.False(t, watching)
	})
}

// TestArtistWatchService_ListWatchedArtists verifies that ListWatchedArtists
// delegates to the repository and returns paginated results.
func TestArtistWatchService_ListWatchedArtists(t *testing.T) {
	t.Run("returns list of watched artists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items: []models.ArtistWatch{
				{UserID: "user-123", ArtistName: "Daft Punk", WatchedAt: time.Now()},
				{UserID: "user-123", ArtistName: "Deadmau5", WatchedAt: time.Now()},
			},
			HasMore: false,
		}
		mockRepo.On("ListWatchedArtists", ctx, "user-123", 20, "").Return(watches, nil)

		svc := NewArtistWatchService(mockRepo)
		result, err := svc.ListWatchedArtists(ctx, "user-123", 20, "")

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, "Daft Punk", result.Items[0].ArtistName)
		assert.Equal(t, "Deadmau5", result.Items[1].ArtistName)
	})

	t.Run("handles pagination", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items: []models.ArtistWatch{
				{UserID: "user-123", ArtistName: "Daft Punk", WatchedAt: time.Now()},
			},
			NextCursor: "cursor-abc",
			HasMore:    true,
		}
		mockRepo.On("ListWatchedArtists", ctx, "user-123", 1, "").Return(watches, nil)

		svc := NewArtistWatchService(mockRepo)
		result, err := svc.ListWatchedArtists(ctx, "user-123", 1, "")

		require.NoError(t, err)
		assert.True(t, result.HasMore)
		assert.Equal(t, "cursor-abc", result.NextCursor)
	})

	t.Run("handles empty list", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items:   []models.ArtistWatch{},
			HasMore: false,
		}
		mockRepo.On("ListWatchedArtists", ctx, "user-123", 20, "").Return(watches, nil)

		svc := NewArtistWatchService(mockRepo)
		result, err := svc.ListWatchedArtists(ctx, "user-123", 20, "")

		require.NoError(t, err)
		assert.Empty(t, result.Items)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistWatchRepository)

		mockRepo.On("ListWatchedArtists", ctx, "user-123", 20, "").Return((*repository.PaginatedResult[models.ArtistWatch])(nil), errors.New("dynamo error"))

		svc := NewArtistWatchService(mockRepo)
		result, err := svc.ListWatchedArtists(ctx, "user-123", 20, "")

		require.Error(t, err)
		assert.Nil(t, result)
	})
}
