package service

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHelloRepository provides mockable repository methods for hello service tests
type MockHelloRepository struct {
	mock.Mock
}

func (m *MockHelloRepository) ListTracksByUser(ctx context.Context, userID string) ([]models.Track, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

func sampleTracks() []models.Track {
	return []models.Track{
		{ID: "1", Title: "Midnight Drift", Artist: "Luna Waves", Genre: "Electronic", Duration: 240},
		{ID: "2", Title: "Solar Wind", Artist: "Luna Waves", Genre: "Electronic", Duration: 195},
		{ID: "3", Title: "Ghost Protocol", Artist: "DJ Phantom", Genre: "House", Duration: 330},
		{ID: "4", Title: "Silk Road", Artist: "Aria Chen", Genre: "Classical Crossover", Duration: 264},
		{ID: "5", Title: "Circuit Break", Artist: "Voltage", Genre: "Techno", Duration: 290},
	}
}

func TestHelloSearchTracks_ReturnsMatchingTracks(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	results, err := svc.SearchTracks(context.Background(), "luna", 20)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Midnight Drift + Solar Wind
	assert.Equal(t, "Midnight Drift", results[0].Title)
	assert.Equal(t, "Solar Wind", results[1].Title)
}

func TestHelloSearchTracks_CaseInsensitive(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	resultsLower, _ := svc.SearchTracks(context.Background(), "luna", 20)
	resultsUpper, _ := svc.SearchTracks(context.Background(), "LUNA", 20)
	assert.Equal(t, len(resultsLower), len(resultsUpper))
}

func TestHelloSearchTracks_MatchesArtist(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	results, err := svc.SearchTracks(context.Background(), "phantom", 20)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Ghost Protocol", results[0].Title)
}

func TestHelloSearchTracks_MatchesGenre(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	results, err := svc.SearchTracks(context.Background(), "electronic", 20)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestHelloSearchTracks_NoResults(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	results, err := svc.SearchTracks(context.Background(), "xyz123nonexistent", 20)
	assert.NoError(t, err)
	assert.Empty(t, results)
	assert.NotNil(t, results) // Should be empty slice, not nil
}

func TestHelloSearchTracks_EmptyQuery(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	results, err := svc.SearchTracks(context.Background(), "", 20)
	assert.NoError(t, err)
	assert.Empty(t, results)
	mockRepo.AssertNotCalled(t, "ListTracksByUser", mock.Anything, mock.Anything)
}

func TestHelloGetFeaturedTracks_ReturnsAllTracks(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	results, err := svc.GetFeaturedTracks(context.Background(), 0)
	assert.NoError(t, err)
	assert.Len(t, results, 5)
}

func TestHelloGetFeaturedTracks_RespectsLimit(t *testing.T) {
	mockRepo := new(MockHelloRepository)
	svc := NewHelloService(mockRepo)

	mockRepo.On("ListTracksByUser", mock.Anything, seedUserID).Return(sampleTracks(), nil)

	results, err := svc.GetFeaturedTracks(context.Background(), 3)
	assert.NoError(t, err)
	assert.Len(t, results, 3)
}
