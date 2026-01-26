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

// MockArtistProfileRepository is a mock implementation for artist profile operations
type MockArtistProfileRepository struct {
	mock.Mock
}

func (m *MockArtistProfileRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockArtistProfileRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistProfile), args.Error(1)
}

func (m *MockArtistProfileRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockArtistProfileRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockArtistProfileRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	args := m.Called(ctx, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.ArtistProfile]), args.Error(1)
}

func (m *MockArtistProfileRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func (m *MockArtistProfileRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestNewArtistProfileService(t *testing.T) {
	t.Run("creates service with repository", func(t *testing.T) {
		mockRepo := new(MockArtistProfileRepository)
		svc := NewArtistProfileService(mockRepo)
		require.NotNil(t, svc)
	})
}

func TestArtistProfileService_CreateProfile(t *testing.T) {
	t.Run("creates profile for artist user", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleArtist,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("CreateArtistProfile", ctx, mock.MatchedBy(func(p models.ArtistProfile) bool {
			return p.UserID == "user-123" && p.DisplayName == "Test Artist"
		})).Return(nil)

		svc := NewArtistProfileService(mockRepo)
		profile, err := svc.CreateProfile(ctx, "user-123", models.CreateArtistProfileRequest{
			DisplayName: "Test Artist",
			Bio:         "Test bio",
		})

		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, "Test Artist", profile.DisplayName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("rejects non-artist user", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewArtistProfileService(mockRepo)
		_, err := svc.CreateProfile(ctx, "user-123", models.CreateArtistProfileRequest{
			DisplayName: "Test Artist",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "artist role required")
	})

	t.Run("returns error when profile already exists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleArtist,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("CreateArtistProfile", ctx, mock.Anything).Return(repository.ErrAlreadyExists)

		svc := NewArtistProfileService(mockRepo)
		_, err := svc.CreateProfile(ctx, "user-123", models.CreateArtistProfileRequest{
			DisplayName: "Test Artist",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestArtistProfileService_GetProfile(t *testing.T) {
	t.Run("returns profile by user ID", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		profile := &models.ArtistProfile{
			UserID:        "user-123",
			DisplayName:   "Test Artist",
			Bio:           "Test bio",
			FollowerCount: 100,
			Timestamps: models.Timestamps{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockRepo.On("GetArtistProfile", ctx, "user-123").Return(profile, nil)

		svc := NewArtistProfileService(mockRepo)
		result, err := svc.GetProfile(ctx, "user-123")

		require.NoError(t, err)
		assert.Equal(t, "Test Artist", result.DisplayName)
		assert.Equal(t, 100, result.FollowerCount)
	})

	t.Run("returns error when profile not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		mockRepo.On("GetArtistProfile", ctx, "nonexistent").Return((*models.ArtistProfile)(nil), repository.ErrNotFound)

		svc := NewArtistProfileService(mockRepo)
		_, err := svc.GetProfile(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestArtistProfileService_UpdateProfile(t *testing.T) {
	t.Run("updates profile as owner", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		existingProfile := &models.ArtistProfile{
			UserID:      "user-123",
			DisplayName: "Old Name",
			Bio:         "Old bio",
		}

		mockRepo.On("GetArtistProfile", ctx, "user-123").Return(existingProfile, nil)
		mockRepo.On("UpdateArtistProfile", ctx, mock.MatchedBy(func(p models.ArtistProfile) bool {
			return p.UserID == "user-123" && p.DisplayName == "New Name" && p.Bio == "New bio"
		})).Return(nil)

		svc := NewArtistProfileService(mockRepo)
		result, err := svc.UpdateProfile(ctx, "user-123", "user-123", models.UpdateArtistProfileRequest{
			DisplayName: stringPtr("New Name"),
			Bio:         stringPtr("New bio"),
		})

		require.NoError(t, err)
		assert.Equal(t, "New Name", result.DisplayName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("rejects update from non-owner", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		svc := NewArtistProfileService(mockRepo)
		_, err := svc.UpdateProfile(ctx, "other-user", "user-123", models.UpdateArtistProfileRequest{
			DisplayName: stringPtr("New Name"),
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})

	t.Run("partially updates profile", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		existingProfile := &models.ArtistProfile{
			UserID:      "user-123",
			DisplayName: "Original Name",
			Bio:         "Original bio",
		}

		mockRepo.On("GetArtistProfile", ctx, "user-123").Return(existingProfile, nil)
		mockRepo.On("UpdateArtistProfile", ctx, mock.MatchedBy(func(p models.ArtistProfile) bool {
			return p.DisplayName == "Original Name" && p.Bio == "Updated bio only"
		})).Return(nil)

		svc := NewArtistProfileService(mockRepo)
		result, err := svc.UpdateProfile(ctx, "user-123", "user-123", models.UpdateArtistProfileRequest{
			Bio: stringPtr("Updated bio only"),
		})

		require.NoError(t, err)
		assert.Equal(t, "Original Name", result.DisplayName)
	})
}

func TestArtistProfileService_DeleteProfile(t *testing.T) {
	t.Run("deletes profile as owner", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		mockRepo.On("DeleteArtistProfile", ctx, "user-123").Return(nil)

		svc := NewArtistProfileService(mockRepo)
		err := svc.DeleteProfile(ctx, "user-123", "user-123")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("rejects delete from non-owner", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		svc := NewArtistProfileService(mockRepo)
		err := svc.DeleteProfile(ctx, "other-user", "user-123")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}

func TestArtistProfileService_ListProfiles(t *testing.T) {
	t.Run("lists artist profiles for discovery", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		profiles := &repository.PaginatedResult[models.ArtistProfile]{
			Items: []models.ArtistProfile{
				{UserID: "user-1", DisplayName: "Artist 1"},
				{UserID: "user-2", DisplayName: "Artist 2"},
			},
			HasMore: false,
		}

		mockRepo.On("ListArtistProfiles", ctx, 20, "").Return(profiles, nil)

		svc := NewArtistProfileService(mockRepo)
		result, err := svc.ListProfiles(ctx, 20, "")

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
	})

	t.Run("handles pagination", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockArtistProfileRepository)

		profiles := &repository.PaginatedResult[models.ArtistProfile]{
			Items: []models.ArtistProfile{
				{UserID: "user-1", DisplayName: "Artist 1"},
			},
			NextCursor: "cursor-123",
			HasMore:    true,
		}

		mockRepo.On("ListArtistProfiles", ctx, 1, "").Return(profiles, nil)

		svc := NewArtistProfileService(mockRepo)
		result, err := svc.ListProfiles(ctx, 1, "")

		require.NoError(t, err)
		assert.True(t, result.HasMore)
		assert.Equal(t, "cursor-123", result.NextCursor)
	})
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
