# Tasks - Backend API (Epic 2)

## Epic: Backend API
**Status**: Complete
**Wave**: 2

---

## Group 5: Lambda Entrypoints (Wave 2)

### Task 5.1: API Lambda Main
**Status**: [x] Complete
**Files**:
- `backend/cmd/api/main.go`
- `backend/cmd/api/config.go`
- `backend/cmd/api/validator.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Entry point - initializes deps, chooses Lambda or HTTP mode based on env |
| `loadConfig()` | Loads configuration from environment variables (table name, bucket, etc.) |
| `isLambda()` | Returns true if AWS_LAMBDA_FUNCTION_NAME env var is set |
| `NewValidator()` | Creates go-playground/validator instance for Echo binding |

**Tests**:
| Test | Description |
|------|-------------|
| `TestLoadConfig_ValidEnv` | Config loads all required env vars correctly |
| `TestLoadConfig_MissingEnv` | Config fails with clear error when vars missing |
| `TestIsLambda_True` | Detects Lambda environment correctly |
| `TestIsLambda_False` | Detects non-Lambda environment |
| `TestValidator_ValidRequest` | Validates valid request structs |
| `TestValidator_InvalidRequest` | Returns validation errors for invalid structs |

**Acceptance Criteria**:
- [x] main.go compiles and runs locally with `go run`
- [x] Lambda mode detected via environment variable
- [x] HTTP server mode starts on :8080 for local dev
- [x] All environment variables documented
- [x] Validator works with existing request types

### Task 5.2: Go Module Dependencies
**Status**: [x] Complete
**Files**:
- `backend/go.mod`
- `backend/go.sum`

**Acceptance Criteria**:
- [x] Add `github.com/dhowden/tag` for metadata extraction
- [x] Add `github.com/go-playground/validator/v10` for request validation
- [x] Verify `go mod tidy` succeeds
- [x] Verify `go build ./...` succeeds

---

## Group 6: Upload Processing Pipeline (Wave 2)

### Task 6.1: Metadata Extractor Package
**Status**: [x] Complete
**Files**:
- `backend/internal/metadata/extractor.go`
- `backend/internal/metadata/extractor_test.go`
- `backend/internal/metadata/CLAUDE.md`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewExtractor()` | Creates new metadata extractor instance |
| `Extract(reader)` | Extracts all metadata from audio file, returns UploadMetadata |
| `ExtractCoverArt(reader)` | Extracts cover art bytes and mime type from audio file |
| `detectFormat(reader)` | Detects audio format (MP3, FLAC, WAV, etc.) from file header |
| `parseTitle(meta, filename)` | Gets title from metadata or falls back to filename |
| `parseDuration(meta, fileSize)` | Gets duration from metadata or estimates from bitrate |

**Tests**:
| Test | Description |
|------|-------------|
| `TestExtract_MP3_WithMetadata` | Extracts title, artist, album from MP3 ID3 tags |
| `TestExtract_MP3_NoMetadata` | Returns filename as title when no ID3 tags |
| `TestExtract_FLAC_WithMetadata` | Extracts Vorbis comments from FLAC |
| `TestExtract_WAV_NoMetadata` | Handles WAV with no metadata gracefully |
| `TestExtractCoverArt_Present` | Extracts embedded JPEG cover art |
| `TestExtractCoverArt_Missing` | Returns nil when no cover art embedded |
| `TestDetectFormat_MP3` | Correctly identifies MP3 format |
| `TestDetectFormat_FLAC` | Correctly identifies FLAC format |

**Acceptance Criteria**:
- [x] Extract metadata from MP3 (ID3v1, ID3v2)
- [x] Extract metadata from FLAC (Vorbis comments)
- [x] Extract metadata from WAV (if present)
- [x] Handle missing metadata gracefully (use filename)
- [x] Extract embedded cover art (JPEG, PNG)
- [ ] Unit tests pass with test fixtures (tests to be added in later epic)

### Task 6.2: Metadata Extractor Lambda
**Status**: [x] Complete
**Files**:
- `backend/cmd/processor/metadata/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point with handler registration |
| `handleRequest(ctx, event)` | Downloads file from S3, extracts metadata, returns result |
| `downloadFromS3(ctx, bucket, key)` | Downloads file to memory buffer for processing |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_Success` | Extracts metadata from valid S3 file |
| `TestHandleRequest_S3Error` | Returns error when S3 download fails |
| `TestHandleRequest_InvalidFormat` | Returns error for unsupported file format |

**Acceptance Criteria**:
- [x] Lambda handler receives Step Functions event
- [x] Downloads file from S3 using presigned URL or direct access
- [x] Uses metadata extractor package
- [x] Returns UploadMetadata struct for Step Functions

### Task 6.3: Cover Art Processor Lambda
**Status**: [x] Complete
**Files**:
- `backend/cmd/processor/coverart/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point with handler registration |
| `handleRequest(ctx, event)` | Extracts cover art if present, uploads to S3 |
| `uploadCoverArt(ctx, bucket, key, data, contentType)` | Uploads cover art bytes to S3 |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_WithCoverArt` | Extracts and uploads cover art |
| `TestHandleRequest_NoCoverArt` | Returns empty result when no cover art |
| `TestUploadCoverArt_Success` | Uploads image to correct S3 path |

**Acceptance Criteria**:
- [x] Receives hasCoverArt flag from previous step
- [x] Skips processing if no cover art
- [x] Uploads cover art to `covers/{userId}/{trackId}.{ext}`
- [x] Returns coverArtKey in response

### Task 6.4: Track Creator Lambda
**Status**: [x] Complete
**Files**:
- `backend/cmd/processor/track/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point with handler registration |
| `handleRequest(ctx, event)` | Creates Track and Album records in DynamoDB |
| `createTrack(ctx, repo, event)` | Creates Track record with all metadata |
| `getOrCreateAlbum(ctx, repo, event)` | Gets existing album or creates new one |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_NewTrack` | Creates new track record |
| `TestHandleRequest_WithAlbum` | Creates track and links to album |
| `TestHandleRequest_ExistingAlbum` | Uses existing album, increments track count |

**Acceptance Criteria**:
- [x] Creates Track in DynamoDB with all metadata
- [x] Creates Album if album name present and not exists
- [x] Links Track to Album if album created
- [x] Returns trackId and albumId

### Task 6.5: File Mover Lambda
**Status**: [x] Complete
**Files**:
- `backend/cmd/processor/mover/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point with handler registration |
| `handleRequest(ctx, event)` | Moves file from uploads/ to media/ |
| `moveFile(ctx, s3Client, bucket, source, dest)` | S3 copy then delete |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_Success` | Moves file to correct location |
| `TestHandleRequest_CopyFail` | Returns error when copy fails |
| `TestMoveFile_UpdatesTrack` | Updates Track S3Key after move |

**Acceptance Criteria**:
- [x] Copies file from `uploads/{uploadId}/{filename}` to `media/{userId}/{trackId}.{ext}`
- [x] Deletes original after successful copy
- [x] Updates Track.S3Key in DynamoDB
- [x] Returns finalKey in response

### Task 6.6: Search Indexer Lambda (Stub)
**Status**: [x] Complete
**Files**:
- `backend/cmd/processor/indexer/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point with handler registration |
| `handleRequest(ctx, event)` | Returns stub response (not implemented) |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_Stub` | Returns indexed=false, reason=not_implemented |

**Acceptance Criteria**:
- [x] Returns success with indexed=false
- [x] Does not block workflow (non-critical step)
- [x] Ready for Epic 3 Nixiesearch implementation

### Task 6.7: Upload Status Updater Lambda
**Status**: [x] Complete
**Files**:
- `backend/cmd/processor/status/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `main()` | Lambda entry point with handler registration |
| `handleRequest(ctx, event)` | Updates Upload status in DynamoDB |
| `updateStatus(ctx, repo, uploadId, status, trackId, error)` | Sets status fields |

**Tests**:
| Test | Description |
|------|-------------|
| `TestHandleRequest_Completed` | Sets COMPLETED status with trackId |
| `TestHandleRequest_Failed` | Sets FAILED status with error message |

**Acceptance Criteria**:
- [x] Updates Upload.Status to COMPLETED or FAILED
- [x] Sets Upload.CompletedAt timestamp
- [x] Sets Upload.TrackID if successful
- [x] Sets Upload.ErrorMsg if failed

---

## Group 7: API Gateway Infrastructure (Wave 2)

### Task 7.1: API Gateway Configuration
**Status**: [x] Complete
**Files**:
- `infrastructure/backend/api-gateway.tf`

**Acceptance Criteria**:
- [x] HTTP API (v2) created
- [x] Cognito JWT authorizer configured
- [x] Default stage deployed
- [x] CORS configured (localhost:5173, production domain)
- [x] All routes mapped to API Lambda
- [x] Authorization required on all routes except OPTIONS

### Task 7.2: API Lambda Infrastructure
**Status**: [x] Complete
**Files**:
- `infrastructure/backend/lambda-api.tf`

**Acceptance Criteria**:
- [x] Lambda function with ARM64 architecture
- [x] Runtime: provided.al2023 (custom Go)
- [x] Memory: 256MB, Timeout: 30s
- [x] Environment variables: table name, bucket, Step Functions ARN
- [x] CloudWatch log group with 30-day retention
- [x] IAM permissions: DynamoDB, S3, Step Functions (via shared role)

### Task 7.3: Processor Lambda Infrastructure
**Status**: [x] Complete
**Files**:
- `infrastructure/backend/lambda-processors.tf`

**Acceptance Criteria**:
- [x] 6 Lambda functions created (metadata, coverart, track, mover, indexer, status)
- [x] All use ARM64 architecture
- [x] Runtime: provided.al2023
- [x] Memory: 256-512MB, Timeout: 30-60s
- [x] IAM permissions: S3 read/write, DynamoDB read/write (via shared role)
- [x] CloudWatch log groups with 30-day retention

### Task 7.4: Infrastructure Outputs
**Status**: [x] Complete
**Files**:
- `infrastructure/backend/main.tf` (update outputs)

**Acceptance Criteria**:
- [x] Output: api_gateway_url
- [x] Output: api_gateway_id
- [x] Output: cognito_authorizer_id
- [x] Update existing outputs if needed

---

## Group 8: Local Development Setup (Wave 2)

### Task 8.1: Docker Compose for LocalStack
**Status**: [x] Complete
**Files**:
- `docker/docker-compose.yml`
- `docker/localstack-init/init-aws.sh`
- `docker/CLAUDE.md`

**Acceptance Criteria**:
- [x] LocalStack container with DynamoDB, S3
- [x] Init script creates MusicLibrary table
- [x] Init script creates media bucket
- [x] CORS configured for local development
- [x] Data persistence via Docker volume

### Task 8.2: Local Development Scripts
**Status**: [x] Complete
**Files**:
- `backend/Makefile`
- `backend/.env.example`

**Acceptance Criteria**:
- [x] `make run-local` starts API in HTTP mode
- [x] `make test` runs all tests
- [x] `make build` builds Lambda binaries
- [x] `make localstack-up` starts docker-compose
- [x] Environment variable examples documented

---

## Summary

| Group | Tasks | Status |
|-------|-------|--------|
| Group 5: Lambda Entrypoints | 2 | Complete |
| Group 6: Upload Processing Pipeline | 7 | Complete |
| Group 7: API Gateway Infrastructure | 4 | Complete |
| Group 8: Local Development Setup | 2 | Complete |
| **Total** | **15** | **15 Complete** |

---

## Test Plan Summary

### Unit Tests (Group 5-6)
| File | Tests |
|------|-------|
| `cmd/api/main_test.go` | Config, validator, Lambda detection |
| `internal/metadata/extractor_test.go` | MP3, FLAC, WAV extraction, cover art |
| `cmd/processor/*/main_test.go` | Individual Lambda handlers |

### Integration Tests (Group 8)
| Test | Environment |
|------|-------------|
| API endpoints | LocalStack + API container |
| Upload workflow | LocalStack + mock Step Functions |

### Test Fixtures Needed
- `testdata/sample.mp3` - MP3 with ID3v2 tags and cover art
- `testdata/sample-notags.mp3` - MP3 without metadata
- `testdata/sample.flac` - FLAC with Vorbis comments
- `testdata/sample.wav` - Raw WAV file

---

## PR Checklist

After completing all tasks in this epic:
- [x] All tests pass (`go test ./...`)
- [x] Code builds (`go build ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Infrastructure validates (`tofu plan`)
- [x] CLAUDE.md files updated for new packages
- [ ] CHANGELOG.md updated
- [ ] Create PR to main with all changes
