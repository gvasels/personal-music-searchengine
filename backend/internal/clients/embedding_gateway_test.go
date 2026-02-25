package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSecretsManager implements SecretsManagerClient for testing.
type mockSecretsManager struct {
	secret string
	err    error
}

func (m *mockSecretsManager) GetSecretString(ctx context.Context, secretName string) (string, error) {
	return m.secret, m.err
}

func TestEmbeddingGatewayClient_GenerateEmbedding_Success(t *testing.T) {
	expectedEmbedding := []float32{0.1, 0.2, 0.3}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/embeddings", r.URL.Path)
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

		var req EmbeddingRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "text-embedding-3-small", req.Model)

		resp := EmbeddingResponse{
			Object: "list",
			Data:   []EmbeddingData{{Embedding: expectedEmbedding}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "secret-name", &mockSecretsManager{secret: "test-api-key"})
	embedding, err := client.GenerateEmbedding(context.Background(), "test text")

	require.NoError(t, err)
	assert.Equal(t, expectedEmbedding, embedding)
}

func TestEmbeddingGatewayClient_GenerateEmbedding_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := EmbeddingResponse{Object: "list", Data: []EmbeddingData{}}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "secret-name", &mockSecretsManager{secret: "key"})
	_, err := client.GenerateEmbedding(context.Background(), "test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty embedding")
}

func TestEmbeddingGatewayClient_GenerateEmbedding_ServerError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "secret-name", &mockSecretsManager{secret: "key"})
	_, err := client.GenerateEmbedding(context.Background(), "test")

	assert.Error(t, err)
	assert.Equal(t, 3, attempts) // 1 initial + 2 retries
}

func TestEmbeddingGatewayClient_GenerateEmbedding_AuthFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "secret-name", &mockSecretsManager{secret: "bad-key"})
	_, err := client.GenerateEmbedding(context.Background(), "test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func TestEmbeddingGatewayClient_GenerateEmbedding_SecretsFetchError(t *testing.T) {
	client := NewEmbeddingGatewayClient("http://localhost", "secret-name", &mockSecretsManager{err: assert.AnError})
	_, err := client.GenerateEmbedding(context.Background(), "test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api key")
}

func TestEmbeddingGatewayClient_CachesAPIKey(t *testing.T) {
	calls := 0
	sm := &mockSecretsManager{secret: "cached-key"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		resp := EmbeddingResponse{
			Object: "list",
			Data:   []EmbeddingData{{Embedding: []float32{0.1}}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewEmbeddingGatewayClient(server.URL, "secret-name", sm)

	_, err := client.GenerateEmbedding(context.Background(), "first")
	require.NoError(t, err)
	_, err = client.GenerateEmbedding(context.Background(), "second")
	require.NoError(t, err)

	assert.Equal(t, 2, calls) // 2 HTTP calls but API key fetched only once (sync.Once)
}
