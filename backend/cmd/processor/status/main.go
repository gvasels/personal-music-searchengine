package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID  string `json:"uploadId"`
	UserID    string `json:"userId"`
	TrackID   string `json:"trackId,omitempty"`
	Status    string `json:"status"`
	Error     *Error `json:"error,omitempty"`
	TableName string `json:"tableName"`
}

// Error represents error information from Step Functions
type Error struct {
	Error string `json:"Error"`
	Cause string `json:"Cause"`
}

// Response represents the output
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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

	var status models.UploadStatus
	var errorMsg string

	switch event.Status {
	case "COMPLETED":
		status = models.UploadStatusCompleted
	case "FAILED":
		status = models.UploadStatusFailed
		if event.Error != nil {
			errorMsg = event.Error.Error
			if event.Error.Cause != "" {
				errorMsg = fmt.Sprintf("%s: %s", errorMsg, event.Error.Cause)
			}
		}
	default:
		status = models.UploadStatus(event.Status)
	}

	// Update upload status
	err := repo.UpdateUploadStatus(ctx, event.UserID, event.UploadID, status, event.TrackID, errorMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to update upload status: %w", err)
	}

	// If completed, also update the completion timestamp and step flags
	if status == models.UploadStatusCompleted {
		upload, err := repo.GetUpload(ctx, event.UserID, event.UploadID)
		if err == nil {
			now := time.Now()
			upload.CompletedAt = &now
			upload.MetadataExtracted = true
			upload.CoverArtExtracted = true // May be false if no cover art, but step completed
			upload.TrackCreated = true
			upload.FileMoved = true
			upload.Indexed = false // Stub for now
			upload.TrackID = event.TrackID
			upload.Status = status

			if err := repo.UpdateUpload(ctx, *upload); err != nil {
				fmt.Printf("Warning: failed to update upload details: %v\n", err)
			}
		}
	}

	return &Response{
		Success: true,
		Message: fmt.Sprintf("Upload %s status updated to %s", event.UploadID, status),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
