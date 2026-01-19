package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListPlaylists returns a list of playlists
func (h *Handlers) ListPlaylists(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var filter models.PlaylistFilter
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Set defaults
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	if filter.SortBy == "" {
		filter.SortBy = "name"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	result, err := h.playlistService.ListPlaylists(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// CreatePlaylist creates a new playlist
func (h *Handlers) CreatePlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.CreatePlaylistRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	playlist, err := h.playlistService.CreatePlaylist(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, playlist)
}

// GetPlaylist returns a playlist with its tracks
func (h *Handlers) GetPlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	playlist, err := h.playlistService.GetPlaylistWithTracks(c.Request().Context(), userID, playlistID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, playlist)
}

// UpdatePlaylist updates a playlist
func (h *Handlers) UpdatePlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.UpdatePlaylistRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	playlist, err := h.playlistService.UpdatePlaylist(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, playlist)
}

// DeletePlaylist deletes a playlist
func (h *Handlers) DeletePlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	if err := h.playlistService.DeletePlaylist(c.Request().Context(), userID, playlistID); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// AddTracksToPlaylist adds tracks to a playlist
func (h *Handlers) AddTracksToPlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.AddTracksToPlaylistRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	playlist, err := h.playlistService.AddTracks(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, playlist)
}

// RemoveTracksFromPlaylist removes tracks from a playlist
func (h *Handlers) RemoveTracksFromPlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.RemoveTracksFromPlaylistRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	playlist, err := h.playlistService.RemoveTracks(c.Request().Context(), userID, playlistID, req.TrackIDs)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, playlist)
}

// ReorderPlaylistTracks reorders tracks in a playlist
func (h *Handlers) ReorderPlaylistTracks(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.ReorderPlaylistTracksRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	playlist, err := h.playlistService.ReorderTracks(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, playlist)
}
