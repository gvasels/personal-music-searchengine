package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListTracks returns a paginated list of tracks
func (h *Handlers) ListTracks(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var filter models.TrackFilter
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Set defaults
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	if filter.SortBy == "" {
		filter.SortBy = "createdAt"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	result, err := h.trackService.ListTracks(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetTrack returns a single track by ID
func (h *Handlers) GetTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("id")
	if trackID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	track, err := h.trackService.GetTrack(c.Request().Context(), userID, trackID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, track)
}

// UpdateTrack updates a track's metadata
func (h *Handlers) UpdateTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("id")
	if trackID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.UpdateTrackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	track, err := h.trackService.UpdateTrack(c.Request().Context(), userID, trackID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, track)
}

// DeleteTrack deletes a track
func (h *Handlers) DeleteTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("id")
	if trackID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	if err := h.trackService.DeleteTrack(c.Request().Context(), userID, trackID); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// AddTagsToTrack adds tags to a track
func (h *Handlers) AddTagsToTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("id")
	if trackID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.AddTagsToTrackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	track, err := h.trackService.AddTagsToTrack(c.Request().Context(), userID, trackID, req.Tags)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, track)
}

// RemoveTagFromTrack removes a tag from a track
func (h *Handlers) RemoveTagFromTrack(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("id")
	tagName := c.Param("tag")
	if trackID == "" || tagName == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	track, err := h.trackService.RemoveTagFromTrack(c.Request().Context(), userID, trackID, tagName)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, track)
}
