package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
)

// Event represents a MediaConvert EventBridge event
type Event = service.MediaConvertEvent

// Response represents the output from the Lambda
type Response struct {
	TrackID string `json:"trackId"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
}

var (
	dynamoClient *dynamodb.Client
	tableName    string
)

func init() {
	tableName = os.Getenv("DYNAMODB_TABLE_NAME")

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	detail := event.Detail

	// Extract track ID and user ID from job tags
	trackID := detail.UserMetadata["trackId"]
	userID := detail.UserMetadata["userId"]

	if trackID == "" || userID == "" {
		return &Response{
			Status: "failed",
			Reason: "missing_metadata",
		}, nil
	}

	// Handle based on job status
	switch detail.Status {
	case "COMPLETE":
		return handleSuccess(ctx, userID, trackID, detail)
	case "ERROR", "CANCELED":
		return handleFailure(ctx, userID, trackID, detail)
	default:
		// Ignore other statuses (SUBMITTED, PROGRESSING)
		return &Response{
			TrackID: trackID,
			Status:  "ignored",
			Reason:  fmt.Sprintf("status_%s", detail.Status),
		}, nil
	}
}

func handleSuccess(ctx context.Context, userID, trackID string, detail service.MediaConvertEventDetail) (*Response, error) {
	// Find the playlist path from output details
	var playlistKey string
	for _, og := range detail.OutputGroupDetails {
		if len(og.PlaylistFilePaths) > 0 {
			// Extract the S3 key from the full path
			// Format: s3://bucket/hls/userId/trackId/master.m3u8
			playlistKey = extractS3Key(og.PlaylistFilePaths[0])
			break
		}
	}

	if playlistKey == "" {
		// Fallback to constructed key
		playlistKey = service.BuildHLSPlaylistKey(userID, trackID)
	}

	// Update track in DynamoDB
	if err := updateTrackHLSStatus(ctx, userID, trackID, models.HLSStatusReady, playlistKey, ""); err != nil {
		return &Response{
			TrackID: trackID,
			Status:  "failed",
			Reason:  fmt.Sprintf("db_update_failed: %v", err),
		}, nil
	}

	return &Response{
		TrackID: trackID,
		Status:  "completed",
	}, nil
}

func handleFailure(ctx context.Context, userID, trackID string, detail service.MediaConvertEventDetail) (*Response, error) {
	errorMsg := detail.ErrorMessage
	if errorMsg == "" {
		errorMsg = fmt.Sprintf("Job failed with code %d", detail.ErrorCode)
	}

	// Update track in DynamoDB
	if err := updateTrackHLSStatus(ctx, userID, trackID, models.HLSStatusFailed, "", errorMsg); err != nil {
		return &Response{
			TrackID: trackID,
			Status:  "failed",
			Reason:  fmt.Sprintf("db_update_failed: %v", err),
		}, nil
	}

	return &Response{
		TrackID: trackID,
		Status:  "transcode_failed",
		Reason:  errorMsg,
	}, nil
}

func updateTrackHLSStatus(ctx context.Context, userID, trackID string, status models.HLSStatus, playlistKey, errorMsg string) error {
	if dynamoClient == nil || tableName == "" {
		return fmt.Errorf("DynamoDB not configured")
	}

	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("TRACK#%s", trackID)

	// Build update expression based on status
	updateExpr := "SET hlsStatus = :status, updatedAt = :now"
	exprValues := map[string]dynamodbtypes.AttributeValue{
		":status": &dynamodbtypes.AttributeValueMemberS{Value: string(status)},
		":now":    &dynamodbtypes.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	if status == models.HLSStatusReady && playlistKey != "" {
		updateExpr += ", hlsPlaylistKey = :playlist, hlsTranscodedAt = :transcodedAt"
		exprValues[":playlist"] = &dynamodbtypes.AttributeValueMemberS{Value: playlistKey}
		exprValues[":transcodedAt"] = &dynamodbtypes.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)}
	}

	if status == models.HLSStatusFailed && errorMsg != "" {
		updateExpr += ", hlsError = :error"
		exprValues[":error"] = &dynamodbtypes.AttributeValueMemberS{Value: errorMsg}
	}

	input := &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]dynamodbtypes.AttributeValue{
			"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"SK": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprValues,
	}

	_, err := dynamoClient.UpdateItem(ctx, input)
	return err
}

// extractS3Key extracts the S3 key from an S3 URI
func extractS3Key(s3URI string) string {
	// Format: s3://bucket/key
	// We want just the key portion
	if len(s3URI) < 6 {
		return ""
	}

	// Remove s3:// prefix
	withoutPrefix := s3URI[5:]

	// Find the first / after bucket name
	for i, c := range withoutPrefix {
		if c == '/' {
			return withoutPrefix[i+1:]
		}
	}

	return ""
}

func main() {
	lambda.Start(handleRequest)
}
