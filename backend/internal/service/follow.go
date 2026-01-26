package service

import (
	"context"
	"fmt"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// FollowRepository defines the repository interface for follow operations.
type FollowRepository interface {
	CreateFollow(ctx context.Context, follow models.Follow) error
	DeleteFollow(ctx context.Context, followerID, followedID string) error
	GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error)
	ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error)
	ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error)
	IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error
	IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error
	GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error)
}

// FollowService handles follow/unfollow operations.
type FollowService interface {
	Follow(ctx context.Context, followerUserID, followedUserID string) error
	Unfollow(ctx context.Context, followerUserID, followedUserID string) error
	IsFollowing(ctx context.Context, followerUserID, followedUserID string) (bool, error)
	GetFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error)
	GetFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error)
}

type followService struct {
	repo FollowRepository
}

// NewFollowService creates a new FollowService.
func NewFollowService(repo FollowRepository) FollowService {
	return &followService{repo: repo}
}

// Follow creates a follow relationship between two users.
func (s *followService) Follow(ctx context.Context, followerUserID, followedUserID string) error {
	// Prevent self-follow
	if followerUserID == followedUserID {
		return fmt.Errorf("cannot follow yourself")
	}

	// Verify the followed user has an artist profile
	_, err := s.repo.GetArtistProfile(ctx, followedUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return fmt.Errorf("artist profile not found")
		}
		return fmt.Errorf("failed to verify artist profile: %w", err)
	}

	// Create the follow relationship
	follow := models.Follow{
		FollowerID: followerUserID,
		FollowedID: followedUserID,
	}

	err = s.repo.CreateFollow(ctx, follow)
	if err != nil {
		if err == repository.ErrAlreadyExists {
			return fmt.Errorf("already following this user")
		}
		return fmt.Errorf("failed to create follow: %w", err)
	}

	// Increment counts (best effort - don't fail if these error)
	_ = s.repo.IncrementArtistFollowerCount(ctx, followedUserID, 1)
	_ = s.repo.IncrementUserFollowingCount(ctx, followerUserID, 1)

	return nil
}

// Unfollow removes a follow relationship between two users.
func (s *followService) Unfollow(ctx context.Context, followerUserID, followedUserID string) error {
	err := s.repo.DeleteFollow(ctx, followerUserID, followedUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return fmt.Errorf("not following this user")
		}
		return fmt.Errorf("failed to delete follow: %w", err)
	}

	// Decrement counts (best effort - don't fail if these error)
	_ = s.repo.IncrementArtistFollowerCount(ctx, followedUserID, -1)
	_ = s.repo.IncrementUserFollowingCount(ctx, followerUserID, -1)

	return nil
}

// IsFollowing checks if a user is following another user.
func (s *followService) IsFollowing(ctx context.Context, followerUserID, followedUserID string) (bool, error) {
	_, err := s.repo.GetFollow(ctx, followerUserID, followedUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get follow: %w", err)
	}
	return true, nil
}

// GetFollowers returns users who follow a given user.
func (s *followService) GetFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return s.repo.ListFollowers(ctx, userID, limit, cursor)
}

// GetFollowing returns users that a given user follows.
func (s *followService) GetFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return s.repo.ListFollowing(ctx, userID, limit, cursor)
}
