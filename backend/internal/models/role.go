package models

// UserRole represents a user's role in the system.
// Roles determine what actions a user can perform.
type UserRole string

const (
	// RoleGuest is for unauthenticated or limited-access users.
	RoleGuest UserRole = "guest"

	// RoleSubscriber is the default role for authenticated users who can listen and create playlists.
	RoleSubscriber UserRole = "subscriber"

	// RoleArtist is for users who can upload and manage their own music.
	RoleArtist UserRole = "artist"

	// RoleAdmin is for users who can moderate content and manage the platform.
	RoleAdmin UserRole = "admin"
)

// validRoles contains all valid role values for validation.
var validRoles = map[UserRole]bool{
	RoleGuest:      true,
	RoleSubscriber: true,
	RoleArtist:     true,
	RoleAdmin:      true,
}

// Permission represents an action that can be performed in the system.
type Permission string

const (
	// PermissionUploadTracks allows uploading audio tracks.
	PermissionUploadTracks Permission = "upload_tracks"

	// PermissionCreatePlaylists allows creating playlists.
	PermissionCreatePlaylists Permission = "create_playlists"

	// PermissionCreatePublicPlaylists allows creating public playlists.
	PermissionCreatePublicPlaylists Permission = "create_public_playlists"

	// PermissionFollowArtists allows following artists.
	PermissionFollowArtists Permission = "follow_artists"

	// PermissionHaveFollowers allows being followed by other users.
	PermissionHaveFollowers Permission = "have_followers"

	// PermissionManageOwnContent allows managing own content (playlists, etc).
	PermissionManageOwnContent Permission = "manage_own_content"

	// PermissionModerateContent allows moderating any content.
	PermissionModerateContent Permission = "moderate_content"

	// PermissionManageUsers allows managing user accounts.
	PermissionManageUsers Permission = "manage_users"
)

// RolePermissions maps roles to their allowed permissions.
var RolePermissions = map[UserRole]map[Permission]bool{
	RoleGuest: {},
	RoleSubscriber: {
		PermissionCreatePlaylists:       true,
		PermissionCreatePublicPlaylists: true,
		PermissionFollowArtists:         true,
		PermissionManageOwnContent:      true,
	},
	RoleArtist: {
		PermissionUploadTracks:          true,
		PermissionCreatePlaylists:       true,
		PermissionCreatePublicPlaylists: true,
		PermissionFollowArtists:         true,
		PermissionHaveFollowers:         true,
		PermissionManageOwnContent:      true,
	},
	RoleAdmin: {
		PermissionUploadTracks:          true,
		PermissionCreatePlaylists:       true,
		PermissionCreatePublicPlaylists: true,
		PermissionFollowArtists:         true,
		PermissionHaveFollowers:         true,
		PermissionManageOwnContent:      true,
		PermissionModerateContent:       true,
		PermissionManageUsers:           true,
	},
}

// IsValid returns true if the role is a valid system role.
func (r UserRole) IsValid() bool {
	return validRoles[r]
}

// CognitoGroupName returns the Cognito group name for this role.
func (r UserRole) CognitoGroupName() string {
	return string(r)
}

// HasPermission returns true if the role has the specified permission.
func (r UserRole) HasPermission(p Permission) bool {
	perms, ok := RolePermissions[r]
	if !ok {
		return false
	}
	return perms[p]
}

// CanUploadTracks returns true if the role allows uploading tracks.
func (r UserRole) CanUploadTracks() bool {
	return r.HasPermission(PermissionUploadTracks)
}

// CanHaveFollowers returns true if the role allows having followers.
func (r UserRole) CanHaveFollowers() bool {
	return r.HasPermission(PermissionHaveFollowers)
}

// CanModerateContent returns true if the role allows moderating content.
func (r UserRole) CanModerateContent() bool {
	return r.HasPermission(PermissionModerateContent)
}

// DefaultUserRole returns the default role for new authenticated users.
func DefaultUserRole() UserRole {
	return RoleSubscriber
}

// AllUserRoles returns a slice of all valid roles.
func AllUserRoles() []UserRole {
	return []UserRole{RoleGuest, RoleSubscriber, RoleArtist, RoleAdmin}
}
