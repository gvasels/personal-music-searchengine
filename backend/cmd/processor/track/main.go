package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID   string                 `json:"uploadId"`
	UserID     string                 `json:"userId"`
	S3Key      string                 `json:"s3Key"`
	FileName   string                 `json:"fileName"`
	Metadata   *models.UploadMetadata `json:"metadata"`
	CoverArt   *CoverArtResult        `json:"coverArt"`
	BucketName string                 `json:"bucketName"`
	TableName  string                 `json:"tableName"`
}

// CoverArtResult represents the cover art extraction result
type CoverArtResult struct {
	CoverArtKey string `json:"coverArtKey"`
}

// Response represents the output to Step Functions
type Response struct {
	TrackID string `json:"trackId"`
	AlbumID string `json:"albumId,omitempty"`
}

var repo repository.Repository

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "MusicLibrary"
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo = repository.NewDynamoDBRepository(dynamoClient, tableName)
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context (5 seconds less than Lambda timeout)
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	// Validate input UUIDs to prevent injection attacks
	if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
		return nil, err
	}
	if err := validation.ValidateUUID(event.UploadID, "uploadId"); err != nil {
		return nil, err
	}

	trackID := uuid.New().String()
	now := time.Now()

	// Determine format from metadata
	format := models.AudioFormatMP3
	if event.Metadata != nil && event.Metadata.Format != "" {
		format = models.AudioFormat(event.Metadata.Format)
	}

	// Create track record
	track := models.Track{
		ID:        trackID,
		UserID:    event.UserID,
		Title:     getOrDefault(event.Metadata, "title", event.FileName),
		Artist:    getOrDefault(event.Metadata, "artist", "Unknown Artist"),
		Album:     getOrDefault(event.Metadata, "album", ""),
		Genre:     getOrDefault(event.Metadata, "genre", ""),
		Year:      getIntOrDefault(event.Metadata, "year", 0),
		Duration:  getIntOrDefault(event.Metadata, "duration", 0),
		Format:    format,
		S3Key:     event.S3Key, // Will be updated after file is moved
		PlayCount: 0,
	}
	track.CreatedAt = now
	track.UpdatedAt = now

	// Set cover art key if available
	if event.CoverArt != nil && event.CoverArt.CoverArtKey != "" {
		track.CoverArtKey = event.CoverArt.CoverArtKey
	}

	// Set additional metadata fields if available
	if event.Metadata != nil {
		track.Bitrate = event.Metadata.Bitrate
	}

	// Create the track
	if err := repo.CreateTrack(ctx, track); err != nil {
		return nil, fmt.Errorf("failed to create track: %w", err)
	}

	response := &Response{TrackID: trackID}

	// Create or update album if album name is present
	if track.Album != "" {
		album, err := repo.GetOrCreateAlbum(ctx, event.UserID, track.Album, track.Artist)
		if err != nil {
			// Log error but don't fail - track is already created
			fmt.Printf("Warning: failed to create/update album: %v\n", err)
		} else {
			response.AlbumID = album.ID
		}
	}

	return response, nil
}

func getOrDefault(meta *models.UploadMetadata, field, defaultVal string) string {
	if meta == nil {
		return defaultVal
	}
	switch field {
	case "title":
		if meta.Title != "" {
			return meta.Title
		}
	case "artist":
		if meta.Artist != "" {
			return meta.Artist
		}
	case "album":
		if meta.Album != "" {
			return meta.Album
		}
	case "genre":
		if meta.Genre != "" {
			return meta.Genre
		}
	}
	return defaultVal
}

func getIntOrDefault(meta *models.UploadMetadata, field string, defaultVal int) int {
	if meta == nil {
		return defaultVal
	}
	switch field {
	case "year":
		if meta.Year != 0 {
			return meta.Year
		}
	case "duration":
		if meta.Duration != 0 {
			return meta.Duration
		}
	}
	return defaultVal
}

func main() {
	lambda.Start(handleRequest)
}
