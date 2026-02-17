package vectors

import (
	"context"
	"fmt"

	"github.com/gvasels/personal-music-searchengine/internal/service"
)

// DualWriteAdapter writes vectors to both S3 Vectors (managed kNN) and S3 (flat JSON).
// Reads and queries go to S3 Vectors. S3 flat files serve as backup and for native testing.
type DualWriteAdapter struct {
	primary *S3VectorsService
	backup  *S3FlatService
}

func NewDualWriteAdapter(vectorBucket, indexName, s3Bucket, s3Prefix string) *DualWriteAdapter {
	return &DualWriteAdapter{
		primary: NewS3VectorsService(vectorBucket, indexName),
		backup:  NewS3FlatService(s3Bucket, s3Prefix),
	}
}

func (d *DualWriteAdapter) PutVector(ctx context.Context, id string, vector []float32, metadata map[string]string) error {
	// Write to both — primary failure is fatal, backup failure is logged
	if err := d.primary.PutVector(ctx, id, vector, metadata); err != nil {
		return err
	}
	if err := d.backup.PutVector(ctx, id, vector, metadata); err != nil {
		fmt.Printf("Warning: S3 flat backup write failed for %s: %v\n", id, err)
	}
	return nil
}

func (d *DualWriteAdapter) GetVector(ctx context.Context, id string) ([]float32, error) {
	return d.primary.GetVector(ctx, id)
}

func (d *DualWriteAdapter) QuerySimilar(ctx context.Context, vector []float32, k int) ([]service.VectorResult, error) {
	results, err := d.primary.QuerySimilar(ctx, vector, k)
	if err != nil {
		return nil, err
	}
	out := make([]service.VectorResult, len(results))
	for i, r := range results {
		out[i] = service.VectorResult{ID: r.ID, Score: r.Score, Metadata: r.Metadata}
	}
	return out, nil
}

func (d *DualWriteAdapter) DeleteVector(ctx context.Context, id string) error {
	if err := d.primary.DeleteVector(ctx, id); err != nil {
		return err
	}
	if err := d.backup.DeleteVector(ctx, id); err != nil {
		fmt.Printf("Warning: S3 flat backup delete failed for %s: %v\n", id, err)
	}
	return nil
}
