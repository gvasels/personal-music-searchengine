# Waveform Processor

## Overview
Generates audio waveform data for visualization in the player UI. Uses FFmpeg for audio processing and peak extraction.

## Files

| File | Description |
|------|-------------|
| `waveform.go` | Waveform generator implementation using FFmpeg |
| `waveform_test.go` | Unit tests for waveform generation |

## Key Types

### WaveformData
```go
type WaveformData struct {
    Peaks      []float64 `json:"peaks"`      // Normalized 0.0-1.0 amplitude values
    SampleRate int       `json:"sampleRate"` // Samples per second (100)
    Duration   float64   `json:"duration"`   // Track duration in seconds
    Version    int       `json:"version"`    // Data format version
}
```

### Generator
Creates waveform data from audio files or bytes.

## Key Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewGenerator` | `() *Generator` | Creates new waveform generator |
| `Generate` | `(ctx, audioPath) (*WaveformData, error)` | Generate from file path |
| `GenerateFromBytes` | `(ctx, data, format) (*WaveformData, error)` | Generate from audio bytes |
| `Validate` | `() bool` | Check if waveform data is valid |

## Dependencies

- **FFmpeg**: Required for audio decoding and peak extraction
- **FFprobe**: Required for duration detection

## Supported Formats

- MP3, FLAC, WAV, AAC, OGG, M4A

## Usage Example

```go
gen := waveform.NewGenerator()
data, err := gen.Generate(ctx, "/path/to/audio.mp3")
if err != nil {
    return err
}
jsonBytes, _ := data.ToJSON()
// Store jsonBytes in S3
```

## Integration Points

- Called by upload processor after transcoding
- Waveform JSON stored in S3: `waveforms/{trackId}.json`
- URL stored in Track.WaveformURL field
