# HLS-Based Audio Analysis Pipeline — Design

## Architecture Change

**Before:**
```
CreateTrackRecord ──[triggers audio pipeline]──► AudioFeatures (reads original file)
         ↓
MoveToMediaStorage → StartTranscode → ... → EventBridge → transcode-complete (DynamoDB only)
```

**After:**
```
CreateTrackRecord ──[NO trigger]──► MoveToMediaStorage → StartTranscode → ...
                                                            ↓ (async)
EventBridge → transcode-complete ──[triggers audio pipeline]──► AudioFeatures (reads HLS segments)
```

## Code Reuse Analysis

| Existing Code | Reuse In |
|---------------|----------|
| `triggerAudioPipeline()` in `track/main.go` | Same pattern replicated in `transcode/complete/main.go` |
| Raw DynamoDB GetItem pattern in `transcode/complete/main.go` | Extended for track read (same file, same client) |
| S3 list/download in `audio-features/handler.py` | Extended with `list_objects_v2` + paginator |
| `service.BuildHLSPlaylistKey()` | Referenced for path construction in audio-features |

## Component Changes

### 1. transcode-complete Lambda

**File:** `backend/cmd/processor/transcode/complete/main.go`

**New dependencies:** `sfn` SDK client, `encoding/json`

**New functions:**
```go
func readTrack(ctx context.Context, userID, trackID string) (s3Key, title, artist string, err error)
func triggerAudioPipeline(ctx context.Context, trackID, userID, s3Key, title, artist string)
```

**Modified function:** `handleSuccess()` — append audio pipeline trigger after DynamoDB update

**Design decision:** Use raw DynamoDB GetItem (not repository package) for consistency with existing code in this file. The repository import would add unnecessary coupling.

### 2. track-creator Lambda

**File:** `backend/cmd/processor/track/main.go`

**Removed:** Lines 179-182 (new-track audio pipeline trigger)
**Modified:** Line 109 — add `existing.HLSStatus == models.HLSStatusReady` guard

### 3. audio-features Lambda

**File:** `backend/cmd/processor/audio-features/handler.py`

**New functions:**
```python
def list_hls_segments(user_id: str, track_id: str) -> dict[str, list[str]]
def download_hls_sample(segments_by_quality: dict) -> tuple[str, str]
```

**Modified function:** `handler()` — try HLS first, fall back to original

**HLS segment parsing:**
- Segment filename format: `{quality}_{index}.ts` (e.g., `320k_00001.ts`)
- Quality extracted by splitting on `_` and taking first part
- Segments sorted by key for temporal ordering before analysis

**ffmpeg conversion:**
- TS segments concatenated into single `.ts` file
- Converted to mono WAV at 22050Hz via `ffmpeg -y -i combined.ts -ac 1 -ar 22050 combined.wav`
- Required because libsndfile (librosa's backend) cannot read MPEG-TS containers

### 4. Infrastructure

**File:** `infrastructure/backend/mediaconvert.tf`

- Add `AUDIO_PIPELINE_ARN` env var to `transcode_complete` Lambda
- Add `states:StartExecution` IAM permission to `transcode_lambda` role
- Resource scoped to `aws_sfn_state_machine.audio_pipeline.arn`

### 5. Wipe Script

**File:** `scripts/wipe-library.sh`

- Follows pattern of existing scripts in `scripts/migrations/`
- DynamoDB: Scan with FilterExpression on SK prefix, BatchWriteItem delete
- S3: `aws s3 rm --recursive` for each prefix

## Data Flow

### New Track Upload (normal flow)
```
1. Upload → ExtractMetadata → ProcessCoverArt → CreateTrackRecord (no audio trigger)
2. → CheckDuplicate → MoveToMediaStorage → StartTranscode (async MediaConvert job)
3. → IndexForSearch → MarkUploadCompleted
4. [async] MediaConvert completes → EventBridge → transcode-complete
5. transcode-complete: update hlsStatus=READY, read track, trigger audio pipeline
6. Audio pipeline: AudioFeatures (HLS segments) → AudioAnalyzer → EmbeddingGenerator → TrackUpdater
```

### Duplicate Upload (delta analysis)
```
1. Upload → CreateTrackRecord → finds duplicate
2. If existing.HLSStatus == READY && needsAnalysis(existing):
   → triggerAudioPipeline(existing track data)
3. Cleanup uploaded file, return isDuplicate=true
4. CheckDuplicate → MarkUploadCompleted (skip remaining steps)
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| readTrack fails in transcode-complete | Log warning, skip audio pipeline, return success |
| StartExecution fails in transcode-complete | Log warning, return success (best-effort) |
| No HLS segments found in audio-features | Fall back to original s3Key |
| ffmpeg conversion fails | Return error in result, pipeline continues to next step |
| S3 list_objects_v2 fails | Return error in result, pipeline continues |
