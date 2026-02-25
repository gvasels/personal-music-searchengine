package service

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Events Service Tests
// =============================================================================

// MockEventsProvider provides a mockable events.EventsProvider for tests.
type MockEventsProvider struct {
	mock.Mock
}

func (m *MockEventsProvider) GetArtistEvents(ctx context.Context, artistName string) ([]models.Event, error) {
	args := m.Called(ctx, artistName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockEventsProvider) SearchArtists(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ArtistSearchResult), args.Error(1)
}

// MockArtistWatchServiceForEvents mocks the ArtistWatchService for events tests.
type MockArtistWatchServiceForEvents struct {
	mock.Mock
}

func (m *MockArtistWatchServiceForEvents) WatchArtist(ctx context.Context, userID, artistName string) error {
	args := m.Called(ctx, userID, artistName)
	return args.Error(0)
}

func (m *MockArtistWatchServiceForEvents) UnwatchArtist(ctx context.Context, userID, artistName string) error {
	args := m.Called(ctx, userID, artistName)
	return args.Error(0)
}

func (m *MockArtistWatchServiceForEvents) IsWatching(ctx context.Context, userID, artistName string) (bool, error) {
	args := m.Called(ctx, userID, artistName)
	return args.Bool(0), args.Error(1)
}

func (m *MockArtistWatchServiceForEvents) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error) {
	args := m.Called(ctx, userID, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.ArtistWatch]), args.Error(1)
}

// TestEventsService_GetArtistEvents verifies that GetArtistEvents delegates to
// the events provider and returns the results.
func TestEventsService_GetArtistEvents(t *testing.T) {
	t.Run("returns events from provider", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		now := time.Now().UTC()
		expectedEvents := []models.Event{
			{
				ID:         "evt-1",
				ArtistName: "Daft Punk",
				Title:      "Daft Punk Live",
				Date:       now.Add(30 * 24 * time.Hour),
				Venue:      "Madison Square Garden",
				City:       "New York",
				Country:    "US",
				Status:     "scheduled",
				Source:     "mock",
			},
			{
				ID:         "evt-2",
				ArtistName: "Daft Punk",
				Title:      "Daft Punk Live",
				Date:       now.Add(60 * 24 * time.Hour),
				Venue:      "The O2",
				City:       "London",
				Country:    "UK",
				Status:     "scheduled",
				Source:     "mock",
			},
		}
		mockProvider.On("GetArtistEvents", ctx, "Daft Punk").Return(expectedEvents, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		events, err := svc.GetArtistEvents(ctx, "Daft Punk")

		require.NoError(t, err)
		assert.Len(t, events, 2)
		assert.Equal(t, "evt-1", events[0].ID)
		assert.Equal(t, "evt-2", events[1].ID)
		mockProvider.AssertExpectations(t)
	})

	t.Run("returns error from provider", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		mockProvider.On("GetArtistEvents", ctx, "Unknown Artist").Return(([]models.Event)(nil), errors.New("provider error"))

		svc := NewEventsService(mockProvider, mockWatch)
		events, err := svc.GetArtistEvents(ctx, "Unknown Artist")

		require.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "provider error")
	})

	t.Run("returns empty slice when no events", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		mockProvider.On("GetArtistEvents", ctx, "No Events Artist").Return([]models.Event{}, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		events, err := svc.GetArtistEvents(ctx, "No Events Artist")

		require.NoError(t, err)
		assert.Empty(t, events)
	})
}

// TestEventsService_SearchArtists verifies that SearchArtists delegates to
// the events provider and returns the results.
func TestEventsService_SearchArtists(t *testing.T) {
	t.Run("returns search results from provider", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		expectedResults := []models.ArtistSearchResult{
			{Name: "Daft Punk", UpcomingEvents: 3, Source: "mock"},
			{Name: "Deadmau5", UpcomingEvents: 2, Source: "mock"},
		}
		mockProvider.On("SearchArtists", ctx, "daft", 10).Return(expectedResults, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.SearchArtists(ctx, "daft", 10)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Daft Punk", results[0].Name)
		mockProvider.AssertExpectations(t)
	})

	t.Run("returns error from provider", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		mockProvider.On("SearchArtists", ctx, "query", 5).Return(([]models.ArtistSearchResult)(nil), errors.New("search error"))

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.SearchArtists(ctx, "query", 5)

		require.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("returns empty results for no matches", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		mockProvider.On("SearchArtists", ctx, "nonexistent", 10).Return([]models.ArtistSearchResult{}, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.SearchArtists(ctx, "nonexistent", 10)

		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		expectedResults := []models.ArtistSearchResult{
			{Name: "Daft Punk", UpcomingEvents: 3, Source: "mock"},
		}
		mockProvider.On("SearchArtists", ctx, "d", 1).Return(expectedResults, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.SearchArtists(ctx, "d", 1)

		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

// TestEventsService_GetWatchedArtistEvents verifies that GetWatchedArtistEvents
// fetches the user's watched artists, then retrieves events for each, and
// returns the results aggregated and sorted by date.
func TestEventsService_GetWatchedArtistEvents(t *testing.T) {
	t.Run("returns aggregated events for watched artists", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		// User watches two artists
		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items: []models.ArtistWatch{
				{UserID: "user-123", ArtistName: "Daft Punk", WatchedAt: time.Now()},
				{UserID: "user-123", ArtistName: "Deadmau5", WatchedAt: time.Now()},
			},
			HasMore: false,
		}
		// ListWatchedArtists is called with a high limit to get all watched artists
		mockWatch.On("ListWatchedArtists", ctx, "user-123", mock.AnythingOfType("int"), "").Return(watches, nil)

		now := time.Now().UTC()
		daftPunkEvents := []models.Event{
			{
				ID:         "dp-1",
				ArtistName: "Daft Punk",
				Title:      "Daft Punk Live",
				Date:       now.Add(60 * 24 * time.Hour),
				Venue:      "Madison Square Garden",
				City:       "New York",
				Country:    "US",
				Status:     "scheduled",
				Source:     "mock",
			},
		}
		deadmau5Events := []models.Event{
			{
				ID:         "dm-1",
				ArtistName: "Deadmau5",
				Title:      "Deadmau5 Live",
				Date:       now.Add(30 * 24 * time.Hour),
				Venue:      "The O2",
				City:       "London",
				Country:    "UK",
				Status:     "scheduled",
				Source:     "mock",
			},
		}
		mockProvider.On("GetArtistEvents", ctx, "Daft Punk").Return(daftPunkEvents, nil)
		mockProvider.On("GetArtistEvents", ctx, "Deadmau5").Return(deadmau5Events, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.GetWatchedArtistEvents(ctx, "user-123")

		require.NoError(t, err)
		require.Len(t, results, 2)

		// Verify both artists are represented
		artistNames := make([]string, len(results))
		for i, r := range results {
			artistNames[i] = r.ArtistName
		}
		sort.Strings(artistNames)
		assert.Equal(t, []string{"Daft Punk", "Deadmau5"}, artistNames)

		// Verify each has events
		for _, r := range results {
			assert.NotEmpty(t, r.Events)
			assert.Equal(t, len(r.Events), r.TotalCount)
		}
	})

	t.Run("returns empty when no watched artists", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items:   []models.ArtistWatch{},
			HasMore: false,
		}
		mockWatch.On("ListWatchedArtists", ctx, "user-123", mock.AnythingOfType("int"), "").Return(watches, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.GetWatchedArtistEvents(ctx, "user-123")

		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("returns error when watch service fails", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		mockWatch.On("ListWatchedArtists", ctx, "user-123", mock.AnythingOfType("int"), "").
			Return((*repository.PaginatedResult[models.ArtistWatch])(nil), errors.New("watch service error"))

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.GetWatchedArtistEvents(ctx, "user-123")

		require.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "watch service error")
	})

	t.Run("skips artist when provider returns error for one artist", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items: []models.ArtistWatch{
				{UserID: "user-123", ArtistName: "Daft Punk", WatchedAt: time.Now()},
				{UserID: "user-123", ArtistName: "Failing Artist", WatchedAt: time.Now()},
			},
			HasMore: false,
		}
		mockWatch.On("ListWatchedArtists", ctx, "user-123", mock.AnythingOfType("int"), "").Return(watches, nil)

		now := time.Now().UTC()
		daftPunkEvents := []models.Event{
			{
				ID:         "dp-1",
				ArtistName: "Daft Punk",
				Title:      "Daft Punk Live",
				Date:       now.Add(30 * 24 * time.Hour),
				Venue:      "Madison Square Garden",
				City:       "New York",
				Country:    "US",
				Status:     "scheduled",
				Source:     "mock",
			},
		}
		mockProvider.On("GetArtistEvents", ctx, "Daft Punk").Return(daftPunkEvents, nil)
		mockProvider.On("GetArtistEvents", ctx, "Failing Artist").Return(([]models.Event)(nil), errors.New("provider timeout"))

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.GetWatchedArtistEvents(ctx, "user-123")

		// Should succeed with partial results — the failing artist is skipped
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, "Daft Punk", results[0].ArtistName)
	})

	t.Run("events within each artist are sorted by date ascending", func(t *testing.T) {
		ctx := context.Background()
		mockProvider := new(MockEventsProvider)
		mockWatch := new(MockArtistWatchServiceForEvents)

		watches := &repository.PaginatedResult[models.ArtistWatch]{
			Items: []models.ArtistWatch{
				{UserID: "user-123", ArtistName: "Daft Punk", WatchedAt: time.Now()},
			},
			HasMore: false,
		}
		mockWatch.On("ListWatchedArtists", ctx, "user-123", mock.AnythingOfType("int"), "").Return(watches, nil)

		now := time.Now().UTC()
		events := []models.Event{
			{ID: "dp-later", ArtistName: "Daft Punk", Date: now.Add(90 * 24 * time.Hour), Status: "scheduled", Source: "mock"},
			{ID: "dp-sooner", ArtistName: "Daft Punk", Date: now.Add(30 * 24 * time.Hour), Status: "scheduled", Source: "mock"},
			{ID: "dp-middle", ArtistName: "Daft Punk", Date: now.Add(60 * 24 * time.Hour), Status: "scheduled", Source: "mock"},
		}
		mockProvider.On("GetArtistEvents", ctx, "Daft Punk").Return(events, nil)

		svc := NewEventsService(mockProvider, mockWatch)
		results, err := svc.GetWatchedArtistEvents(ctx, "user-123")

		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Len(t, results[0].Events, 3)

		// Verify events are sorted by date ascending
		for i := 1; i < len(results[0].Events); i++ {
			assert.True(t, results[0].Events[i-1].Date.Before(results[0].Events[i].Date) || results[0].Events[i-1].Date.Equal(results[0].Events[i].Date),
				"events should be sorted by date ascending: %v should be before %v",
				results[0].Events[i-1].Date, results[0].Events[i].Date)
		}
	})
}
