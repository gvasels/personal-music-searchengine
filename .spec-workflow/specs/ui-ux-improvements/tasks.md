# Tasks Document: UI/UX Improvements

## Task Overview
| Task | Description | Estimated Files |
|------|-------------|-----------------|
| 1 | Keyboard Shortcuts Hook | 2 |
| 2 | Shortcuts Modal | 1 |
| 3 | Selection Context | 2 |
| 4 | Batch Action Bar | 1 |
| 5 | Drag and Drop | 3 |
| 6 | Mobile Navigation | 2 |
| 7 | Mobile Track Cards | 1 |
| 8 | User Preferences | 2 |
| 9 | Tests | 4 |

---

- [x] 1. Create Keyboard Shortcuts hook
  - Files: frontend/src/hooks/useKeyboardShortcuts.ts, frontend/src/lib/shortcuts/shortcuts.ts
  - Define default shortcut map (Space, arrows, /, M, S, R, Escape, ?)
  - Listen to keydown events globally
  - Disable when focus is in input/textarea
  - Integrate with playerStore actions
  - Purpose: Global keyboard control
  - _Leverage: frontend/src/lib/store/playerStore.ts_
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 1.10, 1.11_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer | Task: Create useKeyboardShortcuts hook with configurable shortcut map. Listen to document keydown. Check activeElement for input focus. Map shortcuts to playerStore actions (play/pause, volume, skip). Support modifier keys (Cmd+K) | Restrictions: Don't override browser shortcuts, handle Mac/Windows differences, debounce rapid presses | Success: All shortcuts work, disabled in inputs, no conflicts | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 2. Create Keyboard Shortcuts modal
  - File: frontend/src/components/ui/ShortcutsModal.tsx
  - Display all shortcuts grouped by category (Playback, Volume, Navigation, Modes)
  - Trigger with ? key
  - Show modifier key symbols (⌘, ⌃, ⇧)
  - Accessible focus trap
  - Purpose: Help users discover shortcuts
  - _Leverage: DaisyUI modal component_
  - _Requirements: 1.12_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create ShortcutsModal showing all keyboard shortcuts in grouped layout. Use DaisyUI modal. Show Mac symbols (⌘) vs Windows (Ctrl). Focus trap for accessibility. Close on Escape | Restrictions: Responsive layout, accessible, clear visual hierarchy | Success: Modal displays all shortcuts, grouped clearly, accessible | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 3. Create Selection Context for batch operations
  - Files: frontend/src/lib/context/SelectionContext.tsx, frontend/src/hooks/useSelection.ts
  - Implement SelectionContextValue interface
  - Support single select, shift-click range, Ctrl+A all
  - Persist selection during scroll/pagination
  - Expose selectionCount for UI
  - Purpose: Multi-select state management
  - _Leverage: React Context patterns_
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.9_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create SelectionContext with Set<string> for selectedIds. Implement select, deselect, toggle, selectRange (shift-click), selectAll. Expose via useSelection hook. Handle pagination (keep selection) | Restrictions: Efficient Set operations, memoize context value, clear API | Success: Selection works across interactions, persists during scroll, performant | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 4. Create Batch Action Bar
  - File: frontend/src/components/library/BatchActionBar.tsx
  - Show when tracks selected: "X tracks selected"
  - Actions: Add to playlist (dropdown), Add tags (input), Delete
  - Optimistic updates with rollback on failure
  - Show progress during operations
  - Purpose: Bulk operations UI
  - _Leverage: frontend/src/lib/context/SelectionContext.tsx, TanStack Query mutations_
  - _Requirements: 3.5, 3.6, 3.7, 3.8_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create BatchActionBar that appears when selection exists. Show count, action buttons. Add to playlist dropdown with user's playlists. Tag input for bulk tagging. Delete with confirmation. Show progress, handle partial failures | Restrictions: Optimistic updates, clear error reporting, confirm destructive actions | Success: All batch actions work, progress shown, errors handled gracefully | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 5. Implement Drag and Drop for playlist reorder
  - Files: frontend/src/components/playlist/DraggableTrackList.tsx, frontend/src/components/playlist/DraggableTrackItem.tsx, frontend/src/routes/playlists/$playlistId.tsx (modify)
  - Use @dnd-kit for drag and drop
  - Visual feedback during drag (ghost, drop indicator)
  - Optimistic reorder with API sync
  - Support multi-select drag
  - Touch support for mobile
  - Purpose: Intuitive playlist reordering
  - _Leverage: @dnd-kit/core, @dnd-kit/sortable_
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer with DnD expertise | Task: Create DraggableTrackList using @dnd-kit/sortable. Show drag handle, ghost element during drag, drop indicator. Reorder optimistically, sync to API. Support dragging multiple selected tracks. Touch-friendly | Restrictions: Maintain 60fps during drag, handle large lists (virtualize if needed), accessible | Success: Drag works smoothly, order persists, works on mobile | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 6. Improve Mobile Navigation
  - Files: frontend/src/components/layout/MobileNav.tsx, frontend/src/components/layout/Layout.tsx (modify)
  - Full-screen nav overlay on hamburger click
  - Swipe to close
  - Current route highlighting
  - Mini player visible at bottom
  - Purpose: Mobile-optimized navigation
  - _Leverage: frontend/src/components/layout/Sidebar.tsx_
  - _Requirements: 4.1, 4.4, 4.5, 4.6, 4.8_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create MobileNav as full-screen overlay. Show all nav items from Sidebar. Support swipe-to-close gesture. Highlight current route. Keep mini player visible. Animate open/close | Restrictions: Smooth animations (60fps), accessible, prevent body scroll when open | Success: Nav works on mobile, swipe close works, accessible | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 7. Create Mobile Track Card view
  - File: frontend/src/components/library/TrackCard.tsx
  - Card layout for tracks on mobile (< 768px)
  - Show album art, title, artist, duration
  - Swipe right to play, swipe left to delete
  - Tap to view details
  - Purpose: Mobile-friendly track display
  - _Leverage: frontend/src/components/library/TrackList.tsx_
  - _Requirements: 4.2, 4.7_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create TrackCard component for mobile view. Show cover art, title, artist in card format. Implement swipe gestures (right=play, left=delete). Tap for details. Use in TracksPage when viewport < 768px | Restrictions: Smooth swipe gestures, clear visual feedback, undo for swipe-delete | Success: Cards display on mobile, swipe actions work, responsive | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 8. Create User Preferences system
  - Files: frontend/src/lib/store/preferencesStore.ts, frontend/src/routes/settings.tsx
  - Store: theme, sidebarVisible, trackListColumns, shortcutsEnabled
  - Persist to localStorage
  - Sync to server when authenticated (optional)
  - Settings page for configuration
  - Purpose: User customization
  - _Leverage: frontend/src/lib/store/themeStore.ts_
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/Zustand Developer | Task: Create preferencesStore with Zustand. Store theme, sidebar visibility, visible columns, shortcuts enabled. Persist to localStorage. Create Settings page with toggles/dropdowns. Optionally sync to server | Restrictions: Handle missing localStorage gracefully, type-safe, migration for schema changes | Success: Preferences persist, settings UI works, syncs when authenticated | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 9. Write tests for UI/UX features
  - Files: frontend/src/hooks/__tests__/useKeyboardShortcuts.test.ts, frontend/src/components/library/__tests__/BatchActionBar.test.tsx, frontend/src/components/playlist/__tests__/DraggableTrackList.test.tsx, frontend/src/components/layout/__tests__/MobileNav.test.tsx
  - Test keyboard shortcuts (fire events, verify actions)
  - Test selection context (range select, all)
  - Test drag and drop (reorder verification)
  - Test mobile interactions (swipe, tap)
  - Purpose: Ensure UI features are reliable
  - _Leverage: React Testing Library, user-event_
  - _Requirements: All_
  - _Prompt: Implement the task for spec ui-ux-improvements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Test Engineer | Task: Write tests for keyboard shortcuts (fireEvent.keyDown), selection context (state changes), drag-drop (@dnd-kit testing utils), mobile interactions (touch events). Test accessibility (keyboard navigation, focus management) | Restrictions: Use user-event for realistic interactions, test edge cases, accessibility assertions | Success: 80%+ coverage, all tests pass, accessibility verified | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_
