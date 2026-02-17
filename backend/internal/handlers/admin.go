package handlers

import (
	"net/http"

	"github.com/gvasels/personal-music-searchengine/internal/handlers/middleware"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// AdminHandler handles admin management endpoints.
type AdminHandler struct {
	adminService service.AdminService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// SearchUsers handles GET /api/v1/admin/users?search=query&limit=20
// Admin only - searches for users by email or display name.
func (h *AdminHandler) SearchUsers(c echo.Context) error {
	var req models.AdminSearchUsersRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	if req.Search == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.NewValidationError("search query parameter is required"),
		))
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, err := h.adminService.SearchUsers(c.Request().Context(), req.Search, limit)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, models.AdminSearchUsersResponse{
		Items: users,
	})
}

// GetUserDetails handles GET /api/v1/admin/users/:id
// Admin only - gets full user details including content counts.
func (h *AdminHandler) GetUserDetails(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	details, err := h.adminService.GetUserDetails(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, details)
}

// UpdateUserRole handles PUT /api/v1/admin/users/:id/role
// Admin only - updates a user's role in both DynamoDB and Cognito.
func (h *AdminHandler) UpdateUserRole(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Get admin's own user ID to prevent self-modification
	adminID := middleware.GetUserID(c)
	if adminID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.UpdateRoleRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	// Validate and convert role string
	newRole, valid := models.ValidateRole(req.Role)
	if !valid {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			models.NewValidationError("invalid role"),
		))
	}

	// Use the admin-aware version that prevents self-modification
	err := h.adminService.UpdateUserRoleByAdmin(c.Request().Context(), adminID, userID, newRole)
	if err != nil {
		return handleError(c, err)
	}

	// Return updated user details so frontend cache stays consistent
	details, err := h.adminService.GetUserDetails(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, details)
}

// UpdateUserStatus handles PUT /api/v1/admin/users/:id/status
// Admin only - enables or disables a user in both DynamoDB and Cognito.
func (h *AdminHandler) UpdateUserStatus(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	// Get admin's own user ID to prevent self-disabling
	adminID := middleware.GetUserID(c)
	if adminID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	// Prevent admin from disabling themselves
	if adminID == userID {
		return c.JSON(http.StatusForbidden, models.NewErrorResponse(
			models.NewForbiddenError("cannot modify your own status"),
		))
	}

	var req models.UpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	err := h.adminService.SetUserStatus(c.Request().Context(), userID, req.Disabled)
	if err != nil {
		return handleError(c, err)
	}

	// Return updated user details so frontend cache stays consistent
	details, err := h.adminService.GetUserDetails(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, details)
}

// GetSelfSignUp handles GET /api/v1/admin/settings/self-signup
func (h *AdminHandler) GetSelfSignUp(c echo.Context) error {
	enabled, err := h.adminService.GetSelfSignUpEnabled(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]bool{"enabled": enabled})
}

// SetSelfSignUp handles PUT /api/v1/admin/settings/self-signup
func (h *AdminHandler) SetSelfSignUp(c echo.Context) error {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.adminService.SetSelfSignUpEnabled(c.Request().Context(), req.Enabled); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]bool{"enabled": req.Enabled})
}

// Vector backend mode: "s3vectors", "s3flat", or "auto" (default).
// Stored in-memory — resets to "auto" on cold start.
var vectorBackendMode = "auto"

// GetVectorBackend handles GET /api/v1/admin/settings/vector-backend
func (h *AdminHandler) GetVectorBackend(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"mode": vectorBackendMode,
		"options": []map[string]string{
			{"value": "auto", "label": "Auto", "description": "S3 Vectors for kNN queries, S3 flat files as backup. Falls back gracefully if either is unavailable."},
			{"value": "s3vectors", "label": "S3 Vectors", "description": "AWS managed vector database. Best for kNN similarity search with cosine distance. Scales independently of Lambda memory."},
			{"value": "s3flat", "label": "S3 Flat Files", "description": "Vectors stored as JSON in S3. No kNN support — useful for debugging, data export, and native S3 testing."},
		},
	})
}

// SetVectorBackend handles PUT /api/v1/admin/settings/vector-backend
func (h *AdminHandler) SetVectorBackend(c echo.Context) error {
	var req struct {
		Mode string `json:"mode"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	switch req.Mode {
	case "auto", "s3vectors", "s3flat":
		vectorBackendMode = req.Mode
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mode must be auto, s3vectors, or s3flat"})
	}
	return c.JSON(http.StatusOK, map[string]string{"mode": vectorBackendMode})
}

// GetVectorBackendMode returns the current vector backend mode (used by services).
func GetVectorBackendMode() string {
	return vectorBackendMode
}
