package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// seedTracks returns a fixture of 20 diverse tracks for testing
// 5 artists x 4 tracks each with varied genres (jazz, rock, electronic, ambient, soul)
func seedTracks() []HelloTrack {
	return []HelloTrack{
		// Aurora Waves - Jazz/Ambient artist
		{ID: "seed-t1", Title: "Aurora Borealis", Artist: "Aurora Waves", Album: "Dreamscape", Genre: "jazz", Year: 2023, Duration: 245},
		{ID: "seed-t2", Title: "Midnight Jazz", Artist: "Aurora Waves", Album: "Dreamscape", Genre: "jazz", Year: 2023, Duration: 312},
		{ID: "seed-t3", Title: "Ethereal Dreams", Artist: "Aurora Waves", Album: "Nocturne", Genre: "ambient", Year: 2022, Duration: 428},
		{ID: "seed-t4", Title: "Cosmic Flow", Artist: "Aurora Waves", Album: "Nocturne", Genre: "ambient", Year: 2022, Duration: 356},

		// Electric Pulse - Electronic artist
		{ID: "seed-t5", Title: "Neon Nights", Artist: "Electric Pulse", Album: "Voltage", Genre: "electronic", Year: 2024, Duration: 198},
		{ID: "seed-t6", Title: "Digital Dawn", Artist: "Electric Pulse", Album: "Voltage", Genre: "electronic", Year: 2024, Duration: 223},
		{ID: "seed-t7", Title: "Synthetic Soul", Artist: "Electric Pulse", Album: "Circuits", Genre: "electronic", Year: 2023, Duration: 267},
		{ID: "seed-t8", Title: "Binary Beat", Artist: "Electric Pulse", Album: "Circuits", Genre: "electronic", Year: 2023, Duration: 189},

		// Stone Temple - Rock artist
		{ID: "seed-t9", Title: "Thunder Road", Artist: "Stone Temple", Album: "Granite", Genre: "rock", Year: 2021, Duration: 284},
		{ID: "seed-t10", Title: "Mountain High", Artist: "Stone Temple", Album: "Granite", Genre: "rock", Year: 2021, Duration: 301},
		{ID: "seed-t11", Title: "Desert Wind", Artist: "Stone Temple", Album: "Sandstone", Genre: "rock", Year: 2020, Duration: 256},
		{ID: "seed-t12", Title: "Canyon Echo", Artist: "Stone Temple", Album: "Sandstone", Genre: "rock", Year: 2020, Duration: 342},

		// Velvet Soul - Soul artist
		{ID: "seed-t13", Title: "Soulful Morning", Artist: "Velvet Soul", Album: "Smooth", Genre: "soul", Year: 2022, Duration: 275},
		{ID: "seed-t14", Title: "Heart of Gold", Artist: "Velvet Soul", Album: "Smooth", Genre: "soul", Year: 2022, Duration: 298},
		{ID: "seed-t15", Title: "Rhythm Divine", Artist: "Velvet Soul", Album: "Groove", Genre: "soul", Year: 2021, Duration: 234},
		{ID: "seed-t16", Title: "Sweet Sensation", Artist: "Velvet Soul", Album: "Groove", Genre: "soul", Year: 2021, Duration: 312},

		// Ambient Dreams - Ambient/Jazz crossover artist
		{ID: "seed-t17", Title: "Ocean Waves", Artist: "Ambient Dreams", Album: "Serenity", Genre: "ambient", Year: 2024, Duration: 467},
		{ID: "seed-t18", Title: "Forest Rain", Artist: "Ambient Dreams", Album: "Serenity", Genre: "ambient", Year: 2024, Duration: 389},
		{ID: "seed-t19", Title: "Starlight Jazz", Artist: "Ambient Dreams", Album: "Celestial", Genre: "jazz", Year: 2023, Duration: 345},
		{ID: "seed-t20", Title: "Moonlit Sonata", Artist: "Ambient Dreams", Album: "Celestial", Genre: "jazz", Year: 2023, Duration: 412},
	}
}

// MockHelloRepo is a mock implementation of HelloRepository
type MockHelloRepo struct {
	mock.Mock
}

func (m *MockHelloRepo) GetSeedTracks(ctx context.Context) ([]HelloTrack, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]HelloTrack), args.Error(1)
}

func TestNewHelloService(t *testing.T) {
	t.Run("creates service with repository", func(t *testing.T) {
		mockRepo := new(MockHelloRepo)
		svc := NewHelloService(mockRepo)
		require.NotNil(t, svc)
	})
}

func TestHelloService_Search(t *testing.T) {
	t.Run("matches title case-insensitive", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "aurora")

		require.NoError(t, err)
		require.NotNil(t, results)
		// Should match "Aurora Borealis" title
		assert.True(t, len(results) > 0, "expected at least one match for 'aurora'")
		for _, track := range results {
			assert.True(t,
				strings.Contains(strings.ToLower(track.Title), "aurora") ||
					strings.Contains(strings.ToLower(track.Artist), "aurora") ||
					strings.Contains(strings.ToLower(track.Album), "aurora") ||
					strings.Contains(strings.ToLower(track.Genre), "aurora"),
				"track should match 'aurora' in title, artist, album, or genre")
		}
	})

	t.Run("matches artist case-insensitive", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "ELECTRIC PULSE")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 4, "expected 4 tracks from Electric Pulse")
		for _, track := range results {
			assert.Equal(t, "Electric Pulse", track.Artist)
		}
	})

	t.Run("matches album case-insensitive", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "dreamscape")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 2, "expected 2 tracks from Dreamscape album")
		for _, track := range results {
			assert.Equal(t, "Dreamscape", track.Album)
		}
	})

	t.Run("matches genre case-insensitive", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "JAZZ")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 4, "expected 4 jazz tracks")
		for _, track := range results {
			assert.Equal(t, "jazz", track.Genre)
		}
	})

	t.Run("empty query returns empty slice not nil", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		// Note: with empty query, we might not even call the repo
		// but if we do, set it up
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil).Maybe()

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "")

		require.NoError(t, err)
		require.NotNil(t, results, "results should be empty slice, not nil")
		assert.Len(t, results, 0, "empty query should return zero results")
	})

	t.Run("no matches returns empty slice not nil", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "nonexistent query xyz123")

		require.NoError(t, err)
		require.NotNil(t, results, "results should be empty slice, not nil")
		assert.Len(t, results, 0, "non-matching query should return zero results")
	})

	t.Run("partial match works", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "soul")

		require.NoError(t, err)
		require.NotNil(t, results)
		// Should match:
		// - "Velvet Soul" (artist, 4 tracks)
		// - "Synthetic Soul" (title, 1 track)
		// - "soul" (genre, 4 tracks from Velvet Soul)
		// Total unique: 5 tracks (4 Velvet Soul + 1 Synthetic Soul)
		assert.True(t, len(results) >= 4, "expected at least 4 tracks matching 'soul'")
	})

	t.Run("handles repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(([]HelloTrack)(nil), assert.AnError)

		svc := NewHelloService(mockRepo)
		results, err := svc.Search(ctx, "jazz")

		require.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestHelloService_Featured(t *testing.T) {
	t.Run("returns up to limit tracks", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Featured(ctx, 5)

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 5, "should return exactly 5 tracks when limit is 5")
	})

	t.Run("returns all tracks when limit exceeds count", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Featured(ctx, 100)

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 20, "should return all 20 tracks when limit exceeds count")
	})

	t.Run("limit 0 returns all tracks", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Featured(ctx, 0)

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 20, "limit 0 should return all 20 tracks")
	})

	t.Run("returns empty slice when no tracks", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return([]HelloTrack{}, nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Featured(ctx, 10)

		require.NoError(t, err)
		require.NotNil(t, results, "results should be empty slice, not nil")
		assert.Len(t, results, 0)
	})

	t.Run("handles repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(([]HelloTrack)(nil), assert.AnError)

		svc := NewHelloService(mockRepo)
		results, err := svc.Featured(ctx, 10)

		require.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("negative limit treated as zero (returns all)", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockHelloRepo)
		mockRepo.On("GetSeedTracks", ctx).Return(seedTracks(), nil)

		svc := NewHelloService(mockRepo)
		results, err := svc.Featured(ctx, -5)

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 20, "negative limit should return all tracks")
	})
}
