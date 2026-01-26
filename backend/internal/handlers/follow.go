package handlers

import (
	"net/http"
	"strconv"

	"github.com/gvasels/personal-music-searchengine/internal/handlers/middleware"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// FollowHandler handles follow-related endpoints.
type FollowHandler struct {
	followService service.FollowService
}

// NewFollowHandler creates a new FollowHandler.
func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// Follow handles POST /api/v1/artist-profiles/:id/follow
// Subscriber+ - follows an artist.
func (h *FollowHandler) Follow(c echo.Context) error {
	followerUserID := middleware.GetUserID(c)
	if followerUserID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	followedUserID := c.Param("id")
	if followedUserID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	err := h.followService.Follow(c.Request().Context(), followerUserID, followedUserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"following": true,
		"followedId": followedUserID,
	})
}

// Unfollow handles DELETE /api/v1/artist-profiles/:id/follow
// Subscriber+ - unfollows an artist.
func (h *FollowHandler) Unfollow(c echo.Context) error {
	followerUserID := middleware.GetUserID(c)
	if followerUserID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	followedUserID := c.Param("id")
	if followedUserID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	err := h.followService.Unfollow(c.Request().Context(), followerUserID, followedUserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetFollowers handles GET /api/v1/artist-profiles/:id/followers
// Public - lists followers of an artist.
func (h *FollowHandler) GetFollowers(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	limit := 20
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	cursor := c.QueryParam("cursor")

	result, err := h.followService.GetFollowers(c.Request().Context(), userID, limit, cursor)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// GetFollowing handles GET /api/v1/users/me/following
// Subscriber+ - lists artists the user is following.
func (h *FollowHandler) GetFollowing(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	limit := 20
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	cursor := c.QueryParam("cursor")

	result, err := h.followService.GetFollowing(c.Request().Context(), userID, limit, cursor)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// IsFollowing handles GET /api/v1/artist-profiles/:id/following-status
// Subscriber+ - checks if the user is following an artist.
func (h *FollowHandler) IsFollowing(c echo.Context) error {
	followerUserID := middleware.GetUserID(c)
	if followerUserID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	followedUserID := c.Param("id")
	if followedUserID == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	following, err := h.followService.IsFollowing(c.Request().Context(), followerUserID, followedUserID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"following": following,
		"followedId": followedUserID,
	})
}
