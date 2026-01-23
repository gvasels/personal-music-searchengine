package service

import (
	"context"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock repository for feature flags
type mockFeatureRepository struct {
	flags     map[models.FeatureKey]*models.FeatureFlag
	overrides map[string]map[models.FeatureKey]*models.UserFeatureOverride
}

func newMockFeatureRepository() *mockFeatureRepository {
	return &mockFeatureRepository{
		flags:     make(map[models.FeatureKey]*models.FeatureFlag),
		overrides: make(map[string]map[models.FeatureKey]*models.UserFeatureOverride),
	}
}

func (m *mockFeatureRepository) GetFeatureFlag(ctx context.Context, key models.FeatureKey) (*models.FeatureFlag, error) {
	if flag, ok := m.flags[key]; ok {
		return flag, nil
	}
	return nil, nil
}

func (m *mockFeatureRepository) ListFeatureFlags(ctx context.Context) ([]models.FeatureFlag, error) {
	flags := make([]models.FeatureFlag, 0, len(m.flags))
	for _, flag := range m.flags {
		flags = append(flags, *flag)
	}
	return flags, nil
}

func (m *mockFeatureRepository) PutFeatureFlag(ctx context.Context, flag models.FeatureFlag) error {
	m.flags[flag.Key] = &flag
	return nil
}

func (m *mockFeatureRepository) DeleteFeatureFlag(ctx context.Context, key models.FeatureKey) error {
	delete(m.flags, key)
	return nil
}

func (m *mockFeatureRepository) GetUserFeatureOverride(ctx context.Context, userID string, key models.FeatureKey) (*models.UserFeatureOverride, error) {
	if userOverrides, ok := m.overrides[userID]; ok {
		if override, ok := userOverrides[key]; ok {
			return override, nil
		}
	}
	return nil, nil
}

func (m *mockFeatureRepository) ListUserFeatureOverrides(ctx context.Context, userID string) ([]models.UserFeatureOverride, error) {
	overrides := make([]models.UserFeatureOverride, 0)
	if userOverrides, ok := m.overrides[userID]; ok {
		for _, override := range userOverrides {
			overrides = append(overrides, *override)
		}
	}
	return overrides, nil
}

func (m *mockFeatureRepository) PutUserFeatureOverride(ctx context.Context, override models.UserFeatureOverride) error {
	if _, ok := m.overrides[override.UserID]; !ok {
		m.overrides[override.UserID] = make(map[models.FeatureKey]*models.UserFeatureOverride)
	}
	m.overrides[override.UserID][override.FeatureKey] = &override
	return nil
}

func (m *mockFeatureRepository) DeleteUserFeatureOverride(ctx context.Context, userID string, key models.FeatureKey) error {
	if userOverrides, ok := m.overrides[userID]; ok {
		delete(userOverrides, key)
	}
	return nil
}

func (m *mockFeatureRepository) SeedDefaultFeatureFlags(ctx context.Context) error {
	return nil
}

// Mock user repository
type mockUserRepository struct {
	users map[string]*models.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if user, ok := m.users[userID]; ok {
		return user, nil
	}
	return nil, nil
}

// Test helper to create a feature flag
func createTestFeatureFlag(key models.FeatureKey, globalEnabled bool, enabledTiers []models.SubscriptionTier) models.FeatureFlag {
	return models.FeatureFlag{
		Key:           key,
		Name:          string(key),
		Description:   "Test feature",
		GlobalEnabled: globalEnabled,
		EnabledTiers:  enabledTiers,
	}
}

func TestFeatureService_IsEnabled_GlobalDisabled(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	// Create a globally disabled feature
	flag := createTestFeatureFlag(models.FeatureDJModule, false, []models.SubscriptionTier{models.TierFree, models.TierCreator, models.TierPro})
	featureRepo.flags[models.FeatureDJModule] = &flag

	// Add user
	userRepo.users[userID] = &models.User{ID: userID, Tier: models.TierPro}

	enabled, err := service.IsEnabled(ctx, userID, models.FeatureDJModule)
	require.NoError(t, err)
	assert.False(t, enabled, "globally disabled feature should be disabled")
}

func TestFeatureService_IsEnabled_TierBased(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()

	// Create a feature that requires Creator tier
	flag := models.FeatureFlag{
		Key:           models.FeatureCrates,
		Name:          "Crates",
		Description:   "DJ crates feature",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierCreator, models.TierPro},
	}
	featureRepo.flags[models.FeatureCrates] = &flag

	tests := []struct {
		name     string
		userID   string
		userTier models.SubscriptionTier
		expected bool
	}{
		{"free user", "user-free", models.TierFree, false},
		{"creator user", "user-creator", models.TierCreator, true},
		{"pro user", "user-pro", models.TierPro, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo.users[tt.userID] = &models.User{ID: tt.userID, Tier: tt.userTier}

			enabled, err := service.IsEnabled(ctx, tt.userID, models.FeatureCrates)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, enabled)
		})
	}
}

func TestFeatureService_IsEnabled_UserOverride(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	// Create a feature that requires Pro tier
	flag := models.FeatureFlag{
		Key:           models.FeatureHotCues,
		Name:          "Hot Cues",
		Description:   "Hot cues feature",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierPro},
	}
	featureRepo.flags[models.FeatureHotCues] = &flag

	// User is free tier (normally wouldn't have access)
	userRepo.users[userID] = &models.User{ID: userID, Tier: models.TierFree}

	// Without override, should be disabled
	enabled, err := service.IsEnabled(ctx, userID, models.FeatureHotCues)
	require.NoError(t, err)
	assert.False(t, enabled)

	// Add override to enable
	featureRepo.overrides[userID] = map[models.FeatureKey]*models.UserFeatureOverride{
		models.FeatureHotCues: {
			UserID:     userID,
			FeatureKey: models.FeatureHotCues,
			Enabled:    true,
			Reason:     "Beta tester",
		},
	}

	// With override, should be enabled
	enabled, err = service.IsEnabled(ctx, userID, models.FeatureHotCues)
	require.NoError(t, err)
	assert.True(t, enabled)
}

func TestFeatureService_IsEnabled_ExpiredOverride(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	// Create a feature that requires Pro tier
	flag := models.FeatureFlag{
		Key:           models.FeatureHotCues,
		Name:          "Hot Cues",
		Description:   "Hot cues feature",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierPro},
	}
	featureRepo.flags[models.FeatureHotCues] = &flag

	// User is free tier
	userRepo.users[userID] = &models.User{ID: userID, Tier: models.TierFree}

	// Add expired override
	expiredTime := time.Now().Add(-1 * time.Hour)
	featureRepo.overrides[userID] = map[models.FeatureKey]*models.UserFeatureOverride{
		models.FeatureHotCues: {
			UserID:     userID,
			FeatureKey: models.FeatureHotCues,
			Enabled:    true,
			Reason:     "Trial expired",
			ExpiresAt:  &expiredTime,
		},
	}

	// With expired override, should fall back to tier check
	enabled, err := service.IsEnabled(ctx, userID, models.FeatureHotCues)
	require.NoError(t, err)
	assert.False(t, enabled)
}

func TestFeatureService_IsEnabled_FeatureNotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	userRepo.users[userID] = &models.User{ID: userID, Tier: models.TierPro}

	// Non-existent feature should return false
	enabled, err := service.IsEnabled(ctx, userID, "non-existent-feature")
	require.NoError(t, err)
	assert.False(t, enabled)
}

func TestFeatureService_GetUserFeatures(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	// Create multiple features
	featureRepo.flags[models.FeatureDJModule] = &models.FeatureFlag{
		Key:           models.FeatureDJModule,
		Name:          "DJ Module",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree, models.TierCreator, models.TierPro},
	}
	featureRepo.flags[models.FeatureCrates] = &models.FeatureFlag{
		Key:           models.FeatureCrates,
		Name:          "Crates",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierCreator, models.TierPro},
	}
	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		Name:          "Hot Cues",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierPro},
	}

	// Creator user
	userRepo.users[userID] = &models.User{ID: userID, Tier: models.TierCreator}

	response, err := service.GetUserFeatures(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, models.TierCreator, response.Tier)
	assert.True(t, response.Features[string(models.FeatureDJModule)])
	assert.True(t, response.Features[string(models.FeatureCrates)])
	assert.False(t, response.Features[string(models.FeatureHotCues)])
}

func TestFeatureService_SetUserOverride(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	err := service.SetUserOverride(ctx, userID, models.FeatureHotCues, true, "Beta tester", nil)
	require.NoError(t, err)

	// Verify override was set
	override, err := featureRepo.GetUserFeatureOverride(ctx, userID, models.FeatureHotCues)
	require.NoError(t, err)
	require.NotNil(t, override)
	assert.Equal(t, userID, override.UserID)
	assert.Equal(t, models.FeatureHotCues, override.FeatureKey)
	assert.True(t, override.Enabled)
	assert.Equal(t, "Beta tester", override.Reason)
}

func TestFeatureService_RemoveUserOverride(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	// Set override
	featureRepo.overrides[userID] = map[models.FeatureKey]*models.UserFeatureOverride{
		models.FeatureHotCues: {
			UserID:     userID,
			FeatureKey: models.FeatureHotCues,
			Enabled:    true,
		},
	}

	err := service.RemoveUserOverride(ctx, userID, models.FeatureHotCues)
	require.NoError(t, err)

	// Verify override was removed
	override, err := featureRepo.GetUserFeatureOverride(ctx, userID, models.FeatureHotCues)
	require.NoError(t, err)
	assert.Nil(t, override)
}

func TestFeatureService_Cache(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "user-1"

	// Create a feature
	flag := createTestFeatureFlag(models.FeatureDJModule, true, []models.SubscriptionTier{models.TierFree, models.TierCreator, models.TierPro})
	featureRepo.flags[models.FeatureDJModule] = &flag

	userRepo.users[userID] = &models.User{ID: userID, Tier: models.TierFree}

	// First call - should cache
	enabled1, err := service.IsEnabled(ctx, userID, models.FeatureDJModule)
	require.NoError(t, err)
	assert.True(t, enabled1)

	// Modify flag in repo (simulating external change)
	// Create a new flag object to properly test cache (not modifying the cached pointer)
	disabledFlag := createTestFeatureFlag(models.FeatureDJModule, false, []models.SubscriptionTier{models.TierFree, models.TierCreator, models.TierPro})
	featureRepo.flags[models.FeatureDJModule] = &disabledFlag

	// Second call - should use cache
	enabled2, err := service.IsEnabled(ctx, userID, models.FeatureDJModule)
	require.NoError(t, err)
	assert.True(t, enabled2, "should use cached value")

	// Invalidate cache
	service.InvalidateCache()

	// Third call - should fetch fresh
	enabled3, err := service.IsEnabled(ctx, userID, models.FeatureDJModule)
	require.NoError(t, err)
	assert.False(t, enabled3, "should use fresh value after cache invalidation")
}

func TestFeatureService_DefaultTierForNewUser(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	service := NewFeatureService(featureRepo, userRepo)

	ctx := context.Background()
	userID := "new-user"

	// Create feature requiring creator tier
	flag := models.FeatureFlag{
		Key:           models.FeatureCrates,
		Name:          "Crates",
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierCreator, models.TierPro},
	}
	featureRepo.flags[models.FeatureCrates] = &flag

	// User doesn't exist in repo - should default to Free tier
	enabled, err := service.IsEnabled(ctx, userID, models.FeatureCrates)
	require.NoError(t, err)
	assert.False(t, enabled, "new user should default to free tier")
}
