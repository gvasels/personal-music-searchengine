package models

import (
	"testing"
)

func TestPlaylistVisibilityConstants(t *testing.T) {
	tests := []struct {
		name       string
		visibility PlaylistVisibility
		expected   string
	}{
		{"Private visibility", VisibilityPrivate, "private"},
		{"Unlisted visibility", VisibilityUnlisted, "unlisted"},
		{"Public visibility", VisibilityPublic, "public"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.visibility) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.visibility)
			}
		})
	}
}

func TestPlaylistVisibility_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		visibility PlaylistVisibility
		expected   bool
	}{
		{"Valid private", VisibilityPrivate, true},
		{"Valid unlisted", VisibilityUnlisted, true},
		{"Valid public", VisibilityPublic, true},
		{"Invalid empty", PlaylistVisibility(""), false},
		{"Invalid unknown", PlaylistVisibility("unknown"), false},
		{"Invalid capitalized", PlaylistVisibility("Private"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.visibility.IsValid(); got != tt.expected {
				t.Errorf("PlaylistVisibility(%s).IsValid() = %v, want %v", tt.visibility, got, tt.expected)
			}
		})
	}
}

func TestPlaylistVisibility_IsPubliclyAccessible(t *testing.T) {
	tests := []struct {
		name       string
		visibility PlaylistVisibility
		expected   bool
	}{
		{"Private is not publicly accessible", VisibilityPrivate, false},
		{"Unlisted is publicly accessible", VisibilityUnlisted, true},
		{"Public is publicly accessible", VisibilityPublic, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.visibility.IsPubliclyAccessible(); got != tt.expected {
				t.Errorf("PlaylistVisibility(%s).IsPubliclyAccessible() = %v, want %v", tt.visibility, got, tt.expected)
			}
		})
	}
}

func TestPlaylistVisibility_IsDiscoverable(t *testing.T) {
	tests := []struct {
		name       string
		visibility PlaylistVisibility
		expected   bool
	}{
		{"Private is not discoverable", VisibilityPrivate, false},
		{"Unlisted is not discoverable", VisibilityUnlisted, false},
		{"Public is discoverable", VisibilityPublic, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.visibility.IsDiscoverable(); got != tt.expected {
				t.Errorf("PlaylistVisibility(%s).IsDiscoverable() = %v, want %v", tt.visibility, got, tt.expected)
			}
		})
	}
}

func TestDefaultPlaylistVisibility(t *testing.T) {
	if DefaultPlaylistVisibility() != VisibilityPrivate {
		t.Errorf("DefaultPlaylistVisibility() = %v, want %v", DefaultPlaylistVisibility(), VisibilityPrivate)
	}
}

func TestAllPlaylistVisibilities(t *testing.T) {
	visibilities := AllPlaylistVisibilities()
	if len(visibilities) != 3 {
		t.Errorf("AllPlaylistVisibilities() returned %d visibilities, want 3", len(visibilities))
	}

	expected := map[PlaylistVisibility]bool{
		VisibilityPrivate:  false,
		VisibilityUnlisted: false,
		VisibilityPublic:   false,
	}
	for _, v := range visibilities {
		if _, ok := expected[v]; !ok {
			t.Errorf("Unexpected visibility in AllPlaylistVisibilities(): %s", v)
		}
		expected[v] = true
	}

	for visibility, found := range expected {
		if !found {
			t.Errorf("AllPlaylistVisibilities() missing visibility: %s", visibility)
		}
	}
}

func TestVisibilityFromIsPublic(t *testing.T) {
	tests := []struct {
		name     string
		isPublic bool
		expected PlaylistVisibility
	}{
		{"IsPublic true converts to public", true, VisibilityPublic},
		{"IsPublic false converts to private", false, VisibilityPrivate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VisibilityFromIsPublic(tt.isPublic); got != tt.expected {
				t.Errorf("VisibilityFromIsPublic(%v) = %v, want %v", tt.isPublic, got, tt.expected)
			}
		})
	}
}
