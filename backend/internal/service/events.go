package service

import (
	"context"
	"sort"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
)

// EventsProvider abstracts event data sources for the EventsService.
type EventsProvider interface {
	GetArtistEvents(ctx context.Context, artistName string) ([]models.Event, error)
	SearchArtists(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error)
}

// ArtistWatchServiceForEvents defines the watch service methods needed by EventsService.
type ArtistWatchServiceForEvents interface {
	WatchArtist(ctx context.Context, userID, artistName string) error
	UnwatchArtist(ctx context.Context, userID, artistName string) error
	IsWatching(ctx context.Context, userID, artistName string) (bool, error)
	ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatch], error)
}

// EventsService handles event-related operations.
type EventsService struct {
	provider EventsProvider
	watchSvc ArtistWatchServiceForEvents
}

// NewEventsService creates a new EventsService.
func NewEventsService(provider EventsProvider, watchSvc ArtistWatchServiceForEvents) *EventsService {
	return &EventsService{
		provider: provider,
		watchSvc: watchSvc,
	}
}

// GetArtistEvents returns events for the given artist from the provider.
func (s *EventsService) GetArtistEvents(ctx context.Context, artistName string) ([]models.Event, error) {
	return s.provider.GetArtistEvents(ctx, artistName)
}

// SearchArtists searches for artists via the provider.
func (s *EventsService) SearchArtists(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error) {
	return s.provider.SearchArtists(ctx, query, limit)
}

// GetWatchedArtistEvents returns aggregated events for all artists the user is watching.
// Events within each artist are sorted by date ascending. Artists where the provider
// returns an error are silently skipped (partial failure tolerance).
func (s *EventsService) GetWatchedArtistEvents(ctx context.Context, userID string) ([]models.ArtistEventsResponse, error) {
	// Fetch all watched artists with a high limit
	watches, err := s.watchSvc.ListWatchedArtists(ctx, userID, 1000, "")
	if err != nil {
		return nil, err
	}

	if len(watches.Items) == 0 {
		return []models.ArtistEventsResponse{}, nil
	}

	var results []models.ArtistEventsResponse
	for _, watch := range watches.Items {
		events, err := s.provider.GetArtistEvents(ctx, watch.ArtistName)
		if err != nil {
			// Skip artists where the provider fails (partial failure)
			continue
		}

		if len(events) == 0 {
			continue
		}

		// Sort events by date ascending
		sort.Slice(events, func(i, j int) bool {
			return events[i].Date.Before(events[j].Date)
		})

		results = append(results, models.ArtistEventsResponse{
			ArtistName: watch.ArtistName,
			Events:     events,
			TotalCount: len(events),
			Source:     events[0].Source,
		})
	}

	if results == nil {
		results = []models.ArtistEventsResponse{}
	}

	return results, nil
}
