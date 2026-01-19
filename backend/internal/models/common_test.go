package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEntityTypeConstants verifies all entity type constants exist
func TestEntityTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant EntityType
		expected string
	}{
		{"EntityUser", EntityUser, "USER"},
		{"EntityTrack", EntityTrack, "TRACK"},
		{"EntityAlbum", EntityAlbum, "ALBUM"},
		{"EntityPlaylist", EntityPlaylist, "PLAYLIST"},
		{"EntityPlaylistTrack", EntityPlaylistTrack, "PLAYLIST_TRACK"},
		{"EntityUpload", EntityUpload, "UPLOAD"},
		{"EntityTag", EntityTag, "TAG"},
		{"EntityTrackTag", EntityTrackTag, "TRACK_TAG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.constant))
		})
	}
}

// TestUploadStatusConstants verifies all upload status constants
func TestUploadStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant UploadStatus
		expected string
	}{
		{"UploadStatusPending", UploadStatusPending, "PENDING"},
		{"UploadStatusProcessing", UploadStatusProcessing, "PROCESSING"},
		{"UploadStatusCompleted", UploadStatusCompleted, "COMPLETED"},
		{"UploadStatusFailed", UploadStatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.constant))
		})
	}
}

// TestAudioFormatConstants verifies all audio format constants
func TestAudioFormatConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant AudioFormat
		expected string
	}{
		{"AudioFormatMP3", AudioFormatMP3, "MP3"},
		{"AudioFormatFLAC", AudioFormatFLAC, "FLAC"},
		{"AudioFormatWAV", AudioFormatWAV, "WAV"},
		{"AudioFormatAAC", AudioFormatAAC, "AAC"},
		{"AudioFormatOGG", AudioFormatOGG, "OGG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.constant))
		})
	}
}

// TestFormatDuration verifies duration formatting
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"zero seconds", 0, "0:00"},
		{"30 seconds", 30, "0:30"},
		{"1 minute", 60, "1:00"},
		{"1 minute 30 seconds", 90, "1:30"},
		{"3 minutes 45 seconds", 225, "3:45"},
		{"1 hour", 3600, "1:00:00"},
		{"1 hour 30 minutes", 5400, "1:30:00"},
		{"1 hour 5 minutes 30 seconds", 3930, "1:05:30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.seconds)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatFileSize verifies file size formatting
func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"zero bytes", 0, "0 B"},
		{"500 bytes", 500, "500 B"},
		{"1 KB", 1024, "1.00 KB"},
		{"1.5 KB", 1536, "1.50 KB"},
		{"1 MB", 1048576, "1.00 MB"},
		{"5.25 MB", 5505024, "5.25 MB"},
		{"1 GB", 1073741824, "1.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDynamoDBItemFields verifies DynamoDBItem struct fields
func TestDynamoDBItemFields(t *testing.T) {
	item := DynamoDBItem{
		PK:     "USER#123",
		SK:     "TRACK#456",
		GSI1PK: "USER#123#ARTIST#Test",
		GSI1SK: "TRACK#456",
		Type:   "TRACK",
	}

	assert.Equal(t, "USER#123", item.PK)
	assert.Equal(t, "TRACK#456", item.SK)
	assert.Equal(t, "USER#123#ARTIST#Test", item.GSI1PK)
	assert.Equal(t, "TRACK#456", item.GSI1SK)
	assert.Equal(t, "TRACK", item.Type)
}

// TestTimestampsFields verifies Timestamps struct fields
func TestTimestampsFields(t *testing.T) {
	ts := Timestamps{}
	// Should have CreatedAt and UpdatedAt fields
	assert.NotNil(t, ts.CreatedAt)
	assert.NotNil(t, ts.UpdatedAt)
}

// TestPaginationFields verifies Pagination struct fields
func TestPaginationFields(t *testing.T) {
	p := Pagination{
		Limit:         20,
		LastKey:       "abc123",
		NextKey:       "def456",
		TotalEstimate: 100,
	}

	assert.Equal(t, 20, p.Limit)
	assert.Equal(t, "abc123", p.LastKey)
	assert.Equal(t, "def456", p.NextKey)
	assert.Equal(t, 100, p.TotalEstimate)
}

// TestPaginatedResponse verifies generic PaginatedResponse
func TestPaginatedResponse(t *testing.T) {
	items := []string{"a", "b", "c"}
	pagination := Pagination{Limit: 10}

	response := PaginatedResponse[string]{
		Items:      items,
		Pagination: pagination,
	}

	assert.Equal(t, 3, len(response.Items))
	assert.Equal(t, 10, response.Pagination.Limit)
}
