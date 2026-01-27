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
// Playlist Service Tests (Epic 4)
// =============================================================================

// MockPlaylistRepository provides mockable repository methods for playlist service tests
type MockPlaylistRepository struct {
	mock.Mock
}

// Playlist repository methods
func (m *MockPlaylistRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	args := m.Called(ctx, playlist)
	return args.Error(0)
}

func (m *MockPlaylistRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	args := m.Called(ctx, userID, playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	args := m.Called(ctx, playlist)
	return args.Error(0)
}

func (m *MockPlaylistRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	args := m.Called(ctx, userID, playlistID)
	return args.Error(0)
}

func (m *MockPlaylistRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Playlist]), args.Error(1)
}

func (m *MockPlaylistRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	args := m.Called(ctx, playlistID, trackIDs, position)
	return args.Error(0)
}

func (m *MockPlaylistRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	args := m.Called(ctx, playlistID, trackIDs)
	return args.Error(0)
}

func (m *MockPlaylistRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	args := m.Called(ctx, playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PlaylistTrack), args.Error(1)
}

func (m *MockPlaylistRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	args := m.Called(ctx, playlistID, tracks)
	return args.Error(0)
}

func (m *MockPlaylistRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

// Stub implementations for Repository interface (required but not used in playlist tests)
func (m *MockPlaylistRepository) CreateTrack(ctx context.Context, track models.Track) error    { return nil }
func (m *MockPlaylistRepository) UpdateTrack(ctx context.Context, track models.Track) error    { return nil }
func (m *MockPlaylistRepository) DeleteTrack(ctx context.Context, userID, trackID string) error { return nil }
func (m *MockPlaylistRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}
func (m *MockPlaylistRepository) CreateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockPlaylistRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockPlaylistRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}
func (m *MockPlaylistRepository) CreateTag(ctx context.Context, tag models.Tag) error { return nil }
func (m *MockPlaylistRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateTag(ctx context.Context, tag models.Tag) error { return nil }
func (m *MockPlaylistRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	return nil
}
func (m *MockPlaylistRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	return nil
}
func (m *MockPlaylistRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	return nil
}
func (m *MockPlaylistRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockPlaylistRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockPlaylistRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockPlaylistRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockPlaylistRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// Artist-related methods
func (m *MockPlaylistRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockPlaylistRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockPlaylistRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockPlaylistRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) MergeArtists(ctx context.Context, userID, primaryID string, secondaryIDs []string) error {
	return nil
}
func (m *MockPlaylistRepository) GetOrCreateArtist(ctx context.Context, userID, name string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) IncrementArtistStats(ctx context.Context, userID, artistID string, trackDelta, albumDelta int) error {
	return nil
}
func (m *MockPlaylistRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockPlaylistRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockPlaylistRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockPlaylistRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}

// User role methods
func (m *MockPlaylistRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockPlaylistRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}

// Playlist visibility methods
func (m *MockPlaylistRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	args := m.Called(ctx, userID, playlistID, visibility)
	return args.Error(0)
}
func (m *MockPlaylistRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	args := m.Called(ctx, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Playlist]), args.Error(1)
}

// ArtistProfile methods
func (m *MockPlaylistRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockPlaylistRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockPlaylistRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockPlaylistRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Follow methods
func (m *MockPlaylistRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockPlaylistRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockPlaylistRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Admin-related methods for track visibility
func (m *MockPlaylistRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}
func (m *MockPlaylistRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockPlaylistRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockPlaylistRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

// User settings methods
func (m *MockPlaylistRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]repository.UserSearchResult, string, error) {
	return nil, "", nil
}
func (m *MockPlaylistRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) UpdateUserSettings(ctx context.Context, userID string, update *repository.UserSettingsUpdate) (*models.UserSettings, error) {
	return nil, nil
}

// MockPlaylistS3Repository provides mockable S3 repository methods
type MockPlaylistS3Repository struct {
	mock.Mock
}

func (m *MockPlaylistS3Repository) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, ttl)
	return args.String(0), args.Error(1)
}

func (m *MockPlaylistS3Repository) GeneratePresignedDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	args := m.Called(ctx, key, ttl)
	return args.String(0), args.Error(1)
}

func (m *MockPlaylistS3Repository) DeleteObject(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Stub implementations for S3Repository interface
func (m *MockPlaylistS3Repository) InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	return "", nil
}

func (m *MockPlaylistS3Repository) GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error) {
	return nil, nil
}

func (m *MockPlaylistS3Repository) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error {
	return nil
}

func (m *MockPlaylistS3Repository) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return nil
}

func (m *MockPlaylistS3Repository) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	return nil
}

func (m *MockPlaylistS3Repository) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}

func (m *MockPlaylistS3Repository) ObjectExists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m *MockPlaylistS3Repository) GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, ttl time.Duration, filename string) (string, error) {
	return "", nil
}

// =============================================================================
// CreatePlaylist Tests
// =============================================================================

func TestCreatePlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("CreatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.UserID == "user-123" && p.Name == "My Playlist" && p.Description == "A great playlist"
	})).Return(nil)

	req := models.CreatePlaylistRequest{
		Name:        "My Playlist",
		Description: "A great playlist",
		IsPublic:    false,
	}
	resp, err := svc.CreatePlaylist(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "My Playlist", resp.Name)
	assert.Equal(t, "A great playlist", resp.Description)
	assert.False(t, resp.IsPublic)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// GetPlaylist Tests
// =============================================================================

func TestGetPlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:            "playlist-1",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    2,
		TotalDuration: 360,
		Timestamps:    models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("GetPlaylistTracks", ctx, "playlist-1").Return([]models.PlaylistTrack{
		{PlaylistID: "playlist-1", TrackID: "track-1", Position: 0},
		{PlaylistID: "playlist-1", TrackID: "track-2", Position: 1},
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:       "track-1",
		UserID:   "user-123",
		Title:    "Song 1",
		Duration: 180,
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "track-2").Return(&models.Track{
		ID:       "track-2",
		UserID:   "user-123",
		Title:    "Song 2",
		Duration: 180,
	}, nil)

	resp, err := svc.GetPlaylist(ctx, "user-123", "playlist-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "playlist-1", resp.Playlist.ID)
	assert.Equal(t, "My Playlist", resp.Playlist.Name)
	assert.Len(t, resp.Tracks, 2)
	assert.Equal(t, "Song 1", resp.Tracks[0].Title)
	assert.Equal(t, "Song 2", resp.Tracks[1].Title)
	mockRepo.AssertExpectations(t)
}

func TestGetPlaylist_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	resp, err := svc.GetPlaylist(ctx, "user-123", "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetPlaylist_WithCoverArt(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:          "playlist-1",
		UserID:      "user-123",
		Name:        "My Playlist",
		CoverArtKey: "covers/playlist-1.jpg",
		Timestamps:  models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("GetPlaylistTracks", ctx, "playlist-1").Return([]models.PlaylistTrack{}, nil)

	mockS3.On("GeneratePresignedDownloadURL", ctx, "covers/playlist-1.jpg", mock.Anything).Return("https://s3.example.com/covers/playlist-1.jpg?signed", nil)

	resp, err := svc.GetPlaylist(ctx, "user-123", "playlist-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "https://s3.example.com/covers/playlist-1.jpg?signed", resp.Playlist.CoverArtURL)
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}

// =============================================================================
// UpdatePlaylist Tests
// =============================================================================

func TestUpdatePlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:         "playlist-1",
		UserID:     "user-123",
		Name:       "Old Name",
		IsPublic:   false,
		Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.Name == "New Name" && p.IsPublic == true
	})).Return(nil)

	newName := "New Name"
	isPublic := true
	req := models.UpdatePlaylistRequest{
		Name:     &newName,
		IsPublic: &isPublic,
	}
	resp, err := svc.UpdatePlaylist(ctx, "user-123", "playlist-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Name", resp.Name)
	assert.True(t, resp.IsPublic)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePlaylist_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	newName := "New Name"
	req := models.UpdatePlaylistRequest{Name: &newName}
	resp, err := svc.UpdatePlaylist(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// DeletePlaylist Tests
// =============================================================================

func TestDeletePlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:     "playlist-1",
		UserID: "user-123",
		Name:   "My Playlist",
	}, nil)
	mockRepo.On("DeletePlaylist", ctx, "user-123", "playlist-1").Return(nil)

	err := svc.DeletePlaylist(ctx, "user-123", "playlist-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePlaylist_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	err := svc.DeletePlaylist(ctx, "user-123", "nonexistent")

	assert.Error(t, err)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// ListPlaylists Tests
// =============================================================================

func TestListPlaylists_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("ListPlaylists", ctx, "user-123", mock.Anything).Return(&repository.PaginatedResult[models.Playlist]{
		Items: []models.Playlist{
			{ID: "playlist-1", UserID: "user-123", Name: "Playlist 1", Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now}},
			{ID: "playlist-2", UserID: "user-123", Name: "Playlist 2", Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now}},
		},
		NextCursor: "",
		HasMore:    false,
	}, nil)

	// Mock GetPlaylistTracks calls for track count calculation
	mockRepo.On("GetPlaylistTracks", ctx, "playlist-1").Return([]models.PlaylistTrack{}, nil)
	mockRepo.On("GetPlaylistTracks", ctx, "playlist-2").Return([]models.PlaylistTrack{}, nil)

	filter := models.PlaylistFilter{Limit: 20}
	resp, err := svc.ListPlaylists(ctx, "user-123", filter)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, "Playlist 1", resp.Items[0].Name)
	assert.Equal(t, "Playlist 2", resp.Items[1].Name)
	assert.False(t, resp.HasMore)
	mockRepo.AssertExpectations(t)
}

func TestListPlaylists_Empty(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("ListPlaylists", ctx, "user-123", mock.Anything).Return(&repository.PaginatedResult[models.Playlist]{
		Items:      []models.Playlist{},
		NextCursor: "",
		HasMore:    false,
	}, nil)

	filter := models.PlaylistFilter{Limit: 20}
	resp, err := svc.ListPlaylists(ctx, "user-123", filter)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 0)
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// AddTracks Tests
// =============================================================================

func TestAddTracks_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:            "playlist-1",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    0,
		TotalDuration: 0,
		Timestamps:    models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:       "track-1",
		UserID:   "user-123",
		Duration: 180,
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "track-2").Return(&models.Track{
		ID:       "track-2",
		UserID:   "user-123",
		Duration: 200,
	}, nil)

	mockRepo.On("AddTracksToPlaylist", ctx, "playlist-1", []string{"track-1", "track-2"}, 0).Return(nil)

	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.TrackCount == 2 && p.TotalDuration == 380
	})).Return(nil)

	req := models.AddTracksToPlaylistRequest{
		TrackIDs: []string{"track-1", "track-2"},
	}
	resp, err := svc.AddTracks(ctx, "user-123", "playlist-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.TrackCount)
	assert.Equal(t, 380, resp.TotalDuration)
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_AtPosition(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:            "playlist-1",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    2,
		TotalDuration: 400,
		Timestamps:    models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "track-3").Return(&models.Track{
		ID:       "track-3",
		UserID:   "user-123",
		Duration: 150,
	}, nil)

	// Position 1 (insert in the middle)
	mockRepo.On("AddTracksToPlaylist", ctx, "playlist-1", []string{"track-3"}, 1).Return(nil)

	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.TrackCount == 3 && p.TotalDuration == 550
	})).Return(nil)

	position := 1
	req := models.AddTracksToPlaylistRequest{
		TrackIDs: []string{"track-3"},
		Position: &position,
	}
	resp, err := svc.AddTracks(ctx, "user-123", "playlist-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 3, resp.TrackCount)
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_TrackNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:         "playlist-1",
		UserID:     "user-123",
		Name:       "My Playlist",
		Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.AddTracksToPlaylistRequest{
		TrackIDs: []string{"nonexistent"},
	}
	resp, err := svc.AddTracks(ctx, "user-123", "playlist-1", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_PlaylistNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.AddTracksToPlaylistRequest{
		TrackIDs: []string{"track-1"},
	}
	resp, err := svc.AddTracks(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// RemoveTracks Tests
// =============================================================================

func TestRemoveTracks_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:            "playlist-1",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    3,
		TotalDuration: 600,
		Timestamps:    models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(&models.Track{
		ID:       "track-1",
		UserID:   "user-123",
		Duration: 200,
	}, nil)

	mockRepo.On("RemoveTracksFromPlaylist", ctx, "playlist-1", []string{"track-1"}).Return(nil)

	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.TrackCount == 2 && p.TotalDuration == 400
	})).Return(nil)

	req := models.RemoveTracksFromPlaylistRequest{
		TrackIDs: []string{"track-1"},
	}
	resp, err := svc.RemoveTracks(ctx, "user-123", "playlist-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.TrackCount)
	assert.Equal(t, 400, resp.TotalDuration)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTracks_PlaylistNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.RemoveTracksFromPlaylistRequest{
		TrackIDs: []string{"track-1"},
	}
	resp, err := svc.RemoveTracks(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Playlist Visibility Tests
// =============================================================================

func TestUpdatePlaylistVisibility_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-1").Return(&models.Playlist{
		ID:         "playlist-1",
		UserID:     "user-123",
		Name:       "My Playlist",
		Visibility: models.VisibilityPrivate,
		Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	mockRepo.On("UpdatePlaylistVisibility", ctx, "user-123", "playlist-1", models.VisibilityPublic).Return(nil)

	err := svc.UpdateVisibility(ctx, "user-123", "playlist-1", models.VisibilityPublic)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePlaylistVisibility_NotOwner(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("GetPlaylist", ctx, "other-user", "playlist-1").Return(&models.Playlist{
		ID:         "playlist-1",
		UserID:     "user-123", // Different from requester
		Name:       "My Playlist",
		Visibility: models.VisibilityPrivate,
		Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now},
	}, nil)

	err := svc.UpdateVisibility(ctx, "other-user", "playlist-1", models.VisibilityPublic)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FORBIDDEN")
}

func TestUpdatePlaylistVisibility_InvalidVisibility(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	err := svc.UpdateVisibility(ctx, "user-123", "playlist-1", models.PlaylistVisibility("invalid"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "VALIDATION_ERROR")
}

func TestUpdatePlaylistVisibility_PlaylistNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	err := svc.UpdateVisibility(ctx, "user-123", "nonexistent", models.VisibilityPublic)

	assert.Error(t, err)

	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	}
}

func TestListPublicPlaylists_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("ListPublicPlaylists", ctx, 20, "").Return(&repository.PaginatedResult[models.Playlist]{
		Items: []models.Playlist{
			{
				ID:          "playlist-1",
				UserID:      "user-1",
				Name:        "Public Playlist 1",
				Visibility:  models.VisibilityPublic,
				CreatorName: "Artist 1",
				Timestamps:  models.Timestamps{CreatedAt: now, UpdatedAt: now},
			},
			{
				ID:          "playlist-2",
				UserID:      "user-2",
				Name:        "Public Playlist 2",
				Visibility:  models.VisibilityPublic,
				CreatorName: "Artist 2",
				Timestamps:  models.Timestamps{CreatedAt: now, UpdatedAt: now},
			},
		},
		HasMore: false,
	}, nil)

	// Mock GetPlaylistTracks for track count
	mockRepo.On("GetPlaylistTracks", ctx, "playlist-1").Return([]models.PlaylistTrack{}, nil)
	mockRepo.On("GetPlaylistTracks", ctx, "playlist-2").Return([]models.PlaylistTrack{}, nil)

	result, err := svc.ListPublicPlaylists(ctx, 20, "")

	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, "Public Playlist 1", result.Items[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestListPublicPlaylists_WithPagination(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	mockRepo.On("ListPublicPlaylists", ctx, 1, "cursor-1").Return(&repository.PaginatedResult[models.Playlist]{
		Items: []models.Playlist{
			{
				ID:         "playlist-2",
				UserID:     "user-2",
				Name:       "Public Playlist 2",
				Visibility: models.VisibilityPublic,
				Timestamps: models.Timestamps{CreatedAt: now, UpdatedAt: now},
			},
		},
		NextCursor: "cursor-2",
		HasMore:    true,
	}, nil)

	mockRepo.On("GetPlaylistTracks", ctx, "playlist-2").Return([]models.PlaylistTrack{}, nil)

	result, err := svc.ListPublicPlaylists(ctx, 1, "cursor-1")

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.True(t, result.HasMore)
	assert.Equal(t, "cursor-2", result.NextCursor)
	mockRepo.AssertExpectations(t)
}

func TestListPublicPlaylists_Empty(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("ListPublicPlaylists", ctx, 20, "").Return(&repository.PaginatedResult[models.Playlist]{
		Items:   []models.Playlist{},
		HasMore: false,
	}, nil)

	result, err := svc.ListPublicPlaylists(ctx, 20, "")

	assert.NoError(t, err)
	assert.Empty(t, result.Items)
	mockRepo.AssertExpectations(t)
}
