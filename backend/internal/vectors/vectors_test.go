package vectors

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestService(t *testing.T) *S3VectorsService {
	bucket := os.Getenv("VECTOR_BUCKET_NAME")
	index := os.Getenv("VECTOR_INDEX_NAME")
	if bucket == "" || index == "" {
		t.Skip("VECTOR_BUCKET_NAME and VECTOR_INDEX_NAME required for integration tests")
	}
	return NewS3VectorsService(bucket, index)
}

func TestPutVector(t *testing.T) {
	ctx := context.Background()
	svc := getTestService(t)

	vector := make([]float32, 1024)
	for i := range vector {
		vector[i] = float32(i+1) / 1024.0 // Avoid zero vector for cosine
	}

	err := svc.PutVector(ctx, "test-track-123", vector, map[string]string{
		"trackId": "test-track-123",
		"userId":  "user-456",
	})
	require.NoError(t, err)

	// Cleanup
	_ = svc.DeleteVector(ctx, "test-track-123")
}

func TestQuerySimilar(t *testing.T) {
	ctx := context.Background()
	svc := getTestService(t)

	// Insert test vector first
	vector := make([]float32, 1024)
	for i := range vector {
		vector[i] = float32(i+1) / 1024.0
	}
	err := svc.PutVector(ctx, "test-query-track", vector, nil)
	require.NoError(t, err)

	// Query
	results, err := svc.QuerySimilar(ctx, vector, 5)
	require.NoError(t, err)
	assert.NotEmpty(t, results)

	// Cleanup
	_ = svc.DeleteVector(ctx, "test-query-track")
}

func TestDeleteVector(t *testing.T) {
	ctx := context.Background()
	svc := getTestService(t)

	// Insert then delete
	vector := make([]float32, 1024)
	for i := range vector {
		vector[i] = float32(i+1) / 1024.0
	}
	_ = svc.PutVector(ctx, "test-delete-track", vector, nil)

	err := svc.DeleteVector(ctx, "test-delete-track")
	require.NoError(t, err)
}
