//go:build integration

package testutil

import (
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gvasels/personal-music-searchengine/internal/handlers"
	handlermw "github.com/gvasels/personal-music-searchengine/internal/handlers/middleware"
	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
)

// TestServerContext holds a running HTTP test server backed by LocalStack.
type TestServerContext struct {
	*TestContext                  // Embedded — DynamoDB, S3, Cognito clients for direct verification
	Server       *httptest.Server // Running HTTP test server
	BaseURL      string           // e.g. "http://127.0.0.1:54321"
	Echo         *echo.Echo       // Echo instance
	Services     *service.Services
}

// customValidator implements echo.Validator (mirrors cmd/api/validator.go)
type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// SetupTestServer creates a full Echo HTTP server backed by LocalStack.
// The server has all routes, middleware, and services wired identically to production.
// Returns TestServerContext and cleanup function.
// Skips the test if LocalStack is not running.
func SetupTestServer(t *testing.T) (*TestServerContext, func()) {
	t.Helper()

	// Reuse existing LocalStack setup for clients
	tc, tcCleanup := SetupLocalStack(t)

	// Create repositories using real LocalStack clients
	repo := repository.NewDynamoDBRepository(tc.DynamoDB, tc.TableName)
	presignClient := s3.NewPresignClient(tc.S3)
	s3Repo := repository.NewS3Repository(tc.S3, presignClient, tc.BucketName)

	// Create services (same wiring as cmd/api/main.go setupEcho)
	services := service.NewServices(
		repo,
		s3Repo,
		nil, // No CloudFront signer in tests
		tc.BucketName,
		"", // No Step Functions ARN in tests
	)

	// Wire admin service if Cognito is available
	if tc.UserPoolID != "" {
		cognitoSvc := service.NewCognitoClient(tc.Cognito, tc.UserPoolID)
		services.Admin = service.NewAdminService(repo, cognitoSvc)
	}

	// Create handlers
	h := handlers.NewHandlers(services)

	// Create Echo instance (mirrors cmd/api/main.go setupEcho)
	e := echo.New()
	e.HideBanner = true
	e.Validator = &customValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Recover())

	// Register routes
	h.RegisterRoutes(e)

	// Register admin routes if admin service is configured
	if services.Admin != nil {
		adminHandler := handlers.NewAdminHandler(services.Admin)
		roleResolver := services.User.GetUserRole
		handlers.RegisterAdminRoutes(e, adminHandler, roleResolver)
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Follow system routes (not in RegisterRoutes yet - wire manually)
	registerFollowRoutes(e, repo)

	// Start test server
	server := httptest.NewServer(e)

	tsc := &TestServerContext{
		TestContext: tc,
		Server:     server,
		BaseURL:    server.URL,
		Echo:       e,
		Services:   services,
	}

	cleanup := func() {
		server.Close()
		tcCleanup()
	}

	return tsc, cleanup
}

// registerFollowRoutes adds follow-related routes.
// Follow/ArtistProfile handlers are not yet wired in RegisterRoutes,
// so we create the service and handler here for integration tests.
func registerFollowRoutes(e *echo.Echo, repo *repository.DynamoDBRepository) {
	followSvc := service.NewFollowService(repo)
	followHandler := handlers.NewFollowHandler(followSvc)

	// Follow routes under /api/v1/artists/entity/:id with auth middleware
	artists := e.Group("/api/v1/artists/entity/:id", handlermw.RequireAuth())
	artists.POST("/follow", followHandler.Follow)
	artists.DELETE("/follow", followHandler.Unfollow)
	artists.GET("/following", followHandler.IsFollowing)
	artists.GET("/followers", followHandler.GetFollowers)
}

// RequireRoleForTest is a test helper middleware that checks X-User-Role header.
// This mirrors the production RequireRoleWithDBCheck but uses the header fallback.
func RequireRoleForTest(requiredRole models.UserRole) echo.MiddlewareFunc {
	return handlermw.RequireRoleWithDBCheck(requiredRole, nil)
}
