package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// HelloSearchResponse is the API response for hello search
type HelloSearchResponse struct {
	Tracks interface{} `json:"tracks"`
	Total  int         `json:"total"`
	Query  string      `json:"query"`
}

// HelloHealth returns a health check response
func (h *Handlers) HelloHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HelloSearch searches seed tracks by query
func (h *Handlers) HelloSearch(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "query parameter 'q' is required",
		})
	}

	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	tracks, err := h.services.Hello.SearchTracks(c.Request().Context(), query, limit)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, HelloSearchResponse{
		Tracks: tracks,
		Total:  len(tracks),
		Query:  query,
	})
}

// HelloFeatured returns all featured/seed tracks
func (h *Handlers) HelloFeatured(c echo.Context) error {
	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	tracks, err := h.services.Hello.GetFeaturedTracks(c.Request().Context(), limit)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, HelloSearchResponse{
		Tracks: tracks,
		Total:  len(tracks),
		Query:  "",
	})
}
