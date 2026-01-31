//go:build integration

package service_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTrackService creates real repo + service backed by LocalStack.
func setupTrackService(t *testing.T) (*testutil.TestContext, service.TrackService, repository.Repository, func()) {
	t.Helper()
	tc, cleanup := testutil.SetupLocalStack(t)
	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	presignClient := s3.NewPresignClient(tc.S3)
	s3Repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)
	svc := service.NewTrackService(repo, s3Repo)
	return tc, svc, repo, cleanup
}

func TestIntegration_TrackService_VisibilityEnforcement(t *testing.T) {
	tc, svc, repo, cleanup := setupTrackService(t)
	defer cleanup()
	ctx := context.Background()

	// Create two users
	owner := models.User{ID: "owner-user", Email: "owner@test.com", DisplayName: "Owner", Role: models.RoleSubscriber}
	other := models.User{ID: "other-user", Email: "other@test.com", DisplayName: "Other", Role: models.RoleSubscriber}
	admin := models.User{ID: "admin-user", Email: "admin@test.com", DisplayName: "Admin", Role: models.RoleAdmin}
	for _, u := range []models.User{owner, other, admin} {
		require.NoError(t, repo.CreateUser(ctx, u))
		tc.RegisterCleanup("dynamodb", "USER#"+u.ID, "PROFILE")
	}

	// Create a private track
	track := models.Track{
		ID: "private-track", UserID: "owner-user", Title: "Private Song",
		Artist: "Test", Duration: 120, Format: models.AudioFormatMP3,
		S3Key: "uploads/owner-user/private-track.mp3", Visibility: models.VisibilityPrivate,
	}
	require.NoError(t, repo.CreateTrack(ctx, track))

	t.Run("owner can access own private track", func(t *testing.T) {
		resp, err := svc.GetTrack(ctx, "owner-user", "private-track", false)
		require.NoError(t, err)
		assert.Equal(t, "Private Song", resp.Title)
	})

	t.Run("other user cannot access private track", func(t *testing.T) {
		_, err := svc.GetTrack(ctx, "other-user", "private-track", false)
		require.Error(t, err)
		// Should be a forbidden error, not a not-found
		assert.Contains(t, err.Error(), "forbidden")
	})

	t.Run("admin with hasGlobal can access private track", func(t *testing.T) {
		resp, err := svc.GetTrack(ctx, "admin-user", "private-track", true)
		require.NoError(t, err)
		assert.Equal(t, "Private Song", resp.Title)
	})

	t.Run("public track accessible by anyone", func(t *testing.T) {
		// Make track public
		require.NoError(t, repo.UpdateTrackVisibility(ctx, "owner-user", "private-track", models.VisibilityPublic))

		resp, err := svc.GetTrack(ctx, "other-user", "private-track", false)
		require.NoError(t, err)
		assert.Equal(t, "Private Song", resp.Title)
	})
}

func TestIntegration_TrackService_AdminDelete(t *testing.T) {
	tc, svc, repo, cleanup := setupTrackService(t)
	defer cleanup()
	ctx := context.Background()

	// Create owner and admin users
	owner := models.User{ID: "del-owner", Email: "delowner@test.com", DisplayName: "Del Owner", Role: models.RoleArtist}
	admin := models.User{ID: "del-admin", Email: "deladmin@test.com", DisplayName: "Del Admin", Role: models.RoleAdmin}
	require.NoError(t, repo.CreateUser(ctx, owner))
	require.NoError(t, repo.CreateUser(ctx, admin))
	tc.RegisterCleanup("dynamodb", "USER#del-owner", "PROFILE")
	tc.RegisterCleanup("dynamodb", "USER#del-admin", "PROFILE")

	// Create track with S3 objects (audio + HLS)
	track := models.Track{
		ID: "del-track", UserID: "del-owner", Title: "To Delete",
		Artist: "Test", Duration: 120, Format: models.AudioFormatMP3,
		S3Key: "uploads/del-owner/del-track.mp3", Visibility: models.VisibilityPrivate,
	}
	require.NoError(t, repo.CreateTrack(ctx, track))

	// Create S3 objects for the track
	_, err := tc.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName), Key: aws.String("uploads/del-owner/del-track.mp3"),
		Body: bytes.NewReader([]byte("audio")), ContentType: aws.String("audio/mpeg"),
	})
	require.NoError(t, err)
	_, err = tc.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(tc.BucketName), Key: aws.String("hls/del-owner/del-track/master.m3u8"),
		Body: bytes.NewReader([]byte("#EXTM3U")), ContentType: aws.String("application/x-mpegURL"),
	})
	require.NoError(t, err)

	t.Run("admin can delete another users track", func(t *testing.T) {
		err := svc.DeleteTrack(ctx, "del-admin", "del-track", true)
		require.NoError(t, err)

		// Verify track is gone from DynamoDB
		_, err = repo.GetTrack(ctx, "del-owner", "del-track")
		assert.Error(t, err)
	})
}

func TestIntegration_TrackService_ListWithVisibility(t *testing.T) {
	tc, svc, repo, cleanup := setupTrackService(t)
	defer cleanup()
	ctx := context.Background()

	user1 := models.User{ID: "list-user1", Email: "list1@test.com", DisplayName: "List User 1", Role: models.RoleSubscriber}
	user2 := models.User{ID: "list-user2", Email: "list2@test.com", DisplayName: "List User 2", Role: models.RoleSubscriber}
	require.NoError(t, repo.CreateUser(ctx, user1))
	require.NoError(t, repo.CreateUser(ctx, user2))
	tc.RegisterCleanup("dynamodb", "USER#list-user1", "PROFILE")
	tc.RegisterCleanup("dynamodb", "USER#list-user2", "PROFILE")

	// User1 tracks: 2 private, 1 public
	for i, vis := range []models.TrackVisibility{models.VisibilityPrivate, models.VisibilityPrivate, models.VisibilityPublic} {
		track := models.Track{
			ID: fmt.Sprintf("u1-track-%d", i), UserID: "list-user1", Title: fmt.Sprintf("U1 Track %d", i),
			Artist: "U1 Artist", Duration: 120, Format: models.AudioFormatMP3,
			S3Key: fmt.Sprintf("uploads/list-user1/track-%d.mp3", i), Visibility: vis,
		}
		require.NoError(t, repo.CreateTrack(ctx, track))
	}

	// User2 tracks: 1 private, 1 public
	for i, vis := range []models.TrackVisibility{models.VisibilityPrivate, models.VisibilityPublic} {
		track := models.Track{
			ID: fmt.Sprintf("u2-track-%d", i), UserID: "list-user2", Title: fmt.Sprintf("U2 Track %d", i),
			Artist: "U2 Artist", Duration: 120, Format: models.AudioFormatMP3,
			S3Key: fmt.Sprintf("uploads/list-user2/track-%d.mp3", i), Visibility: vis,
		}
		require.NoError(t, repo.CreateTrack(ctx, track))
	}

	t.Run("user1 sees own tracks", func(t *testing.T) {
		result, err := svc.ListTracks(ctx, "list-user1", models.TrackFilter{Limit: 20})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Items), 3, "user1 should see all own tracks")
	})

	t.Run("admin sees all tracks with hasGlobal", func(t *testing.T) {
		result, err := svc.ListTracks(ctx, "list-user1", models.TrackFilter{Limit: 20, GlobalScope: true})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Items), 5, "admin should see all tracks")
	})
}
