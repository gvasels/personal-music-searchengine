//go:build integration

package service_test

import (
	"context"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/service"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCognitoClient(t *testing.T) (*testutil.TestContext, service.CognitoClient, func()) {
	t.Helper()
	tc, cleanup := testutil.SetupLocalStack(t)
	if tc.UserPoolID == "" {
		t.Skip("Cognito not available in LocalStack")
	}
	client := service.NewCognitoClient(tc.Cognito, tc.UserPoolID)
	return tc, client, cleanup
}

func TestIntegration_CognitoClient_Authentication(t *testing.T) {
	tc, _, cleanup := setupCognitoClient(t)
	defer cleanup()

	t.Run("authenticate admin user", func(t *testing.T) {
		result := tc.AuthenticateTestUser(t, testutil.TestUsers["admin"].Email, testutil.TestUsers["admin"].Password)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.IDToken)
	})

	t.Run("authenticate subscriber user", func(t *testing.T) {
		result := tc.AuthenticateTestUser(t, testutil.TestUsers["subscriber"].Email, testutil.TestUsers["subscriber"].Password)
		assert.NotEmpty(t, result.AccessToken)
	})

	t.Run("authenticate artist user", func(t *testing.T) {
		result := tc.AuthenticateTestUser(t, testutil.TestUsers["artist"].Email, testutil.TestUsers["artist"].Password)
		assert.NotEmpty(t, result.AccessToken)
	})
}

func TestIntegration_CognitoClient_GroupOperations(t *testing.T) {
	_, client, cleanup := setupCognitoClient(t)
	defer cleanup()
	ctx := context.Background()

	adminEmail := testutil.TestUsers["admin"].Email

	t.Run("get user groups for admin", func(t *testing.T) {
		groups, err := client.GetUserGroups(ctx, adminEmail)
		require.NoError(t, err)
		assert.Contains(t, groups, "admin")
	})

	t.Run("add user to group and verify", func(t *testing.T) {
		subscriberEmail := testutil.TestUsers["subscriber"].Email

		// Add subscriber to artist group
		err := client.AddUserToGroup(ctx, subscriberEmail, "artist")
		require.NoError(t, err)

		groups, err := client.GetUserGroups(ctx, subscriberEmail)
		require.NoError(t, err)
		assert.Contains(t, groups, "artist")
		assert.Contains(t, groups, "subscriber")

		// Clean up: remove from artist group
		err = client.RemoveUserFromGroup(ctx, subscriberEmail, "artist")
		require.NoError(t, err)

		groups, err = client.GetUserGroups(ctx, subscriberEmail)
		require.NoError(t, err)
		assert.NotContains(t, groups, "artist")
	})

	t.Run("remove user from group", func(t *testing.T) {
		artistEmail := testutil.TestUsers["artist"].Email

		// Remove from artist group
		err := client.RemoveUserFromGroup(ctx, artistEmail, "artist")
		require.NoError(t, err)

		groups, err := client.GetUserGroups(ctx, artistEmail)
		require.NoError(t, err)
		assert.NotContains(t, groups, "artist")

		// Restore: add back to artist group
		err = client.AddUserToGroup(ctx, artistEmail, "artist")
		require.NoError(t, err)
	})
}

func TestIntegration_CognitoClient_UserManagement(t *testing.T) {
	_, client, cleanup := setupCognitoClient(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("search users by email prefix", func(t *testing.T) {
		users, err := client.SearchUsers(ctx, "admin@", 10)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "admin@local.test", users[0].Email)
		assert.True(t, users[0].Enabled)
	})

	t.Run("search returns multiple results", func(t *testing.T) {
		// All test users share @local.test domain
		users, err := client.SearchUsers(ctx, "", 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3)
	})

	t.Run("get user status", func(t *testing.T) {
		enabled, err := client.GetUserStatus(ctx, testutil.TestUsers["subscriber"].Email)
		require.NoError(t, err)
		assert.True(t, enabled)
	})

	t.Run("disable and enable user", func(t *testing.T) {
		artistEmail := testutil.TestUsers["artist"].Email

		// Disable user
		err := client.DisableUser(ctx, artistEmail)
		require.NoError(t, err)

		enabled, err := client.GetUserStatus(ctx, artistEmail)
		require.NoError(t, err)
		assert.False(t, enabled)

		// Re-enable user
		err = client.EnableUser(ctx, artistEmail)
		require.NoError(t, err)

		enabled, err = client.GetUserStatus(ctx, artistEmail)
		require.NoError(t, err)
		assert.True(t, enabled)
	})

	t.Run("operations on nonexistent user return error", func(t *testing.T) {
		_, err := client.GetUserGroups(ctx, "nonexistent@test.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
