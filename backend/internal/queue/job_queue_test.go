// backend/internal/queue/job_queue_test.go
package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/cache"
	"github.com/your-org/esms/internal/queue"
)

func TestRedisJobQueue_EnqueueDequeue(t *testing.T) {
	// テスト用に時刻とIDを固定
	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	fixedID := "test-job-id"

	originalNow := queue.NowFunc
	originalUUID := queue.UUIDFunc
	defer func() {
		queue.NowFunc = originalNow
		queue.UUIDFunc = originalUUID
	}()

	queue.NowFunc = func() time.Time { return fixedTime }
	queue.UUIDFunc = func() string { return fixedID }

	db, mock := redismock.NewClientMock()
	client := &cache.RedisClient{Client: db}
	q := queue.NewRedisJobQueue(client, "test-queue")
	ctx := context.Background()

	jobType := "SEND_EMAIL"
	payload := map[string]interface{}{"to": "test@example.com"}
	queueKey := "queue:test-queue"

	// Enqueueのモック
	// 固定値を使用するため、正確なJSONを期待できる
	expectedJSON := `{"id":"test-job-id","type":"SEND_EMAIL","payload":{"to":"test@example.com"},"created_at":"2025-01-01T00:00:00Z","retry_count":0}`
	mock.ExpectLPush(queueKey, []byte(expectedJSON)).SetVal(1)

	// テスト実行: Enqueue
	jobID, err := q.Enqueue(ctx, jobType, payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, jobID)

	// Dequeueのモック
	mock.ExpectRPop(queueKey).SetVal(expectedJSON)

	// テスト実行: Dequeue
	job, err := q.Dequeue(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, job)
	assert.Equal(t, fixedID, job.ID)
	assert.Equal(t, jobType, job.Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisJobQueue_Dequeue_Empty(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &cache.RedisClient{Client: db}
	q := queue.NewRedisJobQueue(client, "test-queue")
	ctx := context.Background()
	queueKey := "queue:test-queue"

	mock.ExpectRPop(queueKey).RedisNil()

	job, err := q.Dequeue(ctx)
	assert.NoError(t, err)
	assert.Nil(t, job)

	assert.NoError(t, mock.ExpectationsWereMet())
}
