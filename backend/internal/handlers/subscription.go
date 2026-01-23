package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
)

// SubscriptionHandler handles subscription-related endpoints
type SubscriptionHandler struct {
	subscriptionSvc *service.SubscriptionService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionSvc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionSvc: subscriptionSvc,
	}
}

// GetSubscription retrieves the current user's subscription
// GET /subscription
func (h *SubscriptionHandler) GetSubscription(c echo.Context) error {
	userID := c.Get("userID").(string)

	sub, err := h.subscriptionSvc.GetSubscription(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get subscription")
	}

	return c.JSON(http.StatusOK, sub)
}

// GetTiers returns available subscription tiers
// GET /subscription/tiers
func (h *SubscriptionHandler) GetTiers(c echo.Context) error {
	tiers := h.subscriptionSvc.GetTierConfigs()
	return c.JSON(http.StatusOK, tiers)
}

// GetStorageUsage returns storage usage information
// GET /subscription/storage
func (h *SubscriptionHandler) GetStorageUsage(c echo.Context) error {
	userID := c.Get("userID").(string)

	usage, err := h.subscriptionSvc.GetStorageUsage(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get storage usage")
	}

	return c.JSON(http.StatusOK, usage)
}

// CreateCheckoutSession creates a Stripe checkout session
// POST /subscription/checkout
func (h *SubscriptionHandler) CreateCheckoutSession(c echo.Context) error {
	userID := c.Get("userID").(string)

	var req models.CreateCheckoutSessionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Get success and cancel URLs from query params or use defaults
	successURL := c.QueryParam("successUrl")
	if successURL == "" {
		successURL = c.Request().Header.Get("Origin") + "/subscription/success"
	}
	cancelURL := c.QueryParam("cancelUrl")
	if cancelURL == "" {
		cancelURL = c.Request().Header.Get("Origin") + "/subscription/cancel"
	}

	session, err := h.subscriptionSvc.CreateCheckoutSession(c.Request().Context(), userID, req, successURL, cancelURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, session)
}

// CreatePortalSession creates a Stripe customer portal session
// POST /subscription/portal
func (h *SubscriptionHandler) CreatePortalSession(c echo.Context) error {
	userID := c.Get("userID").(string)

	returnURL := c.QueryParam("returnUrl")
	if returnURL == "" {
		returnURL = c.Request().Header.Get("Origin") + "/subscription"
	}

	session, err := h.subscriptionSvc.CreatePortalSession(c.Request().Context(), userID, returnURL)
	if err != nil {
		if err.Error() == "no subscription found" {
			return echo.NewHTTPError(http.StatusNotFound, "No active subscription found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, session)
}

// HandleStripeWebhook processes Stripe webhook events
// POST /webhooks/stripe
func (h *SubscriptionHandler) HandleStripeWebhook(c echo.Context) error {
	// Parse the webhook payload
	var event struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&event); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook payload")
	}

	// Extract the object from data
	object, ok := event.Data["object"].(map[string]interface{})
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook data")
	}

	if err := h.subscriptionSvc.HandleStripeWebhook(c.Request().Context(), event.Type, object); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// RegisterSubscriptionRoutes registers subscription routes
func RegisterSubscriptionRoutes(g *echo.Group, h *SubscriptionHandler) {
	g.GET("/subscription", h.GetSubscription)
	g.GET("/subscription/tiers", h.GetTiers)
	g.GET("/subscription/storage", h.GetStorageUsage)
	g.POST("/subscription/checkout", h.CreateCheckoutSession)
	g.POST("/subscription/portal", h.CreatePortalSession)
}

// RegisterWebhookRoutes registers webhook routes (usually unprotected)
func RegisterWebhookRoutes(g *echo.Group, h *SubscriptionHandler) {
	g.POST("/webhooks/stripe", h.HandleStripeWebhook)
}
