package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTagStructFields verifies Tag struct has all required fields
func TestTagStructFields(t *testing.T) {
	tag := Tag{
		UserID:     "user-456",
		Name:       "favorites",
		Color:      "#FF5733",
		TrackCount: 25,
	}

	assert.Equal(t, "user-456", tag.UserID)
	assert.Equal(t, "favorites", tag.Name)
	assert.Equal(t, "#FF5733", tag.Color)
	assert.Equal(t, 25, tag.TrackCount)
}

// TestTagJSONTags verifies JSON serialization
func TestTagJSONTags(t *testing.T) {
	tag := Tag{
		UserID:     "user-456",
		Name:       "favorites",
		Color:      "#FF5733",
		TrackCount: 25,
	}

	jsonBytes, err := json.Marshal(tag)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "userId")
	assert.Contains(t, jsonMap, "name")
	assert.Contains(t, jsonMap, "color")
	assert.Contains(t, jsonMap, "trackCount")
}

// TestTagJSONOmitEmpty verifies omitempty behavior
func TestTagJSONOmitEmpty(t *testing.T) {
	tag := Tag{
		UserID: "user-456",
		Name:   "favorites",
		// Leave optional fields empty
	}

	jsonBytes, err := json.Marshal(tag)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Optional fields should not appear
	assert.NotContains(t, jsonMap, "color")
}

// TestNewTagItem verifies DynamoDB item creation
func TestNewTagItem(t *testing.T) {
	tag := Tag{
		UserID: "user-456",
		Name:   "favorites",
		Color:  "#FF5733",
	}

	item := NewTagItem(tag)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-456", item.PK)
	assert.Equal(t, "TAG#favorites", item.SK)
	assert.Equal(t, string(EntityTag), item.Type)
}

// TestTagToResponse verifies API response conversion
func TestTagToResponse(t *testing.T) {
	now := time.Now()
	tag := Tag{
		UserID:     "user-456",
		Name:       "favorites",
		Color:      "#FF5733",
		TrackCount: 25,
	}
	tag.CreatedAt = now
	tag.UpdatedAt = now

	response := tag.ToResponse()

	assert.Equal(t, "favorites", response.Name)
	assert.Equal(t, "#FF5733", response.Color)
	assert.Equal(t, 25, response.TrackCount)
}

// TestTrackTagStructFields verifies TrackTag struct
func TestTrackTagStructFields(t *testing.T) {
	now := time.Now()
	tt := TrackTag{
		UserID:  "user-456",
		TrackID: "track-789",
		TagName: "favorites",
		AddedAt: now,
	}

	assert.Equal(t, "user-456", tt.UserID)
	assert.Equal(t, "track-789", tt.TrackID)
	assert.Equal(t, "favorites", tt.TagName)
	assert.Equal(t, now, tt.AddedAt)
}

// TestNewTrackTagItem verifies DynamoDB item creation
func TestNewTrackTagItem(t *testing.T) {
	now := time.Now()
	tt := TrackTag{
		UserID:  "user-456",
		TrackID: "track-789",
		TagName: "favorites",
		AddedAt: now,
	}

	item := NewTrackTagItem(tt)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-456#TRACK#track-789", item.PK)
	assert.Equal(t, "TAG#favorites", item.SK)
	assert.Equal(t, string(EntityTrackTag), item.Type)

	// Verify GSI1 for tag lookup
	assert.Equal(t, "USER#user-456#TAG#favorites", item.GSI1PK)
	assert.Equal(t, "TRACK#track-789", item.GSI1SK)
}

// TestCreateTagRequestFields verifies create request
func TestCreateTagRequestFields(t *testing.T) {
	req := CreateTagRequest{
		Name:  "new-tag",
		Color: "#00FF00",
	}

	assert.Equal(t, "new-tag", req.Name)
	assert.Equal(t, "#00FF00", req.Color)
}

// TestUpdateTagRequestFields verifies update request
func TestUpdateTagRequestFields(t *testing.T) {
	name := "updated-tag"
	color := "#0000FF"

	req := UpdateTagRequest{
		Name:  &name,
		Color: &color,
	}

	assert.NotNil(t, req.Name)
	assert.Equal(t, "updated-tag", *req.Name)
	assert.NotNil(t, req.Color)
	assert.Equal(t, "#0000FF", *req.Color)
}

// TestTagFilterFields verifies filter struct
func TestTagFilterFields(t *testing.T) {
	filter := TagFilter{
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

// TestAddTagsToTrackRequestFields verifies add tags request
func TestAddTagsToTrackRequestFields(t *testing.T) {
	req := AddTagsToTrackRequest{
		Tags: []string{"favorites", "rock", "2024"},
	}

	assert.Equal(t, 3, len(req.Tags))
	assert.Contains(t, req.Tags, "favorites")
	assert.Contains(t, req.Tags, "rock")
	assert.Contains(t, req.Tags, "2024")
}

// TestTagResponseFields verifies response struct
func TestTagResponseFields(t *testing.T) {
	now := time.Now()
	response := TagResponse{
		Name:       "favorites",
		Color:      "#FF5733",
		TrackCount: 25,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	assert.Equal(t, "favorites", response.Name)
	assert.Equal(t, "#FF5733", response.Color)
	assert.Equal(t, 25, response.TrackCount)
}

// TestTrackTagItemFields verifies TrackTagItem struct for DynamoDB
func TestTrackTagItemFields(t *testing.T) {
	now := time.Now()
	item := TrackTagItem{
		DynamoDBItem: DynamoDBItem{
			PK:     "USER#user-456#TRACK#track-789",
			SK:     "TAG#favorites",
			GSI1PK: "USER#user-456#TAG#favorites",
			GSI1SK: "TRACK#track-789",
			Type:   string(EntityTrackTag),
		},
		TrackTag: TrackTag{
			UserID:  "user-456",
			TrackID: "track-789",
			TagName: "favorites",
			AddedAt: now,
		},
	}

	assert.Equal(t, "USER#user-456#TRACK#track-789", item.PK)
	assert.Equal(t, "TAG#favorites", item.SK)
	assert.Equal(t, "USER#user-456#TAG#favorites", item.GSI1PK)
	assert.Equal(t, "TRACK#track-789", item.GSI1SK)
	assert.Equal(t, "favorites", item.TagName)
}
