package service

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock repository for matching
type mockMatchingRepository struct {
	tracks map[string]*models.Track
}

func newMockMatchingRepository() *mockMatchingRepository {
	return &mockMatchingRepository{
		tracks: make(map[string]*models.Track),
	}
}

func (m *mockMatchingRepository) ListTracksForMatching(ctx context.Context, userID string) ([]models.Track, error) {
	tracks := make([]models.Track, 0)
	for _, track := range m.tracks {
		tracks = append(tracks, *track)
	}
	return tracks, nil
}

func (m *mockMatchingRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	if track, ok := m.tracks[trackID]; ok {
		return track, nil
	}
	return nil, nil
}

func TestCalculateBPMCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		sourceBPM   float64
		targetBPM   float64
		minScore    float64
		maxScore    float64
		description string
	}{
		{"exact match", 128.0, 128.0, 1.0, 1.0, "same BPM"},
		{"within 3%", 128.0, 131.0, 0.7, 1.0, "very compatible"},
		{"double time", 128.0, 64.0, 0.9, 1.0, "half time compatible"},
		{"half time", 64.0, 128.0, 0.9, 1.0, "double time compatible"},
		{"6% difference", 128.0, 136.0, 0.3, 0.7, "somewhat compatible"},
		{"incompatible", 128.0, 150.0, 0.0, 0.3, "too different"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, diff := calculateBPMCompatibility(tt.sourceBPM, tt.targetBPM)
			assert.GreaterOrEqual(t, score, tt.minScore, "score should be at least %f for %s", tt.minScore, tt.description)
			assert.LessOrEqual(t, score, tt.maxScore, "score should be at most %f for %s", tt.maxScore, tt.description)
			assert.GreaterOrEqual(t, diff, 0.0, "diff should be non-negative")
		})
	}
}

func TestCalculateKeyCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		sourceKey string
		targetKey string
		expected  float64
		relation  string
	}{
		{"same key", "Am", "Am", 1.0, "same"},
		{"same key alt notation", "A minor", "Am", 1.0, "same"},
		{"relative major/minor", "Am", "C", 0.9, "relative"},
		{"neighbor key", "Am", "Em", 0.85, "neighbor"},
		{"neighbor key reverse", "Am", "Dm", 0.85, "neighbor"},
		{"two steps away", "Am", "Bm", 0.6, "compatible"},
		{"energy shift", "Am", "G", 0.7, "energy_shift"},
		{"tritone", "Am", "Ebm", 0.4, "tritone"},
		{"incompatible", "Am", "F#m", 0.3, "incompatible"},
		{"unknown key", "Am", "X#m", 0.5, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, relation := calculateKeyCompatibility(tt.sourceKey, tt.targetKey)
			assert.InDelta(t, tt.expected, score, 0.01, "score for %s -> %s", tt.sourceKey, tt.targetKey)
			assert.Equal(t, tt.relation, relation)
		})
	}
}

func TestCamelotDistance(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{1, 1, 0},   // Same position
		{1, 2, 1},   // Adjacent
		{12, 1, 1},  // Wrap around
		{1, 7, 6},   // Opposite
		{6, 12, 6},  // Opposite
		{1, 4, 3},   // 3 steps
		{10, 3, 5},  // 5 steps
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := camelotDistance(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "distance between %d and %d", tt.a, tt.b)
		})
	}
}

func TestMatchingService_FindCompatibleTracks(t *testing.T) {
	// Setup feature service mock
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	// Setup matching repository
	matchingRepo := newMockMatchingRepository()

	// Enable both BPM and key matching features
	featureRepo.flags[models.FeatureBPMMatching] = &models.FeatureFlag{
		Key:           models.FeatureBPMMatching,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree, models.TierCreator, models.TierPro},
	}
	featureRepo.flags[models.FeatureKeyMatching] = &models.FeatureFlag{
		Key:           models.FeatureKeyMatching,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree, models.TierCreator, models.TierPro},
	}

	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierPro}

	// Add tracks
	matchingRepo.tracks["track-source"] = &models.Track{
		ID:         "track-source",
		Title:      "Source Track",
		BPM:        128,
		MusicalKey: "Am",
	}
	matchingRepo.tracks["track-perfect"] = &models.Track{
		ID:         "track-perfect",
		Title:      "Perfect Match",
		BPM:        128,
		MusicalKey: "Am",
	}
	matchingRepo.tracks["track-bpm-match"] = &models.Track{
		ID:         "track-bpm-match",
		Title:      "BPM Match",
		BPM:        130,
		MusicalKey: "F#m", // Incompatible key
	}
	matchingRepo.tracks["track-key-match"] = &models.Track{
		ID:         "track-key-match",
		Title:      "Key Match",
		BPM:        90, // Different BPM
		MusicalKey: "C", // Relative major
	}
	matchingRepo.tracks["track-incompatible"] = &models.Track{
		ID:         "track-incompatible",
		Title:      "Incompatible",
		BPM:        200,
		MusicalKey: "F#",
	}

	svc := NewMatchingService(matchingRepo, featureSvc)
	ctx := context.Background()

	results, err := svc.FindCompatibleTracks(ctx, "user-1", "track-source", 10)
	require.NoError(t, err)
	require.NotEmpty(t, results)

	// Perfect match should be first
	assert.Equal(t, "track-perfect", results[0].Track.ID)
	assert.Equal(t, 1.0, results[0].BPMCompatibility)
	assert.Equal(t, 1.0, results[0].KeyCompatibility)
	assert.Equal(t, "same", results[0].KeyRelation)
}

func TestMatchingService_FindCompatibleTracks_NoFeatures(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	matchingRepo := newMockMatchingRepository()
	matchingRepo.tracks["track-1"] = &models.Track{ID: "track-1", BPM: 128, MusicalKey: "Am"}

	// No features enabled
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	svc := NewMatchingService(matchingRepo, featureSvc)
	ctx := context.Background()

	results, err := svc.FindCompatibleTracks(ctx, "user-1", "track-1", 10)
	require.NoError(t, err)
	assert.Nil(t, results, "should return nil when no features enabled")
}

func TestMatchingService_FindCompatibleTracks_TrackNotFound(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	matchingRepo := newMockMatchingRepository()

	// Enable feature
	featureRepo.flags[models.FeatureBPMMatching] = &models.FeatureFlag{
		Key:           models.FeatureBPMMatching,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	svc := NewMatchingService(matchingRepo, featureSvc)
	ctx := context.Background()

	results, err := svc.FindCompatibleTracks(ctx, "user-1", "non-existent", 10)
	require.NoError(t, err)
	assert.Nil(t, results)
}

func TestMatchingService_FindBPMCompatible(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	matchingRepo := newMockMatchingRepository()

	featureRepo.flags[models.FeatureBPMMatching] = &models.FeatureFlag{
		Key:           models.FeatureBPMMatching,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	matchingRepo.tracks["track-126"] = &models.Track{ID: "track-126", BPM: 126}
	matchingRepo.tracks["track-128"] = &models.Track{ID: "track-128", BPM: 128}
	matchingRepo.tracks["track-130"] = &models.Track{ID: "track-130", BPM: 130}
	matchingRepo.tracks["track-160"] = &models.Track{ID: "track-160", BPM: 160}

	svc := NewMatchingService(matchingRepo, featureSvc)
	ctx := context.Background()

	results, err := svc.FindBPMCompatible(ctx, "user-1", 128.0, 5.0, 10)
	require.NoError(t, err)

	// Should find tracks around 128 BPM
	assert.GreaterOrEqual(t, len(results), 3, "should find compatible tracks")

	// Verify sorted by score descending
	for i := 1; i < len(results); i++ {
		assert.GreaterOrEqual(t, results[i-1].OverallScore, results[i].OverallScore,
			"results should be sorted by score descending")
	}
}

func TestMatchingService_FindKeyCompatible(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	matchingRepo := newMockMatchingRepository()

	featureRepo.flags[models.FeatureKeyMatching] = &models.FeatureFlag{
		Key:           models.FeatureKeyMatching,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	matchingRepo.tracks["track-am"] = &models.Track{ID: "track-am", MusicalKey: "Am"}
	matchingRepo.tracks["track-c"] = &models.Track{ID: "track-c", MusicalKey: "C"}  // Relative major
	matchingRepo.tracks["track-em"] = &models.Track{ID: "track-em", MusicalKey: "Em"} // Neighbor
	matchingRepo.tracks["track-fsharp"] = &models.Track{ID: "track-fsharp", MusicalKey: "F#"}

	svc := NewMatchingService(matchingRepo, featureSvc)
	ctx := context.Background()

	results, err := svc.FindKeyCompatible(ctx, "user-1", "Am", 10)
	require.NoError(t, err)

	// Should find tracks compatible with Am
	assert.GreaterOrEqual(t, len(results), 2, "should find compatible tracks")

	// Find Am track (perfect match)
	var foundAm bool
	for _, r := range results {
		if r.Track.ID == "track-am" {
			foundAm = true
			assert.Equal(t, 1.0, r.KeyCompatibility)
			assert.Equal(t, "same", r.KeyRelation)
		}
	}
	assert.True(t, foundAm, "should find Am track")
}

func TestMatchingService_LimitResults(t *testing.T) {
	featureRepo := newMockFeatureRepository()
	userRepo := newMockUserRepository()
	featureSvc := NewFeatureService(featureRepo, userRepo)

	matchingRepo := newMockMatchingRepository()

	featureRepo.flags[models.FeatureBPMMatching] = &models.FeatureFlag{
		Key:           models.FeatureBPMMatching,
		GlobalEnabled: true,
		EnabledTiers:  []models.SubscriptionTier{models.TierFree},
	}
	userRepo.users["user-1"] = &models.User{ID: "user-1", Tier: models.TierFree}

	// Add many tracks
	for i := 0; i < 20; i++ {
		matchingRepo.tracks[string(rune('a'+i))] = &models.Track{
			ID:  string(rune('a' + i)),
			BPM: 120 + i,
		}
	}

	svc := NewMatchingService(matchingRepo, featureSvc)
	ctx := context.Background()

	results, err := svc.FindBPMCompatible(ctx, "user-1", 128.0, 10.0, 5)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(results), 5, "should respect limit")
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Am", "Am"},
		{"A min", "A minor"},    // "min" → "minor" (no space added)
		{"C maj", "C major"},    // "maj" → "major" (no space added)
		{" Am ", "Am"},          // trimmed
		{"A#min", "A#minor"},    // "min" → "minor" (no space in input)
		{"A minor", "A minor"},  // already has "minor", unchanged
		{"C major", "C major"},  // already has "major", unchanged
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
