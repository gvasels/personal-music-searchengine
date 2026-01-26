package models

// PlaylistVisibility represents the visibility level of a playlist.
type PlaylistVisibility string

const (
	// VisibilityPrivate means only the owner can see the playlist.
	VisibilityPrivate PlaylistVisibility = "private"

	// VisibilityUnlisted means anyone with the link can access, but it's not discoverable.
	VisibilityUnlisted PlaylistVisibility = "unlisted"

	// VisibilityPublic means anyone can find and access the playlist.
	VisibilityPublic PlaylistVisibility = "public"
)

// validVisibilities contains all valid visibility values for validation.
var validVisibilities = map[PlaylistVisibility]bool{
	VisibilityPrivate:  true,
	VisibilityUnlisted: true,
	VisibilityPublic:   true,
}

// IsValid returns true if the visibility is a valid option.
func (v PlaylistVisibility) IsValid() bool {
	return validVisibilities[v]
}

// IsPubliclyAccessible returns true if the playlist can be accessed by non-owners.
// Both public and unlisted playlists are accessible if you have the link.
func (v PlaylistVisibility) IsPubliclyAccessible() bool {
	return v == VisibilityPublic || v == VisibilityUnlisted
}

// IsDiscoverable returns true if the playlist appears in search results and public listings.
// Only public playlists are discoverable.
func (v PlaylistVisibility) IsDiscoverable() bool {
	return v == VisibilityPublic
}

// DefaultPlaylistVisibility returns the default visibility for new playlists.
func DefaultPlaylistVisibility() PlaylistVisibility {
	return VisibilityPrivate
}

// AllPlaylistVisibilities returns a slice of all valid visibility options.
func AllPlaylistVisibilities() []PlaylistVisibility {
	return []PlaylistVisibility{VisibilityPrivate, VisibilityUnlisted, VisibilityPublic}
}

// VisibilityFromIsPublic converts the legacy IsPublic bool to PlaylistVisibility.
// Used for backwards compatibility during migration.
func VisibilityFromIsPublic(isPublic bool) PlaylistVisibility {
	if isPublic {
		return VisibilityPublic
	}
	return VisibilityPrivate
}
