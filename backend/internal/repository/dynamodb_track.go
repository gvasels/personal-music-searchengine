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

// GetTrack retrieves a track by ID
func (r *DynamoDBRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRACK#%s", trackID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	if result.Item == nil {
		return nil, models.NewNotFoundError("Track", trackID)
	}

	var item models.TrackItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track: %w", err)
	}

	return &item.Track, nil
}

// CreateTrack creates a new track
func (r *DynamoDBRepository) CreateTrack(ctx context.Context, track *models.Track) error {
	item := models.NewTrackItem(*track)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal track: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	return nil
}

// UpdateTrack updates an existing track
func (r *DynamoDBRepository) UpdateTrack(ctx context.Context, track *models.Track) error {
	item := models.NewTrackItem(*track)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal track: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	return nil
}

// DeleteTrack deletes a track
func (r *DynamoDBRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRACK#%s", trackID)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	return nil
}

// ListTracks lists tracks with filtering and pagination
func (r *DynamoDBRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.Track], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TRACK#"))

	// Build filter expression
	var filterExpr expression.ConditionBuilder
	hasFilter := false

	if filter.Artist != "" {
		filterExpr = expression.Name("artist").Equal(expression.Value(filter.Artist))
		hasFilter = true
	}
	if filter.Album != "" {
		if hasFilter {
			filterExpr = filterExpr.And(expression.Name("album").Equal(expression.Value(filter.Album)))
		} else {
			filterExpr = expression.Name("album").Equal(expression.Value(filter.Album))
			hasFilter = true
		}
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

	// Handle pagination
	if filter.LastKey != "" {
		startKey, err := decodeLastKey(filter.LastKey)
		if err != nil {
			return nil, fmt.Errorf("invalid last key: %w", err)
		}
		input.ExclusiveStartKey = startKey
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks: %w", err)
	}

	tracks := make([]models.Track, 0, len(result.Items))
	for _, item := range result.Items {
		var trackItem models.TrackItem
		if err := attributevalue.UnmarshalMap(item, &trackItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal track: %w", err)
		}
		tracks = append(tracks, trackItem.Track)
	}

	response := &models.PaginatedResponse[models.Track]{
		Items: tracks,
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

// ListTracksByAlbum lists all tracks in an album
func (r *DynamoDBRepository) ListTracksByAlbum(ctx context.Context, userID, albumID string) ([]models.Track, error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TRACK#"))
	filterExpr := expression.Name("albumId").Equal(expression.Value(albumID))

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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks by album: %w", err)
	}

	tracks := make([]models.Track, 0, len(result.Items))
	for _, item := range result.Items {
		var trackItem models.TrackItem
		if err := attributevalue.UnmarshalMap(item, &trackItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal track: %w", err)
		}
		tracks = append(tracks, trackItem.Track)
	}

	return tracks, nil
}

// ListTracksByTag lists all tracks with a specific tag
func (r *DynamoDBRepository) ListTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	// Query GSI to find track-tag associations
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#TAG#%s", userID, tagName)))

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
		return nil, fmt.Errorf("failed to query tracks by tag: %w", err)
	}

	// Get track IDs from results
	trackIDs := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		var tt models.TrackTagItem
		if err := attributevalue.UnmarshalMap(item, &tt); err != nil {
			continue
		}
		trackIDs = append(trackIDs, tt.TrackID)
	}

	// Batch get tracks
	if len(trackIDs) == 0 {
		return []models.Track{}, nil
	}

	tracks := make([]models.Track, 0, len(trackIDs))
	for _, trackID := range trackIDs {
		track, err := r.GetTrack(ctx, userID, trackID)
		if err == nil {
			tracks = append(tracks, *track)
		}
	}

	return tracks, nil
}

// Helper functions for pagination encoding/decoding
func encodeLastKey(key map[string]types.AttributeValue) (string, error) {
	// Simple implementation - in production, use base64 encoding of JSON
	if pk, ok := key["PK"].(*types.AttributeValueMemberS); ok {
		if sk, ok := key["SK"].(*types.AttributeValueMemberS); ok {
			return pk.Value + "|" + sk.Value, nil
		}
	}
	return "", fmt.Errorf("invalid key format")
}

func decodeLastKey(encoded string) (map[string]types.AttributeValue, error) {
	// Simple implementation - in production, decode from base64 JSON
	// For now, return nil to indicate start from beginning
	return nil, nil
}
