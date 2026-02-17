package search

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockLambdaInvoker implements LambdaInvoker for testing.
type mockLambdaInvoker struct {
	response *lambda.InvokeOutput
	err      error
	captured []byte // captures the request payload
}

func (m *mockLambdaInvoker) Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	m.captured = params.Payload
	return m.response, m.err
}

func TestHybridSearch_Success(t *testing.T) {
	searchResp := NixiesearchResponse{
		Success: true,
		Data: SearchResponse{
			Results: []SearchResult{{ID: "track-1", Title: "Test", Score: 0.95}},
			Total:   1,
		},
	}
	payload, _ := json.Marshal(searchResp)

	mock := &mockLambdaInvoker{
		response: &lambda.InvokeOutput{Payload: payload},
	}

	client := NewClient(mock, "test-function")
	query := SearchQuery{
		Query:          "chill vibes",
		Embedding:      []float32{0.1, 0.2, 0.3},
		SemanticWeight: 0.5,
		Limit:          10,
	}

	resp, err := client.HybridSearch(context.Background(), "user-1", query)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, "track-1", resp.Results[0].ID)

	// Verify the request sent to Lambda includes hybrid_search operation
	var req NixiesearchRequest
	require.NoError(t, json.Unmarshal(mock.captured, &req))
	assert.Equal(t, "hybrid_search", req.Operation)
}

func TestHybridSearch_LambdaError(t *testing.T) {
	errResp := NixiesearchResponse{Success: false, Error: "internal error"}
	payload, _ := json.Marshal(errResp)

	mock := &mockLambdaInvoker{
		response: &lambda.InvokeOutput{Payload: payload},
	}

	client := NewClient(mock, "test-function")
	query := SearchQuery{
		Query:          "test",
		Embedding:      []float32{0.1},
		SemanticWeight: 0.5,
	}

	_, err := client.HybridSearch(context.Background(), "user-1", query)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hybrid search failed")
}

func TestSearch_StillWorksWithoutEmbedding(t *testing.T) {
	searchResp := NixiesearchResponse{
		Success: true,
		Data: SearchResponse{
			Results: []SearchResult{{ID: "track-2", Title: "Keyword Match", Score: 0.8}},
			Total:   1,
		},
	}
	payload, _ := json.Marshal(searchResp)

	mock := &mockLambdaInvoker{
		response: &lambda.InvokeOutput{Payload: payload},
	}

	client := NewClient(mock, "test-function")
	query := SearchQuery{Query: "keyword", Limit: 10}

	resp, err := client.Search(context.Background(), "user-1", query)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Total)

	// Verify it uses "search" operation (not hybrid_search)
	var req NixiesearchRequest
	require.NoError(t, json.Unmarshal(mock.captured, &req))
	assert.Equal(t, "search", req.Operation)
}

func TestDocument_EmbeddingIncludedInJSON(t *testing.T) {
	doc := Document{
		ID:        "test-id",
		Title:     "Test Track",
		Embedding: []float32{0.1, 0.2, 0.3},
	}

	data, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"embedding"`)

	// Without embedding, field should be omitted
	doc2 := Document{ID: "test-id", Title: "No Embedding"}
	data2, _ := json.Marshal(doc2)
	assert.NotContains(t, string(data2), `"embedding"`)
}
