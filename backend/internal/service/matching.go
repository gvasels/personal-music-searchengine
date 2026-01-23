package service

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// MatchingRepository defines the repository interface for matching
type MatchingRepository interface {
	ListTracksForMatching(ctx context.Context, userID string) ([]models.Track, error)
	GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error)
}

// MatchResult represents a track match with compatibility scores
type MatchResult struct {
	Track            *models.Track `json:"track"`
	BPMCompatibility float64       `json:"bpmCompatibility"`  // 0-1 score
	KeyCompatibility float64       `json:"keyCompatibility"`  // 0-1 score
	OverallScore     float64       `json:"overallScore"`      // Combined score
	BPMDiff          float64       `json:"bpmDiff"`           // Actual BPM difference
	KeyRelation      string        `json:"keyRelation"`       // Relation type (same, relative, neighbor, etc.)
}

// MatchingService handles BPM and key matching for DJ features
type MatchingService struct {
	repo       MatchingRepository
	featureSvc *FeatureService
}

// NewMatchingService creates a new matching service
func NewMatchingService(repo MatchingRepository, featureSvc *FeatureService) *MatchingService {
	return &MatchingService{
		repo:       repo,
		featureSvc: featureSvc,
	}
}

// FindCompatibleTracks finds tracks compatible with a source track
func (s *MatchingService) FindCompatibleTracks(ctx context.Context, userID, trackID string, limit int) ([]MatchResult, error) {
	// Check feature access for BPM matching
	bpmEnabled, err := s.featureSvc.IsEnabled(ctx, userID, models.FeatureBPMMatching)
	if err != nil {
		return nil, err
	}

	keyEnabled, err := s.featureSvc.IsEnabled(ctx, userID, models.FeatureKeyMatching)
	if err != nil {
		return nil, err
	}

	if !bpmEnabled && !keyEnabled {
		return nil, nil // No matching features enabled
	}

	// Get source track
	sourceTrack, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, err
	}
	if sourceTrack == nil {
		return nil, nil
	}

	// Get all tracks for matching
	tracks, err := s.repo.ListTracksForMatching(ctx, userID)
	if err != nil {
		return nil, err
	}

	var results []MatchResult

	for i := range tracks {
		track := &tracks[i]

		// Skip the source track
		if track.ID == trackID {
			continue
		}

		var bpmScore float64 = 0.5 // Neutral if not calculating
		var bpmDiff float64 = 0
		var keyScore float64 = 0.5 // Neutral if not calculating
		var keyRelation string = "unknown"

		// Calculate BPM compatibility
		if bpmEnabled && sourceTrack.BPM > 0 && track.BPM > 0 {
			bpmScore, bpmDiff = calculateBPMCompatibility(float64(sourceTrack.BPM), float64(track.BPM))
		}

		// Calculate key compatibility
		if keyEnabled && sourceTrack.MusicalKey != "" && track.MusicalKey != "" {
			keyScore, keyRelation = calculateKeyCompatibility(sourceTrack.MusicalKey, track.MusicalKey)
		}

		// Calculate overall score (weighted average)
		var overallScore float64
		if bpmEnabled && keyEnabled {
			overallScore = (bpmScore*0.5 + keyScore*0.5)
		} else if bpmEnabled {
			overallScore = bpmScore
		} else {
			overallScore = keyScore
		}

		// Only include tracks with reasonable compatibility
		if overallScore >= 0.3 {
			results = append(results, MatchResult{
				Track:            track,
				BPMCompatibility: bpmScore,
				KeyCompatibility: keyScore,
				OverallScore:     overallScore,
				BPMDiff:          bpmDiff,
				KeyRelation:      keyRelation,
			})
		}
	}

	// Sort by overall score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].OverallScore > results[j].OverallScore
	})

	// Apply limit
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// FindBPMCompatible finds tracks within a BPM range
func (s *MatchingService) FindBPMCompatible(ctx context.Context, userID string, targetBPM float64, tolerance float64, limit int) ([]MatchResult, error) {
	enabled, err := s.featureSvc.IsEnabled(ctx, userID, models.FeatureBPMMatching)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, nil
	}

	tracks, err := s.repo.ListTracksForMatching(ctx, userID)
	if err != nil {
		return nil, err
	}

	var results []MatchResult

	for i := range tracks {
		track := &tracks[i]
		if track.BPM <= 0 {
			continue
		}

		score, diff := calculateBPMCompatibility(targetBPM, float64(track.BPM))
		if score >= 0.5 {
			results = append(results, MatchResult{
				Track:            track,
				BPMCompatibility: score,
				KeyCompatibility: 1.0,
				OverallScore:     score,
				BPMDiff:          diff,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].OverallScore > results[j].OverallScore
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// FindKeyCompatible finds tracks with compatible musical keys
func (s *MatchingService) FindKeyCompatible(ctx context.Context, userID, targetKey string, limit int) ([]MatchResult, error) {
	enabled, err := s.featureSvc.IsEnabled(ctx, userID, models.FeatureKeyMatching)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, nil
	}

	tracks, err := s.repo.ListTracksForMatching(ctx, userID)
	if err != nil {
		return nil, err
	}

	var results []MatchResult

	for i := range tracks {
		track := &tracks[i]
		if track.MusicalKey == "" {
			continue
		}

		score, relation := calculateKeyCompatibility(targetKey, track.MusicalKey)
		if score >= 0.5 {
			results = append(results, MatchResult{
				Track:            track,
				BPMCompatibility: 1.0,
				KeyCompatibility: score,
				OverallScore:     score,
				KeyRelation:      relation,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].OverallScore > results[j].OverallScore
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// BPM Compatibility Calculation

// calculateBPMCompatibility calculates BPM compatibility score
// Considers double-time and half-time relationships
func calculateBPMCompatibility(sourceBPM, targetBPM float64) (float64, float64) {
	// Consider multiple BPM relationships
	candidates := []float64{
		targetBPM,
		targetBPM * 2,   // Double time
		targetBPM / 2,   // Half time
	}

	bestDiff := math.MaxFloat64
	for _, bpm := range candidates {
		diff := math.Abs(sourceBPM - bpm)
		if diff < bestDiff {
			bestDiff = diff
		}
	}

	// Calculate percentage difference
	percentDiff := (bestDiff / sourceBPM) * 100

	// Score based on ±3% range for full compatibility
	// 0% diff = 1.0 score
	// 3% diff = 0.7 score
	// 6% diff = 0.4 score
	// >10% diff = 0.0 score
	var score float64
	if percentDiff <= 3 {
		score = 1.0 - (percentDiff / 10)
	} else if percentDiff <= 6 {
		score = 0.7 - ((percentDiff - 3) / 10)
	} else if percentDiff <= 10 {
		score = 0.4 - ((percentDiff - 6) / 10)
	} else {
		score = 0
	}

	if score < 0 {
		score = 0
	}

	return score, bestDiff
}

// Camelot Wheel Key Compatibility

// CamelotPosition represents a position on the Camelot wheel
type CamelotPosition struct {
	Number int    // 1-12
	Mode   string // "A" for minor, "B" for major
}

// camelotWheel maps musical keys to Camelot positions
var camelotWheel = map[string]CamelotPosition{
	// Minor keys (A)
	"Am": {1, "A"}, "A minor": {1, "A"}, "Abm": {4, "A"}, "Ab minor": {4, "A"},
	"A#m": {8, "A"}, "A# minor": {8, "A"}, "Bbm": {8, "A"}, "Bb minor": {8, "A"},
	"Bm": {3, "A"}, "B minor": {3, "A"},
	"Cm": {10, "A"}, "C minor": {10, "A"},
	"C#m": {5, "A"}, "C# minor": {5, "A"}, "Dbm": {5, "A"}, "Db minor": {5, "A"},
	"Dm": {12, "A"}, "D minor": {12, "A"},
	"D#m": {7, "A"}, "D# minor": {7, "A"}, "Ebm": {7, "A"}, "Eb minor": {7, "A"},
	"Em": {2, "A"}, "E minor": {2, "A"},
	"Fm": {9, "A"}, "F minor": {9, "A"},
	"F#m": {4, "A"}, "F# minor": {4, "A"}, "Gbm": {4, "A"}, "Gb minor": {4, "A"},
	"Gm": {11, "A"}, "G minor": {11, "A"},
	"G#m": {6, "A"}, "G# minor": {6, "A"},

	// Major keys (B)
	"A": {4, "B"}, "A major": {4, "B"},
	"Ab": {9, "B"}, "Ab major": {9, "B"},
	"A#": {11, "B"}, "A# major": {11, "B"}, "Bb": {11, "B"}, "Bb major": {11, "B"},
	"B": {6, "B"}, "B major": {6, "B"},
	"C": {1, "B"}, "C major": {1, "B"},
	"C#": {8, "B"}, "C# major": {8, "B"}, "Db": {8, "B"}, "Db major": {8, "B"},
	"D": {3, "B"}, "D major": {3, "B"},
	"D#": {10, "B"}, "D# major": {10, "B"}, "Eb": {10, "B"}, "Eb major": {10, "B"},
	"E": {5, "B"}, "E major": {5, "B"},
	"F": {12, "B"}, "F major": {12, "B"},
	"F#": {7, "B"}, "F# major": {7, "B"}, "Gb": {7, "B"}, "Gb major": {7, "B"},
	"G": {2, "B"}, "G major": {2, "B"},
}

// calculateKeyCompatibility calculates key compatibility using the Camelot wheel
func calculateKeyCompatibility(sourceKey, targetKey string) (float64, string) {
	sourceKey = normalizeKey(sourceKey)
	targetKey = normalizeKey(targetKey)

	sourceCamelot, sourceOk := camelotWheel[sourceKey]
	targetCamelot, targetOk := camelotWheel[targetKey]

	if !sourceOk || !targetOk {
		return 0.5, "unknown" // Unknown keys get neutral score
	}

	// Same key - perfect match
	if sourceCamelot.Number == targetCamelot.Number && sourceCamelot.Mode == targetCamelot.Mode {
		return 1.0, "same"
	}

	// Relative major/minor (same number, different mode)
	if sourceCamelot.Number == targetCamelot.Number && sourceCamelot.Mode != targetCamelot.Mode {
		return 0.9, "relative"
	}

	// Adjacent on wheel (±1), same mode - perfect neighbors
	numDiff := camelotDistance(sourceCamelot.Number, targetCamelot.Number)
	if numDiff == 1 && sourceCamelot.Mode == targetCamelot.Mode {
		return 0.85, "neighbor"
	}

	// Two steps away, same mode - still compatible
	if numDiff == 2 && sourceCamelot.Mode == targetCamelot.Mode {
		return 0.6, "compatible"
	}

	// Energy boost/drop (adjacent number, different mode)
	if numDiff == 1 && sourceCamelot.Mode != targetCamelot.Mode {
		return 0.7, "energy_shift"
	}

	// Tritone/opposite (6 steps away)
	if numDiff == 6 && sourceCamelot.Mode == targetCamelot.Mode {
		return 0.4, "tritone"
	}

	return 0.3, "incompatible"
}

// normalizeKey normalizes key notation
func normalizeKey(key string) string {
	key = strings.TrimSpace(key)
	// Handle common abbreviations - only expand if not already the full word
	if !strings.Contains(key, "minor") {
		key = strings.ReplaceAll(key, "min", "minor")
	}
	if !strings.Contains(key, "major") {
		key = strings.ReplaceAll(key, "maj", "major")
	}
	return key
}

// camelotDistance calculates the circular distance on the Camelot wheel (1-12)
func camelotDistance(a, b int) int {
	diff := abs(a - b)
	if diff > 6 {
		diff = 12 - diff
	}
	return diff
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
