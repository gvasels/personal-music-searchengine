package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// AdminTestValidator implements echo.Validator for testing
type AdminTestValidator struct {
	validator *validator.Validate
}

func (v *AdminTestValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func setupAdminTestEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &AdminTestValidator{validator: validator.New()}
	return e
}

// MockAdminService is a mock implementation of service.AdminService
type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) SearchUsers(ctx context.Context, query string, limit int) ([]models.UserSummary, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSummary), args.Error(1)
}

func (m *MockAdminService) GetUserDetails(ctx context.Context, userID string) (*models.UserDetails, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserDetails), args.Error(1)
}

func (m *MockAdminService) UpdateUserRole(ctx context.Context, userID string, newRole models.UserRole) error {
	args := m.Called(ctx, userID, newRole)
	return args.Error(0)
}

func (m *MockAdminService) UpdateUserRoleByAdmin(ctx context.Context, adminID, userID string, newRole models.UserRole) error {
	args := m.Called(ctx, adminID, userID, newRole)
	return args.Error(0)
}

func (m *MockAdminService) SetUserStatus(ctx context.Context, userID string, disabled bool) error {
	args := m.Called(ctx, userID, disabled)
	return args.Error(0)
}

func resetAdminMock(m *MockAdminService) {
	m.ExpectedCalls = nil
	m.Calls = nil
}

func TestAdminHandler_SearchUsers(t *testing.T) {
	e := setupAdminTestEcho()
	mockService := new(MockAdminService)
	handler := NewAdminHandler(mockService)

	now := time.Now()

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:        "successful search",
			queryParams: "?search=test",
			setupMock: func() {
				mockService.On("SearchUsers", mock.Anything, "test", 20).Return(
					[]models.UserSummary{
						{ID: "user-1", Email: "user1@example.com", DisplayName: "User One", Role: models.RoleSubscriber, Disabled: false, CreatedAt: now},
						{ID: "user-2", Email: "user2@example.com", DisplayName: "User Two", Role: models.RoleAdmin, Disabled: false, CreatedAt: now},
					},
					nil,
				).Once()
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var response models.AdminSearchUsersResponse
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Len(t, response.Items, 2)
			},
		},
		{
			name:        "search with custom limit",
			queryParams: "?search=admin&limit=10",
			setupMock: func() {
				mockService.On("SearchUsers", mock.Anything, "admin", 10).Return(
					[]models.UserSummary{{ID: "user-2", Email: "admin@example.com", DisplayName: "Admin", Role: models.RoleAdmin, Disabled: false, CreatedAt: now}},
					nil,
				).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing search query",
			queryParams:    "",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "service error",
			queryParams: "?search=test",
			setupMock: func() {
				mockService.On("SearchUsers", mock.Anything, "test", 20).Return(nil, errors.New("database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAdminMock(mockService)
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users"+tt.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.SearchUsers(c)
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.checkBody != nil {
					tt.checkBody(t, rec.Body.String())
				}
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestAdminHandler_GetUserDetails(t *testing.T) {
	e := setupAdminTestEcho()
	mockService := new(MockAdminService)
	handler := NewAdminHandler(mockService)

	now := time.Now()

	tests := []struct {
		name           string
		userID         string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "successful retrieval",
			userID: "user-123",
			setupMock: func() {
				mockService.On("GetUserDetails", mock.Anything, "user-123").Return(
					&models.UserDetails{
						UserSummary: models.UserSummary{
							ID:          "user-123",
							Email:       "test@example.com",
							DisplayName: "Test User",
							Role:        models.RoleSubscriber,
							Disabled:    false,
							CreatedAt:   now,
						},
						TrackCount:    42,
						PlaylistCount: 5,
						StorageUsed:   1073741824,
					},
					nil,
				).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMock: func() {
				mockService.On("GetUserDetails", mock.Anything, "nonexistent").Return(nil, models.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing user ID",
			userID:         "",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAdminMock(mockService)
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			err := handler.GetUserDetails(c)
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestAdminHandler_UpdateUserRole(t *testing.T) {
	e := setupAdminTestEcho()
	mockService := new(MockAdminService)
	handler := NewAdminHandler(mockService)

	now := time.Now()
	mockUserDetails := &models.UserDetails{
		UserSummary: models.UserSummary{
			ID:          "user-123",
			Email:       "test@example.com",
			DisplayName: "Test User",
			Role:        models.RoleAdmin,
			Disabled:    false,
			CreatedAt:   now,
		},
		TrackCount:    42,
		PlaylistCount: 5,
		StorageUsed:   1073741824,
	}

	tests := []struct {
		name           string
		userID         string
		adminID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:        "successful role update",
			userID:      "user-123",
			adminID:     "admin-1",
			requestBody: `{"role": "admin"}`,
			setupMock: func() {
				mockService.On("UpdateUserRoleByAdmin", mock.Anything, "admin-1", "user-123", models.RoleAdmin).Return(nil).Once()
				mockService.On("GetUserDetails", mock.Anything, "user-123").Return(mockUserDetails, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid role",
			userID:         "user-123",
			adminID:        "admin-1",
			requestBody:    `{"role": "superadmin"}`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "user not found",
			userID:      "nonexistent",
			adminID:     "admin-1",
			requestBody: `{"role": "admin"}`,
			setupMock: func() {
				mockService.On("UpdateUserRoleByAdmin", mock.Anything, "admin-1", "nonexistent", models.RoleAdmin).Return(models.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing admin auth",
			userID:         "user-123",
			adminID:        "",
			requestBody:    `{"role": "admin"}`,
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAdminMock(mockService)
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/"+tt.userID+"/role", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)
			// Set admin user ID in context
			if tt.adminID != "" {
				c.Set("user_id", tt.adminID)
			}

			err := handler.UpdateUserRole(c)
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestAdminHandler_UpdateUserStatus(t *testing.T) {
	e := setupAdminTestEcho()
	mockService := new(MockAdminService)
	handler := NewAdminHandler(mockService)

	now := time.Now()
	mockDisabledDetails := &models.UserDetails{
		UserSummary: models.UserSummary{
			ID:          "user-123",
			Email:       "test@example.com",
			DisplayName: "Test User",
			Role:        models.RoleSubscriber,
			Disabled:    true,
			CreatedAt:   now,
		},
		TrackCount:    42,
		PlaylistCount: 5,
		StorageUsed:   1073741824,
	}
	mockEnabledDetails := &models.UserDetails{
		UserSummary: models.UserSummary{
			ID:          "user-123",
			Email:       "test@example.com",
			DisplayName: "Test User",
			Role:        models.RoleSubscriber,
			Disabled:    false,
			CreatedAt:   now,
		},
		TrackCount:    42,
		PlaylistCount: 5,
		StorageUsed:   1073741824,
	}

	tests := []struct {
		name           string
		userID         string
		adminID        string
		requestBody    string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:        "successful disable",
			userID:      "user-123",
			adminID:     "admin-1",
			requestBody: `{"disabled": true}`,
			setupMock: func() {
				mockService.On("SetUserStatus", mock.Anything, "user-123", true).Return(nil).Once()
				mockService.On("GetUserDetails", mock.Anything, "user-123").Return(mockDisabledDetails, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "successful enable",
			userID:      "user-123",
			adminID:     "admin-1",
			requestBody: `{"disabled": false}`,
			setupMock: func() {
				mockService.On("SetUserStatus", mock.Anything, "user-123", false).Return(nil).Once()
				mockService.On("GetUserDetails", mock.Anything, "user-123").Return(mockEnabledDetails, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "cannot disable self",
			userID:         "admin-1",
			adminID:        "admin-1",
			requestBody:    `{"disabled": true}`,
			setupMock:      func() {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "user not found",
			userID:      "nonexistent",
			adminID:     "admin-1",
			requestBody: `{"disabled": true}`,
			setupMock: func() {
				mockService.On("SetUserStatus", mock.Anything, "nonexistent", true).Return(models.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAdminMock(mockService)
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/"+tt.userID+"/status", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)
			if tt.adminID != "" {
				c.Set("user_id", tt.adminID)
			}

			err := handler.UpdateUserStatus(c)
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestNewAdminHandler(t *testing.T) {
	mockService := new(MockAdminService)
	handler := NewAdminHandler(mockService)
	assert.NotNil(t, handler)
}
