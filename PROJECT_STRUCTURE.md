# プロジェクトディレクトリ構成（開発者向け）

このドキュメントは開発者向けに、プロジェクトの詳細なディレクトリ構成と設計方針を説明します。

## 完全なディレクトリ構成

```
aiexb2/
├── docker-compose.yml              # 開発環境の全サービス定義
├── docker-compose.prod.yml         # 本番環境用（オプション）
├── .env.example                    # 環境変数のテンプレート
├── .gitignore
├── README.md
├── Makefile                        # よく使うコマンドのショートカット
│
├── docs/                           # ドキュメント（既存）
│   ├── requirements.md
│   ├── usecases.md
│   ├── ieee830.md
│   ├── basic_design.md
│   └── detailed/
│       ├── 01_auth_security.md
│       ├── 02_schedule_resource.md
│       ├── 03_ai_assistant.md
│       ├── 04_collaboration.md
│       └── 05_common_infra.md
│
├── frontend/                       # Next.js フロントエンド
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── package.json
│   ├── package-lock.json
│   ├── next.config.js
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   ├── .eslintrc.json
│   ├── public/                     # 静的ファイル
│   │   ├── favicon.ico
│   │   └── images/
│   ├── src/
│   │   ├── app/                    # Next.js App Router
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   ├── (auth)/             # 認証関連ページ
│   │   │   │   ├── login/
│   │   │   │   └── callback/
│   │   │   ├── dashboard/          # ダッシュボード
│   │   │   ├── events/             # 予定管理
│   │   │   ├── resources/          # リソース管理
│   │   │   └── api/                # API Routes (BFF)
│   │   ├── components/             # UIコンポーネント
│   │   │   ├── ui/                 # 共通UIコンポーネント
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   ├── DatePicker.tsx
│   │   │   │   └── Toast.tsx
│   │   │   ├── layout/             # レイアウトコンポーネント
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   └── Footer.tsx
│   │   │   └── features/           # 機能別コンポーネント
│   │   │       ├── calendar/
│   │   │       ├── reservation/
│   │   │       └── approval/
│   │   ├── lib/                    # ユーティリティ・ヘルパー
│   │   │   ├── api-client.ts       # APIクライアント
│   │   │   ├── auth.ts             # 認証ヘルパー
│   │   │   └── utils.ts
│   │   ├── hooks/                  # カスタムフック
│   │   │   ├── useAuth.ts
│   │   │   ├── useEvents.ts
│   │   │   └── useResources.ts
│   │   ├── types/                  # TypeScript型定義
│   │   │   ├── api.ts
│   │   │   ├── models.ts
│   │   │   └── index.ts
│   │   └── styles/                 # グローバルスタイル
│   │       └── globals.css
│   └── tests/                      # フロントエンドテスト
│       ├── unit/
│       └── e2e/
│
├── backend/                        # Go バックエンド
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── go.mod
│   ├── go.sum
│   ├── .air.toml                   # ホットリロード設定
│   ├── cmd/
│   │   ├── api/                    # APIサーバー
│   │   │   └── main.go
│   │   └── worker/                 # バックグラウンドワーカー
│   │       └── main.go
│   ├── internal/                   # 内部パッケージ
│   │   ├── config/                 # 設定管理
│   │   │   └── config.go
│   │   ├── domain/                 # ドメインモデル
│   │   │   ├── user.go
│   │   │   ├── reservation.go
│   │   │   ├── resource.go
│   │   │   └── audit_log.go
│   │   ├── repository/             # データアクセス層
│   │   │   ├── user_repository.go
│   │   │   ├── reservation_repository.go
│   │   │   └── resource_repository.go
│   │   ├── service/                # ビジネスロジック
│   │   │   ├── auth_service.go
│   │   │   ├── reservation_service.go
│   │   │   ├── approval_service.go
│   │   │   └── notification_service.go
│   │   ├── handler/                # HTTPハンドラー
│   │   │   ├── auth_handler.go
│   │   │   ├── reservation_handler.go
│   │   │   ├── resource_handler.go
│   │   │   └── middleware.go
│   │   ├── cache/                  # Redisキャッシュ
│   │   │   └── redis_client.go
│   │   ├── queue/                  # ジョブキュー
│   │   │   └── job_queue.go
│   │   └── util/                   # ユーティリティ
│   │       ├── logger.go
│   │       ├── validator.go
│   │       └── time.go
│   ├── pkg/                        # 外部公開可能なパッケージ
│   │   └── oidc/                   # OIDC認証
│   │       └── client.go
│   ├── migrations/                 # DBマイグレーション
│   │   ├── 000001_init_schema.up.sql
│   │   ├── 000001_init_schema.down.sql
│   │   ├── 000002_add_approval.up.sql
│   │   └── 000002_add_approval.down.sql
│   ├── scripts/                    # 運用スクリプト
│   │   ├── seed.go                 # テストデータ投入
│   │   └── migrate.sh
│   └── tests/                      # バックエンドテスト
│       ├── unit/
│       ├── integration/
│       └── fixtures/
│
├── database/                       # データベース関連
│   ├── init/                       # 初期化スクリプト
│   │   └── 01_create_extensions.sql
│   └── seed/                       # シードデータ
│       ├── users.sql
│       └── resources.sql
│
├── infra/                          # インフラストラクチャ
│   ├── terraform/                  # Terraform IaC (Future)
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── k8s/                        # Kubernetes manifests (Future)
│       ├── deployment.yaml
│       └── service.yaml
│
├── scripts/                        # プロジェクト全体のスクリプト
│   ├── setup.sh                    # 初回セットアップ
│   ├── dev.sh                      # 開発環境起動
│   ├── test.sh                     # テスト実行
│   └── clean.sh                    # クリーンアップ
│
└── .github/                        # GitHub Actions (CI/CD)
    └── workflows/
        ├── ci.yml
        └── deploy.yml
```

## アーキテクチャ設計方針

### 1. フロントエンド (Next.js)

#### App Router採用
- Next.js 14のApp Routerを採用
- ファイルベースルーティング
- Server ComponentsとClient Componentsの使い分け

#### コンポーネント設計（Atomic Design風）
- **`ui/`**: 汎用UIコンポーネント
  - Button, Modal, DatePicker, Toast等
  - 再利用可能な最小単位
  - ビジネスロジックを含まない
  
- **`layout/`**: レイアウトコンポーネント
  - Header, Sidebar, Footer等
  - ページ全体の構造を定義
  
- **`features/`**: 機能別コンポーネント
  - calendar, reservation, approval等
  - ビジネスロジックを含む
  - 特定の機能に特化

#### 型定義管理
- TypeScript型定義を`types/`に集約
- API型定義とドメインモデル型を分離
- 型の再利用性を重視

#### API通信
- `lib/api-client.ts`で一元管理
- Fetch APIのラッパー
- エラーハンドリングの統一

### 2. バックエンド (Go)

#### Clean Architecture
依存関係を明確に分離し、テスタビリティと保守性を向上：

```
handler → service → repository → database
   ↓         ↓
domain ← domain
```

- **`domain/`**: ドメインモデル
  - ビジネスルールを含むエンティティ
  - 他の層に依存しない
  - 純粋なGoの構造体

- **`repository/`**: データアクセス層
  - DBとの通信を担当
  - インターフェースで抽象化
  - テスト時はモックに差し替え可能

- **`service/`**: ビジネスロジック層
  - ユースケースの実装
  - トランザクション管理
  - 複数のrepositoryを組み合わせ

- **`handler/`**: プレゼンテーション層
  - HTTPリクエスト/レスポンスの処理
  - バリデーション
  - 認証・認可チェック

#### パッケージ可視性制御
- **`internal/`**: プロジェクト内部でのみ使用
  - 外部パッケージからimport不可
  - アプリケーション固有のロジック

- **`pkg/`**: 外部公開可能なパッケージ
  - 他のプロジェクトでも再利用可能
  - 汎用的なユーティリティ

### 3. Docker Compose開発環境

#### サービス構成
- **frontend**: Next.js開発サーバー
- **backend**: Go APIサーバー（Airでホットリロード）
- **worker**: バックグラウンドジョブワーカー
- **postgres**: PostgreSQL 15
- **redis**: Redis 7（キャッシュ + ジョブキュー）
- **mailhog**: 開発用SMTPサーバー

#### ホットリロード
- **Go**: Air（`.air.toml`で設定）
  - ファイル変更を検知して自動再ビルド
  - 開発効率の向上

- **Next.js**: 標準のdev server
  - Fast Refresh対応
  - 即座に変更を反映

#### ボリュームマウント
```yaml
volumes:
  - ./backend:/app          # ソースコードをマウント
  - /app/node_modules       # node_modulesは除外
```

### 4. データベースマイグレーション

#### ツール選定
- golang-migrate または goose
- Up/Down両方のSQLを必ず作成
- バージョン管理とロールバック可能性

#### 命名規則
```
000001_init_schema.up.sql
000001_init_schema.down.sql
000002_add_approval.up.sql
000002_add_approval.down.sql
```

#### 実行方法
```bash
make migrate        # Up実行
make migrate-down   # Down実行（ロールバック）
```

### 5. テスト戦略

#### Unit Test
- 各層で独立したテスト
- モックを使用して依存を排除
- カバレッジ目標: 80%以上

#### Integration Test
- DB接続を含むテスト
- テスト用DBコンテナを使用
- トランザクションロールバックで高速化

#### E2E Test
- Playwright等でブラウザテスト
- 主要なユーザーフローをカバー
- CI/CDで自動実行

## 開発ワークフロー

### 機能開発フロー

```bash
# 1. ブランチ作成
git checkout -b feature/reservation-approval

# 2. 開発環境起動
make up

# 3. コード変更（ホットリロードで即反映）
# - backend/internal/service/approval_service.go を編集
# - frontend/src/components/features/approval/ を編集

# 4. テスト実行
make test

# 5. マイグレーション追加（必要な場合）
# backend/migrations/000003_add_approval_status.up.sql を作成
make migrate

# 6. コミット & プッシュ
git add .
git commit -m "feat: add approval workflow"
git push origin feature/reservation-approval
```

### マイグレーション追加フロー

```bash
# 1. マイグレーションファイル作成
# backend/migrations/000003_add_new_table.up.sql
# backend/migrations/000003_add_new_table.down.sql

# 2. Up SQLを記述
CREATE TABLE new_table (...);

# 3. Down SQLを記述（ロールバック用）
DROP TABLE IF EXISTS new_table;

# 4. マイグレーション実行
make migrate

# 5. 確認
make shell-db
\dt  # テーブル一覧確認
```

## Docker Composeサービス詳細

実際の`docker-compose.yml`を参照してください。主要な設定：

- **ヘルスチェック**: postgres, redisで実装
- **依存関係**: depends_onで起動順序を制御
- **環境変数**: .envファイルから読み込み
- **ボリューム**: データ永続化とホットリロード

## CI/CD パイプライン

### GitHub Actions
- **CI**: Pull Request時に自動テスト
- **CD**: mainブランチマージ時に自動デプロイ（Future）

### テストステージ
1. Lint（golangci-lint, ESLint）
2. Unit Test
3. Integration Test
4. E2E Test（Playwright）

## 参考リソース

- [Go Clean Architecture](https://github.com/bxcodec/go-clean-arch)
- [Next.js App Router](https://nextjs.org/docs/app)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Air - Go Hot Reload](https://github.com/cosmtrek/air)
