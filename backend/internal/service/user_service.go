package service

import (
	"context"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// UserService handles user-related operations
type UserService struct {
	repo repository.Repository
}

// NewUserService creates a new UserService
func NewUserService(repo repository.Repository) *UserService {
	return &UserService{repo: repo}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetUser(ctx, userID)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest, userID string) (*models.User, error) {
	now := time.Now()
	user := &models.User{
		ID:          userID,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
		StorageUsed:   0,
		StorageLimit:  10 * 1024 * 1024 * 1024, // 10 GB default
		TrackCount:    0,
		AlbumCount:    0,
		PlaylistCount: 0,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user's profile
func (s *UserService) UpdateUser(ctx context.Context, userID string, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetOrCreateUser gets an existing user or creates a new one
func (s *UserService) GetOrCreateUser(ctx context.Context, userID, email string) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err == nil {
		return user, nil
	}

	// User doesn't exist, create new
	req := models.CreateUserRequest{
		Email:       email,
		DisplayName: email, // Default display name to email
	}
	return s.CreateUser(ctx, req, userID)
}
