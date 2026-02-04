package service

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock repository for hot cues
type mockHotCueRepository struct {
	tracks map[string]*models.Track
}

func newMockHotCueRepository() *mockHotCueRepository {
	return &mockHotCueRepository{
		tracks: make(map[string]*models.Track),
	}
}

func (m *mockHotCueRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	if track, ok := m.tracks[trackID]; ok {
		return track, nil
	}
	return nil, nil
}

func (m *mockHotCueRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	m.tracks[track.ID] = &track
	return nil
}

func TestHotCueService_SetHotCue(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	// Enable hot cues feature
	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierCreator, models.TierPro},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierCreator}

	// Add a track
	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300, // 5 minutes
		HotCues:  make(map[int]*models.HotCue),
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	// Set a hot cue
	req := models.SetHotCueRequest{
		Position: 60.5,
		Label:    "Drop",
		Color:    "#FF0000",
	}

	cue, err := svc.SetHotCue(ctx, "user-1", "track-1", 1, req)
	require.NoError(t, err)
	require.NotNil(t, cue)

	assert.Equal(t, 1, cue.Slot)
	assert.Equal(t, 60.5, cue.Position)
	assert.Equal(t, "Drop", cue.Label)
	assert.Equal(t, models.HotCueColor("#FF0000"), cue.Color)
	assert.False(t, cue.CreatedAt.IsZero())
	assert.False(t, cue.UpdatedAt.IsZero())
}

func TestHotCueService_SetHotCue_DefaultColor(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues:  make(map[int]*models.HotCue),
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	// Set hot cue without color
	req := models.SetHotCueRequest{
		Position: 30.0,
	}

	cue, err := svc.SetHotCue(ctx, "user-1", "track-1", 3, req)
	require.NoError(t, err)
	require.NotNil(t, cue)

	// Should get default color for slot 3
	assert.NotEmpty(t, cue.Color, "should have default color")
}

func TestHotCueService_SetHotCue_InvalidSlot(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	tests := []struct {
		name string
		slot int
	}{
		{"slot 0", 0},
		{"slot 9", 9},
		{"negative slot", -1},
		{"slot 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.SetHotCue(ctx, "user-1", "track-1", tt.slot, models.SetHotCueRequest{Position: 10.0})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid slot")
		})
	}
}

func TestHotCueService_SetHotCue_FeatureDisabled(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	// Hot cues feature requires Pro
	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierPro},
	}
	// User is only Creator tier
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierCreator}

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	_, err := svc.SetHotCue(ctx, "user-1", "track-1", 1, models.SetHotCueRequest{Position: 10.0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not enabled")
}

func TestHotCueService_SetHotCue_TrackNotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	_, err := svc.SetHotCue(ctx, "user-1", "non-existent", 1, models.SetHotCueRequest{Position: 10.0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "track not found")
}

func TestHotCueService_SetHotCue_PositionExceedsDuration(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 180, // 3 minutes
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	// Position exceeds duration
	_, err := svc.SetHotCue(ctx, "user-1", "track-1", 1, models.SetHotCueRequest{Position: 200.0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds track duration")
}

func TestHotCueService_SetHotCue_UpdateExisting(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues:  make(map[int]*models.HotCue),
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	// Set initial hot cue
	cue1, err := svc.SetHotCue(ctx, "user-1", "track-1", 1, models.SetHotCueRequest{
		Position: 30.0,
		Label:    "Intro",
	})
	require.NoError(t, err)
	originalCreatedAt := cue1.CreatedAt

	// Update the same slot
	cue2, err := svc.SetHotCue(ctx, "user-1", "track-1", 1, models.SetHotCueRequest{
		Position: 45.0,
		Label:    "Verse",
	})
	require.NoError(t, err)

	// Should preserve original creation time
	assert.Equal(t, originalCreatedAt, cue2.CreatedAt, "created_at should be preserved on update")
	assert.Equal(t, 45.0, cue2.Position)
	assert.Equal(t, "Verse", cue2.Label)
}

func TestHotCueService_DeleteHotCue(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues: map[int]*models.HotCue{
			1: {Slot: 1, Position: 30.0},
			2: {Slot: 2, Position: 60.0},
		},
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	err := svc.DeleteHotCue(ctx, "user-1", "track-1", 1)
	require.NoError(t, err)

	// Verify deletion
	track := hotCueRepo.tracks["track-1"]
	assert.Nil(t, track.HotCues[1], "slot 1 should be deleted")
	assert.NotNil(t, track.HotCues[2], "slot 2 should remain")
}

func TestHotCueService_DeleteHotCue_InvalidSlot(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	err := svc.DeleteHotCue(ctx, "user-1", "track-1", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid slot")
}

func TestHotCueService_DeleteHotCue_NoHotCues(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues:  nil, // No hot cues
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	// Should not error when deleting from empty
	err := svc.DeleteHotCue(ctx, "user-1", "track-1", 1)
	assert.NoError(t, err)
}

func TestHotCueService_GetHotCues(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues: map[int]*models.HotCue{
			1: {Slot: 1, Position: 30.0, Label: "Intro", Color: "#FF0000"},
			3: {Slot: 3, Position: 90.0, Label: "Drop", Color: "#00FF00"},
			5: {Slot: 5, Position: 150.0, Color: "#0000FF"},
		},
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	response, err := svc.GetHotCues(ctx, "user-1", "track-1")
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, "track-1", response.TrackID)
	assert.Equal(t, models.MaxHotCuesPerTrack, response.MaxSlots)
	assert.Len(t, response.HotCues, 3)

	// Verify hot cue data
	for _, cue := range response.HotCues {
		switch cue.Slot {
		case 1:
			assert.Equal(t, 30.0, cue.Position)
			assert.Equal(t, "Intro", cue.Label)
		case 3:
			assert.Equal(t, 90.0, cue.Position)
			assert.Equal(t, "Drop", cue.Label)
		case 5:
			assert.Equal(t, 150.0, cue.Position)
		}
	}
}

func TestHotCueService_GetHotCues_TrackNotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	_, err := svc.GetHotCues(ctx, "user-1", "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "track not found")
}

func TestHotCueService_GetHotCues_EmptyHotCues(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues:  nil,
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	response, err := svc.GetHotCues(ctx, "user-1", "track-1")
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, "track-1", response.TrackID)
	assert.Empty(t, response.HotCues)
	assert.Equal(t, models.MaxHotCuesPerTrack, response.MaxSlots)
}

func TestHotCueService_ClearAllHotCues(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 300,
		HotCues: map[int]*models.HotCue{
			1: {Slot: 1, Position: 30.0},
			2: {Slot: 2, Position: 60.0},
			3: {Slot: 3, Position: 90.0},
		},
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	err := svc.ClearAllHotCues(ctx, "user-1", "track-1")
	require.NoError(t, err)

	// Verify all cleared
	track := hotCueRepo.tracks["track-1"]
	assert.Empty(t, track.HotCues)
}

func TestHotCueService_ClearAllHotCues_TrackNotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	err := svc.ClearAllHotCues(ctx, "user-1", "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "track not found")
}

func TestHotCueService_AllSlots(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	hotCueRepo := newMockHotCueRepository()

	featureRepo.flags[models.FeatureHotCues] = &models.FeatureFlag{
		Key:           models.FeatureHotCues,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	hotCueRepo.tracks["track-1"] = &models.Track{
		ID:       "track-1",
		Duration: 600, // 10 minutes
		HotCues:  make(map[int]*models.HotCue),
	}

	svc := NewHotCueService(hotCueRepo, featureSvc)
	ctx := context.Background()

	// Set all 8 slots
	for slot := 1; slot <= 8; slot++ {
		cue, err := svc.SetHotCue(ctx, "user-1", "track-1", slot, models.SetHotCueRequest{
			Position: float64(slot * 60),
			Label:    "Cue " + string(rune('A'+slot-1)),
		})
		require.NoError(t, err)
		require.NotNil(t, cue)
	}

	// Verify all slots set
	response, err := svc.GetHotCues(ctx, "user-1", "track-1")
	require.NoError(t, err)
	assert.Len(t, response.HotCues, 8)
}
