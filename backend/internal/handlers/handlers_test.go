package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// ====================
// Mock Services
// ====================

// MockUserService implements UserServiceInterface for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID string, req models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// MockTrackService implements TrackServiceInterface for testing
type MockTrackService struct {
	mock.Mock
}

func (m *MockTrackService) GetTrack(ctx context.Context, userID, trackID string) (*models.TrackResponse, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TrackResponse), args.Error(1)
}

func (m *MockTrackService) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.TrackResponse], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.TrackResponse]), args.Error(1)
}

func (m *MockTrackService) UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error) {
	args := m.Called(ctx, userID, trackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TrackResponse), args.Error(1)
}

func (m *MockTrackService) DeleteTrack(ctx context.Context, userID, trackID string) error {
	args := m.Called(ctx, userID, trackID)
	return args.Error(0)
}

func (m *MockTrackService) AddTagsToTrack(ctx context.Context, userID, trackID string, tags []string) (*models.TrackResponse, error) {
	args := m.Called(ctx, userID, trackID, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TrackResponse), args.Error(1)
}

func (m *MockTrackService) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) (*models.TrackResponse, error) {
	args := m.Called(ctx, userID, trackID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TrackResponse), args.Error(1)
}

// MockAlbumService implements AlbumServiceInterface for testing
type MockAlbumService struct {
	mock.Mock
}

func (m *MockAlbumService) GetAlbumWithTracks(ctx context.Context, userID, albumID string) (*models.AlbumWithTracks, error) {
	args := m.Called(ctx, userID, albumID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AlbumWithTracks), args.Error(1)
}

func (m *MockAlbumService) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResponse[models.AlbumResponse], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.AlbumResponse]), args.Error(1)
}

func (m *MockAlbumService) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*models.PaginatedResponse[models.ArtistSummary], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.ArtistSummary]), args.Error(1)
}

func (m *MockAlbumService) GetArtist(ctx context.Context, userID, artistName string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, artistName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockPlaylistService implements PlaylistServiceInterface for testing
type MockPlaylistService struct {
	mock.Mock
}

func (m *MockPlaylistService) GetPlaylistWithTracks(ctx context.Context, userID, playlistID string) (*models.PlaylistWithTracks, error) {
	args := m.Called(ctx, userID, playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaylistWithTracks), args.Error(1)
}

func (m *MockPlaylistService) CreatePlaylist(ctx context.Context, userID string, req models.CreatePlaylistRequest) (*models.PlaylistResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaylistResponse), args.Error(1)
}

func (m *MockPlaylistService) UpdatePlaylist(ctx context.Context, userID, playlistID string, req models.UpdatePlaylistRequest) (*models.PlaylistResponse, error) {
	args := m.Called(ctx, userID, playlistID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaylistResponse), args.Error(1)
}

func (m *MockPlaylistService) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	args := m.Called(ctx, userID, playlistID)
	return args.Error(0)
}

func (m *MockPlaylistService) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResponse[models.PlaylistResponse], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.PlaylistResponse]), args.Error(1)
}

func (m *MockPlaylistService) AddTracks(ctx context.Context, userID, playlistID string, req models.AddTracksToPlaylistRequest) (*models.PlaylistResponse, error) {
	args := m.Called(ctx, userID, playlistID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaylistResponse), args.Error(1)
}

func (m *MockPlaylistService) RemoveTracks(ctx context.Context, userID, playlistID string, trackIDs []string) (*models.PlaylistResponse, error) {
	args := m.Called(ctx, userID, playlistID, trackIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaylistResponse), args.Error(1)
}

func (m *MockPlaylistService) ReorderTracks(ctx context.Context, userID, playlistID string, req models.ReorderPlaylistTracksRequest) (*models.PlaylistResponse, error) {
	args := m.Called(ctx, userID, playlistID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaylistResponse), args.Error(1)
}

// MockTagService implements TagServiceInterface for testing
type MockTagService struct {
	mock.Mock
}

func (m *MockTagService) CreateTag(ctx context.Context, userID string, req models.CreateTagRequest) (*models.Tag, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockTagService) UpdateTag(ctx context.Context, userID, tagName string, req models.UpdateTagRequest) (*models.Tag, error) {
	args := m.Called(ctx, userID, tagName, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockTagService) DeleteTag(ctx context.Context, userID, tagName string) error {
	args := m.Called(ctx, userID, tagName)
	return args.Error(0)
}

func (m *MockTagService) ListTags(ctx context.Context, userID string, filter models.TagFilter) (*models.PaginatedResponse[models.TagResponse], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.TagResponse]), args.Error(1)
}

// MockUploadService implements UploadServiceInterface for testing
type MockUploadService struct {
	mock.Mock
}

func (m *MockUploadService) CreatePresignedUpload(ctx context.Context, userID string, req models.PresignedUploadRequest) (*models.PresignedUploadResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PresignedUploadResponse), args.Error(1)
}

func (m *MockUploadService) ConfirmUpload(ctx context.Context, userID, uploadID string) (*models.ConfirmUploadResponse, error) {
	args := m.Called(ctx, userID, uploadID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ConfirmUploadResponse), args.Error(1)
}

func (m *MockUploadService) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.UploadResponse], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.UploadResponse]), args.Error(1)
}

// MockSearchService implements SearchServiceInterface for testing
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Search(ctx context.Context, userID string, req models.SearchRequest) (*models.SearchResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SearchResponse), args.Error(1)
}

func (m *MockSearchService) GetSuggestions(ctx context.Context, userID, query string) (*models.AutocompleteResponse, error) {
	args := m.Called(ctx, userID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AutocompleteResponse), args.Error(1)
}

// MockStreamService implements StreamServiceInterface for testing
type MockStreamService struct {
	mock.Mock
}

func (m *MockStreamService) GetStreamURL(ctx context.Context, userID, trackID, quality string) (*models.StreamResponse, error) {
	args := m.Called(ctx, userID, trackID, quality)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StreamResponse), args.Error(1)
}

func (m *MockStreamService) GetDownloadURL(ctx context.Context, userID, trackID string) (*models.DownloadResponse, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DownloadResponse), args.Error(1)
}

func (m *MockStreamService) RecordPlayback(ctx context.Context, userID string, req models.RecordPlayRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockStreamService) GetQueue(ctx context.Context, userID string) (*models.PlayQueue, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlayQueue), args.Error(1)
}

func (m *MockStreamService) UpdateQueue(ctx context.Context, userID string, req models.UpdateQueueRequest) (*models.PlayQueue, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlayQueue), args.Error(1)
}

func (m *MockStreamService) QueueAction(ctx context.Context, userID string, req models.QueueActionRequest) (*models.PlayQueue, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlayQueue), args.Error(1)
}

// ====================
// Test Helpers
// ====================

func setupTestHandler() (*Handlers, *MockUserService, *MockTrackService, *MockAlbumService, *MockPlaylistService, *MockTagService, *MockUploadService, *MockSearchService, *MockStreamService) {
	mockUserService := new(MockUserService)
	mockTrackService := new(MockTrackService)
	mockAlbumService := new(MockAlbumService)
	mockPlaylistService := new(MockPlaylistService)
	mockTagService := new(MockTagService)
	mockUploadService := new(MockUploadService)
	mockSearchService := new(MockSearchService)
	mockStreamService := new(MockStreamService)

	h := NewHandlersWithServices(
		mockUserService,
		mockTrackService,
		mockAlbumService,
		mockPlaylistService,
		mockTagService,
		mockUploadService,
		mockSearchService,
		mockStreamService,
	)

	return h, mockUserService, mockTrackService, mockAlbumService, mockPlaylistService, mockTagService, mockUploadService, mockSearchService, mockStreamService
}

func createTestContext(method, path string, body string, userID string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user ID via X-User-ID header (simulates API Gateway authorizer)
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}

	return c, rec
}

// ====================
// User Handler Tests
// ====================

func TestGetCurrentUser_Success(t *testing.T) {
	h, mockUserService, _, _, _, _, _, _, _ := setupTestHandler()

	now := time.Now()
	expectedUser := &models.User{
		ID:          "user-123",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockUserService.On("GetUser", mock.Anything, "user-123").Return(expectedUser, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/me", "", "user-123")

	err := h.GetCurrentUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "Test User", response.DisplayName)

	mockUserService.AssertExpectations(t)
}

func TestGetCurrentUser_Unauthorized(t *testing.T) {
	h, _, _, _, _, _, _, _, _ := setupTestHandler()

	// No user ID provided
	c, rec := createTestContext(http.MethodGet, "/api/v1/me", "", "")

	err := h.GetCurrentUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetCurrentUser_NotFound(t *testing.T) {
	h, mockUserService, _, _, _, _, _, _, _ := setupTestHandler()

	mockUserService.On("GetUser", mock.Anything, "user-456").Return(nil, models.NewNotFoundError("User", "user-456"))

	c, rec := createTestContext(http.MethodGet, "/api/v1/me", "", "user-456")

	err := h.GetCurrentUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockUserService.AssertExpectations(t)
}

func TestUpdateCurrentUser_Success(t *testing.T) {
	h, mockUserService, _, _, _, _, _, _, _ := setupTestHandler()

	now := time.Now()
	displayName := "Updated Name"
	updatedUser := &models.User{
		ID:          "user-123",
		Email:       "test@example.com",
		DisplayName: displayName,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockUserService.On("UpdateUser", mock.Anything, "user-123", mock.AnythingOfType("models.UpdateUserRequest")).Return(updatedUser, nil)

	body := `{"displayName":"Updated Name"}`
	c, rec := createTestContext(http.MethodPut, "/api/v1/me", body, "user-123")

	err := h.UpdateCurrentUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.UserResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", response.DisplayName)

	mockUserService.AssertExpectations(t)
}

// ====================
// Track Handler Tests
// ====================

func TestListTracks_Success(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	expectedResult := &models.PaginatedResponse[models.TrackResponse]{
		Items: []models.TrackResponse{
			{ID: "track-1", Title: "Song 1", Artist: "Artist A"},
			{ID: "track-2", Title: "Song 2", Artist: "Artist B"},
		},
		Pagination: models.Pagination{Limit: 50},
	}

	mockTrackService.On("ListTracks", mock.Anything, "user-123", mock.AnythingOfType("models.TrackFilter")).Return(expectedResult, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/tracks", "", "user-123")

	err := h.ListTracks(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.PaginatedResponse[models.TrackResponse]
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Items, 2)
	assert.Equal(t, "track-1", response.Items[0].ID)

	mockTrackService.AssertExpectations(t)
}

func TestListTracks_Unauthorized(t *testing.T) {
	h, _, _, _, _, _, _, _, _ := setupTestHandler()

	c, rec := createTestContext(http.MethodGet, "/api/v1/tracks", "", "")

	err := h.ListTracks(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetTrack_Success(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	expectedTrack := &models.TrackResponse{
		ID:       "track-123",
		Title:    "Test Song",
		Artist:   "Test Artist",
		Album:    "Test Album",
		Duration: 180,
	}

	mockTrackService.On("GetTrack", mock.Anything, "user-123", "track-123").Return(expectedTrack, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/tracks/track-123", "", "user-123")
	c.SetParamNames("id")
	c.SetParamValues("track-123")

	err := h.GetTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.TrackResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "track-123", response.ID)
	assert.Equal(t, "Test Song", response.Title)

	mockTrackService.AssertExpectations(t)
}

func TestGetTrack_NotFound(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	mockTrackService.On("GetTrack", mock.Anything, "user-123", "nonexistent").Return(nil, models.NewNotFoundError("Track", "nonexistent"))

	c, rec := createTestContext(http.MethodGet, "/api/v1/tracks/nonexistent", "", "user-123")
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := h.GetTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockTrackService.AssertExpectations(t)
}

func TestGetTrack_MissingID(t *testing.T) {
	h, _, _, _, _, _, _, _, _ := setupTestHandler()

	c, rec := createTestContext(http.MethodGet, "/api/v1/tracks/", "", "user-123")
	c.SetParamNames("id")
	c.SetParamValues("")

	err := h.GetTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateTrack_Success(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	updatedTrack := &models.TrackResponse{
		ID:     "track-123",
		Title:  "Updated Title",
		Artist: "Test Artist",
	}

	mockTrackService.On("UpdateTrack", mock.Anything, "user-123", "track-123", mock.AnythingOfType("models.UpdateTrackRequest")).Return(updatedTrack, nil)

	body := `{"title":"Updated Title"}`
	c, rec := createTestContext(http.MethodPut, "/api/v1/tracks/track-123", body, "user-123")
	c.SetParamNames("id")
	c.SetParamValues("track-123")

	err := h.UpdateTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.TrackResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", response.Title)

	mockTrackService.AssertExpectations(t)
}

func TestDeleteTrack_Success(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	mockTrackService.On("DeleteTrack", mock.Anything, "user-123", "track-123").Return(nil)

	c, rec := createTestContext(http.MethodDelete, "/api/v1/tracks/track-123", "", "user-123")
	c.SetParamNames("id")
	c.SetParamValues("track-123")

	err := h.DeleteTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	mockTrackService.AssertExpectations(t)
}

func TestAddTagsToTrack_Success(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	updatedTrack := &models.TrackResponse{
		ID:    "track-123",
		Title: "Test Song",
		Tags:  []string{"rock", "favorites"},
	}

	mockTrackService.On("AddTagsToTrack", mock.Anything, "user-123", "track-123", []string{"rock", "favorites"}).Return(updatedTrack, nil)

	body := `{"tags":["rock","favorites"]}`
	c, rec := createTestContext(http.MethodPost, "/api/v1/tracks/track-123/tags", body, "user-123")
	c.SetParamNames("id")
	c.SetParamValues("track-123")

	err := h.AddTagsToTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.TrackResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Tags, "rock")
	assert.Contains(t, response.Tags, "favorites")

	mockTrackService.AssertExpectations(t)
}

func TestRemoveTagFromTrack_Success(t *testing.T) {
	h, _, mockTrackService, _, _, _, _, _, _ := setupTestHandler()

	updatedTrack := &models.TrackResponse{
		ID:    "track-123",
		Title: "Test Song",
		Tags:  []string{"rock"},
	}

	mockTrackService.On("RemoveTagFromTrack", mock.Anything, "user-123", "track-123", "favorites").Return(updatedTrack, nil)

	c, rec := createTestContext(http.MethodDelete, "/api/v1/tracks/track-123/tags/favorites", "", "user-123")
	c.SetParamNames("id", "tag")
	c.SetParamValues("track-123", "favorites")

	err := h.RemoveTagFromTrack(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.TrackResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotContains(t, response.Tags, "favorites")

	mockTrackService.AssertExpectations(t)
}

// ====================
// Album Handler Tests
// ====================

func TestListAlbums_Success(t *testing.T) {
	h, _, _, mockAlbumService, _, _, _, _, _ := setupTestHandler()

	expectedResult := &models.PaginatedResponse[models.AlbumResponse]{
		Items: []models.AlbumResponse{
			{ID: "album-1", Title: "Album 1", Artist: "Artist A"},
			{ID: "album-2", Title: "Album 2", Artist: "Artist B"},
		},
		Pagination: models.Pagination{Limit: 50},
	}

	mockAlbumService.On("ListAlbums", mock.Anything, "user-123", mock.AnythingOfType("models.AlbumFilter")).Return(expectedResult, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/albums", "", "user-123")

	err := h.ListAlbums(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockAlbumService.AssertExpectations(t)
}

func TestGetAlbum_Success(t *testing.T) {
	h, _, _, mockAlbumService, _, _, _, _, _ := setupTestHandler()

	expectedAlbum := &models.AlbumWithTracks{
		Album: models.AlbumResponse{
			ID:     "album-123",
			Title:  "Test Album",
			Artist: "Test Artist",
		},
		Tracks: []models.TrackResponse{
			{ID: "track-1", Title: "Song 1"},
		},
	}

	mockAlbumService.On("GetAlbumWithTracks", mock.Anything, "user-123", "album-123").Return(expectedAlbum, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/albums/album-123", "", "user-123")
	c.SetParamNames("id")
	c.SetParamValues("album-123")

	err := h.GetAlbum(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockAlbumService.AssertExpectations(t)
}

// ====================
// Playlist Handler Tests
// ====================

func TestListPlaylists_Success(t *testing.T) {
	h, _, _, _, mockPlaylistService, _, _, _, _ := setupTestHandler()

	expectedResult := &models.PaginatedResponse[models.PlaylistResponse]{
		Items: []models.PlaylistResponse{
			{ID: "playlist-1", Name: "My Favorites"},
			{ID: "playlist-2", Name: "Chill Vibes"},
		},
		Pagination: models.Pagination{Limit: 50},
	}

	mockPlaylistService.On("ListPlaylists", mock.Anything, "user-123", mock.AnythingOfType("models.PlaylistFilter")).Return(expectedResult, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/playlists", "", "user-123")

	err := h.ListPlaylists(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockPlaylistService.AssertExpectations(t)
}

func TestCreatePlaylist_Success(t *testing.T) {
	h, _, _, _, mockPlaylistService, _, _, _, _ := setupTestHandler()

	expectedPlaylist := &models.PlaylistResponse{
		ID:          "playlist-new",
		Name:        "New Playlist",
		Description: "My new playlist",
	}

	mockPlaylistService.On("CreatePlaylist", mock.Anything, "user-123", mock.AnythingOfType("models.CreatePlaylistRequest")).Return(expectedPlaylist, nil)

	body := `{"name":"New Playlist","description":"My new playlist"}`
	c, rec := createTestContext(http.MethodPost, "/api/v1/playlists", body, "user-123")

	err := h.CreatePlaylist(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	mockPlaylistService.AssertExpectations(t)
}

// ====================
// Upload Handler Tests
// ====================

func TestGetPresignedUploadURL_Success(t *testing.T) {
	h, _, _, _, _, _, mockUploadService, _, _ := setupTestHandler()

	expectedResponse := &models.PresignedUploadResponse{
		UploadID:  "upload-123",
		UploadURL: "https://s3.amazonaws.com/bucket/key?signature=...",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	mockUploadService.On("CreatePresignedUpload", mock.Anything, "user-123", mock.AnythingOfType("models.PresignedUploadRequest")).Return(expectedResponse, nil)

	body := `{"fileName":"test.mp3","fileSize":5242880,"contentType":"audio/mpeg"}`
	c, rec := createTestContext(http.MethodPost, "/api/v1/upload/presigned", body, "user-123")

	err := h.GetPresignedUploadURL(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockUploadService.AssertExpectations(t)
}

// ====================
// Search Handler Tests
// ====================

func TestSearchSimple_Success(t *testing.T) {
	h, _, _, _, _, _, _, mockSearchService, _ := setupTestHandler()

	expectedResponse := &models.SearchResponse{
		Query:        "test",
		TotalResults: 2,
		Tracks: []models.TrackResponse{
			{ID: "track-1", Title: "Test Song"},
		},
	}

	mockSearchService.On("Search", mock.Anything, "user-123", mock.AnythingOfType("models.SearchRequest")).Return(expectedResponse, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/search?q=test", "", "user-123")
	c.QueryParams().Set("q", "test")

	err := h.SearchSimple(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockSearchService.AssertExpectations(t)
}

// ====================
// Stream Handler Tests
// ====================

func TestGetStreamURL_Success(t *testing.T) {
	h, _, _, _, _, _, _, _, mockStreamService := setupTestHandler()

	expectedResponse := &models.StreamResponse{
		TrackID:   "track-123",
		StreamURL: "https://cdn.example.com/stream/track-123",
		ExpiresAt: time.Now().Add(4 * time.Hour),
	}

	mockStreamService.On("GetStreamURL", mock.Anything, "user-123", "track-123", "original").Return(expectedResponse, nil)

	c, rec := createTestContext(http.MethodGet, "/api/v1/stream/track-123", "", "user-123")
	c.SetParamNames("trackId")
	c.SetParamValues("track-123")

	err := h.GetStreamURL(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockStreamService.AssertExpectations(t)
}

func TestRecordPlayback_Success(t *testing.T) {
	h, _, _, _, _, _, _, _, mockStreamService := setupTestHandler()

	mockStreamService.On("RecordPlayback", mock.Anything, "user-123", mock.AnythingOfType("models.RecordPlayRequest")).Return(nil)

	body := `{"trackId":"track-123","duration":120,"completed":false}`
	c, rec := createTestContext(http.MethodPost, "/api/v1/playback/record", body, "user-123")

	err := h.RecordPlayback(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	mockStreamService.AssertExpectations(t)
}

// ====================
// Error Handling Tests
// ====================

func TestHandleError_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleError(c, models.NewNotFoundError("Track", "track-123"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleError_BadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleError(c, models.ErrBadRequest)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleError_InternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleError(c, assert.AnError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ====================
// Helper Function Tests
// ====================

func TestGetUserIDFromContext_FromClaims(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Cognito JWT claims
	claims := map[string]interface{}{
		"sub": "user-from-claims",
	}
	c.Set("claims", claims)

	userID := getUserIDFromContext(c)
	assert.Equal(t, "user-from-claims", userID)
}

func TestGetUserIDFromContext_FromHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-from-header")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := getUserIDFromContext(c)
	assert.Equal(t, "user-from-header", userID)
}

func TestGetUserIDFromContext_Empty(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := getUserIDFromContext(c)
	assert.Empty(t, userID)
}
