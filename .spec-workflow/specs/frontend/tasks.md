# Tasks - Frontend (Epic 5)

## Epic: Frontend
**Status**: Completed
**Wave**: 4

---

## Design Decisions

Based on implementation:
- **State Management**: Zustand for client state (player, theme), TanStack Query for server state
- **Routing**: TanStack Router with file-based routing
- **Styling**: DaisyUI 5 with Tailwind CSS 4, dark/light themes
- **Audio**: Howler.js for playback, Zustand store for queue management
- **Auth**: AWS Amplify v6 with Cognito JWT

---

## Wave 1: Authentication & Core Pages

### Task 1.1: Amplify Auth Configuration
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/amplify.ts`
- `frontend/src/lib/auth.ts`

**Description**: Configure AWS Amplify for Cognito authentication.

**Functions**:
| Function | Description |
|----------|-------------|
| `configureAuth()` | Initialize Amplify with Cognito config from env vars |
| `signIn(email, password)` | Authenticate user with Cognito |
| `signOut()` | Sign out and clear session |
| `getCurrentUser()` | Get current authenticated user |
| `getAccessToken()` | Get JWT token for API calls |

**Acceptance Criteria**:
- [x] Amplify configured with Cognito User Pool
- [x] Auth functions exported and tested (25 tests)

---

### Task 1.4: useAuth Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/hooks/useAuth.ts`
- `frontend/src/hooks/__tests__/useAuth.test.tsx`

**Description**: React hook for authentication state and actions.

**Hook API**:
```typescript
interface UseAuth {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  signIn: (email: string, password: string) => Promise<void>;
  signOut: () => Promise<void>;
  error: Error | null;
}
```

**Acceptance Criteria**:
- [x] Hook manages auth state correctly
- [x] Loading states handled
- [x] 18 tests passing

---

### Task 1.5: Login Page
**Status**: [x] Completed
**Files**:
- `frontend/src/routes/login.tsx`
- `frontend/src/routes/__tests__/login.test.tsx`

**Description**: Login page with form validation.

**Features**:
- Email/password form inputs
- Form validation (required fields, email format)
- Error display
- Redirect to home on success

**Acceptance Criteria**:
- [x] Form validation working
- [x] Error handling implemented
- [x] 19 tests passing

---

### Task 1.11: Home Page
**Status**: [x] Completed
**Files**:
- `frontend/src/routes/index.tsx`
- `frontend/src/routes/__tests__/index.test.tsx`

**Description**: Home page with library statistics.

**Features**:
- Total tracks, albums, artists counts
- Recent uploads section
- Quick access links

**Acceptance Criteria**:
- [x] Stats display correctly
- [x] Loading/error states handled
- [x] 22 tests passing

---

## Wave 2: Library Views

### Task 2.1: Tracks API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/tracks.ts`
- `frontend/src/hooks/useTracks.ts`
- `frontend/src/lib/api/__tests__/tracks.test.ts`
- `frontend/src/hooks/__tests__/useTracks.test.tsx`

**Description**: Track CRUD operations with query key factory.

**API Functions**:
| Function | Description |
|----------|-------------|
| `getTracks(params)` | Fetch paginated tracks |
| `getTrack(id)` | Fetch single track |
| `updateTrack(id, data)` | Update track metadata |
| `deleteTrack(id)` | Delete track |

**Query Keys**:
```typescript
export const trackKeys = {
  all: ['tracks'],
  lists: () => [...trackKeys.all, 'list'],
  list: (params) => [...trackKeys.lists(), params],
  details: () => [...trackKeys.all, 'detail'],
  detail: (id) => [...trackKeys.details(), id],
};
```

**Acceptance Criteria**:
- [x] API functions tested (10 tests)
- [x] Hook tested (13 tests)

---

### Task 2.2: Albums API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/albums.ts`
- `frontend/src/hooks/useAlbums.ts`
- `frontend/src/lib/api/__tests__/albums.test.ts`
- `frontend/src/hooks/__tests__/useAlbums.test.tsx`

**Acceptance Criteria**:
- [x] API functions tested (5 tests)
- [x] Hook tested (11 tests)

---

### Task 2.3: Artists API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/artists.ts`
- `frontend/src/hooks/useArtists.ts`
- `frontend/src/lib/api/__tests__/artists.test.ts`
- `frontend/src/hooks/__tests__/useArtists.test.tsx`

**Acceptance Criteria**:
- [x] API functions tested (5 tests)
- [x] Hook tested (11 tests)

---

### Task 2.4: Library Routes
**Status**: [x] Completed
**Files**:
- `frontend/src/routes/tracks/index.tsx`
- `frontend/src/routes/tracks/$trackId.tsx`
- `frontend/src/routes/albums/index.tsx`
- `frontend/src/routes/artists/index.tsx`
- `frontend/src/routes/__tests__/tracks.test.tsx`
- `frontend/src/routes/__tests__/trackDetail.test.tsx`
- `frontend/src/routes/__tests__/albums.test.tsx`
- `frontend/src/routes/__tests__/artists.test.tsx`

**Features**:
- Track list with sorting
- Track detail with edit/delete
- Album grid view
- Artist list view

**Acceptance Criteria**:
- [x] Tracks route tested (10 tests)
- [x] Track detail tested (17 tests)
- [x] Albums route tested (11 tests)
- [x] Artists route tested (11 tests)

---

## Wave 3: Audio Playback

### Task 3.1: Player Store
**Status**: [x] Completed (pre-existing)
**Files**:
- `frontend/src/lib/store/playerStore.ts`
- `frontend/src/lib/store/__tests__/playerStore.test.ts`

**Description**: Zustand store for audio playback with Howler.js.

**Store State**:
```typescript
interface PlayerState {
  currentTrack: Track | null;
  queue: Track[];
  queueIndex: number;
  isPlaying: boolean;
  volume: number;
  progress: number;
  duration: number;
  repeatMode: 'off' | 'all' | 'one';
  shuffle: boolean;
}
```

**Acceptance Criteria**:
- [x] 14 tests passing
- [x] Queue management working
- [x] Repeat/shuffle modes implemented

---

### Task 3.2: PlayerBar Component
**Status**: [x] Completed (pre-existing)
**Files**:
- `frontend/src/components/player/PlayerBar.tsx`
- `frontend/src/components/player/__tests__/PlayerBar.test.tsx`

**Description**: Fixed bottom player bar with controls.

**Acceptance Criteria**:
- [x] 5 tests passing
- [x] Play/pause, skip, seek, volume controls

---

## Wave 4: Upload & Search

### Task 4.1: Upload API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/upload.ts`
- `frontend/src/hooks/useUpload.ts`
- `frontend/src/lib/api/__tests__/upload.test.ts`
- `frontend/src/hooks/__tests__/useUpload.test.tsx`

**API Functions**:
| Function | Description |
|----------|-------------|
| `getPresignedUploadUrl(data)` | Get S3 presigned URL for upload |
| `confirmUpload(uploadId)` | Confirm upload completion |
| `getUploadStatus(uploadId)` | Poll upload processing status |

**Acceptance Criteria**:
- [x] API functions tested (4 tests)
- [x] Hook tested with progress tracking (6 tests)

---

### Task 4.2: Search API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/search.ts`
- `frontend/src/hooks/useSearch.ts`
- `frontend/src/lib/api/__tests__/search.test.ts`
- `frontend/src/hooks/__tests__/useSearch.test.tsx`

**API Functions**:
| Function | Description |
|----------|-------------|
| `searchTracks(params)` | Search tracks with filters |
| `searchAutocomplete(query)` | Get autocomplete suggestions |

**Acceptance Criteria**:
- [x] API functions tested (4 tests)
- [x] Hook tested with debounced queries (8 tests)

---

### Task 4.3: Upload & Search Routes
**Status**: [x] Completed
**Files**:
- `frontend/src/routes/upload.tsx`
- `frontend/src/routes/search.tsx`
- `frontend/src/routes/__tests__/upload.test.tsx`
- `frontend/src/routes/__tests__/search.test.tsx`

**Features**:
- Upload page with UploadDropzone
- Search results with artist/album filters
- Debounced filter updates (300ms)

**Acceptance Criteria**:
- [x] Upload route tested (9 tests)
- [x] Search route tested (9 tests)

---

## Wave 5: Playlists & Tags

### Task 5.1: Playlists API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/playlists.ts`
- `frontend/src/hooks/usePlaylists.ts`
- `frontend/src/lib/api/__tests__/playlists.test.ts`
- `frontend/src/hooks/__tests__/usePlaylists.test.tsx`

**API Functions**:
| Function | Description |
|----------|-------------|
| `getPlaylists(params)` | Fetch paginated playlists |
| `getPlaylist(id)` | Fetch single playlist with tracks |
| `createPlaylist(data)` | Create new playlist |
| `updatePlaylist(id, data)` | Update playlist |
| `deletePlaylist(id)` | Delete playlist |
| `addTrackToPlaylist(playlistId, trackId)` | Add track |
| `removeTrackFromPlaylist(playlistId, trackId)` | Remove track |

**Acceptance Criteria**:
- [x] API functions tested (7 tests)
- [x] Hook tested (6 tests)

---

### Task 5.2: Tags API & Hook
**Status**: [x] Completed
**Files**:
- `frontend/src/lib/api/tags.ts`
- `frontend/src/hooks/useTags.ts`
- `frontend/src/lib/api/__tests__/tags.test.ts`
- `frontend/src/hooks/__tests__/useTags.test.tsx`

**API Functions**:
| Function | Description |
|----------|-------------|
| `getTags(params)` | Fetch all tags |
| `getTracksByTag(tagName, params)` | Fetch tracks with tag |
| `addTagToTrack(trackId, tagName)` | Add tag to track |
| `removeTagFromTrack(trackId, tagName)` | Remove tag from track |

**Acceptance Criteria**:
- [x] API functions tested (4 tests)
- [x] Hook tested (6 tests)

---

### Task 5.3: Playlists Routes
**Status**: [x] Completed
**Files**:
- `frontend/src/routes/playlists/index.tsx`
- `frontend/src/routes/playlists/$playlistId.tsx`
- `frontend/src/routes/__tests__/playlists.test.tsx`
- `frontend/src/routes/__tests__/playlistDetail.test.tsx`

**Features**:
- Playlist grid with create modal
- Playlist detail with track list
- Play all functionality

**Acceptance Criteria**:
- [x] Playlists route tested (9 tests)
- [x] Playlist detail tested (11 tests)

---

### Task 5.4: Tags Routes
**Status**: [x] Completed
**Files**:
- `frontend/src/routes/tags/index.tsx`
- `frontend/src/routes/tags/$tagName.tsx`
- `frontend/src/routes/__tests__/tags.test.tsx`
- `frontend/src/routes/__tests__/tagDetail.test.tsx`

**Features**:
- Tag cloud with size based on track count
- Tracks filtered by tag
- Play all functionality

**Acceptance Criteria**:
- [x] Tags route tested (8 tests)
- [x] Tag detail tested (9 tests)

---

## Summary

| Wave | Tasks | Description |
|------|-------|-------------|
| Wave 1 | 4 | Auth config, useAuth, Login, Home |
| Wave 2 | 4 | Tracks/Albums/Artists API, hooks, routes |
| Wave 3 | 2 | PlayerStore, PlayerBar (pre-existing) |
| Wave 4 | 3 | Upload API/hook/route, Search API/hook/route |
| Wave 5 | 4 | Playlists API/hook/routes, Tags API/hook/routes |
| **Total** | **17** | |

---

## Test Plan Summary

### Unit Tests Created

| Category | Test Count | Status |
|----------|------------|--------|
| API tests | 50 | ✅ Completed |
| Hook tests | 92 | ✅ Completed |
| Store tests | 18 | ✅ Completed |
| Component tests | 34 | ✅ Completed |
| Route tests | 157 | ✅ Completed |
| **Total** | **351 tests** | ✅ All passing |

### Test Execution
```bash
# Run all tests
cd frontend && npm test -- --run

# Run with coverage
npm test -- --coverage
```

---

## Dependencies

### Between Tasks
- Wave 1 must complete before Wave 2 (auth required)
- Wave 2 provides API patterns for Waves 4-5
- Wave 3 (playerStore) used by routes in Waves 2, 4, 5

### External
- `@tanstack/react-query` - Server state management
- `@tanstack/react-router` - File-based routing
- `zustand` - Client state management
- `aws-amplify` - Cognito authentication
- `howler` - Audio playback
- `react-dropzone` - File upload
- `daisyui` - UI components

---

## PR Checklist

After completing all tasks:
- [x] All tests pass (`npm test -- --run`) - 351 tests passing
- [x] Code builds (`npm run build`)
- [x] CHANGELOG.md updated
- [x] CLAUDE.md files updated
- [x] epics-user-stories.md updated with completion status
- [x] Create PR to main with all changes - PR #8
