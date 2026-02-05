# Design: Hello World Local Development

## Overview

A minimal full-stack feature demonstrating the local development workflow: DynamoDB seed data → Go backend API → React frontend. Validates SDLC process with TDD enforcement.

## Steering Document Alignment

### Technical Standards
Per root `CLAUDE.md`:
- **Backend**: Go 1.22+, Echo v4, handler → service → repository layering
- **Frontend**: React 18, TanStack Router/Query, DaisyUI 5, Zustand
- **Data**: DynamoDB single-table design with `USER#{userId}` / `TRACK#{trackId}` keys

### Project Structure
```
backend/internal/
├── handlers/hello.go      # HTTP handlers
├── service/hello.go       # Business logic + HelloTrack type
├── service/hello_dynamo.go # DynamoDB repository implementation

frontend/src/
├── lib/api/helloSearch.ts # API client
├── hooks/useHelloSearch.ts # TanStack Query hooks
├── components/hello/      # UI components
└── routes/hello-search.tsx # Page route

docker/localstack-init/
└── init-seed-music.sh    # Seed data script
```

## Code Reuse Analysis

### Existing Components to Leverage
- **`backend/internal/handlers/handlers.go`**: `ListResponse[T]` generic for paginated responses
- **`backend/internal/repository/repository.go`**: DynamoDB client patterns (Query, Scan)
- **`frontend/src/lib/api/client.ts`**: `apiClient` axios instance with interceptors
- **`frontend/src/test/test-utils.tsx`**: Test wrapper with QueryClient

### Integration Points
- **DynamoDB**: Single-table `MusicLibrary`, keys `PK=USER#seed-user`, `SK=TRACK#{id}`
- **LocalStack**: Endpoint `http://localhost:4566`, seeded by `docker/localstack-init/` scripts
- **Vite Proxy**: Frontend `/api` → `http://localhost:8080`

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  React Frontend │────▶│   Go Backend    │────▶│   LocalStack    │
│  :5173          │     │   :8080         │     │   DynamoDB      │
│                 │     │                 │     │   :4566         │
│  /hello-search  │     │  /api/v1/hello/ │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

### Data Flow
1. LocalStack starts → `init-seed-music.sh` creates 20 tracks
2. Frontend calls `/api/v1/hello/featured` or `/api/v1/hello/search?q=...`
3. Backend queries DynamoDB for `PK=USER#seed-user, SK begins_with TRACK#`
4. Backend filters results (for search) and returns `{"items": [...], "total": N}`
5. Frontend renders TrackCard components

## Components and Interfaces

### Backend: HelloService

```go
type HelloTrack struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Artist   string `json:"artist"`
    Album    string `json:"album"`
    Genre    string `json:"genre"`
    Year     int    `json:"year"`
    Duration int    `json:"duration"`
}

type HelloRepository interface {
    GetSeedTracks(ctx context.Context) ([]HelloTrack, error)
}

type HelloService struct {
    repo HelloRepository
}

func (s *HelloService) Search(ctx context.Context, query string) ([]HelloTrack, error)
func (s *HelloService) Featured(ctx context.Context, limit int) ([]HelloTrack, error)
```

### Backend: HelloHandler

```go
type HelloHandler struct {
    service HelloServiceInterface
}

func (h *HelloHandler) Health(c echo.Context) error    // GET /api/v1/hello/health
func (h *HelloHandler) Search(c echo.Context) error    // GET /api/v1/hello/search?q=
func (h *HelloHandler) Featured(c echo.Context) error  // GET /api/v1/hello/featured?limit=
func RegisterHelloRoutes(e *echo.Echo, h *HelloHandler)
```

### Frontend: API Client

```typescript
// lib/api/helloSearch.ts
export interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;
}

export interface HelloSearchResponse {
  items: HelloTrack[];
  total: number;
}

export function searchHelloTracks(query: string): Promise<HelloSearchResponse>
export function getFeaturedHelloTracks(limit?: number): Promise<HelloSearchResponse>
export function getHelloHealth(): Promise<{ status: string; service: string }>
```

### Frontend: Hooks

```typescript
// hooks/useHelloSearch.ts
export const helloKeys = {
  all: ['hello'] as const,
  search: (q: string) => [...helloKeys.all, 'search', q] as const,
  featured: () => [...helloKeys.all, 'featured'] as const,
};

export function useHelloSearch(query: string): UseQueryResult<HelloSearchResponse>
export function useHelloFeatured(): UseQueryResult<HelloSearchResponse>
```

### Frontend: Components

| Component | Props | Description |
|-----------|-------|-------------|
| `SearchInput` | `value: string, onChange: (v: string) => void, placeholder?: string` | Controlled text input |
| `TrackCard` | `track: HelloTrack` | Card showing title, artist, album, genre, year, duration (MM:SS) |
| `TrackCardSkeleton` | (none) | Loading placeholder with `.skeleton` elements |
| `HelloNav` | (none) | Navbar with "Hello Music Search" title and home link |

## Data Models

### DynamoDB Seed Track

```
PK: USER#seed-user
SK: TRACK#seed-t1
Title: "Aurora Borealis"
Artist: "Aurora Waves"
Album: "Dreamscape"
Genre: "jazz"
Year: 2023
Duration: 245
```

20 tracks total: 5 artists × 4 tracks each

### API Response Contract

**CRITICAL - All endpoints MUST use this format:**

```json
{
  "items": [
    {
      "id": "seed-t1",
      "title": "Aurora Borealis",
      "artist": "Aurora Waves",
      "album": "Dreamscape",
      "genre": "jazz",
      "year": 2023,
      "duration": 245
    }
  ],
  "total": 1
}
```

- Search: `GET /api/v1/hello/search?q=jazz` (param is `q`, NOT `query`)
- Featured: `GET /api/v1/hello/featured?limit=10`
- Health: `GET /api/v1/hello/health` → `{"status":"ok","service":"hello"}`

## Error Handling

| Scenario | HTTP Status | Response |
|----------|-------------|----------|
| Empty search query | 200 | `{"items":[],"total":0}` |
| No matches | 200 | `{"items":[],"total":0}` |
| DynamoDB error | 500 | `{"error":"Internal server error"}` |

## Testing Strategy

### Unit Testing
- **HelloService**: Mock repository, test Search/Featured logic
- **HelloHandler**: Mock service, test HTTP request/response with httptest
- **Components**: Vitest + React Testing Library

### Contract Testing
- Verify `q` param (not `query`) in frontend API client
- Verify `items` key (not `tracks`) in response handling
- Grep-based verification in Phase 4

### Integration Testing
- `make local` with curl verification
- Manual browser check at `http://localhost:5173/hello-search`
