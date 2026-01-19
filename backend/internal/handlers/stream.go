package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetStreamURL returns a signed URL for streaming a track
func (h *Handlers) GetStreamURL(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("trackId")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	resp, err := h.services.Stream.GetStreamURL(c.Request().Context(), userID, trackID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// GetDownloadURL returns a signed URL for downloading a track
func (h *Handlers) GetDownloadURL(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("trackId")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	resp, err := h.services.Stream.GetDownloadURL(c.Request().Context(), userID, trackID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}
