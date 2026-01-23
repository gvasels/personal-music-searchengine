// Package analysis provides audio analysis functionality for BPM and key detection.
package analysis

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Result contains the audio analysis results
type Result struct {
	BPM        int    // Beats per minute (0 if not detected)
	MusicalKey string // Musical key (e.g., "Am", "C", "F#m")
	KeyMode    string // "major" or "minor"
	KeyCamelot string // Camelot notation (e.g., "8A", "11B")
}

// Analyzer performs audio analysis for BPM and key detection
type Analyzer struct {
	ffmpegPath  string
	ffprobePath string
	sampleRate  int // Sample rate for analysis (lower = faster, less accurate)
}

// NewAnalyzer creates a new audio analyzer
func NewAnalyzer() *Analyzer {
	// Check for FFmpeg path from environment or use default
	ffmpegPath := os.Getenv("FFMPEG_PATH")
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	ffprobePath := os.Getenv("FFPROBE_PATH")
	if ffprobePath == "" {
		ffprobePath = "ffprobe"
	}

	return &Analyzer{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
		sampleRate:  22050, // 22kHz is enough for beat detection
	}
}

// Analyze performs BPM and key detection on an audio stream
func (a *Analyzer) Analyze(ctx context.Context, reader io.Reader, fileName string) (*Result, error) {
	result := &Result{}

	// Create temp file for the audio
	tempDir := os.TempDir()
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".mp3"
	}
	tempFile, err := os.CreateTemp(tempDir, "audio-*"+ext)
	if err != nil {
		return result, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	// Write audio to temp file
	if _, err := io.Copy(tempFile, reader); err != nil {
		tempFile.Close()
		return result, fmt.Errorf("failed to write temp file: %w", err)
	}
	tempFile.Close()

	// Decode audio to raw PCM using FFmpeg
	samples, err := a.decodeToMono(ctx, tempPath)
	if err != nil {
		return result, fmt.Errorf("failed to decode audio: %w", err)
	}

	if len(samples) < a.sampleRate*5 { // Need at least 5 seconds
		return result, fmt.Errorf("audio too short for analysis")
	}

	// Detect BPM
	bpm := a.detectBPM(samples)
	if bpm >= 20 && bpm <= 300 {
		result.BPM = bpm
	}

	// Key detection is more complex - skip for now
	// Would require pitch/chroma analysis

	return result, nil
}

// decodeToMono uses FFmpeg to decode audio to mono 16-bit PCM
func (a *Analyzer) decodeToMono(ctx context.Context, inputPath string) ([]float64, error) {
	// FFmpeg command to decode to raw mono PCM
	args := []string{
		"-i", inputPath,
		"-ac", "1", // Mono
		"-ar", fmt.Sprintf("%d", a.sampleRate), // Sample rate
		"-f", "s16le", // 16-bit signed little-endian
		"-acodec", "pcm_s16le",
		"-", // Output to stdout
	}

	cmd := exec.CommandContext(ctx, a.ffmpegPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	// Parse raw PCM data
	data := stdout.Bytes()
	numSamples := len(data) / 2 // 16-bit = 2 bytes per sample
	samples := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		sample := int16(binary.LittleEndian.Uint16(data[i*2:]))
		samples[i] = float64(sample) / 32768.0 // Normalize to -1.0 to 1.0
	}

	return samples, nil
}

// detectBPM analyzes samples and returns estimated BPM using multi-segment analysis
func (a *Analyzer) detectBPM(samples []float64) int {
	// Parameters for improved detection
	windowSize := a.sampleRate / 20  // 50ms windows (better transient detection)
	hopSize := windowSize / 2        // 25ms hop (50% overlap)

	// Analyze multiple segments of the track for consensus
	segmentDuration := 15 * a.sampleRate // 15 seconds per segment
	numSegments := 4
	segmentStep := (len(samples) - segmentDuration) / numSegments
	if segmentStep < 0 {
		segmentStep = 0
		numSegments = 1
	}

	bpmCandidates := make(map[int]int) // BPM -> vote count

	for seg := 0; seg < numSegments; seg++ {
		startSample := seg * segmentStep
		endSample := startSample + segmentDuration
		if endSample > len(samples) {
			endSample = len(samples)
		}
		if endSample-startSample < 5*a.sampleRate {
			continue // Skip segments shorter than 5 seconds
		}

		segment := samples[startSample:endSample]

		// Apply bass-emphasis filter (simple low-pass to focus on kick drum frequencies)
		filtered := bassEmphasisFilter(segment, a.sampleRate)

		// Calculate energy envelope
		numWindows := (len(filtered) - windowSize) / hopSize
		if numWindows < 20 {
			continue
		}

		energy := make([]float64, numWindows)
		for i := 0; i < numWindows; i++ {
			start := i * hopSize
			end := start + windowSize
			sum := 0.0
			for j := start; j < end; j++ {
				sum += filtered[j] * filtered[j]
			}
			energy[i] = math.Sqrt(sum / float64(windowSize))
		}

		// Normalize energy
		maxEnergy := 0.0
		for _, e := range energy {
			if e > maxEnergy {
				maxEnergy = e
			}
		}
		if maxEnergy > 0 {
			for i := range energy {
				energy[i] /= maxEnergy
			}
		}

		// Apply adaptive threshold for onset detection
		onset := adaptiveOnsetDetection(energy, 8)

		// Use autocorrelation with octave error correction
		bpm, confidence := autocorrelationBPMImproved(onset, hopSize, a.sampleRate)

		if bpm > 0 && confidence > 0.3 {
			// Round to nearest integer BPM
			roundedBPM := int(math.Round(float64(bpm)))
			bpmCandidates[roundedBPM]++

			// Also vote for octave-related tempos with lower weight
			halfBPM := roundedBPM / 2
			doubleBPM := roundedBPM * 2
			if halfBPM >= 60 && halfBPM <= 200 {
				bpmCandidates[halfBPM]++
			}
			if doubleBPM >= 60 && doubleBPM <= 200 {
				bpmCandidates[doubleBPM]++
			}
		}
	}

	if len(bpmCandidates) == 0 {
		return 0
	}

	// Find BPM with most votes, preferring common tempo ranges
	bestBPM := 0
	bestScore := 0
	for bpm, votes := range bpmCandidates {
		// Apply preference for common BPM ranges (120-130 is most common)
		score := votes
		if bpm >= 115 && bpm <= 135 {
			score += 2 // Bonus for house/techno range
		} else if bpm >= 135 && bpm <= 150 {
			score += 1 // Bonus for trance/D&B range
		} else if bpm >= 85 && bpm <= 95 {
			score += 1 // Bonus for hip-hop range
		}

		if score > bestScore || (score == bestScore && bpm > bestBPM) {
			bestScore = score
			bestBPM = bpm
		}
	}

	// Only return if we have reasonable confidence (at least 2 segments agreed)
	if bestScore < 2 && numSegments > 1 {
		return 0
	}

	return bestBPM
}

// bassEmphasisFilter applies a simple low-pass filter to emphasize bass frequencies
func bassEmphasisFilter(samples []float64, sampleRate int) []float64 {
	// Simple 2nd order IIR low-pass filter targeting ~200Hz
	// Cutoff frequency ratio
	fc := 200.0 / float64(sampleRate)
	q := 0.707 // Butterworth

	// Calculate filter coefficients
	w0 := 2.0 * math.Pi * fc
	alpha := math.Sin(w0) / (2.0 * q)

	b0 := (1 - math.Cos(w0)) / 2
	b1 := 1 - math.Cos(w0)
	b2 := (1 - math.Cos(w0)) / 2
	a0 := 1 + alpha
	a1 := -2 * math.Cos(w0)
	a2 := 1 - alpha

	// Normalize coefficients
	b0 /= a0
	b1 /= a0
	b2 /= a0
	a1 /= a0
	a2 /= a0

	// Apply filter
	filtered := make([]float64, len(samples))
	x1, x2, y1, y2 := 0.0, 0.0, 0.0, 0.0

	for i, x := range samples {
		y := b0*x + b1*x1 + b2*x2 - a1*y1 - a2*y2
		filtered[i] = y
		x2, x1 = x1, x
		y2, y1 = y1, y
	}

	return filtered
}

// adaptiveOnsetDetection computes onset strength with adaptive threshold
func adaptiveOnsetDetection(energy []float64, windowLen int) []float64 {
	onset := make([]float64, len(energy))

	for i := 1; i < len(energy); i++ {
		// Calculate local mean
		start := i - windowLen
		if start < 0 {
			start = 0
		}
		sum := 0.0
		for j := start; j < i; j++ {
			sum += energy[j]
		}
		localMean := sum / float64(i-start)

		// Onset is positive deviation from local mean
		diff := energy[i] - localMean
		if diff > 0 {
			onset[i] = diff
		}
	}

	return onset
}

// autocorrelationBPMImproved uses autocorrelation with confidence scoring
func autocorrelationBPMImproved(onset []float64, hopSize, sampleRate int) (int, float64) {
	// Expanded BPM range: 60-200 (covers house to D&B)
	minBPM := 60
	maxBPM := 200

	framesPerSecond := float64(sampleRate) / float64(hopSize)
	minLag := int(60.0 / float64(maxBPM) * framesPerSecond)
	maxLag := int(60.0 / float64(minBPM) * framesPerSecond)

	if maxLag >= len(onset)/2 {
		maxLag = len(onset)/2 - 1
	}
	if minLag < 1 {
		minLag = 1
	}
	if maxLag <= minLag {
		return 0, 0
	}

	// Compute normalized autocorrelation
	correlations := make([]float64, maxLag-minLag+1)
	maxCorr := 0.0

	// First compute the zero-lag correlation for normalization
	zeroLagSum := 0.0
	for i := 0; i < len(onset); i++ {
		zeroLagSum += onset[i] * onset[i]
	}

	for lag := minLag; lag <= maxLag; lag++ {
		sum := 0.0
		for i := 0; i < len(onset)-lag; i++ {
			sum += onset[i] * onset[i+lag]
		}
		// Normalize by zero-lag correlation
		if zeroLagSum > 0 {
			correlations[lag-minLag] = sum / zeroLagSum
		}
		if correlations[lag-minLag] > maxCorr {
			maxCorr = correlations[lag-minLag]
		}
	}

	if maxCorr == 0 {
		return 0, 0
	}

	// Find all significant peaks
	type peak struct {
		lag        int
		value      float64
		prominence float64
	}
	var peaks []peak

	for i := 2; i < len(correlations)-2; i++ {
		if correlations[i] > correlations[i-1] && correlations[i] > correlations[i+1] &&
			correlations[i] > correlations[i-2] && correlations[i] > correlations[i+2] {
			// Calculate prominence (how much peak stands out)
			minNeighbor := math.Min(
				math.Min(correlations[i-1], correlations[i-2]),
				math.Min(correlations[i+1], correlations[i+2]),
			)
			prominence := correlations[i] - minNeighbor
			if correlations[i] > maxCorr*0.3 { // Only consider peaks above 30% of max
				peaks = append(peaks, peak{
					lag:        i + minLag,
					value:      correlations[i],
					prominence: prominence,
				})
			}
		}
	}

	if len(peaks) == 0 {
		return 0, 0
	}

	// Sort peaks by value (descending)
	sort.Slice(peaks, func(i, j int) bool {
		return peaks[i].value > peaks[j].value
	})

	// Take the first peak as primary candidate
	bestLag := peaks[0].lag
	confidence := peaks[0].value

	// Check if there's a peak at half the lag (octave error detection)
	halfLag := bestLag / 2
	for _, p := range peaks {
		if abs(p.lag-halfLag) <= 2 && p.value > peaks[0].value*0.7 {
			// Strong peak at half the tempo - probably double-time detection
			bestLag = p.lag
			break
		}
	}

	// Convert lag to BPM
	bpm := int(math.Round(60.0 * framesPerSecond / float64(bestLag)))

	// Final sanity check
	if bpm < minBPM || bpm > maxBPM {
		return 0, 0
	}

	return bpm, confidence
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CamelotWheel maps musical keys to Camelot notation
var CamelotWheel = map[string]string{
	// Minor keys (A column)
	"Abm": "1A", "G#m": "1A",
	"Ebm": "2A", "D#m": "2A",
	"Bbm": "3A", "A#m": "3A",
	"Fm":  "4A",
	"Cm":  "5A",
	"Gm":  "6A",
	"Dm":  "7A",
	"Am":  "8A",
	"Em":  "9A",
	"Bm":  "10A",
	"F#m": "11A", "Gbm": "11A",
	"Dbm": "12A", "C#m": "12A",

	// Major keys (B column)
	"B":  "1B",
	"F#": "2B", "Gb": "2B",
	"Db": "3B", "C#": "3B",
	"Ab": "4B", "G#": "4B",
	"Eb": "5B", "D#": "5B",
	"Bb": "6B", "A#": "6B",
	"F":  "7B",
	"C":  "8B",
	"G":  "9B",
	"D":  "10B",
	"A":  "11B",
	"E":  "12B",
}

// GetCamelotNotation converts a musical key to Camelot notation
func GetCamelotNotation(key string, mode string) string {
	normalizedKey := key
	if mode == "minor" && !strings.HasSuffix(key, "m") {
		normalizedKey = key + "m"
	}

	if camelot, ok := CamelotWheel[normalizedKey]; ok {
		return camelot
	}
	return ""
}

// detectFormat determines the audio format from filename
func detectFormat(fileName string) string {
	lower := strings.ToLower(fileName)
	switch {
	case strings.HasSuffix(lower, ".mp3"):
		return "mp3"
	case strings.HasSuffix(lower, ".flac"):
		return "flac"
	case strings.HasSuffix(lower, ".wav"):
		return "wav"
	case strings.HasSuffix(lower, ".aac"), strings.HasSuffix(lower, ".m4a"):
		return "aac"
	case strings.HasSuffix(lower, ".ogg"):
		return "ogg"
	default:
		return "unknown"
	}
}
