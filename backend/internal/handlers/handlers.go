package handlers

import (
	"errors"
	"net/http"

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
}

// getUserIDFromContext extracts the user ID from the request context
// In production, this comes from the Cognito JWT authorizer via API Gateway
func getUserIDFromContext(c echo.Context) string {
	// Try to get from API Gateway V2 JWT authorizer claims (Lambda environment)
	if requestCtx, ok := core.GetAPIGatewayV2ContextFromContext(c.Request().Context()); ok {
		if requestCtx.Authorizer != nil && requestCtx.Authorizer.JWT != nil {
			if sub, exists := requestCtx.Authorizer.JWT.Claims["sub"]; exists {
				return sub
			}
		}
	}

	// Try to get from request header (for local development/testing)
	if userID := c.Request().Header.Get("X-User-ID"); userID != "" {
		return userID
	}

	return ""
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

// created returns a JSON response with 201 status
func created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, data)
}

// noContent returns a 204 No Content response
func noContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
