# Requirements - Search & Streaming (Epic 3)

## Epic Overview

**Epic**: Search & Streaming
**Wave**: 3
**Dependencies**: Epic 2 (Backend API)

Implement full-text search using Nixiesearch with DynamoDB for metadata tracking, and HLS adaptive streaming with CloudFront signed URLs for audio playback.

---

## Functional Requirements

### FR-1: Full-Text Search

**FR-1.1**: Users can search their music library by text query
- Search across title, artist, album, genre fields
- Return ranked results by relevance
- Support partial matching and fuzzy search

**FR-1.2**: Search results include metadata from DynamoDB
- Track ID, title, artist, album, duration
- Cover art URL (if available)
- Play count and last played timestamp

**FR-1.3**: Search supports filtering
- Filter by artist
- Filter by album
- Filter by genre
- Filter by year range

**FR-1.4**: Search supports pagination
- Default page size: 20 results
- Maximum page size: 100 results
- Cursor-based pagination for stability

### FR-2: Search Indexing

**FR-2.1**: New tracks are indexed automatically
- Index updated when upload processing completes
- Index includes: title, artist, album, genre, year, filename

**FR-2.2**: Track updates trigger re-indexing
- Metadata edits update the search index
- Deletions remove from search index

**FR-2.3**: Scheduled full re-index
- Daily full re-index to catch any missed updates
- Index stored in S3 and loaded to Lambda /tmp on cold start
- Pure serverless - no VPC or EFS required

### FR-3: HLS Adaptive Streaming

**FR-3.1**: Audio files are transcoded to HLS format
- Generate HLS playlist (.m3u8) and segments (.ts)
- Multiple quality levels for adaptive bitrate
- Support MP3, FLAC, WAV source formats

**FR-3.2**: HLS quality levels
- Low: 96 kbps AAC
- Medium: 192 kbps AAC
- High: 320 kbps AAC
- Original: Preserve original quality (for FLAC)

**FR-3.3**: Transcoding triggered after upload
- Integrated into upload processing pipeline
- Store HLS output alongside original file

### FR-4: Streaming Playback

**FR-4.1**: Stream endpoint returns CloudFront signed URL
- URL points to HLS master playlist
- 24-hour expiration for convenience
- Scoped to user's own tracks only

**FR-4.2**: CloudFront serves HLS content
- Low latency edge delivery
- CORS headers for web player
- Cache HLS segments at edge

**FR-4.3**: Fallback to direct streaming
- If HLS not ready, serve original file directly
- Signed URL with same 24-hour expiration

### FR-5: Download

**FR-5.1**: Download endpoint returns signed URL
- URL points to original (non-transcoded) file
- 24-hour expiration
- Scoped to user's own tracks only

**FR-5.2**: Download tracking
- Increment download count in DynamoDB
- Track last download timestamp

---

## Non-Functional Requirements

### NFR-1: Performance

**NFR-1.1**: Search response time < 500ms (p95)
**NFR-1.2**: Stream URL generation < 200ms
**NFR-1.3**: Transcoding completes within 5 minutes for typical track

### NFR-2: Scalability

**NFR-2.1**: Search index supports up to 100,000 tracks per user
**NFR-2.2**: Concurrent streams limited by CloudFront (no artificial limit)
**NFR-2.3**: Transcoding parallelized for batch uploads

### NFR-3: Reliability

**NFR-3.1**: Search index persisted to S3 (durable storage, no EFS)
**NFR-3.2**: Transcoding failures don't block upload completion
**NFR-3.3**: Graceful degradation when HLS not available
**NFR-3.4**: Index loaded from S3 on Lambda cold start

### NFR-4: Security

**NFR-4.1**: CloudFront signed URLs prevent unauthorized access
**NFR-4.2**: Signed URL key pair managed via Secrets Manager
**NFR-4.3**: Users can only access their own tracks

### NFR-5: Cost

**NFR-5.1**: Use Lambda for search (pay-per-request vs always-on)
**NFR-5.2**: S3 Intelligent-Tiering for HLS segments
**NFR-5.3**: MediaConvert on-demand pricing for transcoding
**NFR-5.4**: No VPC - eliminates NAT Gateway (~$45/month) and EFS costs
**NFR-5.5**: Pure serverless - no idle costs when not in use

---

## User Stories

### US-1: Search
```
As a user, I want to search my music library
So that I can quickly find specific tracks to play
```

**Acceptance Criteria**:
- Search bar accepts text input
- Results appear within 1 second
- Can click result to play track

### US-2: Filter Search
```
As a user, I want to filter search results by artist or album
So that I can narrow down my search
```

**Acceptance Criteria**:
- Filter dropdowns available on search page
- Filters combine with text search
- Clear filters button resets all

### US-3: Stream Track
```
As a user, I want to stream a track with adaptive quality
So that playback is smooth on varying network conditions
```

**Acceptance Criteria**:
- Audio plays without manual quality selection
- Quality adjusts automatically based on bandwidth
- No buffering on stable connection

### US-4: Download Track
```
As a user, I want to download original quality files
So that I can listen offline or backup my library
```

**Acceptance Criteria**:
- Download button on track detail view
- Downloads original file format (MP3, FLAC, etc.)
- Progress shown for large files

---

## API Endpoints

### Search Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/search?q={query}` | Simple text search |
| POST | `/api/v1/search` | Advanced search with filters |

### Streaming Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/stream/{trackId}` | Get HLS streaming URL |
| GET | `/api/v1/download/{trackId}` | Get download URL |

### Index Management (Internal)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/index/rebuild` | Trigger full index rebuild |

---

## Data Models

### Search Document
```json
{
  "id": "track-uuid",
  "userId": "user-uuid",
  "title": "Song Title",
  "artist": "Artist Name",
  "album": "Album Name",
  "genre": "Rock",
  "year": 2023,
  "filename": "original-filename.mp3",
  "duration": 245,
  "indexedAt": "2024-01-15T10:30:00Z"
}
```

### HLS Manifest Reference (DynamoDB)
```json
{
  "PK": "USER#user-uuid",
  "SK": "TRACK#track-uuid",
  "hlsStatus": "READY|PROCESSING|FAILED",
  "hlsPlaylistKey": "hls/user-uuid/track-uuid/master.m3u8",
  "hlsCreatedAt": "2024-01-15T10:35:00Z"
}
```

---

## Dependencies

### AWS Services
- **Nixiesearch**: Full-text search engine (embedded in Lambda container)
- **MediaConvert**: Audio transcoding to HLS
- **CloudFront**: CDN for streaming and downloads
- **S3**: HLS segments, search index storage, and media files
- **Secrets Manager**: CloudFront signing key pair
- **Lambda**: All compute (no VPC, no EFS)

### Architecture
- **Pure Serverless**: No VPC, NAT Gateway, or EFS
- **S3 Index Storage**: Nixiesearch index loaded from S3 on cold start
- **Public Network**: All Lambdas run in AWS public network

### External Libraries
- **Nixiesearch**: Search engine binary (embedded in Lambda container image)
- **aws-sdk-go-v2/service/mediaconvert**: Transcoding API
- **aws-sdk-go-v2/service/cloudfront**: Signed URL generation
- **aws-sdk-go-v2/service/s3**: Index storage operations

---

## Out of Scope

- Real-time search suggestions (autocomplete)
- Lyrics search
- Audio fingerprinting / duplicate detection
- Shared/public playlists
- Social features (following, sharing)
