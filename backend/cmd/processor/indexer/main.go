package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	TrackID   string                 `json:"trackId"`
	UserID    string                 `json:"userId"`
	Metadata  *models.UploadMetadata `json:"metadata"`
	TableName string                 `json:"tableName"`
}

// Response represents the output to Step Functions
type Response struct {
	Indexed bool   `json:"indexed"`
	Reason  string `json:"reason,omitempty"`
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context (5 seconds less than Lambda timeout)
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	// Stub implementation - search indexing will be implemented in Epic 3
	// This allows the Step Functions workflow to continue without blocking
	_ = ctx // Use ctx to prevent unused variable warning

	return &Response{
		Indexed: false,
		Reason:  "not_implemented",
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
