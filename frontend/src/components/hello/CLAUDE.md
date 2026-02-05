# Hello Components - CLAUDE.md

## Overview

React components for the hello-world music search demo page. Provides a navigation bar, search input, track display cards, and loading skeletons using DaisyUI 5 semantic classes.

## File Descriptions

| File | Purpose |
|------|---------|
| `HelloNav.tsx` | Top navigation bar with app title and home link |
| `SearchInput.tsx` | Controlled text input for search queries |
| `TrackCard.tsx` | Displays a single track with title, artist, album, genre, year, and formatted duration |
| `TrackCardSkeleton.tsx` | Skeleton loading placeholder matching TrackCard layout |

## Key Exports

| Export | File | Signature |
|--------|------|-----------|
| `HelloNav` | `HelloNav.tsx` | `() => JSX.Element` |
| `SearchInput` | `SearchInput.tsx` | `({ value, onChange, placeholder? }) => JSX.Element` |
| `TrackCard` | `TrackCard.tsx` | `({ track: HelloTrack }) => JSX.Element` |
| `TrackCardSkeleton` | `TrackCardSkeleton.tsx` | `() => JSX.Element` |

### SearchInput Props

```typescript
interface SearchInputProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string;
}
```

### TrackCard Props

```typescript
interface TrackCardProps {
  track: HelloTrack;  // from lib/api/helloSearch
}
```

## Internal Functions

| Function | File | Purpose |
|----------|------|---------|
| `formatDuration(seconds: number): string` | `TrackCard.tsx` | Formats seconds as `M:SS` |

## Dependencies

### Internal
- `lib/api/helloSearch` - `HelloTrack` type (used by TrackCard)

### External
- DaisyUI 5 classes: `navbar`, `card`, `input`, `skeleton`, `bg-base-200`

## Usage Example

```tsx
import { SearchInput } from '../components/hello/SearchInput';
import { TrackCard } from '../components/hello/TrackCard';
import { TrackCardSkeleton } from '../components/hello/TrackCardSkeleton';

function SearchPage() {
  const [query, setQuery] = useState('');
  const { data, isLoading } = useHelloSearch(query);

  return (
    <>
      <SearchInput value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Search..." />
      {isLoading && <TrackCardSkeleton />}
      {data?.items.map((track) => <TrackCard key={track.id} track={track} />)}
    </>
  );
}
```
