package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// Common service errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrValidation   = errors.New("validation error")
)

// UserSettingsUpdateInput represents a partial update to user settings
type UserSettingsUpdateInput struct {
	Notifications *models.NotificationSettings `json:"notifications,omitempty"`
	Privacy       *models.PrivacySettings      `json:"privacy,omitempty"`
	Player        *models.PlayerSettings       `json:"player,omitempty"`
	Library       *models.LibrarySettings      `json:"library,omitempty"`
}

// userService implements UserService
type userService struct {
	repo repository.Repository
}

// NewUserService creates a new user service
func NewUserService(repo repository.Repository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("User", userID)
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("User", userID)
		}
		return nil, err
	}

	// Apply updates
	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}

	if err := s.repo.UpdateUser(ctx, *user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) CreateUserIfNotExists(ctx context.Context, userID, email, displayName string) (*models.User, error) {
	// Try to get existing user
	user, err := s.repo.GetUser(ctx, userID)
	if err == nil {
		return user, nil
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	// Create new user with default settings
	user = &models.User{
		ID:            userID,
		Email:         email,
		DisplayName:   displayName,
		Role:          models.RoleSubscriber,
		Settings:      models.DefaultUserSettings(),
		StorageUsed:   0,
		StorageLimit:  10 * 1024 * 1024 * 1024, // 10 GB default limit
		TrackCount:    0,
		AlbumCount:    0,
		PlaylistCount: 0,
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := s.repo.CreateUser(ctx, *user); err != nil {
		// If user was created by another request, get it
		existingUser, getErr := s.repo.GetUser(ctx, userID)
		if getErr == nil {
			return existingUser, nil
		}
		return nil, err
	}

	return user, nil
}

// GetSettings returns the user's settings
func (s *userService) GetSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	settings, err := s.repo.GetUserSettings(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound || err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return settings, nil
}

// UpdateSettings performs a partial update of user settings
func (s *userService) UpdateSettings(ctx context.Context, userID string, input *UserSettingsUpdateInput) (*models.UserSettings, error) {
	update := &repository.UserSettingsUpdate{
		Notifications: input.Notifications,
		Privacy:       input.Privacy,
		Player:        input.Player,
		Library:       input.Library,
	}

	settings, err := s.repo.UpdateUserSettings(ctx, userID, update)
	if err != nil {
		if err == repository.ErrNotFound || err == repository.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		// Check for validation errors
		if strings.HasPrefix(err.Error(), "invalid") {
			return nil, ErrValidation
		}
		return nil, err
	}
	return settings, nil
}

// CreateUserFromCognito creates a new user from Cognito signup event
func (s *userService) CreateUserFromCognito(ctx context.Context, cognitoSub, email, displayName string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUser(ctx, cognitoSub)
	if err == nil {
		return existingUser, nil
	}
	if err != repository.ErrNotFound && err != repository.ErrUserNotFound {
		return nil, err
	}

	// Create new user from Cognito data
	user := models.NewUserFromCognito(cognitoSub, email, displayName)
	user.StorageLimit = 10 * 1024 * 1024 * 1024 // 10 GB default limit

	if err := s.repo.CreateUser(ctx, user); err != nil {
		// If user was created by another request, get it
		existingUser, getErr := s.repo.GetUser(ctx, cognitoSub)
		if getErr == nil {
			return existingUser, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetUserRole returns the user's current role from the database.
// This allows real-time role checking without requiring re-login.
func (s *userService) GetUserRole(ctx context.Context, userID string) (models.UserRole, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound || err == repository.ErrUserNotFound {
			// User not in DB yet - default to subscriber
			return models.RoleSubscriber, nil
		}
		return models.RoleGuest, err
	}
	return user.Role, nil
}
