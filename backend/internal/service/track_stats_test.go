package service

import (
	"context"
	"testing"
	"time"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStatsRepository mocks the repository for stats tests.
type MockStatsRepository struct {
	mock.Mock
}

// Track operations - these are the ones we actually test
func (m *MockStatsRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

func (m *MockStatsRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

func (m *MockStatsRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockStatsRepository) GetTrackByID(ctx context.Context, trackID string) (*models.Track, error) {
	args := m.Called(ctx, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockStatsRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockStatsRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	args := m.Called(ctx, userID, trackID, visibility)
	return args.Error(0)
}

func (m *MockStatsRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	args := m.Called(ctx, userID, artist)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Track), args.Error(1)
}

// Stub implementations for remaining Repository interface methods
func (m *MockStatsRepository) CreateTrack(ctx context.Context, track models.Track) error {
	return nil
}
func (m *MockStatsRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return nil
}
func (m *MockStatsRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockStatsRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockStatsRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockStatsRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}
func (m *MockStatsRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockStatsRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockStatsRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockStatsRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockStatsRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockStatsRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockStatsRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockStatsRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockStatsRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockStatsRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockStatsRepository) CreateUser(ctx context.Context, user models.User) error {
	return nil
}
func (m *MockStatsRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockStatsRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockStatsRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateUser(ctx context.Context, user models.User) error {
	return nil
}
func (m *MockStatsRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}
func (m *MockStatsRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockStatsRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}
func (m *MockStatsRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockStatsRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]repository.UserSearchResult, string, error) {
	return nil, "", nil
}
func (m *MockStatsRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockStatsRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockStatsRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
func (m *MockStatsRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateUserSettings(ctx context.Context, userID string, update *repository.UserSettingsUpdate) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockStatsRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockStatsRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockStatsRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return nil
}
func (m *MockStatsRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockStatsRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}
func (m *MockStatsRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	return nil
}
func (m *MockStatsRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	return nil
}
func (m *MockStatsRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockStatsRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	return nil
}
func (m *MockStatsRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	return nil
}
func (m *MockStatsRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockStatsRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockStatsRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockStatsRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockStatsRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockStatsRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}
func (m *MockStatsRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockStatsRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockStatsRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockStatsRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockStatsRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockStatsRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}
func (m *MockStatsRepository) CreateTag(ctx context.Context, tag models.Tag) error {
	return nil
}
func (m *MockStatsRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateTag(ctx context.Context, tag models.Tag) error {
	return nil
}
func (m *MockStatsRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	return nil
}
func (m *MockStatsRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	return nil, nil
}
func (m *MockStatsRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	return nil
}
func (m *MockStatsRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	return nil
}
func (m *MockStatsRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}
func (m *MockStatsRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockStatsRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockStatsRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockStatsRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockStatsRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockStatsRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockStatsRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockStatsRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// MockS3RepoForStats mocks S3 repository for stats tests.
type MockS3RepoForStats struct {
	mock.Mock
}

func (m *MockS3RepoForStats) GeneratePresignedDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockS3RepoForStats) DeleteObject(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Stub implementations for remaining S3Repository interface methods
func (m *MockS3RepoForStats) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	return "", nil
}
func (m *MockS3RepoForStats) GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, expiry time.Duration, filename string) (string, error) {
	return "", nil
}
func (m *MockS3RepoForStats) InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	return "", nil
}
func (m *MockS3RepoForStats) GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error) {
	return nil, nil
}
func (m *MockS3RepoForStats) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error {
	return nil
}
func (m *MockS3RepoForStats) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return nil
}
func (m *MockS3RepoForStats) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	return nil
}
func (m *MockS3RepoForStats) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}
func (m *MockS3RepoForStats) ObjectExists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// createStatsTestService creates a track service with mocked dependencies for testing
func createStatsTestService(mockRepo *MockStatsRepository, mockS3 *MockS3RepoForStats) *trackService {
	return &trackService{
		repo:   mockRepo,
		s3Repo: mockS3,
	}
}

func TestGetLibraryStats_AdminScopeAll(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	// Mock all tracks (admin global view)
	allTracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track1", UserID: "user1", Title: "Song 1", Artist: "Artist A", Album: "Album 1", Duration: 180},
			{ID: "track2", UserID: "user1", Title: "Song 2", Artist: "Artist A", Album: "Album 1", Duration: 200},
			{ID: "track3", UserID: "user2", Title: "Song 3", Artist: "Artist B", Album: "Album 2", Duration: 240},
			{ID: "track4", UserID: "user3", Title: "Song 4", Artist: "Artist C", Album: "Album 3", Duration: 300},
		},
		HasMore: false,
	}

	mockRepo.On("ListTracks", ctx, "admin-user", mock.MatchedBy(func(f models.TrackFilter) bool {
		return f.GlobalScope == true && f.Limit == 10000
	})).Return(allTracks, nil)

	svc := createStatsTestService(mockRepo, mockS3)

	stats, err := svc.GetLibraryStats(ctx, "admin-user", StatsScopeAll, true)

	require.NoError(t, err)
	assert.Equal(t, 4, stats.TotalTracks)
	assert.Equal(t, 3, stats.TotalAlbums)   // Album 1, Album 2, Album 3
	assert.Equal(t, 3, stats.TotalArtists)  // Artist A, Artist B, Artist C
	assert.Equal(t, 920, stats.TotalDuration) // 180+200+240+300

	mockRepo.AssertExpectations(t)
}

func TestGetLibraryStats_AdminScopeAll_RequiresGlobalAccess(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	svc := createStatsTestService(mockRepo, mockS3)

	// Non-admin trying to use scope=all should get forbidden error
	stats, err := svc.GetLibraryStats(ctx, "regular-user", StatsScopeAll, false)

	require.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "admin access required")
}

func TestGetLibraryStats_ScopePublic(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	// Mock public tracks only
	publicTracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track1", UserID: "user1", Title: "Public Song 1", Artist: "Artist A", Album: "Album 1", Duration: 180, Visibility: models.VisibilityPublic},
			{ID: "track2", UserID: "user2", Title: "Public Song 2", Artist: "Artist B", Album: "Album 2", Duration: 220, Visibility: models.VisibilityPublic},
		},
		HasMore: false,
	}

	mockRepo.On("ListPublicTracks", ctx, 10000, "").Return(publicTracks, nil)

	svc := createStatsTestService(mockRepo, mockS3)

	// Regular user using scope=public (subscriber simulation)
	stats, err := svc.GetLibraryStats(ctx, "any-user", StatsScopePublic, false)

	require.NoError(t, err)
	assert.Equal(t, 2, stats.TotalTracks)
	assert.Equal(t, 2, stats.TotalAlbums)
	assert.Equal(t, 2, stats.TotalArtists)
	assert.Equal(t, 400, stats.TotalDuration)

	mockRepo.AssertExpectations(t)
}

func TestGetLibraryStats_ScopeOwn(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	userID := "user123"

	// Mock user's own tracks
	ownTracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track1", UserID: userID, Title: "My Song 1", Artist: "Artist A", Album: "Album 1", Duration: 180},
			{ID: "track2", UserID: userID, Title: "My Song 2", Artist: "Artist A", Album: "Album 1", Duration: 200},
		},
		HasMore: false,
	}

	// Mock public tracks (some from other users)
	publicTracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track3", UserID: "other-user", Title: "Public Song", Artist: "Artist B", Album: "Album 2", Duration: 240, Visibility: models.VisibilityPublic},
		},
		HasMore: false,
	}

	mockRepo.On("ListTracks", ctx, userID, mock.MatchedBy(func(f models.TrackFilter) bool {
		return f.GlobalScope == false && f.Limit == 10000
	})).Return(ownTracks, nil)
	mockRepo.On("ListPublicTracks", ctx, 10000, "").Return(publicTracks, nil)

	svc := createStatsTestService(mockRepo, mockS3)

	stats, err := svc.GetLibraryStats(ctx, userID, StatsScopeOwn, false)

	require.NoError(t, err)
	assert.Equal(t, 3, stats.TotalTracks)  // 2 own + 1 public
	assert.Equal(t, 2, stats.TotalAlbums)  // Album 1, Album 2
	assert.Equal(t, 2, stats.TotalArtists) // Artist A, Artist B
	assert.Equal(t, 620, stats.TotalDuration)

	mockRepo.AssertExpectations(t)
}

func TestGetLibraryStats_ScopeOwn_DeduplicatesPublicTracks(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	userID := "user123"

	// User's own track that is also public
	ownTracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track1", UserID: userID, Title: "My Public Song", Artist: "Artist A", Album: "Album 1", Duration: 180, Visibility: models.VisibilityPublic},
		},
		HasMore: false,
	}

	// Same track appears in public list
	publicTracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track1", UserID: userID, Title: "My Public Song", Artist: "Artist A", Album: "Album 1", Duration: 180, Visibility: models.VisibilityPublic},
			{ID: "track2", UserID: "other-user", Title: "Other Public Song", Artist: "Artist B", Album: "Album 2", Duration: 200, Visibility: models.VisibilityPublic},
		},
		HasMore: false,
	}

	mockRepo.On("ListTracks", ctx, userID, mock.Anything).Return(ownTracks, nil)
	mockRepo.On("ListPublicTracks", ctx, 10000, "").Return(publicTracks, nil)

	svc := createStatsTestService(mockRepo, mockS3)

	stats, err := svc.GetLibraryStats(ctx, userID, StatsScopeOwn, false)

	require.NoError(t, err)
	// Should have 2 tracks (not 3) - track1 should not be counted twice
	assert.Equal(t, 2, stats.TotalTracks)

	mockRepo.AssertExpectations(t)
}

func TestGetLibraryStats_EmptyLibrary(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	emptyResult := &repository.PaginatedResult[models.Track]{
		Items:   []models.Track{},
		HasMore: false,
	}

	mockRepo.On("ListPublicTracks", ctx, 10000, "").Return(emptyResult, nil)

	svc := createStatsTestService(mockRepo, mockS3)

	stats, err := svc.GetLibraryStats(ctx, "user123", StatsScopePublic, false)

	require.NoError(t, err)
	assert.Equal(t, 0, stats.TotalTracks)
	assert.Equal(t, 0, stats.TotalAlbums)
	assert.Equal(t, 0, stats.TotalArtists)
	assert.Equal(t, 0, stats.TotalDuration)

	mockRepo.AssertExpectations(t)
}

func TestGetLibraryStats_TracksWithNoAlbumOrArtist(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockStatsRepository)
	mockS3 := new(MockS3RepoForStats)

	// Tracks with missing metadata
	tracks := &repository.PaginatedResult[models.Track]{
		Items: []models.Track{
			{ID: "track1", Title: "Song 1", Artist: "", Album: "", Duration: 180},
			{ID: "track2", Title: "Song 2", Artist: "Artist A", Album: "", Duration: 200},
			{ID: "track3", Title: "Song 3", Artist: "", Album: "Album 1", Duration: 220},
		},
		HasMore: false,
	}

	mockRepo.On("ListPublicTracks", ctx, 10000, "").Return(tracks, nil)

	svc := createStatsTestService(mockRepo, mockS3)

	stats, err := svc.GetLibraryStats(ctx, "user123", StatsScopePublic, false)

	require.NoError(t, err)
	assert.Equal(t, 3, stats.TotalTracks)
	assert.Equal(t, 1, stats.TotalAlbums)  // Only Album 1 (empty strings not counted)
	assert.Equal(t, 1, stats.TotalArtists) // Only Artist A (empty strings not counted)
	assert.Equal(t, 600, stats.TotalDuration)

	mockRepo.AssertExpectations(t)
}
