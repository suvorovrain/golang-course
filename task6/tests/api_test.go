package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const address = "http://localhost:28080"

var client = http.Client{
	Timeout: 30 * time.Second,
}

type PingService struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string        `json:"status"`
	Services []PingService `json:"services"`
}

type RepositoryInfoResponse struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func waitForAPI(t *testing.T) {
	t.Helper()

	require.Eventually(t, func() bool {
		resp, err := client.Get(address + "/api/ping")
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusServiceUnavailable
	}, 20*time.Second, 500*time.Millisecond, "api did not become ready")
}

func waitForRepositoryInfo(t *testing.T, url string) *RepositoryInfoResponse {
	t.Helper()

	var body RepositoryInfoResponse
	require.Eventually(t, func() bool {
		resp, err := client.Get(url)
		if err != nil {
			t.Logf("Request failed: %v", err)
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Logf("Unexpected status: %d", resp.StatusCode)
			return false
		}

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Logf("Decode failed: %v", err)
			return false
		}

		return body.FullName != "" && body.CreatedAt != ""
	}, 30*time.Second, 500*time.Millisecond, "repository info did not become available")

	return &body
}

func serviceMap(services []PingService) map[string]string {
	res := make(map[string]string, len(services))
	for _, svc := range services {
		res[svc.Name] = svc.Status
	}
	return res
}

func TestPreflight(t *testing.T) {
	require.Equal(t, true, true)
}

func TestPing(t *testing.T) {
	waitForAPI(t)

	resp, err := client.Get(address + "/api/ping")
	require.NoError(t, err, "cannot ping api")
	defer resp.Body.Close()

	var body PingResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body), "cannot decode ping response")

	require.Equal(t, http.StatusOK, resp.StatusCode, "wrong status code")
	require.Equal(t, "ok", body.Status, "wrong overall status")

	services := serviceMap(body.Services)
	require.Equal(t, "up", services["processor"], "processor should be up")
	require.Equal(t, "up", services["subscriber"], "subscriber should be up")
}

func TestRepositoryInfo(t *testing.T) {
	waitForAPI(t)

	url := address + "/api/repositories/info?url=https://github.com/golang/go"

	// Make initial request to trigger Kafka processing
	resp, err := client.Get(url)
	require.NoError(t, err, "cannot request repository info")
	resp.Body.Close()

	// Wait for async processing to complete
	body := waitForRepositoryInfo(t, url)

	require.Equal(t, "golang/go", body.FullName, "wrong full_name")
	require.NotEmpty(t, body.Description, "description should not be empty")
	require.NotEmpty(t, body.CreatedAt, "created_at should not be empty")
	require.True(t, body.Stars > 0, "stars should be positive")
	require.True(t, body.Forks > 0, "forks should be positive")
}

func TestRepositoryInfoWithoutURL(t *testing.T) {
	waitForAPI(t)

	resp, err := client.Get(address + "/api/repositories/info")
	require.NoError(t, err, "cannot request repository info without url")
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "wrong status code")

	var body ErrorResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body), "cannot decode error response")
	require.NotEmpty(t, body.Error, "error message should not be empty")
}

func TestRepositoryInfoInvalidURL(t *testing.T) {
	waitForAPI(t)

	resp, err := client.Get(address + "/api/repositories/info?url=not-a-valid-url")
	require.NoError(t, err, "cannot request repository info with invalid url")
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "wrong status code")

	var body ErrorResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body), "cannot decode error response")
	require.NotEmpty(t, body.Error, "error message should not be empty")
}

func TestPingHelpfulFailureMessage(t *testing.T) {
	waitForAPI(t)

	resp, err := client.Get(address + "/api/ping")
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		var body PingResponse
		_ = json.NewDecoder(resp.Body).Decode(&body)

		services := serviceMap(body.Services)
		t.Fatalf("api is degraded: processor=%s subscriber=%s", services["processor"], services["subscriber"])
	}

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRepositoryInfoHelpfulFailureMessage(t *testing.T) {
	waitForAPI(t)

	url := address + "/api/repositories/info?url=https://github.com/golang/go"

	// Trigger request
	resp, err := client.Get(url)
	require.NoError(t, err)
	resp.Body.Close()

	// Wait for data to be available
	body := waitForRepositoryInfo(t, url)

	require.NotEmpty(t, body.FullName)
	require.NotEmpty(t, body.CreatedAt)
}

func TestRepositoryInfoCreatedAtFormatPresent(t *testing.T) {
	waitForAPI(t)

	url := address + "/api/repositories/info?url=https://github.com/golang/go"

	// Trigger request
	resp, err := client.Get(url)
	require.NoError(t, err)
	resp.Body.Close()

	// Wait for async processing
	body := waitForRepositoryInfo(t, url)

	require.Contains(t, body.CreatedAt, "T", "created_at should look like RFC3339 timestamp")
	require.True(t, strings.Contains(body.CreatedAt, "T") || strings.Contains(body.CreatedAt, "+"),
		"created_at should be in RFC3339 format")
}

func TestRepositoryInfoEndpointStable(t *testing.T) {
	waitForAPI(t)

	urls := []string{
		address + "/api/repositories/info?url=https://github.com/golang/go",
		address + "/api/repositories/info?url=https://github.com/kubernetes/kubernetes",
	}

	for _, u := range urls {
		t.Run(fmt.Sprintf("request_%s", u), func(t *testing.T) {
			// Trigger request to initiate async processing
			resp, err := client.Get(u)
			require.NoError(t, err)
			resp.Body.Close()

			// Wait for async processing to complete
			body := waitForRepositoryInfo(t, u)

			require.NotEmpty(t, body.FullName)
			require.NotEmpty(t, body.CreatedAt)
		})
	}
}
