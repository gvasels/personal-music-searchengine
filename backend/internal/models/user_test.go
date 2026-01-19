package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserStructFields verifies User struct has all required fields
func TestUserStructFields(t *testing.T) {
	user := User{
		ID:            "user-123",
		Email:         "test@example.com",
		DisplayName:   "Test User",
		AvatarURL:     "https://example.com/avatar.jpg",
		StorageUsed:   104857600, // 100 MB
		StorageLimit:  10737418240, // 10 GB
		TrackCount:    50,
		AlbumCount:    10,
		PlaylistCount: 5,
	}

	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.DisplayName)
	assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
	assert.Equal(t, int64(104857600), user.StorageUsed)
	assert.Equal(t, int64(10737418240), user.StorageLimit)
	assert.Equal(t, 50, user.TrackCount)
	assert.Equal(t, 10, user.AlbumCount)
	assert.Equal(t, 5, user.PlaylistCount)
}

// TestUserJSONTags verifies JSON serialization
func TestUserJSONTags(t *testing.T) {
	user := User{
		ID:            "user-123",
		Email:         "test@example.com",
		DisplayName:   "Test User",
		StorageUsed:   104857600,
		StorageLimit:  10737418240,
		TrackCount:    50,
		AlbumCount:    10,
		PlaylistCount: 5,
	}

	jsonBytes, err := json.Marshal(user)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "id")
	assert.Contains(t, jsonMap, "email")
	assert.Contains(t, jsonMap, "displayName")
	assert.Contains(t, jsonMap, "storageUsed")
	assert.Contains(t, jsonMap, "storageLimit")
	assert.Contains(t, jsonMap, "trackCount")
	assert.Contains(t, jsonMap, "albumCount")
	assert.Contains(t, jsonMap, "playlistCount")
}

// TestUserJSONOmitEmpty verifies omitempty behavior
func TestUserJSONOmitEmpty(t *testing.T) {
	user := User{
		ID:    "user-123",
		Email: "test@example.com",
		// Leave optional fields empty
	}

	jsonBytes, err := json.Marshal(user)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Optional fields should not appear
	assert.NotContains(t, jsonMap, "avatarUrl")
}

// TestNewUserItem verifies DynamoDB item creation
func TestNewUserItem(t *testing.T) {
	user := User{
		ID:    "user-123",
		Email: "test@example.com",
	}

	item := NewUserItem(user)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-123", item.PK)
	assert.Equal(t, "PROFILE", item.SK)
	assert.Equal(t, string(EntityUser), item.Type)
}

// TestUserToResponse verifies API response conversion
func TestUserToResponse(t *testing.T) {
	now := time.Now()
	user := User{
		ID:            "user-123",
		Email:         "test@example.com",
		DisplayName:   "Test User",
		AvatarURL:     "https://example.com/avatar.jpg",
		StorageUsed:   5368709120, // 5 GB
		StorageLimit:  10737418240, // 10 GB
		TrackCount:    50,
		AlbumCount:    10,
		PlaylistCount: 5,
	}
	user.CreatedAt = now

	response := user.ToResponse()

	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "Test User", response.DisplayName)
	assert.Equal(t, "https://example.com/avatar.jpg", response.AvatarURL)
	assert.Equal(t, int64(5368709120), response.StorageUsed)
	assert.Equal(t, int64(10737418240), response.StorageLimit)
	assert.Equal(t, 50, response.TrackCount)
	assert.Equal(t, 10, response.AlbumCount)
	assert.Equal(t, 5, response.PlaylistCount)
}

// TestUserResponseStorageValues verifies storage values are passed through
func TestUserResponseStorageValues(t *testing.T) {
	tests := []struct {
		name         string
		storageUsed  int64
		storageLimit int64
	}{
		{
			name:         "MB used, GB limit",
			storageUsed:  104857600,   // 100 MB
			storageLimit: 10737418240, // 10 GB
		},
		{
			name:         "GB used, GB limit",
			storageUsed:  5368709120,  // 5 GB
			storageLimit: 10737418240, // 10 GB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{
				ID:           "user-123",
				StorageUsed:  tt.storageUsed,
				StorageLimit: tt.storageLimit,
			}

			response := user.ToResponse()
			assert.Equal(t, tt.storageUsed, response.StorageUsed)
			assert.Equal(t, tt.storageLimit, response.StorageLimit)
		})
	}
}

// TestUpdateUserRequestFields verifies update request
func TestUpdateUserRequestFields(t *testing.T) {
	displayName := "New Name"
	avatarURL := "https://example.com/new-avatar.jpg"

	req := UpdateUserRequest{
		DisplayName: &displayName,
		AvatarURL:   &avatarURL,
	}

	assert.NotNil(t, req.DisplayName)
	assert.Equal(t, "New Name", *req.DisplayName)
	assert.NotNil(t, req.AvatarURL)
	assert.Equal(t, "https://example.com/new-avatar.jpg", *req.AvatarURL)
}

// TestUserResponseFields verifies UserResponse struct
func TestUserResponseFields(t *testing.T) {
	now := time.Now()
	response := UserResponse{
		ID:            "user-123",
		Email:         "test@example.com",
		DisplayName:   "Test User",
		AvatarURL:     "https://example.com/avatar.jpg",
		StorageUsed:   5368709120,
		StorageLimit:  10737418240,
		TrackCount:    50,
		AlbumCount:    10,
		PlaylistCount: 5,
		CreatedAt:     now,
	}

	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "Test User", response.DisplayName)
	assert.Equal(t, int64(5368709120), response.StorageUsed)
	assert.Equal(t, int64(10737418240), response.StorageLimit)
}
