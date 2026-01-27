# Personal Music Search Engine - Epics and User Stories

## Overview

This plan maps out the epics and user stories for building a personal music library web application. The platform enables users to upload, organize, search, and stream their legally-owned music files.

**Frontend Technology Stack**: React 18 + TanStack Router + TanStack Query + Vite + TypeScript

**Backend Technology Stack**:
| Layer | Technology |
|-------|------------|
| CDN | CloudFront |
| API | API Gateway (OpenAPI spec) |
| Compute | AWS Lambda via Lambda Web Adapter |
| Framework | Echo (Go 1.22+) |
| Database | DynamoDB (single-table design) |
| Search | Nixiesearch (serverless) |
| Storage | S3 with Intelligent-Tiering |
| Auth | Cognito User Pool |

---

## Personas

| Persona | Description |
|---------|-------------|
| **Music Enthusiast** | Primary user who uploads and manages their music library |
| **Casual Listener** | User who primarily browses and plays music |

---

## Epic Overview

| Epic | Name | Phase | Priority | Status |
|------|------|-------|----------|--------|
| E1 | Foundation Infrastructure | 1 | Critical | ✅ Complete |
| E2 | Backend API | 2 | Critical | ✅ Complete |
| E3 | Search & Streaming | 3 | High | ✅ Complete |
| E4 | Tags & Playlists | 4 | High | ✅ Complete |
| E5 | Frontend | 5 | High | ✅ Complete |
| E6 | Distribution & Polish | 6 | Medium | ✅ Complete |
| E7 | Global User Type | 7 | High | ✅ Complete |
| E8 | Access Control & Bug Fixes | 8 | High | 🔄 In Progress |
| E9 | LocalStack Dev Environment | 9 | Medium | ✅ Complete |

---

## Epic 1: Foundation Infrastructure (Phase 1) ✅ **COMPLETE** (2026-01-20)

**Goal**: Establish core infrastructure and data contracts.

**Status**: Complete

**Spec**: [.spec-workflow/specs/foundation-backend/](.spec-workflow/specs/foundation-backend/)

### User Stories

#### US-1.1: Data Models & Contracts
**As a** Developer
**I want to** define Go domain models and OpenAPI specification
**So that** the backend has a consistent contract for implementation

**Acceptance Criteria**:
- [x] Go domain models defined (Track, Album, Artist, Playlist, User, Tag, Upload) _(9 model files in backend/internal/models/)_
- [x] OpenAPI 3.0 specification created _(backend/api/openapi.yaml)_
- [x] Go project structure initialized _(cmd/, internal/ structure)_
- [x] Test fixtures and mocks created _(test files alongside models)_

#### US-1.2: Global Infrastructure
**As a** Developer
**I want to** set up shared infrastructure resources
**So that** all services can use consistent state management and container registries

**Acceptance Criteria**:
- [x] S3 bucket for OpenTofu state with DynamoDB lock table _(music-library-prod-tofu-state)_
- [x] ECR repositories created for Lambda images _(infrastructure/global/)_
- [x] Base IAM roles for Lambda execution _(music-library-prod-lambda-execution)_

#### US-1.3: Shared Services
**As a** Developer
**I want to** deploy authentication, database, and storage services
**So that** the backend API can use them

**Acceptance Criteria**:
- [x] Cognito User Pool with app client and identity pool _(infrastructure/shared/)_
- [x] DynamoDB table with GSIs for single-table design _(MusicLibrary table)_
- [x] S3 media bucket with Intelligent-Tiering _(music-library-prod-media)_

---

## Epic 2: Backend API (Phase 2) ✅ **COMPLETE** (2026-01-20)

**Goal**: Implement the core API Lambda and upload processing pipeline.

**Status**: Complete

**Spec**: [.spec-workflow/specs/backend-api/](.spec-workflow/specs/backend-api/)

### User Stories

#### US-2.1: Core API Lambda
**As a** Music Enthusiast
**I want to** interact with my music library via API
**So that** I can manage my tracks, albums, and playlists

**Acceptance Criteria**:
- [x] Echo-based API Lambda scaffold _(backend/cmd/api/ with Lambda web adapter)_
- [x] DynamoDB repository layer _(backend/internal/repository/)_
- [x] User profile handlers _(backend/internal/handlers/user.go)_
- [x] Track CRUD handlers _(backend/internal/handlers/track.go)_
- [x] Album/artist handlers _(backend/internal/handlers/album.go)_

#### US-2.2: Upload & Processing Pipeline
**As a** Music Enthusiast
**I want to** upload music files and have metadata automatically extracted
**So that** my library is populated without manual data entry

**Acceptance Criteria**:
- [x] Presigned URL generation for S3 upload _(backend/internal/handlers/upload.go)_
- [x] Step Functions state machine for processing _(infrastructure/backend/step-functions.tf)_
- [x] Metadata extraction Lambda _(backend/cmd/processor/metadata/)_
- [x] Cover art extraction Lambda _(backend/cmd/processor/cover-art-processor/)_
- [x] Track creation Lambda _(backend/cmd/processor/track-creator/)_
- [x] File mover Lambda _(backend/cmd/processor/file-mover/)_
- [x] Search indexer Lambda _(backend/cmd/processor/search-indexer/)_
- [x] Status updater Lambda _(backend/cmd/processor/upload-status-updater/)_
- [x] Upload service integration with Step Functions _(backend/internal/service/upload.go)_

#### US-2.3: API Gateway
**As a** Developer
**I want to** expose the API via HTTP API Gateway
**So that** the frontend can communicate with the backend

**Acceptance Criteria**:
- [x] HTTP API with routes _(infrastructure/backend/api-gateway.tf)_
- [x] Cognito JWT authorizer configured _(Cognito authorizer attached to all routes)_
- [x] Lambda integrations configured _(Lambda proxy integrations)_
- [x] CORS configured for frontend _(Allow-Origin, Allow-Headers)_

---

## Epic 3: Search & Streaming (Phase 3) ✅ **COMPLETE** (2026-01-20)

**Goal**: Enable full-text search and audio streaming.

**Status**: Complete

**Spec**: [.spec-workflow/specs/search-streaming/](.spec-workflow/specs/search-streaming/)

### User Stories

#### US-3.1: Full-Text Search
**As a** Music Enthusiast
**I want to** search my music library by title, artist, album, or genre
**So that** I can quickly find songs I want to listen to

**Acceptance Criteria**:
- [x] Nixiesearch Docker image built _(Lucene-based search engine)_
- [x] Nixiesearch Lambda function deployed _(search-indexer Lambda)_
- [x] Index schema defined for music metadata _(title, artist, album, genre, tags)_
- [x] Indexer Lambda for scheduled re-indexing _(search-indexer Lambda via Step Functions)_
- [x] Search handlers implemented _(backend/internal/handlers/search.go)_

#### US-3.2: Audio Streaming
**As a** Music Enthusiast
**I want to** stream and download my music files
**So that** I can listen to my library from anywhere

**Acceptance Criteria**:
- [x] CloudFront distribution for media _(d8wn3lkytn5qe.cloudfront.net)_
- [x] CloudFront signed URLs configured _(backend/internal/service/stream.go)_
- [x] Stream/download handlers _(backend/internal/handlers/stream.go)_
- [x] S3 search index bucket _(music-library-prod-media)_

---

## Epic 4: Tags & Playlists (Phase 4) ✅ **COMPLETE** (2026-01-20)

**Goal**: Enable custom organization with tags and playlists.

**Status**: Complete

**Spec**: [.spec-workflow/specs/tags-playlists/](.spec-workflow/specs/tags-playlists/)

### User Stories

#### US-4.1: Custom Tags
**As a** Music Enthusiast
**I want to** create custom tags and apply them to tracks
**So that** I can organize my library my way

**Acceptance Criteria**:
- [x] Tag repository _(backend/internal/repository/ - tag operations)_
- [x] Tag handlers _(backend/internal/handlers/tag.go)_
- [x] Tag filtering in search _(search supports tag filtering)_

#### US-4.2: Playlists
**As a** Music Enthusiast
**I want to** create and manage playlists
**So that** I can curate collections of songs

**Acceptance Criteria**:
- [x] Playlist repository _(backend/internal/repository/ - playlist operations)_
- [x] Playlist handlers _(backend/internal/handlers/playlist.go)_
- [x] Playlist track ordering service _(position-based ordering with drag-and-drop)_

---

## Epic 5: Frontend (Phase 5) ✅ **COMPLETE** (2025-01-21)

**Goal**: Build the React web application.

**Status**: Complete

**Implementation Summary**: 351 unit tests passing, full TDD implementation across 5 waves.

### User Stories

#### US-5.1: Frontend Foundation
**As a** Developer
**I want to** set up the React project with proper tooling
**So that** UI development can proceed efficiently

**Acceptance Criteria**:
- [x] React + Vite + TypeScript project initialized _(React 18, Vite 5, TypeScript 5)_
- [x] TanStack Query + Router configured _(Query v5 with key factories, Router v1 with file-based routing)_
- [x] Tailwind + DaisyUI with custom themes _(DaisyUI 5, dark/light themes with #50c878 primary)_
- [x] Amplify auth integration _(Cognito JWT via AWS Amplify v6)_

#### US-5.2: Authentication & Layout
**As a** User
**I want to** sign in securely and navigate the application
**So that** I can access my music library

**Acceptance Criteria**:
- [x] Login/signup pages _(Login page with form validation, 19 tests)_
- [x] Protected routes _(Auth guard in __root.tsx with redirect to /login)_
- [x] App shell with navigation _(Header, Sidebar, Layout components with mobile support)_
- [x] Theme switcher (dark/light) _(Zustand themeStore with localStorage persistence)_

#### US-5.3: Library Views
**As a** Music Enthusiast
**I want to** browse my music by tracks, albums, or artists
**So that** I can explore my library visually

**Acceptance Criteria**:
- [x] Track list with sorting/filtering _(/tracks with TrackList, 10 route tests)_
- [x] Album grid view _(/albums with album cards, 11 route tests)_
- [x] Artist list view _(/artists, 11 route tests)_
- [x] Track detail/edit modal _(/tracks/$trackId with TagInput, 17 route tests)_

#### US-5.4: Upload & Search UI
**As a** Music Enthusiast
**I want to** upload files and search my library from the UI
**So that** I can manage my library without using the API directly

**Acceptance Criteria**:
- [x] Drag-drop file upload _(UploadDropzone with react-dropzone, 5 tests)_
- [x] Upload progress tracking _(useUpload hook with progress state, 6 tests)_
- [x] Search bar with autocomplete _(useSearch with searchAutocomplete, 8 tests)_
- [x] Search results page _(/search with debounced filters, 9 tests)_

#### US-5.5: Player & Queue
**As a** Music Enthusiast
**I want to** play music and manage a queue
**So that** I can listen to my library

**Acceptance Criteria**:
- [x] Audio player component _(PlayerBar with Howler.js, 5 tests)_
- [x] Play queue implementation _(Zustand playerStore with queue, repeat, shuffle, 14 tests)_
- [x] Playlist management UI _(/playlists with CreatePlaylistModal, usePlaylists, 20 route tests)_

#### US-5.6: Tags UI (Additional)
**As a** Music Enthusiast
**I want to** manage tags from the UI
**So that** I can organize my library

**Acceptance Criteria**:
- [x] Tag cloud page _(/tags with tag sizing by count, 8 tests)_
- [x] Tracks by tag page _(/tags/$tagName with play all, 9 tests)_
- [x] TagInput component _(add/remove tags with lowercase normalization, 6 tests)_

---

## Epic 6: Distribution & Polish (Phase 6) ✅ **COMPLETE** (2026-01-21)

**Goal**: Deploy frontend, implement CI/CD, and finalize documentation.

**Status**: Complete

**Implementation Summary**: Frontend hosting via S3/CloudFront, GitHub Actions CI/CD with OIDC, comprehensive documentation.

### User Stories

#### US-6.1: Frontend Hosting
**As a** User
**I want to** access the application via a URL
**So that** I can use it from any device

**Acceptance Criteria**:
- [x] S3 bucket for frontend _(music-library-prod-frontend with versioning, encryption, OAC)_
- [x] CloudFront distribution _(SPA routing, security headers, compression)_
- [x] Cache behaviors configured _(no-cache for index.html, 1yr for /assets/*)_

#### US-6.2: Testing & CI/CD
**As a** Developer
**I want to** have automated tests and deployment pipelines
**So that** code quality is maintained

**Acceptance Criteria**:
- [x] Go unit tests (80% coverage) _(ci.yml with coverage threshold check)_
- [x] API integration tests _(existing integration_test.go)_
- [x] Frontend component tests _(351 Vitest tests from Epic 5)_
- [x] GitHub Actions CI/CD _(ci.yml for PRs, deploy.yml for main)_

#### US-6.3: Documentation
**As a** Developer
**I want to** have comprehensive documentation
**So that** the project can be maintained and extended

**Acceptance Criteria**:
- [x] CLAUDE.md updated with project specifics _(infrastructure/frontend/CLAUDE.md)_
- [x] Deployment documentation _(docs/deployment.md)_
- [x] API documentation _(OpenAPI spec from Epic 1)_

---

## Wave Execution Summary

| Wave | Groups | Focus | Status |
|------|--------|-------|--------|
| 0 | 1 | Local: contracts, models, project structure | ✅ Complete |
| 1 | 2, 3 | Infrastructure: state, Cognito, DynamoDB, S3 | ✅ Complete |
| 2 | 4, 5, 6 | Backend: API Lambda, Step Functions, API Gateway | ✅ Complete |
| 3 | 7, 8, 9, 10 | Features: search, streaming, tags, playlists | ✅ Complete |
| 4 | 11-15 | Frontend: full React application | ✅ Complete |
| 5 | 16, 17, 18 | Polish: hosting, tests, documentation | ✅ Complete |

---

## Dependencies

```
Wave 0 ──► Wave 1 ──► Wave 2 ──► Wave 3 ──► Wave 4 ──► Wave 5
  │         │         │         │         │
  └─ Models └─ Infra  └─ API    └─ Features└─ Frontend
```

## Critical Path

1. **Step Functions State Machine** (US-2.2) - Required for upload processing
2. **API Gateway** (US-2.3) - Required for frontend-backend communication
3. **CloudFront Signed URLs** (US-3.2) - Required for streaming
4. **Frontend Hosting** (US-6.1) - Required for deployment

## Design Questions to Resolve

Before implementation, clarify:
1. S3 storage class configuration (Intelligent-Tiering lifecycle rules)
2. Caching strategy (CloudFront TTLs, DynamoDB DAX, in-memory?)
3. Search indexing approach (real-time via DynamoDB Streams vs batch?)
4. Error handling and retry patterns
5. Pagination cursor format (opaque vs structured?)

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Nixiesearch cold start latency | Medium | Use provisioned concurrency or fallback to DynamoDB |
| Large file uploads failing | High | Implement multipart upload with retry |
| Step Functions execution timeout | Medium | Monitor execution time, optimize Lambda code |
| Cover art extraction failures | Low | Continue processing without cover art, log for retry |

---

## Epic 7: Global User Type ✅ **COMPLETE** (2026-01-26)

**Goal**: Transform from single-user to multi-user platform with role-based access control.

**Status**: Complete

**Spec**: [.spec-workflow/specs/global-user-type/](.spec-workflow/specs/global-user-type/)

### User Stories

#### US-7.1: User Roles
**As a** Platform Administrator
**I want to** assign roles to users (admin, artist, subscriber, guest)
**So that** I can control access to features

**Acceptance Criteria**:
- [x] UserRole type (guest, subscriber, artist, admin) _(backend/internal/models/role.go)_
- [x] Cognito groups for roles _(admin, artist, subscriber groups)_
- [x] Role extraction from JWT _(backend/internal/handlers/middleware/auth.go)_
- [x] Permission system _(RolePermissions map)_

#### US-7.2: Public Playlists
**As a** Music Enthusiast
**I want to** make playlists public for others to discover
**So that** I can share my curated collections

**Acceptance Criteria**:
- [x] PlaylistVisibility enum (private, unlisted, public) _(models/visibility.go)_
- [x] Public playlist discovery via GSI2 _(repository queries)_
- [x] VisibilitySelector component _(frontend/src/components/playlist/)_

#### US-7.3: Artist Profiles
**As an** Artist
**I want to** create a profile and link it to my catalog
**So that** fans can follow me and see my work

**Acceptance Criteria**:
- [x] ArtistProfile model _(backend/internal/models/artist_profile.go)_
- [x] Profile CRUD service _(backend/internal/service/artist_profile.go)_
- [x] Catalog linking _(link ArtistProfile to existing Artist entries)_
- [x] Frontend components _(frontend/src/components/artist-profile/)_

#### US-7.4: Follow System
**As a** Music Enthusiast
**I want to** follow artists
**So that** I can stay updated on their activity

**Acceptance Criteria**:
- [x] Follow model _(backend/internal/models/follow.go)_
- [x] Follow service _(backend/internal/service/follow.go)_
- [x] FollowButton component _(frontend/src/components/follow/)_
- [x] Followers/Following lists _(frontend/src/components/follow/)_

#### US-7.5: Admin Panel
**As a** Platform Administrator
**I want to** manage users (search, change roles, enable/disable)
**So that** I can maintain the platform

**Acceptance Criteria**:
- [x] AdminService _(backend/internal/service/admin.go)_
- [x] CognitoClient for user management _(backend/internal/service/cognito_client.go)_
- [x] Admin handlers _(backend/internal/handlers/admin.go)_
- [x] Admin UI _(frontend/src/routes/admin/users.tsx)_

---

## Epic 8: Access Control & Bug Fixes 🔄 **IN PROGRESS**

**Goal**: Fix critical access control bugs and enforce visibility at service layer.

**Status**: In Progress (Tasks 4.1, 4.2 remaining)

**Spec**: [.spec-workflow/specs/access-control-bug-fixes/](.spec-workflow/specs/access-control-bug-fixes/)

### User Stories

#### US-8.1: Track Visibility Enforcement
**As a** Subscriber
**I want** to only see my own tracks and public tracks
**So that** other users' private content is protected

**Acceptance Criteria**:
- [x] GetTrack enforces visibility check _(service layer)_
- [x] ListTracks filters by visibility _(own tracks + public only)_
- [x] 403 returned for unauthorized private track access _(not 404)_
- [ ] End-to-end verification complete

#### US-8.2: Guest Route Protection
**As a** Platform Administrator
**I want** guest users to only access the dashboard
**So that** unauthenticated users can't see protected content

**Acceptance Criteria**:
- [x] Permission denied page created _(frontend/src/routes/permission-denied.tsx)_
- [x] Route guard in root layout _(frontend/src/routes/__root.tsx)_
- [ ] End-to-end verification complete

#### US-8.3: Admin User Modal Fix
**As an** Administrator
**I want** the user detail modal to work without errors
**So that** I can manage users effectively

**Acceptance Criteria**:
- [x] Modal handles null/undefined stats gracefully _(defensive rendering)_
- [x] Tests for edge cases _(frontend/src/components/admin/__tests__/)_

---

## Epic 9: LocalStack Development Environment ✅ **COMPLETE** (2026-01-27)

**Goal**: Enable local development with emulated AWS services for realistic integration testing.

**Status**: Complete

**Spec**: [.spec-workflow/specs/localstack-dev-environment/](.spec-workflow/specs/localstack-dev-environment/)

### User Stories

#### US-9.1: LocalStack Configuration
**As a** Developer
**I want** LocalStack to emulate DynamoDB, S3, and Cognito
**So that** I can test against real AWS APIs locally

**Acceptance Criteria**:
- [x] Docker Compose with LocalStack _(docker/docker-compose.yml)_
- [x] Cognito init script _(docker/localstack-init/init-cognito.sh)_
- [x] Wait-for-healthy script _(scripts/wait-for-localstack.sh)_
- [x] Test users created automatically _(admin, subscriber, artist)_

#### US-9.2: Integration Test Framework
**As a** Developer
**I want** utilities for writing integration tests against LocalStack
**So that** I can test with real DynamoDB operations

**Acceptance Criteria**:
- [x] SetupLocalStack helper _(backend/internal/testutil/localstack.go)_
- [x] Test fixtures _(backend/internal/testutil/fixtures.go)_
- [x] Cleanup utilities _(backend/internal/testutil/cleanup.go)_
- [x] Sample integration test _(backend/internal/service/track_integration_test.go)_

#### US-9.3: Frontend Local Mode
**As a** Developer
**I want** the frontend to work with LocalStack Cognito
**So that** I can test the full stack locally

**Acceptance Criteria**:
- [x] .env.local.example template _(frontend/.env.local.example)_
- [x] LocalStack mode detection _(frontend/src/lib/config.ts)_
- [x] dev:local npm script _(frontend/package.json)_

#### US-9.4: One-Command Setup
**As a** Developer
**I want** a single command to start everything
**So that** I can quickly start developing

**Acceptance Criteria**:
- [x] Makefile with local target _(Makefile)_
- [x] Shell script alternative _(scripts/local-dev.sh)_
- [x] Comprehensive documentation _(LOCAL_DEV.md)_
