package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// CreateArtistWatch creates an artist watch entry. Idempotent — duplicate creates do not error.
func (r *DynamoDBRepository) CreateArtistWatch(ctx context.Context, watch models.ArtistWatch) error {
	item := models.NewArtistWatchItem(watch)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal artist watch: %w", err)
	}

	// Plain PutItem without ConditionExpression — idempotent (overwrites on duplicate)
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to create artist watch: %w", err)
	}

	return nil
}

// DeleteArtistWatch removes an artist watch entry.
// Returns ErrNotFound if the watch does not exist.
func (r *DynamoDBRepository) DeleteArtistWatch(ctx context.Context, userID, artistName string) error {
	normalized := models.NormalizeArtistName(artistName)
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.GetArtistWatchPK(userID)},
			"SK": &types.AttributeValueMemberS{Value: models.GetArtistWatchSK(normalized)},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if ok := isConditionalCheckFailed(err, &condErr); ok {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete artist watch: %w", err)
	}

	return nil
}

// GetArtistWatch retrieves an artist watch by user ID and artist name.
// Returns ErrNotFound if the watch does not exist.
func (r *DynamoDBRepository) GetArtistWatch(ctx context.Context, userID, artistName string) (*models.ArtistWatch, error) {
	normalized := models.NormalizeArtistName(artistName)
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.GetArtistWatchPK(userID)},
			"SK": &types.AttributeValueMemberS{Value: models.GetArtistWatchSK(normalized)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get artist watch: %w", err)
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	var item models.ArtistWatchItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artist watch: %w", err)
	}

	return &item.ArtistWatch, nil
}

// ListWatchedArtists returns a paginated list of artists watched by a user.
// Queries PK=USER#{userID} with SK beginning with ARTIST_WATCH#.
func (r *DynamoDBRepository) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*PaginatedResult[models.ArtistWatch], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(models.GetArtistWatchPK(userID))).
		And(expression.Key("SK").BeginsWith("ARTIST_WATCH#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
	}

	if cursor != "" {
		startKey, err := decodeCursor(cursor)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		input.ExclusiveStartKey = startKey
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list watched artists: %w", err)
	}

	watches := make([]models.ArtistWatch, 0, len(result.Items))
	for _, item := range result.Items {
		var watchItem models.ArtistWatchItem
		if err := attributevalue.UnmarshalMap(item, &watchItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal artist watch: %w", err)
		}
		watches = append(watches, watchItem.ArtistWatch)
	}

	var nextCursor string
	if result.LastEvaluatedKey != nil {
		nextCursor, err = encodeCursor(result.LastEvaluatedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encode cursor: %w", err)
		}
	}

	return &PaginatedResult[models.ArtistWatch]{
		Items:      watches,
		NextCursor: nextCursor,
		HasMore:    result.LastEvaluatedKey != nil,
	}, nil
}
