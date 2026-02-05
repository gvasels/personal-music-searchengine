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
	"github.com/stretchr/testify/require"
)

// mockHelloRepo implements service.HelloRepository for handler tests
type mockHelloRepo struct {
	tracks []service.HelloTrack
	err    error
}

func (m *mockHelloRepo) GetTracksByUser(ctx context.Context, userID string) ([]service.HelloTrack, error) {
	return m.tracks, m.err
}

// helloTestTracks returns sample tracks for handler tests
func helloTestTracks() []service.HelloTrack {
	return []service.HelloTrack{
		{
			ID:       "track-1",
			Title:    "Neon Lights",
			Artist:   "Kraftwerk",
			Album:    "The Man-Machine",
			Genre:    "Electronic",
			Year:     1978,
			Duration: 540,
		},
		{
			ID:       "track-2",
			Title:    "Blue Monday",
			Artist:   "New Order",
			Album:    "Power, Corruption & Lies",
			Genre:    "Synth-Pop",
			Year:     1983,
			Duration: 440,
		},
		{
			ID:       "track-3",
			Title:    "Enjoy the Silence",
			Artist:   "Depeche Mode",
			Album:    "Violator",
			Genre:    "Electronic",
			Year:     1990,
			Duration: 370,
		},
	}
}

// setupHelloTestHandler creates an Echo instance and HelloHandler for testing
func setupHelloTestHandler(tracks []service.HelloTrack) (*echo.Echo, *HelloHandler) {
	repo := &mockHelloRepo{tracks: tracks}
	svc := service.NewHelloService(repo)
	handler := NewHelloHandler(svc)
	e := echo.New()
	return e, handler
}

func TestHelloHealth_ReturnsOK(t *testing.T) {
	e, handler := setupHelloTestHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HelloHealth(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "hello", response["service"])
}

func TestHelloSearch_WithQuery(t *testing.T) {
	e, handler := setupHelloTestHandler(helloTestTracks())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search?q=neon", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HelloSearch(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response should have items array")
	assert.Len(t, items, 1)

	// Verify the matching track
	firstItem := items[0].(map[string]interface{})
	assert.Equal(t, "Neon Lights", firstItem["title"])
}

func TestHelloSearch_EmptyQuery(t *testing.T) {
	e, handler := setupHelloTestHandler(helloTestTracks())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/search", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HelloSearch(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response should have items array")
	assert.Empty(t, items)
}

func TestHelloFeatured_ReturnsAll(t *testing.T) {
	tracks := helloTestTracks()
	e, handler := setupHelloTestHandler(tracks)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HelloFeatured(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response should have items array")
	assert.Len(t, items, len(tracks))
}

func TestHelloFeatured_WithLimit(t *testing.T) {
	e, handler := setupHelloTestHandler(helloTestTracks())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/featured?limit=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HelloFeatured(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	items, ok := response["items"].([]interface{})
	require.True(t, ok, "response should have items array")
	assert.Len(t, items, 2)
}

func TestRegisterHelloRoutes_RoutesExist(t *testing.T) {
	tracks := helloTestTracks()
	repo := &mockHelloRepo{tracks: tracks}
	svc := service.NewHelloService(repo)
	handler := NewHelloHandler(svc)

	e := echo.New()
	RegisterHelloRoutes(e, handler)

	routes := e.Routes()

	// Collect all registered paths and methods
	routeMap := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = true
	}

	// Verify all 3 hello routes are registered
	assert.True(t, routeMap["GET /api/v1/hello/health"], "GET /api/v1/hello/health should be registered")
	assert.True(t, routeMap["GET /api/v1/hello/search"], "GET /api/v1/hello/search should be registered")
	assert.True(t, routeMap["GET /api/v1/hello/featured"], "GET /api/v1/hello/featured should be registered")
}
