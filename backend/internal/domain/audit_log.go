// backend/internal/domain/audit_log.go
// 監査ログドメインモデル
//
// 責務:
// - 監査ログのビジネスルールを定義
// - 改ざん検知のためのハッシュ管理
//
// フィールド:
// - ID: ログID (UUID)
// - UserID: 操作者ID
// - Action: 操作内容 (CREATE_EVENT, CANCEL_WITH_PENALTY等)
// - TargetType: 対象タイプ (Reservation, Resource等)
// - TargetID: 対象ID
// - Timestamp: 操作日時
// - SignatureHash: 改ざん検知用ハッシュ (HMAC-SHA256)

package domain

// TODO: 実装
