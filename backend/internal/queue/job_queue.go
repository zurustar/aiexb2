// backend/internal/queue/job_queue.go
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
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
	MaxRetries int                    `json:"max_retries"`
	LastError  string                 `json:"last_error,omitempty"`
}

// JobQueue はジョブキューのインターフェース
type JobQueue interface {
	Enqueue(ctx context.Context, jobType string, payload map[string]interface{}) (string, error)
	Dequeue(ctx context.Context) (*Job, error)
	Retry(ctx context.Context, job *Job, err error) error
	MoveToDLQ(ctx context.Context, job *Job) error
}

// RedisJobQueue はRedisを使用したジョブキューの実装
type RedisJobQueue struct {
	client     *cache.RedisClient
	queueKey   string
	dlqKey     string
	maxRetries int
}

// NewRedisJobQueue は新しいRedisJobQueueを作成します
func NewRedisJobQueue(client *cache.RedisClient, queueName string) *RedisJobQueue {
	return &RedisJobQueue{
		client:     client,
		queueKey:   fmt.Sprintf("queue:%s", queueName),
		dlqKey:     fmt.Sprintf("queue:%s:dlq", queueName),
		maxRetries: 3,
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
		MaxRetries: q.maxRetries,
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

// Retry はジョブを指数バックオフとジッターでリトライします
func (q *RedisJobQueue) Retry(ctx context.Context, job *Job, err error) error {
	job.RetryCount++
	job.LastError = err.Error()

	// 最大リトライ回数を超えた場合はDLQへ移動
	if job.RetryCount >= job.MaxRetries {
		return q.MoveToDLQ(ctx, job)
	}

	// 指数バックオフ + ジッター計算
	baseDelay := time.Second
	backoff := time.Duration(math.Pow(2, float64(job.RetryCount))) * baseDelay
	jitter := time.Duration(rand.Int63n(int64(baseDelay)))
	delay := backoff + jitter

	// 遅延後に再エンキュー（簡易実装: 即座にエンキューし、ワーカー側で遅延処理）
	jsonBytes, marshalErr := json.Marshal(job)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal job for retry: %w", marshalErr)
	}

	// 遅延キューへの追加（ZADD with score = current_time + delay）
	score := float64(NowFunc().Add(delay).Unix())
	if redisErr := q.client.Client.ZAdd(ctx, q.queueKey+":delayed", redis.Z{
		Score:  score,
		Member: jsonBytes,
	}).Err(); redisErr != nil {
		return fmt.Errorf("failed to retry job: %w", redisErr)
	}

	return nil
}

// MoveToDLQ はジョブをDead Letter Queueへ移動します
func (q *RedisJobQueue) MoveToDLQ(ctx context.Context, job *Job) error {
	jsonBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job for DLQ: %w", err)
	}

	if err := q.client.Client.LPush(ctx, q.dlqKey, jsonBytes).Err(); err != nil {
		return fmt.Errorf("failed to move job to DLQ: %w", err)
	}

	return nil
}
