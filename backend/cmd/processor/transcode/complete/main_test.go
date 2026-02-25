package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock implementations ---
// Interfaces (dynamoDBAPI, sfnAPI) are defined in main.go.

type mockDynamoDB struct {
	getItemOutput *dynamodb.GetItemOutput
	getItemErr    error
	updateItemOut *dynamodb.UpdateItemOutput
	updateItemErr error

	getItemCalls    []dynamodb.GetItemInput
	updateItemCalls []dynamodb.UpdateItemInput
}

func (m *mockDynamoDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	m.getItemCalls = append(m.getItemCalls, *params)
	return m.getItemOutput, m.getItemErr
}

func (m *mockDynamoDB) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	m.updateItemCalls = append(m.updateItemCalls, *params)
	return m.updateItemOut, m.updateItemErr
}

type mockSFN struct {
	startExecOutput *sfn.StartExecutionOutput
	startExecErr    error

	startExecCalls []sfn.StartExecutionInput
}

func (m *mockSFN) StartExecution(ctx context.Context, params *sfn.StartExecutionInput, optFns ...func(*sfn.Options)) (*sfn.StartExecutionOutput, error) {
	m.startExecCalls = append(m.startExecCalls, *params)
	return m.startExecOutput, m.startExecErr
}

// --- Helper to build a MediaConvert COMPLETE event ---

func makeCompleteEvent(userID, trackID string) Event {
	return Event{
		Detail: service.MediaConvertEventDetail{
			Status: "COMPLETE",
			UserMetadata: map[string]string{
				"trackId": trackID,
				"userId":  userID,
			},
			OutputGroupDetails: []service.OutputGroupDetail{
				{
					PlaylistFilePaths: []string{
						"s3://test-bucket/hls/" + userID + "/" + trackID + "/master.m3u8",
					},
				},
			},
		},
	}
}

func makeErrorEvent(userID, trackID string) Event {
	return Event{
		Detail: service.MediaConvertEventDetail{
			Status:       "ERROR",
			ErrorCode:    1234,
			ErrorMessage: "Transcode failed",
			UserMetadata: map[string]string{
				"trackId": trackID,
				"userId":  userID,
			},
		},
	}
}

// --- Tests ---

// TestReadTrack_ExtractsFields tests that readTrack fetches a track from DynamoDB
// and returns the s3Key, title, and artist fields.
// This test MUST FAIL because readTrack does not exist yet.
func TestReadTrack_ExtractsFields(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: map[string]dynamodbtypes.AttributeValue{
				"PK":     &dynamodbtypes.AttributeValueMemberS{Value: "USER#user1"},
				"SK":     &dynamodbtypes.AttributeValueMemberS{Value: "TRACK#track1"},
				"s3Key":  &dynamodbtypes.AttributeValueMemberS{Value: "media/user1/track1.mp3"},
				"title":  &dynamodbtypes.AttributeValueMemberS{Value: "Test Song"},
				"artist": &dynamodbtypes.AttributeValueMemberS{Value: "Test Artist"},
			},
		},
	}

	// readTrack does not exist yet -- this will cause a compile error (Red phase).
	s3Key, title, artist, err := readTrack(ctx, mockDB, "test-table", "user1", "track1")
	require.NoError(t, err)
	assert.Equal(t, "media/user1/track1.mp3", s3Key)
	assert.Equal(t, "Test Song", title)
	assert.Equal(t, "Test Artist", artist)

	// Verify correct DynamoDB key was used
	require.Len(t, mockDB.getItemCalls, 1)
	call := mockDB.getItemCalls[0]
	assert.Equal(t, "test-table", *call.TableName)
	pk := call.Key["PK"].(*dynamodbtypes.AttributeValueMemberS).Value
	sk := call.Key["SK"].(*dynamodbtypes.AttributeValueMemberS).Value
	assert.Equal(t, "USER#user1", pk)
	assert.Equal(t, "TRACK#track1", sk)
}

// TestReadTrack_MissingItem tests that readTrack returns an error when the track
// is not found in DynamoDB (empty item).
func TestReadTrack_MissingItem(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: nil, // No item found
		},
	}

	_, _, _, err := readTrack(ctx, mockDB, "test-table", "user1", "nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestReadTrack_DynamoDBError tests that readTrack propagates DynamoDB errors.
func TestReadTrack_DynamoDBError(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		getItemErr: assert.AnError,
	}

	_, _, _, err := readTrack(ctx, mockDB, "test-table", "user1", "track1")
	require.Error(t, err)
}

// TestTriggerAudioPipeline_CallsSFN tests that triggerAudioPipeline calls
// SFN.StartExecution with the correct input payload.
// This test MUST FAIL because triggerAudioPipeline does not exist in this package.
func TestTriggerAudioPipeline_CallsSFN(t *testing.T) {
	ctx := context.Background()

	mockSfn := &mockSFN{
		startExecOutput: &sfn.StartExecutionOutput{},
	}

	testARN := "arn:aws:states:us-east-1:123456789:stateMachine:audio-pipeline"

	// triggerAudioPipeline does not exist yet in this package -- compile error (Red phase).
	triggerAudioPipeline(ctx, mockSfn, testARN, "track1", "user1", "media/user1/track1.mp3", "Test Song", "Test Artist")

	// Verify StartExecution was called
	require.Len(t, mockSfn.startExecCalls, 1)
	call := mockSfn.startExecCalls[0]

	// Verify the state machine ARN
	assert.Equal(t, testARN, *call.StateMachineArn)

	// Verify the input payload contains expected fields
	var inputPayload map[string]string
	err := json.Unmarshal([]byte(*call.Input), &inputPayload)
	require.NoError(t, err)
	assert.Equal(t, "track1", inputPayload["trackId"])
	assert.Equal(t, "user1", inputPayload["userId"])
	assert.Equal(t, "media/user1/track1.mp3", inputPayload["s3Key"])
	assert.Equal(t, "Test Song", inputPayload["title"])
	assert.Equal(t, "Test Artist", inputPayload["artist"])

	// Verify execution name contains trackId
	assert.Contains(t, *call.Name, "track1")
}

// TestTriggerAudioPipeline_EmptyARN tests that triggerAudioPipeline does nothing
// when the ARN is empty.
func TestTriggerAudioPipeline_EmptyARN(t *testing.T) {
	ctx := context.Background()

	mockSfn := &mockSFN{
		startExecOutput: &sfn.StartExecutionOutput{},
	}

	// With empty ARN, should not call SFN
	triggerAudioPipeline(ctx, mockSfn, "", "track1", "user1", "media/user1/track1.mp3", "Test Song", "Test Artist")

	assert.Empty(t, mockSfn.startExecCalls, "SFN should not be called when ARN is empty")
}

// TestTriggerAudioPipeline_NilClient tests that triggerAudioPipeline does nothing
// when the SFN client is nil.
func TestTriggerAudioPipeline_NilClient(t *testing.T) {
	// Should not panic with nil client
	ctx := context.Background()
	testARN := "arn:aws:states:us-east-1:123456789:stateMachine:audio-pipeline"

	// This must not panic
	triggerAudioPipeline(ctx, nil, testARN, "track1", "user1", "media/user1/track1.mp3", "Test Song", "Test Artist")
}

// TestHandleSuccess_TriggersAudioPipeline tests that after a successful HLS transcode,
// the handler reads the track from DynamoDB and triggers the audio analysis pipeline.
// This test MUST FAIL because handleSuccess currently does not trigger the audio pipeline.
func TestHandleSuccess_TriggersAudioPipeline(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		// First call: GetItem for readTrack
		getItemOutput: &dynamodb.GetItemOutput{
			Item: map[string]dynamodbtypes.AttributeValue{
				"PK":     &dynamodbtypes.AttributeValueMemberS{Value: "USER#user1"},
				"SK":     &dynamodbtypes.AttributeValueMemberS{Value: "TRACK#track1"},
				"s3Key":  &dynamodbtypes.AttributeValueMemberS{Value: "media/user1/track1.mp3"},
				"title":  &dynamodbtypes.AttributeValueMemberS{Value: "Test Song"},
				"artist": &dynamodbtypes.AttributeValueMemberS{Value: "Test Artist"},
			},
		},
		// UpdateItem for HLS status update
		updateItemOut: &dynamodb.UpdateItemOutput{},
	}

	mockSfn := &mockSFN{
		startExecOutput: &sfn.StartExecutionOutput{},
	}

	// Save and restore package-level state
	origDynamoClient := dynamoClient
	origSfnClient := sfnClient
	origTableName := tableName
	origARN := audioPipelineARN
	defer func() {
		dynamoClient = origDynamoClient
		sfnClient = origSfnClient
		tableName = origTableName
		audioPipelineARN = origARN
	}()

	// Set package-level vars
	dynamoClient = mockDB
	sfnClient = mockSfn
	tableName = "test-table"
	audioPipelineARN = "arn:aws:states:us-east-1:123456789:stateMachine:audio-pipeline"

	event := makeCompleteEvent("user1", "track1")
	resp, err := handleSuccess(ctx, "user1", "track1", event.Detail)

	require.NoError(t, err)
	assert.Equal(t, "completed", resp.Status)
	assert.Equal(t, "track1", resp.TrackID)

	// Verify SFN was called to trigger the audio pipeline
	require.Len(t, mockSfn.startExecCalls, 1, "Audio pipeline should be triggered after HLS success")

	call := mockSfn.startExecCalls[0]
	assert.Equal(t, "arn:aws:states:us-east-1:123456789:stateMachine:audio-pipeline", *call.StateMachineArn)

	// Verify the payload
	var inputPayload map[string]string
	err = json.Unmarshal([]byte(*call.Input), &inputPayload)
	require.NoError(t, err)
	assert.Equal(t, "track1", inputPayload["trackId"])
	assert.Equal(t, "user1", inputPayload["userId"])
	assert.Equal(t, "media/user1/track1.mp3", inputPayload["s3Key"])
	assert.Equal(t, "Test Song", inputPayload["title"])
	assert.Equal(t, "Test Artist", inputPayload["artist"])
}

// TestHandleSuccess_NoTriggerWithoutARN tests that handleSuccess does NOT trigger
// the audio pipeline when audioPipelineARN is empty.
func TestHandleSuccess_NoTriggerWithoutARN(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		updateItemOut: &dynamodb.UpdateItemOutput{},
	}

	mockSfn := &mockSFN{
		startExecOutput: &sfn.StartExecutionOutput{},
	}

	origDynamoClient := dynamoClient
	origSfnClient := sfnClient
	origTableName := tableName
	origARN := audioPipelineARN
	defer func() {
		dynamoClient = origDynamoClient
		sfnClient = origSfnClient
		tableName = origTableName
		audioPipelineARN = origARN
	}()

	dynamoClient = mockDB
	sfnClient = mockSfn
	tableName = "test-table"
	audioPipelineARN = "" // Empty ARN -- should not trigger

	event := makeCompleteEvent("user1", "track1")
	resp, err := handleSuccess(ctx, "user1", "track1", event.Detail)

	require.NoError(t, err)
	assert.Equal(t, "completed", resp.Status)

	assert.Empty(t, mockSfn.startExecCalls, "SFN should NOT be called when audioPipelineARN is empty")
}

// TestHandleFailure_NoTrigger tests that handleFailure does NOT trigger the audio pipeline.
func TestHandleFailure_NoTrigger(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		updateItemOut: &dynamodb.UpdateItemOutput{},
	}
	mockSfn := &mockSFN{
		startExecOutput: &sfn.StartExecutionOutput{},
	}

	origDynamoClient := dynamoClient
	origSfnClient := sfnClient
	origTableName := tableName
	origARN := audioPipelineARN
	defer func() {
		dynamoClient = origDynamoClient
		sfnClient = origSfnClient
		tableName = origTableName
		audioPipelineARN = origARN
	}()

	dynamoClient = mockDB
	sfnClient = mockSfn
	tableName = "test-table"
	audioPipelineARN = "arn:aws:states:us-east-1:123456789:stateMachine:audio-pipeline"

	event := makeErrorEvent("user1", "track1")
	resp, err := handleFailure(ctx, "user1", "track1", event.Detail)

	require.NoError(t, err)
	assert.Equal(t, "transcode_failed", resp.Status)

	assert.Empty(t, mockSfn.startExecCalls, "SFN should NOT be called on transcode failure")
}

// TestHandleSuccess_SFNError_DoesNotFailHandler tests that if SFN.StartExecution
// fails, the handleSuccess function still returns success (audio pipeline trigger
// is best-effort, not critical).
func TestHandleSuccess_SFNError_DoesNotFailHandler(t *testing.T) {
	ctx := context.Background()

	mockDB := &mockDynamoDB{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: map[string]dynamodbtypes.AttributeValue{
				"PK":     &dynamodbtypes.AttributeValueMemberS{Value: "USER#user1"},
				"SK":     &dynamodbtypes.AttributeValueMemberS{Value: "TRACK#track1"},
				"s3Key":  &dynamodbtypes.AttributeValueMemberS{Value: "media/user1/track1.mp3"},
				"title":  &dynamodbtypes.AttributeValueMemberS{Value: "Test Song"},
				"artist": &dynamodbtypes.AttributeValueMemberS{Value: "Test Artist"},
			},
		},
		updateItemOut: &dynamodb.UpdateItemOutput{},
	}

	mockSfn := &mockSFN{
		startExecErr: assert.AnError, // SFN will return an error
	}

	origDynamoClient := dynamoClient
	origSfnClient := sfnClient
	origTableName := tableName
	origARN := audioPipelineARN
	defer func() {
		dynamoClient = origDynamoClient
		sfnClient = origSfnClient
		tableName = origTableName
		audioPipelineARN = origARN
	}()

	dynamoClient = mockDB
	sfnClient = mockSfn
	tableName = "test-table"
	audioPipelineARN = "arn:aws:states:us-east-1:123456789:stateMachine:audio-pipeline"

	event := makeCompleteEvent("user1", "track1")
	resp, err := handleSuccess(ctx, "user1", "track1", event.Detail)

	// Handler should still succeed even if SFN fails
	require.NoError(t, err)
	assert.Equal(t, "completed", resp.Status)
}
