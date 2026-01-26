package models

import "time"

// UserSummary represents a user in admin search results.
type UserSummary struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Role        UserRole  `json:"role"`
	Disabled    bool      `json:"disabled"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserDetails represents full user details for admin management.
type UserDetails struct {
	UserSummary
	LastLoginAt    *time.Time `json:"lastLoginAt,omitempty"`
	TrackCount     int        `json:"trackCount"`
	PlaylistCount  int        `json:"playlistCount"`
	AlbumCount     int        `json:"albumCount"`
	StorageUsed    int64      `json:"storageUsed"`
	FollowerCount  int        `json:"followerCount"`
	FollowingCount int        `json:"followingCount"`
}

// UpdateRoleRequest represents a request to update a user's role.
type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=guest subscriber artist admin"`
}

// UpdateStatusRequest represents a request to update a user's status.
type UpdateStatusRequest struct {
	Disabled bool `json:"disabled"`
}

// AdminSearchUsersRequest represents the query parameters for searching users.
type AdminSearchUsersRequest struct {
	Search string `query:"search" validate:"required,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Cursor string `query:"cursor"`
}

// AdminSearchUsersResponse represents the response from searching users.
type AdminSearchUsersResponse struct {
	Items      []UserSummary `json:"items"`
	NextCursor string        `json:"nextCursor,omitempty"`
}

// ToUserSummary converts a User to a UserSummary.
func (u *User) ToUserSummary(disabled bool) UserSummary {
	role := u.Role
	if role == "" {
		role = DefaultUserRole()
	}
	return UserSummary{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Role:        role,
		Disabled:    disabled,
		CreatedAt:   u.CreatedAt,
	}
}

// ToUserDetails converts a User to UserDetails with additional information.
func (u *User) ToUserDetails(disabled bool, lastLoginAt *time.Time, followerCount int) UserDetails {
	return UserDetails{
		UserSummary:    u.ToUserSummary(disabled),
		LastLoginAt:    lastLoginAt,
		TrackCount:     u.TrackCount,
		PlaylistCount:  u.PlaylistCount,
		AlbumCount:     u.AlbumCount,
		StorageUsed:    u.StorageUsed,
		FollowerCount:  followerCount,
		FollowingCount: u.FollowingCount,
	}
}

// ValidateRole checks if the role string is a valid UserRole.
func ValidateRole(role string) (UserRole, bool) {
	r := UserRole(role)
	return r, r.IsValid()
}
