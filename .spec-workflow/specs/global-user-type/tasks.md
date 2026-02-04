# Tasks: Global User Type

## Overview

Implementation tasks for the Global User Type feature, organized by phase following the migration strategy from the design document.

**Last Updated**: 2026-01-26

---

## Implementation Status Summary

| Phase | Status | Completed | Total |
|-------|--------|-----------|-------|
| Phase 1: Backend Models | ✅ Complete | 6/6 | 100% |
| Phase 2: Repository Layer | ✅ Complete | 5/5 | 100% |
| Phase 3: Backend Services | ✅ Complete | 4/4 | 100% |
| Phase 4: Handlers & Middleware | ✅ Complete | 6/6 | 100% |
| Phase 5: Cognito & Infrastructure | ✅ Complete | 3/3 | 100% |
| Phase 6: Data Migration | ✅ Complete | 2/2 | 100% |
| Phase 7: Frontend Updates | ✅ Complete | 8/8 | 100% |
| Phase 8: Testing | 🔄 Partial | 2/4 | 50% |
| **Admin Panel (Added)** | ✅ Complete | 5/5 | 100% |

---

## Phase 1: Backend Models & Types ✅

- [x] 1.1 Create role types and permissions in `backend/internal/models/role.go`
  - File: `backend/internal/models/role.go`
  - ✅ Define `UserRole` type (guest, subscriber, artist, admin)
  - ✅ Define `Permission` type and constants
  - ✅ Create `RolePermissions` map
  - ✅ Add `CognitoGroupName()` method
  - ✅ Added `IsValid()` method for role validation

- [x] 1.2 Create `PlaylistVisibility` type and update Playlist model
  - File: `backend/internal/models/playlist.go`
  - ✅ Add `PlaylistVisibility` enum (private, unlisted, public)
  - ✅ Replace `IsPublic bool` with `Visibility PlaylistVisibility`
  - ✅ Add `CreatorName`, `CreatorAvatar` fields for denormalization
  - ✅ Update `NewPlaylistItem` for GSI2 (public discovery)
  - ✅ Update `PlaylistResponse` and `ToResponse()`

- [x] 1.3 Update User model with Role and FollowingCount
  - File: `backend/internal/models/user.go`
  - ✅ Add `Role UserRole` field (default: subscriber)
  - ✅ Add `FollowingCount int` field
  - ✅ Remove `Tier SubscriptionTier` field
  - ✅ Update `UserResponse` and `ToResponse()`
  - ✅ Add `ToUserDetails()` for admin panel

- [x] 1.4 Create ArtistProfile model
  - File: `backend/internal/models/artist_profile.go`
  - ✅ Define `ArtistProfile` struct with all fields from design
  - ✅ Create `ArtistProfileItem` with DynamoDB keys (PK, SK, GSI1, GSI2)
  - ✅ Add `NewArtistProfileItem()` function
  - ✅ Create request/response DTOs

- [x] 1.5 Create Follow model
  - File: `backend/internal/models/follow.go`
  - ✅ Define `Follow` struct
  - ✅ Create `FollowItem` with DynamoDB keys for both access patterns
  - ✅ Add `NewFollowItem()` function

- [x] 1.6 Remove SubscriptionTier system
  - Files: Various
  - ✅ Replaced SubscriptionTier with UserRole throughout
  - ✅ Updated feature flags to use role-based permissions
  - ✅ Removed tier-based configuration

---

## Phase 2: Backend Repository Layer ✅

- [x] 2.1 Add ArtistProfile repository methods
  - File: `backend/internal/repository/artist_profile.go`
  - ✅ Implement `CreateArtistProfile()`
  - ✅ Implement `GetArtistProfile()` by ID
  - ✅ Implement `GetArtistProfileByUserID()` via GSI1
  - ✅ Implement `UpdateArtistProfile()`
  - ✅ Implement `LinkArtistToProfile()` with uniqueness check via GSI2
  - ✅ Implement `ListArtistProfiles()` for discovery

- [x] 2.2 Add Follow repository methods
  - File: `backend/internal/repository/follow.go`
  - ✅ Implement `CreateFollow()`
  - ✅ Implement `DeleteFollow()`
  - ✅ Implement `GetFollow()` to check if following
  - ✅ Implement `ListFollowers()` via GSI1
  - ✅ Implement `ListFollowing()` by PK prefix
  - ✅ Implement `IncrementFollowerCount()` / `DecrementFollowerCount()`

- [x] 2.3 Update Playlist repository for visibility
  - File: `backend/internal/repository/playlist.go`
  - ✅ Update `CreatePlaylist()` to use Visibility
  - ✅ Update `UpdatePlaylist()` to handle visibility changes
  - ✅ Add `ListPublicPlaylists()` via GSI2
  - ✅ Update item creation to set GSI2 for public playlists

- [x] 2.4 Update User repository for Role
  - File: `backend/internal/repository/dynamodb.go`
  - ✅ Update `CreateUser()` to set default role
  - ✅ Add `UpdateUserRole()`
  - ✅ Update `GetUser()` to include role
  - ✅ Add `SearchUsers()` for admin panel
  - ✅ Add `SetUserDisabled()` for admin panel
  - ✅ Add `GetFollowerCount()` for user details

- [x] 2.5 Add repository interface updates
  - File: `backend/internal/repository/repository.go`
  - ✅ Add `ArtistProfileRepository` interface
  - ✅ Add `FollowRepository` interface
  - ✅ Add `AdminRepository` interface
  - ✅ Update existing interfaces as needed

---

## Phase 3: Backend Services ✅

- [x] 3.1 Create RoleService
  - File: `backend/internal/service/role.go`
  - ✅ Implement `GetUserRole()` - extract from JWT claims
  - ✅ Implement `SetUserRole()` - update Cognito group + DynamoDB
  - ✅ Implement `HasPermission()` - check role permissions map
  - ✅ Add Cognito Admin API integration
  - ✅ Tests in `role_test.go`

- [x] 3.2 Create ArtistProfileService
  - File: `backend/internal/service/artist_profile.go`
  - ✅ Implement `CreateProfile()` - require artist role
  - ✅ Implement `GetProfile()`, `GetProfileByUserID()`
  - ✅ Implement `UpdateProfile()` - owner only
  - ✅ Implement `LinkToArtist()` - claim catalog artist with uniqueness check
  - ✅ Implement `GetProfileWithCatalog()` - include linked artist data
  - ✅ Tests in `artist_profile_test.go`

- [x] 3.3 Create FollowService
  - File: `backend/internal/service/follow.go`
  - ✅ Implement `Follow()` - create follow + increment count
  - ✅ Implement `Unfollow()` - delete follow + decrement count
  - ✅ Implement `GetFollowers()`, `GetFollowing()` with pagination
  - ✅ Implement `IsFollowing()`
  - ✅ Add self-follow prevention
  - ✅ Tests in `follow_test.go`

- [x] 3.4 Update PlaylistService for visibility
  - File: `backend/internal/service/playlist.go`
  - ✅ Update `CreatePlaylist()` to set default visibility
  - ✅ Add `UpdateVisibility()` method
  - ✅ Add `ListPublicPlaylists()` for discovery
  - ✅ Update access checks for visibility

---

## Phase 4: Backend Handlers & Middleware ✅

- [x] 4.1 Create authorization middleware
  - File: `backend/internal/handlers/middleware/auth.go`
  - ✅ Implement `RequireRole()` middleware
  - ✅ Implement `RequireAuth()` middleware
  - ✅ Implement `OptionalAuth()` middleware
  - ✅ Extract role from `cognito:groups` JWT claim
  - ✅ Fix: Handle API Gateway array format `"[admin subscriber]"`
  - ✅ Tests in `auth_test.go`

- [x] 4.2 Create role management handlers
  - File: `backend/internal/handlers/role.go`
  - ✅ Implement `GET /api/v1/users/:id/role` - admin only
  - ✅ Implement `PUT /api/v1/users/:id/role` - admin only

- [x] 4.3 Create artist profile handlers
  - File: `backend/internal/handlers/artist_profile.go`
  - ✅ Implement `POST /api/v1/artists/entity` - artist role
  - ✅ Implement `GET /api/v1/artists/entity/:id` - public
  - ✅ Implement `PUT /api/v1/artists/entity/:id` - owner
  - ✅ Implement `POST /api/v1/artists/entity/:id/link` - owner
  - ✅ Implement `GET /api/v1/artists/entity/:id/catalog` - public

- [x] 4.4 Create follow handlers
  - File: `backend/internal/handlers/follow.go`
  - ✅ Implement `POST /api/v1/artists/entity/:id/follow` - subscriber+
  - ✅ Implement `DELETE /api/v1/artists/entity/:id/follow` - subscriber+
  - ✅ Implement `GET /api/v1/artists/entity/:id/followers` - public
  - ✅ Implement `GET /api/v1/users/me/following` - subscriber+

- [x] 4.5 Update playlist handlers for visibility
  - File: `backend/internal/handlers/playlist.go`
  - ✅ Add `GET /api/v1/playlists/public` - public
  - ✅ Add `PUT /api/v1/playlists/:id/visibility` - owner
  - ✅ Update access checks for visibility

- [x] 4.6 Register new routes
  - File: `backend/cmd/api/main.go`
  - ✅ Register all new handlers
  - ✅ Apply appropriate middleware to routes
  - ✅ Admin routes registered via code-based routing

---

## Phase 5: Cognito & Infrastructure ✅

- [x] 5.1 Create Cognito groups via OpenTofu
  - File: `infrastructure/shared/main.tf`
  - ✅ Add `admin` group resource
  - ✅ Add `artist` group resource
  - ✅ Add `subscriber` group resource
  - ✅ Add `GlobalReaders` group for cross-user content access

- [x] 5.2 Create admin bootstrap script
  - File: `scripts/bootstrap-admin.sh`
  - ✅ Accept email parameter
  - ✅ Look up user by email in Cognito
  - ✅ Add user to admin group
  - ✅ Update DynamoDB user role

- [x] 5.3 Update Lambda authorizer for groups
  - File: API Gateway configuration
  - ✅ Ensure `cognito:groups` claim is included in context
  - ✅ Groups passed as `"[group1 group2]"` format (handled in middleware)

---

## Phase 6: Data Migration ✅

- [x] 6.1 Create playlist visibility migration script
  - File: `scripts/migrations/migrate-playlist-visibility.sh`
  - ✅ Scan all playlists
  - ✅ Convert `IsPublic: true` → `Visibility: public`
  - ✅ Convert `IsPublic: false` → `Visibility: private`
  - ✅ Add GSI2 keys for public playlists

- [x] 6.2 Create user role migration script
  - File: `scripts/migrations/migrate-user-roles.sh`
  - ✅ Scan all users
  - ✅ Set `Role: subscriber` for all existing users
  - ✅ Add users to subscriber Cognito group
  - ✅ Set admin role for gvasels90@gmail.com

---

## Phase 7: Frontend Updates ✅

- [x] 7.1 Create role types and hooks
  - Files: `frontend/src/types/index.ts`, `frontend/src/hooks/useAuth.ts`
  - ✅ Add `UserRole` type
  - ✅ Add `Permission` type
  - ✅ Update `useAuth()` hook to extract role from JWT
  - ✅ Add `isAdmin`, `isArtist` properties

- [x] 7.2 Update playlist components for visibility
  - Files: `frontend/src/components/playlist/*.tsx`
  - ✅ Create `VisibilitySelector` component (private/unlisted/public)
  - ✅ Show creator info on public playlists
  - ✅ Add public playlist discovery page (`/playlists/public`)

- [x] 7.3 Create artist profile components
  - Files: `frontend/src/components/artist-profile/*.tsx`
  - ✅ Create `ArtistProfileCard` component
  - ✅ Create `EditArtistProfileModal` component
  - ✅ Add catalog linking UI

- [x] 7.4 Create follow components
  - Files: `frontend/src/components/follow/*.tsx`
  - ✅ Create `FollowButton` component
  - ✅ Create `FollowersList` component
  - ✅ Create `FollowingList` component

- [x] 7.5 Create artist profile hooks and API
  - Files: `frontend/src/hooks/useArtistProfiles.ts`, `frontend/src/lib/api/artistProfiles.ts`
  - ✅ Add API functions for all endpoints
  - ✅ Create `useArtistProfile()` hook
  - ✅ Create `useArtistProfiles()` hook for discovery

- [x] 7.6 Create follow hooks and API
  - Files: `frontend/src/hooks/useFollows.ts`, `frontend/src/lib/api/follows.ts`
  - ✅ Add API functions for follow/unfollow
  - ✅ Create `useFollow()` mutation hook
  - ✅ Create `useFollowers()` and `useFollowing()` query hooks
  - ✅ Create `useIsFollowing()` hook

- [x] 7.7 Add routes for new pages
  - Files: `frontend/src/routes/artists/entity/*.tsx`, `frontend/src/routes/playlists/public.tsx`
  - ✅ Add `/artists/entity` route
  - ✅ Add `/artists/entity/$artistId` route
  - ✅ Add `/playlists/public` route

- [x] 7.8 Remove subscription tier UI
  - Files: Various frontend files
  - ✅ Remove tier display from user profile
  - ✅ Remove subscription-related components
  - ✅ Replace tier-based feature gating with role-based

---

## Phase 8: Testing 🔄

- [x] 8.1 Backend unit tests for models
  - Files: `backend/internal/models/*_test.go`
  - ✅ Test role permission mappings
  - ✅ Test ArtistProfile and Follow model functions
  - ✅ Test PlaylistVisibility handling

- [x] 8.2 Backend unit tests for services
  - Files: `backend/internal/service/*_test.go`
  - ✅ Test RoleService (`role_test.go`)
  - ✅ Test ArtistProfileService (`artist_profile_test.go`)
  - ✅ Test FollowService (`follow_test.go`)

- [ ] 8.3 Backend integration tests for handlers
  - Files: `backend/internal/handlers/*_test.go`
  - ⬜ Test role-based endpoint access
  - ⬜ Test artist profile CRUD endpoints
  - ⬜ Test follow endpoints
  - ⬜ Test public playlist endpoints

- [ ] 8.4 Frontend unit tests
  - Files: `frontend/src/**/*.test.tsx`
  - ⬜ Test role hooks
  - ⬜ Test artist profile components (partial)
  - ⬜ Test follow components (partial)
  - ✅ Test VisibilitySelector component
  - ✅ Test FollowButton component
  - ✅ Test ArtistProfileCard component

---

## Admin Panel (Added Feature) ✅

- [x] A.1 Create AdminService
  - File: `backend/internal/service/admin.go`
  - ✅ Implement `SearchUsers()` - search Cognito users by email
  - ✅ Implement `GetUserDetails()` - full user details with status
  - ✅ Implement `UpdateUserRole()` - update DynamoDB + Cognito groups
  - ✅ Implement `UpdateUserRoleByAdmin()` - prevent self-modification
  - ✅ Implement `SetUserStatus()` - enable/disable users

- [x] A.2 Create CognitoClient
  - File: `backend/internal/service/cognito_client.go`
  - ✅ Implement `SearchUsers()` - list users by email filter
  - ✅ Implement `GetUserStatus()` - get enabled status
  - ✅ Implement `AddUserToGroup()` / `RemoveUserFromGroup()`
  - ✅ Implement `EnableUser()` / `DisableUser()`
  - ✅ Implement `GetUserGroups()`

- [x] A.3 Create admin handlers
  - File: `backend/internal/handlers/admin.go`
  - ✅ Implement `GET /api/v1/admin/users` - search users
  - ✅ Implement `GET /api/v1/admin/users/:id` - user details
  - ✅ Implement `PUT /api/v1/admin/users/:id/role` - update role
  - ✅ Implement `PUT /api/v1/admin/users/:id/status` - enable/disable

- [x] A.4 Create admin frontend components
  - Files: `frontend/src/components/admin/*.tsx`
  - ✅ Create `UserSearchForm` component
  - ✅ Create `UserCard` component
  - ✅ Create `UserDetailModal` component
  - ✅ Fix: Toggle shows ON when account is active (green)

- [x] A.5 Create admin route
  - File: `frontend/src/routes/admin/users.tsx`
  - ✅ Admin-only access (redirect non-admins)
  - ✅ Search users by email
  - ✅ View user details in modal
  - ✅ Change user roles with confirmation
  - ✅ Enable/disable users with confirmation

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

## Completed Task Count

| Phase | Tasks | Completed | Status |
|-------|-------|-----------|--------|
| 1 | 6 | 6 | ✅ |
| 2 | 5 | 5 | ✅ |
| 3 | 4 | 4 | ✅ |
| 4 | 6 | 6 | ✅ |
| 5 | 3 | 3 | ✅ |
| 6 | 2 | 2 | ✅ |
| 7 | 8 | 8 | ✅ |
| 8 | 4 | 2 | 🔄 |
| Admin | 5 | 5 | ✅ |
| **Total** | **43** | **41** | **95%** |

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

### GitHub Actions CI/CD

- [ ] F.5 Fix GitHub Actions deployment workflow
  - Currently failing on push to main
  - Need to review and fix workflow configuration
  - Purpose: Enable automated deployments
