# Hooks - CLAUDE.md

## Overview

Custom React hooks for the Personal Music Search Engine. Uses TanStack Query for server state management with query key factories for efficient cache invalidation.

## File Descriptions

| File | Purpose |
|------|---------|
| `useAuth.ts` | Authentication state and actions (signIn, signOut, currentUser) |
| `useTracks.ts` | Track CRUD queries with `trackKeys` factory |
| `useAlbums.ts` | Album queries with `albumKeys` factory |
| `useArtists.ts` | Artist queries with `artistKeys` factory |
| `useArtistProfiles.ts` | Artist profile CRUD with `artistProfileKeys` factory |
| `usePlaylists.ts` | Playlist CRUD with `playlistKeys` factory |
| `useTags.ts` | Tag queries with `tagKeys` factory |
| `useSearch.ts` | Search queries with `searchKeys` factory |
| `useUpload.ts` | File upload mutations with progress tracking |
| `useFollow.ts` | Follow/unfollow mutations with `followKeys` factory |
| `useFeatureFlags.ts` | Role-based feature flags, respects simulation |
| `useRoleSimulation.ts` | Admin role simulation for testing |
| `useAdmin.ts` | Admin user management operations |
| `useCrates.ts` | DJ crate management |
| `useHotCues.ts` | Hot cue management for tracks |
| `useWaveform.ts` | Audio waveform generation |
| `useKeyboardShortcuts.ts` | Global keyboard shortcuts |
| `useSelection.ts` | Multi-select state management |
| `useHelloSearch.ts` | Hello World search feature queries |

## Query Key Factory Pattern

All hooks follow the TanStack Query key factory pattern for consistent cache management:

```typescript
export const exampleKeys = {
  all: ['example'] as const,
  lists: () => [...exampleKeys.all, 'list'] as const,
  list: (filters: Filters) => [...exampleKeys.lists(), filters] as const,
  details: () => [...exampleKeys.all, 'detail'] as const,
  detail: (id: string) => [...exampleKeys.details(), id] as const,
};
```

## useHelloSearch

TanStack Query hooks for the Hello World search feature.

### Query Key Factory

```typescript
export const helloKeys = {
  all: ['hello'] as const,
  search: (q: string) => [...helloKeys.all, 'search', q] as const,
  featured: () => [...helloKeys.all, 'featured'] as const,
};
```

### Hooks

```typescript
// Search tracks - disabled when query is empty
export function useHelloSearch(query: string): UseQueryResult<HelloSearchResponse>

// Get featured tracks
export function useHelloFeatured(): UseQueryResult<HelloSearchResponse>
```

### Response Type

```typescript
interface HelloSearchResponse {
  items: HelloTrack[];
  total: number;
}
```

### Usage Example

```tsx
import { useHelloSearch, useHelloFeatured, helloKeys } from '@/hooks/useHelloSearch';

function SearchPage() {
  const [query, setQuery] = useState('');

  // Search is disabled when query is empty
  const { data, isLoading, isError } = useHelloSearch(query);

  // Featured tracks for initial display
  const { data: featured } = useHelloFeatured();

  // Display logic
  const displayData = query ? data : featured;
}
```

## Dependencies

### Internal
- `lib/api/*` - API client functions
- `lib/api/helloSearch.ts` - Hello search API types and functions

### External
- `@tanstack/react-query` - Server state management
- `aws-amplify` - Cognito authentication (useAuth)

## Testing

Hook tests are in `__tests__/` directory:
- Use `@testing-library/react` with `renderHook`
- Mock API functions with `vi.mock`
- Use `QueryClientProvider` wrapper from test utils
