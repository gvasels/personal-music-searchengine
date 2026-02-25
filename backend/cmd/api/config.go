package main

import (
	"fmt"
	"os"
)

// Config holds application configuration loaded from environment variables
type Config struct {
	// AWS
	AWSRegion string

	// DynamoDB
	DynamoDBTableName string

	// S3
	MediaBucketName string

	// Step Functions
	StepFunctionsARN string
	AudioPipelineARN string

	// Nixiesearch
	NixiesearchFunctionName string

	// CloudFront (optional)
	CloudFrontDomain     string
	CloudFrontKeyPairID  string
	CloudFrontPrivateKey string

	// Cognito (for admin operations)
	CognitoUserPoolID string

	// S3 Vectors
	VectorBucketName       string
	VectorIndexName        string
	SearchVectorIndexName  string

	// S3 flat vector backup (for dual-write)
	VectorS3Bucket string
	VectorS3Prefix string

	// Embedding Gateway
	EmbeddingGatewayURL    string
	EmbeddingGatewaySecret string

	// Server (for local development)
	ServerPort string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		AWSRegion:               getEnvOrDefault("AWS_REGION", "us-east-1"),
		DynamoDBTableName:       os.Getenv("DYNAMODB_TABLE_NAME"),
		MediaBucketName:         os.Getenv("MEDIA_BUCKET"),
		StepFunctionsARN:        os.Getenv("STEP_FUNCTIONS_ARN"),
		AudioPipelineARN:        os.Getenv("AUDIO_PIPELINE_ARN"),
		NixiesearchFunctionName: os.Getenv("NIXIESEARCH_FUNCTION_NAME"),
		CloudFrontDomain:        os.Getenv("CLOUDFRONT_DOMAIN"),
		CloudFrontKeyPairID:     os.Getenv("CLOUDFRONT_KEY_PAIR_ID"),
		CloudFrontPrivateKey:    os.Getenv("CLOUDFRONT_PRIVATE_KEY"),
		CognitoUserPoolID:       os.Getenv("COGNITO_USER_POOL_ID"),
		VectorBucketName:        os.Getenv("VECTOR_BUCKET_NAME"),
		VectorIndexName:         getEnvOrDefault("VECTOR_INDEX_NAME", "media-embeddings"),
		SearchVectorIndexName:   getEnvOrDefault("SEARCH_VECTOR_INDEX_NAME", "search-embeddings"),
		VectorS3Bucket:          os.Getenv("VECTOR_S3_BUCKET"),
		VectorS3Prefix:          getEnvOrDefault("VECTOR_S3_PREFIX", "vectors"),
		EmbeddingGatewayURL:    os.Getenv("EMBEDDING_GATEWAY_URL"),
		EmbeddingGatewaySecret: os.Getenv("EMBEDDING_GATEWAY_SECRET"),
		ServerPort:              getEnvOrDefault("PORT", "8080"),
	}

	// Validate required fields
	if cfg.DynamoDBTableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE_NAME environment variable is required")
	}
	if cfg.MediaBucketName == "" {
		return nil, fmt.Errorf("MEDIA_BUCKET environment variable is required")
	}

	return cfg, nil
}

// IsLambda returns true if running in AWS Lambda environment
func IsLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

// getEnvOrDefault returns the environment variable value or a default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
