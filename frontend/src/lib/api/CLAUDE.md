# API Client - CLAUDE.md

## Overview

API client functions for the Personal Music Search Engine frontend. Uses Axios with Cognito JWT interceptor for authenticated requests.

## File Descriptions

| File | Purpose |
|------|---------|
| `client.ts` | Axios client with auth interceptor, base configuration |
| `index.ts` | Barrel export |
| `tracks.ts` | Track CRUD operations |
| `albums.ts` | Album queries |
| `artists.ts` | Artist queries |
| `artistProfiles.ts` | Artist profile CRUD |
| `playlists.ts` | Playlist CRUD and track management |
| `tags.ts` | Tag queries and track-tag associations |
| `search.ts` | Search and autocomplete |
| `upload.ts` | Presigned URL and upload confirmation |
| `follows.ts` | Follow/unfollow operations |
| `features.ts` | Feature flag queries |
| `admin.ts` | Admin user management operations |
| `crates.ts` | DJ crate operations |
| `hotcues.ts` | Hot cue operations |
| `stats.ts` | Library statistics |
| `helloSearch.ts` | Hello World search feature |
| `artistWatch.ts` | Artist watch/unwatch operations |
| `events.ts` | Artist events and search operations |

## helloSearch.ts

API functions for the Hello World local development feature.

### Types

```typescript
export interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number;  // in seconds
}

export interface HelloSearchResponse {
  items: HelloTrack[];
  total: number;
}
```

### Functions

```typescript
// Search tracks by query string
// CRITICAL: Uses param `q` (NOT `query`)
export async function searchHelloTracks(query: string): Promise<HelloSearchResponse>

// Get featured tracks
export async function getFeaturedHelloTracks(limit?: number): Promise<HelloSearchResponse>

// Check service health
export async function getHelloHealth(): Promise<{ status: string; service: string }>
```

### API Endpoints

| Function | Method | Endpoint |
|----------|--------|----------|
| `searchHelloTracks` | GET | `/v1/hello/search?q={query}` |
| `getFeaturedHelloTracks` | GET | `/v1/hello/featured?limit={limit}` |
| `getHelloHealth` | GET | `/v1/hello/health` |

### Contract Critical Notes

- `searchHelloTracks` uses query param `q` (NOT `query`)
- Response uses `items` field (NOT `tracks`)
- These match the backend handler contract in `handlers/hello.go`

### Usage Example

```typescript
import { searchHelloTracks, getFeaturedHelloTracks, HelloTrack } from '@/lib/api/helloSearch';

// Search for tracks
const results = await searchHelloTracks('aurora');
console.log(results.items);  // HelloTrack[]
console.log(results.total);  // number

// Get featured tracks
const featured = await getFeaturedHelloTracks(10);
```

## Dependencies

### Internal
- `client.ts` - `apiClient` Axios instance

### External
- `axios` - HTTP client

## API Base URL

The API client connects to:
- Production: `VITE_API_URL` environment variable
- Local dev: `http://localhost:8080/api` (via Vite proxy)
