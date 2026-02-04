# Requirements Document: Access Control Bug Fixes

## Introduction

This spec addresses critical access control and visibility bugs discovered after the global-user-type feature deployment. The bugs affect how tracks are displayed to users based on their role and the track's visibility settings, guest user navigation, and admin user management functionality.

## Alignment with Product Vision

The Personal Music Search Engine requires proper multi-tenant isolation and role-based access control:
- Subscribers should only see their own content + public tracks from other users
- Admins should have global visibility across all user content
- Guests should have restricted access to public browsing only
- Proper HTTP status codes (403 vs 404) must indicate authorization vs existence

## Immediate Roadmap - Incomplete Items from Other Specs

The following items from other specs have NOT been completed and should be prioritized:

### From `search-streaming` (0/24 tasks)
- Full-text search with Nixiesearch integration
- Search autocomplete functionality
- Advanced search filters (BPM, key, genre)
- Search result pagination
- Stream URL generation with CloudFront
- HLS adaptive bitrate streaming

### From `semantic-search` (0/10 tasks)
- Bedrock Titan embedding generation
- Similar tracks by semantic content
- DJ-compatible track finder (BPM + key)
- Vector similarity search

### From `user-services` (0/13 tasks)
- User profile settings management
- Account preferences
- Storage usage tracking
- Activity history

### From `rights-management` (0/10 tasks)
- Content licensing metadata
- Rights holder attribution
- Usage restrictions enforcement
- License expiration handling

### From `frontend-enhancements` (partial - 31% remaining)
- Waveform visualization component
- Advanced player controls
- Keyboard shortcuts
- Mobile responsive fixes

### From `ui-ux-improvements` (partial - 11% remaining)
- Loading skeleton components
- Error boundary refinements
- Accessibility improvements

---

## Requirements

### REQ-1: Fix Admin User Search Modal Error

**User Story:** As an admin, I want to click on a user from search results to view their details without encountering errors, so that I can effectively manage user accounts.

#### Acceptance Criteria

1. WHEN admin searches for a user AND clicks on a user card THEN the system SHALL display the UserDetailModal with the user's full details
2. WHEN the modal opens AND the backend returns user details THEN the system SHALL populate all stats (tracks, playlists, albums, followers, following, storage)
3. IF the user details API returns an error THEN the system SHALL display an error alert in the modal with a descriptive message
4. WHEN the modal is open AND the user clicks Close or clicks outside the modal THEN the system SHALL close the modal and clear the selected user state

### REQ-2: Enforce Track Visibility for Non-Admin Users

**User Story:** As a subscriber, I want to only see my own tracks and tracks marked as public, so that other users' private content remains hidden from me.

#### Acceptance Criteria

1. WHEN a subscriber calls ListTracks endpoint THEN the system SHALL return only tracks where `userID == requester` OR `visibility == 'public'`
2. WHEN a subscriber has GlobalScope=false THEN the system SHALL NOT return private or unlisted tracks from other users
3. IF the filter includes visibility filters THEN the system SHALL apply them after the ownership/public filter
4. WHEN an admin user (with HasGlobal=true) calls ListTracks THEN the system SHALL return tracks from all users regardless of visibility

### REQ-3: Return 403 Forbidden for Unauthorized Track Access

**User Story:** As a subscriber, when I attempt to access a track I'm not authorized to view, I want to receive a clear "forbidden" error rather than a "not found" error, so I understand the content exists but I cannot access it.

#### Acceptance Criteria

1. WHEN a non-owner requests a track with `visibility == 'private'` THEN the system SHALL return HTTP 403 Forbidden
2. WHEN a non-owner requests a track with `visibility == 'unlisted'` AND they don't have the direct link context THEN the system SHALL return HTTP 404 Not Found
3. IF a track truly does not exist THEN the system SHALL return HTTP 404 Not Found
4. WHEN an admin (HasGlobal=true) requests any track THEN the system SHALL return the track regardless of visibility

### REQ-4: Global Admin Track Access

**User Story:** As an admin, I want to see all tracks in the system regardless of ownership or visibility, so that I can manage content and build playlists across the entire catalog.

#### Acceptance Criteria

1. WHEN an admin views the tracks page THEN the system SHALL display tracks from all users
2. WHEN an admin searches tracks THEN the system SHALL include results from all users' content
3. WHEN an admin adds a track to a playlist THEN the system SHALL allow any track to be added regardless of visibility
4. WHEN a subscriber views another user's playlist containing private tracks THEN the system SHALL hide those tracks from the playlist view (playlist tracks filtered by visibility)

### REQ-5: Guest User Route Restrictions

**User Story:** As a guest user (unauthenticated), I want to only access the public dashboard/discovery page, and see a clear permission error when trying to access other pages, so I understand I need to sign in for full access.

#### Acceptance Criteria

1. WHEN a guest user navigates to the dashboard/home route THEN the system SHALL display the public dashboard with browse-only content
2. WHEN a guest user clicks any navigation link other than dashboard THEN the system SHALL redirect to a "permission denied" page
3. WHEN the permission denied page is shown THEN the system SHALL display a clear message "You do not have permission to view this page"
4. WHEN the permission denied page is shown THEN the system SHALL display a button to return to the dashboard
5. IF a guest user directly navigates to a protected URL THEN the system SHALL redirect to the permission denied page

---

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Visibility checks in a dedicated visibility.go file in service layer
- **Modular Design**: Access control logic separate from business logic
- **Dependency Management**: Service layer enforces visibility, handlers pass auth context
- **Clear Interfaces**: `CanAccessTrack(ctx, userID, trackID, hasGlobal) (bool, error)`

### Performance
- Track visibility checks should not require additional database queries beyond the initial track fetch
- GSI3 (public track discovery) queries should be optimized for filtered results
- ListTracks with visibility filter should use single-query optimization where possible

### Security
- All visibility enforcement must occur at service layer, not just handler
- Frontend should not request data users cannot access (prevent information leakage)
- 403 responses should not leak existence of private content via timing attacks

### Reliability
- Visibility changes (public → private) should immediately reflect in query results
- No race conditions between visibility update and access checks
- Failed visibility checks should log for audit purposes

### Usability
- Clear error messages distinguishing "not found" from "not authorized"
- Guest permission denied page should have clear call-to-action to sign in
- Admin role switcher should continue to work for testing visibility scenarios
