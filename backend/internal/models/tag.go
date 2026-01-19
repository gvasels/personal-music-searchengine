package models

import (
	"fmt"
	"time"
)

// Tag represents a user-defined tag for categorizing tracks
type Tag struct {
	UserID     string `json:"userId" dynamodbav:"userId"`
	Name       string `json:"name" dynamodbav:"name"`
	Color      string `json:"color,omitempty" dynamodbav:"color,omitempty"` // hex color code
	TrackCount int    `json:"trackCount" dynamodbav:"trackCount"`
	Timestamps
}

// TagItem represents a Tag in DynamoDB single-table design
type TagItem struct {
	DynamoDBItem
	Tag
}

// NewTagItem creates a DynamoDB item for a tag
func NewTagItem(tag Tag) TagItem {
	return TagItem{
		DynamoDBItem: DynamoDBItem{
			PK:   fmt.Sprintf("USER#%s", tag.UserID),
			SK:   fmt.Sprintf("TAG#%s", tag.Name),
			Type: string(EntityTag),
		},
		Tag: tag,
	}
}

// TrackTag represents the association between a track and a tag
type TrackTag struct {
	UserID   string    `json:"userId" dynamodbav:"userId"`
	TrackID  string    `json:"trackId" dynamodbav:"trackId"`
	TagName  string    `json:"tagName" dynamodbav:"tagName"`
	AddedAt  time.Time `json:"addedAt" dynamodbav:"addedAt"`
}

// TrackTagItem represents a TrackTag in DynamoDB single-table design
type TrackTagItem struct {
	DynamoDBItem
	TrackTag
}

// NewTrackTagItem creates a DynamoDB item for a track-tag association
func NewTrackTagItem(tt TrackTag) TrackTagItem {
	return TrackTagItem{
		DynamoDBItem: DynamoDBItem{
			PK:     fmt.Sprintf("USER#%s#TRACK#%s", tt.UserID, tt.TrackID),
			SK:     fmt.Sprintf("TAG#%s", tt.TagName),
			GSI1PK: fmt.Sprintf("USER#%s#TAG#%s", tt.UserID, tt.TagName),
			GSI1SK: fmt.Sprintf("TRACK#%s", tt.TrackID),
			Type:   string(EntityTrackTag),
		},
		TrackTag: tt,
	}
}

// CreateTagRequest represents a request to create a tag
type CreateTagRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Color string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

// UpdateTagRequest represents a request to update a tag
type UpdateTagRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	Color *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

// AddTagsToTrackRequest represents a request to add tags to a track
type AddTagsToTrackRequest struct {
	Tags []string `json:"tags" validate:"required,min=1,max=20,dive,min=1,max=50"`
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	Name       string    `json:"name"`
	Color      string    `json:"color,omitempty"`
	TrackCount int       `json:"trackCount"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ToResponse converts a Tag to a TagResponse
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		Name:       t.Name,
		Color:      t.Color,
		TrackCount: t.TrackCount,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

// TagFilter represents filter options for listing tags
type TagFilter struct {
	SortBy    string `query:"sortBy"`    // name, trackCount, createdAt
	SortOrder string `query:"sortOrder"` // asc, desc
	Limit     int    `query:"limit"`
	LastKey   string `query:"lastKey"`
}
