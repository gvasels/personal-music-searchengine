//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_API_Health(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	t.Run("health endpoint returns 200 without auth", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/health")
		testutil.AssertStatus(t, resp, http.StatusOK)
		body := testutil.DecodeJSONBody(t, resp)
		assert.Equal(t, "ok", body["status"])
	})
}

func TestIntegration_API_AuthMiddleware(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	t.Run("protected endpoint without auth returns 401", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/tracks")
		// Without X-User-ID header, the auth middleware should reject
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("subscriber can access tracks", func(t *testing.T) {
		// Create user in DB first
		userID := tsc.CreateTestUser(t, "authtest-sub@test.com", "subscriber")

		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/tracks",
			testutil.AsUser(userID, models.RoleSubscriber),
		)
		// Should get 200 (even if empty)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("admin can access admin routes", func(t *testing.T) {
		if tsc.Services.Admin == nil {
			t.Skip("Admin service not available (no Cognito)")
		}
		adminID := tsc.CreateTestUser(t, "authtest-admin@test.com", "admin")

		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/admin/users",
			testutil.AsUser(adminID, models.RoleAdmin),
			testutil.WithQuery("q", "test"),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("subscriber cannot access admin routes", func(t *testing.T) {
		if tsc.Services.Admin == nil {
			t.Skip("Admin service not available (no Cognito)")
		}
		subID := tsc.CreateTestUser(t, "authtest-nonadmin@test.com", "subscriber")

		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/admin/users",
			testutil.AsUser(subID, models.RoleSubscriber),
			testutil.WithQuery("q", "test"),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		resp.Body.Close()
	})
}
