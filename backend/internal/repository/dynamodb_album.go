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
	"github.com/google/uuid"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// GetAlbum retrieves an album by ID
func (r *DynamoDBRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ALBUM#%s", albumID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	if result.Item == nil {
		return nil, models.NewNotFoundError("Album", albumID)
	}

	var item models.AlbumItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album: %w", err)
	}

	return &item.Album, nil
}

// CreateAlbum creates a new album
func (r *DynamoDBRepository) CreateAlbum(ctx context.Context, album *models.Album) error {
	item := models.NewAlbumItem(*album)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal album: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create album: %w", err)
	}

	return nil
}

// UpdateAlbum updates an existing album
func (r *DynamoDBRepository) UpdateAlbum(ctx context.Context, album *models.Album) error {
	item := models.NewAlbumItem(*album)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal album: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update album: %w", err)
	}

	return nil
}

// DeleteAlbum deletes an album
func (r *DynamoDBRepository) DeleteAlbum(ctx context.Context, userID, albumID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ALBUM#%s", albumID)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}

	return nil
}

// ListAlbums lists albums with filtering and pagination
func (r *DynamoDBRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResponse[models.Album], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("ALBUM#"))

	// Build filter expression
	var filterExpr expression.ConditionBuilder
	hasFilter := false

	if filter.Artist != "" {
		filterExpr = expression.Name("artist").Equal(expression.Value(filter.Artist))
		hasFilter = true
	}
	if filter.Genre != "" {
		if hasFilter {
			filterExpr = filterExpr.And(expression.Name("genre").Equal(expression.Value(filter.Genre)))
		} else {
			filterExpr = expression.Name("genre").Equal(expression.Value(filter.Genre))
			hasFilter = true
		}
	}
	if filter.Year > 0 {
		if hasFilter {
			filterExpr = filterExpr.And(expression.Name("year").Equal(expression.Value(filter.Year)))
		} else {
			filterExpr = expression.Name("year").Equal(expression.Value(filter.Year))
			hasFilter = true
		}
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
		return nil, fmt.Errorf("failed to query albums: %w", err)
	}

	albums := make([]models.Album, 0, len(result.Items))
	for _, item := range result.Items {
		var albumItem models.AlbumItem
		if err := attributevalue.UnmarshalMap(item, &albumItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal album: %w", err)
		}
		albums = append(albums, albumItem.Album)
	}

	response := &models.PaginatedResponse[models.Album]{
		Items: albums,
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

// ListAlbumsByArtist lists all albums by an artist
func (r *DynamoDBRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#ARTIST#%s", userID, artist))).
		And(expression.Key("GSI1SK").BeginsWith("ALBUM#"))

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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query albums by artist: %w", err)
	}

	albums := make([]models.Album, 0, len(result.Items))
	for _, item := range result.Items {
		var albumItem models.AlbumItem
		if err := attributevalue.UnmarshalMap(item, &albumItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal album: %w", err)
		}
		albums = append(albums, albumItem.Album)
	}

	return albums, nil
}

// GetOrCreateAlbum gets an existing album or creates a new one
func (r *DynamoDBRepository) GetOrCreateAlbum(ctx context.Context, userID, title, artist string, year int) (*models.Album, error) {
	// Try to find existing album
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("ALBUM#"))
	filterExpr := expression.Name("title").Equal(expression.Value(title)).
		And(expression.Name("artist").Equal(expression.Value(artist)))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		WithFilter(filterExpr).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query album: %w", err)
	}

	if len(result.Items) > 0 {
		var albumItem models.AlbumItem
		if err := attributevalue.UnmarshalMap(result.Items[0], &albumItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal album: %w", err)
		}
		return &albumItem.Album, nil
	}

	// Create new album
	now := time.Now()
	album := &models.Album{
		ID:     uuid.New().String(),
		UserID: userID,
		Title:  title,
		Artist: artist,
		Year:   year,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := r.CreateAlbum(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}
