package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListAlbums returns a paginated list of albums
func (h *Handlers) ListAlbums(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var filter models.AlbumFilter
	if err := c.Bind(&filter); err != nil {
		return handleError(c, models.ErrBadRequest)
	}

	albums, err := h.services.Album.ListAlbums(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, albums)
}

// GetAlbum returns a single album by ID
func (h *Handlers) GetAlbum(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	albumID := c.Param("id")
	if albumID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	album, err := h.services.Album.GetAlbum(c.Request().Context(), userID, albumID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, album)
}

// ListArtists returns a list of artists with their track/album counts
func (h *Handlers) ListArtists(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var filter models.ArtistFilter
	if err := c.Bind(&filter); err != nil {
		return handleError(c, models.ErrBadRequest)
	}

	artists, err := h.services.Album.ListArtists(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return successList(c, artists)
}

// ListTracksByArtist returns tracks by a specific artist
func (h *Handlers) ListTracksByArtist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	artistName := c.Param("name")
	if artistName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	tracks, err := h.services.Track.ListTracksByArtist(c.Request().Context(), userID, artistName)
	if err != nil {
		return handleError(c, err)
	}

	return successList(c, tracks)
}

// ListAlbumsByArtist returns albums by a specific artist
func (h *Handlers) ListAlbumsByArtist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	artistName := c.Param("name")
	if artistName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	albums, err := h.services.Album.ListAlbumsByArtist(c.Request().Context(), userID, artistName)
	if err != nil {
		return handleError(c, err)
	}

	return successList(c, albums)
}
