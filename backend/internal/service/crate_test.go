package service

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock repository for crates
type mockCrateRepository struct {
	crates      map[string]*models.Crate
	cratesByKey map[string]map[string]*models.Crate // userID -> crateID -> crate
	tracks      map[string]*models.Track
}

func newMockCrateRepository() *mockCrateRepository {
	return &mockCrateRepository{
		crates:      make(map[string]*models.Crate),
		cratesByKey: make(map[string]map[string]*models.Crate),
		tracks:      make(map[string]*models.Track),
	}
}

func (m *mockCrateRepository) CreateCrate(ctx context.Context, crate models.Crate) error {
	m.crates[crate.ID] = &crate
	if _, ok := m.cratesByKey[crate.UserID]; !ok {
		m.cratesByKey[crate.UserID] = make(map[string]*models.Crate)
	}
	m.cratesByKey[crate.UserID][crate.ID] = &crate
	return nil
}

func (m *mockCrateRepository) GetCrate(ctx context.Context, userID, crateID string) (*models.Crate, error) {
	if userCrates, ok := m.cratesByKey[userID]; ok {
		if crate, ok := userCrates[crateID]; ok {
			return crate, nil
		}
	}
	return nil, nil
}

func (m *mockCrateRepository) UpdateCrate(ctx context.Context, crate models.Crate) error {
	m.crates[crate.ID] = &crate
	if _, ok := m.cratesByKey[crate.UserID]; !ok {
		m.cratesByKey[crate.UserID] = make(map[string]*models.Crate)
	}
	m.cratesByKey[crate.UserID][crate.ID] = &crate
	return nil
}

func (m *mockCrateRepository) DeleteCrate(ctx context.Context, userID, crateID string) error {
	delete(m.crates, crateID)
	if userCrates, ok := m.cratesByKey[userID]; ok {
		delete(userCrates, crateID)
	}
	return nil
}

func (m *mockCrateRepository) ListCrates(ctx context.Context, userID string, filter models.CrateFilter) ([]models.Crate, error) {
	crates := make([]models.Crate, 0)
	if userCrates, ok := m.cratesByKey[userID]; ok {
		for _, crate := range userCrates {
			crates = append(crates, *crate)
		}
	}
	return crates, nil
}

func (m *mockCrateRepository) CountUserCrates(ctx context.Context, userID string) (int, error) {
	if userCrates, ok := m.cratesByKey[userID]; ok {
		return len(userCrates), nil
	}
	return 0, nil
}

func (m *mockCrateRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	if track, ok := m.tracks[trackID]; ok {
		return track, nil
	}
	return nil, nil
}

func TestCrateService_CreateCrate(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	// Enable crates feature
	featureRepo.flags[models.FeatureCrates] = &models.FeatureFlag{
		Key:           models.FeatureCrates,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierCreator, models.TierPro},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierCreator}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	req := models.CreateCrateRequest{
		Name:        "House Music",
		Description: "My favorite house tracks",
		Color:       "#FF5733",
	}

	crate, err := svc.CreateCrate(ctx, "user-1", req)
	require.NoError(t, err)
	require.NotNil(t, crate)

	assert.NotEmpty(t, crate.ID)
	assert.Equal(t, "user-1", crate.UserID)
	assert.Equal(t, "House Music", crate.Name)
	assert.Equal(t, "My favorite house tracks", crate.Description)
	assert.Equal(t, "#FF5733", crate.Color)
	assert.Empty(t, crate.TrackIDs)
	assert.Equal(t, 0, crate.TrackCount)
	assert.Equal(t, models.CrateSortCustom, crate.SortOrder)
}

func TestCrateService_CreateCrate_FeatureDisabled(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	// Crates require Creator tier
	featureRepo.flags[models.FeatureCrates] = &models.FeatureFlag{
		Key:           models.FeatureCrates,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierCreator, models.TierPro},
	}
	// User is free tier
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	_, err := svc.CreateCrate(ctx, "user-1", models.CreateCrateRequest{Name: "Test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not enabled")
}

func TestCrateService_CreateCrate_MaxCratesLimit(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	featureRepo.flags[models.FeatureCrates] = &models.FeatureFlag{
		Key:           models.FeatureCrates,
		GlobalEnabled: true,
				EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	// Create max crates
	crateRepo.cratesByKey["user-1"] = make(map[string]*models.Crate)
	for i := 0; i < models.MaxCratesPerUser; i++ {
		crateRepo.cratesByKey["user-1"][string(rune('a'+i))] = &models.Crate{ID: string(rune('a' + i))}
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	_, err := svc.CreateCrate(ctx, "user-1", models.CreateCrateRequest{Name: "One Too Many"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum number of crates")
}

func TestCrateService_CreateSmartCrate(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	featureRepo.flags[models.FeatureCrates] = &models.FeatureFlag{
		Key:           models.FeatureCrates,
		GlobalEnabled: true,
				EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	req := models.CreateCrateRequest{
		Name:         "High Energy",
		IsSmartCrate: true,
		SmartCriteria: &models.SmartCrateCriteria{
			BPMMin: 130,
			BPMMax: 150,
			Keys:   []string{"Am", "Em"},
		},
	}

	crate, err := svc.CreateCrate(ctx, "user-1", req)
	require.NoError(t, err)
	require.NotNil(t, crate)

	assert.True(t, crate.IsSmartCrate)
	assert.NotNil(t, crate.SmartCriteria)
	assert.Equal(t, 130, crate.SmartCriteria.BPMMin)
	assert.Equal(t, 150, crate.SmartCriteria.BPMMax)
}

func TestCrateService_GetCrate(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			Name:       "Test Crate",
			TrackIDs:   []string{"track-1", "track-2"},
			TrackCount: 2,
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	crate, err := svc.GetCrate(ctx, "user-1", "crate-1")
	require.NoError(t, err)
	require.NotNil(t, crate)

	assert.Equal(t, "crate-1", crate.ID)
	assert.Equal(t, "Test Crate", crate.Name)
	assert.Len(t, crate.TrackIDs, 2)
}

func TestCrateService_GetCrate_NotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	_, err := svc.GetCrate(ctx, "user-1", "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCrateService_UpdateCrate(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:       "crate-1",
			UserID:   "user-1",
			Name:     "Old Name",
			Color:    "#000000",
			TrackIDs: []string{"track-1"},
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	newName := "New Name"
	newColor := "#FFFFFF"
	req := models.UpdateCrateRequest{
		Name:  &newName,
		Color: &newColor,
	}

	crate, err := svc.UpdateCrate(ctx, "user-1", "crate-1", req)
	require.NoError(t, err)
	require.NotNil(t, crate)

	assert.Equal(t, "New Name", crate.Name)
	assert.Equal(t, "#FFFFFF", crate.Color)
}

func TestCrateService_DeleteCrate(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:     "crate-1",
			UserID: "user-1",
		},
	}
	crateRepo.crates["crate-1"] = crateRepo.cratesByKey["user-1"]["crate-1"]

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.DeleteCrate(ctx, "user-1", "crate-1")
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetCrate(ctx, "user-1", "crate-1")
	assert.Error(t, err)
}

func TestCrateService_ListCrates(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {ID: "crate-1", UserID: "user-1", Name: "Crate 1"},
		"crate-2": {ID: "crate-2", UserID: "user-1", Name: "Crate 2"},
		"crate-3": {ID: "crate-3", UserID: "user-1", Name: "Crate 3"},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	crates, err := svc.ListCrates(ctx, "user-1", models.CrateFilter{})
	require.NoError(t, err)
	assert.Len(t, crates, 3)
}

func TestCrateService_AddTracks(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1"},
			TrackCount: 1,
		},
	}

	// Add tracks to repo
	crateRepo.tracks["track-2"] = &models.Track{ID: "track-2"}
	crateRepo.tracks["track-3"] = &models.Track{ID: "track-3"}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.AddTracks(ctx, "user-1", "crate-1", models.AddTracksToCrateRequest{
		TrackIDs: []string{"track-2", "track-3"},
	})
	require.NoError(t, err)

	// Verify tracks added
	crate := crateRepo.cratesByKey["user-1"]["crate-1"]
	assert.Len(t, crate.TrackIDs, 3)
	assert.Equal(t, 3, crate.TrackCount)
}

func TestCrateService_AddTracks_AtPosition(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1", "track-4"},
			TrackCount: 2,
		},
	}

	crateRepo.tracks["track-2"] = &models.Track{ID: "track-2"}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.AddTracks(ctx, "user-1", "crate-1", models.AddTracksToCrateRequest{
		TrackIDs: []string{"track-2"},
		Position: 1, // Insert at position 1
	})
	require.NoError(t, err)

	crate := crateRepo.cratesByKey["user-1"]["crate-1"]
	assert.Equal(t, []string{"track-1", "track-2", "track-4"}, crate.TrackIDs)
}

func TestCrateService_AddTracks_MaxLimit(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	// Create crate with nearly max tracks
	trackIDs := make([]string, models.MaxTracksPerCrate-1)
	for i := range trackIDs {
		trackIDs[i] = string(rune('a' + i%26))
	}

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   trackIDs,
			TrackCount: len(trackIDs),
		},
	}

	// Try to add more than limit
	crateRepo.tracks["new-track-1"] = &models.Track{ID: "new-track-1"}
	crateRepo.tracks["new-track-2"] = &models.Track{ID: "new-track-2"}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.AddTracks(ctx, "user-1", "crate-1", models.AddTracksToCrateRequest{
		TrackIDs: []string{"new-track-1", "new-track-2"},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceed the maximum")
}

func TestCrateService_AddTracks_Duplicates(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1"},
			TrackCount: 1,
		},
	}

	crateRepo.tracks["track-1"] = &models.Track{ID: "track-1"} // Already in crate

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	// Add duplicate track
	err := svc.AddTracks(ctx, "user-1", "crate-1", models.AddTracksToCrateRequest{
		TrackIDs: []string{"track-1"},
	})
	require.NoError(t, err)

	// Should not have duplicates
	crate := crateRepo.cratesByKey["user-1"]["crate-1"]
	assert.Len(t, crate.TrackIDs, 1)
}

func TestCrateService_AddTracks_TrackNotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:       "crate-1",
			UserID:   "user-1",
			TrackIDs: []string{},
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.AddTracks(ctx, "user-1", "crate-1", models.AddTracksToCrateRequest{
		TrackIDs: []string{"non-existent"},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCrateService_RemoveTracks(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1", "track-2", "track-3"},
			TrackCount: 3,
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.RemoveTracks(ctx, "user-1", "crate-1", models.RemoveTracksFromCrateRequest{
		TrackIDs: []string{"track-2"},
	})
	require.NoError(t, err)

	crate := crateRepo.cratesByKey["user-1"]["crate-1"]
	assert.Equal(t, []string{"track-1", "track-3"}, crate.TrackIDs)
	assert.Equal(t, 2, crate.TrackCount)
}

func TestCrateService_ReorderTracks(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1", "track-2", "track-3"},
			TrackCount: 3,
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	err := svc.ReorderTracks(ctx, "user-1", "crate-1", models.ReorderTracksRequest{
		TrackIDs: []string{"track-3", "track-1", "track-2"},
	})
	require.NoError(t, err)

	crate := crateRepo.cratesByKey["user-1"]["crate-1"]
	assert.Equal(t, []string{"track-3", "track-1", "track-2"}, crate.TrackIDs)
	assert.Equal(t, models.CrateSortCustom, crate.SortOrder)
}

func TestCrateService_ReorderTracks_MismatchedTracks(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1", "track-2"},
			TrackCount: 2,
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	// Try to reorder with different tracks
	err := svc.ReorderTracks(ctx, "user-1", "crate-1", models.ReorderTracksRequest{
		TrackIDs: []string{"track-1", "track-3"}, // track-3 not in crate
	})
	assert.Error(t, err)
}

func TestCrateService_ReorderTracks_MissingTrack(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-1", "track-2", "track-3"},
			TrackCount: 3,
		},
	}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	// Missing track-2
	err := svc.ReorderTracks(ctx, "user-1", "crate-1", models.ReorderTracksRequest{
		TrackIDs: []string{"track-1", "track-3"},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "same tracks")
}

func TestCrateService_UpdateSortOrder(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	crateRepo := newMockCrateRepository()

	crateRepo.cratesByKey["user-1"] = map[string]*models.Crate{
		"crate-1": {
			ID:         "crate-1",
			UserID:     "user-1",
			TrackIDs:   []string{"track-b", "track-a", "track-c"},
			TrackCount: 3,
			SortOrder:  models.CrateSortCustom,
		},
	}

	// Add tracks with titles for sorting
	crateRepo.tracks["track-a"] = &models.Track{ID: "track-a", Title: "Alpha"}
	crateRepo.tracks["track-b"] = &models.Track{ID: "track-b", Title: "Beta"}
	crateRepo.tracks["track-c"] = &models.Track{ID: "track-c", Title: "Charlie"}

	svc := NewCrateService(crateRepo, featureSvc)
	ctx := context.Background()

	sortByTitle := models.CrateSortTitle
	_, err := svc.UpdateCrate(ctx, "user-1", "crate-1", models.UpdateCrateRequest{
		SortOrder: &sortByTitle,
	})
	require.NoError(t, err)

	crate := crateRepo.cratesByKey["user-1"]["crate-1"]
	assert.Equal(t, models.CrateSortTitle, crate.SortOrder)
	// Note: actual sorting depends on the implementation of sortCrateTracks
}
