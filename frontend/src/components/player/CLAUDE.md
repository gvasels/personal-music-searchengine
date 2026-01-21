# Player Components - CLAUDE.md

## Overview

Audio playback components with Howler.js integration for streaming music playback.

## Files

| File | Description |
|------|-------------|
| `PlayerBar.tsx` | Fixed bottom bar with playback controls and progress |
| `index.ts` | Barrel export for player components |
| `__tests__/PlayerBar.test.tsx` | Unit tests for PlayerBar (5 tests, 95.53% coverage) |

## Key Functions

### PlayerBar.tsx
```typescript
export function PlayerBar(): JSX.Element
```
Renders a fixed bottom player bar with:
- Track info (title, artist)
- Playback controls (previous, play/pause, next)
- Progress bar with seek functionality
- Volume slider
- Repeat and shuffle toggles

**Player Store Integration:**
- `currentTrack` - Currently playing track
- `isPlaying` - Playback state
- `volume` - Volume level (0-1)
- `progress` - Current position in seconds
- `repeat` - Repeat mode ('off' | 'all' | 'one')
- `shuffle` - Shuffle enabled flag

**Actions:**
- `play()` / `pause()` - Toggle playback
- `next()` / `previous()` - Skip tracks
- `seek(seconds)` - Jump to position
- `setVolume(level)` - Adjust volume
- `cycleRepeat()` - Cycle through repeat modes
- `toggleShuffle()` - Toggle shuffle mode

## Dependencies

| Package | Usage |
|---------|-------|
| `@/lib/store/playerStore` | All playback state and actions |
| `@/types` | `Track` interface |

## Usage

```typescript
import { PlayerBar } from '@/components/player';

// Typically rendered in Layout.tsx
<PlayerBar />
```

The PlayerBar automatically shows/hides based on `currentTrack` presence.
