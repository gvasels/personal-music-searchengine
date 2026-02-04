package models

import "time"

// TierConfig defines the configuration for each subscription tier
type TierConfig struct {
	Tier              SubscriptionTier `json:"tier"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	MonthlyPriceCents int              `json:"monthlyPriceCents"`
	YearlyPriceCents  int              `json:"yearlyPriceCents"`
	StorageLimitBytes int64            `json:"storageLimitBytes"` // -1 for unlimited
	Features          []FeatureKey     `json:"features"`
}

// GetTierConfigs returns all tier configurations
func GetTierConfigs() []TierConfig {
	return []TierConfig{
		{
			Tier:              TierFree,
			Name:              "Free",
			Description:       "Basic music library features",
			MonthlyPriceCents: 0,
			YearlyPriceCents:  0,
			StorageLimitBytes: 5 * 1024 * 1024 * 1024, // 5 GB
			Features:          []FeatureKey{},
		},
		{
			Tier:              TierCreator,
			Name:              "Creator",
			Description:       "For DJs and music creators",
			MonthlyPriceCents: 999,                       // $9.99
			YearlyPriceCents:  9900,                      // $99/year (save ~17%)
			StorageLimitBytes: 50 * 1024 * 1024 * 1024,   // 50 GB
			Features: []FeatureKey{
				FeatureDJModule,
				FeatureCrates,
				FeatureHotCues,
				FeatureBPMMatching,
				FeatureKeyMatching,
				FeatureBulkEdit,
				FeatureHQStreaming,
			},
		},
		{
			Tier:              TierPro,
			Name:              "Pro",
			Description:       "Professional features for power users",
			MonthlyPriceCents: 1999, // $19.99
			YearlyPriceCents:  19900, // $199/year (save ~17%)
			StorageLimitBytes: -1,   // Unlimited
			Features: []FeatureKey{
				FeatureDJModule,
				FeatureCrates,
				FeatureHotCues,
				FeatureBPMMatching,
				FeatureKeyMatching,
				FeatureMixRecording,
				FeatureBulkEdit,
				FeatureAdvancedStats,
				FeatureAPIAccess,
				FeatureUnlimitedStorage,
				FeatureHQStreaming,
			},
		},
	}
}

// GetTierConfig returns the configuration for a specific tier
func GetTierConfig(tier SubscriptionTier) *TierConfig {
	for _, config := range GetTierConfigs() {
		if config.Tier == tier {
			return &config
		}
	}
	return nil
}

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
)

// SubscriptionInterval represents billing interval
type SubscriptionInterval string

const (
	SubscriptionIntervalMonthly SubscriptionInterval = "monthly"
	SubscriptionIntervalYearly  SubscriptionInterval = "yearly"
)

// Subscription represents a user's subscription
type Subscription struct {
	UserID              string               `json:"userId" dynamodbav:"userId"`
	Tier                SubscriptionTier     `json:"tier" dynamodbav:"tier"`
	Status              SubscriptionStatus   `json:"status" dynamodbav:"status"`
	Interval            SubscriptionInterval `json:"interval" dynamodbav:"interval"`
	StripeCustomerID    string               `json:"stripeCustomerId,omitempty" dynamodbav:"stripeCustomerId,omitempty"`
	StripeSubscriptionID string              `json:"stripeSubscriptionId,omitempty" dynamodbav:"stripeSubscriptionId,omitempty"`
	StripePriceID       string               `json:"stripePriceId,omitempty" dynamodbav:"stripePriceId,omitempty"`
	CurrentPeriodStart  time.Time            `json:"currentPeriodStart" dynamodbav:"currentPeriodStart"`
	CurrentPeriodEnd    time.Time            `json:"currentPeriodEnd" dynamodbav:"currentPeriodEnd"`
	CancelAtPeriodEnd   bool                 `json:"cancelAtPeriodEnd" dynamodbav:"cancelAtPeriodEnd"`
	TrialEnd            *time.Time           `json:"trialEnd,omitempty" dynamodbav:"trialEnd,omitempty"`
	Timestamps
}

// SubscriptionItem represents a Subscription in DynamoDB single-table design
// PK: USER#{userId}, SK: SUBSCRIPTION
type SubscriptionItem struct {
	DynamoDBItem
	Subscription
}

// NewSubscriptionItem creates a DynamoDB item for a subscription
func NewSubscriptionItem(sub Subscription) SubscriptionItem {
	return SubscriptionItem{
		DynamoDBItem: DynamoDBItem{
			PK:   "USER#" + sub.UserID,
			SK:   "SUBSCRIPTION",
			Type: "SUBSCRIPTION",
		},
		Subscription: sub,
	}
}

// SubscriptionResponse represents subscription in API responses
type SubscriptionResponse struct {
	UserID             string               `json:"userId"`
	Tier               SubscriptionTier     `json:"tier"`
	TierName           string               `json:"tierName"`
	Status             SubscriptionStatus   `json:"status"`
	Interval           SubscriptionInterval `json:"interval"`
	CurrentPeriodStart time.Time            `json:"currentPeriodStart"`
	CurrentPeriodEnd   time.Time            `json:"currentPeriodEnd"`
	CancelAtPeriodEnd  bool                 `json:"cancelAtPeriodEnd"`
	TrialEnd           *time.Time           `json:"trialEnd,omitempty"`
	StorageLimit       int64                `json:"storageLimit"`
	StorageUsed        int64                `json:"storageUsed"`
	Features           []FeatureKey         `json:"features"`
}

// CreateCheckoutSessionRequest represents a request to create a Stripe checkout session
type CreateCheckoutSessionRequest struct {
	Tier     SubscriptionTier     `json:"tier" validate:"required,oneof=creator pro"`
	Interval SubscriptionInterval `json:"interval" validate:"required,oneof=monthly yearly"`
}

// CheckoutSessionResponse represents the response with a checkout URL
type CheckoutSessionResponse struct {
	CheckoutURL string `json:"checkoutUrl"`
	SessionID   string `json:"sessionId"`
}

// PortalSessionResponse represents a Stripe customer portal session
type PortalSessionResponse struct {
	PortalURL string `json:"portalUrl"`
}

// StorageUsageResponse represents storage usage information
type StorageUsageResponse struct {
	StorageUsedBytes  int64   `json:"storageUsedBytes"`
	StorageLimitBytes int64   `json:"storageLimitBytes"` // -1 for unlimited
	UsagePercent      float64 `json:"usagePercent"`      // 0-100, -1 for unlimited
}
