# Design Document: Hello World Local Development Validation

## Overview

Add a standalone hello-world music search feature to validate the full local development stack end-to-end: LocalStack DynamoDB seed data, Go backend API, and React frontend search page. The feature is completely isolated from the main application and requires no authentication. This is iteration 3, incorporating API contract fixes and BUILD phase improvements from iterations 1 and 2.

## Existing Infrastructure Analysis

### Components to Reuse

| Component | Location | How Used |
|-----------|----------|----------|
| Echo router + middleware | `cmd/api/main.go` | Register hello routes alongside existing routes |
| Services struct | `service/service.go` | Add `Hello *HelloService` field |
| `ListResponse[T]` generic | `handlers/handlers.go` | Uses `items`/`total` keys -- matches our API contract |
| DynamoDB repository | `repository/dynamodb.go` | HelloService queries existing table directly |
| Vite proxy | `vite.config.ts` | `/api` proxy to `http://localhost:8080` already configured |
| TanStack Router | `routeTree.gen.ts` | File-based routing auto-generates route tree |
| TanStack Query | hooks pattern | `useHelloSearch` follows existing hook patterns |
| DaisyUI 5 components | throughout frontend | TrackCard uses existing card/badge patterns |
| LocalStack init scripts | `docker/localstack-init/` | Seed script added alongside existing init scripts |
| Makefile | root `Makefile` | Add seed script call to `local-services` target |

### New Components to Create

| Component | Location | Purpose |
|-----------|----------|---------|
| Seed script | `docker/localstack-init/init-seed-music.sh` | Insert 20 tracks into DynamoDB |
| HelloService | `backend/internal/service/hello.go` | Search and featured track business logic |
| HelloHandler | `backend/internal/handlers/hello.go` | HTTP handler for hello endpoints |
| Frontend API client | `frontend/src/lib/api/helloSearch.ts` | API functions for hello endpoints |
| useHelloSearch hook | `frontend/src/hooks/useHelloSearch.ts` | TanStack Query hook for search |
| SearchInput | `frontend/src/components/hello/SearchInput.tsx` | Search input component |
| TrackCard | `frontend/src/components/hello/TrackCard.tsx` | Track display card |
| TrackCardSkeleton | `frontend/src/components/hello/TrackCardSkeleton.tsx` | Loading skeleton |
| HelloNav | `frontend/src/components/hello/HelloNav.tsx` | Navigation bar for hello page |
| Route page | `frontend/src/routes/hello-search.tsx` | Main search page route |

## Architecture

### Request Flow

```
Browser (localhost:5173)
  │
  │  GET /api/v1/hello/search?q=jazz
  │
  ▼
Vite Dev Server (proxy /api → :8080)
  │
  ▼
Go Echo Server (localhost:8080)
  │
  ├── /api/v1/hello/health    → HelloHandler.Health()
  ├── /api/v1/hello/search    → HelloHandler.Search()    → HelloService.Search(q)
  └── /api/v1/hello/featured  → HelloHandler.Featured()  → HelloService.Featured(limit)
                                                                │
                                                                ▼
                                                          DynamoDB (LocalStack :4566)
                                                          Table: MusicLibrary
                                                          PK: USER#seed-user
                                                          SK: TRACK#{trackId}
```

## API Contract (CRITICAL)

**Response Format:** `{"items": [...], "total": N}` -- uses the `items` key (NOT `tracks`).
**Search Query Parameter:** `q` -- the query param name is `q` (NOT `query`).

This matches the existing `ListResponse[T]` generic in `handlers/handlers.go`:
```go
type ListResponse[T any] struct {
    Items []T `json:"items"`
    Total int `json:"total"`
}
```

### Endpoint Specifications

| Endpoint | Method | Query Params | Response Shape |
|----------|--------|--------------|----------------|
| `/api/v1/hello/health` | GET | none | `{"status":"ok","service":"hello"}` |
| `/api/v1/hello/search` | GET | `q` (string) | `{"items": HelloTrack[], "total": number}` |
| `/api/v1/hello/featured` | GET | `limit` (optional int) | `{"items": HelloTrack[], "total": number}` |

### HelloTrack Schema

```go
// Backend
type HelloTrack struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Artist   string `json:"artist"`
    Album    string `json:"album"`
    Genre    string `json:"genre"`
    Year     int    `json:"year"`
    Duration int    `json:"duration"` // seconds
}
```

```typescript
// Frontend
export interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number; // seconds
}
```

### Frontend API Client Contract

```typescript
// MUST match backend exactly:
export interface HelloSearchResponse {
  items: HelloTrack[];  // NOT "tracks"
  total: number;
}

// Search MUST use "q" param:
export async function searchHelloTracks(query: string): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/search', {
    params: { q: query },  // NOT { query }
  });
  return response.data;
}

// Featured endpoint:
export async function getFeaturedHelloTracks(limit?: number): Promise<HelloSearchResponse> {
  const response = await apiClient.get<HelloSearchResponse>('/v1/hello/featured', {
    params: limit ? { limit } : undefined,
  });
  return response.data;
}
```

## Components and Interfaces

### Component 1: Seed Script (`docker/localstack-init/init-seed-music.sh`)

- **Purpose:** Insert 20 seed tracks into DynamoDB for the hello-world feature.
- **Interface:**
  ```bash
  # Called by: make local-services (after init-aws.sh)
  # Idempotent: uses PutItem (overwrites existing)
  # Data: 5 artists x 4 tracks = 20 tracks
  # Schema: PK=USER#seed-user, SK=TRACK#{uuid}
  ```
- **Seed Data:**
  - 5 artists (e.g., "Aurora Waves", "Midnight Echo", "Solar Drift", "Neon Pulse", "Velvet Storm")
  - 4 tracks per artist with varied genres, years, durations
  - All under `USER#seed-user` to keep isolated from real user data
- **Dependencies:** `init-aws.sh` must run first (creates the DynamoDB table)
- **Idempotency:** Uses DynamoDB `PutItem` which overwrites existing items with same PK/SK

### Component 2: HelloService (`backend/internal/service/hello.go`)

- **Purpose:** Business logic for searching and listing hello tracks from DynamoDB.
- **Interface:**
  ```go
  // HelloService defines the hello feature operations
  type HelloService struct {
      repo repository.Repository
  }

  func NewHelloService(repo repository.Repository) *HelloService

  // Search returns tracks matching the query (case-insensitive across title, artist, album, genre)
  // Empty query returns empty result
  func (s *HelloService) Search(ctx context.Context, query string) ([]HelloTrack, error)

  // Featured returns up to limit seed tracks (default 20)
  func (s *HelloService) Featured(ctx context.Context, limit int) ([]HelloTrack, error)
  ```
- **Implementation Details:**
  - Queries DynamoDB with `PK = USER#seed-user` and `SK begins_with TRACK#`
  - Search: scans results and filters by case-insensitive substring match
  - Featured: returns all results up to the limit
  - Converts DynamoDB items to `HelloTrack` structs
- **Dependencies:** `repository.Repository` (existing DynamoDB interface)

### Component 3: HelloRepoAdapter (inline in HelloService)

- **Purpose:** Bridge between existing repository interface and hello service needs.
- **Details:** The HelloService uses the existing `repository.Repository` interface to query DynamoDB directly. No new repository interface is needed -- the service uses `Query` with the `USER#seed-user` partition key.
- **Why no new repo interface:** The hello feature is intentionally simple. Adding a dedicated repository interface would be over-engineering for a validation feature.

### Component 4: HelloHandler (`backend/internal/handlers/hello.go`)

- **Purpose:** HTTP handler for hello endpoints, wired to Echo router.
- **Interface:**
  ```go
  type HelloHandler struct {
      service *service.HelloService
  }

  func NewHelloHandler(svc *service.HelloService) *HelloHandler

  // Health returns service health status
  // GET /api/v1/hello/health
  func (h *HelloHandler) Health(c echo.Context) error

  // Search returns tracks matching query param "q"
  // GET /api/v1/hello/search?q={query}
  func (h *HelloHandler) Search(c echo.Context) error

  // Featured returns all seed tracks
  // GET /api/v1/hello/featured?limit={N}
  func (h *HelloHandler) Featured(c echo.Context) error

  // RegisterHelloRoutes registers hello routes on the Echo instance
  func RegisterHelloRoutes(e *echo.Echo, h *HelloHandler)
  ```
- **Route Registration:**
  ```go
  func RegisterHelloRoutes(e *echo.Echo, h *HelloHandler) {
      hello := e.Group("/api/v1/hello")
      hello.GET("/health", h.Health)
      hello.GET("/search", h.Search)    // uses c.QueryParam("q")
      hello.GET("/featured", h.Featured)
  }
  ```
- **Key Implementation Details:**
  - Search handler reads `c.QueryParam("q")` (NOT `c.QueryParam("query")`)
  - Returns `ListResponse[HelloTrack]{Items: tracks, Total: len(tracks)}` using existing generic
  - No auth middleware on the hello group
- **Dependencies:** `service.HelloService`, `echo.Echo`

### Component 5: Wiring (`cmd/api/main.go` + `service/service.go`)

All 4 wiring steps required by the wiring checklist:

**Step 1: Add to Services struct** (`service/service.go`):
```go
type Services struct {
    // ...existing fields...
    Hello *HelloService  // Hello world validation service
}
```

**Step 2: Initialize in NewServices()** (`service/service.go`):
```go
func NewServices(repo, s3Repo, cloudfront, mediaBucket, sfnARN) *Services {
    return &Services{
        // ...existing services...
        Hello: NewHelloService(repo),
    }
}
```

**Step 3: Create HelloHandler in main.go** (`cmd/api/main.go`):
```go
// After existing handler creation
helloHandler := handlers.NewHelloHandler(services.Hello)
```

**Step 4: Call RegisterHelloRoutes in main.go** (`cmd/api/main.go`):
```go
// After existing route registration
handlers.RegisterHelloRoutes(e, helloHandler)
```

### Component 6: Frontend API Client (`frontend/src/lib/api/helloSearch.ts`)

- **Purpose:** API functions for hello endpoints.
- **Interface:**
  ```typescript
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
    items: HelloTrack[];  // MUST be "items", NOT "tracks"
    total: number;
  }

  // MUST use "q" param, NOT "query"
  export async function searchHelloTracks(query: string): Promise<HelloSearchResponse>;
  export async function getFeaturedHelloTracks(limit?: number): Promise<HelloSearchResponse>;
  export async function getHelloHealth(): Promise<{ status: string; service: string }>;
  ```
- **Dependencies:** `apiClient` from `./client.ts`

### Component 7: useHelloSearch Hook (`frontend/src/hooks/useHelloSearch.ts`)

- **Purpose:** TanStack Query hook for hello search with debouncing.
- **Interface:**
  ```typescript
  export const helloKeys = {
    all: ['hello'] as const,
    search: (query: string) => [...helloKeys.all, 'search', query] as const,
    featured: () => [...helloKeys.all, 'featured'] as const,
  };

  export function useHelloSearch(query: string): UseQueryResult<HelloSearchResponse>;
  export function useHelloFeatured(): UseQueryResult<HelloSearchResponse>;
  ```
- **Dependencies:** `@tanstack/react-query`, `helloSearch.ts` API functions

### Component 8: SearchInput (`frontend/src/components/hello/SearchInput.tsx`)

- **Purpose:** Controlled search input with DaisyUI styling.
- **Interface:**
  ```typescript
  interface SearchInputProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
  }

  export function SearchInput(props: SearchInputProps): JSX.Element;
  ```
- **Dependencies:** DaisyUI input classes

### Component 9: TrackCard + TrackCardSkeleton

**TrackCard** (`frontend/src/components/hello/TrackCard.tsx`):
- **Purpose:** Display a single track with title, artist, album, genre, year, duration.
- **Interface:**
  ```typescript
  interface TrackCardProps {
    track: HelloTrack;
  }

  export function TrackCard(props: TrackCardProps): JSX.Element;
  ```
- **Details:** Duration formatted as `MM:SS` using `Math.floor(duration / 60)` and `duration % 60`

**TrackCardSkeleton** (`frontend/src/components/hello/TrackCardSkeleton.tsx`):
- **Purpose:** Loading placeholder for track cards.
- **Interface:**
  ```typescript
  export function TrackCardSkeleton(): JSX.Element;
  ```

### Component 10: HelloNav (`frontend/src/components/hello/HelloNav.tsx`)

- **Purpose:** Navigation bar for the hello-search page.
- **Interface:**
  ```typescript
  export function HelloNav(): JSX.Element;
  ```
- **Details:** Simple nav with "Hello Music Search" title and link back to home

### Component 11: Route Page (`frontend/src/routes/hello-search.tsx`)

- **Purpose:** Main search page combining all hello components.
- **Interface:**
  ```typescript
  // TanStack Router file-based route at /hello-search
  export const Route = createFileRoute('/hello-search')({
    component: HelloSearchPage,
  });
  ```
- **Behavior:**
  1. On initial load: fetch and display featured tracks
  2. When user types in search input: debounce and search
  3. Show TrackCardSkeleton during loading
  4. Show error alert on API failure
  5. Show "No results found" when search returns empty
- **Public Route:** Must add `/hello-search` to `PUBLIC_ROUTES` array in `__root.tsx`
- **Route Tree:** Must run `npx @tanstack/router-cli generate` after creating this file

## Data Flow

```
User types "jazz" in SearchInput
  │
  ▼
HelloSearchPage: setQuery("jazz")
  │
  ▼
useHelloSearch("jazz") → queryKey: ['hello', 'search', 'jazz']
  │
  ▼
searchHelloTracks("jazz")
  │
  ▼
GET /api/v1/hello/search?q=jazz    ← uses "q" param
  │
  ▼
HelloHandler.Search(c)
  │ q := c.QueryParam("q")         ← reads "q" param
  ▼
HelloService.Search(ctx, "jazz")
  │
  ▼
DynamoDB Query: PK=USER#seed-user, SK begins_with TRACK#
  │ Filter: case-insensitive match on title/artist/album/genre
  ▼
Response: {"items": [...], "total": 3}  ← uses "items" key
  │
  ▼
Frontend receives HelloSearchResponse.items  ← reads "items" key
  │
  ▼
Render TrackCard for each track
```

## Error Handling

| Scenario | Backend Response | Frontend Display |
|----------|-----------------|------------------|
| Empty query | `{"items": [], "total": 0}` | "No results" message |
| No matches | `{"items": [], "total": 0}` | "No results" message |
| DynamoDB error | 500 Internal Server Error | Error alert component |
| Network error | N/A | Error alert component |
| LocalStack down | 500 (connection refused) | Error alert component |

## Testing Strategy

### Backend Tests
- **Unit tests for HelloService:** Mock repository, test search filtering logic, empty query handling, featured with limit
- **Unit tests for HelloHandler:** Mock service, test query param parsing (`q`), response format (`items`), status codes

### Frontend Tests
- **API client tests:** Mock axios, verify `q` param sent (not `query`), verify `items` key parsed (not `tracks`)
- **Hook tests:** Mock API functions, test query key factory, loading/error/success states
- **Component tests:** SearchInput onChange, TrackCard rendering with formatted duration, TrackCardSkeleton structure
- **Route tests:** Page renders, search interaction, featured tracks display, loading states, error display

### Contract Alignment Verification
After both frontend and backend are implemented, verify:
```bash
# Frontend sends "q" param
grep -n "params.*{.*q:" frontend/src/lib/api/helloSearch.ts
# Frontend expects "items" key
grep -n "items:" frontend/src/lib/api/helloSearch.ts
# Backend reads "q" query param
grep -n 'QueryParam("q")' backend/internal/handlers/hello.go
# Backend returns "items" key
grep -n '"items"' backend/internal/handlers/hello.go
```
