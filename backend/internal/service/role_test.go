package service

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRoleRepository is a mock implementation for role operations
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRoleRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockRoleRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	args := m.Called(ctx, role, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.User]), args.Error(1)
}

func TestNewRoleService(t *testing.T) {
	t.Run("creates service with repository", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		svc := NewRoleService(mockRepo)
		require.NotNil(t, svc)
	})
}

func TestRoleService_GetUserRole(t *testing.T) {
	t.Run("returns user role", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewRoleService(mockRepo)
		role, err := svc.GetUserRole(ctx, "user-123")

		require.NoError(t, err)
		assert.Equal(t, models.RoleSubscriber, role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		mockRepo.On("GetUser", ctx, "nonexistent").Return((*models.User)(nil), repository.ErrNotFound)

		svc := NewRoleService(mockRepo)
		_, err := svc.GetUserRole(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("returns default role for user without role", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: "", // Empty role
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewRoleService(mockRepo)
		role, err := svc.GetUserRole(ctx, "user-123")

		require.NoError(t, err)
		assert.Equal(t, models.RoleGuest, role)
	})
}

func TestRoleService_SetUserRole(t *testing.T) {
	t.Run("updates user role", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		mockRepo.On("UpdateUserRole", ctx, "user-123", models.RoleArtist).Return(nil)

		svc := NewRoleService(mockRepo)
		err := svc.SetUserRole(ctx, "user-123", models.RoleArtist)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error for invalid role", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		svc := NewRoleService(mockRepo)
		err := svc.SetUserRole(ctx, "user-123", models.UserRole("invalid"))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		mockRepo.On("UpdateUserRole", ctx, "nonexistent", models.RoleSubscriber).Return(repository.ErrNotFound)

		svc := NewRoleService(mockRepo)
		err := svc.SetUserRole(ctx, "nonexistent", models.RoleSubscriber)

		require.Error(t, err)
	})
}

func TestRoleService_HasPermission(t *testing.T) {
	t.Run("subscriber has listen permission", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewRoleService(mockRepo)
		has, err := svc.HasPermission(ctx, "user-123", models.PermissionListen)

		require.NoError(t, err)
		assert.True(t, has)
	})

	t.Run("subscriber does not have publish permission", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleSubscriber,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewRoleService(mockRepo)
		has, err := svc.HasPermission(ctx, "user-123", models.PermissionPublishTracks)

		require.NoError(t, err)
		assert.False(t, has)
	})

	t.Run("artist has publish permission", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleArtist,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewRoleService(mockRepo)
		has, err := svc.HasPermission(ctx, "user-123", models.PermissionPublishTracks)

		require.NoError(t, err)
		assert.True(t, has)
	})

	t.Run("admin has all permissions", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleAdmin,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil)

		svc := NewRoleService(mockRepo)
		has, err := svc.HasPermission(ctx, "user-123", models.PermissionManageUsers)

		require.NoError(t, err)
		assert.True(t, has)
	})

	t.Run("guest has limited permissions", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		user := &models.User{
			ID:   "user-123",
			Role: models.RoleGuest,
		}

		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil).Once()
		mockRepo.On("GetUser", ctx, "user-123").Return(user, nil).Once()

		svc := NewRoleService(mockRepo)

		// Guest can browse
		has, err := svc.HasPermission(ctx, "user-123", models.PermissionBrowse)
		require.NoError(t, err)
		assert.True(t, has)

		// Guest cannot listen
		has, err = svc.HasPermission(ctx, "user-123", models.PermissionListen)
		require.NoError(t, err)
		assert.False(t, has)
	})
}

func TestRoleService_ListUsersByRole(t *testing.T) {
	t.Run("lists users by role", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		users := &repository.PaginatedResult[models.User]{
			Items: []models.User{
				{ID: "user-1", Role: models.RoleArtist},
				{ID: "user-2", Role: models.RoleArtist},
			},
			HasMore: false,
		}

		mockRepo.On("ListUsersByRole", ctx, models.RoleArtist, 10, "").Return(users, nil)

		svc := NewRoleService(mockRepo)
		result, err := svc.ListUsersByRole(ctx, models.RoleArtist, 10, "")

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("handles pagination", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockRoleRepository)

		users := &repository.PaginatedResult[models.User]{
			Items: []models.User{
				{ID: "user-1", Role: models.RoleSubscriber},
			},
			NextCursor: "next-cursor",
			HasMore:    true,
		}

		mockRepo.On("ListUsersByRole", ctx, models.RoleSubscriber, 1, "cursor").Return(users, nil)

		svc := NewRoleService(mockRepo)
		result, err := svc.ListUsersByRole(ctx, models.RoleSubscriber, 1, "cursor")

		require.NoError(t, err)
		assert.True(t, result.HasMore)
		assert.Equal(t, "next-cursor", result.NextCursor)
	})
}
