package models

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// EntityType represents the type of entity in the single-table design
type EntityType string

const (
	EntityUser          EntityType = "USER"
	EntityTrack         EntityType = "TRACK"
	EntityAlbum         EntityType = "ALBUM"
	EntityArtist        EntityType = "ARTIST"
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
	GSI2PK string `dynamodbav:"GSI2PK,omitempty"` // Used for public playlist discovery
	GSI2SK string `dynamodbav:"GSI2SK,omitempty"` // Used for public playlist discovery
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

// PaginationCursor represents the internal structure of a pagination cursor
// This is encoded to base64 and passed to clients as an opaque string
type PaginationCursor struct {
	PK     string `json:"pk"`
	SK     string `json:"sk"`
	GSI1PK string `json:"gsi1pk,omitempty"`
	GSI1SK string `json:"gsi1sk,omitempty"`
}

// EncodeCursor encodes a PaginationCursor to an opaque base64 string
func EncodeCursor(cursor PaginationCursor) string {
	if cursor.PK == "" && cursor.SK == "" {
		return ""
	}
	data, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeCursor decodes an opaque base64 cursor string to a PaginationCursor
func DecodeCursor(encoded string) (PaginationCursor, error) {
	var cursor PaginationCursor
	if encoded == "" {
		return cursor, nil
	}
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return cursor, err
	}
	err = json.Unmarshal(data, &cursor)
	return cursor, err
}

// NewPaginationCursor creates a new cursor from DynamoDB key components
func NewPaginationCursor(pk, sk string) PaginationCursor {
	return PaginationCursor{
		PK: pk,
		SK: sk,
	}
}

// NewPaginationCursorWithGSI creates a new cursor with GSI key components
func NewPaginationCursorWithGSI(pk, sk, gsi1pk, gsi1sk string) PaginationCursor {
	return PaginationCursor{
		PK:     pk,
		SK:     sk,
		GSI1PK: gsi1pk,
		GSI1SK: gsi1sk,
	}
}
