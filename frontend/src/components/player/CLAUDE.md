# Player Components - CLAUDE.md

## Overview

Audio playback components with Howler.js integration for streaming music playback, waveform visualization, and equalizer controls.

## Files

| File | Description |
|------|-------------|
| `PlayerBar.tsx` | Fixed bottom bar with playback controls, waveform, and EQ |
| `Waveform.tsx` | Canvas-based waveform visualization with seek |
| `Equalizer.tsx` | 3-band EQ with presets and custom settings |
| `index.ts` | Barrel export for player components |
| `__tests__/PlayerBar.test.tsx` | Unit tests for PlayerBar |

## Components

### PlayerBar

```typescript
export function PlayerBar(): JSX.Element
```

Fixed bottom player bar with:
- Track info (title, artist, cover art)
- Playback controls (previous, play/pause, next)
- Waveform visualization with seek functionality
- Volume slider
- Repeat and shuffle toggles
- Compact EQ toggle with expandable panel

**Player Store Integration:**
- `currentTrack` - Currently playing track
- `isPlaying` - Playback state
- `volume` - Volume level (0-1)
- `progress` - Current position in seconds
- `repeat` - Repeat mode ('off' | 'all' | 'one')
- `shuffle` - Shuffle enabled flag

### Waveform

```typescript
interface WaveformProps {
  trackId: string | undefined;
  duration: number;
  progress: number;
  onSeek: (position: number) => void;
  height?: number;
  className?: string;
}

export function Waveform(props: WaveformProps): JSX.Element
```

Canvas-based waveform display with:
- Real waveform data fetching via `useWaveform` hook
- Mock waveform generation for tracks without data
- Progress overlay showing played portion
- Click/touch to seek
- Hover tooltip showing time position
- ResizeObserver for responsive width
- Device pixel ratio handling for crisp rendering

**Colors:**
- Played portion: Primary color (`hsl(var(--p))`)
- Unplayed portion: Base content with opacity

### Equalizer

```typescript
interface EqualizerProps {
  compact?: boolean;
  className?: string;
}

export function Equalizer(props: EqualizerProps): JSX.Element
```

3-band equalizer with:
- Bass (100Hz lowshelf)
- Mid (1kHz peaking)
- Treble (8kHz highshelf)
- Range: -12 to +12 dB per band

**Modes:**
- `compact={true}` - Inline toggle and preset selector for PlayerBar
- `compact={false}` - Full panel with sliders and visual representation

**Presets:**
Built-in: Flat, Bass Boost, Bass Reducer, Treble Boost, Treble Reducer, Vocal, Rock, Electronic, Jazz, Classical

Custom presets can be saved and deleted via localStorage.

## State Management

### EQ Store (`@/lib/store/eqStore`)

```typescript
interface EQState {
  enabled: boolean;
  bands: { bass: number; mid: number; treble: number };
  selectedPreset: string | null;
  customPresets: EQPreset[];

  setEnabled: (enabled: boolean) => void;
  setBand: (band: keyof EQBands, value: number) => void;
  applyPreset: (presetName: string) => void;
  saveCustomPreset: (name: string) => void;
  deleteCustomPreset: (name: string) => void;
  reset: () => void;
}
```

Persisted to localStorage under key `eq-storage`.

## Hooks

### useWaveform (`@/hooks/useWaveform`)

```typescript
export function useWaveform(trackId: string | undefined): UseQueryResult<WaveformData | null>
```

Fetches waveform data from `/tracks/{trackId}/waveform` API endpoint.
- `staleTime: Infinity` - Waveform data doesn't change
- `gcTime: 30 minutes` - Keep in cache

### generateMockWaveform

```typescript
export function generateMockWaveform(duration: number): WaveformData
```

Generates realistic-looking mock waveform for tracks without backend waveform data.

## Dependencies

| Package | Usage |
|---------|-------|
| `@/lib/store/playerStore` | Playback state and actions |
| `@/lib/store/eqStore` | EQ state and presets |
| `@/hooks/useWaveform` | Waveform data fetching |
| `@tanstack/react-query` | Data fetching and caching |
| `zustand` | State management |

## Usage

```typescript
import { PlayerBar, Waveform, Equalizer } from '@/components/player';

// PlayerBar is typically rendered in Layout.tsx
<PlayerBar />

// Standalone waveform (e.g., in track detail page)
<Waveform
  trackId="track-123"
  duration={180}
  progress={45}
  onSeek={(pos) => seek(pos)}
  height={64}
/>

// Standalone equalizer panel
<Equalizer />

// Compact EQ for inline use
<Equalizer compact />
```
