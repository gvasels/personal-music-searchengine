package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// CreateArtist handles POST /api/v1/artists
func (h *Handlers) CreateArtist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.CreateArtistRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	artist, err := h.services.Artist.CreateArtist(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return created(c, artist)
}

// GetArtist handles GET /api/v1/artists/:id
func (h *Handlers) GetArtistByID(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	artistID := c.Param("id")
	if artistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	artist, err := h.services.Artist.GetArtist(c.Request().Context(), userID, artistID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, artist)
}

// UpdateArtist handles PUT /api/v1/artists/:id
func (h *Handlers) UpdateArtist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	artistID := c.Param("id")
	if artistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.UpdateArtistRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	artist, err := h.services.Artist.UpdateArtist(c.Request().Context(), userID, artistID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, artist)
}

// DeleteArtist handles DELETE /api/v1/artists/:id
func (h *Handlers) DeleteArtist(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	artistID := c.Param("id")
	if artistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	if err := h.services.Artist.DeleteArtist(c.Request().Context(), userID, artistID); err != nil {
		return handleError(c, err)
	}

	return noContent(c)
}

// ListArtistsEntity handles GET /api/v1/artists (new entity-based listing)
func (h *Handlers) ListArtistsEntity(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	filter := models.ArtistFilter{}

	// Parse query parameters
	if name := c.QueryParam("name"); name != "" {
		filter.Name = name
	}
	if sortBy := c.QueryParam("sortBy"); sortBy != "" {
		filter.SortBy = sortBy
	}
	if sortOrder := c.QueryParam("sortOrder"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}
	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	if lastKey := c.QueryParam("lastKey"); lastKey != "" {
		filter.LastKey = lastKey
	}

	result, err := h.services.Artist.ListArtists(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items":      result.Items,
		"nextCursor": result.NextCursor,
		"hasMore":    result.HasMore,
	})
}

// SearchArtists handles GET /api/v1/artists/search
func (h *Handlers) SearchArtists(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	query := c.QueryParam("q")
	if query == "" {
		return handleError(c, models.ErrBadRequest)
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	artists, err := h.services.Artist.SearchArtists(c.Request().Context(), userID, query, limit)
	if err != nil {
		return handleError(c, err)
	}

	return successList(c, artists)
}

// GetArtistTracksEntity handles GET /api/v1/artists/:id/tracks (entity-based)
func (h *Handlers) GetArtistTracksEntity(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	artistID := c.Param("id")
	if artistID == "" {
		return handleError(c, models.ErrBadRequest)
	}

	tracks, err := h.services.Artist.GetArtistTracks(c.Request().Context(), userID, artistID)
	if err != nil {
		return handleError(c, err)
	}

	return successList(c, tracks)
}
