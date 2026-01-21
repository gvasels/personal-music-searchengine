# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

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

#### Epic 5: Frontend
- React 18 SPA with Vite, TypeScript, TanStack Router/Query
- TDD Implementation (63 unit tests, 81.56% coverage)
  - Store tests: themeStore (4 tests), playerStore (14 tests)
  - API client tests (11 tests)
  - Component tests: Layout (5), TrackList (6), PlayerBar (5), UploadDropzone (5), CreatePlaylistModal (7), TagInput (6)
- Authentication & Auth Guard
  - AWS Amplify/Cognito integration
  - Protected routes with login redirect
  - useAuth hook for auth state
- Layout Components
  - Header with SearchBar and theme toggle
  - Sidebar with navigation (mobile hamburger menu)
  - Layout app shell with responsive design
- Library Views
  - TrackList with click-to-play integration
  - formatDuration utility for time display
- Audio Playback
  - PlayerBar with Howler.js integration
  - Zustand player store (queue, volume, repeat, shuffle)
  - Play/pause, skip, seek, volume controls
- Upload & Search
  - UploadDropzone with react-dropzone (drag & drop)
  - File type validation (MP3, FLAC, WAV, OGG, M4A, AAC)
  - /upload route with progress tracking
- Playlist & Tag Management
  - CreatePlaylistModal with form validation
  - TagInput component for add/remove tags (lowercase normalization)
- Theming
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
