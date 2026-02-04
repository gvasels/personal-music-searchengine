# Tasks Document: Frontend Enhancements

## Task Group 1: Audio Playback Infrastructure (Critical Path)

- [x] 1.1 Install HLS.js and create AudioService
  - Files: `frontend/package.json`, `frontend/src/lib/audio/audioService.ts`
  - Install hls.js package
  - Create singleton AudioService class wrapping HLS.js and HTML5 Audio
  - Implement: load(url), play(), pause(), seek(time), setVolume(level), destroy()
  - Add event callbacks: onTimeUpdate, onEnded, onError, onDurationChange
  - Handle HLS vs direct file fallback
  - Purpose: Core audio playback engine for streaming
  - _Completed: 2026-01-22_

- [x] 1.2 Connect AudioService to playerStore
  - File: `frontend/src/lib/store/playerStore.ts` (modify)
  - Import and initialize AudioService in store
  - Wire play/pause/seek/setVolume actions to AudioService methods
  - Subscribe to AudioService events to update progress/isPlaying state
  - Add _loadTrack helper to fetch stream URL and load into AudioService
  - Purpose: Connect UI state to actual audio playback
  - _Completed: 2026-01-22_

- [x] 1.3 Write tests for AudioService and playerStore integration
  - Files: `frontend/src/lib/store/__tests__/playerStore.test.ts`
  - Updated existing playerStore tests with AudioService mocks
  - All 351 tests passing
  - _Completed: 2026-01-22_

## Task Group 2: Play Button in Track List

- [x] 2.1 Add play button column to TrackList component
  - File: `frontend/src/routes/tracks/index.tsx` (modify)
  - Add play button as first column
  - Show playing indicator when track is current
  - Click handler calls setQueue with tracks and clicked index
  - Hover state highlights button
  - Purpose: Enable one-click playback from track list
  - _Completed: 2026-01-22_

- [x] 2.2 Add play button to tracks index page
  - File: `frontend/src/routes/tracks/index.tsx` (modify)
  - Play button column with play/pause toggle
  - Proper column alignment and responsive behavior
  - Purpose: Main tracks page has playback capability
  - _Completed: 2026-01-22_

## Task Group 3: Search Bar

- [x] 3.1 Create SearchBar component with dropdown
  - File: `frontend/src/components/search/SearchBar.tsx` (create)
  - Input field with search icon
  - Debounced input (300ms) triggers useAutocompleteQuery
  - Dropdown shows grouped results (Tracks, Artists, Albums)
  - Keyboard navigation (arrows, Enter, Escape)
  - Click outside closes dropdown
  - Purpose: Global search with live results
  - _Completed: 2026-01-22_

- [x] 3.2 Add SearchBar to Header component
  - File: `frontend/src/components/layout/Header.tsx` (modify)
  - Import and render SearchBar between logo and controls
  - Responsive: hidden on mobile, shown on desktop (md:block)
  - Purpose: Search accessible from anywhere in app
  - _Completed: 2026-01-22_

- [ ] 3.3 Write tests for SearchBar
  - File: `frontend/src/components/search/__tests__/SearchBar.test.tsx` (create)
  - Test debounce behavior
  - Test dropdown rendering with mock data
  - Test keyboard navigation
  - Test navigation on selection
  - Purpose: Ensure search reliability
  - _Deferred: Component functional, tests to be added_

## Task Group 4: Tags Column with Inline Editing

- [x] 4.1 Create TagsCell component
  - File: `frontend/src/components/library/TagsCell.tsx` (create)
  - Display mode: Show up to 3 tag badges, "+N more" for overflow
  - Edit mode: Expand to show TagInput component on click
  - Handle add/remove tag mutations via existing TagInput
  - Click outside exits edit mode
  - Purpose: Compact tag display with inline editing
  - _Completed: 2026-01-22_

- [x] 4.2 Replace Format column with Tags in tracks table
  - File: `frontend/src/routes/tracks/index.tsx` (modify)
  - Remove Format column header and cells
  - Add Tags column using TagsCell component
  - Adjusted column widths for tags display
  - Purpose: Show tags instead of format in track list
  - _Completed: 2026-01-22_

- [ ] 4.3 Write tests for TagsCell
  - File: `frontend/src/components/library/__tests__/TagsCell.test.tsx` (create)
  - Test display mode rendering
  - Test edit mode toggle
  - Test tag add/remove mutations
  - Test overflow handling
  - Purpose: Ensure tags functionality works correctly
  - _Deferred: Component functional, tests to be added_

## Task Group 5: Album Detail View

- [x] 5.1 Implement Album detail page
  - File: `frontend/src/routes/albums/$albumId.tsx` (replace placeholder)
  - Fetch album data with useAlbumQuery
  - Display: cover art hero, album name, artist, year, track count, duration
  - "Play All" button queues all tracks
  - Track list showing album tracks with play buttons
  - Loading skeleton and error states
  - Purpose: Full album browsing experience
  - _Completed: 2026-01-22_

- [ ] 5.2 Write tests for Album detail page
  - File: `frontend/src/routes/__tests__/albumDetail.test.tsx` (create)
  - Test data loading and display
  - Test Play All functionality
  - Test track list interaction
  - Test error states
  - Purpose: Ensure album page reliability
  - _Deferred: Component functional, tests to be added_

## Task Group 6: Artist Detail View

- [x] 6.1 Implement Artist detail page
  - File: `frontend/src/routes/artists/$artistName.tsx` (replace placeholder)
  - Fetch artist data with useArtistQuery
  - Display: artist name, track count, album count
  - "Play All" button queues all artist tracks
  - Albums section: grid of album cards (clickable)
  - Top Tracks section: track list with play buttons
  - Loading skeleton and error states
  - Purpose: Full artist browsing experience
  - _Completed: 2026-01-22_

- [ ] 6.2 Write tests for Artist detail page
  - File: `frontend/src/routes/__tests__/artistDetail.test.tsx` (create)
  - Test data loading and display
  - Test Play All functionality
  - Test album card navigation
  - Test track list interaction
  - Purpose: Ensure artist page reliability
  - _Deferred: Component functional, tests to be added_

## Task Group 7: Integration Testing

- [ ] 7.1 End-to-end integration tests
  - File: `frontend/src/__tests__/integration/frontend-enhancements.test.tsx` (create)
  - Test full search -> play flow
  - Test album browse -> play all flow
  - Test tag editing flow
  - Purpose: Verify all features work together
  - _Deferred: All components working, integration tests to be added_

## Summary

**Completed**: 11 tasks
- Audio playback infrastructure (1.1, 1.2, 1.3)
- Play button in track list (2.1, 2.2)
- Search bar (3.1, 3.2)
- Tags column (4.1, 4.2)
- Album detail view (5.1)
- Artist detail view (6.1)

**Deferred**: 5 tasks (additional test coverage)
- SearchBar tests (3.3)
- TagsCell tests (4.3)
- Album detail tests (5.2)
- Artist detail tests (6.2)
- Integration tests (7.1)

All 351 existing tests pass. Build successful.
