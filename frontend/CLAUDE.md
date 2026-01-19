# Frontend - CLAUDE.md

## Overview

React frontend for the Personal Music Search Engine. Built with Vite, TypeScript, TanStack Router, TanStack Query, and DaisyUI for styling. Uses Zustand for state management and Howler.js for audio playback.

## Directory Structure

```
frontend/
├── src/
│   ├── components/         # Reusable React components
│   ├── hooks/              # Custom React hooks
│   ├── lib/                # Utilities and configurations
│   ├── pages/              # Page components (if not using routes)
│   └── routes/             # TanStack Router file-based routes
├── public/                 # Static assets
├── index.html              # HTML entry point
├── package.json            # Dependencies
├── vite.config.ts          # Vite configuration
├── tailwind.config.js      # Tailwind + DaisyUI themes
└── tsconfig.json           # TypeScript configuration
```

## File Descriptions

| File | Purpose |
|------|---------|
| `vite.config.ts` | Vite build configuration with TanStack Router plugin |
| `tailwind.config.js` | Tailwind CSS config with custom dark/light DaisyUI themes |
| `tsconfig.json` | TypeScript compiler options |
| `package.json` | Dependencies and scripts |
| `src/main.tsx` | Application entry point |
| `src/index.css` | Global styles and DaisyUI theme imports |

## Key Components

### Layout Components (`src/components/`)
| Component | File | Description |
|-----------|------|-------------|
| `Layout` | `Layout.tsx` | App shell with sidebar, header, and player |
| `Header` | `Header.tsx` | Top navigation with search and theme toggle |
| `Sidebar` | `Sidebar.tsx` | Navigation menu |
| `Player` | `Player.tsx` | Audio player with Howler.js |

### Library Components (`src/components/`)
| Component | File | Description |
|-----------|------|-------------|
| `TrackList` | `TrackList.tsx` | Track listing with sorting/filtering |
| `AlbumGrid` | `AlbumGrid.tsx` | Album grid view |

## Utilities (`src/lib/`)

### API Client (`api.ts`)
```typescript
// Axios client with Cognito auth interceptor
export const apiClient: AxiosInstance

// Type definitions
export interface Track { ... }
export interface Album { ... }
export interface Playlist { ... }
export interface Upload { ... }

// API functions
export const getTracks: (params?) => Promise<PaginatedResponse<Track>>
export const getTrack: (id: string) => Promise<Track>
export const getPresignedUploadUrl: (data) => Promise<PresignedUploadResponse>
export const confirmUpload: (uploadId: string) => Promise<ConfirmUploadResponse>
export const getStreamUrl: (trackId: string) => Promise<StreamUrlResponse>
```

### Auth (`auth.ts`)
```typescript
// AWS Amplify configuration for Cognito
export const configureAuth: () => void
export const getCurrentUser: () => Promise<User>
export const signIn: (email: string, password: string) => Promise<void>
export const signOut: () => Promise<void>
```

### State Management (`store.ts`)
```typescript
// Zustand stores
export const usePlayerStore: StoreApi<PlayerState & PlayerActions>
export const useThemeStore: StoreApi<ThemeState>
export const useSidebarStore: StoreApi<SidebarState>
```

## Routes (`src/routes/`)

TanStack Router file-based routing:

| Route | File | Description |
|-------|------|-------------|
| `/` | `index.tsx` | Home/library view |
| `/__root` | `__root.tsx` | Root layout wrapper |

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
| `daisyui` | UI components |
| `tailwindcss` | Utility-first CSS |

### Development
| Package | Purpose |
|---------|---------|
| `vite` | Build tool |
| `typescript` | Type checking |
| `vitest` | Testing |
| `@testing-library/react` | Component testing |

## Build Commands

```bash
# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview

# Run tests
npm run test

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

## Testing

Component tests with Vitest and Testing Library:

```typescript
import { render, screen } from '@testing-library/react'
import { TrackList } from './TrackList'

describe('TrackList', () => {
  it('renders tracks', () => {
    render(<TrackList tracks={mockTracks} />)
    expect(screen.getByText('Track Title')).toBeInTheDocument()
  })
})
```

## Usage Examples

### Using the Player Store
```typescript
import { usePlayerStore } from '@/lib/store'

function PlayButton({ track }) {
  const { setQueue, play } = usePlayerStore()

  return (
    <button onClick={() => {
      setQueue([track], 0)
      play()
    }}>
      Play
    </button>
  )
}
```

### Fetching Tracks
```typescript
import { useQuery } from '@tanstack/react-query'
import { getTracks } from '@/lib/api'

function TrackListPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['tracks'],
    queryFn: () => getTracks()
  })

  if (isLoading) return <div>Loading...</div>
  return <TrackList tracks={data.items} />
}
```
