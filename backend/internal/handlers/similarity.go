package handlers

import (
	"net/http"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/labstack/echo/v4"
)

// SimilarTrackResponse represents a similar track result with full track data
type SimilarTrackResponse struct {
	Track models.TrackResponse `json:"track"`
	Score float32              `json:"score"`
}

// GetSimilarTracks returns tracks similar to the given track
func (h *Handlers) GetSimilarTracks(c echo.Context) error {
	trackID := c.Param("id")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	userID := getUserIDFromContext(c)

	if h.services.Vector == nil {
		return handleError(c, models.ErrNotFound)
	}

	// Get the track's embedding
	embedding, err := h.services.Vector.GetVector(c.Request().Context(), trackID)
	if err != nil {
		return handleError(c, err)
	}
	if embedding == nil {
		return handleError(c, models.ErrNotFound)
	}

	// Query for similar tracks
	results, err := h.services.Vector.QuerySimilar(c.Request().Context(), embedding, 10)
	if err != nil {
		return handleError(c, err)
	}

	// Filter out the source track, enrich with track data
	similar := make([]SimilarTrackResponse, 0, len(results))
	for _, r := range results {
		if r.ID == trackID {
			continue
		}
		track, err := h.services.Track.GetTrack(c.Request().Context(), userID, r.ID, false)
		if err != nil {
			continue // skip tracks we can't fetch
		}
		similar = append(similar, SimilarTrackResponse{
			Track: *track,
			Score: r.Score,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"similar": similar})
}
