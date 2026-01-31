//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Track CRUD
// ---------------------------------------------------------------------------

func TestIntegration_TrackCRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	// Setup: create a user first
	user := models.User{
		ID:          "track-user",
		Email:       "trackuser@test.com",
		DisplayName: "Track User",
		Role:        models.RoleSubscriber,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#track-user", "PROFILE")

	track := models.Track{
		ID:         "track-001",
		UserID:     "track-user",
		Title:      "Test Song",
		Artist:     "Test Artist",
		Album:      "Test Album",
		Genre:      "Electronic",
		Duration:   180,
		Format:     models.AudioFormatMP3,
		FileSize:   5242880,
		S3Key:      "uploads/track-user/track-001.mp3",
		Visibility: models.VisibilityPrivate,
	}

	t.Run("create and get track", func(t *testing.T) {
		err := repo.CreateTrack(ctx, track)
		require.NoError(t, err)

		got, err := repo.GetTrack(ctx, "track-user", "track-001")
		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, "track-001", got.ID)
		assert.Equal(t, "track-user", got.UserID)
		assert.Equal(t, "Test Song", got.Title)
		assert.Equal(t, "Test Artist", got.Artist)
		assert.Equal(t, "Test Album", got.Album)
		assert.Equal(t, "Electronic", got.Genre)
		assert.Equal(t, 180, got.Duration)
		assert.Equal(t, models.AudioFormatMP3, got.Format)
		assert.Equal(t, int64(5242880), got.FileSize)
		assert.Equal(t, models.VisibilityPrivate, got.Visibility)
		assert.False(t, got.CreatedAt.IsZero())
		assert.False(t, got.UpdatedAt.IsZero())
	})

	t.Run("get non-existent track", func(t *testing.T) {
		_, err := repo.GetTrack(ctx, "track-user", "nonexistent")
		assert.Error(t, err) // ErrTrackNotFound
	})

	t.Run("GetTrackByID global lookup", func(t *testing.T) {
		got, err := repo.GetTrackByID(ctx, "track-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "track-001", got.ID)
		assert.Equal(t, "track-user", got.UserID)
	})

	t.Run("update track", func(t *testing.T) {
		track.Title = "Updated Song"
		track.Genre = "Ambient"
		track.Visibility = models.VisibilityPublic
		err := repo.UpdateTrack(ctx, track)
		require.NoError(t, err)

		got, err := repo.GetTrack(ctx, "track-user", "track-001")
		require.NoError(t, err)
		assert.Equal(t, "Updated Song", got.Title)
		assert.Equal(t, "Ambient", got.Genre)
		assert.Equal(t, models.VisibilityPublic, got.Visibility)
	})

	t.Run("delete track", func(t *testing.T) {
		err := repo.DeleteTrack(ctx, "track-user", "track-001")
		require.NoError(t, err)

		_, err = repo.GetTrack(ctx, "track-user", "track-001")
		assert.Error(t, err)
	})
}

func TestIntegration_TrackListPagination(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	userID := "paginate-user"
	user := models.User{
		ID:          userID,
		Email:       "paginate@test.com",
		DisplayName: "Paginate User",
		Role:        models.RoleSubscriber,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#"+userID, "PROFILE")

	// Create 12 tracks
	for i := 0; i < 12; i++ {
		track := models.Track{
			ID:         fmt.Sprintf("ptrack-%02d", i),
			UserID:     userID,
			Title:      fmt.Sprintf("Track %02d", i),
			Artist:     "Paginate Artist",
			Duration:   120,
			Format:     models.AudioFormatMP3,
			S3Key:      fmt.Sprintf("uploads/%s/ptrack-%02d.mp3", userID, i),
			Visibility: models.VisibilityPrivate,
		}
		require.NoError(t, repo.CreateTrack(ctx, track))
	}

	t.Run("paginate through all tracks", func(t *testing.T) {
		allTracks := map[string]bool{}
		cursor := ""

		for {
			filter := models.TrackFilter{Limit: 5, LastKey: cursor}
			result, err := repo.ListTracks(ctx, userID, filter)
			require.NoError(t, err)

			for _, tr := range result.Items {
				allTracks[tr.ID] = true
				assert.Equal(t, userID, tr.UserID)
			}

			if !result.HasMore {
				break
			}
			cursor = result.NextCursor
		}

		assert.Len(t, allTracks, 12, "pagination should yield all 12 tracks")
	})
}

// ---------------------------------------------------------------------------
// User CRUD
// ---------------------------------------------------------------------------

func TestIntegration_UserCRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	user := models.User{
		ID:          "user-crud-1",
		Email:       "crud@test.com",
		DisplayName: "CRUD User",
		Role:        models.RoleSubscriber,
	}

	t.Run("create and get user", func(t *testing.T) {
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
		tc.RegisterCleanup("dynamodb", "USER#user-crud-1", "PROFILE")

		got, err := repo.GetUser(ctx, "user-crud-1")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "user-crud-1", got.ID)
		assert.Equal(t, "crud@test.com", got.Email)
		assert.Equal(t, "CRUD User", got.DisplayName)
		assert.Equal(t, models.RoleSubscriber, got.Role)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := repo.GetUser(ctx, "nonexistent-user")
		assert.Error(t, err) // ErrUserNotFound
	})

	t.Run("update user role", func(t *testing.T) {
		err := repo.UpdateUserRole(ctx, "user-crud-1", models.RoleArtist)
		require.NoError(t, err)

		got, err := repo.GetUser(ctx, "user-crud-1")
		require.NoError(t, err)
		assert.Equal(t, models.RoleArtist, got.Role)
	})

	t.Run("update user", func(t *testing.T) {
		user.DisplayName = "Updated CRUD User"
		user.AvatarURL = "https://example.com/avatar.jpg"
		err := repo.UpdateUser(ctx, user)
		require.NoError(t, err)

		got, err := repo.GetUser(ctx, "user-crud-1")
		require.NoError(t, err)
		assert.Equal(t, "Updated CRUD User", got.DisplayName)
	})

	t.Run("GetUserByEmail via GSI", func(t *testing.T) {
		got, err := repo.GetUserByEmail(ctx, "crud@test.com")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "user-crud-1", got.ID)
	})

	t.Run("GetUserByEmail returns error for unknown email", func(t *testing.T) {
		_, err := repo.GetUserByEmail(ctx, "nobody@test.com")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// Tag CRUD
// ---------------------------------------------------------------------------

func TestIntegration_TagCRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	userID := "tag-user"
	user := models.User{
		ID:          userID,
		Email:       "taguser@test.com",
		DisplayName: "Tag User",
		Role:        models.RoleSubscriber,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#"+userID, "PROFILE")

	tag := models.Tag{
		UserID: userID,
		Name:   "chill",
		Color:  "#00FF00",
	}

	t.Run("create and get tag", func(t *testing.T) {
		err := repo.CreateTag(ctx, tag)
		require.NoError(t, err)

		got, err := repo.GetTag(ctx, userID, "chill")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "chill", got.Name)
		assert.Equal(t, "#00FF00", got.Color)
	})

	t.Run("list tags for user", func(t *testing.T) {
		// Create another tag
		tag2 := models.Tag{UserID: userID, Name: "energetic", Color: "#FF0000"}
		require.NoError(t, repo.CreateTag(ctx, tag2))

		tags, err := repo.ListTags(ctx, userID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tags), 2)

		tagNames := map[string]bool{}
		for _, tg := range tags {
			tagNames[tg.Name] = true
		}
		assert.True(t, tagNames["chill"])
		assert.True(t, tagNames["energetic"])
	})

	t.Run("add tags to track and get track tags", func(t *testing.T) {
		// Create a track first
		track := models.Track{
			ID:         "tagged-track",
			UserID:     userID,
			Title:      "Tagged Track",
			Artist:     "Tag Artist",
			Duration:   120,
			Format:     models.AudioFormatMP3,
			S3Key:      "uploads/" + userID + "/tagged-track.mp3",
			Visibility: models.VisibilityPrivate,
		}
		require.NoError(t, repo.CreateTrack(ctx, track))

		err := repo.AddTagsToTrack(ctx, userID, "tagged-track", []string{"chill", "energetic"})
		require.NoError(t, err)

		trackTags, err := repo.GetTrackTags(ctx, userID, "tagged-track")
		require.NoError(t, err)
		assert.Len(t, trackTags, 2)
		assert.Contains(t, trackTags, "chill")
		assert.Contains(t, trackTags, "energetic")
	})

	t.Run("get tracks by tag", func(t *testing.T) {
		tracks, err := repo.GetTracksByTag(ctx, userID, "chill")
		require.NoError(t, err)
		assert.Len(t, tracks, 1)
		assert.Equal(t, "tagged-track", tracks[0].ID)
	})

	t.Run("remove tag from track", func(t *testing.T) {
		err := repo.RemoveTagFromTrack(ctx, userID, "tagged-track", "energetic")
		require.NoError(t, err)

		trackTags, err := repo.GetTrackTags(ctx, userID, "tagged-track")
		require.NoError(t, err)
		assert.Len(t, trackTags, 1)
		assert.Contains(t, trackTags, "chill")
	})

	t.Run("delete tag", func(t *testing.T) {
		err := repo.DeleteTag(ctx, userID, "chill")
		require.NoError(t, err)

		got, err := repo.GetTag(ctx, userID, "chill")
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, got)
	})
}

// ---------------------------------------------------------------------------
// Playlist CRUD
// ---------------------------------------------------------------------------

func TestIntegration_PlaylistCRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	userID := "playlist-user"
	user := models.User{
		ID:          userID,
		Email:       "playlist@test.com",
		DisplayName: "Playlist User",
		Role:        models.RoleSubscriber,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#"+userID, "PROFILE")

	playlist := models.Playlist{
		ID:          "playlist-001",
		UserID:      userID,
		Name:        "Test Playlist",
		Description: "A test playlist",
		Visibility:  models.VisibilityPrivate,
	}

	t.Run("create and get playlist", func(t *testing.T) {
		err := repo.CreatePlaylist(ctx, playlist)
		require.NoError(t, err)

		got, err := repo.GetPlaylist(ctx, userID, "playlist-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "playlist-001", got.ID)
		assert.Equal(t, "Test Playlist", got.Name)
		assert.Equal(t, models.VisibilityPrivate, got.Visibility)
	})

	t.Run("update playlist visibility to public", func(t *testing.T) {
		err := repo.UpdatePlaylistVisibility(ctx, userID, "playlist-001", models.VisibilityPublic)
		require.NoError(t, err)

		got, err := repo.GetPlaylist(ctx, userID, "playlist-001")
		require.NoError(t, err)
		assert.Equal(t, models.VisibilityPublic, got.Visibility)
	})

	t.Run("list public playlists", func(t *testing.T) {
		result, err := repo.ListPublicPlaylists(ctx, 10, "")
		require.NoError(t, err)

		found := false
		for _, p := range result.Items {
			if p.ID == "playlist-001" {
				found = true
				break
			}
		}
		assert.True(t, found, "public playlist should appear in ListPublicPlaylists")
	})

	t.Run("set to private removes from public list", func(t *testing.T) {
		err := repo.UpdatePlaylistVisibility(ctx, userID, "playlist-001", models.VisibilityPrivate)
		require.NoError(t, err)

		result, err := repo.ListPublicPlaylists(ctx, 10, "")
		require.NoError(t, err)

		for _, p := range result.Items {
			assert.NotEqual(t, "playlist-001", p.ID, "private playlist should not appear in public list")
		}
	})

	t.Run("delete playlist", func(t *testing.T) {
		err := repo.DeletePlaylist(ctx, userID, "playlist-001")
		require.NoError(t, err)

		_, err = repo.GetPlaylist(ctx, userID, "playlist-001")
		assert.Error(t, err) // ErrPlaylistNotFound
	})
}

// ---------------------------------------------------------------------------
// Track Visibility
// ---------------------------------------------------------------------------

func TestIntegration_TrackVisibility(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	userID := "vis-user"
	user := models.User{
		ID:          userID,
		Email:       "vis@test.com",
		DisplayName: "Visibility User",
		Role:        models.RoleSubscriber,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#"+userID, "PROFILE")

	track := models.Track{
		ID:         "vis-track",
		UserID:     userID,
		Title:      "Visibility Track",
		Artist:     "Vis Artist",
		Duration:   120,
		Format:     models.AudioFormatMP3,
		S3Key:      "uploads/" + userID + "/vis-track.mp3",
		Visibility: models.VisibilityPrivate,
	}
	require.NoError(t, repo.CreateTrack(ctx, track))

	t.Run("update visibility to public", func(t *testing.T) {
		err := repo.UpdateTrackVisibility(ctx, userID, "vis-track", models.VisibilityPublic)
		require.NoError(t, err)

		got, err := repo.GetTrack(ctx, userID, "vis-track")
		require.NoError(t, err)
		assert.Equal(t, models.VisibilityPublic, got.Visibility)
	})

	t.Run("public tracks appear in ListPublicTracks", func(t *testing.T) {
		result, err := repo.ListPublicTracks(ctx, 10, "")
		require.NoError(t, err)

		found := false
		for _, tr := range result.Items {
			if tr.ID == "vis-track" {
				found = true
				break
			}
		}
		assert.True(t, found, "public track should appear in ListPublicTracks")
	})

	t.Run("private tracks hidden from ListPublicTracks", func(t *testing.T) {
		err := repo.UpdateTrackVisibility(ctx, userID, "vis-track", models.VisibilityPrivate)
		require.NoError(t, err)

		result, err := repo.ListPublicTracks(ctx, 10, "")
		require.NoError(t, err)

		for _, tr := range result.Items {
			assert.NotEqual(t, "vis-track", tr.ID, "private track should not appear in public list")
		}
	})
}
