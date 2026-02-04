package service

import (
	"context"
	"sync"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// FeatureRepository defines the repository interface for feature flags
type FeatureRepository interface {
	GetFeatureFlag(ctx context.Context, key models.FeatureKey) (*models.FeatureFlag, error)
	ListFeatureFlags(ctx context.Context) ([]models.FeatureFlag, error)
	PutFeatureFlag(ctx context.Context, flag models.FeatureFlag) error
	DeleteFeatureFlag(ctx context.Context, key models.FeatureKey) error
	GetUserFeatureOverride(ctx context.Context, userID string, key models.FeatureKey) (*models.UserFeatureOverride, error)
	ListUserFeatureOverrides(ctx context.Context, userID string) ([]models.UserFeatureOverride, error)
	PutUserFeatureOverride(ctx context.Context, override models.UserFeatureOverride) error
	DeleteUserFeatureOverride(ctx context.Context, userID string, key models.FeatureKey) error
	SeedDefaultFeatureFlags(ctx context.Context) error
}

// UserRepository interface for getting user tier
type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

// featureCacheEntry represents a cached feature flag
type featureCacheEntry struct {
	flag      *models.FeatureFlag
	fetchedAt time.Time
}

// FeatureService handles feature flag operations with caching
type FeatureService struct {
	repo     FeatureRepository
	userRepo UserRepository

	// In-memory cache for feature flags
	cache    map[models.FeatureKey]*featureCacheEntry
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

// NewFeatureService creates a new feature service
func NewFeatureService(repo FeatureRepository, userRepo UserRepository) *FeatureService {
	return &FeatureService{
		repo:     repo,
		userRepo: userRepo,
		cache:    make(map[models.FeatureKey]*featureCacheEntry),
		cacheTTL: time.Minute, // 1 minute cache TTL
	}
}

// IsEnabled checks if a feature is enabled for a user
// Priority: user override > tier-based > default
func (s *FeatureService) IsEnabled(ctx context.Context, userID string, key models.FeatureKey) (bool, error) {
	// Check for user-specific override first
	override, err := s.repo.GetUserFeatureOverride(ctx, userID, key)
	if err != nil {
		return false, err
	}

	if override != nil {
		// Check if override has expired
		if override.ExpiresAt != nil && override.ExpiresAt.Before(time.Now()) {
			// Override expired, delete it
			_ = s.repo.DeleteUserFeatureOverride(ctx, userID, key)
		} else {
			return override.Enabled, nil
		}
	}

	// Get feature flag
	flag, err := s.getFeatureFlag(ctx, key)
	if err != nil {
		return false, err
	}

	if flag == nil {
		return false, nil // Feature doesn't exist, default to disabled
	}

	// Check global enable first
	if !flag.GlobalEnabled {
		return false, nil
	}

	// Get user's tier
	tier, err := s.getUserTier(ctx, userID)
	if err != nil {
		return false, err
	}

	return flag.IsEnabled(tier), nil
}

// GetUserFeatures returns all features and their enabled status for a user
func (s *FeatureService) GetUserFeatures(ctx context.Context, userID string) (*models.UserFeaturesResponse, error) {
	// Get user's tier
	tier, err := s.getUserTier(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all feature flags
	flags, err := s.repo.ListFeatureFlags(ctx)
	if err != nil {
		return nil, err
	}

	// Get user's overrides
	overrides, err := s.repo.ListUserFeatureOverrides(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build override map
	overrideMap := make(map[models.FeatureKey]bool)
	for _, o := range overrides {
		if o.ExpiresAt == nil || o.ExpiresAt.After(time.Now()) {
			overrideMap[o.FeatureKey] = o.Enabled
		}
	}

	// Build response
	features := make(map[string]bool)
	for _, flag := range flags {
		// Check override first
		if enabled, hasOverride := overrideMap[flag.Key]; hasOverride {
			features[string(flag.Key)] = enabled
			continue
		}

		// Check tier-based enablement
		features[string(flag.Key)] = flag.IsEnabled(tier)
	}

	return &models.UserFeaturesResponse{
		Tier:     tier,
		Features: features,
	}, nil
}

// SetUserOverride sets a user-specific feature override
func (s *FeatureService) SetUserOverride(ctx context.Context, userID string, key models.FeatureKey, enabled bool, reason string, expiresAt *time.Time) error {
	now := time.Now()
	override := models.UserFeatureOverride{
		UserID:     userID,
		FeatureKey: key,
		Enabled:    enabled,
		Reason:     reason,
		ExpiresAt:  expiresAt,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	return s.repo.PutUserFeatureOverride(ctx, override)
}

// RemoveUserOverride removes a user's feature override
func (s *FeatureService) RemoveUserOverride(ctx context.Context, userID string, key models.FeatureKey) error {
	return s.repo.DeleteUserFeatureOverride(ctx, userID, key)
}

// InvalidateCache invalidates the feature flag cache
func (s *FeatureService) InvalidateCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache = make(map[models.FeatureKey]*featureCacheEntry)
}

// getFeatureFlag retrieves a feature flag with caching
func (s *FeatureService) getFeatureFlag(ctx context.Context, key models.FeatureKey) (*models.FeatureFlag, error) {
	// Check cache first
	s.cacheMu.RLock()
	if entry, ok := s.cache[key]; ok {
		if time.Since(entry.fetchedAt) < s.cacheTTL {
			s.cacheMu.RUnlock()
			return entry.flag, nil
		}
	}
	s.cacheMu.RUnlock()

	// Cache miss or expired, fetch from repository
	flag, err := s.repo.GetFeatureFlag(ctx, key)
	if err != nil {
		return nil, err
	}

	// Update cache
	s.cacheMu.Lock()
	s.cache[key] = &featureCacheEntry{
		flag:      flag,
		fetchedAt: time.Now(),
	}
	s.cacheMu.Unlock()

	return flag, nil
}

// getUserTier retrieves a user's subscription tier
func (s *FeatureService) getUserTier(ctx context.Context, userID string) (models.SubscriptionTier, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return models.TierFree, err
	}

	if user == nil {
		return models.TierFree, nil
	}

	// Check if user has a subscription tier set
	if user.Tier != "" {
		return user.Tier, nil
	}

	return models.TierFree, nil
}

// SeedDefaults seeds default feature flags
func (s *FeatureService) SeedDefaults(ctx context.Context) error {
	return s.repo.SeedDefaultFeatureFlags(ctx)
}
