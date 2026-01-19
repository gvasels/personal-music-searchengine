package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// MockRepository implements Repository interface for testing
type MockRepository struct {
	mock.Mock
}

// Compile-time check that MockRepository implements Repository
var _ Repository = (*MockRepository)(nil)

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

// Album operations
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

// Playlist operations
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

// Tag operations
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

func (m *MockRepository) AddTrackTag(ctx context.Context, tt *models.TrackTag) error {
	args := m.Called(ctx, tt)
	return args.Error(0)
}

func (m *MockRepository) RemoveTrackTag(ctx context.Context, userID, trackID, tagName string) error {
	args := m.Called(ctx, userID, trackID, tagName)
	return args.Error(0)
}

// Upload operations
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
// Interface Compliance Tests
// ====================

func TestRepositoryInterfaceCompliance(t *testing.T) {
	// This test verifies that MockRepository implements Repository
	var repo Repository = &MockRepository{}
	assert.NotNil(t, repo)
}

// ====================
// User Repository Tests
// ====================

func TestMockRepository_GetUser(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-123"

	expectedUser := &models.User{
		ID:           userID,
		Email:        "test@example.com",
		DisplayName:  "Test User",
		StorageUsed:  1024,
		StorageLimit: 10737418240,
	}

	mockRepo.On("GetUser", ctx, userID).Return(expectedUser, nil)

	user, err := mockRepo.GetUser(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetUser_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "nonexistent-user"

	mockRepo.On("GetUser", ctx, userID).Return(nil, models.NewNotFoundError("User", userID))

	user, err := mockRepo.GetUser(ctx, userID)
	assert.Nil(t, user)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_CreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	user := &models.User{
		ID:          "user-123",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	mockRepo.On("CreateUser", ctx, user).Return(nil)

	err := mockRepo.CreateUser(ctx, user)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_UpdateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	user := &models.User{
		ID:          "user-123",
		Email:       "test@example.com",
		DisplayName: "Updated Name",
	}

	mockRepo.On("UpdateUser", ctx, user).Return(nil)

	err := mockRepo.UpdateUser(ctx, user)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ====================
// Track Repository Tests
// ====================

func TestMockRepository_CreateTrack(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	track := &models.Track{
		ID:       "track-123",
		UserID:   "user-456",
		Title:    "Test Song",
		Artist:   "Test Artist",
		Duration: 180,
		Format:   models.AudioFormatMP3,
		FileSize: 5242880,
	}

	mockRepo.On("CreateTrack", ctx, track).Return(nil)

	err := mockRepo.CreateTrack(ctx, track)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetTrack(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	trackID := "track-123"

	expectedTrack := &models.Track{
		ID:       trackID,
		UserID:   userID,
		Title:    "Test Song",
		Artist:   "Test Artist",
		Duration: 180,
		Format:   models.AudioFormatMP3,
	}

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(expectedTrack, nil)

	track, err := mockRepo.GetTrack(ctx, userID, trackID)
	require.NoError(t, err)
	assert.Equal(t, expectedTrack.ID, track.ID)
	assert.Equal(t, expectedTrack.Title, track.Title)
	assert.Equal(t, expectedTrack.Artist, track.Artist)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetTrack_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	trackID := "nonexistent-track"

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(nil, models.NewNotFoundError("Track", trackID))

	track, err := mockRepo.GetTrack(ctx, userID, trackID)
	assert.Nil(t, track)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_UpdateTrack(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	track := &models.Track{
		ID:     "track-123",
		UserID: "user-456",
		Title:  "Updated Title",
		Artist: "Test Artist",
	}

	mockRepo.On("UpdateTrack", ctx, track).Return(nil)

	err := mockRepo.UpdateTrack(ctx, track)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_DeleteTrack(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	trackID := "track-123"

	mockRepo.On("DeleteTrack", ctx, userID, trackID).Return(nil)

	err := mockRepo.DeleteTrack(ctx, userID, trackID)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListTracks(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	filter := models.TrackFilter{
		Limit:     20,
		SortBy:    "title",
		SortOrder: "asc",
	}

	expectedTracks := &models.PaginatedResponse[models.Track]{
		Items: []models.Track{
			{ID: "track-1", Title: "Song A", Artist: "Artist"},
			{ID: "track-2", Title: "Song B", Artist: "Artist"},
		},
		Pagination: models.Pagination{
			Limit: 20,
		},
	}

	mockRepo.On("ListTracks", ctx, userID, filter).Return(expectedTracks, nil)

	result, err := mockRepo.ListTracks(ctx, userID, filter)
	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListTracks_WithFilters(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	filter := models.TrackFilter{
		Artist: "Test Artist",
		Genre:  "Rock",
		Year:   2024,
		Limit:  10,
	}

	expectedTracks := &models.PaginatedResponse[models.Track]{
		Items: []models.Track{
			{ID: "track-1", Title: "Song A", Artist: "Test Artist", Genre: "Rock", Year: 2024},
		},
		Pagination: models.Pagination{
			Limit: 10,
		},
	}

	mockRepo.On("ListTracks", ctx, userID, filter).Return(expectedTracks, nil)

	result, err := mockRepo.ListTracks(ctx, userID, filter)
	require.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "Test Artist", result.Items[0].Artist)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListTracksByAlbum(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	albumID := "album-123"

	expectedTracks := []models.Track{
		{ID: "track-1", Title: "Track 1", AlbumID: albumID},
		{ID: "track-2", Title: "Track 2", AlbumID: albumID},
	}

	mockRepo.On("ListTracksByAlbum", ctx, userID, albumID).Return(expectedTracks, nil)

	tracks, err := mockRepo.ListTracksByAlbum(ctx, userID, albumID)
	require.NoError(t, err)
	assert.Len(t, tracks, 2)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListTracksByTag(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	tagName := "favorites"

	expectedTracks := []models.Track{
		{ID: "track-1", Title: "Fav Song 1"},
		{ID: "track-2", Title: "Fav Song 2"},
	}

	mockRepo.On("ListTracksByTag", ctx, userID, tagName).Return(expectedTracks, nil)

	tracks, err := mockRepo.ListTracksByTag(ctx, userID, tagName)
	require.NoError(t, err)
	assert.Len(t, tracks, 2)
	mockRepo.AssertExpectations(t)
}

// ====================
// Album Repository Tests
// ====================

func TestMockRepository_CreateAlbum(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	album := &models.Album{
		ID:     "album-123",
		UserID: "user-456",
		Title:  "Test Album",
		Artist: "Test Artist",
		Year:   2024,
	}

	mockRepo.On("CreateAlbum", ctx, album).Return(nil)

	err := mockRepo.CreateAlbum(ctx, album)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetAlbum(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	albumID := "album-123"

	expectedAlbum := &models.Album{
		ID:         albumID,
		UserID:     userID,
		Title:      "Test Album",
		Artist:     "Test Artist",
		Year:       2024,
		TrackCount: 12,
	}

	mockRepo.On("GetAlbum", ctx, userID, albumID).Return(expectedAlbum, nil)

	album, err := mockRepo.GetAlbum(ctx, userID, albumID)
	require.NoError(t, err)
	assert.Equal(t, expectedAlbum.ID, album.ID)
	assert.Equal(t, expectedAlbum.Title, album.Title)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetOrCreateAlbum_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	title := "New Album"
	artist := "New Artist"
	year := 2024

	expectedAlbum := &models.Album{
		ID:     "album-new",
		UserID: userID,
		Title:  title,
		Artist: artist,
		Year:   year,
	}

	mockRepo.On("GetOrCreateAlbum", ctx, userID, title, artist, year).Return(expectedAlbum, nil)

	album, err := mockRepo.GetOrCreateAlbum(ctx, userID, title, artist, year)
	require.NoError(t, err)
	assert.Equal(t, title, album.Title)
	assert.Equal(t, artist, album.Artist)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListAlbums(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	filter := models.AlbumFilter{
		Limit:     20,
		SortBy:    "title",
		SortOrder: "asc",
	}

	expectedAlbums := &models.PaginatedResponse[models.Album]{
		Items: []models.Album{
			{ID: "album-1", Title: "Album A"},
			{ID: "album-2", Title: "Album B"},
		},
		Pagination: models.Pagination{Limit: 20},
	}

	mockRepo.On("ListAlbums", ctx, userID, filter).Return(expectedAlbums, nil)

	result, err := mockRepo.ListAlbums(ctx, userID, filter)
	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListAlbumsByArtist(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	artist := "Test Artist"

	expectedAlbums := []models.Album{
		{ID: "album-1", Title: "Album 1", Artist: artist},
		{ID: "album-2", Title: "Album 2", Artist: artist},
	}

	mockRepo.On("ListAlbumsByArtist", ctx, userID, artist).Return(expectedAlbums, nil)

	albums, err := mockRepo.ListAlbumsByArtist(ctx, userID, artist)
	require.NoError(t, err)
	assert.Len(t, albums, 2)
	mockRepo.AssertExpectations(t)
}

// ====================
// Playlist Repository Tests
// ====================

func TestMockRepository_CreatePlaylist(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	playlist := &models.Playlist{
		ID:          "playlist-123",
		UserID:      "user-456",
		Name:        "My Playlist",
		Description: "A test playlist",
	}

	mockRepo.On("CreatePlaylist", ctx, playlist).Return(nil)

	err := mockRepo.CreatePlaylist(ctx, playlist)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetPlaylist(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	playlistID := "playlist-123"

	expectedPlaylist := &models.Playlist{
		ID:          playlistID,
		UserID:      userID,
		Name:        "My Playlist",
		TrackCount:  10,
		IsPublic:    false,
	}

	mockRepo.On("GetPlaylist", ctx, userID, playlistID).Return(expectedPlaylist, nil)

	playlist, err := mockRepo.GetPlaylist(ctx, userID, playlistID)
	require.NoError(t, err)
	assert.Equal(t, expectedPlaylist.ID, playlist.ID)
	assert.Equal(t, expectedPlaylist.Name, playlist.Name)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_AddPlaylistTrack(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	pt := &models.PlaylistTrack{
		PlaylistID: "playlist-123",
		TrackID:    "track-456",
		Position:   5,
		AddedAt:    time.Now(),
	}

	mockRepo.On("AddPlaylistTrack", ctx, pt).Return(nil)

	err := mockRepo.AddPlaylistTrack(ctx, pt)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetPlaylistTracks(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	playlistID := "playlist-123"

	expectedTracks := []models.PlaylistTrack{
		{PlaylistID: playlistID, TrackID: "track-1", Position: 1},
		{PlaylistID: playlistID, TrackID: "track-2", Position: 2},
		{PlaylistID: playlistID, TrackID: "track-3", Position: 3},
	}

	mockRepo.On("GetPlaylistTracks", ctx, playlistID).Return(expectedTracks, nil)

	tracks, err := mockRepo.GetPlaylistTracks(ctx, playlistID)
	require.NoError(t, err)
	assert.Len(t, tracks, 3)
	// Verify order
	for i, track := range tracks {
		assert.Equal(t, i+1, track.Position)
	}
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_RemovePlaylistTrack(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	playlistID := "playlist-123"
	trackID := "track-456"

	mockRepo.On("RemovePlaylistTrack", ctx, playlistID, trackID).Return(nil)

	err := mockRepo.RemovePlaylistTrack(ctx, playlistID, trackID)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ReorderPlaylistTracks(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	playlistID := "playlist-123"
	trackIDs := []string{"track-3", "track-1", "track-2"}

	mockRepo.On("ReorderPlaylistTracks", ctx, playlistID, trackIDs).Return(nil)

	err := mockRepo.ReorderPlaylistTracks(ctx, playlistID, trackIDs)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ====================
// Tag Repository Tests
// ====================

func TestMockRepository_CreateTag(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	tag := &models.Tag{
		UserID: "user-456",
		Name:   "favorites",
		Color:  "#FF5733",
	}

	mockRepo.On("CreateTag", ctx, tag).Return(nil)

	err := mockRepo.CreateTag(ctx, tag)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetTag(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	tagName := "favorites"

	expectedTag := &models.Tag{
		UserID:     userID,
		Name:       tagName,
		Color:      "#FF5733",
		TrackCount: 25,
	}

	mockRepo.On("GetTag", ctx, userID, tagName).Return(expectedTag, nil)

	tag, err := mockRepo.GetTag(ctx, userID, tagName)
	require.NoError(t, err)
	assert.Equal(t, expectedTag.Name, tag.Name)
	assert.Equal(t, expectedTag.Color, tag.Color)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListTags(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	filter := models.TagFilter{
		Limit: 20,
	}

	expectedTags := &models.PaginatedResponse[models.Tag]{
		Items: []models.Tag{
			{Name: "favorites", Color: "#FF5733"},
			{Name: "rock", Color: "#3366FF"},
		},
		Pagination: models.Pagination{Limit: 20},
	}

	mockRepo.On("ListTags", ctx, userID, filter).Return(expectedTags, nil)

	result, err := mockRepo.ListTags(ctx, userID, filter)
	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_AddTrackTag(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	tt := &models.TrackTag{
		UserID:  "user-456",
		TrackID: "track-123",
		TagName: "favorites",
		AddedAt: time.Now(),
	}

	mockRepo.On("AddTrackTag", ctx, tt).Return(nil)

	err := mockRepo.AddTrackTag(ctx, tt)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_RemoveTrackTag(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	trackID := "track-123"
	tagName := "favorites"

	mockRepo.On("RemoveTrackTag", ctx, userID, trackID, tagName).Return(nil)

	err := mockRepo.RemoveTrackTag(ctx, userID, trackID, tagName)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ====================
// Upload Repository Tests
// ====================

func TestMockRepository_CreateUpload(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	upload := &models.Upload{
		ID:          "upload-123",
		UserID:      "user-456",
		FileName:    "song.mp3",
		FileSize:    5242880,
		ContentType: "audio/mpeg",
		Status:      models.UploadStatusPending,
	}

	mockRepo.On("CreateUpload", ctx, upload).Return(nil)

	err := mockRepo.CreateUpload(ctx, upload)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_GetUpload(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	uploadID := "upload-123"

	expectedUpload := &models.Upload{
		ID:          uploadID,
		UserID:      userID,
		FileName:    "song.mp3",
		FileSize:    5242880,
		ContentType: "audio/mpeg",
		Status:      models.UploadStatusCompleted,
		TrackID:     "track-789",
	}

	mockRepo.On("GetUpload", ctx, userID, uploadID).Return(expectedUpload, nil)

	upload, err := mockRepo.GetUpload(ctx, userID, uploadID)
	require.NoError(t, err)
	assert.Equal(t, expectedUpload.ID, upload.ID)
	assert.Equal(t, models.UploadStatusCompleted, upload.Status)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_UpdateUpload(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()

	upload := &models.Upload{
		ID:       "upload-123",
		UserID:   "user-456",
		Status:   models.UploadStatusCompleted,
		TrackID:  "track-789",
	}

	mockRepo.On("UpdateUpload", ctx, upload).Return(nil)

	err := mockRepo.UpdateUpload(ctx, upload)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_ListUploads(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"
	filter := models.UploadFilter{
		Status: models.UploadStatusCompleted,
		Limit:  20,
	}

	expectedUploads := &models.PaginatedResponse[models.Upload]{
		Items: []models.Upload{
			{ID: "upload-1", Status: models.UploadStatusCompleted},
			{ID: "upload-2", Status: models.UploadStatusCompleted},
		},
		Pagination: models.Pagination{Limit: 20},
	}

	mockRepo.On("ListUploads", ctx, userID, filter).Return(expectedUploads, nil)

	result, err := mockRepo.ListUploads(ctx, userID, filter)
	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	mockRepo.AssertExpectations(t)
}

// ====================
// Edge Case Tests
// ====================

func TestMockRepository_EmptyResults(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-with-no-tracks"
	filter := models.TrackFilter{Limit: 20}

	emptyResponse := &models.PaginatedResponse[models.Track]{
		Items:      []models.Track{},
		Pagination: models.Pagination{Limit: 20},
	}

	mockRepo.On("ListTracks", ctx, userID, filter).Return(emptyResponse, nil)

	result, err := mockRepo.ListTracks(ctx, userID, filter)
	require.NoError(t, err)
	assert.Empty(t, result.Items)
	mockRepo.AssertExpectations(t)
}

func TestMockRepository_Pagination(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	userID := "user-456"

	// First page
	filter1 := models.TrackFilter{Limit: 2}
	response1 := &models.PaginatedResponse[models.Track]{
		Items: []models.Track{
			{ID: "track-1", Title: "Song 1"},
			{ID: "track-2", Title: "Song 2"},
		},
		Pagination: models.Pagination{
			Limit:   2,
			NextKey: "eyJQSyI6IlVTRVIjdXNlci00NTYiLCJTSyI6IlRSQUNLI3RyYWNrLTIifQ==",
		},
	}

	mockRepo.On("ListTracks", ctx, userID, filter1).Return(response1, nil)

	result, err := mockRepo.ListTracks(ctx, userID, filter1)
	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.NotEmpty(t, result.Pagination.NextKey)
	mockRepo.AssertExpectations(t)
}
