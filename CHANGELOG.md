# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

#### Access Control & Admin Features
- **Real-Time Role Checking**
  - `getAuthContextWithDBRole` handler method for DB role lookup
  - `RequireRoleWithDBCheck` middleware for admin routes
  - DB role takes precedence over JWT claims for immediate effect
  - Role changes take effect without requiring re-login
- **Admin User Management**
  - Admin routes: `/api/v1/admin/users` (search, details, role, status)
  - Role syncing between DynamoDB and Cognito groups
  - Self-modification prevention (admins cannot modify own role/status)
  - User search via Cognito with DynamoDB role enrichment
- **Track Visibility Enforcement**
  - Service-layer visibility checking (not just handlers)
  - 403 Forbidden for unauthorized access to private tracks
  - 404 Not Found only for truly non-existent resources
  - Admins (`hasGlobal=true`) can access any track
  - Visibility levels: `private`, `unlisted`, `public`
- **"Uploaded By" Column**
  - Added to `/tracks` page for admin users
  - Shows when: user is admin AND preference enabled in settings
  - `OwnerDisplayName` populated in track listings
  - Falls back to email, then "Unknown" if display name empty
- **Guest User Route Protection**
  - Public routes: `/`, `/login`, `/permission-denied`
  - All other routes redirect unauthenticated users to `/permission-denied`
  - `isPublicRoute()` helper in `__root.tsx`

### Fixed

#### Authentication & Session Management
- **JWT Groups Parsing** - API Gateway sends groups as `[admin GlobalReaders]` (bracket-separated), not JSON array
- **Data Leakage Prevention** - Clear React Query cache on sign-in/out to prevent data leakage between users
- **DynamoDB GetTrackByID** - Fixed pagination bug where `Limit: 1` only scanned 1 item instead of limiting returned results
- **Stream/Download Admin Access** - Updated handlers and services to support `hasGlobal` parameter for admin streaming
- **S3 CORS Configuration** - Added `https://music.vasels.com` to allowed origins for media bucket
- **Audio CORS Attribute** - Added `crossOrigin = 'anonymous'` to audio element for proper CORS handling

#### Track Listing & Data Management
- **Admin Track Listing Pagination** - Fixed `listAllTracks` DynamoDB Scan to properly paginate
  - DynamoDB `Limit` applies BEFORE filter, not after
  - Now scans in batches of 100 until enough tracks collected
  - Previously, admin would only see partial tracks when table had many non-track items
- **Admin Track Deletion** - Admins can now delete any track regardless of owner
  - `DeleteTrack` service accepts `hasGlobal` parameter
  - Uses `GetTrackByID` to find track regardless of owner when admin
  - Deletes using actual track owner's userID
  - Previously returned 404 when admin tried to delete other users' tracks
- **Clean Track Deletion** - Added comprehensive S3 cleanup
  - `DeleteByPrefix` method for batch deletions
  - Deletes audio file, cover art, AND all HLS transcoded segments
  - HLS files stored at `hls/{userId}/{trackId}/` now properly cleaned up

#### UI/Layout Fixes
- **Content Overlapping Player Bar** - Fixed layout so tables don't overlap fixed PlayerBar
  - Changed `min-h-screen` to `h-screen overflow-hidden` on root container
  - Added `min-h-0` to flex container to allow shrinking below content height
  - Main content area now scrolls independently with `overflow-auto`
  - Added `pb-28` (7rem) bottom padding to prevent content overlap

#### LocalStack Development Environment
- **Docker & LocalStack Configuration**
  - Added `cognito-idp` service to LocalStack for local authentication
  - Created `docker/localstack-init/init-cognito.sh` for Cognito setup
    - Creates user pool: `music-library-local-pool`
    - Creates app client (no secret for SPA)
    - Creates groups: `admin`, `artist`, `subscriber`
    - Creates test users with known credentials
  - Created `scripts/wait-for-localstack.sh` health check script
- **Integration Test Framework** (`backend/internal/testutil/`)
  - `localstack.go` - SetupLocalStack with TestContext and cleanup
  - `fixtures.go` - Test users, CreateTestTrack, CreateTestUser helpers
  - `cleanup.go` - CleanupUser, CleanupTrack, CleanupAll utilities
  - Sample integration test: `backend/internal/service/track_integration_test.go`
  - Tests use `//go:build integration` tag for separation
- **Frontend Local Mode**
  - `frontend/.env.local.example` - Environment template
  - `frontend/src/lib/config.ts` - LocalStack mode detection
  - `dev:local` npm script for local development
  - Test users: admin@local.test, subscriber@local.test, artist@local.test (password: LocalTest123!)
- **One-Command Setup**
  - Root `Makefile` with targets: `local`, `local-stop`, `test-integration`
  - `scripts/local-dev.sh` shell script alternative
- **Documentation**
  - `LOCAL_DEV.md` - Comprehensive local development guide
  - Updated `docker/CLAUDE.md` with Cognito documentation
  - Created `backend/internal/testutil/CLAUDE.md`

#### Global User Type Feature (Role-Based Access Control)
- **Backend Services**
  - User roles: `guest`, `subscriber`, `artist`, `admin` with Cognito Groups integration
  - Permission system: 12 granular permissions (browse, listen, upload, publish, etc.)
  - Artist profile service: CRUD operations, catalog linking, search
  - Follow service: follow/unfollow artists, followers/following lists
  - Authorization middleware with role extraction from JWT claims
  - Playlist visibility: `private`, `unlisted`, `public` levels
  - Public playlist discovery endpoint
- **Frontend Components**
  - `ArtistProfileCard` - Display artist profile with follow button
  - `EditArtistProfileModal` - Create/edit artist profile form
  - `FollowButton` - Follow/unfollow toggle button
  - `FollowersList` - Display artist followers
  - `FollowingList` - Display artists user is following
  - `VisibilitySelector` - Playlist visibility dropdown/radio group
  - `VisibilityBadge` - Display visibility status
  - `RoleSwitcher` - Admin dropdown for role simulation
  - `SimulationBanner` - Alert banner during role simulation
- **Frontend Hooks & Stores**
  - `useRoleSimulation` - Admin role simulation hook
  - `roleSimulationStore` - Zustand store with localStorage persistence
  - `useFeatureFlags` respects simulated role for `hasRole()` checks
- **Infrastructure**
  - Cognito user groups: `admin`, `artist`, `subscriber`
  - DynamoDB entities: `ARTIST_PROFILE`, `FOLLOW`
- **API Endpoints**
  - `POST/GET/PUT/DELETE /artists/entity` - Artist profile CRUD
  - `POST/DELETE /artists/entity/:id/follow` - Follow/unfollow
  - `GET /artists/entity/:id/followers` - List followers
  - `GET /users/me/following` - List following
  - `PUT /playlists/:id/visibility` - Update playlist visibility
  - `GET /playlists/public` - List public playlists

#### Stream D: DevOps & CI/CD Improvements
- **Backend Lint Job** (`.github/workflows/ci.yml`)
  - Added golangci-lint v1.61 with 5-minute timeout
  - Fixed unchecked json.Unmarshal errors in client_test.go
  - Removed unused getFileExtension function from mover/main.go
  - Fixed empty if branch in search_test.go
- **Security Scan Improvements**
  - Added fetch-depth: 0 for gitleaks full history scan
- **Coverage Threshold**
  - Backend coverage threshold at 19% (temporary, to be increased to 80%)

#### Stream B: Platform Features
- **Playlist Reorder Endpoint** (`PUT /playlists/:id/tracks/reorder`)
  - Reorder tracks within a playlist by specifying new positions
  - Validates track ownership and playlist membership
  - Atomic position updates with batch write
- **Audio Analysis Module** (`backend/internal/analysis/`)
  - BPM detection using multi-segment autocorrelation algorithm
  - Bass-emphasis filter targeting kick drum frequencies (~200Hz)
  - Octave error correction for double/half time detection
  - Genre-aware BPM preference (house 115-135, trance 135-150, hip-hop 85-95)
  - Input validation for command injection prevention
  - File extension whitelist for security
- **Analyzer Lambda** (`backend/cmd/processor/analyzer/`)
  - Step Functions integration for upload processing pipeline
  - Graceful degradation - analysis failures don't block uploads
  - 25-second timeout with context handling
- **Bedrock Gateway** (`infrastructure/backend/bedrock-gateway.tf`)
  - OpenAI-compatible API gateway for AWS Bedrock
  - API key authentication via Secrets Manager
  - CORS restricted to frontend CloudFront and localhost
- **Migration Service** (`backend/internal/service/migration.go`)
  - Idempotent artist entity migration
  - Batch processing with error handling

#### Stream A: Audio Features
- **Waveform Generation** (`backend/cmd/processor/waveform/`)
  - FFmpeg-based waveform peak extraction at 100 samples/second
  - Normalized peak values (0.0-1.0) for consistent visualization
  - Support for MP3, FLAC, WAV, AAC, OGG, M4A formats
  - Fallback synthetic waveform when FFmpeg processing fails
  - JSON serialization for S3 storage
- **Beat Grid Calculation** (`backend/cmd/processor/beatgrid/`)
  - Beat timestamp calculation from BPM (20-300 BPM range)
  - Downbeat markers (every 4th beat) for DJ features
  - Binary search for efficient beat-at-time lookup
  - Variable BPM flag support
- **Track Model Enhancements**
  - `WaveformURL` - S3 URL to waveform JSON data
  - `BeatGrid` - Array of beat timestamps in milliseconds
  - `AnalysisStatus` - PENDING, ANALYZING, COMPLETED, FAILED
  - `AnalyzedAt` - Analysis completion timestamp

#### Stream C: Search Enhancements (Semantic Search)
- EmbeddingService for Bedrock Titan text embeddings
  - `ComposeEmbedText` - Track metadata to embedding text composition
  - `GenerateTrackEmbedding` - 1024-dim embedding from track metadata
  - `GenerateQueryEmbedding` - Search query embedding generation
  - `BatchGenerateEmbeddings` - Batch processing with partial failure handling
  - Truncation to 8000 chars (Titan model limit)
- Camelot key compatibility utilities for DJ mixing
  - `IsKeyCompatible` - Check if two keys can be mixed harmonically
  - `GetCompatibleKeys` - Get all 4 compatible keys per Camelot key
  - `GetKeyTransition` - Describe transition type (Perfect Match, Smooth, Relative)
  - `GetBPMCompatibility` - BPM compatibility with half/double time support
- SimilarityService for finding similar and mixable tracks
  - `FindSimilarTracks` - Semantic + feature-based similarity search
  - `FindMixableTracks` - DJ-compatible tracks by BPM + key
  - `CosineSimilarity` - Vector similarity calculation
  - Three modes: semantic, features, combined (default 60/40 weight)
- Semantic search specifications
  - `.spec-workflow/specs/semantic-search/requirements.md`
  - `.spec-workflow/specs/semantic-search/design.md`
  - `.spec-workflow/specs/semantic-search/tasks.md`
- Unit tests for embedding service (20 tests)
- Unit tests for Camelot utilities

#### Epic 6: Distribution & Polish
- Frontend S3 bucket (`music-library-prod-frontend`) for SPA hosting
  - Versioning enabled for rollback capability
  - AES-256 server-side encryption
  - Public access blocked (OAC only)
- CloudFront distribution for frontend with:
  - Origin Access Control (OAC) for S3 security
  - SPA routing (403/404 → index.html)
  - Cache behaviors: no-cache for index.html, 1 year TTL for /assets/*
  - Gzip and Brotli compression
  - PriceClass_100 (US, Canada, Europe)
- Security headers response policy:
  - X-Content-Type-Options: nosniff
  - X-Frame-Options: DENY
  - Strict-Transport-Security: max-age=31536000; includeSubDomains
  - X-XSS-Protection: 1; mode=block
- GitHub Actions OIDC authentication for AWS deployments
  - IAM role with least-privilege permissions
  - S3, CloudFront, ECR, Lambda deploy permissions
- CI workflow (`.github/workflows/ci.yml`):
  - Go tests with 80% coverage threshold
  - Frontend tests (Vitest) with coverage
  - ESLint and TypeScript type checking
  - OpenTofu validation for all modules
  - Gitleaks secrets scanning
  - Checkov IaC security scanning
- Deploy workflow (`.github/workflows/deploy.yml`):
  - Automated infrastructure deployment on merge to main
  - Frontend build and S3 sync with cache headers
  - CloudFront cache invalidation
  - Backend Docker build and Lambda update
- Deployment documentation (`docs/deployment.md`)

#### Epic 1: Foundation Backend
- Go domain models (Track, Album, Artist, Playlist, User, Tag, Upload)
- DynamoDB single-table design with repository layer
- Service layer with business logic (TrackService, AlbumService, etc.)
- HTTP handlers using Echo framework
- OpenAPI specification for REST API

#### Epic 2: Backend API
- API Lambda entrypoint with dual mode (Lambda/HTTP)
- Audio metadata extraction package (MP3, FLAC, WAV, OGG support)
- Upload processing pipeline Lambdas:
  - Metadata extractor Lambda
  - Cover art processor Lambda
  - Track creator Lambda
  - File mover Lambda
  - Search indexer stub Lambda
  - Upload status updater Lambda
- API Gateway HTTP API with Cognito JWT authorizer
- Infrastructure for all Lambda functions (ARM64, provided.al2023)
- LocalStack docker-compose for local development
- Makefile for build and test commands

#### Epic 3: Search & Streaming
- Nixiesearch integration with Lambda invocation client
  - Full-text search across title, artist, album, filename
  - Weighted scoring (title 3x, artist 2x, album 1.5x)
  - Filter by artist, album, genre, year range
  - Pagination with cursor support
- Nixiesearch Lambda with S3-based index storage (pure serverless, no VPC)
  - Container image (ECR) with embedded search engine
  - Operations: search, index, delete, bulk_index
  - Cold start index loading from S3
- CloudFront distribution for media streaming
  - Signed URLs with RSA-SHA1 canned policy
  - Origin Access Control (OAC) for S3 security
  - CORS configuration for web playback
- HLS adaptive bitrate streaming via MediaConvert
  - 3 quality levels: 96kbps, 192kbps, 320kbps AAC
  - Transcode start/complete Lambdas
  - EventBridge rules for job status handling
- Step Functions workflow updates
  - StartTranscode step after file move
  - Async transcoding via EventBridge completion events
- EventBridge scheduled tasks
  - Daily search index rebuild (3 AM UTC)
- ECR repository for Nixiesearch container images

#### Epic 4: Tags & Playlists
- Tag filtering in search
  - Filter search results by multiple tags (AND logic)
  - Tags stored in DynamoDB, filtered post-search
  - Returns NotFoundError if tag doesn't exist
  - Tag deduplication in filter requests
- Tag name normalization (case-insensitive)
  - All tag names normalized to lowercase throughout
  - "Rock" and "ROCK" match the same tag "rock"
  - Normalization applied in all tag service methods
- Unit tests for tag service (24 tests)
  - Tag CRUD operations
  - Track-tag associations
  - Case normalization tests
- Unit tests for playlist service (16 tests)
  - Playlist CRUD operations
  - Track add/remove with position support
  - Cover art URL generation
- Unit tests for filterByTags (8 tests)
  - Empty tags, tag not found, AND logic
  - Deduplication, case normalization

#### Epic 5: Frontend (Complete)
- React 18 SPA with Vite, TypeScript, TanStack Router/Query
- **TDD Implementation (351 unit tests)**
  - Store tests: themeStore (4 tests), playerStore (14 tests)
  - API client tests: client (11), tracks (10), albums (5), artists (5), upload (4), search (4), playlists (7), tags (4)
  - Hook tests: useAuth (18), useTracks (13), useAlbums (11), useArtists (11), useUpload (6), useSearch (8), usePlaylists (6), useTags (6)
  - Component tests: Layout (5), TrackList (6), PlayerBar (5), UploadDropzone (5), CreatePlaylistModal (7), TagInput (6)
  - Route tests: index (22), login (19), tracks (10), trackDetail (17), albums (11), artists (11), upload (9), search (9), playlists (9), playlistDetail (11), tags (8), tagDetail (9)
- **Wave 1: Authentication & Pages**
  - AWS Amplify/Cognito integration
  - Protected routes with auth guard
  - useAuth hook for auth state
  - Login page with form validation
  - Home page with library statistics
- **Wave 2: Library Views**
  - API layer: tracks, albums, artists with query key factories
  - Hooks: useTracks, useAlbums, useArtists with TanStack Query
  - Routes: /tracks, /tracks/:trackId, /albums, /artists
  - TrackList with click-to-play, sorting, duration formatting
- **Wave 3: Audio Playback**
  - PlayerBar with Howler.js integration
  - Zustand player store (queue, volume, repeat, shuffle)
  - Play/pause, skip, seek, volume controls
- **Wave 4: Upload & Search**
  - API layer: upload (presigned URL, confirm, status), search (query, autocomplete)
  - Hooks: useUpload (file upload with progress), useSearch (debounced queries)
  - UploadDropzone with react-dropzone (drag & drop, file validation)
  - Search page with results and filter inputs (artist, album)
- **Wave 5: Playlists & Tags**
  - API layer: playlists (CRUD, add/remove tracks), tags (list, tracks by tag)
  - Hooks: usePlaylists, useTags with query key factories
  - Routes: /playlists, /playlists/:playlistId, /tags, /tags/:tagName
  - CreatePlaylistModal with form validation
  - TagInput component (add/remove with lowercase normalization)
  - Tag cloud with size based on track count
- **Layout Components**
  - Header with SearchBar and theme toggle
  - Sidebar with navigation (mobile hamburger menu)
  - Layout app shell with responsive design
- **Theming**
  - DaisyUI 5 with Tailwind CSS 4
  - Dark theme: #120612 base, #50c878 primary, #72001c secondary
  - Light theme: #fdfdf8 base with same accent colors
  - Theme persistence via Zustand store

### Changed
- CloudFront URL signing with expiration bounds validation (5 min to 7 days)
- Search service with query length validation (max 500 characters)
- Updated specs for pure serverless architecture (removed VPC/EFS references)

### Security
- Input validation for all Lambda processor payloads
- File size limits on uploads (500MB max)
- Filename sanitization to prevent path traversal
- Content type validation for audio files
- Resource limits on metadata extraction

#### Initial Setup
- Initial project setup
- SDLC workflow plugins (6 phases)
- Specialized agents for development tasks
- Reusable skills for code review and documentation
- Slash commands for common operations
- MCP server configuration
- Project documentation (CLAUDE.md, README.md)

---

## Template Information

This project was created from the Claude Code Starter Project template.

### Included Components

**Plugins (SDLC Workflow)**
- `spec-writer` - Requirements and technical design
- `test-writer` - TDD test creation
- `code-implementer` - Implementation
- `builder` - Build verification
- `security-checker` - Security audit
- `docs-generator` - Documentation generation

**Agents**
- `implementation-agent` - Feature implementation
- `test-engineer` - Test strategy and TDD
- `code-review` - Code quality analysis
- `security-auditor` - Security vulnerability analysis
- `documentation-generator` - Documentation generation

**Skills**
- `code-reviewer` - Code review capabilities
- `documentation-generator` - Documentation capabilities

**Commands**
- `/sdlc` - Full SDLC workflow
- `/update-claudemd` - Update CLAUDE.md from git
- `/code-review` - Code review
- `/test-file` - Generate tests for a file
