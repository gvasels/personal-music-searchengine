package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SecretsManagerClient defines the interface for fetching secrets.
type SecretsManagerClient interface {
	GetSecretString(ctx context.Context, secretName string) (string, error)
}

// EmbeddingGateway defines the interface for generating embeddings via the gateway.
type EmbeddingGateway interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

// EmbeddingGatewayClient calls the Bedrock gateway /v1/embeddings endpoint over HTTP.
type EmbeddingGatewayClient struct {
	gatewayURL string
	secretName string
	secrets    SecretsManagerClient
	httpClient *http.Client

	apiKey     string
	apiKeyOnce sync.Once
	apiKeyErr  error
}

// NewEmbeddingGatewayClient creates a new client.
func NewEmbeddingGatewayClient(gatewayURL, secretName string, secrets SecretsManagerClient) *EmbeddingGatewayClient {
	return &EmbeddingGatewayClient{
		gatewayURL: gatewayURL,
		secretName: secretName,
		secrets:    secrets,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GenerateEmbedding generates an embedding vector for the given text.
func (c *EmbeddingGatewayClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	key, err := c.getAPIKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get api key: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
		}

		embedding, err := c.doRequest(ctx, key, text)
		if err == nil {
			return embedding, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *EmbeddingGatewayClient) doRequest(ctx context.Context, apiKey, text string) ([]float32, error) {
	body, _ := json.Marshal(EmbeddingRequest{
		Model: "text-embedding-3-small",
		Input: text,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", c.gatewayURL+"/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway returned %d", resp.StatusCode)
	}

	var embResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embResp.Data) == 0 || len(embResp.Data[0].Embedding) == 0 {
		return nil, errors.New("empty embedding response")
	}

	return embResp.Data[0].Embedding, nil
}

func (c *EmbeddingGatewayClient) getAPIKey(ctx context.Context) (string, error) {
	c.apiKeyOnce.Do(func() {
		c.apiKey, c.apiKeyErr = c.secrets.GetSecretString(ctx, c.secretName)
	})
	return c.apiKey, c.apiKeyErr
}
