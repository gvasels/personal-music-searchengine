# Design - Search & Streaming (Epic 3)

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CloudFront CDN                                  │
│     ┌──────────────────┐              ┌──────────────────┐                  │
│     │  HLS Streaming   │              │  Original Files  │                  │
│     │  (Signed URLs)   │              │  (Downloads)     │                  │
│     └────────┬─────────┘              └────────┬─────────┘                  │
└──────────────┼─────────────────────────────────┼────────────────────────────┘
               │                                 │
               ▼                                 ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                              S3 Media Bucket                                 │
│   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│   │   media/        │  │   hls/          │  │   uploads/      │             │
│   │   (originals)   │  │   (transcoded)  │  │   (pending)     │             │
│   └─────────────────┘  └─────────────────┘  └─────────────────┘             │
└──────────────────────────────────────────────────────────────────────────────┘
               │                    ▲
               │                    │
               ▼                    │
┌──────────────────────────┐   ┌────────────────────────┐
│   MediaConvert           │   │   S3 Search Index      │
│   (HLS Transcoding)      │   │   (Nixiesearch)        │
└──────────────────────────┘   └────────────────────────┘
                                        ▲
                                        │
┌───────────────────────────────────────┼──────────────────────────────────────┐
│                           API Gateway                                        │
│                                       │                                      │
│   /search ──► Search Lambda ──────────┘                                      │
│   /stream ──► API Lambda ──► CloudFront Signer                               │
│   /download ─► API Lambda ──► CloudFront Signer                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Design

### 1. Nixiesearch Integration

**Deployment Model**: Embedded Lambda with S3 index storage (pure serverless, no VPC)

```
┌─────────────────┐     ┌─────────────────────────────────────┐
│  API Lambda     │────►│  Nixiesearch Lambda (embedded)      │
│  (search req)   │     │  • Loads index from S3 on cold start│
│                 │     │  • Writes index back to S3 on update│
└─────────────────┘     └───────────────────┬─────────────────┘
                                            │
                                            ▼
                                 ┌─────────────────────┐
                                 │  S3 Index Bucket    │
                                 │  (lucene index dir) │
                                 └─────────────────────┘
```

**Index Schema**:
```yaml
name: music_tracks
fields:
  - name: id
    type: keyword
  - name: userId
    type: keyword
  - name: title
    type: text
    analyzer: standard
  - name: artist
    type: text
    analyzer: standard
  - name: album
    type: text
    analyzer: standard
  - name: genre
    type: keyword
  - name: year
    type: integer
  - name: filename
    type: text
    analyzer: standard
settings:
  numberOfShards: 1
  numberOfReplicas: 0
```

### 2. HLS Transcoding Pipeline

**Workflow Integration**:

```
Upload Complete
      │
      ▼
┌─────────────────┐
│ Start Transcode │ (new Step Function state)
│ Lambda          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ MediaConvert    │
│ Job Created     │
└────────┬────────┘
         │
         ▼ (async - EventBridge)
┌─────────────────┐
│ Transcode       │
│ Complete Lambda │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Update Track    │
│ HLS Status      │
└─────────────────┘
```

**MediaConvert Job Settings**:
```json
{
  "OutputGroups": [
    {
      "Name": "HLS",
      "OutputGroupSettings": {
        "Type": "HLS_GROUP_SETTINGS",
        "HlsGroupSettings": {
          "SegmentLength": 6,
          "MinSegmentLength": 0,
          "Destination": "s3://bucket/hls/{userId}/{trackId}/"
        }
      },
      "Outputs": [
        {
          "NameModifier": "_96k",
          "AudioDescriptions": [{
            "CodecSettings": {
              "Codec": "AAC",
              "AacSettings": { "Bitrate": 96000, "SampleRate": 44100 }
            }
          }]
        },
        {
          "NameModifier": "_192k",
          "AudioDescriptions": [{
            "CodecSettings": {
              "Codec": "AAC",
              "AacSettings": { "Bitrate": 192000, "SampleRate": 44100 }
            }
          }]
        },
        {
          "NameModifier": "_320k",
          "AudioDescriptions": [{
            "CodecSettings": {
              "Codec": "AAC",
              "AacSettings": { "Bitrate": 320000, "SampleRate": 48000 }
            }
          }]
        }
      ]
    }
  ]
}
```

### 3. CloudFront Signed URLs

**Signing Configuration**:
```go
type CloudFrontSigner struct {
    keyPairID  string
    privateKey *rsa.PrivateKey
    domain     string
}

func (s *CloudFrontSigner) SignURL(resource string, expiry time.Duration) (string, error) {
    policy := &sign.Policy{
        Statements: []sign.Statement{{
            Resource: fmt.Sprintf("https://%s/%s*", s.domain, resource),
            Condition: sign.Condition{
                DateLessThan: &sign.AWSEpochTime{Time: time.Now().Add(expiry)},
            },
        }},
    }
    return sign.NewURLSigner(s.keyPairID, s.privateKey).SignWithPolicy(resource, policy)
}
```

**URL Structure**:
- Stream: `https://cdn.example.com/hls/{userId}/{trackId}/master.m3u8?Policy=...&Signature=...&Key-Pair-Id=...`
- Download: `https://cdn.example.com/media/{userId}/{trackId}.mp3?Policy=...&Signature=...&Key-Pair-Id=...`

### 4. CloudFront Distribution

**Behaviors**:
| Path Pattern | Origin | Cache Policy | Signed |
|--------------|--------|--------------|--------|
| `/hls/*` | S3 Media | CachingOptimized | Yes |
| `/media/*` | S3 Media | CachingDisabled | Yes |
| `/covers/*` | S3 Media | CachingOptimized | No |

**Origin Access Control**:
- S3 bucket policy restricts access to CloudFront only
- No direct S3 URL access possible

---

## Design Decisions

### DD-1: Nixiesearch vs OpenSearch Serverless

**Decision**: Use Nixiesearch embedded in Lambda with S3 index storage

**Rationale**:
- Lower cost for small-medium libraries (< 100k tracks)
- No minimum hourly charge (unlike OpenSearch Serverless)
- Pure serverless - no VPC, NAT Gateway, or EFS costs
- Sufficient features for music metadata search
- S3 provides durable, cheap storage for index files

**Trade-offs**:
- Cold start latency when loading index from S3
- Index size limited by Lambda ephemeral storage (10GB)
- Eventual consistency for index updates

### DD-2: S3 for Search Index Storage

**Decision**: S3 only for index storage (no EFS)

**Rationale**:
- No VPC required - simpler architecture
- No NAT Gateway costs (~$45/month saved)
- No EFS costs (~$0.30/GB-month saved)
- S3 provides 11 9's durability
- Index loaded to /tmp on cold start, written back on updates

**Trade-offs**:
- Higher cold start latency (index load from S3)
- Must manage index serialization/deserialization
- Index updates require S3 write operations

### DD-3: HLS vs Direct Streaming

**Decision**: HLS adaptive streaming with direct fallback

**Rationale**:
- HLS provides adaptive bitrate for varying network conditions
- Better mobile experience with quality switching
- Industry standard supported by all browsers
- Direct fallback ensures immediate playback for new uploads

**Trade-offs**:
- MediaConvert cost per minute transcoded
- Storage increase (~4x for multiple quality levels)
- Initial delay before HLS is available

### DD-4: MediaConvert vs FFmpeg Lambda

**Decision**: AWS MediaConvert

**Rationale**:
- Managed service with automatic scaling
- No Lambda timeout concerns for long files
- Built-in HLS packaging
- Cost-effective for on-demand transcoding

**Trade-offs**:
- Async processing (not immediate)
- Minimum job charge
- Less control than custom FFmpeg

### DD-5: Transcode Timing

**Decision**: Async transcoding after upload completion

**Rationale**:
- Don't block upload confirmation on transcoding
- User can play original while HLS generates
- Parallel processing for batch uploads

**Trade-offs**:
- Initial plays use direct streaming
- Need to track HLS status per track

### DD-6: Signed URL Expiration

**Decision**: 24-hour expiration

**Rationale**:
- Convenient for long listening sessions
- Supports offline-capable PWA patterns
- Reduces API calls for URL refresh

**Trade-offs**:
- Longer window if URL is compromised
- May need refresh for multi-day sessions

### DD-7: Search Index Updates

**Decision**: Real-time + daily full rebuild

**Rationale**:
- Real-time updates for immediate searchability
- Daily rebuild catches any missed updates
- Full rebuild handles schema changes

**Trade-offs**:
- Dual update path complexity
- Daily rebuild resource cost

### DD-8: Quality Levels

**Decision**: 3 AAC quality levels (96k, 192k, 320k)

**Rationale**:
- 96k: Mobile data saver
- 192k: Standard quality (comparable to Spotify)
- 320k: High quality for good connections
- AAC provides better quality than MP3 at same bitrate

**Trade-offs**:
- No lossless option in HLS (use download for FLAC)
- Storage multiplier for multiple qualities

---

## DynamoDB Schema Updates

### Track Entity (Updated)

| Attribute | Type | Description |
|-----------|------|-------------|
| `hlsStatus` | String | PENDING, PROCESSING, READY, FAILED |
| `hlsPlaylistKey` | String | S3 key for master.m3u8 |
| `hlsCreatedAt` | String | ISO8601 timestamp |
| `hlsError` | String | Error message if failed |
| `searchIndexedAt` | String | Last search index update |

### New GSI for HLS Status

| GSI | PK | SK | Purpose |
|-----|----|----|---------|
| GSI2 | `HLS#STATUS#{status}` | `{createdAt}` | Find tracks by HLS status |

---

## API Design

### Search Request/Response

**GET /api/v1/search?q={query}&artist={artist}&limit={limit}&cursor={cursor}**

Response:
```json
{
  "results": [
    {
      "id": "track-uuid",
      "title": "Song Title",
      "artist": "Artist Name",
      "album": "Album Name",
      "duration": 245,
      "coverArtUrl": "https://cdn.../covers/...",
      "score": 0.95
    }
  ],
  "total": 42,
  "cursor": "eyJsYXN0SWQiOiAiYWJjMTIzIn0="
}
```

**POST /api/v1/search**

Request:
```json
{
  "query": "search text",
  "filters": {
    "artist": "Artist Name",
    "album": "Album Name",
    "genre": "Rock",
    "yearFrom": 2020,
    "yearTo": 2024
  },
  "sort": {
    "field": "title",
    "order": "asc"
  },
  "limit": 20,
  "cursor": "eyJsYXN0SWQiOiAiYWJjMTIzIn0="
}
```

### Stream Response

**GET /api/v1/stream/{trackId}**

Response:
```json
{
  "trackId": "track-uuid",
  "streamUrl": "https://cdn.../hls/.../master.m3u8?Policy=...",
  "format": "hls",
  "expiresAt": "2024-01-16T10:30:00Z",
  "fallbackUrl": "https://cdn.../media/.../track.mp3?Policy=..."
}
```

### Download Response

**GET /api/v1/download/{trackId}**

Response:
```json
{
  "trackId": "track-uuid",
  "downloadUrl": "https://cdn.../media/.../track.flac?Policy=...",
  "filename": "Artist - Title.flac",
  "size": 52428800,
  "format": "flac",
  "expiresAt": "2024-01-16T10:30:00Z"
}
```

---

## Infrastructure Components

### New Terraform Resources

| Resource | Purpose |
|----------|---------|
| `aws_cloudfront_distribution.media` | CDN for streaming/downloads |
| `aws_cloudfront_public_key` | Signed URL verification |
| `aws_cloudfront_key_group` | Key group for signing |
| `aws_cloudfront_origin_access_control` | S3 origin security |
| `aws_secretsmanager_secret.cf_signing_key` | Private key storage |
| `aws_s3_bucket.search_index` | S3 bucket for Nixiesearch index storage |
| `aws_lambda_function.nixiesearch` | Search engine Lambda (embedded, public internet) |
| `aws_lambda_function.transcode_start` | Start MediaConvert job |
| `aws_lambda_function.transcode_complete` | Handle job completion |
| `aws_iam_role.mediaconvert` | MediaConvert job role |
| `aws_media_convert_queue.default` | Transcoding job queue |

### Network Architecture

**Pure Public Serverless** - No VPC required:
- All Lambdas run in AWS public network
- Nixiesearch Lambda loads index from S3 on cold start
- No NAT Gateway, no EFS, no VPC endpoints
- Reduces operational complexity and costs

---

## Error Handling

### Search Errors
| Error | HTTP Code | Response |
|-------|-----------|----------|
| Invalid query | 400 | `{"error": "query must be at least 2 characters"}` |
| Index unavailable | 503 | `{"error": "search temporarily unavailable"}` |

### Stream Errors
| Error | HTTP Code | Response |
|-------|-----------|----------|
| Track not found | 404 | `{"error": "track not found"}` |
| Not owner | 403 | `{"error": "access denied"}` |
| HLS not ready | 200 | Returns fallback URL only |

### Transcode Errors
| Error | Handling |
|-------|----------|
| MediaConvert failure | Set hlsStatus=FAILED, log error |
| Unsupported format | Skip transcode, use direct streaming |
| Timeout | Retry with exponential backoff |

---

## Monitoring

### CloudWatch Metrics
- Search latency (p50, p95, p99)
- Search error rate
- Transcode job success/failure rate
- CloudFront cache hit ratio
- S3 index bucket operations

### CloudWatch Alarms
- Search latency > 1s (p95)
- Transcode failure rate > 5%
- Lambda cold start duration > 5s

### X-Ray Tracing
- End-to-end search request tracing
- Transcode job tracking
