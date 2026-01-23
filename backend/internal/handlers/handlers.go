package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/labstack/echo/v4"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	services *service.Services
}

// NewHandlers creates a new Handlers instance
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		services: services,
	}
}

// RegisterRoutes registers all routes with the Echo instance
func (h *Handlers) RegisterRoutes(e *echo.Echo) {
	// API v1 group
	api := e.Group("/api/v1")

	// User routes
	api.GET("/me", h.GetProfile)
	api.PUT("/me", h.UpdateProfile)

	// Track routes
	api.GET("/tracks", h.ListTracks)
	api.GET("/tracks/:id", h.GetTrack)
	api.PUT("/tracks/:id", h.UpdateTrack)
	api.DELETE("/tracks/:id", h.DeleteTrack)
	api.POST("/tracks/:id/tags", h.AddTagsToTrack)
	api.DELETE("/tracks/:id/tags/:tag", h.RemoveTagFromTrack)
	api.PUT("/tracks/:id/cover", h.UploadCoverArt)

	// Album routes
	api.GET("/albums", h.ListAlbums)
	api.GET("/albums/:id", h.GetAlbum)

	// Artist routes
	api.GET("/artists", h.ListArtists)
	api.GET("/artists/:name", h.GetArtist)
	api.GET("/artists/:name/tracks", h.ListTracksByArtist)
	api.GET("/artists/:name/albums", h.ListAlbumsByArtist)

	// Playlist routes
	api.GET("/playlists", h.ListPlaylists)
	api.POST("/playlists", h.CreatePlaylist)
	api.GET("/playlists/:id", h.GetPlaylist)
	api.PUT("/playlists/:id", h.UpdatePlaylist)
	api.DELETE("/playlists/:id", h.DeletePlaylist)
	api.POST("/playlists/:id/tracks", h.AddTracksToPlaylist)
	api.DELETE("/playlists/:id/tracks", h.RemoveTracksFromPlaylist)

	// Tag routes
	api.GET("/tags", h.ListTags)
	api.POST("/tags", h.CreateTag)
	api.GET("/tags/:name", h.GetTag)
	api.PUT("/tags/:name", h.UpdateTag)
	api.DELETE("/tags/:name", h.DeleteTag)
	api.GET("/tags/:name/tracks", h.GetTracksByTag)

	// Upload routes
	api.POST("/upload/presigned", h.CreatePresignedUpload)
	api.POST("/upload/confirm", h.ConfirmUpload)
	api.POST("/upload/complete-multipart", h.CompleteMultipartUpload)
	api.GET("/uploads", h.ListUploads)
	api.GET("/uploads/:id", h.GetUploadStatus)
	api.POST("/uploads/:id/reprocess", h.ReprocessUpload)

	// Streaming routes
	api.GET("/stream/:trackId", h.GetStreamURL)
	api.GET("/download/:trackId", h.GetDownloadURL)

	// Search routes
	api.GET("/search", h.SimpleSearch)
	api.POST("/search", h.AdvancedSearch)
	api.GET("/search/autocomplete", h.Autocomplete)
}

// AuthContext contains user authentication and permission information
type AuthContext struct {
	UserID      string
	HasGlobal   bool   // True if user has GLOBAL read permission
	Groups      []string
}

// getAuthContext extracts user ID and permissions from the request context
// In production, this comes from the Cognito JWT authorizer via API Gateway
func getAuthContext(c echo.Context) AuthContext {
	ctx := AuthContext{}

	// Try to get from API Gateway V2 JWT authorizer claims (Lambda environment)
	if requestCtx, ok := core.GetAPIGatewayV2ContextFromContext(c.Request().Context()); ok {
		if requestCtx.Authorizer != nil && requestCtx.Authorizer.JWT != nil {
			claims := requestCtx.Authorizer.JWT.Claims

			// Extract user ID (sub claim)
			if sub, exists := claims["sub"]; exists {
				ctx.UserID = sub
			}

			// Extract Cognito groups from cognito:groups claim
			// Groups come as a space-separated string or array
			if groups, exists := claims["cognito:groups"]; exists {
				ctx.Groups = parseGroups(groups)
				ctx.HasGlobal = containsGlobalGroup(ctx.Groups)
			}
		}
	}

	// Fall back to headers for local development/testing
	if ctx.UserID == "" {
		ctx.UserID = c.Request().Header.Get("X-User-ID")
	}

	// Check for global permission header (for testing)
	if c.Request().Header.Get("X-Global-Access") == "true" {
		ctx.HasGlobal = true
	}

	return ctx
}

// parseGroups parses Cognito groups from JWT claim (can be string or []interface{})
func parseGroups(groupsClaim string) []string {
	if groupsClaim == "" {
		return nil
	}
	// Cognito groups come as space-separated in the JWT
	return strings.Split(groupsClaim, " ")
}

// containsGlobalGroup checks if user belongs to a group with global read access
func containsGlobalGroup(groups []string) bool {
	globalGroups := map[string]bool{
		"GlobalReaders": true,
		"Admin":         true,
		"Admins":        true,
	}
	for _, g := range groups {
		if globalGroups[g] {
			return true
		}
	}
	return false
}

// getUserIDFromContext extracts the user ID from the request context (legacy helper)
// In production, this comes from the Cognito JWT authorizer via API Gateway
func getUserIDFromContext(c echo.Context) string {
	return getAuthContext(c).UserID
}

// handleError converts errors to appropriate HTTP responses
func handleError(c echo.Context, err error) error {
	// Check for APIError
	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
	}

	// Default to internal server error
	return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
}

// bindAndValidate binds the request body and validates it
func bindAndValidate(c echo.Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		return models.ErrBadRequest
	}

	if err := c.Validate(v); err != nil {
		return models.NewValidationError(err.Error())
	}

	return nil
}

// success returns a JSON success response
func success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

// ListResponse wraps a slice in a list response with items array
type ListResponse[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

// successList returns a JSON success response for list endpoints
// This wraps plain slices in { items: [...], total: N } format
func successList[T any](c echo.Context, items []T) error {
	return c.JSON(http.StatusOK, ListResponse[T]{
		Items: items,
		Total: len(items),
	})
}

// created returns a JSON response with 201 status
func created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, data)
}

// noContent returns a 204 No Content response
func noContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
