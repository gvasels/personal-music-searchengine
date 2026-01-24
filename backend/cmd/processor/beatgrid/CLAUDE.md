# Beat Grid Processor

## Overview
Calculates beat timing information from BPM for DJ features like beat-synced visualization and mixing assistance.

## Files

| File | Description |
|------|-------------|
| `beatgrid.go` | Beat grid calculator implementation |
| `beatgrid_test.go` | Unit tests for beat grid calculation |

## Key Types

### BeatGrid
```go
type BeatGrid struct {
    BPM        int     `json:"bpm"`        // Beats per minute
    Offset     float64 `json:"offset"`     // First beat offset in ms
    Beats      []int64 `json:"beats"`      // Beat timestamps in ms
    Downbeats  []int   `json:"downbeats"`  // Indices of downbeats (every 4th)
    IsVariable bool    `json:"isVariable"` // True if BPM varies
}
```

### Calculator
Generates beat grids from BPM information.

## Key Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewCalculator` | `() *Calculator` | Creates new beat grid calculator |
| `Calculate` | `(bpm, duration, offset) *BeatGrid` | Calculate beat grid |
| `MsPerBeat` | `(bpm int) float64` | Get milliseconds per beat |
| `Validate` | `() bool` | Check if beat grid is valid |
| `GetBeatAtTime` | `(timeMs) int` | Find nearest beat to time |
| `GetTimeAtBeat` | `(index) int64` | Get timestamp at beat index |
| `IsDownbeat` | `(index) bool` | Check if beat is a downbeat |

## BPM Range

- Minimum: 20 BPM
- Maximum: 300 BPM
- Returns nil for invalid BPM values

## Usage Example

```go
calc := beatgrid.NewCalculator()
// 120 BPM, 180 seconds, no offset
grid := calc.Calculate(120, 180.0, 0.0)

// Find beat at 5 seconds
beatIdx := grid.GetBeatAtTime(5000)
isDownbeat := grid.IsDownbeat(beatIdx)
```

## Integration Points

- Called after BPM detection in upload processor
- Beat grid stored in Track.BeatGrid field ([]int64)
- Used by frontend waveform component for beat overlay
