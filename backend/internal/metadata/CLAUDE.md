# Metadata Package - CLAUDE.md

## Overview

Audio metadata extraction package using dhowden/tag library. Extracts ID3 tags (MP3), Vorbis comments (FLAC, OGG), and M4A metadata.

## File Descriptions

| File | Purpose |
|------|---------|
| `extractor.go` | Main metadata extractor implementation |
| `extractor_test.go` | Unit tests with test fixtures |

## Key Types

### Extractor
The main type for extracting metadata from audio files.

```go
type Extractor struct{}

func NewExtractor() *Extractor
func (e *Extractor) Extract(reader io.ReadSeeker, filename string) (*models.UploadMetadata, error)
func (e *Extractor) ExtractCoverArt(reader io.ReadSeeker) ([]byte, string, error)
func (e *Extractor) DetectFormat(reader io.ReadSeeker) (models.AudioFormat, error)
```

## Functions

| Function | Description |
|----------|-------------|
| `NewExtractor()` | Creates a new metadata extractor instance |
| `Extract(reader, filename)` | Extracts all metadata from audio file, falls back to filename if no tags |
| `ExtractCoverArt(reader)` | Extracts embedded cover art image bytes and MIME type |
| `DetectFormat(reader)` | Detects audio format from file header |

## Supported Formats

| Format | Tags | Cover Art |
|--------|------|-----------|
| MP3 | ID3v1, ID3v2 | Yes |
| FLAC | Vorbis Comments | Yes |
| OGG | Vorbis Comments | Yes |
| M4A/AAC | iTunes Metadata | Yes |
| WAV | Limited | No |

## Dependencies

### Internal
- `github.com/gvasels/personal-music-searchengine/internal/models`

### External
- `github.com/dhowden/tag` - Pure Go audio tag reading library

## Usage Example

```go
extractor := metadata.NewExtractor()

// Open audio file
file, _ := os.Open("song.mp3")
defer file.Close()

// Extract metadata
meta, err := extractor.Extract(file, "song.mp3")
if err != nil {
    log.Printf("Metadata extraction failed: %v", err)
}

// Extract cover art
file.Seek(0, io.SeekStart)
coverData, mimeType, err := extractor.ExtractCoverArt(file)
if coverData != nil {
    // Save cover art to S3
}
```

## Fallback Behavior

When metadata cannot be read (corrupted tags, unsupported format, raw WAV):
- Title: Extracted from filename (without extension)
- Artist: "Unknown Artist"
- Format: Detected from file extension
- Other fields: Empty/zero values

## Testing

Test with fixtures in `testdata/`:
- `sample.mp3` - MP3 with full ID3v2 tags and cover art
- `sample-notags.mp3` - MP3 without any tags
- `sample.flac` - FLAC with Vorbis comments
- `sample.wav` - Raw WAV file
