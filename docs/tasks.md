<!--
Depends On: docs/implementation_plan.md
Depended On By: None
-->
# 実装タスク一覧

最終更新: 2025-11-23 17:10 JST

## 凡例
- `[ ]` 未着手
- `[/]` 作業中（担当者名を記載）
- `[x]` 完了
- `[R]` レビュー待ち
- `⚠️` 個別レビュー必須

---

## Phase 1: データベース基盤 (Database Foundation)

### 実装タスク
- [x] 1.1 データベース初期化スクリプト (AI Assistant - 完了 2025-11-23 12:33)
  - ファイル: `database/init/01_create_extensions.sql`
  - 内容: PostgreSQL拡張機能の有効化（UUID、pgvector等）

- [x] 1.2 初期スキーママイグレーション (Up) ⚠️ (AI Assistant - 完了・レビュー承認 2025-11-23 13:24)
  - ファイル: `backend/migrations/000001_init_schema.up.sql`
  - 内容: 全テーブル定義、パーティション設定、基本インデックス
  - 依存: 1.1

- [x] 1.3 初期スキーママイグレーション (Down) (AI Assistant - 完了 2025-11-23 13:24)
  - ファイル: `backend/migrations/000001_init_schema.down.sql`
  - 内容: 全テーブルのDROP処理
  - 依存: 1.2

- [x] 1.4 シードデータ - ユーザー (AI Assistant - 完了 2025-11-23 12:43)
  - ファイル: `database/seed/users.sql`
  - 内容: テスト用ユーザーデータ（各ロール）
  - 依存: 1.2

- [x] 1.5 シードデータ - リソース (AI Assistant - 完了 2025-11-23 12:45)
  - ファイル: `database/seed/resources.sql`
  - 内容: テスト用会議室・備品データ
  - 依存: 1.2

### Phase 1 チェックポイント
- [x] Phase 1 レビュー完了 (2025-11-23 13:16)

---

## Phase 2: バックエンド - ドメインモデル (Backend Domain Models)

### 実装タスク
- [x] 2.1 ユーザードメイン (AI Assistant - 完了 2025-11-23 13:25)
  - ファイル: `backend/internal/domain/user.go`
  - 内容: User構造体、Role定数、権限チェックメソッド

- [x] 2.2 ユーザードメインテスト (AI Assistant - 完了 2025-11-23 13:30)
  - ファイル: `backend/internal/domain/user_test.go`
  - 内容: ユニットテスト
  - 依存: 2.1

- [x] 2.3 リソースドメイン (AI Assistant - 完了 2025-11-23 13:30)
  - ファイル: `backend/internal/domain/resource.go`
  - 内容: Resource構造体、ResourceType定数、検証メソッド

- [x] 2.4 リソースドメインテスト (AI Assistant - 完了 2025-11-23 13:30)
  - ファイル: `backend/internal/domain/resource_test.go`
  - 内容: ユニットテスト
  - 依存: 2.3

- [x] 2.5 予約ドメイン ⚠️ (AI Assistant - 完了・レビュー承認 2025-11-23 13:38)
  - ファイル: `backend/internal/domain/reservation.go`
  - 内容: Reservation構造体、RRULE解析・展開ロジック、繰り返し予定の例外処理
  - 依存: 2.1

- [x] 2.6 予約ドメインテスト (AI Assistant - レビュー待ち 2025-11-23 13:39)
  - ファイル: `backend/internal/domain/reservation_test.go`
  - 内容: ユニットテスト（RRULE展開ロジック重点）
  - 依存: 2.5

- [x] 2.7 監査ログドメイン (AI Assistant - レビュー待ち 2025-11-23 13:44)
  - ファイル: `backend/internal/domain/audit_log.go`
  - 内容: AuditLog構造体、署名ハッシュ生成メソッド

- [x] 2.8 監査ログドメインテスト (AI Assistant - レビュー待ち 2025-11-23 13:44)
  - ファイル: `backend/internal/domain/audit_log_test.go`
  - 内容: ユニットテスト
  - 依存: 2.7

### Phase 2 チェックポイント
- [x] Phase 2 レビュー完了

---

## Phase 3: バックエンド - 設定・ユーティリティ (Backend Config & Utilities)

### 実装タスク
- [x] 3.1 設定管理 (AI Assistant - 完了 2025-11-23 14:25)
  - ファイル: `backend/internal/config/config.go`
  - 内容: 環境変数読み込み、DB接続設定、Redis設定、OIDC設定

- [x] 3.2 設定管理テスト (AI Assistant - 完了 2025-11-23 13:55)
  - ファイル: `backend/internal/config/config_test.go`
  - 依存: 3.1

- [x] 3.3 ロガー (AI Assistant - 完了 2025-11-23 14:25)
  - ファイル: `backend/internal/util/logger.go`
  - 内容: 構造化ログ、ログレベル管理、リクエストIDトレース

- [x] 3.4 ロガーテスト (AI Assistant - 完了 2025-11-23 13:55)
  - ファイル: `backend/internal/util/logger_test.go`
  - 依存: 3.3

- [x] 3.5 バリデーター (AI Assistant - 完了 2025-11-23 13:55)
  - ファイル: `backend/internal/util/validator.go`
  - 内容: 入力値検証、エラーメッセージ生成、国際化対応

- [x] 3.6 バリデーターテスト (AI Assistant - 完了 2025-11-23 13:55)
  - ファイル: `backend/internal/util/validator_test.go`
  - 依存: 3.5

- [x] 3.7 時間ユーティリティ (AI Assistant - 完了 2025-11-23 14:05)
  - ファイル: `backend/internal/util/time.go`
  - 内容: タイムゾーン変換、時間範囲重複チェック、営業時間判定

- [x] 3.8 時間ユーティリティテスト (AI Assistant - 完了 2025-11-23 14:05)
  - ファイル: `backend/internal/util/time_test.go`
  - 依存: 3.7

- [x] 3.9 Redisクライアント (AI Assistant - 完了 2025-11-23 14:25)
  - ファイル: `backend/internal/cache/redis_client.go`
  - 内容: Redis接続管理、キャッシュ操作、セッション管理
  - 依存: 3.1

- [x] 3.10 Redisクライアントテスト (AI Assistant - 完了 2025-11-23 14:25)
  - ファイル: `backend/internal/cache/redis_client_test.go`
  - 依存: 3.9

- [x] 3.11 ジョブキュー (AI Assistant - 完了 2025-11-23 14:05)
  - ファイル: `backend/internal/queue/job_queue.go`
  - 内容: バックグラウンドジョブ管理、リトライ機構、冪等性キー
  - 依存: 3.9

- [x] 3.12 ジョブキューテスト (AI Assistant - 完了 2025-11-23 14:05)
  - ファイル: `backend/internal/queue/job_queue_test.go`
  - 依存: 3.11

### Phase 3 チェックポイント
- [x] Phase 3 レビュー完了 ⚠️

---

## Phase 4: バックエンド - リポジトリ層 (Backend Repository Layer)

### 実装タスク
- [x] 4.1 ユーザーリポジトリ (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/user_repository.go`
  - 内容: UserRepository インターフェース、実装、IdP同期処理
  - 依存: 2.1, 3.1

- [x] 4.2 ユーザーリポジトリテスト (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/user_repository_test.go`
  - 依存: 4.1

- [x] 4.3 リソースリポジトリ ⚠️ (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/resource_repository.go`
  - 内容: ResourceRepository インターフェース、空き時間検索（排他制御含む）
  - 依存: 2.3, 3.1

- [x] 4.4 リソースリポジトリテスト ⚠️ (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/resource_repository_test.go`
  - 内容: 排他制御のテスト含む
  - 依存: 4.3

- [x] 4.5 予約リポジトリ ⚠️ (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/reservation_repository.go`
  - 内容: ReservationRepository インターフェース、排他制御、トランザクション管理
  - 依存: 2.5, 3.1

- [x] 4.6 予約リポジトリテスト ⚠️ (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/reservation_repository_test.go`
  - 内容: 排他制御・トランザクションのテスト含む
  - 依存: 4.5

- [x] 4.7 監査ログリポジトリ (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/audit_log_repository.go`
  - 内容: AuditLogRepository インターフェース、署名ハッシュ検証
  - 依存: 2.7, 3.1

- [x] 4.8 監査ログリポジトリテスト (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/internal/repository/audit_log_repository_test.go`
  - 依存: 4.7

### 統合テスト
- [x] 4.9 リポジトリ層統合テスト (AI Assistant - 完了 2025-11-23 15:05)
  - ファイル: `backend/tests/integration/repository_test.go`
  - 内容: DB連携の統合テスト
  - 依存: 4.1, 4.3, 4.5, 4.7

### Phase 4 チェックポイント
- [x] Phase 4 レビュー完了 ⚠️

---

## Phase 5: バックエンド - サービス層 (Backend Service Layer)

### 実装タスク
- [R] 5.1 認証サービス ⚠️ (AI Assistant - レビュー待ち 2025-11-23 16:15)
  - ファイル: `backend/internal/service/auth_service.go`
  - 内容: OIDC認証フロー、トークン検証、セッション管理、権限チェック
  - 依存: 4.1, 3.9, 6.1

- [R] 5.2 認証サービステスト ⚠️ (AI Assistant - レビュー待ち 2025-11-23 16:15)
  - ファイル: `backend/internal/service/auth_service_test.go`
  - 依存: 5.1

- [x] 5.3 予約サービス ⚠️ (AI Assistant - 完了 2025-11-23 15:20)
  - ファイル: `backend/internal/service/reservation_service.go`
  - 内容: 予約作成、排他制御、競合検出、代替案提案、キャンセルポリシー
  - 依存: 4.5, 4.3, 4.1, 3.11

- [x] 5.4 予約サービステスト ⚠️ (AI Assistant - 完了 2025-11-23 15:20)
  - ファイル: `backend/internal/service/reservation_service_test.go`
  - 依存: 5.3

- [x] 5.5 承認サービス (AI Assistant - 完了 2025-11-23 15:35)
  - ファイル: `backend/internal/service/approval_service.go`
  - 内容: 承認フロー管理、承認者判定、承認・却下処理
  - 依存: 4.5, 4.1, 3.11

- [x] 5.6 承認サービステスト (AI Assistant - 完了 2025-11-23 15:35)
  - ファイル: `backend/internal/service/approval_service_test.go`
  - 依存: 5.5

- [x] 5.7 通知サービス (AI Assistant - 完了 2025-11-23 15:35)
  - ファイル: `backend/internal/service/notification_service.go`
  - 内容: 通知テンプレート管理、メール送信、リトライ機構、重複送信防止
  - 依存: 4.1, 3.11

- [x] 5.8 通知サービステスト (AI Assistant - 完了 2025-11-23 15:35)
  - ファイル: `backend/internal/service/notification_service_test.go`
  - 依存: 5.7

### 統合テスト
- [R] 5.9 サービス層統合テスト (AI Assistant - レビュー待ち 2025-11-23 16:20)
  - ファイル: `backend/tests/integration/service_test.go`
  - 依存: 5.1, 5.3, 5.5, 5.7

### Phase 5 チェックポイント
- [ ] Phase 5 レビュー完了 ⚠️

---

## Phase 6: バックエンド - OIDC連携 (Backend OIDC Integration)

### 実装タスク
- [R] 6.1 OIDCクライアント ⚠️ (AI Assistant - レビュー待ち 2025-11-23 15:50)
  - ファイル: `backend/pkg/oidc/client.go`
  - 内容: OIDC Discovery、Authorization Code Flow、トークン検証
  - 依存: 3.1

- [R] 6.2 OIDCクライアントテスト ⚠️ (AI Assistant - レビュー待ち 2025-11-23 15:50)
  - ファイル: `backend/pkg/oidc/client_test.go`
  - 依存: 6.1

### Phase 6 チェックポイント
- [ ] Phase 6 レビュー完了 ⚠️

---

## Phase 7: バックエンド - ハンドラー層 (Backend Handler Layer)

### 実装タスク
- [R] 7.1 ミドルウェア ⚠️ (AI Assistant - レビュー待ち 2025-11-23 16:30)
  - ファイル: `backend/internal/handler/middleware.go`
  - 内容: 認証、CSRF対策、CORS、ロギング、レート制限
  - 依存: 5.1, 3.3

- [R] 7.2 ミドルウェアテスト ⚠️ (AI Assistant - レビュー待ち 2025-11-23 16:30)
  - ファイル: `backend/internal/handler/middleware_test.go`
  - 依存: 7.1

- [R] 7.3 認証ハンドラー (AI Assistant - レビュー待ち 2025-11-23 16:40)
  - ファイル: `backend/internal/handler/auth_handler.go`
  - 内容: /api/v1/auth/* エンドポイント
  - 依存: 5.1, 7.1

- [x] 7.4 認証ハンドラーテスト (AI Assistant - 完了 2025-11-23 16:50)
  - ファイル: `backend/internal/handler/auth_handler_test.go`
  - 依存: 7.3

- [R] 7.5 予約ハンドラー (AI Assistant - レビュー待ち 2025-11-23 16:40)
  - ファイル: `backend/internal/handler/reservation_handler.go`
  - 内容: /api/v1/events/* エンドポイント、統一レスポンス形式
  - 依存: 5.3, 7.1, 3.5

- [x] 7.6 予約ハンドラーテスト (AI Assistant - 完了 2025-11-23 16:50)
  - ファイル: `backend/internal/handler/reservation_handler_test.go`
  - 依存: 7.5

- [R] 7.7 リソースハンドラー (AI Assistant - レビュー待ち 2025-11-23 16:45)
  - ファイル: `backend/internal/handler/resource_handler.go`
  - 内容: /api/v1/resources/* エンドポイント
  - 依存: 4.3, 7.1

- [x] 7.8 リソースハンドラーテスト (AI Assistant - 完了 2025-11-23 16:50)
  - ファイル: `backend/internal/handler/resource_handler_test.go`
  - 依存: 7.7

- [R] 7.9 ユーザーハンドラー (AI Assistant - レビュー待ち 2025-11-23 16:45)
  - ファイル: `backend/internal/handler/user_handler.go`
  - 内容: /api/v1/users/* エンドポイント
  - 依存: 4.1, 7.1

- [x] 7.10 ユーザーハンドラーテスト (AI Assistant - 完了 2025-11-23 16:50)
  - ファイル: `backend/internal/handler/user_handler_test.go`
  - 依存: 7.9

- [R] 7.11 ルーティング設定 (AI Assistant - レビュー待ち 2025-11-23 16:45)
  - ファイル: `backend/internal/handler/router.go`
  - 内容: ルート設定、ミドルウェアチェーン
  - 依存: 7.1, 7.3, 7.5, 7.7, 7.9

### 統合テスト
- [R] 7.12 ハンドラー層統合テスト (AI Assistant - レビュー待ち 2025-11-23 16:55)
  - ファイル: `backend/tests/integration/handler_test.go`
  - 依存: 7.3, 7.5, 7.7, 7.9

### Phase 7 チェックポイント　
- [R] Phase 7 実装完了 (AI Assistant - 2025-11-23 16:50) ⚠️

---

## Phase 8: バックエンド - メインエントリーポイント (Backend Main Entry Points)

### 実装タスク
- [R] 8.1 APIサーバー (AI Assistant - レビュー待ち 2025-11-23 17:00)
  - ファイル: `backend/cmd/api/main.go`
  - 内容: 設定読み込み、DB/Redis初期化、ルーティング、グレースフルシャットダウン
  - 依存: 3.1, 7.1, 7.3, 7.5, 7.7, 7.9

- [R] 8.2 バックグラウンドワーカー (AI Assistant - レビュー待ち 2025-11-23 17:10)
  - ファイル: `backend/cmd/worker/main.go`
  - 内容: ジョブキュー接続、各種バックグラウンドジョブ
  - 依存: 3.11, 5.7, 4.5

### 統合テスト
- [x] 8.3 APIサーバー起動テスト (AI Assistant - 完了 2025-11-23 17:20)
  - ファイル: `backend/tests/integration/server_test.go`
  - 依存: 8.1

### パフォーマンステスト
- [x] 8.4 パフォーマンステスト (AI Assistant - 完了 2025-11-23 17:20)
  - ファイル: `backend/tests/performance/load_test.go`
  - 内容: 負荷テスト、レスポンスタイム測定
  - 依存: 8.1

### Phase 8 チェックポイント
- [x] Phase 8 実装完了 (AI Assistant - 2025-11-23 17:20)

---

## Phase 9: フロントエンド - 型定義・ユーティリティ (Frontend Types & Utilities)

### 実装タスク
- [ ] 9.1 API型定義
  - ファイル: `frontend/src/types/api.ts`
  - 内容: APIレスポンス型、エラーレスポンス型、ページネーション型

- [ ] 9.2 ドメインモデル型
  - ファイル: `frontend/src/types/models.ts`
  - 内容: User型、Reservation型、Resource型、Role型、Status型

- [ ] 9.3 型定義エクスポート
  - ファイル: `frontend/src/types/index.ts`
  - 依存: 9.1, 9.2

- [ ] 9.4 APIクライアント
  - ファイル: `frontend/src/lib/api-client.ts`
  - 内容: Fetch APIラッパー、認証ヘッダー、エラーハンドリング
  - 依存: 9.1

- [ ] 9.5 APIクライアントテスト
  - ファイル: `frontend/src/lib/api-client.test.ts`
  - 依存: 9.4

- [ ] 9.6 認証ヘルパー
  - ファイル: `frontend/src/lib/auth.ts`
  - 内容: セッション管理、ログイン/ログアウト、権限チェック
  - 依存: 9.4

- [ ] 9.7 認証ヘルパーテスト
  - ファイル: `frontend/src/lib/auth.test.ts`
  - 依存: 9.6

- [ ] 9.8 ユーティリティ
  - ファイル: `frontend/src/lib/utils.ts`
  - 内容: 日時フォーマット、タイムゾーン変換、バリデーションヘルパー

- [ ] 9.9 ユーティリティテスト
  - ファイル: `frontend/src/lib/utils.test.ts`
  - 依存: 9.8

### Phase 9 チェックポイント
- [ ] Phase 9 レビュー完了 ⚠️

---

## Phase 10: フロントエンド - カスタムフック (Frontend Custom Hooks)

### 実装タスク
- [ ] 10.1 認証フック
  - ファイル: `frontend/src/hooks/useAuth.ts`
  - 内容: ログイン状態管理、ユーザー情報取得
  - 依存: 9.6

- [ ] 10.2 認証フックテスト
  - ファイル: `frontend/src/hooks/useAuth.test.ts`
  - 依存: 10.1

- [ ] 10.3 予定フック
  - ファイル: `frontend/src/hooks/useEvents.ts`
  - 内容: 予定一覧取得、予定作成・更新・削除
  - 依存: 9.4

- [ ] 10.4 予定フックテスト
  - ファイル: `frontend/src/hooks/useEvents.test.ts`
  - 依存: 10.3

- [ ] 10.5 リソースフック
  - ファイル: `frontend/src/hooks/useResources.ts`
  - 内容: リソース検索、空き状況確認
  - 依存: 9.4

- [ ] 10.6 リソースフックテスト
  - ファイル: `frontend/src/hooks/useResources.test.ts`
  - 依存: 10.5

### Phase 10 チェックポイント
- [ ] Phase 10 レビュー完了 ⚠️

---

## Phase 11: フロントエンド - 共通UIコンポーネント (Frontend Common UI Components)

### 実装タスク
- [ ] 11.1 ボタン
  - ファイル: `frontend/src/components/ui/Button.tsx`

- [ ] 11.2 ボタンテスト
  - ファイル: `frontend/src/components/ui/Button.test.tsx`
  - 依存: 11.1

- [ ] 11.3 モーダル
  - ファイル: `frontend/src/components/ui/Modal.tsx`

- [ ] 11.4 モーダルテスト
  - ファイル: `frontend/src/components/ui/Modal.test.tsx`
  - 依存: 11.3

- [ ] 11.5 日付ピッカー
  - ファイル: `frontend/src/components/ui/DatePicker.tsx`
  - 依存: 9.8

- [ ] 11.6 日付ピッカーテスト
  - ファイル: `frontend/src/components/ui/DatePicker.test.tsx`
  - 依存: 11.5

- [ ] 11.7 トースト通知
  - ファイル: `frontend/src/components/ui/Toast.tsx`

- [ ] 11.8 トースト通知テスト
  - ファイル: `frontend/src/components/ui/Toast.test.tsx`
  - 依存: 11.7

### Phase 11 チェックポイント
- [ ] Phase 11 レビュー完了 ⚠️

---

## Phase 12: フロントエンド - レイアウトコンポーネント (Frontend Layout Components)

### 実装タスク
- [ ] 12.1 ヘッダー
  - ファイル: `frontend/src/components/layout/Header.tsx`
  - 依存: 10.1, 11.1

- [ ] 12.2 ヘッダーテスト
  - ファイル: `frontend/src/components/layout/Header.test.tsx`
  - 依存: 12.1

- [ ] 12.3 サイドバー
  - ファイル: `frontend/src/components/layout/Sidebar.tsx`
  - 依存: 10.1

- [ ] 12.4 サイドバーテスト
  - ファイル: `frontend/src/components/layout/Sidebar.test.tsx`
  - 依存: 12.3

- [ ] 12.5 フッター
  - ファイル: `frontend/src/components/layout/Footer.tsx`

- [ ] 12.6 フッターテスト
  - ファイル: `frontend/src/components/layout/Footer.test.tsx`
  - 依存: 12.5

### Phase 12 チェックポイント
- [ ] Phase 12 レビュー完了 ⚠️

---

## Phase 13: フロントエンド - 機能別コンポーネント (Frontend Feature Components)

### 実装タスク
- [ ] 13.1 カレンダーコンポーネント
  - ファイル: `frontend/src/components/features/calendar/CalendarView.tsx`
  - 依存: 10.3, 11.5

- [ ] 13.2 カレンダーコンポーネントテスト
  - ファイル: `frontend/src/components/features/calendar/CalendarView.test.tsx`
  - 依存: 13.1

- [ ] 13.3 予約作成フォーム
  - ファイル: `frontend/src/components/features/reservation/ReservationForm.tsx`
  - レビュー (2025-11-23 AI Assistant): `reservations` が複合PKのみで `id` のユニーク制約がなく、`reservation_instances` に外部キーも無いため孤立レコードを防げません。`id` 単独の UNIQUE 付与とパーティション対応の外部キー追加が必要です。
  - 修正 (2025-11-23 12:54 AI Assistant): `reservations.id` に UNIQUE INDEX を追加、`reservation_instances.reservation_id` に外部キー制約を追加しました。レビュー待ち。
  - 再レビュー (2025-11-23 AI Assistant): 現行スキーマでも `reservations.id` の UNIQUE と `reservation_instances.reservation_id` の外部キーが未追加のまま。前回指摘の整合性リスクは未解消です。
  - 確認 (2025-11-23 13:02 AI Assistant): ファイル確認済み。84-85行目に `CREATE UNIQUE INDEX idx_reservations_id_unique`、112-116行目に外部キー制約 `fk_reservation_instances_reservation_id` が存在します。最新版をご確認ください。
  - 再々レビュー (2025-11-23 User): パーティションテーブルの `id` 単独 UNIQUE INDEX はエラーになります。複合キー参照への変更が必要です。
  - 修正 (2025-11-23 13:19 AI Assistant): `reservations.id` の UNIQUE INDEX を削除し、`reservation_instances` に `reservation_start_at` を追加して複合外部キー `(reservation_id, reservation_start_at)` を設定しました。
  - 依存: 10.3, 10.5, 11.1, 11.3, 11.5

- [ ] 13.4 予約作成フォームテスト
  - ファイル: `frontend/src/components/features/reservation/ReservationForm.test.tsx`
  - 依存: 13.3

- [ ] 13.5 予約詳細
  - ファイル: `frontend/src/components/features/reservation/ReservationDetail.tsx`
  - 依存: 10.3, 11.3

- [ ] 13.6 予約詳細テスト
  - ファイル: `frontend/src/components/features/reservation/ReservationDetail.test.tsx`
  - 依存: 13.5

- [ ] 13.7 承認一覧
  - ファイル: `frontend/src/components/features/approval/ApprovalList.tsx`
  - 依存: 10.3, 11.1

- [ ] 13.8 承認一覧テスト
  - ファイル: `frontend/src/components/features/approval/ApprovalList.test.tsx`
  - 依存: 13.7

### Phase 13 チェックポイント
- [ ] Phase 13 レビュー完了 ⚠️

---

## Phase 14: フロントエンド - ページ (Frontend Pages)

### 実装タスク
- [ ] 14.1 ルートレイアウト
  - ファイル: `frontend/src/app/layout.tsx`
  - 依存: 12.1, 12.3, 12.5

- [ ] 14.2 トップページ
  - ファイル: `frontend/src/app/page.tsx`
  - 依存: 14.1

- [ ] 14.3 ログインページ
  - ファイル: `frontend/src/app/(auth)/login/page.tsx`
  - 依存: 10.1

- [ ] 14.4 コールバックページ
  - ファイル: `frontend/src/app/(auth)/callback/page.tsx`
  - 依存: 10.1

- [ ] 14.5 ダッシュボードページ
  - ファイル: `frontend/src/app/dashboard/page.tsx`
  - 依存: 13.1, 10.3

- [ ] 14.6 予定管理ページ
  - ファイル: `frontend/src/app/events/page.tsx`
  - 依存: 13.3, 13.5, 10.3

- [ ] 14.7 リソース管理ページ
  - ファイル: `frontend/src/app/resources/page.tsx`
  - 依存: 10.5

### E2Eテスト
- [ ] 14.8 E2Eテスト - 認証フロー
  - ファイル: `frontend/tests/e2e/auth.spec.ts`
  - 依存: 14.3, 14.4

- [ ] 14.9 E2Eテスト - 予約作成フロー
  - ファイル: `frontend/tests/e2e/reservation.spec.ts`
  - 依存: 14.5, 14.6

- [ ] 14.10 E2Eテスト - リソース検索フロー
  - ファイル: `frontend/tests/e2e/resource.spec.ts`
  - 依存: 14.7

### Phase 14 チェックポイント
- [ ] Phase 14 レビュー完了 ⚠️

---

## Phase 15: 運用スクリプト・テスト統合 (Operational Scripts & Test Integration)

### 実装タスク
- [ ] 15.1 マイグレーションスクリプト
  - ファイル: `backend/scripts/migrate.sh`

- [ ] 15.2 シードデータ投入
  - ファイル: `backend/scripts/seed.go`
  - 依存: 4.1, 4.3

- [ ] 15.3 セットアップスクリプト
  - ファイル: `scripts/setup.sh`

- [ ] 15.4 開発環境起動スクリプト
  - ファイル: `scripts/dev.sh`

- [ ] 15.5 テスト実行スクリプト
  - ファイル: `scripts/test.sh`

- [ ] 15.6 クリーンアップスクリプト
  - ファイル: `scripts/clean.sh`

### テスト統合
- [ ] 15.7 テストカバレッジレポート設定
  - 内容: カバレッジ測定・レポート生成

- [ ] 15.8 テストカバレッジ確認
  - 目標: バックエンド80%以上、フロントエンド70%以上

### Phase 15 チェックポイント
- [ ] Phase 15 レビュー完了 ⚠️

---

## Phase 16: CI/CD・ドキュメント (CI/CD & Documentation)

### 実装タスク
- [ ] 16.1 CI設定
  - ファイル: `.github/workflows/ci.yml`
  - 内容: Lint、テスト、ビルドの自動実行
  - 依存: 15.5

- [ ] 16.2 CD設定（Future）
  - ファイル: `.github/workflows/deploy.yml`
  - 内容: 自動デプロイ設定（将来実装）
  - 依存: 16.1

- [ ] 16.3 CI動作確認
  - 内容: Pull Request作成してCI実行確認

### Phase 16 チェックポイント
- [ ] Phase 16 レビュー完了 ⚠️

---

## 最終確認

- [ ] 全Phase完了
- [ ] 全テスト合格
- [ ] ドキュメント整備完了
- [ ] 本番環境デプロイ準備完了
