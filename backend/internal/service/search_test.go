package service

import (
	"context"
	"errors"
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
			Text: result.Title,
			Type: "track",
			ID:   result.ID,
		})

		if result.Artist != "" && !seenArtists[result.Artist] {
			seenArtists[result.Artist] = true
			suggestions = append(suggestions, models.SearchSuggestion{
				Text: result.Artist,
				Type: "artist",
			})
		}

		if result.Album != "" && !seenAlbums[result.Album] {
			seenAlbums[result.Album] = true
			suggestions = append(suggestions, models.SearchSuggestion{
				Text: result.Album,
				Type: "album",
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

	if !resp.Deleted {
		// Log but don't fail
	}

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
