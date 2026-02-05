# Design Document

## Overview

The hello-world-local-dev feature adds three components: (1) a seed data script that populates DynamoDB with 20 mock tracks, (2) a Go backend with HelloService + HelloHandler exposing search/featured/health endpoints, and (3) a React frontend search page with components for displaying track cards and a search input.

## Technical Standards

Per root `CLAUDE.md`: Go 1.22+, Echo v4, DynamoDB single-table with `PK`/`SK` key schema, React 18, TanStack Router/Query, DaisyUI 5, Zustand. All backend code follows handler -> service -> repository layering in `backend/internal/`.

## Project Structure

Per root `CLAUDE.md` Repository Structure: backend code in `backend/internal/{handlers,service,repository}/`, frontend in `frontend/src/{components,hooks,lib/api,routes}/`, Docker init scripts in `docker/localstack-init/`.

## Code Reuse Analysis

### Existing Components to Leverage
- **`backend/internal/handlers/handlers.go`**: `Handlers` struct, `RegisterRoutes`, `handleError`, `success`, `getAuthContext` helpers
- **`backend/internal/service/service.go`**: `Services` struct for wiring HelloService
- **`backend/internal/repository/dynamodb.go`**: `DynamoDBRepository` with DynamoDB client and table name — reuse for querying seed tracks
- **`frontend/src/lib/api/client.ts`**: `apiClient` axios instance with auth interceptor
- **`frontend/src/test/test-utils.tsx`**: Custom render with QueryClient wrapper for testing

### Integration Points
- **DynamoDB MusicLibrary table**: Seed data uses existing table with `PK=USER#seed-user`, `SK=TRACK#{trackId}` pattern
- **Makefile `local-services` target**: Add seed script invocation after `init-cognito.sh`
- **Vite proxy**: Frontend `/api` calls proxy to `http://localhost:8080` (existing config)

## Architecture

```
Seed Script (init-seed-music.sh)
    │
    ▼
DynamoDB (MusicLibrary table, PK=USER#seed-user)
    │
    ▼
HelloService (service/hello.go)
    │ ListTracksByUser → filter by query
    ▼
HelloHandler (handlers/hello.go)
    │ /api/v1/hello/{health,search,featured}
    ▼
Frontend API Client (lib/api/helloSearch.ts)
    │
    ▼
TanStack Query Hooks (hooks/useHelloSearch.ts)
    │
    ▼
React Components + Route (routes/hello-search.tsx)
```

## Components and Interfaces

### Seed Script (`docker/localstack-init/init-seed-music.sh`)
- **Purpose**: Insert 20 mock tracks into DynamoDB under `seed-user` partition
- **Interface**: Bash script, no arguments, idempotent (checks before inserting)
- **Dependencies**: AWS CLI, LocalStack running on port 4566
- **Data**: 5 artists (Luna Waves, The Ember Collective, DJ Phantom, Aria Chen, Voltage), 4 tracks each

### Seed Validation Script (`docker/localstack-init/test-seed-data.sh`)
- **Purpose**: Validate seed data correctness (count, fields, artists)
- **Interface**: Bash script, exits non-zero on failure
- **Dependencies**: AWS CLI, jq, seed data already inserted

### HelloService (`backend/internal/service/hello.go`)
- **Purpose**: Business logic for searching and listing seed tracks
- **Interfaces**:
  ```go
  type HelloService interface {
      SearchTracks(ctx context.Context, query string, limit int) ([]models.Track, error)
      GetFeaturedTracks(ctx context.Context, limit int) ([]models.Track, error)
  }
  ```
- **Dependencies**: `HelloRepository` interface (adapter over main Repository)
- **Reuses**: `repository.Repository.ListTracks` via `HelloRepoAdapter`

### HelloRepoAdapter (`backend/internal/service/hello.go`)
- **Purpose**: Adapts main Repository to HelloRepository interface
- **Interface**:
  ```go
  type HelloRepository interface {
      ListTracksByUser(ctx context.Context, userID string) ([]models.Track, error)
  }
  ```

### HelloHandler (`backend/internal/handlers/hello.go`)
- **Purpose**: HTTP handlers for hello endpoints
- **Routes**:
  | Method | Path | Handler | Auth |
  |--------|------|---------|------|
  | GET | `/api/v1/hello/health` | `HelloHealth` | None |
  | GET | `/api/v1/hello/search` | `HelloSearch` | None |
  | GET | `/api/v1/hello/featured` | `HelloFeatured` | None |
- **Dependencies**: `HelloService`
- **Reuses**: `handleError`, `success` from handlers package

### Frontend API Client (`frontend/src/lib/api/helloSearch.ts`)
- **Purpose**: API functions for hello endpoints
- **Interfaces**:
  ```typescript
  export function searchHello(query: string, limit?: number): Promise<{ items: Track[] }>
  export function getFeaturedTracks(limit?: number): Promise<{ items: Track[] }>
  ```
- **Reuses**: `apiClient` from `client.ts`, `Track` type from `types.ts`

### Frontend Hook (`frontend/src/hooks/useHelloSearch.ts`)
- **Purpose**: TanStack Query hooks for hello API
- **Interfaces**:
  ```typescript
  export const helloSearchKeys = {
    all: ['hello-search'] as const,
    search: (query: string) => [...helloSearchKeys.all, 'search', query] as const,
    featured: () => [...helloSearchKeys.all, 'featured'] as const,
  }
  export function useHelloSearch(query: string): UseQueryResult
  export function useFeaturedTracks(): UseQueryResult
  ```
- **Reuses**: TanStack Query `useQuery`, queryKey factory pattern from other hooks

### Frontend Components (`frontend/src/components/hello/`)

| Component | Purpose | Props |
|-----------|---------|-------|
| `SearchInput` | Text input with Enter/Escape handling | `onSearch: (q: string) => void` |
| `TrackCard` | Displays track info (title, artist, album, genre badge, duration) | `track: Track` |
| `TrackCardSkeleton` | Loading placeholder matching TrackCard layout | (none) |
| `HelloNav` | Navbar with search input and theme toggle | `onSearch: (q: string) => void` |

### Frontend Route (`frontend/src/routes/hello-search.tsx`)
- **Purpose**: `/hello-search` page with hero, featured tracks grid, and search results
- **Dependencies**: All hello components + hooks
- **Auth**: Public route (add to `PUBLIC_ROUTES` in `__root.tsx`)

## Data Models

### Seed Track (DynamoDB Item)
```
PK: "USER#seed-user"
SK: "TRACK#{trackId}"
id: string (UUID)
title: string
artist: string
album: string
genre: string
year: number
duration: number (seconds)
Type: "TRACK"
createdAt: ISO8601 string
```

### API Response (search and featured)
```json
{
  "items": [
    {
      "id": "track-001",
      "title": "Neon Dreams",
      "artist": "Luna Waves",
      "album": "Electric Horizons",
      "genre": "Electronic",
      "year": 2024,
      "duration": 234
    }
  ]
}
```

## Error Handling

### Error Scenarios
1. **Missing query parameter on search**
   - Handling: Return 400 `{"error": "query parameter 'q' is required"}`
   - User Impact: Frontend shows validation error

2. **DynamoDB unavailable**
   - Handling: Return 500 with error message
   - User Impact: Frontend shows error state

3. **No search results**
   - Handling: Return 200 with `{"items": []}`
   - User Impact: Frontend shows "no results found" message

## Testing Strategy

### Unit Testing (Backend)
- HelloService: Test SearchTracks with various queries, case-insensitivity, empty query, limit
- HelloService: Test GetFeaturedTracks with default and custom limits
- HelloHandler: Test HTTP responses for all 3 endpoints, error cases

### Unit Testing (Frontend)
- SearchInput: Test Enter/Escape key handling, value updates
- TrackCard: Test rendering all track fields, genre badge
- TrackCardSkeleton: Test skeleton structure renders
- HelloNav: Test search callback, theme toggle presence
- useHelloSearch hook: Test query key factory, enabled flag
- useFeaturedTracks hook: Test data fetching
- helloSearch API: Test URL construction, parameter passing
- hello-search route: Test hero section, featured tracks rendering, search flow, loading/error states

### Integration Testing
- Seed script: Validate 20 items, 5 artists, required fields, idempotency (via test-seed-data.sh)
