//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_API_ArtistProfileCRUD(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	artistUserID := tsc.CreateTestUser(t, "apapi-artist@test.com", "artist")

	t.Run("create artist profile", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPost, "/api/v1/artists/entity",
			testutil.AsUser(artistUserID, models.RoleArtist),
			testutil.WithJSON(map[string]interface{}{
				"displayName": "DJ API Test",
				"bio":         "Integration test artist",
				"genres":      []string{"Electronic", "House"},
			}),
		)
		testutil.AssertStatus(t, resp, http.StatusCreated)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, "DJ API Test", body["displayName"])
	})

	t.Run("get artist profile", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/artists/entity/%s", artistUserID),
			testutil.AsUser(artistUserID, models.RoleArtist),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, "DJ API Test", body["displayName"])
	})

	t.Run("list artist profiles", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/artists/entity",
			testutil.AsUser(artistUserID, models.RoleArtist),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("update artist profile", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/artists/entity/%s", artistUserID),
			testutil.AsUser(artistUserID, models.RoleArtist),
			testutil.WithJSON(map[string]interface{}{
				"bio": "Updated bio from API test",
			}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})
}

func TestIntegration_API_FollowSystem(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	artistID := tsc.CreateTestUser(t, "followapi-artist@test.com", "artist")
	followerID := tsc.CreateTestUser(t, "followapi-follower@test.com", "subscriber")

	// Artist needs a profile
	tsc.CreateTestArtistProfile(t, artistID)

	t.Run("follow artist", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/artists/entity/%s/follow", artistID),
			testutil.AsUser(followerID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("check is following", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/artists/entity/%s/following", artistID),
			testutil.AsUser(followerID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, true, body["following"])
	})

	t.Run("get followers list", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/artists/entity/%s/followers", artistID),
			testutil.AsUser(followerID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("unfollow artist", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodDelete, fmt.Sprintf("/api/v1/artists/entity/%s/follow", artistID),
			testutil.AsUser(followerID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("verify not following anymore", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/artists/entity/%s/following", artistID),
			testutil.AsUser(followerID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, false, body["following"])
	})
}
