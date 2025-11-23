// backend/tests/integration/server_test.go
package integration_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerStartup はAPIサーバーの起動テスト
func TestServerStartup(t *testing.T) {
	if testDB == nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	// Note: 実際のサーバー起動テストは、サーバープロセスを別途起動する必要がある
	// ここでは構造のみ示す
	t.Skip("Requires actual server process")
}

// TestHealthEndpoint はヘルスチェックエンドポイントのテスト
func TestHealthEndpoint(t *testing.T) {
	// サーバーが起動していることを前提とする
	serverURL := "http://localhost:8080"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/health", nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)

	// サーバーが起動していない場合はスキップ
	if err != nil {
		t.Skipf("Server not running: %v", err)
		return
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestAPIEndpoints は主要なAPIエンドポイントのテスト
func TestAPIEndpoints(t *testing.T) {
	serverURL := "http://localhost:8080"

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Health check",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Auth login (no auth)",
			method:         "GET",
			path:           "/api/v1/auth/login",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, tt.method, serverURL+tt.path, nil)
			require.NoError(t, err)

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)

			if err != nil {
				t.Skipf("Server not running: %v", err)
				return
			}
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
