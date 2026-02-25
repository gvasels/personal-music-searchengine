package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gvasels/personal-music-searchengine/internal/clients"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/search"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/validation"
	"github.com/gvasels/personal-music-searchengine/internal/vectors"
)

// Event represents the input from Step Functions
type Event struct {
	TrackID   string                 `json:"trackId"`
	UserID    string                 `json:"userId"`
	UploadID  string                 `json:"uploadId"`
	Metadata  *models.UploadMetadata `json:"metadata"`
	S3Key     string                 `json:"s3Key"`
	TableName string                 `json:"tableName"`
}

// Response represents the output to Step Functions
type Response struct {
	Indexed bool   `json:"indexed"`
	Reason  string `json:"reason,omitempty"`
}

var searchClient *search.Client
var repo repository.Repository
var embeddingGateway clients.EmbeddingGateway
var searchVectorSvc vectors.VectorService

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "MusicLibrary"
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo = repository.NewDynamoDBRepository(dynamoClient, tableName)

	nixieFunctionName := os.Getenv("NIXIESEARCH_FUNCTION_NAME")
	if nixieFunctionName == "" {
		fmt.Println("NIXIESEARCH_FUNCTION_NAME not set, search indexing disabled")
		return
	}

	lambdaClient := awslambda.NewFromConfig(cfg)
	searchClient = search.NewClient(lambdaClient, nixieFunctionName)

	// Initialize embedding gateway client if configured
	gatewayURL := os.Getenv("EMBEDDING_GATEWAY_URL")
	secretName := os.Getenv("EMBEDDING_API_KEY_SECRET")
	if gatewayURL != "" && secretName != "" {
		smClient := clients.NewAWSSecretsManager(secretsmanager.NewFromConfig(cfg))
		embeddingGateway = clients.NewEmbeddingGatewayClient(gatewayURL, secretName, smClient)
	}

	// Initialize S3 Vectors client for storing search text embeddings
	vectorBucket := os.Getenv("VECTOR_BUCKET_NAME")
	searchIndexName := os.Getenv("SEARCH_VECTOR_INDEX_NAME")
	if vectorBucket != "" && searchIndexName != "" {
		searchVectorSvc = vectors.NewS3VectorsService(vectorBucket, searchIndexName)
	}
}

func handleRequest(ctx context.Context, event Event) (*Response, error) {
	// Add timeout to context (5 seconds less than Lambda timeout)
	ctx, cancel := context.WithTimeout(ctx, validation.ProcessorTimeoutSeconds*time.Second)
	defer cancel()

	// Validate required fields
	if err := validation.ValidateUUID(event.TrackID, "trackId"); err != nil {
		return &Response{Indexed: false, Reason: err.Error()}, nil
	}

	if err := validation.ValidateUUID(event.UserID, "userId"); err != nil {
		return &Response{Indexed: false, Reason: err.Error()}, nil
	}

	// If search client not initialized, skip indexing
	if searchClient == nil {
		return &Response{Indexed: false, Reason: "search_disabled"}, nil
	}

	// Validate metadata is present
	if event.Metadata == nil {
		return &Response{Indexed: false, Reason: "missing_metadata"}, nil
	}

	// Build search document from metadata
	doc := search.Document{
		ID:        event.TrackID,
		UserID:    event.UserID,
		Title:     event.Metadata.Title,
		Artist:    event.Metadata.Artist,
		Album:     event.Metadata.Album,
		Genre:     event.Metadata.Genre,
		Year:      event.Metadata.Year,
		Duration:  event.Metadata.Duration,
		Filename:  event.S3Key,
		IndexedAt: time.Now(),
	}

	// Generate embedding if gateway is available (graceful degradation)
	if embeddingGateway != nil {
		text := service.ComposeEmbedTextFromMetadata(event.Metadata)
		if text != "" {
			embedding, err := embeddingGateway.GenerateEmbedding(ctx, text)
			if err != nil {
				fmt.Printf("Warning: embedding generation failed for track %s: %v\n", event.TrackID, err)
			} else {
				doc.Embedding = embedding

				// Store text embedding in S3 Vectors for hybrid search
				if searchVectorSvc != nil {
					metadata := map[string]string{
						"userId":  event.UserID,
						"trackId": event.TrackID,
					}
					if err := searchVectorSvc.PutVector(ctx, event.TrackID, embedding, metadata); err != nil {
						fmt.Printf("Warning: failed to store search embedding for track %s: %v\n", event.TrackID, err)
					}
				}
			}
		}
	}

	// Index the document
	resp, err := searchClient.Index(ctx, doc)
	if err != nil {
		return &Response{Indexed: false, Reason: fmt.Sprintf("index_failed: %v", err)}, nil
	}

	if !resp.Indexed {
		return &Response{Indexed: false, Reason: "index_rejected"}, nil
	}

	// Update step progress
	if event.UploadID != "" && repo != nil {
		if err := repo.UpdateUploadStep(ctx, event.UserID, event.UploadID, models.StepIndex, true); err != nil {
			fmt.Printf("Warning: failed to update step progress: %v\n", err)
		}
	}

	return &Response{Indexed: true}, nil
}

func main() {
	lambda.Start(handleRequest)
}
