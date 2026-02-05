package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// HelloServiceInterface defines the service methods used by HelloHandler.
type HelloServiceInterface interface {
	Search(ctx context.Context, query string) ([]service.HelloTrack, error)
	Featured(ctx context.Context, limit int) ([]service.HelloTrack, error)
}

// HelloHandler handles hello-world API endpoints.
type HelloHandler struct {
	svc HelloServiceInterface
}

// NewHelloHandler creates a new HelloHandler.
func NewHelloHandler(svc HelloServiceInterface) *HelloHandler {
	return &HelloHandler{svc: svc}
}

// Health handles GET /api/v1/hello/health and returns service status.
func (h *HelloHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "hello",
	})
}

// Search handles GET /api/v1/hello/search?q=<query> and returns matching tracks.
func (h *HelloHandler) Search(c echo.Context) error {
	q := c.QueryParam("q")

	tracks, err := h.svc.Search(c.Request().Context(), q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": tracks,
		"total": len(tracks),
	})
}

// Featured handles GET /api/v1/hello/featured?limit=<N> and returns featured tracks.
func (h *HelloHandler) Featured(c echo.Context) error {
	limit := 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	tracks, err := h.svc.Featured(c.Request().Context(), limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": tracks,
		"total": len(tracks),
	})
}

// RegisterHelloRoutes registers hello-world routes on the Echo instance.
func RegisterHelloRoutes(e *echo.Echo, h *HelloHandler) {
	e.GET("/api/v1/hello/health", h.Health)
	e.GET("/api/v1/hello/search", h.Search)
	e.GET("/api/v1/hello/featured", h.Featured)
}
