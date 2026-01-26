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
)

// =============================================================================
// Tag Service Tests (Epic 4)
// =============================================================================

// MockTagRepository provides mockable repository methods for tag service tests
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockTagRepository) CreateTag(ctx context.Context, tag models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepository) UpdateTag(ctx context.Context, tag models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	args := m.Called(ctx, userID, tagName)
	return args.Error(0)
}

func (m *MockTagRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Tag), args.Error(1)
}

func (m *MockTagRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockTagRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockTagRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	args := m.Called(ctx, userID, trackID, tagNames)
	return args.Error(0)
}

func (m *MockTagRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	args := m.Called(ctx, userID, trackID, tagName)
	return args.Error(0)
}

func (m *MockTagRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

// Stub implementations for Repository interface (required but not used in tag tests)
func (m *MockTagRepository) CreateTrack(ctx context.Context, track models.Track) error { return nil }
func (m *MockTagRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return nil
}
func (m *MockTagRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockTagRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockTagRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockTagRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockTagRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockTagRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockTagRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}
func (m *MockTagRepository) CreateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockTagRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockTagRepository) UpdateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockTagRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}
func (m *MockTagRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockTagRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	return nil, nil
}
func (m *MockTagRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockTagRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return nil
}
func (m *MockTagRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockTagRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	return nil
}
func (m *MockTagRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	return nil
}
func (m *MockTagRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockTagRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	return nil
}
func (m *MockTagRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockTagRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockTagRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockTagRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockTagRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockTagRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockTagRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockTagRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// Artist-related methods
func (m *MockTagRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockTagRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockTagRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockTagRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockTagRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockTagRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockTagRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockTagRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockTagRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockTagRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockTagRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockTagRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}

// User role methods
func (m *MockTagRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockTagRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}

// Playlist visibility methods
func (m *MockTagRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	return nil
}
func (m *MockTagRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}

// ArtistProfile methods
func (m *MockTagRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockTagRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockTagRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockTagRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockTagRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockTagRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Follow methods
func (m *MockTagRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockTagRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockTagRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockTagRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockTagRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockTagRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Admin-related methods for track visibility
func (m *MockTagRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockTagRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}
func (m *MockTagRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockTagRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockTagRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockTagRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

// =============================================================================
// CreateTag Tests
// =============================================================================

func TestCreateTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Tag doesn't exist yet
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(nil, repository.ErrNotFound)
	mockRepo.On("CreateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.UserID == "user-123" && tag.Name == "favorites" && tag.Color == "#FF0000"
	})).Return(nil)

	req := models.CreateTagRequest{
		Name:  "favorites",
		Color: "#FF0000",
	}
	resp, err := svc.CreateTag(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "favorites", resp.Name)
	assert.Equal(t, "#FF0000", resp.Color)
	mockRepo.AssertExpectations(t)
}

func TestCreateTag_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Tag already exists
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{
		UserID: "user-123",
		Name:   "favorites",
	}, nil)

	req := models.CreateTagRequest{Name: "favorites"}
	resp, err := svc.CreateTag(ctx, "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "CONFLICT", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// GetTag Tests
// =============================================================================

func TestGetTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{
		UserID:     "user-123",
		Name:       "favorites",
		Color:      "#FF0000",
		TrackCount: 5,
	}, nil)

	resp, err := svc.GetTag(ctx, "user-123", "favorites")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "favorites", resp.Name)
	assert.Equal(t, "#FF0000", resp.Color)
	assert.Equal(t, 5, resp.TrackCount)
	mockRepo.AssertExpectations(t)
}

func TestGetTag_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	resp, err := svc.GetTag(ctx, "user-123", "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// UpdateTag Tests
// =============================================================================

func TestUpdateTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{
		UserID: "user-123",
		Name:   "favorites",
		Color:  "#FF0000",
	}, nil)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Color == "#00FF00"
	})).Return(nil)

	newColor := "#00FF00"
	req := models.UpdateTagRequest{Color: &newColor}
	resp, err := svc.UpdateTag(ctx, "user-123", "favorites", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "#00FF00", resp.Color)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_Rename(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Original tag exists
	mockRepo.On("GetTag", ctx, "user-123", "old-name").Return(&models.Tag{
		UserID: "user-123",
		Name:   "old-name",
	}, nil)
	// New name doesn't exist
	mockRepo.On("GetTag", ctx, "user-123", "new-name").Return(nil, repository.ErrNotFound)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "new-name"
	})).Return(nil)

	newName := "new-name"
	req := models.UpdateTagRequest{Name: &newName}
	resp, err := svc.UpdateTag(ctx, "user-123", "old-name", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new-name", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_RenameConflict(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Original tag exists
	mockRepo.On("GetTag", ctx, "user-123", "old-name").Return(&models.Tag{
		UserID: "user-123",
		Name:   "old-name",
	}, nil)
	// New name already exists
	mockRepo.On("GetTag", ctx, "user-123", "existing-name").Return(&models.Tag{
		UserID: "user-123",
		Name:   "existing-name",
	}, nil)

	newName := "existing-name"
	req := models.UpdateTagRequest{Name: &newName}
	resp, err := svc.UpdateTag(ctx, "user-123", "old-name", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "CONFLICT", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	newColor := "#00FF00"
	req := models.UpdateTagRequest{Color: &newColor}
	resp, err := svc.UpdateTag(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// DeleteTag Tests
// =============================================================================

func TestDeleteTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{
		UserID: "user-123",
		Name:   "favorites",
	}, nil)
	mockRepo.On("DeleteTag", ctx, "user-123", "favorites").Return(nil)

	err := svc.DeleteTag(ctx, "user-123", "favorites")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTag_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	err := svc.DeleteTag(ctx, "user-123", "nonexistent")

	assert.Error(t, err)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// ListTags Tests
// =============================================================================

func TestListTags_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	now := time.Now()
	mockRepo.On("ListTags", ctx, "user-123").Return([]models.Tag{
		{UserID: "user-123", Name: "favorites", Color: "#FF0000", TrackCount: 5, Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now}},
		{UserID: "user-123", Name: "rock", Color: "#00FF00", TrackCount: 10, Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now}},
	}, nil)

	resp, err := svc.ListTags(ctx, "user-123")

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "favorites", resp[0].Name)
	assert.Equal(t, "rock", resp[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestListTags_Empty(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("ListTags", ctx, "user-123").Return([]models.Tag{}, nil)

	resp, err := svc.ListTags(ctx, "user-123")

	assert.NoError(t, err)
	assert.Len(t, resp, 0)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// AddTagsToTrack Tests
// =============================================================================

func TestAddTagsToTrack_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Track exists
	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:     "track-1",
		UserID: "user-123",
		Tags:   []string{"existing-tag"},
	}, nil)

	// One tag doesn't exist (will be created with TrackCount=1)
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(nil, repository.ErrNotFound)
	mockRepo.On("CreateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "favorites" && tag.TrackCount == 1
	})).Return(nil)

	// One tag exists (TrackCount will be incremented)
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{Name: "rock", TrackCount: 5}, nil)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "rock" && tag.TrackCount == 6
	})).Return(nil)

	mockRepo.On("AddTagsToTrack", ctx, "user-123", "track-1", mock.Anything).Return(nil)
	mockRepo.On("UpdateTrack", ctx, mock.Anything).Return(nil)

	req := models.AddTagsToTrackRequest{Tags: []string{"favorites", "rock"}}
	tags, err := svc.AddTagsToTrack(ctx, "user-123", "track-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, tags)
	// Should contain existing-tag, favorites, and rock
	assert.Contains(t, tags, "existing-tag")
	assert.Contains(t, tags, "favorites")
	assert.Contains(t, tags, "rock")
	mockRepo.AssertExpectations(t)
}

func TestAddTagsToTrack_TrackNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTrack", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.AddTagsToTrackRequest{Tags: []string{"favorites"}}
	tags, err := svc.AddTagsToTrack(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, tags)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// RemoveTagFromTrack Tests
// =============================================================================

func TestRemoveTagFromTrack_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:     "track-1",
		UserID: "user-123",
		Tags:   []string{"favorites", "rock"},
	}, nil)
	mockRepo.On("RemoveTagFromTrack", ctx, "user-123", "track-1", "favorites").Return(nil)
	// Tag TrackCount should be decremented
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{Name: "favorites", TrackCount: 3}, nil)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "favorites" && tag.TrackCount == 2
	})).Return(nil)
	mockRepo.On("UpdateTrack", ctx, mock.MatchedBy(func(track models.Track) bool {
		return len(track.Tags) == 1 && track.Tags[0] == "rock"
	})).Return(nil)

	err := svc.RemoveTagFromTrack(ctx, "user-123", "track-1", "favorites")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// GetTracksByTag Tests
// =============================================================================

func TestGetTracksByTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{
		UserID: "user-123",
		Name:   "favorites",
	}, nil)
	mockRepo.On("GetTracksByTag", ctx, "user-123", "favorites").Return([]models.Track{
		{ID: "track-1", UserID: "user-123", Title: "Song 1"},
		{ID: "track-2", UserID: "user-123", Title: "Song 2"},
	}, nil)

	resp, err := svc.GetTracksByTag(ctx, "user-123", "favorites")

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "track-1", resp[0].ID)
	assert.Equal(t, "track-2", resp[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestGetTracksByTag_TagNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	resp, err := svc.GetTracksByTag(ctx, "user-123", "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Tag Name Normalization Tests (Task 10.1)
// =============================================================================

func TestCreateTag_NormalizesName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// When creating "Rock", should check for and create as "rock"
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(nil, repository.ErrNotFound)
	mockRepo.On("CreateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "rock" // stored as lowercase
	})).Return(nil)

	req := models.CreateTagRequest{
		Name:  "Rock", // mixed case input
		Color: "#FF0000",
	}
	resp, err := svc.CreateTag(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "rock", resp.Name) // returned as lowercase
	mockRepo.AssertExpectations(t)
}

func TestGetTag_CaseInsensitive(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// When looking up "ROCK", should query for "rock"
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{
		UserID: "user-123",
		Name:   "rock",
		Color:  "#FF0000",
	}, nil)

	resp, err := svc.GetTag(ctx, "user-123", "ROCK") // uppercase input

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "rock", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTag_NormalizesName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{
		UserID: "user-123",
		Name:   "rock",
	}, nil)
	mockRepo.On("DeleteTag", ctx, "user-123", "rock").Return(nil)

	err := svc.DeleteTag(ctx, "user-123", "ROCK") // uppercase input

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddTagsToTrack_NormalizesNames(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Track exists
	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:     "track-1",
		UserID: "user-123",
		Tags:   []string{},
	}, nil)

	// Tags should be looked up/created as lowercase
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(nil, repository.ErrNotFound)
	mockRepo.On("CreateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "rock" && tag.TrackCount == 1 // lowercase with count
	})).Return(nil)
	mockRepo.On("GetTag", ctx, "user-123", "favorites").Return(&models.Tag{Name: "favorites", TrackCount: 2}, nil)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "favorites" && tag.TrackCount == 3
	})).Return(nil)

	// Should add with normalized names
	mockRepo.On("AddTagsToTrack", ctx, "user-123", "track-1", mock.Anything).Return(nil)
	mockRepo.On("UpdateTrack", ctx, mock.Anything).Return(nil)

	req := models.AddTagsToTrackRequest{Tags: []string{"ROCK", "Favorites"}} // mixed case
	tags, err := svc.AddTagsToTrack(ctx, "user-123", "track-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, tags)
	// Returned tags should be lowercase
	assert.Contains(t, tags, "rock")
	assert.Contains(t, tags, "favorites")
	mockRepo.AssertExpectations(t)
}

func TestRemoveTagFromTrack_NormalizesName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:     "track-1",
		UserID: "user-123",
		Tags:   []string{"rock", "favorites"},
	}, nil)
	mockRepo.On("RemoveTagFromTrack", ctx, "user-123", "track-1", "rock").Return(nil)
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{Name: "rock", TrackCount: 5}, nil)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "rock" && tag.TrackCount == 4
	})).Return(nil)
	mockRepo.On("UpdateTrack", ctx, mock.MatchedBy(func(track models.Track) bool {
		return len(track.Tags) == 1 && track.Tags[0] == "favorites"
	})).Return(nil)

	err := svc.RemoveTagFromTrack(ctx, "user-123", "track-1", "ROCK") // uppercase

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTracksByTag_NormalizesName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{
		UserID: "user-123",
		Name:   "rock",
	}, nil)
	mockRepo.On("GetTracksByTag", ctx, "user-123", "rock").Return([]models.Track{
		{ID: "track-1", UserID: "user-123", Title: "Song 1"},
	}, nil)

	resp, err := svc.GetTracksByTag(ctx, "user-123", "ROCK") // uppercase

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_NormalizesNames(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	// Lookup should be normalized
	mockRepo.On("GetTag", ctx, "user-123", "rock").Return(&models.Tag{
		UserID: "user-123",
		Name:   "rock",
		Color:  "#FF0000",
	}, nil)
	// New name should not exist (normalized)
	mockRepo.On("GetTag", ctx, "user-123", "metal").Return(nil, repository.ErrNotFound)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "metal" // new name normalized
	})).Return(nil)

	newName := "METAL" // uppercase input
	req := models.UpdateTagRequest{Name: &newName}
	resp, err := svc.UpdateTag(ctx, "user-123", "ROCK", req) // uppercase old name

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "metal", resp.Name)
	mockRepo.AssertExpectations(t)
}
