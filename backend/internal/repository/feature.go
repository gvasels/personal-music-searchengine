package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetFeatureFlag retrieves a feature flag by key
func (r *DynamoDBRepository) GetFeatureFlag(ctx context.Context, key models.FeatureKey) (*models.FeatureFlag, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "FEATURE"},
			"SK": &types.AttributeValueMemberS{Value: "FLAG#" + string(key)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get feature flag: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var item models.FeatureFlagItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feature flag: %w", err)
	}

	return &item.FeatureFlag, nil
}

// ListFeatureFlags retrieves all feature flags
func (r *DynamoDBRepository) ListFeatureFlags(ctx context.Context) ([]models.FeatureFlag, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: "FEATURE"},
			":skPrefix": &types.AttributeValueMemberS{Value: "FLAG#"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list feature flags: %w", err)
	}

	var items []models.FeatureFlagItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feature flags: %w", err)
	}

	flags := make([]models.FeatureFlag, len(items))
	for i, item := range items {
		flags[i] = item.FeatureFlag
	}

	return flags, nil
}

// PutFeatureFlag creates or updates a feature flag
func (r *DynamoDBRepository) PutFeatureFlag(ctx context.Context, flag models.FeatureFlag) error {
	item := models.NewFeatureFlagItem(flag)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal feature flag: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put feature flag: %w", err)
	}

	return nil
}

// DeleteFeatureFlag deletes a feature flag
func (r *DynamoDBRepository) DeleteFeatureFlag(ctx context.Context, key models.FeatureKey) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "FEATURE"},
			"SK": &types.AttributeValueMemberS{Value: "FLAG#" + string(key)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete feature flag: %w", err)
	}

	return nil
}

// GetUserFeatureOverride retrieves a user's override for a specific feature
func (r *DynamoDBRepository) GetUserFeatureOverride(ctx context.Context, userID string, key models.FeatureKey) (*models.UserFeatureOverride, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "FEATURE_OVERRIDE#" + string(key)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user feature override: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var item models.UserFeatureOverrideItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user feature override: %w", err)
	}

	return &item.UserFeatureOverride, nil
}

// ListUserFeatureOverrides retrieves all overrides for a user
func (r *DynamoDBRepository) ListUserFeatureOverrides(ctx context.Context, userID string) ([]models.UserFeatureOverride, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: "USER#" + userID},
			":skPrefix": &types.AttributeValueMemberS{Value: "FEATURE_OVERRIDE#"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list user feature overrides: %w", err)
	}

	var items []models.UserFeatureOverrideItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user feature overrides: %w", err)
	}

	overrides := make([]models.UserFeatureOverride, len(items))
	for i, item := range items {
		overrides[i] = item.UserFeatureOverride
	}

	return overrides, nil
}

// PutUserFeatureOverride creates or updates a user's feature override
func (r *DynamoDBRepository) PutUserFeatureOverride(ctx context.Context, override models.UserFeatureOverride) error {
	item := models.NewUserFeatureOverrideItem(override)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal user feature override: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put user feature override: %w", err)
	}

	return nil
}

// DeleteUserFeatureOverride deletes a user's feature override
func (r *DynamoDBRepository) DeleteUserFeatureOverride(ctx context.Context, userID string, key models.FeatureKey) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "FEATURE_OVERRIDE#" + string(key)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete user feature override: %w", err)
	}

	return nil
}

// SeedDefaultFeatureFlags creates default feature flags if they don't exist
func (r *DynamoDBRepository) SeedDefaultFeatureFlags(ctx context.Context) error {
	defaults := models.DefaultFeatureFlags()

	for _, flag := range defaults {
		existing, err := r.GetFeatureFlag(ctx, flag.Key)
		if err != nil {
			return err
		}
		if existing == nil {
			if err := r.PutFeatureFlag(ctx, flag); err != nil {
				return err
			}
		}
	}

	return nil
}
