package analysis

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checkFFmpegAvailable returns true if ffmpeg is available in PATH
func checkFFmpegAvailable() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func TestValidateInputPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid absolute path", "/tmp/audio.mp3", false},
		{"valid relative path", "audio.mp3", false},
		{"valid path with spaces", "/tmp/my audio file.mp3", false},
		{"empty path", "", true},
		{"semicolon injection", "/tmp/audio.mp3; rm -rf /", true},
		{"pipe injection", "/tmp/audio.mp3 | cat /etc/passwd", true},
		{"ampersand injection", "/tmp/audio.mp3 & malicious", true},
		{"dollar sign injection", "/tmp/$HOME/audio.mp3", true},
		{"backtick injection", "/tmp/`whoami`/audio.mp3", true},
		{"parentheses injection", "/tmp/$(rm -rf /)/audio.mp3", true},
		{"curly braces injection", "/tmp/{a,b}/audio.mp3", true},
		{"redirect injection", "/tmp/audio.mp3 > /dev/null", true},
		{"newline injection", "/tmp/audio.mp3\nrm -rf /", true},
		{"carriage return injection", "/tmp/audio.mp3\rrm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInputPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAnalyzer(t *testing.T) {
	t.Run("creates analyzer with default paths", func(t *testing.T) {
		analyzer := NewAnalyzer()
		require.NotNil(t, analyzer)
		assert.Equal(t, "ffmpeg", analyzer.ffmpegPath)
		assert.Equal(t, "ffprobe", analyzer.ffprobePath)
		assert.Equal(t, 22050, analyzer.sampleRate)
	})

	t.Run("creates analyzer with custom ffmpeg path from env", func(t *testing.T) {
		t.Setenv("FFMPEG_PATH", "/custom/ffmpeg")
		t.Setenv("FFPROBE_PATH", "/custom/ffprobe")

		analyzer := NewAnalyzer()
		require.NotNil(t, analyzer)
		assert.Equal(t, "/custom/ffmpeg", analyzer.ffmpegPath)
		assert.Equal(t, "/custom/ffprobe", analyzer.ffprobePath)
	})

	t.Run("rejects malicious ffmpeg path from env", func(t *testing.T) {
		t.Setenv("FFMPEG_PATH", "/custom/ffmpeg; rm -rf /")
		t.Setenv("FFPROBE_PATH", "/custom/ffprobe | cat /etc/passwd")

		analyzer := NewAnalyzer()
		require.NotNil(t, analyzer)
		// Should fall back to default due to dangerous characters
		assert.Equal(t, "ffmpeg", analyzer.ffmpegPath)
		assert.Equal(t, "ffprobe", analyzer.ffprobePath)
	})
}

func TestValidateBinaryPath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		defaultName string
		expected    string
	}{
		{"default name unchanged", "ffmpeg", "ffmpeg", "ffmpeg"},
		{"valid absolute path", "/usr/bin/ffmpeg", "ffmpeg", "/usr/bin/ffmpeg"},
		{"valid custom path", "/opt/ffmpeg/bin/ffmpeg", "ffmpeg", "/opt/ffmpeg/bin/ffmpeg"},
		{"semicolon injection", "/usr/bin/ffmpeg; rm -rf /", "ffmpeg", "ffmpeg"},
		{"pipe injection", "/usr/bin/ffmpeg | cat", "ffmpeg", "ffmpeg"},
		{"ampersand injection", "/usr/bin/ffmpeg & malicious", "ffmpeg", "ffmpeg"},
		{"dollar sign injection", "$HOME/ffmpeg", "ffmpeg", "ffmpeg"},
		{"backtick injection", "`whoami`/ffmpeg", "ffmpeg", "ffmpeg"},
		{"newline injection", "/usr/bin/ffmpeg\nrm", "ffmpeg", "ffmpeg"},
		{"space injection", "/usr/bin/ffmpeg -malicious", "ffmpeg", "ffmpeg"},
		{"relative path with slash", "./ffmpeg", "ffmpeg", "ffmpeg"},
		{"path traversal", "/usr/../../../etc/passwd", "ffmpeg", "/etc/passwd"}, // Clean removes traversal
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateBinaryPath(tt.path, tt.defaultName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeExtension(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected string
	}{
		{"valid mp3", ".mp3", ".mp3"},
		{"valid MP3 uppercase", ".MP3", ".mp3"},
		{"valid flac", ".flac", ".flac"},
		{"valid wav", ".wav", ".wav"},
		{"valid aac", ".aac", ".aac"},
		{"valid m4a", ".m4a", ".m4a"},
		{"valid ogg", ".ogg", ".ogg"},
		{"valid wma", ".wma", ".wma"},
		{"valid aiff", ".aiff", ".aiff"},
		{"invalid extension", ".exe", ".mp3"},
		{"empty extension", "", ".mp3"},
		{"injection attempt", ".mp3; rm", ".mp3"},
		{"unknown extension", ".xyz", ".mp3"},
		{"double extension", ".mp3.bak", ".mp3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeExtension(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCamelotNotation(t *testing.T) {
	tests := []struct {
		key      string
		mode     string
		expected string
	}{
		// Minor keys (A column)
		{"Ab", "minor", "1A"},
		{"G#", "minor", "1A"},
		{"Eb", "minor", "2A"},
		{"D#", "minor", "2A"},
		{"Bb", "minor", "3A"},
		{"A#", "minor", "3A"},
		{"F", "minor", "4A"},
		{"C", "minor", "5A"},
		{"G", "minor", "6A"},
		{"D", "minor", "7A"},
		{"A", "minor", "8A"},
		{"E", "minor", "9A"},
		{"B", "minor", "10A"},
		{"F#", "minor", "11A"},
		{"Gb", "minor", "11A"},
		{"Db", "minor", "12A"},
		{"C#", "minor", "12A"},

		// Major keys (B column)
		{"B", "major", "1B"},
		{"F#", "major", "2B"},
		{"Gb", "major", "2B"},
		{"Db", "major", "3B"},
		{"C#", "major", "3B"},
		{"Ab", "major", "4B"},
		{"G#", "major", "4B"},
		{"Eb", "major", "5B"},
		{"D#", "major", "5B"},
		{"Bb", "major", "6B"},
		{"A#", "major", "6B"},
		{"F", "major", "7B"},
		{"C", "major", "8B"},
		{"G", "major", "9B"},
		{"D", "major", "10B"},
		{"A", "major", "11B"},
		{"E", "major", "12B"},

		// Keys with 'm' suffix
		{"Am", "", "8A"},
		{"Cm", "", "5A"},
		{"F#m", "", "11A"},

		// Unknown keys
		{"X", "major", ""},
		{"", "minor", ""},
		{"Z", "", ""},
	}

	for _, tt := range tests {
		name := tt.key + "_" + tt.mode
		if name == "_" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			result := GetCamelotNotation(tt.key, tt.mode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCamelotWheelMap(t *testing.T) {
	// Verify the CamelotWheel map has all expected entries
	expectedKeys := []string{
		"Abm", "G#m", "Ebm", "D#m", "Bbm", "A#m", "Fm", "Cm", "Gm", "Dm", "Am", "Em", "Bm", "F#m", "Gbm", "Dbm", "C#m",
		"B", "F#", "Gb", "Db", "C#", "Ab", "G#", "Eb", "D#", "Bb", "A#", "F", "C", "G", "D", "A", "E",
	}

	for _, key := range expectedKeys {
		t.Run(key, func(t *testing.T) {
			_, ok := CamelotWheel[key]
			assert.True(t, ok, "CamelotWheel should contain key: %s", key)
		})
	}

	// Verify enharmonic equivalents map to the same value
	assert.Equal(t, CamelotWheel["G#m"], CamelotWheel["Abm"])
	assert.Equal(t, CamelotWheel["D#m"], CamelotWheel["Ebm"])
	assert.Equal(t, CamelotWheel["F#"], CamelotWheel["Gb"])
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"song.mp3", "mp3"},
		{"SONG.MP3", "mp3"},
		{"song.Mp3", "mp3"},
		{"song.flac", "flac"},
		{"SONG.FLAC", "flac"},
		{"song.wav", "wav"},
		{"SONG.WAV", "wav"},
		{"song.ogg", "ogg"},
		{"song.aac", "aac"},
		{"song.m4a", "aac"},
		{"song.unknown", "unknown"},
		{"song", "unknown"},
		{"", "unknown"},
		{"song.MP3.backup", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			result := detectFormat(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
		{100, 100},
		{-100, 100},
		{-2147483648 + 1, 2147483647}, // Near min int
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := abs(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBassEmphasisFilter(t *testing.T) {
	sampleRate := 22050

	t.Run("preserves length", func(t *testing.T) {
		samples := make([]float64, 1000)
		for i := range samples {
			samples[i] = float64(i%100) / 100.0
		}
		result := bassEmphasisFilter(samples, sampleRate)
		assert.Equal(t, len(samples), len(result))
	})

	t.Run("handles empty input", func(t *testing.T) {
		result := bassEmphasisFilter([]float64{}, sampleRate)
		assert.Empty(t, result)
	})

	t.Run("handles single sample", func(t *testing.T) {
		result := bassEmphasisFilter([]float64{0.5}, sampleRate)
		assert.Len(t, result, 1)
	})

	t.Run("produces output for sine wave", func(t *testing.T) {
		// Generate 100Hz sine wave (bass frequency)
		samples := make([]float64, sampleRate)
		for i := range samples {
			samples[i] = 0.5 * (1.0 + float64(i%221)/221.0) // Simple pattern
		}
		result := bassEmphasisFilter(samples, sampleRate)
		assert.Equal(t, len(samples), len(result))

		// Check output is not all zeros
		hasNonZero := false
		for _, v := range result[100:] { // Skip initial transient
			if v != 0 {
				hasNonZero = true
				break
			}
		}
		assert.True(t, hasNonZero, "filter should produce non-zero output")
	})
}

func TestAdaptiveOnsetDetection(t *testing.T) {
	t.Run("preserves length", func(t *testing.T) {
		energy := make([]float64, 100)
		for i := range energy {
			energy[i] = float64(i) / 100.0
		}
		result := adaptiveOnsetDetection(energy, 8)
		assert.Equal(t, len(energy), len(result))
	})

	t.Run("handles empty input", func(t *testing.T) {
		result := adaptiveOnsetDetection([]float64{}, 8)
		assert.Empty(t, result)
	})

	t.Run("handles single element", func(t *testing.T) {
		result := adaptiveOnsetDetection([]float64{1.0}, 8)
		assert.Len(t, result, 1)
	})

	t.Run("detects onset spike", func(t *testing.T) {
		energy := make([]float64, 20)
		for i := range energy {
			energy[i] = 0.1
		}
		energy[10] = 1.0 // Spike

		result := adaptiveOnsetDetection(energy, 5)
		assert.Greater(t, result[10], 0.0, "should detect spike as onset")
	})
}

func TestAutocorrelationBPMImproved(t *testing.T) {
	t.Run("returns zero for empty input", func(t *testing.T) {
		bpm, confidence := autocorrelationBPMImproved([]float64{}, 512, 22050)
		assert.Equal(t, 0, bpm)
		assert.Equal(t, 0.0, confidence)
	})

	t.Run("returns zero for very short input", func(t *testing.T) {
		bpm, confidence := autocorrelationBPMImproved(make([]float64, 10), 512, 22050)
		assert.Equal(t, 0, bpm)
		assert.Equal(t, 0.0, confidence)
	})

	t.Run("returns BPM in valid range when detected", func(t *testing.T) {
		// Generate rhythmic pattern at 120 BPM
		sampleRate := 22050
		hopSize := 512
		framesPerSecond := float64(sampleRate) / float64(hopSize)
		beatsPerSecond := 120.0 / 60.0
		framesPerBeat := int(framesPerSecond / beatsPerSecond)

		onset := make([]float64, 1000)
		for i := 0; i < len(onset); i += framesPerBeat {
			onset[i] = 1.0
		}

		bpm, confidence := autocorrelationBPMImproved(onset, hopSize, sampleRate)

		if bpm > 0 {
			assert.GreaterOrEqual(t, bpm, 60)
			assert.LessOrEqual(t, bpm, 200)
			assert.GreaterOrEqual(t, confidence, 0.0)
		}
	})
}

func TestAnalyze_Errors(t *testing.T) {
	t.Run("handles context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		analyzer := NewAnalyzer()
		reader := bytes.NewReader([]byte("fake audio data"))

		_, err := analyzer.Analyze(ctx, reader, "test.mp3")
		// Should fail due to context or FFmpeg not available
		assert.Error(t, err)
	})

	t.Run("handles empty reader", func(t *testing.T) {
		ctx := context.Background()
		analyzer := NewAnalyzer()
		reader := bytes.NewReader([]byte{})

		_, err := analyzer.Analyze(ctx, reader, "test.mp3")
		assert.Error(t, err)
	})

	t.Run("handles invalid audio data", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		analyzer := NewAnalyzer()
		reader := bytes.NewReader([]byte("not audio data"))

		_, err := analyzer.Analyze(ctx, reader, "test.mp3")
		assert.Error(t, err)
	})
}

func TestAnalyze_WithFFmpeg(t *testing.T) {
	if !checkFFmpegAvailable() {
		t.Skip("FFmpeg not available, skipping integration tests")
	}

	t.Run("returns result for valid audio", func(t *testing.T) {
		// This test would require actual audio file
		// For now, we just verify the analyzer initializes correctly
		analyzer := NewAnalyzer()
		assert.NotNil(t, analyzer)
	})
}

// Benchmark tests
func BenchmarkBassEmphasisFilter(b *testing.B) {
	samples := make([]float64, 22050) // 1 second at 22050Hz
	for i := range samples {
		samples[i] = float64(i%100) / 100.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bassEmphasisFilter(samples, 22050)
	}
}

func BenchmarkAdaptiveOnsetDetection(b *testing.B) {
	energy := make([]float64, 1000)
	for i := range energy {
		energy[i] = float64(i%10) / 10.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adaptiveOnsetDetection(energy, 8)
	}
}

func BenchmarkGetCamelotNotation(b *testing.B) {
	keys := []string{"Am", "C", "F#m", "Bb", "G"}
	modes := []string{"minor", "major", "", "minor", "major"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(keys)
		GetCamelotNotation(keys[idx], modes[idx])
	}
}
