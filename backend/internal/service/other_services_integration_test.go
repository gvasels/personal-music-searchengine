//go:build integration

package service_test

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// TagService integration tests
// ---------------------------------------------------------------------------

func TestIntegration_TagService_CRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()
	ctx := context.Background()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	tagSvc := service.NewTagService(repo)

	// Create user and track
	user := models.User{ID: "tag-user", Email: "taguser@test.com", DisplayName: "Tag User", Role: models.RoleSubscriber}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#tag-user", "PROFILE")

	track := models.Track{
		ID: "tag-track", UserID: "tag-user", Title: "Tagged Song",
		Artist: "Test", Duration: 120, Format: models.AudioFormatMP3,
		S3Key: "uploads/tag-user/tag-track.mp3", Visibility: models.VisibilityPrivate,
	}
	require.NoError(t, repo.CreateTrack(ctx, track))

	t.Run("create tag", func(t *testing.T) {
		resp, err := tagSvc.CreateTag(ctx, "tag-user", models.CreateTagRequest{
			Name:  "Electronic",
			Color: "#00ff00",
		})
		require.NoError(t, err)
		assert.Equal(t, "electronic", resp.Name) // normalized to lowercase
		assert.Equal(t, "#00ff00", resp.Color)
	})

	t.Run("get tag is case-insensitive", func(t *testing.T) {
		resp, err := tagSvc.GetTag(ctx, "tag-user", "ELECTRONIC")
		require.NoError(t, err)
		assert.Equal(t, "electronic", resp.Name)
	})

	t.Run("duplicate tag returns conflict", func(t *testing.T) {
		_, err := tagSvc.CreateTag(ctx, "tag-user", models.CreateTagRequest{Name: "Electronic"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("add tags to track", func(t *testing.T) {
		allTags, err := tagSvc.AddTagsToTrack(ctx, "tag-user", "tag-track", models.AddTagsToTrackRequest{
			Tags: []string{"Electronic", "House"},
		})
		require.NoError(t, err)
		assert.Contains(t, allTags, "electronic")
		assert.Contains(t, allTags, "house")
	})

	t.Run("get tracks by tag", func(t *testing.T) {
		tracks, err := tagSvc.GetTracksByTag(ctx, "tag-user", "electronic")
		require.NoError(t, err)
		assert.Len(t, tracks, 1)
		assert.Equal(t, "Tagged Song", tracks[0].Title)
	})

	t.Run("remove tag from track", func(t *testing.T) {
		err := tagSvc.RemoveTagFromTrack(ctx, "tag-user", "tag-track", "electronic")
		require.NoError(t, err)

		tracks, err := tagSvc.GetTracksByTag(ctx, "tag-user", "electronic")
		require.NoError(t, err)
		assert.Empty(t, tracks)
	})

	t.Run("list tags", func(t *testing.T) {
		tags, err := tagSvc.ListTags(ctx, "tag-user")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tags), 2) // electronic + house
	})

	t.Run("delete tag", func(t *testing.T) {
		err := tagSvc.DeleteTag(ctx, "tag-user", "electronic")
		require.NoError(t, err)

		_, err = tagSvc.GetTag(ctx, "tag-user", "electronic")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// UserService integration tests
// ---------------------------------------------------------------------------

func TestIntegration_UserService_ProfileAndSettings(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()
	ctx := context.Background()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	userSvc := service.NewUserService(repo)

	t.Run("create user if not exists", func(t *testing.T) {
		user, err := userSvc.CreateUserIfNotExists(ctx, "svc-user-1", "svcuser@test.com", "SVC User")
		require.NoError(t, err)
		assert.Equal(t, "svc-user-1", user.ID)
		assert.Equal(t, "svcuser@test.com", user.Email)
		assert.Equal(t, models.RoleSubscriber, user.Role)
		tc.RegisterCleanup("dynamodb", "USER#svc-user-1", "PROFILE")
	})

	t.Run("create user if not exists is idempotent", func(t *testing.T) {
		user, err := userSvc.CreateUserIfNotExists(ctx, "svc-user-1", "svcuser@test.com", "SVC User")
		require.NoError(t, err)
		assert.Equal(t, "svc-user-1", user.ID)
	})

	t.Run("get profile", func(t *testing.T) {
		resp, err := userSvc.GetProfile(ctx, "svc-user-1")
		require.NoError(t, err)
		assert.Equal(t, "SVC User", resp.DisplayName)
	})

	t.Run("update profile", func(t *testing.T) {
		newName := "Updated Name"
		resp, err := userSvc.UpdateProfile(ctx, "svc-user-1", models.UpdateUserRequest{
			DisplayName: &newName,
		})
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", resp.DisplayName)
	})

	t.Run("get non-existent user returns not found", func(t *testing.T) {
		_, err := userSvc.GetProfile(ctx, "nonexistent-user")
		require.Error(t, err)
	})

	t.Run("get user role", func(t *testing.T) {
		role, err := userSvc.GetUserRole(ctx, "svc-user-1")
		require.NoError(t, err)
		assert.Equal(t, models.RoleSubscriber, role)
	})

	t.Run("get role for nonexistent user defaults to subscriber", func(t *testing.T) {
		role, err := userSvc.GetUserRole(ctx, "no-such-user")
		require.NoError(t, err)
		assert.Equal(t, models.RoleSubscriber, role)
	})
}

// ---------------------------------------------------------------------------
// RoleService integration tests
// ---------------------------------------------------------------------------

func TestIntegration_RoleService(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()
	ctx := context.Background()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	roleSvc := service.NewRoleService(repo)

	// Create test user
	user := models.User{ID: "role-user", Email: "roleuser@test.com", DisplayName: "Role User", Role: models.RoleSubscriber}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#role-user", "PROFILE")

	t.Run("get user role", func(t *testing.T) {
		role, err := roleSvc.GetUserRole(ctx, "role-user")
		require.NoError(t, err)
		assert.Equal(t, models.RoleSubscriber, role)
	})

	t.Run("set user role to artist", func(t *testing.T) {
		err := roleSvc.SetUserRole(ctx, "role-user", models.RoleArtist)
		require.NoError(t, err)

		role, err := roleSvc.GetUserRole(ctx, "role-user")
		require.NoError(t, err)
		assert.Equal(t, models.RoleArtist, role)
	})

	t.Run("has permission check", func(t *testing.T) {
		hasPerm, err := roleSvc.HasPermission(ctx, "role-user", models.PermissionUploadTracks)
		require.NoError(t, err)
		assert.True(t, hasPerm, "artist should have upload permission")
	})

	t.Run("get role for nonexistent user returns not found", func(t *testing.T) {
		_, err := roleSvc.GetUserRole(ctx, "no-such-user")
		require.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// FollowService integration tests
// ---------------------------------------------------------------------------

func TestIntegration_FollowService(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()
	ctx := context.Background()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	followSvc := service.NewFollowService(repo)

	// Create users
	follower := models.User{ID: "fol-follower", Email: "follower@test.com", DisplayName: "Follower", Role: models.RoleSubscriber}
	artist := models.User{ID: "fol-artist", Email: "artist@test.com", DisplayName: "Artist", Role: models.RoleArtist}
	require.NoError(t, repo.CreateUser(ctx, follower))
	require.NoError(t, repo.CreateUser(ctx, artist))
	tc.RegisterCleanup("dynamodb", "USER#fol-follower", "PROFILE")
	tc.RegisterCleanup("dynamodb", "USER#fol-artist", "PROFILE")

	// Artist needs a profile for follow to work
	profile := models.ArtistProfile{
		UserID:      "fol-artist",
		DisplayName: "DJ Artist",
		Genres:      []string{"House"},
	}
	require.NoError(t, repo.CreateArtistProfile(ctx, profile))

	t.Run("follow artist", func(t *testing.T) {
		err := followSvc.Follow(ctx, "fol-follower", "fol-artist")
		require.NoError(t, err)
	})

	t.Run("is following returns true", func(t *testing.T) {
		following, err := followSvc.IsFollowing(ctx, "fol-follower", "fol-artist")
		require.NoError(t, err)
		assert.True(t, following)
	})

	t.Run("is following returns false for non-follower", func(t *testing.T) {
		following, err := followSvc.IsFollowing(ctx, "fol-artist", "fol-follower")
		require.NoError(t, err)
		assert.False(t, following)
	})

	t.Run("cannot follow self", func(t *testing.T) {
		err := followSvc.Follow(ctx, "fol-follower", "fol-follower")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot follow yourself")
	})

	t.Run("duplicate follow returns error", func(t *testing.T) {
		err := followSvc.Follow(ctx, "fol-follower", "fol-artist")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already following")
	})

	t.Run("get followers", func(t *testing.T) {
		result, err := followSvc.GetFollowers(ctx, "fol-artist", 10, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, "fol-follower", result.Items[0].FollowerID)
	})

	t.Run("get following", func(t *testing.T) {
		result, err := followSvc.GetFollowing(ctx, "fol-follower", 10, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, "fol-artist", result.Items[0].FollowedID)
	})

	t.Run("unfollow artist", func(t *testing.T) {
		err := followSvc.Unfollow(ctx, "fol-follower", "fol-artist")
		require.NoError(t, err)

		following, err := followSvc.IsFollowing(ctx, "fol-follower", "fol-artist")
		require.NoError(t, err)
		assert.False(t, following)
	})

	t.Run("follow requires artist profile", func(t *testing.T) {
		err := followSvc.Follow(ctx, "fol-artist", "fol-follower")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "artist profile not found")
	})
}

// ---------------------------------------------------------------------------
// ArtistProfileService integration tests
// ---------------------------------------------------------------------------

func TestIntegration_ArtistProfileService(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()
	ctx := context.Background()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	artistSvc := service.NewArtistProfileService(repo)

	// Create artist user
	artistUser := models.User{ID: "ap-artist", Email: "aPartist@test.com", DisplayName: "AP Artist", Role: models.RoleArtist}
	subscriberUser := models.User{ID: "ap-sub", Email: "apsub@test.com", DisplayName: "AP Sub", Role: models.RoleSubscriber}
	require.NoError(t, repo.CreateUser(ctx, artistUser))
	require.NoError(t, repo.CreateUser(ctx, subscriberUser))
	tc.RegisterCleanup("dynamodb", "USER#ap-artist", "PROFILE")
	tc.RegisterCleanup("dynamodb", "USER#ap-sub", "PROFILE")

	t.Run("create artist profile", func(t *testing.T) {
		resp, err := artistSvc.CreateProfile(ctx, "ap-artist", models.CreateArtistProfileRequest{
			DisplayName: "DJ Artist Profile",
			Bio:         "Test artist bio",
			Genres:      []string{"Techno", "House"},
		})
		require.NoError(t, err)
		assert.Equal(t, "ap-artist", resp.UserID)
		assert.Equal(t, "DJ Artist Profile", resp.DisplayName)
		assert.Equal(t, []string{"Techno", "House"}, resp.Genres)
	})

	t.Run("subscriber cannot create artist profile", func(t *testing.T) {
		_, err := artistSvc.CreateProfile(ctx, "ap-sub", models.CreateArtistProfileRequest{
			DisplayName: "Not Allowed",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "artist role required")
	})

	t.Run("get artist profile", func(t *testing.T) {
		resp, err := artistSvc.GetProfile(ctx, "ap-artist")
		require.NoError(t, err)
		assert.Equal(t, "DJ Artist Profile", resp.DisplayName)
		assert.Equal(t, "Test artist bio", resp.Bio)
	})

	t.Run("update artist profile", func(t *testing.T) {
		newBio := "Updated bio"
		resp, err := artistSvc.UpdateProfile(ctx, "ap-artist", "ap-artist", models.UpdateArtistProfileRequest{
			Bio: &newBio,
		})
		require.NoError(t, err)
		assert.Equal(t, "Updated bio", resp.Bio)
	})

	t.Run("non-owner cannot update profile", func(t *testing.T) {
		newBio := "Hacked"
		_, err := artistSvc.UpdateProfile(ctx, "ap-sub", "ap-artist", models.UpdateArtistProfileRequest{
			Bio: &newBio,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})

	t.Run("list artist profiles", func(t *testing.T) {
		result, err := artistSvc.ListProfiles(ctx, 10, "")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Items), 1)
	})

	t.Run("delete artist profile", func(t *testing.T) {
		err := artistSvc.DeleteProfile(ctx, "ap-artist", "ap-artist")
		require.NoError(t, err)

		_, err = artistSvc.GetProfile(ctx, "ap-artist")
		assert.Error(t, err)
	})

	t.Run("non-owner cannot delete profile", func(t *testing.T) {
		// Recreate for this test
		_, err := artistSvc.CreateProfile(ctx, "ap-artist", models.CreateArtistProfileRequest{
			DisplayName: "Recreated",
		})
		require.NoError(t, err)

		err = artistSvc.DeleteProfile(ctx, "ap-sub", "ap-artist")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden")
	})
}
