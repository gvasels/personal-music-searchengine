# Design - Content Pages Expansion

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Video Upload Pipeline                                │
│                                                                              │
│  Frontend ──► S3 Presigned ──► S3 Upload Bucket                             │
│                                      │                                       │
│                                      ▼                                       │
│                              ┌─────────────────┐                            │
│                              │  S3 Event       │                            │
│                              │  Notification   │                            │
│                              └────────┬────────┘                            │
│                                       │                                      │
│                                       ▼                                      │
│                              ┌─────────────────┐                            │
│                              │  MediaConvert   │                            │
│                              │  Job Trigger    │                            │
│                              └────────┬────────┘                            │
│                                       │                                      │
│                                       ▼                                      │
│  ┌─────────────────┐        ┌─────────────────┐        ┌─────────────────┐ │
│  │ Thumbnail Gen   │◄──────│  MediaConvert   │───────►│ HLS Segments    │ │
│  │ (Lambda)        │        │  Transcode      │        │ (S3 Media)      │ │
│  └────────┬────────┘        └────────┬────────┘        └────────┬────────┘ │
│           │                          │                          │           │
│           └──────────────────────────┼──────────────────────────┘           │
│                                      │                                      │
│                                      ▼                                      │
│                              ┌─────────────────┐                            │
│                              │  Update Video   │                            │
│                              │  Status Lambda  │                            │
│                              └────────┬────────┘                            │
│                                       │                                      │
│                                       ▼                                      │
│                              ┌─────────────────┐                            │
│                              │   CloudFront    │───────► Users              │
│                              │   Delivery      │                            │
│                              └─────────────────┘                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                         Live Streaming (IVS)                                 │
│                                                                              │
│  Broadcaster (OBS) ──► RTMP Ingest ──► Amazon IVS Channel                   │
│                                              │                               │
│                                              ├───────► HLS Playback ──► CDN │
│                                              │                               │
│                                              └───────► Recording ──► S3     │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Real-Time Chat                                │   │
│  │                                                                      │   │
│  │  Viewer ──► WebSocket API Gateway ──► Lambda ──► DynamoDB           │   │
│  │                     │                                                │   │
│  │                     └──► Broadcast to all connections               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Design Decisions

### DD-1: Video Transcoding Service

**Decision**: Use AWS MediaConvert

**Rationale**:
- Existing MediaConvert setup for audio HLS transcoding
- Supports video HLS with multiple quality levels
- Thumbnail extraction built-in
- Pay-per-use pricing
- No infrastructure management

**Alternative Considered**: AWS Elastic Transcoder
- Legacy service, limited features
- MediaConvert is recommended replacement

### DD-2: Live Streaming Service

**Decision**: Use Amazon IVS (Interactive Video Service)

**Rationale**:
- Purpose-built for live streaming
- Managed RTMP ingest
- Low-latency HLS delivery (< 5 seconds)
- Built-in recording to S3
- Simple API for channel management

**Alternative Considered**: AWS MediaLive + MediaPackage
- More complex setup
- Better for broadcast-grade features
- Overkill for personal streaming

### DD-3: Chat Architecture

**Decision**: API Gateway WebSocket API + Lambda + DynamoDB

**Rationale**:
- Matches existing API Gateway HTTP API pattern
- Serverless, scales automatically
- Simple implementation
- DynamoDB for message persistence

**Alternative Considered**: AWS AppSync
- Adds GraphQL complexity
- Better if we need GraphQL subscriptions later
- Consider for Phase 2

### DD-4: Video Storage Strategy

**Decision**: Separate buckets for upload vs media

**Rationale**:
- Upload bucket: temporary, lifecycle policy deletes after processing
- Media bucket: permanent, CloudFront origin
- Prevents serving unprocessed files
- Clear separation of concerns

**Storage Structure**:
```
Upload Bucket (music-library-uploads):
  videos/{userId}/{uploadId}/original.mp4

Media Bucket (music-library-media):
  videos/{userId}/{videoId}/
    ├── master.m3u8
    ├── 1080p/
    │   ├── playlist.m3u8
    │   └── segment_*.ts
    ├── 720p/
    │   └── ...
    ├── 360p/
    │   └── ...
    └── thumbnails/
        ├── thumb_001.jpg
        ├── thumb_002.jpg
        └── poster.jpg
```

### DD-5: GameLift Streams (Phase 2)

**Decision**: Defer to Phase 2

**Rationale**:
- Significant complexity and cost
- Requires custom streaming application
- IVS covers most live streaming needs
- Can add later without major changes

---

## Component Design

### 1. Video Model

**File**: `backend/internal/models/video.go`

```go
type Video struct {
    ID           string       `json:"id" dynamodbav:"id"`
    UserID       string       `json:"userId" dynamodbav:"userId"`
    Title        string       `json:"title" dynamodbav:"title"`
    Description  string       `json:"description,omitempty" dynamodbav:"description,omitempty"`
    Artist       string       `json:"artist,omitempty" dynamodbav:"artist,omitempty"`
    Duration     int          `json:"duration" dynamodbav:"duration"`
    Resolution   string       `json:"resolution" dynamodbav:"resolution"`
    Format       string       `json:"format" dynamodbav:"format"`
    FileSize     int64        `json:"fileSize" dynamodbav:"fileSize"`
    ThumbnailKey string       `json:"thumbnailKey" dynamodbav:"thumbnailKey"`
    HLSKey       string       `json:"hlsKey" dynamodbav:"hlsKey"`
    Status       VideoStatus  `json:"status" dynamodbav:"status"`
    TrackID      string       `json:"trackId,omitempty" dynamodbav:"trackId,omitempty"`
    AlbumID      string       `json:"albumId,omitempty" dynamodbav:"albumId,omitempty"`
    Timestamps
}

type VideoStatus string

const (
    VideoStatusUploading   VideoStatus = "UPLOADING"
    VideoStatusProcessing  VideoStatus = "PROCESSING"
    VideoStatusReady       VideoStatus = "READY"
    VideoStatusFailed      VideoStatus = "FAILED"
)
```

### 2. LiveStream Model

**File**: `backend/internal/models/stream.go`

```go
type LiveStream struct {
    ID               string       `json:"id" dynamodbav:"id"`
    UserID           string       `json:"userId" dynamodbav:"userId"`
    Title            string       `json:"title" dynamodbav:"title"`
    Description      string       `json:"description,omitempty" dynamodbav:"description,omitempty"`
    ThumbnailKey     string       `json:"thumbnailKey,omitempty" dynamodbav:"thumbnailKey,omitempty"`
    Status           StreamStatus `json:"status" dynamodbav:"status"`
    ViewerCount      int          `json:"viewerCount" dynamodbav:"viewerCount"`
    IVSChannelArn    string       `json:"ivsChannelArn" dynamodbav:"ivsChannelArn"`
    IVSStreamKey     string       `json:"-" dynamodbav:"ivsStreamKey"` // Never expose
    RTMPIngestUrl    string       `json:"rtmpIngestUrl,omitempty" dynamodbav:"rtmpIngestUrl"`
    PlaybackUrl      string       `json:"playbackUrl" dynamodbav:"playbackUrl"`
    RecordingEnabled bool         `json:"recordingEnabled" dynamodbav:"recordingEnabled"`
    ScheduledAt      *time.Time   `json:"scheduledAt,omitempty" dynamodbav:"scheduledAt,omitempty"`
    StartedAt        *time.Time   `json:"startedAt,omitempty" dynamodbav:"startedAt,omitempty"`
    EndedAt          *time.Time   `json:"endedAt,omitempty" dynamodbav:"endedAt,omitempty"`
    Timestamps
}

type StreamStatus string

const (
    StreamStatusOffline    StreamStatus = "OFFLINE"
    StreamStatusScheduled  StreamStatus = "SCHEDULED"
    StreamStatusLive       StreamStatus = "LIVE"
    StreamStatusEnded      StreamStatus = "ENDED"
)
```

### 3. ChatMessage Model

**File**: `backend/internal/models/chat.go`

```go
type ChatMessage struct {
    ID        string    `json:"id" dynamodbav:"id"`
    StreamID  string    `json:"streamId" dynamodbav:"streamId"`
    UserID    string    `json:"userId" dynamodbav:"userId"`
    Username  string    `json:"username" dynamodbav:"username"`
    Message   string    `json:"message" dynamodbav:"message"`
    Timestamp time.Time `json:"timestamp" dynamodbav:"timestamp"`
}
```

---

## DynamoDB Schema

### Video Entity

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `USER#{userId}` |
| `SK` | String | `VIDEO#{videoId}` |
| `Type` | String | `VIDEO` |
| `GSI1PK` | String | `VIDEO#STATUS#{status}` |
| `GSI1SK` | String | `{createdAt}` |
| All video fields | ... | ... |

### LiveStream Entity

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `USER#{userId}` |
| `SK` | String | `STREAM#{streamId}` |
| `Type` | String | `STREAM` |
| `GSI1PK` | String | `STREAM#STATUS#{status}` |
| `GSI1SK` | String | `{scheduledAt or createdAt}` |
| All stream fields | ... | ... |

### ChatMessage Entity

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `STREAM#{streamId}` |
| `SK` | String | `MSG#{timestamp}#{messageId}` |
| `Type` | String | `CHAT_MESSAGE` |
| `TTL` | Number | Unix timestamp (7 days) |
| All message fields | ... | ... |

---

## Frontend Routes

### New Routes

| Route | Component | Description |
|-------|-----------|-------------|
| `/videos` | `VideosPage` | Video grid listing |
| `/videos/{videoId}` | `VideoDetailPage` | Video player + details |
| `/videos/upload` | `VideoUploadPage` | Upload form |
| `/streams` | `StreamsPage` | Stream listing |
| `/streams/{streamId}` | `StreamDetailPage` | Live player + chat |
| `/streams/create` | `CreateStreamPage` | Create new channel |
| `/streams/{streamId}/broadcast` | `BroadcastPage` | Broadcaster controls |

### New Components

| Component | Purpose |
|-----------|---------|
| `VideoCard` | Video thumbnail with title/duration |
| `VideoPlayer` | HLS video player (video.js or hls.js) |
| `VideoUploader` | Upload form with progress |
| `StreamCard` | Stream thumbnail with status badge |
| `LivePlayer` | IVS player component |
| `ChatPanel` | Real-time chat sidebar |
| `BroadcastControls` | Stream key, status, settings |

---

## Infrastructure

### New Lambda Functions

| Lambda | Runtime | Memory | Timeout | Purpose |
|--------|---------|--------|---------|---------|
| `video-upload-trigger` | Go | 256MB | 30s | Triggers MediaConvert on upload |
| `video-status-update` | Go | 256MB | 30s | Updates video status after transcode |
| `stream-handler` | Go | 256MB | 30s | IVS channel management |
| `chat-connect` | Go | 128MB | 10s | WebSocket connect handler |
| `chat-message` | Go | 128MB | 10s | WebSocket message handler |
| `chat-disconnect` | Go | 128MB | 10s | WebSocket disconnect handler |

### New AWS Resources

| Resource | Purpose |
|----------|---------|
| `aws_media_convert_queue` | Video transcoding queue |
| `aws_ivs_channel` | Live stream channel (per user) |
| `aws_ivs_recording_configuration` | Record to S3 |
| `aws_apigatewayv2_api` | WebSocket API for chat |
| `aws_dynamodb_table` (update) | Add VIDEO, STREAM, CHAT_MESSAGE entities |

### MediaConvert Job Template

```json
{
  "OutputGroups": [
    {
      "Name": "HLS",
      "OutputGroupSettings": {
        "Type": "HLS_GROUP_SETTINGS",
        "HlsGroupSettings": {
          "SegmentLength": 6,
          "Destination": "s3://music-library-media/videos/{userId}/{videoId}/"
        }
      },
      "Outputs": [
        { "Preset": "1080p_5000kbps" },
        { "Preset": "720p_3000kbps" },
        { "Preset": "360p_1000kbps" }
      ]
    },
    {
      "Name": "Thumbnails",
      "OutputGroupSettings": {
        "Type": "FILE_GROUP_SETTINGS"
      },
      "Outputs": [
        {
          "ContainerSettings": { "Container": "RAW" },
          "VideoDescription": {
            "CodecSettings": {
              "Codec": "FRAME_CAPTURE",
              "FrameCaptureSettings": {
                "FramerateNumerator": 1,
                "FramerateDenominator": 30,
                "MaxCaptures": 10
              }
            }
          }
        }
      ]
    }
  ]
}
```

---

## Error Handling

### Video Processing Errors

| Error | Handling | User Impact |
|-------|----------|-------------|
| Upload interrupted | Multipart cleanup, retry | "Upload failed, try again" |
| Invalid format | Reject before transcode | "Unsupported format" |
| Transcode failure | Mark as FAILED, retry option | "Processing failed" |
| S3 access error | Log, alert | "Temporary error" |

### Live Stream Errors

| Error | Handling | User Impact |
|-------|----------|-------------|
| IVS channel limit | Return error | "Stream limit reached" |
| Stream disconnect | Update status | "Stream ended" |
| Chat WebSocket error | Reconnect client | "Reconnecting..." |

---

## Monitoring

### CloudWatch Metrics

- `VideoUploadsPerMinute`
- `TranscodeJobDuration`
- `TranscodeFailureRate`
- `ActiveStreams`
- `StreamViewerCount`
- `ChatMessagesPerMinute`
- `WebSocketConnections`

### CloudWatch Alarms

- MediaConvert spend > $X/day
- IVS viewing hours > Y/day
- Transcode failure rate > 5%
- WebSocket errors > 10/minute

---

## Cost Estimates

### Per Video (5 min, 1080p)

| Service | Cost |
|---------|------|
| MediaConvert | ~$0.10 |
| S3 Storage | ~$0.01/month |
| CloudFront (100 views) | ~$0.05 |

### Per Live Stream (1 hour)

| Service | Cost |
|---------|------|
| IVS Input | ~$0.90 |
| IVS Viewing (100 viewers) | ~$4.00 |
| S3 Recording | ~$0.02 |

### Recommendations

1. Implement cost alerts at $10, $50, $100/day
2. Auto-stop inactive streams after 30 minutes
3. Delete transcoding artifacts after 24 hours
4. Consider IVS basic channel tier for lower costs
