# frontend/src/hooks/

## Purpose
Custom React hooks wrapping TanStack Query for API data fetching and state management.

## Files
- `useAdmin.ts` — Admin operations (user search, role updates)
- `useAlbums.ts` — Album CRUD operations
- `useArtistProfiles.ts` — Artist profile management
- `useArtists.ts` — Artist listing and details
- `useAuth.ts` — Authentication state and Cognito integration
- `useCrates.ts` — Crate management
- `useFeatureFlags.ts` — Feature flag checks
- `useFollow.ts` — Artist follow/unfollow
- `useHotCues.ts` — Hot cue management
- `useKeyboardShortcuts.ts` — Global keyboard shortcuts
- `usePlaylists.ts` — Playlist CRUD operations
- `useRoleSimulation.ts` — Dev role simulation
- `useSearch.ts` — Search queries
- `useSelection.ts` — Multi-select state
- `useSimilarTracks.ts` — Similar tracks via vector similarity search
- `useTags.ts` — Tag management
- `useTracks.ts` — Track CRUD operations
- `useUpload.ts` — File upload with progress
- `useWaveform.ts` — Waveform data fetching
- `useArtistWatch.ts` — Artist watch toggle and list (`watchKeys` factory, `useWatchStatus`, `useWatchToggle`, `useWatchedArtists`)
- `useArtistEvents.ts` — Artist events and search (`eventKeys` factory, `useArtistEvents`, `useSearchArtistEvents`)

## Tests
- `__tests__/` — Unit tests for most hooks (12 test files)
