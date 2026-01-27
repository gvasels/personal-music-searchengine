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

// MockAdminRepository is a mock implementation for admin repository operations.
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAdminRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockAdminRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockAdminRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	args := m.Called(ctx, userID, disabled)
	return args.Error(0)
}

func (m *MockAdminRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

// MockCognitoClient is a mock implementation of CognitoClient.
type MockCognitoClient struct {
	mock.Mock
}

func (m *MockCognitoClient) AddUserToGroup(ctx context.Context, userID string, groupName string) error {
	args := m.Called(ctx, userID, groupName)
	return args.Error(0)
}

func (m *MockCognitoClient) RemoveUserFromGroup(ctx context.Context, userID string, groupName string) error {
	args := m.Called(ctx, userID, groupName)
	return args.Error(0)
}

func (m *MockCognitoClient) GetUserGroups(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCognitoClient) DisableUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCognitoClient) EnableUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCognitoClient) GetUserStatus(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCognitoClient) SearchUsers(ctx context.Context, query string, limit int) ([]CognitoUser, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CognitoUser), args.Error(1)
}

func TestNewAdminService(t *testing.T) {
	t.Run("creates service with dependencies", func(t *testing.T) {
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		svc := NewAdminService(mockRepo, mockCognito)
		require.NotNil(t, svc)
	})
}

func TestAdminService_SearchUsers(t *testing.T) {
	t.Run("returns matching users", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		now := time.Now()
		cognitoUsers := []CognitoUser{
			{
				ID:      "user-1",
				Email:   "john@example.com",
				Name:    "John Doe",
				Status:  "CONFIRMED",
				Enabled: true,
				Created: now.Format(time.RFC3339),
			},
			{
				ID:      "user-2",
				Email:   "jane@example.com",
				Name:    "Jane Smith",
				Status:  "CONFIRMED",
				Enabled: true,
				Created: now.Format(time.RFC3339),
			},
		}

		mockCognito.On("SearchUsers", ctx, "john", 20).Return(cognitoUsers, nil)
		mockRepo.On("GetUser", ctx, "user-1").Return(&models.User{
			ID:          "user-1",
			Email:       "john@example.com",
			DisplayName: "John Doe",
			Role:        models.RoleSubscriber,
		}, nil)
		mockRepo.On("GetUser", ctx, "user-2").Return(&models.User{
			ID:          "user-2",
			Email:       "jane@example.com",
			DisplayName: "Jane Smith",
			Role:        models.RoleArtist,
		}, nil)

		svc := NewAdminService(mockRepo, mockCognito)
		result, err := svc.SearchUsers(ctx, "john", 20)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "John Doe", result[0].DisplayName)
		mockCognito.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty slice when no matches", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		mockCognito.On("SearchUsers", ctx, "nonexistent", 20).Return([]CognitoUser{}, nil)

		svc := NewAdminService(mockRepo, mockCognito)
		result, err := svc.SearchUsers(ctx, "nonexistent", 20)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("applies default limit", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		mockCognito.On("SearchUsers", ctx, "test", 20).Return([]CognitoUser{}, nil)

		svc := NewAdminService(mockRepo, mockCognito)
		_, err := svc.SearchUsers(ctx, "test", 0)

		require.NoError(t, err)
		mockCognito.AssertExpectations(t)
	})

	t.Run("handles cognito error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		mockCognito.On("SearchUsers", ctx, "error", 20).Return(([]CognitoUser)(nil), errors.New("cognito error"))

		svc := NewAdminService(mockRepo, mockCognito)
		_, err := svc.SearchUsers(ctx, "error", 20)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cognito error")
	})
}

func TestAdminService_GetUserDetails(t *testing.T) {
	t.Run("returns full user details", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		now := time.Now()
		user := &models.User{
			ID:            "user-123",
			Email:         "test@example.com",
			DisplayName:   "Test User",
			Role:          models.RoleSubscriber,
			TrackCount:    10,
			PlaylistCount: 3,
			AlbumCount:    2,
			StorageUsed:   1073741824, // 1GB
			Timestamps:    models.Timestamps{CreatedAt: now},
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("GetFollowerCount", ctx, "user-123").Return(5, nil)
		mockCognito.On("GetUserStatus", ctx, "user-123").Return(true, nil)

		svc := NewAdminService(mockRepo, mockCognito)
		details, err := svc.GetUserDetails(ctx, "user-123")

		require.NoError(t, err)
		assert.Equal(t, "user-123", details.ID)
		assert.Equal(t, "Test User", details.DisplayName)
		assert.Equal(t, 10, details.TrackCount)
		assert.Equal(t, 5, details.FollowerCount)
		mockRepo.AssertExpectations(t)
		mockCognito.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		mockRepo.On("GetUser", ctx, "nonexistent").Return((*models.User)(nil), repository.ErrNotFound)

		svc := NewAdminService(mockRepo, mockCognito)
		_, err := svc.GetUserDetails(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("handles follower count error gracefully", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:          "user-123",
			Email:       "test@example.com",
			DisplayName: "Test User",
			Role:        models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("GetFollowerCount", ctx, "user-123").Return(0, errors.New("count error"))
		mockCognito.On("GetUserStatus", ctx, "user-123").Return(true, nil)

		svc := NewAdminService(mockRepo, mockCognito)
		details, err := svc.GetUserDetails(ctx, "user-123")

		// Should still return user details with 0 followers
		require.NoError(t, err)
		assert.Equal(t, 0, details.FollowerCount)
	})
}

func TestAdminService_UpdateUserRole(t *testing.T) {
	t.Run("updates role in both DynamoDB and Cognito", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("UpdateUserRole", ctx, "user-123", models.RoleArtist).Return(nil)
		mockCognito.On("GetUserGroups", ctx, "user-123").Return([]string{"subscriber"}, nil)
		mockCognito.On("RemoveUserFromGroup", ctx, "user-123", "subscriber").Return(nil)
		mockCognito.On("AddUserToGroup", ctx, "user-123", "artist").Return(nil)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.UpdateUserRole(ctx, "user-123", models.RoleArtist)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCognito.AssertExpectations(t)
	})

	t.Run("rollbacks DynamoDB on Cognito failure", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("UpdateUserRole", ctx, "user-123", models.RoleArtist).Return(nil)
		mockCognito.On("GetUserGroups", ctx, "user-123").Return([]string{"subscriber"}, nil)
		mockCognito.On("RemoveUserFromGroup", ctx, "user-123", "subscriber").Return(nil)
		mockCognito.On("AddUserToGroup", ctx, "user-123", "artist").Return(errors.New("cognito error"))
		// Rollback
		mockRepo.On("UpdateUserRole", ctx, "user-123", models.RoleSubscriber).Return(nil)
		mockCognito.On("AddUserToGroup", ctx, "user-123", "subscriber").Return(nil)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.UpdateUserRole(ctx, "user-123", models.RoleArtist)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cognito error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error for invalid role", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.UpdateUserRole(ctx, "user-123", models.UserRole("invalid"))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		mockRepo.On("GetUser", ctx, "nonexistent").Return((*models.User)(nil), repository.ErrNotFound)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.UpdateUserRole(ctx, "nonexistent", models.RoleArtist)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("prevents changing own role", func(t *testing.T) {
		// Note: The userID context key check is done at handler level, not service level.
		// This test verifies the UpdateUserRoleByAdmin method prevents self-modification.
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.UpdateUserRoleByAdmin(ctx, "admin-123", "admin-123", models.RoleSubscriber)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot modify your own role")
	})

	t.Run("skips Cognito update when role unchanged", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleArtist,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.UpdateUserRole(ctx, "user-123", models.RoleArtist)

		require.NoError(t, err)
		// Cognito should not be called when role is unchanged
		mockCognito.AssertNotCalled(t, "AddUserToGroup")
		mockCognito.AssertNotCalled(t, "RemoveUserFromGroup")
	})
}

func TestAdminService_SetUserStatus(t *testing.T) {
	t.Run("disables user in DynamoDB and Cognito", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("SetUserDisabled", ctx, "user-123", true).Return(nil)
		mockCognito.On("DisableUser", ctx, "user-123").Return(nil)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.SetUserStatus(ctx, "user-123", true)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCognito.AssertExpectations(t)
	})

	t.Run("enables user in DynamoDB and Cognito", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("SetUserDisabled", ctx, "user-123", false).Return(nil)
		mockCognito.On("EnableUser", ctx, "user-123").Return(nil)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.SetUserStatus(ctx, "user-123", false)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCognito.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		mockRepo.On("GetUser", ctx, "nonexistent").Return((*models.User)(nil), repository.ErrNotFound)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.SetUserStatus(ctx, "nonexistent", true)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("rollbacks DynamoDB on Cognito failure when disabling", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockAdminRepository)
		mockCognito := new(MockCognitoClient)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)
		mockRepo.On("SetUserDisabled", ctx, "user-123", true).Return(nil)
		mockCognito.On("DisableUser", ctx, "user-123").Return(errors.New("cognito error"))
		// Rollback
		mockRepo.On("SetUserDisabled", ctx, "user-123", false).Return(nil)

		svc := NewAdminService(mockRepo, mockCognito)
		err := svc.SetUserStatus(ctx, "user-123", true)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cognito error")
		mockRepo.AssertExpectations(t)
	})
}
