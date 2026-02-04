package models

import (
	"fmt"
	"time"
)

// Constants for crate limitations
const (
	MaxTracksPerCrate = 1000
	MaxCratesPerUser  = 100
)

// CrateSortOrder represents how tracks are sorted in a crate
type CrateSortOrder string

const (
	CrateSortCustom   CrateSortOrder = "custom"   // Manual drag-and-drop order
	CrateSortBPM      CrateSortOrder = "bpm"      // Sort by BPM
	CrateSortKey      CrateSortOrder = "key"      // Sort by musical key
	CrateSortArtist   CrateSortOrder = "artist"   // Sort by artist name
	CrateSortTitle    CrateSortOrder = "title"    // Sort by track title
	CrateSortAdded    CrateSortOrder = "added"    // Sort by date added to crate
)

// Crate represents a DJ crate (collection of tracks)
type Crate struct {
	ID          string         `json:"id" dynamodbav:"id"`
	UserID      string         `json:"userId" dynamodbav:"userId"`
	Name        string         `json:"name" dynamodbav:"name"`
	Description string         `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Color       string         `json:"color,omitempty" dynamodbav:"color,omitempty"` // Hex color for UI
	TrackIDs    []string       `json:"trackIds" dynamodbav:"trackIds"`               // Ordered list of track IDs
	TrackCount  int            `json:"trackCount" dynamodbav:"trackCount"`
	SortOrder   CrateSortOrder `json:"sortOrder" dynamodbav:"sortOrder"`
	IsSmartCrate bool          `json:"isSmartCrate" dynamodbav:"isSmartCrate"`       // Auto-populated based on criteria
	SmartCriteria *SmartCrateCriteria `json:"smartCriteria,omitempty" dynamodbav:"smartCriteria,omitempty"`
	Timestamps
}

// SmartCrateCriteria defines rules for auto-populating a smart crate
type SmartCrateCriteria struct {
	BPMMin     int      `json:"bpmMin,omitempty" dynamodbav:"bpmMin,omitempty"`
	BPMMax     int      `json:"bpmMax,omitempty" dynamodbav:"bpmMax,omitempty"`
	Keys       []string `json:"keys,omitempty" dynamodbav:"keys,omitempty"`
	Genres     []string `json:"genres,omitempty" dynamodbav:"genres,omitempty"`
	Tags       []string `json:"tags,omitempty" dynamodbav:"tags,omitempty"`
	MinRating  int      `json:"minRating,omitempty" dynamodbav:"minRating,omitempty"`
}

// CrateItem represents a Crate in DynamoDB single-table design
// PK: USER#{userId}, SK: CRATE#{crateId}
type CrateItem struct {
	DynamoDBItem
	Crate
}

// NewCrateItem creates a DynamoDB item for a crate
func NewCrateItem(crate Crate) CrateItem {
	return CrateItem{
		DynamoDBItem: DynamoDBItem{
			PK:   fmt.Sprintf("USER#%s", crate.UserID),
			SK:   fmt.Sprintf("CRATE#%s", crate.ID),
			Type: "CRATE",
		},
		Crate: crate,
	}
}

// CrateResponse represents a crate in API responses
type CrateResponse struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	Color         string         `json:"color,omitempty"`
	TrackCount    int            `json:"trackCount"`
	SortOrder     CrateSortOrder `json:"sortOrder"`
	IsSmartCrate  bool           `json:"isSmartCrate"`
	SmartCriteria *SmartCrateCriteria `json:"smartCriteria,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

// ToResponse converts a Crate to a CrateResponse
func (c *Crate) ToResponse() CrateResponse {
	return CrateResponse{
		ID:            c.ID,
		Name:          c.Name,
		Description:   c.Description,
		Color:         c.Color,
		TrackCount:    c.TrackCount,
		SortOrder:     c.SortOrder,
		IsSmartCrate:  c.IsSmartCrate,
		SmartCriteria: c.SmartCriteria,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

// CreateCrateRequest represents a request to create a crate
type CreateCrateRequest struct {
	Name          string              `json:"name" validate:"required,min=1,max=100"`
	Description   string              `json:"description,omitempty" validate:"omitempty,max=500"`
	Color         string              `json:"color,omitempty" validate:"omitempty,hexcolor"`
	IsSmartCrate  bool                `json:"isSmartCrate,omitempty"`
	SmartCriteria *SmartCrateCriteria `json:"smartCriteria,omitempty"`
}

// UpdateCrateRequest represents a request to update a crate
type UpdateCrateRequest struct {
	Name          *string             `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description   *string             `json:"description,omitempty" validate:"omitempty,max=500"`
	Color         *string             `json:"color,omitempty" validate:"omitempty,hexcolor"`
	SortOrder     *CrateSortOrder     `json:"sortOrder,omitempty"`
	SmartCriteria *SmartCrateCriteria `json:"smartCriteria,omitempty"`
}

// AddTracksToCrateRequest represents a request to add tracks to a crate
type AddTracksToCrateRequest struct {
	TrackIDs []string `json:"trackIds" validate:"required,min=1,max=100"`
	Position int      `json:"position,omitempty"` // -1 or omitted = append at end
}

// RemoveTracksFromCrateRequest represents a request to remove tracks from a crate
type RemoveTracksFromCrateRequest struct {
	TrackIDs []string `json:"trackIds" validate:"required,min=1"`
}

// ReorderTracksRequest represents a request to reorder tracks in a crate
type ReorderTracksRequest struct {
	TrackIDs []string `json:"trackIds" validate:"required"` // New order of all track IDs
}

// CrateFilter represents filter options for listing crates
type CrateFilter struct {
	Limit int `query:"limit"`
}

// CrateWithTracksResponse includes the full track list
type CrateWithTracksResponse struct {
	CrateResponse
	Tracks []TrackResponse `json:"tracks"`
}
