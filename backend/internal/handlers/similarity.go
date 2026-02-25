package handlers

import (
	"net/http"
	"sort"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// SimilarTrackResponse represents a similar track result with full track data
type SimilarTrackResponse struct {
	Track         models.TrackResponse `json:"track"`
	Score         float32              `json:"score"`
	KeyCompatible bool                 `json:"keyCompatible"`
}

// GetSimilarTracks returns tracks similar to the given track.
// Uses vector similarity when embeddings exist, falls back to metadata-based
// similarity (genre, BPM, Camelot key) otherwise.
func (h *Handlers) GetSimilarTracks(c echo.Context) error {
	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	userID := getUserIDFromContext(c)

	empty := map[string]interface{}{"similar": []SimilarTrackResponse{}}

	// Try vector-based similarity first
	if h.services.Vector != nil {
		embedding, err := h.services.Vector.GetVector(c.Request().Context(), trackID)
		if err == nil && embedding != nil {
			return h.similarFromVectors(c, userID, trackID, embedding)
		}
	}

	// Fall back to metadata-based similarity
	if h.services.Similarity == nil {
		return c.JSON(http.StatusOK, empty)
	}

	resp, err := h.services.Similarity.FindSimilarTracks(
		c.Request().Context(), userID, trackID,
		service.SimilarityOptions{
			Limit:            20,
			Mode:             "combined",
			MinSimilarity:    0.1,
			IncludeSameAlbum: true,
		},
	)
	if err != nil {
		return c.JSON(http.StatusOK, empty)
	}

	similar := make([]SimilarTrackResponse, 0, len(resp.Similar))
	for _, s := range resp.Similar {
		similar = append(similar, SimilarTrackResponse{
			Track:         s.Track,
			Score:         float32(s.Similarity),
			KeyCompatible: s.KeyCompatible,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"similar": similar})
}

// similarFromVectors handles the vector-based path with Camelot key re-ranking.
func (h *Handlers) similarFromVectors(c echo.Context, userID, trackID string, embedding []float32) error {
	empty := map[string]interface{}{"similar": []SimilarTrackResponse{}}

	sourceTrack, err := h.services.Track.GetTrack(c.Request().Context(), userID, trackID, false)
	if err != nil {
		return c.JSON(http.StatusOK, empty)
	}

	// Over-fetch candidates so we can re-rank with key compatibility
	results, err := h.services.Vector.QuerySimilar(c.Request().Context(), embedding, 40)
	if err != nil {
		return handleError(c, err)
	}

	similar := make([]SimilarTrackResponse, 0, len(results))
	for _, r := range results {
		if r.ID == trackID {
			continue
		}
		track, err := h.services.Track.GetTrack(c.Request().Context(), userID, r.ID, false)
		if err != nil {
			continue
		}

		score := r.Score
		keyCompat := false
		if sourceTrack.KeyCamelot != "" && track.KeyCamelot != "" {
			if service.IsKeyCompatible(sourceTrack.KeyCamelot, track.KeyCamelot) {
				keyCompat = true
				if sourceTrack.KeyCamelot == track.KeyCamelot {
					score *= 1.15
				} else {
					score *= 1.10
				}
			}
		}

		similar = append(similar, SimilarTrackResponse{
			Track:         *track,
			Score:         score,
			KeyCompatible: keyCompat,
		})
	}

	sort.Slice(similar, func(i, j int) bool {
		return similar[i].Score > similar[j].Score
	})

	if len(similar) > 20 {
		similar = similar[:20]
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"similar": similar})
}
