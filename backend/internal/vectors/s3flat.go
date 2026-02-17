package vectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// storedVector is the JSON format for vectors in S3.
type storedVector struct {
	ID       string            `json:"id"`
	Vector   []float32         `json:"vector"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// S3FlatService stores vectors as JSON files in a regular S3 bucket.
// Path: s3://{bucket}/{prefix}/{id}.json
type S3FlatService struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3FlatService(bucket, prefix string) *S3FlatService {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	return &S3FlatService{client: s3.NewFromConfig(cfg), bucket: bucket, prefix: prefix}
}

func (s *S3FlatService) key(id string) string {
	return fmt.Sprintf("%s/%s.json", s.prefix, id)
}

func (s *S3FlatService) PutVector(ctx context.Context, id string, vector []float32, metadata map[string]string) error {
	data, err := json.Marshal(storedVector{ID: id, Vector: vector, Metadata: metadata})
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         aws.String(s.key(id)),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
	})
	return err
}

func (s *S3FlatService) GetVector(ctx context.Context, id string) ([]float32, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(s.key(id)),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var sv storedVector
	if err := json.Unmarshal(body, &sv); err != nil {
		return nil, err
	}
	return sv.Vector, nil
}

func (s *S3FlatService) DeleteVector(ctx context.Context, id string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(s.key(id)),
	})
	return err
}

// QuerySimilar is not supported for S3 flat files (no kNN index).
// Use S3 Vectors for similarity queries.
func (s *S3FlatService) QuerySimilar(_ context.Context, _ []float32, _ int) ([]SimilarResult, error) {
	return nil, fmt.Errorf("QuerySimilar not supported on S3 flat storage; use S3 Vectors")
}
