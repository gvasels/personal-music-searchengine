# Hooks - CLAUDE.md

## Overview

Custom React hooks for the frontend application. Most hooks wrap TanStack Query for server state management, following a consistent pattern of query key factories and typed return values. All hooks that fetch data use `useQuery` and expose the standard TanStack Query return shape (data, isLoading, isError, etc.).

## Files

| File | Purpose |
|------|---------|
| `useAuth.ts` | Authentication state, sign in/out, role-based access control |
| `useHelloSearch.ts` | Hello-world search and featured tracks queries |
| `useTracks.ts` | Track CRUD operations with query/mutation hooks |
| `useAlbums.ts` | Album listing and detail queries |
| `useArtists.ts` | Legacy artist listing and detail queries (aggregated from tracks) |
| `useArtistProfiles.ts` | Artist profile CRUD with search (entity-based) |
| `useFollow.ts` | Follow/unfollow artists with toggle hook |
| `useSearch.ts` | Full-text search and autocomplete queries |
| `usePlaylists.ts` | Playlist CRUD with track management and reordering |
| `useTags.ts` | Tag listing and tracks-by-tag queries |
| `useUpload.ts` | File upload with presigned URLs, S3 progress, and processing polling |
| `useAdmin.ts` | Admin user search, details, role/status management |
| `useFeatureFlags.ts` | Feature flags with role-based access and simulation support |
| `useRoleSimulation.ts` | Admin role simulation for testing UI as different roles |
| `useKeyboardShortcuts.ts` | Global keyboard shortcuts for player controls |
| `useSelection.ts` | Re-exports SelectionContext for batch operations |
| `useWaveform.ts` | Waveform data fetching (currently disabled, placeholder) |
| `useCrates.ts` | DJ crate CRUD with feature-gated access |
| `useHotCues.ts` | Hot cue management with feature-gated access |

## Key Functions/Exports

### useAuth
```typescript
export function useAuth(): UseAuthReturn
// Returns: user, isLoading, isAuthenticated, isSigningIn, isSigningOut,
//          signIn(email, password), signOut(), error, clearError, refetch,
//          role, groups, can(permission), isAdmin, isArtist, isSubscriber
```
Uses `getCurrentUser` query with 5-minute staleTime. Role derived from `user.role` (defaults to `'guest'`). `can(permission)` checks via `hasPermission(role, permission)`.

### useHelloSearch
```typescript
export const helloKeys = {
  all: ['hello'] as const,
  search: (query: string) => ['hello', 'search', query] as const,
  featured: (limit?: number) => ['hello', 'featured', limit?] as const,
}
export function useHelloSearch(query: string)     // enabled when query.length > 0
export function useFeaturedTracks(limit?: number) // always enabled
```

### useTracks
```typescript
export const trackKeys = { all, lists, list(params?), details, detail(id) }
export function useTracksQuery(params?: GetTracksParams)
export function useTrackQuery(id: string | undefined)
export function useUpdateTrack()            // mutation, invalidates lists
export function useDeleteTrack()            // mutation, removes from cache
export function useUpdateTrackVisibility()  // mutation, optimistic update
```

### useAlbums
```typescript
export const albumKeys = { all, lists, list(params?), details, detail(id) }
export function useAlbumsQuery(params?: GetAlbumsParams)
export function useAlbumQuery(id: string | undefined)
```

### useArtists (Legacy)
```typescript
export const artistKeys = { all, lists, list(params?), details, detail(name) }
export function useArtistsQuery(params?: GetArtistsParams)
export function useArtistQuery(name: string | undefined)
```

### useArtistProfiles
```typescript
export const artistProfileKeys = { all, lists, list(params?), details, detail(id), search(params) }
export function useArtistProfile(profileId: string)
export function useArtistProfiles(params?: ListArtistProfilesParams)
export function useSearchArtists(params: SearchArtistsParams)  // enabled when query.length > 0
export function useCreateArtistProfile()
export function useUpdateArtistProfile(profileId: string)
export function useDeleteArtistProfile()
```

### useFollow
```typescript
export const followKeys = { all, isFollowing(artistId), followers(artistId, params?), following(params?) }
export function useIsFollowing(artistId: string)
export function useFollowers(artistId: string, params?)
export function useFollowing(params?)
export function useFollowArtist()      // mutation, optimistic isFollowing update
export function useUnfollowArtist()    // mutation, optimistic isFollowing update
export function useFollowToggle(artistId: string)  // combined hook: isFollowing, toggle, follow, unfollow
```

### useSearch
```typescript
export const searchKeys = { all, results(params), autocomplete(query) }
export function useSearchQuery(params: SearchParams)     // enabled when query.length > 0
export function useAutocompleteQuery(query: string)      // enabled when query.length >= 3
```

### usePlaylists
```typescript
export const playlistKeys = { all, lists, list(params?), details, detail(id) }
export function usePlaylistsQuery(params?)
export function usePlaylistQuery(id?)
export function useCreatePlaylist()
export function useUpdatePlaylist()
export function useDeletePlaylist()
export function useAddTracksToPlaylist()
export function useRemoveTracksFromPlaylist()
export function useReorderPlaylistTracks()    // optimistic reorder with rollback
```

### useTags
```typescript
export const tagKeys = { all, lists, list(params?), tracks(tagName), trackList(tagName, params?) }
export function useTagsQuery(params?: GetTagsParams)
export function useTracksByTagQuery(tagName: string, params?)
```

### useUpload
```typescript
export function useUpload(): UseUploadReturn
// Returns: upload(files), isUploading, progress (0-100), error, uploads (UploadItem[]), reset
```
Handles full upload flow: presigned URL, S3 upload with XHR progress (0-50%), confirm, poll processing status (50-100%).

### useAdmin
```typescript
export const adminKeys = { all, users, userSearch(params), userDetails, userDetail(id) }
export function useSearchUsers(search, options?, debounceMs?)  // debounced, enabled when search.length >= 1
export function useUserDetails(userId?)
export function useUpdateUserRole()
export function useUpdateUserStatus()
```

### useFeatureFlags
```typescript
export const featureKeys = { all: ['features'], user: () => ['features', 'user'] }
export function useFeatureFlags()
// Returns: role, actualRole, features, isLoading, isError, isLoaded,
//          isEnabled(feature), hasRole(minRole), invalidate, refetch, isSimulating
export function useFeatureGate(feature: FeatureKey)
// Returns: isEnabled, isLoading, role, isLocked
```

### useRoleSimulation
```typescript
export function useRoleSimulation(): UseRoleSimulationReturn
// Returns: isSimulating, simulatedRole, effectiveRole, actualRole,
//          startSimulation(role), stopSimulation, canSimulate
```

### useKeyboardShortcuts
```typescript
export function useKeyboardShortcuts(options?): UseKeyboardShortcutsReturn
// Options: enabled?, onShowShortcuts?, onEscape?, onFocusSearch?
// Returns: enabled, setEnabled, shortcuts
export function useIsMac(): boolean
```

### useWaveform
```typescript
export function useWaveform(trackId?: string)       // Currently disabled (enabled: false)
export function generateMockWaveform(duration): WaveformData
```

### useCrates
```typescript
export const crateKeys = { all, lists, list(filters?), details, detail(id) }
export function useCrates()                    // feature-gated on 'CRATES'
export function useCrate(id: string)
export function useCreateCrate()
export function useUpdateCrate()
export function useDeleteCrate()
export function useAddTracksToCrate()
export function useRemoveTracksFromCrate()
export function useReorderCrateTracks()
```

### useHotCues
```typescript
export const hotCueKeys = { all, track(trackId) }
export function useHotCues(trackId?: string)   // feature-gated on 'HOT_CUES'
export function useSetHotCue()
export function useDeleteHotCue()
export function useClearHotCues()
```

## Dependencies

- **Internal**: `@/lib/api/*` (API clients), `@/lib/store/*` (Zustand stores), `@/lib/auth` (auth functions), `@/types` (TypeScript types)
- **External**: `@tanstack/react-query` (useQuery, useMutation, useQueryClient), React hooks (useState, useCallback, useEffect, useMemo, useRef)

## Usage Examples

### Query Key Factory Pattern
All hooks follow the same key factory pattern for cache management:
```typescript
export const trackKeys = {
  all: ['tracks'] as const,
  lists: () => [...trackKeys.all, 'list'] as const,
  list: (params?) => [...trackKeys.lists(), params] as const,
  details: () => [...trackKeys.all, 'detail'] as const,
  detail: (id: string) => [...trackKeys.details(), id] as const,
};
```

### Using a Query Hook
```typescript
import { useTracksQuery } from '@/hooks/useTracks';

function TrackListPage() {
  const { data, isLoading, isError } = useTracksQuery({ limit: 50 });
  if (isLoading) return <Loading />;
  return <TrackList tracks={data.items} />;
}
```

### Using a Mutation Hook
```typescript
import { useDeleteTrack } from '@/hooks/useTracks';

function DeleteButton({ trackId }) {
  const deleteMutation = useDeleteTrack();
  return <button onClick={() => deleteMutation.mutate(trackId)}>Delete</button>;
}
```
