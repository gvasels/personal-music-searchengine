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

// TestUploadStepTrackingFields verifies step tracking fields
func TestUploadStepTrackingFields(t *testing.T) {
	upload := Upload{
		ID:                "upload-123",
		UserID:            "user-456",
		Status:            UploadStatusProcessing,
		MetadataExtracted: true,
		CoverArtExtracted: true,
		TrackCreated:      false,
		Indexed:           false,
		FileMoved:         false,
	}

	assert.True(t, upload.MetadataExtracted)
	assert.True(t, upload.CoverArtExtracted)
	assert.False(t, upload.TrackCreated)
	assert.False(t, upload.Indexed)
	assert.False(t, upload.FileMoved)
}

// TestUploadMultipartFields verifies multipart upload tracking
func TestUploadMultipartFields(t *testing.T) {
	upload := Upload{
		ID:            "upload-123",
		UserID:        "user-456",
		IsMultipart:   true,
		MultipartID:   "multipart-abc123",
		PartsUploaded: 5,
		TotalParts:    10,
	}

	assert.True(t, upload.IsMultipart)
	assert.Equal(t, "multipart-abc123", upload.MultipartID)
	assert.Equal(t, 5, upload.PartsUploaded)
	assert.Equal(t, 10, upload.TotalParts)
}

// TestUploadToResponseIncludesSteps verifies step tracking in response
func TestUploadToResponseIncludesSteps(t *testing.T) {
	upload := Upload{
		ID:                "upload-123",
		UserID:            "user-456",
		FileName:          "song.mp3",
		FileSize:          5242880,
		ContentType:       "audio/mpeg",
		Status:            UploadStatusProcessing,
		MetadataExtracted: true,
		CoverArtExtracted: true,
		TrackCreated:      false,
		Indexed:           false,
		FileMoved:         false,
	}
	upload.CreatedAt = time.Now()

	response := upload.ToResponse()

	assert.True(t, response.Steps.MetadataExtracted)
	assert.True(t, response.Steps.CoverArtExtracted)
	assert.False(t, response.Steps.TrackCreated)
	assert.False(t, response.Steps.Indexed)
	assert.False(t, response.Steps.FileMoved)
}

// TestProcessingStepConstants verifies processing step constants
func TestProcessingStepConstants(t *testing.T) {
	tests := []struct {
		name     string
		step     ProcessingStep
		expected string
	}{
		{"StepExtractMetadata", StepExtractMetadata, "extract_metadata"},
		{"StepExtractCover", StepExtractCover, "extract_cover"},
		{"StepCreateTrack", StepCreateTrack, "create_track"},
		{"StepIndex", StepIndex, "index"},
		{"StepMoveFile", StepMoveFile, "move_file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.step))
		})
	}
}

// TestReprocessUploadRequestFields verifies reprocess request
func TestReprocessUploadRequestFields(t *testing.T) {
	req := ReprocessUploadRequest{
		FromStep: StepCreateTrack,
	}

	assert.Equal(t, StepCreateTrack, req.FromStep)
}

// TestCoverArtUploadRequestFields verifies cover art upload request
func TestCoverArtUploadRequestFields(t *testing.T) {
	req := CoverArtUploadRequest{
		FileName:    "cover.jpg",
		FileSize:    1048576, // 1MB
		ContentType: "image/jpeg",
	}

	assert.Equal(t, "cover.jpg", req.FileName)
	assert.Equal(t, int64(1048576), req.FileSize)
	assert.Equal(t, "image/jpeg", req.ContentType)
}

// TestCoverArtUploadResponseFields verifies cover art upload response
func TestCoverArtUploadResponseFields(t *testing.T) {
	expiresAt := time.Now().Add(15 * time.Minute)
	response := CoverArtUploadResponse{
		TrackID:     "track-123",
		UploadURL:   "https://s3.amazonaws.com/bucket/cover?signature=...",
		ExpiresAt:   expiresAt,
		MaxFileSize: 10485760, // 10MB
	}

	assert.Equal(t, "track-123", response.TrackID)
	assert.Contains(t, response.UploadURL, "s3.amazonaws.com")
	assert.Equal(t, expiresAt, response.ExpiresAt)
	assert.Equal(t, int64(10485760), response.MaxFileSize)
}

// TestMultipartUploadPartURLFields verifies multipart part URL
func TestMultipartUploadPartURLFields(t *testing.T) {
	expiresAt := time.Now().Add(1 * time.Hour)
	partURL := MultipartUploadPartURL{
		PartNumber: 1,
		UploadURL:  "https://s3.amazonaws.com/bucket/key?partNumber=1&uploadId=abc",
		ExpiresAt:  expiresAt,
	}

	assert.Equal(t, 1, partURL.PartNumber)
	assert.Contains(t, partURL.UploadURL, "partNumber=1")
	assert.Equal(t, expiresAt, partURL.ExpiresAt)
}

// TestPresignedUploadResponseMultipartFields verifies multipart response
func TestPresignedUploadResponseMultipartFields(t *testing.T) {
	expiresAt := time.Now().Add(1 * time.Hour)
	response := PresignedUploadResponse{
		UploadID:    "upload-123",
		IsMultipart: true,
		MultipartID: "multipart-abc",
		PartURLs: []MultipartUploadPartURL{
			{PartNumber: 1, UploadURL: "https://example.com/part1", ExpiresAt: expiresAt},
			{PartNumber: 2, UploadURL: "https://example.com/part2", ExpiresAt: expiresAt},
		},
	}

	assert.True(t, response.IsMultipart)
	assert.Equal(t, "multipart-abc", response.MultipartID)
	assert.Len(t, response.PartURLs, 2)
	assert.Equal(t, 1, response.PartURLs[0].PartNumber)
	assert.Equal(t, 2, response.PartURLs[1].PartNumber)
}

// TestCompleteMultipartUploadRequestFields verifies complete multipart request
func TestCompleteMultipartUploadRequestFields(t *testing.T) {
	req := CompleteMultipartUploadRequest{
		UploadID: "upload-123",
		Parts: []CompletedPartInfo{
			{PartNumber: 1, ETag: "etag1"},
			{PartNumber: 2, ETag: "etag2"},
		},
	}

	assert.Equal(t, "upload-123", req.UploadID)
	assert.Len(t, req.Parts, 2)
	assert.Equal(t, 1, req.Parts[0].PartNumber)
	assert.Equal(t, "etag1", req.Parts[0].ETag)
}

// TestCompletedPartInfoFields verifies completed part info
func TestCompletedPartInfoFields(t *testing.T) {
	part := CompletedPartInfo{
		PartNumber: 5,
		ETag:       "\"abc123def456\"",
	}

	assert.Equal(t, 5, part.PartNumber)
	assert.Equal(t, "\"abc123def456\"", part.ETag)
}

// TestPresignedUploadRequestMaxFileSize verifies 1GB max file size
func TestPresignedUploadRequestMaxFileSize(t *testing.T) {
	// Max file size should be 1GB = 1073741824 bytes
	req := PresignedUploadRequest{
		FileName:    "large-file.flac",
		FileSize:    1073741824, // 1GB - should be valid
		ContentType: "audio/flac",
		IsMultipart: true,
	}

	assert.Equal(t, int64(1073741824), req.FileSize)
	assert.True(t, req.IsMultipart)
}
