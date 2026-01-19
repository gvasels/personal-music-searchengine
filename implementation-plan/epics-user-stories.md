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
| E1 | Foundation Infrastructure | 1 | Critical | Not Started |
| E2 | Backend API | 2 | Critical | Not Started |
| E3 | Search & Streaming | 3 | High | Not Started |
| E4 | Tags & Playlists | 4 | High | Not Started |
| E5 | Frontend | 5 | High | Not Started |
| E6 | Distribution & Polish | 6 | Medium | Not Started |

---

## Epic 1: Foundation Infrastructure (Phase 1)

**Goal**: Establish core infrastructure and data contracts.

**Status**: Not Started

### User Stories

#### US-1.1: Data Models & Contracts
**As a** Developer
**I want to** define Go domain models and OpenAPI specification
**So that** the backend has a consistent contract for implementation

**Acceptance Criteria**:
- [ ] Go domain models defined (Track, Album, Artist, Playlist, User, Tag, Upload)
- [ ] OpenAPI 3.0 specification created
- [ ] Go project structure initialized
- [ ] Test fixtures and mocks created

#### US-1.2: Global Infrastructure
**As a** Developer
**I want to** set up shared infrastructure resources
**So that** all services can use consistent state management and container registries

**Acceptance Criteria**:
- [ ] S3 bucket for OpenTofu state with DynamoDB lock table
- [ ] ECR repositories created for Lambda images
- [ ] Base IAM roles for Lambda execution

#### US-1.3: Shared Services
**As a** Developer
**I want to** deploy authentication, database, and storage services
**So that** the backend API can use them

**Acceptance Criteria**:
- [ ] Cognito User Pool with app client and identity pool
- [ ] DynamoDB table with GSIs for single-table design
- [ ] S3 media bucket with Intelligent-Tiering

---

## Epic 2: Backend API (Phase 2)

**Goal**: Implement the core API Lambda and upload processing pipeline.

**Status**: Not Started

### User Stories

#### US-2.1: Core API Lambda
**As a** Music Enthusiast
**I want to** interact with my music library via API
**So that** I can manage my tracks, albums, and playlists

**Acceptance Criteria**:
- [ ] Echo-based API Lambda scaffold
- [ ] DynamoDB repository layer
- [ ] User profile handlers
- [ ] Track CRUD handlers
- [ ] Album/artist handlers

#### US-2.2: Upload & Processing Pipeline
**As a** Music Enthusiast
**I want to** upload music files and have metadata automatically extracted
**So that** my library is populated without manual data entry

**Acceptance Criteria**:
- [ ] Presigned URL generation for S3 upload
- [ ] Step Functions state machine for processing
- [ ] Metadata extraction Lambda
- [ ] Cover art extraction Lambda
- [ ] Track creation Lambda
- [ ] File mover Lambda
- [ ] Search indexer Lambda
- [ ] Status updater Lambda
- [ ] Upload service integration with Step Functions

#### US-2.3: API Gateway
**As a** Developer
**I want to** expose the API via HTTP API Gateway
**So that** the frontend can communicate with the backend

**Acceptance Criteria**:
- [ ] HTTP API with routes
- [ ] Cognito JWT authorizer configured
- [ ] Lambda integrations configured
- [ ] CORS configured for frontend

---

## Epic 3: Search & Streaming (Phase 3)

**Goal**: Enable full-text search and audio streaming.

**Status**: Not Started

### User Stories

#### US-3.1: Full-Text Search
**As a** Music Enthusiast
**I want to** search my music library by title, artist, album, or genre
**So that** I can quickly find songs I want to listen to

**Acceptance Criteria**:
- [ ] Nixiesearch Docker image built
- [ ] Nixiesearch Lambda function deployed
- [ ] Index schema defined for music metadata
- [ ] Indexer Lambda for scheduled re-indexing
- [ ] Search handlers implemented

#### US-3.2: Audio Streaming
**As a** Music Enthusiast
**I want to** stream and download my music files
**So that** I can listen to my library from anywhere

**Acceptance Criteria**:
- [ ] CloudFront distribution for media
- [ ] CloudFront signed URLs configured
- [ ] Stream/download handlers
- [ ] S3 search index bucket

---

## Epic 4: Tags & Playlists (Phase 4)

**Goal**: Enable custom organization with tags and playlists.

**Status**: Not Started

### User Stories

#### US-4.1: Custom Tags
**As a** Music Enthusiast
**I want to** create custom tags and apply them to tracks
**So that** I can organize my library my way

**Acceptance Criteria**:
- [ ] Tag repository
- [ ] Tag handlers
- [ ] Tag filtering in search

#### US-4.2: Playlists
**As a** Music Enthusiast
**I want to** create and manage playlists
**So that** I can curate collections of songs

**Acceptance Criteria**:
- [ ] Playlist repository
- [ ] Playlist handlers
- [ ] Playlist track ordering service

---

## Epic 5: Frontend (Phase 5)

**Goal**: Build the React web application.

**Status**: Not Started

### User Stories

#### US-5.1: Frontend Foundation
**As a** Developer
**I want to** set up the React project with proper tooling
**So that** UI development can proceed efficiently

**Acceptance Criteria**:
- [ ] React + Vite + TypeScript project initialized
- [ ] TanStack Query + Router configured
- [ ] Tailwind + DaisyUI with custom themes
- [ ] Amplify auth integration

#### US-5.2: Authentication & Layout
**As a** User
**I want to** sign in securely and navigate the application
**So that** I can access my music library

**Acceptance Criteria**:
- [ ] Login/signup pages
- [ ] Protected routes
- [ ] App shell with navigation
- [ ] Theme switcher (dark/light)

#### US-5.3: Library Views
**As a** Music Enthusiast
**I want to** browse my music by tracks, albums, or artists
**So that** I can explore my library visually

**Acceptance Criteria**:
- [ ] Track list with sorting/filtering
- [ ] Album grid view
- [ ] Artist list view
- [ ] Track detail/edit modal

#### US-5.4: Upload & Search UI
**As a** Music Enthusiast
**I want to** upload files and search my library from the UI
**So that** I can manage my library without using the API directly

**Acceptance Criteria**:
- [ ] Drag-drop file upload
- [ ] Upload progress tracking
- [ ] Search bar with autocomplete
- [ ] Search results page

#### US-5.5: Player & Queue
**As a** Music Enthusiast
**I want to** play music and manage a queue
**So that** I can listen to my library

**Acceptance Criteria**:
- [ ] Audio player component
- [ ] Play queue implementation
- [ ] Playlist management UI

---

## Epic 6: Distribution & Polish (Phase 6)

**Goal**: Deploy frontend, implement CI/CD, and finalize documentation.

**Status**: Not Started

### User Stories

#### US-6.1: Frontend Hosting
**As a** User
**I want to** access the application via a URL
**So that** I can use it from any device

**Acceptance Criteria**:
- [ ] S3 bucket for frontend
- [ ] CloudFront distribution
- [ ] Cache behaviors configured

#### US-6.2: Testing & CI/CD
**As a** Developer
**I want to** have automated tests and deployment pipelines
**So that** code quality is maintained

**Acceptance Criteria**:
- [ ] Go unit tests (80% coverage)
- [ ] API integration tests
- [ ] Frontend component tests
- [ ] GitHub Actions CI/CD

#### US-6.3: Documentation
**As a** Developer
**I want to** have comprehensive documentation
**So that** the project can be maintained and extended

**Acceptance Criteria**:
- [ ] CLAUDE.md updated with project specifics
- [ ] Deployment documentation
- [ ] API documentation

---

## Wave Execution Summary

| Wave | Groups | Focus | Status |
|------|--------|-------|--------|
| 0 | 1 | Local: contracts, models, project structure | Not Started |
| 1 | 2, 3 | Infrastructure: state, Cognito, DynamoDB, S3 | Not Started |
| 2 | 4, 5, 6 | Backend: API Lambda, Step Functions, API Gateway | Not Started |
| 3 | 7, 8, 9, 10 | Features: search, streaming, tags, playlists | Not Started |
| 4 | 11-15 | Frontend: full React application | Not Started |
| 5 | 16, 17, 18 | Polish: hosting, tests, documentation | Not Started |

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
