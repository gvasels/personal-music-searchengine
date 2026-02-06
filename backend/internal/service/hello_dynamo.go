package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// HelloDynamoDBClient defines the interface for DynamoDB operations needed by HelloDynamoDBRepository
type HelloDynamoDBClient interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// HelloDynamoDBRepository implements HelloRepository using DynamoDB
type HelloDynamoDBRepository struct {
	client    HelloDynamoDBClient
	tableName string
}

// helloTrackItem represents the DynamoDB item structure for hello seed tracks
type helloTrackItem struct {
	PK       string `dynamodbav:"PK"`
	SK       string `dynamodbav:"SK"`
	Title    string `dynamodbav:"Title"`
	Artist   string `dynamodbav:"Artist"`
	Album    string `dynamodbav:"Album"`
	Genre    string `dynamodbav:"Genre"`
	Year     int    `dynamodbav:"Year"`
	Duration int    `dynamodbav:"Duration"`
}

// NewHelloDynamoDBRepository creates a new HelloDynamoDBRepository with the given client and table name
func NewHelloDynamoDBRepository(client HelloDynamoDBClient, tableName string) *HelloDynamoDBRepository {
	return &HelloDynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

// GetSeedTracks retrieves all seed tracks from DynamoDB
// Queries with PK=USER#seed-user and SK begins_with TRACK#
func (r *HelloDynamoDBRepository) GetSeedTracks(ctx context.Context) ([]HelloTrack, error) {
	// Build key condition: PK = USER#seed-user AND SK begins_with TRACK#
	keyCond := expression.Key("PK").Equal(expression.Value("USER#seed-user")).
		And(expression.Key("SK").BeginsWith("TRACK#"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query seed tracks: %w", err)
	}

	tracks := make([]HelloTrack, 0, len(result.Items))
	for _, item := range result.Items {
		var trackItem helloTrackItem
		if err := attributevalue.UnmarshalMap(item, &trackItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal track: %w", err)
		}

		// Extract ID from SK by stripping "TRACK#" prefix
		id := strings.TrimPrefix(trackItem.SK, "TRACK#")

		tracks = append(tracks, HelloTrack{
			ID:       id,
			Title:    trackItem.Title,
			Artist:   trackItem.Artist,
			Album:    trackItem.Album,
			Genre:    trackItem.Genre,
			Year:     trackItem.Year,
			Duration: trackItem.Duration,
		})
	}

	return tracks, nil
}
