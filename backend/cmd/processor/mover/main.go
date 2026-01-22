package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID   string `json:"uploadId"`
	UserID     string `json:"userId"`
	SourceKey  string `json:"sourceKey"`
	TrackID    string `json:"trackId"` // Direct trackId from Step Functions
	BucketName string `json:"bucketName"`
}

// Response represents the output to Step Functions
type Response struct {
	NewKey string `json:"newKey"` // Matches Step Functions expected output
}

var s3Client *s3.Client
var repo repository.Repository

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	s3Client = s3.NewFromConfig(cfg)

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

	if event.TrackID == "" {
		return nil, fmt.Errorf("track ID is required")
	}

	// Validate UUIDs
	if err := validation.ValidateUUID(event.TrackID, "trackId"); err != nil {
		return nil, err
	}
	if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
		return nil, err
	}

	// Determine file extension from source key
	ext := filepath.Ext(event.SourceKey)
	if ext == "" {
		ext = ".mp3" // Default extension
	}

	// Create destination key
	destKey := fmt.Sprintf("media/%s/%s%s", event.UserID, event.TrackID, ext)

	// Copy file to new location
	copySource := fmt.Sprintf("%s/%s", event.BucketName, event.SourceKey)
	_, err := s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &event.BucketName,
		CopySource: aws.String(copySource),
		Key:        &destKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Delete original file
	_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &event.BucketName,
		Key:    &event.SourceKey,
	})
	if err != nil {
		// Log error but don't fail - file is already copied
		fmt.Printf("Warning: failed to delete original file: %v\n", err)
	}

	// Update track with new S3 key
	track, err := repo.GetTrack(ctx, event.UserID, event.TrackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	track.S3Key = destKey
	if err := repo.UpdateTrack(ctx, *track); err != nil {
		return nil, fmt.Errorf("failed to update track S3 key: %w", err)
	}

	return &Response{NewKey: destKey}, nil
}

// getFileExtension extracts file extension from a path
func getFileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

func main() {
	lambda.Start(handleRequest)
}
