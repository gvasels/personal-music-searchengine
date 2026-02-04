# Analyzer Lambda - CLAUDE.md

## Overview

Lambda function for audio analysis in the upload processing pipeline. Extracts BPM and musical key from uploaded audio files using FFmpeg-based signal processing. Designed for graceful degradation - analysis failures don't block the upload workflow.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.go` | Lambda handler with S3 download and analysis orchestration |
| `main_test.go` | Unit tests for Event/Response serialization |

## Step Functions Integration

This Lambda is invoked as part of the upload processor Step Functions workflow:

```
Upload → Metadata → CoverArt → FileMover → Analyzer → TrackCreator → Indexer → Status
                                              ↑
                                          (this Lambda)
```

### Input Event

Received from Step Functions after file is moved to permanent storage:

```go
type Event struct {
    UploadID   string `json:"uploadId"`   // Upload tracking ID
    UserID     string `json:"userId"`     // Owner's user ID
    S3Key      string `json:"s3Key"`      // Permanent S3 location
    FileName   string `json:"fileName"`   // Original filename (for format detection)
    BucketName string `json:"bucketName"` // Media bucket name
}
```

### Output Response

Returned to Step Functions for track creation:

```go
type Response struct {
    BPM        int    `json:"bpm,omitempty"`        // Detected BPM (0 if failed)
    MusicalKey string `json:"musicalKey,omitempty"` // Musical key (e.g., "C major")
    KeyMode    string `json:"keyMode,omitempty"`    // "major" or "minor"
    KeyCamelot string `json:"keyCamelot,omitempty"` // Camelot notation (e.g., "8B")
    Analyzed   bool   `json:"analyzed"`             // Whether analysis succeeded
    Error      string `json:"error,omitempty"`      // Error message if failed
}
```

## Error Handling

The Lambda uses **graceful degradation** - errors return success responses with `Analyzed: false`:

| Error Type | Behavior |
|------------|----------|
| File too large | Returns error message, continues workflow |
| S3 download failure | Returns error message, continues workflow |
| Analysis failure | Returns error message, continues workflow |
| Timeout (>25s) | Returns error message, continues workflow |

This ensures upload processing completes even if audio analysis fails.

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `AWS_REGION` | AWS region for S3 | From Lambda environment |

### Lambda Settings

| Setting | Value | Reason |
|---------|-------|--------|
| Memory | 1024 MB | FFmpeg processing needs memory |
| Timeout | 30 seconds | Analysis takes up to 25s |
| Architecture | arm64 | Cost efficiency |
| Runtime | provided.al2023 | Go custom runtime |

## Dependencies

### Internal
- `internal/analysis` - BPM detection and key analysis algorithms
- `internal/validation` - File size validation (500MB limit)

### External
- `github.com/aws/aws-lambda-go` - Lambda runtime
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 file download

### Infrastructure
- FFmpeg Lambda layer - Required for audio processing

## Usage Example

Step Functions state definition:

```json
{
  "Analyzer": {
    "Type": "Task",
    "Resource": "arn:aws:lambda:us-east-1:123456789:function:music-library-prod-analyzer",
    "Parameters": {
      "uploadId.$": "$.uploadId",
      "userId.$": "$.userId",
      "s3Key.$": "$.newKey",
      "fileName.$": "$.fileName",
      "bucketName.$": "$.bucketName"
    },
    "ResultPath": "$.analysis",
    "Next": "TrackCreator"
  }
}
```

## Testing

```bash
# Run unit tests
cd backend && go test ./cmd/processor/analyzer/...

# Test with coverage
go test -coverprofile=cover.out ./cmd/processor/analyzer/...
```
