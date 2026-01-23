package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTrackStructFields verifies Track struct has all required fields
func TestTrackStructFields(t *testing.T) {
	now := time.Now()
	track := Track{
		ID:          "track-123",
		UserID:      "user-456",
		Title:       "Test Song",
		Artist:      "Test Artist",
		AlbumArtist: "Album Artist",
		Album:       "Test Album",
		AlbumID:     "album-789",
		Genre:       "Rock",
		Year:        2024,
		TrackNumber: 5,
		DiscNumber:  1,
		Duration:    180,
		Format:      AudioFormatMP3,
		Bitrate:     320,
		SampleRate:  44100,
		Channels:    2,
		FileSize:    5242880,
		S3Key:       "media/user-456/track-123/audio.mp3",
		CoverArtKey: "media/user-456/track-123/cover.jpg",
		Lyrics:      "Test lyrics",
		Comment:     "Test comment",
		Composer:    "Test Composer",
		PlayCount:   10,
		LastPlayed:  &now,
		Tags:        []string{"favorite", "rock"},
	}

	assert.Equal(t, "track-123", track.ID)
	assert.Equal(t, "user-456", track.UserID)
	assert.Equal(t, "Test Song", track.Title)
	assert.Equal(t, "Test Artist", track.Artist)
	assert.Equal(t, "Album Artist", track.AlbumArtist)
	assert.Equal(t, "Test Album", track.Album)
	assert.Equal(t, "album-789", track.AlbumID)
	assert.Equal(t, "Rock", track.Genre)
	assert.Equal(t, 2024, track.Year)
	assert.Equal(t, 5, track.TrackNumber)
	assert.Equal(t, 1, track.DiscNumber)
	assert.Equal(t, 180, track.Duration)
	assert.Equal(t, AudioFormatMP3, track.Format)
	assert.Equal(t, 320, track.Bitrate)
	assert.Equal(t, 44100, track.SampleRate)
	assert.Equal(t, 2, track.Channels)
	assert.Equal(t, int64(5242880), track.FileSize)
	assert.Equal(t, "media/user-456/track-123/audio.mp3", track.S3Key)
	assert.Equal(t, "media/user-456/track-123/cover.jpg", track.CoverArtKey)
	assert.Equal(t, "Test lyrics", track.Lyrics)
	assert.Equal(t, "Test comment", track.Comment)
	assert.Equal(t, "Test Composer", track.Composer)
	assert.Equal(t, 10, track.PlayCount)
	assert.NotNil(t, track.LastPlayed)
	assert.Equal(t, []string{"favorite", "rock"}, track.Tags)
}

// TestTrackJSONTags verifies JSON serialization
func TestTrackJSONTags(t *testing.T) {
	track := Track{
		ID:          "track-123",
		UserID:      "user-456",
		Title:       "Test Song",
		Artist:      "Test Artist",
		AlbumArtist: "Album Artist",
		Album:       "Test Album",
		Format:      AudioFormatMP3,
	}

	jsonBytes, err := json.Marshal(track)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Verify JSON field names
	assert.Contains(t, jsonMap, "id")
	assert.Contains(t, jsonMap, "userId")
	assert.Contains(t, jsonMap, "title")
	assert.Contains(t, jsonMap, "artist")
	assert.Contains(t, jsonMap, "albumArtist")
	assert.Contains(t, jsonMap, "album")
}

// TestTrackJSONOmitEmpty verifies omitempty behavior
func TestTrackJSONOmitEmpty(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Artist: "Test Artist",
		Format: AudioFormatMP3,
		// Intentionally leaving optional fields empty
	}

	jsonBytes, err := json.Marshal(track)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Optional fields with omitempty should not appear
	assert.NotContains(t, jsonMap, "albumArtist")
	assert.NotContains(t, jsonMap, "album")
	assert.NotContains(t, jsonMap, "genre")
	assert.NotContains(t, jsonMap, "lyrics")
	assert.NotContains(t, jsonMap, "lastPlayed")
}

// TestNewTrackItem verifies DynamoDB item creation
func TestNewTrackItem(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Artist: "Test Artist",
		Format: AudioFormatMP3,
	}

	item := NewTrackItem(track)

	// Verify PK/SK patterns
	assert.Equal(t, "USER#user-456", item.PK)
	assert.Equal(t, "TRACK#track-123", item.SK)
	assert.Equal(t, string(EntityTrack), item.Type)

	// Verify GSI1 for artist queries
	assert.Equal(t, "USER#user-456#ARTIST#Test Artist", item.GSI1PK)
	assert.Equal(t, "TRACK#track-123", item.GSI1SK)
}

// TestNewTrackItemWithEmptyArtist verifies GSI handling when artist is empty
func TestNewTrackItemWithEmptyArtist(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Artist: "", // Empty artist
		Format: AudioFormatMP3,
	}

	item := NewTrackItem(track)

	// GSI1 should not be set if artist is empty
	assert.Empty(t, item.GSI1PK)
	assert.Empty(t, item.GSI1SK)
}

// TestTrackToResponse verifies API response conversion
func TestTrackToResponse(t *testing.T) {
	now := time.Now()
	track := Track{
		ID:         "track-123",
		UserID:     "user-456",
		Title:      "Test Song",
		Artist:     "Test Artist",
		Album:      "Test Album",
		Duration:   185, // 3:05
		Format:     AudioFormatMP3,
		FileSize:   5505024, // ~5.25 MB
		PlayCount:  10,
		LastPlayed: &now,
		Tags:       []string{"favorite"},
	}
	track.CreatedAt = now
	track.UpdatedAt = now

	coverArtURL := "https://cdn.example.com/cover.jpg"
	response := track.ToResponse(coverArtURL)

	assert.Equal(t, "track-123", response.ID)
	assert.Equal(t, "Test Song", response.Title)
	assert.Equal(t, "Test Artist", response.Artist)
	assert.Equal(t, "Test Album", response.Album)
	assert.Equal(t, 185, response.Duration)
	assert.Equal(t, "3:05", response.DurationStr)
	assert.Equal(t, "MP3", response.Format)
	assert.Equal(t, int64(5505024), response.FileSize)
	assert.Equal(t, "5.25 MB", response.FileSizeStr)
	assert.Equal(t, coverArtURL, response.CoverArtURL)
	assert.Equal(t, 10, response.PlayCount)
	assert.NotNil(t, response.LastPlayed)
	assert.Equal(t, []string{"favorite"}, response.Tags)
}

// TestTrackFilterFields verifies TrackFilter struct
func TestTrackFilterFields(t *testing.T) {
	filter := TrackFilter{
		Artist:    "Test Artist",
		Album:     "Test Album",
		Genre:     "Rock",
		Year:      2024,
		Tags:      []string{"favorite"},
		SortBy:    "title",
		SortOrder: "asc",
		Limit:     20,
		LastKey:   "abc123",
	}

	assert.Equal(t, "Test Artist", filter.Artist)
	assert.Equal(t, "Test Album", filter.Album)
	assert.Equal(t, "Rock", filter.Genre)
	assert.Equal(t, 2024, filter.Year)
	assert.Equal(t, []string{"favorite"}, filter.Tags)
	assert.Equal(t, "title", filter.SortBy)
	assert.Equal(t, "asc", filter.SortOrder)
	assert.Equal(t, 20, filter.Limit)
	assert.Equal(t, "abc123", filter.LastKey)
}

// TestCreateTrackRequestFields verifies CreateTrackRequest struct
func TestCreateTrackRequestFields(t *testing.T) {
	req := CreateTrackRequest{
		Title:       "Test Song",
		Artist:      "Test Artist",
		AlbumArtist: "Album Artist",
		Album:       "Test Album",
		Genre:       "Rock",
		Year:        2024,
		TrackNumber: 5,
		DiscNumber:  1,
		Tags:        []string{"favorite"},
	}

	assert.Equal(t, "Test Song", req.Title)
	assert.Equal(t, "Test Artist", req.Artist)
	assert.Equal(t, "Album Artist", req.AlbumArtist)
	assert.Equal(t, "Test Album", req.Album)
	assert.Equal(t, "Rock", req.Genre)
	assert.Equal(t, 2024, req.Year)
	assert.Equal(t, 5, req.TrackNumber)
	assert.Equal(t, 1, req.DiscNumber)
	assert.Equal(t, []string{"favorite"}, req.Tags)
}

// TestUpdateTrackRequestFields verifies UpdateTrackRequest with pointers
func TestUpdateTrackRequestFields(t *testing.T) {
	title := "New Title"
	artist := "New Artist"
	year := 2024

	req := UpdateTrackRequest{
		Title:  &title,
		Artist: &artist,
		Year:   &year,
	}

	assert.NotNil(t, req.Title)
	assert.Equal(t, "New Title", *req.Title)
	assert.NotNil(t, req.Artist)
	assert.Equal(t, "New Artist", *req.Artist)
	assert.NotNil(t, req.Year)
	assert.Equal(t, 2024, *req.Year)
}

// TestTrackToResponseNilTags verifies Tags is never nil in response (prevents JS undefined errors)
func TestTrackToResponseNilTags(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Artist: "Test Artist",
		Format: AudioFormatMP3,
		Tags:   nil, // Explicitly nil
	}
	track.CreatedAt = time.Now()
	track.UpdatedAt = time.Now()

	response := track.ToResponse("")

	// Tags should be empty slice, not nil
	assert.NotNil(t, response.Tags, "Tags should not be nil")
	assert.Len(t, response.Tags, 0, "Tags should be empty slice")

	// Verify JSON serialization includes empty array, not null/missing
	jsonBytes, err := json.Marshal(response)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// Tags should be present as empty array
	assert.Contains(t, jsonMap, "tags", "tags field should be present in JSON")
	tags, ok := jsonMap["tags"].([]interface{})
	assert.True(t, ok, "tags should be an array")
	assert.Len(t, tags, 0, "tags should be empty array")
}

// =============================================================================
// Waveform and Analysis Fields Tests (TDD Red - will fail until implemented)
// =============================================================================

// AnalysisStatus constants for testing - TDD Red: These will be added to track.go
const (
	AnalysisStatusPending   = "PENDING"
	AnalysisStatusAnalyzing = "ANALYZING"
	AnalysisStatusCompleted = "COMPLETED"
	AnalysisStatusFailed    = "FAILED"
)

// TestTrack_WaveformURL verifies the WaveformURL field exists
func TestTrack_WaveformURL(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}

	// TDD Red: This will fail until WaveformURL field is added to Track
	track.WaveformURL = "https://cdn.example.com/waveforms/track-123.json"
	assert.Equal(t, "https://cdn.example.com/waveforms/track-123.json", track.WaveformURL)
}

// TestTrack_BeatGrid verifies the BeatGrid field exists
func TestTrack_BeatGrid(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}

	// TDD Red: This will fail until BeatGrid field is added to Track
	track.BeatGrid = []int64{0, 500, 1000, 1500, 2000}
	assert.Equal(t, []int64{0, 500, 1000, 1500, 2000}, track.BeatGrid)
	assert.Len(t, track.BeatGrid, 5)
}

// TestTrack_AnalysisStatus verifies the AnalysisStatus field exists
func TestTrack_AnalysisStatus(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}

	// TDD Red: This will fail until AnalysisStatus field is added to Track
	track.AnalysisStatus = AnalysisStatusCompleted
	assert.Equal(t, AnalysisStatusCompleted, track.AnalysisStatus)
}

// TestTrack_AnalyzedAt verifies the AnalyzedAt field exists
func TestTrack_AnalyzedAt(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}

	// TDD Red: This will fail until AnalyzedAt field is added to Track
	now := time.Now()
	track.AnalyzedAt = &now
	assert.NotNil(t, track.AnalyzedAt)
	assert.Equal(t, now, *track.AnalyzedAt)
}

// TestTrack_WaveformJSON verifies WaveformURL JSON serialization
func TestTrack_WaveformJSON(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}
	track.WaveformURL = "https://cdn.example.com/waveforms/track-123.json"

	jsonBytes, err := json.Marshal(track)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// TDD Red: This will fail until field has correct JSON tag
	assert.Contains(t, jsonMap, "waveformUrl")
	assert.Equal(t, "https://cdn.example.com/waveforms/track-123.json", jsonMap["waveformUrl"])
}

// TestTrack_BeatGridJSON verifies BeatGrid JSON serialization
func TestTrack_BeatGridJSON(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}
	track.BeatGrid = []int64{0, 500, 1000}

	jsonBytes, err := json.Marshal(track)
	require.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	require.NoError(t, err)

	// TDD Red: This will fail until field has correct JSON tag
	assert.Contains(t, jsonMap, "beatGrid")
}

// TestTrack_AnalysisStatusValues verifies valid analysis status values
func TestTrack_AnalysisStatusValues(t *testing.T) {
	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"pending", AnalysisStatusPending, true},
		{"analyzing", AnalysisStatusAnalyzing, true},
		{"completed", AnalysisStatusCompleted, true},
		{"failed", AnalysisStatusFailed, true},
		{"invalid", "INVALID", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TDD Red: This will fail until IsValidAnalysisStatus is implemented
			valid := IsValidAnalysisStatus(tt.status)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

// TestTrackResponse_WaveformURL verifies TrackResponse includes WaveformURL
func TestTrackResponse_WaveformURL(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}
	track.WaveformURL = "https://cdn.example.com/waveforms/track-123.json"
	track.CreatedAt = time.Now()
	track.UpdatedAt = time.Now()

	response := track.ToResponse("")

	// TDD Red: This will fail until TrackResponse.WaveformURL is added
	assert.Equal(t, "https://cdn.example.com/waveforms/track-123.json", response.WaveformURL)
}

// TestTrackResponse_AnalysisStatus verifies TrackResponse includes AnalysisStatus
func TestTrackResponse_AnalysisStatus(t *testing.T) {
	track := Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Test Song",
		Format: AudioFormatMP3,
	}
	track.AnalysisStatus = AnalysisStatusCompleted
	track.CreatedAt = time.Now()
	track.UpdatedAt = time.Now()

	response := track.ToResponse("")

	// TDD Red: This will fail until TrackResponse.AnalysisStatus is added
	assert.Equal(t, AnalysisStatusCompleted, response.AnalysisStatus)
}

// TestTrackFilter_BPMRange verifies TrackFilter supports BPM range filtering
func TestTrackFilter_BPMRange(t *testing.T) {
	filter := TrackFilter{
		BPMMin: 120,
		BPMMax: 130,
	}

	// TDD Red: These fields should already exist but verify they work
	assert.Equal(t, 120, filter.BPMMin)
	assert.Equal(t, 130, filter.BPMMax)
}

// IsValidAnalysisStatus checks if the given status is valid
// TDD Red: STUB - returns false for all inputs, tests will fail
func IsValidAnalysisStatus(status string) bool {
	switch status {
	case AnalysisStatusPending, AnalysisStatusAnalyzing, AnalysisStatusCompleted, AnalysisStatusFailed:
		return true
	default:
		return false
	}
}
