package models

import (
	"fmt"
	"time"
)

// ArtistRole represents the role an artist plays on a track
type ArtistRole string

const (
	RoleMain      ArtistRole = "main"
	RoleFeaturing ArtistRole = "featuring"
	RoleRemixer   ArtistRole = "remixer"
	RoleProducer  ArtistRole = "producer"
)

// Artist represents a music artist in the library
type Artist struct {
	ID            string            `json:"id" dynamodbav:"id"`
	UserID        string            `json:"userId" dynamodbav:"userId"`
	Name          string            `json:"name" dynamodbav:"name"`
	SortName      string            `json:"sortName" dynamodbav:"sortName"`
	Bio           string            `json:"bio,omitempty" dynamodbav:"bio,omitempty"`
	ImageURL      string            `json:"imageUrl,omitempty" dynamodbav:"imageUrl,omitempty"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty" dynamodbav:"externalLinks,omitempty"`
	IsActive      bool              `json:"isActive" dynamodbav:"isActive"`
	Timestamps
}

// ArtistItem represents an Artist in DynamoDB single-table design
type ArtistItem struct {
	DynamoDBItem
	Artist
}

// NewArtistItem creates a DynamoDB item for an artist
// PK: USER#{userId}, SK: ARTIST#{artistId}
// GSI1PK: USER#{userId}#ARTIST, GSI1SK: name (for name lookups and sorting)
func NewArtistItem(artist Artist) ArtistItem {
	item := ArtistItem{
		DynamoDBItem: DynamoDBItem{
			PK:     fmt.Sprintf("USER#%s", artist.UserID),
			SK:     fmt.Sprintf("ARTIST#%s", artist.ID),
			Type:   string(EntityArtist),
			GSI1PK: fmt.Sprintf("USER#%s#ARTIST", artist.UserID),
			GSI1SK: artist.Name,
		},
		Artist: artist,
	}
	return item
}

// ArtistContribution represents an artist's contribution to a track
type ArtistContribution struct {
	ArtistID   string     `json:"artistId" dynamodbav:"artistId"`
	ArtistName string     `json:"artistName,omitempty" dynamodbav:"artistName,omitempty"` // Denormalized for display
	Role       ArtistRole `json:"role" dynamodbav:"role"`
}

// ArtistWithStats represents an artist with aggregated statistics
type ArtistWithStats struct {
	Artist
	TrackCount int `json:"trackCount"`
	AlbumCount int `json:"albumCount"`
	TotalPlays int `json:"totalPlays"`
}

// CreateArtistRequest represents a request to create an artist
type CreateArtistRequest struct {
	Name          string            `json:"name" validate:"required,min=1,max=500"`
	SortName      string            `json:"sortName,omitempty" validate:"omitempty,max=500"`
	Bio           string            `json:"bio,omitempty" validate:"omitempty,max=5000"`
	ImageURL      string            `json:"imageUrl,omitempty" validate:"omitempty,url,max=2000"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"`
}

// UpdateArtistRequest represents a request to update an artist
type UpdateArtistRequest struct {
	Name          *string           `json:"name,omitempty" validate:"omitempty,min=1,max=500"`
	SortName      *string           `json:"sortName,omitempty" validate:"omitempty,max=500"`
	Bio           *string           `json:"bio,omitempty" validate:"omitempty,max=5000"`
	ImageURL      *string           `json:"imageUrl,omitempty" validate:"omitempty,url,max=2000"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"`
}

// ArtistResponse represents an artist in API responses
type ArtistResponse struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	SortName      string            `json:"sortName"`
	Bio           string            `json:"bio,omitempty"`
	ImageURL      string            `json:"imageUrl,omitempty"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"`
	IsActive      bool              `json:"isActive"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

// ArtistWithStatsResponse represents an artist with stats in API responses
type ArtistWithStatsResponse struct {
	ArtistResponse
	TrackCount int `json:"trackCount"`
	AlbumCount int `json:"albumCount"`
	TotalPlays int `json:"totalPlays"`
}

// ToResponse converts an Artist to an ArtistResponse
func (a *Artist) ToResponse() ArtistResponse {
	return ArtistResponse{
		ID:            a.ID,
		Name:          a.Name,
		SortName:      a.SortName,
		Bio:           a.Bio,
		ImageURL:      a.ImageURL,
		ExternalLinks: a.ExternalLinks,
		IsActive:      a.IsActive,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

// ToResponseWithStats converts an ArtistWithStats to an ArtistWithStatsResponse
func (a *ArtistWithStats) ToResponseWithStats() ArtistWithStatsResponse {
	return ArtistWithStatsResponse{
		ArtistResponse: a.Artist.ToResponse(),
		TrackCount:     a.TrackCount,
		AlbumCount:     a.AlbumCount,
		TotalPlays:     a.TotalPlays,
	}
}

// GenerateSortName generates a sort name from the artist name
// Removes common prefixes like "The", "A", "An" for better alphabetical sorting
func GenerateSortName(name string) string {
	if name == "" {
		return ""
	}

	// Common prefixes to strip for sorting
	prefixes := []string{"The ", "A ", "An "}
	sortName := name

	for _, prefix := range prefixes {
		if len(name) > len(prefix) && name[:len(prefix)] == prefix {
			sortName = name[len(prefix):]
			break
		}
	}

	return sortName
}
