# Design Document: Frontend Enhancements

## Overview

This design covers six frontend enhancements: global search bar, tags column with inline editing, HLS audio playback, play buttons per track, and complete album/artist detail views. The implementation leverages existing React components, TanStack Query hooks, and Zustand stores while adding HLS.js for streaming.

## Code Reuse Analysis

### Existing Components to Leverage
- **TagInput** (`components/tag/TagInput.tsx`): Fully functional tag add/remove with optimistic updates
- **TrackList** (`components/library/TrackList.tsx`): Reusable track table component
- **PlayerBar** (`components/player/PlayerBar.tsx`): Player UI with controls (needs audio connection)
- **Layout/Header** (`components/layout/Header.tsx`): Target for search bar placement

### Existing Hooks to Leverage
- **useSearchQuery/useAutocompleteQuery** (`hooks/useSearch.ts`): Search API integration ready
- **useTagsQuery** (`hooks/useTags.ts`): Tag listing and track-by-tag queries
- **useAlbumQuery** (`hooks/useAlbums.ts`): Fetches album with tracks
- **useArtistQuery** (`hooks/useArtists.ts`): Fetches artist with albums/tracks
- **usePlayerStore** (`lib/store/playerStore.ts`): Zustand store for playback state

### Existing API Functions
- **searchTracks/searchAutocomplete** (`lib/api/search.ts`): Search API client
- **addTagToTrack/removeTagFromTrack** (`lib/api/tags.ts`): Tag mutation APIs
- **getStreamUrl** (`lib/api/client.ts`): Returns signed HLS/S3 URL (unused currently)

### Integration Points
- **Backend API**: All endpoints exist (`/search`, `/tracks/:id/tags`, `/stream/:trackId`, `/albums/:id`, `/artists/:name`)
- **MediaConvert HLS**: Streams at `hls/{userId}/{trackId}/master.m3u8`
- **CloudFront**: Signed URLs for media delivery

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Header                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    SearchBar                             │   │
│  │  [Input] ──debounce──> useAutocompleteQuery             │   │
│  │           └──> SearchDropdown (results grouped)          │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      Tracks Page                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ # │ ▶ │ Title │ Artist │ Album │ Tags        │ Duration │  │
│  │───│───│───────│────────│───────│─────────────│──────────│  │
│  │ 1 │ ▷ │ Song  │ Artist │ Album │ [tag][+2]   │ 3:45     │  │
│  │   │   │       │        │       │ └─TagsCell  │          │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      Audio Service                               │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────────┐   │
│  │ playerStore │◄───│ AudioService │◄───│    HLS.js       │   │
│  │  (Zustand)  │    │  (singleton) │    │  + <audio>      │   │
│  └─────────────┘    └──────────────┘    └─────────────────┘   │
│        │                   │                     │              │
│        │ state updates     │ control methods     │ events       │
│        ▼                   ▼                     ▼              │
│  [PlayerBar] ◄────► [play/pause/seek] ◄────► [timeupdate]     │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. SearchBar Component
**File:** `frontend/src/components/search/SearchBar.tsx`
- **Purpose:** Global search input with live dropdown results
- **Interfaces:**
  ```typescript
  // No props - uses internal state and hooks
  const SearchBar: React.FC = () => { ... }
  ```
- **Dependencies:** useAutocompleteQuery, useNavigate, useDebounce
- **Behavior:**
  - Input with 300ms debounce
  - Dropdown appears when query.length >= 2
  - Results grouped: Tracks, Artists, Albums (max 5 each)
  - Keyboard: Arrow keys navigate, Enter selects, Escape closes
  - Click outside closes dropdown

### 2. TagsCell Component
**File:** `frontend/src/components/library/TagsCell.tsx`
- **Purpose:** Compact tag display with inline editing capability
- **Interfaces:**
  ```typescript
  interface TagsCellProps {
    trackId: string;
    tags: string[];
    onTagsChange?: (tags: string[]) => void;
  }
  ```
- **Dependencies:** TagInput (existing), useMutation for tag ops
- **Behavior:**
  - Display mode: Show up to 3 badges + "+N more"
  - Edit mode: Expand to full TagInput on click
  - Click outside closes edit mode

### 3. AudioService
**File:** `frontend/src/lib/audio/audioService.ts`
- **Purpose:** Singleton service managing HLS.js playback
- **Interfaces:**
  ```typescript
  interface AudioService {
    load(url: string): Promise<void>;
    play(): void;
    pause(): void;
    seek(time: number): void;
    setVolume(level: number): void;
    destroy(): void;

    // Event callbacks (set by playerStore)
    onTimeUpdate?: (time: number) => void;
    onEnded?: () => void;
    onError?: (error: Error) => void;
    onDurationChange?: (duration: number) => void;
  }
  ```
- **Dependencies:** hls.js, HTML5 Audio element
- **Behavior:**
  - Detect HLS support (native or via hls.js)
  - Load and play HLS streams (.m3u8)
  - Fallback to direct file URL if HLS unavailable
  - Emit events for time updates, end, errors

### 4. Enhanced PlayerStore
**File:** `frontend/src/lib/store/playerStore.ts` (modify existing)
- **Purpose:** Connect Zustand state to AudioService
- **New Actions:**
  ```typescript
  interface PlayerActions {
    // Existing
    setQueue(tracks: Track[], startIndex: number): void;
    play(): void;
    pause(): void;
    next(): void;
    previous(): void;
    seek(time: number): void;
    setVolume(level: number): void;

    // New - internal
    _initAudio(): void;
    _loadTrack(track: Track): Promise<void>;
  }
  ```
- **Behavior:**
  - On setQueue: fetch stream URL, load into AudioService
  - On play/pause: delegate to AudioService
  - Subscribe to AudioService events for progress updates

### 5. Album Detail Page
**File:** `frontend/src/routes/albums/$albumId.tsx` (replace placeholder)
- **Purpose:** Display album with all tracks
- **Dependencies:** useAlbumQuery, TrackList, usePlayerStore
- **Layout:**
  ```
  ┌──────────────────────────────────────┐
  │ [Cover Art]  Album Name              │
  │              Artist • Year           │
  │              12 tracks • 45 min      │
  │              [▶ Play All]            │
  ├──────────────────────────────────────┤
  │ # │ Title        │ Duration          │
  │ 1 │ Track One    │ 3:45              │
  │ 2 │ Track Two    │ 4:12              │
  └──────────────────────────────────────┘
  ```

### 6. Artist Detail Page
**File:** `frontend/src/routes/artists/$artistName.tsx` (replace placeholder)
- **Purpose:** Display artist with albums and top tracks
- **Dependencies:** useArtistQuery, TrackList, usePlayerStore
- **Layout:**
  ```
  ┌──────────────────────────────────────┐
  │ Artist Name                          │
  │ 45 tracks • 5 albums                 │
  │ [▶ Play All]                         │
  ├──────────────────────────────────────┤
  │ Albums                               │
  │ [Album1] [Album2] [Album3]           │
  ├──────────────────────────────────────┤
  │ Top Tracks                           │
  │ 1. Track One           3:45         │
  │ 2. Track Two           4:12         │
  └──────────────────────────────────────┘
  ```

## Data Models

### SearchResult (existing, for reference)
```typescript
interface AutocompleteSuggestion {
  type: 'track' | 'artist' | 'album';
  value: string;        // Display text
  trackId?: string;     // For track results
  albumId?: string;     // For album results
}

interface AutocompleteResponse {
  suggestions: AutocompleteSuggestion[];
}
```

### StreamResponse
```typescript
interface StreamResponse {
  streamUrl: string;    // Signed HLS URL or direct file URL
  hlsAvailable: boolean;
  expiresAt: string;    // ISO timestamp
}
```

### Track (existing, for reference)
```typescript
interface Track {
  id: string;
  title: string;
  artist: string;
  album: string;
  albumId?: string;
  duration: number;
  tags: string[];
  s3Key: string;
  artworkS3Key?: string;
  // ... other fields
}
```

## Error Handling

### Error Scenarios

1. **Search API Failure**
   - **Handling:** Catch in useAutocompleteQuery, show "Search unavailable"
   - **User Impact:** Dropdown shows error state, user can retry

2. **HLS Stream Load Failure**
   - **Handling:** AudioService catches, attempts direct file fallback
   - **User Impact:** May have brief delay, playback continues with fallback

3. **Tag Add/Remove Failure**
   - **Handling:** Revert optimistic update, show toast error
   - **User Impact:** Tag change reverts, error message displayed

4. **Album/Artist Not Found**
   - **Handling:** API returns 404, show "Not found" UI
   - **User Impact:** Clear message with link back to library

5. **Audio Playback Error**
   - **Handling:** AudioService.onError fires, skip to next track
   - **User Impact:** Toast notification, playback continues with next

## Testing Strategy

### Unit Testing
- **SearchBar:** Mock useAutocompleteQuery, test debounce, keyboard nav
- **TagsCell:** Mock mutations, test display/edit mode toggle
- **AudioService:** Mock HLS.js, test load/play/pause/seek
- **PlayerStore:** Test state transitions, queue management

### Integration Testing
- **Search flow:** Type query → see results → click → navigate
- **Tag editing:** Click tags → add tag → verify API call → see update
- **Playback:** Click play → verify stream URL fetch → audio plays

### End-to-End Testing
- **Full search journey:** Search "beatles" → click track → plays
- **Album playback:** Navigate album → Play All → tracks play in order
- **Tag organization:** Add tag to track → filter by tag → track appears
