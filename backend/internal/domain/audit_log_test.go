// backend/internal/domain/audit_log_test.go
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/domain"
)

func TestAuditLog_GenerateSignature(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	log := domain.AuditLog{
		ID:         id,
		UserID:     userID,
		Action:     domain.AuditActionCreate,
		Resource:   "reservations",
		ResourceID: "123",
		Details:    map[string]interface{}{"foo": "bar"},
		IPAddress:  "127.0.0.1",
		UserAgent:  "Go-Test-Client",
		CreatedAt:  now,
	}

	secretKey := "my-secret-key"

	// 署名生成テスト
	sig1, err := log.GenerateSignature(secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, sig1)

	// 同じデータなら同じ署名になるはず
	sig2, err := log.GenerateSignature(secretKey)
	assert.NoError(t, err)
	assert.Equal(t, sig1, sig2)

	// データが変われば署名も変わるはず
	log.Action = domain.AuditActionUpdate
	sig3, err := log.GenerateSignature(secretKey)
	assert.NoError(t, err)
	assert.NotEqual(t, sig1, sig3)

	// 鍵が変われば署名も変わるはず
	log.Action = domain.AuditActionCreate // 元に戻す
	sig4, err := log.GenerateSignature("another-secret")
	assert.NoError(t, err)
	assert.NotEqual(t, sig1, sig4)
}
