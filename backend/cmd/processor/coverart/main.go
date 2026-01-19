package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gvasels/personal-music-searchengine/internal/metadata"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// Event represents the input from Step Functions
type Event struct {
	UploadID   string                 `json:"uploadId"`
	UserID     string                 `json:"userId"`
	S3Key      string                 `json:"s3Key"`
	Metadata   *models.UploadMetadata `json:"metadata"`
	BucketName string                 `json:"bucketName"`
}

// Response represents the output to Step Functions
type Response struct {
	CoverArtKey string `json:"coverArtKey"`
}

var s3Client *s3.Client
var extractor *metadata.Extractor

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	s3Client = s3.NewFromConfig(cfg)
	extractor = metadata.NewExtractor()
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Check if metadata indicates cover art is present
	if event.Metadata == nil || !event.Metadata.HasCoverArt {
		return &Response{CoverArtKey: ""}, nil
	}

	// Download file from S3
	data, err := downloadFromS3(ctx, event.BucketName, event.S3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	// Extract cover art
	reader := bytes.NewReader(data)
	coverData, mimeType, err := extractor.ExtractCoverArt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to extract cover art: %w", err)
	}

	if coverData == nil {
		return &Response{CoverArtKey: ""}, nil
	}

	// Determine file extension from MIME type
	ext := getExtensionFromMIME(mimeType)

	// Upload cover art to S3
	coverKey := fmt.Sprintf("covers/%s/%s%s", event.UserID, event.UploadID, ext)
	err = uploadToS3(ctx, event.BucketName, coverKey, coverData, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload cover art: %w", err)
	}

	return &Response{CoverArtKey: coverKey}, nil
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

func uploadToS3(ctx context.Context, bucket, key string, data []byte, contentType string) error {
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	return err
}

func getExtensionFromMIME(mimeType string) string {
	switch strings.ToLower(mimeType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func main() {
	lambda.Start(handleRequest)
}
