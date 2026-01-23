# Requirements Document: UI/UX Improvements

## Introduction

This spec covers frontend usability enhancements: keyboard shortcuts, drag-and-drop playlist reordering, batch operations on tracks, and mobile responsiveness improvements. These are quality-of-life features that improve daily usage without requiring backend changes.

## Alignment with Product Vision

This directly supports:
- **UI/UX Improvements** - Keyboard shortcuts, drag/drop, batch ops from roadmap
- **Mobile Responsive** - Improved mobile experience (marked "In Progress")
- **Creator Studio** - Efficient library management for all creator types

## Requirements

### Requirement 1: Keyboard Shortcuts

**User Story:** As a power user, I want keyboard shortcuts for common actions, so that I can control playback and navigate without using the mouse.

#### Acceptance Criteria

1. WHEN user presses `Space` THEN the system SHALL toggle play/pause
2. WHEN user presses `←` (Left Arrow) THEN the system SHALL skip to previous track
3. WHEN user presses `→` (Right Arrow) THEN the system SHALL skip to next track
4. WHEN user presses `↑` (Up Arrow) THEN the system SHALL increase volume by 10%
5. WHEN user presses `↓` (Down Arrow) THEN the system SHALL decrease volume by 10%
6. WHEN user presses `/` or `Cmd+K` THEN the system SHALL focus the search bar
7. WHEN user presses `M` THEN the system SHALL toggle mute
8. WHEN user presses `S` THEN the system SHALL toggle shuffle
9. WHEN user presses `R` THEN the system SHALL cycle repeat mode (off → all → one)
10. WHEN user presses `Escape` THEN the system SHALL close any open modal/dropdown
11. IF focus is in an input field THEN keyboard shortcuts SHALL be disabled (except Escape)
12. WHEN shortcuts are available THEN the system SHALL show `?` modal with all shortcuts

### Requirement 2: Drag and Drop Playlist Reordering

**User Story:** As a user, I want to reorder tracks in my playlist by dragging, so that I can arrange the playback order intuitively.

#### Acceptance Criteria

1. WHEN viewing a playlist detail page THEN each track row SHALL have a drag handle
2. WHEN user drags a track THEN the system SHALL show visual feedback (ghost element, drop zone)
3. WHEN track is dropped THEN the system SHALL update playlist track order via API
4. WHEN reordering THEN the system SHALL use optimistic update (immediate UI, background sync)
5. IF API call fails THEN the system SHALL revert order and show error toast
6. WHEN on mobile THEN the system SHALL support touch drag (long press to initiate)
7. WHEN multiple tracks are selected THEN drag SHALL move all selected tracks together

### Requirement 3: Batch Operations

**User Story:** As a user with a large library, I want to select multiple tracks and perform actions on them at once, so that I can efficiently organize my music.

#### Acceptance Criteria

1. WHEN viewing track list THEN each row SHALL have a checkbox for selection
2. WHEN user clicks checkbox THEN the system SHALL add track to selection
3. WHEN user shift-clicks THEN the system SHALL select range from last selected to clicked
4. WHEN user presses `Cmd/Ctrl+A` THEN the system SHALL select all visible tracks
5. WHEN tracks are selected THEN the system SHALL show batch action bar with options:
   - Add to playlist (dropdown to select playlist)
   - Add tags (input for tag names)
   - Delete selected
6. WHEN batch action is performed THEN the system SHALL show progress indicator
7. WHEN batch action completes THEN the system SHALL show success/error summary
8. IF some items fail THEN the system SHALL report partial success with details
9. WHEN selection exists THEN the system SHALL show "X tracks selected" count

### Requirement 4: Mobile Responsive Improvements

**User Story:** As a mobile user, I want the app to be fully functional on my phone, so that I can manage my music library on the go.

#### Acceptance Criteria

1. WHEN viewport is <768px THEN the sidebar SHALL collapse to hamburger menu
2. WHEN viewport is <768px THEN the track list SHALL switch to card view (not table)
3. WHEN viewport is <768px THEN the player bar SHALL stack controls vertically
4. WHEN on touch device THEN the system SHALL use touch-optimized controls (larger tap targets)
5. WHEN on mobile THEN search dropdown SHALL be full-screen overlay
6. WHEN on mobile THEN modals SHALL be full-screen instead of centered dialogs
7. WHEN on mobile THEN the system SHALL support swipe gestures:
   - Swipe right on track → quick play
   - Swipe left on track → delete/remove
8. WHEN player is minimized THEN the system SHALL show mini player with essential controls

### Requirement 5: Improved Search Experience

**User Story:** As a user, I want a faster and more visual search experience, so that I can find what I'm looking for quickly.

#### Acceptance Criteria

1. WHEN search results appear THEN the system SHALL group by type: Tracks, Albums, Artists, Playlists
2. WHEN displaying track results THEN the system SHALL show album art thumbnails
3. WHEN displaying artist results THEN the system SHALL show artist image (or generated avatar)
4. WHEN user types THEN the system SHALL debounce input (300ms) before querying
5. WHEN results are empty THEN the system SHALL suggest "Create playlist with this name?" option
6. WHEN search dropdown is open THEN arrow keys SHALL navigate results, Enter SHALL select
7. WHEN result is selected THEN the system SHALL navigate to appropriate page (not just play)

### Requirement 6: Accessibility Improvements

**User Story:** As a user with accessibility needs, I want the app to be fully navigable with keyboard and screen reader, so that I can use all features.

#### Acceptance Criteria

1. WHEN navigating THEN all interactive elements SHALL be focusable via Tab
2. WHEN focused THEN the system SHALL show visible focus indicators (outline)
3. WHEN using screen reader THEN all images SHALL have meaningful alt text
4. WHEN using screen reader THEN form inputs SHALL have associated labels
5. WHEN using screen reader THEN dynamic content changes SHALL be announced (aria-live)
6. WHEN color conveys meaning THEN the system SHALL also use icons/text (not color alone)
7. WHEN testing THEN the app SHALL pass WCAG 2.1 AA automated checks

### Requirement 7: Theme and Customization

**User Story:** As a user, I want additional themes and layout customization, so that I can personalize my experience.

#### Acceptance Criteria

1. WHEN viewing settings THEN the user SHALL be able to choose theme: Light, Dark, System
2. WHEN additional themes requested THEN the system SHALL support: Midnight, Sunset, Forest
3. WHEN theme changes THEN the system SHALL apply immediately without page reload
4. WHEN customizing layout THEN the user SHALL toggle sidebar visibility
5. WHEN customizing track list THEN the user SHALL show/hide columns (BPM, Key, Duration, etc.)
6. WHEN preferences change THEN the system SHALL persist to localStorage
7. IF user is authenticated THEN the system SHALL sync preferences to server

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: KeyboardShortcuts hook, DragDrop component, BatchOperations context
- **Modular Design**: All UI enhancements should be opt-in via feature flags
- **Dependency Management**: Use @dnd-kit for drag/drop (not react-dnd)
- **Clear Interfaces**: Define `ShortcutMap`, `BatchAction`, `LayoutConfig` types

### Performance
- Keyboard shortcuts must respond in <50ms
- Drag and drop must maintain 60fps during drag
- Batch operations must process 100 tracks in <5 seconds
- Mobile layout must not increase bundle size by >20KB

### Security
- Batch delete must require confirmation modal
- Drag reorder API calls must validate user ownership
- Theme/preference storage must sanitize values

### Reliability
- Keyboard shortcuts must not interfere with browser shortcuts
- Failed batch operations must not leave partial state
- Mobile gestures must have undo option for destructive actions

### Usability
- Shortcuts reference must be discoverable (? key, help menu)
- Drag handles must be visually distinct
- Batch selection must persist during pagination/scroll
- Mobile layout must be thumb-friendly (important controls at bottom)
