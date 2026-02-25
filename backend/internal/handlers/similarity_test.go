package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVectorService struct {
	mock.Mock
}

func (m *MockVectorService) PutVector(ctx context.Context, id string, vector []float32, metadata map[string]string) error {
	args := m.Called(ctx, id, vector, metadata)
	return args.Error(0)
}

func (m *MockVectorService) GetVector(ctx context.Context, id string) ([]float32, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), args.Error(1)
}

func (m *MockVectorService) QuerySimilar(ctx context.Context, vector []float32, k int) ([]service.VectorResult, error) {
	args := m.Called(ctx, vector, k)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.VectorResult), args.Error(1)
}

func (m *MockVectorService) DeleteVector(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSimilarityTrackService struct {
	mock.Mock
}

func (m *MockSimilarityTrackService) GetTrack(ctx context.Context, requesterID, trackID string, hasGlobal bool) (*models.TrackResponse, error) {
	args := m.Called(ctx, requesterID, trackID, hasGlobal)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TrackResponse), args.Error(1)
}
func (m *MockSimilarityTrackService) UpdateTrack(ctx context.Context, userID, trackID string, req models.UpdateTrackRequest) (*models.TrackResponse, error) {
	return nil, nil
}
func (m *MockSimilarityTrackService) DeleteTrack(ctx context.Context, userID, trackID string, hasGlobal bool) error {
	return nil
}
func (m *MockSimilarityTrackService) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.TrackResponse], error) {
	return nil, nil
}
func (m *MockSimilarityTrackService) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.TrackResponse, error) {
	return nil, nil
}
func (m *MockSimilarityTrackService) IncrementPlayCount(ctx context.Context, userID, trackID string) error {
	return nil
}
func (m *MockSimilarityTrackService) UpdateVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}
func (m *MockSimilarityTrackService) GetLibraryStats(ctx context.Context, userID string, scope service.StatsScope, hasGlobal bool) (*service.LibraryStats, error) {
	return nil, nil
}

func TestGetSimilarTracks(t *testing.T) {
	e := echo.New()

	t.Run("returns similar tracks with key compatibility boost", func(t *testing.T) {
		mockVector := &MockVectorService{}
		mockTrack := &MockSimilarityTrackService{}

		embedding := make([]float32, 1024)
		mockVector.On("GetVector", mock.Anything, "track-1").Return(embedding, nil)
		mockVector.On("QuerySimilar", mock.Anything, embedding, 40).Return([]service.VectorResult{
			{ID: "track-2", Score: 0.95},
			{ID: "track-3", Score: 0.85},
		}, nil)
		// Source track has key 8A
		mockTrack.On("GetTrack", mock.Anything, mock.Anything, "track-1", false).Return(&models.TrackResponse{ID: "track-1", Title: "Source", KeyCamelot: "8A"}, nil)
		// track-2: compatible key (7A is neighbor of 8A)
		mockTrack.On("GetTrack", mock.Anything, mock.Anything, "track-2", false).Return(&models.TrackResponse{ID: "track-2", Title: "Similar 1", KeyCamelot: "7A"}, nil)
		// track-3: incompatible key
		mockTrack.On("GetTrack", mock.Anything, mock.Anything, "track-3", false).Return(&models.TrackResponse{ID: "track-3", Title: "Similar 2", KeyCamelot: "5B"}, nil)

		h := &Handlers{services: &service.Services{Vector: mockVector, Track: mockTrack}}

		req := httptest.NewRequest(http.MethodGet, "/tracks/track-1/similar", nil)
		req.Header.Set("X-User-ID", "user-1")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("track-1")

		err := h.GetSimilarTracks(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		// track-2 should be boosted (compatible key), track-3 not
		assert.Contains(t, rec.Body.String(), `"keyCompatible":true`)
		assert.Contains(t, rec.Body.String(), `"keyCompatible":false`)
	})

	t.Run("returns empty array when track has no embedding", func(t *testing.T) {
		mockVector := &MockVectorService{}
		mockVector.On("GetVector", mock.Anything, "track-no-embed").Return(nil, nil)

		h := &Handlers{services: &service.Services{Vector: mockVector}}

		req := httptest.NewRequest(http.MethodGet, "/tracks/track-no-embed/similar", nil)
		req.Header.Set("X-User-ID", "user-1")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("track-no-embed")

		err := h.GetSimilarTracks(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"similar":[]`)
	})

	t.Run("returns empty array when vector service is nil", func(t *testing.T) {
		h := &Handlers{services: &service.Services{Vector: nil}}

		req := httptest.NewRequest(http.MethodGet, "/tracks/track-1/similar", nil)
		req.Header.Set("X-User-ID", "user-1")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("track-1")

		err := h.GetSimilarTracks(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"similar":[]`)
	})
}
