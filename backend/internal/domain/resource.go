// backend/internal/domain/resource.go
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ResourceType はリソースの種別を表す型
type ResourceType string

const (
	ResourceTypeMeetingRoom ResourceType = "MEETING_ROOM" // 会議室
	ResourceTypeEquipment   ResourceType = "EQUIPMENT"    // 備品
)

// Resource はリソースエンティティを表す構造体
type Resource struct {
	ID           uuid.UUID              // リソースID
	Name         string                 // リソース名
	Type         ResourceType           // リソース種別
	Capacity     *int                   // 収容人数（会議室の場合）
	Location     *string                // 場所
	Equipment    map[string]interface{} // 設備情報（JSON）
	RequiredRole *Role                  // 予約に必要な最低ロール
	IsActive     bool                   // アクティブフラグ
	CreatedAt    time.Time              // 作成日時
	UpdatedAt    time.Time              // 更新日時
}

// IsValid はリソースが有効かどうかを判定します
func (r *Resource) IsValid() bool {
	if r.Name == "" {
		return false
	}
	if r.Type != ResourceTypeMeetingRoom && r.Type != ResourceTypeEquipment {
		return false
	}
	return r.IsActive
}

// Validate はリソースの整合性を検証します
func (r *Resource) Validate() error {
	if r.Name == "" {
		return errors.New("resource name is required")
	}
	if r.Type == "" {
		return errors.New("resource type is required")
	}
	if r.Type == ResourceTypeMeetingRoom && (r.Capacity == nil || *r.Capacity <= 0) {
		return errors.New("capacity is required for meeting rooms")
	}
	return nil
}

// CanBeReservedBy は指定されたユーザーがこのリソースを予約できるかを判定します
func (r *Resource) CanBeReservedBy(user *User) bool {
	if !r.IsActive {
		return false
	}
	return user.CanAccessResource(r.RequiredRole)
}

// IsMeetingRoom は会議室かどうかを判定します
func (r *Resource) IsMeetingRoom() bool {
	return r.Type == ResourceTypeMeetingRoom
}

// IsEquipment は備品かどうかを判定します
func (r *Resource) IsEquipment() bool {
	return r.Type == ResourceTypeEquipment
}
