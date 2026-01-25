# Design - Semantic Search Enhancements

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Search Architecture                                │
│                                                                             │
│  ┌─────────────┐     ┌──────────────────┐     ┌────────────────────┐       │
│  │   Search    │────►│  Hybrid Search   │────►│    Nixiesearch     │       │
│  │   Request   │     │    Service       │     │  (BM25 + Vector)   │       │
│  └─────────────┘     └────────┬─────────┘     └────────────────────┘       │
│                               │                                             │
│                               ▼                                             │
│                      ┌────────────────┐                                     │
│                      │    Bedrock     │                                     │
│                      │    Titan       │                                     │
│                      │  (Embeddings)  │                                     │
│                      └────────────────┘                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                        Track Indexing Pipeline                               │
│                                                                             │
│  ┌─────────────┐     ┌──────────────────┐     ┌────────────────────┐       │
│  │   Track     │────►│   Embedding      │────►│   Index to         │       │
│  │   Metadata  │     │   Service        │     │   Nixiesearch      │       │
│  └─────────────┘     └────────┬─────────┘     └────────────────────┘       │
│                               │                                             │
│                               ▼                                             │
│                      ┌────────────────┐                                     │
│                      │    Bedrock     │                                     │
│                      │    Titan       │                                     │
│                      └────────────────┘                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Design

### 1. Embedding Service

**Purpose**: Generate and manage track metadata embeddings

**File**: `backend/internal/service/embedding.go`

```go
type EmbeddingService interface {
    // GenerateTrackEmbedding creates an embedding for a single track
    GenerateTrackEmbedding(ctx context.Context, track models.Track) (*TrackEmbedding, error)

    // GenerateQueryEmbedding creates an embedding for a search query
    GenerateQueryEmbedding(ctx context.Context, query string) ([]float32, error)

    // BatchGenerateEmbeddings creates embeddings for multiple tracks
    BatchGenerateEmbeddings(ctx context.Context, tracks []models.Track) ([]TrackEmbedding, error)

    // ComposeEmbedText creates the text to embed from track metadata
    ComposeEmbedText(track models.Track) string
}

type embeddingServiceImpl struct {
    bedrockClient *clients.BedrockClient
    modelID       string  // amazon.titan-embed-text-v2:0
    maxTextLen    int     // 8000 for Titan
}
```

**Embedding Text Composition**:
```go
func (s *embeddingServiceImpl) ComposeEmbedText(track models.Track) string {
    var parts []string

    // Core metadata
    parts = append(parts, track.Title)
    if track.Artist != "" {
        parts = append(parts, "by "+track.Artist)
    }
    if track.Album != "" {
        parts = append(parts, "from album "+track.Album)
    }

    // Genre and style
    if track.Genre != "" {
        parts = append(parts, "Genre: "+track.Genre)
    }

    // Tags for semantic context
    if len(track.Tags) > 0 {
        parts = append(parts, "Tags: "+strings.Join(track.Tags, ", "))
    }

    // DJ metadata for mixing context
    if track.BPM > 0 {
        parts = append(parts, fmt.Sprintf("BPM: %d", track.BPM))
    }
    if track.KeyCamelot != "" {
        parts = append(parts, "Key: "+track.KeyCamelot)
    }

    text := strings.Join(parts, ". ")
    if len(text) > s.maxTextLen {
        text = text[:s.maxTextLen]
    }
    return text
}
```

### 2. Enhanced Search Types

**File**: `backend/internal/search/types.go` (additions)

```go
// Document with embedding field
type Document struct {
    ID         string     `json:"id"`
    UserID     string     `json:"userId"`
    Title      string     `json:"title"`
    Artist     string     `json:"artist"`
    Album      string     `json:"album"`
    Genre      string     `json:"genre"`
    Year       int        `json:"year,omitempty"`
    Duration   int        `json:"duration,omitempty"`
    Filename   string     `json:"filename"`
    BPM        int        `json:"bpm,omitempty"`
    KeyCamelot string     `json:"keyCamelot,omitempty"`
    Tags       []string   `json:"tags,omitempty"`
    Embedding  []float32  `json:"embedding,omitempty"` // 1024-dim vector
    IndexedAt  time.Time  `json:"indexedAt"`
}

// SemanticSearchQuery for hybrid/semantic search
type SemanticSearchQuery struct {
    Query          string        `json:"query"`
    QueryEmbedding []float32     `json:"queryEmbedding,omitempty"`
    Mode           SearchMode    `json:"mode"` // hybrid, semantic, keyword
    KeywordWeight  float64       `json:"keywordWeight"`  // default 0.7
    SemanticWeight float64       `json:"semanticWeight"` // default 0.3
    Filters        SearchFilters `json:"filters,omitempty"`
    Limit          int           `json:"limit,omitempty"`
    Cursor         string        `json:"cursor,omitempty"`
}

type SearchMode string

const (
    SearchModeKeyword  SearchMode = "keyword"
    SearchModeSemantic SearchMode = "semantic"
    SearchModeHybrid   SearchMode = "hybrid"
)
```

### 3. Similarity Service

**Purpose**: Find similar tracks based on embedding similarity and audio features

**File**: `backend/internal/service/similarity.go`

```go
type SimilarityService interface {
    // FindSimilarTracks returns tracks similar to the given track
    FindSimilarTracks(ctx context.Context, userID, trackID string, opts SimilarityOptions) (*SimilarTracksResponse, error)

    // FindMixableTracks returns DJ-compatible tracks
    FindMixableTracks(ctx context.Context, userID, trackID string, opts MixingOptions) (*SimilarTracksResponse, error)

    // ComputeSimilarity calculates similarity between two tracks
    ComputeSimilarity(track1, track2 models.Track, embedding1, embedding2 []float32) float64
}

type SimilarityOptions struct {
    Limit            int      `json:"limit"`            // default 10
    Mode             string   `json:"mode"`             // semantic, features, combined
    MinSimilarity    float64  `json:"minSimilarity"`    // default 0.5
    IncludeSameAlbum bool     `json:"includeSameAlbum"` // default true
}

type MixingOptions struct {
    Limit        int    `json:"limit"`        // default 10
    BPMTolerance int    `json:"bpmTolerance"` // default 5
    KeyMode      string `json:"keyMode"`      // exact, harmonic, any
}
```

### 4. Camelot Key Compatibility

**Purpose**: Determine harmonic compatibility for DJ mixing

**File**: `backend/internal/service/camelot.go`

```go
// CamelotWheel represents the Camelot key compatibility system
var CamelotWheel = map[string][]string{
    // Each key maps to harmonically compatible keys
    // Same key + adjacent on wheel + relative major/minor
    "1A":  {"1A", "12A", "2A", "1B"},
    "2A":  {"2A", "1A", "3A", "2B"},
    "3A":  {"3A", "2A", "4A", "3B"},
    "4A":  {"4A", "3A", "5A", "4B"},
    "5A":  {"5A", "4A", "6A", "5B"},
    "6A":  {"6A", "5A", "7A", "6B"},
    "7A":  {"7A", "6A", "8A", "7B"},
    "8A":  {"8A", "7A", "9A", "8B"},
    "9A":  {"9A", "8A", "10A", "9B"},
    "10A": {"10A", "9A", "11A", "10B"},
    "11A": {"11A", "10A", "12A", "11B"},
    "12A": {"12A", "11A", "1A", "12B"},
    "1B":  {"1B", "12B", "2B", "1A"},
    "2B":  {"2B", "1B", "3B", "2A"},
    "3B":  {"3B", "2B", "4B", "3A"},
    "4B":  {"4B", "3B", "5B", "4A"},
    "5B":  {"5B", "4B", "6B", "5A"},
    "6B":  {"6B", "5B", "7B", "6A"},
    "7B":  {"7B", "6B", "8B", "7A"},
    "8B":  {"8B", "7B", "9B", "8A"},
    "9B":  {"9B", "8B", "10B", "9A"},
    "10B": {"10B", "9B", "11B", "10A"},
    "11B": {"11B", "10B", "12B", "11A"},
    "12B": {"12B", "11B", "1B", "12A"},
}

func IsKeyCompatible(key1, key2 string) bool {
    compatibleKeys, ok := CamelotWheel[key1]
    if !ok {
        return false
    }
    for _, k := range compatibleKeys {
        if k == key2 {
            return true
        }
    }
    return false
}
```

---

## API Design

### Semantic Search Endpoint

**POST /api/v1/search/semantic**

Request:
```json
{
  "query": "upbeat electronic tracks for a party",
  "mode": "hybrid",
  "semanticWeight": 0.4,
  "filters": {
    "bpmMin": 120,
    "bpmMax": 140
  },
  "limit": 20
}
```

Response:
```json
{
  "query": "upbeat electronic tracks for a party",
  "mode": "hybrid",
  "totalResults": 45,
  "results": [
    {
      "id": "track-uuid-1",
      "title": "Electric Dreams",
      "artist": "Synthwave Artist",
      "album": "Neon Nights",
      "bpm": 128,
      "keyCamelot": "8A",
      "score": 0.87,
      "keywordScore": 0.65,
      "semanticScore": 0.95
    }
  ],
  "facets": {
    "bpmRanges": [
      {"range": "120-130", "count": 23},
      {"range": "130-140", "count": 22}
    ],
    "keys": [
      {"key": "8A", "count": 12},
      {"key": "11B", "count": 8}
    ]
  }
}
```

### Similar Tracks Endpoint

**GET /api/v1/tracks/:id/similar?mode=combined&limit=10**

Response:
```json
{
  "sourceTrack": {
    "id": "track-uuid",
    "title": "Original Track",
    "artist": "Artist",
    "bpm": 125,
    "keyCamelot": "8A"
  },
  "similar": [
    {
      "track": {
        "id": "similar-uuid-1",
        "title": "Similar Track",
        "artist": "Other Artist",
        "bpm": 124,
        "keyCamelot": "8A"
      },
      "similarity": 0.89,
      "bpmDiff": 1,
      "keyCompatible": true,
      "matchReasons": ["semantic", "bpm", "key"]
    }
  ],
  "totalMatches": 45
}
```

### Mixable Tracks Endpoint

**GET /api/v1/tracks/:id/mixable?bpmTolerance=5&keyMode=harmonic**

Response:
```json
{
  "sourceTrack": {
    "id": "track-uuid",
    "title": "Now Playing",
    "bpm": 128,
    "keyCamelot": "8A"
  },
  "mixable": [
    {
      "track": {
        "id": "mixable-uuid",
        "title": "Great Mix",
        "bpm": 127,
        "keyCamelot": "7A"
      },
      "bpmDiff": 1,
      "keyTransition": "8A → 7A (energy down)",
      "mixScore": 0.95
    }
  ]
}
```

---

## Nixiesearch Index Schema Update

```yaml
name: music_tracks
fields:
  # Existing keyword/text fields
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

  # New fields for semantic search
  - name: bpm
    type: integer
  - name: keyCamelot
    type: keyword
  - name: tags
    type: keyword  # array of keywords

  # Vector field for semantic search
  - name: embedding
    type: dense_vector
    dimension: 1024
    similarity: cosine

settings:
  numberOfShards: 1
  numberOfReplicas: 0
```

---

## Error Handling

| Scenario | Handling | User Impact |
|----------|----------|-------------|
| Bedrock rate limit | Retry with backoff, queue requests | Slight delay |
| Embedding generation fails | Index without embedding, fallback to keyword | Keyword search only |
| Vector search unavailable | Fall back to keyword search | Degraded results |
| Track has no embedding | Exclude from semantic results | Some tracks hidden |
| Query too long | Truncate to 8000 chars | Partial query understanding |

---

## Testing Strategy

### Unit Tests
- Embedding text composition
- Camelot key compatibility
- Similarity score calculation
- Search mode selection

### Integration Tests
- Bedrock embedding generation
- Nixiesearch vector indexing
- Hybrid search scoring
- Similar tracks retrieval

### End-to-End Tests
- Full semantic search flow
- Similar tracks from UI
- DJ mixable tracks feature
