// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

var (
	apiURL       = os.Getenv("API_URL")
	userPoolID   = os.Getenv("COGNITO_USER_POOL_ID")
	clientID     = os.Getenv("COGNITO_CLIENT_ID")
	testEmail    = os.Getenv("TEST_USER_EMAIL")
	testPassword = os.Getenv("TEST_USER_PASSWORD")
	awsRegion    = "us-east-1"
)

// skipIfLocalStack skips tests that require a real deployed AWS environment.
// These tests hit real API Gateway, CloudFront, Cognito, and DynamoDB endpoints
// and cannot run in LocalStack mode.
func skipIfLocalStack(t *testing.T) {
	t.Helper()
	if os.Getenv("LOCALSTACK_ENDPOINT") != "" {
		t.Skip("Skipping deployed-environment test (LOCALSTACK_ENDPOINT is set)")
	}
}

// TestConfig holds the test configuration
type TestConfig struct {
	APIBaseURL string
	AuthToken  string
	HTTPClient *http.Client
}

var testConfig *TestConfig

// setupTestConfig initializes the test configuration
func setupTestConfig(t *testing.T) *TestConfig {
	if testConfig != nil {
		return testConfig
	}

	if apiURL == "" {
		apiURL = "https://r1simytb2i.execute-api.us-east-1.amazonaws.com"
	}

	testConfig = &TestConfig{
		APIBaseURL: apiURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}

	return testConfig
}

// ============================================================================
// EPIC 1: Foundation Infrastructure Tests
// ============================================================================

func TestEpic1_HealthEndpoint(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	resp, err := cfg.HTTPClient.Get(cfg.APIBaseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestEpic1_DynamoDBConnectivity(t *testing.T) {
	skipIfLocalStack(t)
	ctx := context.Background()

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	require.NoError(t, err)

	client := dynamodb.NewFromConfig(awsCfg)

	// List tables to verify connectivity
	result, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	require.NoError(t, err)

	// Check that our table exists
	found := false
	for _, name := range result.TableNames {
		if strings.Contains(name, "music-library") {
			found = true
			break
		}
	}
	assert.True(t, found, "music-library table should exist")
}

func TestEpic1_S3Connectivity(t *testing.T) {
	skipIfLocalStack(t)
	ctx := context.Background()

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	require.NoError(t, err)

	client := s3.NewFromConfig(awsCfg)

	// List buckets to verify connectivity
	result, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	require.NoError(t, err)

	// Check that our media bucket exists
	found := false
	for _, bucket := range result.Buckets {
		if bucket.Name != nil && strings.Contains(*bucket.Name, "music-library") && strings.Contains(*bucket.Name, "media") {
			found = true
			break
		}
	}
	assert.True(t, found, "music-library media bucket should exist")
}

func TestEpic1_CognitoConnectivity(t *testing.T) {
	skipIfLocalStack(t)
	if userPoolID == "" {
		userPoolID = "us-east-1_edxzNkQHw"
	}

	ctx := context.Background()

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	require.NoError(t, err)

	client := cognitoidentityprovider.NewFromConfig(awsCfg)

	// Describe user pool to verify connectivity
	result, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(userPoolID),
	})
	require.NoError(t, err)
	assert.NotNil(t, result.UserPool)
	assert.Contains(t, *result.UserPool.Name, "music-library")
}

// ============================================================================
// EPIC 2: Backend API Tests
// ============================================================================

func TestEpic2_UnauthorizedAccess(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	// Protected endpoints should return 401 without auth
	endpoints := []string{
		"/api/v1/me",
		"/api/v1/tracks",
		"/api/v1/albums",
		"/api/v1/playlists",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			resp, err := cfg.HTTPClient.Get(cfg.APIBaseURL + endpoint)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
				"Endpoint %s should require authentication", endpoint)
		})
	}
}

func TestEpic2_APIRoutes(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	// Test that routes exist (OPTIONS should work)
	routes := []string{
		"/api/v1/tracks",
		"/api/v1/albums",
		"/api/v1/artists",
		"/api/v1/playlists",
		"/api/v1/tags",
		"/api/v1/uploads",
	}

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			req, err := http.NewRequest("OPTIONS", cfg.APIBaseURL+route, nil)
			require.NoError(t, err)

			resp, err := cfg.HTTPClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// OPTIONS should return 200 or 204
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent,
				"OPTIONS %s should succeed, got %d", route, resp.StatusCode)
		})
	}
}

// ============================================================================
// EPIC 3: Search & Streaming Tests
// ============================================================================

func TestEpic3_SearchEndpointExists(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	// Search endpoint should exist (returns 401 without auth)
	resp, err := cfg.HTTPClient.Get(cfg.APIBaseURL + "/api/v1/search?q=test")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should be 401 (unauthorized) not 404 (not found)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
		"Search endpoint should exist but require auth")
}

func TestEpic3_StreamEndpointExists(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	// Stream endpoint should exist
	resp, err := cfg.HTTPClient.Get(cfg.APIBaseURL + "/api/v1/stream/test-track-id")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should be 401 (unauthorized) not 404 (not found)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
		"Stream endpoint should exist but require auth")
}

func TestEpic3_CloudFrontDistribution(t *testing.T) {
	skipIfLocalStack(t)
	// Test CloudFront is accessible
	cloudfrontDomain := "d3039j0t382r36.cloudfront.net"

	resp, err := http.Get("https://" + cloudfrontDomain)
	require.NoError(t, err)
	defer resp.Body.Close()

	// CloudFront should respond (even if 403 for restricted content)
	assert.True(t, resp.StatusCode != 0, "CloudFront should be accessible")
}

// ============================================================================
// EPIC 4: Tags & Playlists Tests
// ============================================================================

func TestEpic4_TagsEndpointExists(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/tags"},
		{"POST", "/api/v1/tags"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			var resp *http.Response
			var err error

			if ep.method == "GET" {
				resp, err = cfg.HTTPClient.Get(cfg.APIBaseURL + ep.path)
			} else {
				resp, err = cfg.HTTPClient.Post(cfg.APIBaseURL+ep.path, "application/json", bytes.NewBuffer([]byte("{}")))
			}
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should be 401 (unauthorized) not 404 (not found)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
				"Tags endpoint should exist but require auth")
		})
	}
}

func TestEpic4_PlaylistsEndpointExists(t *testing.T) {
	skipIfLocalStack(t)
	cfg := setupTestConfig(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/playlists"},
		{"POST", "/api/v1/playlists"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			var resp *http.Response
			var err error

			if ep.method == "GET" {
				resp, err = cfg.HTTPClient.Get(cfg.APIBaseURL + ep.path)
			} else {
				resp, err = cfg.HTTPClient.Post(cfg.APIBaseURL+ep.path, "application/json", bytes.NewBuffer([]byte("{}")))
			}
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should be 401 (unauthorized) not 404 (not found)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
				"Playlists endpoint should exist but require auth")
		})
	}
}

// ============================================================================
// Authenticated Tests (require test user credentials)
// ============================================================================

func getAuthToken(t *testing.T) string {
	if testEmail == "" || testPassword == "" {
		t.Skip("TEST_USER_EMAIL and TEST_USER_PASSWORD required for authenticated tests")
	}

	if clientID == "" {
		clientID = "sfuu14g0un87omvt02lbdfgun"
	}

	ctx := context.Background()

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	require.NoError(t, err)

	client := cognitoidentityprovider.NewFromConfig(awsCfg)

	result, err := client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(clientID),
		AuthParameters: map[string]string{
			"USERNAME": testEmail,
			"PASSWORD": testPassword,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result.AuthenticationResult)

	return *result.AuthenticationResult.IdToken
}

func makeAuthenticatedRequest(t *testing.T, method, url string, body io.Reader) *http.Response {
	cfg := setupTestConfig(t)
	token := getAuthToken(t)

	req, err := http.NewRequest(method, cfg.APIBaseURL+url, body)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfg.HTTPClient.Do(req)
	require.NoError(t, err)

	return resp
}

func TestAuthenticated_GetProfile(t *testing.T) {
	skipIfLocalStack(t)
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/me", nil)
	defer resp.Body.Close()

	// Should get profile or create new user
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound,
		"Profile endpoint should work with auth, got %d", resp.StatusCode)
}

func TestAuthenticated_ListTracks(t *testing.T) {
	skipIfLocalStack(t)
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/tracks", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Should be able to list tracks with auth")

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Should have tracks array (even if empty)
	_, hasItems := result["items"]
	assert.True(t, hasItems, "Response should have items field")
}

func TestAuthenticated_ListPlaylists(t *testing.T) {
	skipIfLocalStack(t)
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/playlists", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Should be able to list playlists with auth")
}

func TestAuthenticated_ListTags(t *testing.T) {
	skipIfLocalStack(t)
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/tags", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Should be able to list tags with auth")
}

func TestAuthenticated_Search(t *testing.T) {
	skipIfLocalStack(t)
	resp := makeAuthenticatedRequest(t, "GET", "/api/v1/search?q=test", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Should be able to search with auth")
}

// ============================================================================
// Helper to print test summary
// ============================================================================

func TestMain(m *testing.M) {
	fmt.Println("==========================================")
	fmt.Println("Music Library Integration Tests")
	fmt.Println("==========================================")
	fmt.Printf("API URL: %s\n", apiURL)
	fmt.Printf("Region: %s\n", awsRegion)
	fmt.Println("==========================================")

	os.Exit(m.Run())
}
