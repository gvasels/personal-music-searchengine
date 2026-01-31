//go:build integration

package service_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPlaylistService(t *testing.T) (*testutil.TestContext, service.PlaylistService, repository.Repository, func()) {
	t.Helper()
	tc, cleanup := testutil.SetupLocalStack(t)
	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	presignClient := s3.NewPresignClient(tc.S3)
	s3Repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	svc := service.NewPlaylistService(repo, s3Repo)
	return tc, svc, repo, cleanup
}

func TestIntegration_PlaylistService_CRUD(t *testing.T) {
	tc, svc, repo, cleanup := setupPlaylistService(t)
	defer cleanup()
	ctx := context.Background()

	// Create user
	user := models.User{ID: "pl-user", Email: "pluser@test.com", DisplayName: "PL User", Role: models.RoleSubscriber}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#pl-user", "PROFILE")

	// Create a track to add to playlist
	track := models.Track{
		ID: "pl-track-1", UserID: "pl-user", Title: "Playlist Track",
		Artist: "Test", Duration: 180, Format: models.AudioFormatMP3,
		S3Key: "uploads/pl-user/pl-track-1.mp3", Visibility: models.VisibilityPrivate,
	}
	require.NoError(t, repo.CreateTrack(ctx, track))

	var playlistID string

	t.Run("create playlist", func(t *testing.T) {
		resp, err := svc.CreatePlaylist(ctx, "pl-user", models.CreatePlaylistRequest{
			Name:        "My Playlist",
			Description: "Integration test playlist",
		})
		require.NoError(t, err)
		assert.Equal(t, "My Playlist", resp.Name)
		assert.Equal(t, "Integration test playlist", resp.Description)
		assert.Equal(t, 0, resp.TrackCount)
		playlistID = resp.ID
	})

	t.Run("get playlist", func(t *testing.T) {
		resp, err := svc.GetPlaylist(ctx, "pl-user", playlistID)
		require.NoError(t, err)
		assert.Equal(t, "My Playlist", resp.Playlist.Name)
		assert.Empty(t, resp.Tracks)
	})

	t.Run("add tracks to playlist", func(t *testing.T) {
		resp, err := svc.AddTracks(ctx, "pl-user", playlistID, models.AddTracksToPlaylistRequest{
			TrackIDs: []string{"pl-track-1"},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, resp.TrackCount)
	})

	t.Run("get playlist with tracks", func(t *testing.T) {
		resp, err := svc.GetPlaylist(ctx, "pl-user", playlistID)
		require.NoError(t, err)
		assert.Len(t, resp.Tracks, 1)
		assert.Equal(t, "Playlist Track", resp.Tracks[0].Title)
	})

	t.Run("update playlist", func(t *testing.T) {
		newName := "Updated Playlist"
		resp, err := svc.UpdatePlaylist(ctx, "pl-user", playlistID, models.UpdatePlaylistRequest{
			Name: &newName,
		})
		require.NoError(t, err)
		assert.Equal(t, "Updated Playlist", resp.Name)
	})

	t.Run("remove tracks from playlist", func(t *testing.T) {
		resp, err := svc.RemoveTracks(ctx, "pl-user", playlistID, models.RemoveTracksFromPlaylistRequest{
			TrackIDs: []string{"pl-track-1"},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, resp.TrackCount)
	})

	t.Run("delete playlist", func(t *testing.T) {
		err := svc.DeletePlaylist(ctx, "pl-user", playlistID)
		require.NoError(t, err)

		_, err = svc.GetPlaylist(ctx, "pl-user", playlistID)
		assert.Error(t, err)
	})
}

func TestIntegration_PlaylistService_ListAndVisibility(t *testing.T) {
	tc, svc, repo, cleanup := setupPlaylistService(t)
	defer cleanup()
	ctx := context.Background()

	user := models.User{ID: "plv-user", Email: "plvuser@test.com", DisplayName: "PLV User", Role: models.RoleSubscriber}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#plv-user", "PROFILE")

	// Create two playlists
	resp1, err := svc.CreatePlaylist(ctx, "plv-user", models.CreatePlaylistRequest{Name: "Private Playlist"})
	require.NoError(t, err)

	resp2, err := svc.CreatePlaylist(ctx, "plv-user", models.CreatePlaylistRequest{Name: "Public Playlist", IsPublic: true})
	require.NoError(t, err)

	t.Run("list user playlists", func(t *testing.T) {
		result, err := svc.ListPlaylists(ctx, "plv-user", models.PlaylistFilter{Limit: 20})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Items), 2)
	})

	t.Run("update visibility to public", func(t *testing.T) {
		err := svc.UpdateVisibility(ctx, "plv-user", resp1.ID, models.VisibilityPublic)
		require.NoError(t, err)
	})

	t.Run("list public playlists", func(t *testing.T) {
		result, err := svc.ListPublicPlaylists(ctx, 20, "")
		require.NoError(t, err)
		// Both playlists should now be public/discoverable
		assert.GreaterOrEqual(t, len(result.Items), 1)
	})

	_ = resp2 // used for creation
}
