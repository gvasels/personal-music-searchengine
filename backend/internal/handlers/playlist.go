package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListPlaylists returns a paginated list of playlists
func (h *Handlers) ListPlaylists(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var filter models.PlaylistFilter
	if err := c.Bind(&filter); err != nil {
		return handleError(c, models.ErrBadRequest)
	}

	playlists, err := h.services.Playlist.ListPlaylists(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, playlists)
}

// CreatePlaylist creates a new playlist
func (h *Handlers) CreatePlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.CreatePlaylistRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	playlist, err := h.services.Playlist.CreatePlaylist(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return created(c, playlist)
}

// GetPlaylist returns a single playlist by ID
func (h *Handlers) GetPlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	playlist, err := h.services.Playlist.GetPlaylist(c.Request().Context(), userID, playlistID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, playlist)
}

// UpdatePlaylist updates a playlist's details
func (h *Handlers) UpdatePlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.UpdatePlaylistRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	playlist, err := h.services.Playlist.UpdatePlaylist(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, playlist)
}

// DeletePlaylist deletes a playlist
func (h *Handlers) DeletePlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	if err := h.services.Playlist.DeletePlaylist(c.Request().Context(), userID, playlistID); err != nil {
		return handleError(c, err)
	}

	return noContent(c)
}

// AddTracksToPlaylist adds tracks to a playlist
func (h *Handlers) AddTracksToPlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.AddTracksToPlaylistRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	playlist, err := h.services.Playlist.AddTracks(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, playlist)
}

// RemoveTracksFromPlaylist removes tracks from a playlist
func (h *Handlers) RemoveTracksFromPlaylist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.RemoveTracksFromPlaylistRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	playlist, err := h.services.Playlist.RemoveTracks(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, playlist)
}

// ReorderPlaylistTracks reorders tracks within a playlist
func (h *Handlers) ReorderPlaylistTracks(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	playlistID := c.Param("id")
	if playlistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.ReorderPlaylistTracksRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	playlist, err := h.services.Playlist.ReorderTracks(c.Request().Context(), userID, playlistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, playlist)
}
