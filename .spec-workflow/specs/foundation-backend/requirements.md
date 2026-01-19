# Requirements - Foundation Backend

## Overview

Foundation backend components for the Personal Music Search Engine. Establishes the core domain models, data access layer, business logic services, and HTTP handlers for the Go backend running on AWS Lambda.

## User Stories

### US-1.1: Domain Models
**As a** developer
**I want** well-defined domain models with proper DynamoDB key patterns
**So that** I can implement the single-table design consistently

**Acceptance Criteria:**
- [x] Common types defined (EntityType, UploadStatus, AudioFormat, Timestamps)
- [x] Track model with JSON tags and DynamoDB item conversion
- [x] Album model with artist aggregation support
- [x] User model with storage tracking
- [x] Playlist and PlaylistTrack models with position ordering
- [x] Tag and TrackTag models for user tagging
- [x] Upload model with status tracking
- [x] All models have ToResponse() methods for API conversion
- [x] Unit tests with 80%+ coverage

### US-1.2: Repository Layer
**As a** developer
**I want** a repository layer for DynamoDB and S3 operations
**So that** data access is abstracted from business logic

**Acceptance Criteria:**
- [ ] Repository interface defined
- [ ] DynamoDB repository implementation
- [ ] S3 repository for media operations
- [ ] Track CRUD operations
- [ ] Album operations with artist queries
- [ ] User profile operations
- [ ] Playlist operations with track ordering
- [ ] Tag operations with track associations
- [ ] Upload status tracking
- [ ] Integration tests with DynamoDB Local

### US-1.3: Service Layer
**As a** developer
**I want** a service layer for business logic
**So that** handlers are kept thin and logic is testable

**Acceptance Criteria:**
- [ ] TrackService for track management
- [ ] AlbumService for album aggregation
- [ ] UserService for profile management
- [ ] PlaylistService for playlist operations
- [ ] TagService for tagging operations
- [ ] UploadService for presigned URLs and status tracking
- [ ] StreamService for CloudFront signed URL generation
- [ ] SearchService for Nixiesearch integration
- [ ] Unit tests with mocked repositories

### US-1.4: HTTP Handlers
**As a** developer
**I want** Echo HTTP handlers for the API
**So that** the Lambda can serve REST requests

**Acceptance Criteria:**
- [ ] Track handlers (GET/PUT/DELETE)
- [ ] Album handlers (GET list, GET by ID)
- [ ] User handlers (GET profile, PUT update)
- [ ] Playlist handlers (CRUD + track management)
- [ ] Tag handlers (CRUD + track associations)
- [ ] Upload handlers (presigned URL, confirm, status)
- [ ] Stream handlers (stream/download URLs)
- [ ] Search handlers (GET query, POST advanced)
- [ ] HTTP tests with mocked services
- [ ] NO custom auth middleware (API Gateway Cognito handles this)

## Non-Functional Requirements

- **Authentication**: API Gateway Cognito JWT Authorizer (NOT custom middleware)
- **Test Coverage**: Minimum 80% for all packages
- **TDD**: All code must be written test-first
- **Documentation**: CLAUDE.md for each package directory
