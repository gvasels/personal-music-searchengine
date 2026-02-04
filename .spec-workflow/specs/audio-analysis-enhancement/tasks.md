# Tasks - Audio Analysis Enhancement

## Feature: Audio Analysis Enhancement
**Status**: Not Started
**Estimated Effort**: 3-4 weeks

---

## Group 14: Backend Analysis Infrastructure

### Task 14.1: Update Track Model for Analysis Fields
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/track.go`
- `backend/internal/models/track_test.go`

**Changes**:
| Field | Type | Description |
|-------|------|-------------|
| `BPM` | `*int` | Beats per minute (nullable) |
| `MusicalKey` | `string` | Key signature (C, D, etc.) |
| `KeyMode` | `string` | major or minor |
| `KeyCamelot` | `string` | Camelot notation |
| `AnalysisStatus` | `string` | PENDING, ANALYZING, COMPLETED, FAILED |
| `AnalyzedAt` | `*time.Time` | Analysis completion timestamp |

**Acceptance Criteria**:
- [ ] Track model has new analysis fields
- [ ] TrackResponse includes analysis data
- [ ] TrackFilter supports BPM range and key
- [ ] Nullable fields handled correctly

### Task 14.2: Similar Artists Cache Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/similar.go`
- `backend/internal/models/similar_test.go`

**Acceptance Criteria**:
- [ ] Similar artists model defined
- [ ] DynamoDB item with proper PK/SK
- [ ] TTL field for cache expiration
- [ ] Unit tests for model

### Task 14.3: Analysis Repository Methods
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/repository.go`
- `backend/internal/repository/repository_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `UpdateTrackAnalysis(ctx, trackID, userID, analysis)` | Update BPM, key fields |
| `GetTracksWithoutAnalysis(ctx, userID, limit)` | Get tracks pending analysis |
| `GetSimilarArtistsCache(ctx, artistName)` | Get cached similar artists |
| `SetSimilarArtistsCache(ctx, cache)` | Cache similar artists |

**Acceptance Criteria**:
- [ ] Repository methods for analysis updates
- [ ] Efficient query for unanalyzed tracks
- [ ] Cache methods with TTL support
- [ ] Integration tests

---

## Group 15: BPM Detection Lambda

### Task 15.1: BPM Analyzer Lambda (Node.js)
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/bpm-analyzer/main.js`
- `backend/cmd/processor/bpm-analyzer/package.json`
- `backend/cmd/processor/bpm-analyzer/Dockerfile`

**Dependencies**:
```json
{
  "dependencies": {
    "realtime-bpm-analyzer": "^3.0.0",
    "@aws-sdk/client-s3": "^3.0.0",
    "music-metadata": "^7.0.0"
  }
}
```

**Acceptance Criteria**:
- [ ] Node.js Lambda with realtime-bpm-analyzer
- [ ] Downloads audio from S3
- [ ] Decodes MP3, FLAC, WAV formats
- [ ] Returns BPM with confidence
- [ ] Handles undetectable BPM gracefully

### Task 15.2: BPM Lambda Infrastructure
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/lambda-bpm-analyzer.tf`

**Acceptance Criteria**:
- [ ] Lambda function with Node.js 20 runtime
- [ ] Container image in ECR
- [ ] Memory: 512MB, Timeout: 60s
- [ ] S3 read permissions

---

## Group 16: Key Detection Lambda

### Task 16.1: Key Detector Lambda (Node.js)
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/key-detector/main.js`
- `backend/cmd/processor/key-detector/package.json`
- `backend/cmd/processor/key-detector/Dockerfile`

**Dependencies**:
```json
{
  "dependencies": {
    "essentia.js": "^0.1.3",
    "@aws-sdk/client-s3": "^3.0.0",
    "music-metadata": "^7.0.0"
  }
}
```

**Acceptance Criteria**:
- [ ] Node.js Lambda with essentia.js WASM
- [ ] Key detection with confidence score
- [ ] Camelot notation conversion
- [ ] Memory optimized for WASM (1024MB)

### Task 16.2: Key Lambda Infrastructure
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/lambda-key-detector.tf`

**Acceptance Criteria**:
- [ ] Lambda function for key detection
- [ ] Container image with essentia.js
- [ ] Memory: 1024MB, Timeout: 60s

---

## Group 17: Analysis Orchestration

### Task 17.1: Analysis Orchestrator Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/analysis-orchestrator/main.go`

**Acceptance Criteria**:
- [ ] Triggers analysis Step Functions
- [ ] Handles tracks uploaded before feature
- [ ] Idempotent (safe to call multiple times)

### Task 17.2: Analysis Step Functions
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/step-functions-analysis.tf`

**Acceptance Criteria**:
- [ ] Parallel BPM and key detection
- [ ] Update track with results
- [ ] Reindex for search
- [ ] Error handling per branch

### Task 17.3: Update Track Analysis Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/update-analysis/main.go`

**Acceptance Criteria**:
- [ ] Updates track with analysis results
- [ ] Sets analysisStatus to COMPLETED
- [ ] Handles partial results (BPM only, key only)

---

## Group 18: Similar Artists Service

### Task 18.1: Last.fm Client
**Status**: [ ] Pending
**Files**:
- `backend/internal/lastfm/client.go`
- `backend/internal/lastfm/client_test.go`
- `backend/internal/lastfm/CLAUDE.md`

**Acceptance Criteria**:
- [ ] HTTP client for Last.fm API
- [ ] API key from environment
- [ ] Parse similar artists response
- [ ] Handle rate limits (429)
- [ ] Handle not found (empty response)

### Task 18.2: Similar Artists Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/similar.go`
- `backend/internal/service/similar_test.go`

**Acceptance Criteria**:
- [ ] Cache-first lookup
- [ ] Fetch from Last.fm on cache miss
- [ ] Check library for each similar artist
- [ ] 30-day cache TTL

### Task 18.3: Similar Artists Handler
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/artists.go` (update)
- `backend/internal/handlers/artists_test.go` (update)

**Acceptance Criteria**:
- [ ] Returns similar artists
- [ ] Indicates which are in library
- [ ] Returns cached or fresh data
- [ ] Handles unknown artists

### Task 18.4: Last.fm API Key Secret
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/secrets.tf`

**Acceptance Criteria**:
- [ ] Secrets Manager secret for API key
- [ ] Lambda permission to read secret
- [ ] Environment variable configuration

---

## Group 19: API Updates

### Task 19.1: Track Filter Extensions
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/tracks.go`
- `backend/internal/service/track.go`
- `backend/internal/repository/repository.go`

**Acceptance Criteria**:
- [ ] Filter by BPM range
- [ ] Filter by key and mode
- [ ] Combined filters work

### Task 19.2: Track Analysis Endpoint
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/tracks.go`

**Acceptance Criteria**:
- [ ] Trigger analysis for existing tracks
- [ ] Batch processing endpoint
- [ ] Progress tracking

### Task 19.3: Search Index Update
**Status**: [ ] Pending
**Files**:
- `backend/internal/search/types.go`
- `backend/cmd/processor/indexer/main.go`

**Acceptance Criteria**:
- [ ] BPM indexed as integer
- [ ] Key indexed as keyword
- [ ] Range queries work for BPM

---

## Group 20: Frontend Updates

### Task 20.1: Track Type Updates
**Status**: [ ] Pending
**Files**:
- `frontend/src/types/index.ts`

**Acceptance Criteria**:
- [ ] Track type includes analysis fields
- [ ] Optional fields properly typed

### Task 20.2: Track Detail Page Updates
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/tracks/$trackId.tsx`
- `frontend/src/routes/tracks/$trackId.test.tsx`

**Acceptance Criteria**:
- [ ] BPM displayed when available
- [ ] Key displayed when available
- [ ] Analysis status shown
- [ ] Loading state during analysis

### Task 20.3: Artist Similar Artists Section
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/artists/$artistName.tsx`
- `frontend/src/hooks/useArtists.ts`
- `frontend/src/lib/api/artists.ts`

**Acceptance Criteria**:
- [ ] Fetch similar artists on artist page
- [ ] Display as clickable badges
- [ ] Indicate which are in library
- [ ] Link to artist page or search

### Task 20.4: Track List BPM/Key Columns
**Status**: [ ] Pending
**Files**:
- `frontend/src/components/library/TrackList.tsx`
- `frontend/src/components/library/TrackRow.tsx`

**Acceptance Criteria**:
- [ ] Optional BPM column in track list
- [ ] Optional Key column in track list
- [ ] Column visibility preference saved

---

## Group 21: Integration & Testing

### Task 21.1: Upload Pipeline Integration
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/step-functions.tf`

**Acceptance Criteria**:
- [ ] Analysis triggered automatically on upload
- [ ] Doesn't block upload completion
- [ ] Handles analysis failures gracefully

### Task 21.2: End-to-End Tests
**Status**: [ ] Pending
**Files**:
- `e2e/analysis.spec.ts`

**Acceptance Criteria**:
- [ ] E2E tests for full analysis flow
- [ ] E2E tests for similar artists
- [ ] Tests for filter functionality

### Task 21.3: Backfill Script
**Status**: [ ] Pending
**Files**:
- `scripts/backfill-analysis.sh`

**Acceptance Criteria**:
- [ ] Script to trigger batch analysis
- [ ] Progress reporting
- [ ] Safe to run multiple times

---

## Summary

| Group | Tasks | Purpose |
|-------|-------|---------|
| Group 14 | 3 | Backend models & repository |
| Group 15 | 2 | BPM detection Lambda |
| Group 16 | 2 | Key detection Lambda |
| Group 17 | 3 | Analysis orchestration |
| Group 18 | 4 | Similar artists service |
| Group 19 | 3 | API updates |
| Group 20 | 4 | Frontend updates |
| Group 21 | 3 | Integration & testing |
| **Total** | **24** | |

---

## Dependencies

### Task Dependencies
```
Task 14.1 → Task 14.3 → Task 17.3
Task 14.2 → Task 18.2
Task 15.1 → Task 15.2 → Task 17.2
Task 16.1 → Task 16.2 → Task 17.2
Task 17.1 → Task 17.2 → Task 21.1
Task 18.1 → Task 18.2 → Task 18.3
Task 20.1 → Tasks 20.2, 20.3, 20.4
```

### External Dependencies
- `realtime-bpm-analyzer` npm package
- `essentia.js` npm package
- Last.fm API key

---

## PR Checklist

After completing all tasks:
- [ ] All tests pass (`go test ./...`, `npm test`)
- [ ] Node.js Lambdas build (`npm run build`)
- [ ] Infrastructure validates (`tofu plan`)
- [ ] CLAUDE.md files updated
- [ ] CHANGELOG.md updated
- [ ] BPM detection works end-to-end
- [ ] Key detection works end-to-end
- [ ] Similar artists works end-to-end
- [ ] Create PR to main
