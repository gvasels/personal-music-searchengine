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
	ArtistID    string      `json:"artistId,omitempty" dynamodbav:"artistId,omitempty"`                 // Reference to Artist entity
	Artists     []ArtistContribution `json:"artists,omitempty" dynamodbav:"artists,omitempty"`          // Multi-artist support
	ArtistLegacy string     `json:"-" dynamodbav:"artistLegacy,omitempty"`                              // Backup during migration
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

	// Audio analysis fields
	BPM         int    `json:"bpm,omitempty" dynamodbav:"bpm,omitempty"`                 // Beats per minute (20-300)
	MusicalKey  string `json:"musicalKey,omitempty" dynamodbav:"musicalKey,omitempty"`   // e.g., "Am", "C", "F#m"
	KeyMode     string `json:"keyMode,omitempty" dynamodbav:"keyMode,omitempty"`         // "major" or "minor"
	KeyCamelot  string `json:"keyCamelot,omitempty" dynamodbav:"keyCamelot,omitempty"`   // e.g., "8A", "11B"

	// HLS streaming fields
	HLSStatus        HLSStatus `json:"hlsStatus,omitempty" dynamodbav:"hlsStatus,omitempty"`
	HLSPlaylistKey   string    `json:"hlsPlaylistKey,omitempty" dynamodbav:"hlsPlaylistKey,omitempty"` // S3 key to master.m3u8
	HLSJobID         string    `json:"hlsJobId,omitempty" dynamodbav:"hlsJobId,omitempty"`             // MediaConvert job ID
	HLSTranscodedAt  *time.Time `json:"hlsTranscodedAt,omitempty" dynamodbav:"hlsTranscodedAt,omitempty"`

	// DJ features
	HotCues map[int]*HotCue `json:"hotCues,omitempty" dynamodbav:"hotCues,omitempty"` // Slot (1-8) -> HotCue

	// Waveform and analysis fields
	WaveformURL    string     `json:"waveformUrl,omitempty" dynamodbav:"waveformUrl,omitempty"`       // S3 URL to waveform JSON
	BeatGrid       []int64    `json:"beatGrid,omitempty" dynamodbav:"beatGrid,omitempty"`             // Beat timestamps in milliseconds
	AnalysisStatus string     `json:"analysisStatus,omitempty" dynamodbav:"analysisStatus,omitempty"` // PENDING, ANALYZING, COMPLETED, FAILED
	AnalyzedAt     *time.Time `json:"analyzedAt,omitempty" dynamodbav:"analyzedAt,omitempty"`         // When analysis completed

	// Visibility fields (admin-panel-track-visibility feature)
	Visibility  TrackVisibility `json:"visibility" dynamodbav:"Visibility"`                   // private, unlisted, public
	PublishedAt *time.Time      `json:"publishedAt,omitempty" dynamodbav:"PublishedAt,omitempty"` // When track was made public

	// For API responses when admin/global views all tracks (not stored in DynamoDB)
	OwnerDisplayName string `json:"ownerDisplayName,omitempty" dynamodbav:"-"`

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

	// Set GSI1 for artist-based queries
	if track.Artist != "" {
		item.GSI1PK = fmt.Sprintf("USER#%s#ARTIST#%s", track.UserID, track.Artist)
		item.GSI1SK = fmt.Sprintf("TRACK#%s", track.ID)
	}

	// Set GSI3 for public track discovery (only when visibility is public)
	if track.Visibility == VisibilityPublic {
		item.GSI3PK = "PUBLIC_TRACK"
		// Sort by creation time for chronological discovery
		item.GSI3SK = fmt.Sprintf("%s#%s", track.CreatedAt.Format("2006-01-02T15:04:05Z"), track.ID)
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

// UpdateTrackVisibilityRequest represents a request to update track visibility
type UpdateTrackVisibilityRequest struct {
	Visibility string `json:"visibility" validate:"required,oneof=private unlisted public"`
}

// TrackResponse represents a track in API responses
type TrackResponse struct {
	ID           string                `json:"id"`
	Title        string                `json:"title"`
	Artist       string                `json:"artist"`
	ArtistID     string                `json:"artistId,omitempty"`
	Artists      []ArtistContribution  `json:"artists,omitempty"`
	AlbumArtist  string                `json:"albumArtist,omitempty"`
	Album        string                `json:"album,omitempty"`
	AlbumID      string                `json:"albumId,omitempty"`
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
	BPM          int       `json:"bpm,omitempty"`
	MusicalKey   string    `json:"musicalKey,omitempty"`
	KeyMode      string    `json:"keyMode,omitempty"`
	KeyCamelot   string    `json:"keyCamelot,omitempty"`
	HLSStatus      string     `json:"hlsStatus,omitempty"`
	HLSReady       bool       `json:"hlsReady"`
	WaveformURL    string     `json:"waveformUrl,omitempty"`
	AnalysisStatus string     `json:"analysisStatus,omitempty"`
	AnalyzedAt     *time.Time `json:"analyzedAt,omitempty"`
	// Visibility fields
	Visibility       string     `json:"visibility"`
	PublishedAt      *time.Time `json:"publishedAt,omitempty"`
	OwnerDisplayName string     `json:"ownerDisplayName,omitempty"` // Populated for admin/global views
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// ToResponse converts a Track to a TrackResponse
func (t *Track) ToResponse(coverArtURL string) TrackResponse {
	// Ensure tags is never nil to avoid undefined in JavaScript
	tags := t.Tags
	if tags == nil {
		tags = []string{}
	}

	// Default visibility to private if not set
	visibility := string(t.Visibility)
	if visibility == "" {
		visibility = string(VisibilityPrivate)
	}

	return TrackResponse{
		ID:           t.ID,
		Title:        t.Title,
		Artist:       t.Artist,
		ArtistID:     t.ArtistID,
		Artists:      t.Artists,
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
		BPM:          t.BPM,
		MusicalKey:   t.MusicalKey,
		KeyMode:      t.KeyMode,
		KeyCamelot:   t.KeyCamelot,
		HLSStatus:      string(t.HLSStatus),
		HLSReady:       t.HLSStatus == HLSStatusReady,
		WaveformURL:    t.WaveformURL,
		AnalysisStatus: t.AnalysisStatus,
		AnalyzedAt:     t.AnalyzedAt,
		Visibility:       visibility,
		PublishedAt:      t.PublishedAt,
		OwnerDisplayName: t.OwnerDisplayName,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

// TrackFilter represents filter options for listing tracks
type TrackFilter struct {
	Artist      string   `query:"artist"`
	Album       string   `query:"album"`
	Genre       string   `query:"genre"`
	Year        int      `query:"year"`
	Tags        []string `query:"tags"`
	BPMMin      int      `query:"bpmMin"`      // Minimum BPM filter
	BPMMax      int      `query:"bpmMax"`      // Maximum BPM filter
	MusicalKey  string   `query:"musicalKey"`  // Filter by musical key (e.g., "Am", "C")
	SortBy      string   `query:"sortBy"`      // title, artist, album, createdAt, playCount, bpm
	SortOrder   string   `query:"sortOrder"`   // asc, desc
	Limit       int      `query:"limit"`
	LastKey     string   `query:"lastKey"`
	GlobalScope bool     `query:"-"` // If true, return tracks from all users (requires GLOBAL permission)

	// Visibility filtering (admin-panel-track-visibility feature)
	IncludePublic bool   `query:"includePublic"` // Include public tracks from other users
	OwnerID       string `query:"ownerId"`       // Filter by specific owner (for admin)
	Visibility    string `query:"visibility"`    // Filter by visibility: private, unlisted, public
}

// Track visibility helper methods

// IsPubliclyAccessible returns true if the track can be accessed by non-owners.
// Both public and unlisted tracks are accessible if you have the link.
func (t *Track) IsPubliclyAccessible() bool {
	return t.Visibility.IsPubliclyAccessible()
}

// IsDiscoverable returns true if the track appears in search results and public listings.
// Only public tracks are discoverable.
func (t *Track) IsDiscoverable() bool {
	return t.Visibility.IsDiscoverable()
}

// GetVisibility returns the track's visibility, defaulting to private if not set.
func (t *Track) GetVisibility() TrackVisibility {
	if t.Visibility == "" {
		return DefaultTrackVisibility()
	}
	return t.Visibility
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
