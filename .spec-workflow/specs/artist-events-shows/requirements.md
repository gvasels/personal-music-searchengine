# Requirements Document: Artist Events & Shows (Phase 1 — Mock)

## Introduction

Users upload music to their personal library. They want to discover upcoming live shows and events for the artists in their collection. This spec covers Phase 1: building the full UI/UX with a pluggable events provider interface backed by mock data. Phase 2 (future spec) will integrate real APIs (Bandsintown, Ticketmaster, SeatGeek, Eventbrite, setlist.fm).

Additionally, the track ingest pipeline already creates Artist entities per attributed artist. This spec leverages those existing catalog records as the foundation for watching artists and querying events.

## Requirements

### Requirement 1: Watch Artists

**User Story:** As a user, I want to watch artists from my library, so that I can track their upcoming shows in one place.

#### Acceptance Criteria

1. WHEN a user views an artist detail page THEN the system SHALL display a "Watch" toggle button
2. WHEN a user clicks "Watch" THEN the system SHALL persist an ArtistWatch record linking the user to that artist name
3. WHEN a user is already watching an artist THEN the button SHALL display "Watching" with an option to unwatch
4. WHEN a user clicks "Unwatch" THEN the system SHALL remove the ArtistWatch record
5. WHEN a user views a track detail page THEN the system SHALL display a small "Watch" icon next to the artist name
6. IF the user is not authenticated THEN watch buttons SHALL NOT be displayed

### Requirement 2: Artist Events Section

**User Story:** As a user, I want to see upcoming events on an artist's page, so I can check if they're playing near me.

#### Acceptance Criteria

1. WHEN a user views an artist detail page THEN the system SHALL display an "Upcoming Events" section below the existing content
2. WHEN events exist THEN each event SHALL display: date, venue name, city/region, country, and a ticket link (if available)
3. WHEN no events exist THEN the system SHALL display "No upcoming events found"
4. WHEN events are loading THEN the system SHALL display a loading skeleton
5. IF the events provider returns an error THEN the system SHALL display a graceful error message without breaking the page

### Requirement 3: My Shows Page

**User Story:** As a user, I want to see upcoming shows for all my watched artists in one place, so I can plan which events to attend.

#### Acceptance Criteria

1. WHEN a user navigates to "/shows" THEN the system SHALL display upcoming events for all watched artists
2. WHEN events are displayed THEN they SHALL be sorted by date (soonest first)
3. WHEN an event is displayed THEN it SHALL show: artist name, date, venue, city, and ticket link
4. IF a watched artist has no upcoming events THEN they SHALL be omitted from the list
5. WHEN a user has no watched artists THEN the system SHALL display an empty state prompting them to browse their artists
6. WHEN a user has watched artists but none have events THEN the system SHALL display "No upcoming shows found for your watched artists"

### Requirement 4: Artist Search for Events

**User Story:** As a user, I want to search for any artist and see their upcoming events, so I can discover shows beyond my library.

#### Acceptance Criteria

1. WHEN a user is on the My Shows page THEN the system SHALL display a search input to look up any artist by name
2. WHEN a user searches for an artist THEN the system SHALL query the events provider and display matching results
3. WHEN a user selects a search result THEN the system SHALL display that artist's upcoming events inline
4. WHEN the events provider is unavailable THEN the system SHALL display a fallback message

### Requirement 5: Events Provider Interface (Backend)

**User Story:** As a developer, I want a pluggable events provider interface, so I can swap between mock data and real APIs without changing business logic.

#### Acceptance Criteria

1. WHEN the backend initializes THEN it SHALL load the configured events provider (mock by default)
2. WHEN mock provider is active THEN it SHALL return deterministic fake event data for any artist name
3. WHEN an API endpoint is called THEN it SHALL delegate to the configured provider through a clean interface
4. IF a new provider is added THEN it SHALL only need to implement the EventsProvider interface (no handler/service changes)

## Non-Functional Requirements

### Code Architecture and Modularity
- **Provider Pattern**: EventsProvider interface isolates external API logic; mock provider for Phase 1
- **Single Responsibility**: Events service, artist watch service, and providers are separate packages
- **Code Reuse**: Leverages existing Artist entity from ingest pipeline, existing DynamoDB patterns, existing frontend hook patterns

### Performance
- Mock provider returns data instantly (no external calls)
- My Shows page fetches events per watched artist in parallel (frontend)
- ArtistWatch queries use DynamoDB GSI for efficient listing

### Security
- Watch/unwatch endpoints require authentication (subscriber+ role)
- Events endpoints are read-only and available to authenticated users
- Mock provider has no external network calls

### Reliability
- Events provider failures are isolated — watch/unwatch always work (DynamoDB only)
- Frontend gracefully handles provider errors per artist (doesn't fail the whole page)

### Usability
- Event dates displayed in user's local timezone
- Ticket links open in new tab
- Watch state visible across artist detail and track detail pages
- Navigation link in sidebar under existing "Artists" section
