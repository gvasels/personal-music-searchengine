package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuth(t *testing.T) {
	t.Run("passes when user is authenticated", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequireAuth()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("rejects when user is not authenticated", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequireAuth()
		h := middleware(handler)
		err := h(c)

		assert.Error(t, err)
		// Check it's an unauthorized error
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})
}

func TestOptionalAuth(t *testing.T) {
	t.Run("passes when user is authenticated", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			userID := GetUserID(c)
			assert.Equal(t, "user-123", userID)
			return c.String(http.StatusOK, "OK")
		}

		middleware := OptionalAuth()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("passes when user is not authenticated", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			userID := GetUserID(c)
			assert.Equal(t, "", userID)
			return c.String(http.StatusOK, "OK")
		}

		middleware := OptionalAuth()
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})
}

func TestRequireRole(t *testing.T) {
	t.Run("passes when user has required role", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		req.Header.Set("X-User-Role", "artist")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequireRole(models.RoleArtist)
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("passes when user has admin role (admin can do anything)", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		req.Header.Set("X-User-Role", "admin")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequireRole(models.RoleArtist)
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("rejects when user does not have required role", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		req.Header.Set("X-User-Role", "subscriber")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequireRole(models.RoleArtist)
		h := middleware(handler)
		err := h(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("rejects when user is not authenticated", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequireRole(models.RoleArtist)
		h := middleware(handler)
		err := h(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})
}

func TestRequirePermission(t *testing.T) {
	t.Run("passes when user has required permission", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		req.Header.Set("X-User-Role", "artist")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlerCalled := false
		handler := func(c echo.Context) error {
			handlerCalled = true
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequirePermission(models.PermissionUploadTracks)
		h := middleware(handler)
		err := h(c)

		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("rejects when user does not have required permission", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		req.Header.Set("X-User-Role", "subscriber")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		}

		middleware := RequirePermission(models.PermissionUploadTracks)
		h := middleware(handler)
		err := h(c)

		assert.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("returns user ID from context", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "user-123")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Simulate middleware setting the context
		c.Set(UserIDKey, "user-123")

		userID := GetUserID(c)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("returns empty string when not set", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		userID := GetUserID(c)
		assert.Equal(t, "", userID)
	})
}

func TestGetUserRole(t *testing.T) {
	t.Run("returns user role from context", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set(UserRoleKey, models.RoleArtist)

		role := GetUserRole(c)
		assert.Equal(t, models.RoleArtist, role)
	})

	t.Run("returns guest when not set", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		role := GetUserRole(c)
		assert.Equal(t, models.RoleGuest, role)
	})
}
