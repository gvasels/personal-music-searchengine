package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// artistWatchService defines the service interface expected by ArtistWatchHandler.
type artistWatchService interface {
	WatchArtist(ctx context.Context, userID, artistName string) (*models.ArtistWatchResponse, error)
	UnwatchArtist(ctx context.Context, userID, artistName string) error
	GetWatchStatus(ctx context.Context, userID, artistName string) (bool, error)
	ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatchResponse], error)
}

// ArtistWatchHandler handles artist watch HTTP requests.
type ArtistWatchHandler struct {
	service artistWatchService
}

// NewArtistWatchHandler creates a new ArtistWatchHandler.
func NewArtistWatchHandler(svc artistWatchService) *ArtistWatchHandler {
	return &ArtistWatchHandler{service: svc}
}

// getUserID extracts user ID from context or header, returns empty string if not found.
func getUserID(c echo.Context) string {
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		userID = c.Request().Header.Get("X-User-ID")
	}
	return userID
}

// WatchArtist handles POST /api/v1/artists/:name/watch
func (h *ArtistWatchHandler) WatchArtist(c echo.Context) error {
	userID := getUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	artistName := c.Param("name")
	if artistName == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	resp, err := h.service.WatchArtist(c.Request().Context(), userID, artistName)
	if err != nil {
		var apiErr *models.APIError
		if errors.As(err, &apiErr) {
			return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
		}
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"artistName": resp.ArtistName,
		"watching":   true,
		"watchedAt":  resp.WatchedAt,
	})
}

// UnwatchArtist handles DELETE /api/v1/artists/:name/watch
func (h *ArtistWatchHandler) UnwatchArtist(c echo.Context) error {
	userID := getUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	artistName := c.Param("name")

	err := h.service.UnwatchArtist(c.Request().Context(), userID, artistName)
	if err != nil {
		var apiErr *models.APIError
		if errors.As(err, &apiErr) {
			return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
		}
		if errors.Is(err, models.ErrNotFound) {
			return c.JSON(http.StatusNotFound, models.NewErrorResponse(models.ErrNotFound))
		}
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
	}

	return c.NoContent(http.StatusNoContent)
}

// GetWatchStatus handles GET /api/v1/artists/:name/watch
func (h *ArtistWatchHandler) GetWatchStatus(c echo.Context) error {
	userID := getUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	artistName := c.Param("name")

	watching, err := h.service.GetWatchStatus(c.Request().Context(), userID, artistName)
	if err != nil {
		var apiErr *models.APIError
		if errors.As(err, &apiErr) {
			return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
		}
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"watching":   watching,
		"artistName": artistName,
	})
}

// ListWatchedArtists handles GET /api/v1/users/me/watched-artists
func (h *ArtistWatchHandler) ListWatchedArtists(c echo.Context) error {
	userID := getUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	limit := 20
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err == nil && parsed > 0 {
			limit = parsed
		}
	}

	cursor := c.QueryParam("cursor")

	result, err := h.service.ListWatchedArtists(c.Request().Context(), userID, limit, cursor)
	if err != nil {
		var apiErr *models.APIError
		if errors.As(err, &apiErr) {
			return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
		}
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
	}

	return c.JSON(http.StatusOK, result)
}

// RegisterArtistWatchRoutes registers artist watch routes on the Echo instance.
func RegisterArtistWatchRoutes(e *echo.Echo, h *ArtistWatchHandler) {
	api := e.Group("/api/v1")
	api.POST("/artists/:name/watch", h.WatchArtist)
	api.DELETE("/artists/:name/watch", h.UnwatchArtist)
	api.GET("/artists/:name/watch", h.GetWatchStatus)
	api.GET("/users/me/watched-artists", h.ListWatchedArtists)
}

// ArtistWatchServiceAdapter adapts service.ArtistWatchService to the handler's artistWatchService interface.
type ArtistWatchServiceAdapter struct {
	svc *service.ArtistWatchService
}

// NewArtistWatchServiceAdapter creates a new ArtistWatchServiceAdapter.
func NewArtistWatchServiceAdapter(svc *service.ArtistWatchService) *ArtistWatchServiceAdapter {
	return &ArtistWatchServiceAdapter{svc: svc}
}

// WatchArtist adapts the service's WatchArtist (returns error) to the handler's (returns *ArtistWatchResponse, error).
func (a *ArtistWatchServiceAdapter) WatchArtist(ctx context.Context, userID, artistName string) (*models.ArtistWatchResponse, error) {
	err := a.svc.WatchArtist(ctx, userID, artistName)
	if err != nil {
		return nil, err
	}
	watch := models.NewArtistWatch(userID, artistName)
	resp := watch.ToResponse()
	return &resp, nil
}

// UnwatchArtist delegates to the service.
func (a *ArtistWatchServiceAdapter) UnwatchArtist(ctx context.Context, userID, artistName string) error {
	return a.svc.UnwatchArtist(ctx, userID, artistName)
}

// GetWatchStatus adapts service's IsWatching to the handler's GetWatchStatus.
func (a *ArtistWatchServiceAdapter) GetWatchStatus(ctx context.Context, userID, artistName string) (bool, error) {
	return a.svc.IsWatching(ctx, userID, artistName)
}

// ListWatchedArtists converts ArtistWatch items to ArtistWatchResponse items.
func (a *ArtistWatchServiceAdapter) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatchResponse], error) {
	result, err := a.svc.ListWatchedArtists(ctx, userID, limit, cursor)
	if err != nil {
		return nil, err
	}

	responses := make([]models.ArtistWatchResponse, len(result.Items))
	for i, w := range result.Items {
		responses[i] = w.ToResponse()
	}

	return &repository.PaginatedResult[models.ArtistWatchResponse]{
		Items:      responses,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}
