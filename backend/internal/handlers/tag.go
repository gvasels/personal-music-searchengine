package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListTags returns a list of user's tags
func (h *Handlers) ListTags(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var filter models.TagFilter
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Set defaults
	if filter.Limit == 0 || filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.SortBy == "" {
		filter.SortBy = "name"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	result, err := h.tagService.ListTags(c.Request().Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// CreateTag creates a new tag
func (h *Handlers) CreateTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	tag, err := h.tagService.CreateTag(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, tag.ToResponse())
}

// UpdateTag updates a tag
func (h *Handlers) UpdateTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	tagName := c.Param("name")
	if tagName == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.UpdateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	tag, err := h.tagService.UpdateTag(c.Request().Context(), userID, tagName, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, tag.ToResponse())
}

// DeleteTag deletes a tag
func (h *Handlers) DeleteTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	tagName := c.Param("name")
	if tagName == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	if err := h.tagService.DeleteTag(c.Request().Context(), userID, tagName); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
