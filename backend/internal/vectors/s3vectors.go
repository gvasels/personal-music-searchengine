package vectors

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3vectors"
	"github.com/aws/aws-sdk-go-v2/service/s3vectors/types"
)

// S3VectorsService implements VectorService using S3 Vectors
type S3VectorsService struct {
	client     *s3vectors.Client
	bucketName string
	indexName  string
}

// NewS3VectorsService creates a new S3 Vectors service
func NewS3VectorsService(bucketName, indexName string) *S3VectorsService {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client := s3vectors.NewFromConfig(cfg)
	return &S3VectorsService{
		client:     client,
		bucketName: bucketName,
		indexName:  indexName,
	}
}

// NewS3VectorsServiceWithClient creates service with custom client (for testing)
func NewS3VectorsServiceWithClient(client *s3vectors.Client, bucketName, indexName string) *S3VectorsService {
	return &S3VectorsService{
		client:     client,
		bucketName: bucketName,
		indexName:  indexName,
	}
}

func (s *S3VectorsService) PutVector(ctx context.Context, id string, vector []float32, metadata map[string]string) error {
	input := &s3vectors.PutVectorsInput{
		VectorBucketName: &s.bucketName,
		IndexName:        &s.indexName,
		Vectors: []types.PutInputVector{
			{
				Key:  aws.String(id),
				Data: &types.VectorDataMemberFloat32{Value: vector},
			},
		},
	}
	// Note: Metadata requires document.Interface - skipping for now as it needs JSON marshaling
	_, err := s.client.PutVectors(ctx, input)
	return err
}

func (s *S3VectorsService) QuerySimilar(ctx context.Context, vector []float32, k int) ([]SimilarResult, error) {
	resp, err := s.client.QueryVectors(ctx, &s3vectors.QueryVectorsInput{
		VectorBucketName: &s.bucketName,
		IndexName:        &s.indexName,
		TopK:             aws.Int32(int32(k)),
		QueryVector:      &types.VectorDataMemberFloat32{Value: vector},
		ReturnDistance:   true,
		ReturnMetadata:   true,
	})
	if err != nil {
		return nil, err
	}

	results := make([]SimilarResult, 0, len(resp.Vectors))
	for _, v := range resp.Vectors {
		result := SimilarResult{}
		if v.Key != nil {
			result.ID = *v.Key
		}
		if v.Distance != nil {
			result.Score = *v.Distance
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *S3VectorsService) DeleteVector(ctx context.Context, id string) error {
	_, err := s.client.DeleteVectors(ctx, &s3vectors.DeleteVectorsInput{
		VectorBucketName: &s.bucketName,
		IndexName:        &s.indexName,
		Keys:             []string{id},
	})
	return err
}

func (s *S3VectorsService) GetVector(ctx context.Context, id string) ([]float32, error) {
	resp, err := s.client.GetVectors(ctx, &s3vectors.GetVectorsInput{
		VectorBucketName: &s.bucketName,
		IndexName:        &s.indexName,
		Keys:             []string{id},
		ReturnData:       true,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Vectors) == 0 {
		return nil, nil
	}
	if data, ok := resp.Vectors[0].Data.(*types.VectorDataMemberFloat32); ok {
		return data.Value, nil
	}
	return nil, nil
}
