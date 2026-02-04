# Requirements - Audio Analysis Enhancement

## Feature Overview

**Feature**: Audio Analysis Enhancement
**Priority**: Medium
**Dependencies**: Backend API (Epic 2), Frontend (Epic 4)

Enhance the music library with advanced audio analysis capabilities including BPM detection, musical key detection, and similar artist recommendations. These features will enrich track metadata and improve music discovery.

---

## Functional Requirements

### FR-1: BPM Detection

**FR-1.1**: System detects BPM for uploaded audio files
- BPM (Beats Per Minute) calculated during upload processing
- Support for MP3, FLAC, WAV, AAC, OGG formats
- BPM range: 20-300 BPM (typical music range)
- Store as integer value in track metadata

**FR-1.2**: BPM displayed in track metadata
- Show BPM on track detail page
- Include BPM in track list columns (optional toggle)
- Format: "120 BPM" display string

**FR-1.3**: Search and filter by BPM
- Filter tracks by BPM range
- Search index includes BPM for range queries
- BPM buckets: Slow (< 90), Medium (90-120), Fast (> 120)

**FR-1.4**: Backfill existing tracks
- On-demand analysis for tracks without BPM
- Batch processing capability via admin endpoint
- Progress tracking for bulk analysis

### FR-2: Musical Key Detection

**FR-2.1**: System detects musical key for uploaded audio
- Detect key (e.g., "C Major", "A Minor", "F# Major")
- Detect mode (Major or Minor)
- Store both key and mode separately

**FR-2.2**: Key displayed in track metadata
- Show key on track detail page
- Camelot notation option (e.g., "8A" for A Minor)
- Include key in track list columns (optional)

**FR-2.3**: Search and filter by key
- Filter tracks by key
- Filter by mode (major/minor)
- Harmonic mixing suggestions (Camelot wheel)

### FR-3: Similar Artists Discovery

**FR-3.1**: System provides similar artist recommendations
- Recommendations based on uploaded artist names
- External API integration (Last.fm, MusicBrainz, or Spotify)
- Cache recommendations to reduce API calls

**FR-3.2**: Similar artists displayed on artist page
- Show "Similar Artists" section on artist detail page
- Display up to 5 similar artists
- Format: "Similar to: Artist1, Artist2, Artist3"

**FR-3.3**: Link to search when clicked
- Clicking similar artist searches library for that artist
- If artist exists in library, navigate to artist page
- If not in library, show search results (may be empty)

---

## Non-Functional Requirements

### NFR-1: Performance

**NFR-1.1**: BPM detection completes within 30 seconds per track
**NFR-1.2**: Key detection completes within 30 seconds per track
**NFR-1.3**: Similar artist lookup < 500ms (with cache hit)
**NFR-1.4**: Combined analysis doesn't delay upload completion

### NFR-2: Accuracy

**NFR-2.1**: BPM accuracy within +/- 2 BPM for typical music
**NFR-2.2**: Key detection accuracy > 75% for clear tonal music
**NFR-2.3**: Graceful handling of atonal/ambient music (return null)

### NFR-3: Reliability

**NFR-3.1**: Analysis failures don't block upload completion
**NFR-3.2**: Retry mechanism for transient failures
**NFR-3.3**: Partial results saved (e.g., BPM found but not key)

### NFR-4: Cost

**NFR-4.1**: Minimize external API calls via caching
**NFR-4.2**: Similar artist cache TTL: 30 days
**NFR-4.3**: Lambda timeout appropriate for audio processing

---

## User Stories

### US-1: DJ/Producer BPM Matching
```
As a DJ, I want to see BPM for all my tracks
So that I can match tempos when mixing
```

**Acceptance Criteria**:
- BPM visible on track detail and list views
- Filter tracks by BPM range
- BPM matches expected value for known tracks

### US-2: Harmonic Mixing
```
As a DJ, I want to know the musical key of my tracks
So that I can create harmonic mixes
```

**Acceptance Criteria**:
- Key displayed in standard notation (C Major, Am)
- Optional Camelot notation (8A, 11B)
- Can filter tracks by key

### US-3: Music Discovery
```
As a user, I want to discover similar artists
So that I can find new music I might like
```

**Acceptance Criteria**:
- Artist page shows similar artists section
- Similar artists based on external recommendations
- Clicking artist searches my library

---

## API Endpoints

### Track Metadata (Updated)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tracks/{id}` | Returns track with BPM, key |
| GET | `/api/v1/tracks?bpmMin=&bpmMax=` | Filter by BPM range |
| GET | `/api/v1/tracks?key=C&mode=major` | Filter by key |

### Analysis Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tracks/{id}/analyze` | Trigger analysis for single track |
| POST | `/api/v1/analysis/batch` | Batch analyze tracks without BPM/key |
| GET | `/api/v1/analysis/status` | Get batch analysis progress |

### Similar Artists

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/artists/{name}/similar` | Get similar artists |

---

## Data Models

### Track (Updated Fields)

```json
{
  "bpm": 128,
  "musicalKey": "Am",
  "keyMode": "minor",
  "keyCamelot": "8A",
  "analysisStatus": "COMPLETED",
  "analyzedAt": "2024-01-15T10:30:00Z"
}
```

### Similar Artists Response

```json
{
  "artistName": "Daft Punk",
  "similarArtists": [
    { "name": "Justice", "inLibrary": true, "matchScore": 0.95 },
    { "name": "Kavinsky", "inLibrary": false, "matchScore": 0.89 },
    { "name": "Breakbot", "inLibrary": true, "matchScore": 0.85 }
  ],
  "source": "lastfm",
  "cachedAt": "2024-01-15T10:30:00Z"
}
```

---

## External Dependencies

### BPM Detection Options

| Option | Type | Pros | Cons |
|--------|------|------|------|
| **realtime-bpm-analyzer** | NPM (browser/Node) | Free, client-side possible | Accuracy varies |
| **essentia.js** | WASM | High accuracy, comprehensive | Larger bundle, complex |
| **AWS Transcribe** | Cloud | Managed, scalable | Cost, overkill for BPM |

**Recommendation**: Use realtime-bpm-analyzer in a Node.js Lambda for server-side processing.

### Key Detection Options

| Option | Type | Pros | Cons |
|--------|------|------|------|
| **essentia.js** | WASM | Accurate key detection | Complex setup |
| **keyfinder** | C++ (via WASM) | Industry standard | Needs compilation |
| **ACRCloud** | Cloud API | High accuracy | Cost per request |
| **AudD** | Cloud API | Includes key detection | Cost per request |

**Recommendation**: Use essentia.js WASM for key detection (free, accurate, runs in Lambda).

### Similar Artists Options

| Option | Type | Rate Limits | Cost |
|--------|------|-------------|------|
| **Last.fm API** | REST | 5 req/sec | Free |
| **MusicBrainz** | REST | 1 req/sec | Free |
| **Spotify Web API** | REST | Varies | Free (with auth) |

**Recommendation**: Use Last.fm API (best balance of simplicity, rate limits, and data quality).

---

## Out of Scope

- Beat grid / downbeat detection
- Audio waveform generation
- Mood/energy classification
- Genre auto-detection
- Lyrics transcription
- Copyright/fingerprint detection
