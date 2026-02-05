package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHelloRepository implements HelloRepository for testing
type MockHelloRepository struct {
	tracks []HelloTrack
	err    error
}

func (m *MockHelloRepository) GetTracksByUser(ctx context.Context, userID string) ([]HelloTrack, error) {
	return m.tracks, m.err
}

// sampleTracks returns a set of test tracks for use across tests
func sampleTracks() []HelloTrack {
	return []HelloTrack{
		{
			ID:       "track-1",
			Title:    "Neon Lights",
			Artist:   "Kraftwerk",
			Album:    "The Man-Machine",
			Genre:    "Electronic",
			Year:     1978,
			Duration: 540,
		},
		{
			ID:       "track-2",
			Title:    "Blue Monday",
			Artist:   "New Order",
			Album:    "Power, Corruption & Lies",
			Genre:    "Synth-Pop",
			Year:     1983,
			Duration: 440,
		},
		{
			ID:       "track-3",
			Title:    "Enjoy the Silence",
			Artist:   "Depeche Mode",
			Album:    "Violator",
			Genre:    "Electronic",
			Year:     1990,
			Duration: 370,
		},
		{
			ID:       "track-4",
			Title:    "Trans-Europe Express",
			Artist:   "Kraftwerk",
			Album:    "Trans-Europe Express",
			Genre:    "Electronic",
			Year:     1977,
			Duration: 405,
		},
		{
			ID:       "track-5",
			Title:    "Personal Jesus",
			Artist:   "Depeche Mode",
			Album:    "Violator",
			Genre:    "Electronic",
			Year:     1990,
			Duration: 295,
		},
	}
}

func TestHelloService_SearchTracks_MatchesTitle(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.SearchTracks(ctx, "neon")

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "track-1", results[0].ID)
	assert.Equal(t, "Neon Lights", results[0].Title)
}

func TestHelloService_SearchTracks_CaseInsensitive(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)

	// Search with uppercase should match lowercase content
	resultsUpper, err := svc.SearchTracks(ctx, "NEON")
	require.NoError(t, err)

	// Search with mixed case
	resultsMixed, err := svc.SearchTracks(ctx, "NeOn")
	require.NoError(t, err)

	// Search with lowercase
	resultsLower, err := svc.SearchTracks(ctx, "neon")
	require.NoError(t, err)

	// All should return same results
	assert.Len(t, resultsUpper, 1)
	assert.Len(t, resultsMixed, 1)
	assert.Len(t, resultsLower, 1)
	assert.Equal(t, resultsUpper[0].ID, resultsMixed[0].ID)
	assert.Equal(t, resultsMixed[0].ID, resultsLower[0].ID)
}

func TestHelloService_SearchTracks_MatchesArtist(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.SearchTracks(ctx, "kraftwerk")

	require.NoError(t, err)
	require.Len(t, results, 2)
	// Both Kraftwerk tracks should be returned
	ids := []string{results[0].ID, results[1].ID}
	assert.Contains(t, ids, "track-1")
	assert.Contains(t, ids, "track-4")
}

func TestHelloService_SearchTracks_MatchesAlbum(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.SearchTracks(ctx, "violator")

	require.NoError(t, err)
	require.Len(t, results, 2)
	// Both tracks from the Violator album
	ids := []string{results[0].ID, results[1].ID}
	assert.Contains(t, ids, "track-3")
	assert.Contains(t, ids, "track-5")
}

func TestHelloService_SearchTracks_MatchesGenre(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.SearchTracks(ctx, "synth-pop")

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "track-2", results[0].ID)
	assert.Equal(t, "Blue Monday", results[0].Title)
}

func TestHelloService_SearchTracks_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.SearchTracks(ctx, "")

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestHelloService_SearchTracks_NoMatch(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.SearchTracks(ctx, "zzzznonexistent")

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestHelloService_GetFeaturedTracks_ReturnsAll(t *testing.T) {
	ctx := context.Background()
	tracks := sampleTracks()
	mockRepo := &MockHelloRepository{
		tracks: tracks,
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.GetFeaturedTracks(ctx, 0)

	require.NoError(t, err)
	// With limit=0, should use default limit (20), which is more than sample tracks
	assert.Len(t, results, len(tracks))
}

func TestHelloService_GetFeaturedTracks_WithLimit(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockHelloRepository{
		tracks: sampleTracks(),
	}

	svc := NewHelloService(mockRepo)
	results, err := svc.GetFeaturedTracks(ctx, 3)

	require.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestHelloService_GetFeaturedTracks_DefaultLimit(t *testing.T) {
	ctx := context.Background()

	// Create 25 tracks to test default limit
	tracks := make([]HelloTrack, 25)
	for i := 0; i < 25; i++ {
		tracks[i] = HelloTrack{
			ID:       "track-" + string(rune('a'+i)),
			Title:    "Track " + string(rune('A'+i)),
			Artist:   "Artist",
			Album:    "Album",
			Genre:    "Genre",
			Year:     2024,
			Duration: 200,
		}
	}

	mockRepo := &MockHelloRepository{
		tracks: tracks,
	}

	svc := NewHelloService(mockRepo)
	// Pass 0 to trigger default limit
	results, err := svc.GetFeaturedTracks(ctx, 0)

	require.NoError(t, err)
	// Default limit should be 20
	assert.Len(t, results, 20)
}
