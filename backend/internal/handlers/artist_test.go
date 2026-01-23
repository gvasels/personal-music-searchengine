package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestValidator for Echo
type TestValidator struct {
	validator *validator.Validate
}

func (tv *TestValidator) Validate(i interface{}) error {
	return tv.validator.Struct(i)
}

// MockArtistService implements service.ArtistService for testing
type MockArtistService struct {
	mock.Mock
}

func (m *MockArtistService) CreateArtist(ctx context.Context, userID string, req models.CreateArtistRequest) (*models.ArtistResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistResponse), args.Error(1)
}

func (m *MockArtistService) GetArtist(ctx context.Context, userID, artistID string) (*models.ArtistWithStatsResponse, error) {
	args := m.Called(ctx, userID, artistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistWithStatsResponse), args.Error(1)
}

func (m *MockArtistService) UpdateArtist(ctx context.Context, userID, artistID string, req models.UpdateArtistRequest) (*models.ArtistResponse, error) {
	args := m.Called(ctx, userID, artistID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtistResponse), args.Error(1)
}

func (m *MockArtistService) DeleteArtist(ctx context.Context, userID, artistID string) error {
	args := m.Called(ctx, userID, artistID)
	return args.Error(0)
}

func (m *MockArtistService) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.ArtistResponse], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.ArtistResponse]), args.Error(1)
}

func (m *MockArtistService) SearchArtists(ctx context.Context, userID, query string, limit int) ([]models.ArtistResponse, error) {
	args := m.Called(ctx, userID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ArtistResponse), args.Error(1)
}

func (m *MockArtistService) GetArtistTracks(ctx context.Context, userID, artistID string) ([]models.TrackResponse, error) {
	args := m.Called(ctx, userID, artistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TrackResponse), args.Error(1)
}

func setupArtistTestHandler(mockArtist *MockArtistService) (*echo.Echo, *Handlers) {
	e := echo.New()
	e.Validator = &TestValidator{validator: validator.New()}
	services := &service.Services{
		Artist: mockArtist,
	}
	h := NewHandlers(services)
	return e, h
}

func TestCreateArtist(t *testing.T) {
	t.Run("creates artist successfully", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		now := time.Now()
		expectedResponse := &models.ArtistResponse{
			ID:        "artist-123",
			Name:      "The Beatles",
			SortName:  "Beatles",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockArtist.On("CreateArtist", mock.Anything, "user-123", mock.AnythingOfType("models.CreateArtistRequest")).Return(expectedResponse, nil)

		reqBody := `{"name": "The Beatles"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/entity", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "artist-123", response["id"])
		assert.Equal(t, "The Beatles", response["name"])

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 401 without user ID", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		reqBody := `{"name": "The Beatles"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/entity", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// No X-User-ID header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("returns 400 for invalid request body", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/entity", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 for service error", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("CreateArtist", mock.Anything, "user-123", mock.AnythingOfType("models.CreateArtistRequest")).Return(nil, errors.New("database error"))

		reqBody := `{"name": "Test Artist"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/artists/entity", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		mockArtist.AssertExpectations(t)
	})
}

func TestGetArtistByID(t *testing.T) {
	t.Run("gets artist successfully", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		now := time.Now()
		expectedResponse := &models.ArtistWithStatsResponse{
			ArtistResponse: models.ArtistResponse{
				ID:        "artist-123",
				Name:      "Pink Floyd",
				SortName:  "Pink Floyd",
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			},
			TrackCount: 50,
			AlbumCount: 15,
			TotalPlays: 10000,
		}

		mockArtist.On("GetArtist", mock.Anything, "user-123", "artist-123").Return(expectedResponse, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/artist-123", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("artist-123")

		err := h.GetArtistByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "artist-123", response["id"])
		assert.Equal(t, float64(50), response["trackCount"])

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent artist", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("GetArtist", mock.Anything, "user-123", "nonexistent").Return(nil, models.NewNotFoundError("artist", "nonexistent"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/nonexistent", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent")

		err := h.GetArtistByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 400 for empty ID", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		err := h.GetArtistByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestUpdateArtist(t *testing.T) {
	t.Run("updates artist successfully", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		now := time.Now()
		expectedResponse := &models.ArtistResponse{
			ID:        "artist-123",
			Name:      "Updated Name",
			SortName:  "Updated Name",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockArtist.On("UpdateArtist", mock.Anything, "user-123", "artist-123", mock.AnythingOfType("models.UpdateArtistRequest")).Return(expectedResponse, nil)

		reqBody := `{"name": "Updated Name"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/artists/entity/artist-123", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("artist-123")

		err := h.UpdateArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", response["name"])

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent artist", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("UpdateArtist", mock.Anything, "user-123", "nonexistent", mock.AnythingOfType("models.UpdateArtistRequest")).Return(nil, models.NewNotFoundError("artist", "nonexistent"))

		reqBody := `{"name": "Updated Name"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/artists/entity/nonexistent", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent")

		err := h.UpdateArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockArtist.AssertExpectations(t)
	})
}

func TestDeleteArtist(t *testing.T) {
	t.Run("deletes artist successfully", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("DeleteArtist", mock.Anything, "user-123", "artist-123").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/artists/entity/artist-123", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("artist-123")

		err := h.DeleteArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent artist", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("DeleteArtist", mock.Anything, "user-123", "nonexistent").Return(models.NewNotFoundError("artist", "nonexistent"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/artists/entity/nonexistent", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent")

		err := h.DeleteArtist(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockArtist.AssertExpectations(t)
	})
}

func TestListArtistsEntity(t *testing.T) {
	t.Run("lists artists with pagination", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		expectedResult := &repository.PaginatedResult[models.ArtistResponse]{
			Items: []models.ArtistResponse{
				{ID: "artist-1", Name: "Artist One"},
				{ID: "artist-2", Name: "Artist Two"},
			},
			NextCursor: "cursor-abc",
			HasMore:    true,
		}

		mockArtist.On("ListArtists", mock.Anything, "user-123", mock.AnythingOfType("models.ArtistFilter")).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity?limit=10", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ListArtistsEntity(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		items := response["items"].([]interface{})
		assert.Len(t, items, 2)
		assert.Equal(t, "cursor-abc", response["nextCursor"])
		assert.True(t, response["hasMore"].(bool))

		mockArtist.AssertExpectations(t)
	})

	t.Run("applies filter parameters", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("ListArtists", mock.Anything, "user-123", mock.MatchedBy(func(filter models.ArtistFilter) bool {
			return filter.Name == "Beatles" && filter.SortBy == "name" && filter.SortOrder == "asc" && filter.Limit == 5
		})).Return(&repository.PaginatedResult[models.ArtistResponse]{
			Items:   []models.ArtistResponse{},
			HasMore: false,
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity?name=Beatles&sortBy=name&sortOrder=asc&limit=5", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ListArtistsEntity(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		mockArtist.AssertExpectations(t)
	})
}

func TestSearchArtists(t *testing.T) {
	t.Run("searches artists successfully", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		expectedResults := []models.ArtistResponse{
			{ID: "artist-1", Name: "The Beatles"},
			{ID: "artist-2", Name: "Beach Boys"},
		}

		mockArtist.On("SearchArtists", mock.Anything, "user-123", "bea", 10).Return(expectedResults, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/search?q=bea", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SearchArtists(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		items := response["items"].([]interface{})
		assert.Len(t, items, 2)

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 400 for empty query", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/search?q=", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SearchArtists(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("uses custom limit", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("SearchArtists", mock.Anything, "user-123", "test", 5).Return([]models.ArtistResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/search?q=test&limit=5", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SearchArtists(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		mockArtist.AssertExpectations(t)
	})
}

func TestGetArtistTracksEntity(t *testing.T) {
	t.Run("gets artist tracks successfully", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		expectedTracks := []models.TrackResponse{
			{ID: "track-1", Title: "Song One"},
			{ID: "track-2", Title: "Song Two"},
		}

		mockArtist.On("GetArtistTracks", mock.Anything, "user-123", "artist-123").Return(expectedTracks, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/artist-123/tracks", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("artist-123")

		err := h.GetArtistTracksEntity(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		items := response["items"].([]interface{})
		assert.Len(t, items, 2)

		mockArtist.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent artist", func(t *testing.T) {
		mockArtist := new(MockArtistService)
		e, h := setupArtistTestHandler(mockArtist)

		mockArtist.On("GetArtistTracks", mock.Anything, "user-123", "nonexistent").Return(nil, models.NewNotFoundError("artist", "nonexistent"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/artists/entity/nonexistent/tracks", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent")

		err := h.GetArtistTracksEntity(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockArtist.AssertExpectations(t)
	})
}
