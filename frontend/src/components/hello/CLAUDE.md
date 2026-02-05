# Hello Components - CLAUDE.md

## Overview

UI components for the hello-world music search validation feature. These are lightweight components used on the `/hello-search` route to verify end-to-end connectivity between the frontend, backend API, and DynamoDB seed data.

## Files

| File | Purpose |
|------|---------|
| `SearchInput.tsx` | Controlled text input for search queries |
| `TrackCard.tsx` | Displays a single track's metadata (title, artist, album, genre, year, duration) |
| `TrackCardSkeleton.tsx` | Loading placeholder skeleton for TrackCard |
| `HelloNav.tsx` | Navigation link to the `/hello-search` route |
| `__tests__/SearchInput.test.tsx` | Tests for SearchInput (render, onChange, value display) |
| `__tests__/TrackCard.test.tsx` | Tests for TrackCard (renders all fields, formats duration) |
| `__tests__/TrackCardSkeleton.test.tsx` | Tests for TrackCardSkeleton (renders 6 skeleton lines) |

## Key Functions/Exports

### SearchInput

```typescript
interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string; // default: 'Search tracks...'
}
export function SearchInput(props: SearchInputProps): JSX.Element
```

Renders an `input` with DaisyUI `input input-bordered` classes. Calls `onChange` with the string value (not the event).

### TrackCard

```typescript
interface HelloTrack {
  id: string;
  title: string;
  artist: string;
  album: string;
  genre: string;
  year: number;
  duration: number; // in seconds
}
export function TrackCard({ track }: { track: HelloTrack }): JSX.Element
```

Renders a DaisyUI card (`bg-base-200`) with:
- Title styled with `card-title text-primary`
- Artist and album as plain text
- Genre as a `badge`
- Year as plain text
- Duration formatted as `M:SS` via internal `formatDuration(seconds)` helper

### TrackCardSkeleton

```typescript
export function TrackCardSkeleton(): JSX.Element
```

Renders 6 skeleton placeholder lines (`data-testid="skeleton-line"`) matching the TrackCard layout (title, artist, album, genre, year, duration).

### HelloNav

```typescript
export function HelloNav(): JSX.Element
```

Renders a `<nav>` element with an anchor link (`btn btn-ghost`) pointing to `/hello-search`.

## Dependencies

- **Internal**: None (self-contained components)
- **External**: React (JSX)
- **Styling**: DaisyUI 5 semantic classes (`card`, `badge`, `skeleton`, `btn-ghost`, `input-bordered`)

## Usage Examples

```tsx
import { SearchInput } from '@/components/hello/SearchInput';
import { TrackCard } from '@/components/hello/TrackCard';
import { TrackCardSkeleton } from '@/components/hello/TrackCardSkeleton';

// Search input with controlled state
const [query, setQuery] = useState('');
<SearchInput value={query} onChange={setQuery} />

// Render a track card
<TrackCard track={{ id: '1', title: 'Song', artist: 'Artist', album: 'Album', genre: 'Rock', year: 2024, duration: 234 }} />

// Loading state
{isLoading && [...Array(6)].map((_, i) => <TrackCardSkeleton key={i} />)}
```
