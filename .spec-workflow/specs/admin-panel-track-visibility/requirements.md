# Requirements Document: Admin Panel & Track Visibility

## Introduction

This feature adds administrative capabilities and enhanced track visibility controls to the Personal Music Search Engine. It enables administrators to manage users directly from the platform, introduces an "Uploaded By" column for track attribution, and implements role-based track visibility where admins see all content while regular users see only their own tracks plus any public tracks.

## Alignment with Product Vision

This feature supports the multi-user platform vision established in the global-user-type feature by:
- Enabling platform administrators to manage users without AWS Console access
- Providing content attribution for a multi-user environment
- Implementing proper access control for a shared music platform

## Requirements

### Requirement 1: Admin User Management Page

**User Story:** As an admin, I want to search for users and manage their roles/status, so that I can moderate the platform without accessing AWS Console.

#### Acceptance Criteria

1. WHEN an admin navigates to `/admin/users` THEN the system SHALL display the admin user management page
2. IF a non-admin user attempts to access `/admin/users` THEN the system SHALL redirect to home with an access denied message
3. WHEN an admin searches for a user by email or user ID THEN the system SHALL return matching users with their current role and status
4. WHEN an admin selects a user THEN the system SHALL display the user's profile details including: email, display name, current role, creation date, last login, and content counts (tracks, playlists)
5. WHEN an admin changes a user's role THEN the system SHALL update both DynamoDB AND the Cognito user group membership
6. WHEN an admin disables a user THEN the system SHALL mark the user as disabled in DynamoDB (not delete)
7. WHEN Cognito group update fails THEN the system SHALL rollback the DynamoDB change and display an error

### Requirement 2: Track List "Uploaded By" Column

**User Story:** As a user, I want to optionally see who uploaded each track, so that I can understand content attribution when viewing shared/public tracks.

#### Acceptance Criteria

1. WHEN the track list is displayed THEN the "Uploaded By" column SHALL be hidden by default
2. WHEN a user enables "Show Uploaded By" in settings THEN the track list SHALL display the uploader's display name
3. WHEN the preference changes THEN the system SHALL persist it to localStorage via the preferences store
4. IF the uploader has no display name THEN the system SHALL show "Unknown" as fallback
5. WHEN viewing own tracks THEN the "Uploaded By" column SHALL show "You" instead of the user's name

### Requirement 3: Admin/Global Track Visibility

**User Story:** As an admin or global reader, I want to see all tracks on the platform, so that I can moderate content and assist users.

#### Acceptance Criteria

1. WHEN an admin user lists tracks THEN the system SHALL return ALL tracks from ALL users
2. WHEN a user with GlobalReaders group lists tracks THEN the system SHALL return ALL tracks from ALL users
3. WHEN a regular user (subscriber/artist) lists tracks THEN the system SHALL return only their own tracks
4. WHEN tracks are returned for admin/global users THEN each track SHALL include the owner's user ID and display name

### Requirement 4: Public Track Visibility

**User Story:** As a subscriber, I want to see public tracks from other users, so that I can discover new music.

#### Acceptance Criteria

1. WHEN a track is marked as "public" THEN any authenticated user SHALL be able to view it in listings
2. WHEN a regular user lists tracks THEN the system SHALL return: their own tracks + all public tracks from other users
3. WHEN a track's visibility is "private" THEN only the owner and admins SHALL see it
4. WHEN a track's visibility is "unlisted" THEN only users with direct link and admins SHALL access it
5. IF no visibility is set on a track THEN the system SHALL default to "private"

### Requirement 5: Backend Admin API Endpoints

**User Story:** As the frontend, I need API endpoints to manage users, so that the admin panel can function.

#### Acceptance Criteria

1. WHEN `GET /api/v1/admin/users?search=query` is called by admin THEN the system SHALL return matching users
2. WHEN `GET /api/v1/admin/users/:userId` is called by admin THEN the system SHALL return full user details
3. WHEN `PUT /api/v1/admin/users/:userId/role` is called by admin with `{"role": "artist"}` THEN the system SHALL update role in DynamoDB and Cognito
4. WHEN `PUT /api/v1/admin/users/:userId/status` is called by admin with `{"disabled": true}` THEN the system SHALL disable the user
5. IF any admin endpoint is called by non-admin THEN the system SHALL return 403 Forbidden

### Requirement 6: Track Visibility Field

**User Story:** As a track owner, I want to set visibility on my tracks, so that I can control who sees my content.

#### Acceptance Criteria

1. WHEN a track is created THEN the visibility field SHALL default to "private"
2. WHEN an artist sets track visibility to "public" THEN the track SHALL appear in all users' listings
3. WHEN the visibility is updated THEN the system SHALL update the track's GSI entry for discoverability
4. WHEN listing tracks for regular users THEN the filter SHALL combine: `(OwnerID = currentUser) OR (Visibility = "public")`

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Admin handlers separate from user handlers
- **Modular Design**: Reuse existing VisibilitySelector component for tracks
- **Dependency Management**: Cognito operations isolated in dedicated service
- **Clear Interfaces**: Admin service interface clearly defined

### Performance
- User search SHALL complete within 500ms for tables up to 10,000 users
- Track listing with visibility filter SHALL not add more than 50ms latency
- Admin operations SHALL have optimistic UI updates with rollback on failure

### Security
- All admin endpoints SHALL require `RoleAdmin` permission
- Cognito operations SHALL use AWS SDK with proper IAM permissions
- User disable SHALL NOT delete data, only prevent access
- Audit log SHALL record all admin role changes

### Reliability
- Cognito and DynamoDB updates SHALL be atomic (rollback on partial failure)
- Admin panel SHALL gracefully handle Cognito service unavailability
- Track visibility filter SHALL default to "private" if field is missing

### Usability
- Admin search SHALL support partial email/name matching
- Role changes SHALL show confirmation dialog
- Track visibility selector SHALL use same UI pattern as playlists
- Settings toggle SHALL take effect immediately without page reload
