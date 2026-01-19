# Tasks - Search & Streaming (Epic 3)

## Epic: Search & Streaming
**Status**: Not Started
**Wave**: 3

---

## Group 8: Epic 2 Security Hardening (Prerequisite)

> **Must be completed before Epic 3 implementation begins**

### Task 8.1: Fix Memory Exhaustion in S3 Downloads
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/metadata/main.go`
- `backend/cmd/processor/coverart/main.go`

**Issue**: `io.ReadAll` loads entire audio file into Lambda memory. Will OOM on files >256MB.

**Fix**:
- Add file size validation via S3 HeadObject before download
- Reject files larger than 100MB
- Return clear error message for oversized files

**Functions**:
| Function | Description |
|----------|-------------|
| `validateFileSize(ctx, bucket, key)` | Check file size via HeadObject, return error if >100MB |

**Tests**:
| Test | Description |
|------|-------------|
| `TestValidateFileSize_Under100MB` | File under limit passes validation |
| `TestValidateFileSize_Over100MB` | File over limit returns error |
| `TestHandleRequest_RejectsLargeFile` | Handler returns error for large files |

**Acceptance Criteria**:
- [ ] Files >100MB rejected before download
- [ ] Clear error message returned to Step Functions
- [ ] HeadObject used for size check (no download needed)

### Task 8.2: Add UUID Validation for Input Parameters
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/track/main.go`
- `backend/internal/validation/uuid.go` (new)
- `backend/internal/validation/uuid_test.go` (new)
- `backend/internal/validation/CLAUDE.md` (new)

**Issue**: No UUID validation for userID/uploadID before DynamoDB/S3 operations. Could allow injection attacks.

**Fix**:
- Create validation package with UUID validator
- Validate userID and uploadID format before any operations
- Return 400-equivalent error for invalid UUIDs

**Functions**:
| Function | Description |
|----------|-------------|
| `IsValidUUID(s string)` | Returns true if string is valid UUID v4 format |
| `ValidateUUID(s, fieldName string)` | Returns error with field name if invalid |

**Tests**:
| Test | Description |
|------|-------------|
| `TestIsValidUUID_Valid` | Valid UUID returns true |
| `TestIsValidUUID_Invalid` | Invalid string returns false |
| `TestIsValidUUID_Empty` | Empty string returns false |
| `TestValidateUUID_Error` | Returns descriptive error |

**Acceptance Criteria**:
- [ ] Validation package created
- [ ] track/main.go validates userID and uploadID
- [ ] Invalid UUIDs rejected with clear error
- [ ] Prevents malformed input from reaching DynamoDB/S3

### Task 8.3: Add Context Timeout Configuration
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/metadata/main.go`
- `backend/cmd/processor/coverart/main.go`
- `backend/cmd/processor/track/main.go`
- `backend/cmd/processor/mover/main.go`
- `backend/cmd/processor/indexer/main.go`
- `backend/cmd/processor/status/main.go`

**Issue**: All processor Lambdas lack explicit context timeouts. Could cause hung executions.

**Fix**:
- Add context timeout slightly less than Lambda timeout
- Use `context.WithTimeout` in handler
- Ensure all operations respect context cancellation

**Functions**:
| Function | Description |
|----------|-------------|
| `withTimeout(ctx, seconds)` | Wraps context with timeout, returns cancel func |

**Tests**:
| Test | Description |
|------|-------------|
| `TestWithTimeout_CancelsAfterDuration` | Context cancelled after timeout |
| `TestHandleRequest_RespectsTimeout` | Long operations cancelled |

**Acceptance Criteria**:
- [ ] All processor Lambdas have explicit timeouts
- [ ] Timeout is 5 seconds less than Lambda config
- [ ] Context passed through to all operations
- [ ] Graceful handling of timeout errors

---

## Group 9: Nixiesearch Integration (Wave 3)

### Task 9.1: Nixiesearch Client Package
**Status**: [ ] Pending
**Files**:
- `backend/internal/search/client.go`
- `backend/internal/search/types.go`
- `backend/internal/search/client_test.go`
- `backend/internal/search/CLAUDE.md`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewClient(endpoint)` | Creates new Nixiesearch client with endpoint URL |
| `Search(ctx, query, filters, opts)` | Executes search query and returns results |
| `Index(ctx, doc)` | Indexes or updates a single document |
| `Delete(ctx, docId)` | Removes document from index |
| `BulkIndex(ctx, docs)` | Indexes multiple documents in batch |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSearch_SimpleQuery` | Basic text search returns ranked results |
| `TestSearch_WithFilters` | Search with artist/album filters works |
| `TestSearch_Pagination` | Cursor-based pagination returns correct pages |
| `TestIndex_NewDocument` | New document is indexed and searchable |
| `TestIndex_UpdateDocument` | Updated document reflects changes in search |
| `TestDelete_RemovesFromIndex` | Deleted document no longer in results |

**Acceptance Criteria**:
- [ ] Client connects to Nixiesearch endpoint
- [ ] Search supports full-text queries
- [ ] Filters apply correctly (artist, album, genre, year)
- [ ] Pagination uses cursor-based approach
- [ ] Index operations work for single and bulk

### Task 9.2: Search Handler
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/search.go`
- `backend/internal/handlers/search_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `(h *Handlers) SearchSimple(c)` | GET /search?q= handler |
| `(h *Handlers) SearchAdvanced(c)` | POST /search handler with filters |
| `parseSearchFilters(c)` | Extracts and validates filter parameters |
| `enrichSearchResults(ctx, results)` | Adds DynamoDB metadata to search results |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSearchSimple_Success` | Simple query returns results |
| `TestSearchSimple_EmptyQuery` | Empty query returns 400 |
| `TestSearchAdvanced_WithFilters` | Filters applied correctly |
| `TestSearchAdvanced_Pagination` | Cursor pagination works |
| `TestEnrichResults_AddsCoverArt` | Results include cover art URLs |

**Acceptance Criteria**:
- [ ] GET /search?q= works for simple searches
- [ ] POST /search works for advanced searches
- [ ] Results enriched with DynamoDB metadata
- [ ] Cover art URLs generated for results
- [ ] Proper error handling for invalid queries

### Task 9.3: Search Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/search.go`
- `backend/internal/service/search_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewSearchService(client, repo)` | Creates search service with dependencies |
| `Search(ctx, userID, query, filters, opts)` | Executes search scoped to user |
| `IndexTrack(ctx, track)` | Indexes track in search engine |
| `RemoveTrack(ctx, trackID)` | Removes track from search index |
| `RebuildIndex(ctx, userID)` | Full index rebuild for user |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSearch_ScopedToUser` | Search only returns user's tracks |
| `TestIndexTrack_Success` | Track indexed after creation |
| `TestRemoveTrack_Success` | Track removed after deletion |
| `TestRebuildIndex_AllTracks` | Rebuild indexes all user tracks |

**Acceptance Criteria**:
- [ ] Search results scoped to authenticated user
- [ ] Index operations integrate with search client
- [ ] Rebuild loads all tracks from DynamoDB
- [ ] Error handling for search service failures

### Task 9.4: Search Indexer Lambda (Update)
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/indexer/main.go` (update from stub)

**Functions**:
| Function | Description |
|----------|-------------|
| `handleRequest(ctx, event)` | Index track after upload completion |
| `indexTrack(ctx, client, track)` | Build search document and index |
| `buildSearchDocument(track)` | Convert Track to search document |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_IndexesTrack` | Track indexed after event |
| `TestBuildSearchDocument_AllFields` | All searchable fields included |

**Acceptance Criteria**:
- [ ] Replace stub with real indexing
- [ ] Uses Nixiesearch client
- [ ] Handles index failures gracefully
- [ ] Returns indexed=true on success

### Task 9.5: Scheduled Index Rebuild Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/indexer/rebuild/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point |
| `handleRequest(ctx, event)` | Triggered by EventBridge schedule |
| `rebuildAllUsers(ctx)` | Iterate all users and rebuild indexes |
| `rebuildUserIndex(ctx, userID)` | Full rebuild for single user |

**Tests**:
| Test | Description |
|------|-------------|
| `TestRebuildAllUsers_Success` | All users processed |
| `TestRebuildUserIndex_AllTracks` | All user tracks indexed |

**Acceptance Criteria**:
- [ ] Triggered by daily EventBridge schedule
- [ ] Processes all users in database
- [ ] Full rebuild of each user's index
- [ ] Handles large libraries (pagination)

---

## Group 10: HLS Transcoding Pipeline (Wave 3)

### Task 10.1: Transcode Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/transcode.go`
- `backend/internal/service/transcode_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewTranscodeService(mc, bucket)` | Creates transcode service |
| `StartTranscode(ctx, track)` | Creates MediaConvert job for track |
| `GetTranscodeStatus(ctx, jobID)` | Gets MediaConvert job status |
| `buildJobSettings(track)` | Builds MediaConvert job settings |

**Tests**:
| Test | Description |
|------|-------------|
| `TestStartTranscode_CreatesJob` | MediaConvert job created |
| `TestBuildJobSettings_ThreeQualities` | Settings include 96k, 192k, 320k |
| `TestBuildJobSettings_CorrectPaths` | Output paths use trackId |

**Acceptance Criteria**:
- [ ] Creates MediaConvert job with correct settings
- [ ] Three quality levels (96k, 192k, 320k AAC)
- [ ] Output to correct S3 paths
- [ ] Returns job ID for tracking

### Task 10.2: Transcode Start Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/transcode/start/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point |
| `handleRequest(ctx, event)` | Triggered after file move |
| `startTranscodeJob(ctx, track)` | Calls transcode service |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_StartsJob` | MediaConvert job started |
| `TestHandleRequest_UpdatesStatus` | Track hlsStatus set to PROCESSING |

**Acceptance Criteria**:
- [ ] Integrated into Step Functions after file mover
- [ ] Creates MediaConvert job
- [ ] Updates track hlsStatus to PROCESSING
- [ ] Returns job ID in response

### Task 10.3: Transcode Complete Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/transcode/complete/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point |
| `handleRequest(ctx, event)` | Triggered by EventBridge |
| `processSuccess(ctx, event)` | Update track with HLS info |
| `processFailure(ctx, event)` | Update track with error |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_Success` | Track updated with HLS info |
| `TestHandleRequest_Failure` | Track marked as FAILED |

**Acceptance Criteria**:
- [ ] Triggered by MediaConvert EventBridge events
- [ ] Updates track hlsStatus to READY on success
- [ ] Sets hlsPlaylistKey to master.m3u8 path
- [ ] Handles failures gracefully

### Task 10.4: MediaConvert IAM Role
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/iam-mediaconvert.tf`

**Acceptance Criteria**:
- [ ] IAM role for MediaConvert jobs
- [ ] S3 read access to media bucket
- [ ] S3 write access to hls prefix
- [ ] CloudWatch Logs access

### Task 10.5: MediaConvert Queue
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/mediaconvert.tf`

**Acceptance Criteria**:
- [ ] On-demand MediaConvert queue
- [ ] EventBridge rule for job completion
- [ ] Lambda trigger for completion events

---

## Group 11: CloudFront Streaming (Wave 3)

### Task 11.1: CloudFront Signer Package
**Status**: [ ] Pending
**Files**:
- `backend/internal/cloudfront/signer.go`
- `backend/internal/cloudfront/signer_test.go`
- `backend/internal/cloudfront/CLAUDE.md`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewSigner(keyPairID, privateKey, domain)` | Creates CloudFront signer |
| `SignURL(resource, expiry)` | Generates signed URL for resource |
| `SignStreamURL(trackID, userID)` | Signs HLS master playlist URL |
| `SignDownloadURL(s3Key)` | Signs direct download URL |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSignURL_ValidSignature` | Signed URL verifiable |
| `TestSignURL_Expiry` | URL includes correct expiration |
| `TestSignStreamURL_HLSPath` | URL points to HLS playlist |
| `TestSignDownloadURL_MediaPath` | URL points to original file |

**Acceptance Criteria**:
- [ ] Signs URLs with RSA private key
- [ ] Supports custom policy with expiration
- [ ] Works for both HLS and direct downloads
- [ ] 24-hour default expiration

### Task 11.2: Stream Handler
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/stream.go`
- `backend/internal/handlers/stream_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `(h *Handlers) GetStreamURL(c)` | GET /stream/{trackId} handler |
| `(h *Handlers) GetDownloadURL(c)` | GET /download/{trackId} handler |
| `buildStreamResponse(track, hlsURL, fallbackURL)` | Build stream response |

**Tests**:
| Test | Description |
|------|-------------|
| `TestGetStreamURL_HLSReady` | Returns HLS URL when available |
| `TestGetStreamURL_Fallback` | Returns direct URL when HLS not ready |
| `TestGetStreamURL_NotFound` | Returns 404 for unknown track |
| `TestGetStreamURL_NotOwner` | Returns 403 for other user's track |
| `TestGetDownloadURL_Success` | Returns signed download URL |

**Acceptance Criteria**:
- [ ] Returns HLS URL when hlsStatus=READY
- [ ] Returns fallback URL when HLS not ready
- [ ] Validates track ownership
- [ ] 24-hour URL expiration

### Task 11.3: Stream Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/stream.go`
- `backend/internal/service/stream_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewStreamService(repo, signer)` | Creates stream service |
| `GetStreamURLs(ctx, userID, trackID)` | Gets HLS and fallback URLs |
| `GetDownloadURL(ctx, userID, trackID)` | Gets download URL |
| `TrackDownload(ctx, trackID)` | Increments download count |

**Tests**:
| Test | Description |
|------|-------------|
| `TestGetStreamURLs_BothURLs` | Returns HLS and fallback |
| `TestGetStreamURLs_FallbackOnly` | Returns only fallback when no HLS |
| `TestGetDownloadURL_OriginalFile` | URL points to original file |
| `TestTrackDownload_IncrementsCount` | Download count updated |

**Acceptance Criteria**:
- [ ] Generates both HLS and fallback URLs
- [ ] Validates user ownership before signing
- [ ] Tracks download activity
- [ ] Handles missing HLS gracefully

### Task 11.4: CloudFront Distribution
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/cloudfront.tf`

**Acceptance Criteria**:
- [ ] CloudFront distribution for media bucket
- [ ] Origin Access Control (OAC) for S3
- [ ] Behaviors for /hls/*, /media/*, /covers/*
- [ ] Signed URL requirement for hls and media
- [ ] CORS headers in response
- [ ] Cache policy optimized for streaming

### Task 11.5: CloudFront Key Pair
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/cloudfront-keys.tf`

**Acceptance Criteria**:
- [ ] CloudFront public key resource
- [ ] CloudFront key group
- [ ] Private key in Secrets Manager
- [ ] Lambda permission to read secret

### Task 11.6: S3 Bucket Policy Update
**Status**: [ ] Pending
**Files**:
- `infrastructure/shared/s3.tf` (update)

**Acceptance Criteria**:
- [ ] Bucket policy allows CloudFront OAC
- [ ] Direct S3 access blocked
- [ ] Separate policy for hls and media prefixes

---

## Group 12: VPC and EFS for Nixiesearch (Wave 3)

### Task 12.1: VPC Infrastructure
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/vpc.tf`

**Acceptance Criteria**:
- [ ] VPC with private subnets (2 AZs)
- [ ] NAT Gateway for Lambda internet access
- [ ] Security group for Lambda-EFS communication
- [ ] VPC endpoints for DynamoDB and S3

### Task 12.2: EFS File System
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/efs.tf`

**Acceptance Criteria**:
- [ ] EFS file system for search index
- [ ] Mount targets in each private subnet
- [ ] Access point for Lambda mount
- [ ] Encryption at rest enabled

### Task 12.3: Nixiesearch Lambda
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/lambda-nixiesearch.tf`
- `backend/cmd/nixiesearch/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point |
| `handleSearch(ctx, req)` | Process search request |
| `handleIndex(ctx, req)` | Process index request |
| `initNixiesearch()` | Initialize Nixiesearch engine |

**Acceptance Criteria**:
- [ ] Lambda with VPC config and EFS mount
- [ ] Nixiesearch binary bundled in layer
- [ ] Handles search and index operations
- [ ] Index persisted to EFS
- [ ] Memory: 1024MB, Timeout: 30s

### Task 12.4: Nixiesearch Docker Image
**Status**: [ ] Pending
**Files**:
- `docker/nixiesearch/Dockerfile`
- `docker/nixiesearch/index-schema.yaml`

**Acceptance Criteria**:
- [ ] Docker image with Nixiesearch binary
- [ ] Index schema for music tracks
- [ ] Lambda layer build script
- [ ] ECR repository for image

---

## Group 13: Step Functions Update (Wave 3)

### Task 13.1: Update Upload Workflow
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/step-functions.tf` (update)

**Acceptance Criteria**:
- [ ] Add TranscodeStart state after FileMover
- [ ] TranscodeStart is async (doesn't wait for completion)
- [ ] Existing states unchanged
- [ ] Error handling for transcode failures

### Task 13.2: EventBridge Rules
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/eventbridge.tf`

**Acceptance Criteria**:
- [ ] Rule for MediaConvert COMPLETE events
- [ ] Rule for MediaConvert ERROR events
- [ ] Daily schedule for index rebuild
- [ ] Lambda targets for each rule

---

## Summary

| Group | Tasks | Status |
|-------|-------|--------|
| Group 8: Epic 2 Security Hardening | 3 | Not Started |
| Group 9: Nixiesearch Integration | 5 | Not Started |
| Group 10: HLS Transcoding Pipeline | 5 | Not Started |
| Group 11: CloudFront Streaming | 6 | Not Started |
| Group 12: VPC and EFS | 4 | Not Started |
| Group 13: Step Functions Update | 2 | Not Started |
| **Total** | **25** | **0 Complete** |

---

## Test Plan Summary

### Unit Tests
| File | Tests |
|------|-------|
| `internal/search/client_test.go` | Search client operations |
| `internal/service/search_test.go` | Search service logic |
| `internal/service/transcode_test.go` | Transcode service logic |
| `internal/service/stream_test.go` | Stream service logic |
| `internal/cloudfront/signer_test.go` | URL signing |
| `internal/handlers/search_test.go` | Search handler HTTP |
| `internal/handlers/stream_test.go` | Stream handler HTTP |

### Integration Tests
| Test | Environment |
|------|-------------|
| Search flow | LocalStack + Mock Nixiesearch |
| Stream URL generation | LocalStack + CloudFront mock |

### E2E Tests
| Test | Description |
|------|-------------|
| Upload → Search | New track appears in search |
| Upload → Stream | HLS available after transcode |
| Search → Play | Search result → stream playback |

---

## Dependencies

### Between Groups
- Group 10 depends on Group 9 (indexer needs search client)
- Group 11 depends on Group 12 (Lambda needs VPC for EFS)
- Group 13 depends on Groups 10, 11 (workflow includes new steps)

### External
- Nixiesearch binary (download or build)
- MediaConvert service (AWS managed)
- CloudFront (AWS managed)

---

## PR Checklist

After completing all tasks in this epic:
- [ ] All tests pass (`go test ./...`)
- [ ] Code builds (`go build ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Infrastructure validates (`tofu plan`)
- [ ] CLAUDE.md files updated for new packages
- [ ] CHANGELOG.md updated
- [ ] Search works end-to-end
- [ ] HLS streaming works
- [ ] Create PR to main with all changes
