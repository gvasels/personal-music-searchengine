package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ArtistWatch represents a user watching an artist for events.
type ArtistWatch struct {
	UserID     string    `json:"userId" dynamodbav:"userId"`
	ArtistName string    `json:"artistName" dynamodbav:"artistName"`
	WatchedAt  time.Time `json:"watchedAt" dynamodbav:"watchedAt"`
}

// ArtistWatchItem represents an ArtistWatch in DynamoDB single-table design
type ArtistWatchItem struct {
	DynamoDBItem
	ArtistWatch
}

// NewArtistWatch creates a new ArtistWatch.
func NewArtistWatch(userID, artistName string) *ArtistWatch {
	return &ArtistWatch{
		UserID:     userID,
		ArtistName: artistName,
		WatchedAt:  time.Now(),
	}
}

// NewArtistWatchItem creates a DynamoDB item for an artist watch.
// Primary key pattern: PK=USER#{userID}, SK=ARTIST_WATCH#{normalizedArtistName}
// GSI1 pattern: GSI1PK=ARTIST_WATCH#{normalizedArtistName}, GSI1SK=USER#{userID}
func NewArtistWatchItem(watch ArtistWatch) ArtistWatchItem {
	normalized := NormalizeArtistName(watch.ArtistName)
	return ArtistWatchItem{
		DynamoDBItem: DynamoDBItem{
			PK:     GetArtistWatchPK(watch.UserID),
			SK:     GetArtistWatchSK(normalized),
			GSI1PK: GetArtistWatchGSI1PK(normalized),
			GSI1SK: fmt.Sprintf("USER#%s", watch.UserID),
			Type:   string(EntityArtistWatch),
		},
		ArtistWatch: watch,
	}
}

// NormalizeArtistName normalizes an artist name for use as a DynamoDB key.
// It trims whitespace and converts to lowercase.
func NormalizeArtistName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Validate checks if the artist watch is valid.
func (w *ArtistWatch) Validate() error {
	if w.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if strings.TrimSpace(w.ArtistName) == "" {
		return errors.New("artist name cannot be empty")
	}
	return nil
}

// GetArtistWatchPK returns the partition key for querying a user's watched artists.
func GetArtistWatchPK(userID string) string {
	return fmt.Sprintf("USER#%s", userID)
}

// GetArtistWatchSK returns the sort key for an artist watch.
func GetArtistWatchSK(normalizedName string) string {
	return fmt.Sprintf("ARTIST_WATCH#%s", normalizedName)
}

// GetArtistWatchGSI1PK returns the GSI1 partition key for querying watchers of an artist.
func GetArtistWatchGSI1PK(normalizedName string) string {
	return fmt.Sprintf("ARTIST_WATCH#%s", normalizedName)
}

// ArtistWatchResponse represents an artist watch in API responses.
type ArtistWatchResponse struct {
	ArtistName string    `json:"artistName"`
	WatchedAt  time.Time `json:"watchedAt"`
}

// ToResponse converts an ArtistWatch to an ArtistWatchResponse.
func (w *ArtistWatch) ToResponse() ArtistWatchResponse {
	return ArtistWatchResponse{
		ArtistName: w.ArtistName,
		WatchedAt:  w.WatchedAt,
	}
}

// WatchedArtistListResponse represents a list of watched artists in API responses.
type WatchedArtistListResponse struct {
	Artists    []ArtistWatchResponse `json:"artists"`
	TotalCount int                   `json:"totalCount"`
	NextKey    string                `json:"nextKey,omitempty"`
}
