package events

import (
	"context"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// EventsProvider abstracts event data sources (mock, Bandsintown, Ticketmaster, etc.)
type EventsProvider interface {
	GetArtistEvents(ctx context.Context, artistName string) ([]models.Event, error)
	SearchArtists(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error)
}
