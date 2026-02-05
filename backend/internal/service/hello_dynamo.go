package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// HelloDynamoDBClient defines the DynamoDB operations needed by HelloDynamoDBRepository.
type HelloDynamoDBClient interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// HelloDynamoDBRepository implements HelloRepository using DynamoDB.
// It queries for seed tracks stored with PK=USER#seed-user and SK begins_with TRACK#.
type HelloDynamoDBRepository struct {
	client    HelloDynamoDBClient
	tableName string
}

// NewHelloDynamoDBRepository creates a new HelloDynamoDBRepository.
func NewHelloDynamoDBRepository(client HelloDynamoDBClient, tableName string) *HelloDynamoDBRepository {
	return &HelloDynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

// GetSeedTracks queries DynamoDB for all seed tracks (PK=USER#seed-user, SK begins_with TRACK#).
func (r *HelloDynamoDBRepository) GetSeedTracks(ctx context.Context) ([]HelloTrack, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value("USER#seed-user")),
		expression.Key("SK").BeginsWith("TRACK#"),
	)

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query seed tracks: %w", err)
	}

	tracks := make([]HelloTrack, 0, len(result.Items))
	for _, item := range result.Items {
		track := mapItemToHelloTrack(item)
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// mapItemToHelloTrack converts a DynamoDB item to a HelloTrack.
func mapItemToHelloTrack(item map[string]types.AttributeValue) HelloTrack {
	track := HelloTrack{}

	if v, ok := item["SK"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			// SK format is "TRACK#seed-t1" — extract the ID after "TRACK#"
			if len(s.Value) > 6 && s.Value[:6] == "TRACK#" {
				track.ID = s.Value[6:]
			} else {
				track.ID = s.Value
			}
		}
	}
	if v, ok := item["Title"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			track.Title = s.Value
		}
	}
	if v, ok := item["Artist"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			track.Artist = s.Value
		}
	}
	if v, ok := item["Album"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			track.Album = s.Value
		}
	}
	if v, ok := item["Genre"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			track.Genre = s.Value
		}
	}
	if v, ok := item["Year"]; ok {
		if n, ok := v.(*types.AttributeValueMemberN); ok {
			if parsed, err := strconv.Atoi(n.Value); err == nil {
				track.Year = parsed
			}
		}
	}
	if v, ok := item["Duration"]; ok {
		if n, ok := v.(*types.AttributeValueMemberN); ok {
			if parsed, err := strconv.Atoi(n.Value); err == nil {
				track.Duration = parsed
			}
		}
	}

	return track
}
