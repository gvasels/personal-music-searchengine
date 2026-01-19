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

// GetUser retrieves a user by ID
func (r *DynamoDBRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if result.Item == nil {
		return nil, models.NewNotFoundError("User", userID)
	}

	var item models.UserItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &item.User, nil
}

// CreateUser creates a new user
func (r *DynamoDBRepository) CreateUser(ctx context.Context, user *models.User) error {
	item := models.NewUserItem(*user)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user
func (r *DynamoDBRepository) UpdateUser(ctx context.Context, user *models.User) error {
	item := models.NewUserItem(*user)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
