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
