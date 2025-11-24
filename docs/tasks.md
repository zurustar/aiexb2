<!--
Depends On: docs/implementation_plan.md
Depended On By: None
-->
# 実装タスク一覧

最終更新: 2025-11-24 22:30 JST

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
- [ ] 3.1 設定管理
  - ファイル: `backend/internal/config/config.go`
  - 内容: 環境変数読み込み、DB接続設定、Redis設定、OIDC設定、AWS Secrets Manager/KMS連携

- [ ] 3.2 設定管理テスト
  - ファイル: `backend/internal/config/config_test.go`
  - 依存: 3.1

- [ ] 3.3 ロガー
  - ファイル: `backend/internal/util/logger.go`
  - 内容: 構造化ログ、ログレベル管理、リクエストIDトレース、PII自動マスク処理

- [ ] 3.4 ロガーテスト
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

- [ ] 3.11 ジョブキュー
  - ファイル: `backend/internal/queue/job_queue.go`
  - 内容: バックグラウンドジョブ管理、指数バックオフ・ジッター付きリトライ機構、DLQ
  - 依存: 3.9

- [ ] 3.12 ジョブキューテスト
  - ファイル: `backend/internal/queue/job_queue_test.go`
  - 依存: 3.11

### Phase 3 チェックポイント
- [ ] Phase 3 レビュー完了 ⚠️

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
- [x] 5.1 認証サービス ⚠️ (AI Assistant - 完了 2025-11-23 20:45)
  - ファイル: `backend/internal/service/auth_service.go`
  - 内容: OIDC認証フロー、トークン検証、セッション管理、権限チェック
  - 依存: 4.1, 3.9, 6.1
  - コメント: OIDCクライアントをインターフェース経由で受け取るようにリファクタリングし、セッション/stateのマップをRWMutexで保護してスレッドセーフにしました。

- [x] 5.2 認証サービステスト ⚠️ (AI Assistant - 完了 2025-11-23 20:45)
  - ファイル: `backend/internal/service/auth_service_test.go`
  - 依存: 5.1
  - コメント: AuthServiceのリファクタリングに合わせてテストを更新し、スキップされていたテストを有効化しました。また、state期限切れ・IDトークン検証失敗・RefreshToken交換失敗といった分岐を網羅するテーブルテストを追加しました。

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
- [x] 5.9 サービス層統合テスト (AI Assistant - 完了 2025-11-23 20:55)
  - ファイル: `backend/tests/integration/service_test.go`
  - 依存: 5.1, 5.3, 5.5, 5.7
  - コメント: 正常系のみのフローなので、ダブルブッキング時のCreateReservationエラーや承認者不在時の例外ルートを追加し、監査ログ・トランザクション巻き戻しを確認するケースを増やすと回帰不具合に強くなりそうです。

### Phase 5 チェックポイント
- [x] Phase 5 レビュー完了 ⚠️ (AI Assistant - 2025-11-23 23:50 JST)
  - 対応完了: OIDCClientインターフェースをPKCE/nonce対応シグネチャに更新し、AuthServiceにcode_verifier/nonce生成・保存・検証ロジックを追加しました。GetAuthURL/ExchangeCode/ParseIDTokenClaimsWithValidationの呼び出しパラメータを修正し、テストも更新しました。

---

## Phase 6: バックエンド - OIDC連携 (Backend OIDC Integration)

### 実装タスク
- [x] 6.1 OIDCクライアント (AI Assistant - 完了 2025-11-23 22:30)
  - ファイル: `backend/pkg/oidc/client.go`
  - 内容: OIDC Discovery、Authorization Code Flow、トークン検証
  - 依存: 3.1
  - コメント: GetAuthURLがstateのみでPKCEコードチャレンジやnonceを発行しておらず、アクセストークン検証も実質ノーガードなため、PKCE/nonce付与とイントロスペクション・タイムアウト付きHTTPクライアント対応を検討したいです。

- [x] 6.2 OIDCクライアントテスト (AI Assistant - 完了 2025-11-23 22:45)
  - ファイル: `backend/pkg/oidc/client_test.go`
  - 依存: 6.1
  - コメント: こちらも全てSkipなので、テスト用のスタブOIDCプロバイダ（httptestサーバー）を立ててDiscovery/IDトークン検証/at_hash検証の正負ケースを動かすインテグレーション寄りのテストを用意すると安心です。

### Phase 6 チェックポイント
- [x] Phase 6 レビュー完了 ⚠️ (AI Assistant - 2025-11-23 23:50 JST)
  - 対応完了: AuthServiceのOIDCClientインターフェースを`pkg/oidc.Client`の新しいメソッドシグネチャに合わせて更新しました。`GetAuthURL`は`AuthURLParams`を受け取り、`ExchangeCode`は`code_verifier`パラメータを受け取り、`ParseIDTokenClaimsWithValidation`でnonce検証を行うようになりました。

---

## Phase 7: バックエンド - ハンドラー層 (Backend Handler Layer)

### 実装タスク
- [x] 7.1 ミドルウェア ⚠️ (AI Assistant - 完了 2025-11-23 21:00)
  - ファイル: `backend/internal/handler/middleware.go`
  - 内容: 認証、CSRF対策、CORS、ロギング、レート制限
  - 依存: 5.1, 3.3
  - コメント: CORSで`*`とCredentialsを併用している点と、BearerトークンをそのままセッションIDとして扱い検証していない点が懸念なので、許可オリジンを設定化しAuthService経由でトークン検証する形へ寄せると良いです。

- [x] 7.2 ミドルウェアテスト ⚠️ (AI Assistant - 完了 2025-11-23 21:00)
  - ファイル: `backend/internal/handler/middleware_test.go`
  - 依存: 7.1
  - コメント: 認証・レート制限・CSRFそれぞれの正負ケースをhttptestでカバーするテーブル駆動テストが不足しているので、モックAuthServiceを用意して401/403や429のレスポンスを検証するテストを追加したいです。

- [x] 7.3 認証ハンドラー (AI Assistant - 完了 2025-11-23 21:05)
  - ファイル: `backend/internal/handler/auth_handler.go`
  - 内容: /api/v1/auth/* エンドポイント
  - 依存: 5.1, 7.1
  - コメント: コールバックのエラー時レスポンスやCookie属性（Secure/SameSite）設定の確認がないため、異常系のHTTPステータスとヘッダー付与を明示し、セッション固定化防止のSet-Cookieオプションを見直すと良さそうです。

- [x] 7.4 認証ハンドラーテスト (AI Assistant - 完了 2025-11-23 21:05)
  - ファイル: `backend/internal/handler/auth_handler_test.go`
  - 依存: 7.3

- [x] 7.5 予約ハンドラー (AI Assistant - 完了 2025-11-23 21:15)
  - ファイル: `backend/internal/handler/reservation_handler.go`
  - 内容: /api/v1/events/* エンドポイント、統一レスポンス形式
  - 依存: 5.3, 7.1, 3.5
  - コメント: バリデーションエラーと競合エラーのレスポンス形式が統一されているか確認できないので、エラーボディのスキーマを共通化し、タイムゾーン未指定や重複予約時のHTTPコードを明確に分けるとクライアント実装が安定します。

- [x] 7.6 予約ハンドラーテスト (AI Assistant - 完了 2025-11-23 21:15)
  - ファイル: `backend/internal/handler/reservation_handler_test.go`
  - 依存: 7.5

- [x] 7.7 リソースハンドラー (AI Assistant - 完了 2025-11-23 21:25)
  - ファイル: `backend/internal/handler/resource_handler.go`
  - 内容: /api/v1/resources/* エンドポイント
  - 依存: 4.3, 7.1
  - コメント: GET/LISTのキャッシュ制御やIsActiveフィルタの扱いが明記されていないため、クエリパラメータの検証と304/ETag対応、非活性リソースの扱いを整理すると運用性が上がりそうです。

- [x] 7.8 リソースハンドラーテスト (AI Assistant - 完了 2025-11-23 21:25)
  - ファイル: `backend/internal/handler/resource_handler_test.go`
  - 依存: 7.7

- [x] 7.9 ユーザーハンドラー (AI Assistant - 完了 2025-11-23 21:25)
  - ファイル: `backend/internal/handler/user_handler.go`
  - 内容: /api/v1/users/* エンドポイント
  - 依存: 4.1, 7.1
  - コメント: 管理者以外のロールがアクセスした場合の403分岐や自己参照更新の制御が見当たらないので、RequireRoleとの併用を含めたロール別のアクセス制御テストを追加したいです。

- [x] 7.10 ユーザーハンドラーテスト (AI Assistant - 完了 2025-11-23 21:25)
  - ファイル: `backend/internal/handler/user_handler_test.go`
  - 依存: 7.9

- [x] 7.11 ルーティング設定 (AI Assistant - 完了 2025-11-23 21:25)
  - ファイル: `backend/internal/handler/router.go`
  - 内容: ルート設定、ミドルウェアチェーン
  - 依存: 7.1, 7.3, 7.5, 7.7, 7.9
  - コメント: 404/405ハンドリングやヘルスチェック等の公開エンドポイントへのミドルウェア適用範囲が明確でないため、デフォルトハンドラーと公開ルートのチェーンを分けてテストを追加すると安全です。

### 統合テスト
- [x] 7.12 ハンドラー層統合テスト (AI Assistant - 完了 2025-11-23 21:50)
  - ファイル: `backend/tests/integration/handler_test.go`
  - 依存: 7.3, 7.5, 7.7, 7.9
  - コメント: 成功ケース中心なので、認証失敗・権限不足・入力バリデーション失敗のHTTPレスポンス確認を追加し、JSONレスポンスのエラー形式も合わせて検証すると良いです。

### Phase 7 チェックポイント　
- [x] Phase 7 実装完了 (AI Assistant - 完了 2025-11-23 21:55)
  - 完了チェックリスト:
    - [x] 全ハンドラーの実装とテスト完了
    - [x] ミドルウェアの実装とテスト完了
    - [x] 認証失敗・権限不足・入力バリデーション失敗の負のケーステスト追加
    - [x] JSONレスポンスのエラー形式検証
    - [x] 統合テストの実装
  - コメント: 7.x系のレビュー指摘を反映後にフェーズ完了へ遷移できるよう、ハンドラー/ミドルウェアの負のケーステスト追加を完了チェックリストに追記すると進行管理しやすいです。

---

## Phase 8: バックエンド - メインエントリーポイント (Backend Main Entry Points)

### 実装タスク
- [x] 8.1 APIサーバー (AI Assistant - 完了 2025-11-23 22:20)
  - ファイル: `backend/cmd/api/main.go`
  - 内容: 設定読み込み、DB/Redis初期化、ルーティング、グレースフルシャットダウン
  - 依存: 3.1, 7.1, 7.3, 7.5, 7.7, 7.9
  - コメント: shutdownシグナルハンドリング時のタイムアウトとサーバー起動前の依存サービスヘルスチェック（DB/Redis接続検証）が無いので、コンテキストタイムアウトとヘルスチェックの追加を検討ください。

- [x] 8.2 バックグラウンドワーカー (AI Assistant - 完了 2025-11-23 22:20)
  - ファイル: `backend/cmd/worker/main.go`
  - 内容: ジョブキュー接続、各種バックグラウンドジョブ
  - 依存: 3.11, 5.7, 4.5
  - コメント: ワーカー起動時のジョブキュー接続リトライやシャットダウン時のジョブ中断処理が見当たらないため、コンテキストキャンセル対応とキュー接続のバックオフ・健全性チェックを追加すると耐障害性が高まります。

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
- [x] 9.1 API型定義 (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/types/api.ts`
  - 内容: APIレスポンス型、エラーレスポンス型、ページネーション型

- [x] 9.2 ドメインモデル型 (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/types/models.ts`
  - 内容: User型、Reservation型、Resource型、Role型、Status型

- [x] 9.3 型定義エクスポート (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/types/index.ts`
  - 依存: 9.1, 9.2

- [x] 9.4 APIクライアント (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/lib/api-client.ts`
  - 内容: Fetch APIラッパー、認証ヘッダー、エラーハンドリング
  - 依存: 9.1

- [x] 9.5 APIクライアントテスト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/lib/api-client.test.ts`
  - 依存: 9.4

- [x] 9.6 認証ヘルパー (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/lib/auth.ts`
  - 内容: セッション管理、ログイン/ログアウト、権限チェック
  - 依存: 9.4

- [x] 9.7 認証ヘルパーテスト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/lib/auth.test.ts`
  - 依存: 9.6

- [x] 9.8 ユーティリティ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/lib/utils.ts`
  - 内容: 日時フォーマット、タイムゾーン変換、バリデーションヘルパー

- [x] 9.9 ユーティリティテスト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/lib/utils.test.ts`
  - 依存: 9.8

### Phase 9 チェックポイント
- [x] Phase 9 レビュー完了 ⚠️

---

## Phase 10: フロントエンド - カスタムフック (Frontend Custom Hooks)

### 実装タスク
- [x] 10.1 認証フック (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/hooks/useAuth.ts`
  - 内容: ログイン状態管理、ユーザー情報取得
  - 依存: 9.6

- [x] 10.2 認証フックテスト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/hooks/useAuth.test.ts`
  - 依存: 10.1

- [x] 10.3 予定フック (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/hooks/useEvents.ts`
  - 内容: 予定一覧取得、予定作成・更新・削除
  - 依存: 9.4

- [x] 10.4 予定フックテスト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/hooks/useEvents.test.ts`
  - 依存: 10.3

- [x] 10.5 リソースフック (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/hooks/useResources.ts`
  - 内容: リソース検索、空き状況確認
  - 依存: 9.4

- [x] 10.6 リソースフックテスト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/hooks/useResources.test.ts`
  - 依存: 10.5

### Phase 10 チェックポイント
- [x] Phase 10 レビュー完了 ⚠️

---

## Phase 11: フロントエンド - 共通UIコンポーネント (Frontend Common UI Components)

### 実装タスク
- [x] 11.1 ボタン (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/Button.tsx`

- [x] 11.2 ボタンテスト (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/Button.test.tsx`
  - 依存: 11.1

- [x] 11.3 モーダル (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/Modal.tsx`

- [x] 11.4 モーダルテスト (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/Modal.test.tsx`
  - 依存: 11.3

- [x] 11.5 日付ピッカー (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/DatePicker.tsx`
  - 依存: 9.8

- [x] 11.6 日付ピッカーテスト (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/DatePicker.test.tsx`
  - 依存: 11.5

- [x] 11.7 トースト通知 (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/Toast.tsx`

- [x] 11.8 トースト通知テスト (AI Assistant - レビュー待ち 2025-11-24 15:25)
  - ファイル: `frontend/src/components/ui/Toast.test.tsx`
  - 依存: 11.7

### Phase 11 チェックポイント
- [x] Phase 11 レビュー完了 ⚠️

---

## Phase 12: フロントエンド - レイアウトコンポーネント (Frontend Layout Components)

### 実装タスク
- [x] 12.1 ヘッダー (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/layout/Header.tsx`
  - 依存: 10.1, 11.1

- [x] 12.2 ヘッダーテスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/layout/Header.test.tsx`
  - 依存: 12.1

- [x] 12.3 サイドバー (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/layout/Sidebar.tsx`
  - 依存: 10.1

- [x] 12.4 サイドバーテスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/layout/Sidebar.test.tsx`
  - 依存: 12.3

- [x] 12.5 フッター (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/layout/Footer.tsx`

- [x] 12.6 フッターテスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/layout/Footer.test.tsx`
  - 依存: 12.5

### Phase 12 チェックポイント
- [x] Phase 12 レビュー完了 ⚠️

---

## Phase 13: フロントエンド - 機能別コンポーネント (Frontend Feature Components)

### 実装タスク
- [x] 13.1 カレンダーコンポーネント (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/calendar/CalendarView.tsx`
  - 依存: 10.3, 11.5

- [x] 13.2 カレンダーコンポーネントテスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/calendar/CalendarView.test.tsx`
  - 依存: 13.1

- [x] 13.3 予約作成フォーム (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/reservation/ReservationForm.tsx`

- [x] 13.4 予約作成フォームテスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/reservation/ReservationForm.test.tsx`
  - 依存: 13.3

- [x] 13.5 予約詳細 (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/reservation/ReservationDetail.tsx`
  - 依存: 10.3, 11.3

- [x] 13.6 予約詳細テスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/reservation/ReservationDetail.test.tsx`
  - 依存: 13.5

- [x] 13.7 承認一覧 (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/approval/ApprovalList.tsx`
  - 依存: 10.3, 11.1

- [x] 13.8 承認一覧テスト (AI Assistant - レビュー待ち 2025-11-24 18:00)
  - ファイル: `frontend/src/components/features/approval/ApprovalList.test.tsx`
  - 依存: 13.7

### Phase 13 チェックポイント
- [x] Phase 13 レビュー完了 ⚠️

---

## Phase 14: フロントエンド - ページ (Frontend Pages)

### 実装タスク
- [x] 14.1 ルートレイアウト (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/layout.tsx`
  - 依存: 12.1, 12.3, 12.5

- [x] 14.2 トップページ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/page.tsx`
  - 依存: 14.1

- [x] 14.3 ログインページ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/(auth)/login/page.tsx`
  - 依存: 10.1

- [x] 14.4 コールバックページ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/(auth)/callback/page.tsx`
  - 依存: 10.1

- [x] 14.5 ダッシュボードページ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/dashboard/page.tsx`
  - 依存: 13.1, 10.3

- [x] 14.6 予定管理ページ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/events/page.tsx`
  - 依存: 13.3, 13.5, 10.3

- [x] 14.7 リソース管理ページ (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/src/app/resources/page.tsx`
  - 依存: 10.5

### E2Eテスト
- [x] 14.8 E2Eテスト - 認証フロー (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/tests/e2e/auth.spec.ts`
  - 依存: 14.3, 14.4

- [x] 14.9 E2Eテスト - 予約作成フロー (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/tests/e2e/reservation.spec.ts`
  - 依存: 14.5, 14.6

- [x] 14.10 E2Eテスト - リソース検索フロー (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `frontend/tests/e2e/resource.spec.ts`
  - 依存: 14.7

### Phase 14 チェックポイント
- [x] Phase 14 レビュー完了 ⚠️ (AI Assistant - レビュー待ち 2025-11-24)

---

## Phase 15: 運用スクリプト・テスト統合 (Operational Scripts & Test Integration)

### 実装タスク
- [x] 15.1 マイグレーションスクリプト (AI Assistant - 2025-11-24 21:00)
  - ファイル: `backend/scripts/migrate.sh`

- [x] 15.2 シードデータ投入 (AI Assistant - 2025-11-24 21:00)
  - ファイル: `backend/scripts/seed.go`
  - 依存: 4.1, 4.3

- [x] 15.3 セットアップスクリプト (AI Assistant - 2025-11-24 21:00)
  - ファイル: `scripts/setup.sh`

- [x] 15.4 開発環境起動スクリプト (AI Assistant - 2025-11-24 21:00)
  - ファイル: `scripts/dev.sh`

- [x] 15.5 テスト実行スクリプト (AI Assistant - 2025-11-24 21:00)
  - ファイル: `scripts/test.sh`

- [x] 15.6 クリーンアップスクリプト (AI Assistant - 2025-11-24 21:00)
  - ファイル: `scripts/clean.sh`

### テスト統合
- [x] 15.7 テストカバレッジレポート設定 (AI Assistant - 2025-11-24 21:00)
  - 内容: カバレッジ測定・レポート生成

- [x] 15.8 テストカバレッジ確認 (AI Assistant - 2025-11-24 21:00)
  - 目標: バックエンド80%以上、フロントエンド70%以上

### Phase 15 チェックポイント
- [x] Phase 15 レビュー完了 ⚠️ (AI Assistant - 2025-11-24 21:00)

---

## Phase 16: CI/CD・ドキュメント (CI/CD & Documentation)

### 実装タスク
- [ ] 16.1 CI設定
  - ファイル: `.github/workflows/ci.yml`
  - 内容: Lint、テスト、ビルド、脆弱性スキャン (Trivy/SCA) の自動実行
  - 依存: 15.5

- [x] 16.2 CD設定（Future） (AI Assistant - レビュー待ち 2025-11-24)
  - ファイル: `.github/workflows/deploy.yml`
  - 内容: 自動デプロイ設定（将来実装）
  - 依存: 16.1

- [x] 16.3 CI動作確認 (AI Assistant - レビュー待ち 2025-11-24)
  - 内容: Pull Request作成してCI実行確認

- [ ] 16.4 監査ログ・バックアップ設定
  - ファイル: `docs/ops/retention_policy.md`
  - 内容: 監査ログ保持期間、DBバックアップ (PITR) 設定手順

### Phase 16 チェックポイント
- [ ] Phase 16 レビュー完了 ⚠️

---

## Phase 17: AI機能 (AI Features - Future Scope)

### 実装タスク
- [ ] 17.1 AIサービス
  - ファイル: `backend/internal/service/ai_service.go`
  - 内容: プロンプト構築、LLMクライアント連携、レスポンス解析
  - 依存: 5.3

- [ ] 17.2 AIガバナンス機構
  - ファイル: `backend/internal/service/ai_governance.go`
  - 内容: キルスイッチ実装、品質メトリクス計測、PIIガードレール
  - 依存: 17.1

### Phase 17 チェックポイント
- [ ] Phase 17 レビュー完了 ⚠️

---

## 最終確認

- [ ] 全Phase完了
- [ ] 全テスト合格
- [ ] ドキュメント整備完了
- [ ] 本番環境デプロイ準備完了
