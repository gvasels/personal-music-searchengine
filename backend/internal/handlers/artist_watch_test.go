package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockArtistWatchService implements the ArtistWatchService interface for testing.
type MockArtistWatchService struct {
	mock.Mock
}

func (m *MockArtistWatchService) WatchArtist(ctx context.Context, userID, artistName string) (*models.ArtistWatchResponse, error) {
	args := m.Called(ctx, userID, artistName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistWatchResponse), args.Error(1)
}

func (m *MockArtistWatchService) UnwatchArtist(ctx context.Context, userID, artistName string) error {
	args := m.Called(ctx, userID, artistName)
	return args.Error(0)
}

func (m *MockArtistWatchService) GetWatchStatus(ctx context.Context, userID, artistName string) (bool, error) {
	args := m.Called(ctx, userID, artistName)
	return args.Bool(0), args.Error(1)
}

func (m *MockArtistWatchService) ListWatchedArtists(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.ArtistWatchResponse], error) {
	args := m.Called(ctx, userID, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.ArtistWatchResponse]), args.Error(1)
}

func setupArtistWatchTestHandler(mockWatch *MockArtistWatchService) (*echo.Echo, *ArtistWatchHandler) {
	e := echo.New()
	h := NewArtistWatchHandler(mockWatch)
	return e, h
}

func TestWatchArtist_Success(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	now := time.Now()
	expectedResponse := &models.ArtistWatchResponse{
		ArtistName: "deadmau5",
		WatchedAt:  now,
	}

	mockWatch.On("WatchArtist", mock.Anything, "user-123", "deadmau5").Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/deadmau5/watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")
	c.Set("user_id", "user-123")

	err := h.WatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "deadmau5", response["artistName"])
	assert.Equal(t, true, response["watching"])

	mockWatch.AssertExpectations(t)
}

func TestWatchArtist_Unauthenticated(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/deadmau5/watch", nil)
	// No X-User-ID header and no user_id in context
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")

	err := h.WatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestWatchArtist_EmptyArtistName(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/artists//watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("")
	c.Set("user_id", "user-123")

	err := h.WatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWatchArtist_ServiceError(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	mockWatch.On("WatchArtist", mock.Anything, "user-123", "deadmau5").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/deadmau5/watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")
	c.Set("user_id", "user-123")

	err := h.WatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockWatch.AssertExpectations(t)
}

func TestUnwatchArtist_Success(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	mockWatch.On("UnwatchArtist", mock.Anything, "user-123", "deadmau5").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/artists/deadmau5/watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")
	c.Set("user_id", "user-123")

	err := h.UnwatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	mockWatch.AssertExpectations(t)
}

func TestUnwatchArtist_Unauthenticated(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/artists/deadmau5/watch", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")

	err := h.UnwatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUnwatchArtist_NotFound(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	mockWatch.On("UnwatchArtist", mock.Anything, "user-123", "unknown-artist").Return(models.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/artists/unknown-artist/watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("unknown-artist")
	c.Set("user_id", "user-123")

	err := h.UnwatchArtist(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockWatch.AssertExpectations(t)
}

func TestGetWatchStatus_Watching(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	mockWatch.On("GetWatchStatus", mock.Anything, "user-123", "deadmau5").Return(true, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/deadmau5/watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")
	c.Set("user_id", "user-123")

	err := h.GetWatchStatus(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, true, response["watching"])
	assert.Equal(t, "deadmau5", response["artistName"])

	mockWatch.AssertExpectations(t)
}

func TestGetWatchStatus_NotWatching(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	mockWatch.On("GetWatchStatus", mock.Anything, "user-123", "skrillex").Return(false, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/skrillex/watch", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("skrillex")
	c.Set("user_id", "user-123")

	err := h.GetWatchStatus(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, false, response["watching"])
	assert.Equal(t, "skrillex", response["artistName"])

	mockWatch.AssertExpectations(t)
}

func TestGetWatchStatus_Unauthenticated(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/deadmau5/watch", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")

	err := h.GetWatchStatus(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListWatchedArtists_Success(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	now := time.Now()
	expectedResult := &repository.PaginatedResult[models.ArtistWatchResponse]{
		Items: []models.ArtistWatchResponse{
			{ArtistName: "deadmau5", WatchedAt: now},
			{ArtistName: "skrillex", WatchedAt: now.Add(-time.Hour)},
			{ArtistName: "daft punk", WatchedAt: now.Add(-2 * time.Hour)},
		},
		NextCursor: "cursor-abc",
		HasMore:    true,
	}

	mockWatch.On("ListWatchedArtists", mock.Anything, "user-123", 20, "").Return(expectedResult, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/watched-artists", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.ListWatchedArtists(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items := response["items"].([]interface{})
	assert.Len(t, items, 3)
	assert.Equal(t, "cursor-abc", response["nextCursor"])
	assert.Equal(t, true, response["hasMore"])

	mockWatch.AssertExpectations(t)
}

func TestListWatchedArtists_WithCustomLimit(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	expectedResult := &repository.PaginatedResult[models.ArtistWatchResponse]{
		Items:   []models.ArtistWatchResponse{},
		HasMore: false,
	}

	mockWatch.On("ListWatchedArtists", mock.Anything, "user-123", 5, "").Return(expectedResult, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/watched-artists?limit=5", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.ListWatchedArtists(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockWatch.AssertExpectations(t)
}

func TestListWatchedArtists_WithCursor(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	now := time.Now()
	expectedResult := &repository.PaginatedResult[models.ArtistWatchResponse]{
		Items: []models.ArtistWatchResponse{
			{ArtistName: "tiesto", WatchedAt: now},
		},
		HasMore: false,
	}

	mockWatch.On("ListWatchedArtists", mock.Anything, "user-123", 20, "some-cursor").Return(expectedResult, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/watched-artists?cursor=some-cursor", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.ListWatchedArtists(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockWatch.AssertExpectations(t)
}

func TestListWatchedArtists_Unauthenticated(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/watched-artists", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListWatchedArtists(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListWatchedArtists_ServiceError(t *testing.T) {
	mockWatch := new(MockArtistWatchService)
	e, h := setupArtistWatchTestHandler(mockWatch)

	mockWatch.On("ListWatchedArtists", mock.Anything, "user-123", 20, "").Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/watched-artists", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.ListWatchedArtists(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockWatch.AssertExpectations(t)
}
