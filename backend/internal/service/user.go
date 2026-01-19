package service

import (
	"context"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

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

	// Create new user
	now := time.Now()
	user = &models.User{
		ID:            userID,
		Email:         email,
		DisplayName:   displayName,
		StorageUsed:   0,
		StorageLimit:  10 * 1024 * 1024 * 1024, // 10 GB default limit
		TrackCount:    0,
		AlbumCount:    0,
		PlaylistCount: 0,
	}
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
