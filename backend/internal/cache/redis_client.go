// backend/internal/cache/redis_client.go
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/your-org/esms/internal/config"
)

// ErrCacheMiss はキャッシュが存在しない場合のエラー
var ErrCacheMiss = errors.New("cache: key not found")

// RedisClient はRedis操作のラッパー
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient は新しいRedisクライアントを作成します
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr(),
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  cfg.RedisConnectTimeout,
		ReadTimeout:  cfg.RedisReadTimeout,
		WriteTimeout: cfg.RedisWriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.RedisConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

// Set は値をJSONとして保存します
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.Client.Set(ctx, key, jsonBytes, expiration).Err()
}

// Get は値をJSONとして取得し、destにデコードします
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss // 独自エラーを返す
		}
		return fmt.Errorf("failed to get value from redis: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	return nil
}

// Delete はキーを削除します
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// Close は接続を閉じます
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
