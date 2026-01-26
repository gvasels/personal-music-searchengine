# Tasks Document: Admin Panel & Track Visibility

## Phase 1: Backend Models & Types

- [x] 1.1 Extend Track model with visibility field
  - File: `backend/internal/models/track.go`
  - Add `Visibility TrackVisibility` field to Track struct
  - Add `PublishedAt *time.Time` field
  - Add `OwnerDisplayName string` for API responses
  - Add helper methods: `IsPubliclyAccessible()`, `IsDiscoverable()`
  - Purpose: Enable track-level visibility control
  - _Leverage: `backend/internal/models/visibility.go` (playlist visibility pattern)_
  - _Requirements: 4.1, 6.1, 6.2_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Extend Track model with visibility field following playlist visibility pattern from models/visibility.go. Add Visibility, PublishedAt, OwnerDisplayName fields and helper methods. | Restrictions: Do not modify existing Track fields, maintain backward compatibility, use same visibility constants as playlists | Success: Track model compiles, has visibility field with proper JSON/DynamoDB tags, helper methods work correctly | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 1.2 Create admin models (UserSummary, UserDetails)
  - File: `backend/internal/models/admin.go`
  - Define `UserSummary` struct for search results
  - Define `UserDetails` struct with full user info including counts
  - Define `UpdateRoleRequest` and `UpdateStatusRequest` DTOs
  - Purpose: Data structures for admin API responses
  - _Leverage: `backend/internal/models/user.go`_
  - _Requirements: 1.4, 5.1, 5.2_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create admin models file with UserSummary, UserDetails, and request DTOs as specified in design document | Restrictions: Follow existing model patterns, use proper JSON tags, include validation tags | Success: Models compile, follow project conventions, include all fields from design | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 2: Backend Services (TDD)

- [x] 2.1 Write CognitoClient tests
  - File: `backend/internal/service/cognito_client_test.go`
  - Test AddUserToGroup, RemoveUserFromGroup success cases
  - Test GetUserGroups returns correct groups
  - Test DisableUser/EnableUser operations
  - Test error handling for AWS errors
  - Purpose: TDD - write failing tests first
  - _Leverage: `backend/internal/service/*_test.go` patterns_
  - _Requirements: 1.5, 1.6_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Write comprehensive unit tests for CognitoClient interface methods using mocked AWS SDK. Tests must fail initially (TDD Red phase). | Restrictions: Use testify/mock for AWS SDK, do not call real Cognito, test both success and error paths | Success: Tests compile, cover all interface methods, fail because implementation doesn't exist | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 2.2 Implement CognitoClient
  - File: `backend/internal/service/cognito_client.go`
  - Implement CognitoClient interface using AWS SDK v2
  - Add AdminAddUserToGroup, AdminRemoveUserFromGroup operations
  - Add AdminListGroupsForUser operation
  - Add AdminDisableUser, AdminEnableUser operations
  - Purpose: TDD Green - make tests pass
  - _Leverage: AWS SDK v2 cognito-idp client patterns_
  - _Requirements: 1.5, 1.6_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with AWS expertise | Task: Implement CognitoClient using AWS SDK v2 cognito-idp client to make tests pass (TDD Green phase) | Restrictions: Use context for cancellation, handle all AWS errors properly, make all tests from 2.1 pass | Success: All tests from 2.1 pass, implementation follows AWS SDK best practices | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 2.3 Write AdminService tests
  - File: `backend/internal/service/admin_test.go`
  - Test SearchUsers with various queries
  - Test GetUserDetails returns correct data
  - Test UpdateUserRole updates both DynamoDB and Cognito
  - Test UpdateUserRole rollback on Cognito failure
  - Test SetUserStatus enable/disable flow
  - Purpose: TDD - write failing tests first
  - _Leverage: `backend/internal/service/*_test.go` patterns, mock repository_
  - _Requirements: 1.3, 1.4, 1.5, 1.6, 1.7_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Write comprehensive unit tests for AdminService with mocked CognitoClient and Repository. Include rollback test for atomic operations. Tests must fail initially. | Restrictions: Mock all dependencies, test both success and failure scenarios, verify atomic behavior | Success: Tests compile, cover all service methods, fail because implementation doesn't exist | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 2.4 Implement AdminService
  - File: `backend/internal/service/admin.go`
  - Implement AdminService interface
  - SearchUsers: query DynamoDB with partial match
  - GetUserDetails: fetch user with content counts
  - UpdateUserRole: atomic update to DynamoDB + Cognito with rollback
  - SetUserStatus: update disabled flag in DynamoDB + Cognito
  - Purpose: TDD Green - make tests pass
  - _Leverage: `backend/internal/service/role.go`, `backend/internal/repository/`_
  - _Requirements: 1.3, 1.4, 1.5, 1.6, 1.7_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Implement AdminService with atomic DynamoDB+Cognito operations and rollback on failure (TDD Green phase) | Restrictions: Must implement rollback pattern, use transactions where possible, make all tests from 2.3 pass | Success: All tests from 2.3 pass, atomic operations work correctly with rollback | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 2.5 Write track visibility filter tests
  - File: `backend/internal/service/track_test.go` (extend existing)
  - Test admin user sees all tracks
  - Test GlobalReaders user sees all tracks
  - Test regular user sees only own tracks + public tracks
  - Test visibility filter logic
  - Purpose: TDD - extend existing tests
  - _Leverage: existing track service tests_
  - _Requirements: 3.1, 3.2, 3.3, 4.2, 4.3, 4.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Extend track service tests with visibility filtering tests. Test admin, global, and regular user track visibility scenarios. | Restrictions: Use existing test patterns, mock repository, test edge cases | Success: New tests compile, cover visibility scenarios, fail because logic not implemented | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 2.6 Implement track visibility filtering in TrackService
  - File: `backend/internal/service/track.go` (modify existing)
  - Modify ListTracks to filter by visibility based on user role
  - Add IncludePublic filter option
  - Query GSI3 for public tracks and merge with user's tracks
  - Include OwnerDisplayName in responses for admin/global users
  - Purpose: TDD Green - make visibility tests pass
  - _Leverage: existing TrackService, playlist visibility patterns_
  - _Requirements: 3.1, 3.2, 3.3, 4.2, 4.3, 4.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Modify TrackService.ListTracks to implement visibility filtering based on user role. Merge own tracks with public tracks for regular users. | Restrictions: Maintain backward compatibility, efficient queries, make visibility tests pass | Success: All visibility tests pass, regular users see own+public, admins see all | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 3: Backend Handlers & Repository

- [x] 3.1 Add GSI3 to DynamoDB for public tracks
  - File: `infrastructure/shared/dynamodb.tf`
  - Add GSI3 with GSI3PK (partition) and GSI3SK (sort)
  - Purpose: Enable efficient public track discovery
  - _Leverage: existing GSI patterns in dynamodb.tf_
  - _Requirements: 4.2, 6.3_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer | Task: Add GSI3 to DynamoDB table for public track discovery. Follow existing GSI patterns. | Restrictions: Do not modify existing GSIs, use projection ALL, follow naming conventions | Success: GSI3 added to terraform, tofu plan shows expected changes | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 3.2 Update repository for track visibility
  - File: `backend/internal/repository/track.go` (modify existing)
  - Add methods to query public tracks via GSI3
  - Update SaveTrack to set GSI3PK/GSI3SK when public
  - Add method to update track visibility
  - Purpose: Data access for visibility feature
  - _Leverage: existing repository patterns, playlist GSI2 pattern_
  - _Requirements: 4.2, 6.3, 6.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Update track repository with GSI3 queries for public tracks and visibility update methods. Follow playlist GSI2 pattern. | Restrictions: Maintain existing method signatures, efficient queries, proper GSI key management | Success: Repository compiles, GSI3 queries work, visibility updates correctly | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 3.3 Create admin handlers
  - File: `backend/internal/handlers/admin.go`
  - Implement SearchUsers handler (GET /api/v1/admin/users)
  - Implement GetUserDetails handler (GET /api/v1/admin/users/:id)
  - Implement UpdateUserRole handler (PUT /api/v1/admin/users/:id/role)
  - Implement UpdateUserStatus handler (PUT /api/v1/admin/users/:id/status)
  - All handlers require RoleAdmin middleware
  - Purpose: HTTP layer for admin API
  - _Leverage: `backend/internal/handlers/handlers.go` patterns, middleware_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create admin handlers following existing handler patterns. All endpoints require admin role. | Restrictions: Use existing error handling, validate inputs, follow handler conventions | Success: Handlers compile, return correct responses, reject non-admin users | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 3.4 Register admin routes
  - File: `backend/internal/handlers/handlers.go` (modify existing)
  - Add admin route group with RequireRole(RoleAdmin) middleware
  - Register all admin handlers
  - Purpose: Wire up admin API routes
  - _Leverage: existing route registration patterns_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Register admin routes in handlers.go with admin middleware protection. | Restrictions: Follow existing route patterns, use proper middleware chain | Success: Routes registered, admin middleware applied, non-admin requests rejected | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 3.5 Add track visibility update handler
  - File: `backend/internal/handlers/track.go` (modify existing)
  - Add UpdateTrackVisibility handler (PUT /api/v1/tracks/:id/visibility)
  - Validate user owns the track (or is admin)
  - Purpose: Allow artists to change track visibility
  - _Leverage: playlist visibility handler pattern_
  - _Requirements: 6.2, 6.3_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Add track visibility update handler following playlist visibility pattern. | Restrictions: Check ownership or admin, validate visibility values, follow existing patterns | Success: Handler works, ownership checked, visibility updates correctly | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 4: Infrastructure

- [x] 4.1 Add admin API Gateway routes
  - File: `infrastructure/backend/api-gateway.tf`
  - Add routes for all admin endpoints
  - Use JWT authorizer
  - Purpose: Expose admin API via API Gateway
  - _Leverage: existing route patterns in api-gateway.tf_
  - _Requirements: 5.1, 5.2, 5.3, 5.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer | Task: Add API Gateway routes for admin endpoints with JWT authorizer. | Restrictions: Follow existing route patterns, use same integration | Success: tofu plan shows new routes, routes configured correctly | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 4.2 Add Cognito admin IAM permissions
  - File: `infrastructure/backend/lambda.tf` (modify existing)
  - Add cognito-idp:Admin* permissions to API Lambda role
  - Scope to Cognito user pool ARN
  - Purpose: Allow Lambda to manage Cognito users
  - _Leverage: existing IAM patterns_
  - _Requirements: 1.5, 1.6_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer | Task: Add Cognito admin IAM permissions to API Lambda role for user group management. | Restrictions: Least privilege, scope to user pool ARN only | Success: tofu plan shows IAM changes, permissions are minimal | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 4.3 Add track visibility route
  - File: `infrastructure/backend/api-gateway.tf`
  - Add PUT /api/v1/tracks/{id}/visibility route
  - Purpose: Expose track visibility endpoint
  - _Leverage: existing route patterns_
  - _Requirements: 6.2_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer | Task: Add track visibility update route to API Gateway. | Restrictions: Follow existing patterns, use JWT authorizer | Success: Route added, tofu plan clean | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 5: Frontend API & Hooks

- [x] 5.1 Create admin API client
  - File: `frontend/src/lib/api/admin.ts`
  - Implement searchUsers, getUserDetails functions
  - Implement updateUserRole, updateUserStatus functions
  - Purpose: API layer for admin functionality
  - _Leverage: `frontend/src/lib/api/client.ts` patterns_
  - _Requirements: 5.1, 5.2, 5.3, 5.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create admin API client following existing API patterns. | Restrictions: Use apiClient, proper TypeScript types, handle errors | Success: API functions compile, types match backend responses | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 5.2 Create useAdmin hook
  - File: `frontend/src/hooks/useAdmin.ts`
  - Create useSearchUsers query with debounce
  - Create useUserDetails query
  - Create useUpdateUserRole mutation
  - Create useUpdateUserStatus mutation
  - Purpose: React Query hooks for admin operations
  - _Leverage: existing hooks patterns, query key factories_
  - _Requirements: 1.3, 1.4, 1.5, 1.6_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create useAdmin hook with TanStack Query patterns for search, details, and mutations. | Restrictions: Use query key factory, proper error handling, optimistic updates for mutations | Success: Hooks compile, follow existing patterns, queries/mutations work | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 5.3 Update Track type with visibility
  - File: `frontend/src/types/index.ts` (modify existing)
  - Add visibility field to Track interface
  - Add ownerDisplayName optional field
  - Purpose: TypeScript types for track visibility
  - _Leverage: existing PlaylistVisibility type_
  - _Requirements: 6.1_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Add visibility and ownerDisplayName fields to Track type. | Restrictions: Maintain backward compatibility, use existing PlaylistVisibility type | Success: Types compile, Track has visibility field | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 5.4 Add showUploadedBy to preferences store
  - File: `frontend/src/lib/store/preferencesStore.ts` (modify existing)
  - Add showUploadedBy: boolean (default false)
  - Add setShowUploadedBy action
  - Purpose: User preference for track list column
  - _Leverage: existing preferences pattern_
  - _Requirements: 2.1, 2.2, 2.3_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Add showUploadedBy preference to store with localStorage persistence. | Restrictions: Follow existing preference patterns, default to false | Success: Preference persists, toggles correctly | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 6: Frontend Components

- [x] 6.1 Create UserSearchForm component
  - File: `frontend/src/components/admin/UserSearchForm.tsx`
  - Search input with debounce (300ms)
  - Display loading state
  - Purpose: Search users by email/ID
  - _Leverage: existing SearchBar patterns_
  - _Requirements: 1.3_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create UserSearchForm with debounced input following SearchBar patterns. | Restrictions: Use DaisyUI, 300ms debounce, accessible | Success: Component renders, search triggers on debounce | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 6.2 Create UserCard component
  - File: `frontend/src/components/admin/UserCard.tsx`
  - Display user summary (email, name, role, status)
  - Click to select user
  - Show disabled badge if disabled
  - Purpose: Display user in search results
  - _Leverage: existing card components_
  - _Requirements: 1.3, 1.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create UserCard component displaying user summary with click handler. | Restrictions: Use DaisyUI card, show role badge, disabled state styling | Success: Component renders user info, click works | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 6.3 Create UserDetailModal component
  - File: `frontend/src/components/admin/UserDetailModal.tsx`
  - Display full user details
  - Include RoleSelector dropdown
  - Include enable/disable toggle
  - Confirmation dialog for changes
  - Purpose: View and edit user details
  - _Leverage: existing modal patterns, EditArtistProfileModal_
  - _Requirements: 1.4, 1.5, 1.6_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create UserDetailModal with role selector and status toggle. Include confirmation dialog. | Restrictions: Follow modal patterns, show loading states, handle errors | Success: Modal shows user details, role/status changes work | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 6.4 Create admin components barrel export
  - File: `frontend/src/components/admin/index.ts`
  - Export all admin components
  - Purpose: Clean imports
  - _Leverage: existing barrel exports_
  - _Requirements: N/A_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create barrel export for admin components. | Restrictions: Follow existing patterns | Success: Components importable from index | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 7: Frontend Pages & Routes

- [x] 7.1 Create admin users page
  - File: `frontend/src/routes/admin/users.tsx`
  - Combine UserSearchForm, UserCard list, UserDetailModal
  - Protect with admin role check (redirect non-admins)
  - Purpose: Main admin user management page
  - _Leverage: existing route patterns, useAuth_
  - _Requirements: 1.1, 1.2_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Create admin users page at /admin/users with search, results, and detail modal. Protect with admin check. | Restrictions: Redirect non-admins to home with toast, follow route patterns | Success: Page renders for admin, redirects non-admin, search and edit work | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 7.2 Add admin link to sidebar (admin only)
  - File: `frontend/src/components/layout/Sidebar.tsx` (modify existing)
  - Add "Admin" section visible only to admins
  - Link to /admin/users
  - Purpose: Navigation to admin panel
  - _Leverage: existing sidebar patterns, useAuth_
  - _Requirements: 1.1_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Add admin section to sidebar visible only to admin users. | Restrictions: Use useAuth.isAdmin, follow existing nav patterns | Success: Admin link shows for admins only | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 7.3 Add "Show Uploaded By" toggle to settings
  - File: `frontend/src/routes/settings.tsx` (modify existing)
  - Add toggle in "Track List" section
  - Connect to preferencesStore
  - Purpose: User preference for uploaded by column
  - _Leverage: existing settings toggles_
  - _Requirements: 2.2_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Add showUploadedBy toggle to settings page in Track List section. | Restrictions: Follow existing toggle patterns, immediate effect | Success: Toggle works, preference persists | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 7.4 Add "Uploaded By" column to TrackList
  - File: `frontend/src/components/library/TrackList.tsx` (modify existing)
  - Add optional "Uploaded By" column
  - Show/hide based on showUploadedBy preference
  - Display "You" for own tracks, display name for others
  - Purpose: Track attribution display
  - _Leverage: existing column patterns, preferencesStore_
  - _Requirements: 2.1, 2.4, 2.5_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Add conditional Uploaded By column to TrackList based on preference. Show "You" for own tracks. | Restrictions: Maintain existing columns, conditional rendering, follow patterns | Success: Column shows when enabled, displays correct values | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 7.5 Add track visibility selector to track detail
  - File: `frontend/src/routes/tracks/$trackId.tsx` (modify existing)
  - Add VisibilitySelector for track owner
  - Add mutation to update visibility
  - Purpose: Allow artists to change track visibility
  - _Leverage: `frontend/src/components/playlist/VisibilitySelector.tsx`_
  - _Requirements: 6.2_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer | Task: Add VisibilitySelector to track detail page for track owner. Use existing component from playlist. | Restrictions: Only show for track owner, reuse existing component | Success: Selector shows for owner, visibility updates work | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 8: Testing & Documentation

- [ ] 8.1 Write admin handler integration tests
  - File: `backend/internal/handlers/admin_test.go`
  - Test all admin endpoints
  - Test permission checks (non-admin rejected)
  - Purpose: Verify admin API works end-to-end
  - _Leverage: existing handler test patterns_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Write integration tests for admin handlers covering all endpoints and permission checks. | Restrictions: Follow existing test patterns, mock services | Success: Tests pass, cover success and permission scenarios | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [ ] 8.2 Write frontend admin component tests
  - File: `frontend/src/components/admin/__tests__/`
  - Test UserSearchForm, UserCard, UserDetailModal
  - Purpose: Verify admin components work correctly
  - _Leverage: existing component test patterns_
  - _Requirements: 1.1, 1.3, 1.4_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Test Engineer | Task: Write component tests for admin components using React Testing Library. | Restrictions: Follow existing test patterns, mock API | Success: Tests pass, cover key interactions | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 8.3 Create admin components CLAUDE.md
  - File: `frontend/src/components/admin/CLAUDE.md`
  - Document all admin components
  - Include usage examples
  - Purpose: Documentation per HARD REQUIREMENT #4
  - _Leverage: existing component CLAUDE.md files_
  - _Requirements: N/A (documentation)_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer | Task: Create CLAUDE.md for admin components following existing patterns. | Restrictions: Follow existing CLAUDE.md format | Success: Documentation complete, follows project standards | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [ ] 8.4 Update CHANGELOG.md
  - File: `CHANGELOG.md`, `backend/CHANGELOG.md`
  - Document admin panel feature
  - Document track visibility feature
  - Purpose: Documentation per HARD REQUIREMENT #4
  - _Leverage: existing CHANGELOG format_
  - _Requirements: N/A (documentation)_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer | Task: Update CHANGELOG.md files with admin panel and track visibility features. | Restrictions: Follow Keep a Changelog format | Success: Changelog updated, follows format | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

## Phase 9: Deployment & Migration

- [ ] 9.1 Deploy infrastructure changes
  - Run tofu apply for shared (GSI3) and backend (routes, IAM)
  - Purpose: Deploy infrastructure for new feature
  - _Leverage: existing deployment process_
  - _Requirements: All infrastructure_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer | Task: Deploy infrastructure changes using tofu apply for shared and backend layers. | Restrictions: Review plan before apply, backup state | Success: GSI3 created, routes added, IAM updated | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [x] 9.2 Create track visibility migration script
  - File: `scripts/migrations/migrate-track-visibility.sh`
  - Set Visibility="private" for all existing tracks
  - Idempotent (safe to run multiple times)
  - Purpose: Backfill visibility for existing tracks
  - _Leverage: existing migration scripts pattern_
  - _Requirements: 4.5_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer | Task: Create migration script to set visibility=private for existing tracks following existing migration patterns. | Restrictions: Idempotent, use AWS CLI, paginate for large tables | Success: Script runs, sets visibility on all tracks | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_

- [ ] 9.3 Run migrations and deploy
  - Run track visibility migration
  - Deploy backend Lambda
  - Deploy frontend
  - Purpose: Complete deployment
  - _Leverage: existing deployment process_
  - _Requirements: All_
  - _Prompt: Implement the task for spec admin-panel-track-visibility, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer | Task: Run migrations and deploy backend/frontend. | Restrictions: Run migration before deploying new code | Success: Migration complete, backend/frontend deployed, feature working | Instructions: Mark task [-] in tasks.md before starting, use log-implementation tool after completion, mark [x] when done_
