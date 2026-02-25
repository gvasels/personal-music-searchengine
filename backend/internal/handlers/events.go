package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// eventsService defines the service interface expected by EventsHandler.
type eventsService interface {
	GetArtistEvents(ctx context.Context, artistName string) (*models.ArtistEventsResponse, error)
	SearchArtistEvents(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error)
}

// EventsHandler handles event-related HTTP requests.
type EventsHandler struct {
	service eventsService
}

// NewEventsHandler creates a new EventsHandler.
func NewEventsHandler(svc eventsService) *EventsHandler {
	return &EventsHandler{service: svc}
}

// GetArtistEvents handles GET /api/v1/artists/:name/events
func (h *EventsHandler) GetArtistEvents(c echo.Context) error {
	userID := getUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	artistName := c.Param("name")
	if artistName == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	resp, err := h.service.GetArtistEvents(c.Request().Context(), artistName)
	if err != nil {
		var apiErr *models.APIError
		if errors.As(err, &apiErr) {
			return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
		}
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
	}

	return c.JSON(http.StatusOK, resp)
}

// SearchArtistEvents handles GET /api/v1/events/search
func (h *EventsHandler) SearchArtistEvents(c echo.Context) error {
	userID := getUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.ErrUnauthorized))
	}

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.ErrBadRequest))
	}

	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err == nil && parsed > 0 {
			limit = parsed
		}
	}

	results, err := h.service.SearchArtistEvents(c.Request().Context(), query, limit)
	if err != nil {
		var apiErr *models.APIError
		if errors.As(err, &apiErr) {
			return c.JSON(apiErr.StatusCode, models.NewErrorResponse(apiErr))
		}
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.ErrInternalServer))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": results,
		"total": len(results),
	})
}

// RegisterEventsRoutes registers event routes on the Echo instance.
func RegisterEventsRoutes(e *echo.Echo, h *EventsHandler) {
	api := e.Group("/api/v1")
	api.GET("/artists/:name/events", h.GetArtistEvents)
	api.GET("/events/search", h.SearchArtistEvents)
}

// EventsServiceAdapter adapts service.EventsService to the handler's eventsService interface.
type EventsServiceAdapter struct {
	svc *service.EventsService
}

// NewEventsServiceAdapter creates a new EventsServiceAdapter.
func NewEventsServiceAdapter(svc *service.EventsService) *EventsServiceAdapter {
	return &EventsServiceAdapter{svc: svc}
}

// GetArtistEvents adapts the service's GetArtistEvents (returns []Event) to the handler's (returns *ArtistEventsResponse).
func (a *EventsServiceAdapter) GetArtistEvents(ctx context.Context, artistName string) (*models.ArtistEventsResponse, error) {
	events, err := a.svc.GetArtistEvents(ctx, artistName)
	if err != nil {
		return nil, err
	}

	source := ""
	if len(events) > 0 {
		source = events[0].Source
	}

	return &models.ArtistEventsResponse{
		ArtistName: artistName,
		Events:     events,
		TotalCount: len(events),
		Source:     source,
	}, nil
}

// SearchArtistEvents adapts the service's SearchArtists to the handler's SearchArtistEvents.
func (a *EventsServiceAdapter) SearchArtistEvents(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error) {
	return a.svc.SearchArtists(ctx, query, limit)
}
