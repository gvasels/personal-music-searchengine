package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetStreamURL returns a signed URL for streaming a track
func (h *Handlers) GetStreamURL(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("trackId")
	if trackID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	quality := c.QueryParam("quality")
	if quality == "" {
		quality = "original"
	}

	result, err := h.streamService.GetStreamURL(c.Request().Context(), userID, trackID, quality)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetDownloadURL returns a signed URL for downloading a track
func (h *Handlers) GetDownloadURL(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	trackID := c.Param("trackId")
	if trackID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	result, err := h.streamService.GetDownloadURL(c.Request().Context(), userID, trackID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// RecordPlayback records a playback event
func (h *Handlers) RecordPlayback(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.RecordPlayRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	if err := h.streamService.RecordPlayback(c.Request().Context(), userID, req); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetQueue returns the user's play queue
func (h *Handlers) GetQueue(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	queue, err := h.streamService.GetQueue(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, queue)
}

// UpdateQueue updates the user's play queue
func (h *Handlers) UpdateQueue(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.UpdateQueueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	queue, err := h.streamService.UpdateQueue(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, queue)
}

// QueueAction performs an action on the play queue
func (h *Handlers) QueueAction(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.QueueActionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	queue, err := h.streamService.QueueAction(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, queue)
}
