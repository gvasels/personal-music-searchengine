package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
)

// GetLibraryStats returns library statistics based on user role and requested scope
// GET /api/v1/stats?scope=own|public|all
func (h *Handlers) GetLibraryStats(c echo.Context) error {
	// Use DB role for real-time permission checking (role changes take effect immediately)
	authCtx := h.getAuthContextWithDBRole(c)
	if authCtx.UserID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	// Get scope from query parameter
	scopeParam := c.QueryParam("scope")
	scope := service.StatsScopeOwn // default

	switch scopeParam {
	case "all":
		scope = service.StatsScopeAll
	case "public":
		scope = service.StatsScopePublic
	case "own", "":
		scope = service.StatsScopeOwn
	default:
		return handleError(c, models.NewValidationError("invalid scope: must be 'own', 'public', or 'all'"))
	}

	stats, err := h.services.Track.GetLibraryStats(c.Request().Context(), authCtx.UserID, scope, authCtx.HasGlobal)
	if err != nil {
		return handleError(c, err)
	}

	return success(c, stats)
}

// FeaturesResponse represents the user's features based on their role
type FeaturesResponse struct {
	Tier     string          `json:"tier"`
	Role     string          `json:"role"`
	Features map[string]bool `json:"features"`
}

// GetFeatures returns the current user's features based on their role
// GET /api/v1/features
func (h *Handlers) GetFeatures(c echo.Context) error {
	authCtx := getAuthContext(c)
	if authCtx.UserID == "" {
		return handleError(c, models.ErrUnauthorized)
	}

	// Get role from DB for real-time permission checking (role changes take effect immediately)
	dbRole, err := h.services.User.GetUserRole(c.Request().Context(), authCtx.UserID)
	if err != nil {
		// Fall back to JWT groups if DB lookup fails
		dbRole = determineRoleFromGroups(authCtx.Groups)
	}
	role := dbRole

	// Map role to tier for backward compatibility
	tier := mapRoleToTier(role)

	// Get features based on role
	features := getFeaturesForRole(role)

	return success(c, FeaturesResponse{
		Tier:     tier,
		Role:     string(role),
		Features: features,
	})
}

// determineRoleFromGroups extracts the user role from Cognito groups
func determineRoleFromGroups(groups []string) models.UserRole {
	for _, g := range groups {
		lowerGroup := strings.ToLower(g)
		if lowerGroup == "admin" || lowerGroup == "admins" {
			return models.RoleAdmin
		}
	}
	for _, g := range groups {
		lowerGroup := strings.ToLower(g)
		if lowerGroup == "artist" || lowerGroup == "artists" {
			return models.RoleArtist
		}
	}
	for _, g := range groups {
		lowerGroup := strings.ToLower(g)
		if lowerGroup == "subscriber" || lowerGroup == "subscribers" {
			return models.RoleSubscriber
		}
	}
	return models.RoleSubscriber // Default role
}

// mapRoleToTier maps role to subscription tier for backward compatibility
func mapRoleToTier(role models.UserRole) string {
	switch role {
	case models.RoleAdmin:
		return "creator"
	case models.RoleArtist:
		return "pro"
	default:
		return "free"
	}
}

// getFeaturesForRole returns feature flags based on user role
func getFeaturesForRole(role models.UserRole) map[string]bool {
	// Base features for all authenticated users
	features := map[string]bool{
		"CRATES":            true,
		"PLAYLISTS":         true,
		"TAGS":              true,
		"SEARCH":            true,
		"HQ_STREAMING":      true,
	}

	// Artist features
	if role == models.RoleArtist || role == models.RoleAdmin {
		features["HOT_CUES"] = true
		features["BPM_MATCHING"] = true
		features["KEY_MATCHING"] = true
		features["WAVEFORMS"] = true
		features["BULK_EDIT"] = true
	}

	// Admin features (everything)
	if role == models.RoleAdmin {
		features["MIX_RECORDING"] = true
		features["ADVANCED_STATS"] = true
		features["API_ACCESS"] = true
		features["UNLIMITED_STORAGE"] = true
	}

	return features
}

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
