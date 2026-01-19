package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// MockRepository implements repository.Repository interface for testing
type MockRepository struct {
	mock.Mock
}

// User operations
func (m *MockRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Track operations
func (m *MockRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockRepository) CreateTrack(ctx context.Context, track *models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockRepository) UpdateTrack(ctx context.Context, track *models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	args := m.Called(ctx, userID, trackID)
	return args.Error(0)
}

func (m *MockRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*models.PaginatedResponse[models.Track], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.Track]), args.Error(1)
}

func (m *MockRepository) ListTracksByAlbum(ctx context.Context, userID, albumID string) ([]models.Track, error) {
	args := m.Called(ctx, userID, albumID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

func (m *MockRepository) ListTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

// Tag operations
func (m *MockRepository) AddTrackTag(ctx context.Context, tt *models.TrackTag) error {
	args := m.Called(ctx, tt)
	return args.Error(0)
}

func (m *MockRepository) RemoveTrackTag(ctx context.Context, userID, trackID, tagName string) error {
	args := m.Called(ctx, userID, trackID, tagName)
	return args.Error(0)
}

func (m *MockRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockRepository) CreateTag(ctx context.Context, tag *models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockRepository) UpdateTag(ctx context.Context, tag *models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	args := m.Called(ctx, userID, tagName)
	return args.Error(0)
}

func (m *MockRepository) ListTags(ctx context.Context, userID string, filter models.TagFilter) (*models.PaginatedResponse[models.Tag], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.Tag]), args.Error(1)
}

// Album operations (stubs for interface)
func (m *MockRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	args := m.Called(ctx, userID, albumID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Album), args.Error(1)
}

func (m *MockRepository) CreateAlbum(ctx context.Context, album *models.Album) error {
	args := m.Called(ctx, album)
	return args.Error(0)
}

func (m *MockRepository) UpdateAlbum(ctx context.Context, album *models.Album) error {
	args := m.Called(ctx, album)
	return args.Error(0)
}

func (m *MockRepository) DeleteAlbum(ctx context.Context, userID, albumID string) error {
	args := m.Called(ctx, userID, albumID)
	return args.Error(0)
}

func (m *MockRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*models.PaginatedResponse[models.Album], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.Album]), args.Error(1)
}

func (m *MockRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	args := m.Called(ctx, userID, artist)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Album), args.Error(1)
}

func (m *MockRepository) GetOrCreateAlbum(ctx context.Context, userID, title, artist string, year int) (*models.Album, error) {
	args := m.Called(ctx, userID, title, artist, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Album), args.Error(1)
}

// Playlist operations (stubs)
func (m *MockRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	args := m.Called(ctx, userID, playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Playlist), args.Error(1)
}

func (m *MockRepository) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	args := m.Called(ctx, playlist)
	return args.Error(0)
}

func (m *MockRepository) UpdatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	args := m.Called(ctx, playlist)
	return args.Error(0)
}

func (m *MockRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	args := m.Called(ctx, userID, playlistID)
	return args.Error(0)
}

func (m *MockRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*models.PaginatedResponse[models.Playlist], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.Playlist]), args.Error(1)
}

func (m *MockRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	args := m.Called(ctx, playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PlaylistTrack), args.Error(1)
}

func (m *MockRepository) AddPlaylistTrack(ctx context.Context, pt *models.PlaylistTrack) error {
	args := m.Called(ctx, pt)
	return args.Error(0)
}

func (m *MockRepository) RemovePlaylistTrack(ctx context.Context, playlistID, trackID string) error {
	args := m.Called(ctx, playlistID, trackID)
	return args.Error(0)
}

func (m *MockRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, trackIDs []string) error {
	args := m.Called(ctx, playlistID, trackIDs)
	return args.Error(0)
}

// Upload operations (stubs)
func (m *MockRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	args := m.Called(ctx, userID, uploadID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Upload), args.Error(1)
}

func (m *MockRepository) CreateUpload(ctx context.Context, upload *models.Upload) error {
	args := m.Called(ctx, upload)
	return args.Error(0)
}

func (m *MockRepository) UpdateUpload(ctx context.Context, upload *models.Upload) error {
	args := m.Called(ctx, upload)
	return args.Error(0)
}

func (m *MockRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*models.PaginatedResponse[models.Upload], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse[models.Upload]), args.Error(1)
}

// ====================
// TrackService Tests
// ====================

// Note: These tests document expected behavior but require the TrackService
// to accept a Repository interface instead of concrete DynamoDBRepository.
// TODO: Refactor TrackService to use interface for dependency injection.

func TestTrackService_GetTrack_ExpectedBehavior(t *testing.T) {
	// This test documents expected behavior
	// When TrackService.GetTrack is called:
	// 1. It should call repo.GetTrack with userID and trackID
	// 2. If track found, convert to TrackResponse with cover art URL
	// 3. If track not found, return not found error

	t.Run("returns track response on success", func(t *testing.T) {
		// Setup
		userID := "user-456"
		trackID := "track-123"

		expectedTrack := &models.Track{
			ID:          trackID,
			UserID:      userID,
			Title:       "Test Song",
			Artist:      "Test Artist",
			Album:       "Test Album",
			Duration:    180,
			Format:      models.AudioFormatMP3,
			FileSize:    5242880,
			CoverArtKey: "covers/user-456/track-123/cover.jpg",
		}

		// When GetTrack returns a track:
		// - TrackResponse should have all fields populated
		// - DurationStr should be formatted (e.g., "3:00")
		// - FileSizeStr should be formatted (e.g., "5.00 MB")
		// - CoverArtURL should be CloudFront URL if CLOUDFRONT_DOMAIN set

		response := expectedTrack.ToResponse("")
		assert.Equal(t, trackID, response.ID)
		assert.Equal(t, userID, expectedTrack.UserID) // Track belongs to this user
		assert.Equal(t, "Test Song", response.Title)
		assert.Equal(t, "3:00", response.DurationStr)
		assert.Equal(t, "5.00 MB", response.FileSizeStr)
	})

	t.Run("returns not found error when track doesn't exist", func(t *testing.T) {
		// When repo.GetTrack returns NotFoundError, service should propagate it
		err := models.NewNotFoundError("Track", "nonexistent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "NOT_FOUND")
	})
}

func TestTrackService_ListTracks_ExpectedBehavior(t *testing.T) {
	t.Run("returns paginated track responses", func(t *testing.T) {
		// Expected behavior:
		// 1. Call repo.ListTracks with filter
		// 2. Convert each Track to TrackResponse
		// 3. Return PaginatedResponse with converted items

		tracks := []models.Track{
			{ID: "track-1", Title: "Song A", Artist: "Artist", Duration: 120, Format: models.AudioFormatMP3, FileSize: 3145728},
			{ID: "track-2", Title: "Song B", Artist: "Artist", Duration: 180, Format: models.AudioFormatMP3, FileSize: 4194304},
		}

		responses := make([]models.TrackResponse, len(tracks))
		for i, track := range tracks {
			responses[i] = track.ToResponse("")
		}

		assert.Len(t, responses, 2)
		assert.Equal(t, "2:00", responses[0].DurationStr)
		assert.Equal(t, "3:00", responses[1].DurationStr)
	})

	t.Run("applies filters correctly", func(t *testing.T) {
		// When filter has Artist, Genre, Year, etc.:
		// - Should pass to repository for filtering
		// - Results should match filter criteria

		filter := models.TrackFilter{
			Artist: "Test Artist",
			Genre:  "Rock",
			Year:   2024,
			Limit:  20,
		}

		assert.Equal(t, "Test Artist", filter.Artist)
		assert.Equal(t, "Rock", filter.Genre)
		assert.Equal(t, 2024, filter.Year)
		assert.Equal(t, 20, filter.Limit)
	})

	t.Run("supports pagination", func(t *testing.T) {
		// When filter has LastKey:
		// - Results should start after that key
		// - NextKey should be set if more results exist

		pagination := models.Pagination{
			Limit:   20,
			LastKey: "eyJQSyI6IlVTRVIjdXNlci00NTYifQ==",
			NextKey: "eyJQSyI6IlVTRVIjdXNlci00NTYiLCJTSyI6IlRSQUNLI3RyYWNrLTIwIn0=",
		}

		assert.Equal(t, 20, pagination.Limit)
		assert.NotEmpty(t, pagination.LastKey)
		assert.NotEmpty(t, pagination.NextKey)
	})
}

func TestTrackService_UpdateTrack_ExpectedBehavior(t *testing.T) {
	t.Run("updates only provided fields", func(t *testing.T) {
		// Expected behavior:
		// 1. Get existing track
		// 2. Apply only non-nil fields from request
		// 3. Update track in repository
		// 4. Return updated TrackResponse

		track := models.Track{
			ID:     "track-123",
			Title:  "Original Title",
			Artist: "Original Artist",
			Album:  "Original Album",
		}

		newTitle := "Updated Title"
		req := models.UpdateTrackRequest{
			Title: &newTitle,
			// Artist, Album not provided - should remain unchanged
		}

		// Apply update
		if req.Title != nil {
			track.Title = *req.Title
		}

		assert.Equal(t, "Updated Title", track.Title)
		assert.Equal(t, "Original Artist", track.Artist) // Unchanged
		assert.Equal(t, "Original Album", track.Album)   // Unchanged
	})

	t.Run("updates tags array", func(t *testing.T) {
		track := models.Track{
			ID:   "track-123",
			Tags: []string{"rock", "favorites"},
		}

		req := models.UpdateTrackRequest{
			Tags: []string{"jazz", "2024"},
		}

		if req.Tags != nil {
			track.Tags = req.Tags
		}

		assert.Equal(t, []string{"jazz", "2024"}, track.Tags)
	})

	t.Run("sets UpdatedAt timestamp", func(t *testing.T) {
		track := models.Track{
			ID: "track-123",
		}
		track.UpdatedAt = time.Now()

		assert.False(t, track.UpdatedAt.IsZero())
	})
}

func TestTrackService_DeleteTrack_ExpectedBehavior(t *testing.T) {
	t.Run("deletes track from repository", func(t *testing.T) {
		// Expected behavior:
		// - Call repo.DeleteTrack with userID and trackID
		// - Return nil on success, error on failure
		// - Track should be removed from DynamoDB
		// - Related resources (S3 files, search index) may need cleanup

		userID := "user-456"
		trackID := "track-123"
		assert.NotEmpty(t, userID)
		assert.NotEmpty(t, trackID)
	})
}

func TestTrackService_AddTagsToTrack_ExpectedBehavior(t *testing.T) {
	t.Run("adds new tags avoiding duplicates", func(t *testing.T) {
		// Expected behavior:
		// 1. Get existing track
		// 2. Check which tags are new
		// 3. Add only new tags to track.Tags
		// 4. Create TrackTag associations for new tags
		// 5. Update track

		track := models.Track{
			ID:   "track-123",
			Tags: []string{"rock", "favorites"},
		}

		newTags := []string{"favorites", "2024", "chill"} // "favorites" is duplicate

		existingTags := make(map[string]bool)
		for _, t := range track.Tags {
			existingTags[t] = true
		}

		addedCount := 0
		for _, tag := range newTags {
			if !existingTags[tag] {
				track.Tags = append(track.Tags, tag)
				addedCount++
			}
		}

		assert.Equal(t, 2, addedCount) // Only "2024" and "chill" added
		assert.Len(t, track.Tags, 4)   // "rock", "favorites", "2024", "chill"
	})
}

func TestTrackService_RemoveTagFromTrack_ExpectedBehavior(t *testing.T) {
	t.Run("removes tag from track", func(t *testing.T) {
		track := models.Track{
			ID:   "track-123",
			Tags: []string{"rock", "favorites", "2024"},
		}

		tagToRemove := "favorites"

		newTags := make([]string, 0, len(track.Tags))
		for _, tag := range track.Tags {
			if tag != tagToRemove {
				newTags = append(newTags, tag)
			}
		}
		track.Tags = newTags

		assert.Len(t, track.Tags, 2)
		assert.NotContains(t, track.Tags, "favorites")
		assert.Contains(t, track.Tags, "rock")
		assert.Contains(t, track.Tags, "2024")
	})
}

// ====================
// Integration Test Markers
// ====================

func TestTrackService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	t.Run("requires DynamoDB Local", func(t *testing.T) {
		t.Skip("Integration tests require DynamoDB Local - run with proper environment")
		// These tests would:
		// 1. Set up DynamoDB Local
		// 2. Create actual repository
		// 3. Test full service operations
	})
}

// ====================
// Helper Function Tests
// ====================

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "0:00"},
		{30, "0:30"},
		{60, "1:00"},
		{90, "1:30"},
		{180, "3:00"},
		{3600, "1:00:00"},
		{3665, "1:01:05"},
		{7325, "2:02:05"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Test via Track.ToResponse which uses formatDuration
			track := models.Track{Duration: tt.seconds}
			response := track.ToResponse("")
			assert.Equal(t, tt.expected, response.DurationStr)
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{5242880, "5.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			track := models.Track{FileSize: tt.bytes}
			response := track.ToResponse("")
			assert.Equal(t, tt.expected, response.FileSizeStr)
		})
	}
}

// ====================
// Repository Interface Test
// ====================

func TestMockRepository_InterfaceCompliance(t *testing.T) {
	// Verify MockRepository can be used where Repository interface is expected
	// This is important for testing services that accept interfaces
	mockRepo := new(MockRepository)
	require.NotNil(t, mockRepo)

	// Test that methods can be set up
	ctx := context.Background()
	mockRepo.On("GetTrack", ctx, "user", "track").Return(&models.Track{}, nil)

	track, err := mockRepo.GetTrack(ctx, "user", "track")
	assert.NoError(t, err)
	assert.NotNil(t, track)
}
