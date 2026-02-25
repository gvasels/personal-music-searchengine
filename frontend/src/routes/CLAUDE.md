# Routes - CLAUDE.md

## Overview

TanStack Router file-based routes for the Personal Music Search Engine. Routes are auto-generated from file structure into `routeTree.gen.ts`.

## Route Files

| File | Route | Description |
|------|-------|-------------|
| `__root.tsx` | - | Root layout with auth guard and guest route protection |
| `index.tsx` | `/` | Home page with library stats |
| `login.tsx` | `/login` | Login page with Cognito auth |
| `permission-denied.tsx` | `/permission-denied` | Access denied page |
| `hello.tsx` | `/hello` | Hello World validation page |
| `hello-search.tsx` | `/hello-search` | Hello World full-stack search demo |
| `shows.tsx` | `/shows` | My Shows page - watched artist events and artist search |
| `search.tsx` | `/search` | Search results page |
| `upload.tsx` | `/upload` | File upload page |

### Subdirectory Routes

| Directory | Routes | Description |
|-----------|--------|-------------|
| `tracks/` | `/tracks`, `/tracks/$trackId` | Track listing and detail |
| `albums/` | `/albums`, `/albums/$albumId` | Album grid and detail |
| `artists/` | `/artists/*` | Artist pages |
| `playlists/` | `/playlists`, `/playlists/$playlistId` | Playlist management |
| `tags/` | `/tags`, `/tags/$tagName` | Tag cloud and tracks by tag |

## hello-search.tsx

Full-stack demonstration page for the Hello World local dev feature.

### Route Definition

```typescript
export const Route = createFileRoute('/hello-search')({
  component: HelloSearchPage,
});
```

### Features

- Displays featured tracks on initial load (`useHelloFeatured`)
- Supports real-time search (`useHelloSearch`)
- Shows loading skeletons during fetch
- Displays error alert on failure
- Shows "No results found" for empty search results

### Components Used

- `HelloNav` - Navigation bar
- `SearchInput` - Controlled search input
- `TrackCard` - Track display cards
- `TrackCardSkeleton` - Loading placeholders

### State Management

```typescript
const [query, setQuery] = useState('');
const searchResult = useHelloSearch(query);
const featuredResult = useHelloFeatured();

// Display search results when searching, otherwise featured
const isSearching = query.length > 0 || searchResult.isLoading || ...;
const data = isSearching ? searchResult.data : featuredResult.data;
```

### Layout

- Full-height background (`min-h-screen bg-base-100`)
- Navigation bar at top
- Container with search input
- Responsive grid: 1 column (mobile), 2 columns (md), 3 columns (lg)

## Route Protection

The `__root.tsx` implements route protection:

**Public Routes** (no auth required):
- `/` - Dashboard
- `/login` - Login page
- `/permission-denied` - Access denied
- `/hello-search` - Hello search demo

**Protected Routes** (redirect to `/permission-denied` if not authenticated):
- All other routes

## Dependencies

### Internal
- `hooks/useHelloSearch.ts` - Search and featured hooks
- `components/hello/*` - Hello components

### External
- `@tanstack/react-router` - Routing framework
- React hooks (`useState`)
