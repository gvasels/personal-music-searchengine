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

// BedrockClient provides access to Amazon Bedrock models
type BedrockClient struct {
	client *bedrockruntime.Client
}

// NewBedrockClient creates a new BedrockClient
func NewBedrockClient(client *bedrockruntime.Client) *BedrockClient {
	return &BedrockClient{client: client}
}

// OpenAI-compatible request/response types

// ChatCompletionRequest represents an OpenAI-compatible chat completion request
type ChatCompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []ChatMessage   `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
	Stop        []string        `json:"stop,omitempty"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents an OpenAI-compatible chat completion response
type ChatCompletionResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []ChatChoice   `json:"choices"`
	Usage   *UsageInfo     `json:"usage,omitempty"`
}

// ChatChoice represents a choice in the response
type ChatChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason string       `json:"finish_reason"`
}

// UsageInfo represents token usage information
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// EmbeddingRequest represents an OpenAI-compatible embedding request
type EmbeddingRequest struct {
	Model string      `json:"model"`
	Input interface{} `json:"input"` // string or []string
}

// EmbeddingResponse represents an OpenAI-compatible embedding response
type EmbeddingResponse struct {
	Object string          `json:"object"`
	Model  string          `json:"model"`
	Data   []EmbeddingData `json:"data"`
	Usage  *UsageInfo      `json:"usage,omitempty"`
}

// EmbeddingData represents a single embedding
type EmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

// ModelInfo represents model information
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsResponse represents the list of available models
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// Model mapping from OpenAI to Bedrock
var modelMapping = map[string]string{
	"gpt-4":                     "anthropic.claude-3-5-sonnet-20241022-v2:0",
	"gpt-4-turbo":               "anthropic.claude-3-5-sonnet-20241022-v2:0",
	"gpt-4o":                    "anthropic.claude-3-5-sonnet-20241022-v2:0",
	"gpt-3.5-turbo":             "anthropic.claude-3-haiku-20240307-v1:0",
	"claude-3-sonnet":           "anthropic.claude-3-5-sonnet-20241022-v2:0",
	"claude-3-haiku":            "anthropic.claude-3-haiku-20240307-v1:0",
	"claude-3-opus":             "anthropic.claude-3-opus-20240229-v1:0",
	"text-embedding-ada-002":    "amazon.titan-embed-text-v2:0",
	"text-embedding-3-small":    "amazon.titan-embed-text-v2:0",
	"text-embedding-3-large":    "amazon.titan-embed-text-v2:0",
}

// mapModel translates OpenAI model names to Bedrock model IDs
func mapModel(model string) string {
	if bedrockModel, ok := modelMapping[model]; ok {
		return bedrockModel
	}
	// If already a Bedrock model ID, use as-is
	return model
}

// CreateChatCompletion creates a chat completion using Bedrock
func (c *BedrockClient) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	bedrockModel := mapModel(req.Model)

	// Build Claude request
	claudeReq := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        req.MaxTokens,
		"messages":          req.Messages,
	}

	if req.MaxTokens == 0 {
		claudeReq["max_tokens"] = 4096
	}
	if req.Temperature > 0 {
		claudeReq["temperature"] = req.Temperature
	}
	if len(req.Stop) > 0 {
		claudeReq["stop_sequences"] = req.Stop
	}

	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bedrockModel),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse Claude response
	var claudeResp struct {
		ID           string `json:"id"`
		Type         string `json:"type"`
		Role         string `json:"role"`
		Content      []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Model        string `json:"model"`
		StopReason   string `json:"stop_reason"`
		Usage        struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(output.Body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Build content from response
	var content string
	for _, c := range claudeResp.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	// Convert to OpenAI format
	return &ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%s", uuid.New().String()[:8]),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChatChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: mapStopReason(claudeResp.StopReason),
			},
		},
		Usage: &UsageInfo{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
	}, nil
}

// CreateEmbedding creates text embeddings using Bedrock Titan
func (c *BedrockClient) CreateEmbedding(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	bedrockModel := mapModel(req.Model)

	// Handle single string or array of strings
	var inputs []string
	switch v := req.Input.(type) {
	case string:
		inputs = []string{v}
	case []string:
		inputs = v
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				inputs = append(inputs, s)
			}
		}
	default:
		return nil, fmt.Errorf("invalid input type: expected string or []string")
	}

	embeddings := make([]EmbeddingData, 0, len(inputs))
	var totalTokens int

	for i, input := range inputs {
		// Build Titan request
		titanReq := map[string]interface{}{
			"inputText": input,
		}
		body, err := json.Marshal(titanReq)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		output, err := c.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(bedrockModel),
			ContentType: aws.String("application/json"),
			Accept:      aws.String("application/json"),
			Body:        body,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to invoke model: %w", err)
		}

		// Parse Titan response
		var titanResp struct {
			Embedding      []float32 `json:"embedding"`
			InputTextTokenCount int  `json:"inputTextTokenCount"`
		}
		if err := json.Unmarshal(output.Body, &titanResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		embeddings = append(embeddings, EmbeddingData{
			Object:    "embedding",
			Index:     i,
			Embedding: titanResp.Embedding,
		})
		totalTokens += titanResp.InputTextTokenCount
	}

	return &EmbeddingResponse{
		Object: "list",
		Model:  req.Model,
		Data:   embeddings,
		Usage: &UsageInfo{
			PromptTokens: totalTokens,
			TotalTokens:  totalTokens,
		},
	}, nil
}

// ListModels returns available models
func (c *BedrockClient) ListModels(ctx context.Context) (*ModelsResponse, error) {
	models := []ModelInfo{
		{ID: "gpt-4", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "gpt-4-turbo", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "gpt-4o", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "gpt-3.5-turbo", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "claude-3-sonnet", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "claude-3-haiku", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "claude-3-opus", Object: "model", Created: time.Now().Unix(), OwnedBy: "anthropic"},
		{ID: "text-embedding-ada-002", Object: "model", Created: time.Now().Unix(), OwnedBy: "amazon"},
		{ID: "text-embedding-3-small", Object: "model", Created: time.Now().Unix(), OwnedBy: "amazon"},
		{ID: "text-embedding-3-large", Object: "model", Created: time.Now().Unix(), OwnedBy: "amazon"},
	}

	return &ModelsResponse{
		Object: "list",
		Data:   models,
	}, nil
}

// StreamChatCompletion creates a streaming chat completion
// Note: For simplicity, this falls back to non-streaming and sends the complete response
// A full streaming implementation would use InvokeModelWithResponseStream with proper event handling
func (c *BedrockClient) StreamChatCompletion(ctx context.Context, req ChatCompletionRequest) (<-chan StreamChunk, error) {
	chunks := make(chan StreamChunk, 100)

	go func() {
		defer close(chunks)

		// For now, use non-streaming and simulate streaming response
		resp, err := c.CreateChatCompletion(ctx, req)
		if err != nil {
			chunks <- StreamChunk{Error: err}
			return
		}

		// Send the content in chunks (simulate streaming)
		if len(resp.Choices) > 0 {
			content := resp.Choices[0].Message.Content
			// Send content in chunks of ~50 chars to simulate streaming
			chunkSize := 50
			for i := 0; i < len(content); i += chunkSize {
				end := i + chunkSize
				if end > len(content) {
					end = len(content)
				}
				chunks <- StreamChunk{
					ID:      resp.ID,
					Object:  "chat.completion.chunk",
					Created: resp.Created,
					Model:   resp.Model,
					Delta: DeltaContent{
						Content: content[i:end],
					},
				}
			}
		}

		// Send done marker
		chunks <- StreamChunk{Done: true}
	}()

	return chunks, nil
}

// StreamChunk represents a streaming response chunk
type StreamChunk struct {
	ID      string       `json:"id,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int64        `json:"created,omitempty"`
	Model   string       `json:"model,omitempty"`
	Delta   DeltaContent `json:"choices,omitempty"`
	Done    bool         `json:"-"`
	Error   error        `json:"-"`
}

// DeltaContent represents incremental content
type DeltaContent struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

// mapStopReason maps Claude stop reasons to OpenAI format
func mapStopReason(reason string) string {
	switch reason {
	case "end_turn":
		return "stop"
	case "max_tokens":
		return "length"
	case "stop_sequence":
		return "stop"
	default:
		return "stop"
	}
}
