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

// GetPlaylist retrieves a playlist by ID
func (r *DynamoDBRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAYLIST#%s", playlistID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	if result.Item == nil {
		return nil, models.NewNotFoundError("Playlist", playlistID)
	}

	var item models.PlaylistItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist: %w", err)
	}

	return &item.Playlist, nil
}

// CreatePlaylist creates a new playlist
func (r *DynamoDBRepository) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	item := models.NewPlaylistItem(*playlist)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal playlist: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}

	return nil
}

// UpdatePlaylist updates an existing playlist
func (r *DynamoDBRepository) UpdatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	item := models.NewPlaylistItem(*playlist)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal playlist: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}

	return nil
}

// DeletePlaylist deletes a playlist and all its tracks
func (r *DynamoDBRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	// First delete all playlist tracks
	tracks, err := r.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		if err := r.RemovePlaylistTrack(ctx, playlistID, track.TrackID); err != nil {
			return err
		}
	}

	// Then delete the playlist
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAYLIST#%s", playlistID)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	return nil
}

// ListPlaylists lists playlists with filtering and pagination
func (r *DynamoDBRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResponse[models.Playlist], error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("PLAYLIST#"))

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
		return nil, fmt.Errorf("failed to query playlists: %w", err)
	}

	playlists := make([]models.Playlist, 0, len(result.Items))
	for _, item := range result.Items {
		var playlistItem models.PlaylistItem
		if err := attributevalue.UnmarshalMap(item, &playlistItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal playlist: %w", err)
		}
		playlists = append(playlists, playlistItem.Playlist)
	}

	response := &models.PaginatedResponse[models.Playlist]{
		Items: playlists,
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

// GetPlaylistTracks retrieves all tracks in a playlist
func (r *DynamoDBRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("PLAYLIST#%s", playlistID))).
		And(expression.Key("SK").BeginsWith("POSITION#"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query playlist tracks: %w", err)
	}

	tracks := make([]models.PlaylistTrack, 0, len(result.Items))
	for _, item := range result.Items {
		var ptItem models.PlaylistTrackItem
		if err := attributevalue.UnmarshalMap(item, &ptItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal playlist track: %w", err)
		}
		tracks = append(tracks, ptItem.PlaylistTrack)
	}

	return tracks, nil
}

// AddPlaylistTrack adds a track to a playlist
func (r *DynamoDBRepository) AddPlaylistTrack(ctx context.Context, pt *models.PlaylistTrack) error {
	item := models.NewPlaylistTrackItem(*pt)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal playlist track: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to add playlist track: %w", err)
	}

	return nil
}

// RemovePlaylistTrack removes a track from a playlist
func (r *DynamoDBRepository) RemovePlaylistTrack(ctx context.Context, playlistID, trackID string) error {
	// Find the track position first
	tracks, err := r.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	var position int = -1
	for _, t := range tracks {
		if t.TrackID == trackID {
			position = t.Position
			break
		}
	}

	if position == -1 {
		return nil // Track not in playlist
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAYLIST#%s", playlistID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("POSITION#%08d", position)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to remove playlist track: %w", err)
	}

	return nil
}

// ReorderPlaylistTracks reorders tracks in a playlist
func (r *DynamoDBRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, trackIDs []string) error {
	// Get existing tracks
	existingTracks, err := r.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	// Create a map of track ID to added time
	addedTimes := make(map[string]models.PlaylistTrack)
	for _, t := range existingTracks {
		addedTimes[t.TrackID] = t
	}

	// Delete all existing tracks
	for _, t := range existingTracks {
		_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAYLIST#%s", playlistID)},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("POSITION#%08d", t.Position)},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete track during reorder: %w", err)
		}
	}

	// Re-add tracks in new order
	for i, trackID := range trackIDs {
		if existing, ok := addedTimes[trackID]; ok {
			pt := &models.PlaylistTrack{
				PlaylistID: playlistID,
				TrackID:    trackID,
				Position:   i,
				AddedAt:    existing.AddedAt,
			}
			if err := r.AddPlaylistTrack(ctx, pt); err != nil {
				return err
			}
		}
	}

	return nil
}
