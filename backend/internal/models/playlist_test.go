package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlaylistStructFields verifies Playlist struct has all required fields
func TestPlaylistStructFields(t *testing.T) {
	playlist := Playlist{
		ID:            "playlist-123",
		UserID:        "user-456",
		Name:          "My Favorites",
		Description:   "Collection of my favorite songs",
		CoverArtKey:   "playlists/user-456/playlist-123/cover.jpg",
		TrackCount:    25,
		TotalDuration: 5400, // 90 minutes
		IsPublic:      false,
	}

	assert.Equal(t, "playlist-123", playlist.ID)
	assert.Equal(t, "user-456", playlist.UserID)
	assert.Equal(t, "My Favorites", playlist.Name)
	assert.Equal(t, "Collection of my favorite songs", playlist.Description)
	assert.Equal(t, "playlists/user-456/playlist-123/cover.jpg", playlist.CoverArtKey)
	assert.Equal(t, 25, playlist.TrackCount)
	assert.Equal(t, 5400, playlist.TotalDuration)
	assert.False(t, playlist.IsPublic)
}

// TestPlaylistJSONTags verifies JSON serialization
func TestPlaylistJSONTags(t *testing.T) {
	playlist := Playlist{
		ID:            "playlist-123",
		UserID:        "user-456",
		Name:          "My Favorites",
		Description:   "Collection of my favorite songs",
		TrackCount:    25,
		TotalDuration: 5400,
		IsPublic:      true,
	}

	jsonBytes, err := json.Marshal(playlist)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "id")
	assert.Contains(t, jsonMap, "userId")
	assert.Contains(t, jsonMap, "name")
	assert.Contains(t, jsonMap, "description")
	assert.Contains(t, jsonMap, "trackCount")
	assert.Contains(t, jsonMap, "totalDuration")
	assert.Contains(t, jsonMap, "isPublic")
}

// TestPlaylistJSONOmitEmpty verifies omitempty behavior
func TestPlaylistJSONOmitEmpty(t *testing.T) {
	playlist := Playlist{
		ID:     "playlist-123",
		UserID: "user-456",
		Name:   "My Favorites",
		// Leave optional fields empty
	}

	jsonBytes, err := json.Marshal(playlist)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Optional fields should not appear
	assert.NotContains(t, jsonMap, "description")
	assert.NotContains(t, jsonMap, "coverArtKey")
}

// TestNewPlaylistItem verifies DynamoDB item creation
func TestNewPlaylistItem(t *testing.T) {
	playlist := Playlist{
		ID:     "playlist-123",
		UserID: "user-456",
		Name:   "My Favorites",
	}

	item := NewPlaylistItem(playlist)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-456", item.PK)
	assert.Equal(t, "PLAYLIST#playlist-123", item.SK)
	assert.Equal(t, string(EntityPlaylist), item.Type)
}

// TestPlaylistToResponse verifies API response conversion
func TestPlaylistToResponse(t *testing.T) {
	now := time.Now()
	playlist := Playlist{
		ID:            "playlist-123",
		UserID:        "user-456",
		Name:          "My Favorites",
		Description:   "Collection of my favorite songs",
		TrackCount:    25,
		TotalDuration: 5465, // 1:31:05
		IsPublic:      true,
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	coverArtURL := "https://cdn.example.com/cover.jpg"
	response := playlist.ToResponse(coverArtURL)

	assert.Equal(t, "playlist-123", response.ID)
	assert.Equal(t, "My Favorites", response.Name)
	assert.Equal(t, "Collection of my favorite songs", response.Description)
	assert.Equal(t, coverArtURL, response.CoverArtURL)
	assert.Equal(t, 25, response.TrackCount)
	assert.Equal(t, 5465, response.TotalDuration)
	assert.Equal(t, "1:31:05", response.DurationStr)
	assert.True(t, response.IsPublic)
}

// TestPlaylistTrackStructFields verifies PlaylistTrack struct
func TestPlaylistTrackStructFields(t *testing.T) {
	now := time.Now()
	pt := PlaylistTrack{
		PlaylistID: "playlist-123",
		TrackID:    "track-456",
		Position:   5,
		AddedAt:    now,
	}

	assert.Equal(t, "playlist-123", pt.PlaylistID)
	assert.Equal(t, "track-456", pt.TrackID)
	assert.Equal(t, 5, pt.Position)
	assert.Equal(t, now, pt.AddedAt)
}

// TestNewPlaylistTrackItem verifies DynamoDB item creation
func TestNewPlaylistTrackItem(t *testing.T) {
	now := time.Now()
	pt := PlaylistTrack{
		PlaylistID: "playlist-123",
		TrackID:    "track-456",
		Position:   5,
		AddedAt:    now,
	}

	item := NewPlaylistTrackItem(pt)

	// Verify PK/SK patterns - position is zero-padded
	assert.Equal(t, "PLAYLIST#playlist-123", item.PK)
	assert.Equal(t, "POSITION#00000005", item.SK)
	assert.Equal(t, string(EntityPlaylistTrack), item.Type)
}

// TestNewPlaylistTrackItemPositionPadding verifies position padding
func TestNewPlaylistTrackItemPositionPadding(t *testing.T) {
	tests := []struct {
		name       string
		position   int
		expectedSK string
	}{
		{"single digit", 1, "POSITION#00000001"},
		{"double digit", 12, "POSITION#00000012"},
		{"triple digit", 123, "POSITION#00000123"},
		{"four digit", 1234, "POSITION#00001234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := PlaylistTrack{
				PlaylistID: "playlist-123",
				TrackID:    "track-456",
				Position:   tt.position,
			}

			item := NewPlaylistTrackItem(pt)
			assert.Equal(t, tt.expectedSK, item.SK)
		})
	}
}

// TestCreatePlaylistRequestFields verifies create request
func TestCreatePlaylistRequestFields(t *testing.T) {
	req := CreatePlaylistRequest{
		Name:        "My New Playlist",
		Description: "A great playlist",
		IsPublic:    true,
	}

	assert.Equal(t, "My New Playlist", req.Name)
	assert.Equal(t, "A great playlist", req.Description)
	assert.True(t, req.IsPublic)
}

// TestUpdatePlaylistRequestFields verifies update request
func TestUpdatePlaylistRequestFields(t *testing.T) {
	name := "Updated Name"
	description := "Updated description"
	isPublic := true

	req := UpdatePlaylistRequest{
		Name:        &name,
		Description: &description,
		IsPublic:    &isPublic,
	}

	assert.NotNil(t, req.Name)
	assert.Equal(t, "Updated Name", *req.Name)
	assert.NotNil(t, req.Description)
	assert.Equal(t, "Updated description", *req.Description)
	assert.NotNil(t, req.IsPublic)
	assert.True(t, *req.IsPublic)
}

// TestPlaylistFilterFields verifies filter struct
func TestPlaylistFilterFields(t *testing.T) {
	filter := PlaylistFilter{
		SortBy:    "name",
		SortOrder: "asc",
		Limit:     20,
		LastKey:   "abc123",
	}

	assert.Equal(t, "name", filter.SortBy)
	assert.Equal(t, "asc", filter.SortOrder)
	assert.Equal(t, 20, filter.Limit)
	assert.Equal(t, "abc123", filter.LastKey)
}

// TestAddTracksToPlaylistRequestFields verifies add tracks request
func TestAddTracksToPlaylistRequestFields(t *testing.T) {
	position := 5
	req := AddTracksToPlaylistRequest{
		TrackIDs: []string{"track-1", "track-2", "track-3"},
		Position: &position,
	}

	assert.Equal(t, 3, len(req.TrackIDs))
	assert.Contains(t, req.TrackIDs, "track-1")
	assert.NotNil(t, req.Position)
	assert.Equal(t, 5, *req.Position)
}

// TestRemoveTracksFromPlaylistRequestFields verifies remove tracks request
func TestRemoveTracksFromPlaylistRequestFields(t *testing.T) {
	req := RemoveTracksFromPlaylistRequest{
		TrackIDs: []string{"track-1", "track-2"},
	}

	assert.Equal(t, 2, len(req.TrackIDs))
	assert.Contains(t, req.TrackIDs, "track-1")
	assert.Contains(t, req.TrackIDs, "track-2")
}

// TestPlaylistResponseFields verifies response struct
func TestPlaylistResponseFields(t *testing.T) {
	now := time.Now()
	response := PlaylistResponse{
		ID:            "playlist-123",
		Name:          "My Favorites",
		Description:   "Collection of my favorite songs",
		CoverArtURL:   "https://cdn.example.com/cover.jpg",
		TrackCount:    25,
		TotalDuration: 5400,
		DurationStr:   "1:30:00",
		IsPublic:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	assert.Equal(t, "playlist-123", response.ID)
	assert.Equal(t, "My Favorites", response.Name)
	assert.Equal(t, "1:30:00", response.DurationStr)
	assert.Equal(t, 25, response.TrackCount)
}
