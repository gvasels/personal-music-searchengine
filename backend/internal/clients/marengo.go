package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/google/uuid"
)

// MarengoClient provides access to TwelveLabs Marengo video embedding models via Bedrock
type MarengoClient struct {
	client       *bedrockruntime.Client
	outputBucket string
}

// NewMarengoClient creates a new MarengoClient
func NewMarengoClient(client *bedrockruntime.Client) *MarengoClient {
	return &MarengoClient{client: client}
}

// SetOutputBucket sets the S3 bucket for async output
func (c *MarengoClient) SetOutputBucket(bucket string) {
	c.outputBucket = bucket
}

// MarengoModelID is the Bedrock model ID for TwelveLabs Marengo Embed 3.0
const MarengoModelID = "twelvelabs.marengo-embed-3-0-v1:0"

// VideoEmbeddingRequest represents a video embedding request
type VideoEmbeddingRequest struct {
	VideoURI    string              `json:"video_uri"`              // S3 URI or presigned URL
	EmbedType   VideoEmbedType      `json:"embed_type,omitempty"`   // visual, audio, speech, combined
	Segments    []VideoSegment      `json:"segments,omitempty"`     // Optional segments for temporal embeddings
	Options     *VideoEmbedOptions  `json:"options,omitempty"`
}

// VideoEmbedType represents the type of embedding to generate
type VideoEmbedType string

const (
	EmbedTypeVisual   VideoEmbedType = "visual"
	EmbedTypeAudio    VideoEmbedType = "audio"
	EmbedTypeSpeech   VideoEmbedType = "speech"
	EmbedTypeCombined VideoEmbedType = "combined"
)

// VideoSegment defines a time range for segment-based embedding
type VideoSegment struct {
	StartTime float64 `json:"start_time"` // Start time in seconds
	EndTime   float64 `json:"end_time"`   // End time in seconds
}

// VideoEmbedOptions contains optional parameters for video embedding
type VideoEmbedOptions struct {
	IncludeThumbnail bool `json:"include_thumbnail,omitempty"`
	MaxDuration      int  `json:"max_duration,omitempty"` // Max video duration in seconds (default 14400 = 4 hours)
}

// VideoEmbeddingResponse represents the video embedding response
type VideoEmbeddingResponse struct {
	ID                string             `json:"id"`
	Object            string             `json:"object"`
	Model             string             `json:"model"`
	Embedding         []float32          `json:"embedding,omitempty"`       // 1024-dimensional vector (for sync)
	InvocationARN     string             `json:"invocationArn,omitempty"`   // For async invocations
	SegmentEmbeddings []SegmentEmbedding `json:"segment_embeddings,omitempty"`
	Duration          float64            `json:"duration,omitempty"`        // Video duration in seconds
	CreatedAt         time.Time          `json:"created_at"`
}

// SegmentEmbedding represents an embedding for a specific video segment
type SegmentEmbedding struct {
	StartTime float64   `json:"start_time"`
	EndTime   float64   `json:"end_time"`
	Embedding []float32 `json:"embedding"`
}

// CreateVideoEmbedding generates embeddings for a video using Marengo
// NOTE: Marengo 3.0 requires StartAsyncInvoke for video/audio files.
// This is a placeholder that returns an error - async implementation needed.
func (c *MarengoClient) CreateVideoEmbedding(ctx context.Context, req VideoEmbeddingRequest) (*VideoEmbeddingResponse, error) {
	// Marengo 3.0 requires StartAsyncInvoke API for video/audio
	// which involves polling and S3 output parsing.
	// For now, return a placeholder response indicating async is needed.
	return &VideoEmbeddingResponse{
		ID:        fmt.Sprintf("emb-%s", uuid.New().String()[:8]),
		Object:    "video.embedding.pending",
		Model:     MarengoModelID,
		CreatedAt: time.Now(),
	}, fmt.Errorf("Marengo 3.0 requires async invocation - not yet implemented")
}

// OpenAI-compatible wrapper types for /v1/embeddings/video endpoint

// VideoEmbeddingOpenAIRequest is the OpenAI-compatible request format
type VideoEmbeddingOpenAIRequest struct {
	Model      string         `json:"model"`
	VideoURI   string         `json:"video_uri"`
	EmbedType  string         `json:"embed_type,omitempty"`
	Segments   []VideoSegment `json:"segments,omitempty"`
}

// VideoEmbeddingOpenAIResponse is the OpenAI-compatible response format
type VideoEmbeddingOpenAIResponse struct {
	Object string                  `json:"object"`
	Model  string                  `json:"model"`
	Data   []VideoEmbeddingData    `json:"data"`
	Usage  *VideoEmbeddingUsage    `json:"usage,omitempty"`
}

// VideoEmbeddingData represents a single video embedding in OpenAI format
type VideoEmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
	Duration  float64   `json:"duration,omitempty"`
	Segment   *VideoSegment `json:"segment,omitempty"`
}

// VideoEmbeddingUsage represents usage information for video embedding
type VideoEmbeddingUsage struct {
	DurationSeconds float64 `json:"duration_seconds"`
	SegmentCount    int     `json:"segment_count"`
}

// CreateVideoEmbeddingOpenAI creates video embeddings in OpenAI-compatible format
func (c *MarengoClient) CreateVideoEmbeddingOpenAI(ctx context.Context, req VideoEmbeddingOpenAIRequest) (*VideoEmbeddingOpenAIResponse, error) {
	// Convert to internal request format
	internalReq := VideoEmbeddingRequest{
		VideoURI:  req.VideoURI,
		EmbedType: VideoEmbedType(req.EmbedType),
		Segments:  req.Segments,
	}

	resp, err := c.CreateVideoEmbedding(ctx, internalReq)
	if err != nil {
		return nil, err
	}

	// Build OpenAI-compatible response
	data := []VideoEmbeddingData{
		{
			Object:    "video.embedding",
			Index:     0,
			Embedding: resp.Embedding,
			Duration:  resp.Duration,
		},
	}

	// Add segment embeddings
	for i, seg := range resp.SegmentEmbeddings {
		data = append(data, VideoEmbeddingData{
			Object:    "video.embedding.segment",
			Index:     i + 1,
			Embedding: seg.Embedding,
			Segment: &VideoSegment{
				StartTime: seg.StartTime,
				EndTime:   seg.EndTime,
			},
		})
	}

	return &VideoEmbeddingOpenAIResponse{
		Object: "list",
		Model:  req.Model,
		Data:   data,
		Usage: &VideoEmbeddingUsage{
			DurationSeconds: resp.Duration,
			SegmentCount:    len(resp.SegmentEmbeddings),
		},
	}, nil
}

// TextToVideoSearch generates text embedding for searching against video embeddings
// This uses Marengo's text encoder to generate embeddings compatible with video embeddings
func (c *MarengoClient) TextToVideoSearch(ctx context.Context, text string) ([]float32, error) {
	// Build Marengo text embedding request
	marengoReq := map[string]interface{}{
		"text": text,
		"mode": "text_to_video",
	}

	body, err := json.Marshal(marengoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(MarengoModelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke Marengo model: %w", err)
	}

	var resp struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Embedding, nil
}

// AudioToVideoSearch generates audio embedding for searching against video embeddings
func (c *MarengoClient) AudioToVideoSearch(ctx context.Context, audioURI string) ([]float32, error) {
	// Build Marengo audio embedding request
	marengoReq := map[string]interface{}{
		"audio_url": audioURI,
		"mode":      "audio_to_video",
	}

	body, err := json.Marshal(marengoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(MarengoModelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke Marengo model: %w", err)
	}

	var resp struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Embedding, nil
}
