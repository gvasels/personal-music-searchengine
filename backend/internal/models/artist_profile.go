package models

import (
	"fmt"
	"time"
)

// EntityArtistProfile represents the entity type for artist profiles
const EntityArtistProfile EntityType = "ARTIST_PROFILE"

// ArtistProfile represents extended profile information for users with the Artist role.
type ArtistProfile struct {
	UserID        string            `json:"userId" dynamodbav:"userId"`
	DisplayName   string            `json:"displayName" dynamodbav:"displayName"`
	Bio           string            `json:"bio,omitempty" dynamodbav:"bio,omitempty"`
	AvatarURL     string            `json:"avatarUrl,omitempty" dynamodbav:"avatarUrl,omitempty"`
	BannerURL     string            `json:"bannerUrl,omitempty" dynamodbav:"bannerUrl,omitempty"`
	SocialLinks   map[string]string `json:"socialLinks,omitempty" dynamodbav:"socialLinks,omitempty"`
	Genres        []string          `json:"genres,omitempty" dynamodbav:"genres,omitempty"`
	FollowerCount int               `json:"followerCount" dynamodbav:"followerCount"`
	TrackCount    int               `json:"trackCount" dynamodbav:"trackCount"`
	AlbumCount    int               `json:"albumCount" dynamodbav:"albumCount"`
	TotalPlays    int64             `json:"totalPlays" dynamodbav:"totalPlays"`
	IsVerified    bool              `json:"isVerified" dynamodbav:"isVerified"`
	LinkedArtist  string            `json:"linkedArtist,omitempty" dynamodbav:"linkedArtist,omitempty"` // Link to catalog artist
	Timestamps
}

// ArtistProfileItem represents an ArtistProfile in DynamoDB single-table design
type ArtistProfileItem struct {
	DynamoDBItem
	ArtistProfile
}

// NewArtistProfile creates a new ArtistProfile with default values.
func NewArtistProfile(userID string) *ArtistProfile {
	now := time.Now()
	return &ArtistProfile{
		UserID:        userID,
		Bio:           "",
		SocialLinks:   make(map[string]string),
		Genres:        make([]string, 0),
		FollowerCount: 0,
		TrackCount:    0,
		AlbumCount:    0,
		TotalPlays:    0,
		IsVerified:    false,
		Timestamps: Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// NewArtistProfileItem creates a DynamoDB item for an artist profile.
func NewArtistProfileItem(profile ArtistProfile) ArtistProfileItem {
	item := ArtistProfileItem{
		DynamoDBItem: DynamoDBItem{
			PK:     fmt.Sprintf("USER#%s", profile.UserID),
			SK:     "ARTIST_PROFILE",
			GSI1PK: "ARTIST_PROFILE", // For listing all artist profiles
			GSI1SK: fmt.Sprintf("USER#%s", profile.UserID),
			Type:   string(EntityArtistProfile),
		},
		ArtistProfile: profile,
	}

	// GSI2 for linked artist lookup (if linked)
	if profile.LinkedArtist != "" {
		item.GSI2PK = fmt.Sprintf("LINKED_ARTIST#%s", profile.LinkedArtist)
		item.GSI2SK = fmt.Sprintf("USER#%s", profile.UserID)
	}

	return item
}

// IncrementFollowerCount increments the follower count by 1.
func (ap *ArtistProfile) IncrementFollowerCount() {
	ap.FollowerCount++
	ap.UpdatedAt = time.Now()
}

// DecrementFollowerCount decrements the follower count by 1, but not below 0.
func (ap *ArtistProfile) DecrementFollowerCount() {
	if ap.FollowerCount > 0 {
		ap.FollowerCount--
		ap.UpdatedAt = time.Now()
	}
}

// IncrementTrackCount increments the track count by 1.
func (ap *ArtistProfile) IncrementTrackCount() {
	ap.TrackCount++
	ap.UpdatedAt = time.Now()
}

// DecrementTrackCount decrements the track count by 1, but not below 0.
func (ap *ArtistProfile) DecrementTrackCount() {
	if ap.TrackCount > 0 {
		ap.TrackCount--
		ap.UpdatedAt = time.Now()
	}
}

// AddPlays adds the specified number of plays to the total.
func (ap *ArtistProfile) AddPlays(count int64) {
	ap.TotalPlays += count
	ap.UpdatedAt = time.Now()
}

// ArtistProfileResponse represents an artist profile in API responses.
type ArtistProfileResponse struct {
	UserID        string            `json:"userId"`
	DisplayName   string            `json:"displayName"`
	Bio           string            `json:"bio,omitempty"`
	AvatarURL     string            `json:"avatarUrl,omitempty"`
	BannerURL     string            `json:"bannerUrl,omitempty"`
	SocialLinks   map[string]string `json:"socialLinks,omitempty"`
	Genres        []string          `json:"genres,omitempty"`
	FollowerCount int               `json:"followerCount"`
	TrackCount    int               `json:"trackCount"`
	AlbumCount    int               `json:"albumCount"`
	TotalPlays    int64             `json:"totalPlays"`
	IsVerified    bool              `json:"isVerified"`
	LinkedArtist  string            `json:"linkedArtist,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
}

// ToResponse converts an ArtistProfile to an ArtistProfileResponse.
func (ap *ArtistProfile) ToResponse() ArtistProfileResponse {
	return ArtistProfileResponse{
		UserID:        ap.UserID,
		DisplayName:   ap.DisplayName,
		Bio:           ap.Bio,
		AvatarURL:     ap.AvatarURL,
		BannerURL:     ap.BannerURL,
		SocialLinks:   ap.SocialLinks,
		Genres:        ap.Genres,
		FollowerCount: ap.FollowerCount,
		TrackCount:    ap.TrackCount,
		AlbumCount:    ap.AlbumCount,
		TotalPlays:    ap.TotalPlays,
		IsVerified:    ap.IsVerified,
		LinkedArtist:  ap.LinkedArtist,
		CreatedAt:     ap.CreatedAt,
	}
}

// CreateArtistProfileRequest represents a request to create an artist profile.
type CreateArtistProfileRequest struct {
	DisplayName string            `json:"displayName" validate:"required,min=1,max=100"`
	Bio         string            `json:"bio,omitempty" validate:"omitempty,max=2000"`
	SocialLinks map[string]string `json:"socialLinks,omitempty"`
	Genres      []string          `json:"genres,omitempty" validate:"omitempty,max=10,dive,max=50"`
}

// UpdateArtistProfileRequest represents a request to update an artist profile.
type UpdateArtistProfileRequest struct {
	DisplayName *string            `json:"displayName,omitempty" validate:"omitempty,min=1,max=100"`
	Bio         *string            `json:"bio,omitempty" validate:"omitempty,max=2000"`
	AvatarURL   *string            `json:"avatarUrl,omitempty" validate:"omitempty,url"`
	BannerURL   *string            `json:"bannerUrl,omitempty" validate:"omitempty,url"`
	SocialLinks *map[string]string `json:"socialLinks,omitempty"`
	Genres      *[]string          `json:"genres,omitempty" validate:"omitempty,max=10,dive,max=50"`
}

// LinkArtistRequest represents a request to link an artist profile to a catalog artist.
type LinkArtistRequest struct {
	ArtistID string `json:"artistId" validate:"required,uuid"`
}
