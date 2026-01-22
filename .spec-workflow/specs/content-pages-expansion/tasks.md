# Tasks - Content Pages Expansion

## Feature: Content Pages Expansion (Videos & Live Streams)
**Status**: Not Started
**Estimated Effort**: 4-5 weeks

---

## Group 22: Video Infrastructure

### Task 22.1: Video Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/video.go`
- `backend/internal/models/video_test.go`

**Acceptance Criteria**:
- [ ] Video struct with all fields
- [ ] VideoItem DynamoDB struct (PK/SK)
- [ ] VideoStatus constants
- [ ] ToResponse method
- [ ] Unit tests

### Task 22.2: Video Repository
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/video.go`
- `backend/internal/repository/video_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `CreateVideo(ctx, video)` | Create new video record |
| `GetVideo(ctx, userID, videoID)` | Get video by ID |
| `ListVideos(ctx, userID, filter)` | List user's videos |
| `UpdateVideo(ctx, video)` | Update video metadata/status |
| `DeleteVideo(ctx, userID, videoID)` | Delete video |

**Acceptance Criteria**:
- [ ] CRUD operations for videos
- [ ] Filter by status
- [ ] Pagination support

### Task 22.3: Video Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/video.go`
- `backend/internal/service/video_test.go`

**Acceptance Criteria**:
- [ ] Business logic for video operations
- [ ] Presigned URL generation
- [ ] Status transitions
- [ ] Association with tracks/albums

### Task 22.4: Video Handlers
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/video.go`
- `backend/internal/handlers/video_test.go`

**Endpoints**:
| Endpoint | Handler |
|----------|---------|
| `GET /videos` | ListVideos |
| `GET /videos/{id}` | GetVideo |
| `POST /videos/upload/presigned` | GetPresignedUpload |
| `POST /videos/upload/confirm` | ConfirmUpload |
| `PUT /videos/{id}` | UpdateVideo |
| `DELETE /videos/{id}` | DeleteVideo |
| `GET /videos/{id}/stream` | GetStreamURL |

**Acceptance Criteria**:
- [ ] All endpoints implemented
- [ ] Validation
- [ ] Error handling
- [ ] Integration tests

---

## Group 23: Video Transcoding Pipeline

### Task 23.1: MediaConvert Job Template
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/mediaconvert-video.tf`

**Acceptance Criteria**:
- [ ] HLS output (1080p, 720p, 360p)
- [ ] Thumbnail extraction
- [ ] Job queue configuration
- [ ] IAM roles

### Task 23.2: Video Upload Trigger Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/video-upload/main.go`

**Acceptance Criteria**:
- [ ] Triggers on S3 video upload
- [ ] Submits MediaConvert job
- [ ] Updates video status to PROCESSING

### Task 23.3: Video Status Update Lambda
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/video-status/main.go`

**Acceptance Criteria**:
- [ ] Triggers on MediaConvert completion
- [ ] Updates video with HLS/thumbnail paths
- [ ] Sets status to READY or FAILED

### Task 23.4: Video Infrastructure
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/lambda-video.tf`
- `infrastructure/backend/s3-video.tf`

**Acceptance Criteria**:
- [ ] Lambda functions deployed
- [ ] S3 event notifications
- [ ] EventBridge rules for MediaConvert
- [ ] CloudFront configuration

---

## Group 24: Live Stream Infrastructure

### Task 24.1: LiveStream Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/stream.go`
- `backend/internal/models/stream_test.go`

**Acceptance Criteria**:
- [ ] LiveStream struct
- [ ] StreamStatus constants
- [ ] IVS integration fields
- [ ] Unit tests

### Task 24.2: IVS Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/ivs/client.go`
- `backend/internal/ivs/client_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `CreateChannel(ctx, name)` | Create IVS channel |
| `DeleteChannel(ctx, arn)` | Delete IVS channel |
| `GetStreamKey(ctx, arn)` | Get stream key |
| `GetChannel(ctx, arn)` | Get channel details |

**Acceptance Criteria**:
- [ ] IVS SDK integration
- [ ] Channel lifecycle management
- [ ] Error handling

### Task 24.3: Stream Repository
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/stream.go`
- `backend/internal/repository/stream_test.go`

**Acceptance Criteria**:
- [ ] CRUD operations for streams
- [ ] Filter by status
- [ ] Pagination support

### Task 24.4: Stream Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/stream.go`
- `backend/internal/service/stream_test.go`

**Acceptance Criteria**:
- [ ] Create channel + IVS channel
- [ ] Start/stop stream
- [ ] Viewer count updates
- [ ] Recording management

### Task 24.5: Stream Handlers
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/stream.go`
- `backend/internal/handlers/stream_test.go`

**Endpoints**:
| Endpoint | Handler |
|----------|---------|
| `GET /streams` | ListStreams |
| `GET /streams/{id}` | GetStream |
| `POST /streams` | CreateStream |
| `PUT /streams/{id}` | UpdateStream |
| `DELETE /streams/{id}` | DeleteStream |
| `GET /streams/{id}/ingest` | GetIngestURL |
| `POST /streams/{id}/start` | StartStream |
| `POST /streams/{id}/stop` | StopStream |

**Acceptance Criteria**:
- [ ] All endpoints implemented
- [ ] IVS integration working
- [ ] Error handling

### Task 24.6: IVS Infrastructure
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/ivs.tf`

**Resources**:
- [ ] IVS channel (on-demand via Lambda)
- [ ] IVS recording configuration
- [ ] S3 bucket for recordings
- [ ] IAM roles

---

## Group 25: Live Chat

### Task 25.1: ChatMessage Model
**Status**: [ ] Pending
**Files**:
- `backend/internal/models/chat.go`
- `backend/internal/models/chat_test.go`

**Acceptance Criteria**:
- [ ] ChatMessage struct
- [ ] DynamoDB item with TTL
- [ ] Unit tests

### Task 25.2: Chat Repository
**Status**: [ ] Pending
**Files**:
- `backend/internal/repository/chat.go`
- `backend/internal/repository/chat_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `SaveMessage(ctx, message)` | Store chat message |
| `GetMessages(ctx, streamID, limit)` | Get recent messages |

**Acceptance Criteria**:
- [ ] Message persistence
- [ ] TTL for auto-deletion
- [ ] Recent messages query

### Task 25.3: WebSocket Connect Handler
**Status**: [ ] Pending
**Files**:
- `backend/cmd/chat/connect/main.go`

**Acceptance Criteria**:
- [ ] Store connection ID in DynamoDB
- [ ] Associate with stream ID
- [ ] Validate stream exists

### Task 25.4: WebSocket Message Handler
**Status**: [ ] Pending
**Files**:
- `backend/cmd/chat/message/main.go`

**Acceptance Criteria**:
- [ ] Parse incoming message
- [ ] Store in DynamoDB
- [ ] Broadcast to all connections

### Task 25.5: WebSocket Disconnect Handler
**Status**: [ ] Pending
**Files**:
- `backend/cmd/chat/disconnect/main.go`

**Acceptance Criteria**:
- [ ] Remove connection from DynamoDB
- [ ] Update viewer count

### Task 25.6: WebSocket Infrastructure
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/websocket.tf`

**Resources**:
- [ ] API Gateway WebSocket API
- [ ] Routes ($connect, $disconnect, message)
- [ ] Lambda integrations
- [ ] DynamoDB table for connections

---

## Group 26: Frontend - Videos

### Task 26.1: Video Types
**Status**: [ ] Pending
**Files**:
- `frontend/src/types/index.ts`

**Acceptance Criteria**:
- [ ] Video interface
- [ ] VideoStatus type
- [ ] API response types

### Task 26.2: Video API Functions
**Status**: [ ] Pending
**Files**:
- `frontend/src/lib/api/videos.ts`

**Functions**:
- [ ] getVideos, getVideo
- [ ] getPresignedUpload, confirmUpload
- [ ] updateVideo, deleteVideo
- [ ] getStreamURL

### Task 26.3: Video Hooks
**Status**: [ ] Pending
**Files**:
- `frontend/src/hooks/useVideos.ts`

**Acceptance Criteria**:
- [ ] useVideosQuery
- [ ] useVideoQuery
- [ ] useUploadVideoMutation
- [ ] useDeleteVideoMutation

### Task 26.4: Videos List Page
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/videos/index.tsx`

**Acceptance Criteria**:
- [ ] Video grid with thumbnails
- [ ] Status badges
- [ ] Link to detail page
- [ ] Upload button

### Task 26.5: Video Detail Page
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/videos/$videoId.tsx`

**Acceptance Criteria**:
- [ ] HLS video player
- [ ] Video metadata display
- [ ] Edit/delete actions
- [ ] Associated track/album links

### Task 26.6: Video Upload Page
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/videos/upload.tsx`

**Acceptance Criteria**:
- [ ] Drag-and-drop upload
- [ ] Progress bar
- [ ] Metadata form
- [ ] Track/album association

### Task 26.7: Video Player Component
**Status**: [ ] Pending
**Files**:
- `frontend/src/components/video/VideoPlayer.tsx`

**Dependencies**: `hls.js` or `video.js`

**Acceptance Criteria**:
- [ ] HLS playback
- [ ] Quality selection
- [ ] Fullscreen support
- [ ] Responsive design

---

## Group 27: Frontend - Live Streams

### Task 27.1: Stream Types
**Status**: [ ] Pending
**Files**:
- `frontend/src/types/index.ts`

**Acceptance Criteria**:
- [ ] LiveStream interface
- [ ] StreamStatus type
- [ ] ChatMessage interface

### Task 27.2: Stream API Functions
**Status**: [ ] Pending
**Files**:
- `frontend/src/lib/api/streams.ts`

**Functions**:
- [ ] getStreams, getStream
- [ ] createStream, updateStream, deleteStream
- [ ] getIngestURL
- [ ] startStream, stopStream

### Task 27.3: Stream Hooks
**Status**: [ ] Pending
**Files**:
- `frontend/src/hooks/useStreams.ts`

**Acceptance Criteria**:
- [ ] useStreamsQuery
- [ ] useStreamQuery
- [ ] useCreateStreamMutation
- [ ] useStartStreamMutation

### Task 27.4: Streams List Page
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/streams/index.tsx`

**Acceptance Criteria**:
- [ ] Stream cards with status
- [ ] Live badge for active streams
- [ ] Scheduled streams section
- [ ] Create stream button

### Task 27.5: Stream Detail/Watch Page
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/streams/$streamId.tsx`

**Acceptance Criteria**:
- [ ] IVS player
- [ ] Live chat panel
- [ ] Viewer count
- [ ] Stream info

### Task 27.6: Create Stream Page
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/streams/create.tsx`

**Acceptance Criteria**:
- [ ] Stream title/description form
- [ ] Thumbnail upload
- [ ] Schedule option
- [ ] Recording toggle

### Task 27.7: Broadcast Dashboard
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/streams/$streamId/broadcast.tsx`

**Acceptance Criteria**:
- [ ] RTMP URL and stream key display
- [ ] Copy buttons
- [ ] Go live/end stream buttons
- [ ] Stream health status
- [ ] Chat moderation

### Task 27.8: Chat Component
**Status**: [ ] Pending
**Files**:
- `frontend/src/components/stream/ChatPanel.tsx`
- `frontend/src/hooks/useChat.ts`

**Acceptance Criteria**:
- [ ] WebSocket connection
- [ ] Message display
- [ ] Send message input
- [ ] Auto-scroll
- [ ] Reconnection logic

### Task 27.9: IVS Player Component
**Status**: [ ] Pending
**Files**:
- `frontend/src/components/stream/LivePlayer.tsx`

**Dependencies**: `amazon-ivs-player`

**Acceptance Criteria**:
- [ ] IVS player integration
- [ ] Low-latency playback
- [ ] Quality selection
- [ ] Fullscreen support

---

## Group 28: Integration & Testing

### Task 28.1: API Route Registration
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/handlers.go`
- `infrastructure/backend/api-gateway.tf`

**Acceptance Criteria**:
- [ ] Video routes registered
- [ ] Stream routes registered
- [ ] API Gateway updated

### Task 28.2: Frontend Route Registration
**Status**: [ ] Pending
**Files**:
- `frontend/src/routes/__root.tsx`
- `frontend/src/components/layout/Sidebar.tsx`

**Acceptance Criteria**:
- [ ] Videos and Streams in sidebar
- [ ] Routes configured
- [ ] Navigation working

### Task 28.3: End-to-End Tests
**Status**: [ ] Pending
**Files**:
- `e2e/videos.spec.ts`
- `e2e/streams.spec.ts`

**Acceptance Criteria**:
- [ ] Video upload and playback
- [ ] Stream creation
- [ ] Chat functionality

### Task 28.4: Cost Monitoring
**Status**: [ ] Pending
**Files**:
- `infrastructure/backend/cloudwatch-alarms.tf`

**Acceptance Criteria**:
- [ ] MediaConvert spend alarm
- [ ] IVS viewing hours alarm
- [ ] Budget alerts configured

---

## Summary

| Group | Tasks | Purpose |
|-------|-------|---------|
| Group 22 | 4 | Video backend infrastructure |
| Group 23 | 4 | Video transcoding pipeline |
| Group 24 | 6 | Live stream infrastructure |
| Group 25 | 6 | Live chat system |
| Group 26 | 7 | Frontend videos |
| Group 27 | 9 | Frontend live streams |
| Group 28 | 4 | Integration & testing |
| **Total** | **40** | |

---

## Dependencies

### Task Dependencies
```
Group 22 → Group 23 (video processing needs model)
Group 24 → Group 25 (chat needs stream model)
Groups 22-25 → Group 26-27 (frontend needs backend)
All → Group 28 (integration last)
```

### External Dependencies
- `hls.js` or `video.js` for video playback
- `amazon-ivs-player` for live streams
- AWS MediaConvert access
- AWS IVS access

---

## PR Checklist

After completing all tasks:
- [ ] All tests pass
- [ ] Video upload and playback works
- [ ] Live streaming works
- [ ] Chat works in real-time
- [ ] Cost alarms configured
- [ ] CLAUDE.md files updated
- [ ] CHANGELOG.md updated
- [ ] Create PR to main
