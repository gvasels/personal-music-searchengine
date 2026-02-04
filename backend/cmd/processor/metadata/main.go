package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gvasels/personal-music-searchengine/internal/metadata"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID   string `json:"uploadId"`
	UserID     string `json:"userId"`
	S3Key      string `json:"s3Key"`
	FileName   string `json:"fileName"`
	BucketName string `json:"bucketName"`
}

// Response represents the output to Step Functions
type Response struct {
	*models.UploadMetadata
}

var s3Client *s3.Client
var extractor *metadata.Extractor
var repo repository.Repository

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	s3Client = s3.NewFromConfig(cfg)
	extractor = metadata.NewExtractor()

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

	// Validate file size before download to prevent OOM
	if err := validation.ValidateFileSize(ctx, s3Client, event.BucketName, event.S3Key); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// Download file from S3
	data, err := downloadFromS3(ctx, event.BucketName, event.S3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	// Extract metadata
	reader := bytes.NewReader(data)
	meta, err := extractor.Extract(reader, event.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Update step progress
	if err := repo.UpdateUploadStep(ctx, event.UserID, event.UploadID, models.StepExtractMetadata, true); err != nil {
		fmt.Printf("Warning: failed to update step progress: %v\n", err)
	}

	return &Response{UploadMetadata: meta}, nil
}

func downloadFromS3(ctx context.Context, bucket, key string) ([]byte, error) {
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func main() {
	lambda.Start(handleRequest)
}
