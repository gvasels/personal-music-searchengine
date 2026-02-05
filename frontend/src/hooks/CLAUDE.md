# Hooks - CLAUDE.md

## Overview

Custom React hooks providing TanStack Query wrappers for API calls, authentication state, feature flags, role simulation, and UI utilities. Each hook follows the query key factory pattern for cache management.

## File Descriptions

| File | Purpose |
|------|---------|
| `useAuth.ts` | Authentication state, sign-in/out, role-based access (`UseAuthReturn`) |
| `useHelloSearch.ts` | Hello demo search and featured track queries |
| `useTracks.ts` | Track CRUD queries and mutations (list, detail, update, delete, visibility) |
| `useAlbums.ts` | Album list and detail queries |
| `useArtists.ts` | Legacy artist list and detail queries (aggregated from tracks) |
| `useArtistProfiles.ts` | Artist profile entity CRUD (create, update, delete, search) |
| `useFollow.ts` | Follow/unfollow artists, check following status, list followers |
| `usePlaylists.ts` | Playlist CRUD with optimistic reorder updates |
| `useTags.ts` | Tag list and tracks-by-tag queries |
| `useSearch.ts` | Full-text search and autocomplete queries |
| `useUpload.ts` | File upload with S3 presigned URLs and processing progress polling |
| `useFeatureFlags.ts` | Feature flags with role-based gating, respects simulation |
| `useRoleSimulation.ts` | Admin role simulation for testing UI as different roles |
| `useAdmin.ts` | Admin user management (search users, update role/status) |
| `useCrates.ts` | DJ crate CRUD with feature gating |
| `useHotCues.ts` | Hot cue management with feature gating |
| `useWaveform.ts` | Waveform data fetching (currently disabled, pending backend) |
| `useKeyboardShortcuts.ts` | Global keyboard shortcuts for player controls |
| `useSelection.ts` | Re-export of SelectionContext for batch operations |

## Query Key Factories

Each hook exports a key factory for cache management:

| Factory | Keys | File |
|---------|------|------|
| `helloKeys` | `all`, `search(query)`, `featured()` | `useHelloSearch.ts` |
| `trackKeys` | `all`, `lists()`, `list(params)`, `details()`, `detail(id)` | `useTracks.ts` |
| `albumKeys` | `all`, `lists()`, `list(params)`, `details()`, `detail(id)` | `useAlbums.ts` |
| `artistKeys` | `all`, `lists()`, `list(params)`, `details()`, `detail(name)` | `useArtists.ts` |
| `artistProfileKeys` | `all`, `lists()`, `list(params)`, `details()`, `detail(id)`, `search(params)` | `useArtistProfiles.ts` |
| `followKeys` | `all`, `isFollowing(id)`, `followers(id, params)`, `following(params)` | `useFollow.ts` |
| `playlistKeys` | `all`, `lists()`, `list(params)`, `details()`, `detail(id)` | `usePlaylists.ts` |
| `tagKeys` | `all`, `lists()`, `list(params)`, `tracks(name)`, `trackList(name, params)` | `useTags.ts` |
| `searchKeys` | `all`, `results(params)`, `autocomplete(query)` | `useSearch.ts` |
| `featureKeys` | `all`, `user()` | `useFeatureFlags.ts` |
| `adminKeys` | `all`, `users()`, `userSearch(params)`, `userDetails()`, `userDetail(id)` | `useAdmin.ts` |
| `crateKeys` | `all`, `lists()`, `list(filters)`, `details()`, `detail(id)` | `useCrates.ts` |
| `hotCueKeys` | `all`, `track(trackId)` | `useHotCues.ts` |

## Dependencies

### Internal
- `lib/api/*` - API client functions
- `lib/store/*` - Zustand stores (player, featureFlag, roleSimulation)
- `types` - TypeScript interfaces (`UserRole`, `Permission`, `FeatureKey`, etc.)

### External
- `@tanstack/react-query` - `useQuery`, `useMutation`, `useQueryClient`
- `aws-amplify/auth` - Cognito session (via `lib/auth`)

## Usage Example

```tsx
import { useHelloSearch } from '../hooks/useHelloSearch';
import { useAuth } from '../hooks/useAuth';

function MyComponent() {
  const { user, isAdmin, can } = useAuth();
  const { data, isLoading } = useHelloSearch('jazz');

  if (can('upload_tracks')) {
    // show upload button
  }
}
```
