//go:build integration

package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gvasels/personal-music-searchengine/internal/models"
	"github.com/stretchr/testify/require"
)

// RequestOption configures an HTTP request before it is sent.
type RequestOption func(*http.Request)

// AsUser sets X-User-ID and X-User-Role headers to simulate an authenticated user.
// The auth middleware falls back to these headers when API Gateway JWT claims are absent.
func AsUser(userID string, role models.UserRole) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("X-User-ID", userID)
		req.Header.Set("X-User-Role", string(role))
		// Admin users get global access
		if role == models.RoleAdmin {
			req.Header.Set("X-Global-Access", "true")
		}
	}
}

// WithJSON sets the Content-Type to application/json and marshals the body.
func WithJSON(body interface{}) RequestOption {
	return func(req *http.Request) {
		data, err := json.Marshal(body)
		if err != nil {
			panic("WithJSON: failed to marshal body: " + err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(bytes.NewReader(data))
		req.ContentLength = int64(len(data))
	}
}

// WithHeader sets a custom header on the request.
func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

// WithQuery adds a query parameter to the request URL.
func WithQuery(key, value string) RequestOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
	}
}

// DoRequest makes an HTTP request to the test server and returns the response.
// The path should start with "/" (e.g., "/api/v1/tracks").
func (tsc *TestServerContext) DoRequest(t *testing.T, method, path string, opts ...RequestOption) *http.Response {
	t.Helper()

	url := tsc.BaseURL + path
	req, err := http.NewRequest(method, url, nil)
	require.NoError(t, err, "Failed to create request: %s %s", method, path)

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "Failed to execute request: %s %s", method, path)

	return resp
}

// AssertStatus verifies the response has the expected status code.
// On failure, it dumps the response body for debugging.
func AssertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(body)) // Re-wrap for further reads
		t.Fatalf("Expected status %d, got %d. Body: %s", expected, resp.StatusCode, string(body))
	}
}

// DecodeJSON reads and unmarshals the response body into the target type.
// Fails the test if unmarshaling fails.
func DecodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	defer resp.Body.Close()
	var result T
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode JSON response")
	return result
}

// DecodeJSONBody reads the response body as a map for flexible assertions.
func DecodeJSONBody(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()

	defer resp.Body.Close()
	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode JSON response body")
	return result
}
