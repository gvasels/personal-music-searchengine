package vectors

import (
	"context"

	"github.com/gvasels/personal-music-searchengine/internal/service"
)

// ServiceAdapter wraps S3VectorsService to implement service.VectorService
type ServiceAdapter struct {
	svc *S3VectorsService
}

// NewServiceAdapter creates a new adapter
func NewServiceAdapter(bucketName, indexName string) *ServiceAdapter {
	return &ServiceAdapter{svc: NewS3VectorsService(bucketName, indexName)}
}

func (a *ServiceAdapter) PutVector(ctx context.Context, id string, vector []float32, metadata map[string]string) error {
	return a.svc.PutVector(ctx, id, vector, metadata)
}

func (a *ServiceAdapter) QuerySimilar(ctx context.Context, vector []float32, k int) ([]service.VectorResult, error) {
	results, err := a.svc.QuerySimilar(ctx, vector, k)
	if err != nil {
		return nil, err
	}
	out := make([]service.VectorResult, len(results))
	for i, r := range results {
		out[i] = service.VectorResult{
			ID:       r.ID,
			Score:    r.Score,
			Metadata: r.Metadata,
		}
	}
	return out, nil
}

func (a *ServiceAdapter) DeleteVector(ctx context.Context, id string) error {
	return a.svc.DeleteVector(ctx, id)
}

func (a *ServiceAdapter) GetVector(ctx context.Context, id string) ([]float32, error) {
	return a.svc.GetVector(ctx, id)
}
