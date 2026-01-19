package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ListTags returns a list of all user tags
func (h *Handlers) ListTags(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	tags, err := h.services.Tag.ListTags(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, tags)
}

// CreateTag creates a new tag
func (h *Handlers) CreateTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.CreateTagRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	tag, err := h.services.Tag.CreateTag(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return created(c, tag)
}

// GetTag returns a single tag by name
func (h *Handlers) GetTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	tagName := c.Param("name")
	if tagName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	tag, err := h.services.Tag.GetTag(c.Request().Context(), userID, tagName)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, tag)
}

// UpdateTag updates a tag
func (h *Handlers) UpdateTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	tagName := c.Param("name")
	if tagName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	var req models.UpdateTagRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	tag, err := h.services.Tag.UpdateTag(c.Request().Context(), userID, tagName, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, tag)
}

// DeleteTag deletes a tag
func (h *Handlers) DeleteTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	tagName := c.Param("name")
	if tagName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	if err := h.services.Tag.DeleteTag(c.Request().Context(), userID, tagName); err != nil {
		return handleError(c, err)
	}

	return noContent(c)
}

// GetTracksByTag returns tracks with a specific tag
func (h *Handlers) GetTracksByTag(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	tagName := c.Param("name")
	if tagName == "" {
		return handleError(c, models.ErrBadRequest)
	}

	tracks, err := h.services.Tag.GetTracksByTag(c.Request().Context(), userID, tagName)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, tracks)
}
