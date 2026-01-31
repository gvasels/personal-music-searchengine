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
// ArtistProfile CRUD
// ---------------------------------------------------------------------------

func TestIntegration_ArtistProfileCRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	// Create a user first (artist profiles are keyed by userID)
	user := models.User{
		ID:          "artist-user-1",
		Email:       "artist1@test.com",
		DisplayName: "Artist User 1",
		Role:        models.RoleArtist,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#artist-user-1", "PROFILE")

	profile := models.ArtistProfile{
		UserID:      "artist-user-1",
		DisplayName: "DJ Test",
		Bio:         "Integration test artist profile",
		Genres:      []string{"Electronic", "House"},
		SocialLinks: map[string]string{"twitter": "https://twitter.com/djtest"},
	}

	t.Run("create and get artist profile", func(t *testing.T) {
		err := repo.CreateArtistProfile(ctx, profile)
		require.NoError(t, err)

		got, err := repo.GetArtistProfile(ctx, "artist-user-1")
		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, "artist-user-1", got.UserID)
		assert.Equal(t, "DJ Test", got.DisplayName)
		assert.Equal(t, "Integration test artist profile", got.Bio)
		assert.Equal(t, []string{"Electronic", "House"}, got.Genres)
		assert.Equal(t, "https://twitter.com/djtest", got.SocialLinks["twitter"])
		assert.Equal(t, 0, got.FollowerCount)
		assert.False(t, got.CreatedAt.IsZero())
	})

	t.Run("get non-existent profile returns nil", func(t *testing.T) {
		got, err := repo.GetArtistProfile(ctx, "nonexistent-user")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("update artist profile", func(t *testing.T) {
		profile.DisplayName = "DJ Updated"
		profile.Bio = "Updated bio"
		profile.Genres = []string{"Techno", "Minimal"}
		err := repo.UpdateArtistProfile(ctx, profile)
		require.NoError(t, err)

		got, err := repo.GetArtistProfile(ctx, "artist-user-1")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "DJ Updated", got.DisplayName)
		assert.Equal(t, "Updated bio", got.Bio)
		assert.Equal(t, []string{"Techno", "Minimal"}, got.Genres)
	})

	t.Run("delete artist profile", func(t *testing.T) {
		err := repo.DeleteArtistProfile(ctx, "artist-user-1")
		require.NoError(t, err)

		got, err := repo.GetArtistProfile(ctx, "artist-user-1")
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestIntegration_ListArtistProfiles(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	// Create users and profiles
	for i := 1; i <= 5; i++ {
		userID := fmt.Sprintf("list-artist-user-%d", i)
		user := models.User{
			ID:          userID,
			Email:       fmt.Sprintf("artist%d@test.com", i),
			DisplayName: fmt.Sprintf("Artist %d", i),
			Role:        models.RoleArtist,
		}
		require.NoError(t, repo.CreateUser(ctx, user))
		tc.RegisterCleanup("dynamodb", "USER#"+userID, "PROFILE")

		profile := models.ArtistProfile{
			UserID:      userID,
			DisplayName: fmt.Sprintf("DJ Artist %d", i),
			Bio:         fmt.Sprintf("Bio for artist %d", i),
		}
		require.NoError(t, repo.CreateArtistProfile(ctx, profile))
	}

	t.Run("paginate through all profiles", func(t *testing.T) {
		allProfiles := make(map[string]bool)
		cursor := ""

		for {
			result, err := repo.ListArtistProfiles(ctx, 2, cursor)
			require.NoError(t, err)

			for _, p := range result.Items {
				allProfiles[p.UserID] = true
			}

			if !result.HasMore {
				break
			}
			cursor = result.NextCursor
		}

		assert.Len(t, allProfiles, 5, "pagination should yield all 5 profiles")
	})
}

func TestIntegration_IncrementArtistFollowerCount(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	user := models.User{
		ID:          "count-artist-user",
		Email:       "countartist@test.com",
		DisplayName: "Count Artist",
		Role:        models.RoleArtist,
	}
	require.NoError(t, repo.CreateUser(ctx, user))
	tc.RegisterCleanup("dynamodb", "USER#count-artist-user", "PROFILE")

	profile := models.ArtistProfile{
		UserID:      "count-artist-user",
		DisplayName: "Count DJ",
	}
	require.NoError(t, repo.CreateArtistProfile(ctx, profile))

	t.Run("increment by 1", func(t *testing.T) {
		err := repo.IncrementArtistFollowerCount(ctx, "count-artist-user", 1)
		require.NoError(t, err)

		got, err := repo.GetArtistProfile(ctx, "count-artist-user")
		require.NoError(t, err)
		assert.Equal(t, 1, got.FollowerCount)
	})

	t.Run("increment by 5", func(t *testing.T) {
		err := repo.IncrementArtistFollowerCount(ctx, "count-artist-user", 5)
		require.NoError(t, err)

		got, err := repo.GetArtistProfile(ctx, "count-artist-user")
		require.NoError(t, err)
		assert.Equal(t, 6, got.FollowerCount)
	})

	t.Run("decrement", func(t *testing.T) {
		err := repo.IncrementArtistFollowerCount(ctx, "count-artist-user", -2)
		require.NoError(t, err)

		got, err := repo.GetArtistProfile(ctx, "count-artist-user")
		require.NoError(t, err)
		assert.Equal(t, 4, got.FollowerCount)
	})
}

// ---------------------------------------------------------------------------
// Follow CRUD + GSI queries
// ---------------------------------------------------------------------------

func TestIntegration_FollowCRUD(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	// Create users
	for _, uid := range []string{"follower-1", "follower-2", "followed-artist"} {
		user := models.User{
			ID:          uid,
			Email:       uid + "@test.com",
			DisplayName: uid,
			Role:        models.RoleArtist,
		}
		require.NoError(t, repo.CreateUser(ctx, user))
		tc.RegisterCleanup("dynamodb", "USER#"+uid, "PROFILE")
	}

	// Create artist profile for the followed user
	profile := models.ArtistProfile{
		UserID:      "followed-artist",
		DisplayName: "Followed Artist",
	}
	require.NoError(t, repo.CreateArtistProfile(ctx, profile))

	follow := models.Follow{
		FollowerID: "follower-1",
		FollowedID: "followed-artist",
	}

	t.Run("create follow", func(t *testing.T) {
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	})

	t.Run("get follow returns record", func(t *testing.T) {
		got, err := repo.GetFollow(ctx, "follower-1", "followed-artist")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "follower-1", got.FollowerID)
		assert.Equal(t, "followed-artist", got.FollowedID)
		assert.False(t, got.CreatedAt.IsZero())
	})

	t.Run("get follow returns nil for nonexistent", func(t *testing.T) {
		got, err := repo.GetFollow(ctx, "follower-1", "nobody")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("list followers via GSI query", func(t *testing.T) {
		// Add a second follower
		follow2 := models.Follow{FollowerID: "follower-2", FollowedID: "followed-artist"}
		require.NoError(t, repo.CreateFollow(ctx, follow2))

		result, err := repo.ListFollowers(ctx, "followed-artist", 10, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 2)

		followerIDs := map[string]bool{}
		for _, f := range result.Items {
			followerIDs[f.FollowerID] = true
		}
		assert.True(t, followerIDs["follower-1"])
		assert.True(t, followerIDs["follower-2"])
	})

	t.Run("list following for a user", func(t *testing.T) {
		result, err := repo.ListFollowing(ctx, "follower-1", 10, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, "followed-artist", result.Items[0].FollowedID)
	})

	t.Run("delete follow", func(t *testing.T) {
		err := repo.DeleteFollow(ctx, "follower-1", "followed-artist")
		require.NoError(t, err)

		got, err := repo.GetFollow(ctx, "follower-1", "followed-artist")
		require.NoError(t, err)
		assert.Nil(t, got, "follow should be deleted")

		// Verify follower list reflects deletion
		result, err := repo.ListFollowers(ctx, "followed-artist", 10, "")
		require.NoError(t, err)
		assert.Len(t, result.Items, 1, "should have 1 follower after unfollow")
	})

	t.Run("delete follow on nonexistent is idempotent", func(t *testing.T) {
		err := repo.DeleteFollow(ctx, "follower-1", "followed-artist")
		assert.NoError(t, err)
	})
}

func TestIntegration_FollowPagination(t *testing.T) {
	tc, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	ctx := context.Background()

	// Create artist user with profile
	artistUser := models.User{
		ID:          "popular-artist",
		Email:       "popular@test.com",
		DisplayName: "Popular Artist",
		Role:        models.RoleArtist,
	}
	require.NoError(t, repo.CreateUser(ctx, artistUser))
	tc.RegisterCleanup("dynamodb", "USER#popular-artist", "PROFILE")

	profile := models.ArtistProfile{
		UserID:      "popular-artist",
		DisplayName: "Popular DJ",
	}
	require.NoError(t, repo.CreateArtistProfile(ctx, profile))

	// Create 8 followers
	for i := 1; i <= 8; i++ {
		uid := fmt.Sprintf("page-follower-%d", i)
		user := models.User{
			ID:          uid,
			Email:       fmt.Sprintf("f%d@test.com", i),
			DisplayName: uid,
			Role:        models.RoleSubscriber,
		}
		require.NoError(t, repo.CreateUser(ctx, user))
		tc.RegisterCleanup("dynamodb", "USER#"+uid, "PROFILE")

		follow := models.Follow{FollowerID: uid, FollowedID: "popular-artist"}
		require.NoError(t, repo.CreateFollow(ctx, follow))
	}

	t.Run("paginate through all followers", func(t *testing.T) {
		allFollowers := map[string]bool{}
		cursor := ""

		for {
			result, err := repo.ListFollowers(ctx, "popular-artist", 3, cursor)
			require.NoError(t, err)

			for _, f := range result.Items {
				allFollowers[f.FollowerID] = true
			}

			if !result.HasMore {
				break
			}
			cursor = result.NextCursor
		}

		assert.Len(t, allFollowers, 8, "pagination should yield all 8 followers")
	})
}
