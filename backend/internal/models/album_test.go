package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlbumStructFields verifies Album struct has all required fields
func TestAlbumStructFields(t *testing.T) {
	album := Album{
		ID:            "album-123",
		UserID:        "user-456",
		Title:         "Test Album",
		Artist:        "Test Artist",
		AlbumArtist:   "Album Artist",
		Genre:         "Rock",
		Year:          2024,
		CoverArtKey:   "covers/user-456/album-123/cover.jpg",
		TrackCount:    12,
		TotalDuration: 3600, // 1 hour
		DiscCount:     2,
	}

	assert.Equal(t, "album-123", album.ID)
	assert.Equal(t, "user-456", album.UserID)
	assert.Equal(t, "Test Album", album.Title)
	assert.Equal(t, "Test Artist", album.Artist)
	assert.Equal(t, "Album Artist", album.AlbumArtist)
	assert.Equal(t, "Rock", album.Genre)
	assert.Equal(t, 2024, album.Year)
	assert.Equal(t, "covers/user-456/album-123/cover.jpg", album.CoverArtKey)
	assert.Equal(t, 12, album.TrackCount)
	assert.Equal(t, 3600, album.TotalDuration)
	assert.Equal(t, 2, album.DiscCount)
}

// TestAlbumJSONTags verifies JSON serialization
func TestAlbumJSONTags(t *testing.T) {
	album := Album{
		ID:            "album-123",
		UserID:        "user-456",
		Title:         "Test Album",
		Artist:        "Test Artist",
		AlbumArtist:   "Album Artist",
		Genre:         "Rock",
		Year:          2024,
		TrackCount:    12,
		TotalDuration: 3600,
	}

	jsonBytes, err := json.Marshal(album)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "id")
	assert.Contains(t, jsonMap, "userId")
	assert.Contains(t, jsonMap, "title")
	assert.Contains(t, jsonMap, "artist")
	assert.Contains(t, jsonMap, "albumArtist")
	assert.Contains(t, jsonMap, "genre")
	assert.Contains(t, jsonMap, "year")
	assert.Contains(t, jsonMap, "trackCount")
	assert.Contains(t, jsonMap, "totalDuration")
}

// TestAlbumJSONOmitEmpty verifies omitempty behavior
func TestAlbumJSONOmitEmpty(t *testing.T) {
	album := Album{
		ID:     "album-123",
		UserID: "user-456",
		Title:  "Test Album",
		Artist: "Test Artist",
		// Leave optional fields empty
	}

	jsonBytes, err := json.Marshal(album)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Optional fields should not appear
	assert.NotContains(t, jsonMap, "albumArtist")
	assert.NotContains(t, jsonMap, "genre")
	assert.NotContains(t, jsonMap, "coverArtKey")
}

// TestNewAlbumItem verifies DynamoDB item creation
func TestNewAlbumItem(t *testing.T) {
	album := Album{
		ID:     "album-123",
		UserID: "user-456",
		Title:  "Test Album",
		Artist: "Test Artist",
		Year:   2024,
	}

	item := NewAlbumItem(album)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-456", item.PK)
	assert.Equal(t, "ALBUM#album-123", item.SK)
	assert.Equal(t, string(EntityAlbum), item.Type)

	// Verify GSI1 for artist queries
	assert.Equal(t, "USER#user-456#ARTIST#Test Artist", item.GSI1PK)
	assert.Equal(t, "ALBUM#2024", item.GSI1SK)
}

// TestNewAlbumItemWithEmptyArtist verifies GSI handling when artist is empty
func TestNewAlbumItemWithEmptyArtist(t *testing.T) {
	album := Album{
		ID:     "album-123",
		UserID: "user-456",
		Title:  "Test Album",
		Artist: "", // Empty artist
	}

	item := NewAlbumItem(album)

	// GSI1 should not be set if artist is empty
	assert.Empty(t, item.GSI1PK)
	assert.Empty(t, item.GSI1SK)
}

// TestAlbumToResponse verifies API response conversion
func TestAlbumToResponse(t *testing.T) {
	now := time.Now()
	album := Album{
		ID:            "album-123",
		UserID:        "user-456",
		Title:         "Test Album",
		Artist:        "Test Artist",
		AlbumArtist:   "Album Artist",
		Genre:         "Rock",
		Year:          2024,
		TrackCount:    12,
		TotalDuration: 3665, // 1:01:05
		DiscCount:     2,
	}
	album.CreatedAt = now
	album.UpdatedAt = now

	coverArtURL := "https://cdn.example.com/cover.jpg"
	response := album.ToResponse(coverArtURL)

	assert.Equal(t, "album-123", response.ID)
	assert.Equal(t, "Test Album", response.Title)
	assert.Equal(t, "Test Artist", response.Artist)
	assert.Equal(t, "Album Artist", response.AlbumArtist)
	assert.Equal(t, "Rock", response.Genre)
	assert.Equal(t, 2024, response.Year)
	assert.Equal(t, coverArtURL, response.CoverArtURL)
	assert.Equal(t, 12, response.TrackCount)
	assert.Equal(t, 3665, response.TotalDuration)
	assert.Equal(t, "1:01:05", response.DurationStr)
	assert.Equal(t, 2, response.DiscCount)
}

// TestAlbumFilterFields verifies AlbumFilter struct
func TestAlbumFilterFields(t *testing.T) {
	filter := AlbumFilter{
		Artist:    "Test Artist",
		Genre:     "Rock",
		Year:      2024,
		SortBy:    "title",
		SortOrder: "asc",
		Limit:     20,
		LastKey:   "abc123",
	}

	assert.Equal(t, "Test Artist", filter.Artist)
	assert.Equal(t, "Rock", filter.Genre)
	assert.Equal(t, 2024, filter.Year)
	assert.Equal(t, "title", filter.SortBy)
	assert.Equal(t, "asc", filter.SortOrder)
	assert.Equal(t, 20, filter.Limit)
	assert.Equal(t, "abc123", filter.LastKey)
}

// TestArtistSummaryFields verifies ArtistSummary struct
func TestArtistSummaryFields(t *testing.T) {
	artist := ArtistSummary{
		Name:        "Test Artist",
		TrackCount:  50,
		AlbumCount:  5,
		CoverArtURL: "https://cdn.example.com/artist.jpg",
	}

	assert.Equal(t, "Test Artist", artist.Name)
	assert.Equal(t, 50, artist.TrackCount)
	assert.Equal(t, 5, artist.AlbumCount)
	assert.Equal(t, "https://cdn.example.com/artist.jpg", artist.CoverArtURL)
}

// TestAlbumResponseFields verifies AlbumResponse struct
func TestAlbumResponseFields(t *testing.T) {
	now := time.Now()
	response := AlbumResponse{
		ID:            "album-123",
		Title:         "Test Album",
		Artist:        "Test Artist",
		AlbumArtist:   "Album Artist",
		Genre:         "Rock",
		Year:          2024,
		CoverArtURL:   "https://cdn.example.com/cover.jpg",
		TrackCount:    12,
		TotalDuration: 3600,
		DurationStr:   "1:00:00",
		DiscCount:     2,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	assert.Equal(t, "album-123", response.ID)
	assert.Equal(t, "Test Album", response.Title)
	assert.Equal(t, "Test Artist", response.Artist)
	assert.Equal(t, "1:00:00", response.DurationStr)
	assert.Equal(t, 12, response.TrackCount)
}
