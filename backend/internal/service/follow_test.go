package service

import (
	"context"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFollowRepository is a mock implementation for follow operations
type MockFollowRepository struct {
	mock.Mock
}

func (m *MockFollowRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	args := m.Called(ctx, follow)
	return args.Error(0)
}

func (m *MockFollowRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	args := m.Called(ctx, followerID, followedID)
	return args.Error(0)
}

func (m *MockFollowRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	args := m.Called(ctx, followerID, followedID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Follow), args.Error(1)
}

func (m *MockFollowRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	args := m.Called(ctx, userID, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Follow]), args.Error(1)
}

func (m *MockFollowRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	args := m.Called(ctx, userID, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Follow]), args.Error(1)
}

func (m *MockFollowRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func (m *MockFollowRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func (m *MockFollowRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistProfile), args.Error(1)
}

func TestNewFollowService(t *testing.T) {
	t.Run("creates service with repository", func(t *testing.T) {
		mockRepo := new(MockFollowRepository)
		svc := NewFollowService(mockRepo)
		require.NotNil(t, svc)
	})
}

func TestFollowService_Follow(t *testing.T) {
	t.Run("creates follow relationship", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		profile := &models.ArtistProfile{
			UserID:      "artist-123",
			DisplayName: "Test Artist",
		}

		mockRepo.On("GetArtistProfile", ctx, "artist-123").Return(profile, nil)
		mockRepo.On("CreateFollow", ctx, mock.MatchedBy(func(f models.Follow) bool {
			return f.FollowerID == "user-123" && f.FollowedID == "artist-123"
		})).Return(nil)
		mockRepo.On("IncrementArtistFollowerCount", ctx, "artist-123", 1).Return(nil)
		mockRepo.On("IncrementUserFollowingCount", ctx, "user-123", 1).Return(nil)

		svc := NewFollowService(mockRepo)
		err := svc.Follow(ctx, "user-123", "artist-123")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("prevents self-follow", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		svc := NewFollowService(mockRepo)
		err := svc.Follow(ctx, "user-123", "user-123")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot follow yourself")
	})

	t.Run("returns error when artist profile not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		mockRepo.On("GetArtistProfile", ctx, "nonexistent").Return((*models.ArtistProfile)(nil), repository.ErrNotFound)

		svc := NewFollowService(mockRepo)
		err := svc.Follow(ctx, "user-123", "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "artist profile not found")
	})

	t.Run("returns error when already following", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		profile := &models.ArtistProfile{
			UserID:      "artist-123",
			DisplayName: "Test Artist",
		}

		mockRepo.On("GetArtistProfile", ctx, "artist-123").Return(profile, nil)
		mockRepo.On("CreateFollow", ctx, mock.Anything).Return(repository.ErrAlreadyExists)

		svc := NewFollowService(mockRepo)
		err := svc.Follow(ctx, "user-123", "artist-123")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "already following")
	})
}

func TestFollowService_Unfollow(t *testing.T) {
	t.Run("deletes follow relationship", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		mockRepo.On("DeleteFollow", ctx, "user-123", "artist-123").Return(nil)
		mockRepo.On("IncrementArtistFollowerCount", ctx, "artist-123", -1).Return(nil)
		mockRepo.On("IncrementUserFollowingCount", ctx, "user-123", -1).Return(nil)

		svc := NewFollowService(mockRepo)
		err := svc.Unfollow(ctx, "user-123", "artist-123")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when not following", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		mockRepo.On("DeleteFollow", ctx, "user-123", "artist-123").Return(repository.ErrNotFound)

		svc := NewFollowService(mockRepo)
		err := svc.Unfollow(ctx, "user-123", "artist-123")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not following")
	})
}

func TestFollowService_IsFollowing(t *testing.T) {
	t.Run("returns true when following", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		follow := &models.Follow{
			FollowerID: "user-123",
			FollowedID: "artist-123",
			CreatedAt:      time.Now(),
		}

		mockRepo.On("GetFollow", ctx, "user-123", "artist-123").Return(follow, nil)

		svc := NewFollowService(mockRepo)
		following, err := svc.IsFollowing(ctx, "user-123", "artist-123")

		require.NoError(t, err)
		assert.True(t, following)
	})

	t.Run("returns false when not following", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		mockRepo.On("GetFollow", ctx, "user-123", "artist-123").Return((*models.Follow)(nil), repository.ErrNotFound)

		svc := NewFollowService(mockRepo)
		following, err := svc.IsFollowing(ctx, "user-123", "artist-123")

		require.NoError(t, err)
		assert.False(t, following)
	})
}

func TestFollowService_GetFollowers(t *testing.T) {
	t.Run("returns list of followers", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		follows := &repository.PaginatedResult[models.Follow]{
			Items: []models.Follow{
				{FollowerID: "user-1", FollowedID: "artist-123"},
				{FollowerID: "user-2", FollowedID: "artist-123"},
			},
			HasMore: false,
		}

		mockRepo.On("ListFollowers", ctx, "artist-123", 20, "").Return(follows, nil)

		svc := NewFollowService(mockRepo)
		result, err := svc.GetFollowers(ctx, "artist-123", 20, "")

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
	})

	t.Run("handles pagination", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		follows := &repository.PaginatedResult[models.Follow]{
			Items: []models.Follow{
				{FollowerID: "user-1", FollowedID: "artist-123"},
			},
			NextCursor: "cursor-123",
			HasMore:    true,
		}

		mockRepo.On("ListFollowers", ctx, "artist-123", 1, "").Return(follows, nil)

		svc := NewFollowService(mockRepo)
		result, err := svc.GetFollowers(ctx, "artist-123", 1, "")

		require.NoError(t, err)
		assert.True(t, result.HasMore)
		assert.Equal(t, "cursor-123", result.NextCursor)
	})
}

func TestFollowService_GetFollowing(t *testing.T) {
	t.Run("returns list of users being followed", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		follows := &repository.PaginatedResult[models.Follow]{
			Items: []models.Follow{
				{FollowerID: "user-123", FollowedID: "artist-1"},
				{FollowerID: "user-123", FollowedID: "artist-2"},
			},
			HasMore: false,
		}

		mockRepo.On("ListFollowing", ctx, "user-123", 20, "").Return(follows, nil)

		svc := NewFollowService(mockRepo)
		result, err := svc.GetFollowing(ctx, "user-123", 20, "")

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
	})

	t.Run("handles empty following list", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockFollowRepository)

		follows := &repository.PaginatedResult[models.Follow]{
			Items:   []models.Follow{},
			HasMore: false,
		}

		mockRepo.On("ListFollowing", ctx, "user-123", 20, "").Return(follows, nil)

		svc := NewFollowService(mockRepo)
		result, err := svc.GetFollowing(ctx, "user-123", 20, "")

		require.NoError(t, err)
		assert.Empty(t, result.Items)
	})
}
