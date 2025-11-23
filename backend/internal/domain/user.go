// backend/internal/domain/user.go
// ユーザードメインモデル
//
// 責務:
// - ユーザーのビジネスルールを定義
// - ユーザーエンティティの構造定義
// - ロール・権限の定義

package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role はユーザーのロールを表す型
type Role string

// ロール定数
const (
	RoleGeneral   Role = "GENERAL"   // 一般ユーザー
	RoleSecretary Role = "SECRETARY" // 秘書
	RoleManager   Role = "MANAGER"   // マネージャー
	RoleAdmin     Role = "ADMIN"     // 管理者
	RoleAuditor   Role = "AUDITOR"   // 監査者
)

// User はユーザーエンティティを表す構造体
type User struct {
	ID                    uuid.UUID  // ユーザーID
	Sub                   string     // IdPから取得したユーザー識別子（不変）
	Email                 string     // メールアドレス
	Name                  string     // 表示名
	Role                  Role       // ロール
	ManagerID             *uuid.UUID // 上長のユーザーID（承認フロー用）
	PenaltyScore          int        // キャンセルペナルティスコア
	PenaltyScoreExpireAt  *time.Time // ペナルティスコア有効期限
	IsActive              bool       // アクティブフラグ
	CreatedAt             time.Time  // 作成日時
	UpdatedAt             time.Time  // 更新日時
	DeletedAt             *time.Time // 論理削除日時
}

// IsValid はユーザーが有効かどうかを判定します
func (u *User) IsValid() bool {
	return u.IsActive && u.DeletedAt == nil
}

// HasRole は指定されたロールを持っているかを判定します
func (u *User) HasRole(role Role) bool {
	return u.Role == role
}

// IsAdmin は管理者かどうかを判定します
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsManager はマネージャーかどうかを判定します
func (u *User) IsManager() bool {
	return u.Role == RoleManager
}

// IsSecretary は秘書かどうかを判定します
func (u *User) IsSecretary() bool {
	return u.Role == RoleSecretary
}

// IsAuditor は監査者かどうかを判定します
func (u *User) IsAuditor() bool {
	return u.Role == RoleAuditor
}

// CanManage は指定されたユーザーを管理できるかを判定します
// 管理者は全員を管理でき、マネージャーは自分の部下を管理できます
func (u *User) CanManage(targetUser *User) bool {
	if u.IsAdmin() {
		return true
	}
	
	if u.IsManager() && targetUser.ManagerID != nil {
		return *targetUser.ManagerID == u.ID
	}
	
	return false
}

// CanApproveReservation は予約を承認できるかを判定します
// マネージャー以上のロールが承認可能です
func (u *User) CanApproveReservation() bool {
	return u.Role == RoleManager || u.Role == RoleAdmin
}

// CanAccessResource は指定されたリソースにアクセスできるかを判定します
// required_role が nil の場合は全員アクセス可能
// required_role が設定されている場合は、そのロール以上が必要
func (u *User) CanAccessResource(requiredRole *Role) bool {
	if requiredRole == nil {
		return true
	}
	
	// ロールの階層: GENERAL < SECRETARY < MANAGER < ADMIN
	// AUDITOR は特殊なロールで、リソース予約には関与しない
	roleHierarchy := map[Role]int{
		RoleGeneral:   1,
		RoleSecretary: 2,
		RoleManager:   3,
		RoleAdmin:     4,
	}
	
	userLevel, userExists := roleHierarchy[u.Role]
	requiredLevel, requiredExists := roleHierarchy[*requiredRole]
	
	if !userExists || !requiredExists {
		return false
	}
	
	return userLevel >= requiredLevel
}

// HasActivePenalty はアクティブなペナルティを持っているかを判定します
func (u *User) HasActivePenalty() bool {
	if u.PenaltyScore <= 0 {
		return false
	}
	
	if u.PenaltyScoreExpireAt == nil {
		return false
	}
	
	return time.Now().Before(*u.PenaltyScoreExpireAt)
}

// GetActivePenaltyScore はアクティブなペナルティスコアを取得します
func (u *User) GetActivePenaltyScore() int {
	if u.HasActivePenalty() {
		return u.PenaltyScore
	}
	return 0
}
