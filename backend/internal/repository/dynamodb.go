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

// DynamoDBClient interface for testability
type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

// DynamoDBRepository implements Repository using DynamoDB
type DynamoDBRepository struct {
	client    DynamoDBClient
	tableName string
}

// NewDynamoDBRepository creates a new DynamoDB repository
func NewDynamoDBRepository(client DynamoDBClient, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

// ============================================================================
// Track Operations
// ============================================================================

func (r *DynamoDBRepository) CreateTrack(ctx context.Context, track models.Track) error {
	track.CreatedAt = time.Now()
	track.UpdatedAt = track.CreatedAt

	item := models.NewTrackItem(track)
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
		return nil, ErrNotFound
	}

	var item models.TrackItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track: %w", err)
	}

	return &item.Track, nil
}

func (r *DynamoDBRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	track.UpdatedAt = time.Now()

	item := models.NewTrackItem(track)
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

func (r *DynamoDBRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRACK#%s", trackID)},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*PaginatedResult[models.Track], error) {
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}

	// Global scope: scan all tracks across all users
	if filter.GlobalScope {
		return r.listAllTracks(ctx, limit, filter)
	}

	// User-scoped query (default behavior)
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TRACK#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit + 1)), // Get one extra to check hasMore
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
		return nil, fmt.Errorf("failed to query tracks: %w", err)
	}

	var items []models.TrackItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tracks: %w", err)
	}

	tracks := make([]models.Track, 0, len(items))
	for _, item := range items {
		tracks = append(tracks, item.Track)
	}

	// Determine if there are more results
	hasMore := len(tracks) > limit
	if hasMore {
		tracks = tracks[:limit]
	}

	// Build next cursor
	var nextCursor string
	if hasMore && len(tracks) > 0 {
		lastTrack := tracks[len(tracks)-1]
		cursor := models.NewPaginationCursor(
			fmt.Sprintf("USER#%s", userID),
			fmt.Sprintf("TRACK#%s", lastTrack.ID),
		)
		nextCursor = models.EncodeCursor(cursor)
	}

	return &PaginatedResult[models.Track]{
		Items:      tracks,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// listAllTracks returns tracks from all users (requires GLOBAL permission)
func (r *DynamoDBRepository) listAllTracks(ctx context.Context, limit int, filter models.TrackFilter) (*PaginatedResult[models.Track], error) {
	// Filter for TRACK# SK prefix
	filterExpr := expression.Name("SK").BeginsWith("TRACK#")
	builder := expression.NewBuilder().WithFilter(filterExpr)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(r.tableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit + 1)),
	}

	// Handle pagination cursor for scan
	if filter.LastKey != "" {
		cursor, err := models.DecodeCursor(filter.LastKey)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		input.ExclusiveStartKey = cursorToAttributeValue(cursor)
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan tracks: %w", err)
	}

	var items []models.TrackItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tracks: %w", err)
	}

	tracks := make([]models.Track, 0, len(items))
	for _, item := range items {
		tracks = append(tracks, item.Track)
	}

	// Determine if there are more results
	hasMore := len(tracks) > limit
	if hasMore {
		tracks = tracks[:limit]
	}

	// Build next cursor
	var nextCursor string
	if hasMore && len(tracks) > 0 {
		lastTrack := tracks[len(tracks)-1]
		cursor := models.NewPaginationCursor(
			fmt.Sprintf("USER#%s", lastTrack.UserID),
			fmt.Sprintf("TRACK#%s", lastTrack.ID),
		)
		nextCursor = models.EncodeCursor(cursor)
	}

	return &PaginatedResult[models.Track]{
		Items:      tracks,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *DynamoDBRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#ARTIST#%s", userID, artist)))
	// Filter to only return tracks (not albums which share the same GSI1PK)
	filter := expression.Name("Type").Equal(expression.Value("TRACK"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithFilter(filter)
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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks by artist: %w", err)
	}

	var items []models.TrackItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tracks: %w", err)
	}

	tracks := make([]models.Track, 0, len(items))
	for _, item := range items {
		tracks = append(tracks, item.Track)
	}

	return tracks, nil
}

// ============================================================================
// Album Operations
// ============================================================================

func (r *DynamoDBRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	// Generate consistent album ID from name and artist
	albumID := generateAlbumID(albumName, artist)

	// Try to get existing album
	album, err := r.GetAlbum(ctx, userID, albumID)
	if err == nil {
		return album, nil
	}
	if err != ErrNotFound {
		return nil, err
	}

	// Create new album
	now := time.Now()
	album = &models.Album{
		ID:            albumID,
		UserID:        userID,
		Title:         albumName,
		Artist:        artist,
		Year:          0, // Will be updated when tracks are added
		TrackCount:    0,
		TotalDuration: 0,
	}
	album.CreatedAt = now
	album.UpdatedAt = now

	item := models.NewAlbumItem(*album)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal album: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})
	if err != nil {
		// If album was created by another request, get it
		existingAlbum, getErr := r.GetAlbum(ctx, userID, albumID)
		if getErr == nil {
			return existingAlbum, nil
		}
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	return album, nil
}

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
		return nil, ErrNotFound
	}

	var item models.AlbumItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album: %w", err)
	}

	return &item.Album, nil
}

func (r *DynamoDBRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*PaginatedResult[models.Album], error) {
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}

	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("ALBUM#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit + 1)),
	}

	if filter.LastKey != "" {
		cursor, err := models.DecodeCursor(filter.LastKey)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		input.ExclusiveStartKey = cursorToAttributeValue(cursor)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query albums: %w", err)
	}

	var items []models.AlbumItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal albums: %w", err)
	}

	albums := make([]models.Album, 0, len(items))
	for _, item := range items {
		albums = append(albums, item.Album)
	}

	hasMore := len(albums) > limit
	if hasMore {
		albums = albums[:limit]
	}

	var nextCursor string
	if hasMore && len(albums) > 0 {
		lastAlbum := albums[len(albums)-1]
		cursor := models.NewPaginationCursor(
			fmt.Sprintf("USER#%s", userID),
			fmt.Sprintf("ALBUM#%s", lastAlbum.ID),
		)
		nextCursor = models.EncodeCursor(cursor)
	}

	return &PaginatedResult[models.Album]{
		Items:      albums,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *DynamoDBRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#ARTIST#%s", userID, artist)))

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
		return nil, fmt.Errorf("failed to query albums by artist: %w", err)
	}

	var items []models.AlbumItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal albums: %w", err)
	}

	albums := make([]models.Album, 0, len(items))
	for _, item := range items {
		albums = append(albums, item.Album)
	}

	return albums, nil
}

func (r *DynamoDBRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	update := expression.Set(
		expression.Name("trackCount"), expression.Value(trackCount),
	).Set(
		expression.Name("totalDuration"), expression.Value(totalDuration),
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
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ALBUM#%s", albumID)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update album stats: %w", err)
	}

	return nil
}

// ============================================================================
// User Operations
// ============================================================================

func (r *DynamoDBRepository) CreateUser(ctx context.Context, user models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	item := models.NewUserItem(user)
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
		return nil, ErrNotFound
	}

	var item models.UserItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &item.User, nil
}

func (r *DynamoDBRepository) UpdateUser(ctx context.Context, user models.User) error {
	user.UpdatedAt = time.Now()

	item := models.NewUserItem(user)
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

func (r *DynamoDBRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	update := expression.Set(
		expression.Name("storageUsed"), expression.Value(storageUsed),
	).Set(
		expression.Name("trackCount"), expression.Value(trackCount),
	).Set(
		expression.Name("albumCount"), expression.Value(albumCount),
	).Set(
		expression.Name("playlistCount"), expression.Value(playlistCount),
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
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}

	return nil
}

// ============================================================================
// Playlist Operations
// ============================================================================

func (r *DynamoDBRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	playlist.CreatedAt = time.Now()
	playlist.UpdatedAt = playlist.CreatedAt

	item := models.NewPlaylistItem(playlist)
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
		return nil, ErrNotFound
	}

	var item models.PlaylistItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist: %w", err)
	}

	return &item.Playlist, nil
}

func (r *DynamoDBRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	playlist.UpdatedAt = time.Now()

	item := models.NewPlaylistItem(playlist)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal playlist: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	// Delete playlist tracks first
	tracks, err := r.GetPlaylistTracks(ctx, playlistID)
	if err != nil && err != ErrNotFound {
		return fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	// Delete in batches of 25
	for i := 0; i < len(tracks); i += 25 {
		end := i + 25
		if end > len(tracks) {
			end = len(tracks)
		}

		writeRequests := make([]types.WriteRequest, 0, end-i)
		for _, track := range tracks[i:end] {
			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAYLIST#%s", playlistID)},
						"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("POSITION#%08d", track.Position)},
					},
				},
			})
		}

		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete playlist tracks: %w", err)
		}
	}

	// Delete the playlist itself
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

func (r *DynamoDBRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*PaginatedResult[models.Playlist], error) {
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}

	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("PLAYLIST#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit + 1)),
	}

	if filter.LastKey != "" {
		cursor, err := models.DecodeCursor(filter.LastKey)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		input.ExclusiveStartKey = cursorToAttributeValue(cursor)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query playlists: %w", err)
	}

	var items []models.PlaylistItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlists: %w", err)
	}

	playlists := make([]models.Playlist, 0, len(items))
	for _, item := range items {
		playlists = append(playlists, item.Playlist)
	}

	hasMore := len(playlists) > limit
	if hasMore {
		playlists = playlists[:limit]
	}

	var nextCursor string
	if hasMore && len(playlists) > 0 {
		lastPlaylist := playlists[len(playlists)-1]
		cursor := models.NewPaginationCursor(
			fmt.Sprintf("USER#%s", userID),
			fmt.Sprintf("PLAYLIST#%s", lastPlaylist.ID),
		)
		nextCursor = models.EncodeCursor(cursor)
	}

	return &PaginatedResult[models.Playlist]{
		Items:      playlists,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *DynamoDBRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	if limit <= 0 {
		limit = 10
	}

	// Query all playlists for the user
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("PLAYLIST#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query playlists: %w", err)
	}

	var items []models.PlaylistItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlists: %w", err)
	}

	// Filter by name (case-insensitive contains)
	queryLower := strings.ToLower(query)
	playlists := make([]models.Playlist, 0)
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), queryLower) {
			playlists = append(playlists, item.Playlist)
			if len(playlists) >= limit {
				break
			}
		}
	}

	return playlists, nil
}

func (r *DynamoDBRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	writeRequests := make([]types.WriteRequest, 0, len(trackIDs))
	now := time.Now()

	for i, trackID := range trackIDs {
		track := models.PlaylistTrack{
			PlaylistID: playlistID,
			TrackID:    trackID,
			Position:   position + i,
			AddedAt:    now,
		}

		item := models.NewPlaylistTrackItem(track)
		av, err := attributevalue.MarshalMap(item)
		if err != nil {
			return fmt.Errorf("failed to marshal playlist track: %w", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		})
	}

	// Write in batches of 25
	for i := 0; i < len(writeRequests); i += 25 {
		end := i + 25
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests[i:end],
			},
		})
		if err != nil {
			return fmt.Errorf("failed to add tracks to playlist: %w", err)
		}
	}

	return nil
}

func (r *DynamoDBRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	// Get all playlist tracks to find positions
	tracks, err := r.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	// Find positions for tracks to remove
	trackIDSet := make(map[string]bool)
	for _, id := range trackIDs {
		trackIDSet[id] = true
	}

	writeRequests := make([]types.WriteRequest, 0)
	for _, track := range tracks {
		if trackIDSet[track.TrackID] {
			writeRequests = append(writeRequests, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAYLIST#%s", playlistID)},
						"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("POSITION#%08d", track.Position)},
					},
				},
			})
		}
	}

	// Delete in batches of 25
	for i := 0; i < len(writeRequests); i += 25 {
		end := i + 25
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests[i:end],
			},
		})
		if err != nil {
			return fmt.Errorf("failed to remove tracks from playlist: %w", err)
		}
	}

	return nil
}

func (r *DynamoDBRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("PLAYLIST#%s", playlistID)))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query playlist tracks: %w", err)
	}

	var items []models.PlaylistTrackItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist tracks: %w", err)
	}

	tracks := make([]models.PlaylistTrack, 0, len(items))
	for _, item := range items {
		tracks = append(tracks, item.PlaylistTrack)
	}

	return tracks, nil
}

// ============================================================================
// Tag Operations
// ============================================================================

func (r *DynamoDBRepository) CreateTag(ctx context.Context, tag models.Tag) error {
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = tag.CreatedAt

	item := models.NewTagItem(tag)
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
		return nil, ErrNotFound
	}

	var item models.TagItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tag: %w", err)
	}

	return &item.Tag, nil
}

func (r *DynamoDBRepository) UpdateTag(ctx context.Context, tag models.Tag) error {
	tag.UpdatedAt = time.Now()

	item := models.NewTagItem(tag)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal tag: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
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

func (r *DynamoDBRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("TAG#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}

	var items []models.TagItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	tags := make([]models.Tag, 0, len(items))
	for _, item := range items {
		tags = append(tags, item.Tag)
	}

	return tags, nil
}

func (r *DynamoDBRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	now := time.Now()
	writeRequests := make([]types.WriteRequest, 0, len(tagNames))

	for _, tagName := range tagNames {
		trackTag := models.TrackTag{
			UserID:  userID,
			TrackID: trackID,
			TagName: tagName,
			AddedAt: now,
		}

		item := models.NewTrackTagItem(trackTag)
		av, err := attributevalue.MarshalMap(item)
		if err != nil {
			return fmt.Errorf("failed to marshal track tag: %w", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		})
	}

	// Write in batches of 25
	for i := 0; i < len(writeRequests); i += 25 {
		end := i + 25
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: writeRequests[i:end],
			},
		})
		if err != nil {
			return fmt.Errorf("failed to add tags to track: %w", err)
		}
	}

	return nil
}

func (r *DynamoDBRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s#TRACK#%s", userID, trackID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TAG#%s", tagName)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to remove tag from track: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s#TRACK#%s", userID, trackID)))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query track tags: %w", err)
	}

	var items []models.TrackTagItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track tags: %w", err)
	}

	tags := make([]string, 0, len(items))
	for _, item := range items {
		tags = append(tags, item.TrackTag.TagName)
	}

	return tags, nil
}

func (r *DynamoDBRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("USER#%s#TAG#%s", userID, tagName)))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	// First, get the track IDs from the tag associations
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

	var items []models.TrackTagItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track tags: %w", err)
	}

	// Get each track
	tracks := make([]models.Track, 0, len(items))
	for _, item := range items {
		track, err := r.GetTrack(ctx, userID, item.TrackTag.TrackID)
		if err != nil {
			if err == ErrNotFound {
				continue // Track was deleted
			}
			return nil, err
		}
		tracks = append(tracks, *track)
	}

	return tracks, nil
}

// ============================================================================
// Upload Operations
// ============================================================================

func (r *DynamoDBRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	upload.CreatedAt = time.Now()
	upload.UpdatedAt = upload.CreatedAt

	item := models.NewUploadItem(upload)
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
		return nil, ErrNotFound
	}

	var item models.UploadItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal upload: %w", err)
	}

	return &item.Upload, nil
}

func (r *DynamoDBRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	upload.UpdatedAt = time.Now()

	item := models.NewUploadItem(upload)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal upload: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update upload: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	now := time.Now()
	update := expression.Set(
		expression.Name("status"), expression.Value(string(status)),
	).Set(
		expression.Name("updatedAt"), expression.Value(now.Format(time.RFC3339)),
	).Set(
		expression.Name("GSI1PK"), expression.Value(fmt.Sprintf("UPLOAD#STATUS#%s", status)),
	).Set(
		expression.Name("GSI1SK"), expression.Value(now.Format(time.RFC3339)),
	)

	if errorMsg != "" {
		update = update.Set(expression.Name("errorMsg"), expression.Value(errorMsg))
	}
	if trackID != "" {
		update = update.Set(expression.Name("trackId"), expression.Value(trackID))
	}
	if status == models.UploadStatusCompleted {
		update = update.Set(expression.Name("completedAt"), expression.Value(now.Format(time.RFC3339)))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("UPLOAD#%s", uploadID)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update upload status: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	var fieldName string
	switch step {
	case models.StepExtractMetadata:
		fieldName = "metadataExtracted"
	case models.StepExtractCover:
		fieldName = "coverArtExtracted"
	case models.StepCreateTrack:
		fieldName = "trackCreated"
	case models.StepIndex:
		fieldName = "indexed"
	case models.StepMoveFile:
		fieldName = "fileMoved"
	default:
		return fmt.Errorf("unknown processing step: %s", step)
	}

	update := expression.Set(
		expression.Name(fieldName), expression.Value(success),
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
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("UPLOAD#%s", uploadID)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       aws.String("attribute_exists(PK)"),
	})
	if err != nil {
		return fmt.Errorf("failed to update upload step: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*PaginatedResult[models.Upload], error) {
	limit := filter.Limit
	if limit == 0 {
		limit = 20
	}

	keyCondition := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("USER#%s", userID))).
		And(expression.Key("SK").BeginsWith("UPLOAD#"))

	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit + 1)),
		ScanIndexForward:          aws.Bool(false), // Most recent first
	}

	if filter.LastKey != "" {
		cursor, err := models.DecodeCursor(filter.LastKey)
		if err != nil {
			return nil, ErrInvalidCursor
		}
		input.ExclusiveStartKey = cursorToAttributeValue(cursor)
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query uploads: %w", err)
	}

	var items []models.UploadItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal uploads: %w", err)
	}

	uploads := make([]models.Upload, 0, len(items))
	for _, item := range items {
		// Filter by status if specified
		if filter.Status != "" && item.Upload.Status != filter.Status {
			continue
		}
		uploads = append(uploads, item.Upload)
	}

	hasMore := len(uploads) > limit
	if hasMore {
		uploads = uploads[:limit]
	}

	var nextCursor string
	if hasMore && len(uploads) > 0 {
		lastUpload := uploads[len(uploads)-1]
		cursor := models.NewPaginationCursor(
			fmt.Sprintf("USER#%s", userID),
			fmt.Sprintf("UPLOAD#%s", lastUpload.ID),
		)
		nextCursor = models.EncodeCursor(cursor)
	}

	return &PaginatedResult[models.Upload]{
		Items:      uploads,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *DynamoDBRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	keyCondition := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("UPLOAD#STATUS#%s", status)))

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
		return nil, fmt.Errorf("failed to query uploads by status: %w", err)
	}

	var items []models.UploadItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal uploads: %w", err)
	}

	uploads := make([]models.Upload, 0, len(items))
	for _, item := range items {
		uploads = append(uploads, item.Upload)
	}

	return uploads, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// cursorToAttributeValue converts a PaginationCursor to DynamoDB ExclusiveStartKey
func cursorToAttributeValue(cursor models.PaginationCursor) map[string]types.AttributeValue {
	av := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: cursor.PK},
		"SK": &types.AttributeValueMemberS{Value: cursor.SK},
	}

	if cursor.GSI1PK != "" {
		av["GSI1PK"] = &types.AttributeValueMemberS{Value: cursor.GSI1PK}
	}
	if cursor.GSI1SK != "" {
		av["GSI1SK"] = &types.AttributeValueMemberS{Value: cursor.GSI1SK}
	}

	return av
}

// generateAlbumID creates a consistent album ID from name and artist
func generateAlbumID(albumName, artist string) string {
	// Use a simple hash for consistent ID generation
	combined := fmt.Sprintf("%s-%s", albumName, artist)
	// For now, just use the combined string as ID
	// In production, you might want to use a proper hash function
	return combined
}
