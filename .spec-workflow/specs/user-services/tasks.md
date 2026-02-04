# Tasks: User Services

## Overview

This document breaks down the User Services feature into implementable tasks following the TDD workflow. This establishes DynamoDB as the source of truth for user data with Cognito synchronization.

---

## Task Groups

### Group 1: Models and Types (Backend)

#### Task 1.1: Add UserSettings Model
- **Description**: Create UserSettings model with nested types
- **Files**:
  - `backend/internal/models/user_settings.go`
  - `backend/internal/models/user_settings_test.go`
- **Acceptance Criteria**:
  - UserSettings struct with Theme, Notifications, Privacy, Player
  - NotificationSettings, PrivacySettings, PlayerSettings nested structs
  - DefaultUserSettings() returns proper defaults
  - Validation for theme values
  - JSON and DynamoDB tags

#### Task 1.2: Extend User Model with Settings
- **Description**: Add Settings field to User model
- **Files**:
  - `backend/internal/models/user.go` (modify)
  - `backend/internal/models/user_test.go` (modify)
- **Acceptance Criteria**:
  - User has Settings field
  - User has LastLoginAt field
  - User has count fields (TrackCount, etc.)
  - Backward compatible with existing data

---

### Group 2: Repository Layer (Backend)

#### Task 2.1: Add User Email Search to Repository
- **Description**: Implement SearchUsersByEmail using GSI1
- **Files**:
  - `backend/internal/repository/user_search.go`
  - `backend/internal/repository/user_search_test.go`
- **Acceptance Criteria**:
  - Query GSI1 with email prefix
  - Support pagination with cursor
  - Return PaginatedResult[User]
  - Handle empty results

#### Task 2.2: Add User Settings Repository Methods
- **Description**: CRUD operations for user settings
- **Files**:
  - `backend/internal/repository/user_settings.go`
  - `backend/internal/repository/user_settings_test.go`
- **Acceptance Criteria**:
  - GetUserSettings returns settings or defaults
  - UpdateUserSettings partial update support
  - Atomic update with UpdatedAt

---

### Group 3: Service Layer (Backend)

#### Task 3.1: Create UserService with Settings Management
- **Description**: Service for user operations with settings
- **Files**:
  - `backend/internal/service/user_service.go`
  - `backend/internal/service/user_service_test.go`
- **Acceptance Criteria**:
  - GetUserSettings(userID) method
  - UpdateUserSettings(userID, settings) method
  - Settings validation before save
  - Uses repository for data access

#### Task 3.2: Add CreateUserFromCognito Method
- **Description**: Create DynamoDB user from Cognito event data
- **Files**:
  - `backend/internal/service/user_service.go` (extend)
  - `backend/internal/service/user_service_test.go` (extend)
- **Acceptance Criteria**:
  - Creates user with default settings
  - Extracts email, name from Cognito attributes
  - Sets default role to subscriber
  - Idempotent (doesn't overwrite existing)

#### Task 3.3: Update AdminService for DynamoDB Search
- **Description**: Replace Cognito search with DynamoDB
- **Files**:
  - `backend/internal/service/admin_service.go` (modify)
  - `backend/internal/service/admin_service_test.go` (modify)
- **Acceptance Criteria**:
  - SearchUsers uses repository instead of Cognito
  - Supports email prefix and name search
  - Supports role and status filtering
  - Returns within 500ms for 1000 users

---

### Group 4: Lambda Trigger (Backend)

#### Task 4.1: Create Post-Confirmation Lambda Handler
- **Description**: Lambda triggered by Cognito post-confirmation
- **Files**:
  - `backend/cmd/triggers/post-confirmation/main.go`
  - `backend/cmd/triggers/post-confirmation/main_test.go`
- **Acceptance Criteria**:
  - Handles CognitoEventUserPoolsPostConfirmation
  - Creates DynamoDB user profile
  - Adds user to subscriber Cognito group
  - Returns event (doesn't block signup on error)
  - Logs errors for debugging

---

### Group 5: API Handlers (Backend)

#### Task 5.1: Add Settings API Endpoints
- **Description**: GET and PATCH endpoints for user settings
- **Files**:
  - `backend/internal/handlers/user_settings.go`
  - `backend/internal/handlers/user_settings_test.go`
- **Acceptance Criteria**:
  - GET /api/v1/users/me/settings
  - PATCH /api/v1/users/me/settings (partial update)
  - Validation errors return 400
  - Authentication required

#### Task 5.2: Register Settings Routes
- **Description**: Add routes to Echo router
- **Files**:
  - `backend/internal/handlers/routes.go` (modify)
- **Acceptance Criteria**:
  - Routes registered under /api/v1/users/me/settings
  - Auth middleware applied
  - Integrated with existing router

---

### Group 6: Infrastructure (OpenTofu)

#### Task 6.1: Add GSI1 to DynamoDB Table
- **Description**: Add email search GSI to existing table
- **Files**:
  - `infrastructure/shared/dynamodb.tf` (modify)
- **Acceptance Criteria**:
  - GSI1 with GSI1PK hash key, GSI1SK range key
  - Projection type ALL
  - On-demand billing mode preserved

#### Task 6.2: Create Post-Confirmation Lambda Infrastructure
- **Description**: OpenTofu for Lambda trigger
- **Files**:
  - `infrastructure/shared/cognito-triggers.tf` (new)
- **Acceptance Criteria**:
  - Lambda function resource
  - IAM role with DynamoDB and Cognito permissions
  - Lambda permission for Cognito invoke
  - User pool updated with lambda_config

---

### Group 7: Migration Script

#### Task 7.1: Create Backfill Migration Script
- **Description**: Script to backfill Cognito users to DynamoDB
- **Files**:
  - `scripts/migrations/backfill-cognito-users.go`
- **Acceptance Criteria**:
  - Lists all Cognito users with pagination
  - Creates DynamoDB profile if missing
  - Dry-run mode support
  - Progress logging
  - Idempotent operation

---

## Implementation Order

```
Phase 1 (Models):     1.1 → 1.2
Phase 2 (Repository): 2.1 → 2.2
Phase 3 (Service):    3.1 → 3.2 → 3.3
Phase 4 (Lambda):     4.1
Phase 5 (Handlers):   5.1 → 5.2
Phase 6 (Infra):      6.1 → 6.2
Phase 7 (Migration):  7.1
```

Each phase follows TDD: write failing tests first, then implement.

---

## Estimated Effort

| Group | Tasks | Complexity |
|-------|-------|------------|
| Models | 2 | Low |
| Repository | 2 | Medium |
| Service | 3 | Medium |
| Lambda | 1 | Medium |
| Handlers | 2 | Low |
| Infrastructure | 2 | Medium |
| Migration | 1 | Low |
| **Total** | **13 tasks** | |

---

## Dependencies

- Existing Cognito User Pool
- Existing DynamoDB table
- Existing CognitoClient service
- Existing User model and repository
