package service

import (
	"context"
	"fmt"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

const (
	// DefaultSearchLimit is the default limit for user searches.
	DefaultSearchLimit = 20
)

// AdminRepository defines the repository operations needed by AdminService.
type AdminRepository interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error
	SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error)
	SetUserDisabled(ctx context.Context, userID string, disabled bool) error
	GetFollowerCount(ctx context.Context, userID string) (int, error)
}

// AdminService provides administrative user management operations.
type AdminService interface {
	// SearchUsers searches for users by email or display name.
	SearchUsers(ctx context.Context, query string, limit int) ([]models.UserSummary, error)

	// GetUserDetails returns full user details including content counts.
	GetUserDetails(ctx context.Context, userID string) (*models.UserDetails, error)

	// UpdateUserRole updates a user's role in both DynamoDB and Cognito.
	UpdateUserRole(ctx context.Context, userID string, newRole models.UserRole) error

	// UpdateUserRoleByAdmin updates a user's role, preventing self-modification.
	UpdateUserRoleByAdmin(ctx context.Context, adminID, userID string, newRole models.UserRole) error

	// SetUserStatus enables or disables a user in both DynamoDB and Cognito.
	SetUserStatus(ctx context.Context, userID string, disabled bool) error
}

// adminService implements AdminService.
type adminService struct {
	repo    AdminRepository
	cognito CognitoClient
}

// NewAdminService creates a new AdminService.
func NewAdminService(repo AdminRepository, cognito CognitoClient) AdminService {
	return &adminService{
		repo:    repo,
		cognito: cognito,
	}
}

// SearchUsers searches for users by email or display name.
func (s *adminService) SearchUsers(ctx context.Context, query string, limit int) ([]models.UserSummary, error) {
	if limit <= 0 {
		limit = DefaultSearchLimit
	}

	users, err := s.repo.SearchUsers(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	summaries := make([]models.UserSummary, len(users))
	for i, user := range users {
		summaries[i] = user.ToUserSummary(false) // Disabled status is not fetched in search
	}

	return summaries, nil
}

// GetUserDetails returns full user details including content counts.
func (s *adminService) GetUserDetails(ctx context.Context, userID string) (*models.UserDetails, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("user", userID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get follower count (optional, don't fail if error)
	followerCount, _ := s.repo.GetFollowerCount(ctx, userID)

	// Get disabled status from Cognito groups (optional)
	// For now, assume not disabled if we can't determine
	disabled := false

	details := user.ToUserDetails(disabled, nil, followerCount)
	return &details, nil
}

// UpdateUserRole updates a user's role in both DynamoDB and Cognito.
func (s *adminService) UpdateUserRole(ctx context.Context, userID string, newRole models.UserRole) error {
	// Validate the new role
	if !newRole.IsValid() {
		return fmt.Errorf("invalid role: %s", newRole)
	}

	// Get the current user
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("user", userID)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Skip if role is unchanged
	if user.Role == newRole {
		return nil
	}

	oldRole := user.Role

	// Step 1: Update DynamoDB
	if err := s.repo.UpdateUserRole(ctx, userID, newRole); err != nil {
		return fmt.Errorf("failed to update role in DynamoDB: %w", err)
	}

	// Step 2: Get current Cognito groups
	currentGroups, err := s.cognito.GetUserGroups(ctx, userID)
	if err != nil {
		// Rollback DynamoDB
		_ = s.repo.UpdateUserRole(ctx, userID, oldRole)
		return fmt.Errorf("failed to get user groups from Cognito: %w", err)
	}

	// Step 3: Remove from old groups
	for _, group := range currentGroups {
		if err := s.cognito.RemoveUserFromGroup(ctx, userID, group); err != nil {
			// Rollback DynamoDB
			_ = s.repo.UpdateUserRole(ctx, userID, oldRole)
			return fmt.Errorf("failed to remove user from group %s: %w", group, err)
		}
	}

	// Step 4: Add to new group
	newGroupName := newRole.CognitoGroupName()
	if err := s.cognito.AddUserToGroup(ctx, userID, newGroupName); err != nil {
		// Rollback DynamoDB and restore old Cognito groups
		_ = s.repo.UpdateUserRole(ctx, userID, oldRole)
		for _, group := range currentGroups {
			_ = s.cognito.AddUserToGroup(ctx, userID, group)
		}
		return fmt.Errorf("failed to add user to group %s: %w", newGroupName, err)
	}

	return nil
}

// UpdateUserRoleByAdmin updates a user's role, preventing self-modification.
func (s *adminService) UpdateUserRoleByAdmin(ctx context.Context, adminID, userID string, newRole models.UserRole) error {
	if adminID == userID {
		return fmt.Errorf("cannot modify your own role")
	}

	return s.UpdateUserRole(ctx, userID, newRole)
}

// SetUserStatus enables or disables a user in both DynamoDB and Cognito.
func (s *adminService) SetUserStatus(ctx context.Context, userID string, disabled bool) error {
	// Verify user exists
	_, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("user", userID)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Step 1: Update DynamoDB
	if err := s.repo.SetUserDisabled(ctx, userID, disabled); err != nil {
		return fmt.Errorf("failed to update user status in DynamoDB: %w", err)
	}

	// Step 2: Update Cognito
	var cognitoErr error
	if disabled {
		cognitoErr = s.cognito.DisableUser(ctx, userID)
	} else {
		cognitoErr = s.cognito.EnableUser(ctx, userID)
	}

	if cognitoErr != nil {
		// Rollback DynamoDB
		_ = s.repo.SetUserDisabled(ctx, userID, !disabled)
		return fmt.Errorf("failed to update user status in Cognito: %w", cognitoErr)
	}

	return nil
}
