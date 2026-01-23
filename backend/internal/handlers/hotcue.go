package handlers

import (
	"net/http"
	"strconv"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// HotCueHandler handles hot cue endpoints
type HotCueHandler struct {
	hotCueSvc *service.HotCueService
}

// NewHotCueHandler creates a new hot cue handler
func NewHotCueHandler(hotCueSvc *service.HotCueService) *HotCueHandler {
	return &HotCueHandler{hotCueSvc: hotCueSvc}
}

// SetHotCue sets or updates a hot cue
// PUT /tracks/:id/hotcues/:slot
func (h *HotCueHandler) SetHotCue(c echo.Context) error {
	userID := c.Get("userID").(string)
	trackID := c.Param("id")

	slot, err := strconv.Atoi(c.Param("slot"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid slot number")
	}

	var req models.SetHotCueRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	hotCue, err := h.hotCueSvc.SetHotCue(c.Request().Context(), userID, trackID, slot, req)
	if err != nil {
		if err.Error() == "track not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if err.Error() == "hot cues feature is not enabled for your subscription tier" {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, models.HotCueResponse{
		Slot:      hotCue.Slot,
		Position:  hotCue.Position,
		Label:     hotCue.Label,
		Color:     hotCue.Color,
		CreatedAt: hotCue.CreatedAt,
		UpdatedAt: hotCue.UpdatedAt,
	})
}

// DeleteHotCue removes a hot cue
// DELETE /tracks/:id/hotcues/:slot
func (h *HotCueHandler) DeleteHotCue(c echo.Context) error {
	userID := c.Get("userID").(string)
	trackID := c.Param("id")

	slot, err := strconv.Atoi(c.Param("slot"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid slot number")
	}

	if err := h.hotCueSvc.DeleteHotCue(c.Request().Context(), userID, trackID, slot); err != nil {
		if err.Error() == "track not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// GetHotCues returns all hot cues for a track
// GET /tracks/:id/hotcues
func (h *HotCueHandler) GetHotCues(c echo.Context) error {
	userID := c.Get("userID").(string)
	trackID := c.Param("id")

	hotCues, err := h.hotCueSvc.GetHotCues(c.Request().Context(), userID, trackID)
	if err != nil {
		if err.Error() == "track not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get hot cues")
	}

	return c.JSON(http.StatusOK, hotCues)
}

// ClearHotCues removes all hot cues from a track
// DELETE /tracks/:id/hotcues
func (h *HotCueHandler) ClearHotCues(c echo.Context) error {
	userID := c.Get("userID").(string)
	trackID := c.Param("id")

	if err := h.hotCueSvc.ClearAllHotCues(c.Request().Context(), userID, trackID); err != nil {
		if err.Error() == "track not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clear hot cues")
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterHotCueRoutes registers hot cue routes
func RegisterHotCueRoutes(g *echo.Group, h *HotCueHandler) {
	g.GET("/tracks/:id/hotcues", h.GetHotCues)
	g.PUT("/tracks/:id/hotcues/:slot", h.SetHotCue)
	g.DELETE("/tracks/:id/hotcues/:slot", h.DeleteHotCue)
	g.DELETE("/tracks/:id/hotcues", h.ClearHotCues)
}
