package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetCurrentUser returns the current authenticated user's profile
func (h *Handlers) GetCurrentUser(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	user, err := h.userService.GetUser(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateCurrentUser updates the current user's profile
func (h *Handlers) UpdateCurrentUser(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

// Helper functions

func getUserIDFromContext(c echo.Context) string {
	// Extract user ID from Cognito JWT claims
	// In Lambda, this comes from the API Gateway authorizer context
	if claims, ok := c.Get("claims").(map[string]interface{}); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub
		}
	}

	// Fallback for local development
	if userID := c.Request().Header.Get("X-User-ID"); userID != "" {
		return userID
	}

	return ""
}

func handleError(c echo.Context, err error) error {
	if apiErr, ok := err.(*models.APIError); ok {
		return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
	}
	return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
}
