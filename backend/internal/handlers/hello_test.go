package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockHelloService implements the HelloService interface for testing.
type MockHelloService struct {
	mock.Mock
}

func (m *MockHelloService) Search(ctx context.Context, query string) ([]service.HelloTrack, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]service.HelloTrack), args.Error(1)
}

func (m *MockHelloService) Featured(ctx context.Context, limit int) ([]service.HelloTrack, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]service.HelloTrack), args.Error(1)
}

func setupHelloTestHandler(mockHello *MockHelloService) (*echo.Echo, *HelloHandler) {
	e := echo.New()
	h := NewHelloHandler(mockHello)
	return e, h
}

// TestHelloHandler_Health verifies GET /api/v1/hello/health returns 200
// with {"status":"ok","service":"hello"}.
func TestHelloHandler_Health(t *testing.T) {
	mockHello := new(MockHelloService)
	e, h := setupHelloTestHandler(mockHello)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Health(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "hello", response["service"])
}

// TestHelloHandler_Search_WithQuery verifies GET /api/v1/hello/search?q=jazz
// returns 200 with {"items":[...],"total":N}.
func TestHelloHandler_Search_WithQuery(t *testing.T) {
	mockHello := new(MockHelloService)
	e, h := setupHelloTestHandler(mockHello)

	tracks := []service.HelloTrack{
		{ID: "track-1", Title: "Jazz Vibes", Artist: "Artist A"},
		{ID: "track-2", Title: "Jazz Fusion", Artist: "Artist B"},
	}
	mockHello.On("Search", mock.Anything, "jazz").Return(tracks, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=jazz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Search(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response must contain 'items' array")
	assert.Len(t, items, 2)
	assert.Equal(t, float64(2), response["total"])

	mockHello.AssertExpectations(t)
}

// TestHelloHandler_Search_ReadsQParam verifies the handler specifically reads
// the "q" query parameter (NOT "query").
func TestHelloHandler_Search_ReadsQParam(t *testing.T) {
	mockHello := new(MockHelloService)
	e, h := setupHelloTestHandler(mockHello)

	// Only set "q" param; the service should receive "rock"
	tracks := []service.HelloTrack{
		{ID: "track-1", Title: "Rock Anthem", Artist: "Band X"},
	}
	mockHello.On("Search", mock.Anything, "rock").Return(tracks, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=rock", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Search(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the mock was called with "rock" (from q param)
	mockHello.AssertCalled(t, "Search", mock.Anything, "rock")

	// Also verify that "query" param is NOT used by testing that a request
	// with only "query" param (no "q") yields empty results.
	mockHello2 := new(MockHelloService)
	_, h2 := setupHelloTestHandler(mockHello2)

	// When q is empty, service should receive ""
	mockHello2.On("Search", mock.Anything, "").Return([]service.HelloTrack{}, nil)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?query=rock", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	err = h2.Search(c2)
	require.NoError(t, err)

	// The handler should have called Search with empty string (reading "q", not "query")
	mockHello2.AssertCalled(t, "Search", mock.Anything, "")
}

// TestHelloHandler_Search_EmptyQuery verifies GET /api/v1/hello/search
// (no q param) returns {"items":[],"total":0}.
func TestHelloHandler_Search_EmptyQuery(t *testing.T) {
	mockHello := new(MockHelloService)
	e, h := setupHelloTestHandler(mockHello)

	mockHello.On("Search", mock.Anything, "").Return([]service.HelloTrack{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Search(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response must contain 'items' array")
	assert.Len(t, items, 0)
	assert.Equal(t, float64(0), response["total"])

	mockHello.AssertExpectations(t)
}

// TestHelloHandler_Featured_ReturnsItems verifies GET /api/v1/hello/featured
// returns 200 with {"items":[...],"total":N}.
func TestHelloHandler_Featured_ReturnsItems(t *testing.T) {
	mockHello := new(MockHelloService)
	e, h := setupHelloTestHandler(mockHello)

	tracks := []service.HelloTrack{
		{ID: "track-1", Title: "Featured Hit", Artist: "Star"},
		{ID: "track-2", Title: "Top Song", Artist: "Legend"},
		{ID: "track-3", Title: "Classic", Artist: "Master"},
	}
	// Default limit when not specified (expect 10 as default)
	mockHello.On("Featured", mock.Anything, 10).Return(tracks, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Featured(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response must contain 'items' array")
	assert.Len(t, items, 3)
	assert.Equal(t, float64(3), response["total"])

	mockHello.AssertExpectations(t)
}

// TestHelloHandler_Featured_WithLimit verifies GET /api/v1/hello/featured?limit=5
// passes limit=5 to the service.
func TestHelloHandler_Featured_WithLimit(t *testing.T) {
	mockHello := new(MockHelloService)
	e, h := setupHelloTestHandler(mockHello)

	tracks := []service.HelloTrack{
		{ID: "track-1", Title: "Song One", Artist: "Artist 1"},
		{ID: "track-2", Title: "Song Two", Artist: "Artist 2"},
	}
	mockHello.On("Featured", mock.Anything, 5).Return(tracks, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured?limit=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Featured(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the service was called with limit=5
	mockHello.AssertCalled(t, "Featured", mock.Anything, 5)

	mockHello.AssertExpectations(t)
}

// TestHelloHandler_ResponseUsesItemsKey verifies that all list responses
// use the JSON key "items" (NOT "tracks" or any other key).
func TestHelloHandler_ResponseUsesItemsKey(t *testing.T) {
	t.Run("search response uses items key", func(t *testing.T) {
		mockHello := new(MockHelloService)
		e, h := setupHelloTestHandler(mockHello)

		tracks := []service.HelloTrack{
			{ID: "track-1", Title: "Test", Artist: "Test"},
		}
		mockHello.On("Search", mock.Anything, "test").Return(tracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Search(c)
		require.NoError(t, err)

		var raw map[string]json.RawMessage
		err = json.Unmarshal(rec.Body.Bytes(), &raw)
		require.NoError(t, err)

		_, hasItems := raw["items"]
		assert.True(t, hasItems, "response JSON must contain 'items' key")

		_, hasTracks := raw["tracks"]
		assert.False(t, hasTracks, "response JSON must NOT contain 'tracks' key")

		_, hasTotal := raw["total"]
		assert.True(t, hasTotal, "response JSON must contain 'total' key")
	})

	t.Run("featured response uses items key", func(t *testing.T) {
		mockHello := new(MockHelloService)
		e, h := setupHelloTestHandler(mockHello)

		tracks := []service.HelloTrack{
			{ID: "track-1", Title: "Featured", Artist: "Star"},
		}
		mockHello.On("Featured", mock.Anything, 10).Return(tracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Featured(c)
		require.NoError(t, err)

		var raw map[string]json.RawMessage
		err = json.Unmarshal(rec.Body.Bytes(), &raw)
		require.NoError(t, err)

		_, hasItems := raw["items"]
		assert.True(t, hasItems, "response JSON must contain 'items' key")

		_, hasTracks := raw["tracks"]
		assert.False(t, hasTracks, "response JSON must NOT contain 'tracks' key")

		_, hasTotal := raw["total"]
		assert.True(t, hasTotal, "response JSON must contain 'total' key")
	})
}
