# API Client Modules - CLAUDE.md

## Overview

TypeScript modules that provide typed API client functions for communicating with the backend. All modules use a shared Axios instance (`apiClient`) configured with Cognito JWT auth interceptor and structured error handling. Each module corresponds to a backend domain.

## Files

| File | Purpose |
|------|---------|
| `client.ts` | Shared Axios instance with auth interceptor, error handling, and legacy API functions |
| `index.ts` | Barrel export (re-exports `client.ts`) |
| `helloSearch.ts` | Hello-world search, featured tracks, and health check endpoints |
| `tracks.ts` | Track CRUD operations and visibility management |
| `albums.ts` | Album listing and detail with tracks |
| `artists.ts` | Legacy artist listing/detail (aggregated from tracks) |
| `artistProfiles.ts` | Artist profile entity CRUD and search |
| `follows.ts` | Follow/unfollow artists, followers/following lists |
| `playlists.ts` | Playlist CRUD, track management, reordering, visibility, public playlists |
| `tags.ts` | Tag listing, tracks-by-tag, add/remove tag operations |
| `search.ts` | Full-text search and autocomplete |
| `upload.ts` | Presigned URL generation, upload confirmation, status polling |
| `features.ts` | Feature flags, subscription, tier config, storage usage |
| `admin.ts` | Admin user search, details, role/status management |
| `stats.ts` | Library statistics with role-based scope |
| `crates.ts` | DJ crate CRUD with track management |
| `hotcues.ts` | Hot cue CRUD per track |

## Key Functions/Exports

### client.ts (Core)

```typescript
export const apiClient: AxiosInstance
// baseURL: VITE_API_URL || '/api'
// Request interceptor: attaches Cognito Bearer token from fetchAuthSession()
// Response interceptor: converts errors to ApiError instances

export class ApiError extends Error {
  code: string;
  statusCode: number;
  details?: unknown;
}
```

Also exports convenience functions for many endpoints (legacy pattern):
`getTracks`, `getTrack`, `getAlbums`, `getArtists`, `getArtistEntities`, `getArtistEntity`, `createArtist`, `updateArtist`, `deleteArtist`, `searchArtistEntities`, `getArtistEntityTracks`, `getPlaylists`, `getPlaylist`, `createPlaylist`, `deletePlaylist`, `addTrackToPlaylist`, `removeTrackFromPlaylist`, `reorderPlaylistTracks`, `addTagToTrack`, `removeTagFromTrack`, `searchTracks`, `getPresignedUploadUrl`, `getStreamUrl`, `getDownloadUrl`, `deleteTrack`

### helloSearch.ts

```typescript
export interface HelloTrack { id, title, artist, album, genre, year, duration }
export interface HelloSearchResponse { items: HelloTrack[], total: number }
export async function searchHelloTracks(query: string): Promise<HelloSearchResponse>   // GET /v1/hello/search?q=
export async function getFeaturedTracks(limit?: number): Promise<HelloSearchResponse>  // GET /v1/hello/featured
export async function getHelloHealth(): Promise<{ status: string; service: string }>   // GET /v1/hello/health
```

### tracks.ts

```typescript
export interface GetTracksParams { page?, limit?, sortBy?, sortOrder?, search?, artist?, album?, tags? }
export interface UpdateTrackData { title?, artist?, album?, tags? }
export async function getTracks(params?): Promise<PaginatedResponse<Track>>
export async function getTrack(id): Promise<Track>
export async function updateTrack(id, data): Promise<Track>       // PATCH /tracks/:id
export async function deleteTrack(id): Promise<void>
export async function updateTrackVisibility(id, visibility): Promise<{ trackId, visibility }>
```

### albums.ts

```typescript
export interface GetAlbumsParams { page?, limit?, sortBy?, sortOrder?, artist? }
export interface AlbumWithTracks extends Album { tracks: Track[] }
export async function getAlbums(params?): Promise<PaginatedResponse<Album>>
export async function getAlbum(id): Promise<AlbumWithTracks>
```

### artists.ts (Legacy)

```typescript
export interface GetArtistsParams { page?, limit?, sortBy?, sortOrder?, search? }
export interface ArtistWithDetails extends ArtistSummary { albums, recentTracks }
export async function getArtists(params?): Promise<PaginatedResponse<ArtistSummary>>
export async function getArtist(name): Promise<ArtistWithDetails>
```

### artistProfiles.ts

```typescript
export interface CreateArtistProfileData { displayName, bio?, location?, website?, socialLinks? }
export interface UpdateArtistProfileData { displayName?, bio?, avatarUrl?, headerImageUrl?, location?, website?, socialLinks? }
export async function createArtistProfile(data): Promise<ArtistProfile>
export async function getArtistProfile(profileId): Promise<ArtistProfile>
export async function updateArtistProfile(profileId, data): Promise<ArtistProfile>
export async function deleteArtistProfile(profileId): Promise<void>
export async function listArtistProfiles(params?): Promise<ArtistProfilesResponse>
export async function searchArtists(params): Promise<ArtistProfilesResponse>
```

### follows.ts

```typescript
export async function followArtist(artistId): Promise<void>
export async function unfollowArtist(artistId): Promise<void>
export async function isFollowing(artistId): Promise<boolean>
export async function getFollowers(artistId, params?): Promise<FollowersResponse>
export async function getFollowing(params?): Promise<FollowingResponse>
```

### playlists.ts

```typescript
export async function getPlaylists(params?): Promise<PaginatedResponse<Playlist>>
export async function getPlaylist(id): Promise<PlaylistWithTracks>
export async function createPlaylist(data): Promise<Playlist>
export async function updatePlaylist(id, data): Promise<Playlist>
export async function deletePlaylist(id): Promise<void>
export async function addTrackToPlaylist(playlistId, trackId): Promise<Playlist>
export async function removeTrackFromPlaylist(playlistId, trackId): Promise<Playlist>
export async function reorderPlaylistTracks(playlistId, data): Promise<Playlist>
export async function updatePlaylistVisibility(playlistId, visibility): Promise<UpdateVisibilityResponse>
export async function getPublicPlaylists(params?): Promise<PublicPlaylistsResponse>
```

### tags.ts

```typescript
export async function getTags(params?): Promise<PaginatedResponse<Tag>>
export async function getTracksByTag(tagName, params?): Promise<PaginatedResponse<Track>>
export async function addTagToTrack(trackId, tagName): Promise<Track>
export async function removeTagFromTrack(trackId, tagName): Promise<Track>
```

### search.ts

```typescript
export async function searchTracks(params: SearchParams): Promise<SearchResponse>
export async function searchAutocomplete(query): Promise<AutocompleteResponse>
```

### upload.ts

```typescript
export async function getPresignedUploadUrl(data): Promise<PresignedUploadResponse>
export async function confirmUpload(uploadId): Promise<UploadConfirmResponse>
export async function getUploadStatus(uploadId): Promise<UploadStatusResponse>
```

### features.ts

```typescript
export async function getUserFeatures(): Promise<UserFeaturesResponse>
export async function getSubscription(): Promise<SubscriptionResponse>
export async function getTierConfigs(): Promise<TierConfig[]>
export async function getStorageUsage(): Promise<StorageUsageResponse>
export async function createCheckoutSession(tier, interval, successUrl?, cancelUrl?): Promise<{ checkoutUrl, sessionId }>
export async function createPortalSession(returnUrl?): Promise<{ portalUrl }>
```

### admin.ts

```typescript
export async function searchUsers(params): Promise<SearchUsersResponse>
export async function getUserDetails(userId): Promise<UserDetails>
export async function updateUserRole(userId, role): Promise<UserDetails>
export async function updateUserStatus(userId, disabled): Promise<UserDetails>
```

### stats.ts

```typescript
export async function getLibraryStats(scope?: StatsScope): Promise<LibraryStats>
// scope: 'own' | 'public' | 'all'
```

### crates.ts

```typescript
export async function getCrates(): Promise<Crate[]>
export async function getCrate(id): Promise<CrateWithTracks>
export async function createCrate(data): Promise<Crate>
export async function updateCrate(id, data): Promise<Crate>
export async function deleteCrate(id): Promise<void>
export async function addTracksToCrate(crateId, trackIds, position?): Promise<void>
export async function removeTracksFromCrate(crateId, trackIds): Promise<void>
export async function reorderCrateTracks(crateId, trackIds): Promise<void>
export async function getCrateTracks(crateId): Promise<Track[]>
```

### hotcues.ts

```typescript
export async function getTrackHotCues(trackId): Promise<TrackHotCuesResponse>
export async function setHotCue(trackId, slot, data): Promise<HotCue>
export async function deleteHotCue(trackId, slot): Promise<void>
export async function clearHotCues(trackId): Promise<void>
```

## Dependencies

- **Internal**: `@/types` (TypeScript interfaces for Track, Album, Artist, Playlist, etc.)
- **External**: `axios` (HTTP client), `aws-amplify/auth` (fetchAuthSession for JWT tokens)

## Usage Examples

```typescript
// Direct API call
import { searchHelloTracks } from '@/lib/api/helloSearch';
const results = await searchHelloTracks('neon');

// Using the shared client for custom requests
import { apiClient } from '@/lib/api/client';
const response = await apiClient.get('/custom-endpoint');

// Error handling
import { ApiError } from '@/lib/api/client';
try {
  await deleteTrack(id);
} catch (err) {
  if (err instanceof ApiError && err.statusCode === 403) {
    // Handle forbidden
  }
}
```
