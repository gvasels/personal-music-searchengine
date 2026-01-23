package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// SubscriptionRepository defines the repository interface for subscriptions
type SubscriptionRepository interface {
	GetSubscription(ctx context.Context, userID string) (*models.Subscription, error)
	PutSubscription(ctx context.Context, sub models.Subscription) error
	DeleteSubscription(ctx context.Context, userID string) error
	UpdateUserTier(ctx context.Context, userID string, tier models.SubscriptionTier) error
	GetUserStorageUsage(ctx context.Context, userID string) (int64, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

// StripeClient defines the interface for Stripe operations (mock-able for testing)
type StripeClient interface {
	CreateCustomer(ctx context.Context, userID, email string) (string, error)
	CreateCheckoutSession(ctx context.Context, customerID string, priceID string, successURL, cancelURL string) (string, string, error)
	CreatePortalSession(ctx context.Context, customerID string, returnURL string) (string, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
}

// SubscriptionService handles subscription operations
type SubscriptionService struct {
	repo   SubscriptionRepository
	stripe StripeClient
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(repo SubscriptionRepository, stripe StripeClient) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		stripe: stripe,
	}
}

// GetSubscription retrieves a user's subscription
func (s *SubscriptionService) GetSubscription(ctx context.Context, userID string) (*models.SubscriptionResponse, error) {
	sub, err := s.repo.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If no subscription exists, return free tier defaults
	if sub == nil {
		user, err := s.repo.GetUser(ctx, userID)
		if err != nil {
			return nil, err
		}

		storageUsed := int64(0)
		if user != nil {
			storageUsed = user.StorageUsed
		}

		tierConfig := models.GetTierConfig(models.TierFree)
		return &models.SubscriptionResponse{
			UserID:             userID,
			Tier:               models.TierFree,
			TierName:           tierConfig.Name,
			Status:             models.SubscriptionStatusActive,
			Interval:           models.SubscriptionIntervalMonthly,
			CurrentPeriodStart: time.Now(),
			CurrentPeriodEnd:   time.Now().AddDate(100, 0, 0), // Far future
			CancelAtPeriodEnd:  false,
			StorageLimit:       tierConfig.StorageLimitBytes,
			StorageUsed:        storageUsed,
			Features:           tierConfig.Features,
		}, nil
	}

	tierConfig := models.GetTierConfig(sub.Tier)
	storageUsed, _ := s.repo.GetUserStorageUsage(ctx, userID)

	return &models.SubscriptionResponse{
		UserID:             sub.UserID,
		Tier:               sub.Tier,
		TierName:           tierConfig.Name,
		Status:             sub.Status,
		Interval:           sub.Interval,
		CurrentPeriodStart: sub.CurrentPeriodStart,
		CurrentPeriodEnd:   sub.CurrentPeriodEnd,
		CancelAtPeriodEnd:  sub.CancelAtPeriodEnd,
		TrialEnd:           sub.TrialEnd,
		StorageLimit:       tierConfig.StorageLimitBytes,
		StorageUsed:        storageUsed,
		Features:           tierConfig.Features,
	}, nil
}

// GetStorageUsage retrieves storage usage information
func (s *SubscriptionService) GetStorageUsage(ctx context.Context, userID string) (*models.StorageUsageResponse, error) {
	sub, err := s.repo.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	tier := models.TierFree
	if sub != nil {
		tier = sub.Tier
	}

	tierConfig := models.GetTierConfig(tier)
	storageUsed, err := s.repo.GetUserStorageUsage(ctx, userID)
	if err != nil {
		return nil, err
	}

	usagePercent := float64(-1)
	if tierConfig.StorageLimitBytes > 0 {
		usagePercent = float64(storageUsed) / float64(tierConfig.StorageLimitBytes) * 100
	}

	return &models.StorageUsageResponse{
		StorageUsedBytes:  storageUsed,
		StorageLimitBytes: tierConfig.StorageLimitBytes,
		UsagePercent:      usagePercent,
	}, nil
}

// CreateCheckoutSession creates a Stripe checkout session for upgrading
func (s *SubscriptionService) CreateCheckoutSession(ctx context.Context, userID string, req models.CreateCheckoutSessionRequest, successURL, cancelURL string) (*models.CheckoutSessionResponse, error) {
	// Get existing subscription to get customer ID
	sub, err := s.repo.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	var customerID string
	if sub != nil && sub.StripeCustomerID != "" {
		customerID = sub.StripeCustomerID
	} else {
		// Create a new Stripe customer
		user, err := s.repo.GetUser(ctx, userID)
		if err != nil {
			return nil, err
		}

		email := ""
		if user != nil {
			email = user.Email
		}

		customerID, err = s.stripe.CreateCustomer(ctx, userID, email)
		if err != nil {
			return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
		}
	}

	// Get price ID for the tier and interval
	priceID := s.getPriceID(req.Tier, req.Interval)

	checkoutURL, sessionID, err := s.stripe.CreateCheckoutSession(ctx, customerID, priceID, successURL, cancelURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return &models.CheckoutSessionResponse{
		CheckoutURL: checkoutURL,
		SessionID:   sessionID,
	}, nil
}

// CreatePortalSession creates a Stripe customer portal session
func (s *SubscriptionService) CreatePortalSession(ctx context.Context, userID, returnURL string) (*models.PortalSessionResponse, error) {
	sub, err := s.repo.GetSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	if sub == nil || sub.StripeCustomerID == "" {
		return nil, fmt.Errorf("no subscription found")
	}

	portalURL, err := s.stripe.CreatePortalSession(ctx, sub.StripeCustomerID, returnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create portal session: %w", err)
	}

	return &models.PortalSessionResponse{
		PortalURL: portalURL,
	}, nil
}

// HandleStripeWebhook processes Stripe webhook events
func (s *SubscriptionService) HandleStripeWebhook(ctx context.Context, eventType string, data map[string]interface{}) error {
	switch eventType {
	case "customer.subscription.created", "customer.subscription.updated":
		return s.handleSubscriptionChange(ctx, data)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, data)
	default:
		// Ignore other events
		return nil
	}
}

func (s *SubscriptionService) handleSubscriptionChange(ctx context.Context, data map[string]interface{}) error {
	// Extract metadata to get userID
	metadata, ok := data["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing metadata in subscription event")
	}

	userID, ok := metadata["user_id"].(string)
	if !ok || userID == "" {
		return fmt.Errorf("missing user_id in subscription metadata")
	}

	// Extract subscription details
	stripeSubID, _ := data["id"].(string)
	stripeCustomerID, _ := data["customer"].(string)
	status, _ := data["status"].(string)
	cancelAtPeriodEnd, _ := data["cancel_at_period_end"].(bool)
	currentPeriodStart, _ := data["current_period_start"].(float64)
	currentPeriodEnd, _ := data["current_period_end"].(float64)

	// Determine tier from price
	tier := s.getTierFromPrice(data)

	now := time.Now()
	sub := models.Subscription{
		UserID:               userID,
		Tier:                 tier,
		Status:               models.SubscriptionStatus(status),
		Interval:             s.getIntervalFromPrice(data),
		StripeCustomerID:     stripeCustomerID,
		StripeSubscriptionID: stripeSubID,
		CurrentPeriodStart:   time.Unix(int64(currentPeriodStart), 0),
		CurrentPeriodEnd:     time.Unix(int64(currentPeriodEnd), 0),
		CancelAtPeriodEnd:    cancelAtPeriodEnd,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := s.repo.PutSubscription(ctx, sub); err != nil {
		return err
	}

	// Update user's tier
	return s.repo.UpdateUserTier(ctx, userID, tier)
}

func (s *SubscriptionService) handleSubscriptionDeleted(ctx context.Context, data map[string]interface{}) error {
	metadata, ok := data["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing metadata")
	}

	userID, ok := metadata["user_id"].(string)
	if !ok || userID == "" {
		return fmt.Errorf("missing user_id")
	}

	// Delete subscription and reset to free tier
	if err := s.repo.DeleteSubscription(ctx, userID); err != nil {
		return err
	}

	return s.repo.UpdateUserTier(ctx, userID, models.TierFree)
}

// GetTierConfigs returns all tier configurations
func (s *SubscriptionService) GetTierConfigs() []models.TierConfig {
	return models.GetTierConfigs()
}

// Helper functions

func (s *SubscriptionService) getPriceID(tier models.SubscriptionTier, interval models.SubscriptionInterval) string {
	// In production, these would be actual Stripe price IDs
	// For now, return mock IDs
	return fmt.Sprintf("price_%s_%s", tier, interval)
}

func (s *SubscriptionService) getTierFromPrice(data map[string]interface{}) models.SubscriptionTier {
	// In production, map Stripe price IDs to tiers
	// For now, check metadata
	metadata, _ := data["metadata"].(map[string]interface{})
	tier, _ := metadata["tier"].(string)

	switch tier {
	case "creator":
		return models.TierCreator
	case "pro":
		return models.TierPro
	default:
		return models.TierFree
	}
}

func (s *SubscriptionService) getIntervalFromPrice(data map[string]interface{}) models.SubscriptionInterval {
	// In production, inspect the price object
	// For now, default to monthly
	return models.SubscriptionIntervalMonthly
}

// MockStripeClient is a mock implementation for development/testing
type MockStripeClient struct{}

func NewMockStripeClient() *MockStripeClient {
	return &MockStripeClient{}
}

func (m *MockStripeClient) CreateCustomer(ctx context.Context, userID, email string) (string, error) {
	return "cus_mock_" + userID, nil
}

func (m *MockStripeClient) CreateCheckoutSession(ctx context.Context, customerID, priceID, successURL, cancelURL string) (string, string, error) {
	sessionID := "cs_mock_session"
	return successURL + "?session_id=" + sessionID, sessionID, nil
}

func (m *MockStripeClient) CreatePortalSession(ctx context.Context, customerID, returnURL string) (string, error) {
	return returnURL + "?portal=mock", nil
}

func (m *MockStripeClient) CancelSubscription(ctx context.Context, subscriptionID string) error {
	return nil
}
