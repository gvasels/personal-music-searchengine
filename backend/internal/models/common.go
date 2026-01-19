package models

import "time"

// EntityType represents the type of entity in the single-table design
type EntityType string

const (
	EntityUser          EntityType = "USER"
	EntityTrack         EntityType = "TRACK"
	EntityAlbum         EntityType = "ALBUM"
	EntityPlaylist      EntityType = "PLAYLIST"
	EntityPlaylistTrack EntityType = "PLAYLIST_TRACK"
	EntityUpload        EntityType = "UPLOAD"
	EntityTag           EntityType = "TAG"
	EntityTrackTag      EntityType = "TRACK_TAG"
)

// UploadStatus represents the status of a file upload
type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "PENDING"
	UploadStatusProcessing UploadStatus = "PROCESSING"
	UploadStatusCompleted  UploadStatus = "COMPLETED"
	UploadStatusFailed     UploadStatus = "FAILED"
)

// AudioFormat represents supported audio formats
type AudioFormat string

const (
	AudioFormatMP3  AudioFormat = "MP3"
	AudioFormatFLAC AudioFormat = "FLAC"
	AudioFormatWAV  AudioFormat = "WAV"
	AudioFormatAAC  AudioFormat = "AAC"
	AudioFormatOGG  AudioFormat = "OGG"
)

// Timestamps provides common timestamp fields
type Timestamps struct {
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

// DynamoDBItem represents a base item for single-table design
type DynamoDBItem struct {
	PK     string `dynamodbav:"PK"`
	SK     string `dynamodbav:"SK"`
	GSI1PK string `dynamodbav:"GSI1PK,omitempty"`
	GSI1SK string `dynamodbav:"GSI1SK,omitempty"`
	Type   string `dynamodbav:"Type"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Limit         int    `json:"limit"`
	LastKey       string `json:"lastKey,omitempty"`
	NextKey       string `json:"nextKey,omitempty"`
	TotalEstimate int    `json:"totalEstimate,omitempty"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}
