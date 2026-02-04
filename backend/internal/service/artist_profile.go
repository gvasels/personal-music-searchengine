package service

import (
	"context"
	"fmt"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// ArtistProfileRepository defines the repository interface for artist profile operations.
type ArtistProfileRepository interface {
	CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error
	GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error)
	UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error
	DeleteArtistProfile(ctx context.Context, userID string) error
	ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error)
	IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

// ArtistProfileService handles artist profile business logic.
type ArtistProfileService interface {
	CreateProfile(ctx context.Context, userID string, req models.CreateArtistProfileRequest) (*models.ArtistProfileResponse, error)
	GetProfile(ctx context.Context, userID string) (*models.ArtistProfileResponse, error)
	UpdateProfile(ctx context.Context, requestingUserID, profileUserID string, req models.UpdateArtistProfileRequest) (*models.ArtistProfileResponse, error)
	DeleteProfile(ctx context.Context, requestingUserID, profileUserID string) error
	ListProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfileResponse], error)
}

type artistProfileService struct {
	repo ArtistProfileRepository
}

// NewArtistProfileService creates a new ArtistProfileService.
func NewArtistProfileService(repo ArtistProfileRepository) ArtistProfileService {
	return &artistProfileService{repo: repo}
}

// CreateProfile creates a new artist profile for a user with artist role.
func (s *artistProfileService) CreateProfile(ctx context.Context, userID string, req models.CreateArtistProfileRequest) (*models.ArtistProfileResponse, error) {
	// Check user has artist role
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("user", userID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Role != models.RoleArtist && user.Role != models.RoleAdmin {
		return nil, fmt.Errorf("artist role required to create profile")
	}

	// Create profile
	profile := models.NewArtistProfile(userID)
	profile.DisplayName = req.DisplayName
	profile.Bio = req.Bio
	if req.SocialLinks != nil {
		profile.SocialLinks = req.SocialLinks
	}
	if req.Genres != nil {
		profile.Genres = req.Genres
	}

	err = s.repo.CreateArtistProfile(ctx, *profile)
	if err != nil {
		if err == repository.ErrAlreadyExists {
			return nil, fmt.Errorf("artist profile already exists for user")
		}
		return nil, fmt.Errorf("failed to create artist profile: %w", err)
	}

	response := profile.ToResponse()
	return &response, nil
}

// GetProfile retrieves an artist profile by user ID.
func (s *artistProfileService) GetProfile(ctx context.Context, userID string) (*models.ArtistProfileResponse, error) {
	profile, err := s.repo.GetArtistProfile(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("artist profile", userID)
		}
		return nil, fmt.Errorf("failed to get artist profile: %w", err)
	}

	response := profile.ToResponse()
	return &response, nil
}

// UpdateProfile updates an artist profile. Only the owner can update their profile.
func (s *artistProfileService) UpdateProfile(ctx context.Context, requestingUserID, profileUserID string, req models.UpdateArtistProfileRequest) (*models.ArtistProfileResponse, error) {
	// Check ownership
	if requestingUserID != profileUserID {
		return nil, fmt.Errorf("forbidden: only the profile owner can update their profile")
	}

	// Get existing profile
	profile, err := s.repo.GetArtistProfile(ctx, profileUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, models.NewNotFoundError("artist profile", profileUserID)
		}
		return nil, fmt.Errorf("failed to get artist profile: %w", err)
	}

	// Apply updates
	if req.DisplayName != nil {
		profile.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		profile.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = *req.AvatarURL
	}
	if req.BannerURL != nil {
		profile.BannerURL = *req.BannerURL
	}
	if req.SocialLinks != nil {
		profile.SocialLinks = *req.SocialLinks
	}
	if req.Genres != nil {
		profile.Genres = *req.Genres
	}

	err = s.repo.UpdateArtistProfile(ctx, *profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update artist profile: %w", err)
	}

	response := profile.ToResponse()
	return &response, nil
}

// DeleteProfile deletes an artist profile. Only the owner can delete their profile.
func (s *artistProfileService) DeleteProfile(ctx context.Context, requestingUserID, profileUserID string) error {
	// Check ownership
	if requestingUserID != profileUserID {
		return fmt.Errorf("forbidden: only the profile owner can delete their profile")
	}

	err := s.repo.DeleteArtistProfile(ctx, profileUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return models.NewNotFoundError("artist profile", profileUserID)
		}
		return fmt.Errorf("failed to delete artist profile: %w", err)
	}

	return nil
}

// ListProfiles lists all artist profiles for discovery.
func (s *artistProfileService) ListProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfileResponse], error) {
	result, err := s.repo.ListArtistProfiles(ctx, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to list artist profiles: %w", err)
	}

	responses := make([]models.ArtistProfileResponse, len(result.Items))
	for i, profile := range result.Items {
		responses[i] = profile.ToResponse()
	}

	return &repository.PaginatedResult[models.ArtistProfileResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}
