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

// ListTracksForMatching retrieves all tracks for a user with BPM/Key data for matching
// This returns a simplified list optimized for matching operations
func (r *DynamoDBRepository) ListTracksForMatching(ctx context.Context, userID string) ([]models.Track, error) {
	pk := fmt.Sprintf("USER#%s", userID)

	var tracks []models.Track
	var lastKey map[string]types.AttributeValue

	for {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":       &types.AttributeValueMemberS{Value: pk},
				":skPrefix": &types.AttributeValueMemberS{Value: "TRACK#"},
			},
			// Only fetch fields needed for matching
			ProjectionExpression: aws.String("id, title, artist, album, bpm, musicalKey, keyCamelot, keyMode, #dur"),
			ExpressionAttributeNames: map[string]string{
				"#dur": "duration",
			},
		}

		if lastKey != nil {
			input.ExclusiveStartKey = lastKey
		}

		result, err := r.client.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list tracks for matching: %w", err)
		}

		var items []models.TrackItem
		if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tracks: %w", err)
		}

		for _, item := range items {
			tracks = append(tracks, item.Track)
		}

		if result.LastEvaluatedKey == nil {
			break
		}
		lastKey = result.LastEvaluatedKey
	}

	return tracks, nil
}

// GetTracksWithBPMRange retrieves tracks within a BPM range
func (r *DynamoDBRepository) GetTracksWithBPMRange(ctx context.Context, userID string, minBPM, maxBPM int) ([]models.Track, error) {
	pk := fmt.Sprintf("USER#%s", userID)

	var tracks []models.Track
	var lastKey map[string]types.AttributeValue

	for {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
			FilterExpression:       aws.String("bpm BETWEEN :minBPM AND :maxBPM"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":       &types.AttributeValueMemberS{Value: pk},
				":skPrefix": &types.AttributeValueMemberS{Value: "TRACK#"},
				":minBPM":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", minBPM)},
				":maxBPM":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", maxBPM)},
			},
		}

		if lastKey != nil {
			input.ExclusiveStartKey = lastKey
		}

		result, err := r.client.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to query tracks by BPM: %w", err)
		}

		var items []models.TrackItem
		if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tracks: %w", err)
		}

		for _, item := range items {
			tracks = append(tracks, item.Track)
		}

		if result.LastEvaluatedKey == nil {
			break
		}
		lastKey = result.LastEvaluatedKey
	}

	return tracks, nil
}
