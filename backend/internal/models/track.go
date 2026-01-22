package models

import (
	"fmt"
	"time"
)

// HLSStatus represents the transcoding status for HLS streaming
type HLSStatus string

const (
	HLSStatusPending    HLSStatus = "PENDING"
	HLSStatusProcessing HLSStatus = "PROCESSING"
	HLSStatusReady      HLSStatus = "READY"
	HLSStatusFailed     HLSStatus = "FAILED"
)

// Track represents a music track in the library
type Track struct {
	ID          string      `json:"id" dynamodbav:"id"`
	UserID      string      `json:"userId" dynamodbav:"userId"`
	Title       string      `json:"title" dynamodbav:"title"`
	Artist      string      `json:"artist" dynamodbav:"artist"`
	AlbumArtist string      `json:"albumArtist,omitempty" dynamodbav:"albumArtist,omitempty"`
	Album       string      `json:"album,omitempty" dynamodbav:"album,omitempty"`
	AlbumID     string      `json:"albumId,omitempty" dynamodbav:"albumId,omitempty"`
	Genre       string      `json:"genre,omitempty" dynamodbav:"genre,omitempty"`
	Year        int         `json:"year,omitempty" dynamodbav:"year,omitempty"`
	TrackNumber int         `json:"trackNumber,omitempty" dynamodbav:"trackNumber,omitempty"`
	DiscNumber  int         `json:"discNumber,omitempty" dynamodbav:"discNumber,omitempty"`
	Duration    int         `json:"duration" dynamodbav:"duration"` // Duration in seconds
	Format      AudioFormat `json:"format" dynamodbav:"format"`
	Bitrate     int         `json:"bitrate,omitempty" dynamodbav:"bitrate,omitempty"`   // kbps
	SampleRate  int         `json:"sampleRate,omitempty" dynamodbav:"sampleRate,omitempty"` // Hz
	Channels    int         `json:"channels,omitempty" dynamodbav:"channels,omitempty"`
	FileSize    int64       `json:"fileSize" dynamodbav:"fileSize"` // bytes
	S3Key       string      `json:"s3Key" dynamodbav:"s3Key"`
	CoverArtKey string      `json:"coverArtKey,omitempty" dynamodbav:"coverArtKey,omitempty"`
	Lyrics      string      `json:"lyrics,omitempty" dynamodbav:"lyrics,omitempty"`
	Comment     string      `json:"comment,omitempty" dynamodbav:"comment,omitempty"`
	Composer    string      `json:"composer,omitempty" dynamodbav:"composer,omitempty"`
	PlayCount   int         `json:"playCount" dynamodbav:"playCount"`
	LastPlayed  *time.Time  `json:"lastPlayed,omitempty" dynamodbav:"lastPlayed,omitempty"`
	Tags        []string    `json:"tags,omitempty" dynamodbav:"tags,omitempty"`

	// HLS streaming fields
	HLSStatus        HLSStatus `json:"hlsStatus,omitempty" dynamodbav:"hlsStatus,omitempty"`
	HLSPlaylistKey   string    `json:"hlsPlaylistKey,omitempty" dynamodbav:"hlsPlaylistKey,omitempty"` // S3 key to master.m3u8
	HLSJobID         string    `json:"hlsJobId,omitempty" dynamodbav:"hlsJobId,omitempty"`             // MediaConvert job ID
	HLSTranscodedAt  *time.Time `json:"hlsTranscodedAt,omitempty" dynamodbav:"hlsTranscodedAt,omitempty"`

	Timestamps
}

// TrackItem represents a Track in DynamoDB single-table design
type TrackItem struct {
	DynamoDBItem
	Track
}

// NewTrackItem creates a DynamoDB item for a track
func NewTrackItem(track Track) TrackItem {
	item := TrackItem{
		DynamoDBItem: DynamoDBItem{
			PK:   fmt.Sprintf("USER#%s", track.UserID),
			SK:   fmt.Sprintf("TRACK#%s", track.ID),
			Type: string(EntityTrack),
		},
		Track: track,
	}

	// Set GSI for artist-based queries
	if track.Artist != "" {
		item.GSI1PK = fmt.Sprintf("USER#%s#ARTIST#%s", track.UserID, track.Artist)
		item.GSI1SK = fmt.Sprintf("TRACK#%s", track.ID)
	}

	return item
}

// CreateTrackRequest represents a request to create a track (typically from upload)
type CreateTrackRequest struct {
	Title       string   `json:"title" validate:"required,min=1,max=500"`
	Artist      string   `json:"artist" validate:"required,min=1,max=500"`
	AlbumArtist string   `json:"albumArtist,omitempty" validate:"omitempty,max=500"`
	Album       string   `json:"album,omitempty" validate:"omitempty,max=500"`
	Genre       string   `json:"genre,omitempty" validate:"omitempty,max=100"`
	Year        int      `json:"year,omitempty" validate:"omitempty,min=1,max=9999"`
	TrackNumber int      `json:"trackNumber,omitempty" validate:"omitempty,min=0"`
	DiscNumber  int      `json:"discNumber,omitempty" validate:"omitempty,min=0"`
	Tags        []string `json:"tags,omitempty" validate:"omitempty,dive,min=1,max=50"`
}

// UpdateTrackRequest represents a request to update track metadata
type UpdateTrackRequest struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=1,max=500"`
	Artist      *string  `json:"artist,omitempty" validate:"omitempty,min=1,max=500"`
	AlbumArtist *string  `json:"albumArtist,omitempty" validate:"omitempty,max=500"`
	Album       *string  `json:"album,omitempty" validate:"omitempty,max=500"`
	Genre       *string  `json:"genre,omitempty" validate:"omitempty,max=100"`
	Year        *int     `json:"year,omitempty" validate:"omitempty,min=1,max=9999"`
	TrackNumber *int     `json:"trackNumber,omitempty" validate:"omitempty,min=0"`
	DiscNumber  *int     `json:"discNumber,omitempty" validate:"omitempty,min=0"`
	Lyrics      *string  `json:"lyrics,omitempty"`
	Comment     *string  `json:"comment,omitempty" validate:"omitempty,max=1000"`
	Tags        []string `json:"tags,omitempty" validate:"omitempty,dive,min=1,max=50"`
}

// TrackResponse represents a track in API responses
type TrackResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Artist       string    `json:"artist"`
	AlbumArtist  string    `json:"albumArtist,omitempty"`
	Album        string    `json:"album,omitempty"`
	AlbumID      string    `json:"albumId,omitempty"`
	Genre        string    `json:"genre,omitempty"`
	Year         int       `json:"year,omitempty"`
	TrackNumber  int       `json:"trackNumber,omitempty"`
	DiscNumber   int       `json:"discNumber,omitempty"`
	Duration     int       `json:"duration"`
	DurationStr  string    `json:"durationStr"`
	Format       string    `json:"format"`
	FileSize     int64     `json:"fileSize"`
	FileSizeStr  string    `json:"fileSizeStr"`
	CoverArtURL  string    `json:"coverArtUrl,omitempty"`
	PlayCount    int       `json:"playCount"`
	LastPlayed   *time.Time `json:"lastPlayed,omitempty"`
	Tags         []string  `json:"tags"`
	HLSStatus    string    `json:"hlsStatus,omitempty"`
	HLSReady     bool      `json:"hlsReady"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ToResponse converts a Track to a TrackResponse
func (t *Track) ToResponse(coverArtURL string) TrackResponse {
	// Ensure tags is never nil to avoid undefined in JavaScript
	tags := t.Tags
	if tags == nil {
		tags = []string{}
	}

	return TrackResponse{
		ID:           t.ID,
		Title:        t.Title,
		Artist:       t.Artist,
		AlbumArtist:  t.AlbumArtist,
		Album:        t.Album,
		AlbumID:      t.AlbumID,
		Genre:        t.Genre,
		Year:         t.Year,
		TrackNumber:  t.TrackNumber,
		DiscNumber:   t.DiscNumber,
		Duration:     t.Duration,
		DurationStr:  formatDuration(t.Duration),
		Format:       string(t.Format),
		FileSize:     t.FileSize,
		FileSizeStr:  formatFileSize(t.FileSize),
		CoverArtURL:  coverArtURL,
		PlayCount:    t.PlayCount,
		LastPlayed:   t.LastPlayed,
		Tags:         tags,
		HLSStatus:    string(t.HLSStatus),
		HLSReady:     t.HLSStatus == HLSStatusReady,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

// TrackFilter represents filter options for listing tracks
type TrackFilter struct {
	Artist      string   `query:"artist"`
	Album       string   `query:"album"`
	Genre       string   `query:"genre"`
	Year        int      `query:"year"`
	Tags        []string `query:"tags"`
	SortBy      string   `query:"sortBy"`    // title, artist, album, createdAt, playCount
	SortOrder   string   `query:"sortOrder"` // asc, desc
	Limit       int      `query:"limit"`
	LastKey     string   `query:"lastKey"`
	GlobalScope bool     `query:"-"` // If true, return tracks from all users (requires GLOBAL permission)
}

// Helper functions

func formatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
