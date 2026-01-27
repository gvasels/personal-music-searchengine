package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchClient mocks the search.Client
type MockSearchClient struct {
	mock.Mock
}

func (m *MockSearchClient) Search(ctx context.Context, userID string, query search.SearchQuery) (*search.SearchResponse, error) {
	args := m.Called(ctx, userID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.SearchResponse), args.Error(1)
}

func (m *MockSearchClient) Index(ctx context.Context, doc search.Document) (*search.IndexResponse, error) {
	args := m.Called(ctx, doc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.IndexResponse), args.Error(1)
}

func (m *MockSearchClient) Delete(ctx context.Context, docID string) (*search.DeleteResponse, error) {
	args := m.Called(ctx, docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.DeleteResponse), args.Error(1)
}

func (m *MockSearchClient) BulkIndex(ctx context.Context, docs []search.Document) (*search.BulkIndexResponse, error) {
	args := m.Called(ctx, docs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.BulkIndexResponse), args.Error(1)
}

// MockRepository mocks the repository.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTrack(ctx context.Context, track models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	args := m.Called(ctx, userID, trackID)
	return args.Error(0)
}

func (m *MockRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

func (m *MockRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	args := m.Called(ctx, userID, artist)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

// Stub implementations for other Repository methods
func (m *MockRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}
func (m *MockRepository) CreateUser(ctx context.Context, user models.User) error       { return nil }
func (m *MockRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockRepository) UpdateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}
func (m *MockRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	return nil, nil
}
func (m *MockRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return nil
}
func (m *MockRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	return nil
}
func (m *MockRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	return nil
}
func (m *MockRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	return nil
}
func (m *MockRepository) CreateTag(ctx context.Context, tag models.Tag) error { return nil }
func (m *MockRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	return nil, nil
}
func (m *MockRepository) UpdateTag(ctx context.Context, tag models.Tag) error { return nil }
func (m *MockRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	return nil
}
func (m *MockRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	return nil, nil
}
func (m *MockRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	return nil
}
func (m *MockRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	return nil
}
func (m *MockRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// Artist-related methods for Repository interface
func (m *MockRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}

// User role methods
func (m *MockRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}

// Playlist visibility methods
func (m *MockRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	return nil
}
func (m *MockRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}

// ArtistProfile methods
func (m *MockRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Follow methods
func (m *MockRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Admin-related methods for track visibility
func (m *MockRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}
func (m *MockRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

// User settings methods
func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return nil, nil
}
func (m *MockRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]repository.UserSearchResult, string, error) {
	return nil, "", nil
}
func (m *MockRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockRepository) UpdateUserSettings(ctx context.Context, userID string, update *repository.UserSettingsUpdate) (*models.UserSettings, error) {
	return nil, nil
}

// MockS3Repository mocks the repository.S3Repository
type MockS3Repository struct {
	mock.Mock
}

func (m *MockS3Repository) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockS3Repository) GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockS3Repository) InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	return "", nil
}
func (m *MockS3Repository) GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error) {
	return nil, nil
}
func (m *MockS3Repository) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error {
	return nil
}
func (m *MockS3Repository) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return nil
}
func (m *MockS3Repository) DeleteObject(ctx context.Context, key string) error { return nil }
func (m *MockS3Repository) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	return nil
}
func (m *MockS3Repository) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}
func (m *MockS3Repository) GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, expiry time.Duration, filename string) (string, error) {
	return "", nil
}
func (m *MockS3Repository) ObjectExists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// SearchClient interface for mocking
type SearchClient interface {
	Search(ctx context.Context, userID string, query search.SearchQuery) (*search.SearchResponse, error)
	Index(ctx context.Context, doc search.Document) (*search.IndexResponse, error)
	Delete(ctx context.Context, docID string) (*search.DeleteResponse, error)
	BulkIndex(ctx context.Context, docs []search.Document) (*search.BulkIndexResponse, error)
}

// testSearchService allows injecting mock client
type testSearchService struct {
	client SearchClient
	repo   repository.Repository
	s3Repo repository.S3Repository
}

func newTestSearchService(client SearchClient, repo repository.Repository, s3Repo repository.S3Repository) *testSearchService {
	return &testSearchService{
		client: client,
		repo:   repo,
		s3Repo: s3Repo,
	}
}

func (s *testSearchService) Search(ctx context.Context, userID string, req models.SearchRequest) (*models.SearchResponse, error) {
	if req.Query == "" {
		return nil, models.NewValidationError("search query cannot be empty")
	}
	if len(req.Query) > MaxQueryLength {
		return nil, models.NewValidationError(fmt.Sprintf("search query too long (maximum %d characters)", MaxQueryLength))
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	searchQuery := search.SearchQuery{
		Query:  req.Query,
		Limit:  limit,
		Cursor: req.Cursor,
	}

	// Convert filters
	if len(req.Filters.Artists) > 0 {
		searchQuery.Filters.Artist = req.Filters.Artists[0]
	}
	if len(req.Filters.Albums) > 0 {
		searchQuery.Filters.Album = req.Filters.Albums[0]
	}
	if len(req.Filters.Genres) > 0 {
		searchQuery.Filters.Genre = req.Filters.Genres[0]
	}

	if req.Sort.Field != "" {
		searchQuery.Sort = &search.SortOption{
			Field: req.Sort.Field,
			Order: req.Sort.Order,
		}
	}

	resp, err := s.client.Search(ctx, userID, searchQuery)
	if err != nil {
		return nil, err
	}

	tracks := make([]models.TrackResponse, 0, len(resp.Results))
	for _, result := range resp.Results {
		tracks = append(tracks, models.TrackResponse{
			ID:          result.ID,
			Title:       result.Title,
			Artist:      result.Artist,
			Album:       result.Album,
			Genre:       result.Genre,
			Year:        result.Year,
			Duration:    result.Duration,
			DurationStr: formatDuration(result.Duration),
		})
	}

	return &models.SearchResponse{
		Query:        req.Query,
		TotalResults: resp.Total,
		Tracks:       tracks,
		Limit:        limit,
		NextCursor:   resp.NextCursor,
		HasMore:      resp.NextCursor != "",
	}, nil
}

func (s *testSearchService) Autocomplete(ctx context.Context, userID, query string) (*models.AutocompleteResponse, error) {
	if query == "" {
		return &models.AutocompleteResponse{
			Query:       query,
			Suggestions: []models.SearchSuggestion{},
		}, nil
	}

	searchQuery := search.SearchQuery{
		Query: query,
		Limit: 10,
	}

	resp, err := s.client.Search(ctx, userID, searchQuery)
	if err != nil {
		return nil, err
	}

	suggestions := make([]models.SearchSuggestion, 0)
	seenArtists := make(map[string]bool)
	seenAlbums := make(map[string]bool)

	for _, result := range resp.Results {
		suggestions = append(suggestions, models.SearchSuggestion{
			Value:   result.Title,
			Type:    "track",
			TrackID: result.ID,
		})

		if result.Artist != "" && !seenArtists[result.Artist] {
			seenArtists[result.Artist] = true
			suggestions = append(suggestions, models.SearchSuggestion{
				Value: result.Artist,
				Type:  "artist",
			})
		}

		if result.Album != "" && !seenAlbums[result.Album] {
			seenAlbums[result.Album] = true
			suggestions = append(suggestions, models.SearchSuggestion{
				Value: result.Album,
				Type:  "album",
			})
		}

		if len(suggestions) >= 10 {
			break
		}
	}

	return &models.AutocompleteResponse{
		Query:       query,
		Suggestions: suggestions,
	}, nil
}

func (s *testSearchService) IndexTrack(ctx context.Context, track models.Track) error {
	doc := search.Document{
		ID:        track.ID,
		UserID:    track.UserID,
		Title:     track.Title,
		Artist:    track.Artist,
		Album:     track.Album,
		Genre:     track.Genre,
		Year:      track.Year,
		Duration:  track.Duration,
		Filename:  track.S3Key,
		IndexedAt: time.Now(),
	}

	resp, err := s.client.Index(ctx, doc)
	if err != nil {
		return err
	}

	if !resp.Indexed {
		return errors.New("track was not indexed")
	}

	return nil
}

func (s *testSearchService) RemoveTrack(ctx context.Context, trackID string) error {
	resp, err := s.client.Delete(ctx, trackID)
	if err != nil {
		return err
	}

	// Note: resp.Deleted is not checked - delete is best-effort in tests
	_ = resp.Deleted

	return nil
}

func TestSearch_SimpleQuery(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	expectedResults := []search.SearchResult{
		{ID: "track-1", Title: "Hey Jude", Artist: "The Beatles", Album: "Past Masters", Duration: 180},
		{ID: "track-2", Title: "Let It Be", Artist: "The Beatles", Album: "Let It Be", Duration: 240},
	}

	mockClient.On("Search", ctx, "user-123", mock.MatchedBy(func(q search.SearchQuery) bool {
		return q.Query == "beatles" && q.Limit == 20
	})).Return(&search.SearchResponse{
		Results:    expectedResults,
		Total:      2,
		NextCursor: "",
	}, nil)

	req := models.SearchRequest{Query: "beatles"}
	resp, err := svc.Search(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "beatles", resp.Query)
	assert.Equal(t, 2, resp.TotalResults)
	assert.Len(t, resp.Tracks, 2)
	assert.Equal(t, "Hey Jude", resp.Tracks[0].Title)
	assert.Equal(t, "The Beatles", resp.Tracks[0].Artist)
	assert.False(t, resp.HasMore)

	mockClient.AssertExpectations(t)
}

func TestSearch_WithFilters(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	mockClient.On("Search", ctx, "user-123", mock.MatchedBy(func(q search.SearchQuery) bool {
		return q.Query == "love" && q.Filters.Artist == "The Beatles" && q.Filters.Genre == "Rock"
	})).Return(&search.SearchResponse{
		Results: []search.SearchResult{
			{ID: "track-1", Title: "All You Need Is Love", Artist: "The Beatles", Genre: "Rock"},
		},
		Total: 1,
	}, nil)

	req := models.SearchRequest{
		Query: "love",
		Filters: models.SearchFilters{
			Artists: []string{"The Beatles"},
			Genres:  []string{"Rock"},
		},
	}
	resp, err := svc.Search(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Tracks, 1)
	assert.Equal(t, "All You Need Is Love", resp.Tracks[0].Title)

	mockClient.AssertExpectations(t)
}

func TestSearch_WithPagination(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	mockClient.On("Search", ctx, "user-123", mock.MatchedBy(func(q search.SearchQuery) bool {
		return q.Query == "rock" && q.Limit == 10 && q.Cursor == "cursor-abc"
	})).Return(&search.SearchResponse{
		Results:    []search.SearchResult{{ID: "track-3", Title: "Track 3"}},
		Total:      30,
		NextCursor: "cursor-xyz",
	}, nil)

	req := models.SearchRequest{
		Query:  "rock",
		Limit:  10,
		Cursor: "cursor-abc",
	}
	resp, err := svc.Search(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 30, resp.TotalResults)
	assert.True(t, resp.HasMore)
	assert.Equal(t, "cursor-xyz", resp.NextCursor)

	mockClient.AssertExpectations(t)
}

func TestSearch_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	req := models.SearchRequest{Query: ""}
	resp, err := svc.Search(ctx, "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestSearch_QueryTooLong(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	// Create a query that exceeds MaxQueryLength (500 characters)
	longQuery := strings.Repeat("a", MaxQueryLength+1)
	req := models.SearchRequest{Query: longQuery}
	resp, err := svc.Search(ctx, "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	// NewValidationError stores message in Details field; Error() returns Code: Message
	assert.Contains(t, err.Error(), "VALIDATION_ERROR")
	// Verify the details contain the actual message
	apiErr, ok := err.(*models.APIError)
	assert.True(t, ok, "error should be APIError")
	if ok {
		assert.Contains(t, apiErr.Details.(string), "too long")
	}
}

func TestSearch_QueryAtMaxLength(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	// Create a query exactly at MaxQueryLength (500 characters)
	maxQuery := strings.Repeat("a", MaxQueryLength)

	mockClient.On("Search", ctx, "user-123", mock.MatchedBy(func(q search.SearchQuery) bool {
		return len(q.Query) == MaxQueryLength
	})).Return(&search.SearchResponse{Results: []search.SearchResult{}, Total: 0}, nil)

	req := models.SearchRequest{Query: maxQuery}
	_, err := svc.Search(ctx, "user-123", req)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSearch_LimitClamping(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	// Test limit clamping to 100
	mockClient.On("Search", ctx, "user-123", mock.MatchedBy(func(q search.SearchQuery) bool {
		return q.Limit == 100
	})).Return(&search.SearchResponse{Results: []search.SearchResult{}, Total: 0}, nil)

	req := models.SearchRequest{Query: "test", Limit: 500}
	_, err := svc.Search(ctx, "user-123", req)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAutocomplete_Success(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	mockClient.On("Search", ctx, "user-123", mock.MatchedBy(func(q search.SearchQuery) bool {
		return q.Query == "beat" && q.Limit == 10
	})).Return(&search.SearchResponse{
		Results: []search.SearchResult{
			{ID: "t1", Title: "Beat It", Artist: "Michael Jackson", Album: "Thriller"},
			{ID: "t2", Title: "Heart Beat", Artist: "The Beatles", Album: "Love"},
		},
		Total: 2,
	}, nil)

	resp, err := svc.Autocomplete(ctx, "user-123", "beat")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "beat", resp.Query)

	// Should have track suggestions + unique artist/album suggestions
	trackCount := 0
	artistCount := 0
	albumCount := 0
	for _, s := range resp.Suggestions {
		switch s.Type {
		case "track":
			trackCount++
		case "artist":
			artistCount++
		case "album":
			albumCount++
		}
	}
	assert.Equal(t, 2, trackCount)
	assert.Equal(t, 2, artistCount) // Michael Jackson, The Beatles
	assert.Equal(t, 2, albumCount)  // Thriller, Love

	mockClient.AssertExpectations(t)
}

func TestAutocomplete_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	resp, err := svc.Autocomplete(ctx, "user-123", "")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Suggestions)
}

func TestIndexTrack_Success(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	track := models.Track{
		ID:       "track-123",
		UserID:   "user-123",
		Title:    "Test Track",
		Artist:   "Test Artist",
		Album:    "Test Album",
		Genre:    "Rock",
		Year:     2024,
		Duration: 180,
		S3Key:    "audio/track-123.mp3",
	}

	mockClient.On("Index", ctx, mock.MatchedBy(func(doc search.Document) bool {
		return doc.ID == "track-123" && doc.UserID == "user-123" && doc.Title == "Test Track"
	})).Return(&search.IndexResponse{ID: "track-123", Indexed: true}, nil)

	err := svc.IndexTrack(ctx, track)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestIndexTrack_Failure(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	track := models.Track{ID: "track-123", UserID: "user-123", Title: "Test"}

	mockClient.On("Index", ctx, mock.Anything).Return(&search.IndexResponse{ID: "track-123", Indexed: false}, nil)

	err := svc.IndexTrack(ctx, track)

	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestRemoveTrack_Success(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	mockClient.On("Delete", ctx, "track-123").Return(&search.DeleteResponse{ID: "track-123", Deleted: true}, nil)

	err := svc.RemoveTrack(ctx, "track-123")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRemoveTrack_NotFound(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockSearchClient)
	mockRepo := new(MockRepository)
	mockS3 := new(MockS3Repository)

	svc := newTestSearchService(mockClient, mockRepo, mockS3)

	// Even if not found, should not error
	mockClient.On("Delete", ctx, "track-999").Return(&search.DeleteResponse{ID: "track-999", Deleted: false}, nil)

	err := svc.RemoveTrack(ctx, "track-999")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "0:00"},
		{-1, "0:00"},
		{30, "0:30"},
		{60, "1:00"},
		{90, "1:30"},
		{180, "3:00"},
		{3661, "61:01"}, // 1 hour, 1 minute, 1 second
	}

	for _, test := range tests {
		result := formatDuration(test.seconds)
		assert.Equal(t, test.expected, result, "formatDuration(%d)", test.seconds)
	}
}

// =============================================================================
// filterByTags Tests (Epic 4)
// =============================================================================

// MockFilterTagsRepository provides mockable GetTag and GetTracksByTag for filterByTags tests
type MockFilterTagsRepository struct {
	mock.Mock
}

func (m *MockFilterTagsRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockFilterTagsRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

// Stub implementations for Repository interface (required but not used in filterByTags tests)
func (m *MockFilterTagsRepository) CreateTrack(ctx context.Context, track models.Track) error {
	return nil
}
func (m *MockFilterTagsRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	return nil
}
func (m *MockFilterTagsRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return nil
}
func (m *MockFilterTagsRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}
func (m *MockFilterTagsRepository) CreateUser(ctx context.Context, user models.User) error {
	return nil
}
func (m *MockFilterTagsRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateUser(ctx context.Context, user models.User) error {
	return nil
}
func (m *MockFilterTagsRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}
func (m *MockFilterTagsRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockFilterTagsRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockFilterTagsRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return nil
}
func (m *MockFilterTagsRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	return nil
}
func (m *MockFilterTagsRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	return nil
}
func (m *MockFilterTagsRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	return nil
}
func (m *MockFilterTagsRepository) CreateTag(ctx context.Context, tag models.Tag) error {
	return nil
}
func (m *MockFilterTagsRepository) UpdateTag(ctx context.Context, tag models.Tag) error {
	return nil
}
func (m *MockFilterTagsRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	return nil
}
func (m *MockFilterTagsRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	return nil
}
func (m *MockFilterTagsRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	return nil
}
func (m *MockFilterTagsRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockFilterTagsRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockFilterTagsRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockFilterTagsRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockFilterTagsRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// Artist-related methods
func (m *MockFilterTagsRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockFilterTagsRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockFilterTagsRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockFilterTagsRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockFilterTagsRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockFilterTagsRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockFilterTagsRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}

// User role methods
func (m *MockFilterTagsRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockFilterTagsRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}

// Playlist visibility methods
func (m *MockFilterTagsRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	return nil
}
func (m *MockFilterTagsRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}

// ArtistProfile methods
func (m *MockFilterTagsRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockFilterTagsRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockFilterTagsRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockFilterTagsRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Follow methods
func (m *MockFilterTagsRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockFilterTagsRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockFilterTagsRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Admin-related methods for track visibility
func (m *MockFilterTagsRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}
func (m *MockFilterTagsRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockFilterTagsRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockFilterTagsRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

// User settings methods
func (m *MockFilterTagsRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]repository.UserSearchResult, string, error) {
	return nil, "", nil
}
func (m *MockFilterTagsRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockFilterTagsRepository) UpdateUserSettings(ctx context.Context, userID string, update *repository.UserSettingsUpdate) (*models.UserSettings, error) {
	return nil, nil
}

// TestFilterByTags_EmptyTags verifies that empty tags array returns all results unchanged
func TestFilterByTags_EmptyTags(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
		{ID: "track-2", Title: "Track Two"},
	}

	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{})

	assert.NoError(t, err)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "track-1", filtered[0].ID)
	assert.Equal(t, "track-2", filtered[1].ID)
}

// TestFilterByTags_TagNotFound verifies NotFoundError when tag doesn't exist
func TestFilterByTags_TagNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
	}

	// Tag doesn't exist
	mockRepo.On("GetTag", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"nonexistent"})

	assert.Error(t, err)
	assert.Nil(t, filtered)

	// Verify it's a NotFoundError
	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
		assert.Contains(t, apiErr.Message, "nonexistent")
	}
	mockRepo.AssertExpectations(t)
}

// TestFilterByTags_SingleTag_Success verifies single tag filters correctly
func TestFilterByTags_SingleTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
		{ID: "track-2", Title: "Track Two"},
		{ID: "track-3", Title: "Track Three"},
	}

	// Tag exists
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{
		UserID: "user-123",
		Name:   "favorites",
	}, nil)

	// Only track-1 and track-3 have the tag
	mockRepo.On("GetTracksByTag", ctx, "user-123", "favorites").Return([]models.Track{
		{ID: "track-1"},
		{ID: "track-3"},
	}, nil)

	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"favorites"})

	assert.NoError(t, err)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "track-1", filtered[0].ID)
	assert.Equal(t, "track-3", filtered[1].ID)
	mockRepo.AssertExpectations(t)
}

// TestFilterByTags_MultipleTags_ANDLogic verifies multiple tags use AND logic
func TestFilterByTags_MultipleTags_ANDLogic(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
		{ID: "track-2", Title: "Track Two"},
		{ID: "track-3", Title: "Track Three"},
	}

	// Both tags exist
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{UserID: "user-123", Name: "favorites"}, nil)
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{UserID: "user-123", Name: "rock"}, nil)

	// track-1 and track-2 have "favorites"
	mockRepo.On("GetTracksByTag", ctx, "user-123", "favorites").Return([]models.Track{
		{ID: "track-1"},
		{ID: "track-2"},
	}, nil)

	// track-1 and track-3 have "rock"
	mockRepo.On("GetTracksByTag", ctx, "user-123", "rock").Return([]models.Track{
		{ID: "track-1"},
		{ID: "track-3"},
	}, nil)

	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"favorites", "rock"})

	assert.NoError(t, err)
	// Only track-1 has BOTH tags
	assert.Len(t, filtered, 1)
	assert.Equal(t, "track-1", filtered[0].ID)
	mockRepo.AssertExpectations(t)
}

// TestFilterByTags_SecondTagNotFound verifies error on second tag returns NotFoundError
func TestFilterByTags_SecondTagNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
	}

	// First tag exists
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{UserID: "user-123", Name: "favorites"}, nil)
	// Second tag doesn't exist
	mockRepo.On("GetTag", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"favorites", "nonexistent"})

	assert.Error(t, err)
	assert.Nil(t, filtered)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
		assert.Contains(t, apiErr.Message, "nonexistent")
	}
	mockRepo.AssertExpectations(t)
}

// TestFilterByTags_NoMatchingTracks verifies empty array when no tracks match
func TestFilterByTags_NoMatchingTracks(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
		{ID: "track-2", Title: "Track Two"},
	}

	// Tag exists but has no tracks
	mockRepo.On("GetTag", ctx, "user-123", "empty-tag").Return(&models.Tag{UserID: "user-123", Name: "empty-tag"}, nil)
	mockRepo.On("GetTracksByTag", ctx, "user-123", "empty-tag").Return([]models.Track{}, nil)

	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"empty-tag"})

	assert.NoError(t, err)
	assert.Len(t, filtered, 0)
	mockRepo.AssertExpectations(t)
}

// TestFilterByTags_DeduplicatesTags verifies duplicate tags are deduplicated silently
func TestFilterByTags_DeduplicatesTags(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
	}

	// Tag exists - should only be called ONCE despite duplicates in input
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{UserID: "user-123", Name: "rock"}, nil).Once()
	mockRepo.On("GetTracksByTag", ctx, "user-123", "rock").Return([]models.Track{
		{ID: "track-1"},
	}, nil).Once()

	// Input has duplicate tags
	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"rock", "rock", "rock"})

	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	mockRepo.AssertExpectations(t)
}

// TestFilterByTags_NormalizesCase verifies tags are normalized to lowercase
func TestFilterByTags_NormalizesCase(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockFilterTagsRepository)

	svc := &searchServiceImpl{
		client: nil,
		repo:   mockRepo,
		s3Repo: nil,
	}

	results := []search.SearchResult{
		{ID: "track-1", Title: "Track One"},
	}

	// Tag stored as lowercase - lookup should also be lowercase
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{UserID: "user-123", Name: "rock"}, nil)
	mockRepo.On("GetTracksByTag", ctx, "user-123", "rock").Return([]models.Track{
		{ID: "track-1"},
	}, nil)

	// Input is uppercase - should be normalized to lowercase
	filtered, err := svc.filterByTags(ctx, "user-123", results, []string{"ROCK"})

	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	mockRepo.AssertExpectations(t)
}
