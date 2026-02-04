# Requirements - Backend API (Epic 2)

## Overview

Deploy the backend API as AWS Lambda with API Gateway, implement upload processing pipeline using Step Functions, and create metadata extraction capabilities. This epic builds on Foundation Backend (Epic 1) which implemented models, repository, service, and handler layers.

## Goals

1. **API Deployment**: Deploy Go API as Lambda function accessible via API Gateway
2. **Upload Processing**: Implement Step Functions workflow for processing uploaded audio files
3. **Metadata Extraction**: Extract audio metadata (title, artist, album, etc.) from uploaded files
4. **Infrastructure**: Create API Gateway with Cognito JWT authorization

## Functional Requirements

### FR-1: API Lambda Deployment
- API Lambda receives HTTP requests via API Gateway
- Echo framework handles routing (already implemented in handlers)
- Lambda adapter converts between API Gateway events and Echo requests
- Support both local development (HTTP server) and Lambda deployment

### FR-2: Upload Processing Workflow
The Step Functions workflow (already defined in infrastructure) requires these Lambda handlers:

| Step | Lambda | Function |
|------|--------|----------|
| 1 | `metadata-extractor` | Extract audio metadata from S3 file |
| 2 | `cover-art-processor` | Extract embedded cover art, resize, save to S3 |
| 3 | `track-creator` | Create Track record in DynamoDB |
| 4 | `file-mover` | Move file from uploads/ to media/{userId}/{trackId} |
| 5 | `search-indexer` | Index track for search (stub for now) |
| 6 | `upload-status-updater` | Update Upload status to COMPLETED or FAILED |

### FR-3: Metadata Extraction
- Extract from MP3, FLAC, WAV, AAC, OGG files using `dhowden/tag`
- Extract: title, artist, album, album artist, genre, year, track number, disc number
- Extract: duration, bitrate, sample rate, channels
- Detect and extract embedded cover art (JPEG/PNG)

### FR-4: API Gateway
- HTTP API (v2) with Cognito JWT authorizer
- Routes matching handlers (29 endpoints under /api/v1)
- CORS configured for frontend origin
- Lambda integration with API Lambda

## Non-Functional Requirements

### NFR-1: Performance
- API Lambda cold start < 3s
- Metadata extraction < 10s for typical file
- Step Functions total processing < 60s for typical upload

### NFR-2: Reliability
- Step Functions retry on transient failures
- Partial success continues processing (cover art, indexing failures don't block)
- Upload status always updated (success or failure)

### NFR-3: Security
- API Gateway uses Cognito JWT authorizer (not custom middleware)
- User ID extracted from JWT claims (sub)
- S3 presigned URLs for all file access

### NFR-4: Cost
- Lambda ARM64 architecture for cost efficiency
- Minimal memory allocation (adjust based on profiling)
- S3 Intelligent-Tiering for media storage

## Constraints

1. **Existing Infrastructure**: Step Functions state machine already defined
2. **Existing Code**: Models, repository, service, handlers already implemented
3. **AWS Profile**: `gvasels-muza` for all AWS operations
4. **Region**: us-east-1 for all resources

## Dependencies

### From Epic 1 (Foundation Backend)
- Models: internal/models/*
- Repository: internal/repository/*
- Service: internal/service/*
- Handlers: internal/handlers/*

### From Infrastructure
- DynamoDB table (shared)
- S3 media bucket (shared)
- Cognito User Pool (shared)
- Lambda execution role (global)
- Step Functions state machine (backend)

## Success Criteria

1. [ ] API Lambda deploys and handles all 29 endpoints
2. [ ] Upload confirmation triggers Step Functions execution
3. [ ] Metadata extraction works for MP3, FLAC, WAV files
4. [ ] Cover art extracted and resized (if present)
5. [ ] Track created in DynamoDB after successful processing
6. [ ] Upload status updated to COMPLETED or FAILED
7. [ ] API Gateway authorizes requests with valid Cognito JWT
8. [ ] Local development works with HTTP server mode

## Out of Scope

- Nixiesearch integration (Epic 3)
- CloudFront signed URLs for streaming (Epic 3)
- Frontend (Epic 5)
