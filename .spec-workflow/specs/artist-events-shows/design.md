# Design Document: Artist Events & Shows (Phase 1 — Mock)

## Overview

Add artist watching and event discovery to the music library. Users can watch artists from their catalog and see upcoming events aggregated on a "My Shows" page. Phase 1 uses a mock events provider; Phase 2 will integrate real APIs (Bandsintown, Ticketmaster, SeatGeek, etc.) via the same provider interface.

## Code Reuse Analysis

### Existing Components to Leverage

- **Artist model** (`backend/internal/models/artist.go`): Artist entity already created on ingest via `getOrCreateArtist()` in track processor. PK=`USER#{userId}`, SK=`ARTIST#{artistId}`, GSI1 for name lookups.
- **Follow system** (`backend/internal/service/follow.go`, `repository/follow.go`): Pattern for watch/unwatch — DynamoDB PutItem/DeleteItem with GSI queries. ArtistWatch follows same pattern.
- **FollowButton component** (`frontend/src/components/follow/FollowButton.tsx`): UI pattern for toggle buttons with optimistic updates.
- **useFollow hook** (`frontend/src/hooks/useFollow.ts`): TanStack Query mutation pattern with cache invalidation.
- **Artist detail page** (`frontend/src/routes/artists/$artistName.tsx`): Existing page to extend with events section and watch button.
- **Track detail page** (`frontend/src/routes/tracks/$trackId.tsx`): Existing page to add watch icon next to artist link.
- **Sidebar** (`frontend/src/components/layout/Sidebar.tsx`): Add "Shows" nav link.
- **DynamoDB repository patterns** (`backend/internal/repository/dynamodb.go`): Existing CRUD, pagination, GSI query patterns.

### Integration Points

- **Ingest pipeline**: No changes needed — Artist entities already created per attributed artist
- **DynamoDB table**: New ArtistWatch entity in same single-table
- **API router** (`backend/cmd/api/main.go`): Register new handler routes
- **Services struct** (`backend/internal/service/`): Add EventsService and ArtistWatchService

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Frontend                             │
│                                                          │
│  /shows (My Shows page)     /artists/$name (detail)     │
│  ├── WatchedArtistEvents    ├── ArtistEventsSection      │
│  ├── ArtistEventSearch      ├── WatchButton              │
│  └── EventCard              └── EventCard                │
│                                                          │
│  Hooks: useArtistWatch, useArtistEvents                  │
│  API:   artistWatch.ts, events.ts                        │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTP
┌──────────────────────▼──────────────────────────────────┐
│                     Backend API                          │
│                                                          │
│  Handlers:                                               │
│  ├── WatchArtist    POST   /artists/:name/watch          │
│  ├── UnwatchArtist  DELETE /artists/:name/watch          │
│  ├── GetWatchStatus GET    /artists/:name/watch          │
│  ├── ListWatched    GET    /users/me/watched-artists     │
│  ├── GetEvents      GET    /artists/:name/events         │
│  └── SearchEvents   GET    /events/search?q=             │
│                                                          │
│  Services:                                               │
│  ├── ArtistWatchService → DynamoDB (ArtistWatch entity)  │
│  └── EventsService      → EventsProvider interface       │
│                             ├── MockProvider (Phase 1)    │
│                             ├── BandsintownProvider       │
│                             ├── TicketmasterProvider      │
│                             └── ... (Phase 2)             │
└──────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Backend: EventsProvider Interface

- **Purpose:** Abstract events data source so providers can be swapped without changing business logic
- **Interface:**
  ```go
  type EventsProvider interface {
      GetArtistEvents(ctx context.Context, artistName string) ([]Event, error)
      SearchArtists(ctx context.Context, query string, limit int) ([]ArtistSearchResult, error)
  }
  ```
- **Dependencies:** None (pure interface)
- **Reuses:** None — new abstraction

### Backend: MockProvider

- **Purpose:** Return deterministic fake event data for development/testing
- **Interface:** Implements `EventsProvider`
- **Dependencies:** None
- **Behavior:** Generates 2-4 fake events per artist using hash of artist name for determinism. Events are set 1-6 months in the future. Venues use a pool of real city names.

### Backend: ArtistWatchService

- **Purpose:** Manage user artist watch relationships
- **Interface:**
  ```go
  WatchArtist(ctx context.Context, userID, artistName string) error
  UnwatchArtist(ctx context.Context, userID, artistName string) error
  IsWatching(ctx context.Context, userID, artistName string) (bool, error)
  ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*PaginatedResult, error)
  ```
- **Dependencies:** Repository
- **Reuses:** Follow service pattern (watch/unwatch is structurally identical to follow/unfollow)

### Backend: EventsService

- **Purpose:** Orchestrate events queries through the provider interface
- **Interface:**
  ```go
  GetArtistEvents(ctx context.Context, artistName string) ([]EventResponse, error)
  SearchArtists(ctx context.Context, query string, limit int) ([]ArtistSearchResultResponse, error)
  GetWatchedArtistEvents(ctx context.Context, userID string) ([]ArtistEventsResponse, error)
  ```
- **Dependencies:** EventsProvider, ArtistWatchService
- **Reuses:** Service layer patterns from existing services

### Backend: EventsHandler

- **Purpose:** HTTP handler for events and watch endpoints
- **Interface:** Echo handler methods
- **Dependencies:** EventsService, ArtistWatchService
- **Reuses:** Existing handler patterns (auth middleware, pagination parsing, error responses)

### Frontend: WatchButton Component

- **Purpose:** Toggle watch/unwatch for an artist
- **Props:** `artistName: string`, `size?: 'sm' | 'md'`
- **Reuses:** FollowButton pattern (toggle state, optimistic updates, auth check)

### Frontend: EventCard Component

- **Purpose:** Display a single event with date, venue, city, ticket link
- **Props:** `event: ArtistEvent`
- **Reuses:** DaisyUI card pattern from existing components

### Frontend: ArtistEventsSection Component

- **Purpose:** Show upcoming events for one artist (used on artist detail page)
- **Props:** `artistName: string`
- **Reuses:** SimilarTracks component pattern (fetch-on-mount, loading skeleton, error state)

### Frontend: My Shows Page (`/shows`)

- **Purpose:** Aggregated view of events for all watched artists
- **Reuses:** Track listing page patterns (pagination, empty states, search input)

## Data Models

### ArtistWatch (DynamoDB)

```
PK:     USER#{userId}
SK:     ARTIST_WATCH#{normalizedArtistName}
GSI1PK: ARTIST_WATCH#{normalizedArtistName}
GSI1SK: USER#{userId}
Type:   ARTIST_WATCH

Fields:
- userId:     string   (user who is watching)
- artistName: string   (original casing for display)
- watchedAt:  string   (ISO 8601 timestamp)
```

Normalized name = lowercase, trimmed. Used as SK to prevent duplicate watches for same artist with different casing.

GSI1 enables "who is watching this artist?" queries (not needed for Phase 1 but future-proofs).

### Event (provider response model)

```go
type Event struct {
    ID          string    `json:"id"`
    ArtistName  string    `json:"artistName"`
    Title       string    `json:"title"`       // e.g. "Summer Tour 2026"
    Date        time.Time `json:"date"`
    Venue       string    `json:"venue"`
    City        string    `json:"city"`
    Region      string    `json:"region"`      // state/province
    Country     string    `json:"country"`
    TicketURL   string    `json:"ticketUrl,omitempty"`
    Status      string    `json:"status"`      // scheduled, cancelled, postponed
    Source      string    `json:"source"`      // "mock", "bandsintown", "ticketmaster"
}
```

### ArtistSearchResult (provider response model)

```go
type ArtistSearchResult struct {
    Name           string `json:"name"`
    ImageURL       string `json:"imageUrl,omitempty"`
    UpcomingEvents int    `json:"upcomingEvents"`
    Source         string `json:"source"`
}
```

### Frontend Types

```typescript
interface ArtistEvent {
  id: string;
  artistName: string;
  title: string;
  date: string;           // ISO 8601
  venue: string;
  city: string;
  region: string;
  country: string;
  ticketUrl?: string;
  status: 'scheduled' | 'cancelled' | 'postponed';
  source: string;
}

interface ArtistSearchResult {
  name: string;
  imageUrl?: string;
  upcomingEvents: number;
  source: string;
}

interface WatchedArtist {
  artistName: string;
  watchedAt: string;
}
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/artists/:name/watch` | Subscriber+ | Watch an artist |
| DELETE | `/api/v1/artists/:name/watch` | Subscriber+ | Unwatch an artist |
| GET | `/api/v1/artists/:name/watch` | Subscriber+ | Check watch status |
| GET | `/api/v1/users/me/watched-artists` | Subscriber+ | List watched artists |
| GET | `/api/v1/artists/:name/events` | Subscriber+ | Get events for artist |
| GET | `/api/v1/events/search?q=` | Subscriber+ | Search artists for events |

## Error Handling

### Error Scenarios

1. **Events provider failure (mock won't fail, but future APIs will)**
   - **Handling:** Return empty events list with error flag, log error
   - **User Impact:** "Unable to load events" message, rest of page works

2. **Watch duplicate artist**
   - **Handling:** DynamoDB conditional put — idempotent, returns success
   - **User Impact:** No visible change (already watching)

3. **Unwatch non-watched artist**
   - **Handling:** DynamoDB delete — idempotent, returns success
   - **User Impact:** No visible change

4. **Artist name not found in provider**
   - **Handling:** Return empty events list
   - **User Impact:** "No upcoming events found"

## Testing Strategy

### Unit Testing
- **Models**: ArtistWatch entity creation, key normalization
- **Services**: ArtistWatchService CRUD with mocked repository, EventsService with mocked provider
- **Handlers**: HTTP request/response with mocked services
- **MockProvider**: Deterministic output verification
- **Frontend hooks**: useArtistWatch, useArtistEvents with mocked API
- **Frontend components**: WatchButton, EventCard, ArtistEventsSection, MyShowsPage

### Integration Testing
- ArtistWatch DynamoDB CRUD + pagination
- Watch/unwatch API endpoints with auth middleware
- Events endpoints returning mock data

### End-to-End Testing
- Watch artist from artist detail page → appears on My Shows page
- Search for artist on My Shows page → see events
- Unwatch → removed from My Shows
