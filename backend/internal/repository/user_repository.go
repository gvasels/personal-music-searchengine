package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetUserByEmail retrieves a user by their email address using GSI1
func (r *DynamoDBRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value("EMAIL#" + email))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, ErrUserNotFound
	}

	var item models.UserItem
	if err := attributevalue.UnmarshalMap(result.Items[0], &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &item.User, nil
}

// GetUserByCognitoID retrieves a user by their Cognito sub ID
// Since we use Cognito sub as the user ID, this is equivalent to GetUser
func (r *DynamoDBRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return r.GetUser(ctx, cognitoID)
}

// SearchUsersByEmail searches for users by email prefix using GSI1
func (r *DynamoDBRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]UserSearchResult, string, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Build key condition for begins_with on GSI1PK
	keyCondition := expression.Key("GSI1PK").BeginsWith("EMAIL#" + strings.ToLower(emailPrefix))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, "", fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
	}

	// Handle pagination
	if cursor != "" {
		decodedCursor, err := models.DecodeCursor(cursor)
		if err != nil {
			return nil, "", ErrInvalidCursor
		}
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"PK":     &types.AttributeValueMemberS{Value: decodedCursor.PK},
			"SK":     &types.AttributeValueMemberS{Value: decodedCursor.SK},
			"GSI1PK": &types.AttributeValueMemberS{Value: decodedCursor.GSI1PK},
			"GSI1SK": &types.AttributeValueMemberS{Value: decodedCursor.GSI1SK},
		}
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to search users: %w", err)
	}

	var results []UserSearchResult
	for _, item := range result.Items {
		var userItem models.UserItem
		if err := attributevalue.UnmarshalMap(item, &userItem); err != nil {
			continue // Skip malformed items
		}

		results = append(results, UserSearchResult{
			ID:          userItem.ID,
			Email:       userItem.Email,
			DisplayName: userItem.DisplayName,
			Role:        userItem.Role,
			Disabled:    userItem.Disabled,
			CreatedAt:   userItem.CreatedAt,
		})
	}

	var nextCursor string
	if result.LastEvaluatedKey != nil {
		pk, _ := result.LastEvaluatedKey["PK"].(*types.AttributeValueMemberS)
		sk, _ := result.LastEvaluatedKey["SK"].(*types.AttributeValueMemberS)
		gsi1pk, _ := result.LastEvaluatedKey["GSI1PK"].(*types.AttributeValueMemberS)
		gsi1sk, _ := result.LastEvaluatedKey["GSI1SK"].(*types.AttributeValueMemberS)

		if pk != nil && sk != nil {
			cursor := models.PaginationCursor{
				PK: pk.Value,
				SK: sk.Value,
			}
			if gsi1pk != nil {
				cursor.GSI1PK = gsi1pk.Value
			}
			if gsi1sk != nil {
				cursor.GSI1SK = gsi1sk.Value
			}
			nextCursor = models.EncodeCursor(cursor)
		}
	}

	return results, nextCursor, nil
}

// GetUserSettings retrieves just the settings for a user
func (r *DynamoDBRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	user, err := r.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user.Settings, nil
}

// UpdateUserSettings performs a partial update of user settings
func (r *DynamoDBRepository) UpdateUserSettings(ctx context.Context, userID string, update *UserSettingsUpdate) (*models.UserSettings, error) {
	// First get current user to merge settings
	user, err := r.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if update.Notifications != nil {
		user.Settings.Notifications = *update.Notifications
	}
	if update.Privacy != nil {
		user.Settings.Privacy = *update.Privacy
	}
	if update.Player != nil {
		user.Settings.Player = *update.Player
	}
	if update.Library != nil {
		user.Settings.Library = *update.Library
	}

	// Validate the updated settings
	if err := user.Settings.Validate(); err != nil {
		return nil, fmt.Errorf("invalid settings: %w", err)
	}

	// Update the user
	user.UpdatedAt = time.Now()

	// Build update expression
	updateExpr := expression.Set(
		expression.Name("settings"),
		expression.Value(user.Settings),
	).Set(
		expression.Name("updatedAt"),
		expression.Value(user.UpdatedAt),
	)

	expr, err := expression.NewBuilder().WithUpdate(updateExpr).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return &user.Settings, nil
}
