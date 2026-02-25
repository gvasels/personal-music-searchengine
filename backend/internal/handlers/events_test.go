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
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockEventsService implements the EventsService interface for testing.
type MockEventsService struct {
	mock.Mock
}

func (m *MockEventsService) GetArtistEvents(ctx context.Context, artistName string) (*models.ArtistEventsResponse, error) {
	args := m.Called(ctx, artistName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistEventsResponse), args.Error(1)
}

func (m *MockEventsService) SearchArtistEvents(ctx context.Context, query string, limit int) ([]models.ArtistSearchResult, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ArtistSearchResult), args.Error(1)
}

func setupEventsTestHandler(mockEvents *MockEventsService) (*echo.Echo, *EventsHandler) {
	e := echo.New()
	h := NewEventsHandler(mockEvents)
	return e, h
}

func TestGetArtistEvents_Success(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	eventDate := time.Date(2026, 6, 15, 20, 0, 0, 0, time.UTC)
	expectedResponse := &models.ArtistEventsResponse{
		ArtistName: "deadmau5",
		Events: []models.Event{
			{
				ID:         "evt-1",
				ArtistName: "deadmau5",
				Title:      "deadmau5 Live at Red Rocks",
				Date:       eventDate,
				Venue:      "Red Rocks Amphitheatre",
				City:       "Morrison",
				Region:     "CO",
				Country:    "US",
				TicketURL:  "https://tickets.example.com/evt-1",
				Status:     "confirmed",
				Source:     "mock",
			},
			{
				ID:         "evt-2",
				ArtistName: "deadmau5",
				Title:      "deadmau5 @ Ultra Music Festival",
				Date:       eventDate.Add(30 * 24 * time.Hour),
				Venue:      "Bayfront Park",
				City:       "Miami",
				Region:     "FL",
				Country:    "US",
				Status:     "confirmed",
				Source:     "mock",
			},
		},
		TotalCount: 2,
		Source:     "mock",
	}

	mockEvents.On("GetArtistEvents", mock.Anything, "deadmau5").Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/deadmau5/events", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")
	c.Set("user_id", "user-123")

	err := h.GetArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.ArtistEventsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "deadmau5", response.ArtistName)
	assert.Len(t, response.Events, 2)
	assert.Equal(t, 2, response.TotalCount)
	assert.Equal(t, "deadmau5 Live at Red Rocks", response.Events[0].Title)

	mockEvents.AssertExpectations(t)
}

func TestGetArtistEvents_EmptyResult(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	expectedResponse := &models.ArtistEventsResponse{
		ArtistName: "unknown-artist",
		Events:     []models.Event{},
		TotalCount: 0,
		Source:     "mock",
	}

	mockEvents.On("GetArtistEvents", mock.Anything, "unknown-artist").Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/unknown-artist/events", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("unknown-artist")
	c.Set("user_id", "user-123")

	err := h.GetArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.ArtistEventsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "unknown-artist", response.ArtistName)
	assert.Empty(t, response.Events)
	assert.Equal(t, 0, response.TotalCount)

	mockEvents.AssertExpectations(t)
}

func TestGetArtistEvents_Unauthenticated(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/deadmau5/events", nil)
	// No auth headers or context
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")

	err := h.GetArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetArtistEvents_EmptyArtistName(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists//events", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("")
	c.Set("user_id", "user-123")

	err := h.GetArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetArtistEvents_ServiceError(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	mockEvents.On("GetArtistEvents", mock.Anything, "deadmau5").Return(nil, errors.New("provider unavailable"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/deadmau5/events", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("deadmau5")
	c.Set("user_id", "user-123")

	err := h.GetArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockEvents.AssertExpectations(t)
}

func TestSearchArtistEvents_Success(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	expectedResults := []models.ArtistSearchResult{
		{
			Name:           "deadmau5",
			ImageURL:       "https://example.com/deadmau5.jpg",
			UpcomingEvents: 5,
			Source:         "mock",
		},
		{
			Name:           "Daft Punk",
			ImageURL:       "https://example.com/daftpunk.jpg",
			UpcomingEvents: 2,
			Source:         "mock",
		},
	}

	mockEvents.On("SearchArtistEvents", mock.Anything, "dea", 10).Return(expectedResults, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search?q=dea", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	items := response["items"].([]interface{})
	assert.Len(t, items, 2)

	firstItem := items[0].(map[string]interface{})
	assert.Equal(t, "deadmau5", firstItem["name"])
	assert.Equal(t, float64(5), firstItem["upcomingEvents"])

	mockEvents.AssertExpectations(t)
}

func TestSearchArtistEvents_DefaultLimit(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	// When no limit is specified, handler should default to 10
	mockEvents.On("SearchArtistEvents", mock.Anything, "tiesto", 10).Return([]models.ArtistSearchResult{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search?q=tiesto", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockEvents.AssertExpectations(t)
}

func TestSearchArtistEvents_CustomLimit(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	mockEvents.On("SearchArtistEvents", mock.Anything, "avicii", 5).Return([]models.ArtistSearchResult{
		{Name: "Avicii", UpcomingEvents: 0, Source: "mock"},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search?q=avicii&limit=5", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockEvents.AssertExpectations(t)
}

func TestSearchArtistEvents_EmptyQuery(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search?q=", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSearchArtistEvents_MissingQuery(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSearchArtistEvents_Unauthenticated(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search?q=deadmau5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSearchArtistEvents_ServiceError(t *testing.T) {
	mockEvents := new(MockEventsService)
	e, h := setupEventsTestHandler(mockEvents)

	mockEvents.On("SearchArtistEvents", mock.Anything, "deadmau5", 10).Return(nil, errors.New("provider unavailable"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/search?q=deadmau5", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	err := h.SearchArtistEvents(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockEvents.AssertExpectations(t)
}
