package handlers

import (
	"net/http"

	"github.com/gvasels/personal-music-searchengine/internal/handlers/middleware"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// ArtistProfileHandler handles artist profile endpoints.
type ArtistProfileHandler struct {
	profileService service.ArtistProfileService
}

// NewArtistProfileHandler creates a new ArtistProfileHandler.
func NewArtistProfileHandler(profileService service.ArtistProfileService) *ArtistProfileHandler {
	return &ArtistProfileHandler{profileService: profileService}
}

// CreateProfile handles POST /api/v1/artist-profiles
// Requires artist role - creates a new artist profile.
func (h *ArtistProfileHandler) CreateProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	var req models.CreateArtistProfileRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	profile, err := h.profileService.CreateProfile(c.Request().Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, profile)
}

// GetProfile handles GET /api/v1/artist-profiles/:id
// Public - retrieves an artist profile by user ID.
func (h *ArtistProfileHandler) GetProfile(c echo.Context) error {
	profileUserID := c.Param("id")
	if profileUserID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	profile, err := h.profileService.GetProfile(c.Request().Context(), profileUserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles PUT /api/v1/artist-profiles/:id
// Owner only - updates an artist profile.
func (h *ArtistProfileHandler) UpdateProfile(c echo.Context) error {
	requestingUserID := middleware.GetUserID(c)
	if requestingUserID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	profileUserID := c.Param("id")
	if profileUserID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	var req models.UpdateArtistProfileRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	profile, err := h.profileService.UpdateProfile(c.Request().Context(), requestingUserID, profileUserID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, profile)
}

// DeleteProfile handles DELETE /api/v1/artist-profiles/:id
// Owner only - deletes an artist profile.
func (h *ArtistProfileHandler) DeleteProfile(c echo.Context) error {
	requestingUserID := middleware.GetUserID(c)
	if requestingUserID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	profileUserID := c.Param("id")
	if profileUserID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	err := h.profileService.DeleteProfile(c.Request().Context(), requestingUserID, profileUserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListProfiles handles GET /api/v1/artist-profiles
// Public - lists artist profiles for discovery.
func (h *ArtistProfileHandler) ListProfiles(c echo.Context) error {
	limit := 20
	cursor := c.QueryParam("cursor")

	result, err := h.profileService.ListProfiles(c.Request().Context(), limit, cursor)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetMyProfile handles GET /api/v1/artist-profiles/me
// Authenticated - retrieves the authenticated user's artist profile.
func (h *ArtistProfileHandler) GetMyProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	profile, err := h.profileService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, profile)
}
