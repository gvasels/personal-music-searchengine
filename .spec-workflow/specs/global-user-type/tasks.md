# Tasks: Global User Type

## Overview

Implementation tasks for the Global User Type feature, organized by phase following the migration strategy from the design document.

---

## Phase 1: Backend Models & Types

- [ ] 1.1 Create role types and permissions in `backend/internal/models/role.go`
  - File: `backend/internal/models/role.go`
  - Define `UserRole` type (guest, subscriber, artist, admin)
  - Define `Permission` type and constants
  - Create `RolePermissions` map
  - Add `CognitoGroupName()` method
  - Purpose: Establish role-based access control foundation
  - _Leverage: `backend/internal/models/common.go` patterns_
  - _Requirements: REQ-1_

- [ ] 1.2 Create `PlaylistVisibility` type and update Playlist model
  - File: `backend/internal/models/playlist.go`
  - Add `PlaylistVisibility` enum (private, unlisted, public)
  - Replace `IsPublic bool` with `Visibility PlaylistVisibility`
  - Add `CreatorName`, `CreatorAvatar` fields for denormalization
  - Update `NewPlaylistItem` for GSI2 (public discovery)
  - Update `PlaylistResponse` and `ToResponse()`
  - Purpose: Enable playlist visibility options
  - _Leverage: Existing playlist model_
  - _Requirements: REQ-2_

- [ ] 1.3 Update User model with Role and FollowingCount
  - File: `backend/internal/models/user.go`
  - Add `Role UserRole` field (default: subscriber)
  - Add `FollowingCount int` field
  - Remove `Tier SubscriptionTier` field
  - Update `UserResponse` and `ToResponse()`
  - Purpose: Replace subscription tier with role system
  - _Leverage: Existing user model_
  - _Requirements: REQ-1_

- [ ] 1.4 Create ArtistProfile model
  - File: `backend/internal/models/artist_profile.go`
  - Define `ArtistProfile` struct with all fields from design
  - Create `ArtistProfileItem` with DynamoDB keys (PK, SK, GSI1, GSI2)
  - Add `NewArtistProfileItem()` function
  - Create request/response DTOs
  - Purpose: Enable artist profile management
  - _Leverage: `backend/internal/models/artist.go` patterns_
  - _Requirements: REQ-3_

- [ ] 1.5 Create Follow model
  - File: `backend/internal/models/follow.go`
  - Define `Follow` struct
  - Create `FollowItem` with DynamoDB keys for both access patterns
  - Add `NewFollowItem()` function
  - Purpose: Enable follow system data layer
  - _Leverage: `backend/internal/models/common.go` patterns_
  - _Requirements: REQ-4_

- [ ] 1.6 Remove SubscriptionTier system
  - Files: `backend/internal/models/feature.go`, `backend/internal/models/subscription.go`
  - Remove `SubscriptionTier` type (keep for migration reference)
  - Remove `TierConfig` and related functions
  - Update `FeatureFlag` to not depend on tier
  - Purpose: Clean up replaced subscription system
  - _Requirements: REQ-1 (replaces tier)_

---

## Phase 2: Backend Repository Layer

- [ ] 2.1 Add ArtistProfile repository methods
  - File: `backend/internal/repository/artist_profile.go`
  - Implement `CreateArtistProfile()`
  - Implement `GetArtistProfile()` by ID
  - Implement `GetArtistProfileByUserID()` via GSI1
  - Implement `UpdateArtistProfile()`
  - Implement `LinkArtistToProfile()` with uniqueness check via GSI2
  - Implement `ListArtistProfiles()` for discovery
  - Purpose: Data access for artist profiles
  - _Leverage: `backend/internal/repository/repository.go` patterns_
  - _Requirements: REQ-3_

- [ ] 2.2 Add Follow repository methods
  - File: `backend/internal/repository/follow.go`
  - Implement `CreateFollow()`
  - Implement `DeleteFollow()`
  - Implement `GetFollow()` to check if following
  - Implement `ListFollowers()` via GSI1
  - Implement `ListFollowing()` by PK prefix
  - Implement `IncrementFollowerCount()` / `DecrementFollowerCount()`
  - Purpose: Data access for follow relationships
  - _Leverage: `backend/internal/repository/repository.go` patterns_
  - _Requirements: REQ-4_

- [ ] 2.3 Update Playlist repository for visibility
  - File: `backend/internal/repository/playlist.go`
  - Update `CreatePlaylist()` to use Visibility
  - Update `UpdatePlaylist()` to handle visibility changes
  - Add `ListPublicPlaylists()` via GSI2
  - Update item creation to set GSI2 for public playlists
  - Purpose: Enable public playlist discovery
  - _Leverage: Existing playlist repository_
  - _Requirements: REQ-2_

- [ ] 2.4 Update User repository for Role
  - File: `backend/internal/repository/user.go`
  - Update `CreateUser()` to set default role
  - Add `UpdateUserRole()`
  - Update `GetUser()` to include role
  - Purpose: Persist role in DynamoDB
  - _Leverage: Existing user repository_
  - _Requirements: REQ-1_

- [ ] 2.5 Add repository interface updates
  - File: `backend/internal/repository/repository.go`
  - Add `ArtistProfileRepository` interface
  - Add `FollowRepository` interface
  - Update existing interfaces as needed
  - Purpose: Define repository contracts
  - _Requirements: All_

---

## Phase 3: Backend Services

- [ ] 3.1 Create RoleService
  - File: `backend/internal/service/role_service.go`
  - Implement `GetUserRole()` - extract from JWT claims
  - Implement `SetUserRole()` - update Cognito group + DynamoDB
  - Implement `HasPermission()` - check role permissions map
  - Add Cognito Admin API integration
  - Purpose: Role management and permission checking
  - _Leverage: `backend/internal/service/service.go` patterns_
  - _Requirements: REQ-1_

- [ ] 3.2 Create ArtistProfileService
  - File: `backend/internal/service/artist_profile_service.go`
  - Implement `CreateProfile()` - require artist role
  - Implement `GetProfile()`, `GetProfileByUserID()`
  - Implement `UpdateProfile()` - owner only
  - Implement `LinkToArtist()` - claim catalog artist with uniqueness check
  - Implement `GetProfileWithCatalog()` - include linked artist data
  - Purpose: Artist profile business logic
  - _Leverage: Existing service patterns_
  - _Requirements: REQ-3_

- [ ] 3.3 Create FollowService
  - File: `backend/internal/service/follow_service.go`
  - Implement `Follow()` - create follow + increment count
  - Implement `Unfollow()` - delete follow + decrement count
  - Implement `GetFollowers()`, `GetFollowing()` with pagination
  - Implement `IsFollowing()`
  - Add self-follow prevention
  - Purpose: Follow system business logic
  - _Leverage: Existing service patterns_
  - _Requirements: REQ-4_

- [ ] 3.4 Update PlaylistService for visibility
  - File: `backend/internal/service/playlist_service.go`
  - Update `CreatePlaylist()` to set default visibility
  - Add `UpdateVisibility()` method
  - Add `ListPublicPlaylists()` for discovery
  - Update access checks for visibility
  - Purpose: Public playlist functionality
  - _Leverage: Existing playlist service_
  - _Requirements: REQ-2_

---

## Phase 4: Backend Handlers & Middleware

- [ ] 4.1 Create authorization middleware
  - File: `backend/internal/handlers/middleware/auth.go`
  - Implement `RequireRole()` middleware
  - Implement `RequireAuth()` middleware
  - Implement `OptionalAuth()` middleware
  - Extract role from `cognito:groups` JWT claim
  - Purpose: Role-based endpoint protection
  - _Leverage: Existing auth middleware_
  - _Requirements: REQ-1_

- [ ] 4.2 Create role management handlers
  - File: `backend/internal/handlers/role_handler.go`
  - Implement `GET /api/v1/users/:id/role` - admin only
  - Implement `PUT /api/v1/users/:id/role` - admin only
  - Purpose: Admin role management API
  - _Leverage: Existing handler patterns_
  - _Requirements: REQ-1_

- [ ] 4.3 Create artist profile handlers
  - File: `backend/internal/handlers/artist_profile_handler.go`
  - Implement `POST /api/v1/artist-profiles` - artist role
  - Implement `GET /api/v1/artist-profiles/:id` - public
  - Implement `PUT /api/v1/artist-profiles/:id` - owner
  - Implement `POST /api/v1/artist-profiles/:id/link` - owner
  - Implement `GET /api/v1/artist-profiles/:id/catalog` - public
  - Implement `GET /api/v1/artist-profiles/discover` - public
  - Purpose: Artist profile API endpoints
  - _Leverage: Existing handler patterns_
  - _Requirements: REQ-3_

- [ ] 4.4 Create follow handlers
  - File: `backend/internal/handlers/follow_handler.go`
  - Implement `POST /api/v1/artist-profiles/:id/follow` - subscriber+
  - Implement `DELETE /api/v1/artist-profiles/:id/follow` - subscriber+
  - Implement `GET /api/v1/artist-profiles/:id/followers` - public
  - Implement `GET /api/v1/users/me/following` - subscriber+
  - Purpose: Follow system API endpoints
  - _Leverage: Existing handler patterns_
  - _Requirements: REQ-4_

- [ ] 4.5 Update playlist handlers for visibility
  - File: `backend/internal/handlers/playlist_handler.go`
  - Add `GET /api/v1/playlists/public` - public
  - Add `PUT /api/v1/playlists/:id/visibility` - owner
  - Update access checks for visibility
  - Purpose: Public playlist API endpoints
  - _Leverage: Existing playlist handlers_
  - _Requirements: REQ-2_

- [ ] 4.6 Register new routes
  - File: `backend/cmd/api/main.go`
  - Register all new handlers
  - Apply appropriate middleware to routes
  - Purpose: Wire up new endpoints
  - _Requirements: All_

---

## Phase 5: Cognito & Infrastructure

- [ ] 5.1 Create Cognito groups via OpenTofu
  - File: `infrastructure/shared/cognito.tf`
  - Add `admin` group resource
  - Add `artist` group resource
  - Add `subscriber` group resource
  - Purpose: Create role groups in Cognito
  - _Leverage: Existing Cognito configuration_
  - _Requirements: REQ-1_

- [ ] 5.2 Create admin bootstrap script
  - File: `scripts/bootstrap-admin.sh`
  - Accept email parameter
  - Look up user by email in Cognito
  - Add user to admin group
  - Update DynamoDB user role
  - Purpose: Bootstrap first admin user
  - _Requirements: REQ-1_

- [ ] 5.3 Update Lambda authorizer for groups
  - File: `infrastructure/backend/api-gateway.tf` or Lambda code
  - Ensure `cognito:groups` claim is included in context
  - Purpose: Pass role info to handlers
  - _Leverage: Existing authorizer_
  - _Requirements: REQ-1_

---

## Phase 6: Data Migration

- [ ] 6.1 Create playlist visibility migration script
  - File: `scripts/migrate-playlist-visibility.go`
  - Scan all playlists
  - Convert `IsPublic: true` → `Visibility: public`
  - Convert `IsPublic: false` → `Visibility: private`
  - Add GSI2 keys for public playlists
  - Purpose: Migrate existing playlist data
  - _Requirements: REQ-2_

- [ ] 6.2 Create user role migration script
  - File: `scripts/migrate-user-roles.go`
  - Scan all users
  - Set `Role: subscriber` for all existing users
  - Add users to subscriber Cognito group
  - Set admin role for gvasels90@gmail.com
  - Purpose: Migrate existing user data
  - _Requirements: REQ-1_

---

## Phase 7: Frontend Updates

- [ ] 7.1 Create role types and hooks
  - Files: `frontend/src/lib/api/types.ts`, `frontend/src/hooks/useRole.ts`
  - Add `UserRole` type
  - Add `Permission` type
  - Create `useRole()` hook to extract role from auth
  - Create `useHasPermission()` hook
  - Purpose: Frontend role utilities
  - _Leverage: Existing hooks patterns_
  - _Requirements: REQ-1_

- [ ] 7.2 Update playlist components for visibility
  - Files: `frontend/src/components/playlist/*.tsx`
  - Add visibility selector (private/unlisted/public)
  - Show creator info on public playlists
  - Add public playlist discovery page
  - Purpose: Playlist visibility UI
  - _Leverage: Existing playlist components_
  - _Requirements: REQ-2_

- [ ] 7.3 Create artist profile components
  - Files: `frontend/src/components/artist-profile/*.tsx`
  - Create `ArtistProfileCard` component
  - Create `ArtistProfilePage` component
  - Create `EditArtistProfileModal` component
  - Add catalog linking UI
  - Purpose: Artist profile UI
  - _Leverage: Existing component patterns_
  - _Requirements: REQ-3_

- [ ] 7.4 Create follow components
  - Files: `frontend/src/components/follow/*.tsx`
  - Create `FollowButton` component
  - Create `FollowersList` component
  - Create `FollowingList` component
  - Purpose: Follow system UI
  - _Leverage: Existing component patterns_
  - _Requirements: REQ-4_

- [ ] 7.5 Create artist profile hooks and API
  - Files: `frontend/src/hooks/useArtistProfile.ts`, `frontend/src/lib/api/artistProfiles.ts`
  - Add API functions for all endpoints
  - Create `useArtistProfile()` hook
  - Create `useArtistProfiles()` hook for discovery
  - Purpose: Artist profile data fetching
  - _Leverage: Existing TanStack Query patterns_
  - _Requirements: REQ-3_

- [ ] 7.6 Create follow hooks and API
  - Files: `frontend/src/hooks/useFollow.ts`, `frontend/src/lib/api/follows.ts`
  - Add API functions for follow/unfollow
  - Create `useFollow()` mutation hook
  - Create `useFollowers()` and `useFollowing()` query hooks
  - Create `useIsFollowing()` hook
  - Purpose: Follow system data fetching
  - _Leverage: Existing TanStack Query patterns_
  - _Requirements: REQ-4_

- [ ] 7.7 Add routes for new pages
  - Files: `frontend/src/routes/artist-profiles/*.tsx`, `frontend/src/routes/playlists/public.tsx`
  - Add `/artist-profiles` route
  - Add `/artist-profiles/$profileId` route
  - Add `/playlists/public` route
  - Purpose: Navigation to new features
  - _Leverage: Existing TanStack Router patterns_
  - _Requirements: REQ-2, REQ-3_

- [ ] 7.8 Remove subscription tier UI
  - Files: Various frontend files
  - Remove tier display from user profile
  - Remove subscription-related components
  - Remove tier-based feature gating UI
  - Purpose: Clean up replaced tier system
  - _Requirements: REQ-1 (replaces tier)_

---

## Phase 8: Testing

- [ ] 8.1 Backend unit tests for models
  - Files: `backend/internal/models/*_test.go`
  - Test role permission mappings
  - Test ArtistProfile and Follow model functions
  - Test PlaylistVisibility handling
  - Purpose: Model layer test coverage
  - _Requirements: All_

- [ ] 8.2 Backend unit tests for services
  - Files: `backend/internal/service/*_test.go`
  - Test RoleService with mocked Cognito
  - Test ArtistProfileService
  - Test FollowService
  - Test updated PlaylistService
  - Purpose: Service layer test coverage
  - _Requirements: All_

- [ ] 8.3 Backend integration tests for handlers
  - Files: `backend/internal/handlers/*_test.go`
  - Test role-based endpoint access
  - Test artist profile CRUD endpoints
  - Test follow endpoints
  - Test public playlist endpoints
  - Purpose: API endpoint test coverage
  - _Requirements: All_

- [ ] 8.4 Frontend unit tests
  - Files: `frontend/src/**/*.test.tsx`
  - Test role hooks
  - Test artist profile components
  - Test follow components
  - Test visibility components
  - Purpose: Frontend test coverage
  - _Requirements: All_

---

## Task Dependencies

```
Phase 1 (Models) → Phase 2 (Repository) → Phase 3 (Services) → Phase 4 (Handlers)
                                                                      ↓
Phase 5 (Cognito) ────────────────────────────────────────────────────┤
                                                                      ↓
Phase 6 (Migration) ←─────────────────────────────────────────────────┤
                                                                      ↓
Phase 7 (Frontend) ←──────────────────────────────────────────────────┘
                                                                      ↓
Phase 8 (Testing) ←───────────────────────────────────────────────────┘
```

---

## Estimated Task Count

| Phase | Tasks | Description |
|-------|-------|-------------|
| 1 | 6 | Backend models |
| 2 | 5 | Repository layer |
| 3 | 4 | Service layer |
| 4 | 6 | Handlers & middleware |
| 5 | 3 | Cognito & infrastructure |
| 6 | 2 | Data migration |
| 7 | 8 | Frontend updates |
| 8 | 4 | Testing |
| **Total** | **38** | |

---

## Future Enhancements (Roadmap)

### User Management Architecture Refactor

**Current State**: Admin user search queries Cognito directly, with role info stored in both Cognito groups and DynamoDB. This creates data synchronization challenges.

**Target State**: DynamoDB becomes the source of truth for user data, with Cognito sync as needed for authentication.

#### Tasks:
- [ ] F.1 Create DynamoDB user profile on Cognito post-confirmation trigger
  - Add Lambda trigger for `PostConfirmation_ConfirmSignUp`
  - Create user profile in DynamoDB with default role (subscriber)
  - Sync display name from Cognito attributes
  - Purpose: Ensure all Cognito users have DynamoDB profiles

- [ ] F.2 Migrate admin search from Cognito to DynamoDB
  - Update `AdminService.SearchUsers` to query DynamoDB only
  - Add GSI for email prefix search if needed
  - Remove Cognito ListUsers dependency
  - Purpose: Simplify architecture, reduce Cognito API calls

- [ ] F.3 Create Cognito sync service for role changes
  - On role change in DynamoDB, sync to Cognito groups
  - Handle group membership atomically with DynamoDB updates
  - Add retry/rollback logic for consistency
  - Purpose: Keep Cognito groups in sync for JWT claims

- [ ] F.4 Backfill existing Cognito users to DynamoDB
  - Migration script to scan Cognito users
  - Create DynamoDB profiles for users without them
  - Preserve existing data where profiles exist
  - Purpose: One-time migration to complete architecture
