package models

import "fmt"

// ProfileVisibility represents visibility options for user profile
type ProfileVisibility string

const (
	ProfileVisibilityPrivate ProfileVisibility = "private"
	ProfileVisibilityPublic  ProfileVisibility = "public"
)

// AudioQuality represents audio streaming quality preferences
type AudioQuality string

const (
	QualityLow      AudioQuality = "low"      // 128kbps
	QualityMedium   AudioQuality = "medium"   // 256kbps
	QualityHigh     AudioQuality = "high"     // 320kbps
	QualityLossless AudioQuality = "lossless" // FLAC
)

// DuplicateHandling represents how to handle duplicate uploads
type DuplicateHandling string

const (
	DuplicateSkip    DuplicateHandling = "skip"
	DuplicateReplace DuplicateHandling = "replace"
	DuplicateKeep    DuplicateHandling = "keep" // Keep both
)

// UserSettings represents all user preferences and settings
type UserSettings struct {
	Notifications NotificationSettings `json:"notifications" dynamodbav:"notifications"`
	Privacy       PrivacySettings      `json:"privacy" dynamodbav:"privacy"`
	Player        PlayerSettings       `json:"player" dynamodbav:"player"`
	Library       LibrarySettings      `json:"library" dynamodbav:"library"`
}

// NotificationSettings represents notification preferences
type NotificationSettings struct {
	EmailDigest     bool `json:"emailDigest" dynamodbav:"emailDigest"`
	PlaylistUpdates bool `json:"playlistUpdates" dynamodbav:"playlistUpdates"`
	NewFeatures     bool `json:"newFeatures" dynamodbav:"newFeatures"`
	MarketingEmails bool `json:"marketingEmails" dynamodbav:"marketingEmails"`
}

// PrivacySettings represents privacy preferences
type PrivacySettings struct {
	ProfileVisibility     ProfileVisibility `json:"profileVisibility" dynamodbav:"profileVisibility"`
	ShowListeningActivity bool              `json:"showListeningActivity" dynamodbav:"showListeningActivity"`
	AllowPlaylistSharing  bool              `json:"allowPlaylistSharing" dynamodbav:"allowPlaylistSharing"`
}

// PlayerSettings represents audio player preferences
type PlayerSettings struct {
	DefaultVolume     float64      `json:"defaultVolume" dynamodbav:"defaultVolume"`
	CrossfadeEnabled  bool         `json:"crossfadeEnabled" dynamodbav:"crossfadeEnabled"`
	CrossfadeDuration int          `json:"crossfadeDuration" dynamodbav:"crossfadeDuration"` // seconds
	AudioQuality      AudioQuality `json:"audioQuality" dynamodbav:"audioQuality"`
	NormalizeVolume   bool         `json:"normalizeVolume" dynamodbav:"normalizeVolume"`
}

// LibrarySettings represents library organization preferences
type LibrarySettings struct {
	AutoOrganize      bool              `json:"autoOrganize" dynamodbav:"autoOrganize"`
	DuplicateHandling DuplicateHandling `json:"duplicateHandling" dynamodbav:"duplicateHandling"`
	ExtractMetadata   bool              `json:"extractMetadata" dynamodbav:"extractMetadata"`
}

// DefaultUserSettings returns the default settings for a new user
func DefaultUserSettings() UserSettings {
	return UserSettings{
		Notifications: NotificationSettings{
			EmailDigest:     true,
			PlaylistUpdates: true,
			NewFeatures:     true,
			MarketingEmails: false,
		},
		Privacy: PrivacySettings{
			ProfileVisibility:     ProfileVisibilityPrivate,
			ShowListeningActivity: false,
			AllowPlaylistSharing:  false,
		},
		Player: PlayerSettings{
			DefaultVolume:     0.8,
			CrossfadeEnabled:  false,
			CrossfadeDuration: 0,
			AudioQuality:      QualityHigh,
			NormalizeVolume:   false,
		},
		Library: LibrarySettings{
			AutoOrganize:      true,
			DuplicateHandling: DuplicateSkip,
			ExtractMetadata:   true,
		},
	}
}

// Validate validates the user settings
func (s *UserSettings) Validate() error {
	// Validate volume range
	if s.Player.DefaultVolume < 0 || s.Player.DefaultVolume > 1 {
		return fmt.Errorf("defaultVolume must be between 0 and 1")
	}

	// Validate crossfade duration
	if s.Player.CrossfadeDuration < 0 || s.Player.CrossfadeDuration > 12 {
		return fmt.Errorf("crossfadeDuration must be between 0 and 12 seconds")
	}

	// Validate audio quality
	validQualities := map[AudioQuality]bool{
		QualityLow:      true,
		QualityMedium:   true,
		QualityHigh:     true,
		QualityLossless: true,
	}
	if !validQualities[s.Player.AudioQuality] {
		return fmt.Errorf("invalid audioQuality: %s", s.Player.AudioQuality)
	}

	// Validate visibility
	validProfileVisibilities := map[ProfileVisibility]bool{
		ProfileVisibilityPrivate: true,
		ProfileVisibilityPublic:  true,
	}
	if !validProfileVisibilities[s.Privacy.ProfileVisibility] {
		return fmt.Errorf("invalid profileVisibility: %s", s.Privacy.ProfileVisibility)
	}

	// Validate duplicate handling
	validDuplicateHandling := map[DuplicateHandling]bool{
		DuplicateSkip:    true,
		DuplicateReplace: true,
		DuplicateKeep:    true,
	}
	if !validDuplicateHandling[s.Library.DuplicateHandling] {
		return fmt.Errorf("invalid duplicateHandling: %s", s.Library.DuplicateHandling)
	}

	return nil
}
