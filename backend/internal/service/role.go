package service

import (
	"context"
	"fmt"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// RoleRepository defines the repository interface for role operations.
type RoleRepository interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error
	ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error)
}

// RoleService handles user role management and permission checking.
type RoleService interface {
	GetUserRole(ctx context.Context, userID string) (models.UserRole, error)
	SetUserRole(ctx context.Context, userID string, role models.UserRole) error
	HasPermission(ctx context.Context, userID string, permission models.Permission) (bool, error)
	ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error)
}

type roleService struct {
	repo RoleRepository
}

// NewRoleService creates a new RoleService.
func NewRoleService(repo RoleRepository) RoleService {
	return &roleService{repo: repo}
}

// GetUserRole retrieves the role for a user.
func (s *roleService) GetUserRole(ctx context.Context, userID string) (models.UserRole, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return "", models.NewNotFoundError("user", userID)
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Return default role if not set
	if user.Role == "" {
		return models.RoleGuest, nil
	}

	return user.Role, nil
}

// SetUserRole updates the role for a user.
func (s *roleService) SetUserRole(ctx context.Context, userID string, role models.UserRole) error {
	// Validate role
	if !role.IsValid() {
		return fmt.Errorf("invalid role: %s", role)
	}

	err := s.repo.UpdateUserRole(ctx, userID, role)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("user", userID)
		}
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

// HasPermission checks if a user has a specific permission.
func (s *roleService) HasPermission(ctx context.Context, userID string, permission models.Permission) (bool, error) {
	role, err := s.GetUserRole(ctx, userID)
	if err != nil {
		return false, err
	}

	return role.HasPermission(permission), nil
}

// ListUsersByRole lists users with a specific role.
func (s *roleService) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return s.repo.ListUsersByRole(ctx, role, limit, cursor)
}
