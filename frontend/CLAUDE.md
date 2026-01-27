# Frontend - CLAUDE.md

## Overview

React frontend for the Personal Music Search Engine. Built with Vite, TypeScript, TanStack Router, TanStack Query, and DaisyUI 5 for styling. Uses Zustand for state management and Howler.js for audio playback.

## Directory Structure

```
frontend/
├── src/
│   ├── components/         # Reusable React components
│   │   ├── layout/         # Layout components (Header, Sidebar, Layout)
│   │   ├── library/        # Music library (TrackList, AlbumCard, ArtistCard)
│   │   ├── player/         # Audio player (PlayerBar)
│   │   ├── playlist/       # Playlist management (PlaylistCard, CreatePlaylistModal)
│   │   ├── search/         # Search functionality (SearchBar)
│   │   ├── tag/            # Tag management (TagInput)
│   │   └── upload/         # File upload (UploadDropzone)
│   ├── hooks/              # Custom React hooks (useAuth, useTracks, useAlbums, etc.)
│   ├── lib/                # Utilities and configurations
│   │   ├── api/            # API client and types
│   │   └── store/          # Zustand stores (player, theme)
│   └── routes/             # TanStack Router file-based routes
├── public/                 # Static assets
├── index.html              # HTML entry point
├── package.json            # Dependencies
├── vite.config.ts          # Vite configuration
├── tailwind.config.js      # Tailwind + DaisyUI themes
└── tsconfig.json           # TypeScript configuration
```

## Components

### Layout (`src/components/layout/`)
| Component | File | Description |
|-----------|------|-------------|
| `Layout` | `Layout.tsx` | App shell with sidebar, header, and PlayerBar |
| `Header` | `Header.tsx` | Top navigation with SearchBar and theme toggle |
| `Sidebar` | `Sidebar.tsx` | Navigation menu with mobile hamburger support |

### Library (`src/components/library/`)
| Component | File | Description |
|-----------|------|-------------|
| `TrackList` | `TrackList.tsx` | Track table with click-to-play integration |
| `TrackRow` | `TrackRow.tsx` | Individual track row component |
| `AlbumCard` | `AlbumCard.tsx` | Album card with cover art |
| `AlbumGrid` | `AlbumGrid.tsx` | Grid layout for albums |
| `ArtistCard` | `ArtistCard.tsx` | Artist card component |

### Player (`src/components/player/`)
| Component | File | Description |
|-----------|------|-------------|
| `PlayerBar` | `PlayerBar.tsx` | Fixed player bar with Howler.js integration |

### Playlist (`src/components/playlist/`)
| Component | File | Description |
|-----------|------|-------------|
| `PlaylistCard` | `PlaylistCard.tsx` | Playlist card with track count |
| `CreatePlaylistModal` | `CreatePlaylistModal.tsx` | Modal for creating new playlists |

### Search (`src/components/search/`)
| Component | File | Description |
|-----------|------|-------------|
| `SearchBar` | `SearchBar.tsx` | Search with autocomplete dropdown |

### Tag (`src/components/tag/`)
| Component | File | Description |
|-----------|------|-------------|
| `TagInput` | `TagInput.tsx` | Tag add/remove component with badges |

### Upload (`src/components/upload/`)
| Component | File | Description |
|-----------|------|-------------|
| `UploadDropzone` | `UploadDropzone.tsx` | Drag-and-drop file upload with progress |

## Hooks (`src/hooks/`)

| Hook | File | Purpose |
|------|------|---------|
| `useAuth` | `useAuth.ts` | Authentication state and actions (signIn, signOut, currentUser) |
| `useTracks` | `useTracks.ts` | Track CRUD queries with `trackKeys` factory |
| `useAlbums` | `useAlbums.ts` | Album queries with `albumKeys` factory |
| `useArtists` | `useArtists.ts` | Artist queries with `artistKeys` factory |
| `useUpload` | `useUpload.ts` | File upload mutations with progress tracking |
| `useSearch` | `useSearch.ts` | Search queries and autocomplete with `searchKeys` factory |
| `usePlaylists` | `usePlaylists.ts` | Playlist CRUD queries with `playlistKeys` factory |
| `useTags` | `useTags.ts` | Tag queries and tracks-by-tag with `tagKeys` factory |
| `useFeatureFlags` | `useFeatureFlags.ts` | Feature flags with role-based access, respects simulation |
| `useRoleSimulation` | `useRoleSimulation.ts` | Admin role simulation for testing UI as different roles |

## Utilities (`src/lib/`)

### API Modules (`api/`)

| File | Functions |
|------|-----------|
| `client.ts` | `apiClient` (Axios with Cognito auth interceptor) |
| `tracks.ts` | `getTracks`, `getTrack`, `updateTrack`, `deleteTrack` |
| `albums.ts` | `getAlbums`, `getAlbum` |
| `artists.ts` | `getArtists`, `getArtist` |
| `upload.ts` | `getPresignedUploadUrl`, `confirmUpload`, `getUploadStatus` |
| `search.ts` | `searchTracks`, `searchAutocomplete` |
| `playlists.ts` | `getPlaylists`, `getPlaylist`, `createPlaylist`, `updatePlaylist`, `deletePlaylist`, `addTrackToPlaylist`, `removeTrackFromPlaylist` |
| `tags.ts` | `getTags`, `getTracksByTag`, `addTagToTrack`, `removeTagFromTrack` |

### Types (`api/types.ts`)
```typescript
export interface Track { id, title, artist, album, duration, format, ... }
export interface Album { id, name, artist, year, trackCount, ... }
export interface Artist { id, name, trackCount, albumCount }
export interface Playlist { id, name, description, trackIds, trackCount }
export interface Tag { name, trackCount }
export interface Upload { id, status, filename, progress, ... }
export interface PaginatedResponse<T> { items, total, limit, offset }
```

### Auth (`auth.ts`, `amplify.ts`)
```typescript
// AWS Amplify configuration for Cognito
export const amplifyConfig: ResourcesConfig
export const configureAuth: () => void
```

### State Management (`store/`)
```typescript
// Player store with Howler.js (playerStore.ts)
export const usePlayerStore: {
  currentTrack, queue, isPlaying, volume, progress,
  play, pause, next, previous, seek, setVolume, setQueue
}

// Theme store with persistence (themeStore.ts)
export const useThemeStore: {
  theme, toggleTheme
}

// Role simulation store with persistence (roleSimulationStore.ts)
export const useRoleSimulationStore: {
  simulatedRole, startedAt,
  setSimulatedRole, clearSimulation
}
```

## Routes (`src/routes/`)

TanStack Router file-based routing:

| Route | File | Description |
|-------|------|-------------|
| `/` | `index.tsx` | Home page with library stats |
| `/__root` | `__root.tsx` | Root layout with auth guard |
| `/login` | `login.tsx` | Login page with Cognito auth |
| `/tracks` | `tracks/index.tsx` | Track listing |
| `/tracks/$trackId` | `tracks/$trackId.tsx` | Track detail with TagInput |
| `/albums` | `albums/index.tsx` | Album grid view |
| `/albums/$albumId` | `albums/$albumId.tsx` | Album detail with tracks |
| `/artists` | `artists/index.tsx` | Artist listing |
| `/artists/$artistId` | `artists/$artistId.tsx` | Artist detail |
| `/playlists` | `playlists/index.tsx` | Playlist grid |
| `/playlists/$playlistId` | `playlists/$playlistId.tsx` | Playlist detail |
| `/tags` | `tags/index.tsx` | Tag cloud |
| `/tags/$tagName` | `tags/$tagName.tsx` | Tracks by tag |
| `/upload` | `upload.tsx` | File upload page |
| `/search` | `search.tsx` | Search results page |

## DaisyUI Themes

### Dark Theme
- Base: `#120612` (deep purple-black)
- Primary: `#50c878` (emerald green)
- Secondary: `#72001c` (deep crimson)

### Light Theme
- Base: `#fdfdf8` (warm white)
- Primary: `#50c878` (emerald green)
- Secondary: `#72001c` (deep crimson)

## Dependencies

### Production
| Package | Purpose |
|---------|---------|
| `react` | UI framework |
| `@tanstack/react-router` | Type-safe routing |
| `@tanstack/react-query` | Server state management |
| `zustand` | Client state management |
| `axios` | HTTP client |
| `aws-amplify` | Cognito authentication |
| `howler` | Audio playback |
| `daisyui` | UI components (v5) |
| `tailwindcss` | Utility-first CSS (v4) |
| `react-dropzone` | Drag-and-drop file upload |
| `react-hot-toast` | Toast notifications |
| `clsx` | Classname utilities |

### Development
| Package | Purpose |
|---------|---------|
| `vite` | Build tool |
| `typescript` | Type checking |
| `vitest` | Testing |
| `@tanstack/router-vite-plugin` | Route generation |

## Testing

### Framework
- **Vitest** with React Testing Library
- **jsdom** environment for DOM simulation
- **@testing-library/user-event** for user interaction simulation

### Test Coverage (351 tests)
| Area | Tests |
|------|-------|
| API tests (client, tracks, albums, artists, upload, search, playlists, tags) | 50 |
| Hook tests (useAuth, useTracks, useAlbums, useArtists, useUpload, useSearch, usePlaylists, useTags) | 92 |
| Store tests (playerStore, themeStore) | 18 |
| Component tests (Layout, TrackList, PlayerBar, UploadDropzone, CreatePlaylistModal, TagInput) | 34 |
| Route tests (index, login, tracks, trackDetail, albums, artists, upload, search, playlists, playlistDetail, tags, tagDetail) | 157 |

### Running Tests
```bash
# Run all tests
npm test

# Run with coverage
npm run test -- --coverage

# Run in watch mode
npm run test -- --watch
```

### Test Utilities (`src/test/`)
- `setup.ts` - Global mocks (matchMedia, localStorage, etc.)
- `test-utils.tsx` - Custom render with QueryClient wrapper

## Build Commands

```bash
# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview

# Run tests
npm test

# Run tests with coverage
npm test -- --coverage

# Type check
npm run typecheck

# Lint
npm run lint
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VITE_API_URL` | API Gateway URL |
| `VITE_COGNITO_USER_POOL_ID` | Cognito User Pool ID |
| `VITE_COGNITO_CLIENT_ID` | Cognito App Client ID |
| `VITE_COGNITO_REGION` | AWS Region (us-east-1) |

## Usage Examples

### Using the Player Store
```typescript
import { usePlayerStore } from '@/lib/store'

function PlayButton({ track }) {
  const { setQueue } = usePlayerStore()

  return (
    <button onClick={() => setQueue([track], 0)}>
      Play
    </button>
  )
}
```

### Fetching Tracks with TanStack Query
```typescript
import { useQuery } from '@tanstack/react-query'
import { getTracks } from '@/lib/api/client'

function TrackListPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['tracks'],
    queryFn: () => getTracks({ limit: 50 })
  })

  if (isLoading) return <div>Loading...</div>
  return <TrackList tracks={data.items} />
}
```

### Adding Tags to a Track
```typescript
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { addTagToTrack } from '@/lib/api/client'

function TagButton({ trackId, tagName }) {
  const queryClient = useQueryClient()
  const mutation = useMutation({
    mutationFn: () => addTagToTrack(trackId, tagName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['track', trackId] })
    }
  })

  return <button onClick={() => mutation.mutate()}>Add Tag</button>
}
```
