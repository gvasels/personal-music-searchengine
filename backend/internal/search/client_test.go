package search

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockLambdaClient implements LambdaInvoker for testing
type mockLambdaClient struct {
	response   *lambda.InvokeOutput
	err        error
	lastInput  *lambda.InvokeInput
	invokeFunc func(ctx context.Context, params *lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

func (m *mockLambdaClient) Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	m.lastInput = params
	if m.invokeFunc != nil {
		return m.invokeFunc(ctx, params)
	}
	return m.response, m.err
}

func TestSearch_SimpleQuery(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: true,
		Data: SearchResponse{
			Results: []SearchResult{
				{ID: "track-1", Title: "Test Song", Artist: "Test Artist", Score: 0.95},
				{ID: "track-2", Title: "Another Song", Artist: "Test Artist", Score: 0.85},
			},
			Total: 2,
		},
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	resp, err := client.Search(context.Background(), "user-123", SearchQuery{
		Query: "test",
		Limit: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "track-1", resp.Results[0].ID)
}

func TestSearch_WithFilters(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: true,
		Data: SearchResponse{
			Results: []SearchResult{
				{ID: "track-1", Title: "Rock Song", Artist: "Rock Band", Genre: "Rock", Score: 0.9},
			},
			Total: 1,
		},
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	resp, err := client.Search(context.Background(), "user-123", SearchQuery{
		Query: "rock",
		Filters: SearchFilters{
			Genre:    "Rock",
			YearFrom: 2020,
			YearTo:   2024,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, "Rock", resp.Results[0].Genre)

	// Verify user filter was added
	var req NixiesearchRequest
	err = json.Unmarshal(mockClient.lastInput.Payload, &req)
	require.NoError(t, err)
	queryPayload, err := json.Marshal(req.Payload)
	require.NoError(t, err)
	var query SearchQuery
	err = json.Unmarshal(queryPayload, &query)
	require.NoError(t, err)
	assert.Equal(t, "user-123", query.Filters.UserID)
}

func TestSearch_Pagination(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: true,
		Data: SearchResponse{
			Results:    []SearchResult{{ID: "track-21", Title: "Song 21", Score: 0.5}},
			Total:      50,
			NextCursor: "cursor-page-3",
		},
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	resp, err := client.Search(context.Background(), "user-123", SearchQuery{
		Query:  "song",
		Limit:  20,
		Cursor: "cursor-page-2",
	})

	require.NoError(t, err)
	assert.Equal(t, 50, resp.Total)
	assert.Equal(t, "cursor-page-3", resp.NextCursor)
}

func TestIndex_NewDocument(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: true,
		Data: IndexResponse{
			ID:      "track-new",
			Indexed: true,
		},
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	resp, err := client.Index(context.Background(), Document{
		ID:     "track-new",
		UserID: "user-123",
		Title:  "New Song",
		Artist: "New Artist",
	})

	require.NoError(t, err)
	assert.True(t, resp.Indexed)
	assert.Equal(t, "track-new", resp.ID)
}

func TestDelete_RemovesFromIndex(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: true,
		Data: DeleteResponse{
			ID:      "track-delete",
			Deleted: true,
		},
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	resp, err := client.Delete(context.Background(), "track-delete")

	require.NoError(t, err)
	assert.True(t, resp.Deleted)
}

func TestBulkIndex_Success(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: true,
		Data: BulkIndexResponse{
			Indexed: 5,
			Failed:  0,
		},
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	docs := []Document{
		{ID: "track-1", Title: "Song 1"},
		{ID: "track-2", Title: "Song 2"},
		{ID: "track-3", Title: "Song 3"},
		{ID: "track-4", Title: "Song 4"},
		{ID: "track-5", Title: "Song 5"},
	}

	resp, err := client.BulkIndex(context.Background(), docs)

	require.NoError(t, err)
	assert.Equal(t, 5, resp.Indexed)
	assert.Equal(t, 0, resp.Failed)
}

func TestSearch_LambdaError(t *testing.T) {
	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			FunctionError: aws.String("Unhandled"),
			Payload:       []byte(`{"errorMessage": "timeout"}`),
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	_, err := client.Search(context.Background(), "user-123", SearchQuery{Query: "test"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lambda function error")
}

func TestSearch_OperationError(t *testing.T) {
	mockResp := NixiesearchResponse{
		Success: false,
		Error:   "index not found",
	}
	payload, _ := json.Marshal(mockResp)

	mockClient := &mockLambdaClient{
		response: &lambda.InvokeOutput{
			Payload: payload,
		},
	}

	client := NewClient(mockClient, "nixiesearch-lambda")
	_, err := client.Search(context.Background(), "user-123", SearchQuery{Query: "test"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "index not found")
}
