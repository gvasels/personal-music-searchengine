package service

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMediaConvertClient mocks MediaConvert operations
type MockMediaConvertClient struct {
	mock.Mock
}

func (m *MockMediaConvertClient) CreateJob(ctx context.Context, params *mediaconvert.CreateJobInput, optFns ...func(*mediaconvert.Options)) (*mediaconvert.CreateJobOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mediaconvert.CreateJobOutput), args.Error(1)
}

func (m *MockMediaConvertClient) GetJob(ctx context.Context, params *mediaconvert.GetJobInput, optFns ...func(*mediaconvert.Options)) (*mediaconvert.GetJobOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mediaconvert.GetJobOutput), args.Error(1)
}

func TestStartTranscode_CreatesJob(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockMediaConvertClient)

	svc := NewTranscodeService(mockClient, "my-bucket", "arn:aws:iam::123456789012:role/MediaConvertRole", "arn:aws:mediaconvert:us-east-1:123456789012:queues/Default")

	mockClient.On("CreateJob", ctx, mock.MatchedBy(func(input *mediaconvert.CreateJobInput) bool {
		return *input.Role == "arn:aws:iam::123456789012:role/MediaConvertRole" &&
			input.Tags["trackId"] == "track-123" &&
			input.Tags["userId"] == "user-456"
	})).Return(&mediaconvert.CreateJobOutput{
		Job: &types.Job{
			Id:     aws.String("job-789"),
			Status: types.JobStatusSubmitted,
		},
	}, nil)

	req := TranscodeRequest{
		TrackID: "track-123",
		UserID:  "user-456",
		S3Key:   "audio/user-456/track-123/original.mp3",
	}

	resp, err := svc.StartTranscode(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "job-789", resp.JobID)
	assert.Equal(t, "SUBMITTED", resp.Status)
	assert.Equal(t, "hls/user-456/track-123/master.m3u8", resp.PlaylistKey)

	mockClient.AssertExpectations(t)
}

func TestStartTranscode_MissingTrackID(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockMediaConvertClient)

	svc := NewTranscodeService(mockClient, "my-bucket", "role-arn", "queue-arn")

	req := TranscodeRequest{
		UserID: "user-456",
		S3Key:  "audio/file.mp3",
	}

	resp, err := svc.StartTranscode(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "trackID")
}

func TestStartTranscode_MissingS3Key(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockMediaConvertClient)

	svc := NewTranscodeService(mockClient, "my-bucket", "role-arn", "queue-arn")

	req := TranscodeRequest{
		TrackID: "track-123",
		UserID:  "user-456",
	}

	resp, err := svc.StartTranscode(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "s3Key")
}

func TestBuildJobSettings_ThreeQualities(t *testing.T) {
	mockClient := new(MockMediaConvertClient)
	svc := NewTranscodeService(mockClient, "my-bucket", "role-arn", "queue-arn")

	req := TranscodeRequest{
		TrackID: "track-123",
		UserID:  "user-456",
		S3Key:   "audio/file.mp3",
	}

	settings := svc.buildJobSettings(req)

	assert.NotNil(t, settings)
	assert.Len(t, settings.Inputs, 1)
	assert.Len(t, settings.OutputGroups, 1)

	outputs := settings.OutputGroups[0].Outputs
	assert.Len(t, outputs, 3, "Should have 3 quality levels")

	// Check bitrates
	bitrates := map[string]int32{
		"96k":  96000,
		"192k": 192000,
		"320k": 320000,
	}

	for _, output := range outputs {
		modifier := *output.NameModifier
		expectedBitrate, ok := bitrates[modifier]
		assert.True(t, ok, "Unknown name modifier: %s", modifier)

		if len(output.AudioDescriptions) > 0 {
			actualBitrate := *output.AudioDescriptions[0].CodecSettings.AacSettings.Bitrate
			assert.Equal(t, expectedBitrate, actualBitrate, "Bitrate mismatch for %s", modifier)
		}
	}
}

func TestBuildJobSettings_CorrectPaths(t *testing.T) {
	mockClient := new(MockMediaConvertClient)
	svc := NewTranscodeService(mockClient, "my-bucket", "role-arn", "queue-arn")

	req := TranscodeRequest{
		TrackID: "track-123",
		UserID:  "user-456",
		S3Key:   "audio/user-456/track-123/original.mp3",
	}

	settings := svc.buildJobSettings(req)

	// Check input path
	inputPath := *settings.Inputs[0].FileInput
	assert.Equal(t, "s3://my-bucket/audio/user-456/track-123/original.mp3", inputPath)

	// Check output path
	hlsSettings := settings.OutputGroups[0].OutputGroupSettings.HlsGroupSettings
	outputPath := *hlsSettings.Destination
	assert.Equal(t, "s3://my-bucket/hls/user-456/track-123/", outputPath)
}

func TestGetTranscodeStatus_Success(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockMediaConvertClient)

	svc := NewTranscodeService(mockClient, "my-bucket", "role-arn", "queue-arn")

	mockClient.On("GetJob", ctx, mock.MatchedBy(func(input *mediaconvert.GetJobInput) bool {
		return *input.Id == "job-789"
	})).Return(&mediaconvert.GetJobOutput{
		Job: &types.Job{
			Id:     aws.String("job-789"),
			Status: types.JobStatusComplete,
		},
	}, nil)

	status, err := svc.GetTranscodeStatus(ctx, "job-789")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "job-789", status.JobID)
	assert.Equal(t, "COMPLETE", status.Status)

	mockClient.AssertExpectations(t)
}

func TestGetTranscodeStatus_WithError(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockMediaConvertClient)

	svc := NewTranscodeService(mockClient, "my-bucket", "role-arn", "queue-arn")

	errorCode := int32(1001)
	mockClient.On("GetJob", ctx, mock.Anything).Return(&mediaconvert.GetJobOutput{
		Job: &types.Job{
			Id:           aws.String("job-789"),
			Status:       types.JobStatusError,
			ErrorCode:    &errorCode,
			ErrorMessage: aws.String("Invalid input file"),
		},
	}, nil)

	status, err := svc.GetTranscodeStatus(ctx, "job-789")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "ERROR", status.Status)
	assert.Equal(t, 1001, status.ErrorCode)
	assert.Equal(t, "Invalid input file", status.ErrorMessage)

	mockClient.AssertExpectations(t)
}

func TestGetHLSQualities(t *testing.T) {
	qualities := GetHLSQualities()

	assert.Len(t, qualities, 3)
	assert.Equal(t, "Low", qualities[0].Name)
	assert.Equal(t, 96000, qualities[0].Bitrate)
	assert.Equal(t, "Medium", qualities[1].Name)
	assert.Equal(t, 192000, qualities[1].Bitrate)
	assert.Equal(t, "High", qualities[2].Name)
	assert.Equal(t, 320000, qualities[2].Bitrate)
}

func TestBuildHLSPlaylistKey(t *testing.T) {
	key := BuildHLSPlaylistKey("user-123", "track-456")
	assert.Equal(t, "hls/user-123/track-456/master.m3u8", key)
}
