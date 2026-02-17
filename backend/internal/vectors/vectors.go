package vectors

import "context"

// SimilarResult represents a similarity search result
type SimilarResult struct {
	ID       string            `json:"id"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// VectorService defines vector storage and similarity search operations
type VectorService interface {
	// PutVector stores a vector with metadata
	PutVector(ctx context.Context, id string, vector []float32, metadata map[string]string) error

	// GetVector retrieves a vector by ID
	GetVector(ctx context.Context, id string) ([]float32, error)

	// QuerySimilar finds k most similar vectors
	QuerySimilar(ctx context.Context, vector []float32, k int) ([]SimilarResult, error)

	// DeleteVector removes a vector by ID
	DeleteVector(ctx context.Context, id string) error
}
