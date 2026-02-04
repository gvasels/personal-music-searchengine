# Tasks Document: Access Control Bug Fixes

## Task Group 1: Backend Track Visibility Enforcement

- [x] 1.1 Add visibility check to GetTrack service method
  - File: `backend/internal/service/track.go`
  - Modify `GetTrack` to accept `hasGlobal bool` parameter
  - Add visibility check: if not owner and not admin, check track.Visibility
  - Return `models.NewForbiddenError` for unauthorized private track access
  - Return 404 for unlisted tracks when accessed by non-owner without direct context
  - Purpose: Enforce visibility at service layer for single track retrieval
  - _Leverage: `models.NewForbiddenError`, existing `GetTrack` method_
  - _Requirements: REQ-3, REQ-4_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: Go Backend Developer specializing in service layer and access control

    Task: Modify GetTrack in backend/internal/service/track.go to enforce visibility rules:
    1. Add `hasGlobal bool` parameter to GetTrack signature
    2. After fetching track, check if requester can access it:
       - Owner (userID == track.UserID): Allow
       - Admin (hasGlobal == true): Allow
       - Public track (visibility == "public"): Allow
       - Private track and not owner: Return models.NewForbiddenError("you do not have permission to access this track")
       - Track not found: Return models.NewNotFoundError (existing behavior)
    3. Update the service interface in service.go to match new signature

    Restrictions:
    - Do not modify repository layer
    - Keep existing cover art URL generation logic
    - Maintain backward compatibility with tests (update test calls)

    _Leverage: models.NewForbiddenError (already exists), existing GetTrack implementation

    Success:
    - GetTrack returns 403 for non-owner accessing private tracks
    - GetTrack returns track for owners, admins, and public tracks
    - All existing tests updated and passing

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 1.2 Update GetTrack handler to pass HasGlobal
  - File: `backend/internal/handlers/track.go`
  - Change `GetTrack` handler to use `getAuthContext(c)` instead of `getUserIDFromContext(c)`
  - Pass `auth.HasGlobal` to service layer
  - Purpose: Provide full auth context to service layer for visibility checks
  - _Leverage: Existing `getAuthContext` helper, `ListTracks` handler pattern_
  - _Requirements: REQ-3, REQ-4_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: Go Backend Developer specializing in HTTP handlers

    Task: Update GetTrack handler in backend/internal/handlers/track.go:
    1. Replace `getUserIDFromContext(c)` with `getAuthContext(c)`
    2. Pass `auth.HasGlobal` as third parameter to `h.services.Track.GetTrack`
    3. Follow the pattern used in ListTracks handler

    Restrictions:
    - Do not modify other handlers in this task
    - Keep existing error handling pattern

    _Leverage: getAuthContext helper, ListTracks handler as reference

    Success:
    - GetTrack handler passes HasGlobal to service
    - Handler compiles and existing tests pass

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 1.3 Add visibility filtering to ListTracks service
  - File: `backend/internal/service/track.go`
  - When `filter.GlobalScope == false`, filter results to only include:
    - Tracks where `userID == requesterID`, OR
    - Tracks where `visibility == "public"`
  - Purpose: Prevent subscribers from seeing other users' private tracks in lists
  - _Leverage: Existing ListTracks method, TrackFilter.GlobalScope field_
  - _Requirements: REQ-2, REQ-4_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: Go Backend Developer specializing in data filtering and access control

    Task: Modify ListTracks in backend/internal/service/track.go to filter by visibility:
    1. After getting results from repository, if filter.GlobalScope is false:
       - Filter results to include only tracks where:
         a) track.UserID == requesterID (user's own tracks), OR
         b) track.Visibility == models.VisibilityPublic
    2. This filtering happens in memory after repository query (simpler than modifying repository)
    3. Adjust pagination counts if filtering removes items

    Alternative approach (if repository supports it):
    - Pass visibility filter to repository and let DynamoDB handle it
    - Check if repository.ListTracks can accept visibility criteria

    Restrictions:
    - Do not break existing GlobalScope=true behavior for admins
    - Maintain pagination cursor integrity

    _Leverage: Existing ListTracks, filter.GlobalScope field

    Success:
    - Subscribers only see their own tracks + public tracks
    - Admins continue to see all tracks
    - Pagination works correctly with filtered results

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 1.4 Write unit tests for visibility enforcement
  - File: `backend/internal/service/track_visibility_test.go` (new or extend existing)
  - Test cases:
    - Owner can access private track
    - Admin can access any track
    - Non-owner gets 403 for private track
    - Non-owner can access public track
    - Subscriber list only shows own + public tracks
  - Purpose: Ensure visibility rules are correctly enforced with test coverage
  - _Leverage: Existing test patterns in service tests, MockRepository_
  - _Requirements: REQ-2, REQ-3, REQ-4_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: Go Test Engineer specializing in unit testing

    Task: Create/extend tests in backend/internal/service/track_visibility_test.go:

    Test cases for GetTrack:
    1. TestGetTrack_OwnerCanAccessPrivateTrack
    2. TestGetTrack_AdminCanAccessAnyTrack
    3. TestGetTrack_NonOwnerGetsForbiddenForPrivate
    4. TestGetTrack_NonOwnerCanAccessPublic
    5. TestGetTrack_NotFoundReturns404

    Test cases for ListTracks:
    6. TestListTracks_SubscriberSeesOnlyOwnAndPublic
    7. TestListTracks_AdminSeesAllTracks

    Use table-driven tests where appropriate.

    Restrictions:
    - Use existing mock patterns from the codebase
    - Test service layer only, not handlers

    _Leverage: Existing track_visibility_test.go, MockRepository patterns

    Success:
    - All test cases pass
    - Tests cover happy paths and error cases
    - 80%+ coverage on visibility logic

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

---

## Task Group 2: Frontend Guest Route Protection

- [x] 2.1 Create Permission Denied page
  - File: `frontend/src/routes/permission-denied.tsx`
  - Create new route component with:
    - Clear "You do not have permission" message
    - Explanation that sign-in is required
    - Button/link to return to dashboard (/)
    - DaisyUI styling consistent with rest of app
  - Purpose: Provide clear feedback when guests try to access protected routes
  - _Leverage: Existing route patterns, DaisyUI card/alert components_
  - _Requirements: REQ-5_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: React Frontend Developer specializing in TanStack Router

    Task: Create frontend/src/routes/permission-denied.tsx:
    1. Create route component for /permission-denied path
    2. Display centered card with:
       - Icon (lock or shield)
       - Heading: "Access Denied"
       - Message: "You do not have permission to view this page. Please sign in to continue."
       - Primary button: "Go to Dashboard" linking to /
       - Secondary button: "Sign In" linking to /login
    3. Use DaisyUI components (card, btn, alert)
    4. Follow existing route file patterns

    Restrictions:
    - Keep styling consistent with existing pages
    - No additional dependencies

    _Leverage: Existing route patterns (e.g., login.tsx), DaisyUI components

    Success:
    - Route accessible at /permission-denied
    - Page renders correctly with message and buttons
    - Navigation works to / and /login

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 2.2 Add route guard for guest users in root layout
  - File: `frontend/src/routes/__root.tsx`
  - Add logic to check if user is authenticated
  - Define list of public routes: `/`, `/login`, `/permission-denied`
  - Redirect unauthenticated users to `/permission-denied` for non-public routes
  - Purpose: Enforce guest access restrictions at routing level
  - _Leverage: useAuth hook, TanStack Router navigation_
  - _Requirements: REQ-5_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: React Frontend Developer specializing in route protection

    Task: Modify frontend/src/routes/__root.tsx to add guest route guard:
    1. Import useAuth hook and useLocation/Navigate from router
    2. Define PUBLIC_ROUTES constant: ['/', '/login', '/permission-denied']
    3. In the root component, check:
       - If not authenticated AND current path not in PUBLIC_ROUTES
       - Redirect to /permission-denied
    4. Handle loading state while auth is being determined
    5. Ensure the check runs on every navigation

    Restrictions:
    - Do not block authenticated users from any route
    - Keep existing layout structure intact
    - Handle auth loading state to avoid flicker

    _Leverage: useAuth hook (isAuthenticated, isLoading), router navigation

    Success:
    - Guest users can access /, /login, /permission-denied
    - Guest users redirected to /permission-denied for all other routes
    - Authenticated users can access all routes
    - No flash of protected content during auth loading

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 2.3 Write tests for guest route protection
  - File: `frontend/src/routes/__tests__/permission-denied.test.tsx` (new)
  - Test cases:
    - Permission denied page renders correctly
    - Guest accessing /tracks is redirected
    - Guest can access dashboard
    - Authenticated user can access all routes
  - Purpose: Ensure route protection works correctly
  - _Leverage: Existing route test patterns, test-utils_
  - _Requirements: REQ-5_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: React Test Engineer specializing in route testing

    Task: Create frontend/src/routes/__tests__/permission-denied.test.tsx:

    Test cases:
    1. TestPermissionDeniedPage_RendersMessage - page shows access denied message
    2. TestPermissionDeniedPage_HasDashboardButton - button links to /
    3. TestPermissionDeniedPage_HasSignInButton - button links to /login
    4. TestRouteGuard_GuestRedirectedFromTracks - mock unauthenticated, navigate to /tracks
    5. TestRouteGuard_GuestCanAccessDashboard - mock unauthenticated, navigate to /
    6. TestRouteGuard_AuthenticatedCanAccessTracks - mock authenticated, navigate to /tracks

    Restrictions:
    - Mock useAuth appropriately for each test case
    - Use existing test utilities

    _Leverage: Existing route test patterns, @testing-library/react

    Success:
    - All test cases pass
    - Tests properly mock auth state
    - Route guard behavior verified

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

---

## Task Group 3: Admin User Modal Fix

- [x] 3.1 Debug and fix admin user detail modal error
  - Files: `frontend/src/components/admin/UserDetailModal.tsx`, `frontend/src/hooks/useAdmin.ts`
  - Investigate the error when clicking on a user in search results
  - Possible issues:
    - API endpoint returning unexpected shape
    - Null/undefined stats fields causing render errors
    - Modal state not properly initialized
  - Add defensive null checks for all user stats
  - Purpose: Fix the error preventing admin user management
  - _Leverage: Browser dev tools, existing modal component_
  - _Requirements: REQ-1_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: React Frontend Developer specializing in debugging

    Task: Debug and fix UserDetailModal error:

    Investigation steps:
    1. Check browser console for error messages when clicking user
    2. Verify API response shape matches UserDetails interface
    3. Check if any stats fields (trackCount, playlistCount, etc.) are null/undefined
    4. Verify modal state initialization

    Likely fixes:
    1. Add null checks/default values for stats: `user.trackCount ?? 0`
    2. Ensure API error handling displays properly
    3. Check if userId is properly passed to useUserDetails hook

    In UserDetailModal.tsx:
    - Add defensive rendering for all numeric fields
    - Ensure stats grid handles undefined values

    Restrictions:
    - Do not change the API contract
    - Maintain existing functionality when data is complete

    _Leverage: Existing UserDetailModal, useUserDetails hook

    Success:
    - Clicking a user in search results opens modal without error
    - Modal displays all available user data
    - Missing data shows graceful fallback (0 or "N/A")

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 3.2 Add tests for admin user modal
  - File: `frontend/src/components/admin/__tests__/UserDetailModal.test.tsx`
  - Test cases:
    - Modal renders with complete user data
    - Modal handles missing/null stats gracefully
    - Modal displays error when API fails
    - Modal closes properly
  - Purpose: Prevent regressions in admin user management
  - _Leverage: Existing admin component tests, mock patterns_
  - _Requirements: REQ-1_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: React Test Engineer

    Task: Create/extend frontend/src/components/admin/__tests__/UserDetailModal.test.tsx:

    Test cases:
    1. TestUserDetailModal_RendersUserInfo - displays name, email, role
    2. TestUserDetailModal_RendersStats - displays all stat values
    3. TestUserDetailModal_HandlesNullStats - graceful fallback for missing stats
    4. TestUserDetailModal_HandlesAPIError - shows error alert
    5. TestUserDetailModal_CloseButtonWorks - onClose called when closed
    6. TestUserDetailModal_LoadingState - shows spinner while loading

    Restrictions:
    - Mock useUserDetails hook appropriately
    - Test component in isolation

    _Leverage: Existing admin test patterns, @testing-library/react

    Success:
    - All test cases pass
    - Edge cases covered (null stats, API errors)
    - Component behavior fully tested

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

---

## Task Group 4: Integration and Verification

- [ ] 4.1 End-to-end verification of track visibility
  - Manual testing steps:
    1. Login as subscriber → verify only own tracks + public visible
    2. Login as admin → verify all tracks visible
    3. As subscriber, try to access private track URL directly → verify 403
  - Update any integration tests if needed
  - Purpose: Verify full flow of visibility enforcement
  - _Requirements: REQ-2, REQ-3, REQ-4_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: QA Engineer specializing in integration testing

    Task: Verify track visibility enforcement end-to-end:

    Manual verification:
    1. Deploy backend with visibility fixes
    2. Test as subscriber:
       - GET /api/v1/tracks - should not show other users' private tracks
       - GET /api/v1/tracks/:id for own track - should succeed
       - GET /api/v1/tracks/:id for other's public track - should succeed
       - GET /api/v1/tracks/:id for other's private track - should get 403
    3. Test as admin:
       - GET /api/v1/tracks - should show ALL tracks
       - GET /api/v1/tracks/:id for any track - should succeed

    Document any issues found and create follow-up tasks if needed.

    Success:
    - All scenarios behave as expected
    - No unauthorized data leakage
    - Error messages are clear

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [ ] 4.2 End-to-end verification of guest restrictions
  - Manual testing steps:
    1. Visit site without logging in
    2. Verify dashboard displays
    3. Click navigation links → verify redirected to permission-denied
    4. Verify permission-denied page has working buttons
  - Purpose: Verify guest route protection works correctly
  - _Requirements: REQ-5_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: QA Engineer specializing in frontend testing

    Task: Verify guest route restrictions end-to-end:

    Manual verification:
    1. Open site in incognito/private browser (no session)
    2. Verify dashboard (/) loads and shows public content
    3. Click "Tracks" in navigation - verify redirect to /permission-denied
    4. Click "Albums" in navigation - verify redirect to /permission-denied
    5. On permission-denied page:
       - Verify message displays
       - Click "Go to Dashboard" - verify navigation to /
       - Click "Sign In" - verify navigation to /login
    6. Log in, then verify all routes accessible

    Document any issues found.

    Success:
    - Guest users blocked from all non-public routes
    - Permission denied page works correctly
    - Authenticated users unaffected

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_

- [x] 4.3 Update CLAUDE.md and documentation
  - Files: `backend/internal/service/CLAUDE.md`, `frontend/CLAUDE.md`
  - Document new visibility enforcement behavior
  - Document guest route restrictions
  - Update API documentation for 403 responses
  - Purpose: Keep documentation current with implementation
  - _Requirements: All_
  - _Prompt: |
    Implement the task for spec access-control-bug-fixes, first run spec-workflow-guide to get the workflow guide then implement the task:

    Role: Technical Writer

    Task: Update documentation to reflect access control changes:

    Backend CLAUDE.md updates:
    1. Document GetTrack visibility behavior (hasGlobal parameter)
    2. Document ListTracks visibility filtering
    3. Note 403 vs 404 response distinction

    Frontend CLAUDE.md updates:
    1. Document permission-denied route
    2. Document guest route restrictions (PUBLIC_ROUTES)
    3. Update hooks section to note auth requirements

    Restrictions:
    - Keep existing documentation format
    - Be concise but complete

    Success:
    - Documentation accurately reflects new behavior
    - New developers can understand access control from docs

    After implementation:
    1. Mark task as in-progress in tasks.md: change [ ] to [-]
    2. After completion, log implementation with log-implementation tool
    3. Mark task complete: change [-] to [x]_
