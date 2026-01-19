package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUploadStructFields verifies Upload struct has all required fields
func TestUploadStructFields(t *testing.T) {
	now := time.Now()
	upload := Upload{
		ID:          "upload-123",
		UserID:      "user-456",
		FileName:    "song.mp3",
		FileSize:    5242880,
		ContentType: "audio/mpeg",
		S3Key:       "uploads/user-456/upload-123/song.mp3",
		Status:      UploadStatusPending,
		ErrorMsg:    "",
		TrackID:     "",
		CompletedAt: &now,
	}

	assert.Equal(t, "upload-123", upload.ID)
	assert.Equal(t, "user-456", upload.UserID)
	assert.Equal(t, "song.mp3", upload.FileName)
	assert.Equal(t, int64(5242880), upload.FileSize)
	assert.Equal(t, "audio/mpeg", upload.ContentType)
	assert.Equal(t, "uploads/user-456/upload-123/song.mp3", upload.S3Key)
	assert.Equal(t, UploadStatusPending, upload.Status)
	assert.Equal(t, "", upload.ErrorMsg)
	assert.Equal(t, "", upload.TrackID)
	assert.NotNil(t, upload.CompletedAt)
}

// TestUploadJSONTags verifies JSON serialization
func TestUploadJSONTags(t *testing.T) {
	upload := Upload{
		ID:          "upload-123",
		UserID:      "user-456",
		FileName:    "song.mp3",
		FileSize:    5242880,
		ContentType: "audio/mpeg",
		Status:      UploadStatusPending,
	}

	jsonBytes, err := json.Marshal(upload)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	assert.Contains(t, jsonMap, "id")
	assert.Contains(t, jsonMap, "userId")
	assert.Contains(t, jsonMap, "fileName")
	assert.Contains(t, jsonMap, "fileSize")
	assert.Contains(t, jsonMap, "contentType")
	assert.Contains(t, jsonMap, "status")
}

// TestUploadJSONOmitEmpty verifies omitempty behavior
func TestUploadJSONOmitEmpty(t *testing.T) {
	upload := Upload{
		ID:     "upload-123",
		UserID: "user-456",
		Status: UploadStatusPending,
		// Leave optional fields empty
	}

	jsonBytes, err := json.Marshal(upload)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Optional fields with omitempty should not appear
	assert.NotContains(t, jsonMap, "errorMsg")
	assert.NotContains(t, jsonMap, "trackId")
	assert.NotContains(t, jsonMap, "completedAt")
}

// TestNewUploadItem verifies DynamoDB item creation
func TestNewUploadItem(t *testing.T) {
	now := time.Now()
	upload := Upload{
		ID:     "upload-123",
		UserID: "user-456",
		Status: UploadStatusPending,
	}
	upload.CreatedAt = now

	item := NewUploadItem(upload)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-456", item.PK)
	assert.Equal(t, "UPLOAD#upload-123", item.SK)
	assert.Equal(t, string(EntityUpload), item.Type)

	// Verify GSI1 for status queries
	assert.Equal(t, "UPLOAD#STATUS#PENDING", item.GSI1PK)
	assert.Equal(t, now.Format(time.RFC3339), item.GSI1SK)
}

// TestNewUploadItemStatusGSI verifies GSI updates with different statuses
func TestNewUploadItemStatusGSI(t *testing.T) {
	tests := []struct {
		name           string
		status         UploadStatus
		expectedGSI1PK string
	}{
		{"PENDING", UploadStatusPending, "UPLOAD#STATUS#PENDING"},
		{"PROCESSING", UploadStatusProcessing, "UPLOAD#STATUS#PROCESSING"},
		{"COMPLETED", UploadStatusCompleted, "UPLOAD#STATUS#COMPLETED"},
		{"FAILED", UploadStatusFailed, "UPLOAD#STATUS#FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upload := Upload{
				ID:     "upload-123",
				UserID: "user-456",
				Status: tt.status,
			}
			upload.CreatedAt = time.Now()

			item := NewUploadItem(upload)
			assert.Equal(t, tt.expectedGSI1PK, item.GSI1PK)
		})
	}
}

// TestUploadToResponse verifies API response conversion
func TestUploadToResponse(t *testing.T) {
	now := time.Now()
	upload := Upload{
		ID:          "upload-123",
		UserID:      "user-456",
		FileName:    "song.mp3",
		FileSize:    5505024, // ~5.25 MB
		ContentType: "audio/mpeg",
		Status:      UploadStatusCompleted,
		TrackID:     "track-789",
		CompletedAt: &now,
	}
	upload.CreatedAt = now

	response := upload.ToResponse()

	assert.Equal(t, "upload-123", response.ID)
	assert.Equal(t, "song.mp3", response.FileName)
	assert.Equal(t, int64(5505024), response.FileSize)
	assert.Equal(t, "5.25 MB", response.FileSizeStr)
	assert.Equal(t, "audio/mpeg", response.ContentType)
	assert.Equal(t, UploadStatusCompleted, response.Status)
	assert.Equal(t, "track-789", response.TrackID)
	assert.NotNil(t, response.CompletedAt)
}

// TestPresignedUploadRequestFields verifies request struct
func TestPresignedUploadRequestFields(t *testing.T) {
	req := PresignedUploadRequest{
		FileName:    "song.mp3",
		FileSize:    5242880,
		ContentType: "audio/mpeg",
	}

	assert.Equal(t, "song.mp3", req.FileName)
	assert.Equal(t, int64(5242880), req.FileSize)
	assert.Equal(t, "audio/mpeg", req.ContentType)
}

// TestPresignedUploadResponseFields verifies response struct
func TestPresignedUploadResponseFields(t *testing.T) {
	expiresAt := time.Now().Add(15 * time.Minute)
	response := PresignedUploadResponse{
		UploadID:    "upload-123",
		UploadURL:   "https://s3.amazonaws.com/bucket/key?signature=...",
		Fields:      map[string]string{"key": "value"},
		ExpiresAt:   expiresAt,
		MaxFileSize: 524288000,
	}

	assert.Equal(t, "upload-123", response.UploadID)
	assert.Contains(t, response.UploadURL, "s3.amazonaws.com")
	assert.Equal(t, map[string]string{"key": "value"}, response.Fields)
	assert.Equal(t, expiresAt, response.ExpiresAt)
	assert.Equal(t, int64(524288000), response.MaxFileSize)
}

// TestConfirmUploadRequestFields verifies confirm request
func TestConfirmUploadRequestFields(t *testing.T) {
	req := ConfirmUploadRequest{
		UploadID: "upload-123",
	}

	assert.Equal(t, "upload-123", req.UploadID)
}

// TestConfirmUploadResponseFields verifies confirm response
func TestConfirmUploadResponseFields(t *testing.T) {
	response := ConfirmUploadResponse{
		UploadID: "upload-123",
		Status:   UploadStatusProcessing,
		Message:  "Upload processing started",
	}

	assert.Equal(t, "upload-123", response.UploadID)
	assert.Equal(t, UploadStatusProcessing, response.Status)
	assert.Equal(t, "Upload processing started", response.Message)
}

// TestUploadFilterFields verifies filter struct
func TestUploadFilterFields(t *testing.T) {
	filter := UploadFilter{
		Status:    UploadStatusCompleted,
		SortBy:    "createdAt",
		SortOrder: "desc",
		Limit:     20,
		LastKey:   "abc123",
	}

	assert.Equal(t, UploadStatusCompleted, filter.Status)
	assert.Equal(t, "createdAt", filter.SortBy)
	assert.Equal(t, "desc", filter.SortOrder)
	assert.Equal(t, 20, filter.Limit)
	assert.Equal(t, "abc123", filter.LastKey)
}

// TestUploadMetadataFields verifies metadata extraction result
func TestUploadMetadataFields(t *testing.T) {
	metadata := UploadMetadata{
		Title:       "Test Song",
		Artist:      "Test Artist",
		AlbumArtist: "Album Artist",
		Album:       "Test Album",
		Genre:       "Rock",
		Year:        2024,
		TrackNumber: 5,
		DiscNumber:  1,
		Duration:    180,
		Bitrate:     320,
		SampleRate:  44100,
		Channels:    2,
		Format:      "MP3",
		HasCoverArt: true,
		Composer:    "Test Composer",
		Comment:     "Test comment",
		Lyrics:      "Test lyrics",
	}

	assert.Equal(t, "Test Song", metadata.Title)
	assert.Equal(t, "Test Artist", metadata.Artist)
	assert.Equal(t, "Album Artist", metadata.AlbumArtist)
	assert.Equal(t, "Test Album", metadata.Album)
	assert.Equal(t, "Rock", metadata.Genre)
	assert.Equal(t, 2024, metadata.Year)
	assert.Equal(t, 5, metadata.TrackNumber)
	assert.Equal(t, 1, metadata.DiscNumber)
	assert.Equal(t, 180, metadata.Duration)
	assert.Equal(t, 320, metadata.Bitrate)
	assert.Equal(t, 44100, metadata.SampleRate)
	assert.Equal(t, 2, metadata.Channels)
	assert.Equal(t, "MP3", metadata.Format)
	assert.True(t, metadata.HasCoverArt)
	assert.Equal(t, "Test Composer", metadata.Composer)
	assert.Equal(t, "Test comment", metadata.Comment)
	assert.Equal(t, "Test lyrics", metadata.Lyrics)
}
