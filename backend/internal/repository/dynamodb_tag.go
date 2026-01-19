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

// GetTag retrieves a tag by name
func (r *DynamoDBRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TAG#%s", tagName)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	if result.Item == nil {
		return nil, models.NewNotFoundError("Tag", tagName)
	}

	var item models.TagItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tag: %w", err)
	}

	return &item.Tag, nil
}

// CreateTag creates a new tag
func (r *DynamoDBRepository) CreateTag(ctx context.Context, tag *models.Tag) error {
	item := models.NewTagItem(*tag)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal tag: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// UpdateTag updates an existing tag
func (r *DynamoDBRepository) UpdateTag(ctx context.Context, tag *models.Tag) error {
	item := models.NewTagItem(*tag)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal tag: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

// DeleteTag deletes a tag and all its associations
func (r *DynamoDBRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	// Delete the tag
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TAG#%s", tagName)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// ListTags lists tags with filtering and pagination
func (r *DynamoDBRepository) ListTags(ctx context.Context, userID string, filter models.TagFilter) (*models.PaginatedResponse[models.Tag], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TAG#"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(filter.Limit)),
		ScanIndexForward:          aws.Bool(filter.SortOrder == "asc"),
	}

	if filter.LastKey != "" {
		startKey, err := decodeLastKey(filter.LastKey)
		if err != nil {
			return nil, fmt.Errorf("invalid last key: %w", err)
		}
		input.ExclusiveStartKey = startKey
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}

	tags := make([]models.Tag, 0, len(result.Items))
	for _, item := range result.Items {
		var tagItem models.TagItem
		if err := attributevalue.UnmarshalMap(item, &tagItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tag: %w", err)
		}
		tags = append(tags, tagItem.Tag)
	}

	response := &models.PaginatedResponse[models.Tag]{
		Items: tags,
		Pagination: models.Pagination{
			Limit: filter.Limit,
		},
	}

	if result.LastEvaluatedKey != nil {
		nextKey, err := encodeLastKey(result.LastEvaluatedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encode next key: %w", err)
		}
		response.Pagination.NextKey = nextKey
	}

	return response, nil
}

// AddTrackTag adds a tag to a track
func (r *DynamoDBRepository) AddTrackTag(ctx context.Context, tt *models.TrackTag) error {
	item := models.NewTrackTagItem(*tt)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal track tag: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to add track tag: %w", err)
	}

	return nil
}

// RemoveTrackTag removes a tag from a track
func (r *DynamoDBRepository) RemoveTrackTag(ctx context.Context, userID, trackID, tagName string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s#TRACK#%s", userID, trackID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TAG#%s", tagName)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to remove track tag: %w", err)
	}

	return nil
}
