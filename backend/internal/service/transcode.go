package service

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
)

// MediaConvertClient defines the interface for MediaConvert operations.
type MediaConvertClient interface {
	CreateJob(ctx context.Context, params *mediaconvert.CreateJobInput, optFns ...func(*mediaconvert.Options)) (*mediaconvert.CreateJobOutput, error)
	GetJob(ctx context.Context, params *mediaconvert.GetJobInput, optFns ...func(*mediaconvert.Options)) (*mediaconvert.GetJobOutput, error)
}

// TranscodeService provides HLS transcoding operations.
type TranscodeService struct {
	mcClient     MediaConvertClient
	bucket       string
	role         string
	queue        string
	outputPrefix string
}

// NewTranscodeService creates a new transcode service.
func NewTranscodeService(mcClient MediaConvertClient, bucket, role, queue string) *TranscodeService {
	return &TranscodeService{
		mcClient:     mcClient,
		bucket:       bucket,
		role:         role,
		queue:        queue,
		outputPrefix: "hls",
	}
}

// TranscodeRequest represents a request to transcode a track.
type TranscodeRequest struct {
	TrackID string
	UserID  string
	S3Key   string // Source audio file key
}

// TranscodeResponse represents the response from starting a transcode job.
type TranscodeResponse struct {
	JobID       string
	Status      string
	PlaylistKey string // S3 key where master.m3u8 will be created
}

// StartTranscode creates a MediaConvert job to transcode audio to HLS.
func (s *TranscodeService) StartTranscode(ctx context.Context, req TranscodeRequest) (*TranscodeResponse, error) {
	if req.TrackID == "" || req.UserID == "" || req.S3Key == "" {
		return nil, fmt.Errorf("trackID, userID, and s3Key are required")
	}

	// Build job settings
	jobSettings := s.buildJobSettings(req)

	input := &mediaconvert.CreateJobInput{
		Role:     aws.String(s.role),
		Queue:    aws.String(s.queue),
		Settings: jobSettings,
		Tags: map[string]string{
			"trackId": req.TrackID,
			"userId":  req.UserID,
		},
	}

	output, err := s.mcClient.CreateJob(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create MediaConvert job: %w", err)
	}

	playlistKey := fmt.Sprintf("%s/%s/%s/master.m3u8", s.outputPrefix, req.UserID, req.TrackID)

	return &TranscodeResponse{
		JobID:       *output.Job.Id,
		Status:      string(output.Job.Status),
		PlaylistKey: playlistKey,
	}, nil
}

// GetTranscodeStatus retrieves the status of a MediaConvert job.
func (s *TranscodeService) GetTranscodeStatus(ctx context.Context, jobID string) (*TranscodeJobStatus, error) {
	input := &mediaconvert.GetJobInput{
		Id: aws.String(jobID),
	}

	output, err := s.mcClient.GetJob(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get MediaConvert job: %w", err)
	}

	job := output.Job
	status := &TranscodeJobStatus{
		JobID:  *job.Id,
		Status: string(job.Status),
	}

	if job.ErrorCode != nil {
		status.ErrorCode = int(*job.ErrorCode)
	}
	if job.ErrorMessage != nil {
		status.ErrorMessage = *job.ErrorMessage
	}

	return status, nil
}

// TranscodeJobStatus represents the status of a transcode job.
type TranscodeJobStatus struct {
	JobID        string
	Status       string // SUBMITTED, PROGRESSING, COMPLETE, CANCELED, ERROR
	ErrorCode    int
	ErrorMessage string
}

// buildJobSettings creates MediaConvert job settings for HLS output.
func (s *TranscodeService) buildJobSettings(req TranscodeRequest) *types.JobSettings {
	inputS3URI := fmt.Sprintf("s3://%s/%s", s.bucket, req.S3Key)
	outputS3Path := fmt.Sprintf("s3://%s/%s/%s/%s/", s.bucket, s.outputPrefix, req.UserID, req.TrackID)

	return &types.JobSettings{
		Inputs: []types.Input{
			{
				FileInput: aws.String(inputS3URI),
				AudioSelectors: map[string]types.AudioSelector{
					"Audio Selector 1": {
						DefaultSelection: types.AudioDefaultSelectionDefault,
					},
				},
			},
		},
		OutputGroups: []types.OutputGroup{
			{
				Name: aws.String("HLS Group"),
				OutputGroupSettings: &types.OutputGroupSettings{
					Type: types.OutputGroupTypeHlsGroupSettings,
					HlsGroupSettings: &types.HlsGroupSettings{
						Destination:         aws.String(outputS3Path),
						SegmentLength:       aws.Int32(6),
						MinSegmentLength:    aws.Int32(0),
						OutputSelection:     types.HlsOutputSelectionManifestsAndSegments,
						SegmentControl:      types.HlsSegmentControlSegmentedFiles,
						ManifestDurationFormat: types.HlsManifestDurationFormatFloatingPoint,
					},
				},
				Outputs: []types.Output{
					// 96 kbps AAC (low quality for poor connections)
					s.buildAACOutput("96k", 96000),
					// 192 kbps AAC (medium quality)
					s.buildAACOutput("192k", 192000),
					// 320 kbps AAC (high quality)
					s.buildAACOutput("320k", 320000),
				},
			},
		},
	}
}

// buildAACOutput creates an HLS output configuration for a specific bitrate.
func (s *TranscodeService) buildAACOutput(nameModifier string, bitrate int32) types.Output {
	return types.Output{
		NameModifier: aws.String(nameModifier),
		ContainerSettings: &types.ContainerSettings{
			Container: types.ContainerTypeM3u8,
			M3u8Settings: &types.M3u8Settings{
				AudioFramesPerPes: aws.Int32(4),
				PcrControl:        types.M3u8PcrControlPcrEveryPesPacket,
			},
		},
		AudioDescriptions: []types.AudioDescription{
			{
				AudioSourceName: aws.String("Audio Selector 1"),
				CodecSettings: &types.AudioCodecSettings{
					Codec: types.AudioCodecAac,
					AacSettings: &types.AacSettings{
						Bitrate:         aws.Int32(bitrate),
						CodingMode:      types.AacCodingModeCodingMode20,
						SampleRate:      aws.Int32(48000),
						RateControlMode: types.AacRateControlModeCbr,
					},
				},
			},
		},
	}
}

// HLSQualityLevel represents an HLS quality level.
type HLSQualityLevel struct {
	Name    string
	Bitrate int
}

// GetHLSQualities returns the available HLS quality levels.
func GetHLSQualities() []HLSQualityLevel {
	return []HLSQualityLevel{
		{Name: "Low", Bitrate: 96000},
		{Name: "Medium", Bitrate: 192000},
		{Name: "High", Bitrate: 320000},
	}
}

// BuildHLSPlaylistKey builds the S3 key for the master HLS playlist.
func BuildHLSPlaylistKey(userID, trackID string) string {
	return path.Join("hls", userID, trackID, "master.m3u8")
}

// ParseMediaConvertEvent parses a MediaConvert EventBridge event.
type MediaConvertEvent struct {
	Version    string                 `json:"version"`
	ID         string                 `json:"id"`
	DetailType string                 `json:"detail-type"`
	Source     string                 `json:"source"`
	Account    string                 `json:"account"`
	Time       time.Time              `json:"time"`
	Region     string                 `json:"region"`
	Detail     MediaConvertEventDetail `json:"detail"`
}

// MediaConvertEventDetail contains the detail of a MediaConvert event.
type MediaConvertEventDetail struct {
	Timestamp   int64             `json:"timestamp"`
	AccountID   string            `json:"accountId"`
	Queue       string            `json:"queue"`
	JobID       string            `json:"jobId"`
	Status      string            `json:"status"`
	ErrorCode   int               `json:"errorCode,omitempty"`
	ErrorMessage string           `json:"errorMessage,omitempty"`
	UserMetadata map[string]string `json:"userMetadata,omitempty"`
	OutputGroupDetails []OutputGroupDetail `json:"outputGroupDetails,omitempty"`
}

// OutputGroupDetail contains details about an output group.
type OutputGroupDetail struct {
	OutputDetails []OutputDetail `json:"outputDetails,omitempty"`
	PlaylistFilePaths []string `json:"playlistFilePaths,omitempty"`
}

// OutputDetail contains details about an output.
type OutputDetail struct {
	OutputFilePaths []string `json:"outputFilePaths,omitempty"`
	DurationInMs    int64    `json:"durationInMs,omitempty"`
}
