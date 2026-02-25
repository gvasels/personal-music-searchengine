# Events Package - CLAUDE.md

## Overview

Events provider abstraction for artist event data. Defines the `EventsProvider` interface and includes a `MockProvider` for development/testing that generates deterministic fake events.

## File Descriptions

| File | Purpose |
|------|---------|
| `provider.go` | `EventsProvider` interface definition |
| `mock_provider.go` | `MockProvider` - deterministic fake event data for dev/test |
| `provider_test.go` | Unit tests for MockProvider (8 tests) |

## Key Types and Functions

### EventsProvider Interface (`provider.go`)
```go
type EventsProvider interface {
    GetArtistEvents(ctx context.Context, artistName string) ([]models.Event, error)
    SearchArtists(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error)
}
```

### MockProvider (`mock_provider.go`)
| Function | Description |
|----------|-------------|
| `NewMockProvider()` | Creates a new MockProvider |
| `GetArtistEvents(ctx, artistName)` | Generates 2-4 deterministic events using FNV-1a hash of artist name |
| `SearchArtists(ctx, query, limit)` | Case-insensitive substring search against built-in artist list |

**Determinism**: Same artist name always produces same events (same venues, date offsets). Dates are relative to `time.Now()` so they shift forward naturally.

**Venue Pool**: 6 real venues (Madison Square Garden, The O2, Ziggo Dome, Olympiastadion, Forum Melbourne, Budokan).

**Mock Artists**: 10 built-in names for search (Kylie Minogue, Deadmau5, Armin van Buuren, etc.).

## Dependencies

### Internal
- `github.com/gvasels/personal-music-searchengine/internal/models` - Event, ArtistSearchResult types

### External
- `hash/fnv` - FNV-1a hash for deterministic output
