package models

import (
	"fmt"
	"time"
)

// Album represents an album in the library
type Album struct {
	ID           string `json:"id" dynamodbav:"id"`
	UserID       string `json:"userId" dynamodbav:"userId"`
	Title        string `json:"title" dynamodbav:"title"`
	Artist       string `json:"artist" dynamodbav:"artist"`
	AlbumArtist  string `json:"albumArtist,omitempty" dynamodbav:"albumArtist,omitempty"`
	Genre        string `json:"genre,omitempty" dynamodbav:"genre,omitempty"`
	Year         int    `json:"year,omitempty" dynamodbav:"year,omitempty"`
	CoverArtKey  string `json:"coverArtKey,omitempty" dynamodbav:"coverArtKey,omitempty"`
	TrackCount   int    `json:"trackCount" dynamodbav:"trackCount"`
	TotalDuration int   `json:"totalDuration" dynamodbav:"totalDuration"` // seconds
	DiscCount    int    `json:"discCount" dynamodbav:"discCount"`
	Timestamps
}

// AlbumItem represents an Album in DynamoDB single-table design
type AlbumItem struct {
	DynamoDBItem
	Album
}

// NewAlbumItem creates a DynamoDB item for an album
func NewAlbumItem(album Album) AlbumItem {
	item := AlbumItem{
		DynamoDBItem: DynamoDBItem{
			PK:   fmt.Sprintf("USER#%s", album.UserID),
			SK:   fmt.Sprintf("ALBUM#%s", album.ID),
			Type: string(EntityAlbum),
		},
		Album: album,
	}

	// Set GSI for artist-based queries
	if album.Artist != "" {
		item.GSI1PK = fmt.Sprintf("USER#%s#ARTIST#%s", album.UserID, album.Artist)
		item.GSI1SK = fmt.Sprintf("ALBUM#%d", album.Year)
	}

	return item
}

// AlbumResponse represents an album in API responses
type AlbumResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Artist        string    `json:"artist"`
	AlbumArtist   string    `json:"albumArtist,omitempty"`
	Genre         string    `json:"genre,omitempty"`
	Year          int       `json:"year,omitempty"`
	CoverArtURL   string    `json:"coverArtUrl,omitempty"`
	TrackCount    int       `json:"trackCount"`
	TotalDuration int       `json:"totalDuration"`
	DurationStr   string    `json:"durationStr"`
	DiscCount     int       `json:"discCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ToResponse converts an Album to an AlbumResponse
func (a *Album) ToResponse(coverArtURL string) AlbumResponse {
	return AlbumResponse{
		ID:            a.ID,
		Title:         a.Title,
		Artist:        a.Artist,
		AlbumArtist:   a.AlbumArtist,
		Genre:         a.Genre,
		Year:          a.Year,
		CoverArtURL:   coverArtURL,
		TrackCount:    a.TrackCount,
		TotalDuration: a.TotalDuration,
		DurationStr:   formatDuration(a.TotalDuration),
		DiscCount:     a.DiscCount,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

// AlbumWithTracks represents an album with its tracks
type AlbumWithTracks struct {
	Album  AlbumResponse   `json:"album"`
	Tracks []TrackResponse `json:"tracks"`
}

// AlbumFilter represents filter options for listing albums
type AlbumFilter struct {
	Artist    string `query:"artist"`
	Genre     string `query:"genre"`
	Year      int    `query:"year"`
	SortBy    string `query:"sortBy"`    // title, artist, year, createdAt
	SortOrder string `query:"sortOrder"` // asc, desc
	Limit     int    `query:"limit"`
	LastKey   string `query:"lastKey"`
}

// ArtistSummary represents an artist with aggregated stats
type ArtistSummary struct {
	Name       string `json:"name"`
	TrackCount int    `json:"trackCount"`
	AlbumCount int    `json:"albumCount"`
	CoverArtURL string `json:"coverArtUrl,omitempty"`
}

// ArtistFilter represents filter options for listing artists
type ArtistFilter struct {
	SortBy    string `query:"sortBy"`    // name, trackCount, albumCount
	SortOrder string `query:"sortOrder"` // asc, desc
	Limit     int    `query:"limit"`
	LastKey   string `query:"lastKey"`
}
