<!--
Depends On: docs/basic_design.md
Depended On By: None
-->
# 共通機能・インフラ詳細設計書

**対象機能:** FN-08 (通知), 共通エラーハンドリング, ログ設計

## 1. はじめに
本ドキュメントは、システム全体で共通利用される機能およびインフラストラクチャの詳細設計を記述する。

## 2. 通知配信基盤 (Notification Infrastructure)

### 2.1 アーキテクチャ
通知機能は、送信チャネル（メール、チャット等）を抽象化した `NotificationService` と、非同期処理のためのジョブキューで構成する。

- **NotificationService Interface**:
    ```go
    type NotificationService interface {
        Send(ctx context.Context, recipient User, message Message) error
    }
    ```
- **非同期処理**:
    - APIサーバーは通知リクエストを **Redis Job Queue** (e.g., Asynq, Bull) にエンキューし、即座にレスポンスを返す。
    - **Idempotency:** ジョブIDは `ReservationID + JobType` で生成し、重複実行を防止する。
    - **Retry Strategy:** 指数バックオフ (Exponential Backoff) とジッター (Jitter) を組み合わせ、最大3回のリトライを行う。
    - **Dead Letter Queue (DLQ):** 3回失敗したジョブはDLQへ送られ、手動調査対象とする。
    - バックグラウンドワーカーがキューからジョブを取り出し、実際の送信処理（SMTP, Webhook）を実行する。

### 2.2 通知チャネル
1.  **Email**: SMTPサーバー (AWS SES / SendGrid 等) を利用。
    - HTMLメールとプレーンテキストのマルチパート送信。
2.  **Chat Tools**: Slack / Microsoft Teams への Incoming Webhook。
    - ユーザー設定でWebhook URLを登録可能とする。

## 3. 共通APIレスポンス形式 (API Standards)

### 3.1 JSONエンベロープ
全てのREST APIは以下の統一フォーマットでレスポンスを返す。

**成功時 (HTTP 200/201):**
```json
{
  "success": true,
  "data": { ... },  // 実際のリソースデータ
  "meta": {         // ページネーション等
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

**エラー時 (HTTP 4xx/5xx):**
```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "The requested meeting room does not exist.",
    "details": [ ... ] // バリデーションエラーの詳細等
  }
}
```

### 3.2 エラーコード体系
| HTTP Status | Error Code | Description |
| :--- | :--- | :--- |
| 400 | `INVALID_ARGUMENT` | 入力値不正。詳細は `details` に記述。 |
| 401 | `UNAUTHENTICATED` | 未認証、またはトークン無効。 |
| 403 | `PERMISSION_DENIED` | 権限不足。 |
| 404 | `RESOURCE_NOT_FOUND` | リソースが存在しない。 |
| 409 | `CONFLICT` | リソース競合（排他エラー等）。 |
| 500 | `INTERNAL_ERROR` | サーバー内部エラー。 |

## 4. ログ設計 (Logging)

### 4.1 ログフォーマット
CloudWatch Logs / Datadog 等での解析を容易にするため、**JSON構造化ログ**を採用する。
また、ログ出力時にPII（個人情報）が含まれる場合は、自動的にマスク処理を行う（例: メールアドレスの部分隠蔽）。

```json
{
  "timestamp": "2025-11-23T10:00:00Z",
  "level": "INFO",
  "trace_id": "a1b2c3d4e5",
  "user_id": "user_123",
  "message": "Reservation created successfully",
  "context": {
    "reservation_id": "res_999",
    "resource_id": "room_a"
  }
}
```

### 4.2 ログレベル運用
-   **ERROR**: システムの動作継続に影響する異常（DB接続断、パニック等）。即時アラート対象。
-   **WARN**: 正常ではないが動作継続可能（APIリトライ発生、バリデーションエラー多発等）。
-   **INFO**: 正常系の主要イベント（ログイン、予約作成、バッチ開始/終了）。
-   **DEBUG**: 開発・デバッグ用の詳細情報（SQLクエリ、内部変数値）。本番環境では原則出力しない。

### 4.3 可観測性 (Observability)
-   **Logging:** ECS FireLens (Fluent Bit) を使用し、ログを集約・転送する。
-   **Tracing:** OpenTelemetry を導入し、リクエストの分散トレーシングを行う。
-   **Metrics:** Prometheus 等で以下の主要メトリクスを収集し、Grafana で可視化する。
    -   `http_request_duration_seconds`: APIレイテンシ (p95, p99)
    -   `reservation_success_total`: 予約作成成功数
    -   `job_failure_total`: バックグラウンドジョブ失敗数
-   **Alerting:** SLO (Service Level Objective) 違反時にアラートを発報する（例: APIエラー率 > 0.5%）。
-   **Cost Management:** 環境（dev/stg/prod）ごとにコストをタグ付けし、月次予算超過時にアラートを通知する。

## 5. 可用性・災害復旧 (Availability & DR)
-   **RPO (Recovery Point Objective):** 15分
-   **RTO (Recovery Time Objective):** 1時間
-   **Backup:** RDSの自動スナップショットに加え、PITR (Point-in-Time Recovery) を有効化する。
-   **Redundancy:** クロスリージョンレプリケーション、またはIaCによるDR環境の迅速な立ち上げを準備する。
-   **Drill:** 四半期ごとに復旧演習を実施する。

## 6. 品質保証・CI/CD (Quality & CI/CD)
-   **Quality Gates:**
    -   Lint (Static Analysis) 通過
    -   Unit/Integration Test Pass (Coverage > 80%)
    -   Security Scan (SCA/Container) Pass
-   **Deployment:**
    -   DBマイグレーションは後方互換性を維持する。
    -   Blue/Greenデプロイメント等により、無停止でのリリースを実現する。

## 7. フロントエンド共通コンポーネント

### 5.1 デザインシステム
-   **Color Palette**:
    -   Primary: Brand Blue (`#0052CC`)
    -   Success: Green (`#36B37E`)
    -   Warning: Orange (`#FFAB00`)
    -   Error: Red (`#FF5630`)
-   **Typography**:
    -   Font Family: Inter, Noto Sans JP
    -   Scale: 12px, 14px (Body), 16px, 20px, 24px (Headings)

### 5.2 共通UIコンポーネント
Atomic Designに基づき、再利用可能なコンポーネントを整備する。
-   `Button`: Primary, Secondary, Ghost, Danger のバリエーション。
-   `Modal`: ダイアログ表示用。
-   `Toast`: 成功/エラーメッセージの通知用。
-   `DatePicker`: カレンダー操作用（日付・時刻選択）。

### 5.3 アクセシビリティ方針
-   **フォーカス管理:** キーボード操作のみで全コンポーネントを操作可能にし、モーダル表示時はフォーカストラップを有効化する。
-   **ARIA属性:** フォーム部品・ダイアログ・トーストに適切な`aria-*`属性を付与し、トーストは `role="status"` で読み上げ対応させる。
-   **コントラスト:** 主要カラーについてWCAG AA準拠のコントラスト比（4.5:1以上）を確認し、非準拠のテーマは利用禁止とする。

## 8. フロントエンド - API クライアント

### 8.1 タイムアウト設定
-   **デフォルトタイムアウト:** 30秒
-   **設定可能:** エンドポイントごとにタイムアウト値を変更可能
-   **実装:** AbortControllerを使用してリクエストをキャンセル
-   **エラーハンドリング:** タイムアウト時は408ステータスとTIMEOUTエラーコードを返却

### 8.2 リクエストキャンセル
-   **AbortSignal:** 外部からのキャンセルシグナルをサポート
-   **タイムアウト管理:** setTimeout + AbortControllerで実装
-   **クリーンアップ:** 成功時・エラー時ともにタイマーをクリア

## 9. フロントエンド - セッション管理

### 9.1 クロスタブ同期
-   **実装方法:** `storage` イベントリスナーを使用
-   **同期対象:** ログイン/ログアウト状態
-   **遅延:** 1秒以内に他のタブへ反映
-   **対象キー:** `esms.session`

### 9.2 自動トークンリフレッシュ
-   **チェック間隔:** 60秒ごと
-   **リフレッシュタイミング:** 有効期限の5分前
-   **エンドポイント:** `/api/v1/auth/refresh`
-   **失敗時の処理:** セッションをクリアしてログアウト状態へ遷移

### 9.3 トークン有効期限管理
-   **expiresAt:** セッションに有効期限（Unix timestamp）を保存
-   **チェック:** `shouldRefreshToken()` メソッドで判定
-   **自動実行:** useAuthフック内で定期的にチェック

### 5.4 コンポーネントのテスト方針
-   **スナップショット:** 原子的なUIコンポーネントはStorybookまたはJestのスナップショットテストで意図しないUI差分を検知する。
-   **挙動テスト:** モーダルやトーストなど状態を持つコンポーネントはReact Testing Library等でDOMイベントをシミュレートし、フォーカス移動・閉じる操作・ARIA属性の有無を検証する。
-   **アクセシビリティLint:** `eslint-plugin-jsx-a11y` をCIで必須化し、違反があるPRはブロックする。
