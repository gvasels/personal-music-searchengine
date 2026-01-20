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

// MockTagRepository mocks the repository for tag tests
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) CreateTag(ctx context.Context, tag models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
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

func (m *MockTagRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	args := m.Called(ctx, userID, trackID, tagNames)
	return args.Error(0)
}

func (m *MockTagRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	args := m.Called(ctx, userID, trackID, tagName)
	return args.Error(0)
}

func (m *MockTagRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockTagRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	args := m.Called(ctx, userID, tagName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
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

// Stub implementations for Repository interface methods not used in tag tests
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

// Tests

func TestCreateTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(nil, repository.ErrNotFound)
	mockRepo.On("CreateTag", ctx, mock.MatchedBy(func(tag models.Tag) bool {
		return tag.Name == "Rock" && tag.UserID == "user-123"
	})).Return(nil)

	req := models.CreateTagRequest{Name: "Rock", Color: "#FF0000"}
	resp, err := svc.CreateTag(ctx, "user-123", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Rock", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateTag_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	existingTag := &models.Tag{UserID: "user-123", Name: "Rock"}
	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(existingTag, nil)

	req := models.CreateTagRequest{Name: "Rock"}
	resp, err := svc.CreateTag(ctx, "user-123", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestGetTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tag := &models.Tag{
		UserID:     "user-123",
		Name:       "Rock",
		Color:      "#FF0000",
		TrackCount: 5,
	}
	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(tag, nil)

	resp, err := svc.GetTag(ctx, "user-123", "Rock")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Rock", resp.Name)
	assert.Equal(t, "#FF0000", resp.Color)
	assert.Equal(t, 5, resp.TrackCount)
	mockRepo.AssertExpectations(t)
}

func TestGetTag_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "NonExistent").Return(nil, repository.ErrNotFound)

	resp, err := svc.GetTag(ctx, "user-123", "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tag := &models.Tag{UserID: "user-123", Name: "Rock", Color: "#FF0000"}
	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(tag, nil)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(t models.Tag) bool {
		return t.Color == "#00FF00"
	})).Return(nil)

	newColor := "#00FF00"
	req := models.UpdateTagRequest{Color: &newColor}
	resp, err := svc.UpdateTag(ctx, "user-123", "Rock", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "#00FF00", resp.Color)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_Rename(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tag := &models.Tag{UserID: "user-123", Name: "Rock", Color: "#FF0000"}
	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(tag, nil)
	mockRepo.On("GetTag", ctx, "user-123", "Metal").Return(nil, repository.ErrNotFound)
	mockRepo.On("UpdateTag", ctx, mock.MatchedBy(func(t models.Tag) bool {
		return t.Name == "Metal"
	})).Return(nil)

	newName := "Metal"
	req := models.UpdateTagRequest{Name: &newName}
	resp, err := svc.UpdateTag(ctx, "user-123", "Rock", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Metal", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTag_RenameConflict(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tag := &models.Tag{UserID: "user-123", Name: "Rock"}
	existingTag := &models.Tag{UserID: "user-123", Name: "Metal"}
	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(tag, nil)
	mockRepo.On("GetTag", ctx, "user-123", "Metal").Return(existingTag, nil)

	newName := "Metal"
	req := models.UpdateTagRequest{Name: &newName}
	resp, err := svc.UpdateTag(ctx, "user-123", "Rock", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestDeleteTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tag := &models.Tag{UserID: "user-123", Name: "Rock"}
	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(tag, nil)
	mockRepo.On("DeleteTag", ctx, "user-123", "Rock").Return(nil)

	err := svc.DeleteTag(ctx, "user-123", "Rock")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTag_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "NonExistent").Return(nil, repository.ErrNotFound)

	err := svc.DeleteTag(ctx, "user-123", "NonExistent")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListTags_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tags := []models.Tag{
		{UserID: "user-123", Name: "Rock", TrackCount: 5},
		{UserID: "user-123", Name: "Jazz", TrackCount: 3},
	}
	mockRepo.On("ListTags", ctx, "user-123").Return(tags, nil)

	resp, err := svc.ListTags(ctx, "user-123")

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Rock", resp[0].Name)
	assert.Equal(t, "Jazz", resp[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestListTags_Empty(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("ListTags", ctx, "user-123").Return([]models.Tag{}, nil)

	resp, err := svc.ListTags(ctx, "user-123")

	assert.NoError(t, err)
	assert.Empty(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestAddTagsToTrack_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	track := &models.Track{ID: "track-123", UserID: "user-123", Tags: []string{"Rock"}}
	mockRepo.On("GetTrack", ctx, "user-123", "track-123").Return(track, nil)
	mockRepo.On("GetTag", ctx, "user-123", "Jazz").Return(nil, repository.ErrNotFound)
	mockRepo.On("CreateTag", ctx, mock.MatchedBy(func(t models.Tag) bool {
		return t.Name == "Jazz"
	})).Return(nil)
	mockRepo.On("AddTagsToTrack", ctx, "user-123", "track-123", []string{"Jazz"}).Return(nil)
	mockRepo.On("UpdateTrack", ctx, mock.Anything).Return(nil)

	req := models.AddTagsToTrackRequest{Tags: []string{"Jazz"}}
	resp, err := svc.AddTagsToTrack(ctx, "user-123", "track-123", req)

	assert.NoError(t, err)
	assert.Contains(t, resp, "Rock")
	assert.Contains(t, resp, "Jazz")
	mockRepo.AssertExpectations(t)
}

func TestAddTagsToTrack_TrackNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTrack", ctx, "user-123", "nonexistent").Return(nil, repository.ErrNotFound)

	req := models.AddTagsToTrackRequest{Tags: []string{"Jazz"}}
	resp, err := svc.AddTagsToTrack(ctx, "user-123", "nonexistent", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTagFromTrack_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	track := &models.Track{ID: "track-123", UserID: "user-123", Tags: []string{"Rock", "Jazz"}}
	mockRepo.On("GetTrack", ctx, "user-123", "track-123").Return(track, nil)
	mockRepo.On("RemoveTagFromTrack", ctx, "user-123", "track-123", "Jazz").Return(nil)
	mockRepo.On("UpdateTrack", ctx, mock.MatchedBy(func(t models.Track) bool {
		return len(t.Tags) == 1 && t.Tags[0] == "Rock"
	})).Return(nil)

	err := svc.RemoveTagFromTrack(ctx, "user-123", "track-123", "Jazz")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTracksByTag_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	tag := &models.Tag{UserID: "user-123", Name: "Rock"}
	tracks := []models.Track{
		{ID: "track-1", Title: "Song 1", Artist: "Artist 1"},
		{ID: "track-2", Title: "Song 2", Artist: "Artist 2"},
	}
	tracks[0].CreatedAt = time.Now()
	tracks[0].UpdatedAt = time.Now()
	tracks[1].CreatedAt = time.Now()
	tracks[1].UpdatedAt = time.Now()

	mockRepo.On("GetTag", ctx, "user-123", "Rock").Return(tag, nil)
	mockRepo.On("GetTracksByTag", ctx, "user-123", "Rock").Return(tracks, nil)

	resp, err := svc.GetTracksByTag(ctx, "user-123", "Rock")

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Song 1", resp[0].Title)
	assert.Equal(t, "Song 2", resp[1].Title)
	mockRepo.AssertExpectations(t)
}

func TestGetTracksByTag_TagNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTagRepository)
	svc := NewTagService(mockRepo)

	mockRepo.On("GetTag", ctx, "user-123", "NonExistent").Return(nil, repository.ErrNotFound)

	resp, err := svc.GetTracksByTag(ctx, "user-123", "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}
