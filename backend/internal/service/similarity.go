package service

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/search"
)

// SimilarityOptions configures the similar tracks search.
type SimilarityOptions struct {
	Limit            int     `json:"limit"`            // Maximum number of similar tracks to return
	Mode             string  `json:"mode"`             // "semantic", "features", "combined"
	MinSimilarity    float64 `json:"minSimilarity"`    // Minimum similarity score (0.0-1.0)
	IncludeSameAlbum bool    `json:"includeSameAlbum"` // Whether to include tracks from same album
}

// DefaultSimilarityOptions returns sensible defaults.
func DefaultSimilarityOptions() SimilarityOptions {
	return SimilarityOptions{
		Limit:            10,
		Mode:             "combined",
		MinSimilarity:    0.5,
		IncludeSameAlbum: true,
	}
}

// MixingOptions configures the DJ-compatible tracks search.
type MixingOptions struct {
	Limit        int    `json:"limit"`        // Maximum number of mixable tracks to return
	BPMTolerance int    `json:"bpmTolerance"` // BPM tolerance (default 5)
	KeyMode      string `json:"keyMode"`      // "exact", "harmonic", "any"
}

// DefaultMixingOptions returns sensible defaults for DJ mixing.
func DefaultMixingOptions() MixingOptions {
	return MixingOptions{
		Limit:        10,
		BPMTolerance: 5,
		KeyMode:      "harmonic",
	}
}

// SimilarTrack represents a track with similarity information.
type SimilarTrack struct {
	Track         models.TrackResponse `json:"track"`
	Similarity    float64              `json:"similarity"`
	BPMDiff       int                  `json:"bpmDiff"`
	KeyCompatible bool                 `json:"keyCompatible"`
	MatchReasons  []string             `json:"matchReasons"`
}

// SimilarTracksResponse contains similar tracks search results.
type SimilarTracksResponse struct {
	SourceTrack  models.TrackResponse `json:"sourceTrack"`
	Similar      []SimilarTrack       `json:"similar"`
	TotalMatches int                  `json:"totalMatches"`
}

// MixableTrack represents a DJ-compatible track.
type MixableTrack struct {
	Track         models.TrackResponse `json:"track"`
	BPMDiff       int                  `json:"bpmDiff"`
	KeyTransition string               `json:"keyTransition"`
	MixScore      float64              `json:"mixScore"`
}

// MixableTracksResponse contains DJ-mixable tracks search results.
type MixableTracksResponse struct {
	SourceTrack models.TrackResponse `json:"sourceTrack"`
	Mixable     []MixableTrack       `json:"mixable"`
}

// SimilarityService finds similar and mixable tracks.
type SimilarityService struct {
	searchClient     *search.Client
	repo             repository.Repository
	embeddingService *EmbeddingService
}

// NewSimilarityService creates a new SimilarityService.
func NewSimilarityService(
	searchClient *search.Client,
	repo repository.Repository,
	embeddingService *EmbeddingService,
) *SimilarityService {
	return &SimilarityService{
		searchClient:     searchClient,
		repo:             repo,
		embeddingService: embeddingService,
	}
}

// FindSimilarTracks finds tracks similar to the given track.
func (s *SimilarityService) FindSimilarTracks(
	ctx context.Context,
	userID, trackID string,
	opts SimilarityOptions,
) (*SimilarTracksResponse, error) {
	// Get the source track
	sourceTrack, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source track: %w", err)
	}

	// Set defaults
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.MinSimilarity <= 0 {
		opts.MinSimilarity = 0.5
	}

	// Get all user tracks for comparison
	// In a production system, this would use vector search
	allTracks, err := s.getAllUserTracks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tracks: %w", err)
	}

	// Calculate similarity for each track
	var candidates []SimilarTrack
	for _, track := range allTracks {
		// Skip the source track
		if track.ID == sourceTrack.ID {
			continue
		}

		// Skip same album if not wanted
		if !opts.IncludeSameAlbum && track.Album == sourceTrack.Album && track.Album != "" {
			continue
		}

		// Calculate similarity based on mode
		var similarity float64
		var matchReasons []string

		switch opts.Mode {
		case "semantic":
			similarity, matchReasons = s.calculateSemanticSimilarity(sourceTrack, &track)
		case "features":
			similarity, matchReasons = s.calculateFeatureSimilarity(sourceTrack, &track)
		default: // "combined"
			semanticSim, semanticReasons := s.calculateSemanticSimilarity(sourceTrack, &track)
			featureSim, featureReasons := s.calculateFeatureSimilarity(sourceTrack, &track)
			// Weight: 60% semantic, 40% features
			similarity = semanticSim*0.6 + featureSim*0.4
			matchReasons = append(matchReasons, semanticReasons...)
			matchReasons = append(matchReasons, featureReasons...)
		}

		if similarity >= opts.MinSimilarity {
			bpmDiff := 0
			if sourceTrack.BPM > 0 && track.BPM > 0 {
				bpmDiff = sourceTrack.BPM - track.BPM
				if bpmDiff < 0 {
					bpmDiff = -bpmDiff
				}
			}

			candidates = append(candidates, SimilarTrack{
				Track:         track.ToResponse(""),
				Similarity:    similarity,
				BPMDiff:       bpmDiff,
				KeyCompatible: IsKeyCompatible(sourceTrack.KeyCamelot, track.KeyCamelot),
				MatchReasons:  matchReasons,
			})
		}
	}

	// Sort by similarity (descending)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Similarity > candidates[j].Similarity
	})

	// Limit results
	if len(candidates) > opts.Limit {
		candidates = candidates[:opts.Limit]
	}

	return &SimilarTracksResponse{
		SourceTrack:  sourceTrack.ToResponse(""),
		Similar:      candidates,
		TotalMatches: len(candidates),
	}, nil
}

// FindMixableTracks finds tracks that can be DJ-mixed with the given track.
func (s *SimilarityService) FindMixableTracks(
	ctx context.Context,
	userID, trackID string,
	opts MixingOptions,
) (*MixableTracksResponse, error) {
	// Get the source track
	sourceTrack, err := s.repo.GetTrack(ctx, userID, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source track: %w", err)
	}

	// Set defaults
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.BPMTolerance <= 0 {
		opts.BPMTolerance = 5
	}
	if opts.KeyMode == "" {
		opts.KeyMode = "harmonic"
	}

	// Get all user tracks for comparison
	allTracks, err := s.getAllUserTracks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tracks: %w", err)
	}

	// Filter for mixable tracks
	var candidates []MixableTrack
	for _, track := range allTracks {
		// Skip the source track
		if track.ID == sourceTrack.ID {
			continue
		}

		// Check BPM compatibility
		bpmDiff, bpmCompatible := GetBPMCompatibility(sourceTrack.BPM, track.BPM, opts.BPMTolerance)
		if !bpmCompatible && sourceTrack.BPM > 0 && track.BPM > 0 {
			continue
		}

		// Check key compatibility
		keyCompatible := true
		keyTransition := ""
		if opts.KeyMode != "any" && sourceTrack.KeyCamelot != "" && track.KeyCamelot != "" {
			if opts.KeyMode == "exact" {
				keyCompatible = sourceTrack.KeyCamelot == track.KeyCamelot
				if keyCompatible {
					keyTransition = "Same Key"
				}
			} else { // "harmonic"
				keyCompatible = IsKeyCompatible(sourceTrack.KeyCamelot, track.KeyCamelot)
				keyTransition = GetKeyTransition(sourceTrack.KeyCamelot, track.KeyCamelot)
			}
		}

		if !keyCompatible {
			continue
		}

		// Calculate mix score (0.0-1.0)
		mixScore := s.calculateMixScore(sourceTrack, &track, bpmDiff)

		candidates = append(candidates, MixableTrack{
			Track:         track.ToResponse(""),
			BPMDiff:       bpmDiff,
			KeyTransition: keyTransition,
			MixScore:      mixScore,
		})
	}

	// Sort by mix score (descending)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MixScore > candidates[j].MixScore
	})

	// Limit results
	if len(candidates) > opts.Limit {
		candidates = candidates[:opts.Limit]
	}

	return &MixableTracksResponse{
		SourceTrack: sourceTrack.ToResponse(""),
		Mixable:     candidates,
	}, nil
}

// calculateSemanticSimilarity calculates similarity based on metadata text.
// In a full implementation, this would use actual vector embeddings.
func (s *SimilarityService) calculateSemanticSimilarity(track1, track2 *models.Track) (float64, []string) {
	var similarity float64
	var reasons []string

	// Same artist is a strong signal
	if track1.Artist != "" && track1.Artist == track2.Artist {
		similarity += 0.4
		reasons = append(reasons, "same artist")
	}

	// Same genre
	if track1.Genre != "" && track1.Genre == track2.Genre {
		similarity += 0.3
		reasons = append(reasons, "same genre")
	}

	// Overlapping tags
	tagOverlap := countOverlappingTags(track1.Tags, track2.Tags)
	if tagOverlap > 0 {
		tagScore := float64(tagOverlap) / float64(max(len(track1.Tags), len(track2.Tags)))
		similarity += tagScore * 0.3
		reasons = append(reasons, "shared tags")
	}

	// Cap at 1.0
	if similarity > 1.0 {
		similarity = 1.0
	}

	return similarity, reasons
}

// calculateFeatureSimilarity calculates similarity based on audio features.
func (s *SimilarityService) calculateFeatureSimilarity(track1, track2 *models.Track) (float64, []string) {
	var similarity float64
	var reasons []string
	featureCount := 0

	// BPM similarity
	if track1.BPM > 0 && track2.BPM > 0 {
		bpmDiff := track1.BPM - track2.BPM
		if bpmDiff < 0 {
			bpmDiff = -bpmDiff
		}
		if bpmDiff <= 5 {
			similarity += 0.5
			reasons = append(reasons, "similar BPM")
		} else if bpmDiff <= 10 {
			similarity += 0.3
		}
		featureCount++
	}

	// Key compatibility
	if track1.KeyCamelot != "" && track2.KeyCamelot != "" {
		if IsKeyCompatible(track1.KeyCamelot, track2.KeyCamelot) {
			similarity += 0.5
			reasons = append(reasons, "harmonic key")
		}
		featureCount++
	}

	// Normalize by feature count if any features were compared
	if featureCount > 0 {
		similarity = similarity / float64(featureCount) * 2 // Scale up since max is 1.0
		if similarity > 1.0 {
			similarity = 1.0
		}
	}

	return similarity, reasons
}

// calculateMixScore calculates how well two tracks can be mixed for DJing.
func (s *SimilarityService) calculateMixScore(track1, track2 *models.Track, bpmDiff int) float64 {
	var score float64

	// BPM closeness (50% weight)
	if track1.BPM > 0 && track2.BPM > 0 {
		bpmScore := 1.0 - (float64(bpmDiff) / 10.0)
		if bpmScore < 0 {
			bpmScore = 0
		}
		score += bpmScore * 0.5
	} else {
		score += 0.25 // Neutral score if BPM unknown
	}

	// Key compatibility (40% weight)
	if track1.KeyCamelot != "" && track2.KeyCamelot != "" {
		if track1.KeyCamelot == track2.KeyCamelot {
			score += 0.4 // Perfect key match
		} else if IsKeyCompatible(track1.KeyCamelot, track2.KeyCamelot) {
			score += 0.35 // Harmonic match
		}
	} else {
		score += 0.2 // Neutral score if key unknown
	}

	// Genre match (10% weight)
	if track1.Genre != "" && track1.Genre == track2.Genre {
		score += 0.1
	}

	return score
}

// getAllUserTracks fetches all tracks for a user.
// In production, this would use pagination and possibly caching.
func (s *SimilarityService) getAllUserTracks(ctx context.Context, userID string) ([]models.Track, error) {
	var allTracks []models.Track
	cursor := ""

	for {
		filter := models.TrackFilter{
			Limit:   100,
			LastKey: cursor,
		}

		result, err := s.repo.ListTracks(ctx, userID, filter)
		if err != nil {
			return nil, err
		}

		allTracks = append(allTracks, result.Items...)

		if !result.HasMore || result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}

	return allTracks, nil
}

// CosineSimilarity calculates the cosine similarity between two vectors.
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// countOverlappingTags counts how many tags appear in both lists.
func countOverlappingTags(tags1, tags2 []string) int {
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0
	}

	tagSet := make(map[string]bool)
	for _, tag := range tags1 {
		tagSet[tag] = true
	}

	count := 0
	for _, tag := range tags2 {
		if tagSet[tag] {
			count++
		}
	}

	return count
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
