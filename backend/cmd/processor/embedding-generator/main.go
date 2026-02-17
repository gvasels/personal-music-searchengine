package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3vectors"
	s3vtypes "github.com/aws/aws-sdk-go-v2/service/s3vectors/types"
	"github.com/google/uuid"
)

const MarengoModelID = "twelvelabs.marengo-embed-3-0-v1:0"

type Event struct {
	TrackID string `json:"trackId"`
	UserID  string `json:"userId"`
	S3Key   string `json:"s3Key"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
}

type Result struct {
	TrackID     string `json:"trackId"`
	UserID      string `json:"userId"`
	EmbeddingID string `json:"embeddingId,omitempty"`
	Error       string `json:"error,omitempty"`
}

var (
	bedrockClient  *bedrockruntime.Client
	s3Client       *s3.Client
	s3VectorsClient *s3vectors.Client
	mediaBucket    string
	vectorBucket   string
	vectorIndex    string
	accountID      string
)

func init() {
	mediaBucket = os.Getenv("MEDIA_BUCKET")
	accountID = os.Getenv("AWS_ACCOUNT_ID")
	vectorBucket = os.Getenv("VECTOR_BUCKET_NAME")
	vectorIndex = os.Getenv("VECTOR_INDEX_NAME")
	if vectorIndex == "" {
		vectorIndex = "media-embeddings"
	}

	if mediaBucket == "" || accountID == "" {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return
	}

	bedrockClient = bedrockruntime.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)
	if vectorBucket != "" {
		s3VectorsClient = s3vectors.NewFromConfig(cfg)
	}
}

func handler(ctx context.Context, event Event) (Result, error) {
	result := Result{TrackID: event.TrackID, UserID: event.UserID}

	if bedrockClient == nil {
		result.Error = "bedrock client not configured"
		return result, nil
	}

	s3URI := fmt.Sprintf("s3://%s/%s", mediaBucket, event.S3Key)
	inferenceID := uuid.New().String()
	outputURI := fmt.Sprintf("s3://%s/embeddings/%s", mediaBucket, inferenceID)

	// Build Marengo request for audio embedding
	marengoReq := document.NewLazyDocument(map[string]interface{}{
		"inputType": "audio",
		"audio": map[string]interface{}{
			"mediaSource": map[string]interface{}{
				"s3Location": map[string]interface{}{
					"uri":         s3URI,
					"bucketOwner": accountID,
				},
			},
			"embeddingOption": []string{"audio"},
			"embeddingScope":  []string{"asset"},
		},
	})

	asyncResp, err := bedrockClient.StartAsyncInvoke(ctx, &bedrockruntime.StartAsyncInvokeInput{
		ModelId:    aws.String(MarengoModelID),
		ModelInput: marengoReq,
		OutputDataConfig: &types.AsyncInvokeOutputDataConfigMemberS3OutputDataConfig{
			Value: types.AsyncInvokeS3OutputDataConfig{
				S3Uri:       aws.String(outputURI),
				BucketOwner: aws.String(accountID),
			},
		},
	})
	if err != nil {
		result.Error = fmt.Sprintf("failed to start async invoke: %v", err)
		return result, nil
	}

	invocationARN := aws.ToString(asyncResp.InvocationArn)

	// Poll for completion (max 100s)
	var status string
	var failureMsg string
	for i := 0; i < 20; i++ {
		time.Sleep(5 * time.Second)
		getResp, err := bedrockClient.GetAsyncInvoke(ctx, &bedrockruntime.GetAsyncInvokeInput{
			InvocationArn: aws.String(invocationARN),
		})
		if err != nil {
			continue
		}
		status = string(getResp.Status)
		if getResp.FailureMessage != nil {
			failureMsg = *getResp.FailureMessage
		}
		if status == "Completed" || status == "Failed" {
			break
		}
	}

	if status == "Failed" {
		result.Error = fmt.Sprintf("async invoke failed: %s", failureMsg)
		return result, nil
	}
	if status != "Completed" {
		result.Error = fmt.Sprintf("async invoke did not complete: status=%s", status)
		return result, nil
	}

	// Extract invocation ID from ARN for output path
	parts := strings.Split(invocationARN, "/")
	invocationID := parts[len(parts)-1]

	// Read embedding from Marengo output and store in S3 Vectors
	if s3VectorsClient != nil && vectorBucket != "" {
		embedding, err := readMarengoEmbedding(ctx, outputURI, invocationID)
		if err != nil {
			result.Error = fmt.Sprintf("failed to read embedding: %v", err)
			return result, nil
		}
		if err := storeInS3Vectors(ctx, event.TrackID, embedding); err != nil {
			result.Error = fmt.Sprintf("failed to store in S3 Vectors: %v", err)
			return result, nil
		}
	}

	result.EmbeddingID = invocationID
	return result, nil
}

// MarengoOutput represents the Marengo embedding output format
type MarengoOutput struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func readMarengoEmbedding(ctx context.Context, outputURI, invocationID string) ([]float32, error) {
	// Marengo stores output at: {outputURI}/{invocationID}/output.json
	outputKey := strings.TrimPrefix(outputURI, fmt.Sprintf("s3://%s/", mediaBucket))
	outputKey = fmt.Sprintf("%s/%s/output.json", outputKey, invocationID)

	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(mediaBucket),
		Key:    aws.String(outputKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get output.json: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read output: %w", err)
	}

	var output MarengoOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}

	if len(output.Data) == 0 || len(output.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("no embedding in output")
	}

	return output.Data[0].Embedding, nil
}

func storeInS3Vectors(ctx context.Context, trackID string, embedding []float32) error {
	_, err := s3VectorsClient.PutVectors(ctx, &s3vectors.PutVectorsInput{
		VectorBucketName: &vectorBucket,
		IndexName:        &vectorIndex,
		Vectors: []s3vtypes.PutInputVector{
			{
				Key:  aws.String(trackID),
				Data: &s3vtypes.VectorDataMemberFloat32{Value: embedding},
			},
		},
	})
	return err
}

func main() {
	lambda.Start(handler)
}
