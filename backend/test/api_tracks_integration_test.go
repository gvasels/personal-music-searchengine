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

func TestIntegration_API_TracksCRUD(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	// Create test user and track
	userID := tsc.CreateTestUser(t, "trackapi@test.com", "subscriber")
	trackID := tsc.CreateTestTrack(t, userID, testutil.WithTrackTitle("API Test Track"), testutil.WithTrackArtist("API Artist"))

	t.Run("list tracks returns created track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/tracks",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		items, ok := body["items"].([]interface{})
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(items), 1)
	})

	t.Run("get track by ID", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, "API Test Track", body["title"])
	})

	t.Run("update track metadata", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(userID, models.RoleSubscriber),
			testutil.WithJSON(map[string]string{"title": "Updated Title"}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, "Updated Title", body["title"])
	})

	t.Run("delete track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodDelete, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusNoContent)

		// Verify gone
		resp = tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	})
}

func TestIntegration_API_TrackVisibility(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	ownerID := tsc.CreateTestUser(t, "track-owner@test.com", "subscriber")
	otherID := tsc.CreateTestUser(t, "track-other@test.com", "subscriber")
	adminID := tsc.CreateTestUser(t, "track-admin@test.com", "admin")
	trackID := tsc.CreateTestTrack(t, ownerID, testutil.WithTrackTitle("Private Track"))

	t.Run("owner can access own private track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(ownerID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("other user cannot access private track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(otherID, models.RoleSubscriber),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("admin can access private track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(adminID, models.RoleAdmin),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("make track public then other user can access", func(t *testing.T) {
		// Change visibility to public
		resp := tsc.DoRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/tracks/%s/visibility", trackID),
			testutil.AsUser(ownerID, models.RoleSubscriber),
			testutil.WithJSON(map[string]string{"visibility": "public"}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)

		// Other user can now access
		resp = tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(otherID, models.RoleSubscriber),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})
}

func TestIntegration_API_AdminTrackDelete(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	ownerID := tsc.CreateTestUser(t, "admindel-owner@test.com", "artist")
	adminID := tsc.CreateTestUser(t, "admindel-admin@test.com", "admin")
	trackID := tsc.CreateTestTrack(t, ownerID, testutil.WithTrackTitle("Admin Delete Me"))

	t.Run("admin can delete another users track", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodDelete, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(adminID, models.RoleAdmin),
		)
		testutil.AssertStatus(t, resp, http.StatusNoContent)

		// Verify track is gone
		resp = tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/tracks/%s", trackID),
			testutil.AsUser(ownerID, models.RoleSubscriber),
		)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	})
}
