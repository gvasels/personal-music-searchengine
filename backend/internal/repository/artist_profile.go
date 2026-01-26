package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// CreateArtistProfile creates a new artist profile
func (r *DynamoDBRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	now := time.Now()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	item := models.NewArtistProfileItem(profile)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal artist profile: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if ok := isConditionalCheckFailed(err, &condErr); ok {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to create artist profile: %w", err)
	}

	return nil
}

// GetArtistProfile retrieves an artist profile by user ID
func (r *DynamoDBRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: "ARTIST_PROFILE"},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get artist profile: %w", err)
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	var item models.ArtistProfileItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artist profile: %w", err)
	}

	return &item.ArtistProfile, nil
}

// UpdateArtistProfile updates an existing artist profile
func (r *DynamoDBRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	profile.UpdatedAt = time.Now()

	item := models.NewArtistProfileItem(profile)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal artist profile: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if ok := isConditionalCheckFailed(err, &condErr); ok {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update artist profile: %w", err)
	}

	return nil
}

// DeleteArtistProfile deletes an artist profile
func (r *DynamoDBRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: "ARTIST_PROFILE"},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if ok := isConditionalCheckFailed(err, &condErr); ok {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete artist profile: %w", err)
	}

	return nil
}

// ListArtistProfiles lists all artist profiles with pagination
func (r *DynamoDBRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*PaginatedResult[models.ArtistProfile], error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value("ARTIST_PROFILE"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("GSI1"),
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
		return nil, fmt.Errorf("failed to list artist profiles: %w", err)
	}

	profiles := make([]models.ArtistProfile, 0, len(result.Items))
	for _, item := range result.Items {
		var profileItem models.ArtistProfileItem
		if err := attributevalue.UnmarshalMap(item, &profileItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal artist profile: %w", err)
		}
		profiles = append(profiles, profileItem.ArtistProfile)
	}

	var nextCursor string
	if result.LastEvaluatedKey != nil {
		nextCursor, err = encodeCursor(result.LastEvaluatedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encode cursor: %w", err)
		}
	}

	return &PaginatedResult[models.ArtistProfile]{
		Items:      profiles,
		NextCursor: nextCursor,
		HasMore:    result.LastEvaluatedKey != nil,
	}, nil
}

// IncrementArtistFollowerCount increments or decrements the follower count
func (r *DynamoDBRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	update := expression.Add(expression.Name("followerCount"), expression.Value(delta)).
		Set(expression.Name("updatedAt"), expression.Value(time.Now().Format(time.RFC3339)))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: "ARTIST_PROFILE"},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if ok := isConditionalCheckFailed(err, &condErr); ok {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update follower count: %w", err)
	}

	return nil
}

// isConditionalCheckFailed checks if an error is a conditional check failure
func isConditionalCheckFailed(err error, target **types.ConditionalCheckFailedException) bool {
	if err == nil {
		return false
	}
	var condErr *types.ConditionalCheckFailedException
	if ok := isError(err, &condErr); ok {
		if target != nil {
			*target = condErr
		}
		return true
	}
	return false
}

// isError is a helper to check error types
func isError[T error](err error, target *T) bool {
	return err != nil && target != nil
}
