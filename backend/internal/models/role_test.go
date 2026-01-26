package models

import (
	"testing"
)

func TestUserRoleConstants(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected string
	}{
		{"Guest role", RoleGuest, "guest"},
		{"Subscriber role", RoleSubscriber, "subscriber"},
		{"Artist role", RoleArtist, "artist"},
		{"Admin role", RoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.role)
			}
		})
	}
}

func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"Valid guest", RoleGuest, true},
		{"Valid subscriber", RoleSubscriber, true},
		{"Valid artist", RoleArtist, true},
		{"Valid admin", RoleAdmin, true},
		{"Invalid empty", UserRole(""), false},
		{"Invalid unknown", UserRole("unknown"), false},
		{"Invalid capitalized", UserRole("Subscriber"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.expected {
				t.Errorf("UserRole(%s).IsValid() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestUserRole_CognitoGroupName(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected string
	}{
		{"Guest cognito group", RoleGuest, "guest"},
		{"Subscriber cognito group", RoleSubscriber, "subscriber"},
		{"Artist cognito group", RoleArtist, "artist"},
		{"Admin cognito group", RoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.CognitoGroupName(); got != tt.expected {
				t.Errorf("UserRole(%s).CognitoGroupName() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestPermissionConstants(t *testing.T) {
	// Verify all expected permissions exist
	permissions := []Permission{
		PermissionUploadTracks,
		PermissionCreatePlaylists,
		PermissionCreatePublicPlaylists,
		PermissionFollowArtists,
		PermissionHaveFollowers,
		PermissionManageOwnContent,
		PermissionModerateContent,
		PermissionManageUsers,
	}

	if len(permissions) != 8 {
		t.Errorf("Expected 8 permissions, got %d", len(permissions))
	}
}

func TestRolePermissions(t *testing.T) {
	tests := []struct {
		name       string
		role       UserRole
		permission Permission
		expected   bool
	}{
		// Guest permissions
		{"Guest cannot upload", RoleGuest, PermissionUploadTracks, false},
		{"Guest cannot create playlists", RoleGuest, PermissionCreatePlaylists, false},
		{"Guest cannot follow", RoleGuest, PermissionFollowArtists, false},

		// Subscriber permissions
		{"Subscriber cannot upload", RoleSubscriber, PermissionUploadTracks, false},
		{"Subscriber can create playlists", RoleSubscriber, PermissionCreatePlaylists, true},
		{"Subscriber can create public playlists", RoleSubscriber, PermissionCreatePublicPlaylists, true},
		{"Subscriber can follow artists", RoleSubscriber, PermissionFollowArtists, true},
		{"Subscriber cannot have followers", RoleSubscriber, PermissionHaveFollowers, false},
		{"Subscriber can manage own content", RoleSubscriber, PermissionManageOwnContent, true},

		// Artist permissions
		{"Artist can upload", RoleArtist, PermissionUploadTracks, true},
		{"Artist can create playlists", RoleArtist, PermissionCreatePlaylists, true},
		{"Artist can create public playlists", RoleArtist, PermissionCreatePublicPlaylists, true},
		{"Artist can follow artists", RoleArtist, PermissionFollowArtists, true},
		{"Artist can have followers", RoleArtist, PermissionHaveFollowers, true},
		{"Artist can manage own content", RoleArtist, PermissionManageOwnContent, true},
		{"Artist cannot moderate", RoleArtist, PermissionModerateContent, false},

		// Admin permissions
		{"Admin can upload", RoleAdmin, PermissionUploadTracks, true},
		{"Admin can moderate", RoleAdmin, PermissionModerateContent, true},
		{"Admin can manage users", RoleAdmin, PermissionManageUsers, true},
		{"Admin can have followers", RoleAdmin, PermissionHaveFollowers, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.HasPermission(tt.permission); got != tt.expected {
				t.Errorf("UserRole(%s).HasPermission(%s) = %v, want %v", tt.role, tt.permission, got, tt.expected)
			}
		})
	}
}

func TestUserRole_CanUploadTracks(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"Guest cannot upload", RoleGuest, false},
		{"Subscriber cannot upload", RoleSubscriber, false},
		{"Artist can upload", RoleArtist, true},
		{"Admin can upload", RoleAdmin, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.CanUploadTracks(); got != tt.expected {
				t.Errorf("UserRole(%s).CanUploadTracks() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestUserRole_CanHaveFollowers(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"Guest cannot have followers", RoleGuest, false},
		{"Subscriber cannot have followers", RoleSubscriber, false},
		{"Artist can have followers", RoleArtist, true},
		{"Admin can have followers", RoleAdmin, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.CanHaveFollowers(); got != tt.expected {
				t.Errorf("UserRole(%s).CanHaveFollowers() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestUserRole_CanModerateContent(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"Guest cannot moderate", RoleGuest, false},
		{"Subscriber cannot moderate", RoleSubscriber, false},
		{"Artist cannot moderate", RoleArtist, false},
		{"Admin can moderate", RoleAdmin, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.CanModerateContent(); got != tt.expected {
				t.Errorf("UserRole(%s).CanModerateContent() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestDefaultUserRole(t *testing.T) {
	if DefaultUserRole() != RoleSubscriber {
		t.Errorf("DefaultUserRole() = %v, want %v", DefaultUserRole(), RoleSubscriber)
	}
}

func TestAllUserRoles(t *testing.T) {
	roles := AllUserRoles()
	if len(roles) != 4 {
		t.Errorf("AllUserRoles() returned %d roles, want 4", len(roles))
	}

	expected := map[UserRole]bool{RoleGuest: false, RoleSubscriber: false, RoleArtist: false, RoleAdmin: false}
	for _, r := range roles {
		if _, ok := expected[r]; !ok {
			t.Errorf("Unexpected role in AllUserRoles(): %s", r)
		}
		expected[r] = true
	}

	for role, found := range expected {
		if !found {
			t.Errorf("AllUserRoles() missing role: %s", role)
		}
	}
}
