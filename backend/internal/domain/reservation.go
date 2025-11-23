// backend/internal/domain/reservation.go
// 予約ドメインモデル
//
// 責務:
// - 予約のビジネスルールを定義
// - 予約エンティティの構造定義
// - 繰り返し予定のルール管理
//
// フィールド:
// - ID: 予約ID (UUID)
// - OrganizerID: 主催者ID
// - Title: 件名
// - Description: 詳細説明
// - StartAt, EndAt: 開始・終了日時
// - RRule: 繰り返しルール (iCalendar形式)
// - IsPrivate: 非公開フラグ
// - ApprovalStatus: 承認ステータス (Pending, Confirmed, Rejected)
// - Timezone: タイムゾーン
// - UpdatedBy: 最終更新者
// - Version: 楽観的ロック用バージョン
// - DeletedAt: 論理削除日時

package domain

// TODO: 実装
