package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// SearchSimple performs a simple text search
func (h *Handlers) SearchSimple(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, models.NewValidationError(map[string]string{
			"q": "search query is required",
		}))
	}

	// Build search request from query params
	req := models.SearchRequest{
		Query: query,
		Limit: 50,
	}

	result, err := h.searchService.Search(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// SearchAdvanced performs an advanced search with filters
func (h *Handlers) SearchAdvanced(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	if req.Query == "" {
		return c.JSON(http.StatusBadRequest, models.NewValidationError(map[string]string{
			"query": "search query is required",
		}))
	}

	// Set defaults
	if req.Limit == 0 || req.Limit > 100 {
		req.Limit = 50
	}

	result, err := h.searchService.Search(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetSearchSuggestions returns autocomplete suggestions
func (h *Handlers) GetSearchSuggestions(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	query := c.QueryParam("q")
	if query == "" || len(query) < 2 {
		return c.JSON(http.StatusOK, models.AutocompleteResponse{
			Query:       query,
			Suggestions: []models.SearchSuggestion{},
		})
	}

	result, err := h.searchService.GetSuggestions(c.Request().Context(), userID, query)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
