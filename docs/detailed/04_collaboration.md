# コラボレーション機能詳細設計書

**対象機能:** FN-05 (代理操作), FN-06 (外部連携), FN-01 (権限管理の一部)

## 1. はじめに
本ドキュメントは、代理人操作（秘書機能）、社外ゲストとの連携、およびプライバシー制御機能の詳細設計を記述する。

## 2. 代理操作 (Secretary Mode)

### 2.1 権限委譲モデル
ユーザー（委譲者）は、他のユーザー（受任者）に対して自身のカレンダー操作権限を付与できる。

- **データ構造 (Delegations Table)**:
    - `delegator_id` (UUID): 権限を与えるユーザー（役員等）
    - `delegatee_id` (UUID): 権限を受けるユーザー（秘書等）
    - `permissions` (String[]): 付与する権限リスト
        - `READ_PRIVATE`: 非公開予定の閲覧
        - `EDIT`: 予定の作成・編集・削除
        - `RESPOND`: 招待への回答

### 2.2 操作フロー (Acting As)
1.  **モード切替**: 受任者がUI上の「代理操作モード」トグルをONにし、対象委譲者を選択。
2.  **コンテキスト切替**: フロントエンドは `X-Act-As-User: {delegator_id}` ヘッダーを付与してAPIをリクエスト。
3.  **バックエンド検証**:
    - リクエストユーザー(`delegatee`)が、指定された`delegator`に対して有効な委譲設定を持っているか確認。
    - 権限があれば、`delegator`として振る舞い処理を実行。

### 2.3 監査ログ (Audit Logging)
代理操作の透明性を担保するため、**「誰が(Actor)」「誰の(Subject)」「何を(Target)」**操作したかを明確に記録する。

| Field | Value Example | Description |
| :--- | :--- | :--- |
| `actor_id` | `user_secretary_01` | 実際に操作したユーザー |
| `subject_id` | `user_executive_01` | 操作対象（なりかわり先）のユーザー |
| `action` | `CREATE_EVENT` | 操作内容 |
| `target_id` | `event_12345` | 作成された予定ID |
| `metadata` | `{"is_proxy": true}` | 付加情報 |

## 3. 外部ゲスト連携 (External Collaboration)

### 3.1 招待フロー
社外メールアドレスが含まれる予定が作成された場合、システムは標準的なiCalendar形式で招待状を送付する。

1.  **iCal生成**: `text/calendar` MIMEタイプの `.ics` ファイルを生成。
    - `METHOD:REQUEST`
    - `UID`: システム内でユニークなID
    - `ORGANIZER`: 主催者のメールアドレス
    - `ATTENDEE`: ゲストのメールアドレス
2.  **メール送信**: SMTPサーバー経由で送信。
    - 件名: `[Invitation] {Event Title} @ {Date}`
    - 本文: プレーンテキスト + HTML（参加ボタン付き）

### 3.2 参加可否 (RSVP) の取り込み
外部ゲストがメールクライアント（Outlook, Gmail等）で「承諾/辞退」を押下した際のレスポンスを処理する。

-   **方式**: **IMAP Polling** または **専用受信アドレス (invite-reply@esms.corp)** へのWebhook連携。
-   **処理**:
    1.  受信メールの `text/calendar` パートを解析 (`METHOD:REPLY`)。
    2.  `UID` で対象の予定を特定。
    3.  `ATTENDEE` と送信元アドレスを照合。
    4.  `PARTSTAT` (ACCEPTED/DECLINED/TENTATIVE) をDBに反映。

## 4. プライバシー制御 (Privacy Control)

### 4.1 公開範囲レベル
予定ごとに以下の公開レベルを設定可能とする。

| レベル | 名称 | 他者からの見え方 | 用途 |
| :--- | :--- | :--- | :--- |
| **PUBLIC** | 公開 | タイトル、場所、詳細、参加者がすべて見える。 | 一般的な会議、チーム定例 |
| **BUSY_ONLY** | 時間枠のみ | 「予定あり」とだけ表示され、詳細は隠される。 | 面接、評価面談、集中作業 |
| **PRIVATE** | 非公開 | 完全に不可視（空き時間に見える）。 | 個人的な用事（※業務時間外推奨） |

### 4.2 閲覧ロジック
予定取得API (`GET /events`) において、リクエストユーザーと予定所有者の関係性に基づきフィルタリングを行う。

-   **本人 or 代理人(READ_PRIVATE権限あり)**: すべてのフィールドを返す。
-   **一般ユーザー**:
    -   `PUBLIC`: すべて返す。
    -   `BUSY_ONLY`: `start_at`, `end_at` のみを返し、タイトルは「予定あり」に置換。
    -   `PRIVATE`: レスポンスに含めない。
