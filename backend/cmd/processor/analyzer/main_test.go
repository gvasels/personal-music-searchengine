package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvent_Structure(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected Event
	}{
		{
			name: "complete event",
			json: `{
				"uploadId": "upload-123",
				"userId": "user-456",
				"s3Key": "uploads/user-456/track.mp3",
				"fileName": "track.mp3",
				"bucketName": "music-bucket"
			}`,
			expected: Event{
				UploadID:   "upload-123",
				UserID:     "user-456",
				S3Key:      "uploads/user-456/track.mp3",
				FileName:   "track.mp3",
				BucketName: "music-bucket",
			},
		},
		{
			name: "minimal event",
			json: `{
				"uploadId": "u1",
				"userId": "u2",
				"s3Key": "key",
				"fileName": "f.mp3",
				"bucketName": "bucket"
			}`,
			expected: Event{
				UploadID:   "u1",
				UserID:     "u2",
				S3Key:      "key",
				FileName:   "f.mp3",
				BucketName: "bucket",
			},
		},
		{
			name: "event with special characters in filename",
			json: `{
				"uploadId": "upload-789",
				"userId": "user-abc",
				"s3Key": "uploads/user-abc/My Song (feat. Artist) [Remix].mp3",
				"fileName": "My Song (feat. Artist) [Remix].mp3",
				"bucketName": "music-uploads"
			}`,
			expected: Event{
				UploadID:   "upload-789",
				UserID:     "user-abc",
				S3Key:      "uploads/user-abc/My Song (feat. Artist) [Remix].mp3",
				FileName:   "My Song (feat. Artist) [Remix].mp3",
				BucketName: "music-uploads",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event Event
			err := json.Unmarshal([]byte(tt.json), &event)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, event)
		})
	}
}

func TestEvent_EmptyFields(t *testing.T) {
	jsonStr := `{}`
	var event Event
	err := json.Unmarshal([]byte(jsonStr), &event)
	require.NoError(t, err)
	assert.Empty(t, event.UploadID)
	assert.Empty(t, event.UserID)
	assert.Empty(t, event.S3Key)
	assert.Empty(t, event.FileName)
	assert.Empty(t, event.BucketName)
}

func TestResponse_Structure(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		contains []string
	}{
		{
			name: "successful analysis response",
			response: Response{
				BPM:        120,
				MusicalKey: "A",
				KeyMode:    "minor",
				KeyCamelot: "8A",
				Analyzed:   true,
				Error:      "",
			},
			contains: []string{
				`"bpm":120`,
				`"analyzed":true`,
			},
		},
		{
			name: "failed analysis response",
			response: Response{
				BPM:        0,
				MusicalKey: "",
				KeyMode:    "",
				KeyCamelot: "",
				Analyzed:   false,
				Error:      "analysis failed: FFmpeg error",
			},
			contains: []string{
				`"analyzed":false`,
				`"error":"analysis failed: FFmpeg error"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			require.NoError(t, err)
			jsonStr := string(data)
			for _, expected := range tt.contains {
				assert.Contains(t, jsonStr, expected)
			}
		})
	}
}

func TestResponse_AnalyzedFalse(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
	}{
		{
			name:         "file size validation error",
			errorMessage: "file too large: exceeds 500MB limit",
		},
		{
			name:         "download error",
			errorMessage: "failed to download from S3: access denied",
		},
		{
			name:         "analysis error",
			errorMessage: "FFmpeg analysis failed: unsupported codec",
		},
		{
			name:         "timeout error",
			errorMessage: "analysis timed out after 25 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := Response{
				BPM:        0,
				MusicalKey: "",
				KeyMode:    "",
				KeyCamelot: "",
				Analyzed:   false,
				Error:      tt.errorMessage,
			}

			assert.False(t, response.Analyzed)
			assert.Equal(t, tt.errorMessage, response.Error)
			assert.Equal(t, 0, response.BPM)
			assert.Empty(t, response.MusicalKey)

			// Verify JSON marshaling
			data, err := json.Marshal(response)
			require.NoError(t, err)
			assert.Contains(t, string(data), `"analyzed":false`)
		})
	}
}

func TestResponse_AnalyzedTrue(t *testing.T) {
	tests := []struct {
		name       string
		bpm        int
		musicalKey string
		keyMode    string
		keyCamelot string
	}{
		{
			name:       "major key analysis",
			bpm:        128,
			musicalKey: "C",
			keyMode:    "major",
			keyCamelot: "8B",
		},
		{
			name:       "minor key analysis",
			bpm:        120,
			musicalKey: "A",
			keyMode:    "minor",
			keyCamelot: "8A",
		},
		{
			name:       "high BPM track",
			bpm:        175,
			musicalKey: "F#",
			keyMode:    "minor",
			keyCamelot: "2A",
		},
		{
			name:       "low BPM track",
			bpm:        70,
			musicalKey: "Bb",
			keyMode:    "major",
			keyCamelot: "6B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := Response{
				BPM:        tt.bpm,
				MusicalKey: tt.musicalKey,
				KeyMode:    tt.keyMode,
				KeyCamelot: tt.keyCamelot,
				Analyzed:   true,
				Error:      "",
			}

			assert.True(t, response.Analyzed)
			assert.Empty(t, response.Error)
			assert.Equal(t, tt.bpm, response.BPM)
			assert.Equal(t, tt.musicalKey, response.MusicalKey)
			assert.Equal(t, tt.keyMode, response.KeyMode)
			assert.Equal(t, tt.keyCamelot, response.KeyCamelot)

			// Verify JSON marshaling
			data, err := json.Marshal(response)
			require.NoError(t, err)
			assert.Contains(t, string(data), `"analyzed":true`)
		})
	}
}

func TestResponse_JSONRoundTrip(t *testing.T) {
	original := Response{
		BPM:        120,
		MusicalKey: "A",
		KeyMode:    "minor",
		KeyCamelot: "8A",
		Analyzed:   true,
		Error:      "",
	}

	// Marshal
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal
	var decoded Response
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEvent_JSONRoundTrip(t *testing.T) {
	original := Event{
		UploadID:   "upload-123",
		UserID:     "user-456",
		S3Key:      "uploads/user-456/track.mp3",
		FileName:   "track.mp3",
		BucketName: "music-bucket",
	}

	// Marshal
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal
	var decoded Event
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEvent_JSONTags(t *testing.T) {
	event := Event{
		UploadID:   "test-upload",
		UserID:     "test-user",
		S3Key:      "test-key",
		FileName:   "test-file",
		BucketName: "test-bucket",
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"uploadId"`)
	assert.Contains(t, jsonStr, `"userId"`)
	assert.Contains(t, jsonStr, `"s3Key"`)
	assert.Contains(t, jsonStr, `"fileName"`)
	assert.Contains(t, jsonStr, `"bucketName"`)
}

func TestResponse_JSONTags(t *testing.T) {
	response := Response{
		BPM:        100,
		MusicalKey: "C",
		KeyMode:    "major",
		KeyCamelot: "8B",
		Analyzed:   true,
		Error:      "",
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"bpm"`)
	assert.Contains(t, jsonStr, `"analyzed"`)
}
