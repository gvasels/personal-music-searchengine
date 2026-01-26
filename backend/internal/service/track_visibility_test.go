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
