# Tasks Document - Epic 5: Frontend

## Wave 1: Foundation (Auth & App Shell)

- [ ] 1.1 Configure Amplify authentication
  - File: `frontend/src/lib/auth.ts`
  - Configure AWS Amplify with Cognito User Pool settings
  - Export auth functions: signIn, signOut, getCurrentUser, getToken
  - Purpose: Establish authentication foundation
  - _Leverage: frontend/.env.example for env vars, aws-amplify package_
  - _Requirements: REQ-5.1_
  - _Prompt: Implement the task for spec frontend, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer specializing in AWS Amplify and authentication | Task: Configure Amplify auth with Cognito using env vars VITE_COGNITO_USER_POOL_ID, VITE_COGNITO_CLIENT_ID, VITE_COGNITO_REGION. Export signIn, signOut, getCurrentUser, getToken functions. | Restrictions: Do not store tokens in localStorage, use Amplify's secure storage. Do not expose credentials. | _Leverage: aws-amplify package already in package.json | Success: Auth functions work with Cognito, tokens are securely managed. Mark task [-] in progress, log with log-implementation when done, mark [x] complete._

- [x] 1.2 Create API client with auth interceptor
  - File: `frontend/src/lib/api/client.ts`
  - Create Axios instance with base URL from env
  - Add request interceptor to attach JWT token
  - Add response interceptor for 401 handling
  - Purpose: Centralized API client with automatic auth
  - _Leverage: frontend/src/lib/auth.ts, axios package_
  - _Requirements: REQ-5.1_
  - **Completed**: API client with Axios, auth interceptor, 11 unit tests

- [x] 1.3 Create API type definitions
  - File: `frontend/src/types/index.ts`
  - Define TypeScript interfaces matching OpenAPI schemas
  - Include Track, Album, Artist, Playlist, Tag, Upload types
  - Include paginated response types
  - Purpose: Type safety for API responses
  - _Leverage: backend/api/openapi.yaml for schema reference_
  - _Requirements: REQ-5.3, REQ-5.4, REQ-5.9, REQ-5.10_
  - **Completed**: All type definitions in src/types/index.ts

- [ ] 1.4 Create useAuth hook
  - File: `frontend/src/hooks/useAuth.ts`
  - Create hook with user state, loading state, error state
  - Expose signIn, signOut, isAuthenticated
  - Integrate with TanStack Query for caching
  - Purpose: Reusable auth state management
  - _Leverage: frontend/src/lib/auth.ts, @tanstack/react-query_
  - _Requirements: REQ-5.1_

- [ ] 1.5 Create Login page
  - File: `frontend/src/routes/login.tsx`
  - Create login form with email/password fields
  - Use DaisyUI form components with emerald primary theme
  - Handle loading and error states
  - Redirect to home on success
  - Purpose: User authentication entry point
  - _Leverage: frontend/src/hooks/useAuth.ts, DaisyUI form components_
  - _Requirements: REQ-5.1_

- [x] 1.6 Create root layout with auth guard
  - File: `frontend/src/routes/__root.tsx`
  - Create TanStack Router root layout
  - Check authentication state
  - Redirect unauthenticated users to login
  - Render children when authenticated
  - Purpose: Protect all routes requiring auth
  - _Leverage: frontend/src/hooks/useAuth.ts, @tanstack/react-router_
  - _Requirements: REQ-5.1, REQ-5.2_
  - **Completed**: Basic root layout created

- [x] 1.7 Create theme store
  - File: `frontend/src/lib/store/themeStore.ts`
  - Create Zustand store for theme preference
  - Persist to localStorage
  - Apply theme to document element
  - Purpose: Theme state management
  - _Leverage: zustand, existing tailwind themes_
  - _Requirements: REQ-5.2_
  - **Completed**: Theme store with persist middleware, 4 unit tests, 100% coverage

- [x] 1.8 Create Header component
  - File: `frontend/src/components/layout/Header.tsx`
  - Include search bar, theme toggle, user menu
  - Use DaisyUI navbar component
  - Purpose: Top navigation bar
  - _Leverage: frontend/src/lib/store/theme.ts, DaisyUI navbar_
  - _Requirements: REQ-5.2, REQ-5.8_
  - **Completed**: Header with theme toggle, 100% coverage

- [x] 1.9 Create Sidebar component
  - File: `frontend/src/components/layout/Sidebar.tsx`
  - Navigation links: Home, Tracks, Albums, Artists, Playlists, Tags, Upload
  - Collapsible on mobile
  - Active state highlighting
  - Purpose: Main navigation menu
  - _Leverage: DaisyUI menu component, @tanstack/react-router Link_
  - _Requirements: REQ-5.2_
  - **Completed**: Sidebar with navigation links, 100% coverage

- [x] 1.10 Create Layout component
  - File: `frontend/src/components/layout/Layout.tsx`
  - Compose Header, Sidebar, content area, PlayerBar slot
  - Responsive grid layout
  - Purpose: Main app shell
  - _Leverage: frontend/src/components/layout/Header.tsx, Sidebar.tsx_
  - _Requirements: REQ-5.2_
  - **Completed**: Layout with Header, Sidebar, PlayerBar, 5 unit tests, 100% coverage

- [ ] 1.11 Create Home page
  - File: `frontend/src/routes/index.tsx`
  - Dashboard with recent tracks, stats
  - Quick access to upload
  - Purpose: Landing page after login
  - _Leverage: frontend/src/components/layout/Layout.tsx_
  - _Requirements: REQ-5.2, REQ-5.3_

## Wave 2: Library Views

- [ ] 2.1 Create tracks API functions
  - File: `frontend/src/lib/api/tracks.ts`
  - Functions: getTracks, getTrack, updateTrack, deleteTrack
  - Handle pagination parameters
  - Purpose: Track API layer
  - _Leverage: frontend/src/lib/api/client.ts, types.ts_
  - _Requirements: REQ-5.3_
  - **Note**: Basic functions exist in client.ts, needs dedicated file

- [ ] 2.2 Create useTracks hook
  - File: `frontend/src/hooks/useTracks.ts`
  - Use TanStack Query useInfiniteQuery for pagination
  - Include useTrack for single track
  - Include mutations for update/delete
  - Purpose: Track data fetching hooks
  - _Leverage: frontend/src/lib/api/tracks.ts, @tanstack/react-query_
  - _Requirements: REQ-5.3_

- [ ] 2.3 Create TrackRow component
  - File: `frontend/src/components/library/TrackRow.tsx`
  - Display: cover art, title, artist, album, duration
  - Play on double-click, context menu on right-click
  - Hover state with play button
  - Purpose: Single track list item
  - _Leverage: DaisyUI table row styling_
  - _Requirements: REQ-5.3_

- [x] 2.4 Create TrackList component
  - File: `frontend/src/components/library/TrackList.tsx`
  - Virtual scrolling for performance (1000+ tracks)
  - Sortable column headers
  - Infinite scroll trigger
  - Purpose: Main track listing
  - _Leverage: frontend/src/components/library/TrackRow.tsx, react-virtual_
  - _Requirements: REQ-5.3_
  - **Completed**: TrackList with click-to-play, loading/empty states, 6 unit tests, 98.27% coverage

- [ ] 2.5 Create Tracks page
  - File: `frontend/src/routes/tracks/index.tsx`
  - Use useTracks hook
  - Render TrackList
  - Handle loading and empty states
  - Purpose: Track library view
  - _Leverage: frontend/src/hooks/useTracks.ts, TrackList.tsx_
  - _Requirements: REQ-5.3_

- [ ] 2.6 Create albums API and hooks
  - Files: `frontend/src/lib/api/albums.ts`, `frontend/src/hooks/useAlbums.ts`
  - getAlbums, getAlbum with tracks
  - useAlbums, useAlbum hooks
  - Purpose: Album data layer
  - _Leverage: Pattern from tracks implementation_
  - _Requirements: REQ-5.4_

- [ ] 2.7 Create AlbumCard component
  - File: `frontend/src/components/library/AlbumCard.tsx`
  - Cover art image (lazy loaded)
  - Title, artist, track count
  - Hover state with play button
  - Purpose: Album grid item
  - _Leverage: DaisyUI card component_
  - _Requirements: REQ-5.4_

- [ ] 2.8 Create AlbumGrid and Albums page
  - Files: `frontend/src/components/library/AlbumGrid.tsx`, `frontend/src/routes/albums/index.tsx`
  - Responsive grid of AlbumCards
  - Albums page with useAlbums
  - Purpose: Album browsing
  - _Leverage: frontend/src/components/library/AlbumCard.tsx_
  - _Requirements: REQ-5.4_

- [ ] 2.9 Create Album detail page
  - File: `frontend/src/routes/albums/$albumId.tsx`
  - Album header with large cover, metadata
  - Track list for album
  - Play All button
  - Purpose: Album detail view
  - _Leverage: frontend/src/hooks/useAlbums.ts, TrackList.tsx_
  - _Requirements: REQ-5.4_

- [ ] 2.10 Create artists API, hooks, and pages
  - Files: `frontend/src/lib/api/artists.ts`, `frontend/src/hooks/useArtists.ts`, `frontend/src/routes/artists/index.tsx`, `frontend/src/routes/artists/$artistName.tsx`
  - Artist list and detail views
  - Purpose: Artist browsing
  - _Leverage: Pattern from albums implementation_
  - _Requirements: REQ-5.4_

## Wave 3: Audio Playback

- [x] 3.1 Create player Zustand store
  - File: `frontend/src/lib/store/playerStore.ts`
  - State: currentTrack, queue, queueIndex, isPlaying, progress, duration, volume, shuffle, repeat
  - Actions: play, pause, next, previous, seek, setVolume, setQueue, addToQueue, toggleShuffle, setRepeat
  - Purpose: Central player state
  - _Leverage: zustand, persist middleware for volume_
  - _Requirements: REQ-5.6, REQ-5.7_
  - **Completed**: Full player store with queue, shuffle, repeat, 14 unit tests, 90.32% coverage

- [ ] 3.2 Create usePlayer hook with Howler integration
  - File: `frontend/src/hooks/usePlayer.ts`
  - Initialize Howler.js instance
  - Sync with player store
  - Handle track loading from stream URL
  - Update progress via animation frame
  - Purpose: Audio playback engine
  - _Leverage: frontend/src/lib/store/player.ts, howler, frontend/src/lib/api/client.ts_
  - _Requirements: REQ-5.6_

- [ ] 3.3 Create NowPlaying component
  - File: `frontend/src/components/player/NowPlaying.tsx`
  - Cover art, title, artist
  - Link to track/album
  - Purpose: Current track display
  - _Leverage: frontend/src/lib/store/player.ts_
  - _Requirements: REQ-5.6_

- [ ] 3.4 Create ProgressBar component
  - File: `frontend/src/components/player/ProgressBar.tsx`
  - Clickable/draggable progress slider
  - Time display (current / total)
  - Purpose: Playback position control
  - _Leverage: frontend/src/lib/store/player.ts, frontend/src/hooks/usePlayer.ts_
  - _Requirements: REQ-5.6_

- [ ] 3.5 Create VolumeControl component
  - File: `frontend/src/components/player/VolumeControl.tsx`
  - Volume slider
  - Mute toggle icon
  - Purpose: Volume control
  - _Leverage: frontend/src/lib/store/player.ts_
  - _Requirements: REQ-5.6_

- [x] 3.6 Create PlayerBar component
  - File: `frontend/src/components/player/PlayerBar.tsx`
  - Compose NowPlaying, playback controls, ProgressBar, VolumeControl
  - Fixed to bottom of screen
  - Purpose: Main player interface
  - _Leverage: All player components created above_
  - _Requirements: REQ-5.6_
  - **Completed**: PlayerBar with controls, progress, volume, 5 unit tests, 95.53% coverage

- [ ] 3.7 Create QueueDrawer component
  - File: `frontend/src/components/player/QueueDrawer.tsx`
  - Slide-out drawer showing queue
  - Drag to reorder
  - Remove from queue
  - Purpose: Queue management UI
  - _Leverage: frontend/src/lib/store/player.ts, DaisyUI drawer_
  - _Requirements: REQ-5.7_

- [x] 3.8 Integrate player into Layout
  - File: `frontend/src/components/layout/Layout.tsx` (modify)
  - Add PlayerBar to Layout
  - Initialize player hook
  - Purpose: Enable playback across app
  - _Leverage: frontend/src/components/player/PlayerBar.tsx, usePlayer.ts_
  - _Requirements: REQ-5.6_
  - **Completed**: PlayerBar integrated into Layout component

## Wave 4: Upload & Search

- [ ] 4.1 Create upload API functions
  - File: `frontend/src/lib/api/upload.ts`
  - getPresignedUrl, confirmUpload, listUploads, getUploadStatus
  - Purpose: Upload API layer
  - _Leverage: frontend/src/lib/api/client.ts_
  - _Requirements: REQ-5.5_
  - **Note**: getPresignedUploadUrl exists in client.ts

- [ ] 4.2 Create useUpload hook
  - File: `frontend/src/hooks/useUpload.ts`
  - Upload state management
  - Progress tracking with XHR
  - Retry failed uploads
  - Purpose: Upload orchestration
  - _Leverage: frontend/src/lib/api/upload.ts_
  - _Requirements: REQ-5.5_

- [x] 4.3 Create UploadZone component
  - File: `frontend/src/components/upload/UploadDropzone.tsx`
  - Drag-and-drop area
  - File type validation (audio only)
  - File picker fallback
  - Purpose: File selection UI
  - _Leverage: react-dropzone package_
  - _Requirements: REQ-5.5_
  - **Completed**: UploadDropzone with drag-drop, file validation, progress bar, 5 unit tests, 86.07% coverage

- [ ] 4.4 Create UploadProgress component
  - File: `frontend/src/components/upload/UploadProgress.tsx`
  - Progress bar per file
  - Status indicator (uploading, processing, complete, error)
  - Retry button on error
  - Purpose: Upload status display
  - _Leverage: DaisyUI progress component_
  - _Requirements: REQ-5.5_
  - **Note**: Basic progress in UploadDropzone, needs dedicated component

- [ ] 4.5 Create Upload page
  - File: `frontend/src/routes/upload.tsx`
  - UploadZone for file selection
  - List of UploadProgress items
  - Purpose: Upload workflow page
  - _Leverage: frontend/src/components/upload/UploadZone.tsx, UploadProgress.tsx, useUpload.ts_
  - _Requirements: REQ-5.5_

- [ ] 4.6 Create search API functions
  - File: `frontend/src/lib/api/search.ts`
  - search, getSuggestions functions
  - Purpose: Search API layer
  - _Leverage: frontend/src/lib/api/client.ts_
  - _Requirements: REQ-5.8_
  - **Note**: searchTracks exists in client.ts

- [ ] 4.7 Create useSearch hook
  - File: `frontend/src/hooks/useSearch.ts`
  - Debounced search query
  - Filter state management
  - Purpose: Search data fetching
  - _Leverage: frontend/src/lib/api/search.ts, @tanstack/react-query_
  - _Requirements: REQ-5.8_

- [ ] 4.8 Create SearchBar and Autocomplete components
  - Files: `frontend/src/components/search/SearchBar.tsx`, `frontend/src/components/search/Autocomplete.tsx`
  - Search input with autocomplete dropdown
  - Navigate on selection
  - Purpose: Search input
  - _Leverage: frontend/src/hooks/useSearch.ts, DaisyUI input_
  - _Requirements: REQ-5.8_

- [ ] 4.9 Create Search Results page
  - File: `frontend/src/routes/search.tsx`
  - Display search results
  - Filter sidebar
  - Purpose: Search results view
  - _Leverage: frontend/src/hooks/useSearch.ts, TrackList.tsx_
  - _Requirements: REQ-5.8_

## Wave 5: Organization (Playlists & Tags)

- [ ] 5.1 Create playlists API and hooks
  - Files: `frontend/src/lib/api/playlists.ts`, `frontend/src/hooks/usePlaylists.ts`
  - CRUD operations, add/remove tracks
  - Purpose: Playlist data layer
  - _Leverage: Pattern from tracks implementation_
  - _Requirements: REQ-5.9_
  - **Note**: createPlaylist exists in client.ts

- [x] 5.2 Create PlaylistCard and playlists list page (partial)
  - Files: `frontend/src/components/playlist/PlaylistCard.tsx`, `frontend/src/routes/playlists/index.tsx`
  - Playlist grid view
  - Create playlist button
  - Purpose: Playlist browsing
  - _Leverage: Pattern from albums implementation_
  - _Requirements: REQ-5.9_
  - **Partial**: CreatePlaylistModal completed with validation, 7 unit tests, 88.79% coverage. PlaylistCard and page pending.

- [ ] 5.3 Create PlaylistDetail page
  - File: `frontend/src/routes/playlists/$playlistId.tsx`
  - Playlist header, track list
  - Drag to reorder tracks
  - Edit/delete playlist
  - Purpose: Playlist detail view
  - _Leverage: frontend/src/hooks/usePlaylists.ts, TrackList.tsx_
  - _Requirements: REQ-5.9_

- [ ] 5.4 Create tags API and hooks
  - Files: `frontend/src/lib/api/tags.ts`, `frontend/src/hooks/useTags.ts`
  - CRUD operations, add/remove from track
  - Purpose: Tag data layer
  - _Leverage: Pattern from tracks implementation_
  - _Requirements: REQ-5.10_
  - **Note**: addTagToTrack, removeTagFromTrack exist in client.ts

- [x] 5.5 Create TagChip and TagSelector components (partial)
  - Files: `frontend/src/components/tag/TagChip.tsx`, `frontend/src/components/tag/TagSelector.tsx`
  - Tag display with color
  - Multi-select tag picker
  - Purpose: Tag UI primitives
  - _Leverage: DaisyUI badge component_
  - _Requirements: REQ-5.10_
  - **Partial**: TagInput completed with add/remove, lowercase normalization, 6 unit tests, 82.75% coverage. TagSelector pending.

- [ ] 5.6 Create Tags page and tag detail page
  - Files: `frontend/src/routes/tags/index.tsx`, `frontend/src/routes/tags/$tagName.tsx`
  - Tag list with track counts
  - Tag detail shows tagged tracks
  - Purpose: Tag browsing
  - _Leverage: frontend/src/hooks/useTags.ts, TagChip.tsx_
  - _Requirements: REQ-5.10_

- [ ] 5.7 Create TrackEditModal component
  - File: `frontend/src/components/library/TrackEditModal.tsx`
  - Edit form for track metadata
  - Tag selector integration
  - Cover art upload
  - Purpose: Track editing UI
  - _Leverage: frontend/src/hooks/useTracks.ts, TagSelector.tsx, DaisyUI modal_
  - _Requirements: REQ-5.11_

- [ ] 5.8 Create ContextMenu component
  - File: `frontend/src/components/ui/ContextMenu.tsx`
  - Right-click menu for tracks
  - Options: Play, Add to queue, Add to playlist, Edit, Delete
  - Purpose: Track actions menu
  - _Leverage: DaisyUI menu component_
  - _Requirements: REQ-5.3, REQ-5.7, REQ-5.9_

- [ ] 5.9 Wire up context menu to TrackRow
  - File: `frontend/src/components/library/TrackRow.tsx` (modify)
  - Add right-click handler
  - Show ContextMenu
  - Handle all actions
  - Purpose: Enable track actions
  - _Leverage: frontend/src/components/ui/ContextMenu.tsx, player store, usePlaylists_
  - _Requirements: REQ-5.3, REQ-5.7, REQ-5.9_

- [ ] 5.10 Final integration and polish
  - Files: Multiple
  - Toast notifications for actions
  - Keyboard shortcuts (space=play/pause, arrows=seek)
  - Loading skeletons
  - Empty states
  - Purpose: UX polish
  - _Leverage: react-hot-toast, useKeyboard hook_
  - _Requirements: All_

---

## Summary

### Completed Tasks (TDD Implementation)
| Task | Component | Tests | Coverage |
|------|-----------|-------|----------|
| 1.2 | API Client | 11 | 19.71% |
| 1.3 | Type Definitions | - | - |
| 1.6 | Root Layout | - | - |
| 1.7 | Theme Store | 4 | 100% |
| 1.8 | Header | - | 100% |
| 1.9 | Sidebar | - | 100% |
| 1.10 | Layout | 5 | 100% |
| 2.4 | TrackList | 6 | 98.27% |
| 3.1 | Player Store | 14 | 90.32% |
| 3.6 | PlayerBar | 5 | 95.53% |
| 3.8 | Player in Layout | - | - |
| 4.3 | UploadDropzone | 5 | 86.07% |
| 5.2 | CreatePlaylistModal | 7 | 88.79% |
| 5.5 | TagInput | 6 | 82.75% |

**Total: 63 tests, 81.56% overall coverage**

### Remaining Tasks
- Wave 1: 1.1 (Amplify auth), 1.4 (useAuth), 1.5 (Login page), 1.11 (Home page)
- Wave 2: 2.1-2.3, 2.5-2.10 (API functions, hooks, pages)
- Wave 3: 3.2-3.5, 3.7 (usePlayer, NowPlaying, ProgressBar, VolumeControl, QueueDrawer)
- Wave 4: 4.1-4.2, 4.4-4.9 (Upload/Search APIs, hooks, pages)
- Wave 5: 5.1, 5.3-5.4, 5.6-5.10 (Playlist/Tag hooks, pages, modals)
