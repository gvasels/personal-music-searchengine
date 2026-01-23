package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetSubscription retrieves a user's subscription
func (r *DynamoDBRepository) GetSubscription(ctx context.Context, userID string) (*models.Subscription, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "SUBSCRIPTION"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var item models.SubscriptionItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return &item.Subscription, nil
}

// PutSubscription creates or updates a subscription
func (r *DynamoDBRepository) PutSubscription(ctx context.Context, sub models.Subscription) error {
	item := models.NewSubscriptionItem(sub)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put subscription: %w", err)
	}

	return nil
}

// DeleteSubscription deletes a user's subscription
func (r *DynamoDBRepository) DeleteSubscription(ctx context.Context, userID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "SUBSCRIPTION"},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// UpdateUserTier updates the user's tier field
func (r *DynamoDBRepository) UpdateUserTier(ctx context.Context, userID string, tier models.SubscriptionTier) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
		UpdateExpression: aws.String("SET tier = :tier, updatedAt = :updatedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":tier":      &types.AttributeValueMemberS{Value: string(tier)},
			":updatedAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update user tier: %w", err)
	}

	return nil
}

// GetUserStorageUsage retrieves the user's storage usage
func (r *DynamoDBRepository) GetUserStorageUsage(ctx context.Context, userID string) (int64, error) {
	user, err := r.GetUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, nil
	}
	return user.StorageUsed, nil
}
