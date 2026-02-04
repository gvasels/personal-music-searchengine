// Package beatgrid provides beat grid calculation from BPM and duration
package beatgrid

import (
	"encoding/json"
	"sort"
)

// BeatGrid contains beat timing information for DJ features
type BeatGrid struct {
	// BPM is the beats per minute of the track
	BPM int `json:"bpm"`
	// Offset is the time offset to the first beat in milliseconds
	Offset float64 `json:"offset"`
	// Beats contains timestamps of each beat in milliseconds
	Beats []int64 `json:"beats"`
	// Downbeats contains indices of downbeats (every 4th beat, starting at 0)
	Downbeats []int `json:"downbeats"`
	// IsVariable indicates if the BPM varies throughout the track
	IsVariable bool `json:"isVariable"`
}

// Calculator generates beat grids from BPM information
type Calculator struct{}

// NewCalculator creates a new beat grid calculator
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Calculate generates a beat grid from BPM, duration, and offset
func (c *Calculator) Calculate(bpm int, durationSeconds float64, offsetMs float64) *BeatGrid {
	// Validate inputs
	if bpm < 20 || bpm > 300 {
		return nil
	}
	if durationSeconds <= 0 {
		return nil
	}
	if offsetMs < 0 {
		offsetMs = 0
	}

	// Calculate milliseconds per beat
	msPerBeat := MsPerBeat(bpm)
	if msPerBeat == 0 {
		return nil
	}

	// Calculate total duration in milliseconds
	durationMs := durationSeconds * 1000.0

	// Calculate number of beats
	numBeats := int((durationMs - offsetMs) / msPerBeat)
	if numBeats < 0 {
		numBeats = 0
	}

	// Generate beat timestamps
	beats := make([]int64, numBeats)
	downbeats := make([]int, 0)

	for i := 0; i < numBeats; i++ {
		beatTime := offsetMs + float64(i)*msPerBeat
		beats[i] = int64(beatTime)

		// Mark every 4th beat as downbeat (0, 4, 8, 12, ...)
		if i%4 == 0 {
			downbeats = append(downbeats, i)
		}
	}

	return &BeatGrid{
		BPM:        bpm,
		Offset:     offsetMs,
		Beats:      beats,
		Downbeats:  downbeats,
		IsVariable: false,
	}
}

// Validate checks if the beat grid is valid
func (bg *BeatGrid) Validate() bool {
	// Check BPM range (typical music range)
	if bg.BPM < 20 || bg.BPM > 300 {
		return false
	}

	// Check offset is not negative
	if bg.Offset < 0 {
		return false
	}

	// Check beats exist
	if len(bg.Beats) == 0 {
		return false
	}

	// Check beats are in ascending order
	for i := 1; i < len(bg.Beats); i++ {
		if bg.Beats[i] <= bg.Beats[i-1] {
			return false
		}
	}

	// Check downbeats are valid indices
	for _, db := range bg.Downbeats {
		if db < 0 || db >= len(bg.Beats) {
			return false
		}
	}

	return true
}

// GetBeatAtTime returns the index of the beat nearest to the given time
func (bg *BeatGrid) GetBeatAtTime(timeMs int64) int {
	if len(bg.Beats) == 0 {
		return -1
	}

	// Handle time before first beat
	if timeMs < bg.Beats[0] {
		return 0
	}

	// Handle time after last beat
	if timeMs >= bg.Beats[len(bg.Beats)-1] {
		return len(bg.Beats) - 1
	}

	// Binary search for closest beat
	idx := sort.Search(len(bg.Beats), func(i int) bool {
		return bg.Beats[i] >= timeMs
	})

	// Check if previous beat is closer
	if idx > 0 {
		prevDiff := timeMs - bg.Beats[idx-1]
		currDiff := bg.Beats[idx] - timeMs
		if prevDiff <= currDiff {
			return idx - 1
		}
	}

	return idx
}

// GetTimeAtBeat returns the timestamp of the beat at the given index
func (bg *BeatGrid) GetTimeAtBeat(beatIndex int) int64 {
	if beatIndex < 0 || beatIndex >= len(bg.Beats) {
		return -1
	}
	return bg.Beats[beatIndex]
}

// IsDownbeat returns true if the beat at the given index is a downbeat
func (bg *BeatGrid) IsDownbeat(beatIndex int) bool {
	if beatIndex < 0 || beatIndex >= len(bg.Beats) {
		return false
	}

	for _, db := range bg.Downbeats {
		if db == beatIndex {
			return true
		}
	}

	return false
}

// MsPerBeat returns the number of milliseconds per beat
func MsPerBeat(bpm int) float64 {
	if bpm <= 0 {
		return 0
	}
	// 60 seconds * 1000 ms / BPM
	return 60000.0 / float64(bpm)
}

// ToJSON serializes beat grid to JSON
func (bg *BeatGrid) ToJSON() ([]byte, error) {
	return json.Marshal(bg)
}

// FromJSON deserializes beat grid from JSON
func FromJSON(data []byte) (*BeatGrid, error) {
	var bg BeatGrid
	if err := json.Unmarshal(data, &bg); err != nil {
		return nil, err
	}
	return &bg, nil
}
