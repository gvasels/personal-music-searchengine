# Requirements Document: Frontend Enhancements

## Introduction

This spec covers six key frontend enhancements to improve the music library user experience: global search with live results, tags-based track organization with inline editing, audio playback with HLS streaming, and complete album/artist detail views. These features transform the app from a basic library into a fully functional music streaming application.

## Alignment with Product Vision

These enhancements align with the core product goal of providing a personal music library with powerful search, organization, and streaming capabilities. The features address gaps in the current MVP implementation where placeholder views and missing functionality limit the user experience.

## Requirements

### REQ-1: Global Search Bar with Live Results

**User Story:** As a user, I want to search my music library from anywhere in the app using a search bar in the header, so that I can quickly find tracks, albums, or artists without navigating to a dedicated search page.

#### Acceptance Criteria

1. WHEN the user types in the header search bar THEN the system SHALL display search results in a dropdown after 300ms debounce
2. WHEN search results are displayed THEN the system SHALL group results by type (Tracks, Artists, Albums)
3. WHEN the user clicks a search result THEN the system SHALL navigate to the appropriate detail page
4. WHEN the user presses Enter THEN the system SHALL navigate to `/search?q={query}` with full results
5. WHEN the search input is empty THEN the system SHALL hide the dropdown
6. IF no results match the query THEN the system SHALL display "No results found"

### REQ-2: Tags Column in Tracks List

**User Story:** As a user, I want to see tags for each track in the tracks list instead of the file format, so that I can understand how my tracks are organized at a glance.

#### Acceptance Criteria

1. WHEN viewing the tracks list THEN the system SHALL display a Tags column instead of Format
2. WHEN a track has tags THEN the system SHALL display up to 3 tag badges with "+N more" for additional tags
3. WHEN a track has no tags THEN the system SHALL display an "Add tags" placeholder
4. IF tags overflow the column width THEN the system SHALL truncate with ellipsis and show full list on hover

### REQ-3: Inline Tag Editing

**User Story:** As a user, I want to add and remove tags directly from the tracks list, so that I can organize my music without navigating to a separate page.

#### Acceptance Criteria

1. WHEN the user clicks on the tags cell THEN the system SHALL expand an inline tag editor
2. WHEN the user types a tag name and presses Enter THEN the system SHALL add the tag immediately
3. WHEN the user clicks the X on a tag badge THEN the system SHALL remove the tag immediately
4. WHEN tag operations complete THEN the system SHALL show success/error feedback via toast
5. IF a tag operation fails THEN the system SHALL revert the optimistic update and show error message

### REQ-4: Audio Playback with HLS Streaming

**User Story:** As a user, I want to play my uploaded music directly in the browser, so that I can listen to my library without downloading files.

#### Acceptance Criteria

1. WHEN the user clicks a play button THEN the system SHALL start audio playback using HLS.js
2. WHEN audio is playing THEN the system SHALL show real-time progress in the player bar
3. WHEN the user clicks pause THEN the system SHALL pause playback and retain position
4. WHEN the user clicks next/previous THEN the system SHALL play the corresponding track in the queue
5. WHEN the user adjusts volume THEN the system SHALL change audio output level
6. WHEN the user seeks via progress bar THEN the system SHALL jump to that position
7. IF HLS stream is not available THEN the system SHALL fall back to direct file playback
8. WHEN a track ends THEN the system SHALL play the next track in queue (if repeat is off)

### REQ-5: Play Button per Track

**User Story:** As a user, I want a play button on each track row, so that I can start playback of any track with a single click.

#### Acceptance Criteria

1. WHEN viewing the tracks list THEN the system SHALL display a play button on each row
2. WHEN hovering over a track row THEN the system SHALL highlight the play button
3. WHEN the user clicks the play button THEN the system SHALL queue all visible tracks and start from the clicked track
4. WHEN a track is currently playing THEN the system SHALL show a playing indicator instead of play button
5. IF the user clicks the playing indicator THEN the system SHALL pause playback

### REQ-6: Album Detail View

**User Story:** As a user, I want to view an album's full details including all tracks, so that I can browse and play albums as a cohesive unit.

#### Acceptance Criteria

1. WHEN navigating to `/albums/:albumId` THEN the system SHALL display full album information
2. WHEN album data loads THEN the system SHALL show: cover art, album name, artist, year, track count, total duration
3. WHEN album tracks load THEN the system SHALL display all tracks in a sortable list
4. WHEN the user clicks "Play All" THEN the system SHALL queue all album tracks and start playback
5. WHEN the user clicks a track THEN the system SHALL queue album tracks starting from that track
6. IF album has no cover art THEN the system SHALL display a placeholder image

### REQ-7: Artist Detail View

**User Story:** As a user, I want to view an artist's complete profile including albums and tracks, so that I can explore all music by a specific artist.

#### Acceptance Criteria

1. WHEN navigating to `/artists/:artistName` THEN the system SHALL display full artist information
2. WHEN artist data loads THEN the system SHALL show: artist name, track count, album count
3. WHEN artist albums load THEN the system SHALL display albums in a grid layout
4. WHEN artist tracks load THEN the system SHALL display top tracks in a list
5. WHEN the user clicks "Play All" THEN the system SHALL queue all artist tracks and start playback
6. WHEN the user clicks an album THEN the system SHALL navigate to that album's detail page

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Each component handles one feature (SearchBar, TagsCell, AudioService)
- **Modular Design**: Audio service is decoupled from UI, reusable across components
- **Dependency Management**: HLS.js loaded only when needed, tree-shakeable
- **Clear Interfaces**: Zustand store actions define clean contracts for audio control

### Performance
- Search debounce of 300ms to prevent excessive API calls
- Audio service preloads next track for gapless playback
- Tag updates use optimistic UI for instant feedback
- Lazy loading of album/artist data on route entry

### Security
- All streaming URLs are signed S3 URLs with expiration
- Audio service validates URL format before playback
- No user credentials stored in audio service

### Reliability
- Audio playback falls back to direct file if HLS unavailable
- Tag operations retry on transient failures
- Search gracefully handles API errors with user feedback

### Usability
- Keyboard navigation in search dropdown (arrow keys, Enter, Escape)
- Touch-friendly play buttons (44px tap targets)
- Accessible ARIA labels on all interactive elements
- Loading skeletons during data fetch
