// backend/internal/domain/audit_log.go
package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AuditAction は監査アクションの種類を表す型
type AuditAction string

const (
	AuditActionCreate      AuditAction = "CREATE"
	AuditActionUpdate      AuditAction = "UPDATE"
	AuditActionDelete      AuditAction = "DELETE"
	AuditActionLogin       AuditAction = "LOGIN"
	AuditActionLogout      AuditAction = "LOGOUT"
	AuditActionApprove     AuditAction = "APPROVE"
	AuditActionReject      AuditAction = "REJECT"
	AuditActionCheckIn     AuditAction = "CHECK_IN"
	AuditActionCancel      AuditAction = "CANCEL"
	AuditActionForceCancel AuditAction = "FORCE_CANCEL"
)

// AuditLog は監査ログエンティティを表す構造体
type AuditLog struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Action     AuditAction
	Resource   string // 対象リソース名（例: "reservations", "users"）
	ResourceID string // 対象リソースID
	Details    map[string]interface{}
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time
}

// GenerateSignature は監査ログの改ざん検知用署名を生成します
// secretKey はサーバー設定から注入されることを想定
func (a *AuditLog) GenerateSignature(secretKey string) (string, error) {
	detailsJSON, err := json.Marshal(a.Details)
	if err != nil {
		return "", fmt.Errorf("failed to marshal details: %w", err)
	}

	// 署名対象のデータを結合
	// 順序: ID:UserID:Action:Resource:ResourceID:Details:CreatedAt:IPAddress:UserAgent
	data := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%d:%s:%s",
		a.ID.String(),
		a.UserID.String(),
		string(a.Action),
		a.Resource,
		a.ResourceID,
		string(detailsJSON),
		a.CreatedAt.UnixNano(),
		a.IPAddress,
		a.UserAgent,
	)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}
