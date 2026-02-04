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

// CreateCrate creates a new crate
func (r *DynamoDBRepository) CreateCrate(ctx context.Context, crate models.Crate) error {
	item := models.NewCrateItem(crate)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal crate: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create crate: %w", err)
	}

	return nil
}

// GetCrate retrieves a crate by ID
func (r *DynamoDBRepository) GetCrate(ctx context.Context, userID, crateID string) (*models.Crate, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("CRATE#%s", crateID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get crate: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var item models.CrateItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal crate: %w", err)
	}

	return &item.Crate, nil
}

// UpdateCrate updates an existing crate
func (r *DynamoDBRepository) UpdateCrate(ctx context.Context, crate models.Crate) error {
	item := models.NewCrateItem(crate)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal crate: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update crate: %w", err)
	}

	return nil
}

// DeleteCrate deletes a crate
func (r *DynamoDBRepository) DeleteCrate(ctx context.Context, userID, crateID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("CRATE#%s", crateID)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete crate: %w", err)
	}

	return nil
}

// ListCrates retrieves all crates for a user
func (r *DynamoDBRepository) ListCrates(ctx context.Context, userID string, filter models.CrateFilter) ([]models.Crate, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			":skPrefix": &types.AttributeValueMemberS{Value: "CRATE#"},
		},
		Limit: aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list crates: %w", err)
	}

	var items []models.CrateItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal crates: %w", err)
	}

	crates := make([]models.Crate, len(items))
	for i, item := range items {
		crates[i] = item.Crate
	}

	return crates, nil
}

// CountUserCrates counts the number of crates a user has
func (r *DynamoDBRepository) CountUserCrates(ctx context.Context, userID string) (int, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			":skPrefix": &types.AttributeValueMemberS{Value: "CRATE#"},
		},
		Select: types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count crates: %w", err)
	}

	return int(result.Count), nil
}
