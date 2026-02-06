# Hello Components - CLAUDE.md

## Overview

React components for the Hello World search feature. Demonstrates full-stack integration with the LocalStack-seeded music data.

## File Descriptions

| File | Purpose |
|------|---------|
| `HelloNav.tsx` | Navigation bar with branding and home link |
| `SearchInput.tsx` | Controlled text input for search functionality |
| `TrackCard.tsx` | Displays track information in a card format |
| `TrackCardSkeleton.tsx` | Loading skeleton placeholder for TrackCard |

## Components

### HelloNav

Navigation bar for the Hello Music Search page.

```typescript
export function HelloNav(): JSX.Element
```

**Features:**
- Brand text "Hello Music Search"
- Home link button

**DaisyUI classes:** `navbar`, `btn btn-ghost`

---

### SearchInput

Controlled text input for search functionality.

```typescript
interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;  // default: "Search tracks..."
}

export function SearchInput(props: SearchInputProps): JSX.Element
```

**Features:**
- Full-width bordered input
- Controlled component pattern
- Customizable placeholder

**DaisyUI classes:** `input input-bordered w-full`

---

### TrackCard

Displays track information in a card format.

```typescript
interface TrackCardProps {
  track: HelloTrack;  // from lib/api/helloSearch.ts
}

export function TrackCard({ track }: TrackCardProps): JSX.Element
```

**Features:**
- Displays title (card-title), artist, album
- Shows genre badge, year, formatted duration
- Duration formatted as "M:SS" (e.g., 245 seconds -> "4:05")

**Internal function:**
```typescript
const formatDuration = (seconds: number): string => {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${String(secs).padStart(2, '0')}`;
};
```

**DaisyUI classes:** `card bg-base-200`, `card-body`, `card-title`, `badge`

---

### TrackCardSkeleton

Loading skeleton placeholder matching TrackCard layout.

```typescript
export function TrackCardSkeleton(): JSX.Element
```

**Features:**
- Mimics TrackCard structure with skeleton placeholders
- Shows 4 skeleton lines matching title, artist, album, metadata

**DaisyUI classes:** `skeleton`, `card bg-base-200`

## Dependencies

### Internal
- `lib/api/helloSearch.ts` - `HelloTrack` type

### External
- React (JSX)
- DaisyUI 5 component classes

## Usage Example

```tsx
import { HelloNav } from '@/components/hello/HelloNav';
import { SearchInput } from '@/components/hello/SearchInput';
import { TrackCard } from '@/components/hello/TrackCard';
import { TrackCardSkeleton } from '@/components/hello/TrackCardSkeleton';

function HelloPage() {
  const [query, setQuery] = useState('');
  const { data, isLoading } = useHelloSearch(query);

  return (
    <>
      <HelloNav />
      <SearchInput value={query} onChange={setQuery} />
      {isLoading && <TrackCardSkeleton />}
      {data?.items.map(track => <TrackCard key={track.id} track={track} />)}
    </>
  );
}
```
