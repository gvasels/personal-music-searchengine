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

// CreateFollow creates a follow relationship
func (r *DynamoDBRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	if err := follow.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	follow.CreatedAt = time.Now()
	item := models.NewFollowItem(follow)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal follow: %w", err)
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
		return fmt.Errorf("failed to create follow: %w", err)
	}

	return nil
}

// DeleteFollow removes a follow relationship
func (r *DynamoDBRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.GetFollowingPK(followerID)},
			"SK": &types.AttributeValueMemberS{Value: models.GetFollowingSK(followedID)},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if ok := isConditionalCheckFailed(err, &condErr); ok {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete follow: %w", err)
	}

	return nil
}

// GetFollow retrieves a follow relationship if it exists
func (r *DynamoDBRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: models.GetFollowingPK(followerID)},
			"SK": &types.AttributeValueMemberS{Value: models.GetFollowingSK(followedID)},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get follow: %w", err)
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	var item models.FollowItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal follow: %w", err)
	}

	return &item.Follow, nil
}

// ListFollowers lists users who follow a given user (via GSI1)
func (r *DynamoDBRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*PaginatedResult[models.Follow], error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(models.GetFollowersGSI1PK(userID)))

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
		return nil, fmt.Errorf("failed to list followers: %w", err)
	}

	follows := make([]models.Follow, 0, len(result.Items))
	for _, item := range result.Items {
		var followItem models.FollowItem
		if err := attributevalue.UnmarshalMap(item, &followItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal follow: %w", err)
		}
		follows = append(follows, followItem.Follow)
	}

	var nextCursor string
	if result.LastEvaluatedKey != nil {
		nextCursor, err = encodeCursor(result.LastEvaluatedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encode cursor: %w", err)
		}
	}

	return &PaginatedResult[models.Follow]{
		Items:      follows,
		NextCursor: nextCursor,
		HasMore:    result.LastEvaluatedKey != nil,
	}, nil
}

// ListFollowing lists users that a given user follows
func (r *DynamoDBRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*PaginatedResult[models.Follow], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(models.GetFollowingPK(userID))).
		And(expression.Key("SK").BeginsWith("FOLLOWING#"))

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
		return nil, fmt.Errorf("failed to list following: %w", err)
	}

	follows := make([]models.Follow, 0, len(result.Items))
	for _, item := range result.Items {
		var followItem models.FollowItem
		if err := attributevalue.UnmarshalMap(item, &followItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal follow: %w", err)
		}
		follows = append(follows, followItem.Follow)
	}

	var nextCursor string
	if result.LastEvaluatedKey != nil {
		nextCursor, err = encodeCursor(result.LastEvaluatedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encode cursor: %w", err)
		}
	}

	return &PaginatedResult[models.Follow]{
		Items:      follows,
		NextCursor: nextCursor,
		HasMore:    result.LastEvaluatedKey != nil,
	}, nil
}

// IncrementUserFollowingCount increments or decrements the following count on a user
func (r *DynamoDBRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	update := expression.Add(expression.Name("followingCount"), expression.Value(delta)).
		Set(expression.Name("updatedAt"), expression.Value(time.Now().Format(time.RFC3339)))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
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
		return fmt.Errorf("failed to update following count: %w", err)
	}

	return nil
}
