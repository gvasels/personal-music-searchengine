# Requirements - Content Pages Expansion

## Feature Overview

**Feature**: Content Pages Expansion (Videos & Live Streams)
**Priority**: Medium
**Dependencies**: Backend API, Frontend, Infrastructure

Expand the music library with video content and live streaming capabilities. This includes uploading music videos, managing live streams with chat, and supporting both standard streaming (IVS) and interactive experiences (GameLift - Phase 2).

---

## Functional Requirements

### FR-1: Videos Page

**FR-1.1**: Video Upload
- Upload video files up to 10GB
- Supported formats: MP4, WebM, MOV, MKV
- Multipart upload for large files
- Progress tracking during upload

**FR-1.2**: Video Metadata
- Title, description, artist association
- Duration (auto-detected)
- Resolution (auto-detected)
- Thumbnail (auto-generated or custom)
- Associate with tracks, albums, or artists

**FR-1.3**: Video Transcoding
- Automatic HLS transcoding via MediaConvert
- Multiple quality levels (360p, 720p, 1080p)
- Thumbnail generation at key frames
- Processing status tracking

**FR-1.4**: Video Playback
- HLS adaptive bitrate streaming
- Standard video controls (play, pause, seek, volume)
- Fullscreen support
- Quality selection
- CloudFront CDN delivery

**FR-1.5**: Video Management
- List all videos with thumbnails
- Edit video metadata
- Delete videos
- Associate/disassociate with tracks/albums

### FR-2: Live Streams Page (Regular - IVS)

**FR-2.1**: Stream Creation
- Create new live stream channel
- Generate RTMP ingest URL and stream key
- Configure stream settings (title, description, thumbnail)
- Schedule future streams

**FR-2.2**: Stream Broadcasting
- RTMP ingest from OBS, Streamlabs, etc.
- Low-latency HLS playback
- Automatic recording option
- Stream health monitoring

**FR-2.3**: Stream Viewing
- Live video player with chat
- Viewer count display
- Stream status (live, offline, scheduled)
- Past recordings playback

**FR-2.4**: Live Chat
- Real-time chat during streams
- WebSocket-based messaging
- Username display
- Basic moderation (delete messages)

### FR-3: Live Streams (GameLift - Phase 2)

**FR-3.1**: GameLift Integration
- Low-latency game streaming infrastructure
- Interactive music experiences (VJ software, rhythm games)
- Real-time input handling
- Session management

**FR-3.2**: Audience Participation
- Interactive controls for viewers
- Real-time voting/reactions
- Synchronized audio/visual experiences

*Note: GameLift Streams is Phase 2 and requires separate detailed spec*

---

## Non-Functional Requirements

### NFR-1: Performance

**NFR-1.1**: Video upload completes within 10 minutes for 1GB file
**NFR-1.2**: Transcoding starts within 30 seconds of upload
**NFR-1.3**: Live stream latency < 5 seconds (IVS low-latency mode)
**NFR-1.4**: Chat message delivery < 500ms

### NFR-2: Scalability

**NFR-2.1**: Support 1000 concurrent video viewers
**NFR-2.2**: Support 100 concurrent stream viewers per channel
**NFR-2.3**: Support 1000 chat messages per minute

### NFR-3: Reliability

**NFR-3.1**: Video playback 99.9% availability (CloudFront)
**NFR-3.2**: Stream recordings persisted after broadcast
**NFR-3.3**: Graceful handling of stream disconnections

### NFR-4: Cost

**NFR-4.1**: MediaConvert charged per minute of output
**NFR-4.2**: IVS charged per hour of streaming + viewing
**NFR-4.3**: CloudWatch alarms for cost monitoring
**NFR-4.4**: Automatic cleanup of old transcoding jobs

---

## User Stories

### US-1: Upload Music Video
```
As a user, I want to upload music videos
So that I can associate visual content with my music
```

**Acceptance Criteria**:
- Can upload MP4, WebM, MOV files
- Progress bar during upload
- Automatic thumbnail generation
- Can associate with track/album

### US-2: Watch Music Video
```
As a user, I want to watch my uploaded videos
So that I can enjoy visual content with my music
```

**Acceptance Criteria**:
- Adaptive quality streaming
- Standard playback controls
- Fullscreen support

### US-3: Start Live Stream
```
As a DJ/Artist, I want to go live
So that I can perform for my audience
```

**Acceptance Criteria**:
- Get RTMP URL and stream key
- Connect with OBS/Streamlabs
- See live viewer count
- Interact via chat

### US-4: Watch Live Stream
```
As a viewer, I want to watch live streams
So that I can enjoy live performances
```

**Acceptance Criteria**:
- Low-latency video playback
- Live chat participation
- See other viewers' messages

---

## API Endpoints

### Videos

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/videos` | List all videos |
| GET | `/api/v1/videos/{id}` | Get video details |
| POST | `/api/v1/videos/upload/presigned` | Get presigned upload URL |
| POST | `/api/v1/videos/upload/confirm` | Confirm upload, start transcode |
| PUT | `/api/v1/videos/{id}` | Update video metadata |
| DELETE | `/api/v1/videos/{id}` | Delete video |
| GET | `/api/v1/videos/{id}/stream` | Get HLS playback URL |

### Live Streams

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/streams` | List all streams |
| GET | `/api/v1/streams/{id}` | Get stream details |
| POST | `/api/v1/streams` | Create new stream channel |
| PUT | `/api/v1/streams/{id}` | Update stream settings |
| DELETE | `/api/v1/streams/{id}` | Delete stream channel |
| GET | `/api/v1/streams/{id}/ingest` | Get RTMP ingest URL |
| POST | `/api/v1/streams/{id}/start` | Start stream (go live) |
| POST | `/api/v1/streams/{id}/stop` | Stop stream |

### Chat (WebSocket)

| Action | Direction | Description |
|--------|-----------|-------------|
| `connect` | Client→Server | Join stream chat room |
| `message` | Client→Server | Send chat message |
| `message` | Server→Client | Broadcast message |
| `viewer_count` | Server→Client | Update viewer count |
| `stream_status` | Server→Client | Stream went live/offline |

---

## Data Models

### Video

```json
{
  "id": "uuid",
  "userId": "uuid",
  "title": "Music Video Title",
  "description": "Description text",
  "artist": "Artist Name",
  "duration": 240,
  "resolution": "1920x1080",
  "format": "MP4",
  "fileSize": 524288000,
  "thumbnailUrl": "https://cdn.../thumb.jpg",
  "hlsUrl": "https://cdn.../video.m3u8",
  "status": "READY",
  "trackId": "uuid",
  "albumId": "uuid",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

### LiveStream

```json
{
  "id": "uuid",
  "userId": "uuid",
  "title": "Friday Night DJ Set",
  "description": "Live performance",
  "thumbnailUrl": "https://cdn.../thumb.jpg",
  "status": "LIVE",
  "viewerCount": 42,
  "ivsChannelArn": "arn:aws:ivs:...",
  "ivsStreamKey": "sk_...",
  "rtmpIngestUrl": "rtmps://...",
  "playbackUrl": "https://...",
  "recordingEnabled": true,
  "scheduledAt": "2024-01-15T20:00:00Z",
  "startedAt": "2024-01-15T20:05:00Z",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

### ChatMessage

```json
{
  "id": "uuid",
  "streamId": "uuid",
  "userId": "uuid",
  "username": "DJ_Fan_42",
  "message": "Great track!",
  "timestamp": "2024-01-15T20:15:30Z"
}
```

---

## External Dependencies

### Video Processing

| Service | Purpose | Cost Model |
|---------|---------|------------|
| **AWS MediaConvert** | Video transcoding | Per minute of output |
| **S3** | Video storage | Storage + transfer |
| **CloudFront** | Video delivery | Per GB transferred |

### Live Streaming

| Service | Purpose | Cost Model |
|---------|---------|------------|
| **Amazon IVS** | Live streaming | Per hour input + viewing |
| **API Gateway WebSocket** | Real-time chat | Per message + connection |
| **DynamoDB** | Chat message storage | Read/write units |

### Phase 2 (GameLift)

| Service | Purpose | Cost Model |
|---------|---------|------------|
| **Amazon GameLift Streams** | Interactive streaming | Per stream hour |

---

## Out of Scope (Phase 1)

- GameLift Streams integration (Phase 2)
- Video editing tools
- Stream overlays/graphics
- Monetization features
- Multi-camera streams
- DVR/rewind during live
- Chat emoji/reactions
- Chat spam filtering
