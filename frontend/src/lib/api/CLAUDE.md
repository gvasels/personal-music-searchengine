# API Client - CLAUDE.md

## Overview

Axios-based API client modules for communicating with the backend. The shared `apiClient` instance automatically attaches Cognito JWT tokens and transforms error responses into typed `ApiError` objects. Each domain module provides typed request/response functions.

## File Descriptions

| File | Purpose |
|------|---------|
| `client.ts` | Shared Axios instance with auth interceptor, legacy API functions, `ApiError` class |
| `index.ts` | Barrel re-export of `client.ts` |
| `helloSearch.ts` | Hello demo: search, featured tracks, health check |
| `tracks.ts` | Track CRUD, visibility updates |
| `albums.ts` | Album list and detail |
| `artists.ts` | Legacy artist list and detail (aggregated) |
| `artistProfiles.ts` | Artist profile entity CRUD and search |
| `follows.ts` | Follow/unfollow, followers list, following list |
| `playlists.ts` | Playlist CRUD, track add/remove/reorder |
| `tags.ts` | Tag list, tracks-by-tag |
| `search.ts` | Full-text search, autocomplete |
| `upload.ts` | Presigned upload URLs, confirm upload, status polling |
| `admin.ts` | Admin user search, role/status management |
| `features.ts` | Feature flags, tier config, subscription, storage usage |
| `stats.ts` | Library statistics with role-based scope |
| `crates.ts` | DJ crate CRUD, track add/remove/reorder |
| `hotcues.ts` | Hot cue set/delete/clear per track |

## Key Exports

### Shared (`client.ts`)

| Export | Type | Description |
|--------|------|-------------|
| `apiClient` | `AxiosInstance` | Pre-configured Axios with auth interceptor; base URL from `VITE_API_URL` or `/api` |
| `ApiError` | `class` | Typed error with `code`, `message`, `statusCode`, `details` |
| `getTracks(params?)` | `function` | Paginated track list |
| `getTrack(id)` | `function` | Single track by ID |
| `searchTracks(query)` | `function` | Search tracks by query string |
| `getPresignedUploadUrl(data)` | `function` | Get S3 presigned URL for upload |
| `getStreamUrl(trackId)` | `function` | Get streaming URL |
| `deleteTrack(trackId)` | `function` | Delete a track |

### Hello Search (`helloSearch.ts`)

| Export | Signature |
|--------|-----------|
| `searchHelloTracks(query)` | `(query: string) => Promise<HelloSearchResponse>` |
| `getFeaturedHelloTracks(limit?)` | `(limit?: number) => Promise<HelloSearchResponse>` |
| `getHelloHealth()` | `() => Promise<HelloHealthResponse>` |
| `HelloTrack` | Interface: `{ id, title, artist, album, genre?, year?, duration }` |
| `HelloSearchResponse` | Interface: `{ items: HelloTrack[], total: number }` |

## Auth Interceptor

The request interceptor in `client.ts` automatically attaches the Cognito access token:

```typescript
apiClient.interceptors.request.use(async (config) => {
  const session = await fetchAuthSession();
  const token = session.tokens?.accessToken?.toString();
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});
```

## Dependencies

### Internal
- `@/types` - Shared TypeScript interfaces

### External
- `axios` - HTTP client
- `aws-amplify/auth` - `fetchAuthSession` for JWT tokens

## Usage Example

```typescript
import { apiClient } from './client';
import { searchHelloTracks } from './helloSearch';

// Direct client usage
const { data } = await apiClient.get('/v1/hello/health');

// Typed function usage
const results = await searchHelloTracks('jazz');
console.log(results.items); // HelloTrack[]
```
