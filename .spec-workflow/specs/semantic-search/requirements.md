# Requirements - Semantic Search Enhancements

## Epic Overview

**Epic**: Semantic Search with Bedrock Embeddings + Nixiesearch
**Stream**: C - Search Enhancements
**Dependencies**: Epic 3 (Search & Streaming), Bedrock Gateway

Enhance search capabilities with semantic/vector search using Bedrock text embeddings stored alongside Nixiesearch. Both are OpenAI API compatible, enabling seamless integration.

---

## Functional Requirements

### FR-1: Track Embedding Generation

**FR-1.1**: Generate embeddings from track metadata
- Combine title, artist, album, genre, and tags into embedding text
- Use Bedrock Titan text embeddings (1024 dimensions)
- Generate embedding during track indexing
- Store embedding in search index alongside keyword data

**FR-1.2**: Embedding text composition
- Format: "{title} by {artist} from album {album}. Genre: {genre}. Tags: {tags}. BPM: {bpm}. Key: {keyCamelot}"
- Handle missing fields gracefully
- Normalize text (lowercase, trim whitespace)
- Maximum 8000 characters for Titan model limit

**FR-1.3**: Batch embedding generation
- Support bulk embedding for index rebuilds
- Process tracks in batches of 25 (Titan rate limits)
- Retry failed embeddings with exponential backoff

### FR-2: Semantic Search

**FR-2.1**: Hybrid search (keyword + semantic)
- Accept natural language queries
- Generate query embedding via Bedrock
- Combine BM25 (keyword) score with cosine similarity score
- Configurable weights: `keyword_weight + semantic_weight = 1.0`

**FR-2.2**: Search API enhancement
- Add `semantic: true` query parameter for semantic search
- Add `POST /api/v1/search/semantic` for advanced semantic queries
- Support pure semantic search or hybrid mode
- Default: hybrid with 0.7 keyword / 0.3 semantic weights

**FR-2.3**: Semantic query understanding
- Support natural language: "upbeat electronic tracks for a party"
- Support descriptive queries: "relaxing jazz for late night coding"
- Support similarity queries: "songs similar to Daft Punk style"

### FR-3: Similar Tracks

**FR-3.1**: "More like this" feature
- GET `/api/v1/tracks/:id/similar` endpoint
- Return top-N tracks with highest embedding similarity
- Exclude the source track from results
- Support filtering by user scope

**FR-3.2**: Similarity by audio features
- Find tracks with similar BPM (within ±5 BPM range)
- Find tracks in compatible Camelot keys (adjacent keys)
- Combine embedding similarity with feature matching

**FR-3.3**: DJ-specific similarity
- "Harmonic mixing" mode: prioritize key compatibility
- "Energy match" mode: prioritize BPM + genre similarity
- "Vibe match" mode: prioritize semantic similarity

### FR-4: Advanced Filtering

**FR-4.1**: Faceted search enhancement
- Add facets for BPM ranges: 60-90, 90-120, 120-150, 150+
- Add facets for Camelot keys: 1A-12A, 1B-12B
- Return facet counts in search response

**FR-4.2**: Filter combinations
- BPM range queries: `bpmMin=120&bpmMax=130`
- Key filtering: `keyCamelot=8A` or `keyCamelot=7A,8A,9A`
- Multiple filter AND logic

**FR-4.3**: Smart filters for DJs
- "Mixable with {trackId}" - find harmonically compatible tracks
- "Energy level: high/medium/low" - BPM-based energy classification

### FR-5: Smart Playlists (Future)

**FR-5.1**: Rule-based playlist generation
- Define rules: "BPM 120-130 AND genre=electronic"
- Auto-populate playlist from rule matches
- Update playlist when new tracks match

---

## Non-Functional Requirements

### NFR-1: Performance

**NFR-1.1**: Embedding generation < 500ms per track
**NFR-1.2**: Semantic search < 1s (including embedding + search)
**NFR-1.3**: Similar tracks query < 500ms
**NFR-1.4**: Batch embedding: 100 tracks/minute throughput

### NFR-2: Cost Optimization

**NFR-2.1**: Cache query embeddings for repeated queries
**NFR-2.2**: Only regenerate track embedding if metadata changes
**NFR-2.3**: Use Titan v2 (lowest cost Bedrock embedding model)

### NFR-3: Reliability

**NFR-3.1**: Graceful degradation if embedding fails
**NFR-3.2**: Keyword search always available as fallback
**NFR-3.3**: Track without embedding still searchable by keywords

---

## API Endpoints

### New Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/search?q=...&semantic=true` | Hybrid search with semantic |
| POST | `/api/v1/search/semantic` | Advanced semantic search |
| GET | `/api/v1/tracks/:id/similar` | Similar tracks |
| GET | `/api/v1/tracks/:id/mixable` | DJ-compatible tracks |

### Enhanced Endpoints

| Method | Path | Enhancement |
|--------|------|-------------|
| GET | `/api/v1/search` | Add `semantic` param, facets |
| POST | `/api/v1/search` | Add embedding-based ranking |

---

## Data Models

### TrackEmbedding

```go
type TrackEmbedding struct {
    TrackID     string     `json:"trackId"`
    UserID      string     `json:"userId"`
    Embedding   []float32  `json:"embedding"` // 1024 dimensions
    EmbedText   string     `json:"embedText"` // Text that was embedded
    ModelID     string     `json:"modelId"`   // amazon.titan-embed-text-v2:0
    GeneratedAt time.Time  `json:"generatedAt"`
}
```

### SemanticSearchRequest

```go
type SemanticSearchRequest struct {
    Query           string         `json:"query"`
    Mode            string         `json:"mode"`           // hybrid, semantic, keyword
    SemanticWeight  float64        `json:"semanticWeight"` // 0.0-1.0, default 0.3
    Filters         SearchFilters  `json:"filters,omitempty"`
    Limit           int            `json:"limit,omitempty"`
    Cursor          string         `json:"cursor,omitempty"`
}
```

### SimilarTracksResponse

```go
type SimilarTracksResponse struct {
    SourceTrack  TrackResponse   `json:"sourceTrack"`
    Similar      []SimilarTrack  `json:"similar"`
    TotalMatches int             `json:"totalMatches"`
}

type SimilarTrack struct {
    Track           TrackResponse `json:"track"`
    Similarity      float64       `json:"similarity"`      // 0.0-1.0
    BPMDiff         int           `json:"bpmDiff"`
    KeyCompatible   bool          `json:"keyCompatible"`
    MatchReasons    []string      `json:"matchReasons"`    // ["semantic", "bpm", "key"]
}
```

---

## Integration Points

### Existing Components

- **BedrockClient** (`backend/internal/clients/bedrock.go`): Add batch embedding support
- **SearchService** (`backend/internal/service/search.go`): Add semantic search methods
- **SearchClient** (`backend/internal/search/client.go`): Add vector search support
- **Search Types** (`backend/internal/search/types.go`): Add embedding fields

### New Components

- **EmbeddingService** (`backend/internal/service/embedding.go`): Track embedding generation
- **SimilarityService** (`backend/internal/service/similarity.go`): Similar tracks logic

---

## Out of Scope

- Real-time audio analysis for embedding (use metadata only)
- Cross-user recommendation (privacy boundary)
- Collaborative filtering (no user behavior tracking)
- GPU-accelerated embedding (use Bedrock API)
