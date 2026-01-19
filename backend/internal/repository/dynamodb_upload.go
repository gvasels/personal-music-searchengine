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

// GetUpload retrieves an upload by ID
func (r *DynamoDBRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("UPLOAD#%s", uploadID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get upload: %w", err)
	}

	if result.Item == nil {
		return nil, models.NewNotFoundError("Upload", uploadID)
	}

	var item models.UploadItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal upload: %w", err)
	}

	return &item.Upload, nil
}

// CreateUpload creates a new upload record
func (r *DynamoDBRepository) CreateUpload(ctx context.Context, upload *models.Upload) error {
	item := models.NewUploadItem(*upload)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal upload: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create upload: %w", err)
	}

	return nil
}

// UpdateUpload updates an existing upload record
func (r *DynamoDBRepository) UpdateUpload(ctx context.Context, upload *models.Upload) error {
	item := models.NewUploadItem(*upload)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal upload: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update upload: %w", err)
	}

	return nil
}

// ListUploads lists uploads with filtering and pagination
func (r *DynamoDBRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.Upload], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("UPLOAD#"))

	var filterExpr expression.ConditionBuilder
	hasFilter := false

	if filter.Status != "" {
		filterExpr = expression.Name("status").Equal(expression.Value(filter.Status))
		hasFilter = true
	}

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	if hasFilter {
		builder = builder.WithFilter(filterExpr)
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
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
		return nil, fmt.Errorf("failed to query uploads: %w", err)
	}

	uploads := make([]models.Upload, 0, len(result.Items))
	for _, item := range result.Items {
		var uploadItem models.UploadItem
		if err := attributevalue.UnmarshalMap(item, &uploadItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal upload: %w", err)
		}
		uploads = append(uploads, uploadItem.Upload)
	}

	response := &models.PaginatedResponse[models.Upload]{
		Items: uploads,
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
