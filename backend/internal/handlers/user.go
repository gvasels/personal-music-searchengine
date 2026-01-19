package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetProfile returns the current user's profile
func (h *Handlers) GetProfile(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	profile, err := h.services.User.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, profile)
}

// UpdateProfile updates the current user's profile
func (h *Handlers) UpdateProfile(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var req models.UpdateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	profile, err := h.services.User.UpdateProfile(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, profile)
}
