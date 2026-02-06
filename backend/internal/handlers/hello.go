package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// helloServiceSearcher is a generic interface for searching tracks
type helloServiceSearcher[T any] interface {
	Search(ctx context.Context, query string) ([]T, error)
	Featured(ctx context.Context, limit int) ([]T, error)
}

// HelloHandler handles hello world API endpoints
// It uses a generic service interface to work with different track types
type HelloHandler[T any] struct {
	service helloServiceSearcher[T]
}

// NewHelloHandler creates a new HelloHandler with the given service
func NewHelloHandler[T any](svc helloServiceSearcher[T]) *HelloHandler[T] {
	return &HelloHandler[T]{
		service: svc,
	}
}

// Health returns the health status of the hello service
// GET /api/v1/hello/health
func (h *HelloHandler[T]) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "hello",
	})
}

// Search searches tracks by query parameter "q"
// GET /api/v1/hello/search?q=<query>
// CRITICAL CONTRACT: reads "q" parameter, NOT "query"
// Response uses "items" key, NOT "tracks"
func (h *HelloHandler[T]) Search(c echo.Context) error {
	// CRITICAL: Use "q" parameter, NOT "query"
	query := c.QueryParam("q")

	tracks, err := h.service.Search(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}

	// Ensure we return empty slice not nil
	if tracks == nil {
		tracks = make([]T, 0)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"items": tracks,
		"total": len(tracks),
	})
}

// Featured returns featured tracks with optional limit
// GET /api/v1/hello/featured?limit=<n>
// Default limit is 10
func (h *HelloHandler[T]) Featured(c echo.Context) error {
	// Default limit is 10
	limit := 10

	// Parse limit parameter if provided
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err == nil && parsed > 0 {
			limit = parsed
		}
		// Invalid or negative limit uses default (10)
	}

	tracks, err := h.service.Featured(c.Request().Context(), limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}

	// Ensure we return empty slice not nil
	if tracks == nil {
		tracks = make([]T, 0)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"items": tracks,
		"total": len(tracks),
	})
}

// RegisterHelloRoutes registers the hello routes with the Echo instance
// Routes are registered at /api/v1/hello/* without authentication middleware
func RegisterHelloRoutes[T any](e *echo.Echo, h *HelloHandler[T]) {
	hello := e.Group("/api/v1/hello")
	hello.GET("/health", h.Health)
	hello.GET("/search", h.Search)
	hello.GET("/featured", h.Featured)
}
