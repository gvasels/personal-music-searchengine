package handlers

import (
	"net/http"

	"github.com/gvasels/personal-music-searchengine/internal/handlers/middleware"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// RoleHandler handles role management endpoints.
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GetUserRole handles GET /api/v1/users/:id/role
// Admin only - gets the role for any user.
func (h *RoleHandler) GetUserRole(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	role, err := h.roleService.GetUserRole(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"userId": userID,
		"role":   role,
	})
}

// SetUserRoleRequest represents the request body for setting a user's role.
type SetUserRoleRequest struct {
	Role models.UserRole `json:"role" validate:"required"`
}

// SetUserRole handles PUT /api/v1/users/:id/role
// Admin only - sets the role for any user.
func (h *RoleHandler) SetUserRole(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req SetUserRoleRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	err := h.roleService.SetUserRole(c.Request().Context(), userID, req.Role)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"userId": userID,
		"role":   req.Role,
	})
}

// GetMyRole handles GET /api/v1/me/role
// Returns the role of the authenticated user.
func (h *RoleHandler) GetMyRole(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	role, err := h.roleService.GetUserRole(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"userId": userID,
		"role":   role,
	})
}

// ListUsersByRole handles GET /api/v1/admin/users
// Admin only - lists users by role with pagination.
func (h *RoleHandler) ListUsersByRole(c echo.Context) error {
	roleStr := c.QueryParam("role")
	if roleStr == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.NewValidationError("role query parameter is required"),
		))
	}

	role := models.UserRole(roleStr)
	if !role.IsValid() {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.NewValidationError("invalid role"),
		))
	}

	limit := 20
	cursor := c.QueryParam("cursor")

	result, err := h.roleService.ListUsersByRole(c.Request().Context(), role, limit, cursor)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
