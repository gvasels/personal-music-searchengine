//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_API_PlaylistCRUD(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	userID := tsc.CreateTestUser(t, "plapi@test.com", "subscriber")
	trackID := tsc.CreateTestTrack(t, userID, testutil.WithTrackTitle("Playlist API Track"))
	var playlistID string

	t.Run("create playlist", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPost, "/api/v1/playlists",
			testutil.AsUser(userID, models.RoleSubscriber),
			testutil.WithJSON(map[string]interface{}{
				"name":        "API Test Playlist",
				"description": "Created via API test",
			}),
		)
		testutil.AssertStatus(t, resp, http.StatusCreated)
		body := testutil.DecodeJSONBody(t, resp)
		playlistID = body["id"].(string)
		assert.Equal(t, "API Test Playlist", body["name"])
	})

	t.Run("list playlists", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/playlists",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		items, ok := body["items"].([]interface{})
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(items), 1)
	})

	t.Run("add tracks to playlist", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/playlists/%s/tracks", playlistID),
			testutil.AsUser(userID, models.RoleSubscriber),
			testutil.WithJSON(map[string]interface{}{
				"trackIds": []string{trackID},
			}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("get playlist with tracks", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/playlists/%s", playlistID),
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		tracks, ok := body["tracks"].([]interface{})
		require.True(t, ok)
		assert.Len(t, tracks, 1)
	})

	t.Run("update playlist visibility to public", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/playlists/%s/visibility", playlistID),
			testutil.AsUser(userID, models.RoleSubscriber),
			testutil.WithJSON(map[string]string{"visibility": "public"}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("public playlist appears in discovery", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/playlists/public",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		items, ok := body["items"].([]interface{})
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(items), 1)
	})

	t.Run("delete playlist", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodDelete, fmt.Sprintf("/api/v1/playlists/%s", playlistID),
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusNoContent)
	})
}

func TestIntegration_API_TagsCRUD(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	userID := tsc.CreateTestUser(t, "tagapi@test.com", "subscriber")
	trackID := tsc.CreateTestTrack(t, userID, testutil.WithTrackTitle("Tag API Track"))

	t.Run("create tag", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPost, "/api/v1/tags",
			testutil.AsUser(userID, models.RoleSubscriber),
			testutil.WithJSON(map[string]string{
				"name":  "Chill",
				"color": "#0000ff",
			}),
		)
		testutil.AssertStatus(t, resp, http.StatusCreated)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, "chill", body["name"]) // normalized to lowercase
	})

	t.Run("add tags to track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/tracks/%s/tags", trackID),
			testutil.AsUser(userID, models.RoleSubscriber),
			testutil.WithJSON(map[string]interface{}{
				"tags": []string{"Chill", "Ambient"},
			}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("get tracks by tag", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/tags/chill/tracks",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("list tags", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/tags",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("remove tag from track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodDelete, fmt.Sprintf("/api/v1/tracks/%s/tags/chill", trackID),
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusNoContent)
	})

	t.Run("delete tag", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodDelete, "/api/v1/tags/chill",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusNoContent)
	})
}
