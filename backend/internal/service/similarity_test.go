package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSimilarityRepository is a mock implementation of repository.Repository for similarity tests
type MockSimilarityRepository struct {
	mock.Mock
}

func (m *MockSimilarityRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockSimilarityRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

// Implement other repository methods as stubs
func (m *MockSimilarityRepository) CreateTrack(ctx context.Context, track models.Track) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	return nil
}
func (m *MockSimilarityRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return nil
}
func (m *MockSimilarityRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) CreateAlbum(ctx context.Context, album models.Album) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateAlbum(ctx context.Context, album models.Album) error {
	return nil
}
func (m *MockSimilarityRepository) DeleteAlbum(ctx context.Context, userID, albumID string) error {
	return nil
}
func (m *MockSimilarityRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) CreateUser(ctx context.Context, user models.User) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateUser(ctx context.Context, user models.User) error {
	return nil
}
func (m *MockSimilarityRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockSimilarityRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockSimilarityRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return nil
}
func (m *MockSimilarityRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockSimilarityRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) CreateTag(ctx context.Context, tag models.Tag) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateTag(ctx context.Context, tag models.Tag) error {
	return nil
}
func (m *MockSimilarityRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	return nil
}
func (m *MockSimilarityRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	return nil
}
func (m *MockSimilarityRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	return nil
}
func (m *MockSimilarityRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockSimilarityRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockSimilarityRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockSimilarityRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockSimilarityRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockSimilarityRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}
func (m *MockSimilarityRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	return nil
}
func (m *MockSimilarityRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	return nil
}
func (m *MockSimilarityRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	return nil
}
func (m *MockSimilarityRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockSimilarityRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockSimilarityRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// User role methods
func (m *MockSimilarityRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockSimilarityRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}

// Playlist visibility methods
func (m *MockSimilarityRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	return nil
}
func (m *MockSimilarityRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}

// ArtistProfile methods
func (m *MockSimilarityRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockSimilarityRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockSimilarityRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockSimilarityRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Follow methods
func (m *MockSimilarityRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockSimilarityRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockSimilarityRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Admin-related methods for track visibility
func (m *MockSimilarityRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}
func (m *MockSimilarityRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}
func (m *MockSimilarityRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockSimilarityRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockSimilarityRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

// User settings methods
func (m *MockSimilarityRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]repository.UserSearchResult, string, error) {
	return nil, "", nil
}
func (m *MockSimilarityRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockSimilarityRepository) UpdateUserSettings(ctx context.Context, userID string, update *repository.UserSettingsUpdate) (*models.UserSettings, error) {
	return nil, nil
}

// Helper function to create test tracks
func createSimilarityTestTrack(id, artist, album, genre, keyCamelot string, bpm int, tags []string) models.Track {
	return models.Track{
		ID:         id,
		UserID:     "user-123",
		Title:      "Track " + id,
		Artist:     artist,
		Album:      album,
		Genre:      genre,
		BPM:        bpm,
		KeyCamelot: keyCamelot,
		Tags:       tags,
	}
}

func TestDefaultSimilarityOptions(t *testing.T) {
	opts := DefaultSimilarityOptions()

	assert.Equal(t, 10, opts.Limit)
	assert.Equal(t, "combined", opts.Mode)
	assert.Equal(t, 0.5, opts.MinSimilarity)
	assert.True(t, opts.IncludeSameAlbum)
}

func TestDefaultMixingOptions(t *testing.T) {
	opts := DefaultMixingOptions()

	assert.Equal(t, 10, opts.Limit)
	assert.Equal(t, 5, opts.BPMTolerance)
	assert.Equal(t, "harmonic", opts.KeyMode)
}

func TestNewSimilarityService(t *testing.T) {
	mockRepo := new(MockSimilarityRepository)

	svc := NewSimilarityService(nil, mockRepo, nil)

	require.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.repo)
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 1, 0},
			b:        []float32{1, 0, 0},
			expected: 0.7071067811865475, // 1/sqrt(2)
		},
		{
			name:     "empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0.0,
		},
		{
			name:     "different length vectors",
			a:        []float32{1, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
		},
		{
			name:     "zero vector",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}

func TestCountOverlappingTags(t *testing.T) {
	tests := []struct {
		name     string
		tags1    []string
		tags2    []string
		expected int
	}{
		{
			name:     "no overlap",
			tags1:    []string{"rock", "metal"},
			tags2:    []string{"pop", "dance"},
			expected: 0,
		},
		{
			name:     "partial overlap",
			tags1:    []string{"rock", "metal", "classic"},
			tags2:    []string{"rock", "pop"},
			expected: 1,
		},
		{
			name:     "full overlap",
			tags1:    []string{"rock", "metal"},
			tags2:    []string{"rock", "metal"},
			expected: 2,
		},
		{
			name:     "empty first list",
			tags1:    []string{},
			tags2:    []string{"rock"},
			expected: 0,
		},
		{
			name:     "empty second list",
			tags1:    []string{"rock"},
			tags2:    []string{},
			expected: 0,
		},
		{
			name:     "both empty",
			tags1:    []string{},
			tags2:    []string{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countOverlappingTags(tt.tags1, tt.tags2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMax(t *testing.T) {
	assert.Equal(t, 5, max(3, 5))
	assert.Equal(t, 5, max(5, 3))
	assert.Equal(t, 5, max(5, 5))
	assert.Equal(t, 0, max(-5, 0))
	assert.Equal(t, -3, max(-5, -3))
}

func TestFindSimilarTracks_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, []string{"energetic"})
	similarTrack := createSimilarityTestTrack("track-2", "Artist A", "Album 2", "Rock", "8A", 122, []string{"energetic"})
	differentTrack := createSimilarityTestTrack("track-3", "Artist B", "Album 3", "Jazz", "1B", 80, []string{"calm"})

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, similarTrack, differentTrack},
		HasMore: false,
	}, nil)

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, DefaultSimilarityOptions())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, sourceTrack.ID, result.SourceTrack.ID)
	assert.NotEmpty(t, result.Similar)
	mockRepo.AssertExpectations(t)
}

func TestFindSimilarTracks_SourceTrackNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "nonexistent"
	mockRepo := new(MockSimilarityRepository)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(nil, errors.New("track not found"))

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, DefaultSimilarityOptions())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get source track")
}

func TestFindSimilarTracks_ExcludeSameAlbum(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, nil)
	sameAlbumTrack := createSimilarityTestTrack("track-2", "Artist A", "Album 1", "Rock", "8A", 122, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, sameAlbumTrack},
		HasMore: false,
	}, nil)

	opts := DefaultSimilarityOptions()
	opts.IncludeSameAlbum = false

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Same album track should be excluded
	for _, similar := range result.Similar {
		assert.NotEqual(t, sameAlbumTrack.ID, similar.Track.ID)
	}
}

func TestFindSimilarTracks_SemanticMode(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, []string{"rock"})
	similarTrack := createSimilarityTestTrack("track-2", "Artist A", "Album 2", "Rock", "1B", 80, []string{"rock"})

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, similarTrack},
		HasMore: false,
	}, nil)

	opts := DefaultSimilarityOptions()
	opts.Mode = "semantic"
	opts.MinSimilarity = 0.3

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Similar track should match due to same artist and genre
	assert.NotEmpty(t, result.Similar)
}

func TestFindSimilarTracks_FeaturesMode(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, nil)
	similarTrack := createSimilarityTestTrack("track-2", "Artist B", "Album 2", "Jazz", "8A", 122, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, similarTrack},
		HasMore: false,
	}, nil)

	opts := DefaultSimilarityOptions()
	opts.Mode = "features"
	opts.MinSimilarity = 0.3

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Similar track should match due to similar BPM and harmonic key
	assert.NotEmpty(t, result.Similar)
}

func TestFindMixableTracks_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)
	mixableTrack := createSimilarityTestTrack("track-2", "Artist B", "Album 2", "House", "8A", 130, nil)
	incompatibleTrack := createSimilarityTestTrack("track-3", "Artist C", "Album 3", "Ambient", "1B", 70, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, mixableTrack, incompatibleTrack},
		HasMore: false,
	}, nil)

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, DefaultMixingOptions())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, sourceTrack.ID, result.SourceTrack.ID)
	// Should include mixable track (close BPM, same key)
	found := false
	for _, mixable := range result.Mixable {
		if mixable.Track.ID == mixableTrack.ID {
			found = true
			assert.LessOrEqual(t, mixable.BPMDiff, 5)
		}
	}
	assert.True(t, found, "mixable track should be in results")
}

func TestFindMixableTracks_SourceTrackNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "nonexistent"
	mockRepo := new(MockSimilarityRepository)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(nil, errors.New("track not found"))

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, DefaultMixingOptions())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get source track")
}

func TestFindMixableTracks_ExactKeyMode(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)
	sameKeyTrack := createSimilarityTestTrack("track-2", "Artist B", "Album 2", "House", "8A", 130, nil)
	harmonicKeyTrack := createSimilarityTestTrack("track-3", "Artist C", "Album 3", "House", "7A", 128, nil) // Harmonic but not exact

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, sameKeyTrack, harmonicKeyTrack},
		HasMore: false,
	}, nil)

	opts := DefaultMixingOptions()
	opts.KeyMode = "exact"

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Only exact key match should be included
	for _, mixable := range result.Mixable {
		assert.Equal(t, "8A", mixable.Track.KeyCamelot)
	}
}

func TestFindMixableTracks_AnyKeyMode(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)
	differentKeyTrack := createSimilarityTestTrack("track-2", "Artist B", "Album 2", "House", "1B", 130, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, differentKeyTrack},
		HasMore: false,
	}, nil)

	opts := DefaultMixingOptions()
	opts.KeyMode = "any"

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Any key should be accepted
	assert.NotEmpty(t, result.Mixable)
}

func TestFindMixableTracks_BPMTolerance(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)
	closeTrack := createSimilarityTestTrack("track-2", "Artist B", "Album 2", "House", "8A", 130, nil)    // 2 BPM diff
	farTrack := createSimilarityTestTrack("track-3", "Artist C", "Album 3", "House", "8A", 140, nil)      // 12 BPM diff
	halfTimeTrack := createSimilarityTestTrack("track-4", "Artist D", "Album 4", "House", "8A", 64, nil)  // Half time

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, closeTrack, farTrack, halfTimeTrack},
		HasMore: false,
	}, nil)

	opts := DefaultMixingOptions()
	opts.BPMTolerance = 5

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Close track should be included, far track excluded
	foundClose := false
	foundFar := false
	for _, mixable := range result.Mixable {
		if mixable.Track.ID == closeTrack.ID {
			foundClose = true
		}
		if mixable.Track.ID == farTrack.ID {
			foundFar = true
		}
	}
	assert.True(t, foundClose, "close BPM track should be included")
	assert.False(t, foundFar, "far BPM track should be excluded")
}

func TestFindSimilarTracks_Limit(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, nil)

	// Create many similar tracks
	tracks := []models.Track{sourceTrack}
	for i := 2; i <= 20; i++ {
		tracks = append(tracks, createSimilarityTestTrack(
			"track-"+string(rune('0'+i)),
			"Artist A", // Same artist for high similarity
			"Album "+string(rune('0'+i)),
			"Rock",
			"8A",
			120+i,
			nil,
		))
	}

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   tracks,
		HasMore: false,
	}, nil)

	opts := DefaultSimilarityOptions()
	opts.Limit = 5
	opts.MinSimilarity = 0.1

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.LessOrEqual(t, len(result.Similar), 5)
}

func TestFindMixableTracks_Limit(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)

	// Create many mixable tracks
	tracks := []models.Track{sourceTrack}
	for i := 2; i <= 20; i++ {
		tracks = append(tracks, createSimilarityTestTrack(
			"track-"+string(rune('0'+i)),
			"Artist "+string(rune('A'+i)),
			"Album "+string(rune('0'+i)),
			"House",
			"8A",
			128+i%5, // Close BPM
			nil,
		))
	}

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   tracks,
		HasMore: false,
	}, nil)

	opts := DefaultMixingOptions()
	opts.Limit = 3

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.LessOrEqual(t, len(result.Mixable), 3)
}

func TestFindSimilarTracks_ListTracksError(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(nil, errors.New("database error"))

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, DefaultSimilarityOptions())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user tracks")
}

func TestFindMixableTracks_ListTracksError(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(nil, errors.New("database error"))

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, DefaultMixingOptions())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user tracks")
}

func TestSimilarTrack_KeyCompatibleFlag(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, nil)
	harmonicTrack := createSimilarityTestTrack("track-2", "Artist A", "Album 2", "Rock", "7A", 122, nil) // Harmonic
	incompatibleTrack := createSimilarityTestTrack("track-3", "Artist A", "Album 3", "Rock", "1B", 120, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack, harmonicTrack, incompatibleTrack},
		HasMore: false,
	}, nil)

	opts := DefaultSimilarityOptions()
	opts.MinSimilarity = 0.1

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)

	for _, similar := range result.Similar {
		if similar.Track.ID == harmonicTrack.ID {
			assert.True(t, similar.KeyCompatible, "harmonic key should be compatible")
		}
		if similar.Track.ID == incompatibleTrack.ID {
			assert.False(t, similar.KeyCompatible, "non-harmonic key should not be compatible")
		}
	}
}

func TestFindSimilarTracks_DefaultOptions(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "Rock", "8A", 120, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack},
		HasMore: false,
	}, nil)

	// Pass zero-value options to test defaults
	opts := SimilarityOptions{}

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindSimilarTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestFindMixableTracks_DefaultOptions(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	trackID := "track-1"
	mockRepo := new(MockSimilarityRepository)

	sourceTrack := createSimilarityTestTrack("track-1", "Artist A", "Album 1", "House", "8A", 128, nil)

	mockRepo.On("GetTrack", ctx, userID, trackID).Return(&sourceTrack, nil)
	mockRepo.On("ListTracks", ctx, userID, mock.AnythingOfType("models.TrackFilter")).Return(&repository.PaginatedResult[models.Track]{
		Items:   []models.Track{sourceTrack},
		HasMore: false,
	}, nil)

	// Pass zero-value options to test defaults
	opts := MixingOptions{}

	svc := NewSimilarityService(nil, mockRepo, nil)
	result, err := svc.FindMixableTracks(ctx, userID, trackID, opts)

	require.NoError(t, err)
	require.NotNil(t, result)
}
