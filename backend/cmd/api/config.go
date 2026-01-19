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

	// CloudFront (optional)
	CloudFrontDomain     string
	CloudFrontKeyPairID  string
	CloudFrontPrivateKey string

	// Server (for local development)
	ServerPort string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		AWSRegion:            getEnvOrDefault("AWS_REGION", "us-east-1"),
		DynamoDBTableName:    os.Getenv("DYNAMODB_TABLE_NAME"),
		MediaBucketName:      os.Getenv("MEDIA_BUCKET"),
		StepFunctionsARN:     os.Getenv("STEP_FUNCTIONS_ARN"),
		CloudFrontDomain:     os.Getenv("CLOUDFRONT_DOMAIN"),
		CloudFrontKeyPairID:  os.Getenv("CLOUDFRONT_KEY_PAIR_ID"),
		CloudFrontPrivateKey: os.Getenv("CLOUDFRONT_PRIVATE_KEY"),
		ServerPort:           getEnvOrDefault("PORT", "8080"),
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
