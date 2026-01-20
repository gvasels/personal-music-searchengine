package service

import (
	"context"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlaylistRepository mocks the repository for playlist tests
type MockPlaylistRepository struct {
	mock.Mock
}

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

func (m *MockPlaylistRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

// Stub implementations for Repository interface methods not used in playlist tests
func (m *MockPlaylistRepository) CreateTrack(ctx context.Context, track models.Track) error { return nil }
func (m *MockPlaylistRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	return nil
}
func (m *MockPlaylistRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return nil
}
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
func (m *MockPlaylistRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockPlaylistRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
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

// MockPlaylistS3Repository mocks S3 repository for playlist tests
type MockPlaylistS3Repository struct {
	mock.Mock
}

func (m *MockPlaylistS3Repository) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	return "", nil
}

func (m *MockPlaylistS3Repository) GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

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
func (m *MockPlaylistS3Repository) DeleteObject(ctx context.Context, key string) error { return nil }
func (m *MockPlaylistS3Repository) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	return nil
}
func (m *MockPlaylistS3Repository) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}
func (m *MockPlaylistS3Repository) ObjectExists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// Tests

func TestCreatePlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("CreatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.Name == "My Playlist" && p.UserID == "user-123"
	})).Return(nil)

	req := models.CreatePlaylistRequest{Name: "My Playlist", Description: "Test playlist"}
	resp, err := svc.CreatePlaylist(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "My Playlist", resp.Name)
	assert.Equal(t, "Test playlist", resp.Description)
	assert.NotEmpty(t, resp.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetPlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlist := &models.Playlist{
		ID:            "playlist-123",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    2,
		TotalDuration: 400,
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	playlistTracks := []models.PlaylistTrack{
		{PlaylistID: "playlist-123", TrackID: "track-1", Position: 0},
		{PlaylistID: "playlist-123", TrackID: "track-2", Position: 1},
	}

	track1 := &models.Track{ID: "track-1", Title: "Song 1", Duration: 200}
	track1.CreatedAt = now
	track1.UpdatedAt = now
	track2 := &models.Track{ID: "track-2", Title: "Song 2", Duration: 200}
	track2.CreatedAt = now
	track2.UpdatedAt = now

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("GetPlaylistTracks", ctx, "playlist-123").Return(playlistTracks, nil)
	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(track1, nil)
	mockRepo.On("GetTrack", ctx, "user-123", "track-2").Return(track2, nil)

	resp, err := svc.GetPlaylist(ctx, "user-123", "playlist-123")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
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
	mockRepo.AssertExpectations(t)
}

func TestUpdatePlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlist := &models.Playlist{
		ID:     "playlist-123",
		UserID: "user-123",
		Name:   "Old Name",
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.Name == "New Name"
	})).Return(nil)

	newName := "New Name"
	req := models.UpdatePlaylistRequest{Name: &newName}
	resp, err := svc.UpdatePlaylist(ctx, "user-123", "playlist-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Name", resp.Name)
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
	mockRepo.AssertExpectations(t)
}

func TestDeletePlaylist_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	playlist := &models.Playlist{ID: "playlist-123", UserID: "user-123"}
	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("DeletePlaylist", ctx, "user-123", "playlist-123").Return(nil)

	err := svc.DeletePlaylist(ctx, "user-123", "playlist-123")

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
	mockRepo.AssertExpectations(t)
}

func TestListPlaylists_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlists := []models.Playlist{
		{ID: "p1", Name: "Playlist 1", UserID: "user-123"},
		{ID: "p2", Name: "Playlist 2", UserID: "user-123"},
	}
	playlists[0].CreatedAt = now
	playlists[0].UpdatedAt = now
	playlists[1].CreatedAt = now
	playlists[1].UpdatedAt = now

	result := &repository.PaginatedResult[models.Playlist]{
		Items:   playlists,
		HasMore: false,
	}
	mockRepo.On("ListPlaylists", ctx, "user-123", models.PlaylistFilter{}).Return(result, nil)

	resp, err := svc.ListPlaylists(ctx, "user-123", models.PlaylistFilter{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, "Playlist 1", resp.Items[0].Name)
	assert.Equal(t, "Playlist 2", resp.Items[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlist := &models.Playlist{
		ID:            "playlist-123",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    0,
		TotalDuration: 0,
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	track := &models.Track{ID: "track-1", Duration: 200}
	track.CreatedAt = now
	track.UpdatedAt = now

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(track, nil)
	mockRepo.On("AddTracksToPlaylist", ctx, "playlist-123", []string{"track-1"}, 0).Return(nil)
	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.TrackCount == 1 && p.TotalDuration == 200
	})).Return(nil)

	req := models.AddTracksToPlaylistRequest{TrackIDs: []string{"track-1"}}
	resp, err := svc.AddTracks(ctx, "user-123", "playlist-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.TrackCount)
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_PlaylistNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.AddTracksToPlaylistRequest{TrackIDs: []string{"track-1"}}
	resp, err := svc.AddTracks(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_TrackNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlist := &models.Playlist{ID: "playlist-123", UserID: "user-123"}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("GetTrack", ctx, "user-123", "nonexistent-track").Return(nil, repository.ErrNotFound)

	req := models.AddTracksToPlaylistRequest{TrackIDs: []string{"nonexistent-track"}}
	resp, err := svc.AddTracks(ctx, "user-123", "playlist-123", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTracks_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlist := &models.Playlist{
		ID:            "playlist-123",
		UserID:        "user-123",
		Name:          "My Playlist",
		TrackCount:    2,
		TotalDuration: 400,
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	track := &models.Track{ID: "track-1", Duration: 200}
	track.CreatedAt = now
	track.UpdatedAt = now

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("GetTrack", ctx, "user-123", "track-1").Return(track, nil)
	mockRepo.On("RemoveTracksFromPlaylist", ctx, "playlist-123", []string{"track-1"}).Return(nil)
	mockRepo.On("UpdatePlaylist", ctx, mock.MatchedBy(func(p models.Playlist) bool {
		return p.TrackCount == 1 && p.TotalDuration == 200
	})).Return(nil)

	req := models.RemoveTracksFromPlaylistRequest{TrackIDs: []string{"track-1"}}
	resp, err := svc.RemoveTracks(ctx, "user-123", "playlist-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.TrackCount)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTracks_PlaylistNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	mockRepo.On("GetPlaylist", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.RemoveTracksFromPlaylistRequest{TrackIDs: []string{"track-1"}}
	resp, err := svc.RemoveTracks(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestAddTracks_WithPosition(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPlaylistRepository)
	mockS3 := new(MockPlaylistS3Repository)
	svc := NewPlaylistService(mockRepo, mockS3)

	now := time.Now()
	playlist := &models.Playlist{
		ID:            "playlist-123",
		UserID:        "user-123",
		TrackCount:    5,
		TotalDuration: 1000,
	}
	playlist.CreatedAt = now
	playlist.UpdatedAt = now

	track := &models.Track{ID: "track-new", Duration: 180}
	track.CreatedAt = now
	track.UpdatedAt = now

	mockRepo.On("GetPlaylist", ctx, "user-123", "playlist-123").Return(playlist, nil)
	mockRepo.On("GetTrack", ctx, "user-123", "track-new").Return(track, nil)
	// Verify position is used correctly
	mockRepo.On("AddTracksToPlaylist", ctx, "playlist-123", []string{"track-new"}, 2).Return(nil)
	mockRepo.On("UpdatePlaylist", ctx, mock.Anything).Return(nil)

	position := 2
	req := models.AddTracksToPlaylistRequest{TrackIDs: []string{"track-new"}, Position: &position}
	resp, err := svc.AddTracks(ctx, "user-123", "playlist-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockRepo.AssertExpectations(t)
}
