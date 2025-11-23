// backend/internal/domain/user.go
// ユーザードメインモデル
//
// 責務:
// - ユーザーのビジネスルールを定義
// - ユーザーエンティティの構造定義
// - ロール・権限の定義
//
// フィールド:
// - ID: ユーザーID (UUID)
// - Sub: IdPから取得したユーザー識別子
// - Email: メールアドレス
// - Name: 表示名
// - Role: ロール (GENERAL, SECRETARY, MANAGER, ADMIN, AUDITOR)
// - ManagerID: 上長のユーザーID (承認フロー用)
// - CreatedAt, UpdatedAt: タイムスタンプ

package domain

// TODO: 実装
