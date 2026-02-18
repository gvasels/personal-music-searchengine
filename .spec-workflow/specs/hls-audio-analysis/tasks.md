# HLS-Based Audio Analysis Pipeline — Tasks

## Task 1: transcode-complete audio pipeline trigger

### 1.1 (Red) Write failing tests for transcode-complete changes
- [ ] Create `backend/cmd/processor/transcode/complete/main_test.go`
- [ ] `TestHandleSuccess_TriggersAudioPipeline` — mock DynamoDB GetItem + SFN StartExecution, verify trigger
- [ ] `TestHandleSuccess_NoTriggerWithoutARN` — empty AUDIO_PIPELINE_ARN, verify no trigger
- [ ] `TestHandleFailure_NoTrigger` — ERROR status, verify no audio pipeline trigger
- [ ] `TestReadTrack_ExtractsFields` — verify s3Key, title, artist extraction from DynamoDB response
- [ ] Run tests, verify all FAIL

### 1.2 (Green) Implement transcode-complete changes
- [ ] Add SFN client + `audioPipelineARN` env var to `init()`
- [ ] Add `readTrack()` function (raw DynamoDB GetItem)
- [ ] Add `triggerAudioPipeline()` function
- [ ] Modify `handleSuccess()` to trigger audio pipeline after HLS update
- [ ] Run tests, verify all PASS

## Task 2: track-creator trigger changes

### 2.1 (Red) Write failing tests for track-creator changes
- [ ] Add tests to `backend/cmd/processor/track/main_test.go` (or create)
- [ ] `TestNewTrack_DoesNotTriggerAudioPipeline` — verify SFN NOT called for new tracks
- [ ] `TestDuplicate_WithHLSReady_Triggers` — verify SFN called when HLS ready + needs analysis
- [ ] `TestDuplicate_WithoutHLS_NoTrigger` — verify SFN NOT called when HLS not ready
- [ ] Run tests, verify all FAIL

### 2.2 (Green) Implement track-creator changes
- [ ] Remove new-track audio pipeline trigger (lines 179-182)
- [ ] Add `existing.HLSStatus == models.HLSStatusReady` guard to duplicate trigger (line 109)
- [ ] Run tests, verify all PASS

## Task 3: audio-features HLS sampling

### 3.1 (Red) Write failing tests for audio-features HLS logic
- [ ] Create `backend/cmd/processor/audio-features/test_handler.py`
- [ ] `test_list_hls_segments_groups_by_quality` — mock S3, verify grouping by quality
- [ ] `test_list_hls_segments_empty` — no objects, returns empty dict
- [ ] `test_download_hls_sample_random_selection` — verify random quality + segment selection
- [ ] `test_handler_uses_hls_when_available` — verify HLS preferred over original
- [ ] `test_handler_falls_back_to_original` — no HLS segments, uses s3Key
- [ ] `test_handler_returns_source_field` — verify `source` in result
- [ ] Run tests, verify all FAIL

### 3.2 (Green) Implement audio-features HLS sampling
- [ ] Add `list_hls_segments()` function
- [ ] Add `download_hls_sample()` function with ffmpeg conversion
- [ ] Update `handler()` to try HLS first, fall back to original
- [ ] Add `source` field to result
- [ ] Update Dockerfile to install `ffmpeg-free`
- [ ] Run tests, verify all PASS

## Task 4: Infrastructure + deploy

### 4.1 Update infrastructure
- [ ] Add `AUDIO_PIPELINE_ARN` env var to `transcode_complete` Lambda in `mediaconvert.tf`
- [ ] Add `states:StartExecution` IAM permission to `transcode_lambda` role
- [ ] Run `tofu plan` to verify changes
- [ ] Run `tofu apply`

### 4.2 Build and deploy Lambdas
- [ ] Build transcode-complete: `GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap`
- [ ] Build track-creator: same build command
- [ ] Build audio-features container: `docker build` + ECR push
- [ ] Deploy all 3 Lambda function code updates

## Task 5: Data wipe + verification

### 5.1 Create wipe script
- [ ] Create `scripts/wipe-library.sh`
- [ ] Implement DynamoDB scan + batch delete for TRACK# and UPLOAD# items
- [ ] Implement S3 recursive delete for media/, hls/, uploads/, coverart/
- [ ] Add --dry-run and --confirm safety flags

### 5.2 Wipe and verify
- [ ] Run wipe script with --dry-run to confirm scope
- [ ] Run wipe script with --confirm to execute
- [ ] Verify DynamoDB empty (no TRACK# items)
- [ ] Verify S3 prefixes empty

### 5.3 End-to-end verification
- [ ] Re-upload test track through UI
- [ ] Verify upload pipeline completes (track created, file moved, HLS transcoded)
- [ ] Verify transcode-complete triggers audio pipeline (CloudWatch logs)
- [ ] Verify audio-features uses HLS segments (source: hls in result)
- [ ] Verify DynamoDB track has: bpm, musicalKey, keyCamelot, genre, mood, embeddingId
- [ ] Compare BPM/key against known values for test tracks
