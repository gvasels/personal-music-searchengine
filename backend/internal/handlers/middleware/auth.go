package middleware

import (
	"net/http"
	"strings"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/labstack/echo/v4"
)

// Context keys for auth data
const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
	UserGroupsKey = "user_groups"
)

// RequireAuth middleware ensures the user is authenticated.
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, role, groups := extractAuthFromContext(c)

			if userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			// Store auth info in context for handlers
			c.Set(UserIDKey, userID)
			c.Set(UserRoleKey, role)
			c.Set(UserGroupsKey, groups)

			return next(c)
		}
	}
}

// OptionalAuth middleware extracts auth info if present but doesn't require it.
func OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, role, groups := extractAuthFromContext(c)

			// Store auth info in context (may be empty)
			c.Set(UserIDKey, userID)
			c.Set(UserRoleKey, role)
			c.Set(UserGroupsKey, groups)

			return next(c)
		}
	}
}

// RequireRole middleware ensures the user has a specific role (or admin).
func RequireRole(requiredRole models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, role, groups := extractAuthFromContext(c)

			if userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			// Store auth info in context
			c.Set(UserIDKey, userID)
			c.Set(UserRoleKey, role)
			c.Set(UserGroupsKey, groups)

			// Admin can do anything
			if role == models.RoleAdmin {
				return next(c)
			}

			// Check if user has the required role
			if role != requiredRole {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}

// RequirePermission middleware ensures the user has a specific permission.
func RequirePermission(permission models.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, role, groups := extractAuthFromContext(c)

			if userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}

			// Store auth info in context
			c.Set(UserIDKey, userID)
			c.Set(UserRoleKey, role)
			c.Set(UserGroupsKey, groups)

			// Check if role has the required permission
			if !role.HasPermission(permission) {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}

// extractAuthFromContext extracts user ID, role, and groups from the request.
func extractAuthFromContext(c echo.Context) (string, models.UserRole, []string) {
	var userID string
	var role models.UserRole = models.RoleGuest
	var groups []string

	// Try to get from API Gateway V2 JWT authorizer claims (Lambda environment)
	if requestCtx, ok := core.GetAPIGatewayV2ContextFromContext(c.Request().Context()); ok {
		if requestCtx.Authorizer != nil && requestCtx.Authorizer.JWT != nil {
			claims := requestCtx.Authorizer.JWT.Claims

			// Extract user ID (sub claim)
			if sub, exists := claims["sub"]; exists {
				userID = sub
			}

			// Extract Cognito groups from cognito:groups claim
			if groupsClaim, exists := claims["cognito:groups"]; exists {
				groups = parseGroups(groupsClaim)
				role = roleFromGroups(groups)
			}
		}
	}

	// Fall back to headers for local development/testing
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}

	// Get role from header if not set from JWT
	if role == models.RoleGuest {
		if roleHeader := c.Request().Header.Get("X-User-Role"); roleHeader != "" {
			role = models.UserRole(roleHeader)
			if !role.IsValid() {
				role = models.RoleGuest
			}
		}
	}

	return userID, role, groups
}

// parseGroups parses Cognito groups from JWT claim.
// API Gateway passes array claims as "[value1 value2]" format.
func parseGroups(groupsClaim string) []string {
	if groupsClaim == "" {
		return nil
	}

	// Strip brackets if present (API Gateway formats arrays as "[a b c]")
	groupsClaim = strings.TrimPrefix(groupsClaim, "[")
	groupsClaim = strings.TrimSuffix(groupsClaim, "]")
	groupsClaim = strings.TrimSpace(groupsClaim)

	if groupsClaim == "" {
		return nil
	}

	return strings.Split(groupsClaim, " ")
}

// roleFromGroups determines the highest-priority role from Cognito groups.
func roleFromGroups(groups []string) models.UserRole {
	// Priority: admin > artist > subscriber > guest
	hasAdmin := false
	hasArtist := false
	hasSubscriber := false

	for _, g := range groups {
		switch strings.ToLower(g) {
		case "admin", "admins":
			hasAdmin = true
		case "artist", "artists":
			hasArtist = true
		case "subscriber", "subscribers":
			hasSubscriber = true
		}
	}

	if hasAdmin {
		return models.RoleAdmin
	}
	if hasArtist {
		return models.RoleArtist
	}
	if hasSubscriber {
		return models.RoleSubscriber
	}

	return models.RoleGuest
}

// GetUserID retrieves the user ID from the Echo context.
func GetUserID(c echo.Context) string {
	if userID, ok := c.Get(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserRole retrieves the user role from the Echo context.
func GetUserRole(c echo.Context) models.UserRole {
	if role, ok := c.Get(UserRoleKey).(models.UserRole); ok {
		return role
	}
	return models.RoleGuest
}

// GetUserGroups retrieves the user's Cognito groups from the Echo context.
func GetUserGroups(c echo.Context) []string {
	if groups, ok := c.Get(UserGroupsKey).([]string); ok {
		return groups
	}
	return nil
}

// IsAuthenticated returns true if a user is authenticated.
func IsAuthenticated(c echo.Context) bool {
	return GetUserID(c) != ""
}

// HasPermission checks if the current user has a specific permission.
func HasPermission(c echo.Context, permission models.Permission) bool {
	return GetUserRole(c).HasPermission(permission)
}
