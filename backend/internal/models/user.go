package models

import "time"

// User represents a user profile in the system
type User struct {
	ID             string           `json:"id" dynamodbav:"id"`
	Email          string           `json:"email" dynamodbav:"email"`
	DisplayName    string           `json:"displayName" dynamodbav:"displayName"`
	AvatarURL      string           `json:"avatarUrl,omitempty" dynamodbav:"avatarUrl,omitempty"`
	Role           UserRole         `json:"role" dynamodbav:"role"`
	Tier           SubscriptionTier `json:"tier" dynamodbav:"tier"` // Deprecated: Use Role instead
	FollowingCount int              `json:"followingCount" dynamodbav:"followingCount"`
	Timestamps
	StorageUsed   int64 `json:"storageUsed" dynamodbav:"storageUsed"`
	StorageLimit  int64 `json:"storageLimit" dynamodbav:"storageLimit"`
	TrackCount    int   `json:"trackCount" dynamodbav:"trackCount"`
	AlbumCount    int   `json:"albumCount" dynamodbav:"albumCount"`
	PlaylistCount int   `json:"playlistCount" dynamodbav:"playlistCount"`
}

// UserItem represents a User in DynamoDB single-table design
type UserItem struct {
	DynamoDBItem
	User
}

// NewUserItem creates a DynamoDB item for a user
func NewUserItem(user User) UserItem {
	return UserItem{
		DynamoDBItem: DynamoDBItem{
			PK:   "USER#" + user.ID,
			SK:   "PROFILE",
			Type: string(EntityUser),
		},
		User: user,
	}
}

// UserPreferences represents user preferences
type UserPreferences struct {
	UserID            string `json:"userId" dynamodbav:"userId"`
	Theme             string `json:"theme" dynamodbav:"theme"`
	DefaultSortField  string `json:"defaultSortField" dynamodbav:"defaultSortField"`
	DefaultSortOrder  string `json:"defaultSortOrder" dynamodbav:"defaultSortOrder"`
	DefaultViewType   string `json:"defaultViewType" dynamodbav:"defaultViewType"`
	AudioQuality      string `json:"audioQuality" dynamodbav:"audioQuality"`
	CrossfadeEnabled  bool   `json:"crossfadeEnabled" dynamodbav:"crossfadeEnabled"`
	CrossfadeDuration int    `json:"crossfadeDuration" dynamodbav:"crossfadeDuration"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"displayName" validate:"required,min=1,max=100"`
}

// UpdateUserRequest represents a request to update a user profile
type UpdateUserRequest struct {
	DisplayName *string `json:"displayName,omitempty" validate:"omitempty,min=1,max=100"`
	AvatarURL   *string `json:"avatarUrl,omitempty" validate:"omitempty,url"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID             string           `json:"id"`
	Email          string           `json:"email"`
	DisplayName    string           `json:"displayName"`
	AvatarURL      string           `json:"avatarUrl,omitempty"`
	Role           UserRole         `json:"role"`
	Tier           SubscriptionTier `json:"tier"` // Deprecated: Use Role instead
	FollowingCount int              `json:"followingCount"`
	CreatedAt      time.Time        `json:"createdAt"`
	StorageUsed    int64            `json:"storageUsed"`
	StorageLimit   int64            `json:"storageLimit"`
	TrackCount     int              `json:"trackCount"`
	AlbumCount     int              `json:"albumCount"`
	PlaylistCount  int              `json:"playlistCount"`
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() UserResponse {
	tier := u.Tier
	if tier == "" {
		tier = TierFree
	}
	role := u.Role
	if role == "" {
		role = DefaultUserRole()
	}
	return UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		DisplayName:    u.DisplayName,
		AvatarURL:      u.AvatarURL,
		Role:           role,
		Tier:           tier,
		FollowingCount: u.FollowingCount,
		CreatedAt:      u.CreatedAt,
		StorageUsed:    u.StorageUsed,
		StorageLimit:   u.StorageLimit,
		TrackCount:     u.TrackCount,
		AlbumCount:     u.AlbumCount,
		PlaylistCount:  u.PlaylistCount,
	}
}
