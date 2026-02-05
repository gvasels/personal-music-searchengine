package handlers

import (
	"net/http"
	"strconv"

	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// HelloHandler handles hello-world API endpoints
type HelloHandler struct {
	service *service.HelloService
}

// NewHelloHandler creates a new HelloHandler
func NewHelloHandler(svc *service.HelloService) *HelloHandler {
	return &HelloHandler{service: svc}
}

// HelloHealth returns the health status of the hello service
func (h *HelloHandler) HelloHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "hello",
	})
}

// HelloSearch searches tracks by query parameter "q"
func (h *HelloHandler) HelloSearch(c echo.Context) error {
	query := c.QueryParam("q")

	results, err := h.service.SearchTracks(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": results,
		"total": len(results),
	})
}

// HelloFeatured returns featured tracks with optional limit parameter
func (h *HelloHandler) HelloFeatured(c echo.Context) error {
	limit := 0
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid limit parameter",
			})
		}
	}

	results, err := h.service.GetFeaturedTracks(c.Request().Context(), limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": results,
		"total": len(results),
	})
}

// RegisterHelloRoutes registers hello-world routes on the Echo instance
func RegisterHelloRoutes(e *echo.Echo, h *HelloHandler) {
	e.GET("/api/v1/hello/health", h.HelloHealth)
	e.GET("/api/v1/hello/search", h.HelloSearch)
	e.GET("/api/v1/hello/featured", h.HelloFeatured)
}
