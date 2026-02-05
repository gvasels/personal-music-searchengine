package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// HelloTrack mirrors service.HelloTrack for testing purposes.
// This allows tests to compile before the service package is implemented.
// The handler will use the real service.HelloTrack type.
type HelloTrack struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	Year     int    `json:"year"`
	Duration int    `json:"duration"`
}

// MockHelloService is a mock implementation of HelloServiceInterface
type MockHelloService struct {
	mock.Mock
}

func (m *MockHelloService) Search(ctx context.Context, query string) ([]HelloTrack, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]HelloTrack), args.Error(1)
}

func (m *MockHelloService) Featured(ctx context.Context, limit int) ([]HelloTrack, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]HelloTrack), args.Error(1)
}

// TestHelloHandler_Health tests the health endpoint
func TestHelloHandler_Health(t *testing.T) {
	e := echo.New()
	mockService := new(MockHelloService)
	handler := NewHelloHandler(mockService)

	t.Run("returns status ok and service name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Health(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "ok", response["status"])
		assert.Equal(t, "hello", response["service"])
	})
}

// TestHelloHandler_Search tests the search endpoint
func TestHelloHandler_Search(t *testing.T) {
	e := echo.New()

	t.Run("calls service with q parameter value", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		tracks := []HelloTrack{
			{ID: "track-1", Title: "Jazz Song", Artist: "Jazz Artist", Album: "Jazz Album", Genre: "jazz", Year: 2023, Duration: 180},
		}
		mockService.On("Search", mock.Anything, "jazz").Return(tracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=jazz", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		mockService.AssertCalled(t, "Search", mock.Anything, "jazz")
	})

	t.Run("ignores query parameter - only reads q", func(t *testing.T) {
		// CRITICAL CONTRACT TEST: The param is "q", NOT "query"
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		// Service should be called with empty string because "query" param is ignored
		mockService.On("Search", mock.Anything, "").Return([]HelloTrack{}, nil)

		// Using "query" param instead of "q" - should be ignored
		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?query=jazz", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)

		// Verify service was called with empty string (query param ignored)
		mockService.AssertCalled(t, "Search", mock.Anything, "")
	})

	t.Run("returns items key not tracks", func(t *testing.T) {
		// CRITICAL CONTRACT TEST: Response must have "items" key, NOT "tracks"
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		tracks := []HelloTrack{
			{ID: "track-1", Title: "Test Song", Artist: "Test Artist", Album: "Test Album", Genre: "pop", Year: 2024, Duration: 200},
		}
		mockService.On("Search", mock.Anything, "test").Return(tracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Must have "items" key
		_, hasItems := response["items"]
		assert.True(t, hasItems, "Response must have 'items' key")

		// Must NOT have "tracks" key
		_, hasTracks := response["tracks"]
		assert.False(t, hasTracks, "Response must NOT have 'tracks' key")

		// Verify total is present
		total, hasTotal := response["total"]
		assert.True(t, hasTotal, "Response must have 'total' key")
		assert.Equal(t, float64(1), total)
	})

	t.Run("empty query returns empty items array", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Search", mock.Anything, "").Return([]HelloTrack{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			Items []HelloTrack `json:"items"`
			Total int          `json:"total"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Empty(t, response.Items)
		assert.Equal(t, 0, response.Total)
	})

	t.Run("no q parameter returns empty items array", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Search", mock.Anything, "").Return([]HelloTrack{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			Items []HelloTrack `json:"items"`
			Total int          `json:"total"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Empty(t, response.Items)
		assert.Equal(t, 0, response.Total)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Search", mock.Anything, "error").Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=error", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
	})

	t.Run("search results contain all track fields", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		tracks := []HelloTrack{
			{
				ID:       "seed-t1",
				Title:    "Aurora Borealis",
				Artist:   "Aurora Waves",
				Album:    "Dreamscape",
				Genre:    "jazz",
				Year:     2023,
				Duration: 245,
			},
		}
		mockService.On("Search", mock.Anything, "aurora").Return(tracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=aurora", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Search(c)
		require.NoError(t, err)

		var response struct {
			Items []struct {
				ID       string `json:"id"`
				Title    string `json:"title"`
				Artist   string `json:"artist"`
				Album    string `json:"album"`
				Genre    string `json:"genre"`
				Year     int    `json:"year"`
				Duration int    `json:"duration"`
			} `json:"items"`
			Total int `json:"total"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Len(t, response.Items, 1)
		item := response.Items[0]
		assert.Equal(t, "seed-t1", item.ID)
		assert.Equal(t, "Aurora Borealis", item.Title)
		assert.Equal(t, "Aurora Waves", item.Artist)
		assert.Equal(t, "Dreamscape", item.Album)
		assert.Equal(t, "jazz", item.Genre)
		assert.Equal(t, 2023, item.Year)
		assert.Equal(t, 245, item.Duration)
	})
}

// TestHelloHandler_Featured tests the featured endpoint
func TestHelloHandler_Featured(t *testing.T) {
	e := echo.New()

	t.Run("returns featured tracks with items key", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		tracks := []HelloTrack{
			{ID: "track-1", Title: "Featured Song 1", Artist: "Artist 1", Album: "Album 1", Genre: "pop", Year: 2024, Duration: 180},
			{ID: "track-2", Title: "Featured Song 2", Artist: "Artist 2", Album: "Album 2", Genre: "rock", Year: 2023, Duration: 200},
		}
		mockService.On("Featured", mock.Anything, 10).Return(tracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Featured(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			Items []HelloTrack `json:"items"`
			Total int          `json:"total"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Items, 2)
		assert.Equal(t, 2, response.Total)
	})

	t.Run("default limit is 10 when no param", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Featured", mock.Anything, 10).Return([]HelloTrack{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Featured(c)
		require.NoError(t, err)

		// Verify service was called with default limit of 10
		mockService.AssertCalled(t, "Featured", mock.Anything, 10)
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Featured", mock.Anything, 5).Return([]HelloTrack{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured?limit=5", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Featured(c)
		require.NoError(t, err)

		mockService.AssertCalled(t, "Featured", mock.Anything, 5)
	})

	t.Run("invalid limit uses default", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Featured", mock.Anything, 10).Return([]HelloTrack{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured?limit=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Featured(c)
		require.NoError(t, err)

		// Invalid limit should fall back to default of 10
		mockService.AssertCalled(t, "Featured", mock.Anything, 10)
	})

	t.Run("negative limit uses default", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Featured", mock.Anything, 10).Return([]HelloTrack{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured?limit=-5", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Featured(c)
		require.NoError(t, err)

		// Negative limit should fall back to default of 10
		mockService.AssertCalled(t, "Featured", mock.Anything, 10)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		mockService := new(MockHelloService)
		handler := NewHelloHandler(mockService)

		mockService.On("Featured", mock.Anything, 10).Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Featured(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
	})
}

// TestRegisterHelloRoutes tests route registration
func TestRegisterHelloRoutes(t *testing.T) {
	e := echo.New()
	mockService := new(MockHelloService)
	handler := NewHelloHandler(mockService)

	RegisterHelloRoutes(e, handler)

	// Get all registered routes
	routes := e.Routes()

	// Build a map for easier lookup
	routeMap := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = true
	}

	t.Run("health route is registered", func(t *testing.T) {
		assert.True(t, routeMap["GET /api/v1/hello/health"], "GET /api/v1/hello/health should be registered")
	})

	t.Run("search route is registered", func(t *testing.T) {
		assert.True(t, routeMap["GET /api/v1/hello/search"], "GET /api/v1/hello/search should be registered")
	})

	t.Run("featured route is registered", func(t *testing.T) {
		assert.True(t, routeMap["GET /api/v1/hello/featured"], "GET /api/v1/hello/featured should be registered")
	})
}
