package models

import "time"

// HotCue constants
const (
	MaxHotCuesPerTrack = 8
)

// HotCueColor represents preset colors for hot cues
type HotCueColor string

const (
	HotCueColorRed    HotCueColor = "#FF0000"
	HotCueColorOrange HotCueColor = "#FF8C00"
	HotCueColorYellow HotCueColor = "#FFFF00"
	HotCueColorGreen  HotCueColor = "#00FF00"
	HotCueColorCyan   HotCueColor = "#00FFFF"
	HotCueColorBlue   HotCueColor = "#0000FF"
	HotCueColorPurple HotCueColor = "#800080"
	HotCueColorPink   HotCueColor = "#FF69B4"
)

// DefaultHotCueColors returns the default color palette
func DefaultHotCueColors() []HotCueColor {
	return []HotCueColor{
		HotCueColorRed,
		HotCueColorOrange,
		HotCueColorYellow,
		HotCueColorGreen,
		HotCueColorCyan,
		HotCueColorBlue,
		HotCueColorPurple,
		HotCueColorPink,
	}
}

// GetDefaultColorForSlot returns the default color for a slot number (1-8)
func GetDefaultColorForSlot(slot int) HotCueColor {
	colors := DefaultHotCueColors()
	if slot >= 1 && slot <= len(colors) {
		return colors[slot-1]
	}
	return HotCueColorRed
}

// HotCue represents a hot cue point on a track
type HotCue struct {
	Slot      int         `json:"slot" dynamodbav:"slot"`                         // 1-8
	Position  float64     `json:"position" dynamodbav:"position"`                 // Position in seconds
	Label     string      `json:"label,omitempty" dynamodbav:"label,omitempty"`   // Optional label
	Color     HotCueColor `json:"color" dynamodbav:"color"`                       // Display color
	CreatedAt time.Time   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt" dynamodbav:"updatedAt"`
}

// IsValidSlot checks if the slot number is valid (1-8)
func IsValidSlot(slot int) bool {
	return slot >= 1 && slot <= MaxHotCuesPerTrack
}

// SetHotCueRequest represents a request to set a hot cue
type SetHotCueRequest struct {
	Position float64     `json:"position" validate:"required,gte=0"`
	Label    string      `json:"label,omitempty" validate:"omitempty,max=50"`
	Color    HotCueColor `json:"color,omitempty"`
}

// HotCueResponse represents a hot cue in API responses
type HotCueResponse struct {
	Slot      int         `json:"slot"`
	Position  float64     `json:"position"`
	Label     string      `json:"label,omitempty"`
	Color     HotCueColor `json:"color"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// TrackHotCuesResponse represents all hot cues for a track
type TrackHotCuesResponse struct {
	TrackID  string           `json:"trackId"`
	HotCues  []HotCueResponse `json:"hotCues"`
	MaxSlots int              `json:"maxSlots"`
}
