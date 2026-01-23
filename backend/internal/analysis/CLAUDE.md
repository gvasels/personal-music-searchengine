# Analysis Package - CLAUDE.md

## Overview

Audio analysis package for BPM (tempo) and musical key detection. Used during upload processing to enrich track metadata with audio analysis data.

## File Descriptions

| File | Purpose |
|------|---------|
| `analyzer.go` | Main analyzer implementation with BPM and key detection |

## Key Types

### Result
Contains audio analysis results:
```go
type Result struct {
    BPM        int    // Beats per minute (0 if not detected)
    MusicalKey string // Musical key (e.g., "Am", "C", "F#m")
    KeyMode    string // "major" or "minor"
    KeyCamelot string // Camelot notation (e.g., "8A", "11B")
}
```

### Analyzer
Performs audio analysis on uploaded tracks.

## Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewAnalyzer` | `() *Analyzer` | Creates new analyzer instance |
| `Analyze` | `(ctx, reader, fileName) (*Result, error)` | Analyzes audio for BPM and key |
| `GetCamelotNotation` | `(key, mode string) string` | Converts key to Camelot notation |
| `detectFormat` | `(fileName string) string` | Detects audio format from filename |

## Camelot Wheel

The package includes a complete Camelot wheel mapping for harmonic mixing:
- Minor keys map to "A" column (1A-12A)
- Major keys map to "B" column (1B-12B)

Example mappings:
- Am → 8A
- C → 8B
- F#m → 11A
- D → 10B

## Implementation Status

**Current**: Placeholder implementation that returns empty results.

**Future enhancements**:
1. BPM detection using beat tracking algorithm
2. Key detection using chroma/pitch analysis
3. FFmpeg integration for audio decoding
4. Support for all audio formats (MP3, FLAC, WAV, AAC, OGG)

## Usage

```go
analyzer := analysis.NewAnalyzer()
result, err := analyzer.Analyze(ctx, audioReader, "track.mp3")
if err != nil {
    // Handle error
}
fmt.Printf("BPM: %d, Key: %s (%s)\n", result.BPM, result.MusicalKey, result.KeyCamelot)
```

## Dependencies

- Standard library only (for now)
- Future: ffmpeg binary or go-audio libraries

## Integration

Called by the `analyzer` Lambda processor during upload workflow:
```
ExtractMetadata → Parallel ─┬─► ProcessCoverArt ─┬─► CreateTrackRecord
                            └─► AnalyzeAudio ────┘
```

Analysis failures are non-blocking - upload continues even if analysis fails.
