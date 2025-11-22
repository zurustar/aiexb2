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

## 5. フロントエンド共通コンポーネント

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
