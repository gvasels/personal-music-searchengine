package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/clients"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockBedrockEmbeddingClient implements BedrockEmbeddingClient for testing
type MockBedrockEmbeddingClient struct {
	CreateEmbeddingFunc func(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error)
	CallCount           int
	LastRequest         clients.EmbeddingRequest
}

func (m *MockBedrockEmbeddingClient) CreateEmbedding(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error) {
	m.CallCount++
	m.LastRequest = req
	if m.CreateEmbeddingFunc != nil {
		return m.CreateEmbeddingFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

// Helper to create a mock that returns a valid 1024-dim embedding
func newSuccessEmbeddingMock() *MockBedrockEmbeddingClient {
	return &MockBedrockEmbeddingClient{
		CreateEmbeddingFunc: func(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error) {
			embedding := make([]float32, 1024)
			for i := range embedding {
				embedding[i] = float32(i) * 0.001
			}
			return &clients.EmbeddingResponse{
				Object: "list",
				Model:  "text-embedding-3-small",
				Data: []clients.EmbeddingData{
					{
						Object:    "embedding",
						Index:     0,
						Embedding: embedding,
					},
				},
			}, nil
		},
	}
}

// Helper to create a mock that returns an error
func newErrorEmbeddingMock(err error) *MockBedrockEmbeddingClient {
	return &MockBedrockEmbeddingClient{
		CreateEmbeddingFunc: func(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error) {
			return nil, err
		},
	}
}

// createTestTrack creates a track with all fields populated for testing
func createFullTestTrackForEmbedding() models.Track {
	return models.Track{
		ID:         "track-123",
		UserID:     "user-456",
		Title:      "Midnight Dreams",
		Artist:     "Synthwave Master",
		Album:      "Neon Nights",
		Genre:      "Electronic",
		Tags:       []string{"synthwave", "retro", "80s", "driving"},
		BPM:        128,
		KeyCamelot: "8A",
		Duration:   240,
		Year:       2023,
	}
}

// createMinimalTestTrack creates a track with only required fields
func createMinimalTestTrackForEmbedding() models.Track {
	return models.Track{
		ID:     "track-minimal",
		UserID: "user-456",
		Title:  "Untitled Track",
	}
}

// ============================================================================
// ComposeEmbedText Tests
// ============================================================================

func TestComposeEmbedText_AllFields(t *testing.T) {
	track := createFullTestTrackForEmbedding()
	svc := NewEmbeddingService(newSuccessEmbeddingMock())

	text := svc.ComposeEmbedText(track)

	require.NotEmpty(t, text, "composed text should not be empty")
	assert.Contains(t, text, "Midnight Dreams", "should contain title")
	assert.Contains(t, text, "Synthwave Master", "should contain artist")
	assert.Contains(t, text, "Neon Nights", "should contain album")
	assert.Contains(t, text, "Electronic", "should contain genre")
	assert.Contains(t, text, "synthwave", "should contain tags")
	assert.Contains(t, text, "128", "should contain BPM")
	assert.Contains(t, text, "8A", "should contain key")
}

func TestComposeEmbedText_MinimalFields(t *testing.T) {
	track := createMinimalTestTrackForEmbedding()
	svc := NewEmbeddingService(newSuccessEmbeddingMock())

	text := svc.ComposeEmbedText(track)

	require.NotEmpty(t, text, "composed text should not be empty even with minimal fields")
	assert.Contains(t, text, "Untitled Track", "should contain title")
	assert.NotContains(t, text, "nil", "should not contain nil")
	assert.NotContains(t, text, "<nil>", "should not contain <nil>")
}

func TestComposeEmbedText_MaxLength(t *testing.T) {
	track := models.Track{
		ID:     "track-long",
		UserID: "user-456",
		Title:  strings.Repeat("Very Long Title ", 500),
		Artist: strings.Repeat("Long Artist Name ", 500),
		Album:  strings.Repeat("Long Album Name ", 500),
		Genre:  "Electronic",
		Tags:   []string{strings.Repeat("tag", 1000)},
	}
	svc := NewEmbeddingService(newSuccessEmbeddingMock())

	text := svc.ComposeEmbedText(track)

	maxLength := 8000
	assert.LessOrEqual(t, len(text), maxLength,
		"composed text should be truncated to max %d characters, got %d", maxLength, len(text))
	assert.NotEmpty(t, text, "truncated text should not be empty")
}

func TestComposeEmbedText_WithTags(t *testing.T) {
	track := models.Track{
		ID:     "track-tags",
		UserID: "user-456",
		Title:  "Tagged Track",
		Tags:   []string{"chill", "ambient", "focus", "study music"},
	}
	svc := NewEmbeddingService(newSuccessEmbeddingMock())

	text := svc.ComposeEmbedText(track)

	assert.Contains(t, text, "chill", "should contain first tag")
	assert.Contains(t, text, "ambient", "should contain second tag")
	assert.Contains(t, text, "focus", "should contain third tag")
	assert.Contains(t, text, "study music", "should contain multi-word tag")
}

func TestComposeEmbedText_WithDJMetadata(t *testing.T) {
	track := models.Track{
		ID:         "track-dj",
		UserID:     "user-456",
		Title:      "DJ Track",
		Artist:     "DJ Producer",
		BPM:        140,
		KeyCamelot: "11B",
	}
	svc := NewEmbeddingService(newSuccessEmbeddingMock())

	text := svc.ComposeEmbedText(track)

	assert.Contains(t, text, "140", "should contain BPM value")
	assert.Contains(t, text, "11B", "should contain Camelot key")
}

// ============================================================================
// GenerateTrackEmbedding Tests
// ============================================================================

func TestGenerateTrackEmbedding_Success(t *testing.T) {
	track := createFullTestTrackForEmbedding()
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	embedding, err := svc.GenerateTrackEmbedding(ctx, track)

	require.NoError(t, err, "should not return error on success")
	require.NotNil(t, embedding, "embedding should not be nil")
	assert.Len(t, embedding, 1024, "Bedrock Titan returns 1024-dimensional embeddings")
	assert.Equal(t, 1, mockClient.CallCount, "should call Bedrock exactly once")
}

func TestGenerateTrackEmbedding_BedrockError(t *testing.T) {
	track := createFullTestTrackForEmbedding()
	bedrockErr := errors.New("bedrock service unavailable")
	mockClient := newErrorEmbeddingMock(bedrockErr)
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	embedding, err := svc.GenerateTrackEmbedding(ctx, track)

	require.Error(t, err, "should return error when Bedrock fails")
	assert.Nil(t, embedding, "embedding should be nil on error")
	assert.Equal(t, 1, mockClient.CallCount, "should attempt to call Bedrock")
}

func TestGenerateTrackEmbedding_ContextCancelled(t *testing.T) {
	track := createFullTestTrackForEmbedding()
	mockClient := &MockBedrockEmbeddingClient{
		CreateEmbeddingFunc: func(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error) {
			return nil, ctx.Err()
		},
	}
	svc := NewEmbeddingService(mockClient)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	embedding, err := svc.GenerateTrackEmbedding(ctx, track)

	require.Error(t, err, "should return error when context is cancelled")
	assert.Nil(t, embedding, "embedding should be nil when context is cancelled")
}

// ============================================================================
// GenerateQueryEmbedding Tests
// ============================================================================

func TestGenerateQueryEmbedding_Success(t *testing.T) {
	query := "upbeat electronic music for workout"
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	embedding, err := svc.GenerateQueryEmbedding(ctx, query)

	require.NoError(t, err, "should not return error on success")
	require.NotNil(t, embedding, "embedding should not be nil")
	assert.Len(t, embedding, 1024, "query embedding should be 1024-dimensional")
	assert.Equal(t, 1, mockClient.CallCount, "should call Bedrock exactly once")
}

func TestGenerateQueryEmbedding_EmptyQuery(t *testing.T) {
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	embedding, err := svc.GenerateQueryEmbedding(ctx, "")

	require.Error(t, err, "should return error for empty query")
	assert.Nil(t, embedding, "embedding should be nil for empty query")
	assert.Equal(t, 0, mockClient.CallCount, "should not call Bedrock for empty query")
}

func TestGenerateQueryEmbedding_WhitespaceQuery(t *testing.T) {
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	embedding, err := svc.GenerateQueryEmbedding(ctx, "   \t\n  ")

	require.Error(t, err, "should return error for whitespace-only query")
	assert.Nil(t, embedding, "embedding should be nil for whitespace query")
	assert.Equal(t, 0, mockClient.CallCount, "should not call Bedrock for whitespace query")
}

func TestGenerateQueryEmbedding_LongQuery(t *testing.T) {
	longQuery := strings.Repeat("music ", 2000)
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	embedding, err := svc.GenerateQueryEmbedding(ctx, longQuery)

	require.NoError(t, err, "should handle long queries")
	require.NotNil(t, embedding, "should return embedding for long query")
}

// ============================================================================
// BatchGenerateEmbeddings Tests
// ============================================================================

func TestBatchGenerateEmbeddings_Success(t *testing.T) {
	tracks := []models.Track{
		{ID: "track-1", UserID: "user-1", Title: "Track One", Artist: "Artist A"},
		{ID: "track-2", UserID: "user-1", Title: "Track Two", Artist: "Artist B"},
		{ID: "track-3", UserID: "user-1", Title: "Track Three", Artist: "Artist C"},
	}
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	results, err := svc.BatchGenerateEmbeddings(ctx, tracks)

	require.NoError(t, err, "should not return error on success")
	require.NotNil(t, results, "results should not be nil")
	assert.Len(t, results, 3, "should return embedding for each track")

	for trackID, embedding := range results {
		assert.NotEmpty(t, trackID, "track ID should not be empty")
		assert.Len(t, embedding, 1024, "each embedding should be 1024-dimensional")
	}

	assert.Equal(t, 3, mockClient.CallCount, "should call Bedrock for each track")
}

func TestBatchGenerateEmbeddings_EmptySlice(t *testing.T) {
	mockClient := newSuccessEmbeddingMock()
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	results, err := svc.BatchGenerateEmbeddings(ctx, []models.Track{})

	require.NoError(t, err, "should not return error for empty slice")
	assert.Empty(t, results, "should return empty map for empty input")
	assert.Equal(t, 0, mockClient.CallCount, "should not call Bedrock for empty input")
}

func TestBatchGenerateEmbeddings_PartialFailure(t *testing.T) {
	tracks := []models.Track{
		{ID: "track-1", UserID: "user-1", Title: "Track One"},
		{ID: "track-2", UserID: "user-1", Title: "Track Two"},
		{ID: "track-3", UserID: "user-1", Title: "Track Three"},
	}

	callCount := 0
	mockClient := &MockBedrockEmbeddingClient{
		CreateEmbeddingFunc: func(ctx context.Context, req clients.EmbeddingRequest) (*clients.EmbeddingResponse, error) {
			callCount++
			if callCount == 2 {
				return nil, errors.New("temporary bedrock error")
			}
			embedding := make([]float32, 1024)
			return &clients.EmbeddingResponse{
				Data: []clients.EmbeddingData{{Embedding: embedding}},
			}, nil
		},
	}
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	results, err := svc.BatchGenerateEmbeddings(ctx, tracks)

	require.NoError(t, err, "should not return error on partial failure")
	require.NotNil(t, results, "results should not be nil")
	assert.Len(t, results, 2, "should return embeddings for successful tracks")
	assert.Contains(t, results, "track-1", "should have embedding for track-1")
	assert.Contains(t, results, "track-3", "should have embedding for track-3")
	assert.NotContains(t, results, "track-2", "should not have embedding for failed track-2")
	assert.Equal(t, 3, callCount, "should attempt all tracks")
}

func TestBatchGenerateEmbeddings_AllFailures(t *testing.T) {
	tracks := []models.Track{
		{ID: "track-1", UserID: "user-1", Title: "Track One"},
		{ID: "track-2", UserID: "user-1", Title: "Track Two"},
	}
	mockClient := newErrorEmbeddingMock(errors.New("bedrock down"))
	svc := NewEmbeddingService(mockClient)
	ctx := context.Background()

	results, err := svc.BatchGenerateEmbeddings(ctx, tracks)

	require.Error(t, err, "should return error when all embeddings fail")
	assert.Empty(t, results, "should return empty results on total failure")
}

// ============================================================================
// NewEmbeddingService Tests
// ============================================================================

func TestNewEmbeddingService_NilClient(t *testing.T) {
	assert.Panics(t, func() {
		NewEmbeddingService(nil)
	}, "should panic when client is nil")
}

func TestNewEmbeddingService_ValidClient(t *testing.T) {
	mockClient := newSuccessEmbeddingMock()

	svc := NewEmbeddingService(mockClient)

	assert.NotNil(t, svc, "service should not be nil")
}
