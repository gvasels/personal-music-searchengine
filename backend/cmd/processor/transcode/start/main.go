package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents the input from Step Functions
type Event struct {
	TrackID   string `json:"trackId"`
	UserID    string `json:"userId"`
	S3Key     string `json:"s3Key"`
	TableName string `json:"tableName"`
}

// Response represents the output to Step Functions
type Response struct {
	JobID       string `json:"jobId,omitempty"`
	PlaylistKey string `json:"playlistKey,omitempty"`
	Status      string `json:"status"`
	Reason      string `json:"reason,omitempty"`
}

var (
	transcodeSvc *service.TranscodeService
	dynamoClient *dynamodb.Client
	tableName    string
)

func init() {
	mediaConvertEndpoint := os.Getenv("MEDIACONVERT_ENDPOINT")
	mediaConvertRole := os.Getenv("MEDIACONVERT_ROLE")
	mediaConvertQueue := os.Getenv("MEDIACONVERT_QUEUE")
	mediaBucket := os.Getenv("MEDIA_BUCKET")
	tableName = os.Getenv("DYNAMODB_TABLE_NAME")

	if mediaConvertEndpoint == "" || mediaConvertRole == "" || mediaBucket == "" {
		fmt.Println("MediaConvert configuration incomplete, transcoding disabled")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}

	// Create MediaConvert client with custom endpoint
	mcClient := mediaconvert.NewFromConfig(cfg, func(o *mediaconvert.Options) {
		o.BaseEndpoint = &mediaConvertEndpoint
	})

	transcodeSvc = service.NewTranscodeService(mcClient, mediaBucket, mediaConvertRole, mediaConvertQueue)
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	// Validate required fields
	if err := validation.ValidateUUID(event.TrackID, "trackId"); err != nil {
		return &Response{
			Status: "failed",
			Reason: err.Error(),
		}, nil
	}

	if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
		return &Response{
			Status: "failed",
			Reason: err.Error(),
		}, nil
	}

	if event.S3Key == "" {
		return &Response{
			Status: "failed",
			Reason: "s3Key is required",
		}, nil
	}

	// Check if transcode service is available
	if transcodeSvc == nil {
		return &Response{
			Status: "skipped",
			Reason: "transcode_disabled",
		}, nil
	}

	// Start transcode job
	req := service.TranscodeRequest{
		TrackID: event.TrackID,
		UserID:  event.UserID,
		S3Key:   event.S3Key,
	}

	resp, err := transcodeSvc.StartTranscode(ctx, req)
	if err != nil {
		return &Response{
			Status: "failed",
			Reason: fmt.Sprintf("transcode_failed: %v", err),
		}, nil
	}

	// Update track HLS status in DynamoDB
	if dynamoClient != nil && tableName != "" {
		if err := updateTrackHLSStatus(ctx, event.UserID, event.TrackID, models.HLSStatusProcessing, resp.JobID, resp.PlaylistKey); err != nil {
			fmt.Printf("Warning: failed to update track HLS status: %v\n", err)
			// Continue - job was created successfully
		}
	}

	return &Response{
		JobID:       resp.JobID,
		PlaylistKey: resp.PlaylistKey,
		Status:      "started",
	}, nil
}

func updateTrackHLSStatus(ctx context.Context, userID, trackID string, status models.HLSStatus, jobID, playlistKey string) error {
	if dynamoClient == nil || tableName == "" {
		return fmt.Errorf("DynamoDB not configured")
	}

	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("TRACK#%s", trackID)

	updateExpr := "SET hlsStatus = :status, hlsJobId = :jobId, hlsPlaylistKey = :playlist, updatedAt = :now"
	exprValues := map[string]dynamodbtypes.AttributeValue{
		":status":   &dynamodbtypes.AttributeValueMemberS{Value: string(status)},
		":jobId":    &dynamodbtypes.AttributeValueMemberS{Value: jobID},
		":playlist": &dynamodbtypes.AttributeValueMemberS{Value: playlistKey},
		":now":      &dynamodbtypes.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key: map[string]dynamodbtypes.AttributeValue{
			"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"SK": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:          stringPtr(updateExpr),
		ExpressionAttributeValues: exprValues,
	}

	_, err := dynamoClient.UpdateItem(ctx, input)
	return err
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	lambda.Start(handleRequest)
}
