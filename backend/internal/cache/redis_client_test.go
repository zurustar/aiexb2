// backend/internal/cache/redis_client_test.go
package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/cache"
)

type TestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRedisClient_SetGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &cache.RedisClient{Client: db}
	ctx := context.Background()

	key := "test-key"
	data := TestData{Name: "Alice", Age: 30}
	jsonStr := `{"name":"Alice","age":30}`
	expiration := 1 * time.Hour

	// Setのモック
	mock.ExpectSet(key, []byte(jsonStr), expiration).SetVal("OK")

	// Getのモック
	mock.ExpectGet(key).SetVal(jsonStr)

	// テスト実行: Set
	err := client.Set(ctx, key, data, expiration)
	assert.NoError(t, err)

	// テスト実行: Get
	var result TestData
	err = client.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, data, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisClient_Get_NotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &cache.RedisClient{Client: db}
	ctx := context.Background()

	key := "not-found"
	mock.ExpectGet(key).RedisNil()

	var result TestData
	err := client.Get(ctx, key, &result)
	assert.Error(t, err)
	// redis.Nil エラーが返ることを期待（ラップされている可能性も考慮）

	assert.NoError(t, mock.ExpectationsWereMet())
}
