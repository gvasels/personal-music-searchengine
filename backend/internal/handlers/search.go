package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// SimpleSearch performs a simple text search
func (h *Handlers) SimpleSearch(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	query := c.QueryParam("q")
	if query == "" {
		return handleError(c, models.NewValidationError("query parameter 'q' is required"))
	}

	// Build a simple search request from query params
	req := models.SearchRequest{
		Query:  query,
		Limit:  20, // Default limit
		Cursor: c.QueryParam("cursor"),
	}

	// Parse optional limit
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		// Let the service handle limit parsing/validation
		req.Limit = 20 // Keep default, service will validate
	}

	resp, err := h.services.Search.Search(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// AdvancedSearch performs an advanced search with filters
func (h *Handlers) AdvancedSearch(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.SearchRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.services.Search.Search(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}

// Autocomplete provides search suggestions for the query
func (h *Handlers) Autocomplete(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	query := c.QueryParam("q")
	if query == "" {
		return handleError(c, models.NewValidationError("query parameter 'q' is required"))
	}

	resp, err := h.services.Search.Autocomplete(c.Request().Context(), userID, query)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, resp)
}
