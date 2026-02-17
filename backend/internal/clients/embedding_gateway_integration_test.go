package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddingGatewayClient_Integration_FullFlow(t *testing.T) {
	expectedEmbedding := make([]float32, 1024)
	for i := range expectedEmbedding {
		expectedEmbedding[i] = float32(i) * 0.001
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request format
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/embeddings", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer integration-test-key", r.Header.Get("Authorization"))

		var req EmbeddingRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "text-embedding-3-small", req.Model)
		assert.NotEmpty(t, req.Input)

		resp := EmbeddingResponse{
			Object: "list",
			Model:  "text-embedding-3-small",
			Data:   []EmbeddingData{{Object: "embedding", Index: 0, Embedding: expectedEmbedding}},
			Usage:  &UsageInfo{PromptTokens: 10, TotalTokens: 10},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "test-secret", &mockSecretsManager{secret: "integration-test-key"})

	embedding, err := client.GenerateEmbedding(context.Background(), "Bohemian Rhapsody Queen A Night at the Opera Rock")
	require.NoError(t, err)
	assert.Len(t, embedding, 1024)
	assert.InDelta(t, 0.0, embedding[0], 0.001)
}

func TestEmbeddingGatewayClient_Integration_GatewayDown(t *testing.T) {
	// Server that immediately closes connections
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "test-secret", &mockSecretsManager{secret: "key"})

	_, err := client.GenerateEmbedding(context.Background(), "test query")
	assert.Error(t, err)
}

func TestEmbeddingGatewayClient_Integration_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // Slow response
		resp := EmbeddingResponse{
			Object: "list",
			Data:   []EmbeddingData{{Embedding: []float32{0.1}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "test-secret", &mockSecretsManager{secret: "key"})

	// With a very short context timeout, should fail
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.GenerateEmbedding(ctx, "test")
	assert.Error(t, err)
}
