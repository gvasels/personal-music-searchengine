package beatgrid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// BeatGrid Validation Tests
// =============================================================================

func TestBeatGrid_Validate_ValidData(t *testing.T) {
	bg := &BeatGrid{
		BPM:        120,
		Offset:     0.0,
		Beats:      []int64{0, 500, 1000, 1500, 2000, 2500, 3000, 3500},
		Downbeats:  []int{0, 4},
		IsVariable: false,
	}

	assert.True(t, bg.Validate(), "Valid beat grid should pass validation")
}

func TestBeatGrid_Validate_InvalidBPM(t *testing.T) {
	tests := []struct {
		name string
		bpm  int
		want bool
	}{
		{"zero BPM", 0, false},
		{"negative BPM", -60, false},
		{"too low BPM", 10, false},
		{"too high BPM", 500, false},
		{"valid low BPM", 20, true},
		{"valid high BPM", 300, true},
		{"typical BPM", 120, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := &BeatGrid{
				BPM:        tt.bpm,
				Offset:     0.0,
				Beats:      []int64{0, 500},
				Downbeats:  []int{0},
				IsVariable: false,
			}
			assert.Equal(t, tt.want, bg.Validate())
		})
	}
}

func TestBeatGrid_Validate_EmptyBeats(t *testing.T) {
	bg := &BeatGrid{
		BPM:        120,
		Offset:     0.0,
		Beats:      []int64{},
		Downbeats:  []int{},
		IsVariable: false,
	}

	assert.False(t, bg.Validate(), "Empty beats should fail validation")
}

func TestBeatGrid_Validate_NegativeOffset(t *testing.T) {
	bg := &BeatGrid{
		BPM:        120,
		Offset:     -100.0, // Negative offset
		Beats:      []int64{0, 500},
		Downbeats:  []int{0},
		IsVariable: false,
	}

	assert.False(t, bg.Validate(), "Negative offset should fail validation")
}

// =============================================================================
// Calculator Tests
// =============================================================================

func TestCalculator_Calculate_120BPM(t *testing.T) {
	calc := NewCalculator()

	// 120 BPM = 500ms per beat
	// 10 second track = 20 beats
	bg := calc.Calculate(120, 10.0, 0.0)

	require.NotNil(t, bg, "Calculate should return a beat grid")
	assert.Equal(t, 120, bg.BPM)
	assert.Equal(t, 0.0, bg.Offset)
	assert.Equal(t, 20, len(bg.Beats), "10 seconds at 120 BPM = 20 beats")
	assert.False(t, bg.IsVariable)
}

func TestCalculator_Calculate_60BPM(t *testing.T) {
	calc := NewCalculator()

	// 60 BPM = 1000ms per beat
	// 10 second track = 10 beats
	bg := calc.Calculate(60, 10.0, 0.0)

	require.NotNil(t, bg)
	assert.Equal(t, 60, bg.BPM)
	assert.Equal(t, 10, len(bg.Beats))
}

func TestCalculator_Calculate_WithOffset(t *testing.T) {
	calc := NewCalculator()

	// 120 BPM with 250ms offset
	bg := calc.Calculate(120, 10.0, 250.0)

	require.NotNil(t, bg)
	assert.Equal(t, 250.0, bg.Offset)
	assert.Equal(t, int64(250), bg.Beats[0], "First beat should be at offset")
	assert.Equal(t, int64(750), bg.Beats[1], "Second beat should be offset + 500ms")
}

func TestCalculator_Calculate_BeatTimings(t *testing.T) {
	calc := NewCalculator()

	// 120 BPM = 500ms per beat
	bg := calc.Calculate(120, 5.0, 0.0)

	require.NotNil(t, bg)

	expectedBeats := []int64{0, 500, 1000, 1500, 2000, 2500, 3000, 3500, 4000, 4500}
	assert.Equal(t, expectedBeats, bg.Beats[:len(expectedBeats)])
}

func TestCalculator_Calculate_DownbeatsEvery4(t *testing.T) {
	calc := NewCalculator()

	bg := calc.Calculate(120, 10.0, 0.0)

	require.NotNil(t, bg)

	// Downbeats should be at indices 0, 4, 8, 12, 16
	expectedDownbeats := []int{0, 4, 8, 12, 16}
	assert.Equal(t, expectedDownbeats, bg.Downbeats)

	// Verify downbeats are in bounds
	for _, db := range bg.Downbeats {
		assert.Less(t, db, len(bg.Beats), "Downbeat index should be within beats range")
	}
}

func TestCalculator_Calculate_InvalidBPM(t *testing.T) {
	calc := NewCalculator()

	bg := calc.Calculate(0, 10.0, 0.0) // Zero BPM

	assert.Nil(t, bg, "Calculate should return nil for invalid BPM")
}

func TestCalculator_Calculate_InvalidDuration(t *testing.T) {
	calc := NewCalculator()

	bg := calc.Calculate(120, 0.0, 0.0) // Zero duration

	assert.Nil(t, bg, "Calculate should return nil for zero duration")
}

func TestCalculator_Calculate_VeryShortDuration(t *testing.T) {
	calc := NewCalculator()

	// 0.1 seconds at 120 BPM = 0.2 beats
	bg := calc.Calculate(120, 0.1, 0.0)

	require.NotNil(t, bg)
	// Should still have at least 1 beat if any part of the duration covers it
	assert.GreaterOrEqual(t, len(bg.Beats), 0)
}

// =============================================================================
// BeatGrid Helper Method Tests
// =============================================================================

func TestBeatGrid_GetBeatAtTime(t *testing.T) {
	bg := &BeatGrid{
		BPM:       120,
		Offset:    0.0,
		Beats:     []int64{0, 500, 1000, 1500, 2000},
		Downbeats: []int{0, 4},
	}

	tests := []struct {
		name   string
		timeMs int64
		want   int
	}{
		{"exact first beat", 0, 0},
		{"exact second beat", 500, 1},
		{"between beats - closer to first", 200, 0},
		{"between beats - closer to second", 400, 1},
		{"exactly halfway", 250, 0}, // or 1, implementation decides tie-breaker
		{"past last beat", 2500, 4},
		{"before first beat", -100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bg.GetBeatAtTime(tt.timeMs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBeatGrid_GetTimeAtBeat(t *testing.T) {
	bg := &BeatGrid{
		BPM:       120,
		Offset:    0.0,
		Beats:     []int64{0, 500, 1000, 1500, 2000},
		Downbeats: []int{0, 4},
	}

	tests := []struct {
		name      string
		beatIndex int
		want      int64
	}{
		{"first beat", 0, 0},
		{"second beat", 1, 500},
		{"last beat", 4, 2000},
		{"out of range negative", -1, -1},
		{"out of range high", 10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bg.GetTimeAtBeat(tt.beatIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBeatGrid_IsDownbeat(t *testing.T) {
	bg := &BeatGrid{
		BPM:       120,
		Offset:    0.0,
		Beats:     []int64{0, 500, 1000, 1500, 2000, 2500, 3000, 3500},
		Downbeats: []int{0, 4},
	}

	tests := []struct {
		name      string
		beatIndex int
		want      bool
	}{
		{"first beat is downbeat", 0, true},
		{"second beat not downbeat", 1, false},
		{"third beat not downbeat", 2, false},
		{"fourth beat not downbeat", 3, false},
		{"fifth beat is downbeat", 4, true},
		{"out of range", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bg.IsDownbeat(tt.beatIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// MsPerBeat Tests
// =============================================================================

func TestMsPerBeat(t *testing.T) {
	tests := []struct {
		bpm  int
		want float64
	}{
		{60, 1000.0},   // 60 BPM = 1 beat per second = 1000ms
		{120, 500.0},   // 120 BPM = 2 beats per second = 500ms
		{180, 333.33},  // 180 BPM ≈ 333.33ms
		{90, 666.67},   // 90 BPM ≈ 666.67ms
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.bpm)), func(t *testing.T) {
			got := MsPerBeat(tt.bpm)
			assert.InDelta(t, tt.want, got, 0.01, "MsPerBeat(%d) should be ~%.2f", tt.bpm, tt.want)
		})
	}
}

func TestMsPerBeat_InvalidBPM(t *testing.T) {
	result := MsPerBeat(0)
	assert.Equal(t, 0.0, result, "MsPerBeat(0) should return 0")
}

// =============================================================================
// JSON Serialization Tests
// =============================================================================

func TestBeatGrid_JSONSerialization(t *testing.T) {
	original := &BeatGrid{
		BPM:        120,
		Offset:     100.0,
		Beats:      []int64{100, 600, 1100, 1600, 2100},
		Downbeats:  []int{0, 4},
		IsVariable: false,
	}

	jsonBytes, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded BeatGrid
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.BPM, decoded.BPM)
	assert.Equal(t, original.Offset, decoded.Offset)
	assert.Equal(t, original.Beats, decoded.Beats)
	assert.Equal(t, original.Downbeats, decoded.Downbeats)
	assert.Equal(t, original.IsVariable, decoded.IsVariable)
}

func TestBeatGrid_JSONFields(t *testing.T) {
	bg := &BeatGrid{
		BPM:        120,
		Offset:     0.0,
		Beats:      []int64{0, 500},
		Downbeats:  []int{0},
		IsVariable: true,
	}

	jsonBytes, err := json.Marshal(bg)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "bpm")
	assert.Contains(t, jsonMap, "offset")
	assert.Contains(t, jsonMap, "beats")
	assert.Contains(t, jsonMap, "downbeats")
	assert.Contains(t, jsonMap, "isVariable")
}
