package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListAlbums returns a paginated list of albums
func (h *Handlers) ListAlbums(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var filter models.AlbumFilter
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Set defaults
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	if filter.SortBy == "" {
		filter.SortBy = "title"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	result, err := h.albumService.ListAlbums(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetAlbum returns a single album with its tracks
func (h *Handlers) GetAlbum(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	albumID := c.Param("id")
	if albumID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	album, err := h.albumService.GetAlbumWithTracks(c.Request().Context(), userID, albumID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, album)
}

// ListArtists returns a list of artists
func (h *Handlers) ListArtists(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var filter models.ArtistFilter
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

	result, err := h.albumService.ListArtists(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetArtist returns an artist with their albums and tracks
func (h *Handlers) GetArtist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	artistName := c.Param("name")
	if artistName == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	result, err := h.albumService.GetArtist(c.Request().Context(), userID, artistName)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
