# Tasks Document: Artist Events & Shows (Phase 1 — Mock)

## Group 1: Backend — Models, Provider Interface, and Mock Provider

- [-] 1.1 Unit tests for ArtistWatch model and Event types
  - Files: `backend/internal/models/artist_watch_test.go`
  - Test ArtistWatch DynamoDB key generation (PK, SK, GSI1), name normalization, Event struct serialization
  - _Leverage: `backend/internal/models/follow_test.go` for pattern, `backend/internal/models/common.go` for DynamoDBItem_
  - _Requirements: 5.1, 1.1_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer specializing in domain models and DynamoDB single-table design | Task: Write failing unit tests for ArtistWatch model (DynamoDB key generation with PK=USER#{userId} SK=ARTIST_WATCH#{normalizedName}, GSI1PK=ARTIST_WATCH#{normalizedName} GSI1SK=USER#{userId}, name normalization to lowercase/trimmed) and Event types (JSON serialization). Follow existing patterns in `backend/internal/models/follow_test.go`. Tests MUST fail initially. | Restrictions: Do NOT write implementation code. Only test files. | Success: Tests exist, compile, and FAIL because model types don't exist yet. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [-] 1.2 Unit tests for EventsProvider interface and MockProvider
  - Files: `backend/internal/events/provider_test.go`
  - Test MockProvider returns deterministic events for a given artist name, SearchArtists returns results
  - _Leverage: None (new package)_
  - _Requirements: 5.2, 5.3_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Write failing unit tests for MockProvider implementing EventsProvider interface. Test: GetArtistEvents returns 2-4 events with valid fields (date in future, venue/city/country populated, source="mock"), SearchArtists returns results matching query, deterministic output (same artist name → same events). Create new package `backend/internal/events/`. Tests MUST fail initially. | Restrictions: Do NOT write implementation code. Only test files. | Success: Tests exist, compile (with stub interface file if needed), and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 1.3 Implement ArtistWatch model and Event types
  - Files: `backend/internal/models/artist_watch.go`, `backend/internal/models/event.go`
  - Create ArtistWatch struct with DynamoDB keys, Event and ArtistSearchResult structs, EventResponse types
  - _Leverage: `backend/internal/models/follow.go` for DynamoDB entity pattern, `backend/internal/models/common.go` for DynamoDBItem_
  - _Requirements: 1.1, 5.1_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Implement ArtistWatch model (PK=USER#{userId}, SK=ARTIST_WATCH#{normalizedName}, GSI1PK=ARTIST_WATCH#{normalizedName}, GSI1SK=USER#{userId}, fields: userId, artistName, watchedAt) and Event/ArtistSearchResult types per design.md. Read test files FIRST. Follow existing patterns in follow.go. | Restrictions: ONLY write code needed to pass task 1.1 tests. | Success: All task 1.1 tests pass. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 1.4 Implement EventsProvider interface and MockProvider
  - Files: `backend/internal/events/provider.go`, `backend/internal/events/mock_provider.go`
  - Create EventsProvider interface, implement MockProvider with deterministic fake data
  - _Leverage: None (new package, clean interface)_
  - _Requirements: 5.2, 5.3_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Implement EventsProvider interface (GetArtistEvents, SearchArtists) and MockProvider. Mock generates 2-4 events per artist using hash of name for determinism, dates 1-6 months in future, pool of real venue/city/country combos, source="mock". Read test files FIRST. | Restrictions: ONLY write code needed to pass task 1.2 tests. | Success: All task 1.2 tests pass. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

## Group 2: Backend — Repository, Service, and Handlers

- [ ] 2.1 Unit tests for ArtistWatch repository methods
  - Files: `backend/internal/repository/artist_watch_test.go`
  - Test CreateArtistWatch, DeleteArtistWatch, GetArtistWatch, ListWatchedArtists with mocked DynamoDB
  - _Leverage: `backend/internal/repository/follow_test.go` for pattern_
  - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Write failing unit tests for ArtistWatch repository methods: CreateArtistWatch (PutItem with condition), DeleteArtistWatch, GetArtistWatch, ListWatchedArtists (Query on PK with SK prefix ARTIST_WATCH#, pagination). Follow patterns in follow_test.go or existing repository tests. Tests MUST fail initially. | Restrictions: Only test files. | Success: Tests compile and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 2.2 Unit tests for ArtistWatchService and EventsService
  - Files: `backend/internal/service/artist_watch_test.go`, `backend/internal/service/events_test.go`
  - Test watch/unwatch/isWatching/listWatched, test GetArtistEvents/SearchArtists/GetWatchedArtistEvents
  - _Leverage: `backend/internal/service/follow_test.go` for pattern_
  - _Requirements: 1.1-1.4, 2.1-2.4, 3.1-3.6_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Write failing unit tests for ArtistWatchService (WatchArtist, UnwatchArtist, IsWatching, ListWatchedArtists with mocked repo) and EventsService (GetArtistEvents delegates to provider, SearchArtists delegates to provider, GetWatchedArtistEvents calls ListWatchedArtists then GetArtistEvents per artist). Tests MUST fail. | Restrictions: Only test files. | Success: Tests compile and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 2.3 Unit tests for events and watch HTTP handlers
  - Files: `backend/internal/handlers/events_test.go`, `backend/internal/handlers/artist_watch_test.go`
  - Test all 6 endpoints: watch, unwatch, watch status, list watched, get events, search events
  - _Leverage: `backend/internal/handlers/follow_test.go` for pattern_
  - _Requirements: All_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Write failing unit tests for HTTP handlers: POST/DELETE/GET /artists/:name/watch, GET /users/me/watched-artists, GET /artists/:name/events, GET /events/search?q=. Test auth requirements (subscriber+), request parsing, response format. Follow patterns in follow_test.go. Tests MUST fail. | Restrictions: Only test files. | Success: Tests compile and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 2.4 Implement ArtistWatch repository methods
  - Files: `backend/internal/repository/dynamodb.go` (or new `artist_watch.go`), `backend/internal/repository/repository.go`
  - Add CreateArtistWatch, DeleteArtistWatch, GetArtistWatch, ListWatchedArtists to repository interface and DynamoDB implementation
  - _Leverage: `backend/internal/repository/follow.go` for DynamoDB patterns_
  - _Requirements: 1.1-1.4_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Implement ArtistWatch repository methods. Add to Repository interface: CreateArtistWatch, DeleteArtistWatch, GetArtistWatch, ListWatchedArtists. DynamoDB implementation: PutItem with condition (no overwrite), DeleteItem, GetItem, Query PK=USER#{userId} with SK begins_with ARTIST_WATCH#. Read tests FIRST. Follow follow.go patterns. | Restrictions: ONLY code to pass 2.1 tests. | Success: Task 2.1 tests pass. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 2.5 Implement ArtistWatchService and EventsService
  - Files: `backend/internal/service/artist_watch.go`, `backend/internal/service/events.go`
  - Business logic for watching artists and querying events
  - _Leverage: `backend/internal/service/follow.go` for pattern_
  - _Requirements: 1.1-1.4, 2.1-2.4, 3.1-3.6_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Implement ArtistWatchService (WatchArtist, UnwatchArtist, IsWatching, ListWatchedArtists) and EventsService (GetArtistEvents via provider, SearchArtists via provider, GetWatchedArtistEvents lists watched then fetches events per artist). Read tests FIRST. | Restrictions: ONLY code to pass 2.2 tests. | Success: Task 2.2 tests pass. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 2.6 Implement events and watch HTTP handlers + wire routes
  - Files: `backend/internal/handlers/events.go`, `backend/internal/handlers/artist_watch.go`, `backend/cmd/api/main.go`
  - HTTP handlers + route registration in main.go
  - _Leverage: `backend/internal/handlers/follow.go` for pattern, `backend/cmd/api/main.go` for wiring_
  - _Requirements: All backend_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Implement HTTP handlers for all 6 endpoints per design.md. Wire into main.go: create EventsService (with MockProvider) and ArtistWatchService in NewServices, create handler, register routes. CRITICAL: Follow wiring checklist — services initialized, handlers created, routes registered. Read tests FIRST. | Restrictions: ONLY code to pass 2.3 tests. Verify with `go build ./cmd/api`. | Success: Task 2.3 tests pass AND `go build ./cmd/api` succeeds. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

## Group 3: Frontend — Types, API Client, Hooks

- [ ] 3.1 Unit tests for artist watch and events API clients and hooks
  - Files: `frontend/src/lib/api/__tests__/artistWatch.test.ts`, `frontend/src/lib/api/__tests__/events.test.ts`, `frontend/src/hooks/__tests__/useArtistWatch.test.ts`, `frontend/src/hooks/__tests__/useArtistEvents.test.ts`
  - Test API functions and hook behavior with mocked axios
  - _Leverage: `frontend/src/lib/api/__tests__/follows.test.ts`, `frontend/src/hooks/__tests__/useFollow.test.ts` for patterns_
  - _Requirements: 1.1-1.6, 2.1-2.5, 3.1-3.6_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: TypeScript/React Developer | Task: Write failing tests for: (1) artistWatch.ts API client (watchArtist, unwatchArtist, isWatching, getWatchedArtists), (2) events.ts API client (getArtistEvents, searchArtistsForEvents), (3) useArtistWatch hook (watch toggle, loading states), (4) useArtistEvents hook (data fetching, error handling). Follow patterns in follows.test.ts and useFollow.test.ts. Tests MUST fail. | Restrictions: Only test files. | Success: Tests compile and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 3.2 Implement frontend types, API clients, and hooks
  - Files: `frontend/src/types/index.ts` (add types), `frontend/src/lib/api/artistWatch.ts`, `frontend/src/lib/api/events.ts`, `frontend/src/hooks/useArtistWatch.ts`, `frontend/src/hooks/useArtistEvents.ts`
  - ArtistEvent, ArtistSearchResult, WatchedArtist types; API functions; TanStack Query hooks
  - _Leverage: `frontend/src/lib/api/follows.ts`, `frontend/src/hooks/useFollow.ts` for patterns_
  - _Requirements: 1.1-1.6, 2.1-2.5, 3.1-3.6_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: TypeScript/React Developer | Task: Add ArtistEvent, ArtistSearchResult, WatchedArtist types to types/index.ts. Implement artistWatch.ts (watchArtist, unwatchArtist, isWatching, getWatchedArtists) and events.ts (getArtistEvents, searchArtistsForEvents) API clients. Implement useArtistWatch hook (useWatchToggle with optimistic updates) and useArtistEvents hook (query with enabled flag). Read tests FIRST. Follow follows.ts/useFollow.ts patterns. | Restrictions: ONLY code to pass 3.1 tests. | Success: Task 3.1 tests pass. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

## Group 4: Frontend — Components and Pages

- [ ] 4.1 Unit tests for WatchButton, EventCard, and ArtistEventsSection components
  - Files: `frontend/src/components/events/__tests__/WatchButton.test.tsx`, `frontend/src/components/events/__tests__/EventCard.test.tsx`, `frontend/src/components/events/__tests__/ArtistEventsSection.test.tsx`
  - Test rendering, interactions, loading/error states
  - _Leverage: `frontend/src/components/follow/__tests__/FollowButton.test.tsx` for pattern_
  - _Requirements: 1.1-1.6, 2.1-2.5_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer | Task: Write failing tests for: (1) WatchButton — renders watch/unwatch toggle, requires auth, handles click, (2) EventCard — renders date/venue/city/country/ticket link, handles missing ticketUrl, (3) ArtistEventsSection — loading skeleton, renders events list, empty state, error state. Follow FollowButton.test.tsx patterns. Tests MUST fail. | Restrictions: Only test files. | Success: Tests compile and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 4.2 Unit tests for My Shows page
  - Files: `frontend/src/routes/__tests__/shows.test.tsx`
  - Test: renders watched artist events, empty state (no watched artists), search functionality, sorted by date
  - _Leverage: `frontend/src/routes/__tests__/artists.test.tsx` for pattern_
  - _Requirements: 3.1-3.6, 4.1-4.3_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer | Task: Write failing tests for My Shows page (/shows route): (1) renders events for watched artists sorted by date, (2) empty state when no watched artists, (3) "no upcoming shows" state, (4) artist search input and results, (5) loading states. Follow artists.test.tsx patterns. Tests MUST fail. | Restrictions: Only test files. | Success: Tests compile and FAIL. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 4.3 Implement WatchButton, EventCard, and ArtistEventsSection components
  - Files: `frontend/src/components/events/WatchButton.tsx`, `frontend/src/components/events/EventCard.tsx`, `frontend/src/components/events/ArtistEventsSection.tsx`, `frontend/src/components/events/index.ts`
  - DaisyUI-styled components for event display and artist watching
  - _Leverage: `frontend/src/components/follow/FollowButton.tsx`, `frontend/src/components/track/SimilarTracks.tsx` for patterns_
  - _Requirements: 1.1-1.6, 2.1-2.5_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer | Task: Implement WatchButton (toggle with auth check, DaisyUI btn), EventCard (date formatted, venue/city/country, ticket link opens new tab, status badge for cancelled/postponed), ArtistEventsSection (fetches events, loading skeleton, error state, empty state). Read tests FIRST. Follow FollowButton and SimilarTracks patterns. | Restrictions: ONLY code to pass 4.1 tests. | Success: Task 4.1 tests pass. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

- [ ] 4.4 Implement My Shows page and integrate into existing pages
  - Files: `frontend/src/routes/shows.tsx`, `frontend/src/routes/artists/$artistName.tsx` (modify), `frontend/src/routes/tracks/$trackId.tsx` (modify), `frontend/src/components/layout/Sidebar.tsx` (modify)
  - New /shows route, add events section to artist detail, add watch icon to track detail, add nav link
  - _Leverage: Existing artist detail page, track detail page, Sidebar component_
  - _Requirements: All_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer | Task: (1) Create /shows route with watched artist events + search, (2) Add ArtistEventsSection to artist detail page below tracks, (3) Add WatchButton to artist detail header, (4) Add small watch icon next to artist link on track detail page, (5) Add "Shows" link in Sidebar nav. Read tests FIRST. CRITICAL: regenerate route tree after adding shows.tsx. | Restrictions: ONLY code to pass 4.2 tests + visual integration. | Success: Task 4.2 tests pass, `npm run build` succeeds. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._

## Group 5: Verification

- [ ] 5.1 Build validation and contract alignment
  - Run: `go vet ./...`, `go build ./cmd/api`, `go test ./...`, `tsc --noEmit`, `eslint .`, `npm run build`, `npm test -- --run`
  - Verify param names match between FE and BE (artistName param, query params)
  - _Requirements: All_
  - _Prompt: Implement the task for spec artist-events-shows, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer | Task: Run full build validation for both backend and frontend. Verify contract alignment: FE API client param names match BE handler query param names. Fix any build/lint/type errors. | Restrictions: Only fix issues found during validation. | Success: All builds pass with zero errors, all tests pass, contract alignment verified. Set task to in-progress in tasks.md before starting, log implementation after, mark complete._
