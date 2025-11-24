# Enterprise Schedule Management System (ESMS)

企業向け次世代統合スケジュール管理システム

## プロジェクト概要

- **ユーザー規模**: 10,000人未満
- **開発手法**: AI-Driven Development (AIDD)
- **アーキテクチャ**: モノレポ構成（Frontend + Backend）

## 技術スタック

- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS
- **Backend**: Go (Golang) + Clean Architecture
- **Database**: PostgreSQL 15
- **Cache/Queue**: Redis 7
- **Development**: Docker Compose

## クイックスタート

### 前提条件
- Docker Desktop
- Make

### セットアップ

```bash
# リポジトリをクローン
git clone <repository-url>
cd aiexb2

# 環境変数を設定（オプション: デフォルト値で動作します）
cp .env.example .env

# セットアップ実行（Docker イメージビルド、DB作成、マイグレーション、シードデータ投入）
make setup

# 開発環境起動
make up
```

> **Note**: `.env`ファイルはオプションです。作成しなくてもデフォルト値で動作します。
> Keycloakを含むすべてのサービスが自動的に設定されます。


### アクセスURL

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Keycloak (認証)**: http://localhost:8180
- **MailHog (開発用メール)**: http://localhost:8025

### テストユーザー

開発環境では以下のテストユーザーが自動的に作成されます:

| ユーザー名 | パスワード | ロール |
|-----------|----------|--------|
| admin@esms.local | admin123 | admin, user |
| user@esms.local | user123 | user |
| manager@esms.local | manager123 | admin, user |

**Keycloak管理コンソール**: http://localhost:8180 (admin / admin)


## 開発コマンド

```bash
make help      # コマンド一覧表示
make up        # 開発環境起動
make down      # 開発環境停止
make logs      # ログ確認
make test      # テスト実行
make migrate   # DBマイグレーション
make clean     # クリーンアップ
```

詳細は `make help` を参照してください。

## テスト計画（開発環境）

開発用の Docker Compose 環境を前提に、以下のコマンドで品質チェックを実施します。テストデータは `make setup` または `make seed` で投入される `database/seed/*.sql`（ユーザー・会議室など）を利用します。

- **コード整形/静的解析**
  - Backend Format: `docker-compose run --rm backend gofmt -w .`
  - Backend Lint: `docker-compose run --rm backend go vet ./...`
  - Frontend Lint: `docker-compose run --rm frontend npm run lint`
- **ユニット/統合テスト**
  - Backend Unit/Integration: `docker-compose run --rm backend go test ./... -v`
  - Frontend Unit: `docker-compose run --rm frontend npm test`
- **E2E テスト**（Playwright 追加時想定）
  - `docker-compose run --rm frontend npx playwright test`
  - 前提: `make up` でフロント/バックエンドを起動し、シードデータ投入済みであること。

> メモ: E2E は Playwright セットアップ後に有効化を想定しています。追加のモック/フィクスチャが必要な場合は `frontend/e2e` 配下に配置してください。

## ディレクトリ構成

```
aiexb2/
├── frontend/          # Next.js フロントエンド
├── backend/           # Go バックエンド
├── database/          # DB初期化・シードデータ
├── docs/              # プロジェクトドキュメント
├── infra/             # インフラストラクチャ (Future)
└── scripts/           # 運用スクリプト
```

詳細な構成と設計方針は [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) を参照してください。

## ドキュメント

### 設計ドキュメント
- [要件定義書](./docs/requirements.md)
- [ユースケース定義書](./docs/usecases.md)
- [ソフトウェア要求仕様書 (SRS)](./docs/ieee830.md)
- [基本設計書](./docs/basic_design.md)
- [詳細設計書](./docs/detailed/)

### 開発者向けドキュメント
- [プロジェクト構造](./PROJECT_STRUCTURE.md) - ディレクトリ構成と設計方針

## トラブルシューティング

### ポートが既に使用されている

```bash
lsof -i :3000  # 使用中のポート確認
make down && make up
```

### データベースをリセットしたい

```bash
make clean
make setup
```

### ログを確認したい

```bash
make logs              # 全サービス
docker-compose logs backend   # 特定のサービス
```

### Keycloakが起動しない

```bash
# Keycloakのログを確認
docker-compose logs keycloak

# Keycloakデータベースを確認
docker-compose exec postgres psql -U esms_user -d keycloak -c "\dt"

# 完全リセット
make clean
make setup
```
