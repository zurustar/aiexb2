// backend/tests/performance/load_test.go
package performance_test

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLoadBasic は基本的な負荷テスト
func TestLoadBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	serverURL := "http://localhost:8080"

	// サーバーが起動しているか確認
	if !isServerRunning(serverURL) {
		t.Skip("Server not running, skipping load test")
	}

	concurrentUsers := 10
	requestsPerUser := 100

	results := runLoadTest(t, serverURL, concurrentUsers, requestsPerUser)

	// 結果の検証
	assert.Greater(t, results.SuccessCount, 0, "Should have successful requests")
	assert.Less(t, results.ErrorCount, results.TotalRequests/10, "Error rate should be less than 10%")
	assert.Less(t, results.AverageResponseTime, 500*time.Millisecond, "Average response time should be less than 500ms")

	t.Logf("Load Test Results:")
	t.Logf("  Total Requests: %d", results.TotalRequests)
	t.Logf("  Successful: %d", results.SuccessCount)
	t.Logf("  Failed: %d", results.ErrorCount)
	t.Logf("  Average Response Time: %v", results.AverageResponseTime)
	t.Logf("  Min Response Time: %v", results.MinResponseTime)
	t.Logf("  Max Response Time: %v", results.MaxResponseTime)
}

// TestResponseTime はレスポンスタイムの測定
func TestResponseTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping response time test in short mode")
	}

	serverURL := "http://localhost:8080"

	if !isServerRunning(serverURL) {
		t.Skip("Server not running, skipping response time test")
	}

	endpoints := []string{
		"/health",
		"/api/v1/auth/login",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			responseTimes := make([]time.Duration, 0, 100)

			for i := 0; i < 100; i++ {
				start := time.Now()
				resp, err := http.Get(serverURL + endpoint)
				duration := time.Since(start)

				if err == nil {
					resp.Body.Close()
					responseTimes = append(responseTimes, duration)
				}
			}

			if len(responseTimes) > 0 {
				avg := averageDuration(responseTimes)
				t.Logf("Endpoint %s - Average response time: %v", endpoint, avg)
				assert.Less(t, avg, 200*time.Millisecond, "Response time should be less than 200ms")
			}
		})
	}
}

// LoadTestResults は負荷テストの結果
type LoadTestResults struct {
	TotalRequests       int
	SuccessCount        int
	ErrorCount          int
	AverageResponseTime time.Duration
	MinResponseTime     time.Duration
	MaxResponseTime     time.Duration
}

// runLoadTest は負荷テストを実行します
func runLoadTest(t *testing.T, serverURL string, concurrentUsers, requestsPerUser int) *LoadTestResults {
	var wg sync.WaitGroup
	var mu sync.Mutex

	responseTimes := make([]time.Duration, 0)
	successCount := 0
	errorCount := 0

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			client := &http.Client{Timeout: 10 * time.Second}

			for j := 0; j < requestsPerUser; j++ {
				start := time.Now()

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/health", nil)

				if err != nil {
					cancel()
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}

				resp, err := client.Do(req)
				duration := time.Since(start)
				cancel()

				mu.Lock()
				if err != nil || resp.StatusCode != http.StatusOK {
					errorCount++
				} else {
					successCount++
					responseTimes = append(responseTimes, duration)
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()

	results := &LoadTestResults{
		TotalRequests: concurrentUsers * requestsPerUser,
		SuccessCount:  successCount,
		ErrorCount:    errorCount,
	}

	if len(responseTimes) > 0 {
		results.AverageResponseTime = averageDuration(responseTimes)
		results.MinResponseTime = minDuration(responseTimes)
		results.MaxResponseTime = maxDuration(responseTimes)
	}

	return results
}

// isServerRunning はサーバーが起動しているか確認します
func isServerRunning(serverURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/health", nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// averageDuration は平均時間を計算します
func averageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

// minDuration は最小時間を返します
func minDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	min := durations[0]
	for _, d := range durations[1:] {
		if d < min {
			min = d
		}
	}
	return min
}

// maxDuration は最大時間を返します
func maxDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	max := durations[0]
	for _, d := range durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

// BenchmarkHealthEndpoint はヘルスチェックエンドポイントのベンチマーク
func BenchmarkHealthEndpoint(b *testing.B) {
	serverURL := "http://localhost:8080"

	if !isServerRunning(serverURL) {
		b.Skip("Server not running")
	}

	client := &http.Client{Timeout: 5 * time.Second}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(serverURL + "/health")
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}
