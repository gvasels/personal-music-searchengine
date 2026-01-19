package models

import "time"

// StreamRequest represents a request for a streaming URL
type StreamRequest struct {
	TrackID string `param:"trackId" validate:"required,uuid"`
	Quality string `query:"quality,omitempty"` // original, high, medium, low
}

// StreamResponse represents a response with streaming URL
type StreamResponse struct {
	TrackID   string    `json:"trackId"`
	StreamURL string    `json:"streamUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
	Format    string    `json:"format"`
	Bitrate   int       `json:"bitrate,omitempty"`
}

// DownloadRequest represents a request for a download URL
type DownloadRequest struct {
	TrackID string `param:"trackId" validate:"required,uuid"`
}

// DownloadResponse represents a response with download URL
type DownloadResponse struct {
	TrackID     string    `json:"trackId"`
	DownloadURL string    `json:"downloadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
	FileName    string    `json:"fileName"`
	FileSize    int64     `json:"fileSize"`
	Format      string    `json:"format"`
}

// PlaybackEvent represents a playback event for analytics
type PlaybackEvent struct {
	UserID    string    `json:"userId"`
	TrackID   string    `json:"trackId"`
	EventType string    `json:"eventType"` // play, pause, stop, seek, complete
	Position  int       `json:"position"`  // playback position in seconds
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId,omitempty"`
	Device    string    `json:"device,omitempty"`
}

// RecordPlayRequest represents a request to record a play
type RecordPlayRequest struct {
	TrackID   string `json:"trackId" validate:"required,uuid"`
	Duration  int    `json:"duration" validate:"required,min=1"` // listened duration in seconds
	Completed bool   `json:"completed"`
}

// PlayQueue represents a user's play queue
type PlayQueue struct {
	UserID       string   `json:"userId"`
	CurrentIndex int      `json:"currentIndex"`
	TrackIDs     []string `json:"trackIds"`
	ShuffleMode  bool     `json:"shuffleMode"`
	RepeatMode   string   `json:"repeatMode"` // none, one, all
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UpdateQueueRequest represents a request to update the play queue
type UpdateQueueRequest struct {
	TrackIDs     []string `json:"trackIds,omitempty"`
	CurrentIndex *int     `json:"currentIndex,omitempty"`
	ShuffleMode  *bool    `json:"shuffleMode,omitempty"`
	RepeatMode   *string  `json:"repeatMode,omitempty"`
}

// QueueAction represents an action on the play queue
type QueueAction string

const (
	QueueActionAddNext   QueueAction = "add_next"
	QueueActionAddLast   QueueAction = "add_last"
	QueueActionRemove    QueueAction = "remove"
	QueueActionClear     QueueAction = "clear"
	QueueActionShuffle   QueueAction = "shuffle"
	QueueActionUnshuffle QueueAction = "unshuffle"
)

// QueueActionRequest represents a request to perform a queue action
type QueueActionRequest struct {
	Action   QueueAction `json:"action" validate:"required"`
	TrackIDs []string    `json:"trackIds,omitempty" validate:"omitempty,dive,uuid"`
}
