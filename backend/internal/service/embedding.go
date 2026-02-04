package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/clients"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

const (
	// MaxEmbedTextLength is the maximum length for embedding text (Titan model limit).
	MaxEmbedTextLength = 8000
	// EmbeddingModelID is the model identifier for embeddings.
	EmbeddingModelID = "text-embedding-3-small"
)

// BedrockEmbeddingClient defines the interface for embedding generation.
// This allows for easy mocking in tests.
type BedrockEmbeddingClient interface {
	CreateEmbedding(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error)
}

// EmbeddingService handles generating embeddings for tracks and queries
// using AWS Bedrock Titan text embeddings model.
type EmbeddingService struct {
	client BedrockEmbeddingClient
}

// NewEmbeddingService creates a new EmbeddingService with the given Bedrock client.
// Panics if client is nil.
func NewEmbeddingService(client BedrockEmbeddingClient) *EmbeddingService {
	if client == nil {
		panic("bedrock client cannot be nil")
	}
	return &EmbeddingService{
		client: client,
	}
}

// ComposeEmbedText creates a text representation of track metadata for embedding generation.
// The text includes title, artist, album, genre, tags, BPM, and key.
// Truncates to 8000 characters max (Titan model limit).
func (s *EmbeddingService) ComposeEmbedText(track models.Track) string {
	var parts []string

	// Add title if present
	if track.Title != "" {
		parts = append(parts, track.Title)
	}

	// Add artist if present
	if track.Artist != "" {
		parts = append(parts, track.Artist)
	}

	// Add album if present
	if track.Album != "" {
		parts = append(parts, track.Album)
	}

	// Add genre if present
	if track.Genre != "" {
		parts = append(parts, track.Genre)
	}

	// Add tags if present (comma-separated)
	if len(track.Tags) > 0 {
		parts = append(parts, strings.Join(track.Tags, ", "))
	}

	// Add BPM if present (non-zero, labeled)
	if track.BPM > 0 {
		parts = append(parts, fmt.Sprintf("BPM: %d", track.BPM))
	}

	// Add KeyCamelot if present (labeled)
	if track.KeyCamelot != "" {
		parts = append(parts, fmt.Sprintf("Key: %s", track.KeyCamelot))
	}

	// Join all parts with space
	result := strings.Join(parts, " ")

	// Truncate if necessary
	if len(result) > MaxEmbedTextLength {
		result = result[:MaxEmbedTextLength]
	}

	return result
}

// GenerateTrackEmbedding generates a 1024-dimensional embedding vector for a track
// using AWS Bedrock Titan Embeddings model.
// Returns error if Bedrock call fails.
func (s *EmbeddingService) GenerateTrackEmbedding(ctx context.Context, track models.Track) ([]float32, error) {
	// Compose the text to embed
	text := s.ComposeEmbedText(track)

	// Call the Bedrock client
	req := clients.EmbeddingRequest{
		Model: EmbeddingModelID,
		Input: text,
	}

	resp, err := s.client.CreateEmbedding(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate track embedding: %w", err)
	}

	// Extract embedding from response
	if len(resp.Data) == 0 || len(resp.Data[0].Embedding) == 0 {
		return nil, errors.New("empty embedding response")
	}

	return resp.Data[0].Embedding, nil
}

// GenerateQueryEmbedding generates a 1024-dimensional embedding vector for a search query.
// Returns error if query is empty/whitespace or Bedrock call fails.
func (s *EmbeddingService) GenerateQueryEmbedding(ctx context.Context, query string) ([]float32, error) {
	// Validate query - reject empty or whitespace-only queries
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, errors.New("query cannot be empty")
	}

	// Truncate if necessary
	if len(trimmed) > MaxEmbedTextLength {
		trimmed = trimmed[:MaxEmbedTextLength]
	}

	// Call the Bedrock client
	req := clients.EmbeddingRequest{
		Model: EmbeddingModelID,
		Input: trimmed,
	}

	resp, err := s.client.CreateEmbedding(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Extract embedding from response
	if len(resp.Data) == 0 || len(resp.Data[0].Embedding) == 0 {
		return nil, errors.New("empty embedding response")
	}

	return resp.Data[0].Embedding, nil
}

// BatchGenerateEmbeddings generates embeddings for multiple tracks.
// Returns a map of trackID to embedding vector.
// Continues on individual failures - returns partial results if some tracks fail.
// Returns error only if all tracks fail or context is cancelled.
func (s *EmbeddingService) BatchGenerateEmbeddings(ctx context.Context, tracks []models.Track) (map[string][]float32, error) {
	result := make(map[string][]float32)

	// Return empty map for empty/nil input
	if len(tracks) == 0 {
		return result, nil
	}

	var lastErr error
	successCount := 0

	for _, track := range tracks {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		embedding, err := s.GenerateTrackEmbedding(ctx, track)
		if err != nil {
			lastErr = err
			continue
		}

		result[track.ID] = embedding
		successCount++
	}

	// Only return error if ALL tracks failed
	if successCount == 0 && lastErr != nil {
		return result, fmt.Errorf("all embeddings failed: %w", lastErr)
	}

	return result, nil
}
