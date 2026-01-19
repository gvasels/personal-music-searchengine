package models

import (
	"fmt"
	"time"
)

// Playlist represents a user-created playlist
type Playlist struct {
	ID            string `json:"id" dynamodbav:"id"`
	UserID        string `json:"userId" dynamodbav:"userId"`
	Name          string `json:"name" dynamodbav:"name"`
	Description   string `json:"description,omitempty" dynamodbav:"description,omitempty"`
	CoverArtKey   string `json:"coverArtKey,omitempty" dynamodbav:"coverArtKey,omitempty"`
	TrackCount    int    `json:"trackCount" dynamodbav:"trackCount"`
	TotalDuration int    `json:"totalDuration" dynamodbav:"totalDuration"` // seconds
	IsPublic      bool   `json:"isPublic" dynamodbav:"isPublic"`
	Timestamps
}

// PlaylistItem represents a Playlist in DynamoDB single-table design
type PlaylistItem struct {
	DynamoDBItem
	Playlist
}

// NewPlaylistItem creates a DynamoDB item for a playlist
func NewPlaylistItem(playlist Playlist) PlaylistItem {
	return PlaylistItem{
		DynamoDBItem: DynamoDBItem{
			PK:   fmt.Sprintf("USER#%s", playlist.UserID),
			SK:   fmt.Sprintf("PLAYLIST#%s", playlist.ID),
			Type: string(EntityPlaylist),
		},
		Playlist: playlist,
	}
}

// PlaylistTrack represents a track within a playlist (with position)
type PlaylistTrack struct {
	PlaylistID string    `json:"playlistId" dynamodbav:"playlistId"`
	TrackID    string    `json:"trackId" dynamodbav:"trackId"`
	Position   int       `json:"position" dynamodbav:"position"`
	AddedAt    time.Time `json:"addedAt" dynamodbav:"addedAt"`
}

// PlaylistTrackItem represents a PlaylistTrack in DynamoDB single-table design
type PlaylistTrackItem struct {
	DynamoDBItem
	PlaylistTrack
}

// NewPlaylistTrackItem creates a DynamoDB item for a playlist track
func NewPlaylistTrackItem(pt PlaylistTrack) PlaylistTrackItem {
	return PlaylistTrackItem{
		DynamoDBItem: DynamoDBItem{
			PK:   fmt.Sprintf("PLAYLIST#%s", pt.PlaylistID),
			SK:   fmt.Sprintf("POSITION#%08d", pt.Position),
			Type: string(EntityPlaylistTrack),
		},
		PlaylistTrack: pt,
	}
}

// CreatePlaylistRequest represents a request to create a playlist
type CreatePlaylistRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=200"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsPublic    bool   `json:"isPublic"`
}

// UpdatePlaylistRequest represents a request to update a playlist
type UpdatePlaylistRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsPublic    *bool   `json:"isPublic,omitempty"`
}

// AddTracksToPlaylistRequest represents a request to add tracks to a playlist
type AddTracksToPlaylistRequest struct {
	TrackIDs []string `json:"trackIds" validate:"required,min=1,max=100,dive,uuid"`
	Position *int     `json:"position,omitempty" validate:"omitempty,min=0"`
}

// RemoveTracksFromPlaylistRequest represents a request to remove tracks from a playlist
type RemoveTracksFromPlaylistRequest struct {
	TrackIDs []string `json:"trackIds" validate:"required,min=1,max=100,dive,uuid"`
}

// ReorderPlaylistTracksRequest represents a request to reorder tracks in a playlist
type ReorderPlaylistTracksRequest struct {
	TrackID      string `json:"trackId" validate:"required,uuid"`
	NewPosition  int    `json:"newPosition" validate:"required,min=0"`
}

// PlaylistResponse represents a playlist in API responses
type PlaylistResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	CoverArtURL   string    `json:"coverArtUrl,omitempty"`
	TrackCount    int       `json:"trackCount"`
	TotalDuration int       `json:"totalDuration"`
	DurationStr   string    `json:"durationStr"`
	IsPublic      bool      `json:"isPublic"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ToResponse converts a Playlist to a PlaylistResponse
func (p *Playlist) ToResponse(coverArtURL string) PlaylistResponse {
	return PlaylistResponse{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		CoverArtURL:   coverArtURL,
		TrackCount:    p.TrackCount,
		TotalDuration: p.TotalDuration,
		DurationStr:   formatDuration(p.TotalDuration),
		IsPublic:      p.IsPublic,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// PlaylistWithTracks represents a playlist with its tracks
type PlaylistWithTracks struct {
	Playlist PlaylistResponse `json:"playlist"`
	Tracks   []TrackResponse  `json:"tracks"`
}

// PlaylistFilter represents filter options for listing playlists
type PlaylistFilter struct {
	SortBy    string `query:"sortBy"`    // name, createdAt, updatedAt, trackCount
	SortOrder string `query:"sortOrder"` // asc, desc
	Limit     int    `query:"limit"`
	LastKey   string `query:"lastKey"`
}
