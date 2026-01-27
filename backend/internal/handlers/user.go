package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
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

// GetSettings returns the current user's settings
// GET /api/v1/users/me/settings
func (h *Handlers) GetSettings(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	settings, err := h.services.User.GetSettings(c.Request().Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return handleError(c, models.NewNotFoundError("User", userID))
		}
		return handleError(c, err)
	}

	return success(c, settings)
}

// UpdateSettings performs a partial update of the current user's settings
// PATCH /api/v1/users/me/settings
func (h *Handlers) UpdateSettings(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	var input service.UserSettingsUpdateInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.NewAPIError("INVALID_REQUEST", "Invalid request body", http.StatusBadRequest),
		))
	}

	settings, err := h.services.User.UpdateSettings(c.Request().Context(), userID, &input)
	if err != nil {
		if err == service.ErrUserNotFound {
			return handleError(c, models.NewNotFoundError("User", userID))
		}
		if err == service.ErrValidation {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
				models.NewAPIError("VALIDATION_ERROR", "Invalid settings values", http.StatusBadRequest),
			))
		}
		return handleError(c, err)
	}

	return success(c, settings)
}
