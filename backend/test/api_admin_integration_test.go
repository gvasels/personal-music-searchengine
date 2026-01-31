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

func TestIntegration_API_AdminRoutes(t *testing.T) {
	tsc, cleanup := testutil.SetupTestServer(t)
	defer cleanup()

	if tsc.Services.Admin == nil {
		t.Skip("Admin service not available (Cognito not configured in LocalStack)")
	}

	adminID := tsc.CreateTestUser(t, "adminapi-admin@test.com", "admin")
	subID := tsc.CreateTestUser(t, "adminapi-sub@test.com", "subscriber")

	t.Run("non-admin gets 403 on admin routes", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/admin/users",
			testutil.AsUser(subID, models.RoleSubscriber),
			testutil.WithQuery("q", "admin"),
		)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("admin can search users", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/admin/users",
			testutil.AsUser(adminID, models.RoleAdmin),
			testutil.WithQuery("q", "admin"),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("admin can get user details", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/admin/users/%s", subID),
			testutil.AsUser(adminID, models.RoleAdmin),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("admin can update user role", func(t *testing.T) {
		resp := tsc.DoRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/admin/users/%s/role", subID),
			testutil.AsUser(adminID, models.RoleAdmin),
			testutil.WithJSON(map[string]string{"role": "artist"}),
		)
		testutil.AssertStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	})

	t.Run("DB role check overrides header role", func(t *testing.T) {
		// User has subscriber in DB but sends admin in headers
		// The RequireRoleWithDBCheck middleware should check DB and reject
		resp := tsc.DoRequest(t, http.MethodGet, "/api/v1/admin/users",
			testutil.AsUser(subID, models.RoleAdmin), // Lie in headers
			testutil.WithQuery("q", "test"),
		)
		// The middleware queries DB for real role, which should not be admin
		// (subID was upgraded to artist above, not admin)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		resp.Body.Close()
	})
}
