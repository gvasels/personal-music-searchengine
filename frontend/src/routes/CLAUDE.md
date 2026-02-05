# Routes - CLAUDE.md

## Overview

TanStack Router file-based routes for the frontend application. The root layout enforces authentication, redirecting unauthenticated users from protected routes to `/permission-denied`. Route files export a `Route` object created via `createFileRoute`.

## File Descriptions

| File | Route | Auth | Description |
|------|-------|------|-------------|
| `__root.tsx` | (layout) | -- | Root layout with auth guard, guest detection, role simulation support |
| `index.tsx` | `/` | Public | Home page with library stats, recent tracks, quick actions |
| `login.tsx` | `/login` | Public | Login page with Cognito authentication |
| `signup.tsx` | `/signup` | Public | User registration |
| `permission-denied.tsx` | `/permission-denied` | Public | Access denied page for unauthenticated users |
| `hello-search.tsx` | `/hello-search` | Public | Hello-world demo search page with featured/search tracks |
| `search.tsx` | `/search` | Protected | Full-text search results |
| `upload.tsx` | `/upload` | Protected | File upload with drag-and-drop |
| `settings.tsx` | `/settings` | Protected | User settings |

### Subdirectory Routes

| Directory | Routes | Description |
|-----------|--------|-------------|
| `tracks/` | `/tracks`, `/tracks/$trackId` | Track listing and detail |
| `albums/` | `/albums`, `/albums/$albumId` | Album grid and detail |
| `artists/` | `/artists`, `/artists/$artistName` | Legacy artist listing and detail |
| `artists/entity/` | `/artists/entity`, `/artists/entity/$artistId` | Artist profile entity pages |
| `playlists/` | `/playlists`, `/playlists/$playlistId`, `/playlists/public` | Playlist management and public discovery |
| `tags/` | `/tags`, `/tags/$tagName` | Tag cloud and tracks-by-tag |
| `admin/` | `/admin/users` | Admin user management panel |
| `studio/` | `/studio` | DJ studio interface |
| `__tests__/` | -- | Route unit tests (Vitest) |

## Route Protection

Defined in `__root.tsx`:

```typescript
const PUBLIC_ROUTES = ['/', '/login', '/permission-denied', '/hello-search'];
```

All other routes redirect unauthenticated or guest-role users to `/permission-denied`. Admin role simulation is also respected -- an admin simulating `guest` is treated as unauthenticated.

## Key Functions

| Function | File | Description |
|----------|------|-------------|
| `isPublicRoute(pathname)` | `__root.tsx` | Checks if a path is in the public routes list |
| `RootComponent()` | `__root.tsx` | Root layout with auth loading, guest detection, debug logging |
| `HelloSearchPage()` | `hello-search.tsx` | Search page with featured tracks fallback and loading states |
| `HomePage()` | `index.tsx` | Dashboard with role-aware stats scope (`own`, `public`, `all`) |

## Dependencies

### Internal
- `hooks/useAuth` - Authentication state
- `hooks/useFeatureFlags` - Role simulation awareness
- `hooks/useHelloSearch` - Hello search and featured queries
- `components/hello/*` - SearchInput, TrackCard, TrackCardSkeleton

### External
- `@tanstack/react-router` - `createRootRoute`, `createFileRoute`, `Outlet`, `Navigate`
- `@tanstack/react-query` - Server state (via hooks)

## Usage Example

```tsx
// Creating a new route file (e.g., routes/my-page.tsx)
import { createFileRoute } from '@tanstack/react-router';

function MyPage() {
  return <div>My Page Content</div>;
}

export const Route = createFileRoute('/my-page')({
  component: MyPage,
});
```

To make a route public, add it to `PUBLIC_ROUTES` in `__root.tsx`.
