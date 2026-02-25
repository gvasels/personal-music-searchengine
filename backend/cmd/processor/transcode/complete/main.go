package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
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

// Interfaces for testability
type dynamoDBAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type sfnAPI interface {
	StartExecution(ctx context.Context, params *sfn.StartExecutionInput, optFns ...func(*sfn.Options)) (*sfn.StartExecutionOutput, error)
}

var (
	dynamoClient     dynamoDBAPI
	sfnClient        sfnAPI
	tableName        string
	audioPipelineARN string
)

func init() {
	tableName = os.Getenv("DYNAMODB_TABLE_NAME")
	audioPipelineARN = os.Getenv("AUDIO_PIPELINE_ARN")

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	sfnClient = sfn.NewFromConfig(cfg)
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
		fmt.Printf("Missing metadata in event: trackId=%q, userId=%q, jobId=%q\n", trackID, userID, detail.JobID)
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

	// Start audio analysis pipeline
	if audioPipelineARN != "" {
		if err := startAudioPipeline(ctx, userID, trackID); err != nil {
			fmt.Printf("Warning: failed to start audio pipeline for track %s: %v\n", trackID, err)
			// Don't fail the response — HLS is ready, analysis is best-effort
		}
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

// readTrack fetches a track from DynamoDB and returns its s3Key, title, and artist.
func readTrack(ctx context.Context, db dynamoDBAPI, table, userID, trackID string) (s3Key, title, artist string, err error) {
	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("TRACK#%s", trackID)

	result, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &table,
		Key: map[string]dynamodbtypes.AttributeValue{
			"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"SK": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get track: %w", err)
	}
	if result.Item == nil {
		return "", "", "", fmt.Errorf("track not found: %s/%s", userID, trackID)
	}

	getStr := func(key string) string {
		if v, ok := result.Item[key]; ok {
			if sv, ok := v.(*dynamodbtypes.AttributeValueMemberS); ok {
				return sv.Value
			}
		}
		return ""
	}

	return getStr("s3Key"), getStr("title"), getStr("artist"), nil
}

// triggerAudioPipeline starts the audio analysis Step Functions execution.
// It is a no-op if the ARN is empty or the client is nil.
func triggerAudioPipeline(ctx context.Context, client sfnAPI, arn, trackID, userID, s3Key, title, artist string) {
	if arn == "" || client == nil {
		return
	}

	pipelineInput := map[string]string{
		"trackId": trackID,
		"userId":  userID,
		"s3Key":   s3Key,
		"title":   title,
		"artist":  artist,
	}

	inputBytes, err := json.Marshal(pipelineInput)
	if err != nil {
		fmt.Printf("Warning: failed to marshal pipeline input: %v\n", err)
		return
	}

	executionName := fmt.Sprintf("transcode-%s-%d", trackID, time.Now().Unix())
	_, err = client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: &arn,
		Name:            &executionName,
		Input:           aws.String(string(inputBytes)),
	})
	if err != nil {
		fmt.Printf("Warning: failed to start audio pipeline: %v\n", err)
	} else {
		fmt.Printf("Started audio pipeline for track %s: execution=%s\n", trackID, executionName)
	}
}

// startAudioPipeline reads track metadata and triggers the audio analysis pipeline.
func startAudioPipeline(ctx context.Context, userID, trackID string) error {
	if sfnClient == nil {
		return fmt.Errorf("Step Functions client not configured")
	}

	s3Key, title, artist, err := readTrack(ctx, dynamoClient, tableName, userID, trackID)
	if err != nil {
		return err
	}

	triggerAudioPipeline(ctx, sfnClient, audioPipelineARN, trackID, userID, s3Key, title, artist)
	return nil
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
