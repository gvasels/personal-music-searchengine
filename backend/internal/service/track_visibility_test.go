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

// MockTrackVisibilityRepository mocks the repository for track visibility tests.
type MockTrackVisibilityRepository struct {
	mock.Mock
}

func (m *MockTrackVisibilityRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

func (m *MockTrackVisibilityRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	args := m.Called(ctx, limit, cursor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[models.Track]), args.Error(1)
}

func (m *MockTrackVisibilityRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

// MockS3RepoForVisibility mocks S3 repository for visibility tests.
type MockS3RepoForVisibility struct {
	mock.Mock
}

func (m *MockS3RepoForVisibility) GeneratePresignedDownloadURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	args := m.Called(ctx, key, duration)
	return args.String(0), args.Error(1)
}

// MockRoleServiceForVisibility mocks the role service.
type MockRoleServiceForVisibility struct {
	mock.Mock
}

func (m *MockRoleServiceForVisibility) GetUserRole(ctx context.Context, userID string) (models.UserRole, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.UserRole), args.Error(1)
}

func (m *MockRoleServiceForVisibility) HasPermission(ctx context.Context, userID string, permission models.Permission) (bool, error) {
	args := m.Called(ctx, userID, permission)
	return args.Bool(0), args.Error(1)
}

func TestListTracksWithVisibility_AdminSeesAll(t *testing.T) {
	t.Run("admin user sees all tracks from all users", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		adminID := "admin-123"
		now := time.Now()

		// All tracks from all users
		allTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: "user-1", Title: "Track 1", Visibility: models.VisibilityPrivate, Timestamps: models.Timestamps{CreatedAt: now}},
				{ID: "track-2", UserID: "user-2", Title: "Track 2", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
				{ID: "track-3", UserID: "admin-123", Title: "Track 3", Visibility: models.VisibilityPrivate, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		mockRole.On("GetUserRole", ctx, adminID).Return(models.RoleAdmin, nil)
		mockRepo.On("ListTracks", ctx, "", mock.MatchedBy(func(f models.TrackFilter) bool {
			return f.GlobalScope == true
		})).Return(allTracks, nil)
		mockRepo.On("GetUserDisplayName", ctx, "user-1").Return("User One", nil)
		mockRepo.On("GetUserDisplayName", ctx, "user-2").Return("User Two", nil)
		mockRepo.On("GetUserDisplayName", ctx, "admin-123").Return("Admin", nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		result, err := svc.ListTracksWithVisibility(ctx, adminID, models.TrackFilter{})

		require.NoError(t, err)
		assert.Len(t, result.Items, 3)
		// Admin should see all tracks with owner display names
		mockRepo.AssertExpectations(t)
	})
}

func TestListTracksWithVisibility_GlobalReaderSeesAll(t *testing.T) {
	t.Run("global reader sees all tracks from all users", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		globalReaderID := "global-reader-123"
		now := time.Now()

		allTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: "user-1", Title: "Track 1", Visibility: models.VisibilityPrivate, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		// Global reader has GlobalReaders permission
		mockRole.On("GetUserRole", ctx, globalReaderID).Return(models.RoleSubscriber, nil)
		mockRole.On("HasPermission", ctx, globalReaderID, models.PermissionViewGlobal).Return(true, nil)
		mockRepo.On("ListTracks", ctx, "", mock.MatchedBy(func(f models.TrackFilter) bool {
			return f.GlobalScope == true
		})).Return(allTracks, nil)
		mockRepo.On("GetUserDisplayName", ctx, "user-1").Return("User One", nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		result, err := svc.ListTracksWithVisibility(ctx, globalReaderID, models.TrackFilter{})

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
	})
}

func TestListTracksWithVisibility_RegularUserSeesOwnAndPublic(t *testing.T) {
	t.Run("regular user sees own tracks plus public tracks", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		userID := "user-123"
		now := time.Now()

		// User's own tracks
		ownTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: userID, Title: "My Private Track", Visibility: models.VisibilityPrivate, Timestamps: models.Timestamps{CreatedAt: now}},
				{ID: "track-2", UserID: userID, Title: "My Public Track", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		// Public tracks from other users
		publicTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-3", UserID: "other-user", Title: "Other Public Track", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		mockRole.On("GetUserRole", ctx, userID).Return(models.RoleSubscriber, nil)
		mockRole.On("HasPermission", ctx, userID, models.PermissionViewGlobal).Return(false, nil)
		mockRepo.On("ListTracks", ctx, userID, mock.Anything).Return(ownTracks, nil)
		mockRepo.On("ListPublicTracks", ctx, mock.Anything, mock.Anything).Return(publicTracks, nil)
		mockRepo.On("GetUserDisplayName", ctx, "other-user").Return("Other User", nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		filter := models.TrackFilter{IncludePublic: true}
		result, err := svc.ListTracksWithVisibility(ctx, userID, filter)

		require.NoError(t, err)
		// Should see own tracks (2) + public tracks from others (1)
		assert.Len(t, result.Items, 3)
	})

	t.Run("regular user without IncludePublic sees only own tracks", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		userID := "user-123"
		now := time.Now()

		ownTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: userID, Title: "My Track", Visibility: models.VisibilityPrivate, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		mockRole.On("GetUserRole", ctx, userID).Return(models.RoleSubscriber, nil)
		mockRole.On("HasPermission", ctx, userID, models.PermissionViewGlobal).Return(false, nil)
		mockRepo.On("ListTracks", ctx, userID, mock.Anything).Return(ownTracks, nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		result, err := svc.ListTracksWithVisibility(ctx, userID, models.TrackFilter{})

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		// Should NOT call ListPublicTracks
		mockRepo.AssertNotCalled(t, "ListPublicTracks")
	})
}

func TestListTracksWithVisibility_DeduplicatesPublicTracks(t *testing.T) {
	t.Run("deduplicates when user's public track appears in both queries", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		userID := "user-123"
		now := time.Now()

		// User's own public track
		ownTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: userID, Title: "My Public Track", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		// Public tracks includes the same track
		publicTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: userID, Title: "My Public Track", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
				{ID: "track-2", UserID: "other-user", Title: "Other Public Track", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		mockRole.On("GetUserRole", ctx, userID).Return(models.RoleSubscriber, nil)
		mockRole.On("HasPermission", ctx, userID, models.PermissionViewGlobal).Return(false, nil)
		mockRepo.On("ListTracks", ctx, userID, mock.Anything).Return(ownTracks, nil)
		mockRepo.On("ListPublicTracks", ctx, mock.Anything, mock.Anything).Return(publicTracks, nil)
		mockRepo.On("GetUserDisplayName", ctx, "other-user").Return("Other User", nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		filter := models.TrackFilter{IncludePublic: true}
		result, err := svc.ListTracksWithVisibility(ctx, userID, filter)

		require.NoError(t, err)
		// Should have 2 unique tracks (track-1 deduplicated)
		assert.Len(t, result.Items, 2)
	})
}

func TestListTracksWithVisibility_OwnerDisplayName(t *testing.T) {
	t.Run("sets OwnerDisplayName to 'You' for own tracks", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		userID := "user-123"
		now := time.Now()

		ownTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: userID, Title: "My Track", Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		publicTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-2", UserID: "other-user", Title: "Other Track", Visibility: models.VisibilityPublic, Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		mockRole.On("GetUserRole", ctx, userID).Return(models.RoleSubscriber, nil)
		mockRole.On("HasPermission", ctx, userID, models.PermissionViewGlobal).Return(false, nil)
		mockRepo.On("ListTracks", ctx, userID, mock.Anything).Return(ownTracks, nil)
		mockRepo.On("ListPublicTracks", ctx, mock.Anything, mock.Anything).Return(publicTracks, nil)
		mockRepo.On("GetUserDisplayName", ctx, "other-user").Return("Other User", nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		filter := models.TrackFilter{IncludePublic: true}
		result, err := svc.ListTracksWithVisibility(ctx, userID, filter)

		require.NoError(t, err)
		assert.Len(t, result.Items, 2)

		// Find own track and verify OwnerDisplayName
		for _, track := range result.Items {
			if track.ID == "track-1" {
				assert.Equal(t, "You", track.OwnerDisplayName)
			} else if track.ID == "track-2" {
				assert.Equal(t, "Other User", track.OwnerDisplayName)
			}
		}
	})
}

func TestListTracksWithVisibility_DefaultsPrivateVisibility(t *testing.T) {
	t.Run("defaults visibility to private for tracks without visibility set", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockTrackVisibilityRepository)
		mockRole := new(MockRoleServiceForVisibility)

		userID := "user-123"
		now := time.Now()

		// Track without visibility set
		ownTracks := &repository.PaginatedResult[models.Track]{
			Items: []models.Track{
				{ID: "track-1", UserID: userID, Title: "Old Track", Visibility: "", Timestamps: models.Timestamps{CreatedAt: now}},
			},
			HasMore: false,
		}

		mockRole.On("GetUserRole", ctx, userID).Return(models.RoleSubscriber, nil)
		mockRole.On("HasPermission", ctx, userID, models.PermissionViewGlobal).Return(false, nil)
		mockRepo.On("ListTracks", ctx, userID, mock.Anything).Return(ownTracks, nil)

		svc := NewTrackVisibilityService(mockRepo, mockRole)
		result, err := svc.ListTracksWithVisibility(ctx, userID, models.TrackFilter{})

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		// Visibility should default to private in the response
		assert.Equal(t, "private", result.Items[0].Visibility)
	})
}

// =============================================================================
// TrackService.GetTrack Visibility Tests (Access Control Bug Fixes Task 1.4)
// =============================================================================

// MockTrackServiceRepository provides a mock for testing trackService.GetTrack
type MockTrackServiceRepository struct {
	mock.Mock
}

func (m *MockTrackServiceRepository) GetTrack(ctx context.Context, userID, trackID string) (*models.Track, error) {
	args := m.Called(ctx, userID, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

func (m *MockTrackServiceRepository) GetTrackByID(ctx context.Context, trackID string) (*models.Track, error) {
	args := m.Called(ctx, trackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Track), args.Error(1)
}

// Required Repository interface stubs
func (m *MockTrackServiceRepository) UpdateTrack(ctx context.Context, track models.Track) error {
	args := m.Called(ctx, track)
	return args.Error(0)
}

func (m *MockTrackServiceRepository) DeleteTrack(ctx context.Context, userID, trackID string) error {
	return nil
}

func (m *MockTrackServiceRepository) CreateTrack(ctx context.Context, track models.Track) error {
	return nil
}

func (m *MockTrackServiceRepository) ListTracks(ctx context.Context, userID string, filter models.TrackFilter) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}

func (m *MockTrackServiceRepository) ListTracksByArtist(ctx context.Context, userID, artist string) ([]models.Track, error) {
	return nil, nil
}

func (m *MockTrackServiceRepository) ListPublicTracks(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Track], error) {
	return nil, nil
}

func (m *MockTrackServiceRepository) UpdateTrackVisibility(ctx context.Context, userID, trackID string, visibility models.TrackVisibility) error {
	return nil
}

// Album stubs
func (m *MockTrackServiceRepository) GetOrCreateAlbum(ctx context.Context, userID, albumName, artist string) (*models.Album, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) GetAlbum(ctx context.Context, userID, albumID string) (*models.Album, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ListAlbums(ctx context.Context, userID string, filter models.AlbumFilter) (*repository.PaginatedResult[models.Album], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ListAlbumsByArtist(ctx context.Context, userID, artist string) ([]models.Album, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateAlbumStats(ctx context.Context, userID, albumID string, trackCount, totalDuration int) error {
	return nil
}

// User stubs
func (m *MockTrackServiceRepository) CreateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockTrackServiceRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateUser(ctx context.Context, user models.User) error { return nil }
func (m *MockTrackServiceRepository) UpdateUserStats(ctx context.Context, userID string, storageUsed int64, trackCount, albumCount, playlistCount int) error {
	return nil
}

// Playlist stubs
func (m *MockTrackServiceRepository) CreatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockTrackServiceRepository) GetPlaylist(ctx context.Context, userID, playlistID string) (*models.Playlist, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdatePlaylist(ctx context.Context, playlist models.Playlist) error {
	return nil
}
func (m *MockTrackServiceRepository) DeletePlaylist(ctx context.Context, userID, playlistID string) error {
	return nil
}
func (m *MockTrackServiceRepository) ListPlaylists(ctx context.Context, userID string, filter models.PlaylistFilter) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) AddTracksToPlaylist(ctx context.Context, playlistID string, trackIDs []string, position int) error {
	return nil
}
func (m *MockTrackServiceRepository) RemoveTracksFromPlaylist(ctx context.Context, playlistID string, trackIDs []string) error {
	return nil
}
func (m *MockTrackServiceRepository) GetPlaylistTracks(ctx context.Context, playlistID string) ([]models.PlaylistTrack, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ReorderPlaylistTracks(ctx context.Context, playlistID string, tracks []models.PlaylistTrack) error {
	return nil
}

// Tag stubs
func (m *MockTrackServiceRepository) CreateTag(ctx context.Context, tag models.Tag) error { return nil }
func (m *MockTrackServiceRepository) GetTag(ctx context.Context, userID, tagName string) (*models.Tag, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateTag(ctx context.Context, tag models.Tag) error { return nil }
func (m *MockTrackServiceRepository) DeleteTag(ctx context.Context, userID, tagName string) error {
	return nil
}
func (m *MockTrackServiceRepository) ListTags(ctx context.Context, userID string) ([]models.Tag, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) AddTagsToTrack(ctx context.Context, userID, trackID string, tagNames []string) error {
	return nil
}
func (m *MockTrackServiceRepository) RemoveTagFromTrack(ctx context.Context, userID, trackID, tagName string) error {
	return nil
}
func (m *MockTrackServiceRepository) GetTracksByTag(ctx context.Context, userID, tagName string) ([]models.Track, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) GetTrackTags(ctx context.Context, userID, trackID string) ([]string, error) {
	return nil, nil
}

// Artist stubs
func (m *MockTrackServiceRepository) CreateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockTrackServiceRepository) GetArtist(ctx context.Context, userID, artistID string) (*models.Artist, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) GetArtistByName(ctx context.Context, userID, name string) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ListArtists(ctx context.Context, userID string, filter models.ArtistFilter) (*repository.PaginatedResult[models.Artist], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateArtist(ctx context.Context, artist models.Artist) error {
	return nil
}
func (m *MockTrackServiceRepository) DeleteArtist(ctx context.Context, userID, artistID string) error {
	return nil
}
func (m *MockTrackServiceRepository) BatchGetArtists(ctx context.Context, userID string, artistIDs []string) (map[string]*models.Artist, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) SearchArtists(ctx context.Context, userID, query string, limit int) ([]*models.Artist, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) GetArtistTrackCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockTrackServiceRepository) GetArtistAlbumCount(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}
func (m *MockTrackServiceRepository) GetArtistTotalPlays(ctx context.Context, userID, artistID string) (int, error) {
	return 0, nil
}

// User extended stubs
func (m *MockTrackServiceRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (*models.User, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return nil
}
func (m *MockTrackServiceRepository) ListUsersByRole(ctx context.Context, role models.UserRole, limit int, cursor string) (*repository.PaginatedResult[models.User], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) SearchUsersByEmail(ctx context.Context, emailPrefix string, limit int, cursor string) ([]repository.UserSearchResult, string, error) {
	return nil, "", nil
}
func (m *MockTrackServiceRepository) SetUserDisabled(ctx context.Context, userID string, disabled bool) error {
	return nil
}
func (m *MockTrackServiceRepository) GetUserDisplayName(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockTrackServiceRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
func (m *MockTrackServiceRepository) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateUserSettings(ctx context.Context, userID string, update *repository.UserSettingsUpdate) (*models.UserSettings, error) {
	return nil, nil
}

// Playlist extended stubs
func (m *MockTrackServiceRepository) SearchPlaylists(ctx context.Context, userID, query string, limit int) ([]models.Playlist, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdatePlaylistVisibility(ctx context.Context, userID, playlistID string, visibility models.PlaylistVisibility) error {
	return nil
}
func (m *MockTrackServiceRepository) ListPublicPlaylists(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.Playlist], error) {
	return nil, nil
}

// ArtistProfile stubs
func (m *MockTrackServiceRepository) CreateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockTrackServiceRepository) GetArtistProfile(ctx context.Context, userID string) (*models.ArtistProfile, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateArtistProfile(ctx context.Context, profile models.ArtistProfile) error {
	return nil
}
func (m *MockTrackServiceRepository) DeleteArtistProfile(ctx context.Context, userID string) error {
	return nil
}
func (m *MockTrackServiceRepository) ListArtistProfiles(ctx context.Context, limit int, cursor string) (*repository.PaginatedResult[models.ArtistProfile], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) IncrementArtistFollowerCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Follow stubs
func (m *MockTrackServiceRepository) CreateFollow(ctx context.Context, follow models.Follow) error {
	return nil
}
func (m *MockTrackServiceRepository) DeleteFollow(ctx context.Context, followerID, followedID string) error {
	return nil
}
func (m *MockTrackServiceRepository) GetFollow(ctx context.Context, followerID, followedID string) (*models.Follow, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ListFollowers(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ListFollowing(ctx context.Context, userID string, limit int, cursor string) (*repository.PaginatedResult[models.Follow], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) IncrementUserFollowingCount(ctx context.Context, userID string, delta int) error {
	return nil
}

// Upload stubs
func (m *MockTrackServiceRepository) CreateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockTrackServiceRepository) GetUpload(ctx context.Context, userID, uploadID string) (*models.Upload, error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) UpdateUpload(ctx context.Context, upload models.Upload) error {
	return nil
}
func (m *MockTrackServiceRepository) UpdateUploadStatus(ctx context.Context, userID, uploadID string, status models.UploadStatus, errorMsg string, trackID string) error {
	return nil
}
func (m *MockTrackServiceRepository) UpdateUploadStep(ctx context.Context, userID, uploadID string, step models.ProcessingStep, success bool) error {
	return nil
}
func (m *MockTrackServiceRepository) ListUploads(ctx context.Context, userID string, filter models.UploadFilter) (*repository.PaginatedResult[models.Upload], error) {
	return nil, nil
}
func (m *MockTrackServiceRepository) ListUploadsByStatus(ctx context.Context, status models.UploadStatus) ([]models.Upload, error) {
	return nil, nil
}

// S3 extended stubs
func (m *MockS3RepoForTrackService) GeneratePresignedDownloadURLWithFilename(ctx context.Context, key string, expiry time.Duration, filename string) (string, error) {
	return "", nil
}
func (m *MockS3RepoForTrackService) InitiateMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	return "", nil
}
func (m *MockS3RepoForTrackService) GenerateMultipartUploadURLs(ctx context.Context, key, uploadID string, numParts int, expiry time.Duration) ([]models.MultipartUploadPartURL, error) {
	return nil, nil
}
func (m *MockS3RepoForTrackService) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []models.CompletedPartInfo) error {
	return nil
}
func (m *MockS3RepoForTrackService) AbortMultipartUpload(ctx context.Context, key, uploadID string) error {
	return nil
}
func (m *MockS3RepoForTrackService) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	return nil
}
func (m *MockS3RepoForTrackService) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}
func (m *MockS3RepoForTrackService) ObjectExists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// MockS3RepoForTrackService mocks S3 repository for track service tests.
type MockS3RepoForTrackService struct {
	mock.Mock
}

func (m *MockS3RepoForTrackService) GeneratePresignedDownloadURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	args := m.Called(ctx, key, duration)
	return args.String(0), args.Error(1)
}

func (m *MockS3RepoForTrackService) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, duration time.Duration) (string, error) {
	return "", nil
}

func (m *MockS3RepoForTrackService) DeleteObject(ctx context.Context, key string) error {
	return nil
}

func TestTrackService_GetTrack_OwnerCanAccessPrivate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTrackServiceRepository)
	mockS3 := new(MockS3RepoForTrackService)

	ownerID := "user-123"
	trackID := "track-456"
	now := time.Now()

	privateTrack := &models.Track{
		ID:         trackID,
		UserID:     ownerID,
		Title:      "My Private Track",
		Visibility: models.VisibilityPrivate,
		Timestamps: models.Timestamps{CreatedAt: now},
	}

	// Owner requests their own private track
	mockRepo.On("GetTrack", ctx, ownerID, trackID).Return(privateTrack, nil)

	svc := NewTrackService(mockRepo, mockS3)
	result, err := svc.GetTrack(ctx, ownerID, trackID, false) // hasGlobal=false

	require.NoError(t, err)
	assert.Equal(t, trackID, result.ID)
	assert.Equal(t, "My Private Track", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestTrackService_GetTrack_AdminCanAccessAnyTrack(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTrackServiceRepository)
	mockS3 := new(MockS3RepoForTrackService)

	adminID := "admin-123"
	ownerID := "user-456"
	trackID := "track-789"
	now := time.Now()

	privateTrack := &models.Track{
		ID:         trackID,
		UserID:     ownerID,
		Title:      "Other User Private Track",
		Visibility: models.VisibilityPrivate,
		Timestamps: models.Timestamps{CreatedAt: now},
	}

	// Admin requests someone else's private track
	mockRepo.On("GetTrack", ctx, adminID, trackID).Return(nil, repository.ErrNotFound)
	mockRepo.On("GetTrackByID", ctx, trackID).Return(privateTrack, nil)

	svc := NewTrackService(mockRepo, mockS3)
	result, err := svc.GetTrack(ctx, adminID, trackID, true) // hasGlobal=true (admin)

	require.NoError(t, err)
	assert.Equal(t, trackID, result.ID)
	assert.Equal(t, "Other User Private Track", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestTrackService_GetTrack_NonOwnerForbiddenForPrivate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTrackServiceRepository)
	mockS3 := new(MockS3RepoForTrackService)

	requesterID := "subscriber-123"
	ownerID := "owner-456"
	trackID := "track-789"
	now := time.Now()

	privateTrack := &models.Track{
		ID:         trackID,
		UserID:     ownerID,
		Title:      "Private Track",
		Visibility: models.VisibilityPrivate,
		Timestamps: models.Timestamps{CreatedAt: now},
	}

	// Non-owner requests someone else's private track
	mockRepo.On("GetTrack", ctx, requesterID, trackID).Return(nil, repository.ErrNotFound)
	mockRepo.On("GetTrackByID", ctx, trackID).Return(privateTrack, nil)

	svc := NewTrackService(mockRepo, mockS3)
	result, err := svc.GetTrack(ctx, requesterID, trackID, false) // hasGlobal=false

	require.Error(t, err)
	assert.Nil(t, result)

	// Verify it's a 403 Forbidden error
	var apiErr *models.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 403, apiErr.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestTrackService_GetTrack_NonOwnerCanAccessPublic(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTrackServiceRepository)
	mockS3 := new(MockS3RepoForTrackService)

	requesterID := "subscriber-123"
	ownerID := "owner-456"
	trackID := "track-789"
	now := time.Now()

	publicTrack := &models.Track{
		ID:         trackID,
		UserID:     ownerID,
		Title:      "Public Track",
		Visibility: models.VisibilityPublic,
		Timestamps: models.Timestamps{CreatedAt: now},
	}

	// Non-owner requests someone else's public track
	mockRepo.On("GetTrack", ctx, requesterID, trackID).Return(nil, repository.ErrNotFound)
	mockRepo.On("GetTrackByID", ctx, trackID).Return(publicTrack, nil)

	svc := NewTrackService(mockRepo, mockS3)
	result, err := svc.GetTrack(ctx, requesterID, trackID, false) // hasGlobal=false

	require.NoError(t, err)
	assert.Equal(t, trackID, result.ID)
	assert.Equal(t, "Public Track", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestTrackService_GetTrack_NonOwnerCanAccessUnlisted(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTrackServiceRepository)
	mockS3 := new(MockS3RepoForTrackService)

	requesterID := "subscriber-123"
	ownerID := "owner-456"
	trackID := "track-789"
	now := time.Now()

	unlistedTrack := &models.Track{
		ID:         trackID,
		UserID:     ownerID,
		Title:      "Unlisted Track",
		Visibility: models.VisibilityUnlisted,
		Timestamps: models.Timestamps{CreatedAt: now},
	}

	// Non-owner requests someone else's unlisted track via direct link
	mockRepo.On("GetTrack", ctx, requesterID, trackID).Return(nil, repository.ErrNotFound)
	mockRepo.On("GetTrackByID", ctx, trackID).Return(unlistedTrack, nil)

	svc := NewTrackService(mockRepo, mockS3)
	result, err := svc.GetTrack(ctx, requesterID, trackID, false) // hasGlobal=false

	require.NoError(t, err)
	assert.Equal(t, trackID, result.ID)
	assert.Equal(t, "Unlisted Track", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestTrackService_GetTrack_NotFoundReturns404(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTrackServiceRepository)
	mockS3 := new(MockS3RepoForTrackService)

	requesterID := "user-123"
	trackID := "nonexistent-track"

	// Track doesn't exist at all
	mockRepo.On("GetTrack", ctx, requesterID, trackID).Return(nil, repository.ErrNotFound)
	mockRepo.On("GetTrackByID", ctx, trackID).Return(nil, repository.ErrNotFound)

	svc := NewTrackService(mockRepo, mockS3)
	result, err := svc.GetTrack(ctx, requesterID, trackID, false)

	require.Error(t, err)
	assert.Nil(t, result)

	// Verify it's a 404 Not Found error
	var apiErr *models.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 404, apiErr.StatusCode)
	mockRepo.AssertExpectations(t)
}
