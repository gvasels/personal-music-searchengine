package events

import (
	"context"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Compile-time check: MockProvider must implement EventsProvider.
var _ EventsProvider = &MockProvider{}

func TestMockProvider_GetArtistEvents(t *testing.T) {
	provider := NewMockProvider()
	ctx := context.Background()

	events, err := provider.GetArtistEvents(ctx, "Kylie Minogue")
	require.NoError(t, err)

	t.Run("returns 2-4 events", func(t *testing.T) {
		assert.GreaterOrEqual(t, len(events), 2, "should return at least 2 events")
		assert.LessOrEqual(t, len(events), 4, "should return at most 4 events")
	})

	t.Run("all events have correct artist name", func(t *testing.T) {
		for i, event := range events {
			assert.Equal(t, "Kylie Minogue", event.ArtistName,
				"event[%d] should have ArtistName = 'Kylie Minogue'", i)
		}
	})

	t.Run("all events are in the future", func(t *testing.T) {
		now := time.Now()
		for i, event := range events {
			assert.True(t, event.Date.After(now),
				"event[%d] date %v should be after now %v", i, event.Date, now)
		}
	})

	t.Run("all events have source mock", func(t *testing.T) {
		for i, event := range events {
			assert.Equal(t, "mock", event.Source,
				"event[%d] should have Source = 'mock'", i)
		}
	})

	t.Run("all events have non-empty venue city country", func(t *testing.T) {
		for i, event := range events {
			assert.NotEmpty(t, event.Venue,
				"event[%d] should have a non-empty Venue", i)
			assert.NotEmpty(t, event.City,
				"event[%d] should have a non-empty City", i)
			assert.NotEmpty(t, event.Country,
				"event[%d] should have a non-empty Country", i)
		}
	})

	t.Run("all events have status scheduled", func(t *testing.T) {
		for i, event := range events {
			assert.Equal(t, "scheduled", event.Status,
				"event[%d] should have Status = 'scheduled'", i)
		}
	})

	t.Run("all events have non-empty ID and Title", func(t *testing.T) {
		for i, event := range events {
			assert.NotEmpty(t, event.ID,
				"event[%d] should have a non-empty ID", i)
			assert.NotEmpty(t, event.Title,
				"event[%d] should have a non-empty Title", i)
		}
	})

	t.Run("event dates are 1-6 months in the future", func(t *testing.T) {
		now := time.Now()
		oneMonth := now.AddDate(0, 1, 0)
		sixMonths := now.AddDate(0, 6, 0)

		for i, event := range events {
			assert.True(t, event.Date.After(oneMonth) || event.Date.Equal(oneMonth),
				"event[%d] date %v should be at least 1 month from now", i, event.Date)
			assert.True(t, event.Date.Before(sixMonths) || event.Date.Equal(sixMonths),
				"event[%d] date %v should be at most 6 months from now", i, event.Date)
		}
	})

	// Suppress unused import warning for models package by using Event type explicitly
	_ = models.Event{}
}

func TestMockProvider_GetArtistEvents_Deterministic(t *testing.T) {
	provider := NewMockProvider()
	ctx := context.Background()

	events1, err1 := provider.GetArtistEvents(ctx, "Kylie Minogue")
	require.NoError(t, err1)

	events2, err2 := provider.GetArtistEvents(ctx, "Kylie Minogue")
	require.NoError(t, err2)

	t.Run("same number of events", func(t *testing.T) {
		require.Equal(t, len(events1), len(events2),
			"two calls with the same artist name should return the same number of events")
	})

	t.Run("same event IDs", func(t *testing.T) {
		for i := range events1 {
			assert.Equal(t, events1[i].ID, events2[i].ID,
				"event[%d] ID should be identical across calls", i)
		}
	})

	t.Run("same event dates", func(t *testing.T) {
		for i := range events1 {
			assert.True(t, events1[i].Date.Equal(events2[i].Date),
				"event[%d] Date should be identical across calls: got %v and %v",
				i, events1[i].Date, events2[i].Date)
		}
	})

	t.Run("same venues", func(t *testing.T) {
		for i := range events1 {
			assert.Equal(t, events1[i].Venue, events2[i].Venue,
				"event[%d] Venue should be identical across calls", i)
			assert.Equal(t, events1[i].City, events2[i].City,
				"event[%d] City should be identical across calls", i)
			assert.Equal(t, events1[i].Country, events2[i].Country,
				"event[%d] Country should be identical across calls", i)
		}
	})
}

func TestMockProvider_GetArtistEvents_DifferentArtists(t *testing.T) {
	provider := NewMockProvider()
	ctx := context.Background()

	events1, err1 := provider.GetArtistEvents(ctx, "Kylie Minogue")
	require.NoError(t, err1)
	require.NotEmpty(t, events1, "should return events for Kylie Minogue")

	events2, err2 := provider.GetArtistEvents(ctx, "MÖWE")
	require.NoError(t, err2)
	require.NotEmpty(t, events2, "should return events for MÖWE")

	t.Run("different artist names in events", func(t *testing.T) {
		assert.Equal(t, "Kylie Minogue", events1[0].ArtistName)
		assert.Equal(t, "MÖWE", events2[0].ArtistName)
	})

	t.Run("different event sets", func(t *testing.T) {
		// Collect all event IDs from each artist
		ids1 := make(map[string]bool)
		for _, e := range events1 {
			ids1[e.ID] = true
		}

		// At least one event ID from artist2 should differ
		hasDifferentID := false
		for _, e := range events2 {
			if !ids1[e.ID] {
				hasDifferentID = true
				break
			}
		}
		assert.True(t, hasDifferentID,
			"different artists should produce different event IDs")
	})
}

func TestMockProvider_SearchArtists(t *testing.T) {
	provider := NewMockProvider()
	ctx := context.Background()

	results, err := provider.SearchArtists(ctx, "kylie", 10)
	require.NoError(t, err)

	t.Run("returns results for matching query", func(t *testing.T) {
		require.NotEmpty(t, results, "should return results for 'kylie'")
	})

	t.Run("all results have source mock", func(t *testing.T) {
		for i, result := range results {
			assert.Equal(t, "mock", result.Source,
				"result[%d] should have Source = 'mock'", i)
		}
	})

	t.Run("all results have a name", func(t *testing.T) {
		for i, result := range results {
			assert.NotEmpty(t, result.Name,
				"result[%d] should have a non-empty Name", i)
		}
	})

	t.Run("all results have non-negative upcoming events count", func(t *testing.T) {
		for i, result := range results {
			assert.GreaterOrEqual(t, result.UpcomingEvents, 0,
				"result[%d] UpcomingEvents should be >= 0", i)
		}
	})

	t.Run("results match query case-insensitively", func(t *testing.T) {
		// Search is case-insensitive substring match, so all returned names
		// should contain "kylie" (case-insensitive)
		for i, result := range results {
			assert.Contains(t, toLower(result.Name), "kylie",
				"result[%d] Name '%s' should contain 'kylie' (case-insensitive)", i, result.Name)
		}
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		limited, err := provider.SearchArtists(ctx, "kylie", 1)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(limited), 1,
			"should return at most 1 result when limit=1")
	})

	// Ensure the ArtistSearchResult type from models is referenced
	_ = models.ArtistSearchResult{}
}

func TestMockProvider_SearchArtists_EmptyQuery(t *testing.T) {
	provider := NewMockProvider()
	ctx := context.Background()

	results, err := provider.SearchArtists(ctx, "", 10)
	require.NoError(t, err)

	t.Run("returns empty or all results for empty query", func(t *testing.T) {
		// Empty query may return empty slice or all available artists --
		// either behavior is acceptable.
		assert.NotNil(t, results, "results should not be nil (empty slice is OK)")
	})

	t.Run("all results still have source mock", func(t *testing.T) {
		for i, result := range results {
			assert.Equal(t, "mock", result.Source,
				"result[%d] should have Source = 'mock'", i)
		}
	})
}

func TestMockProvider_ImplementsInterface(t *testing.T) {
	// This test verifies at runtime that MockProvider satisfies EventsProvider.
	// The compile-time check (var _ EventsProvider = &MockProvider{}) at the
	// top of this file is the primary enforcement; this test provides explicit
	// documentation of the contract.
	var provider EventsProvider = NewMockProvider()
	assert.NotNil(t, provider, "NewMockProvider should return a non-nil EventsProvider")
}

// toLower is a helper to lowercase strings for case-insensitive comparison.
func toLower(s string) string {
	// Using a simple byte-level lowering; sufficient for ASCII artist names in mock data.
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}
