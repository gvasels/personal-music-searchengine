# Design Document: Hello World Local Development Validation

## Overview

This design adds a standalone hello-world music search feature to validate the local development stack end-to-end. It introduces a HelloService, HelloHandler, seed script, and React search page -- all following existing codebase patterns. The feature is isolated from the main application logic and requires no authentication.

## Existing Infrastructure Analysis

### What Already Exists

| Component | Location | Reuse Strategy |
|-----------|----------|----------------|
| Handler pattern | `backend/internal/handlers/handlers.go` | Follow `Handlers` struct, `RegisterRoutes`, `success()`, `successList()` patterns |
| Service interfaces | `backend/internal/service/service.go` | Follow `Services` struct, `NewServices()` initialization pattern |
| Repository interface | `backend/internal/repository/repository.go` | Use existing `Repository` interface methods (`ListTracks`, `GetTrack`) |
| API client | `frontend/src/lib/api/client.ts` | Use `apiClient` axios instance for API calls |
| TanStack Query hooks | `frontend/src/hooks/useTracks.ts` | Follow `queryKey` factory + `useQuery` pattern |
| DaisyUI card components | `frontend/src/components/library/` | Follow existing card styling patterns |
| Docker Compose | `docker/docker-compose.yml` | LocalStack already configured with DynamoDB |
| Init scripts | `docker/localstack-init/init-aws.sh` | Follow idempotent init script pattern |
| Makefile | `Makefile` | Add seed script call to `local-services` target |

### What Needs to Be Added

| Component | Purpose |
|-----------|---------|
| Seed script | `docker/localstack-init/init-seed-music.sh` -- insert 20 mock tracks |
| HelloService | Business logic for search and featured tracks |
| HelloHandler | HTTP handlers for `/api/v1/hello/*` routes |
| HelloRepoAdapter | Bridge `repository.Repository` to `HelloRepository` interface |
| Frontend API module | `frontend/src/lib/api/helloSearch.ts` |
| Frontend hooks | `frontend/src/hooks/useHelloSearch.ts` |
| Frontend components | SearchInput, TrackCard, TrackCardSkeleton, HelloNav |
| Frontend route | `frontend/src/routes/hello-search.tsx` |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 Frontend (localhost:5173)                     │
│  /hello-search route                                         │
│  ┌──────────┐  ┌──────────┐  ┌────────────────┐            │
│  │SearchInput│  │TrackCard │  │TrackCardSkeleton│            │
│  └─────┬────┘  └──────────┘  └────────────────┘            │
│        │                                                     │
│  ┌─────▼─────────────────────────────────┐                  │
│  │ useHelloSearch / useFeaturedTracks     │                  │
│  │ (TanStack Query hooks)                │                  │
│  └─────┬─────────────────────────────────┘                  │
│        │                                                     │
│  ┌─────▼─────────────────────────────────┐                  │
│  │ helloSearch.ts (API client)           │                  │
│  │ Uses apiClient from client.ts         │                  │
│  └─────┬─────────────────────────────────┘                  │
│        │ GET /api/v1/hello/search?q=...                      │
│        │ GET /api/v1/hello/featured                           │
└────────┼────────────────────────────────────────────────────┘
         │ Vite proxy → localhost:8080
┌────────▼────────────────────────────────────────────────────┐
│                 Backend (localhost:8080)                      │
│  ┌──────────────────────────────────────┐                   │
│  │ HelloHandler                         │                   │
│  │  - HelloHealth()                     │                   │
│  │  - HelloSearch(q)                    │                   │
│  │  - HelloFeatured(limit)             │                   │
│  └─────┬────────────────────────────────┘                   │
│        │                                                     │
│  ┌─────▼────────────────────────────────┐                   │
│  │ HelloService                         │                   │
│  │  - SearchTracks(ctx, query)          │                   │
│  │  - GetFeaturedTracks(ctx, limit)     │                   │
│  └─────┬────────────────────────────────┘                   │
│        │                                                     │
│  ┌─────▼────────────────────────────────┐                   │
│  │ HelloRepoAdapter                     │                   │
│  │  wraps Repository interface          │                   │
│  │  - GetTracksByUser(ctx, userID)      │                   │
│  └─────┬────────────────────────────────┘                   │
│        │                                                     │
│  ┌─────▼────────────────────────────────┐                   │
│  │ DynamoDB (LocalStack :4566)          │                   │
│  │  MusicLibrary table                   │                   │
│  │  PK=USER#seed-user, SK=TRACK#...     │                   │
│  └──────────────────────────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Component 1: Seed Script

**Purpose:** Insert 20 mock tracks into DynamoDB MusicLibrary table for local development.

**File:** `docker/localstack-init/init-seed-music.sh`

**Interface:**
```bash
#!/bin/bash
# Seeds 20 mock tracks into LocalStack DynamoDB
# Idempotent: uses condition-expression to skip existing items
# Uses: aws dynamodb put-item --endpoint-url http://localhost:4566

ENDPOINT="${AWS_ENDPOINT:-http://localhost:4566}"
TABLE="MusicLibrary"
USER_ID="seed-user"
```

**Data:** 20 tracks across 5 artists:

| Artist | Tracks | Albums | Genre |
|--------|--------|--------|-------|
| Aurora Waves | 4 | Neon Horizons | Synthwave |
| The Midnight Echo | 4 | Velvet Thunder | Indie Rock |
| Luna Noir | 4 | Shadow Dance | Electronic |
| Crimson Tide | 4 | Ocean Drive | Funk |
| Stellar Drift | 4 | Cosmic Bloom | Ambient |

Each track item:
```json
{
  "PK": {"S": "USER#seed-user"},
  "SK": {"S": "TRACK#hello-001"},
  "Title": {"S": "Neon Dreams"},
  "Artist": {"S": "Aurora Waves"},
  "Album": {"S": "Neon Horizons"},
  "Genre": {"S": "Synthwave"},
  "Year": {"N": "2024"},
  "Duration": {"N": "234"},
  "TrackID": {"S": "hello-001"},
  "UserID": {"S": "seed-user"},
  "CreatedAt": {"S": "2024-01-01T00:00:00Z"},
  "UpdatedAt": {"S": "2024-01-01T00:00:00Z"},
  "Visibility": {"S": "public"}
}
```

**Idempotency:** Each `put-item` uses `--condition-expression "attribute_not_exists(PK)"` to skip existing items without error (exit code suppressed with `|| true`).

**Dependencies:** AWS CLI, LocalStack running on port 4566

### Component 2: HelloService (Backend)

**Purpose:** Business logic for searching and listing seed tracks.

**File:** `backend/internal/service/hello.go`

**Interface:**
```go
package service

import "context"

// HelloTrack represents a track in hello search results
type HelloTrack struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Artist   string `json:"artist"`
    Album    string `json:"album"`
    Genre    string `json:"genre"`
    Year     int    `json:"year"`
    Duration int    `json:"duration"`
}

// HelloRepository defines the data access interface for hello service
type HelloRepository interface {
    GetTracksByUser(ctx context.Context, userID string) ([]HelloTrack, error)
}

// HelloService handles hello-world search operations
type HelloService struct {
    repo   HelloRepository
    userID string // hardcoded "seed-user"
}

// NewHelloService creates a new HelloService
func NewHelloService(repo HelloRepository) *HelloService {
    return &HelloService{
        repo:   repo,
        userID: "seed-user",
    }
}

// SearchTracks searches seed tracks by query (case-insensitive match on title, artist, album, genre)
func (s *HelloService) SearchTracks(ctx context.Context, query string) ([]HelloTrack, error)

// GetFeaturedTracks returns all seed tracks up to limit
func (s *HelloService) GetFeaturedTracks(ctx context.Context, limit int) ([]HelloTrack, error)
```

**Search Logic:**
- Fetch all tracks for `seed-user` from repository
- Filter in-memory using `strings.Contains(strings.ToLower(field), strings.ToLower(query))`
- Match against title, artist, album, and genre fields
- Return empty slice (not nil) when no matches or empty query

**Dependencies:** `HelloRepository` interface (satisfied by `HelloRepoAdapter`)

### Component 3: HelloRepoAdapter (Backend)

**Purpose:** Adapts the existing `repository.Repository` interface to the `HelloRepository` interface used by HelloService.

**File:** `backend/internal/service/hello.go` (same file as HelloService)

**Interface:**
```go
// HelloRepoAdapter adapts Repository to HelloRepository
type HelloRepoAdapter struct {
    repo repository.Repository
}

// NewHelloRepoAdapter creates a new adapter
func NewHelloRepoAdapter(repo repository.Repository) *HelloRepoAdapter {
    return &HelloRepoAdapter{repo: repo}
}

// GetTracksByUser fetches all tracks for a user using ListTracks with a large limit
func (a *HelloRepoAdapter) GetTracksByUser(ctx context.Context, userID string) ([]HelloTrack, error)
```

**Implementation Notes:**
- Calls `repo.ListTracks(ctx, userID, models.TrackFilter{Limit: 100})` to fetch tracks
- Converts `models.Track` to `HelloTrack` (mapping relevant fields)
- The adapter ensures HelloService does not directly depend on the repository package

### Component 4: HelloHandler (Backend)

**Purpose:** HTTP handlers for hello-world API endpoints.

**File:** `backend/internal/handlers/hello.go`

**Interface:**
```go
package handlers

import "github.com/labstack/echo/v4"

// HelloHandler handles hello-world HTTP endpoints
type HelloHandler struct {
    service *service.HelloService
}

// NewHelloHandler creates a new HelloHandler
func NewHelloHandler(svc *service.HelloService) *HelloHandler {
    return &HelloHandler{service: svc}
}

// HelloHealth returns hello service health status
// GET /api/v1/hello/health
func (h *HelloHandler) HelloHealth(c echo.Context) error {
    return c.JSON(200, map[string]string{
        "status":  "ok",
        "service": "hello",
    })
}

// HelloSearch searches seed tracks by query parameter
// GET /api/v1/hello/search?q={query}
func (h *HelloHandler) HelloSearch(c echo.Context) error

// HelloFeatured returns featured (all) seed tracks
// GET /api/v1/hello/featured?limit={limit}
func (h *HelloHandler) HelloFeatured(c echo.Context) error

// RegisterHelloRoutes registers hello routes on the Echo instance
// Routes are PUBLIC (no auth middleware)
func RegisterHelloRoutes(e *echo.Echo, h *HelloHandler) {
    hello := e.Group("/api/v1/hello")
    hello.GET("/health", h.HelloHealth)
    hello.GET("/search", h.HelloSearch)
    hello.GET("/featured", h.HelloFeatured)
}
```

**Response Format:** Uses the existing `ListResponse[T]` pattern with `{ "items": [...], "total": N }`.

**Dependencies:** `service.HelloService`

### Component 5: Wiring (CRITICAL)

**Purpose:** Wire HelloService and HelloHandler into the application bootstrap.

This is the **most critical** part -- unit tests with mocks cannot catch missing wiring. Follow the wiring checklist (`.claude/docs/wiring-checklist.md`).

#### Step 5a: Add HelloService to Services struct

**File:** `backend/internal/service/service.go`

```go
type Services struct {
    // ...existing services...
    Hello  *HelloService  // <- ADD THIS
}
```

#### Step 5b: Initialize HelloService in NewServices()

**File:** `backend/internal/service/service.go`

```go
func NewServices(...) *Services {
    return &Services{
        // ...existing initializations...
        Hello: NewHelloService(NewHelloRepoAdapter(repo)),  // <- ADD THIS
    }
}
```

#### Step 5c: Create HelloHandler and register routes in setupEcho()

**File:** `backend/cmd/api/main.go`

```go
func setupEcho() (*echo.Echo, error) {
    // ... existing code ...

    // Register routes
    h.RegisterRoutes(e)

    // Register hello routes (public, no auth) -- ADD THIS BLOCK
    helloHandler := handlers.NewHelloHandler(services.Hello)
    handlers.RegisterHelloRoutes(e, helloHandler)

    // ... rest of existing code ...
}
```

### Component 6: Frontend API Client

**Purpose:** API functions for hello search endpoints.

**File:** `frontend/src/lib/api/helloSearch.ts`

**Interface:**
```typescript
import { apiClient } from './client';

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

// Search tracks by query
export async function searchHelloTracks(query: string): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/search', {
    params: { q: query },
  });
  return response.data;
}

// Get featured tracks
export async function getFeaturedTracks(limit?: number): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/featured', {
    params: { limit },
  });
  return response.data;
}

// Health check
export async function getHelloHealth(): Promise<{ status: string; service: string }> {
  const response = await apiClient.get<{ status: string; service: string }>('/v1/hello/health');
  return response.data;
}
```

**Note:** Uses `apiClient` from `client.ts` which has `baseURL: '/api'`. The path `/v1/hello/search` becomes `/api/v1/hello/search` via the Vite proxy.

### Component 7: Frontend Hooks

**Purpose:** TanStack Query hooks for hello search data fetching.

**File:** `frontend/src/hooks/useHelloSearch.ts`

**Interface:**
```typescript
import { useQuery } from '@tanstack/react-query';
import { searchHelloTracks, getFeaturedTracks } from '@/lib/api/helloSearch';

// Query key factory
export const helloKeys = {
  all: ['hello'] as const,
  search: (query: string) => [...helloKeys.all, 'search', query] as const,
  featured: (limit?: number) => [...helloKeys.all, 'featured', limit] as const,
};

// Search tracks hook with debounce-friendly enabled flag
export function useHelloSearch(query: string) {
  return useQuery({
    queryKey: helloKeys.search(query),
    queryFn: () => searchHelloTracks(query),
    enabled: query.length > 0,
  });
}

// Featured tracks hook
export function useFeaturedTracks(limit?: number) {
  return useQuery({
    queryKey: helloKeys.featured(limit),
    queryFn: () => getFeaturedTracks(limit),
  });
}
```

### Component 8: Frontend Components

**Purpose:** UI components for the hello search page.

#### SearchInput

**File:** `frontend/src/components/hello/SearchInput.tsx`

```typescript
interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}
```

- Text input with DaisyUI `input input-bordered` classes
- Magnifying glass icon
- Controlled component with value/onChange props

#### TrackCard

**File:** `frontend/src/components/hello/TrackCard.tsx`

```typescript
interface TrackCardProps {
  track: HelloTrack;
}
```

- DaisyUI `card bg-base-200` styling
- Displays: title (bold), artist, album, genre badge, year, formatted duration (MM:SS)
- Uses `base-200` background, `text-primary` for title

#### TrackCardSkeleton

**File:** `frontend/src/components/hello/TrackCardSkeleton.tsx`

- DaisyUI `skeleton` class for loading placeholders
- Matches TrackCard layout dimensions
- No props required

#### HelloNav

**File:** `frontend/src/components/hello/HelloNav.tsx`

- Simple navigation link back to home
- Link to `/hello-search` for sidebar integration if desired

### Component 9: Frontend Route

**Purpose:** The hello search page route.

**File:** `frontend/src/routes/hello-search.tsx`

**Structure:**
```typescript
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/hello-search')({
  component: HelloSearchPage,
});

function HelloSearchPage() {
  // State: searchQuery (string), debouncedQuery (string)
  // Hooks: useHelloSearch(debouncedQuery), useFeaturedTracks()
  // Layout:
  //   - Hero section with title + subtitle
  //   - SearchInput
  //   - If searching: search results grid (or skeletons, or error)
  //   - If not searching: featured tracks grid (or skeletons, or error)
}
```

**Public Route:** Must be added to the `PUBLIC_ROUTES` array in `__root.tsx` so it does not redirect unauthenticated users to `/permission-denied`.

### Component 10: Root Route Update

**File:** `frontend/src/routes/__root.tsx`

Add `/hello-search` to the `PUBLIC_ROUTES` constant:

```typescript
const PUBLIC_ROUTES = ['/', '/login', '/permission-denied', '/hello-search'];
```

## Data Model

### DynamoDB Item Schema (Seed Tracks)

| Attribute | Type | Description |
|-----------|------|-------------|
| PK | S | `USER#seed-user` |
| SK | S | `TRACK#{trackId}` (e.g., `TRACK#hello-001`) |
| TrackID | S | Unique ID (e.g., `hello-001`) |
| UserID | S | `seed-user` |
| Title | S | Track title |
| Artist | S | Artist name |
| Album | S | Album name |
| Genre | S | Genre |
| Year | N | Release year |
| Duration | N | Duration in seconds |
| CreatedAt | S | ISO 8601 timestamp |
| UpdatedAt | S | ISO 8601 timestamp |
| Visibility | S | `public` |

## Error Handling

### Backend Errors

| Scenario | HTTP Code | Response |
|----------|-----------|----------|
| DynamoDB unavailable | 500 | `{ "error": { "code": "INTERNAL_ERROR", "message": "..." } }` |
| Invalid limit parameter | 400 | `{ "error": { "code": "BAD_REQUEST", "message": "..." } }` |
| Missing query param on search | 200 | `{ "items": [], "total": 0 }` (not an error) |

### Frontend Errors

| Scenario | UI Behavior |
|----------|-------------|
| API request fails | Show error alert with retry button |
| Network timeout | Show error alert with retry button |
| Loading state | Show TrackCardSkeleton grid |

## File Structure

```
personal-music-searchengine/
├── docker/localstack-init/
│   └── init-seed-music.sh              # NEW: Seed 20 mock tracks
├── Makefile                            # MODIFY: Add seed script call
├── backend/internal/
│   ├── service/
│   │   ├── service.go                  # MODIFY: Add Hello to Services struct + NewServices()
│   │   └── hello.go                    # NEW: HelloService, HelloRepoAdapter, HelloRepository
│   └── handlers/
│       └── hello.go                    # NEW: HelloHandler, RegisterHelloRoutes
├── backend/cmd/api/
│   └── main.go                         # MODIFY: Wire HelloHandler in setupEcho()
└── frontend/src/
    ├── lib/api/
    │   └── helloSearch.ts              # NEW: API client functions
    ├── hooks/
    │   └── useHelloSearch.ts           # NEW: TanStack Query hooks
    ├── components/hello/
    │   ├── SearchInput.tsx             # NEW: Search input component
    │   ├── TrackCard.tsx               # NEW: Track card component
    │   ├── TrackCardSkeleton.tsx       # NEW: Loading skeleton
    │   └── HelloNav.tsx               # NEW: Navigation component
    └── routes/
        ├── __root.tsx                  # MODIFY: Add /hello-search to PUBLIC_ROUTES
        └── hello-search.tsx            # NEW: Hello search page route
```

## Implementation Notes

### Why a Separate HelloRepository Interface?

The existing `Repository` interface has many methods. HelloService only needs to list tracks for a single user. The `HelloRepository` interface keeps HelloService decoupled and testable with a simple mock, while `HelloRepoAdapter` bridges the gap to the real repository.

### Why In-Memory Search Instead of DynamoDB Query?

With only 20 tracks, fetching all and filtering in Go is simpler and faster than building DynamoDB filter expressions. This keeps the hello-world feature minimal while still validating the full DynamoDB read path.

### Why No Authentication?

The purpose is to validate the local stack works, not to test auth flows. Auth is already tested by the existing Cognito integration tests and the main application. Keeping hello routes public simplifies the validation.
