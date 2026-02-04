package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetStreamURL returns a signed URL for streaming a track
func (h *Handlers) GetStreamURL(c echo.Context) error {
	// Use DB role for real-time permission checking
	auth := h.getAuthContextWithDBRole(c)
	if auth.UserID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("trackId")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	resp, err := h.services.Stream.GetStreamURL(c.Request().Context(), auth.UserID, trackID, auth.HasGlobal)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// GetDownloadURL returns a signed URL for downloading a track
func (h *Handlers) GetDownloadURL(c echo.Context) error {
	// Use DB role for real-time permission checking
	auth := h.getAuthContextWithDBRole(c)
	if auth.UserID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	trackID := c.Param("trackId")
	if trackID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	resp, err := h.services.Stream.GetDownloadURL(c.Request().Context(), auth.UserID, trackID, auth.HasGlobal)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}
