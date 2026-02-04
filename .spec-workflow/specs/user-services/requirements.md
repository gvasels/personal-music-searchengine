# Requirements: User Services

## Introduction

This feature establishes DynamoDB as the source of truth for user data, with automatic synchronization from Cognito. When users sign up, a DynamoDB profile is automatically created. When roles/groups change, both systems stay in sync. User settings and permissions are managed through DynamoDB, making the admin panel and user management more efficient.

This replaces the current architecture where admin searches query Cognito directly, creating a more maintainable and performant system.

## Alignment with Product Vision

This feature completes the Global User Type foundation by:
- Ensuring all users have consistent DynamoDB profiles
- Enabling efficient user search without Cognito API calls
- Providing a foundation for user settings and preferences
- Supporting future features like notifications, preferences, and audit logging

---

## Requirements

### REQ-1: Automatic User Creation on Signup

**User Story:** As a platform operator, I want user profiles automatically created in DynamoDB when users sign up so that all users have consistent data records.

#### Acceptance Criteria

1. WHEN a user completes Cognito signup (post-confirmation) THEN the system SHALL create a DynamoDB user profile
2. The created profile SHALL include:
   - `userId` (Cognito sub)
   - `email` (from Cognito attributes)
   - `displayName` (from Cognito name attribute, or email prefix as fallback)
   - `role` (default: `subscriber`)
   - `createdAt` (timestamp)
   - `settings` (default settings object)
3. IF the user profile already exists THEN the system SHALL NOT overwrite existing data
4. IF DynamoDB write fails THEN the system SHALL log the error but NOT block signup
5. WHEN a user is created THEN the system SHALL add them to the `subscriber` Cognito group

---

### REQ-2: Role Change Synchronization

**User Story:** As an admin, I want role changes to automatically sync between DynamoDB and Cognito so that permissions are consistent across the system.

#### Acceptance Criteria

1. WHEN an admin changes a user's role via the admin panel THEN the system SHALL:
   - Update the role in DynamoDB
   - Update Cognito group membership (remove old, add new)
   - Both updates happen atomically (rollback on failure)

2. WHEN a role change succeeds THEN the system SHALL:
   - Return success to the admin
   - The user's next token refresh reflects the new role

3. IF Cognito group update fails THEN the system SHALL:
   - Rollback the DynamoDB change
   - Return an error to the admin
   - Log the failure for debugging

4. WHEN checking a user's role THEN the system SHALL read from DynamoDB (source of truth)

---

### REQ-3: User Settings Management

**User Story:** As a user, I want my preferences and settings stored persistently so that my experience is consistent across sessions and devices.

#### Acceptance Criteria

1. WHEN a user is created THEN the system SHALL initialize default settings:
   ```json
   {
     "theme": "system",
     "notifications": {
       "email": true,
       "push": false
     },
     "privacy": {
       "showActivity": true,
       "allowFollows": true
     },
     "player": {
       "autoplay": true,
       "crossfade": 0,
       "normalizeVolume": false
     }
   }
   ```

2. WHEN a user updates a setting THEN the system SHALL persist it to DynamoDB
3. WHEN a user logs in THEN the system SHALL load their settings
4. IF settings are corrupted/missing THEN the system SHALL use defaults
5. Settings SHALL be accessible via `GET /api/v1/users/me/settings`
6. Settings SHALL be updatable via `PATCH /api/v1/users/me/settings`

---

### REQ-4: Admin User Search (DynamoDB-powered)

**User Story:** As an admin, I want to search users efficiently without hitting Cognito API limits so that user management is fast and reliable.

#### Acceptance Criteria

1. WHEN an admin searches users THEN the system SHALL query DynamoDB (not Cognito)
2. Search SHALL support:
   - Email prefix matching
   - Display name contains matching
   - Role filtering
   - Disabled status filtering
3. Search results SHALL include:
   - User ID, email, display name
   - Role, disabled status
   - Created date, last login date
4. Search SHALL return results within 500ms for up to 1000 users
5. Search SHALL support pagination with cursor-based navigation

---

### REQ-5: User Status Management

**User Story:** As an admin, I want to enable/disable users with the change reflected immediately so that I can manage access effectively.

#### Acceptance Criteria

1. WHEN an admin disables a user THEN the system SHALL:
   - Set `disabled: true` in DynamoDB
   - Disable the user in Cognito (prevents login)
   - Active sessions continue until token expiry

2. WHEN an admin enables a user THEN the system SHALL:
   - Set `disabled: false` in DynamoDB
   - Enable the user in Cognito

3. WHEN a disabled user tries to access the API THEN the system SHALL return 403
4. The admin panel SHALL show real-time disabled status from DynamoDB

---

### REQ-6: Backfill Migration

**User Story:** As a platform operator, I want existing Cognito users backfilled to DynamoDB so that all users have consistent profiles.

#### Acceptance Criteria

1. A migration script SHALL scan all Cognito users
2. For each Cognito user without a DynamoDB profile, the script SHALL create one
3. The script SHALL preserve existing DynamoDB data (no overwrites)
4. The script SHALL be idempotent (safe to run multiple times)
5. The script SHALL log progress and any errors
6. The script SHALL support dry-run mode

---

## Non-Functional Requirements

### Code Architecture and Modularity
- **Lambda Trigger**: Separate Lambda function for Cognito post-confirmation
- **Service Layer**: UserService handles all user operations
- **Repository Pattern**: User repository abstracts DynamoDB operations
- **Atomic Operations**: Role changes use DynamoDB transactions where possible

### Performance
- User creation trigger SHALL complete within 3 seconds
- User search SHALL return within 500ms
- Settings load SHALL complete within 100ms
- GSI on email for efficient prefix searches

### Security
- Only admins can modify other users' roles
- Users can only modify their own settings
- Cognito remains the authentication authority
- DynamoDB is the authorization/data authority
- All user mutations logged for audit

### Reliability
- Cognito signup succeeds even if DynamoDB write fails (graceful degradation)
- Role changes are atomic (both succeed or both fail)
- Settings corruption falls back to defaults
- Migration script handles partial failures

### Usability
- Settings changes reflect immediately
- Role changes reflect on next page load
- Admin search is responsive and filterable

---

## Out of Scope

- **User deletion**: Soft delete vs hard delete strategy (future spec)
- **Audit logging**: Detailed action logging (future enhancement)
- **User merge**: Merging duplicate accounts
- **SSO integration**: Third-party identity providers
- **Rate limiting**: API rate limits per user

---

## Dependencies

- Existing Cognito User Pool
- Existing DynamoDB table with User entity
- Admin panel (from Global User Type)
- CognitoClient service (existing)

---

## Technical Notes

### DynamoDB User Entity (Updated)

```
PK: USER#{userId}
SK: PROFILE
GSI1PK: EMAIL#{email}  (for email search)
GSI1SK: USER

Attributes:
- userId: string
- email: string
- displayName: string
- avatarUrl: string (optional)
- role: UserRole
- disabled: boolean
- settings: map
- createdAt: string (ISO)
- updatedAt: string (ISO)
- lastLoginAt: string (ISO, optional)
```

### Lambda Trigger Configuration

```
Cognito User Pool → Post Confirmation → Lambda
                 → Pre Token Generation → Lambda (add custom claims)
```
