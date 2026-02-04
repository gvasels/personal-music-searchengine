package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/search"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	TrackID   string                 `json:"trackId"`
	UserID    string                 `json:"userId"`
	UploadID  string                 `json:"uploadId"`
	Metadata  *models.UploadMetadata `json:"metadata"`
	S3Key     string                 `json:"s3Key"`
	TableName string                 `json:"tableName"`
}

// Response represents the output to Step Functions
type Response struct {
	Indexed bool   `json:"indexed"`
	Reason  string `json:"reason,omitempty"`
}

var searchClient *search.Client
var repo repository.Repository

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "MusicLibrary"
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo = repository.NewDynamoDBRepository(dynamoClient, tableName)

	nixieFunctionName := os.Getenv("NIXIESEARCH_FUNCTION_NAME")
	if nixieFunctionName == "" {
		fmt.Println("NIXIESEARCH_FUNCTION_NAME not set, search indexing disabled")
		return
	}

	lambdaClient := awslambda.NewFromConfig(cfg)
	searchClient = search.NewClient(lambdaClient, nixieFunctionName)
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context (5 seconds less than Lambda timeout)
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	// Validate required fields
	if err := validation.ValidateUUID(event.TrackID, "trackId"); err != nil {
		return &Response{
			Indexed: false,
			Reason:  err.Error(),
		}, nil
	}

	if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
		return &Response{
			Indexed: false,
			Reason:  err.Error(),
		}, nil
	}

	// If search client not initialized, skip indexing
	if searchClient == nil {
		return &Response{
			Indexed: false,
			Reason:  "search_disabled",
		}, nil
	}

	// Validate metadata is present
	if event.Metadata == nil {
		return &Response{
			Indexed: false,
			Reason:  "missing_metadata",
		}, nil
	}

	// Build search document from metadata
	doc := search.Document{
		ID:        event.TrackID,
		UserID:    event.UserID,
		Title:     event.Metadata.Title,
		Artist:    event.Metadata.Artist,
		Album:     event.Metadata.Album,
		Genre:     event.Metadata.Genre,
		Year:      event.Metadata.Year,
		Duration:  event.Metadata.Duration,
		Filename:  event.S3Key,
		IndexedAt: time.Now(),
	}

	// Index the document
	resp, err := searchClient.Index(ctx, doc)
	if err != nil {
		return &Response{
			Indexed: false,
			Reason:  fmt.Sprintf("index_failed: %v", err),
		}, nil
	}

	if !resp.Indexed {
		return &Response{
			Indexed: false,
			Reason:  "index_rejected",
		}, nil
	}

	// Update step progress
	if event.UploadID != "" && repo != nil {
		if err := repo.UpdateUploadStep(ctx, event.UserID, event.UploadID, models.StepIndex, true); err != nil {
			fmt.Printf("Warning: failed to update step progress: %v\n", err)
		}
	}

	return &Response{
		Indexed: true,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
