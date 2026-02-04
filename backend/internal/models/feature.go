package models

import "time"

// FeatureKey represents available feature flags
type FeatureKey string

const (
	// DJ Module features
	FeatureDJModule     FeatureKey = "DJ_MODULE"
	FeatureCrates       FeatureKey = "CRATES"
	FeatureHotCues      FeatureKey = "HOT_CUES"
	FeatureBPMMatching  FeatureKey = "BPM_MATCHING"
	FeatureKeyMatching  FeatureKey = "KEY_MATCHING"
	FeatureMixRecording FeatureKey = "MIX_RECORDING"

	// Creator features
	FeatureBulkEdit      FeatureKey = "BULK_EDIT"
	FeatureAdvancedStats FeatureKey = "ADVANCED_STATS"
	FeatureAPIAccess     FeatureKey = "API_ACCESS"

	// Storage features
	FeatureUnlimitedStorage FeatureKey = "UNLIMITED_STORAGE"
	FeatureHQStreaming      FeatureKey = "HQ_STREAMING"
)

// SubscriptionTier represents user subscription levels
type SubscriptionTier string

const (
	TierFree    SubscriptionTier = "free"
	TierCreator SubscriptionTier = "creator"
	TierPro     SubscriptionTier = "pro"
)

// FeatureFlag represents a feature flag configuration
type FeatureFlag struct {
	Key           FeatureKey         `json:"key" dynamodbav:"key"`
	Name          string             `json:"name" dynamodbav:"name"`
	Description   string             `json:"description" dynamodbav:"description"`
	EnabledTiers  []SubscriptionTier `json:"enabledTiers" dynamodbav:"enabledTiers"`   // Tiers that have this feature
	DefaultValue  bool               `json:"defaultValue" dynamodbav:"defaultValue"`   // Default if no tier match
	GlobalEnabled bool               `json:"globalEnabled" dynamodbav:"globalEnabled"` // Master switch
	Timestamps
}

// FeatureFlagItem represents a FeatureFlag in DynamoDB single-table design
type FeatureFlagItem struct {
	DynamoDBItem
	FeatureFlag
}

// NewFeatureFlagItem creates a DynamoDB item for a feature flag
// PK: FEATURE, SK: FLAG#{key}
func NewFeatureFlagItem(flag FeatureFlag) FeatureFlagItem {
	return FeatureFlagItem{
		DynamoDBItem: DynamoDBItem{
			PK:   "FEATURE",
			SK:   "FLAG#" + string(flag.Key),
			Type: "FEATURE_FLAG",
		},
		FeatureFlag: flag,
	}
}

// UserFeatureOverride allows per-user feature overrides
type UserFeatureOverride struct {
	UserID     string     `json:"userId" dynamodbav:"userId"`
	FeatureKey FeatureKey `json:"featureKey" dynamodbav:"featureKey"`
	Enabled    bool       `json:"enabled" dynamodbav:"enabled"`
	Reason     string     `json:"reason,omitempty" dynamodbav:"reason,omitempty"` // Why override exists
	ExpiresAt  *time.Time `json:"expiresAt,omitempty" dynamodbav:"expiresAt,omitempty"`
	Timestamps
}

// UserFeatureOverrideItem represents a user override in DynamoDB
// PK: USER#{userId}, SK: FEATURE_OVERRIDE#{key}
type UserFeatureOverrideItem struct {
	DynamoDBItem
	UserFeatureOverride
}

// NewUserFeatureOverrideItem creates a DynamoDB item for a user override
func NewUserFeatureOverrideItem(override UserFeatureOverride) UserFeatureOverrideItem {
	return UserFeatureOverrideItem{
		DynamoDBItem: DynamoDBItem{
			PK:   "USER#" + override.UserID,
			SK:   "FEATURE_OVERRIDE#" + string(override.FeatureKey),
			Type: "FEATURE_OVERRIDE",
		},
		UserFeatureOverride: override,
	}
}

// UserFeaturesResponse represents the API response for user's enabled features
type UserFeaturesResponse struct {
	Tier     SubscriptionTier `json:"tier"`
	Features map[string]bool  `json:"features"`
}

// DefaultFeatureFlags returns the default feature flag configurations
func DefaultFeatureFlags() []FeatureFlag {
	now := time.Now()
	timestamps := Timestamps{CreatedAt: now, UpdatedAt: now}

	return []FeatureFlag{
		{Key: FeatureDJModule, Name: "DJ Module", Description: "Access to DJ features", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureCrates, Name: "Crates", Description: "DJ crate organization", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureHotCues, Name: "Hot Cues", Description: "Hot cue points on tracks", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureBPMMatching, Name: "BPM Matching", Description: "Find tracks with compatible BPM", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureKeyMatching, Name: "Key Matching", Description: "Find tracks with compatible keys", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureMixRecording, Name: "Mix Recording", Description: "Record DJ mixes", EnabledTiers: []SubscriptionTier{TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureBulkEdit, Name: "Bulk Edit", Description: "Edit multiple tracks at once", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureAdvancedStats, Name: "Advanced Statistics", Description: "Detailed listening analytics", EnabledTiers: []SubscriptionTier{TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureAPIAccess, Name: "API Access", Description: "External API access", EnabledTiers: []SubscriptionTier{TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureUnlimitedStorage, Name: "Unlimited Storage", Description: "No storage limits", EnabledTiers: []SubscriptionTier{TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
		{Key: FeatureHQStreaming, Name: "HQ Streaming", Description: "High quality audio streaming", EnabledTiers: []SubscriptionTier{TierCreator, TierPro}, DefaultValue: false, GlobalEnabled: true, Timestamps: timestamps},
	}
}

// IsEnabled checks if a feature is enabled for a given tier
func (f *FeatureFlag) IsEnabled(tier SubscriptionTier) bool {
	if !f.GlobalEnabled {
		return false
	}
	for _, t := range f.EnabledTiers {
		if t == tier {
			return true
		}
	}
	return f.DefaultValue
}
