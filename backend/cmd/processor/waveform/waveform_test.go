package waveform

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skipIfNoFFmpeg skips the test if FFmpeg/FFprobe is not available
func skipIfNoFFmpeg(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("Skipping test: ffprobe not available in PATH")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("Skipping test: ffmpeg not available in PATH")
	}
}

// =============================================================================
// WaveformData Tests
// =============================================================================

func TestWaveformData_Validate_ValidData(t *testing.T) {
	data := &WaveformData{
		Peaks:      []float64{0.5, 0.7, 0.9, 0.6, 0.4, 0.3, 0.8, 0.5},
		SampleRate: 100,
		Duration:   180.5,
		Version:    1,
	}

	assert.True(t, data.Validate(), "Valid waveform data should pass validation")
}

func TestWaveformData_Validate_EmptyPeaks(t *testing.T) {
	data := &WaveformData{
		Peaks:      []float64{},
		SampleRate: 100,
		Duration:   180.5,
		Version:    1,
	}

	assert.False(t, data.Validate(), "Empty peaks should fail validation")
}

func TestWaveformData_Validate_InvalidPeakRange(t *testing.T) {
	tests := []struct {
		name   string
		peaks  []float64
		expect bool
	}{
		{"negative peak", []float64{0.5, -0.1, 0.7}, false},
		{"peak above 1", []float64{0.5, 1.5, 0.7}, false},
		{"all valid", []float64{0.0, 0.5, 1.0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &WaveformData{
				Peaks:      tt.peaks,
				SampleRate: 100,
				Duration:   10.0,
				Version:    1,
			}
			assert.Equal(t, tt.expect, data.Validate())
		})
	}
}

func TestWaveformData_Validate_InvalidSampleRate(t *testing.T) {
	data := &WaveformData{
		Peaks:      []float64{0.5, 0.7},
		SampleRate: 0, // Invalid
		Duration:   10.0,
		Version:    1,
	}

	assert.False(t, data.Validate(), "Zero sample rate should fail validation")
}

func TestWaveformData_Validate_InvalidDuration(t *testing.T) {
	data := &WaveformData{
		Peaks:      []float64{0.5, 0.7},
		SampleRate: 100,
		Duration:   0.0, // Invalid
		Version:    1,
	}

	assert.False(t, data.Validate(), "Zero duration should fail validation")
}

func TestWaveformData_JSONSerialization(t *testing.T) {
	original := &WaveformData{
		Peaks:      []float64{0.5, 0.7, 0.9, 0.6, 0.4},
		SampleRate: 100,
		Duration:   180.5,
		Version:    1,
	}

	jsonBytes, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded WaveformData
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.Peaks, decoded.Peaks)
	assert.Equal(t, original.SampleRate, decoded.SampleRate)
	assert.Equal(t, original.Duration, decoded.Duration)
	assert.Equal(t, original.Version, decoded.Version)
}

func TestWaveformData_JSONFields(t *testing.T) {
	data := &WaveformData{
		Peaks:      []float64{0.5, 0.7},
		SampleRate: 100,
		Duration:   180.5,
		Version:    1,
	}

	jsonBytes, err := json.Marshal(data)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "peaks")
	assert.Contains(t, jsonMap, "sampleRate")
	assert.Contains(t, jsonMap, "duration")
	assert.Contains(t, jsonMap, "version")
}

// =============================================================================
// Generator Tests
// =============================================================================

func TestGenerator_Generate_ValidMP3(t *testing.T) {
	skipIfNoFFmpeg(t)

	// Create a test fixture MP3 file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")

	// Write minimal valid MP3 header (ID3v2 tag)
	// Real implementation would use actual MP3 files
	mp3Header := []byte{
		0x49, 0x44, 0x33, // ID3 magic
		0x04, 0x00, // Version
		0x00,       // Flags
		0x00, 0x00, 0x00, 0x00, // Size
	}
	err := os.WriteFile(testFile, mp3Header, 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, testFile)

	// TDD Red: This should fail because implementation returns error
	require.NoError(t, err, "Generate should not return error for valid MP3")
	require.NotNil(t, data, "Generate should return waveform data")
	assert.Equal(t, 100, data.SampleRate, "Sample rate should be 100")
	assert.NotEmpty(t, data.Peaks, "Peaks should not be empty")
	assert.Greater(t, data.Duration, 0.0, "Duration should be positive")
}

func TestGenerator_Generate_ValidFLAC(t *testing.T) {
	skipIfNoFFmpeg(t)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.flac")

	// FLAC magic bytes
	flacHeader := []byte{0x66, 0x4c, 0x61, 0x43} // "fLaC"
	err := os.WriteFile(testFile, flacHeader, 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, testFile)

	require.NoError(t, err, "Generate should not return error for valid FLAC")
	require.NotNil(t, data)
}

func TestGenerator_Generate_ValidWAV(t *testing.T) {
	skipIfNoFFmpeg(t)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.wav")

	// WAV RIFF header
	wavHeader := []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x00, 0x00, 0x00, 0x00, // File size
		0x57, 0x41, 0x56, 0x45, // "WAVE"
	}
	err := os.WriteFile(testFile, wavHeader, 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, testFile)

	require.NoError(t, err, "Generate should not return error for valid WAV")
	require.NotNil(t, data)
}

func TestGenerator_Generate_NonExistentFile(t *testing.T) {
	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, "/nonexistent/path/audio.mp3")

	assert.Error(t, err, "Generate should return error for non-existent file")
	assert.Nil(t, data)
}

func TestGenerator_Generate_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.txt")

	err := os.WriteFile(testFile, []byte("not audio data"), 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, testFile)

	assert.Error(t, err, "Generate should return error for invalid audio file")
	assert.True(t, err == ErrInvalidAudioFile || err == ErrUnsupportedFormat,
		"Should return ErrInvalidAudioFile or ErrUnsupportedFormat")
	assert.Nil(t, data)
}

func TestGenerator_Generate_PeaksNormalized(t *testing.T) {
	skipIfNoFFmpeg(t)

	// This test verifies that generated peaks are normalized 0.0-1.0
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")

	// Write test audio file (would use real fixture in practice)
	mp3Header := []byte{0x49, 0x44, 0x33, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := os.WriteFile(testFile, mp3Header, 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, testFile)

	require.NoError(t, err)
	require.NotNil(t, data)

	for i, peak := range data.Peaks {
		assert.GreaterOrEqual(t, peak, 0.0, "Peak %d should be >= 0.0", i)
		assert.LessOrEqual(t, peak, 1.0, "Peak %d should be <= 1.0", i)
	}
}

func TestGenerator_Generate_SampleRateCorrect(t *testing.T) {
	skipIfNoFFmpeg(t)

	// Verify waveform has ~100 samples per second
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")

	mp3Header := []byte{0x49, 0x44, 0x33, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := os.WriteFile(testFile, mp3Header, 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.Generate(ctx, testFile)

	require.NoError(t, err)
	require.NotNil(t, data)

	assert.Equal(t, 100, data.SampleRate, "Sample rate should be 100 samples/second")

	// Expected peaks count should be roughly duration * sampleRate
	expectedPeaks := int(data.Duration * float64(data.SampleRate))
	actualPeaks := len(data.Peaks)

	// Allow 1% tolerance
	tolerance := expectedPeaks / 100
	if tolerance < 1 {
		tolerance = 1
	}
	assert.InDelta(t, expectedPeaks, actualPeaks, float64(tolerance),
		"Peak count should match duration * sampleRate")
}

func TestGenerator_Generate_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")

	mp3Header := []byte{0x49, 0x44, 0x33, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := os.WriteFile(testFile, mp3Header, 0644)
	require.NoError(t, err)

	gen := NewGenerator()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	data, err := gen.Generate(ctx, testFile)

	assert.Error(t, err, "Generate should return error when context is cancelled")
	assert.Nil(t, data)
}

// =============================================================================
// GenerateFromBytes Tests
// =============================================================================

func TestGenerator_GenerateFromBytes_ValidMP3(t *testing.T) {
	skipIfNoFFmpeg(t)

	// Minimal MP3 header bytes
	mp3Data := []byte{0x49, 0x44, 0x33, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.GenerateFromBytes(ctx, mp3Data, "mp3")

	require.NoError(t, err, "GenerateFromBytes should not return error for valid MP3")
	require.NotNil(t, data)
	assert.Equal(t, 100, data.SampleRate)
}

func TestGenerator_GenerateFromBytes_UnsupportedFormat(t *testing.T) {
	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.GenerateFromBytes(ctx, []byte("data"), "unknown")

	assert.Error(t, err)
	assert.Equal(t, ErrUnsupportedFormat, err)
	assert.Nil(t, data)
}

func TestGenerator_GenerateFromBytes_EmptyData(t *testing.T) {
	gen := NewGenerator()
	ctx := context.Background()

	data, err := gen.GenerateFromBytes(ctx, []byte{}, "mp3")

	assert.Error(t, err, "GenerateFromBytes should return error for empty data")
	assert.Nil(t, data)
}
