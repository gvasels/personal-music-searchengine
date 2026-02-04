//go:build integration

// Package testutil provides utilities for integration testing with LocalStack.
package testutil

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DefaultLocalStackEndpoint is the default LocalStack endpoint.
const DefaultLocalStackEndpoint = "http://localhost:4566"

// DefaultTableName is the default DynamoDB table name for tests.
const DefaultTableName = "MusicLibrary"

// DefaultBucketName is the default S3 bucket name for tests.
const DefaultBucketName = "music-library-local-media"

// TestContext holds LocalStack clients and configuration for integration tests.
type TestContext struct {
	DynamoDB   *dynamodb.Client
	S3         *s3.Client
	Cognito    *cognitoidentityprovider.Client
	TableName  string
	BucketName string
	UserPoolID string
	ClientID   string
	Endpoint   string
	Region     string

	// Cleanup tracking
	cleanupItems []cleanupItem
}

type cleanupItem struct {
	itemType string
	pk       string
	sk       string
}

// LocalStackConfig holds configuration for LocalStack connection.
type LocalStackConfig struct {
	Endpoint    string
	Region      string
	TableName   string
	BucketName  string
	UserPoolID  string
	ClientID    string
}

// DefaultConfig returns the default LocalStack configuration.
func DefaultConfig() LocalStackConfig {
	return LocalStackConfig{
		Endpoint:   getEnvOrDefault("LOCALSTACK_ENDPOINT", DefaultLocalStackEndpoint),
		Region:     getEnvOrDefault("AWS_REGION", "us-east-1"),
		TableName:  getEnvOrDefault("DYNAMODB_TABLE_NAME", DefaultTableName),
		BucketName: getEnvOrDefault("MEDIA_BUCKET", DefaultBucketName),
		UserPoolID: os.Getenv("COGNITO_USER_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
	}
}

// IsLocalStackRunning checks if LocalStack is available at the configured endpoint.
func IsLocalStackRunning() bool {
	endpoint := getEnvOrDefault("LOCALSTACK_ENDPOINT", DefaultLocalStackEndpoint)
	healthURL := fmt.Sprintf("%s/_localstack/health", endpoint)

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// SetupLocalStack initializes clients connected to LocalStack.
// Returns TestContext and a cleanup function.
// Skips the test if LocalStack is not running.
func SetupLocalStack(t *testing.T) (*TestContext, func()) {
	t.Helper()

	if !IsLocalStackRunning() {
		t.Skip("LocalStack not running. Start with: docker-compose -f docker/docker-compose.yml up -d")
	}

	cfg := DefaultConfig()
	ctx := context.Background()

	// Create AWS config for LocalStack
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create clients with LocalStack endpoint
	dynamoClient := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})

	cognitoClient := cognitoidentityprovider.NewFromConfig(awsCfg, func(o *cognitoidentityprovider.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	// Try to get user pool ID if not configured
	userPoolID := cfg.UserPoolID
	clientID := cfg.ClientID
	if userPoolID == "" {
		userPoolID, clientID = discoverCognitoConfig(ctx, cognitoClient)
	}

	tc := &TestContext{
		DynamoDB:     dynamoClient,
		S3:           s3Client,
		Cognito:      cognitoClient,
		TableName:    cfg.TableName,
		BucketName:   cfg.BucketName,
		UserPoolID:   userPoolID,
		ClientID:     clientID,
		Endpoint:     cfg.Endpoint,
		Region:       cfg.Region,
		cleanupItems: make([]cleanupItem, 0),
	}

	cleanup := func() {
		tc.runCleanup(t)
	}

	return tc, cleanup
}

// discoverCognitoConfig attempts to find the user pool and client ID from LocalStack.
func discoverCognitoConfig(ctx context.Context, client *cognitoidentityprovider.Client) (string, string) {
	pools, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int32(10),
	})
	if err != nil || len(pools.UserPools) == 0 {
		return "", ""
	}

	// Find the music-library pool
	var poolID string
	for _, pool := range pools.UserPools {
		if pool.Name != nil && *pool.Name == "music-library-local-pool" {
			poolID = *pool.Id
			break
		}
	}

	if poolID == "" && len(pools.UserPools) > 0 {
		poolID = *pools.UserPools[0].Id
	}

	if poolID == "" {
		return "", ""
	}

	// Get client ID
	clients, err := client.ListUserPoolClients(ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: aws.String(poolID),
		MaxResults: aws.Int32(10),
	})
	if err != nil || len(clients.UserPoolClients) == 0 {
		return poolID, ""
	}

	return poolID, *clients.UserPoolClients[0].ClientId
}

// RegisterCleanup registers an item for cleanup at test end.
func (tc *TestContext) RegisterCleanup(itemType, pk, sk string) {
	tc.cleanupItems = append(tc.cleanupItems, cleanupItem{
		itemType: itemType,
		pk:       pk,
		sk:       sk,
	})
}

// runCleanup removes all registered items (DynamoDB and S3).
func (tc *TestContext) runCleanup(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	for _, item := range tc.cleanupItems {
		if item.itemType == "s3" {
			err := tc.deleteS3Object(ctx, item.pk)
			if err != nil {
				t.Logf("Cleanup warning: failed to delete S3 object %s: %v", item.pk, err)
			}
		} else {
			err := tc.deleteItem(ctx, item.pk, item.sk)
			if err != nil {
				t.Logf("Cleanup warning: failed to delete %s/%s: %v", item.pk, item.sk, err)
			}
		}
	}
}

// deleteItem deletes an item from DynamoDB.
func (tc *TestContext) deleteItem(ctx context.Context, pk, sk string) error {
	_, err := tc.DynamoDB.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tc.TableName),
		Key: map[string]dynamodbtypes.AttributeValue{
			"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"SK": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		},
	})
	return err
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
