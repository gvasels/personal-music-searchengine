package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockArtistRepository is a mock implementation of ArtistRepository
type MockArtistRepository struct {
	mock.Mock
}

func (m *MockArtistRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	args := m.Called(ctx, artist)
	return args.Error(0)
}

func (m *MockArtistRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	args := m.Called(ctx, userID, artistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Artist), args.Error(1)
}

func (m *MockArtistRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Artist), args.Error(1)
}

func (m *MockArtistRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Artist]), args.Error(1)
}

func (m *MockArtistRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	args := m.Called(ctx, artist)
	return args.Error(0)
}

func (m *MockArtistRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	args := m.Called(ctx, userID, artistID)
	return args.Error(0)
}

func (m *MockArtistRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	args := m.Called(ctx, userID, artistIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*models.Artist), args.Error(1)
}

func (m *MockArtistRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	args := m.Called(ctx, userID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Artist), args.Error(1)
}

func (m *MockArtistRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	args := m.Called(ctx, userID, artistID)
	return args.Int(0), args.Error(1)
}

func (m *MockArtistRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	args := m.Called(ctx, userID, artistID)
	return args.Int(0), args.Error(1)
}

func (m *MockArtistRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	args := m.Called(ctx, userID, artistID)
	return args.Int(0), args.Error(1)
}

func (m *MockArtistRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	args := m.Called(ctx, userID, artist)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

// ArtistMockS3Repository for artist service tests (separate from search_test mock)
type ArtistMockS3Repository struct {
	mock.Mock
}

func (m *ArtistMockS3Repository) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, expiry)
	return args.String(0), args.Error(1)
}

func (m *ArtistMockS3Repository) GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

func (m *ArtistMockS3Repository) GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, expiry time.Duration, filename string) (string, error) {
	args := m.Called(ctx, key, expiry, filename)
	return args.String(0), args.Error(1)
}

func (m *ArtistMockS3Repository) InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	args := m.Called(ctx, key, contentType)
	return args.String(0), args.Error(1)
}

func (m *ArtistMockS3Repository) GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error) {
	args := m.Called(ctx, key, uploadID, numParts, expiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.MultipartUploadPartURL), args.Error(1)
}

func (m *ArtistMockS3Repository) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error {
	args := m.Called(ctx, key, uploadID, parts)
	return args.Error(0)
}

func (m *ArtistMockS3Repository) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	args := m.Called(ctx, key, uploadID)
	return args.Error(0)
}

func (m *ArtistMockS3Repository) DeleteObject(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *ArtistMockS3Repository) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	args := m.Called(ctx, sourceKey, destKey)
	return args.Error(0)
}

func (m *ArtistMockS3Repository) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *ArtistMockS3Repository) GetObject(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *ArtistMockS3Repository) PutObject(ctx context.Context, key string, data []byte, contentType string) error {
	args := m.Called(ctx, key, data, contentType)
	return args.Error(0)
}

func (m *ArtistMockS3Repository) GetObjectSize(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ArtistMockS3Repository) ObjectExists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func TestArtistService_CreateArtist(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"

	t.Run("creates artist successfully", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		req := models.CreateArtistRequest{
			Name: "The Beatles",
			Bio:  "British rock band",
		}

		mockRepo.On("CreateArtist", ctx, mock.AnythingOfType("models.Artist")).Return(nil)

		result, err := service.CreateArtist(ctx, userID, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "The Beatles", result.Name)
		assert.Equal(t, "Beatles", result.SortName) // Generated sort name
		assert.Equal(t, "British rock band", result.Bio)
		assert.True(t, result.IsActive)
		assert.NotEmpty(t, result.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("uses provided sort name", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		req := models.CreateArtistRequest{
			Name:     "The Beatles",
			SortName: "Beatles, The",
		}

		mockRepo.On("CreateArtist", ctx, mock.AnythingOfType("models.Artist")).Return(nil)

		result, err := service.CreateArtist(ctx, userID, req)

		require.NoError(t, err)
		assert.Equal(t, "Beatles, The", result.SortName)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		req := models.CreateArtistRequest{
			Name: "Test Artist",
		}

		mockRepo.On("CreateArtist", ctx, mock.AnythingOfType("models.Artist")).Return(errors.New("database error"))

		result, err := service.CreateArtist(ctx, userID, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create artist")

		mockRepo.AssertExpectations(t)
	})
}

func TestArtistService_GetArtist(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	artistID := "artist-456"

	t.Run("returns artist with stats", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artist := &models.Artist{
			ID:       artistID,
			UserID:   userID,
			Name:     "Pink Floyd",
			IsActive: true,
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(artist, nil)
		mockRepo.On("GetArtistTrackCount", ctx, userID, artistID).Return(50, nil)
		mockRepo.On("GetArtistAlbumCount", ctx, userID, artistID).Return(15, nil)
		mockRepo.On("GetArtistTotalPlays", ctx, userID, artistID).Return(5000, nil)

		result, err := service.GetArtist(ctx, userID, artistID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Pink Floyd", result.Name)
		assert.Equal(t, 50, result.TrackCount)
		assert.Equal(t, 15, result.AlbumCount)
		assert.Equal(t, 5000, result.TotalPlays)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns not found error", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(nil, repository.ErrNotFound)

		result, err := service.GetArtist(ctx, userID, artistID)

		require.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		require.True(t, errors.As(err, &apiErr))
		assert.Equal(t, "NOT_FOUND", apiErr.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles stats fetch errors gracefully", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artist := &models.Artist{
			ID:     artistID,
			UserID: userID,
			Name:   "Test Artist",
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(artist, nil)
		mockRepo.On("GetArtistTrackCount", ctx, userID, artistID).Return(0, errors.New("stats error"))
		mockRepo.On("GetArtistAlbumCount", ctx, userID, artistID).Return(0, errors.New("stats error"))
		mockRepo.On("GetArtistTotalPlays", ctx, userID, artistID).Return(0, errors.New("stats error"))

		result, err := service.GetArtist(ctx, userID, artistID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.TrackCount)
		assert.Equal(t, 0, result.AlbumCount)
		assert.Equal(t, 0, result.TotalPlays)

		mockRepo.AssertExpectations(t)
	})
}

func TestArtistService_UpdateArtist(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	artistID := "artist-456"

	t.Run("updates artist successfully", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		existingArtist := &models.Artist{
			ID:       artistID,
			UserID:   userID,
			Name:     "Old Name",
			SortName: "Old Name",
		}

		newName := "New Name"
		req := models.UpdateArtistRequest{
			Name: &newName,
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(existingArtist, nil)
		mockRepo.On("UpdateArtist", ctx, mock.AnythingOfType("models.Artist")).Return(nil)

		result, err := service.UpdateArtist(ctx, userID, artistID, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New Name", result.SortName) // Auto-generated from name

		mockRepo.AssertExpectations(t)
	})

	t.Run("updates specific fields only", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		existingArtist := &models.Artist{
			ID:       artistID,
			UserID:   userID,
			Name:     "Artist Name",
			SortName: "Artist Name",
			Bio:      "Original bio",
		}

		newBio := "Updated bio"
		req := models.UpdateArtistRequest{
			Bio: &newBio,
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(existingArtist, nil)
		mockRepo.On("UpdateArtist", ctx, mock.AnythingOfType("models.Artist")).Return(nil)

		result, err := service.UpdateArtist(ctx, userID, artistID, req)

		require.NoError(t, err)
		assert.Equal(t, "Artist Name", result.Name) // Unchanged
		assert.Equal(t, "Updated bio", result.Bio)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns not found error", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		newName := "New Name"
		req := models.UpdateArtistRequest{
			Name: &newName,
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(nil, repository.ErrNotFound)

		result, err := service.UpdateArtist(ctx, userID, artistID, req)

		require.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		require.True(t, errors.As(err, &apiErr))
		assert.Equal(t, "NOT_FOUND", apiErr.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestArtistService_DeleteArtist(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	artistID := "artist-456"

	t.Run("deletes artist successfully", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artist := &models.Artist{
			ID:     artistID,
			UserID: userID,
			Name:   "Artist to delete",
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(artist, nil)
		mockRepo.On("DeleteArtist", ctx, userID, artistID).Return(nil)

		err := service.DeleteArtist(ctx, userID, artistID)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns not found error", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(nil, repository.ErrNotFound)

		err := service.DeleteArtist(ctx, userID, artistID)

		require.Error(t, err)

		var apiErr *models.APIError
		require.True(t, errors.As(err, &apiErr))
		assert.Equal(t, "NOT_FOUND", apiErr.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestArtistService_ListArtists(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"

	t.Run("returns paginated artists", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artists := []models.Artist{
			{ID: "artist-1", UserID: userID, Name: "Artist One"},
			{ID: "artist-2", UserID: userID, Name: "Artist Two"},
		}

		filter := models.ArtistFilter{
			Limit: 10,
		}

		mockRepo.On("ListArtists", ctx, userID, filter).Return(&repository.PaginatedResult[models.Artist]{
			Items:      artists,
			NextCursor: "cursor-abc",
			HasMore:    true,
		}, nil)

		result, err := service.ListArtists(ctx, userID, filter)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, "cursor-abc", result.NextCursor)
		assert.True(t, result.HasMore)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty list", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		filter := models.ArtistFilter{}

		mockRepo.On("ListArtists", ctx, userID, filter).Return(&repository.PaginatedResult[models.Artist]{
			Items:   []models.Artist{},
			HasMore: false,
		}, nil)

		result, err := service.ListArtists(ctx, userID, filter)

		require.NoError(t, err)
		assert.Len(t, result.Items, 0)
		assert.False(t, result.HasMore)

		mockRepo.AssertExpectations(t)
	})
}

func TestArtistService_SearchArtists(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"

	t.Run("searches artists successfully", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artists := []*models.Artist{
			{ID: "artist-1", Name: "The Beatles"},
			{ID: "artist-2", Name: "Beach Boys"},
		}

		mockRepo.On("SearchArtists", ctx, userID, "bea", 10).Return(artists, nil)

		result, err := service.SearchArtists(ctx, userID, "bea", 10)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "The Beatles", result[0].Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("uses default limit when zero", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		mockRepo.On("SearchArtists", ctx, userID, "test", 10).Return([]*models.Artist{}, nil)

		_, err := service.SearchArtists(ctx, userID, "test", 0)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestArtistService_GetArtistTracks(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	artistID := "artist-456"

	t.Run("returns artist tracks", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artist := &models.Artist{
			ID:     artistID,
			UserID: userID,
			Name:   "Test Artist",
		}

		tracks := []models.Track{
			{ID: "track-1", Title: "Song One", Artist: "Test Artist"},
			{ID: "track-2", Title: "Song Two", Artist: "Test Artist"},
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(artist, nil)
		mockRepo.On("ListTracksByArtist", ctx, userID, "Test Artist").Return(tracks, nil)

		result, err := service.GetArtistTracks(ctx, userID, artistID)

		require.NoError(t, err)
		assert.Len(t, result, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns not found for non-existent artist", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(nil, repository.ErrNotFound)

		result, err := service.GetArtistTracks(ctx, userID, artistID)

		require.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		require.True(t, errors.As(err, &apiErr))
		assert.Equal(t, "NOT_FOUND", apiErr.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("generates cover art URLs", func(t *testing.T) {
		mockRepo := new(MockArtistRepository)
		mockS3 := new(ArtistMockS3Repository)
		service := NewArtistService(mockRepo, mockS3)

		artist := &models.Artist{
			ID:     artistID,
			UserID: userID,
			Name:   "Test Artist",
		}

		tracks := []models.Track{
			{ID: "track-1", Title: "Song One", CoverArtKey: "covers/track-1.jpg"},
		}

		mockRepo.On("GetArtist", ctx, userID, artistID).Return(artist, nil)
		mockRepo.On("ListTracksByArtist", ctx, userID, "Test Artist").Return(tracks, nil)
		mockS3.On("GeneratePresignedDownloadURL", ctx, "covers/track-1.jpg", 24*time.Hour).Return("https://example.com/cover.jpg", nil)

		result, err := service.GetArtistTracks(ctx, userID, artistID)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "https://example.com/cover.jpg", result[0].CoverArtURL)

		mockRepo.AssertExpectations(t)
		mockS3.AssertExpectations(t)
	})
}

func TestGenerateDeterministicArtistID(t *testing.T) {
	t.Run("generates consistent UUID for same inputs", func(t *testing.T) {
		userID := "user-123"
		artistName := "The Beatles"

		id1 := GenerateDeterministicArtistID(userID, artistName)
		id2 := GenerateDeterministicArtistID(userID, artistName)

		assert.Equal(t, id1, id2)
		assert.Len(t, id1, 36) // UUID format length
	})

	t.Run("generates different UUIDs for different users", func(t *testing.T) {
		artistName := "The Beatles"

		id1 := GenerateDeterministicArtistID("user-1", artistName)
		id2 := GenerateDeterministicArtistID("user-2", artistName)

		assert.NotEqual(t, id1, id2)
	})

	t.Run("generates different UUIDs for different artist names", func(t *testing.T) {
		userID := "user-123"

		id1 := GenerateDeterministicArtistID(userID, "The Beatles")
		id2 := GenerateDeterministicArtistID(userID, "Pink Floyd")

		assert.NotEqual(t, id1, id2)
	})

	t.Run("generates valid UUID format", func(t *testing.T) {
		id := GenerateDeterministicArtistID("user-123", "Test Artist")

		// Check UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, id)
	})

	t.Run("sets UUID version 5 bits", func(t *testing.T) {
		id := GenerateDeterministicArtistID("user-123", "Test Artist")

		// Version bits are at position 14-15 (0-indexed char 14)
		// For version 5, the first nibble of the 7th byte should be 5
		assert.Equal(t, '5', rune(id[14]))
	})
}
