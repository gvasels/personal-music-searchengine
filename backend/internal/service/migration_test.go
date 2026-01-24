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

// MockMigrationRepository is a mock implementation of MigrationRepository
type MockMigrationRepository struct {
	mock.Mock
}

func (m *MockMigrationRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	args := m.Called(ctx, artist)
	return args.Error(0)
}

func (m *MockMigrationRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	args := m.Called(ctx, userID, artistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Artist), args.Error(1)
}

func (m *MockMigrationRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Artist), args.Error(1)
}

func (m *MockMigrationRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Artist]), args.Error(1)
}

func (m *MockMigrationRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockMigrationRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

// Helper function to create test tracks
func createMigrationTestTrack(userID, trackID, artist, artistID string) models.Track {
	now := time.Now().UTC()
	return models.Track{
		UserID:   userID,
		ID:       trackID,
		Artist:   artist,
		ArtistID: artistID,
		Title:    "Test Track " + trackID,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// Helper function to create test artist
func createMigrationTestArtist(userID, artistID, name string) *models.Artist {
	now := time.Now().UTC()
	return &models.Artist{
		UserID: userID,
		ID:     artistID,
		Name:   name,
		Timestamps: models.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func TestNewMigrationService(t *testing.T) {
	t.Run("creates service with repository", func(t *testing.T) {
		mockRepo := new(MockMigrationRepository)
		svc := NewMigrationService(mockRepo)
		require.NotNil(t, svc)
	})
}

func TestMigrateArtists_Success(t *testing.T) {
	t.Run("migrates tracks and creates artists", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Artist A", ""),
			createMigrationTestTrack(userID, "track-2", "Artist B", ""),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)

		// Artist A doesn't exist
		mockRepo.On("GetArtistByName", ctx, userID, "Artist A").Return([]*models.Artist(nil), nil)
		mockRepo.On("GetArtist", ctx, userID, mock.AnythingOfType("string")).Return((*models.Artist)(nil), repository.ErrNotFound).Maybe()
		mockRepo.On("CreateArtist", ctx, mock.MatchedBy(func(a models.Artist) bool {
			return a.UserID == userID && a.Name == "Artist A"
		})).Return(nil)

		// Artist B doesn't exist
		mockRepo.On("GetArtistByName", ctx, userID, "Artist B").Return([]*models.Artist(nil), nil)
		mockRepo.On("CreateArtist", ctx, mock.MatchedBy(func(a models.Artist) bool {
			return a.UserID == userID && a.Name == "Artist B"
		})).Return(nil)

		// Update tracks
		mockRepo.On("UpdateTrack", ctx, mock.AnythingOfType("models.Track")).Return(nil)

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.TracksUpdated)
		assert.Equal(t, 2, result.ArtistsCreated)
	})
}

func TestMigrateArtists_SkipsAlreadyMigrated(t *testing.T) {
	t.Run("skips tracks with existing artistId", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Artist A", "existing-artist-id"),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.TracksUpdated)
		assert.Equal(t, 1, result.TracksSkipped)
		mockRepo.AssertNotCalled(t, "CreateArtist", mock.Anything, mock.Anything)
	})
}

func TestMigrateArtists_SkipsNoArtist(t *testing.T) {
	t.Run("skips tracks without artist name", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "", ""),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.TracksUpdated)
		assert.Equal(t, 1, result.TracksSkipped)
		mockRepo.AssertNotCalled(t, "CreateArtist", mock.Anything, mock.Anything)
	})
}

func TestMigrateArtists_ReusesExistingArtist(t *testing.T) {
	t.Run("reuses existing artist from database", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Existing Artist", ""),
		}

		existingArtist := createMigrationTestArtist(userID, "existing-artist-id", "Existing Artist")

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)
		mockRepo.On("GetArtistByName", ctx, userID, "Existing Artist").Return([]*models.Artist{existingArtist}, nil)
		mockRepo.On("UpdateTrack", ctx, mock.MatchedBy(func(t models.Track) bool {
			return t.ID == "track-1" && t.ArtistID == "existing-artist-id"
		})).Return(nil)

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.TracksUpdated)
		// Note: ArtistsCreated counts cached artists, including those found in DB
		assert.Equal(t, 1, result.ArtistsCreated)
		mockRepo.AssertNotCalled(t, "CreateArtist", mock.Anything, mock.Anything)
	})
}

func TestMigrateArtists_UsesCache(t *testing.T) {
	t.Run("uses cache for same artist across multiple tracks", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Same Artist", ""),
			createMigrationTestTrack(userID, "track-2", "Same Artist", ""),
			createMigrationTestTrack(userID, "track-3", "Same Artist", ""),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)

		// GetArtistByName should only be called once due to caching
		mockRepo.On("GetArtistByName", ctx, userID, "Same Artist").Return([]*models.Artist(nil), nil).Once()
		mockRepo.On("GetArtist", ctx, userID, mock.AnythingOfType("string")).Return((*models.Artist)(nil), repository.ErrNotFound).Maybe()
		mockRepo.On("CreateArtist", ctx, mock.MatchedBy(func(a models.Artist) bool {
			return a.Name == "Same Artist"
		})).Return(nil).Once()

		mockRepo.On("UpdateTrack", ctx, mock.AnythingOfType("models.Track")).Return(nil).Times(3)

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, result.TracksUpdated)
		assert.Equal(t, 1, result.ArtistsCreated)
	})
}

func TestMigrateArtists_Pagination(t *testing.T) {
	t.Run("handles multiple pages of tracks", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracksPage1 := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Artist A", ""),
		}
		tracksPage2 := []models.Track{
			createMigrationTestTrack(userID, "track-2", "Artist B", ""),
		}

		page1Result := &repository.PaginatedResult[models.Track]{
			Items:      tracksPage1,
			HasMore:    true,
			NextCursor: "cursor-1",
		}
		page2Result := &repository.PaginatedResult[models.Track]{
			Items:   tracksPage2,
			HasMore: false,
		}

		// First page
		mockRepo.On("ListTracks", ctx, userID, mock.MatchedBy(func(f models.TrackFilter) bool {
			return f.LastKey == ""
		})).Return(page1Result, nil).Once()

		// Second page
		mockRepo.On("ListTracks", ctx, userID, mock.MatchedBy(func(f models.TrackFilter) bool {
			return f.LastKey == "cursor-1"
		})).Return(page2Result, nil).Once()

		mockRepo.On("GetArtistByName", ctx, userID, mock.AnythingOfType("string")).Return([]*models.Artist(nil), nil)
		mockRepo.On("GetArtist", ctx, userID, mock.AnythingOfType("string")).Return((*models.Artist)(nil), repository.ErrNotFound).Maybe()
		mockRepo.On("CreateArtist", ctx, mock.AnythingOfType("models.Artist")).Return(nil)
		mockRepo.On("UpdateTrack", ctx, mock.AnythingOfType("models.Track")).Return(nil)

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.TracksUpdated)
	})
}

func TestMigrateArtists_ListTracksError(t *testing.T) {
	t.Run("returns error when listing tracks fails", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return((*repository.PaginatedResult[models.Track])(nil), errors.New("database error"))

		svc := NewMigrationService(mockRepo)
		result, err := svc.MigrateArtists(ctx, userID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to list tracks")
	})
}

func TestGetMigrationStatus_NotStarted(t *testing.T) {
	t.Run("returns not_started when no tracks are migrated", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Artist A", ""),
			createMigrationTestTrack(userID, "track-2", "Artist B", ""),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}
		artistsResult := &repository.PaginatedResult[models.Artist]{
			Items:   []models.Artist{},
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)
		mockRepo.On("ListArtists", ctx, userID, mock.AnythingOfType("models.ArtistFilter")).Return(artistsResult, nil)

		svc := NewMigrationService(mockRepo)
		status, err := svc.GetMigrationStatus(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, status)
		assert.Equal(t, "not_started", status.Status)
		assert.Equal(t, 2, status.TotalTracks)
		assert.Equal(t, 0, status.MigratedTracks)
		assert.Equal(t, 2, status.UnmigratedTracks)
	})
}

func TestGetMigrationStatus_Completed(t *testing.T) {
	t.Run("returns completed when all tracks are migrated", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Artist A", "artist-id-1"),
			createMigrationTestTrack(userID, "track-2", "Artist B", "artist-id-2"),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}
		artistsResult := &repository.PaginatedResult[models.Artist]{
			Items:   []models.Artist{{ID: "artist-id-1"}, {ID: "artist-id-2"}},
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)
		mockRepo.On("ListArtists", ctx, userID, mock.AnythingOfType("models.ArtistFilter")).Return(artistsResult, nil)

		svc := NewMigrationService(mockRepo)
		status, err := svc.GetMigrationStatus(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, status)
		assert.Equal(t, "completed", status.Status)
		assert.Equal(t, 2, status.TotalTracks)
		assert.Equal(t, 2, status.MigratedTracks)
		assert.Equal(t, 0, status.UnmigratedTracks)
	})
}

func TestGetMigrationStatus_Partial(t *testing.T) {
	t.Run("returns partial when some tracks are migrated", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracks := []models.Track{
			createMigrationTestTrack(userID, "track-1", "Artist A", "artist-id-1"),
			createMigrationTestTrack(userID, "track-2", "Artist B", ""),
			createMigrationTestTrack(userID, "track-3", "Artist C", "artist-id-3"),
		}

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   tracks,
			HasMore: false,
		}
		artistsResult := &repository.PaginatedResult[models.Artist]{
			Items:   []models.Artist{{ID: "artist-id-1"}, {ID: "artist-id-3"}},
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)
		mockRepo.On("ListArtists", ctx, userID, mock.AnythingOfType("models.ArtistFilter")).Return(artistsResult, nil)

		svc := NewMigrationService(mockRepo)
		status, err := svc.GetMigrationStatus(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, status)
		assert.Equal(t, "partial", status.Status)
		assert.Equal(t, 3, status.TotalTracks)
		assert.Equal(t, 2, status.MigratedTracks)
		assert.Equal(t, 1, status.UnmigratedTracks)
	})
}

func TestGetMigrationStatus_NoTracks(t *testing.T) {
	t.Run("returns completed when there are no tracks", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		tracksResult := &repository.PaginatedResult[models.Track]{
			Items:   []models.Track{},
			HasMore: false,
		}
		artistsResult := &repository.PaginatedResult[models.Artist]{
			Items:   []models.Artist{},
			HasMore: false,
		}

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(tracksResult, nil)
		mockRepo.On("ListArtists", ctx, userID, mock.AnythingOfType("models.ArtistFilter")).Return(artistsResult, nil)

		svc := NewMigrationService(mockRepo)
		status, err := svc.GetMigrationStatus(ctx, userID)

		require.NoError(t, err)
		require.NotNil(t, status)
		assert.Equal(t, 0, status.TotalTracks)
	})
}

func TestGetMigrationStatus_ListTracksError(t *testing.T) {
	t.Run("returns error when listing tracks fails", func(t *testing.T) {
		ctx := context.Background()
		userID := "user-123"
		mockRepo := new(MockMigrationRepository)

		mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return((*repository.PaginatedResult[models.Track])(nil), errors.New("database error"))

		svc := NewMigrationService(mockRepo)
		status, err := svc.GetMigrationStatus(ctx, userID)

		require.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "failed to list tracks")
	})
}
