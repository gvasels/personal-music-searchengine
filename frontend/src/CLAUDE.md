# Frontend Source - CLAUDE.md

## Overview

Source code for the React frontend application. Contains components, utilities, routes, and hooks organized by feature.

## Directory Structure

```
src/
├── components/
│   ├── layout/       # App shell: Header, Sidebar, Layout
│   ├── library/      # Music: TrackList, TrackRow, AlbumCard, AlbumGrid, ArtistCard
│   ├── player/       # Audio: PlayerBar
│   ├── playlist/     # Playlists: PlaylistCard, CreatePlaylistModal
│   ├── search/       # Search: SearchBar
│   ├── tag/          # Tags: TagInput
│   └── upload/       # Upload: UploadDropzone
├── hooks/
│   └── useAuth.ts    # Authentication hook
├── lib/
│   ├── api/          # API client and types
│   │   ├── client.ts # Axios client with auth interceptor
│   │   ├── types.ts  # TypeScript interfaces
│   │   └── index.ts  # Barrel export
│   ├── store/        # Zustand stores
│   │   ├── playerStore.ts  # Audio player state
│   │   ├── themeStore.ts   # Theme persistence
│   │   └── index.ts        # Barrel export
│   ├── amplify.ts    # AWS Amplify config
│   └── auth.ts       # Auth configuration
├── routes/           # TanStack Router file-based routes
│   ├── __root.tsx    # Root layout with auth guard
│   ├── index.tsx     # Home page
│   ├── login.tsx     # Login page
│   ├── search.tsx    # Search results
│   ├── upload.tsx    # File upload
│   ├── tracks/       # /tracks routes
│   ├── albums/       # /albums routes
│   ├── artists/      # /artists routes
│   ├── playlists/    # /playlists routes
│   └── tags/         # /tags routes
├── main.tsx          # Application entry point
├── index.css         # Global styles (Tailwind + DaisyUI)
├── vite-env.d.ts     # Vite type declarations
└── routeTree.gen.ts  # Generated route tree (auto-generated)
```

## Entry Point (`main.tsx`)

```typescript
import { createRouter, RouterProvider } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import { routeTree } from './routeTree.gen'
import './lib/amplify'  // Configure Cognito

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 1000 * 60 * 5 }  // 5 min cache
  }
})
const router = createRouter({ routeTree })

ReactDOM.createRoot(document.getElementById('root')!).render(
  <QueryClientProvider client={queryClient}>
    <RouterProvider router={router} />
    <Toaster position="bottom-right" />
  </QueryClientProvider>
)
```

## Key Files

### API Client (`lib/api/client.ts`)
- Axios instance with Cognito JWT interceptor
- All API functions: getTracks, getAlbums, getPlaylists, etc.
- Presigned upload URL handling
- Stream URL generation

### Types (`lib/api/types.ts`)
- Track, Album, Artist, Playlist, Tag, Upload interfaces
- PaginatedResponse<T> for list endpoints
- SearchResponse for search results

### Player Store (`lib/store/playerStore.ts`)
- Howler.js audio instance management
- Queue management (setQueue, next, previous)
- Playback controls (play, pause, seek, setVolume)
- State: currentTrack, queue, isPlaying, volume, progress

### Theme Store (`lib/store/themeStore.ts`)
- Dark/light theme toggle
- localStorage persistence
- System preference detection

## Testing

Component tests should be placed alongside components:
- `Component.tsx` → `Component.test.tsx`

Run tests: `npm run test`

## Route Pattern

Each route folder follows this pattern:
```
routes/tracks/
├── index.tsx       # /tracks - List view
└── $trackId.tsx    # /tracks/:trackId - Detail view
```

TanStack Router generates the route tree automatically from file structure.
