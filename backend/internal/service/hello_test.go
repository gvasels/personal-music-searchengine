package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =============================================================================
// HelloService Tests (TDD Red Phase)
// =============================================================================

// HelloTrack is the expected return type from HelloService methods.
// This struct must be defined in the implementation (hello.go) to make tests pass.
// Fields: ID, Title, Artist, Album, Genre (string); Year, Duration (int)

// MockHelloRepo provides a mock for the data access interface that HelloService needs.
// HelloService will query DynamoDB for seed tracks (PK=USER#seed-user, SK begins_with TRACK#)
// and return them as []HelloTrack.
type MockHelloRepo struct {
	mock.Mock
}

// GetSeedTracks returns all tracks for the seed user.
// The implementation will call this to get tracks to search/filter over.
func (m *MockHelloRepo) GetSeedTracks(ctx context.Context) ([]HelloTrack, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]HelloTrack), args.Error(1)
}

// seedTracks returns a realistic set of 20 seed tracks for test fixtures.
func seedTracks() []HelloTrack {
	return []HelloTrack{
		{ID: "t1", Title: "Aurora Borealis", Artist: "Midnight Echo", Album: "Northern Lights", Genre: "Electronic", Year: 2023, Duration: 245},
		{ID: "t2", Title: "Velvet Dreams", Artist: "Luna Wave", Album: "Dreamscape", Genre: "Jazz", Year: 2022, Duration: 312},
		{ID: "t3", Title: "Crimson Tide", Artist: "Storm Riders", Album: "Ocean Storm", Genre: "Rock", Year: 2021, Duration: 198},
		{ID: "t4", Title: "Silver Lining", Artist: "Midnight Echo", Album: "Northern Lights", Genre: "Electronic", Year: 2023, Duration: 267},
		{ID: "t5", Title: "Neon Pulse", Artist: "Cyber Drift", Album: "Digital Dawn", Genre: "Synthwave", Year: 2024, Duration: 220},
		{ID: "t6", Title: "Whisper in the Wind", Artist: "Folk Tales", Album: "Countryside", Genre: "Folk", Year: 2020, Duration: 185},
		{ID: "t7", Title: "Jazz Nocturne", Artist: "Blue Note Trio", Album: "Late Night Sessions", Genre: "Jazz", Year: 2019, Duration: 340},
		{ID: "t8", Title: "Thunderstrike", Artist: "Storm Riders", Album: "Ocean Storm", Genre: "Rock", Year: 2021, Duration: 210},
		{ID: "t9", Title: "Celestial Aurora", Artist: "Stargazer", Album: "Cosmos", Genre: "Ambient", Year: 2023, Duration: 420},
		{ID: "t10", Title: "Midnight Serenade", Artist: "Luna Wave", Album: "Dreamscape", Genre: "Jazz", Year: 2022, Duration: 290},
		{ID: "t11", Title: "Pixel Storm", Artist: "Cyber Drift", Album: "Digital Dawn", Genre: "Synthwave", Year: 2024, Duration: 195},
		{ID: "t12", Title: "Desert Mirage", Artist: "Sand Nomads", Album: "Sahara", Genre: "World", Year: 2022, Duration: 275},
		{ID: "t13", Title: "Frozen Lake", Artist: "Arctic Sound", Album: "Polar", Genre: "Ambient", Year: 2023, Duration: 360},
		{ID: "t14", Title: "Funky Grooves", Artist: "Bass Kitchen", Album: "Cookout", Genre: "Funk", Year: 2021, Duration: 230},
		{ID: "t15", Title: "Rainy Day Blues", Artist: "Blue Note Trio", Album: "Late Night Sessions", Genre: "Jazz", Year: 2019, Duration: 305},
		{ID: "t16", Title: "Electric Sunrise", Artist: "Neon Flux", Album: "Voltage", Genre: "Electronic", Year: 2024, Duration: 215},
		{ID: "t17", Title: "Mountain Echo", Artist: "Folk Tales", Album: "Countryside", Genre: "Folk", Year: 2020, Duration: 200},
		{ID: "t18", Title: "Deep Blue Ocean", Artist: "Aqua Sound", Album: "Depths", Genre: "Ambient", Year: 2022, Duration: 380},
		{ID: "t19", Title: "City Lights", Artist: "Urban Jazz Collective", Album: "Metropolitan", Genre: "Jazz", Year: 2023, Duration: 255},
		{ID: "t20", Title: "Solar Flare", Artist: "Stargazer", Album: "Cosmos", Genre: "Electronic", Year: 2023, Duration: 235},
	}
}

// =============================================================================
// Search Tests
// =============================================================================

func TestHelloService_Search_MatchesTitle(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Search(ctx, "aurora")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// Should match "Aurora Borealis" (t1) and "Celestial Aurora" (t9)
	for _, track := range results {
		assert.Contains(t, lower(track.Title+track.Artist+track.Album+track.Genre), "aurora",
			"Every result should contain 'aurora' in at least one field")
	}
	assert.GreaterOrEqual(t, len(results), 2, "Should match at least 'Aurora Borealis' and 'Celestial Aurora'")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Search_MatchesArtist(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Search(ctx, "midnight")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// Should match artist "Midnight Echo" tracks (t1, t4) and title "Midnight Serenade" (t10)
	foundArtistMatch := false
	for _, track := range results {
		if lower(track.Artist) == "midnight echo" {
			foundArtistMatch = true
		}
	}
	assert.True(t, foundArtistMatch, "Should find tracks by 'Midnight Echo' artist")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Search_MatchesAlbum(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Search(ctx, "dreamscape")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// Should match album "Dreamscape" tracks (t2, t10)
	for _, track := range results {
		assert.Equal(t, "Dreamscape", track.Album,
			"Every result should be from the Dreamscape album")
	}
	assert.Equal(t, 2, len(results), "Should find exactly 2 tracks from Dreamscape album")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Search_MatchesGenre(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Search(ctx, "jazz")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// Should match genre "Jazz" tracks (t2, t7, t10, t15, t19)
	for _, track := range results {
		assert.Equal(t, "Jazz", track.Genre,
			"Every result should have Jazz genre")
	}
	assert.Equal(t, 5, len(results), "Should find exactly 5 Jazz tracks")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Search_CaseInsensitive(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Search(ctx, "JAZZ")

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	// Should match genre "Jazz" even though query is "JAZZ"
	for _, track := range results {
		assert.Equal(t, "Jazz", track.Genre,
			"Case-insensitive search should match Jazz genre")
	}
	assert.Equal(t, 5, len(results), "Should find exactly 5 Jazz tracks with uppercase query")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Search_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	// Should NOT call GetSeedTracks for empty query
	results, err := svc.Search(ctx, "")

	assert.NoError(t, err)
	assert.NotNil(t, results, "Empty query should return empty slice, not nil")
	assert.Empty(t, results, "Empty query should return no results")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Search_NoMatches(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Search(ctx, "nonexistent")

	assert.NoError(t, err)
	assert.NotNil(t, results, "No matches should return empty slice, not nil")
	assert.Empty(t, results, "Query for 'nonexistent' should return no results")
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Featured Tests
// =============================================================================

func TestHelloService_Featured_ReturnsAll(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Featured(ctx, 20)

	assert.NoError(t, err)
	assert.Len(t, results, 20, "Featured(20) should return all 20 seed tracks")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Featured_WithLimit(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Featured(ctx, 5)

	assert.NoError(t, err)
	assert.Len(t, results, 5, "Featured(5) should return exactly 5 tracks")
	mockRepo.AssertExpectations(t)
}

func TestHelloService_Featured_DefaultLimit(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockHelloRepo)
	svc := NewHelloService(mockRepo)

	mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

	results, err := svc.Featured(ctx, 0)

	assert.NoError(t, err)
	assert.Len(t, results, 20, "Featured(0) should treat 0 as default and return all tracks")
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Helper
// =============================================================================

// lower is a test helper for case-insensitive comparison.
func lower(s string) string {
	return strings.ToLower(s)
}
