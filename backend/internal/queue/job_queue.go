// backend/internal/queue/job_queue.go
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/your-org/esms/internal/cache"
)

// テスト用にオーバーライド可能な関数
var (
	NowFunc  = time.Now
	UUIDFunc = func() string { return uuid.New().String() }
)

// Job はキューに入れられるジョブの構造体
type Job struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	CreatedAt  time.Time              `json:"created_at"`
	RetryCount int                    `json:"retry_count"`
}

// JobQueue はジョブキューのインターフェース
type JobQueue interface {
	Enqueue(ctx context.Context, jobType string, payload map[string]interface{}) (string, error)
	Dequeue(ctx context.Context) (*Job, error)
}

// RedisJobQueue はRedisを使用したジョブキューの実装
type RedisJobQueue struct {
	client   *cache.RedisClient
	queueKey string
}

// NewRedisJobQueue は新しいRedisJobQueueを作成します
func NewRedisJobQueue(client *cache.RedisClient, queueName string) *RedisJobQueue {
	return &RedisJobQueue{
		client:   client,
		queueKey: fmt.Sprintf("queue:%s", queueName),
	}
}

// Enqueue はジョブをキューに追加します
func (q *RedisJobQueue) Enqueue(ctx context.Context, jobType string, payload map[string]interface{}) (string, error) {
	job := Job{
		ID:         UUIDFunc(),
		Type:       jobType,
		Payload:    payload,
		CreatedAt:  NowFunc(),
		RetryCount: 0,
	}

	jsonBytes, err := json.Marshal(job)
	if err != nil {
		return "", fmt.Errorf("failed to marshal job: %w", err)
	}

	// LPUSH でリストの先頭に追加
	if err := q.client.Client.LPush(ctx, q.queueKey, jsonBytes).Err(); err != nil {
		return "", fmt.Errorf("failed to enqueue job: %w", err)
	}

	return job.ID, nil
}

// Dequeue はジョブをキューから取り出します（ブロッキングなし）
// ジョブがない場合は nil, nil を返します
func (q *RedisJobQueue) Dequeue(ctx context.Context) (*Job, error) {
	// RPOP でリストの末尾から取得
	val, err := q.client.Client.RPop(ctx, q.queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // キューが空
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(val), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}
