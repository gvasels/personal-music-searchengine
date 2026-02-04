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

// ============================================================================
// Artist Operations
// ============================================================================

// CreateArtist creates a new artist in DynamoDB
func (r *DynamoDBRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	artist.CreatedAt = time.Now()
	artist.UpdatedAt = artist.CreatedAt

	item := models.NewArtistItem(artist)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal artist: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create artist: %w", err)
	}

	return nil
}

// GetArtist retrieves an artist by ID
func (r *DynamoDBRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ARTIST#%s", artistID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	var item models.ArtistItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artist: %w", err)
	}

	return &item.Artist, nil
}

// GetArtistByName retrieves artists by name using GSI1
func (r *DynamoDBRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#ARTIST", userID))).
		And(expression.Key("GSI1SK").Equal(expression.Value(name)))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
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
		return nil, fmt.Errorf("failed to query artists by name: %w", err)
	}

	var items []models.ArtistItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
	}

	artists := make([]*models.Artist, 0, len(items))
	for _, item := range items {
		artist := item.Artist
		artists = append(artists, &artist)
	}

	return artists, nil
}

// ListArtists returns paginated list of artists for a user
func (r *DynamoDBRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*PaginatedResult[models.Artist], error) {
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}

	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("ARTIST#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)

	// Filter only active artists by default
	filterExpr := expression.Name("isActive").Equal(expression.Value(true))
	builder = builder.WithFilter(filterExpr)

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
		Limit:                     aws.Int32(int32(limit + 1)),
	}

	// Handle pagination cursor
	if filter.LastKey != "" {
		cursor, err := models.DecodeCursor(filter.LastKey)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		input.ExclusiveStartKey = cursorToAttributeValue(cursor)
	}

	// Handle sort order
	if filter.SortOrder == "desc" {
		input.ScanIndexForward = aws.Bool(false)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query artists: %w", err)
	}

	var items []models.ArtistItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
	}

	artists := make([]models.Artist, 0, len(items))
	for _, item := range items {
		artists = append(artists, item.Artist)
	}

	// Determine if there are more results
	hasMore := len(artists) > limit
	if hasMore {
		artists = artists[:limit]
	}

	// Build next cursor
	var nextCursor string
	if hasMore && len(artists) > 0 {
		lastArtist := artists[len(artists)-1]
		cursor := models.NewPaginationCursor(
			fmt.Sprintf("USER#%s", userID),
			fmt.Sprintf("ARTIST#%s", lastArtist.ID),
		)
		nextCursor = models.EncodeCursor(cursor)
	}

	return &PaginatedResult[models.Artist]{
		Items:      artists,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// UpdateArtist updates an existing artist
func (r *DynamoDBRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	artist.UpdatedAt = time.Now()

	item := models.NewArtistItem(artist)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal artist: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update artist: %w", err)
	}

	return nil
}

// DeleteArtist soft-deletes an artist by setting IsActive to false
func (r *DynamoDBRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	update := expression.Set(
		expression.Name("isActive"), expression.Value(false),
	).Set(
		expression.Name("updatedAt"), expression.Value(time.Now().Format(time.RFC3339)),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ARTIST#%s", artistID)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to delete artist: %w", err)
	}

	return nil
}

// BatchGetArtists retrieves multiple artists by their IDs
func (r *DynamoDBRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	if len(artistIDs) == 0 {
		return make(map[string]*models.Artist), nil
	}

	// Build keys for batch get
	keys := make([]map[string]types.AttributeValue, 0, len(artistIDs))
	for _, artistID := range artistIDs {
		keys = append(keys, map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ARTIST#%s", artistID)},
		})
	}

	result := make(map[string]*models.Artist)

	// Process in batches of 100 (DynamoDB limit)
	for i := 0; i < len(keys); i += 100 {
		end := i + 100
		if end > len(keys) {
			end = len(keys)
		}

		batchResult, err := r.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				r.tableName: {
					Keys: keys[i:end],
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to batch get artists: %w", err)
		}

		// Process results
		if items, ok := batchResult.Responses[r.tableName]; ok {
			var artistItems []models.ArtistItem
			if err := attributevalue.UnmarshalListOfMaps(items, &artistItems); err != nil {
				return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
			}

			for _, item := range artistItems {
				artist := item.Artist
				result[artist.ID] = &artist
			}
		}
	}

	return result, nil
}

// SearchArtists searches artists by name prefix using GSI1
func (r *DynamoDBRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	if limit <= 0 {
		limit = 10
	}

	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#ARTIST", userID))).
		And(expression.Key("GSI1SK").BeginsWith(query))

	// Filter only active artists
	filterExpr := expression.Name("isActive").Equal(expression.Value(true))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithFilter(filterExpr)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search artists: %w", err)
	}

	var items []models.ArtistItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artists: %w", err)
	}

	artists := make([]*models.Artist, 0, len(items))
	for _, item := range items {
		artist := item.Artist
		artists = append(artists, &artist)
	}

	return artists, nil
}

// GetArtistTrackCount returns the number of tracks for an artist
func (r *DynamoDBRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	// Count tracks where artistId matches
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TRACK#"))

	filterExpr := expression.Name("artistId").Equal(expression.Value(artistID))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithFilter(filterExpr)
	expr, err := builder.Build()
	if err != nil {
		return 0, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Select:                    types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count artist tracks: %w", err)
	}

	return int(result.Count), nil
}

// GetArtistAlbumCount returns the number of albums for an artist
func (r *DynamoDBRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	// Query GSI1 for albums by artist
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#ARTIST#%s", userID, artistID)))

	// Filter for only ALBUM type
	filterExpr := expression.Name("Type").Equal(expression.Value("ALBUM"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithFilter(filterExpr)
	expr, err := builder.Build()
	if err != nil {
		return 0, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Select:                    types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count artist albums: %w", err)
	}

	return int(result.Count), nil
}

// GetArtistTotalPlays returns the total play count for an artist's tracks
func (r *DynamoDBRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	// Query all tracks for this artist
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TRACK#"))

	filterExpr := expression.Name("artistId").Equal(expression.Value(artistID))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithFilter(filterExpr)
	expr, err := builder.Build()
	if err != nil {
		return 0, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      aws.String("playCount"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to query artist tracks: %w", err)
	}

	totalPlays := 0
	for _, item := range result.Items {
		if playCount, ok := item["playCount"]; ok {
			if n, ok := playCount.(*types.AttributeValueMemberN); ok {
				var count int
				if err := attributevalue.Unmarshal(n, &count); err == nil {
					totalPlays += count
				}
			}
		}
	}

	return totalPlays, nil
}

// BatchGetItem interface addition for DynamoDBClient
type BatchGetItemAPI interface {
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
}
