//go:build integration

package service_test

import (
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_TrackCreation(t *testing.T) {
	// Setup LocalStack connection
	ctx, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	// Create a test user
	userID := ctx.CreateTestUser(t, "test@example.com", "subscriber")

	// Verify user was created
	userItem := ctx.GetItem(t, "USER#"+userID, "PROFILE")
	require.NotNil(t, userItem, "User should exist in DynamoDB")

	// Create a test track for the user
	trackID := ctx.CreateTestTrack(t, userID,
		testutil.WithTrackTitle("Integration Test Track"),
		testutil.WithTrackArtist("Test Artist"),
		testutil.WithTrackVisibility("public"),
	)

	// Verify track was created
	trackItem := ctx.GetItem(t, "USER#"+userID, "TRACK#"+trackID)
	require.NotNil(t, trackItem, "Track should exist in DynamoDB")

	// Verify track attributes
	titleAttr := trackItem["Title"]
	require.NotNil(t, titleAttr, "Track should have Title")

	// Check visibility
	visAttr := trackItem["Visibility"]
	require.NotNil(t, visAttr, "Track should have Visibility")
}

func TestIntegration_TrackCleanup(t *testing.T) {
	// Setup LocalStack connection
	ctx, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	// Create test data
	userID := ctx.CreateTestUser(t, "cleanup-test@example.com", "subscriber")
	trackID := ctx.CreateTestTrack(t, userID)

	// Verify data exists before cleanup
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "PROFILE"), "User should exist")
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+trackID), "Track should exist")

	// Manually cleanup the user (simulating test teardown)
	ctx.CleanupUser(t, userID)

	// Verify data is cleaned up
	assert.False(t, ctx.ItemExists(t, "USER#"+userID, "PROFILE"), "User should be deleted")
	assert.False(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+trackID), "Track should be deleted")
}

func TestIntegration_MultipleTracksPerUser(t *testing.T) {
	// Setup LocalStack connection
	ctx, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	// Create a test user
	userID := ctx.CreateTestUser(t, "multi-track@example.com", "artist")

	// Create multiple tracks
	track1ID := ctx.CreateTestTrack(t, userID, testutil.WithTrackTitle("Track 1"))
	track2ID := ctx.CreateTestTrack(t, userID, testutil.WithTrackTitle("Track 2"))
	track3ID := ctx.CreateTestTrack(t, userID, testutil.WithTrackTitle("Track 3"))

	// Verify all tracks exist
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+track1ID), "Track 1 should exist")
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+track2ID), "Track 2 should exist")
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+track3ID), "Track 3 should exist")

	// Cleanup single track
	ctx.CleanupTrack(t, userID, track2ID)
	assert.False(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+track2ID), "Track 2 should be deleted")
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+track1ID), "Track 1 should still exist")
	assert.True(t, ctx.ItemExists(t, "USER#"+userID, "TRACK#"+track3ID), "Track 3 should still exist")
}

func TestIntegration_PlaylistCreation(t *testing.T) {
	// Setup LocalStack connection
	ctx, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	// Create a test user
	userID := ctx.CreateTestUser(t, "playlist-test@example.com", "subscriber")

	// Create a playlist
	playlistID := ctx.CreateTestPlaylist(t, userID, "My Test Playlist")

	// Verify playlist exists
	playlistItem := ctx.GetItem(t, "USER#"+userID, "PLAYLIST#"+playlistID)
	require.NotNil(t, playlistItem, "Playlist should exist in DynamoDB")

	// Verify playlist name
	nameAttr := playlistItem["Name"]
	require.NotNil(t, nameAttr, "Playlist should have Name")
}

func TestIntegration_CognitoAuthentication(t *testing.T) {
	// Setup LocalStack connection
	ctx, cleanup := testutil.SetupLocalStack(t)
	defer cleanup()

	// Skip if Cognito not configured
	if ctx.UserPoolID == "" {
		t.Skip("Cognito not configured. Run init-cognito.sh first.")
	}

	// Get token for admin user
	token := ctx.GetTestUserToken(t, "admin")
	assert.NotEmpty(t, token, "Should get valid token for admin user")

	// Get token for subscriber user
	token = ctx.GetTestUserToken(t, "subscriber")
	assert.NotEmpty(t, token, "Should get valid token for subscriber user")

	// Get token for artist user
	token = ctx.GetTestUserToken(t, "artist")
	assert.NotEmpty(t, token, "Should get valid token for artist user")
}
