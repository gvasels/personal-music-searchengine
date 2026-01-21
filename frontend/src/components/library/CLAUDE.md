# Library Components - CLAUDE.md

## Overview

Music library display components for rendering track lists with playback integration.

## Files

| File | Description |
|------|-------------|
| `TrackList.tsx` | Table component displaying tracks with click-to-play |
| `index.ts` | Barrel export for library components |
| `__tests__/TrackList.test.tsx` | Unit tests for TrackList (6 tests, 98.27% coverage) |

## Key Functions

### TrackList.tsx
```typescript
interface TrackListProps {
  tracks: Track[];
  isLoading?: boolean;
}

export function TrackList({ tracks, isLoading }: TrackListProps): JSX.Element
```
Renders a table of tracks with columns:
- Index number (or play indicator for current track)
- Title
- Artist
- Album
- Duration (formatted as mm:ss)

**States:**
- Loading: Shows skeleton placeholders
- Empty: Shows "No tracks found" message
- Populated: Shows track table

```typescript
function formatDuration(seconds: number): string
```
Converts seconds to `m:ss` format (e.g., 180 → "3:00").

## Dependencies

| Package | Usage |
|---------|-------|
| `@/types` | `Track` interface |
| `@/lib/store/playerStore` | `setQueue`, `currentTrack`, `isPlaying` |

## Usage

```typescript
import { TrackList } from '@/components/library';

<TrackList tracks={tracks} isLoading={isLoading} />
```

Clicking a track row calls `setQueue(tracks, index)` to start playback.
