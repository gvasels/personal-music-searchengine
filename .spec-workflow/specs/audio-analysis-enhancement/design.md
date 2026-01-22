# Design - Audio Analysis Enhancement

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Upload Processing Pipeline                           │
│                                                                              │
│  ExtractMetadata → CoverArt → CreateTrack → MoveFile → StartTranscode       │
│                                    │                           │             │
│                                    │                           ▼             │
│                                    │                    ┌─────────────┐     │
│                                    │                    │ StartAnalysis│     │
│                                    │                    │   Lambda    │     │
│                                    │                    └──────┬──────┘     │
│                                    │                           │             │
│                                    ▼                           │ Async       │
│                             IndexForSearch ◄───────────────────┘             │
│                                    │                                         │
│                                    ▼                                         │
│                           MarkUploadCompleted                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                      Async Audio Analysis Pipeline                           │
│                                                                              │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                   │
│  │ BPM Lambda  │     │ Key Lambda  │     │ Update Track│                   │
│  │ (Node.js)   │────►│ (Node.js)   │────►│   Lambda    │                   │
│  └─────────────┘     └─────────────┘     └─────────────┘                   │
│         │                   │                   │                            │
│         └───────────────────┴───────────────────┘                            │
│                             │                                                │
│                             ▼                                                │
│                    ┌─────────────────┐                                      │
│                    │  Update Search  │                                      │
│                    │     Index       │                                      │
│                    └─────────────────┘                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                     Similar Artists Service                                  │
│                                                                              │
│  Artist Page ────► API Lambda ────► Similar Artists Service                 │
│                         │                    │                               │
│                         │                    ▼                               │
│                         │           ┌─────────────────┐                     │
│                         │           │  DynamoDB Cache │                     │
│                         │           └────────┬────────┘                     │
│                         │                    │ Cache Miss                    │
│                         │                    ▼                               │
│                         │           ┌─────────────────┐                     │
│                         │           │  Last.fm API    │                     │
│                         │           └─────────────────┘                     │
│                         ▼                                                    │
│                   Return Response                                           │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Design Decisions

### DD-1: Analysis Execution Location

**Decision**: Server-side analysis in Node.js Lambda (not browser)

**Rationale**:
- Consistent results regardless of client
- No client-side compute burden
- Works for batch processing
- Supports files uploaded via any method

**Trade-offs**:
- Lambda cold starts for analysis
- Memory requirements for audio processing
- Cannot leverage client CPU for distributed processing

### DD-2: BPM Detection Library

**Decision**: Use `realtime-bpm-analyzer` npm package

**Rationale**:
- Pure JavaScript, runs in Node.js Lambda
- No native dependencies or WASM compilation
- Sufficient accuracy for music library use case
- Small bundle size (~50KB)

**Alternative Considered**: essentia.js
- Higher accuracy but much larger (~20MB WASM)
- Complex setup for Lambda
- Overkill for BPM-only use case

**Implementation**:
```javascript
import { analyzeFullBuffer } from 'realtime-bpm-analyzer';

async function detectBPM(audioBuffer) {
  const { bpm } = await analyzeFullBuffer(audioBuffer);
  return Math.round(bpm);
}
```

### DD-3: Key Detection Library

**Decision**: Use `essentia.js` WASM module

**Rationale**:
- Industry-standard algorithm (KeyExtractor)
- High accuracy for tonal music
- Includes Camelot wheel mapping
- Free and open source

**Implementation**:
```javascript
import { EssentiaWASM } from 'essentia.js';

const essentia = new EssentiaWASM();
const keyData = essentia.KeyExtractor(audioVector);
// Returns: { key: "A", scale: "minor", strength: 0.85 }
```

**Lambda Configuration**:
- Runtime: Node.js 20.x
- Memory: 1024MB (for WASM)
- Timeout: 60 seconds
- Ephemeral storage: 512MB

### DD-4: Similar Artists API

**Decision**: Use Last.fm API

**Rationale**:
- Free API with generous rate limits (5 req/sec)
- Excellent music metadata database
- Simple REST API
- No OAuth required (API key only)

**Alternative Considered**: Spotify API
- Requires OAuth flow and user authentication
- Better data but more complex setup
- Rate limits tied to user tokens

**API Endpoint**:
```
GET https://ws.audioscrobbler.com/2.0/
  ?method=artist.getsimilar
  &artist={artist}
  &api_key={API_KEY}
  &format=json
  &limit=5
```

### DD-5: Analysis Timing

**Decision**: Async analysis AFTER upload completion

**Rationale**:
- Upload confirmation not delayed by analysis
- User can play track immediately
- Analysis runs in background
- Retry doesn't affect upload status

**Workflow**:
1. Upload completes (status: COMPLETED)
2. Analysis triggered asynchronously
3. Track shows "Analyzing..." status
4. Track updated when analysis completes

### DD-6: Caching Strategy

**Decision**: DynamoDB cache for similar artists with 30-day TTL

**Rationale**:
- Similar artists rarely change
- Reduces Last.fm API calls
- Low storage cost
- Single-table design compatible

**Cache Schema**:
```
PK: SIMILAR#NormalizedArtistName
SK: CACHE
TTL: Unix timestamp (30 days)
```

### DD-7: Handling Analysis Failures

**Decision**: Graceful degradation with null values

**Rationale**:
- Some tracks have no detectable BPM (ambient, classical)
- Key detection may fail for atonal music
- Partial success is acceptable

**Status Values**:
- `PENDING`: Awaiting analysis
- `ANALYZING`: In progress
- `COMPLETED`: Successful (may have null BPM/key)
- `FAILED`: Unrecoverable error

---

## Component Design

### 1. BPM Analyzer Lambda

**File**: `backend/cmd/processor/bpm-analyzer/main.js` (Node.js)

**Event Input**:
```json
{
  "trackId": "uuid",
  "userId": "uuid",
  "s3Key": "media/userId/trackId.mp3",
  "format": "MP3"
}
```

**Event Output**:
```json
{
  "trackId": "uuid",
  "bpm": 128,
  "confidence": 0.92
}
```

### 2. Key Detector Lambda

**File**: `backend/cmd/processor/key-detector/main.js` (Node.js)

**Event Input**: Same as BPM Analyzer

**Event Output**:
```json
{
  "trackId": "uuid",
  "key": "A",
  "mode": "minor",
  "camelot": "8A",
  "confidence": 0.85
}
```

### 3. Analysis Orchestrator

**File**: `backend/cmd/processor/analysis-orchestrator/main.go`

**Purpose**: Triggers BPM and Key analysis in parallel, updates track

**Step Functions State**:
```json
{
  "StartAnalysis": {
    "Type": "Parallel",
    "Branches": [
      { "StartAt": "DetectBPM" },
      { "StartAt": "DetectKey" }
    ],
    "Next": "UpdateTrackAnalysis",
    "Catch": [{ "ErrorEquals": ["States.ALL"], "Next": "MarkAnalysisFailed" }]
  }
}
```

### 4. Similar Artists Service

**File**: `backend/internal/service/similar.go`

**Functions**:
```go
type SimilarArtistsService struct {
    repo      *repository.Repository
    lastfmKey string
    httpClient *http.Client
}

func (s *SimilarArtistsService) GetSimilarArtists(ctx context.Context, artistName string) (*SimilarArtistsResponse, error)
func (s *SimilarArtistsService) checkLibrary(ctx context.Context, userID string, artists []string) map[string]bool
func (s *SimilarArtistsService) fetchFromLastFM(artistName string) ([]LastFMSimilarArtist, error)
func (s *SimilarArtistsService) getCachedSimilar(artistName string) (*CachedSimilarArtists, error)
func (s *SimilarArtistsService) cacheSimilar(artistName string, artists []SimilarArtist) error
```

### 5. Similar Artists Handler

**File**: `backend/internal/handlers/artists.go` (update)

**New Endpoint**:
```go
func (h *Handlers) GetSimilarArtists(c echo.Context) error {
    artistName := c.Param("name")
    userID := getUserID(c)

    similar, err := h.similarService.GetSimilarArtists(c.Request().Context(), artistName)
    if err != nil {
        return handleError(c, err)
    }

    // Check which artists exist in user's library
    enriched := h.similarService.checkLibrary(c.Request().Context(), userID, similar.Artists)

    return c.JSON(http.StatusOK, enriched)
}
```

---

## DynamoDB Schema Updates

### Track Entity (Updated Fields)

| Attribute | Type | Description |
|-----------|------|-------------|
| `bpm` | Number | Beats per minute (20-300) |
| `musicalKey` | String | Key signature (C, D, E, F, G, A, B) |
| `keyMode` | String | major or minor |
| `keyCamelot` | String | Camelot notation (1A-12B) |
| `analysisStatus` | String | PENDING, ANALYZING, COMPLETED, FAILED |
| `analyzedAt` | String | ISO8601 timestamp |

### Similar Artists Cache Entity

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `SIMILAR#{normalizedName}` |
| `SK` | String | `CACHE` |
| `Type` | String | `SIMILAR_CACHE` |
| `artistName` | String | Original artist name |
| `similarArtists` | List | List of similar artist objects |
| `source` | String | `lastfm` |
| `cachedAt` | String | ISO8601 timestamp |
| `TTL` | Number | Unix timestamp for expiry |

---

## Frontend Updates

### Track Type Update

**File**: `frontend/src/types/index.ts`

```typescript
export interface Track {
  // ... existing fields ...
  bpm?: number;
  musicalKey?: string;
  keyMode?: 'major' | 'minor';
  keyCamelot?: string;
  analysisStatus?: 'PENDING' | 'ANALYZING' | 'COMPLETED' | 'FAILED';
}
```

### Artist Type Update

```typescript
export interface SimilarArtist {
  name: string;
  inLibrary: boolean;
  matchScore: number;
}

export interface ArtistWithDetails extends Artist {
  albums: Album[];
  recentTracks: Track[];
  similarArtists?: SimilarArtist[];
}
```

### Track Detail Component Update

**File**: `frontend/src/routes/tracks/$trackId.tsx`

Add BPM and Key display:
```tsx
{track.bpm && (
  <div className="stat">
    <div className="stat-title">BPM</div>
    <div className="stat-value text-lg">{track.bpm}</div>
  </div>
)}
{track.musicalKey && (
  <div className="stat">
    <div className="stat-title">Key</div>
    <div className="stat-value text-lg">
      {track.musicalKey} {track.keyMode}
      {track.keyCamelot && <span className="text-sm ml-1">({track.keyCamelot})</span>}
    </div>
  </div>
)}
```

### Artist Detail Component Update

**File**: `frontend/src/routes/artists/$artistName.tsx`

Add Similar Artists section:
```tsx
{artist.similarArtists && artist.similarArtists.length > 0 && (
  <section className="mb-8">
    <h2 className="text-xl font-bold mb-4">Similar Artists</h2>
    <div className="flex flex-wrap gap-2">
      {artist.similarArtists.map((similar) => (
        <Link
          key={similar.name}
          to={similar.inLibrary ? `/artists/${similar.name}` : `/search?q=${similar.name}`}
          className={`badge badge-lg ${similar.inLibrary ? 'badge-primary' : 'badge-ghost'}`}
        >
          {similar.name}
          {similar.inLibrary && <CheckIcon className="w-4 h-4 ml-1" />}
        </Link>
      ))}
    </div>
  </section>
)}
```

---

## Infrastructure Updates

### New Lambda Functions

| Lambda | Runtime | Memory | Timeout | Purpose |
|--------|---------|--------|---------|---------|
| `bpm-analyzer` | Node.js 20.x | 512MB | 60s | BPM detection |
| `key-detector` | Node.js 20.x | 1024MB | 60s | Key detection |
| `analysis-orchestrator` | Go | 256MB | 10s | Triggers analysis |

### Step Functions Update

Add parallel analysis after existing workflow completes:

```hcl
MarkUploadCompleted = {
  Type = "Task"
  ...
  Next = "StartAnalysisAsync"
}

StartAnalysisAsync = {
  Type = "Task"
  Resource = aws_lambda_function.analysis_orchestrator.arn
  End = true
}
```

### Secrets Manager

| Secret | Purpose |
|--------|---------|
| `music-library/lastfm-api-key` | Last.fm API key |

---

## Error Handling

### Analysis Errors

| Error | Handling | User Impact |
|-------|----------|-------------|
| Audio decode failure | Mark as FAILED, log | BPM/Key shown as "—" |
| BPM not detectable | Set bpm=null, status=COMPLETED | "No rhythm detected" |
| Key not detectable | Set key=null, status=COMPLETED | "Key unknown" |
| Lambda timeout | Retry once, then FAILED | May retry later |
| Last.fm rate limit | Return cached or empty | Fewer recommendations |
| Last.fm unavailable | Return cached or empty | "Unavailable" message |

---

## Monitoring

### CloudWatch Metrics

- `AnalysisLatency` - Time to analyze track
- `BPMDetectionRate` - % of tracks with detected BPM
- `KeyDetectionRate` - % of tracks with detected key
- `SimilarArtistsCacheHitRate` - Cache effectiveness
- `LastFMAPIErrors` - External API failures

### CloudWatch Alarms

- Analysis latency > 45s (p95)
- BPM detection rate < 80%
- Last.fm error rate > 10%
